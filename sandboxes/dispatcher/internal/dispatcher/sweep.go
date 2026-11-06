package dispatcher

import (
	"context"
	"fmt"
	"time"

	apierrors "k8s.io/apimachinery/pkg/api/errors"

	"github.com/anthropics/anthropic-cli/sandboxes/dispatcher/internal/sandbox"
)

func (d *Dispatcher) sweepLoop(ctx context.Context) {
	ticker := time.NewTicker(d.Options.SweepInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := d.Sweep(ctx); err != nil && ctx.Err() == nil {
				d.logger.Error("sweep", "error", err)
			}
		}
	}
}

// Sweep deletes every managed claim whose session is over, running pod
// included, and, under a retention limit, suspended claims idle past it;
// then it checks every watched session whose pod change the watch may have
// missed. A session is claimed for the duration of its check, so a handoff
// never runs against it at the same time.
func (d *Dispatcher) Sweep(ctx context.Context) error {
	claims, err := d.sandboxes.ListClaims(ctx)
	if err != nil {
		return fmt.Errorf("listing sandbox claims: %w", err)
	}
	claimed := map[string]bool{}
	for i := range claims {
		claim := &claims[i]
		claimed[claim.SessionID] = true
		if !d.claimInflight(claim.SessionID) {
			continue
		}
		if !d.deleteIfSessionOver(ctx, claim) && d.Options.Retention > 0 && d.suspendedPastRetention(ctx, claim) {
			d.deleteClaim(ctx, claim, "idle past retention")
		}
		d.releaseInflight(claim.SessionID)
	}
	for sessionID, h := range d.watchedHandoffs() {
		if !claimed[sessionID] {
			d.forgetIfClaimGone(ctx, h)
			continue
		}
		pod, err := d.sandboxes.PodOf(ctx, h.sb)
		switch {
		case apierrors.IsNotFound(err):
			if h.podUID != "" {
				d.PodGone(h.podUID)
			}
		case err != nil:
			d.logger.Warn("could not read a watched session's pod", "pod", h.sb.PodName, "session_id", sessionID, "error", err)
		default:
			if h.podUID != "" && h.podUID != pod.UID {
				d.PodGone(h.podUID)
			}
			d.PodChanged(ctx, pod)
		}
	}
	return nil
}

// suspendedPastRetention reads the sandbox and then the claim's last-active
// stamp fresh, in that order: the listing may predate a suspend that stamped
// the claim and then suspended the sandbox.
func (d *Dispatcher) suspendedPastRetention(ctx context.Context, claim *sandbox.SandboxClaim) bool {
	sb, err := d.sandboxes.SandboxOf(ctx, claim)
	if err != nil {
		d.warnUnlessAbsent(err, "could not read the claim's sandbox; leaving it", claim)
		return false
	}
	if !sb.Suspended {
		return false
	}
	fresh, err := d.sandboxes.GetClaim(ctx, claim.Name)
	if err != nil {
		d.logger.Warn("could not re-read the claim; leaving it", "claim", claim.Name, "error", err)
		return false
	}
	return !fresh.LastActive.IsZero() && d.Now().Sub(fresh.LastActive) > d.Options.Retention
}

// forgetIfClaimGone stops watching a session the sweep's listing lacks once a
// fresh read finds its claim gone. It also drops the session's workItems
// record, which nothing else drops for a claim deleted out of band.
func (d *Dispatcher) forgetIfClaimGone(ctx context.Context, h handoff) {
	if !d.claimInflight(h.claim.SessionID) {
		return
	}
	defer d.releaseInflight(h.claim.SessionID)
	if _, err := d.sandboxes.GetClaim(ctx, h.claim.Name); apierrors.IsNotFound(err) {
		d.forgetHandoff(h.claim.SessionID)
		d.workItems.ForgetSession(h.claim.SessionID)
	}
}
