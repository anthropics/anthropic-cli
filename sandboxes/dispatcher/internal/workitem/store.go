// Package workitem holds the work item each sandbox pod collects, keyed by the
// pod's UID. WorkItems live only in memory.
package workitem

import (
	"sync"

	"k8s.io/apimachinery/pkg/types"
)

// SessionWorkItem is what a sandbox pod receives for its session. Secret is
// the queue's secret payload, verbatim; the worker decodes it.
type SessionWorkItem struct {
	SessionID     string `json:"session_id"`
	WorkID        string `json:"work_id"`
	Secret        string `json:"secret"`
	BaseURL       string `json:"base_url"`
	EnvironmentID string `json:"environment_id"`
}

// Status is the outcome of a pod's Collect call.
type Status int

const (
	// Collected means the work item was delivered for the first time and is
	// now handed off.
	Collected Status = iota
	// CollectedAgain means the pod collected the work item before and was
	// given it again, because it is still open.
	CollectedAgain
	// NotPending means nothing was ever put for the pod, or it was dropped
	// before the pod collected it.
	NotPending
	// AlreadyCollected means the pod collected its work item and the
	// work item has since been closed.
	AlreadyCollected
)

type openWorkItem struct {
	item      SessionWorkItem
	delivered bool
	taken     chan struct{}
	// delivering is non-nil while a delivery is in flight and closes when
	// it ends; every other operation on the pod waits for it.
	delivering chan struct{}
}

// Store is safe for concurrent use. A delivery runs outside the lock, so
// one pod's slow connection never holds up the others.
type Store struct {
	mu        sync.Mutex
	openByPod map[types.UID]*openWorkItem
	// handed remembers the work id each session's pod last collected, by
	// session id, so it outlives the pod.
	handed map[string]string
	// collectedBy marks the pods that collected a work item, so a request
	// from one after its work item closed is answered 410.
	collectedBy map[types.UID]bool
}

func NewStore() *Store {
	return &Store{openByPod: map[types.UID]*openWorkItem{}, handed: map[string]string{}, collectedBy: map[types.UID]bool{}}
}

// Put opens the work item for podUID and returns a channel that closes once
// the pod has first collected it, replacing an earlier work item only when
// the pod never collected it or Close ended it: an open delivered work item
// means the worker is still running that one, and handing the pod new work
// would be acked without ever running. The record of what the pod collected
// before stays, so a redelivery of that work is still recognized.
func (s *Store) Put(podUID types.UID, a SessionWorkItem) (<-chan struct{}, bool) {
	s.lockSettled(podUID)
	defer s.mu.Unlock()
	if o, ok := s.openByPod[podUID]; ok && o.delivered {
		return nil, false
	}
	o := &openWorkItem{item: a, taken: make(chan struct{})}
	s.openByPod[podUID] = o
	return o.taken, true
}

// Collect delivers the pod's open work item through deliver. The first
// delivery that succeeds hands the work item off; the work item stays open
// for the same pod until Close. A failed first delivery keeps it pending for
// the pod's retry.
func (s *Store) Collect(podUID types.UID, deliver func(SessionWorkItem) error) (Status, error) {
	s.lockSettled(podUID)
	o, ok := s.openByPod[podUID]
	if !ok {
		defer s.mu.Unlock()
		if s.collectedBy[podUID] {
			return AlreadyCollected, nil
		}
		return NotPending, nil
	}
	o.delivering = make(chan struct{})
	s.mu.Unlock()

	err := deliver(o.item)

	s.mu.Lock()
	defer s.mu.Unlock()
	close(o.delivering)
	o.delivering = nil
	if o.delivered {
		return CollectedAgain, err
	}
	if err != nil {
		return Collected, err
	}
	o.delivered = true
	s.handed[o.item.SessionID] = o.item.WorkID
	s.collectedBy[podUID] = true
	close(o.taken)
	return Collected, nil
}

// Close ends a work item the pod has collected, so its later requests are
// answered 410. A work item the pod has not collected stays open: only
// Drop, run by the handoff that waits for it, ends that one.
func (s *Store) Close(podUID types.UID) {
	s.lockSettled(podUID)
	defer s.mu.Unlock()
	s.deleteCollected(podUID)
}

// deleteCollected removes the pod's work item if the pod has collected it.
// The caller holds the lock.
func (s *Store) deleteCollected(podUID types.UID) {
	if o, ok := s.openByPod[podUID]; ok && o.delivered {
		delete(s.openByPod, podUID)
	}
}

// Drop closes a work item the pod has not collected and reports whether
// there was one. false means the pod collected it, since nothing else
// removes an uncollected work item, and the work item then stays open.
func (s *Store) Drop(podUID types.UID) (wasPending bool) {
	s.lockSettled(podUID)
	defer s.mu.Unlock()
	o, ok := s.openByPod[podUID]
	if !ok || o.delivered {
		return false
	}
	delete(s.openByPod, podUID)
	return true
}

// HandedWork returns the work id last collected for the session, if any.
func (s *Store) HandedWork(sessionID string) (string, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	workID, ok := s.handed[sessionID]
	return workID, ok
}

// ForgetSession closes the session's collected workItems and drops the
// record of what its pods collected, once the session's sandbox is gone.
func (s *Store) ForgetSession(sessionID string) {
	for _, podUID := range s.podsOfSession(sessionID) {
		s.Close(podUID)
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.handed, sessionID)
}

func (s *Store) podsOfSession(sessionID string) []types.UID {
	s.mu.Lock()
	defer s.mu.Unlock()
	var pods []types.UID
	for podUID, o := range s.openByPod {
		if o.item.SessionID == sessionID {
			pods = append(pods, podUID)
		}
	}
	return pods
}

// ForgetPod closes the collected work item of a pod that no longer exists
// and drops the record that it collected one. A work item the pod never
// collected stays for Drop, as with Close.
func (s *Store) ForgetPod(podUID types.UID) {
	s.lockSettled(podUID)
	defer s.mu.Unlock()
	s.deleteCollected(podUID)
	delete(s.collectedBy, podUID)
}

// lockSettled takes the lock once no delivery to podUID is in flight, so
// the caller sees the delivery's outcome.
func (s *Store) lockSettled(podUID types.UID) {
	for {
		s.mu.Lock()
		o, ok := s.openByPod[podUID]
		if !ok || o.delivering == nil {
			return
		}
		delivering := o.delivering
		s.mu.Unlock()
		<-delivering
	}
}
