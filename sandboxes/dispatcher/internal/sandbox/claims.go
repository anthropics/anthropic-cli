package sandbox

import (
	"context"
	"fmt"
	"strconv"
	"time"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/dynamic"
)

// EnsureClaim returns the session's sandbox handle, creating the backing
// object when absent. The upstream path creates a SandboxClaim; the add-on
// path creates a Sandbox directly.
func (c *Client) EnsureClaim(ctx context.Context, sessionID string) (*SandboxClaim, error) {
	if c.shape.Direct() {
		return c.ensureSandbox(ctx, sessionID)
	}
	return c.ensureClaim(ctx, sessionID)
}

func (c *Client) GetClaim(ctx context.Context, name string) (*SandboxClaim, error) {
	obj, err := c.handleObjects().Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return c.handleFromObject(obj), nil
}

// ListClaims returns every session handle of this dispatcher's environment.
func (c *Client) ListClaims(ctx context.Context) ([]SandboxClaim, error) {
	list, err := c.handleObjects().List(ctx, metav1.ListOptions{LabelSelector: c.managedSelector()})
	if err != nil {
		return nil, err
	}
	claims := make([]SandboxClaim, 0, len(list.Items))
	for i := range list.Items {
		claims = append(claims, *c.handleFromObject(&list.Items[i]))
	}
	return claims, nil
}

// DeleteClaim deletes the session's handle and its pod and volumes. On the
// upstream path the controller cascades from the deleted claim; on the add-on
// path the Sandbox is deleted directly and its nonce forgotten. A handle
// already gone is not an error.
func (c *Client) DeleteClaim(ctx context.Context, name string) error {
	if c.shape.Direct() {
		c.nonces.Clear(name)
	}
	err := c.handleObjects().Delete(ctx, name, metav1.DeleteOptions{})
	if apierrors.IsNotFound(err) {
		return nil
	}
	return err
}

func (c *Client) StampLastActive(ctx context.Context, claimName string, at time.Time) error {
	return c.annotateHandle(ctx, claimName, AnnotationLastActive, at.UTC().Format(time.RFC3339))
}

// StampHandoff records on the session's handle which pod the session was
// handed to and the worker's restart count at that moment.
func (c *Client) StampHandoff(ctx context.Context, claimName string, podUID types.UID, restarts int32) error {
	patch := fmt.Sprintf(`{"metadata":{"annotations":{%q:%q,%q:%q}}}`,
		AnnotationHandoffRestarts, strconv.FormatInt(int64(restarts), 10), AnnotationHandoffPodUID, string(podUID))
	_, err := c.handleObjects().Patch(ctx, claimName, types.MergePatchType, []byte(patch), metav1.PatchOptions{FieldManager: fieldManager})
	return err
}

func (c *Client) ensureClaim(ctx context.Context, sessionID string) (*SandboxClaim, error) {
	name := ClaimName(sessionID)
	obj := &unstructured.Unstructured{Object: map[string]any{
		"apiVersion": c.shape.Claims().GroupVersion().String(),
		"kind":       "SandboxClaim",
		"metadata": map[string]any{
			"name":      name,
			"namespace": c.namespace,
			"labels": map[string]any{
				LabelSessionID:     sessionID,
				LabelEnvironmentID: c.environmentID,
				LabelManagedBy:     ManagedByValue,
			},
		},
		"spec": c.shape.claimSpec(c.warmPoolName),
	}}
	created, err := c.claims().Create(ctx, obj, metav1.CreateOptions{})
	if apierrors.IsAlreadyExists(err) {
		return c.adoptExistingClaim(ctx, name, sessionID)
	}
	if err != nil {
		return nil, fmt.Errorf("create sandbox claim %s: %w", name, err)
	}
	return claimFromObject(created), nil
}

// adoptExistingClaim returns the claim stored under the session's name. A claim
// labelled for another session, environment or manager is refused: adopting it
// would hand the session's token to whatever pod that claim resolves to.
func (c *Client) adoptExistingClaim(ctx context.Context, name, sessionID string) (*SandboxClaim, error) {
	obj, err := c.claims().Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	labels := obj.GetLabels()
	if labels[LabelSessionID] != sessionID || labels[LabelEnvironmentID] != c.environmentID || labels[LabelManagedBy] != ManagedByValue {
		return nil, fmt.Errorf("a claim named %s already exists but belongs to something other than session %s in environment %s; remove it before the session can run here", name, sessionID, c.environmentID)
	}
	return claimFromObject(obj), nil
}

// handleObjects is the resource holding each session's handle: SandboxClaims
// on the upstream path, Sandboxes on the add-on path.
func (c *Client) handleObjects() dynamic.ResourceInterface {
	if c.shape.Direct() {
		return c.sandboxes()
	}
	return c.claims()
}

// handleFromObject reads a session handle from its stored object, per path.
func (c *Client) handleFromObject(obj *unstructured.Unstructured) *SandboxClaim {
	if c.shape.Direct() {
		return sandboxAsClaim(obj)
	}
	return claimFromObject(obj)
}

func (c *Client) annotateHandle(ctx context.Context, claimName, key, value string) error {
	patch := fmt.Sprintf(`{"metadata":{"annotations":{%q:%q}}}`, key, value)
	_, err := c.handleObjects().Patch(ctx, claimName, types.MergePatchType, []byte(patch), metav1.PatchOptions{FieldManager: fieldManager})
	return err
}

func claimFromObject(obj *unstructured.Unstructured) *SandboxClaim {
	sandboxName, _, _ := unstructured.NestedString(obj.Object, "status", "sandbox", "name")
	claim := &SandboxClaim{
		Name:        obj.GetName(),
		SessionID:   obj.GetLabels()[LabelSessionID],
		SandboxName: sandboxName,
	}
	claim.LastActive = annotatedTime(obj, AnnotationLastActive)
	claim.HandoffRestarts = annotatedCount(obj, AnnotationHandoffRestarts)
	claim.HandoffPodUID = types.UID(obj.GetAnnotations()[AnnotationHandoffPodUID])
	return claim
}

// sandboxAsClaim reads a session handle from a Sandbox the dispatcher created
// directly: the Sandbox is both the handle and its own bound sandbox, and
// carries the same handoff annotations a claim would.
func sandboxAsClaim(obj *unstructured.Unstructured) *SandboxClaim {
	claim := &SandboxClaim{
		Name:        obj.GetName(),
		SessionID:   obj.GetLabels()[LabelSessionID],
		SandboxName: obj.GetName(),
	}
	claim.LastActive = annotatedTime(obj, AnnotationLastActive)
	claim.HandoffRestarts = annotatedCount(obj, AnnotationHandoffRestarts)
	claim.HandoffPodUID = types.UID(obj.GetAnnotations()[AnnotationHandoffPodUID])
	return claim
}

// annotatedTime is the RFC 3339 time in the annotation, or zero when the
// annotation is absent or unreadable.
func annotatedTime(obj *unstructured.Unstructured, key string) time.Time {
	t, _ := time.Parse(time.RFC3339, obj.GetAnnotations()[key])
	return t
}

// annotatedCount is the number in the annotation, or nil when the
// annotation is absent or unreadable.
func annotatedCount(obj *unstructured.Unstructured, key string) *int32 {
	n, err := strconv.ParseInt(obj.GetAnnotations()[key], 10, 32)
	if err != nil {
		return nil
	}
	count := int32(n)
	return &count
}
