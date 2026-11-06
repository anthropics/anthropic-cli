package preflight

import (
	"context"
	"errors"
	"fmt"
	"slices"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/anthropics/anthropic-cli/sandboxes/dispatcher/internal/sandbox"
)

// FetchContainer collects the work item for the pod; it alone holds the
// dispatcher's CA, and on the upstream path the projected token.
const FetchContainer = sandbox.FetchContainer

// dispatcherVariables are the variables only the fetch container may carry.
var dispatcherVariables = []string{"CMA_DISPATCHER_URL", "CMA_SA_TOKEN_FILE", "CMA_CA_FILE"}

const serviceAccountTokenSource = "service account token"

// containers requires fetch and worker, and lets no other container mount what
// fetch's volumes read or set the dispatcher variables. It refuses a shared
// process namespace, which opens fetch's files, and an automounted API token.
func (c *Checker) containers(ctx context.Context) error {
	spec, err := c.templatePodSpec(ctx)
	if err != nil {
		return err
	}
	if spec.ShareProcessNamespace != nil && *spec.ShareProcessNamespace {
		return fmt.Errorf("the SandboxTemplate sets shareProcessNamespace, which lets every container read %q's token", FetchContainer)
	}
	if spec.AutomountServiceAccountToken == nil || *spec.AutomountServiceAccountToken {
		return errors.New("the SandboxTemplate does not set automountServiceAccountToken: false, so every container would get a token for the Kubernetes API")
	}
	fetch := containerNamed(spec, FetchContainer)
	if fetch == nil || containerNamed(spec, sandbox.WorkerContainer) == nil {
		return fmt.Errorf("the SandboxTemplate must have containers named %q and %q", FetchContainer, sandbox.WorkerContainer)
	}
	volumes, err := comparableVolumes(spec)
	if err != nil {
		return err
	}
	for _, container := range slices.Concat(spec.InitContainers, spec.Containers) {
		if err := c.environmentKeyStaysOut(ctx, container, volumes); err != nil {
			return err
		}
		if err := c.serviceAccountTokenSecretsStayOut(ctx, container, volumes); err != nil {
			return err
		}
	}
	fetchOnly := map[string]bool{}
	for _, m := range fetch.VolumeMounts {
		if v := volumes[m.Name]; v.EmptyDir == nil {
			fetchOnly[m.Name] = true
			for _, src := range volumeSources(v) {
				fetchOnly[src] = true
			}
			if err := tokenAudienceMatches(v, c.Options.TokenAudience); err != nil {
				return err
			}
		}
	}
	if !fetchOnly[serviceAccountTokenSource] {
		return fmt.Errorf("container %q mounts no projected service-account token", FetchContainer)
	}
	for _, other := range slices.Concat(spec.InitContainers, spec.Containers) {
		if other.Name == FetchContainer {
			continue
		}
		for _, m := range other.VolumeMounts {
			v := volumes[m.Name]
			if fetchOnly[m.Name] || slices.ContainsFunc(volumeSources(v), func(src string) bool { return fetchOnly[src] }) {
				return fmt.Errorf("container %q mounts volume %q, which reads what only %q may use", other.Name, m.Name, FetchContainer)
			}
		}
		if err := dispatcherVariablesStayOut(other); err != nil {
			return err
		}
		for _, e := range other.EnvFrom {
			if src := envFromSource(e); fetchOnly[src] {
				return fmt.Errorf("container %q reads %s through envFrom, which only %q may use", other.Name, src, FetchContainer)
			}
		}
		for _, e := range other.Env {
			if src := envValueSource(e); fetchOnly[src] {
				return fmt.Errorf("container %q reads %s through an env valueFrom, which only %q may use", other.Name, src, FetchContainer)
			}
		}
	}
	return nil
}

// templatePodSpec reads the SandboxTemplate's pod spec as a typed value.
func (c *Checker) templatePodSpec(ctx context.Context) (corev1.PodSpec, error) {
	tpl, err := c.template(ctx)
	if err != nil {
		return corev1.PodSpec{}, err
	}
	var spec corev1.PodSpec
	podSpec, _, _ := unstructured.NestedMap(tpl.Object, "spec", "podTemplate", "spec")
	if err := runtime.DefaultUnstructuredConverter.FromUnstructured(podSpec, &spec); err != nil {
		return corev1.PodSpec{}, fmt.Errorf("reading the SandboxTemplate's pod spec: %w", err)
	}
	return spec, nil
}

func containerNamed(spec corev1.PodSpec, name string) *corev1.Container {
	for i := range spec.Containers {
		if spec.Containers[i].Name == name {
			return &spec.Containers[i]
		}
	}
	return nil
}

// comparableVolumes maps the spec's volumes by name for the environment-key and
// token scans. It rejects any source volumeSources cannot inspect, since that
// could carry the key unseen, and a projected volume with no source it reads.
func comparableVolumes(spec corev1.PodSpec) (map[string]corev1.Volume, error) {
	volumes := map[string]corev1.Volume{}
	for _, v := range spec.Volumes {
		if v.HostPath != nil {
			return nil, fmt.Errorf("volume %q uses a hostPath, which the restricted pod security standard forbids and this check cannot confine", v.Name)
		}
		if v.EmptyDir == nil && len(volumeSources(v)) == 0 {
			return nil, fmt.Errorf("volume %q has a source this check cannot compare; use an emptyDir, configMap, secret or persistentVolumeClaim volume, or a projected volume with a configMap, secret or service-account-token source", v.Name)
		}
		volumes[v.Name] = v
	}
	return volumes, nil
}

// volumeSources names what a volume reads, so two volumes that read the same
// ConfigMap, Secret, claim or host path match whatever they are called.
func volumeSources(v corev1.Volume) []string {
	var sources []string
	switch {
	case v.ConfigMap != nil:
		sources = append(sources, "configmap "+v.ConfigMap.Name)
	case v.Secret != nil:
		sources = append(sources, "secret "+v.Secret.SecretName)
	case v.PersistentVolumeClaim != nil:
		sources = append(sources, "claim "+v.PersistentVolumeClaim.ClaimName)
	case v.Projected != nil:
		for _, src := range v.Projected.Sources {
			switch {
			case src.ServiceAccountToken != nil:
				sources = append(sources, serviceAccountTokenSource)
			case src.ConfigMap != nil:
				sources = append(sources, "configmap "+src.ConfigMap.Name)
			case src.Secret != nil:
				sources = append(sources, "secret "+src.Secret.Name)
			}
		}
	}
	return sources
}

// tokenAudienceMatches fails when a projected service-account token in v is
// minted for any audience but the dispatcher's: an unset audience is a
// credential for the Kubernetes API itself.
func tokenAudienceMatches(v corev1.Volume, audience string) error {
	if v.Projected == nil {
		return nil
	}
	for _, src := range v.Projected.Sources {
		if src.ServiceAccountToken != nil && src.ServiceAccountToken.Audience != audience {
			return fmt.Errorf("volume %q projects a service-account token for audience %q, not %q; an unset audience is a token for the Kubernetes API itself", v.Name, src.ServiceAccountToken.Audience, audience)
		}
	}
	return nil
}

// dispatcherVariablesStayOut fails when container sets any of
// dispatcherVariables, by value or valueFrom. Only env names are read: an
// envFrom Secret or ConfigMap that holds such a key is not seen.
func dispatcherVariablesStayOut(container corev1.Container) error {
	for _, e := range container.Env {
		if slices.Contains(dispatcherVariables, e.Name) {
			return fmt.Errorf("container %q sets %s, which only %q may carry", container.Name, e.Name, FetchContainer)
		}
	}
	return nil
}

func envFromSource(e corev1.EnvFromSource) string {
	switch {
	case e.ConfigMapRef != nil:
		return "configmap " + e.ConfigMapRef.Name
	case e.SecretRef != nil:
		return "secret " + e.SecretRef.Name
	}
	return ""
}

func envValueSource(e corev1.EnvVar) string {
	if e.ValueFrom == nil {
		return ""
	}
	switch {
	case e.ValueFrom.ConfigMapKeyRef != nil:
		return "configmap " + e.ValueFrom.ConfigMapKeyRef.Name
	case e.ValueFrom.SecretKeyRef != nil:
		return "secret " + e.ValueFrom.SecretKeyRef.Name
	}
	return ""
}
