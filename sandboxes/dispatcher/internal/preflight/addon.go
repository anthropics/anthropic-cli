package preflight

import (
	"context"
	"errors"
	"fmt"
	"slices"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/anthropics/anthropic-cli/sandboxes/dispatcher/internal/sandbox"
)

// directContainers checks the add-on template's containers: fetch and worker
// present, cpu and memory limits on each, and no forbidden credential,
// environment key or long-lived Kubernetes token in reach of any of them.
func (c *Checker) directContainers(ctx context.Context) error {
	spec, err := c.templatePodSpec(ctx)
	if err != nil {
		return err
	}
	if spec.ShareProcessNamespace != nil && *spec.ShareProcessNamespace {
		return fmt.Errorf("the SandboxTemplate sets shareProcessNamespace, which lets every container read %q's nonce", FetchContainer)
	}
	if spec.AutomountServiceAccountToken == nil || *spec.AutomountServiceAccountToken {
		return errors.New("the SandboxTemplate does not set automountServiceAccountToken: false, which the add-on's policy requires")
	}
	if containerNamed(spec, FetchContainer) == nil || containerNamed(spec, sandbox.WorkerContainer) == nil {
		return fmt.Errorf("the SandboxTemplate must have containers named %q and %q", FetchContainer, sandbox.WorkerContainer)
	}
	for _, v := range spec.Volumes {
		if bannedProjectedSource(v) {
			return fmt.Errorf("volume %q projects a service-account token, ClusterTrustBundle, or PodCertificate, which the add-on's policy forbids in a sandbox pod", v.Name)
		}
	}
	volumes, err := comparableVolumes(spec)
	if err != nil {
		return err
	}
	for _, container := range slices.Concat(spec.InitContainers, spec.Containers) {
		if err := requireResourceLimits(container); err != nil {
			return err
		}
		if err := c.environmentKeyStaysOut(ctx, container, volumes); err != nil {
			return err
		}
		if err := c.serviceAccountTokenSecretsStayOut(ctx, container, volumes); err != nil {
			return err
		}
		if container.Name == FetchContainer {
			continue
		}
		if err := dispatcherVariablesStayOut(container); err != nil {
			return err
		}
	}
	return nil
}

// directWorkspace reports how /workspace is backed, since only a claim survives
// a suspend. The dispatcher's claim needs its StorageClass to exist: the dry
// run in admitsSandbox provisions nothing, so without it pods stay Pending.
func (c *Checker) directWorkspace(ctx context.Context) error {
	spec, err := c.templatePodSpec(ctx)
	if err != nil {
		return err
	}
	for _, v := range spec.Volumes {
		if v.Name != sandbox.WorkspaceVolume {
			continue
		}
		if v.PersistentVolumeClaim != nil {
			return nil
		}
		return fmt.Errorf("%w: the template's %s volume is not claim-backed, so a suspend discards /workspace", ErrWarning, sandbox.WorkspaceVolume)
	}
	if c.Options.WorkspaceClaim != nil {
		return c.claimStorageClass(ctx, c.Options.WorkspaceClaim.StorageClass)
	}
	return fmt.Errorf("%w: the stored template carries no %s volume and no workspace claim size is configured, so the dispatcher backs /workspace with an emptyDir that a suspend discards", ErrWarning, sandbox.WorkspaceVolume)
}

// admitsSandbox dry-run-creates the Sandbox the dispatcher would build, so the
// add-on's admission policy sees the real shape before any session runs. A
// field the served schema drops fails it: that is a mismatched -api-version.
func (c *Checker) admitsSandbox(ctx context.Context) error {
	tpl, err := c.template(ctx)
	if err != nil {
		return err
	}
	obj, err := sandbox.SandboxFromTemplate(c.Options.Shape, tpl, c.Options.Namespace, "cma-preflight-probe", nil, c.Options.WorkspaceClaim)
	if err != nil {
		return err
	}
	created, err := c.Dyn.Resource(c.Options.Shape.Sandboxes()).Namespace(c.Options.Namespace).Create(ctx, obj, metav1.CreateOptions{DryRun: []string{metav1.DryRunAll}})
	if err != nil {
		return fmt.Errorf("a dry-run sandbox create was refused, so the add-on would reject a real one: %w", err)
	}
	if created == nil {
		return nil
	}
	wrote, _, _ := unstructured.NestedMap(obj.Object, "spec")
	stored, _, _ := unstructured.NestedMap(created.Object, "spec")
	for key := range wrote {
		if _, ok := stored[key]; !ok {
			return fmt.Errorf("the served Sandbox schema dropped spec.%s; check that -api-version %s matches the cluster", key, c.Options.Shape.Version())
		}
	}
	return nil
}

// bannedProjectedSource reports whether v projects a credential the add-on's
// policy forbids in a sandbox pod.
func bannedProjectedSource(v corev1.Volume) bool {
	if v.Projected == nil {
		return false
	}
	for _, src := range v.Projected.Sources {
		if src.ServiceAccountToken != nil || src.ClusterTrustBundle != nil || src.PodCertificate != nil {
			return true
		}
	}
	return false
}

// requireResourceLimits fails when container sets no cpu or no memory limit,
// which the add-on's policy requires of every container.
func requireResourceLimits(container corev1.Container) error {
	limits := container.Resources.Limits
	if limits.Cpu().IsZero() || limits.Memory().IsZero() {
		return fmt.Errorf("container %q sets no cpu or memory limit, which the add-on's policy requires of every container", container.Name)
	}
	return nil
}
