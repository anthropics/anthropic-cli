package dispatcher

import (
	"context"
	"testing"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/anthropics/anthropic-cli/sandboxes/dispatcher/internal/cmaapi"
	"github.com/anthropics/anthropic-cli/sandboxes/dispatcher/internal/sandbox"
	"github.com/anthropics/anthropic-cli/sandboxes/dispatcher/internal/workitem"
)

func TestPodChangedParksOrDeletesAfterTheWorkerExits(t *testing.T) {
	sessions := fakeSessions{"session_alive": cmaapi.SessionAlive, "session_ended": cmaapi.SessionEnded}
	d, _, sandboxes := newDispatcher(t, sessions)
	now := d.Now()
	for id := range sessions {
		claim, _ := sandboxes.EnsureClaim(context.Background(), id)
		claim.LastActive = now
		pod := sandboxes.runningPod(claim, corev1.ContainerStatus{State: corev1.ContainerState{Running: &corev1.ContainerStateRunning{}}})
		d.watchHandoff(handoff{claim: *claim, sb: sandboxes.sandboxes[claim.SandboxName], podUID: pod.UID})
		d.PodChanged(context.Background(), pod)
		if _, watched := d.handoffOfPod(pod.UID); !watched {
			t.Errorf("%s: a running worker must stay watched", id)
		}
		sandboxes.mu.Lock()
		pod.Status.ContainerStatuses[0].State = exitedAt(now.Add(time.Minute))
		sandboxes.mu.Unlock()
		d.PodChanged(context.Background(), pod)
		if _, watched := d.handoffOfPod(pod.UID); watched {
			t.Errorf("%s: an exited worker must be dealt with and forgotten", id)
		}
	}
	alive, ended := sandbox.ClaimName("session_alive"), sandbox.ClaimName("session_ended")
	if len(sandboxes.suspended) != 1 || sandboxes.suspended[0] != "sb-"+alive || len(sandboxes.stamped) != 1 {
		t.Errorf("alive session: suspended = %v, stamped = %v; want its sandbox suspended", sandboxes.suspended, sandboxes.stamped)
	}
	if len(sandboxes.deleted) != 1 || sandboxes.deleted[0] != ended {
		t.Errorf("ended session: deleted = %v", sandboxes.deleted)
	}
}

func TestAReplacedPodEndsTheSession(t *testing.T) {
	d, _, sandboxes := newDispatcher(t, fakeSessions{"session_1": cmaapi.SessionAlive})
	claim, _ := sandboxes.EnsureClaim(context.Background(), "session_1")
	pod := sandboxes.runningPod(claim, corev1.ContainerStatus{})
	d.watchHandoff(handoff{claim: *claim, sb: sandboxes.sandboxes[claim.SandboxName], podUID: pod.UID})
	d.workItems.Put(pod.UID, workitem.SessionWorkItem{SessionID: "session_1", WorkID: "work_1"})
	_, _ = d.workItems.Collect(pod.UID, func(workitem.SessionWorkItem) error { return nil })

	// The node drains: the pod goes, and the controller makes another.
	delete(sandboxes.pods, claim.SandboxName)
	d.PodGone(pod.UID)
	if h, watched := d.watchedHandoffs()["session_1"]; !watched || h.podUID != "" {
		t.Fatalf("a gone pod must leave its session watched, awaiting the next pod: %+v, %v", h, watched)
	}
	if _, handed := d.workItems.HandedWork("session_1"); !handed {
		t.Error("the session's handoff record must outlive its pod")
	}
	stranger := &corev1.Pod{ObjectMeta: metav1.ObjectMeta{Name: claim.SandboxName, UID: "stranger"}}
	d.PodChanged(context.Background(), stranger)
	if len(sandboxes.suspended) != 0 {
		t.Error("a same-named pod the Sandbox does not name must be ignored")
	}
	replacement := sandboxes.runningPod(claim, corev1.ContainerStatus{})
	replacement.UID = "replacement"
	d.PodChanged(context.Background(), replacement)
	if len(sandboxes.suspended) != 1 || len(d.watchedHandoffs()) != 0 {
		t.Errorf("the replacement holds no session, so the sandbox is suspended: suspended = %v, watched = %v", sandboxes.suspended, d.watchedHandoffs())
	}
}

func TestPodChangedSeesARestartedWorker(t *testing.T) {
	d, _, sandboxes := newDispatcher(t, fakeSessions{"session_1": cmaapi.SessionAlive})
	claim, _ := sandboxes.EnsureClaim(context.Background(), "session_1")
	pod := sandboxes.runningPod(claim, corev1.ContainerStatus{State: corev1.ContainerState{Running: &corev1.ContainerStateRunning{}}, RestartCount: 2})
	d.watchHandoff(handoff{claim: *claim, sb: sandboxes.sandboxes[claim.SandboxName], podUID: pod.UID, restarts: 2})
	sandboxes.mu.Lock()
	pod.Status.ContainerStatuses[0].RestartCount = 3
	sandboxes.mu.Unlock()
	d.PodChanged(context.Background(), pod)
	if len(sandboxes.suspended) != 1 {
		t.Errorf("a restart above the handoff's count must suspend the sandbox: suspended = %v", sandboxes.suspended)
	}
}

func TestTheAssignmentClosesWhenTheWorkerIsDone(t *testing.T) {
	d, _, sandboxes := newDispatcher(t, fakeSessions{"session_1": cmaapi.SessionAlive})
	claim, _ := sandboxes.EnsureClaim(context.Background(), "session_1")
	pod := sandboxes.runningPod(claim, corev1.ContainerStatus{State: corev1.ContainerState{Running: &corev1.ContainerStateRunning{}}})
	d.watchHandoff(handoff{claim: *claim, sb: sandboxes.sandboxes[claim.SandboxName], podUID: pod.UID})
	d.workItems.Put(pod.UID, workitem.SessionWorkItem{SessionID: "session_1", WorkID: "work_1"})
	collect := func() workitem.Status {
		status, _ := d.workItems.Collect(pod.UID, func(workitem.SessionWorkItem) error { return nil })
		return status
	}
	_ = collect()

	d.PodChanged(context.Background(), pod)
	if got := collect(); got != workitem.CollectedAgain {
		t.Fatalf("while the worker runs: status %v, want the work item still open", got)
	}
	sandboxes.mu.Lock()
	pod.Status.ContainerStatuses[0].State = exitedAt(d.Now())
	sandboxes.mu.Unlock()
	d.PodChanged(context.Background(), pod)
	if got := collect(); got != workitem.AlreadyCollected {
		t.Errorf("after the worker exited: status %v, want the work item closed", got)
	}
}

func TestPodChangedWaitsWhileTheSessionIsBusy(t *testing.T) {
	d, _, sandboxes := newDispatcher(t, fakeSessions{"session_1": cmaapi.SessionAlive})
	claim, _ := sandboxes.EnsureClaim(context.Background(), "session_1")
	pod := sandboxes.runningPod(claim, corev1.ContainerStatus{State: exitedAt(d.Now())})
	d.watchHandoff(handoff{claim: *claim, sb: sandboxes.sandboxes[claim.SandboxName], podUID: pod.UID})
	d.claimInflight("session_1")
	d.PodChanged(context.Background(), pod)
	if _, watched := d.handoffOfPod(pod.UID); !watched || len(sandboxes.suspended) != 0 {
		t.Errorf("a busy session must be left for later: watched = %v, suspended = %v", watched, sandboxes.suspended)
	}
	d.releaseInflight("session_1")
	if err := d.Sweep(context.Background()); err != nil {
		t.Fatal(err)
	}
	if _, watched := d.handoffOfPod(pod.UID); watched || len(sandboxes.suspended) != 1 {
		t.Errorf("the sweep must pick up the missed exit: watched = %v, suspended = %v", watched, sandboxes.suspended)
	}
}

// A worker exit observed while a handoff holds the session must not close
// the pod's delivered work item. The open item is what makes Put refuse new
// work for the pod, whose worker is no longer running to take it.
func TestABusySessionKeepsItsDeliveredWorkItem(t *testing.T) {
	d, _, sandboxes := newDispatcher(t, fakeSessions{"session_1": cmaapi.SessionAlive})
	claim, _ := sandboxes.EnsureClaim(context.Background(), "session_1")
	pod := sandboxes.runningPod(claim, corev1.ContainerStatus{State: corev1.ContainerState{Running: &corev1.ContainerStateRunning{}}})
	d.watchHandoff(handoff{claim: *claim, sb: sandboxes.sandboxes[claim.SandboxName], podUID: pod.UID})
	mustPut(t, d, pod.UID, workitem.SessionWorkItem{SessionID: "session_1", WorkID: "work_1"})
	if _, err := d.workItems.Collect(pod.UID, func(workitem.SessionWorkItem) error { return nil }); err != nil {
		t.Fatal(err)
	}

	sandboxes.mu.Lock()
	pod.Status.ContainerStatuses[0].State = exitedAt(d.Now())
	sandboxes.mu.Unlock()
	d.claimInflight("session_1")
	d.PodChanged(context.Background(), pod)
	if _, opened := d.workItems.Put(pod.UID, workitem.SessionWorkItem{SessionID: "session_1", WorkID: "work_2"}); opened {
		t.Fatal("Put opened new work for the pod while its delivered work item should still be open")
	}
	d.releaseInflight("session_1")

	if err := d.Sweep(context.Background()); err != nil {
		t.Fatal(err)
	}
	if status, _ := d.workItems.Collect(pod.UID, func(workitem.SessionWorkItem) error { return nil }); status != workitem.AlreadyCollected {
		t.Errorf("after the sweep: status %v, want the work item closed", status)
	}
	if len(sandboxes.suspended) != 1 {
		t.Errorf("after the sweep: suspended = %v, want the sandbox suspended", sandboxes.suspended)
	}
}

func TestAStaleExitEventNeverSuspendsARunningWorker(t *testing.T) {
	d, _, sandboxes := newDispatcher(t, fakeSessions{"session_1": cmaapi.SessionAlive})
	claim, _ := sandboxes.EnsureClaim(context.Background(), "session_1")
	// The handoff tolerated a mid-wait restart: the watch's baseline is 1.
	pod := sandboxes.runningPod(claim, corev1.ContainerStatus{State: corev1.ContainerState{Running: &corev1.ContainerStateRunning{}}, RestartCount: 1})
	d.watchHandoff(handoff{claim: *claim, sb: sandboxes.sandboxes[claim.SandboxName], podUID: pod.UID, restarts: 1})
	mustPut(t, d, pod.UID, workitem.SessionWorkItem{SessionID: "session_1", WorkID: "work_1"})
	if _, err := d.workItems.Collect(pod.UID, func(workitem.SessionWorkItem) error { return nil }); err != nil {
		t.Fatal(err)
	}
	// A stale watch event from before the restart was tolerated.
	stale := pod.DeepCopy()
	stale.Status.ContainerStatuses[0] = corev1.ContainerStatus{Name: sandbox.WorkerContainer, State: exitedAt(d.Now()), RestartCount: 0}
	d.PodChanged(context.Background(), stale)
	if len(sandboxes.suspended) != 0 {
		t.Fatalf("suspended = %v; a stale exit snapshot must not suspend the running worker", sandboxes.suspended)
	}
	if status, _ := d.workItems.Collect(pod.UID, func(workitem.SessionWorkItem) error { return nil }); status != workitem.CollectedAgain {
		t.Errorf("Collect = %v, want the work item still open", status)
	}
	// The worker really exits; the same event now confirms and suspends.
	sandboxes.mu.Lock()
	sandboxes.pods[claim.SandboxName].Status.ContainerStatuses[0] = corev1.ContainerStatus{Name: sandbox.WorkerContainer, State: exitedAt(d.Now()), RestartCount: 1}
	sandboxes.mu.Unlock()
	d.PodChanged(context.Background(), stale)
	if len(sandboxes.suspended) != 1 {
		t.Errorf("suspended = %v, want the sandbox suspended after a confirmed exit", sandboxes.suspended)
	}
}
