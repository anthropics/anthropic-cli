package preflight

import (
	"context"
	"errors"
	"fmt"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func (c *Checker) apiServed(ctx context.Context) error {
	want := map[string][]string{
		c.Options.Shape.Sandboxes().GroupVersion().String(): {"sandboxes"},
		c.Options.Shape.Claims().GroupVersion().String():    {"sandboxclaims", "sandboxtemplates", "sandboxwarmpools"},
	}
	if c.Options.Shape.Direct() {
		want = map[string][]string{
			c.Options.Shape.Sandboxes().GroupVersion().String(): {"sandboxes"},
			c.Options.Shape.Templates().GroupVersion().String(): {"sandboxtemplates"},
		}
	}
	for groupVersion, resources := range want {
		list, err := c.Kube.Discovery().ServerResourcesForGroupVersion(groupVersion)
		if err != nil {
			return fmt.Errorf("%s is not served: %w", groupVersion, err)
		}
		served := map[string]bool{}
		for _, r := range list.APIResources {
			served[r.Name] = true
		}
		for _, r := range resources {
			if !served[r] {
				return fmt.Errorf("%s does not serve %s", groupVersion, r)
			}
		}
	}
	return nil
}

// warmPoolTemplate fails when session pods, which come from the warm pool's
// template, would not be the template the following checks inspect.
func (c *Checker) warmPoolTemplate(ctx context.Context) error {
	pool, err := c.Dyn.Resource(c.Options.Shape.WarmPools()).Namespace(c.Options.Namespace).Get(ctx, c.Options.WarmPoolName, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("SandboxWarmPool %s: %w", c.Options.WarmPoolName, err)
	}
	name, _, _ := unstructured.NestedString(pool.Object, "spec", "sandboxTemplateRef", "name")
	if name != c.Options.SandboxTemplateName {
		return fmt.Errorf("SandboxWarmPool %s makes its pods from SandboxTemplate %q, not %q; point -warm-pool and -sandbox-template at the same template", c.Options.WarmPoolName, name, c.Options.SandboxTemplateName)
	}
	return nil
}

func (c *Checker) runtimeClass(ctx context.Context) error {
	tpl, err := c.template(ctx)
	if err != nil {
		return err
	}
	name, _, _ := unstructured.NestedString(tpl.Object, "spec", "podTemplate", "spec", "runtimeClassName")
	if name == "" {
		return errors.New("the SandboxTemplate sets no runtimeClassName; sandboxes would run on the host kernel")
	}
	rc, err := c.Kube.NodeV1().RuntimeClasses().Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("RuntimeClass %s: %w", name, err)
	}
	if rc.Handler != "runsc" && rc.Handler != "gvisor" && !strings.HasPrefix(rc.Handler, "kata") {
		return fmt.Errorf("RuntimeClass %s uses handler %q; only runsc or gvisor (gVisor) or a Kata handler isolates sandboxes", name, rc.Handler)
	}
	return nil
}

// managedPolicyOff reads the template as the cluster stores it: a schema
// that predates networkPolicyManagement prunes the field on apply, and the
// controller's own policy then leaves 100.64.0.0/10 open.
func (c *Checker) managedPolicyOff(ctx context.Context) error {
	tpl, err := c.template(ctx)
	if err != nil {
		return err
	}
	mode, _, _ := unstructured.NestedString(tpl.Object, "spec", "networkPolicyManagement")
	if mode != "Unmanaged" {
		return fmt.Errorf("the stored SandboxTemplate has networkPolicyManagement %q, not Unmanaged: the template omits it, or the served schema predates the field and pruned it; either way the controller's own policy leaves 100.64.0.0/10 open", mode)
	}
	return nil
}

func (c *Checker) networkProbe(context.Context) error {
	byHand := fmt.Sprintf("verify by hand that a sandbox pod cannot reach another sandbox pod, 169.254.169.254 or the cluster's API server and can reach the dispatcher at %s", c.Options.ServiceHost)
	if c.Options.SkipNetworkProbe {
		return fmt.Errorf("%w: skipped by -skip-network-probe; %s", ErrWarning, byHand)
	}
	return errors.New("not implemented in this version; " + byHand)
}

func (c *Checker) podSecurity(ctx context.Context) error {
	ns, err := c.Kube.CoreV1().Namespaces().Get(ctx, c.Options.Namespace, metav1.GetOptions{})
	if err != nil {
		return err
	}
	if ns.Labels["pod-security.kubernetes.io/enforce"] != "restricted" {
		return fmt.Errorf("namespace %s lacks the label pod-security.kubernetes.io/enforce=restricted", ns.Name)
	}
	return nil
}

func (c *Checker) deployment(ctx context.Context) error {
	dep, err := c.Kube.AppsV1().Deployments(c.Options.Namespace).Get(ctx, c.Options.DeploymentName, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("deployment %s: %w", c.Options.DeploymentName, err)
	}
	if dep.Spec.Replicas != nil && *dep.Spec.Replicas != 1 {
		return fmt.Errorf("deployment %s has %d replicas; the dispatcher must run exactly one", dep.Name, *dep.Spec.Replicas)
	}
	if dep.Spec.Strategy.Type != "Recreate" {
		return fmt.Errorf("deployment %s uses strategy %q; two dispatchers would overlap during a rollout without Recreate", dep.Name, dep.Spec.Strategy.Type)
	}
	sa := dep.Spec.Template.Spec.ServiceAccountName
	if sa == "" {
		sa = "default"
	}
	if sa != c.Options.DispatcherServiceAccount {
		return fmt.Errorf("deployment %s runs as service account %q, but the RBAC check inspects %q; the checked identity must be the deployed one", dep.Name, sa, c.Options.DispatcherServiceAccount)
	}
	return nil
}

func (c *Checker) environmentKey(ctx context.Context) error {
	return c.API.VerifyEnvironment(ctx)
}

func (c *Checker) template(ctx context.Context) (*unstructured.Unstructured, error) {
	tpl, err := c.Dyn.Resource(c.Options.Shape.Templates()).Namespace(c.Options.Namespace).Get(ctx, c.Options.SandboxTemplateName, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("SandboxTemplate %s: %w", c.Options.SandboxTemplateName, err)
	}
	return tpl, nil
}
