package workitem

import (
	"errors"
	"testing"
	"time"

	"k8s.io/apimachinery/pkg/types"
)

func mustPut(t *testing.T, s *Store, podUID types.UID, a SessionWorkItem) {
	t.Helper()
	if _, ok := s.Put(podUID, a); !ok {
		t.Fatalf("Put for pod %s refused", podUID)
	}
}

func collect(s *Store, podUID types.UID) (SessionWorkItem, Status, error) {
	var got SessionWorkItem
	status, err := s.Collect(podUID, func(a SessionWorkItem) error { got = a; return nil })
	return got, status, err
}

func TestStoreHandsOffOnTheFirstCollect(t *testing.T) {
	s := NewStore()
	if _, status, _ := collect(s, "pod-1"); status != NotPending {
		t.Fatalf("before Put: status %v", status)
	}
	taken, _ := s.Put("pod-1", SessionWorkItem{SessionID: "session_1", WorkID: "work_1"})
	select {
	case <-taken:
		t.Fatal("taken closed before Collect")
	default:
	}
	got, status, err := collect(s, "pod-1")
	if err != nil || status != Collected || got.SessionID != "session_1" {
		t.Fatalf("Collect = %+v, %v, %v", got, status, err)
	}
	select {
	case <-taken:
	default:
		t.Error("taken not closed after Collect")
	}
	if workID, ok := s.HandedWork("session_1"); !ok || workID != "work_1" {
		t.Errorf("HandedWork = %q, %v", workID, ok)
	}
	if s.Drop("pod-1") {
		t.Error("Drop after a collection must report nothing pending")
	}
	if got, status, _ := collect(s, "pod-1"); status != CollectedAgain || got.WorkID != "work_1" {
		t.Errorf("Collect after Drop of a collected work item = %+v, %v; want it still open", got, status)
	}
}

func TestSamePodCollectsAgainUntilClose(t *testing.T) {
	s := NewStore()
	want := SessionWorkItem{SessionID: "session_1", WorkID: "work_1", Secret: "secret_1"}
	mustPut(t, s, "pod-1", want)
	if _, status, _ := collect(s, "pod-1"); status != Collected {
		t.Fatalf("first Collect: status %v", status)
	}
	got, status, err := collect(s, "pod-1")
	if err != nil || status != CollectedAgain || got != want {
		t.Errorf("second Collect = %+v, %v, %v; want the same work item again", got, status, err)
	}
	if got, status, _ := collect(s, "pod-2"); status != NotPending || got != (SessionWorkItem{}) {
		t.Errorf("another pod's Collect = %+v, %v; want nothing", got, status)
	}
	s.Close("pod-1")
	if _, status, _ := collect(s, "pod-1"); status != AlreadyCollected {
		t.Errorf("Collect after Close: status %v, want AlreadyCollected", status)
	}
}

func TestCloseLeavesAnUncollectedWorkItemForDrop(t *testing.T) {
	s := NewStore()
	taken, _ := s.Put("pod-1", SessionWorkItem{SessionID: "session_1"})
	s.Close("pod-1")
	if _, status, _ := collect(s, "pod-1"); status != Collected {
		t.Fatalf("Collect after Close of an uncollected work item: status %v, want Collected", status)
	}
	<-taken
}

func TestForgetPodLeavesAnUncollectedWorkItemForDrop(t *testing.T) {
	s := NewStore()
	mustPut(t, s, "pod-1", SessionWorkItem{SessionID: "session_1"})
	s.ForgetPod("pod-1")
	if !s.Drop("pod-1") {
		t.Error("Drop after ForgetPod of an uncollected work item reported nothing pending; the handoff would ack work the pod never received")
	}
}

func TestForgetSessionClosesItsWorkItems(t *testing.T) {
	s := NewStore()
	mustPut(t, s, "pod-1", SessionWorkItem{SessionID: "session_1"})
	mustPut(t, s, "pod-2", SessionWorkItem{SessionID: "session_2"})
	_, _, _ = collect(s, "pod-1")
	_, _, _ = collect(s, "pod-2")
	s.ForgetSession("session_1")
	if _, status, _ := collect(s, "pod-1"); status != AlreadyCollected {
		t.Errorf("pod-1 after ForgetSession: status %v", status)
	}
	if _, status, _ := collect(s, "pod-2"); status != CollectedAgain {
		t.Errorf("pod-2 of another session: status %v", status)
	}
}

func TestStoreKeepsTheWorkItemWhenDeliveryFails(t *testing.T) {
	s := NewStore()
	taken, _ := s.Put("pod-1", SessionWorkItem{SessionID: "session_1"})
	status, err := s.Collect("pod-1", func(SessionWorkItem) error { return errors.New("write failed") })
	if err == nil || status != Collected {
		t.Fatalf("Collect = %v, %v", status, err)
	}
	select {
	case <-taken:
		t.Fatal("taken closed after a failed delivery")
	default:
	}
	if _, status, _ := collect(s, "pod-1"); status != Collected {
		t.Errorf("retry status %v", status)
	}
}

func TestStoreDrop(t *testing.T) {
	s := NewStore()
	mustPut(t, s, "pod-1", SessionWorkItem{SessionID: "session_1"})
	if !s.Drop("pod-1") {
		t.Error("Drop of a pending work item must report it")
	}
	if _, status, _ := collect(s, "pod-1"); status != NotPending {
		t.Errorf("Collect after Drop: status %v", status)
	}
}

func TestDeliveryRunsOutsideTheLockAndDropWaitsForIt(t *testing.T) {
	s := NewStore()
	mustPut(t, s, "pod-1", SessionWorkItem{WorkID: "work_1"})
	mustPut(t, s, "pod-2", SessionWorkItem{WorkID: "work_2"})
	started, release := make(chan struct{}), make(chan struct{})
	go func() {
		_, _ = s.Collect("pod-1", func(SessionWorkItem) error { close(started); <-release; return nil })
	}()
	<-started
	// Another pod's collection must not wait on pod-1's delivery.
	if _, status, _ := collect(s, "pod-2"); status != Collected {
		t.Errorf("pod-2 status %v", status)
	}
	dropped := make(chan bool, 1)
	go func() { dropped <- s.Drop("pod-1") }()
	select {
	case got := <-dropped:
		t.Fatalf("Drop returned %v while the delivery was in flight", got)
	case <-time.After(20 * time.Millisecond):
	}
	close(release)
	if got := <-dropped; got {
		t.Error("Drop after the delivery succeeded must report nothing pending")
	}
}

func TestHandoffRecordIsKeptBySession(t *testing.T) {
	s := NewStore()
	mustPut(t, s, "pod-1", SessionWorkItem{SessionID: "session_1", WorkID: "work_1"})
	_, _, _ = collect(s, "pod-1")
	s.ForgetPod("pod-1")
	mustPut(t, s, "pod-2", SessionWorkItem{SessionID: "session_1", WorkID: "work_2"})
	if workID, ok := s.HandedWork("session_1"); !ok || workID != "work_1" {
		t.Errorf("HandedWork after the pod was replaced = %q, %v; want the first collection still recorded", workID, ok)
	}
	if !s.Drop("pod-2") {
		t.Error("the second work item was pending and must be reported so")
	}
	s.ForgetSession("session_1")
	if _, ok := s.HandedWork("session_1"); ok {
		t.Error("ForgetSession must drop the record")
	}
}

func TestPutLeavesAnOpenDeliveredWorkItem(t *testing.T) {
	s := NewStore()
	mustPut(t, s, "pod-1", SessionWorkItem{SessionID: "session_1", WorkID: "work_1"})
	if _, _, err := collect(s, "pod-1"); err != nil {
		t.Fatal(err)
	}
	if _, ok := s.Put("pod-1", SessionWorkItem{SessionID: "session_1", WorkID: "work_2"}); ok {
		t.Fatal("Put over an open delivered work item must refuse; the worker is still running the first one")
	}
	if got, status, _ := collect(s, "pod-1"); status != CollectedAgain || got.WorkID != "work_1" {
		t.Errorf("Collect = %+v, %v; want the first work item still served", got, status)
	}
	s.Close("pod-1")
	if _, ok := s.Put("pod-1", SessionWorkItem{SessionID: "session_1", WorkID: "work_2"}); !ok {
		t.Error("Put after Close must open the new work item")
	}
}
