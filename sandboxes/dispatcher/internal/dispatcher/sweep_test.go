package dispatcher

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	corev1 "k8s.io/api/core/v1"

	"github.com/anthropics/anthropic-cli/sandboxes/dispatcher/internal/cmaapi"
	"github.com/anthropics/anthropic-cli/sandboxes/dispatcher/internal/sandbox"
	"github.com/anthropics/anthropic-cli/sandboxes/dispatcher/internal/workitem"
)

func TestSweepRetention(t *testing.T) {
	sessions := fakeSessions{
		"session_old": cmaapi.SessionAlive, "session_new": cmaapi.SessionAlive, "session_unstamped": cmaapi.SessionAlive,
		"session_running": cmaapi.SessionAlive, "session_running_ended": cmaapi.SessionEnded, "session_gone": cmaapi.SessionGone,
		"session_inflight": cmaapi.SessionGone, "session_just_stamped": cmaapi.SessionAlive,
	}
	d, _, sandboxes := newDispatcher(t, sessions)
	d.Options.Retention = 24 * time.Hour
	d.inflight["session_inflight"] = struct{}{}
	now := d.Now()
	running := map[string]bool{"session_running": true, "session_running_ended": true}
	for id, lastActive := range map[string]time.Time{
		"session_old":           now.Add(-48 * time.Hour),
		"session_new":           now.Add(-time.Hour),
		"session_unstamped":     {},
		"session_running":       now.Add(-48 * time.Hour),
		"session_running_ended": now,
		"session_gone":          now,
		"session_inflight":      now.Add(-48 * time.Hour),
		"session_just_stamped":  now.Add(-48 * time.Hour),
	} {
		claim, _ := sandboxes.EnsureClaim(context.Background(), id)
		claim.LastActive = lastActive
		sandboxes.sandboxes[claim.SandboxName] = &sandbox.Sandbox{Name: claim.SandboxName, Suspended: !running[id]}
	}
	// An exit report that stamps and suspends after the sweep listed the claims.
	sandboxes.onSandboxOf = func(claim *sandbox.SandboxClaim) {
		if claim.SessionID == "session_just_stamped" {
			claim.LastActive = now
		}
	}
	if err := d.Sweep(context.Background()); err != nil {
		t.Fatal(err)
	}
	want := map[string]bool{sandbox.ClaimName("session_old"): true, sandbox.ClaimName("session_gone"): true, sandbox.ClaimName("session_running_ended"): true}
	if len(sandboxes.deleted) != len(want) {
		t.Fatalf("deleted = %v, want %v", sandboxes.deleted, want)
	}
	for _, name := range sandboxes.deleted {
		if !want[name] {
			t.Errorf("deleted %s unexpectedly", name)
		}
	}
}

func TestSweepResolvesAReplacedPodTheWatchMissed(t *testing.T) {
	d, _, sandboxes := newDispatcher(t, fakeSessions{"session_1": cmaapi.SessionAlive})
	claim, _ := sandboxes.EnsureClaim(context.Background(), "session_1")
	pod := sandboxes.runningPod(claim, corev1.ContainerStatus{})
	d.watchHandoff(handoff{claim: *claim, sb: sandboxes.sandboxes[claim.SandboxName], podUID: "the-old-pod"})
	if err := d.Sweep(context.Background()); err != nil {
		t.Fatal(err)
	}
	if len(sandboxes.suspended) != 1 || sandboxes.suspended[0] != claim.SandboxName || len(d.watchedHandoffs()) != 0 {
		t.Errorf("the sweep must resolve %s in place of the old pod and suspend the sandbox: suspended = %v", pod.Name, sandboxes.suspended)
	}
}

func TestDeletingAClaimForgetsTheSession(t *testing.T) {
	d, _, sandboxes := newDispatcher(t, fakeSessions{"session_1": cmaapi.SessionGone})
	claim, _ := sandboxes.EnsureClaim(context.Background(), "session_1")
	d.watchHandoff(handoff{claim: *claim, sb: &sandbox.Sandbox{Name: claim.SandboxName}, podUID: "p"})
	d.workItems.Put("p", workitem.SessionWorkItem{SessionID: "session_1", WorkID: "work_1"})
	_, _ = d.workItems.Collect("p", func(workitem.SessionWorkItem) error { return nil })
	if err := d.Sweep(context.Background()); err != nil {
		t.Fatal(err)
	}
	_, handed := d.workItems.HandedWork("session_1")
	if len(d.watchedHandoffs()) != 0 || handed || len(sandboxes.deleted) != 1 {
		t.Errorf("watched = %v, handed = %v, deleted = %v", d.watchedHandoffs(), handed, sandboxes.deleted)
	}
}

func TestAClaimDeletedOutOfBandForgetsTheSession(t *testing.T) {
	d, _, sandboxes := newDispatcher(t, fakeSessions{"session_1": cmaapi.SessionAlive})
	claim, _ := sandboxes.EnsureClaim(context.Background(), "session_1")
	d.watchHandoff(handoff{claim: *claim, sb: &sandbox.Sandbox{Name: claim.SandboxName}, podUID: "p"})
	d.workItems.Put("p", workitem.SessionWorkItem{SessionID: "session_1", WorkID: "work_1"})
	_, _ = d.workItems.Collect("p", func(workitem.SessionWorkItem) error { return nil })
	// An operator deletes the claim; the dispatcher only finds it gone.
	delete(sandboxes.claims, claim.Name)
	if err := d.Sweep(context.Background()); err != nil {
		t.Fatal(err)
	}
	_, handed := d.workItems.HandedWork("session_1")
	if len(d.watchedHandoffs()) != 0 || handed || len(sandboxes.deleted) != 0 {
		t.Errorf("watched = %v, handed = %v, deleted = %v", d.watchedHandoffs(), handed, sandboxes.deleted)
	}
}

func TestSweepKeepsTheWatchOfAClaimCreatedAfterItsListing(t *testing.T) {
	d, _, sandboxes := newDispatcher(t, fakeSessions{"session_old": cmaapi.SessionAlive, "session_new": cmaapi.SessionAlive})
	d.Options.Retention = time.Hour
	old, _ := sandboxes.EnsureClaim(context.Background(), "session_old")
	sandboxes.runningPod(old, corev1.ContainerStatus{})
	d.watchHandoff(handoff{claim: sandbox.SandboxClaim{Name: sandbox.ClaimName("session_deleted"), SessionID: "session_deleted"}, sb: &sandbox.Sandbox{Name: "sb-deleted"}, podUID: "p-deleted"})
	newName := sandbox.ClaimName("session_new")
	sandboxes.onSandboxOf = func(*sandbox.SandboxClaim) {
		if _, ok := sandboxes.claims[newName]; ok {
			return
		}
		claim := &sandbox.SandboxClaim{Name: newName, SessionID: "session_new", SandboxName: "sb-" + newName}
		sandboxes.claims[newName] = claim
		d.watchHandoff(handoff{claim: *claim, sb: &sandbox.Sandbox{Name: claim.SandboxName}, podUID: "p-new"})
	}
	if err := d.Sweep(context.Background()); err != nil {
		t.Fatal(err)
	}
	watched := d.watchedHandoffs()
	if _, ok := watched["session_new"]; !ok {
		t.Error("a session handed off during the sweep must stay watched")
	}
	if _, ok := watched["session_deleted"]; ok {
		t.Error("a session whose claim is gone must be forgotten")
	}
}

// A proxy in front of the API that has no route to the server answers every
// request 404 with a page of its own. The sessions are still alive, so neither
// the start-up pass nor the sweep deletes a sandbox on that answer.
func TestAProxys404DeletesNoSandbox(t *testing.T) {
	// Given: two claimed sessions, and a CMA API client whose every request a
	// proxy answers 404.
	proxy := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusNotFound)
		_, _ = io.WriteString(w, "404 page not found\n")
	}))
	t.Cleanup(proxy.Close)
	d, _, sandboxes := newDispatcher(t, nil)
	d.sessions = cmaapi.New(proxy.URL, "env_1", "key-1", hangAfter)
	for _, id := range []string{"session_1", "session_2"} {
		_, _ = sandboxes.EnsureClaim(context.Background(), id)
	}

	// When: the dispatcher reconciles and then sweeps.
	if err := d.Reconcile(context.Background()); err != nil {
		t.Fatal(err)
	}
	if err := d.Sweep(context.Background()); err != nil {
		t.Fatal(err)
	}

	// Then: both claims are kept.
	if len(sandboxes.deleted) != 0 {
		t.Errorf("deleted = %v, want none", sandboxes.deleted)
	}
}
