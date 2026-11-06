package dispatcher

import "k8s.io/apimachinery/pkg/types"

// restoreHandoff puts back the watch a session had before a handoff that
// timed out, so a pod handed earlier work stays watched for its worker's
// exit.
func (d *Dispatcher) restoreHandoff(sessionID string, earlier *handoff) {
	if earlier == nil {
		d.forgetHandoff(sessionID)
		return
	}
	d.watchHandoff(*earlier)
}

func (d *Dispatcher) watchHandoff(h handoff) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.handoffs[h.claim.SessionID] = h
}

func (d *Dispatcher) forgetHandoff(sessionID string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	delete(d.handoffs, sessionID)
}

// handoffOfSession returns the session's watch, or nil when it has none.
func (d *Dispatcher) handoffOfSession(sessionID string) *handoff {
	d.mu.Lock()
	defer d.mu.Unlock()
	h, ok := d.handoffs[sessionID]
	if !ok {
		return nil
	}
	return &h
}

func (d *Dispatcher) handoffOfPod(podUID types.UID) (handoff, bool) {
	d.mu.Lock()
	defer d.mu.Unlock()
	for _, h := range d.handoffs {
		if h.podUID == podUID {
			return h, true
		}
	}
	return handoff{}, false
}

func (d *Dispatcher) watchedHandoffs() map[string]handoff {
	d.mu.Lock()
	defer d.mu.Unlock()
	out := make(map[string]handoff, len(d.handoffs))
	for sessionID, h := range d.handoffs {
		out[sessionID] = h
	}
	return out
}
