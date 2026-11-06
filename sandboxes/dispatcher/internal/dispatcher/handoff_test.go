package dispatcher

import (
	"context"
	"errors"
	"testing"
	"testing/synctest"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"

	"github.com/anthropics/anthropic-cli/sandboxes/dispatcher/internal/cmaapi"
	"github.com/anthropics/anthropic-cli/sandboxes/dispatcher/internal/sandbox"
	"github.com/anthropics/anthropic-cli/sandboxes/dispatcher/internal/workitem"
)

func TestHandleAcksOnlyAfterHandoff(t *testing.T) {
	d, queue, sandboxes := newDispatcher(t, nil)
	item := &cmaapi.WorkItem{ID: "work_1", SessionID: "session_1", Secret: "st"}
	podUID := types.UID("pod-" + sandbox.ClaimName("session_1"))

	done := make(chan error, 1)
	go func() { done <- d.Handle(context.Background(), item) }()

	var got workitem.SessionWorkItem
	deadline := time.Now().Add(hangAfter)
	for {
		status, _ := d.workItems.Collect(podUID, func(a workitem.SessionWorkItem) error { got = a; return nil })
		if status == workitem.Collected {
			break
		}
		if time.Now().After(deadline) {
			t.Fatal("no work item appeared for the pod")
		}
		if len(queue.ackedIDs()) != 0 {
			t.Fatal("acked before the pod collected its work item")
		}
		time.Sleep(time.Millisecond)
	}
	if got.Secret != "st" || got.WorkID != "work_1" || got.BaseURL != "https://api" || got.EnvironmentID != "env_1" {
		t.Errorf("work item = %+v", got)
	}
	if err := <-done; err != nil {
		t.Fatal(err)
	}
	if acked := queue.ackedIDs(); len(acked) != 1 || acked[0] != "work_1" {
		t.Errorf("acked = %v", acked)
	}
	if len(sandboxes.stamped) != 1 {
		t.Errorf("last-active stamped %d times", len(sandboxes.stamped))
	}
	if h, watched := d.handoffOfPod(podUID); !watched || h.claim.SessionID != "session_1" {
		t.Error("the handed-off pod must be watched for its worker's exit")
	}
}

func TestHandleTimesOutWithoutAck(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		d, queue, _ := newDispatcher(t, nil)
		d.Options.HandoffWait = expiringHandoffWait
		item := &cmaapi.WorkItem{ID: "work_1", SessionID: "session_1", Secret: "st"}
		err := d.Handle(context.Background(), item)
		if !errors.Is(err, ErrHandoffTimeout) {
			t.Fatalf("err = %v, want ErrHandoffTimeout", err)
		}
		if len(queue.ackedIDs()) != 0 {
			t.Error("an uncollected work item must not be acked")
		}
		if status, _ := d.workItems.Collect(types.UID("pod-"+sandbox.ClaimName("session_1")), func(workitem.SessionWorkItem) error { return nil }); status != workitem.NotPending {
			t.Error("the work item must be dropped after the timeout")
		}
		if len(d.watchedHandoffs()) != 0 {
			t.Error("a pod that never collected must not stay watched")
		}
	})
}

func TestHandleAcksWhenCollectionRacesTheDeadline(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		d, queue, _ := newDispatcher(t, nil)
		d.Options.HandoffWait = expiringHandoffWait
		item := &cmaapi.WorkItem{ID: "work_1", SessionID: "session_1"}
		podUID := types.UID("pod-" + sandbox.ClaimName("session_1"))
		collected := make(chan struct{})
		go func() {
			defer close(collected)
			// Collect at the instant the handoff wait ends. The collection then races the deadline.
			time.Sleep(d.Options.HandoffWait)
			_, _ = d.workItems.Collect(podUID, func(workitem.SessionWorkItem) error { return nil })
		}()
		err := d.Handle(context.Background(), item)
		<-collected
		acked := queue.ackedIDs()
		_, handed := d.workItems.HandedWork("session_1")
		if handed && (err != nil || len(acked) != 1) {
			t.Errorf("collected but not acked: err = %v, acked = %v", err, acked)
		}
		if !handed && (!errors.Is(err, ErrHandoffTimeout) || len(acked) != 0) {
			t.Errorf("not collected but err = %v, acked = %v", err, acked)
		}
	})
}

// A delivery still in flight when the handoff wait ends finishes before
// Handle decides. The pod has then collected the work item, and Handle must ack it.
func TestHandleAcksWhenTheWaitEndsMidDelivery(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		d, queue, _ := newDispatcher(t, nil)
		d.Options.HandoffWait = expiringHandoffWait
		item := &cmaapi.WorkItem{ID: "work_1", SessionID: "session_1"}
		podUID := types.UID("pod-" + sandbox.ClaimName("session_1"))
		delivered := make(chan struct{})
		go func() {
			defer close(delivered)
			// The delivery starts 1 ms before the wait ends and takes 2 ms.
			time.Sleep(d.Options.HandoffWait - time.Millisecond)
			_, _ = d.workItems.Collect(podUID, func(workitem.SessionWorkItem) error {
				time.Sleep(2 * time.Millisecond)
				return nil
			})
		}()
		err := d.Handle(context.Background(), item)
		<-delivered
		if err != nil {
			t.Fatalf("err = %v, want nil; the pod has collected the work item, and Handle must ack it", err)
		}
		if acked := queue.ackedIDs(); len(acked) != 1 || acked[0] != "work_1" {
			t.Errorf("acked = %v, want [work_1]", acked)
		}
	})
}

func TestHandleHandsARedeliveryToTheReplacementPod(t *testing.T) {
	d, queue, sandboxes := newDispatcher(t, nil)
	item := &cmaapi.WorkItem{ID: "work_1", SessionID: "session_1"}
	// The work was collected by a pod the sandbox has since replaced, so
	// its worker went with it and the item is run by the replacement.
	mustPut(t, d, "an-earlier-pod", workitem.SessionWorkItem{SessionID: "session_1", WorkID: "work_1"})
	if status, _ := d.workItems.Collect("an-earlier-pod", func(workitem.SessionWorkItem) error { return nil }); status != workitem.Collected {
		t.Fatal("setup: collect")
	}
	claim, _ := sandboxes.EnsureClaim(context.Background(), "session_1")
	claim.HandoffPodUID = "an-earlier-pod"
	claim.HandoffRestarts = ptr(int32(0))
	sandboxes.sandboxes[claim.SandboxName] = &sandbox.Sandbox{Name: claim.SandboxName, Suspended: false}
	done := make(chan error, 1)
	go func() { done <- d.Handle(context.Background(), item) }()
	collectWhenPending(t, d, types.UID("pod-"+claim.Name), done)
	if err := <-done; err != nil {
		t.Fatal(err)
	}
	if acked := queue.ackedIDs(); len(acked) != 1 || acked[0] != "work_1" {
		t.Errorf("acked = %v", acked)
	}
	if err := d.Handle(context.Background(), &cmaapi.WorkItem{ID: "work_2", SessionID: "session_1"}); !failedWithoutTimeout(err) {
		t.Errorf("err = %v, want an error other than a timeout; a new work item for the session must be left for redelivery while its pod runs the earlier one", err)
	}
	if acked := queue.ackedIDs(); len(acked) != 1 {
		t.Errorf("acked = %v; the undelivered new item must not be acked", acked)
	}
}

func TestHandleAcksAfterCancel(t *testing.T) {
	d, queue, _ := newDispatcher(t, nil)
	item := &cmaapi.WorkItem{ID: "work_1", SessionID: "session_1"}
	podUID := types.UID("pod-" + sandbox.ClaimName("session_1"))
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() { done <- d.Handle(ctx, item) }()
	if err := deliverWhenPending(t, d, podUID, done, func(workitem.SessionWorkItem) error { cancel(); return nil }); err != nil {
		t.Fatal(err)
	}
	if err := <-done; err != nil {
		t.Fatal(err)
	}
	if acked := queue.ackedIDs(); len(acked) != 1 {
		t.Errorf("a pod holding the token must be acked even during shutdown: acked = %v", acked)
	}
}

func TestHandleStampsTheHandoffAfterCancel(t *testing.T) {
	d, queue, sandboxes := newDispatcher(t, nil)
	item := &cmaapi.WorkItem{ID: "work_1", SessionID: "session_1"}
	podUID := types.UID("pod-" + sandbox.ClaimName("session_1"))
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() { done <- d.Handle(ctx, item) }()
	if err := deliverWhenPending(t, d, podUID, done, func(workitem.SessionWorkItem) error { cancel(); return nil }); err != nil {
		t.Fatal(err)
	}
	if err := <-done; err != nil {
		t.Fatal(err)
	}
	if acked := queue.ackedIDs(); len(acked) != 1 {
		t.Fatalf("acked = %v", acked)
	}
	claim, _ := sandboxes.GetClaim(context.Background(), sandbox.ClaimName("session_1"))
	if claim.HandoffPodUID != podUID || claim.HandoffRestarts == nil {
		t.Errorf("HandoffPodUID = %q, HandoffRestarts = %v; the stamp must survive a shutdown like the ack, or the next dispatcher judges the session by its previous handoff and suspends the live pod mid-turn", claim.HandoffPodUID, claim.HandoffRestarts)
	}
}

func TestHandleFailsWhenPodNeverRuns(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		d, queue, sandboxes := newDispatcher(t, nil)
		d.Options.HandoffWait = expiringHandoffWait
		sandboxes.podReady = false
		err := d.Handle(context.Background(), &cmaapi.WorkItem{ID: "work_1", SessionID: "session_1"})
		if err == nil || len(queue.ackedIDs()) != 0 {
			t.Errorf("err = %v, acked = %v", err, queue.ackedIDs())
		}
	})
}

func TestHandleResumesSuspendedSandbox(t *testing.T) {
	d, _, sandboxes := newDispatcher(t, nil)
	claim, _ := sandboxes.EnsureClaim(context.Background(), "session_1")
	sandboxes.sandboxes[claim.SandboxName] = &sandbox.Sandbox{Name: claim.SandboxName, Suspended: true}
	done := make(chan error, 1)
	go func() { done <- d.Handle(context.Background(), &cmaapi.WorkItem{ID: "work_1", SessionID: "session_1"}) }()
	collectWhenPending(t, d, types.UID("pod-"+claim.Name), done)
	if err := <-done; err != nil {
		t.Fatal(err)
	}
	if len(sandboxes.resumed) != 1 || sandboxes.resumed[0] != claim.SandboxName {
		t.Errorf("resumed = %v", sandboxes.resumed)
	}
	if len(sandboxes.stamped) != 2 || sandboxes.stamped[0] != claim.Name {
		t.Errorf("stamped = %v, want last-active stamped at the resume and at the ack", sandboxes.stamped)
	}
}

func TestHandleResumesTheSandboxForARedeliveryAfterItsWorkerExited(t *testing.T) {
	d, queue, sandboxes := newDispatcher(t, nil)
	d.workItems.Put("the-suspended-pod", workitem.SessionWorkItem{SessionID: "session_1", WorkID: "work_1"})
	if status, _ := d.workItems.Collect("the-suspended-pod", func(workitem.SessionWorkItem) error { return nil }); status != workitem.Collected {
		t.Fatal("setup: collect")
	}
	claim, _ := sandboxes.EnsureClaim(context.Background(), "session_1")
	claim.HandoffPodUID = "the-suspended-pod"
	claim.HandoffRestarts = ptr(int32(0))
	sandboxes.sandboxes[claim.SandboxName] = &sandbox.Sandbox{Name: claim.SandboxName, Suspended: true}
	done := make(chan error, 1)
	go func() { done <- d.Handle(context.Background(), &cmaapi.WorkItem{ID: "work_1", SessionID: "session_1"}) }()
	collectWhenPending(t, d, types.UID("pod-"+claim.Name), done)
	if err := <-done; err != nil {
		t.Fatal(err)
	}
	if acked := queue.ackedIDs(); len(acked) != 1 || acked[0] != "work_1" {
		t.Errorf("acked = %v", acked)
	}
	if len(sandboxes.resumed) != 1 {
		t.Errorf("resumed = %v; the item is a message to the idle session and must wake it", sandboxes.resumed)
	}
}

// A failed suspend keeps the watch of a pod whose worker has exited. A message
// sent then carries the id that pod's item was acked under. No pod holds the
// message, so acking it would drop the user's turn.
func TestAMessageToASessionAwaitingItsSuspendIsLeftForRedelivery(t *testing.T) {
	d, queue, sandboxes := newDispatcher(t, fakeSessions{"session_1": cmaapi.SessionAlive})
	claim, _ := sandboxes.EnsureClaim(context.Background(), "session_1")
	pod := sandboxes.runningPod(claim, corev1.ContainerStatus{State: corev1.ContainerState{Running: &corev1.ContainerStateRunning{}}})
	item := &cmaapi.WorkItem{ID: "session_1", SessionID: "session_1"}
	done := make(chan error, 1)
	go func() { done <- d.Handle(context.Background(), item) }()
	collectWhenPending(t, d, pod.UID, done)
	if err := <-done; err != nil {
		t.Fatal(err)
	}
	// The worker exits on idle and the suspend fails once.
	sandboxes.mu.Lock()
	pod.Status.ContainerStatuses[0].State = exitedAt(d.Now())
	sandboxes.suspendFailures = 1
	sandboxes.mu.Unlock()
	d.PodChanged(context.Background(), pod)
	if len(queue.ackedIDs()) != 1 || len(d.watchedHandoffs()) != 1 {
		t.Fatalf("setup: acked = %v, watched = %v", queue.ackedIDs(), d.watchedHandoffs())
	}

	err := d.Handle(context.Background(), item)
	if len(queue.ackedIDs()) != 1 {
		t.Errorf("acked = %v, err = %v; the new message was acked with nothing handed to any pod", queue.ackedIDs(), err)
	}
}

func TestATimedOutHandoffKeepsTheEarlierWatch(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		d, _, sandboxes := newDispatcher(t, nil)
		d.Options.HandoffWait = expiringHandoffWait
		claim, _ := sandboxes.EnsureClaim(context.Background(), "session_1")
		// The watched pod is gone and the Sandbox's next pod never collects.
		sandboxes.runningPod(claim, corev1.ContainerStatus{State: corev1.ContainerState{Running: &corev1.ContainerStateRunning{}}})
		earlier := handoff{claim: *claim, sb: sandboxes.sandboxes[claim.SandboxName], podUID: "the-old-pod", restarts: 1}
		d.watchHandoff(earlier)
		if err := d.Handle(context.Background(), &cmaapi.WorkItem{ID: "work_2", SessionID: "session_1"}); !errors.Is(err, ErrHandoffTimeout) {
			t.Fatalf("err = %v, want ErrHandoffTimeout", err)
		}
		if h, ok := d.watchedHandoffs()["session_1"]; !ok || h.podUID != earlier.podUID || h.restarts != earlier.restarts {
			t.Errorf("watch = %+v, %v; want the earlier handoff's watch back after the timeout", h, ok)
		}
	})
}

func TestHandleRecordsTheRestartCountOnTheClaim(t *testing.T) {
	d, _, sandboxes := newDispatcher(t, nil)
	podUID := types.UID("pod-" + sandbox.ClaimName("session_1"))
	done := make(chan error, 1)
	go func() { done <- d.Handle(context.Background(), &cmaapi.WorkItem{ID: "work_1", SessionID: "session_1"}) }()
	collectWhenPending(t, d, podUID, done)
	if err := <-done; err != nil {
		t.Fatal(err)
	}
	claim, _ := sandboxes.GetClaim(context.Background(), sandbox.ClaimName("session_1"))
	if claim.HandoffRestarts == nil || claim.HandoffPodUID == "" {
		t.Errorf("HandoffRestarts = %v, HandoffPodUID = %q; the handoff must record the worker's restart count and pod on the claim", claim.HandoffRestarts, claim.HandoffPodUID)
	}
}

func TestHandleRetriesTheRestartReadAfterCollection(t *testing.T) {
	d, queue, sandboxes := newDispatcher(t, nil)
	podUID := types.UID("pod-" + sandbox.ClaimName("session_1"))
	done := make(chan error, 1)
	go func() { done <- d.Handle(context.Background(), &cmaapi.WorkItem{ID: "work_1", SessionID: "session_1"}) }()
	if err := deliverWhenPending(t, d, podUID, done, func(workitem.SessionWorkItem) error {
		sandboxes.mu.Lock()
		sandboxes.pods["sb-"+sandbox.ClaimName("session_1")].Status.ContainerStatuses = []corev1.ContainerStatus{{Name: sandbox.WorkerContainer, RestartCount: 1}}
		sandboxes.podOfFailures = 2
		sandboxes.mu.Unlock()
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	if err := <-done; err != nil {
		t.Fatal(err)
	}
	if len(queue.ackedIDs()) != 1 {
		t.Fatal("the item must be acked after the handoff")
	}
	if h, watched := d.handoffOfPod(podUID); !watched || h.restarts != 1 {
		t.Errorf("watch = %+v, %v; want the restart count re-read after the transient failures", h, watched)
	}
}

func TestHandleLeavesAPodRunningItsEarlierWorkItem(t *testing.T) {
	d, queue, sandboxes := newDispatcher(t, nil)
	podUID := types.UID("pod-" + sandbox.ClaimName("session_1"))
	claim, _ := sandboxes.EnsureClaim(context.Background(), "session_1")
	sandboxes.runningPod(claim, corev1.ContainerStatus{State: corev1.ContainerState{Running: &corev1.ContainerStateRunning{}}})
	mustPut(t, d, podUID, workitem.SessionWorkItem{SessionID: "session_1", WorkID: "work_1"})
	if _, err := d.workItems.Collect(podUID, func(workitem.SessionWorkItem) error { return nil }); err != nil {
		t.Fatal(err)
	}
	if err := d.Handle(context.Background(), &cmaapi.WorkItem{ID: "work_2", SessionID: "session_1"}); !failedWithoutTimeout(err) {
		t.Fatalf("err = %v, want an error other than a timeout; new work for a pod still running its earlier work item must be left for redelivery", err)
	}
	if len(queue.ackedIDs()) != 0 {
		t.Error("the undelivered item must not be acked")
	}
	var got workitem.SessionWorkItem
	if status, _ := d.workItems.Collect(podUID, func(w workitem.SessionWorkItem) error { got = w; return nil }); status != workitem.CollectedAgain || got.WorkID != "work_1" {
		t.Errorf("Collect = %v %+v; want the earlier work item still served", status, got)
	}
}

// After a dispatcher restart the work-item store is empty, so only the watch
// Reconcile rebuilt can refuse the pod — including one whose worker exited
// after the adoption read and has not been parked yet.
func TestHandleRefusesAWatchedPodWhoseWorkerExited(t *testing.T) {
	d, queue, sandboxes := newDispatcher(t, fakeSessions{"session_1": cmaapi.SessionAlive})
	claim, _ := sandboxes.EnsureClaim(context.Background(), "session_1")
	pod := sandboxes.runningPod(claim, corev1.ContainerStatus{State: exitedAt(d.Now())})
	d.watchHandoff(handoff{claim: *claim, sb: sandboxes.sandboxes[claim.SandboxName], podUID: pod.UID, restarts: 0})
	// The pod's fetch sidecar outlives the worker and keeps polling for work.
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
	if err := d.Handle(context.Background(), &cmaapi.WorkItem{ID: "work_2", SessionID: "session_1"}); !failedWithoutTimeout(err) {
		t.Fatalf("err = %v, want an error other than a timeout; new work for a pod whose worker exited must be left for redelivery", err)
	}
	if acked := queue.ackedIDs(); len(acked) != 0 {
		t.Errorf("acked = %v; work the exited worker will never run must not be acked", acked)
	}
}

func TestAFailedHandoffLeavesTheClaimStampUntouched(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		d, _, sandboxes := newDispatcher(t, nil)
		d.Options.HandoffWait = expiringHandoffWait
		claim, _ := sandboxes.EnsureClaim(context.Background(), "session_1")
		if err := d.Handle(context.Background(), &cmaapi.WorkItem{ID: "work_1", SessionID: "session_1"}); !errors.Is(err, ErrHandoffTimeout) {
			t.Fatalf("err = %v, want ErrHandoffTimeout", err)
		}
		claim, _ = sandboxes.GetClaim(context.Background(), claim.Name)
		if claim.HandoffRestarts != nil || claim.HandoffPodUID != "" {
			t.Errorf("HandoffRestarts = %v, HandoffPodUID = %q; a handoff that did not take must stamp nothing", claim.HandoffRestarts, claim.HandoffPodUID)
		}
	})
}

func TestARestartedDispatcherLeavesAMidTurnPodItsWork(t *testing.T) {
	d, queue, sandboxes := newDispatcher(t, fakeSessions{"session_1": cmaapi.SessionAlive})
	claim, _ := sandboxes.EnsureClaim(context.Background(), "session_1")
	pod := sandboxes.runningPod(claim, corev1.ContainerStatus{State: corev1.ContainerState{Running: &corev1.ContainerStateRunning{}}})
	// As Reconcile rebuilds it after a restart: the watch exists, the
	// work-item store is empty.
	d.watchHandoff(handoff{claim: *claim, sb: sandboxes.sandboxes[claim.SandboxName], podUID: pod.UID, restarts: 0})
	// The pod's fetch re-requests its work item, as it may at any time.
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
	if err := d.Handle(context.Background(), &cmaapi.WorkItem{ID: "work_2", SessionID: "session_1"}); !failedWithoutTimeout(err) {
		t.Fatalf("err = %v, want an error other than a timeout; new work for a pod still running its pre-restart turn must be left for redelivery", err)
	}
	if acked := queue.ackedIDs(); len(acked) != 0 {
		t.Errorf("acked = %v; work the mid-turn pod's worker will never read must not be acked", acked)
	}
}
