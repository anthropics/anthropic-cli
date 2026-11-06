package dispatcher

import (
	"context"
	"errors"
	"testing"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"

	"github.com/anthropics/anthropic-cli/sandboxes/dispatcher/internal/cmaapi"
	"github.com/anthropics/anthropic-cli/sandboxes/dispatcher/internal/sandbox"
	"github.com/anthropics/anthropic-cli/sandboxes/dispatcher/internal/workitem"
)

func TestReconcileDeletesOnlyAuthoritativeAnswers(t *testing.T) {
	sessions := fakeSessions{
		"session_alive": cmaapi.SessionAlive,
		"session_gone":  cmaapi.SessionGone,
		"session_ended": cmaapi.SessionEnded,
		"session_err":   errors.New("401 key rotated"),
	}
	d, _, sandboxes := newDispatcher(t, sessions)
	for id := range sessions {
		_, _ = sandboxes.EnsureClaim(context.Background(), id)
	}
	sandboxes.claims["unlabeled"] = &sandbox.SandboxClaim{Name: "unlabeled"}
	if err := d.Reconcile(context.Background()); err != nil {
		t.Fatal(err)
	}
	want := map[string]bool{sandbox.ClaimName("session_gone"): true, sandbox.ClaimName("session_ended"): true}
	if len(sandboxes.deleted) != len(want) {
		t.Fatalf("deleted = %v, want %v", sandboxes.deleted, want)
	}
	for _, name := range sandboxes.deleted {
		if !want[name] {
			t.Errorf("deleted %s unexpectedly", name)
		}
	}
}

func TestReconcileAdoptsPodsFromBeforeTheRestart(t *testing.T) {
	sessions := fakeSessions{"session_running": cmaapi.SessionAlive, "session_exited": cmaapi.SessionAlive, "session_fresh": cmaapi.SessionAlive, "session_gone": cmaapi.SessionGone}
	d, _, sandboxes := newDispatcher(t, sessions)
	now := d.Now()
	handoff := now.Add(-time.Hour)
	running := corev1.ContainerStatus{State: corev1.ContainerState{Running: &corev1.ContainerStateRunning{}}}
	for id, tc := range map[string]struct {
		lastActive time.Time
		// stamped marks a claim whose handoff was recorded: its work item
		// was collected before the restart.
		stamped bool
		status  corev1.ContainerStatus
	}{
		"session_running": {handoff, true, running},
		"session_exited":  {handoff, true, corev1.ContainerStatus{State: running.State, LastTerminationState: exitedAt(now.Add(-time.Minute)), RestartCount: 1}},
		"session_fresh":   {time.Time{}, false, running},
		"session_gone":    {handoff, true, running},
	} {
		claim, _ := sandboxes.EnsureClaim(context.Background(), id)
		claim.LastActive = tc.lastActive
		pod := sandboxes.runningPod(claim, tc.status)
		if tc.stamped {
			claim.HandoffPodUID = pod.UID
			claim.HandoffRestarts = ptr(int32(0))
		}
	}
	if err := d.Reconcile(context.Background()); err != nil {
		t.Fatal(err)
	}
	watched := d.watchedHandoffs()
	if h := watched["session_running"]; h.podUID != types.UID("pod-"+sandbox.ClaimName("session_running")) {
		t.Errorf("session_running: watched = %+v, want its running pod adopted", h)
	}
	if h, ok := watched["session_fresh"]; !ok || h.podUID != "" {
		t.Errorf("session_fresh: watched = %+v, %v; a claim with no recorded handoff must await its next pod", h, ok)
	}
	if len(watched) != 2 {
		t.Errorf("watched = %v, want only the two running sessions", watched)
	}
	if len(sandboxes.suspended) != 1 || sandboxes.suspended[0] != "sb-"+sandbox.ClaimName("session_exited") {
		t.Errorf("suspended = %v, want the session whose worker exited after its handoff", sandboxes.suspended)
	}
	if len(sandboxes.deleted) != 1 || sandboxes.deleted[0] != sandbox.ClaimName("session_gone") {
		t.Errorf("deleted = %v", sandboxes.deleted)
	}
}

func TestReconcileWatchesASandboxWhosePodIsBeingReplaced(t *testing.T) {
	d, _, sandboxes := newDispatcher(t, fakeSessions{"session_recreating": cmaapi.SessionAlive, "session_creating": cmaapi.SessionAlive})
	recreating, _ := sandboxes.EnsureClaim(context.Background(), "session_recreating")
	recreating.LastActive = d.Now()
	sandboxes.sandboxes[recreating.SandboxName] = &sandbox.Sandbox{Name: recreating.SandboxName, PodName: recreating.SandboxName}
	creating, _ := sandboxes.EnsureClaim(context.Background(), "session_creating")
	creating.LastActive = d.Now()
	starting := sandboxes.runningPod(creating, corev1.ContainerStatus{})
	starting.Status.ContainerStatuses = nil

	if err := d.Reconcile(context.Background()); err != nil {
		t.Fatal(err)
	}
	watched := d.watchedHandoffs()
	for _, id := range []string{"session_recreating", "session_creating"} {
		if h, ok := watched[id]; !ok || h.podUID != "" {
			t.Errorf("%s: watched = %+v, %v; want it watched as awaiting its next pod", id, h, ok)
		}
	}
	replacement := sandboxes.runningPod(recreating, corev1.ContainerStatus{})
	d.PodChanged(context.Background(), replacement)
	if len(sandboxes.suspended) != 1 || sandboxes.suspended[0] != recreating.SandboxName {
		t.Errorf("suspended = %v, want the sandbox whose pod was replaced suspended", sandboxes.suspended)
	}
}

// Reconcile runs once per process, so a session it fails to adopt stays
// unwatched for the life of the dispatcher.
func TestReconcileRetriesTransientReadsBeforeAdopting(t *testing.T) {
	d, _, sandboxes := newDispatcher(t, fakeSessions{"session_1": cmaapi.SessionAlive})
	claim, _ := sandboxes.EnsureClaim(context.Background(), "session_1")
	claim.LastActive = d.Now()
	pod := sandboxes.runningPod(claim, corev1.ContainerStatus{State: corev1.ContainerState{Running: &corev1.ContainerStateRunning{}}})
	claim.HandoffPodUID = pod.UID
	claim.HandoffRestarts = ptr(int32(0))
	sandboxes.mu.Lock()
	sandboxes.sandboxOfFailures = 2
	sandboxes.podOfFailures = 2
	sandboxes.mu.Unlock()
	if err := d.Reconcile(context.Background()); err != nil {
		t.Fatal(err)
	}
	if h, ok := d.watchedHandoffs()["session_1"]; !ok || h.podUID != pod.UID {
		t.Errorf("watch = %+v, %v; a transient read failure during the boot reconcile must not leave the session unwatched", h, ok)
	}
}

func TestReconcileKeepsWatchingAfterAFailedSuspend(t *testing.T) {
	d, _, sandboxes := newDispatcher(t, fakeSessions{"session_1": cmaapi.SessionAlive})
	claim, _ := sandboxes.EnsureClaim(context.Background(), "session_1")
	claim.LastActive = d.Now().Add(-time.Hour)
	pod := sandboxes.runningPod(claim, corev1.ContainerStatus{State: exitedAt(d.Now())})
	claim.HandoffPodUID = pod.UID
	claim.HandoffRestarts = ptr(int32(0))
	sandboxes.mu.Lock()
	sandboxes.suspendFailures = 1
	sandboxes.mu.Unlock()
	if err := d.Reconcile(context.Background()); err != nil {
		t.Fatal(err)
	}
	if _, ok := d.watchedHandoffs()["session_1"]; !ok {
		t.Fatal("a session whose suspend failed must stay watched, or its sandbox is never reclaimed")
	}
	d.PodChanged(context.Background(), pod)
	if len(sandboxes.suspended) != 1 {
		t.Errorf("suspended = %v; the next pod event must retry the suspend", sandboxes.suspended)
	}
}

// A worker whose exit shows only as a restart, on a claim whose restart
// annotation did not survive, must still be re-judged done after a failed
// suspend: a baseline at the post-restart count would never read as exceeded
// while the restarted run lives, leaving the sandbox running forever.
func TestReconcileKeepsWatchingAfterAFailedSuspendOfARestartedWorker(t *testing.T) {
	d, _, sandboxes := newDispatcher(t, fakeSessions{"session_1": cmaapi.SessionAlive})
	claim, _ := sandboxes.EnsureClaim(context.Background(), "session_1")
	claim.LastActive = d.Now().Add(-time.Hour)
	pod := sandboxes.runningPod(claim, corev1.ContainerStatus{
		State:                corev1.ContainerState{Running: &corev1.ContainerStateRunning{}},
		LastTerminationState: exitedAt(d.Now()),
		RestartCount:         1,
	})
	claim.HandoffPodUID = pod.UID
	sandboxes.mu.Lock()
	sandboxes.suspendFailures = 1
	sandboxes.mu.Unlock()
	if err := d.Reconcile(context.Background()); err != nil {
		t.Fatal(err)
	}
	d.PodChanged(context.Background(), pod)
	if len(sandboxes.suspended) != 1 {
		t.Errorf("suspended = %v; the next pod event must retry the failed suspend", sandboxes.suspended)
	}
}

func TestReconcileJudgesTheWorkerByTheRecordedRestartCount(t *testing.T) {
	d, _, sandboxes := newDispatcher(t, fakeSessions{"session_restarted": cmaapi.SessionAlive, "session_tied": cmaapi.SessionAlive})
	now := d.Now()
	running := corev1.ContainerState{Running: &corev1.ContainerStateRunning{}}
	for id, tc := range map[string]struct {
		recorded int32
		status   corev1.ContainerStatus
	}{
		// The worker restarted after the recorded count: its run is over.
		"session_restarted": {0, corev1.ContainerStatus{State: running, LastTerminationState: exitedAt(now), RestartCount: 1}},
		// A pre-handoff restart in the same second as the stamp: still running.
		"session_tied": {1, corev1.ContainerStatus{State: running, LastTerminationState: exitedAt(now), RestartCount: 1}},
	} {
		claim, _ := sandboxes.EnsureClaim(context.Background(), id)
		claim.LastActive = now
		claim.HandoffRestarts = ptr(tc.recorded)
		claim.HandoffPodUID = sandboxes.runningPod(claim, tc.status).UID
	}
	if err := d.Reconcile(context.Background()); err != nil {
		t.Fatal(err)
	}
	if len(sandboxes.suspended) != 1 || sandboxes.suspended[0] != "sb-"+sandbox.ClaimName("session_restarted") {
		t.Errorf("suspended = %v, want only the session whose worker restarted past the recorded count", sandboxes.suspended)
	}
	if _, watched := d.watchedHandoffs()["session_tied"]; !watched {
		t.Error("a worker whose only exit predates the recorded count must stay watched")
	}
}

// A restart inside a session's first handoff window leaves a claim with no
// recorded handoff and an idle Running worker that never got its token.
// Adopting that pod as mid-turn would refuse the redelivered item forever.
func TestReconcileLeavesAnUnhandedPodAwaiting(t *testing.T) {
	d, queue, sandboxes := newDispatcher(t, fakeSessions{"session_1": cmaapi.SessionAlive})
	claim, _ := sandboxes.EnsureClaim(context.Background(), "session_1")
	pod := sandboxes.runningPod(claim, corev1.ContainerStatus{State: corev1.ContainerState{Running: &corev1.ContainerStateRunning{}}})
	if err := d.Reconcile(context.Background()); err != nil {
		t.Fatal(err)
	}
	if h, ok := d.watchedHandoffs()["session_1"]; !ok || h.podUID != "" {
		t.Fatalf("watch = %+v, %v; a pod with no recorded handoff must be watched as awaiting the session's next pod", h, ok)
	}
	// The never-acked item redelivers; the pod's fetch still polls for it.
	stop := make(chan struct{})
	defer close(stop)
	go func() {
		for {
			select {
			case <-stop:
				return
			default:
			}
			_, _ = d.workItems.Collect(pod.UID, func(workitem.SessionWorkItem) error { return nil })
			time.Sleep(time.Millisecond)
		}
	}()
	if err := d.Handle(context.Background(), &cmaapi.WorkItem{ID: "work_1", SessionID: "session_1"}); err != nil {
		t.Fatalf("the redelivered item must reach the session again, got: %v", err)
	}
	if acked := queue.ackedIDs(); len(acked) != 1 {
		t.Errorf("acked = %v, want the redelivered item handed off and acked", acked)
	}
}

func TestReconcileNeverAdoptsAReplacementPod(t *testing.T) {
	d, _, sandboxes := newDispatcher(t, fakeSessions{"session_1": cmaapi.SessionAlive})
	claim, _ := sandboxes.EnsureClaim(context.Background(), "session_1")
	claim.LastActive = d.Now()
	claim.HandoffRestarts = ptr(int32(0))
	claim.HandoffPodUID = "the-handed-pod"
	// The handed pod was replaced: same Sandbox, fresh pod, worker Running at
	// restart count zero.
	replacement := sandboxes.runningPod(claim, corev1.ContainerStatus{State: corev1.ContainerState{Running: &corev1.ContainerStateRunning{}}})
	if err := d.Reconcile(context.Background()); err != nil {
		t.Fatal(err)
	}
	if h, ok := d.watchedHandoffs()["session_1"]; !ok || h.podUID != "" {
		t.Fatalf("watch = %+v, %v; a replacement pod must be watched as awaiting the next pod, never adopted", h, ok)
	}
	d.PodChanged(context.Background(), replacement)
	if len(sandboxes.suspended) != 1 {
		t.Errorf("suspended = %v, want the replaced session's sandbox suspended", sandboxes.suspended)
	}
}
