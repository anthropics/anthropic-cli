package dispatcher

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"

	"github.com/anthropics/anthropic-cli/sandboxes/dispatcher/internal/cmaapi"
	"github.com/anthropics/anthropic-cli/sandboxes/dispatcher/internal/sandbox"
)

// PodChanged acts on a pod of a watched session. Once the worker of the pod
// the session was handed to has exited since the handoff, the pod's
// work item is closed and the sandbox is suspended or deleted. A pod that
// replaced it is resolved from the Sandbox side, and since the session's
// worker went with the old pod, the same rule applies. A session busy with
// a handoff is left for the next change or sweep.
func (d *Dispatcher) PodChanged(ctx context.Context, pod *corev1.Pod) {
	h, handedTo := d.handoffOfPod(pod.UID)
	if !handedTo {
		h, replaced := d.handoffAwaitingPod(pod.Name)
		if replaced {
			d.endSession(ctx, h, "the pod the session was handed to is gone", func() bool { return d.isSandboxPod(ctx, h.sb, pod.UID) })
		}
		return
	}
	if status, ok := sandbox.WorkerStatusOf(pod); ok && status.ExitedPast(h.restarts) {
		d.endSession(ctx, h, "worker exited", func() bool { return d.workerConfirmedDone(ctx, h) })
	}
}

// PodGone closes the work item of a pod that is gone and marks it gone in
// its session's watch; the session stays watched until the Sandbox's next
// pod is resolved.
func (d *Dispatcher) PodGone(podUID types.UID) {
	d.workItems.ForgetPod(podUID)
	d.mu.Lock()
	defer d.mu.Unlock()
	for sessionID, h := range d.handoffs {
		if h.podUID == podUID {
			h.podUID = ""
			d.handoffs[sessionID] = h
		}
	}
}

// handoffAwaitingPod finds the watched session whose pod is gone and whose
// Sandbox names podName.
func (d *Dispatcher) handoffAwaitingPod(podName string) (handoff, bool) {
	d.mu.Lock()
	defer d.mu.Unlock()
	for _, h := range d.handoffs {
		if h.podUID == "" && h.sb.PodName == podName {
			return h, true
		}
	}
	return handoff{}, false
}

// endSession applies the exit rule to a watched session once confirm holds
// under the session's claim, and stops watching it when the rule took.
func (d *Dispatcher) endSession(ctx context.Context, h handoff, why string, confirm func() bool) {
	if !d.claimInflight(h.claim.SessionID) {
		return
	}
	defer d.releaseInflight(h.claim.SessionID)
	if !confirm() {
		return
	}
	d.logger.Info("session's worker is gone", "reason", why, "session_id", h.claim.SessionID, "sandbox", h.sb.Name)
	if d.deleteIfSessionOver(ctx, &h.claim) || d.suspendSandbox(ctx, &h.claim, h.sb) {
		d.forgetHandoff(h.claim.SessionID)
	}
}

// isSandboxPod reports whether the Sandbox's own pod, read from the Sandbox
// side, is podUID.
func (d *Dispatcher) isSandboxPod(ctx context.Context, sb *sandbox.Sandbox, podUID types.UID) bool {
	pod, err := d.sandboxes.PodOf(ctx, sb)
	return err == nil && pod.UID == podUID
}

// workerConfirmedDone re-reads the pod and reports whether its worker exited
// since the handoff, closing the pod's work item when it has. A watch event
// can be a stale snapshot, and acting on one would suspend a pod mid-turn.
func (d *Dispatcher) workerConfirmedDone(ctx context.Context, h handoff) bool {
	pod, err := d.sandboxes.PodOf(ctx, h.sb)
	if apierrors.IsNotFound(err) || (err == nil && pod.UID != h.podUID) {
		d.workItems.Close(h.podUID)
		return true
	}
	if err != nil {
		return false
	}
	status, ok := sandbox.WorkerStatusOf(pod)
	if !ok || !status.ExitedPast(h.restarts) {
		return false
	}
	if h.restartsPreCollection && status.Terminated.IsZero() {
		// A count above a pre-collection baseline may be the mid-wait
		// restart the handoff tolerated; the run now going may hold the
		// token, so only an ended run confirms.
		return false
	}
	d.workItems.Close(h.podUID)
	return true
}

// deleteIfSessionOver reports whether it deleted the claim. Only an
// authoritative answer from the API deletes anything.
func (d *Dispatcher) deleteIfSessionOver(ctx context.Context, claim *sandbox.SandboxClaim) bool {
	if claim.SessionID == "" {
		d.logger.Warn("managed claim carries no session label; leaving it", "claim", claim.Name)
		return false
	}
	state, err := d.sessions.SessionState(ctx, claim.SessionID)
	if err != nil {
		d.logger.Warn("could not check the claim's session; leaving it", "claim", claim.Name, "session_id", claim.SessionID, "error", err)
		return false
	}
	switch state {
	case cmaapi.SessionGone:
		d.deleteClaim(ctx, claim, "session gone")
	case cmaapi.SessionEnded:
		d.deleteClaim(ctx, claim, "session terminated")
	default:
		return false
	}
	return true
}

// suspendSandbox stamps the claim's last activity and suspends the sandbox, which
// deletes the pod and keeps the workspace; it reports whether both took.
func (d *Dispatcher) suspendSandbox(ctx context.Context, claim *sandbox.SandboxClaim, sb *sandbox.Sandbox) bool {
	if err := d.sandboxes.StampLastActive(ctx, claim.Name, d.Now()); err != nil {
		d.logger.Error("stamping last-active before suspending", "claim", claim.Name, "error", err)
		return false
	}
	if err := d.sandboxes.Suspend(ctx, sb.Name); err != nil {
		d.logger.Error("suspending the sandbox", "sandbox", sb.Name, "claim", claim.Name, "error", err)
		return false
	}
	d.logger.Info("suspended sandbox", "sandbox", sb.Name, "session_id", claim.SessionID)
	return true
}

func (d *Dispatcher) deleteClaim(ctx context.Context, claim *sandbox.SandboxClaim, why string) {
	if err := d.sandboxes.DeleteClaim(ctx, claim.Name); err != nil {
		d.logger.Error("deleting sandbox claim", "claim", claim.Name, "session_id", claim.SessionID, "error", err)
		return
	}
	d.forgetHandoff(claim.SessionID)
	d.workItems.ForgetSession(claim.SessionID)
	d.logger.Info("deleted sandbox claim", "claim", claim.Name, "session_id", claim.SessionID, "reason", why)
}
