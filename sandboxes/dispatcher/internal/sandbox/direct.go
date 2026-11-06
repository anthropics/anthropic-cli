package sandbox

import (
	"context"
	"encoding/json"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/types"
)

// ensureSandbox returns the session's Sandbox, creating it with its nonce
// Secret when absent. A new Sandbox that lacks the Secret, or whose workspace
// claim the served schema dropped, is deleted: its pod could never run.
func (c *Client) ensureSandbox(ctx context.Context, sessionID string) (*SandboxClaim, error) {
	name := SandboxName(sessionID)
	obj, err := c.sandboxFromTemplate(ctx, name, sessionID)
	if err != nil {
		return nil, err
	}
	created, err := c.sandboxes().Create(ctx, obj, metav1.CreateOptions{})
	if apierrors.IsAlreadyExists(err) {
		return c.adoptExistingSandbox(ctx, name, sessionID)
	}
	if err != nil {
		return nil, fmt.Errorf("create sandbox %s: %w", name, err)
	}
	if _, wrote, _ := unstructured.NestedSlice(obj.Object, "spec", "volumeClaimTemplates"); wrote {
		if _, stored, _ := unstructured.NestedSlice(created.Object, "spec", "volumeClaimTemplates"); !stored {
			_ = c.sandboxes().Delete(ctx, name, metav1.DeleteOptions{})
			return nil, fmt.Errorf("sandbox %s did not record its workspace claim; check that -api-version %s matches the cluster", name, c.shape.version)
		}
	}
	if err := c.createNonceSecret(ctx, created, NonceSecretName(name)); err != nil {
		_ = c.sandboxes().Delete(ctx, name, metav1.DeleteOptions{})
		c.nonces.Clear(name)
		return nil, err
	}
	return sandboxAsClaim(created), nil
}

// createNonceSecret stores a fresh nonce for the Sandbox in secretName, which
// its fetch container reads. The Sandbox owns the Secret, so deleting the
// Sandbox removes it: the dispatcher holds no Secrets verb but create.
func (c *Client) createNonceSecret(ctx context.Context, sb *unstructured.Unstructured, secretName string) error {
	name := sb.GetName()
	nonce, err := c.nonces.Mint(name)
	if err != nil {
		return fmt.Errorf("mint assignment nonce for sandbox %s: %w", name, err)
	}
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      secretName,
			Namespace: c.namespace,
			Labels:    map[string]string{LabelManagedBy: ManagedByValue},
			OwnerReferences: []metav1.OwnerReference{{
				APIVersion: c.shape.Sandboxes().GroupVersion().String(),
				Kind:       "Sandbox",
				Name:       name,
				UID:        sb.GetUID(),
			}},
		},
		Type:       corev1.SecretTypeOpaque,
		StringData: map[string]string{NonceSecretKey: nonce},
	}
	if _, err := c.kube.CoreV1().Secrets(c.namespace).Create(ctx, secret, metav1.CreateOptions{}); err != nil {
		return fmt.Errorf("create the nonce secret for sandbox %s: %w", name, err)
	}
	return nil
}

// adoptExistingSandbox returns the stored Sandbox while the registry holds its
// nonce, and passes one from before a restart to adoptAfterRestart. Foreign
// labels are refused: an unstamped pod could then pass as the session's.
func (c *Client) adoptExistingSandbox(ctx context.Context, name, sessionID string) (*SandboxClaim, error) {
	obj, err := c.sandboxes().Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	labels := obj.GetLabels()
	if labels[LabelSessionID] != sessionID || labels[LabelEnvironmentID] != c.environmentID || labels[LabelManagedBy] != ManagedByValue {
		c.nonces.Clear(name)
		return nil, fmt.Errorf("a sandbox named %s already exists but belongs to something other than session %s in environment %s; remove it before the session can run here", name, sessionID, c.environmentID)
	}
	if !c.nonces.Known(name) {
		return c.adoptAfterRestart(ctx, obj)
	}
	return sandboxAsClaim(obj), nil
}

// adoptAfterRestart resolves a Sandbox whose nonce a restart lost. It fails
// while the pod exists, since that pod may be running a turn. With no pod it
// deletes a Sandbox that holds nothing of the session and re-arms any other.
func (c *Client) adoptAfterRestart(ctx context.Context, obj *unstructured.Unstructured) (*SandboxClaim, error) {
	name := obj.GetName()
	_, err := c.PodOf(ctx, c.sandboxFromObject(obj))
	if err == nil {
		return nil, fmt.Errorf("sandbox %s predates this dispatcher's nonce registry and its pod may still be running a turn; it is settled once the pod is gone", name)
	}
	if !apierrors.IsNotFound(err) {
		return nil, err
	}
	if sandboxAsClaim(obj).HandoffPodUID == "" && !claimBacked(obj) {
		if err := c.sandboxes().Delete(ctx, name, metav1.DeleteOptions{}); err != nil && !apierrors.IsNotFound(err) {
			return nil, fmt.Errorf("delete sandbox %s, whose nonce this dispatcher does not hold: %w", name, err)
		}
		return nil, fmt.Errorf("sandbox %s predates this dispatcher's nonce registry; it was deleted so the next handoff recreates it", name)
	}
	if !c.shape.isSuspended(obj) {
		if err := c.Suspend(ctx, name); err != nil {
			return nil, fmt.Errorf("suspend podless sandbox %s before re-arming it: %w", name, err)
		}
	}
	return c.rearmNonce(ctx, obj)
}

// claimBacked reports whether obj owns a workspace volume claim: a non-empty
// spec.volumeClaimTemplates.
func claimBacked(obj *unstructured.Unstructured) bool {
	templates, _, _ := unstructured.NestedSlice(obj.Object, "spec", "volumeClaimTemplates")
	return len(templates) > 0
}

// rearmNonce gives a suspended Sandbox a fresh nonce in a new Secret and points
// its fetch container there, so the resumed pod starts with it. A failure
// clears the nonce, or handoffs would wait on a pod that cannot authenticate.
func (c *Client) rearmNonce(ctx context.Context, obj *unstructured.Unstructured) (*SandboxClaim, error) {
	name := obj.GetName()
	secretName, err := RearmedNonceSecretName(name)
	if err != nil {
		return nil, err
	}
	if err := c.createNonceSecret(ctx, obj, secretName); err != nil {
		c.nonces.Clear(name)
		return nil, err
	}
	patched, err := c.pointFetchAtSecret(ctx, obj, secretName)
	if err != nil {
		c.nonces.Clear(name)
		return nil, err
	}
	return sandboxAsClaim(patched), nil
}

// pointFetchAtSecret patches the fetch container to read EnvAssignmentNonce
// from secretName. It fails unless the stored Sandbox reads back that way: a
// schema without the field accepts the patch silently.
func (c *Client) pointFetchAtSecret(ctx context.Context, obj *unstructured.Unstructured, secretName string) (*unstructured.Unstructured, error) {
	name := obj.GetName()
	podTemplate, found, err := unstructured.NestedMap(obj.Object, "spec", "podTemplate")
	if err != nil || !found {
		return nil, fmt.Errorf("sandbox %s has no podTemplate to point at its new nonce secret", name)
	}
	if err := stampNonce(podTemplate, secretName); err != nil {
		return nil, fmt.Errorf("sandbox %s: %w", name, err)
	}
	patch, err := json.Marshal(map[string]any{"spec": map[string]any{"podTemplate": podTemplate}})
	if err != nil {
		return nil, err
	}
	patched, err := c.sandboxes().Patch(ctx, name, types.MergePatchType, patch, metav1.PatchOptions{FieldManager: fieldManager})
	if err != nil {
		return nil, fmt.Errorf("point sandbox %s at its new nonce secret: %w", name, err)
	}
	if fetchNonceSecret(patched) != secretName {
		return nil, fmt.Errorf("sandbox %s did not record its new nonce secret; check that -api-version %s matches the cluster", name, c.shape.version)
	}
	return patched, nil
}

// fetchNonceSecret reads the Secret name the Sandbox's fetch container takes
// EnvAssignmentNonce from, or "" when it takes none.
func fetchNonceSecret(obj *unstructured.Unstructured) string {
	containers, _, _ := unstructured.NestedSlice(obj.Object, "spec", "podTemplate", "spec", "containers")
	for _, raw := range containers {
		container, ok := raw.(map[string]any)
		if !ok || container["name"] != FetchContainer {
			continue
		}
		env, _, _ := unstructured.NestedSlice(container, "env")
		for _, e := range env {
			entry, ok := e.(map[string]any)
			if !ok || entry["name"] != EnvAssignmentNonce {
				continue
			}
			secretName, _, _ := unstructured.NestedString(entry, "valueFrom", "secretKeyRef", "name")
			return secretName
		}
	}
	return ""
}
