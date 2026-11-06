package dispatcher

import (
	"context"
	"errors"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"

	"github.com/anthropics/anthropic-cli/sandboxes/dispatcher/internal/sandbox"
)

// Reconcile deletes every managed claim whose session the API reports gone
// or terminated, suspends every running sandbox whose worker has exited since
// its last handoff, and starts watching the rest. A claim whose session
// cannot be checked is kept.
func (d *Dispatcher) Reconcile(ctx context.Context) error {
	claims, err := d.sandboxes.ListClaims(ctx)
	if err != nil {
		return fmt.Errorf("listing sandbox claims: %w", err)
	}
	for i := range claims {
		claim := &claims[i]
		if d.deleteIfSessionOver(ctx, claim) {
			continue
		}
		d.adoptRunningPod(ctx, claim)
	}
	d.logger.Info("reconciled sandbox claims", "count", len(claims))
	return nil
}

// adoptRunningPod suspends a claim whose worker exited before this dispatcher
// started and watches any other. Reads are retried and a failed suspend keeps
// the watch: this runs once, and a session it drops is never watched again.
func (d *Dispatcher) adoptRunningPod(ctx context.Context, claim *sandbox.SandboxClaim) {
	var sb *sandbox.Sandbox
	err := d.retryTransient(ctx, func() (err error) {
		sb, err = d.sandboxes.SandboxOf(ctx, claim)
		return err
	})
	if err != nil {
		d.warnUnlessAbsent(err, "could not read the claim's sandbox; leaving it", claim)
		return
	}
	if sb.Suspended {
		return
	}
	var pod *corev1.Pod
	err = d.retryTransient(ctx, func() (err error) {
		pod, err = d.sandboxes.PodOf(ctx, sb)
		return err
	})
	if apierrors.IsNotFound(err) {
		d.watchHandoff(handoff{claim: *claim, sb: sb})
		return
	}
	if err != nil {
		d.warnUnlessAbsent(err, "could not read the sandbox's pod; leaving it", claim)
		return
	}
	if claim.HandoffPodUID != pod.UID {
		d.watchHandoff(handoff{claim: *claim, sb: sb})
		return
	}
	status, ok := sandbox.WorkerStatusOf(pod)
	if !ok {
		d.watchHandoff(handoff{claim: *claim, sb: sb})
		return
	}
	if workerDoneSince(status, claim) {
		if d.suspendSandbox(ctx, claim, sb) {
			return
		}
		// The watch retries the suspend, so its baseline must be a count status
		// already exceeds. With no recorded handoff, an exit seen only as a
		// restart exceeds just the count from before that restart.
		restarts := status.RestartCount
		if claim.HandoffRestarts != nil {
			restarts = *claim.HandoffRestarts
		} else if status.Terminated.IsZero() && restarts > 0 {
			restarts--
		}
		d.watchHandoff(handoff{claim: *claim, sb: sb, podUID: pod.UID, restarts: restarts})
		return
	}
	d.watchHandoff(handoff{claim: *claim, sb: sb, podUID: pod.UID, restarts: status.RestartCount})
}

// warnUnlessAbsent logs a read failure, except for an unbound claim or a
// pod that is gone, which are ordinary states.
func (d *Dispatcher) warnUnlessAbsent(err error, msg string, claim *sandbox.SandboxClaim) {
	if errors.Is(err, sandbox.ErrNotBound) || apierrors.IsNotFound(err) {
		return
	}
	d.logger.Warn(msg, "claim", claim.Name, "error", err)
}

// workerDoneSince reports whether the worker exited after the claim's last
// handoff. The restart count recorded at the handoff is exact; a claim from
// before a handoff recorded one is judged by the last-active clock.
func workerDoneSince(status sandbox.WorkerStatus, claim *sandbox.SandboxClaim) bool {
	if claim.HandoffRestarts != nil {
		return status.ExitedPast(*claim.HandoffRestarts)
	}
	return status.ExitedSince(claim.LastActive)
}
