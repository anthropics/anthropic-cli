package sandbox

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

var errPodStillRunning = errors.New("the suspended sandbox's pod is still running")

// Suspend asks the controller to delete the Sandbox's pod and keep its volumes.
func (c *Client) Suspend(ctx context.Context, sandboxName string) error {
	return c.setSuspended(ctx, sandboxName, true)
}

// Resume asks the controller to recreate the Sandbox's pod once the pod from
// before the suspend is gone or terminating, polling every interval. A
// controller that saw the suspend and the resume together would keep that
// pod, whose worker has already exited.
func (c *Client) Resume(ctx context.Context, sandboxName string, interval time.Duration) error {
	sb, err := c.GetSandbox(ctx, sandboxName)
	if err != nil {
		return err
	}
	if err := poll(ctx, interval, func() error { return c.podLeaving(ctx, sb) }); err != nil {
		return fmt.Errorf("waiting for the pod of suspended sandbox %s to stop: %w", sandboxName, err)
	}
	return c.setSuspended(ctx, sandboxName, false)
}

func (c *Client) podLeaving(ctx context.Context, sb *Sandbox) error {
	pod, err := c.PodOf(ctx, sb)
	switch {
	case apierrors.IsNotFound(err):
		return nil
	case err != nil:
		return err
	case pod.DeletionTimestamp != nil:
		return nil
	}
	return errPodStillRunning
}

// setSuspended fails when the patched object does not read back in the
// requested state, which is what a schema the configured API version does
// not have looks like: the server accepts unknown fields silently.
func (c *Client) setSuspended(ctx context.Context, sandboxName string, suspended bool) error {
	patch, err := json.Marshal(c.shape.suspendPatch(suspended))
	if err != nil {
		return err
	}
	patched, err := c.sandboxes().Patch(ctx, sandboxName, types.MergePatchType, patch, metav1.PatchOptions{FieldManager: fieldManager})
	if err != nil {
		return err
	}
	if recorded, found := c.shape.suspendState(patched); !found || recorded != suspended {
		return fmt.Errorf("sandbox %s did not record the suspend change; check that -api-version %s matches the cluster", sandboxName, c.shape.version)
	}
	return nil
}
