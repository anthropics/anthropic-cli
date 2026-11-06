package preflight

import (
	"bytes"
	"cmp"
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"
	"testing"

	appsv1 "k8s.io/api/apps/v1"
	authzv1 "k8s.io/api/authorization/v1"
	corev1 "k8s.io/api/core/v1"
	nodev1 "k8s.io/api/node/v1"
	storagev1 "k8s.io/api/storage/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	fakediscovery "k8s.io/client-go/discovery/fake"
	dynamicfake "k8s.io/client-go/dynamic/fake"
	kubefake "k8s.io/client-go/kubernetes/fake"
	k8stesting "k8s.io/client-go/testing"

	"github.com/anthropics/anthropic-cli/sandboxes/dispatcher/internal/sandbox"
)

type verifierFunc func(context.Context) error

func (f verifierFunc) VerifyEnvironment(ctx context.Context) error { return f(ctx) }

// access is one (resource, verb) pair the SubjectAccessReview fake denies, so
// a test can fail the role check on that pair alone.
type access struct{ resource, verb string }

type cluster struct {
	handler      string
	psaLabel     string
	replicas     int32
	strategy     appsv1.DeploymentStrategyType
	reclaim      corev1.PersistentVolumeReclaimPolicy
	denyVerb     string
	denyResource string
	denyAccess   access
	// reviewsForbidden makes the cluster refuse every SubjectAccessReview.
	reviewsForbidden bool
	// unservedGroupVersion, when set, names the one GroupVersion the served API
	// omits.
	unservedGroupVersion string
	// missingResource, when set, names a resource the served API omits.
	missingResource string
	storageClass    string
	// defaultClassAnnotation, when set, marks the fast class default under that annotation.
	defaultClassAnnotation string
	// secondStorageClass, when set, names a second volume claim template's class.
	secondStorageClass string
	noClaimTemplates   bool
	// noTemplate, noWarmPool and noDeployment drop that object from the
	// cluster. The checker still looks for it under its healthy name.
	noTemplate   bool
	noWarmPool   bool
	noDeployment bool
	poolTemplate string
	// runtimeClassName is the class the template names. The cluster has only
	// gvisor.
	runtimeClassName  string
	networkPolicyMode string
	// worker is the worker container's spec. fetchEnv and fetchMounts extend
	// the fetch container's env and volumeMounts. fetchWithoutToken drops
	// fetch's projected token mount. noFetch drops the whole fetch container.
	worker            map[string]any
	fetchEnv          []any
	fetchMounts       []any
	fetchWithoutToken bool
	noFetch           bool
	sidecars          []any
	initContainers    []any
	// shareProcesses and automountToken are the pod spec's values of those
	// fields. Either takes any type, so a case can set a value the pod spec
	// cannot convert.
	shareProcesses any
	automountToken any
	extraVolumes   []any
	deploymentSA   string
	tokenAudience  string
	serviceHost    string
	// dispatcherSA is the account -dispatcher-service-account names.
	dispatcherSA string
	// networkProbe runs the network check instead of skipping it.
	networkProbe bool
	// verifyErr, when set, is the error the environment verifier returns.
	verifyErr error
}

// directCluster describes the add-on install: a v1alpha1 SandboxTemplate with
// no kubelet credential and cpu+memory limits on both containers.
type directCluster struct {
	// handler, psaLabel, replicas and strategy take the healthy values when
	// unset.
	handler, psaLabel         string
	replicas                  int32
	strategy                  appsv1.DeploymentStrategyType
	deploymentSA              string
	verifyErr                 error
	missingResource           string
	fetchLimits, workerLimits map[string]any
	fetchProjectsToken        bool
	automountToken            any
	denyVerb, denyResource    string
	denyAccess                access
	initContainers            []any
	rejectSandbox             bool
	networkPolicyMode         string
	extraVolumes              []any
	workerName                string
	workerEnv, workerEnvFrom  []any
	workerMounts              []any
	fetchEnv, fetchMounts     []any
	shareProcesses            bool
	// dropsSpecField, when set, names a spec field the served schema prunes
	// from a created Sandbox.
	dropsSpecField string
	// dispatcherSA is the account -dispatcher-service-account names.
	dispatcherSA string
	// workspaceClaim, when set, is the per-session claim the checker is
	// configured with.
	workspaceClaim *sandbox.WorkspaceClaim
	// storageClass, when set, is the one StorageClass the cluster has.
	storageClass *storagev1.StorageClass
	// noFetch drops the fetch container from the template.
	noFetch bool
}

// byHand ends both network check lines. healthy sets serviceHost to the host
// byHand names.
const byHand = "verify by hand that a sandbox pod cannot reach another sandbox pod, 169.254.169.254 or the cluster's API server and can reach the dispatcher at cma-dispatcher.cma.svc:8443"

// skippedProbe is the network check's line when -skip-network-probe is set.
const skippedProbe = "warn  network policy isolates sandboxes: warning: skipped by -skip-network-probe; " + byHand

func healthy() cluster {
	return cluster{handler: "runsc", psaLabel: "restricted", replicas: 1, strategy: appsv1.RecreateDeploymentStrategyType, reclaim: corev1.PersistentVolumeReclaimDelete, storageClass: "fast", poolTemplate: "cma-sandbox", runtimeClassName: "gvisor", networkPolicyMode: "Unmanaged", worker: healthyWorker(), shareProcesses: false, automountToken: false, deploymentSA: "cma-dispatcher", dispatcherSA: "cma-dispatcher", tokenAudience: "cma-dispatcher", serviceHost: "cma-dispatcher.cma.svc:8443"}
}

func healthyWorker() map[string]any {
	return map[string]any{
		"name":         "worker",
		"volumeMounts": []any{map[string]any{"name": "workspace", "mountPath": "/workspace"}, map[string]any{"name": "shared", "mountPath": "/run/cma"}},
	}
}

func run(t *testing.T, c cluster) string {
	out, _ := runWithCount(t, c)
	return out
}

// runWithCount returns the output and the failure count Run reports.
func runWithCount(t *testing.T, c cluster) (string, int) {
	t.Helper()
	shape, _ := sandbox.ShapeFor("v1beta1")
	classAnnotations := map[string]string{}
	if c.defaultClassAnnotation != "" {
		classAnnotations[c.defaultClassAnnotation] = "true"
	}
	objects := []runtime.Object{
		&nodev1.RuntimeClass{ObjectMeta: metav1.ObjectMeta{Name: "gvisor"}, Handler: c.handler},
		&storagev1.StorageClass{ObjectMeta: metav1.ObjectMeta{Name: "fast", Annotations: classAnnotations}, ReclaimPolicy: ptr(c.reclaim)},
		&corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "ns", Labels: map[string]string{"pod-security.kubernetes.io/enforce": c.psaLabel}}},
		&corev1.Secret{ObjectMeta: metav1.ObjectMeta{Name: "env-key", Namespace: "ns"}, Data: map[string][]byte{"ANTHROPIC_ENVIRONMENT_KEY": []byte("k")}},
		&corev1.Secret{ObjectMeta: metav1.ObjectMeta{Name: "legacy-token", Namespace: "ns"}, Type: corev1.SecretTypeServiceAccountToken, Data: map[string][]byte{"token": []byte("t")}},
		&corev1.Secret{ObjectMeta: metav1.ObjectMeta{Name: "ordinary", Namespace: "ns"}, Type: corev1.SecretTypeOpaque, Data: map[string][]byte{"value": []byte("v")}},
		&corev1.ConfigMap{ObjectMeta: metav1.ObjectMeta{Name: "ca", Namespace: "ns"}, Data: map[string]string{"ca.crt": "pem"}},
		&corev1.ConfigMap{ObjectMeta: metav1.ObjectMeta{Name: "key-config", Namespace: "ns"}, Data: map[string]string{"ANTHROPIC_ENVIRONMENT_KEY": "k"}},
		&corev1.ConfigMap{ObjectMeta: metav1.ObjectMeta{Name: "binary-key-config", Namespace: "ns"}, BinaryData: map[string][]byte{"ANTHROPIC_ENVIRONMENT_KEY": []byte("k")}},
		&corev1.ConfigMap{ObjectMeta: metav1.ObjectMeta{Name: "ordinary-config", Namespace: "ns"}, Data: map[string]string{"setting": "on"}},
	}
	if !c.noDeployment {
		objects = append(objects, &appsv1.Deployment{ObjectMeta: metav1.ObjectMeta{Name: "cma-dispatcher", Namespace: "ns"},
			Spec: appsv1.DeploymentSpec{Replicas: ptr(c.replicas), Strategy: appsv1.DeploymentStrategy{Type: c.strategy},
				Template: corev1.PodTemplateSpec{Spec: corev1.PodSpec{ServiceAccountName: c.deploymentSA}}}})
	}
	kube := kubefake.NewSimpleClientset(objects...)
	served := []*metav1.APIResourceList{
		{GroupVersion: "agents.x-k8s.io/v1beta1", APIResources: []metav1.APIResource{{Name: "sandboxes"}}},
		{GroupVersion: "extensions.agents.x-k8s.io/v1beta1", APIResources: []metav1.APIResource{{Name: "sandboxclaims"}, {Name: "sandboxtemplates"}, {Name: "sandboxwarmpools"}}},
	}
	served = slices.DeleteFunc(served, func(list *metav1.APIResourceList) bool { return list.GroupVersion == c.unservedGroupVersion })
	for _, list := range served {
		list.APIResources = slices.DeleteFunc(list.APIResources, func(r metav1.APIResource) bool { return r.Name == c.missingResource })
	}
	kube.Discovery().(*fakediscovery.FakeDiscovery).Resources = served
	kube.PrependReactor("create", "subjectaccessreviews", reviewReactor(c.dispatcherSA, c.reviewsForbidden, c.denyVerb, c.denyResource, c.denyAccess))
	volumeClaimTemplates := []any{map[string]any{"spec": map[string]any{"storageClassName": c.storageClass}}}
	if c.secondStorageClass != "" {
		volumeClaimTemplates = append(volumeClaimTemplates, map[string]any{"spec": map[string]any{"storageClassName": c.secondStorageClass}})
	}
	if c.noClaimTemplates {
		volumeClaimTemplates = nil
	}
	fetchMounts := append([]any{map[string]any{"name": "dispatcher-token", "mountPath": "/var/run/secrets/cma-dispatcher"}, map[string]any{"name": "dispatcher-ca", "mountPath": "/etc/cma-dispatcher"}, map[string]any{"name": "shared", "mountPath": "/run/cma"}}, c.fetchMounts...)
	if c.fetchWithoutToken {
		fetchMounts = slices.DeleteFunc(fetchMounts, func(m any) bool { return m.(map[string]any)["name"] == "dispatcher-token" })
	}
	fetch := map[string]any{
		"name":         "fetch",
		"env":          append([]any{map[string]any{"name": "CMA_DISPATCHER_URL", "value": "https://d"}}, c.fetchEnv...),
		"volumeMounts": fetchMounts,
	}
	containers := append([]any{fetch, c.worker}, c.sidecars...)
	if c.noFetch {
		containers = containers[1:]
	}
	tpl := &unstructured.Unstructured{Object: map[string]any{
		"apiVersion": "extensions.agents.x-k8s.io/v1beta1", "kind": "SandboxTemplate",
		"metadata": map[string]any{"name": "cma-sandbox", "namespace": "ns"},
		"spec": map[string]any{
			"networkPolicyManagement": c.networkPolicyMode,
			"podTemplate": map[string]any{
				"metadata": map[string]any{"labels": map[string]any{"anthropic.com/role": "sandbox"}},
				"spec": map[string]any{
					"runtimeClassName":             c.runtimeClassName,
					"shareProcessNamespace":        c.shareProcesses,
					"automountServiceAccountToken": c.automountToken,
					"initContainers":               c.initContainers,
					"containers":                   containers,
					"volumes": append([]any{
						map[string]any{"name": "dispatcher-token", "projected": map[string]any{"sources": []any{map[string]any{"serviceAccountToken": map[string]any{"audience": c.tokenAudience, "path": "token"}}}}},
						map[string]any{"name": "dispatcher-ca", "configMap": map[string]any{"name": "ca"}},
						map[string]any{"name": "shared", "emptyDir": map[string]any{"medium": "Memory"}},
					}, c.extraVolumes...),
				},
			},
			"volumeClaimTemplates": volumeClaimTemplates,
		},
	}}
	pool := &unstructured.Unstructured{Object: map[string]any{
		"apiVersion": "extensions.agents.x-k8s.io/v1beta1", "kind": "SandboxWarmPool",
		"metadata": map[string]any{"name": "cma-sandbox", "namespace": "ns"},
		"spec":     map[string]any{"sandboxTemplateRef": map[string]any{"name": c.poolTemplate}},
	}}
	var stored []runtime.Object
	if !c.noTemplate {
		stored = append(stored, tpl)
	}
	if !c.noWarmPool {
		stored = append(stored, pool)
	}
	dyn := dynamicfake.NewSimpleDynamicClient(runtime.NewScheme(), stored...)
	checker := &Checker{Kube: kube, Dyn: dyn, API: verifierFunc(func(context.Context) error { return c.verifyErr }), Options: Options{
		Namespace: "ns", Shape: shape, SandboxTemplateName: "cma-sandbox", WarmPoolName: "cma-sandbox", DispatcherServiceAccount: c.dispatcherSA, DeploymentName: "cma-dispatcher", TokenAudience: "cma-dispatcher", ServiceHost: c.serviceHost, SkipNetworkProbe: !c.networkProbe,
	}}
	var out bytes.Buffer
	failures := checker.Run(context.Background(), &out)
	return out.String(), failures
}

func healthyDirect() directCluster {
	limits := map[string]any{"limits": map[string]any{"cpu": "10m", "memory": "64Mi"}}
	return directCluster{fetchLimits: limits, workerLimits: limits, workerName: "worker", automountToken: false, networkPolicyMode: "Unmanaged", deploymentSA: "cma-dispatcher", dispatcherSA: "cma-dispatcher", storageClass: defaultStorageClass(corev1.PersistentVolumeReclaimDelete)}
}

func defaultStorageClass(reclaim corev1.PersistentVolumeReclaimPolicy) *storagev1.StorageClass {
	return &storagev1.StorageClass{ObjectMeta: metav1.ObjectMeta{Name: "standard", Annotations: map[string]string{"storageclass.kubernetes.io/is-default-class": "true"}}, ReclaimPolicy: ptr(reclaim)}
}

func runDirect(t *testing.T, c directCluster) string {
	out, _ := runDirectWithCount(t, c)
	return out
}

// runDirectWithCount returns the output and the failure count Run reports.
func runDirectWithCount(t *testing.T, c directCluster) (string, int) {
	t.Helper()
	shape, _ := sandbox.ShapeFor("v1alpha1")
	kube := kubefake.NewSimpleClientset(
		&nodev1.RuntimeClass{ObjectMeta: metav1.ObjectMeta{Name: "gvisor"}, Handler: cmp.Or(c.handler, "runsc")},
		&corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "ns", Labels: map[string]string{"pod-security.kubernetes.io/enforce": cmp.Or(c.psaLabel, "restricted")}}},
		&corev1.ConfigMap{ObjectMeta: metav1.ObjectMeta{Name: "cma-dispatcher-ca", Namespace: "ns"}, Data: map[string]string{"ca.crt": "pem"}},
		&corev1.Secret{ObjectMeta: metav1.ObjectMeta{Name: "legacy-token", Namespace: "ns"}, Type: corev1.SecretTypeServiceAccountToken, Data: map[string][]byte{"token": []byte("t")}},
		&appsv1.Deployment{ObjectMeta: metav1.ObjectMeta{Name: "cma-dispatcher", Namespace: "ns"},
			Spec: appsv1.DeploymentSpec{Replicas: ptr(cmp.Or(c.replicas, 1)), Strategy: appsv1.DeploymentStrategy{Type: cmp.Or(c.strategy, appsv1.RecreateDeploymentStrategyType)},
				Template: corev1.PodTemplateSpec{Spec: corev1.PodSpec{ServiceAccountName: c.deploymentSA}}}},
	)
	if c.storageClass != nil {
		if _, err := kube.StorageV1().StorageClasses().Create(context.Background(), c.storageClass, metav1.CreateOptions{}); err != nil {
			t.Fatal(err)
		}
	}
	kube.Discovery().(*fakediscovery.FakeDiscovery).Resources = []*metav1.APIResourceList{
		{GroupVersion: "agents.x-k8s.io/v1alpha1", APIResources: []metav1.APIResource{{Name: "sandboxes"}}},
		{GroupVersion: "extensions.agents.x-k8s.io/v1alpha1", APIResources: []metav1.APIResource{{Name: "sandboxtemplates"}}},
	}
	for _, list := range kube.Discovery().(*fakediscovery.FakeDiscovery).Resources {
		list.APIResources = slices.DeleteFunc(list.APIResources, func(r metav1.APIResource) bool { return r.Name == c.missingResource })
	}
	kube.PrependReactor("create", "subjectaccessreviews", reviewReactor(c.dispatcherSA, false, c.denyVerb, c.denyResource, c.denyAccess))
	fetch := map[string]any{
		"name":         "fetch",
		"env":          append([]any{map[string]any{"name": "CMA_DISPATCHER_URL", "value": "https://d"}}, c.fetchEnv...),
		"volumeMounts": append([]any{map[string]any{"name": "ca", "mountPath": "/etc/cma-dispatcher"}}, c.fetchMounts...),
	}
	for k, v := range c.fetchLimits {
		fetch["resources"] = map[string]any{k: v}
	}
	worker := map[string]any{"name": c.workerName}
	for k, v := range c.workerLimits {
		worker["resources"] = map[string]any{k: v}
	}
	if c.workerEnv != nil {
		worker["env"] = c.workerEnv
	}
	if c.workerEnvFrom != nil {
		worker["envFrom"] = c.workerEnvFrom
	}
	if c.workerMounts != nil {
		worker["volumeMounts"] = c.workerMounts
	}
	containers := []any{fetch, worker}
	if c.noFetch {
		containers = containers[1:]
	}
	volumes := append([]any{map[string]any{"name": "ca", "configMap": map[string]any{"name": "cma-dispatcher-ca"}}}, c.extraVolumes...)
	if c.fetchProjectsToken {
		volumes = append(volumes, map[string]any{"name": "tok", "projected": map[string]any{"sources": []any{map[string]any{"serviceAccountToken": map[string]any{"audience": "cma-dispatcher", "path": "token"}}}}})
	}
	tpl := &unstructured.Unstructured{Object: map[string]any{
		"apiVersion": "extensions.agents.x-k8s.io/v1alpha1", "kind": "SandboxTemplate",
		"metadata": map[string]any{"name": "cma-sandbox", "namespace": "ns"},
		"spec": map[string]any{
			"networkPolicyManagement": c.networkPolicyMode,
			"podTemplate": map[string]any{
				"spec": map[string]any{
					"runtimeClassName":             "gvisor",
					"shareProcessNamespace":        c.shareProcesses,
					"automountServiceAccountToken": c.automountToken,
					"containers":                   containers,
					"volumes":                      volumes,
					"initContainers":               c.initContainers,
				},
			},
		},
	}}
	dyn := dynamicfake.NewSimpleDynamicClientWithCustomListKinds(runtime.NewScheme(), map[schema.GroupVersionResource]string{
		shape.Templates(): "SandboxTemplateList",
		shape.Sandboxes(): "SandboxList",
	})
	if _, err := dyn.Resource(shape.Templates()).Namespace("ns").Create(context.Background(), tpl, metav1.CreateOptions{}); err != nil {
		t.Fatal(err)
	}
	if c.rejectSandbox {
		dyn.PrependReactor("create", "sandboxes", func(k8stesting.Action) (bool, runtime.Object, error) {
			return true, nil, errors.New("Security Violation: All containers must define resource limits")
		})
	}
	if c.dropsSpecField != "" {
		dyn.PrependReactor("create", "sandboxes", func(action k8stesting.Action) (bool, runtime.Object, error) {
			stored := action.(k8stesting.CreateAction).GetObject().DeepCopyObject().(*unstructured.Unstructured)
			unstructured.RemoveNestedField(stored.Object, "spec", c.dropsSpecField)
			return true, stored, nil
		})
	}
	checker := &Checker{Kube: kube, Dyn: dyn, API: verifierFunc(func(context.Context) error { return nil }), Options: Options{
		Namespace: "ns", Shape: shape, SandboxTemplateName: "cma-sandbox", DispatcherServiceAccount: c.dispatcherSA, DeploymentName: "cma-dispatcher", TokenAudience: "cma-dispatcher", SkipNetworkProbe: true,
		WorkspaceClaim: c.workspaceClaim,
	}}
	if c.verifyErr != nil {
		checker.API = verifierFunc(func(context.Context) error { return c.verifyErr })
	}
	var out bytes.Buffer
	failures := checker.Run(context.Background(), &out)
	return out.String(), failures
}

// reviewReactor refuses every SubjectAccessReview with forbidden set, else
// allows one about dispatcherSA unless it names a denied verb, resource or
// pair. One about another identity gets an error that fails the case's run.
func reviewReactor(dispatcherSA string, forbidden bool, denyVerb, denyResource string, denyAccess access) k8stesting.ReactionFunc {
	return func(action k8stesting.Action) (bool, runtime.Object, error) {
		if forbidden {
			return true, nil, apierrors.NewForbidden(authzv1.Resource("subjectaccessreviews"), "", errors.New("denied"))
		}
		review := action.(k8stesting.CreateAction).GetObject().(*authzv1.SubjectAccessReview)
		if err := reviewAsksAbout(review, dispatcherSA); err != nil {
			return true, nil, err
		}
		attrs := review.Spec.ResourceAttributes
		review.Status.Allowed = attrs.Verb != denyVerb && attrs.Resource != denyResource && access{attrs.Resource, attrs.Verb} != denyAccess
		return true, review, nil
	}
}

// reviewAsksAbout returns an error unless review asks about dispatcherSA in
// namespace ns, with the three groups every service account carries. A review
// about tokenreviews names no namespace. TokenReview is cluster-scoped.
func reviewAsksAbout(review *authzv1.SubjectAccessReview, dispatcherSA string) error {
	wantUser := "system:serviceaccount:ns:" + dispatcherSA
	if review.Spec.User != wantUser {
		return fmt.Errorf("the review asks about user %q, not %q", review.Spec.User, wantUser)
	}
	wantGroups := []string{"system:serviceaccounts", "system:serviceaccounts:ns", "system:authenticated"}
	if !slices.Equal(review.Spec.Groups, wantGroups) {
		return fmt.Errorf("the review asks about groups %q, not %q", review.Spec.Groups, wantGroups)
	}
	attrs := review.Spec.ResourceAttributes
	wantNamespace := "ns"
	if attrs.Resource == "tokenreviews" {
		wantNamespace = ""
	}
	if attrs.Namespace != wantNamespace {
		return fmt.Errorf("the review asks about %s in namespace %q, not %q", attrs.Resource, attrs.Namespace, wantNamespace)
	}
	return nil
}

// wantLine fails unless out has a line equal to want. Matching the whole line
// fails a case whose check keeps its name but changes its reason. A want that
// is not a FAIL line also asserts that no check failed.
func wantLine(t *testing.T, name, out, want string) {
	t.Helper()
	if !slices.Contains(strings.Split(strings.TrimSuffix(out, "\n"), "\n"), want) {
		t.Errorf("%s: want the line %q in:\n%s", name, want, out)
	}
	if !strings.HasPrefix(want, "FAIL  ") && strings.Contains(out, "FAIL  ") {
		t.Errorf("%s: want no FAIL line beside %q, got:\n%s", name, want, out)
	}
}

func ptr[T any](v T) *T { return &v }
