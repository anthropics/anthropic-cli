package sandbox

import (
	"context"
	"errors"
	"fmt"
	"slices"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

// workspaceSizeLimit caps the emptyDir standing in for a claim-backed
// workspace, so a runaway session fills its own allotment rather than the
// node's disk. It matches the recipe's default claim size.
const workspaceSizeLimit = "10Gi"

// sandboxFromTemplate reads the configured SandboxTemplate and builds the
// session's Sandbox from it.
func (c *Client) sandboxFromTemplate(ctx context.Context, name, sessionID string) (*unstructured.Unstructured, error) {
	tpl, err := c.dyn.Resource(c.shape.Templates()).Namespace(c.namespace).Get(ctx, c.templateName, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("read sandbox template %s: %w", c.templateName, err)
	}
	labels := map[string]any{
		LabelSessionID:     sessionID,
		LabelEnvironmentID: c.environmentID,
		LabelManagedBy:     ManagedByValue,
	}
	return SandboxFromTemplate(c.shape, tpl, c.namespace, name, labels, c.workspaceClaim)
}

// SandboxFromTemplate builds the Sandbox stored under name on the direct
// path: tpl's pod template with the fetch container's nonce variable pointed
// at the sandbox's nonce Secret, one replica, and Retain so a delete is the
// only thing that drops the workspace. The add-on's SandboxTemplate schema
// stores no volumeClaimTemplates, so none are copied from tpl; the Sandbox
// schema does store them, so claim, when set, becomes the Sandbox's own
// volumeClaimTemplates entry and the controller binds the worker's workspace
// mount to that claim, which outlives a suspend. Without a claim the mount
// is backed by an emptyDir a suspend discards. The add-on's Sandbox schema
// stores no networkPolicyManagement, so tpl's Unmanaged setting cannot be
// carried; sandbox isolation on this path is the recipe's own
// NetworkPolicies, which select the pod labels the copied pod template
// carries.
func SandboxFromTemplate(shape Shape, tpl *unstructured.Unstructured, namespace, name string, labels map[string]any, claim *WorkspaceClaim) (*unstructured.Unstructured, error) {
	podTemplate, found, err := unstructured.NestedMap(tpl.Object, "spec", "podTemplate")
	if err != nil || !found {
		return nil, fmt.Errorf("sandbox template %s has no podTemplate", tpl.GetName())
	}
	if err := stampNonce(podTemplate, NonceSecretName(name)); err != nil {
		return nil, fmt.Errorf("sandbox template %s: %w", tpl.GetName(), err)
	}
	spec := map[string]any{
		"podTemplate":    podTemplate,
		"replicas":       int64(1),
		"shutdownPolicy": "Retain",
	}
	if err := ensureWorkspaceVolume(podTemplate, spec, claim); err != nil {
		return nil, fmt.Errorf("sandbox template %s: %w", tpl.GetName(), err)
	}
	metadata := map[string]any{"name": name, "namespace": namespace}
	if len(labels) > 0 {
		metadata["labels"] = labels
	}
	return &unstructured.Unstructured{Object: map[string]any{
		"apiVersion": shape.Sandboxes().GroupVersion().String(),
		"kind":       "Sandbox",
		"metadata":   metadata,
		"spec":       spec,
	}}, nil
}

// stampNonce points the fetch container's EnvAssignmentNonce at secretName,
// replacing any entry of that name. The kubelet resolves duplicate env names
// last-wins, so a kept entry could differ from the nonce the registry accepts.
func stampNonce(podTemplate map[string]any, secretName string) error {
	containers, found, err := unstructured.NestedSlice(podTemplate, "spec", "containers")
	if err != nil || !found {
		return errors.New("the podTemplate has no containers")
	}
	for i, raw := range containers {
		container, ok := raw.(map[string]any)
		if !ok || container["name"] != FetchContainer {
			continue
		}
		env, _, _ := unstructured.NestedSlice(container, "env")
		env = slices.DeleteFunc(env, func(e any) bool {
			entry, ok := e.(map[string]any)
			return ok && entry["name"] == EnvAssignmentNonce
		})
		env = append(env, map[string]any{
			"name":      EnvAssignmentNonce,
			"valueFrom": map[string]any{"secretKeyRef": map[string]any{"name": secretName, "key": NonceSecretKey}},
		})
		if err := unstructured.SetNestedSlice(container, env, "env"); err != nil {
			return err
		}
		containers[i] = container
		return unstructured.SetNestedSlice(podTemplate, containers, "spec", "containers")
	}
	return fmt.Errorf("the podTemplate has no %q container to stamp the nonce into", FetchContainer)
}

// ensureWorkspaceVolume backs the worker's WorkspaceVolume mount, since a mount
// that names no volume is refused. A volume podTemplate defines is kept, else
// claim becomes the Sandbox's volumeClaimTemplates entry, else an emptyDir.
func ensureWorkspaceVolume(podTemplate, spec map[string]any, claim *WorkspaceClaim) error {
	volumes, _, _ := unstructured.NestedSlice(podTemplate, "spec", "volumes")
	for _, raw := range volumes {
		if v, ok := raw.(map[string]any); ok && v["name"] == WorkspaceVolume {
			return nil
		}
	}
	if claim != nil {
		spec["volumeClaimTemplates"] = []any{claim.template()}
		return nil
	}
	volumes = append(volumes, map[string]any{
		"name":     WorkspaceVolume,
		"emptyDir": map[string]any{"sizeLimit": workspaceSizeLimit},
	})
	return unstructured.SetNestedSlice(podTemplate, volumes, "spec", "volumes")
}

// template renders the claim as a Sandbox volumeClaimTemplates entry.
func (w *WorkspaceClaim) template() map[string]any {
	spec := map[string]any{
		"accessModes": []any{"ReadWriteOnce"},
		"resources":   map[string]any{"requests": map[string]any{"storage": w.Size.String()}},
	}
	if w.StorageClass != "" {
		spec["storageClassName"] = w.StorageClass
	}
	return map[string]any{
		"metadata": map[string]any{"name": WorkspaceVolume},
		"spec":     spec,
	}
}
