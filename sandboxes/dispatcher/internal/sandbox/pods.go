package sandbox

import (
	"context"
	"errors"
	"fmt"
	"time"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/types"
)

// podNameAnnotation names the backing pod when it differs from the Sandbox
// name (pods adopted from older warm pools).
const podNameAnnotation = "agents.x-k8s.io/pod-name"

var errPodNotRunning = errors.New("pod is not running yet")

// SandboxOf returns the Sandbox the claim's status names, or ErrNotBound.
func (c *Client) SandboxOf(ctx context.Context, claim *SandboxClaim) (*Sandbox, error) {
	if claim.SandboxName == "" {
		return nil, ErrNotBound
	}
	return c.GetSandbox(ctx, claim.SandboxName)
}

func (c *Client) GetSandbox(ctx context.Context, name string) (*Sandbox, error) {
	obj, err := c.sandboxes().Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return c.sandboxFromObject(obj), nil
}

// PodOf returns the Sandbox's backing pod, named by the Sandbox's pod-name
// annotation or, without one, after the Sandbox itself. A pod under that
// name that sb does not control is answered as not found: anyone able to
// create pods could wear the name while the real pod is absent, and the
// session's token must never be keyed to such a pod.
func (c *Client) PodOf(ctx context.Context, sb *Sandbox) (*corev1.Pod, error) {
	pod, err := c.kube.CoreV1().Pods(c.namespace).Get(ctx, sb.PodName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	if owner := metav1.GetControllerOf(pod); owner == nil || owner.UID != sb.UID {
		notFound := apierrors.NewNotFound(corev1.Resource("pods"), sb.PodName)
		notFound.ErrStatus.Message = fmt.Sprintf("pod %s exists but is not controlled by sandbox %s; treating the sandbox's pod as absent", sb.PodName, sb.Name)
		return nil, notFound
	}
	return pod, nil
}

// SandboxPod reports the sandbox's backing pod's UID, address, and whether it
// is Running, from the sandbox and pod as read now. It is the dispatcher's
// answer to "where is this sandbox's pod", used to authenticate a nonce-
// bearing request against the pod's network position. A pod that is absent,
// terminating, or not yet Running is running=false; any read error is
// returned so the caller fails closed.
func (c *Client) SandboxPod(ctx context.Context, sandboxName string) (podUID types.UID, podIP string, running bool, err error) {
	sb, err := c.GetSandbox(ctx, sandboxName)
	if err != nil {
		return "", "", false, err
	}
	pod, err := c.PodOf(ctx, sb)
	if err != nil {
		return "", "", false, err
	}
	up := pod.Status.Phase == corev1.PodRunning && pod.DeletionTimestamp == nil
	return pod.UID, pod.Status.PodIP, up, nil
}

// WaitForRunningPod polls until the claim is bound and its pod is Running,
// or ctx ends, and returns the Sandbox and that pod. A terminating pod is
// waited out like a missing one: after a resume the old pod may still be
// terminating under the name the new one will take.
func (c *Client) WaitForRunningPod(ctx context.Context, claimName string, interval time.Duration) (*Sandbox, *corev1.Pod, error) {
	var sb *Sandbox
	var pod *corev1.Pod
	err := poll(ctx, interval, func() (err error) {
		sb, pod, err = c.runningPod(ctx, claimName)
		return err
	})
	if err != nil {
		return nil, nil, fmt.Errorf("waiting for the pod of claim %s: %w", claimName, err)
	}
	return sb, pod, nil
}

func (c *Client) sandboxFromObject(obj *unstructured.Unstructured) *Sandbox {
	podName := obj.GetAnnotations()[podNameAnnotation]
	if podName == "" {
		podName = obj.GetName()
	}
	return &Sandbox{
		Name:      obj.GetName(),
		UID:       obj.GetUID(),
		PodName:   podName,
		Suspended: c.shape.isSuspended(obj),
		Ready:     c.shape.isReady(obj),
	}
}

func (c *Client) runningPod(ctx context.Context, claimName string) (*Sandbox, *corev1.Pod, error) {
	claim, err := c.GetClaim(ctx, claimName)
	if err != nil {
		return nil, nil, err
	}
	sb, err := c.SandboxOf(ctx, claim)
	if err != nil {
		return nil, nil, err
	}
	pod, err := c.PodOf(ctx, sb)
	if err != nil {
		return nil, nil, err
	}
	if pod.Status.Phase != corev1.PodRunning || pod.DeletionTimestamp != nil || !sb.Ready {
		return nil, nil, errPodNotRunning
	}
	return sb, pod, nil
}

// poll calls try every interval until it succeeds, fails with an error that
// retryable does not accept, or ctx ends.
func poll(ctx context.Context, interval time.Duration, try func() error) error {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		err := try()
		if err == nil || !retryable(err) {
			return err
		}
		select {
		case <-ctx.Done():
			return fmt.Errorf("%w (last: %v)", ctx.Err(), err)
		case <-ticker.C:
		}
	}
}

func retryable(err error) bool {
	return errors.Is(err, ErrNotBound) || errors.Is(err, errPodNotRunning) || errors.Is(err, errPodStillRunning) || apierrors.IsNotFound(err)
}
