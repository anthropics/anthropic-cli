package dispatcher

import (
	"context"
	"errors"
	"fmt"
	"time"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"

	"github.com/anthropics/anthropic-cli/sandboxes/dispatcher/internal/cmaapi"
	"github.com/anthropics/anthropic-cli/sandboxes/dispatcher/internal/sandbox"
	"github.com/anthropics/anthropic-cli/sandboxes/dispatcher/internal/workitem"
)

// ackTimeout bounds the ack that follows a handoff. The ack runs even when
// the loop is shutting down, because the pod already holds the token.
const ackTimeout = 10 * time.Second

// stampTimeout bounds the post-collection handoff stamp, which runs even
// when the loop is shutting down, like the ack. It is wider than ackTimeout
// because the stamp retries transient failures with pauses in between.
const stampTimeout = 30 * time.Second

// ErrHandoffTimeout means the pod never collected its work item within
// Options.HandoffWait; the item is left unacked for the queue to redeliver.
var ErrHandoffTimeout = errors.New("the sandbox pod did not collect its work item in time")

// Handle gives item's session a running pod, acks the item once the pod has
// collected its work item, and watches the pod for its worker's exit.
func (d *Dispatcher) Handle(ctx context.Context, item *cmaapi.WorkItem) error {
	waitCtx, cancel := context.WithTimeout(ctx, d.Options.HandoffWait)
	defer cancel()

	claim, err := d.sandboxes.EnsureClaim(waitCtx, item.SessionID)
	if err != nil {
		return err
	}
	// A work item's id is its session's id, so the same id can come again.
	// Until the watched pod's delivery is acked, it repeats that delivery.
	// After the ack it is a new message to the session, handed off below.
	if handedWorkID, ok := d.workItems.HandedWork(item.SessionID); ok && handedWorkID == item.ID {
		if h := d.handoffOfSession(item.SessionID); h != nil && h.podUID != "" && !h.acked {
			if !d.restampHandoff(ctx, claim, h) {
				return fmt.Errorf("the handoff of session %s is still not recorded on its claim; leaving the redelivered item unacked", item.SessionID)
			}
			d.logger.Info("work item redelivered after its handoff; acking", "work_id", item.ID, "claim", claim.Name)
			return d.ack(ctx, item, claim)
		}
		d.logger.Info("work item redelivered to a session whose worker is gone; handing it off again", "work_id", item.ID, "session_id", item.SessionID)
	}
	if err := d.resumeIfSuspended(waitCtx, claim); err != nil {
		return err
	}
	sb, pod, err := d.sandboxes.WaitForRunningPod(waitCtx, claim.Name, d.Options.PodPollInterval)
	if err != nil {
		return err
	}
	earlier := d.handoffOfSession(item.SessionID)
	if handedEarlierWork(earlier, pod) {
		return fmt.Errorf("session %s's pod still holds its earlier work item; leaving the item for redelivery", item.SessionID)
	}
	if err := d.recordHandoff(waitCtx, claim, sb, pod.UID); err != nil {
		return err
	}
	taken, opened := d.workItems.Put(pod.UID, workitem.SessionWorkItem{
		SessionID:     item.SessionID,
		WorkID:        item.ID,
		Secret:        item.Secret,
		BaseURL:       d.Options.BaseURL,
		EnvironmentID: d.Options.EnvironmentID,
	})
	if !opened {
		d.restoreHandoff(item.SessionID, earlier)
		return fmt.Errorf("session %s's pod is still running its earlier work item; leaving the item for redelivery", item.SessionID)
	}
	d.logger.Info("pod running; work item pending", "claim", claim.Name, "pod", pod.Name, "session_id", item.SessionID, "work_id", item.ID)
	select {
	case <-taken:
	case <-waitCtx.Done():
		// Drop reports nothing pending when the pod took the work item in
		// the same instant the wait ended; the pod then holds the token, so
		// the item is acked like any handoff.
		if d.workItems.Drop(pod.UID) {
			d.restoreHandoff(item.SessionID, earlier)
			return ErrHandoffTimeout
		}
	}
	// The ack waits for the stamp, since an acked item never comes again. A
	// dispatcher restarted with the claim unstamped would judge the pod that
	// holds the token stale and suspend it mid-turn.
	if !d.recordHandoffRetrying(ctx, claim, sb, pod.UID) {
		return fmt.Errorf("the handoff of session %s could not be recorded on its claim; leaving the item for redelivery", item.SessionID)
	}
	return d.ack(ctx, item, claim)
}

// restampHandoff stamps the watched handoff on the claim before a redelivered
// item is acked. A baseline from before the collection is re-read first, since
// the run it counted may not be the one that holds the token.
func (d *Dispatcher) restampHandoff(ctx context.Context, claim *sandbox.SandboxClaim, h *handoff) bool {
	if h.restartsPreCollection {
		return d.recordHandoffRetrying(ctx, claim, h.sb, h.podUID)
	}
	return d.stampRecordedHandoff(ctx, claim, h)
}

// ack confirms the item after its pod holds the token. It outlives a
// cancelled ctx: leaving the item unacked would only redeliver it.
func (d *Dispatcher) ack(ctx context.Context, item *cmaapi.WorkItem, claim *sandbox.SandboxClaim) error {
	ackCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), ackTimeout)
	defer cancel()
	if err := d.queue.Ack(ackCtx, item.ID); err != nil {
		return fmt.Errorf("acking work item after handoff: %w", err)
	}
	d.mu.Lock()
	if h, ok := d.handoffs[item.SessionID]; ok {
		h.acked = true
		d.handoffs[item.SessionID] = h
	}
	d.mu.Unlock()
	if err := d.sandboxes.StampLastActive(ackCtx, claim.Name, d.Now()); err != nil {
		d.logger.Warn("stamping last-active on the claim", "claim", claim.Name, "error", err)
	}
	d.logger.Info("work item acked", "work_id", item.ID, "session_id", item.SessionID, "claim", claim.Name)
	return nil
}

func (d *Dispatcher) resumeIfSuspended(ctx context.Context, claim *sandbox.SandboxClaim) error {
	sb, err := d.sandboxes.SandboxOf(ctx, claim)
	if errors.Is(err, sandbox.ErrNotBound) {
		return nil
	}
	if err != nil {
		return err
	}
	if !sb.Suspended {
		return nil
	}
	d.logger.Info("resuming suspended sandbox", "sandbox", sb.Name, "session_id", claim.SessionID)
	if err := d.sandboxes.StampLastActive(ctx, claim.Name, d.Now()); err != nil {
		return err
	}
	return d.sandboxes.Resume(ctx, sb.Name, d.Options.PodPollInterval)
}

// handedEarlierWork reports whether pod took the session's watched handoff.
// Such a pod never takes new work. workItems refuses it too but starts empty
// after a restart, while Reconcile rebuilds the watch this reads.
func handedEarlierWork(earlier *handoff, pod *corev1.Pod) bool {
	return earlier != nil && earlier.podUID == pod.UID
}

// recordHandoff watches the session's pod with the worker's restart count read
// fresh. Handle records before the pod can collect the work item and again
// after, so a restart while the worker waited is not counted as an exit.
func (d *Dispatcher) recordHandoff(ctx context.Context, claim *sandbox.SandboxClaim, sb *sandbox.Sandbox, podUID types.UID) error {
	pod, err := d.sandboxes.PodOf(ctx, sb)
	if err != nil {
		return err
	}
	if pod.UID != podUID {
		return fmt.Errorf("sandbox %s replaced its pod during the handoff", sb.Name)
	}
	status, _ := sandbox.WorkerStatusOf(pod)
	d.watchHandoff(handoff{claim: *claim, sb: sb, podUID: podUID, restarts: status.RestartCount})
	return nil
}

// recordHandoffRetrying re-reads the worker's restart count after the
// collection, so a restart during the token wait is not later judged an exit,
// and reports whether the claim was stamped. It outlives a cancelled ctx.
func (d *Dispatcher) recordHandoffRetrying(ctx context.Context, claim *sandbox.SandboxClaim, sb *sandbox.Sandbox, podUID types.UID) bool {
	ctx, cancel := context.WithTimeout(context.WithoutCancel(ctx), stampTimeout)
	defer cancel()
	if err := d.retryTransient(ctx, func() error { return d.recordHandoff(ctx, claim, sb, podUID) }); err != nil {
		d.markRestartsPreCollection(claim.SessionID)
		d.logger.Warn("re-reading the worker's restart count after the collection", "pod", sb.PodName, "error", err)
		return false
	}
	return d.stampRecordedHandoff(ctx, claim, d.handoffOfSession(claim.SessionID))
}

// markRestartsPreCollection flags the session's watch as holding a restart
// baseline from before the collection.
func (d *Dispatcher) markRestartsPreCollection(sessionID string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if h, ok := d.handoffs[sessionID]; ok {
		h.restartsPreCollection = true
		d.handoffs[sessionID] = h
	}
}

// stampRecordedHandoff stamps the handoff on the claim as recorded. On a
// redelivery the current restart count may include the exit that ends the
// turn, and re-reading it would hide that exit. It outlives a cancelled ctx.
func (d *Dispatcher) stampRecordedHandoff(ctx context.Context, claim *sandbox.SandboxClaim, h *handoff) bool {
	if h == nil {
		return false
	}
	ctx, cancel := context.WithTimeout(context.WithoutCancel(ctx), stampTimeout)
	defer cancel()
	if err := d.retryTransient(ctx, func() error { return d.sandboxes.StampHandoff(ctx, claim.Name, h.podUID, h.restarts) }); err != nil {
		d.logger.Warn("recording the handoff on the claim", "claim", claim.Name, "error", err)
		return false
	}
	return true
}

// retryTransient calls try up to five times, pausing PodPollInterval between
// failures, and returns the last error. An absent object ends the retries
// early: no retry brings it back.
func (d *Dispatcher) retryTransient(ctx context.Context, try func() error) error {
	const attempts = 5
	var err error
	for range attempts {
		err = try()
		if err == nil || errors.Is(err, sandbox.ErrNotBound) || apierrors.IsNotFound(err) {
			return err
		}
		sleep(ctx, d.Options.PodPollInterval)
	}
	return err
}
