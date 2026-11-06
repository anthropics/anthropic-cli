// Package dispatcher runs the serve loop: reconcile the namespace's claims
// against the API, poll the work queue, give each session a sandbox pod,
// hand the pod its work item, ack the work item, and suspend or delete the
// sandbox when its worker exits.
package dispatcher

import (
	"context"
	"log/slog"
	"math/rand/v2"
	"sync"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"

	"github.com/anthropics/anthropic-cli/sandboxes/dispatcher/internal/cmaapi"
	"github.com/anthropics/anthropic-cli/sandboxes/dispatcher/internal/sandbox"
	"github.com/anthropics/anthropic-cli/sandboxes/dispatcher/internal/workitem"
)

// WorkQueue is the environment's work queue.
type WorkQueue interface {
	Poll(ctx context.Context, workerID string, reclaimAfter time.Duration) (*cmaapi.WorkItem, error)
	Ack(ctx context.Context, workID string) error
}

// Sessions answers whether a session still exists.
type Sessions interface {
	SessionState(ctx context.Context, sessionID string) (cmaapi.SessionState, error)
}

// Sandboxes is the sandbox client surface the loop uses.
type Sandboxes interface {
	EnsureClaim(ctx context.Context, sessionID string) (*sandbox.SandboxClaim, error)
	GetClaim(ctx context.Context, name string) (*sandbox.SandboxClaim, error)
	ListClaims(ctx context.Context) ([]sandbox.SandboxClaim, error)
	DeleteClaim(ctx context.Context, name string) error
	StampLastActive(ctx context.Context, claimName string, at time.Time) error
	StampHandoff(ctx context.Context, claimName string, podUID types.UID, restarts int32) error
	SandboxOf(ctx context.Context, claim *sandbox.SandboxClaim) (*sandbox.Sandbox, error)
	PodOf(ctx context.Context, sb *sandbox.Sandbox) (*corev1.Pod, error)
	Suspend(ctx context.Context, sandboxName string) error
	Resume(ctx context.Context, sandboxName string, interval time.Duration) error
	WaitForRunningPod(ctx context.Context, claimName string, interval time.Duration) (*sandbox.Sandbox, *corev1.Pod, error)
}

// Options are the loop's settings.
type Options struct {
	EnvironmentID string
	BaseURL       string
	WorkerID      string
	// HandoffWait bounds provisioning plus the pod's work item fetch; the
	// queue's reclaim window, after which an unacked item is redelivered,
	// is twice it.
	HandoffWait     time.Duration
	PodPollInterval time.Duration
	SweepInterval   time.Duration
	// Retention deletes suspended claims idle for longer than this. Zero
	// keeps them forever.
	Retention time.Duration
	// MaxConcurrentHandoffs bounds the sessions being provisioned at once;
	// the loop stops polling while that many are in flight.
	MaxConcurrentHandoffs int
}

type Dispatcher struct {
	queue     WorkQueue
	sessions  Sessions
	sandboxes Sandboxes
	workItems *workitem.Store
	logger    *slog.Logger
	Options   Options
	Now       func() time.Time

	mu       sync.Mutex
	inflight map[string]struct{}
	// handoffs are the sessions handed to a pod, watched for the worker's
	// exit, keyed by session id.
	handoffs map[string]handoff
}

// handoff is what the worker's exit is judged against.
type handoff struct {
	claim sandbox.SandboxClaim
	sb    *sandbox.Sandbox
	// podUID is the pod the session was handed to; empty once that pod is
	// gone, until the Sandbox's next pod is resolved.
	podUID types.UID
	// restarts is the worker's restart count when the handoff was recorded;
	// a higher count later means it exited.
	restarts int32
	// restartsPreCollection means restarts dates from before the collection and
	// was not re-read after. A higher count may then be a restart the handoff
	// tolerated, so only an ended run confirms an exit.
	restartsPreCollection bool
	// acked means the queue took the ack of the work this handoff delivered.
	// The queue never redelivers acked work. The same id after that is a new
	// message, not a redelivery.
	acked bool
}

// New builds a Dispatcher with everything its event handlers read, so pod
// changes may be fed to it before Run starts.
func New(queue WorkQueue, sessions Sessions, sandboxes Sandboxes, workItems *workitem.Store, logger *slog.Logger, opts Options) *Dispatcher {
	if queue == nil || sessions == nil || sandboxes == nil || workItems == nil || logger == nil {
		panic("dispatcher.New: nil collaborator")
	}
	if opts.MaxConcurrentHandoffs <= 0 {
		panic("dispatcher.New: MaxConcurrentHandoffs must be positive")
	}
	return &Dispatcher{
		queue:     queue,
		sessions:  sessions,
		sandboxes: sandboxes,
		workItems: workItems,
		logger:    logger,
		Options:   opts,
		Now:       time.Now,
		inflight:  map[string]struct{}{},
		handoffs:  map[string]handoff{},
	}
}

// Run reconciles once, then polls until ctx ends, sweeping the claims on
// Options.SweepInterval.
func (d *Dispatcher) Run(ctx context.Context) error {
	if err := d.Reconcile(ctx); err != nil {
		return err
	}
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		d.sweepLoop(ctx)
	}()
	d.pollLoop(ctx)
	wg.Wait()
	return nil
}

// pollLoop takes a slot before each poll, so it asks for no work it cannot
// provision yet. The reclaim window is twice the handoff wait, so an item comes
// again only after its first delivery has given up.
func (d *Dispatcher) pollLoop(ctx context.Context) {
	var wg sync.WaitGroup
	defer wg.Wait()
	slots := make(chan struct{}, d.Options.MaxConcurrentHandoffs)
	failures := 0
	for ctx.Err() == nil {
		select {
		case slots <- struct{}{}:
		case <-ctx.Done():
			return
		}
		item, err := d.queue.Poll(ctx, d.Options.WorkerID, 2*d.Options.HandoffWait)
		if err != nil {
			<-slots
			if ctx.Err() != nil {
				return
			}
			failures++
			d.logger.Warn("polling the work queue", "error", err, "attempt", failures)
			sleep(ctx, backoff(failures))
			continue
		}
		failures = 0
		if item == nil {
			<-slots
			sleep(ctx, jitter(time.Second, 3*time.Second))
			continue
		}
		if !d.claimInflight(item.SessionID) {
			<-slots
			d.logger.Info("work item redelivered while its session is busy; ignoring", "work_id", item.ID, "session_id", item.SessionID)
			continue
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			defer func() { <-slots }()
			defer d.releaseInflight(item.SessionID)
			if err := d.Handle(ctx, item); err != nil && ctx.Err() == nil {
				d.logger.Error("handling work item", "work_id", item.ID, "session_id", item.SessionID, "error", err)
			}
		}()
	}
}

func (d *Dispatcher) claimInflight(sessionID string) bool {
	d.mu.Lock()
	defer d.mu.Unlock()
	if _, busy := d.inflight[sessionID]; busy {
		return false
	}
	d.inflight[sessionID] = struct{}{}
	return true
}

func (d *Dispatcher) releaseInflight(sessionID string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	delete(d.inflight, sessionID)
}

func backoff(failures int) time.Duration {
	const capDuration = 60 * time.Second
	if failures > 6 {
		return jitter(capDuration/2, capDuration)
	}
	base := time.Duration(1<<failures) * time.Second
	return jitter(base/2, base)
}

func jitter(low, high time.Duration) time.Duration {
	if high <= low {
		return low
	}
	return low + time.Duration(rand.Int64N(int64(high-low)))
}

func sleep(ctx context.Context, d time.Duration) {
	t := time.NewTimer(d)
	defer t.Stop()
	select {
	case <-ctx.Done():
	case <-t.C:
	}
}
