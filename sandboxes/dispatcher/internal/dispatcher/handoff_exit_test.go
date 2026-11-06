package dispatcher

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"testing"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"

	"github.com/anthropics/anthropic-cli/sandboxes/dispatcher/internal/cmaapi"
	"github.com/anthropics/anthropic-cli/sandboxes/dispatcher/internal/workitem"
)

func TestHandleCountsRestartsFromTheHandoffOn(t *testing.T) {
	d, _, sandboxes := newDispatcher(t, fakeSessions{"session_1": cmaapi.SessionAlive})
	claim, _ := sandboxes.EnsureClaim(context.Background(), "session_1")
	pod := sandboxes.runningPod(claim, corev1.ContainerStatus{State: corev1.ContainerState{Running: &corev1.ContainerStateRunning{}}})
	done := make(chan error, 1)
	go func() { done <- d.Handle(context.Background(), &cmaapi.WorkItem{ID: "work_1", SessionID: "session_1"}) }()
	if err := deliverWhenPending(t, d, pod.UID, done, func(workitem.SessionWorkItem) error {
		// The worker restarted while it waited for its token.
		sandboxes.mu.Lock()
		pod.Status.ContainerStatuses[0].RestartCount = 1
		sandboxes.mu.Unlock()
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	if err := <-done; err != nil {
		t.Fatal(err)
	}
	d.PodChanged(context.Background(), pod)
	if len(sandboxes.suspended) != 0 {
		t.Errorf("a restart before the collection must not count as an exit: suspended = %v", sandboxes.suspended)
	}
}

func TestARestartBeforeTheCollectionKeepsTheAssignment(t *testing.T) {
	d, queue, sandboxes := newDispatcher(t, fakeSessions{"session_1": cmaapi.SessionAlive})
	claim, _ := sandboxes.EnsureClaim(context.Background(), "session_1")
	pod := sandboxes.runningPod(claim, corev1.ContainerStatus{State: corev1.ContainerState{Running: &corev1.ContainerStateRunning{}}})
	done := make(chan error, 1)
	go func() { done <- d.Handle(context.Background(), &cmaapi.WorkItem{ID: "work_1", SessionID: "session_1"}) }()
	notYet := errors.New("not yet")
	if err := deliverWhenPending(t, d, pod.UID, done, func(workitem.SessionWorkItem) error { return notYet }); !errors.Is(err, notYet) {
		t.Fatalf("collecting with a delivery that fails returned err = %v, want the delivery's error %v", err, notYet)
	}

	restarted := pod.DeepCopy()
	restarted.Status.ContainerStatuses[0].RestartCount = 1
	d.PodChanged(context.Background(), restarted)
	if status, _ := d.workItems.Collect(pod.UID, func(workitem.SessionWorkItem) error { return nil }); status != workitem.Collected {
		t.Fatalf("collection after a restart during the handoff: status %v, want Collected", status)
	}
	if err := <-done; err != nil {
		t.Fatal(err)
	}
	if acked := queue.ackedIDs(); len(acked) != 1 {
		t.Errorf("acked = %v, want the collected item acked once", acked)
	}
}

// An acked item never redelivers, so the ack must wait for the handoff stamp
// to land on the claim: a dispatcher restarted before the stamp would judge
// the token-holding pod stale and suspend it mid-turn, losing the turn.
func TestHandleAcksOnlyOnceTheStampLands(t *testing.T) {
	d, queue, sandboxes := newDispatcher(t, fakeSessions{"session_1": cmaapi.SessionAlive})
	claim, _ := sandboxes.EnsureClaim(context.Background(), "session_1")
	pod := sandboxes.runningPod(claim, corev1.ContainerStatus{State: corev1.ContainerState{Running: &corev1.ContainerStateRunning{}}})
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
	// The claim store stays unreachable through this delivery and the first
	// redelivery.
	sandboxes.mu.Lock()
	sandboxes.stampFailures = 10
	sandboxes.mu.Unlock()
	item := &cmaapi.WorkItem{ID: "work_1", SessionID: "session_1"}
	for range 2 {
		if err := d.Handle(context.Background(), item); !failedWithoutTimeout(err) {
			t.Fatalf("err = %v, want an error other than a timeout; an item whose handoff stamp did not land must be left for redelivery", err)
		}
	}
	if acked := queue.ackedIDs(); len(acked) != 0 {
		t.Fatalf("acked = %v, want no ack while the claim carries no handoff stamp", acked)
	}
	// The outage ends: the next redelivery stamps the claim, then acks.
	if err := d.Handle(context.Background(), item); err != nil {
		t.Fatal(err)
	}
	if acked := queue.ackedIDs(); len(acked) != 1 {
		t.Fatalf("acked = %v, want the redelivered item acked once the stamp landed", acked)
	}
	stamped, _ := sandboxes.GetClaim(context.Background(), claim.Name)
	if stamped.HandoffPodUID != pod.UID || stamped.HandoffRestarts == nil {
		t.Fatalf("HandoffPodUID = %q, HandoffRestarts = %v; want the handoff on the claim before the ack", stamped.HandoffPodUID, stamped.HandoffRestarts)
	}
	// A dispatcher restarted now adopts the mid-turn pod instead of suspending it.
	restarted := New(queue, fakeSessions{"session_1": cmaapi.SessionAlive}, sandboxes, workitem.NewStore(), slog.New(slog.NewTextHandler(io.Discard, nil)), d.Options)
	restarted.Now = d.Now
	if err := restarted.Reconcile(context.Background()); err != nil {
		t.Fatal(err)
	}
	if len(sandboxes.suspended) != 0 {
		t.Errorf("suspended = %v; the stamped claim's mid-turn pod must be adopted, not suspended", sandboxes.suspended)
	}
	if h, ok := restarted.watchedHandoffs()["session_1"]; !ok || h.podUID != pod.UID {
		t.Errorf("watch = %+v, %v; want the mid-turn pod adopted after the restart", h, ok)
	}
}

// A redelivery arrives minutes after the collection, when the worker may
// already have run the turn and exited into a kubelet restart. Its stamp
// must carry the baseline recorded at the collection: re-reading the count
// would absorb that exit and the watch would never end the session.
func TestARedeliveryStampsTheRecordedBaseline(t *testing.T) {
	d, queue, sandboxes := newDispatcher(t, fakeSessions{"session_1": cmaapi.SessionAlive})
	claim, _ := sandboxes.EnsureClaim(context.Background(), "session_1")
	pod := sandboxes.runningPod(claim, corev1.ContainerStatus{State: corev1.ContainerState{Running: &corev1.ContainerStateRunning{}}})
	sandboxes.mu.Lock()
	sandboxes.stampFailures = 5
	sandboxes.mu.Unlock()
	item := &cmaapi.WorkItem{ID: "work_1", SessionID: "session_1"}
	done := make(chan error, 1)
	go func() { done <- d.Handle(context.Background(), item) }()
	collectWhenPending(t, d, pod.UID, done)
	if err := <-done; !failedWithoutTimeout(err) {
		t.Fatalf("err = %v, want an error other than a timeout; the item must be left for redelivery while the stamp cannot land", err)
	}
	// The worker runs the turn, exits, and the kubelet restarts it.
	sandboxes.mu.Lock()
	pod.Status.ContainerStatuses[0].RestartCount = 1
	pod.Status.ContainerStatuses[0].LastTerminationState = exitedAt(d.Now())
	sandboxes.mu.Unlock()
	if err := d.Handle(context.Background(), item); err != nil {
		t.Fatal(err)
	}
	if acked := queue.ackedIDs(); len(acked) != 1 {
		t.Fatalf("acked = %v", acked)
	}
	stamped, _ := sandboxes.GetClaim(context.Background(), claim.Name)
	if stamped.HandoffRestarts == nil {
		t.Fatal("HandoffRestarts = nil, want the handoff stamped before the ack")
	}
	if *stamped.HandoffRestarts != 0 {
		t.Fatalf("HandoffRestarts = %d, want the baseline recorded at the collection, not the post-exit count", *stamped.HandoffRestarts)
	}
	// The exit is still detected: the next pod event ends the session.
	d.PodChanged(context.Background(), pod)
	if len(sandboxes.suspended) != 1 {
		t.Errorf("suspended = %v; the worker's completed exit must still suspend the sandbox", sandboxes.suspended)
	}
	if status, _ := d.workItems.Collect(pod.UID, func(workitem.SessionWorkItem) error { return nil }); status != workitem.AlreadyCollected {
		t.Errorf("Collect = %v; the restarted worker must not be handed the acked turn again", status)
	}
}

// A worker restart during the token wait plus a failed post-collection
// re-read leaves the watch with a baseline from before the restart. That
// baseline must not read as the worker's exit: the pod is mid-turn, and a
// false suspend would drop the watch and let the redelivery ack the turn
// without it ever running.
func TestAStaleBaselineNeverSuspendsTheMidTurnPod(t *testing.T) {
	d, queue, sandboxes := newDispatcher(t, fakeSessions{"session_1": cmaapi.SessionAlive})
	claim, _ := sandboxes.EnsureClaim(context.Background(), "session_1")
	pod := sandboxes.runningPod(claim, corev1.ContainerStatus{State: corev1.ContainerState{Running: &corev1.ContainerStateRunning{}}})
	item := &cmaapi.WorkItem{ID: "work_1", SessionID: "session_1"}
	done := make(chan error, 1)
	go func() { done <- d.Handle(context.Background(), item) }()
	if err := deliverWhenPending(t, d, pod.UID, done, func(workitem.SessionWorkItem) error {
		// The worker restarts while it waits for its token, and the
		// re-read after the collection cannot reach the apiserver.
		sandboxes.mu.Lock()
		pod.Status.ContainerStatuses[0].RestartCount = 1
		pod.Status.ContainerStatuses[0].LastTerminationState = exitedAt(d.Now())
		sandboxes.podOfFailures = 5
		sandboxes.mu.Unlock()
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	if err := <-done; !failedWithoutTimeout(err) {
		t.Fatalf("err = %v, want an error other than a timeout; the item must be left for redelivery while the handoff is not recorded", err)
	}
	// The restart event arrives; the worker is mid-turn and must be left alone.
	d.PodChanged(context.Background(), pod)
	if len(sandboxes.suspended) != 0 {
		t.Fatalf("suspended = %v; a baseline from before the collection must not read a running worker as exited", sandboxes.suspended)
	}
	// The redelivered item re-reads the healed apiserver, stamps, and acks.
	if err := d.Handle(context.Background(), item); err != nil {
		t.Fatal(err)
	}
	if acked := queue.ackedIDs(); len(acked) != 1 {
		t.Fatalf("acked = %v, want the redelivered item acked while its turn runs", acked)
	}
	stamped, _ := sandboxes.GetClaim(context.Background(), claim.Name)
	if stamped.HandoffRestarts == nil || *stamped.HandoffRestarts != 1 {
		t.Fatalf("HandoffRestarts = %v, want the re-read count of the run that holds the token", stamped.HandoffRestarts)
	}
	// The worker really exits; the exit is detected against the healed baseline.
	sandboxes.mu.Lock()
	pod.Status.ContainerStatuses[0].State = exitedAt(d.Now())
	sandboxes.mu.Unlock()
	d.PodChanged(context.Background(), pod)
	if len(sandboxes.suspended) != 1 {
		t.Errorf("suspended = %v, want the sandbox suspended after the real exit", sandboxes.suspended)
	}
}

// The queue's work id is its session id, so a message to a session whose
// worker exited on idle delivers the id the session already ran. That
// delivery must wake the session, not be acked as already done.
func TestAMessageToAnIdleSessionWakesIt(t *testing.T) {
	d, queue, sandboxes := newDispatcher(t, fakeSessions{"session_1": cmaapi.SessionAlive})
	claim, _ := sandboxes.EnsureClaim(context.Background(), "session_1")
	pod := sandboxes.runningPod(claim, corev1.ContainerStatus{State: corev1.ContainerState{Running: &corev1.ContainerStateRunning{}}})
	pod.UID = "the-old-pod"
	claim.HandoffPodUID, claim.HandoffRestarts = pod.UID, ptr(int32(0))
	d.watchHandoff(handoff{claim: *claim, sb: sandboxes.sandboxes[claim.SandboxName], podUID: pod.UID})
	mustPut(t, d, pod.UID, workitem.SessionWorkItem{SessionID: "session_1", WorkID: "session_1"})
	if status, _ := d.workItems.Collect(pod.UID, func(workitem.SessionWorkItem) error { return nil }); status != workitem.Collected {
		t.Fatal("setup: collect")
	}
	// The worker exits on idle and the sandbox is suspended.
	sandboxes.mu.Lock()
	pod.Status.ContainerStatuses[0].State = exitedAt(d.Now())
	sandboxes.mu.Unlock()
	d.PodChanged(context.Background(), pod)
	if len(sandboxes.suspended) != 1 || len(d.watchedHandoffs()) != 0 {
		t.Fatalf("setup: suspended = %v, watched = %v", sandboxes.suspended, d.watchedHandoffs())
	}

	// A message to the session puts its item back on the queue.
	done := make(chan error, 1)
	go func() {
		done <- d.Handle(context.Background(), &cmaapi.WorkItem{ID: "session_1", SessionID: "session_1"})
	}()
	collectWhenPending(t, d, types.UID("pod-"+claim.Name), done)
	if err := <-done; err != nil {
		t.Fatal(err)
	}
	if len(sandboxes.resumed) != 1 || len(queue.ackedIDs()) != 1 {
		t.Errorf("resumed = %v, acked = %v; want the sandbox resumed and the item acked once its pod holds it", sandboxes.resumed, queue.ackedIDs())
	}
}
