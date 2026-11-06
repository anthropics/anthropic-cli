package kubernetes

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"maps"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"testing"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	utilyaml "k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/yaml"
)

// chartDir is the Helm chart these tests render.
const chartDir = "chart"

// chartRelease and chartNamespace are what helm is told on every render.
const (
	chartRelease   = "cma"
	chartNamespace = "cma"
)

// minimalSettings are the choices every install has to make besides its
// platform, made the simplest way: no cluster range is added to the internet
// policy's except list, and no policy opens a private CMA server.
var minimalSettings = []string{
	"egress.clusterIsPrivate=true",
	"egress.cmaServer.public=true",
}

// kubeDNSUpstreamIP stands for the address of a cluster's kube-dns-upstream
// Service.
const kubeDNSUpstreamIP = "192.0.2.10"

const (
	dispatcherDigest = "sha256:1111111111111111111111111111111111111111111111111111111111111111"
	sandboxDigest    = "sha256:2222222222222222222222222222222222222222222222222222222222222222"
)

const (
	sandboxAPIBeta  = "extensions.agents.x-k8s.io/v1beta1"
	sandboxAPIAlpha = "extensions.agents.x-k8s.io/v1alpha1"
)

// install is one set of choices a customer can install with.
type install struct {
	platform string
	// settings are the chart values set besides platform and minimalSettings.
	settings []string
	// sandboxAPIVersion is what the SandboxTemplate and SandboxWarmPool must
	// carry after rendering.
	sandboxAPIVersion string
}

func (i install) name() string {
	return strings.Join(append([]string{i.platform}, i.settings...), "+")
}

var installs = []install{
	{platform: "gke", sandboxAPIVersion: sandboxAPIAlpha},
	{platform: "eks", sandboxAPIVersion: sandboxAPIBeta},
	{platform: "generic", sandboxAPIVersion: sandboxAPIBeta},
	{platform: "gke", settings: []string{"gke.dns=cloudDns"}, sandboxAPIVersion: sandboxAPIAlpha},
	{platform: "gke", settings: []string{"egress.mode=fqdnGke"}, sandboxAPIVersion: sandboxAPIAlpha},
}

// clusterScopedKinds are the rendered kinds that take no namespace.
var clusterScopedKinds = []string{"ClusterRoleBinding", "RuntimeClass"}

// rendered is one object the chart renders.
type rendered struct {
	metav1.PartialObjectMetadata
	// obj is the typed form, or nil for a kind outside client-go's scheme.
	obj runtime.Object
	// sandboxPod is spec.podTemplate, set for a SandboxTemplate only.
	sandboxPod *corev1.PodTemplateSpec
	// specKeys is the top-level spec keys, set for a SandboxTemplate or
	// SandboxWarmPool only.
	specKeys []string
}

// addonRejectedFields are the spec fields upstream's v1beta1 stores and the
// GKE add-on's v1alpha1 rejects as unknown, by kind.
var addonRejectedFields = map[string][]string{
	"SandboxTemplate": {"envVarsInjectionPolicy", "volumeClaimTemplatesPolicy", "volumeClaimTemplates"},
	"SandboxWarmPool": {"updateStrategy"},
}

func TestChartRendersAnInstallEachPlatformAccepts(t *testing.T) {
	helm := helmBinary(t)
	appVersion := chartMetadata(t).AppVersion
	for _, in := range installs {
		t.Run(in.name(), func(t *testing.T) {
			// When: the chart is rendered with the install's choices.
			objects := helmTemplate(t, helm, withMinimalSettings(append([]string{"platform=" + in.platform}, in.settings...)...)...)

			// Then: every object is well formed and sits in the release's namespace.
			for _, r := range objects {
				wantNS := chartNamespace
				if slices.Contains(clusterScopedKinds, r.Kind) {
					wantNS = ""
				}
				if r.Namespace != wantNS {
					t.Errorf("%s %s/%s: namespace %q, want %q", r.APIVersion, r.Kind, r.Name, r.Namespace, wantNS)
				}
				if (r.Kind == "SandboxTemplate" || r.Kind == "SandboxWarmPool") && r.APIVersion != in.sandboxAPIVersion {
					t.Errorf("%s %s: apiVersion %q, want %q", r.Kind, r.Name, r.APIVersion, in.sandboxAPIVersion)
				}
			}
			assertPublishedImages(t, objects, appVersion)
			assertBindingFollowsNamespace(t, objects)
			assertDispatcherAPIVersion(t, objects, in.sandboxAPIVersion)
			assertSandboxSchema(t, objects, in.sandboxAPIVersion)

			// Kubernetes lets traffic through when any policy allows it, so the
			// hostname policy narrows nothing unless the public-internet policy
			// is gone.
			wantFQDN := slices.Contains(in.settings, "egress.mode=fqdnGke")
			hasInternet := find(objects, "NetworkPolicy", "cma-sandbox-internet") != nil
			hasFQDN := find(objects, "FQDNNetworkPolicy", "cma-sandbox-fqdn-egress") != nil
			if hasInternet != !wantFQDN {
				t.Errorf("NetworkPolicy cma-sandbox-internet present: %v, want %v", hasInternet, !wantFQDN)
			}
			if hasFQDN != wantFQDN {
				t.Errorf("FQDNNetworkPolicy cma-sandbox-fqdn-egress present: %v, want %v", hasFQDN, wantFQDN)
			}
			wantCloudDNS := slices.Contains(in.settings, "gke.dns=cloudDns")
			if hasCloudDNS := find(objects, "NetworkPolicy", "cma-sandbox-dns-cloud") != nil; hasCloudDNS != wantCloudDNS {
				t.Errorf("NetworkPolicy cma-sandbox-dns-cloud present: %v, want %v", hasCloudDNS, wantCloudDNS)
			}

			// Every platform has exactly one way for a sandbox to prove who it
			// is. On gke the dispatcher may create the per-sandbox nonce Secret,
			// and the pod carries no token, which the add-on rejects. Elsewhere
			// the pod mounts the token the dispatcher verifies, and the
			// dispatcher holds no verb on Secrets.
			wantNonce := in.platform == "gke"
			if hasNonce := dispatcherRoleGrants(t, objects, "", "secrets", "create"); hasNonce != wantNonce {
				t.Errorf("Role cma-dispatcher grants create on secrets: %v, want %v", hasNonce, wantNonce)
			}
			if hasToken := fetchContainerMounts(t, objects, "/var/run/secrets/cma-dispatcher"); hasToken != !wantNonce {
				t.Errorf("fetch container mounts /var/run/secrets/cma-dispatcher: %v, want %v", hasToken, !wantNonce)
			}
		})
	}
}

func TestChartPullsBothImagesByDigestWhenGivenTheirDigests(t *testing.T) {
	// When: the chart is rendered with the two digests a release's notes list.
	objects := helmTemplate(t, helmBinary(t), withMinimalSettings("platform=eks",
		"image.dispatcherDigest="+dispatcherDigest, "image.sandboxDigest="+sandboxDigest)...)

	// Then: every container names its image by its digest and by no tag.
	for _, spec := range podSpecs(objects) {
		for _, c := range slices.Concat(spec.InitContainers, spec.Containers) {
			want := "us-docker.pkg.dev/anthropic-public-registry/sandboxes/cma-sandbox@" + sandboxDigest
			if c.Name == "dispatcher" {
				want = "us-docker.pkg.dev/anthropic-public-registry/sandboxes/cma-dispatcher@" + dispatcherDigest
			}
			if c.Image != want {
				t.Errorf("container %s uses image %q, want %q", c.Name, c.Image, want)
			}
		}
	}
}

func TestChartKeepsSandboxesOffTheClusterRangesItIsGiven(t *testing.T) {
	// Given: a GKE cluster with a custom Service range and a public API endpoint.
	const apiEndpoint, serviceCIDR = "203.0.113.10/32", "198.51.100.0/24"

	// When: the chart is rendered with both.
	objects := helmTemplate(t, helmBinary(t), "platform=gke", "egress.cmaServer.public=true",
		"egress.clusterCIDRs={"+apiEndpoint+"}", "gke.serviceCIDR="+serviceCIDR)

	// Then: the internet policy's except list holds both.
	except := networkPolicy(t, objects, "cma-sandbox-internet").Spec.Egress[0].To[0].IPBlock.Except
	for _, cidr := range []string{apiEndpoint, serviceCIDR} {
		if !slices.Contains(except, cidr) {
			t.Errorf("cma-sandbox-internet: the except list %v does not hold %s", except, cidr)
		}
	}
}

func TestChartOpensEgressToAPrivateServerAndNothingElse(t *testing.T) {
	helm := helmBinary(t)
	private := []string{"platform=generic", "egress.clusterIsPrivate=true", "egress.cmaServer.port=8400"}

	// A public server needs no policy of its own.
	if find(helmTemplate(t, helm, withMinimalSettings("platform=generic")...), "NetworkPolicy", "cma-sandbox-to-server") != nil {
		t.Error("cma-sandbox-to-server is rendered for a public server")
	}

	// A server on a private address is opened by that address and port alone.
	byCIDR := networkPolicy(t, helmTemplate(t, helm, append(private, "egress.cmaServer.cidr=10.20.0.5/32")...), "cma-sandbox-to-server").Spec.Egress
	if len(byCIDR) != 1 || len(byCIDR[0].To) != 1 || byCIDR[0].To[0].IPBlock == nil || byCIDR[0].To[0].IPBlock.CIDR != "10.20.0.5/32" || byCIDR[0].Ports[0].Port.IntValue() != 8400 {
		t.Errorf("cma-sandbox-to-server by cidr: %+v", byCIDR)
	}

	// A server in the cluster is opened by one peer that sets both selectors,
	// which selects only the labeled pods in that namespace. Split into two
	// peers, the first would select every pod in the namespace.
	byPods := networkPolicy(t, helmTemplate(t, helm, append(private, "egress.cmaServer.pods.namespace=cma-backend",
		"egress.cmaServer.pods.matchLabels.app=server", "egress.cmaServer.pods.matchLabels.tier=api")...), "cma-sandbox-to-server").Spec.Egress
	if len(byPods) != 1 || len(byPods[0].To) != 1 || byPods[0].To[0].NamespaceSelector == nil || byPods[0].To[0].PodSelector == nil {
		t.Fatalf("cma-sandbox-to-server by pods: %+v, want one peer with both a namespace and a pod selector", byPods)
	}
	peer := byPods[0].To[0]
	if !maps.Equal(peer.NamespaceSelector.MatchLabels, map[string]string{"kubernetes.io/metadata.name": "cma-backend"}) ||
		!maps.Equal(peer.PodSelector.MatchLabels, map[string]string{"app": "server", "tier": "api"}) ||
		byPods[0].Ports[0].Port.IntValue() != 8400 {
		t.Errorf("cma-sandbox-to-server by pods: namespace %v, pods %v, port %v", peer.NamespaceSelector.MatchLabels, peer.PodSelector.MatchLabels, byPods[0].Ports[0].Port)
	}
}

func TestChartNodeLocalCachePointsSandboxesAtKubeDNSUpstream(t *testing.T) {
	// When: the chart is told that the cluster runs NodeLocal DNSCache.
	objects := helmTemplate(t, helmBinary(t), withMinimalSettings("platform=gke",
		"gke.dns=nodeLocalCache", "gke.kubeDnsUpstreamIP="+kubeDNSUpstreamIP)...)

	// Then: the sandbox pod asks that address, and names the search list that
	// resolves the dispatcher's short Service name, which dnsPolicy None
	// would otherwise leave empty.
	pod := find(objects, "SandboxTemplate", "cma-sandbox").sandboxPod.Spec
	if pod.DNSPolicy != corev1.DNSNone || pod.DNSConfig == nil || !slices.Equal(pod.DNSConfig.Nameservers, []string{kubeDNSUpstreamIP}) {
		t.Fatalf("dnsPolicy %q, dnsConfig %+v", pod.DNSPolicy, pod.DNSConfig)
	}
	if len(pod.DNSConfig.Searches) == 0 || !strings.HasPrefix(pod.DNSConfig.Searches[0], chartNamespace+".svc.") {
		t.Errorf("searches %v do not start with the domain of Services in namespace %s", pod.DNSConfig.Searches, chartNamespace)
	}
}

func TestChartRefusesValuesItCannotRender(t *testing.T) {
	helm := helmBinary(t)
	for _, refusal := range []struct {
		name     string
		settings []string
		// wantInMessage names the key the customer has to fix, or the command
		// that finds its value.
		wantInMessage string
	}{
		{"platform missing", minimalSettings, "platform"},
		{"platform unknown", withMinimalSettings("platform=aks"), "platform"},
		{"misspelled key", withMinimalSettings("platform=eks", "sandbox.idleTimeuot=5m"), "idleTimeuot"},
		{"warm pool on gke", withMinimalSettings("platform=gke", "sandbox.warmPoolSize=2"), "sandbox.warmPoolSize"},
		{"fqdnGke off GKE", withMinimalSettings("platform=eks", "egress.mode=fqdnGke"), "egress.mode"},
		{"cloudDns off GKE", withMinimalSettings("platform=generic", "gke.dns=cloudDns"), "gke.dns"},
		{
			"nodeLocalCache off GKE",
			withMinimalSettings("platform=eks", "gke.dns=nodeLocalCache", "gke.kubeDnsUpstreamIP="+kubeDNSUpstreamIP),
			"gke.dns is nodeLocalCache",
		},
		{
			"nodeLocalCache without the kube-dns-upstream address",
			withMinimalSettings("platform=gke", "gke.dns=nodeLocalCache"),
			"kubectl -n kube-system get svc kube-dns-upstream -o jsonpath='{.spec.clusterIP}'",
		},
		{
			"kube-dns-upstream address without nodeLocalCache",
			withMinimalSettings("platform=gke", "gke.kubeDnsUpstreamIP="+kubeDNSUpstreamIP),
			"gke.kubeDnsUpstreamIP is " + kubeDNSUpstreamIP + " and gke.dns is kubeDns",
		},
		{
			"no cluster range and not declared private",
			[]string{"platform=eks", "egress.cmaServer.public=true"},
			"egress.clusterCIDRs",
		},
		{
			"no CMA server location",
			[]string{"platform=eks", "egress.clusterIsPrivate=true"},
			"egress.cmaServer names no location",
		},
		{
			"two CMA server locations",
			withMinimalSettings("platform=eks", "egress.cmaServer.cidr=10.20.0.5/32"),
			"egress.cmaServer names 2 locations for the CMA server: public and cidr",
		},
		{
			"CMA server pods without their labels",
			[]string{"platform=eks", "egress.clusterIsPrivate=true", "egress.cmaServer.pods.namespace=cma-server"},
			"egress.cmaServer.pods.matchLabels",
		},
		{
			"registry empty",
			withMinimalSettings("platform=eks", "image.registry="),
			"image.registry is empty",
		},
		{
			"one image pinned and the other not",
			withMinimalSettings("platform=eks", "image.sandboxDigest="+sandboxDigest),
			"Only one of image.dispatcherDigest and image.sandboxDigest is set",
		},
		{
			"a tag where a digest goes",
			withMinimalSettings("platform=eks", "image.dispatcherDigest=1.36.0", "image.sandboxDigest="+sandboxDigest),
			"dispatcherDigest",
		},
	} {
		t.Run(refusal.name, func(t *testing.T) {
			// When: the chart is rendered with the refused values.
			out, err := runHelm(helm, templateArgs(refusal.settings...)...)

			// Then: helm fails and names the key.
			if err == nil {
				t.Fatalf("rendered without error:\n%s", out)
			}
			if !strings.Contains(out, refusal.wantInMessage) {
				t.Errorf("message does not name %s:\n%s", refusal.wantInMessage, out)
			}
			t.Log(out)
		})
	}
}

func TestChartVersionIsItsAppVersion(t *testing.T) {
	// appVersion is the image tag when image.tag is empty, and a release sets
	// both to its own version.
	if chart := chartMetadata(t); chart.Version != chart.AppVersion {
		t.Errorf("version is %q and appVersion is %q; a release sets both", chart.Version, chart.AppVersion)
	}
}

func TestChartVersionIsTheChangelogsTopHeading(t *testing.T) {
	changelog, err := os.ReadFile(filepath.Join("..", "CHANGELOG.md"))
	if err != nil {
		t.Fatal(err)
	}
	plainOrLinkedHeading := regexp.MustCompile(`(?m)^## \[?([^\]\s]+)`)
	topHeading := plainOrLinkedHeading.FindSubmatch(changelog)
	if topHeading == nil {
		t.Fatal("CHANGELOG.md has no version heading")
	}
	if chart := chartMetadata(t); chart.Version != string(topHeading[1]) {
		t.Errorf("version is %q and the top heading of CHANGELOG.md is %q", chart.Version, topHeading[1])
	}
}

func TestChartExamplesRenderAndLint(t *testing.T) {
	// Given: the worked values files, one per platform.
	helm := helmBinary(t)
	examples, err := filepath.Glob(filepath.Join(chartDir, "examples", "values-*.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	var platforms []string
	for _, example := range examples {
		t.Run(filepath.Base(example), func(t *testing.T) {
			// When: the chart is linted and rendered with the example.
			if out, err := runHelm(helm, "lint", chartDir, "--strict", "-f", example); err != nil {
				t.Errorf("helm lint: %v\n%s", err, out)
			}
			out, err := runHelm(helm, "template", chartRelease, chartDir, "--namespace", chartNamespace, "-f", example)
			if err != nil {
				t.Fatalf("helm template: %v\n%s", err, out)
			}

			// Then: every object is well formed.
			decodeObjects(t, []byte(out))
		})
		platforms = append(platforms, examplePlatform(t, example))
	}

	// Then: every platform has a worked example.
	slices.Sort(platforms)
	if want := []string{"eks", "generic", "gke"}; !slices.Equal(platforms, want) {
		t.Errorf("examples cover platforms %v, want %v", platforms, want)
	}
}

func TestDecodeStrictRejectsBrokenObjects(t *testing.T) {
	for name, doc := range map[string]string{
		"misspelled field": `
apiVersion: apps/v1
kind: Deployment
metadata: {name: d}
spec:
  template:
    spec:
      containerz: []
`,
		"wrong type": `
apiVersion: apps/v1
kind: Deployment
metadata: {name: d}
spec:
  replicas: one
`,
		"duplicate field": `
apiVersion: v1
kind: Service
metadata: {name: s}
spec:
  ports: []
  ports: []
`,
		"unknown kind": `
apiVersion: apps/v1
kind: Deploymnet
metadata: {name: d}
`,
		"unknown group": `
apiVersion: extensions.agents.x-k8s.io/v1beta1
kind: SandboxClaim
metadata: {name: c}
`,
		"sandbox template with a misspelled pod field": `
apiVersion: extensions.agents.x-k8s.io/v1beta1
kind: SandboxTemplate
metadata: {name: t}
spec:
  podTemplate:
    spec:
      runtimeClasName: gvisor
      containers: []
`,
		"sandbox template without a pod template": `
apiVersion: extensions.agents.x-k8s.io/v1beta1
kind: SandboxTemplate
metadata: {name: t}
spec:
  podTemplat: {}
`,
		"sandbox template with a misspelled claim field": `
apiVersion: extensions.agents.x-k8s.io/v1beta1
kind: SandboxTemplate
metadata: {name: t}
spec:
  podTemplate: {}
  volumeClaimTemplates:
    - spec:
        accessMode: [ReadWriteOnce]
`,
	} {
		t.Run(name, func(t *testing.T) {
			// Given: an object that is not well formed.
			// When: it is decoded the way the rendered manifests are.
			_, err := decodeStrict([]byte(doc))
			// Then: the decode fails.
			if err == nil {
				t.Fatal("decoded without error")
			}
			t.Log(err)
		})
	}
}

// helmBinary returns the helm binary exec.LookPath finds. Without one it
// skips the test, or fails it when CI is set. CI's runner image includes helm.
// A skip there would hide that these tests stopped running.
func helmBinary(t *testing.T) string {
	t.Helper()
	helm, err := exec.LookPath("helm")
	if err != nil {
		if os.Getenv("CI") != "" {
			t.Fatal("helm is not on PATH and CI is set, so this test fails instead of skipping; install helm, or unset CI to skip (the test renders sandboxes/kubernetes/chart with helm)")
		}
		t.Skip("helm is not on PATH; these tests render sandboxes/kubernetes/chart with the helm binary")
	}
	return helm
}

func withMinimalSettings(settings ...string) []string {
	return slices.Concat(minimalSettings, settings)
}

// chartMetadata returns the two versions Chart.yaml states.
func chartMetadata(t *testing.T) (chart struct{ Version, AppVersion string }) {
	t.Helper()
	raw, err := os.ReadFile(filepath.Join(chartDir, "Chart.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	if err := yaml.Unmarshal(raw, &chart); err != nil {
		t.Fatal(err)
	}
	return chart
}

// helmTemplate renders the chart with settings and returns its objects.
func helmTemplate(t *testing.T, helm string, settings ...string) []rendered {
	t.Helper()
	args := templateArgs(settings...)
	out, err := runHelm(helm, args...)
	if err != nil {
		t.Fatalf("helm %s: %v\n%s", strings.Join(args, " "), err, out)
	}
	return decodeObjects(t, []byte(out))
}

func templateArgs(settings ...string) []string {
	args := []string{"template", chartRelease, chartDir, "--namespace", chartNamespace}
	for _, s := range settings {
		args = append(args, "--set", s)
	}
	return args
}

// runHelm returns helm's standard output, or on failure its standard error.
func runHelm(helm string, args ...string) (string, error) {
	var stdout, stderr bytes.Buffer
	cmd := exec.Command(helm, args...)
	cmd.Stdout, cmd.Stderr = &stdout, &stderr
	if err := cmd.Run(); err != nil {
		return stderr.String(), err
	}
	return stdout.String(), nil
}

// decodeObjects returns every object of a YAML stream, and fails the test for
// each one that is not well formed. A document with no object, which helm
// emits for a template that renders nothing, is dropped.
func decodeObjects(t *testing.T, stream []byte) []rendered {
	t.Helper()
	var objects []rendered
	reader := utilyaml.NewYAMLReader(bufio.NewReader(bytes.NewReader(stream)))
	for {
		doc, err := reader.Read()
		if errors.Is(err, io.EOF) {
			return objects
		}
		if err != nil {
			t.Fatal(err)
		}
		if len(bytes.TrimSpace(stripComments(doc))) == 0 {
			continue
		}
		r, err := decodeStrict(doc)
		if err != nil {
			t.Errorf("%s %s/%s: %v", r.APIVersion, r.Kind, r.Name, err)
			continue
		}
		objects = append(objects, r)
	}
}

// stripComments returns doc without its comment lines.
func stripComments(doc []byte) []byte {
	var kept [][]byte
	for _, line := range bytes.Split(doc, []byte("\n")) {
		if !bytes.HasPrefix(bytes.TrimSpace(line), []byte("#")) {
			kept = append(kept, line)
		}
	}
	return bytes.Join(kept, []byte("\n"))
}

// networkPolicy returns the typed form of the NetworkPolicy called name.
func networkPolicy(t *testing.T, objects []rendered, name string) *networkingv1.NetworkPolicy {
	t.Helper()
	r := find(objects, "NetworkPolicy", name)
	if r == nil {
		t.Fatalf("NetworkPolicy %s: not rendered", name)
	}
	return r.obj.(*networkingv1.NetworkPolicy)
}

// strictDecoder rejects unknown and duplicated fields for every kind in
// client-go's scheme.
var strictDecoder = serializer.NewCodecFactory(scheme.Scheme, serializer.EnableStrict).UniversalDeserializer()

// decodeStrict returns raw as a well-formed object or the reason it is not.
// A kind the scheme does not know must be one of the recipe's custom kinds.
// For those, only the embedded core types are checked: the controllers'
// own schemas are not modelled here.
func decodeStrict(raw []byte) (rendered, error) {
	var r rendered
	if err := yaml.Unmarshal(raw, &r.PartialObjectMetadata); err != nil {
		return r, err
	}
	obj, _, err := strictDecoder.Decode(raw, nil, nil)
	if err == nil {
		r.obj = obj
		return r, nil
	}
	if !runtime.IsNotRegisteredError(err) {
		return r, err
	}
	switch r.APIVersion + " " + r.Kind {
	case sandboxAPIBeta + " SandboxTemplate", sandboxAPIAlpha + " SandboxTemplate":
		if r.sandboxPod, err = decodeSandboxTemplateSpec(raw); err != nil {
			return r, err
		}
		r.specKeys, err = specKeys(raw)
		return r, err
	case sandboxAPIBeta + " SandboxWarmPool", sandboxAPIAlpha + " SandboxWarmPool":
		r.specKeys, err = specKeys(raw)
		return r, err
	case "networking.gke.io/v1alpha1 FQDNNetworkPolicy":
		return r, nil
	}
	return r, err
}

// specKeys returns an object's top-level spec keys, sorted.
func specKeys(raw []byte) ([]string, error) {
	var doc struct {
		Spec map[string]json.RawMessage `json:"spec"`
	}
	if err := yaml.Unmarshal(raw, &doc); err != nil {
		return nil, err
	}
	keys := slices.Collect(maps.Keys(doc.Spec))
	slices.Sort(keys)
	return keys, nil
}

// decodeSandboxTemplateSpec strict-decodes the pod template and the volume
// claim templates of a SandboxTemplate and returns the pod template. Those
// are core types, where a typo in the sandbox pod spec would otherwise hide.
func decodeSandboxTemplateSpec(raw []byte) (*corev1.PodTemplateSpec, error) {
	var doc struct {
		Spec map[string]json.RawMessage `json:"spec"`
	}
	if err := yaml.Unmarshal(raw, &doc); err != nil {
		return nil, err
	}
	podTemplate, ok := doc.Spec["podTemplate"]
	if !ok {
		return nil, fmt.Errorf("spec.podTemplate is missing")
	}
	var pts corev1.PodTemplateSpec
	if err := yaml.UnmarshalStrict(podTemplate, &pts); err != nil {
		return nil, fmt.Errorf("spec.podTemplate: %w", err)
	}
	var claims []corev1.PersistentVolumeClaim
	if err := yaml.UnmarshalStrict(doc.Spec["volumeClaimTemplates"], &claims); err != nil {
		return nil, fmt.Errorf("spec.volumeClaimTemplates: %w", err)
	}
	return &pts, nil
}

// podSpecs returns the pod spec of every Deployment and SandboxTemplate.
func podSpecs(objects []rendered) []corev1.PodSpec {
	var specs []corev1.PodSpec
	for _, r := range objects {
		if d, ok := r.obj.(*appsv1.Deployment); ok {
			specs = append(specs, d.Spec.Template.Spec)
		}
		if r.sandboxPod != nil {
			specs = append(specs, r.sandboxPod.Spec)
		}
	}
	return specs
}

// assertPublishedImages checks that, left to the chart's defaults, every
// container runs one of the two images a release publishes, at the chart's
// appVersion.
func assertPublishedImages(t *testing.T, objects []rendered, appVersion string) {
	t.Helper()
	for _, spec := range podSpecs(objects) {
		for _, c := range slices.Concat(spec.InitContainers, spec.Containers) {
			want := "us-docker.pkg.dev/anthropic-public-registry/sandboxes/cma-sandbox:" + appVersion
			if c.Name == "dispatcher" {
				want = "us-docker.pkg.dev/anthropic-public-registry/sandboxes/cma-dispatcher:" + appVersion
			}
			if c.Image != want {
				t.Errorf("container %s uses image %q, want %q", c.Name, c.Image, want)
			}
		}
	}
}

// assertBindingFollowsNamespace checks that the cluster-scoped binding's
// name ends in the namespace. Two installs in one
// cluster then hold separate bindings. Both bindings' subjects name the
// namespace too.
func assertBindingFollowsNamespace(t *testing.T, objects []rendered) {
	t.Helper()
	crb := find(objects, "ClusterRoleBinding", "cma-dispatcher-auth-delegator-cma")
	if crb == nil {
		t.Fatal("no ClusterRoleBinding cma-dispatcher-auth-delegator-cma")
	}
	rb := find(objects, "RoleBinding", "cma-dispatcher")
	if rb == nil {
		t.Fatal("no RoleBinding cma-dispatcher")
	}
	for _, subjects := range [][]rbacv1.Subject{
		crb.obj.(*rbacv1.ClusterRoleBinding).Subjects,
		rb.obj.(*rbacv1.RoleBinding).Subjects,
	} {
		for _, s := range subjects {
			if s.Namespace != "cma" {
				t.Errorf("binding subject %s in namespace %q, want cma", s.Name, s.Namespace)
			}
		}
	}
}

// assertSandboxSchema checks that the template and the warm pool carry the
// fields the add-on rejects exactly when they are rendered for upstream's
// v1beta1. The add-on's API server refuses an object with an unknown field
// instead of dropping it, so a v1alpha1 install that keeps one cannot be
// applied.
func assertSandboxSchema(t *testing.T, objects []rendered, sandboxAPIVersion string) {
	t.Helper()
	want := sandboxAPIVersion == sandboxAPIBeta
	for kind, fields := range addonRejectedFields {
		r := find(objects, kind, "cma-sandbox")
		if r == nil {
			t.Errorf("%s cma-sandbox: not rendered", kind)
			continue
		}
		for _, field := range fields {
			if has := slices.Contains(r.specKeys, field); has != want {
				t.Errorf("%s cma-sandbox spec.%s present: %v, want %v for %s", kind, field, has, want, sandboxAPIVersion)
			}
		}
	}
}

// assertDispatcherAPIVersion checks the dispatcher is told the API version
// the SandboxTemplate carries. The env value is the version without the
// group.
func assertDispatcherAPIVersion(t *testing.T, objects []rendered, sandboxAPIVersion string) {
	t.Helper()
	d := find(objects, "Deployment", "cma-dispatcher")
	if d == nil {
		t.Fatal("no Deployment cma-dispatcher")
	}
	got := containerEnv(d.obj.(*appsv1.Deployment), "dispatcher", "CMA_DISPATCHER_API_VERSION")
	want := path.Base(sandboxAPIVersion)
	if got != want {
		t.Errorf("dispatcher env CMA_DISPATCHER_API_VERSION %q, want %q", got, want)
	}
}

// dispatcherRoleGrants reports whether Role cma-dispatcher has a rule that
// names resource and verb under apiGroup.
func dispatcherRoleGrants(t *testing.T, objects []rendered, apiGroup, resource, verb string) bool {
	t.Helper()
	r := find(objects, "Role", "cma-dispatcher")
	if r == nil {
		t.Fatal("no Role cma-dispatcher")
	}
	for _, rule := range r.obj.(*rbacv1.Role).Rules {
		if slices.Contains(rule.APIGroups, apiGroup) && slices.Contains(rule.Resources, resource) && slices.Contains(rule.Verbs, verb) {
			return true
		}
	}
	return false
}

// fetchContainerMounts reports whether the fetch container of SandboxTemplate
// cma-sandbox has a volume mount at mountPath. fetch is the template's first container; any
// other container there fails the test.
func fetchContainerMounts(t *testing.T, objects []rendered, mountPath string) bool {
	t.Helper()
	st := find(objects, "SandboxTemplate", "cma-sandbox")
	if st == nil || st.sandboxPod == nil || len(st.sandboxPod.Spec.Containers) == 0 {
		t.Fatal("no SandboxTemplate cma-sandbox with a container")
	}
	fetch := st.sandboxPod.Spec.Containers[0]
	if fetch.Name != "fetch" {
		t.Fatalf("SandboxTemplate cma-sandbox: first container is %q, want fetch", fetch.Name)
	}
	for _, m := range fetch.VolumeMounts {
		if m.MountPath == mountPath {
			return true
		}
	}
	return false
}

// containerEnv returns the value of env name in the named container, or ""
// when either is absent.
func containerEnv(d *appsv1.Deployment, container, name string) string {
	for _, c := range d.Spec.Template.Spec.Containers {
		if c.Name != container {
			continue
		}
		for _, e := range c.Env {
			if e.Name == name {
				return e.Value
			}
		}
	}
	return ""
}

func find(objects []rendered, kind, name string) *rendered {
	for i := range objects {
		if objects[i].Kind == kind && objects[i].Name == name {
			return &objects[i]
		}
	}
	return nil
}

// examplePlatform returns the platform a values file sets.
func examplePlatform(t *testing.T, valuesFile string) string {
	t.Helper()
	raw, err := os.ReadFile(valuesFile)
	if err != nil {
		t.Fatal(err)
	}
	var values struct {
		Platform string `json:"platform"`
	}
	if err := yaml.Unmarshal(raw, &values); err != nil {
		t.Fatalf("%s: %v", valuesFile, err)
	}
	return values.Platform
}
