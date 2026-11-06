package dispatcher

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"sync"
	"testing"
	"time"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	"github.com/anthropics/anthropic-cli/sandboxes/dispatcher/internal/cmaapi"
	"github.com/anthropics/anthropic-cli/sandboxes/dispatcher/internal/sandbox"
	"github.com/anthropics/anthropic-cli/sandboxes/dispatcher/internal/workitem"
)

func ptr[T any](v T) *T { return &v }

type fakeQueue struct {
	mu    sync.Mutex
	acked []string
}

func (q *fakeQueue) Poll(context.Context, string, time.Duration) (*cmaapi.WorkItem, error) {
	return nil, nil
}

func (q *fakeQueue) Ack(_ context.Context, workID string) error {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.acked = append(q.acked, workID)
	return nil
}

func (q *fakeQueue) ackedIDs() []string {
	q.mu.Lock()
	defer q.mu.Unlock()
	return append([]string(nil), q.acked...)
}

type fakeSessions map[string]any // cmaapi.SessionState or error

func (f fakeSessions) SessionState(_ context.Context, id string) (cmaapi.SessionState, error) {
	switch v := f[id].(type) {
	case cmaapi.SessionState:
		return v, nil
	case error:
		return cmaapi.SessionAlive, v
	}
	return cmaapi.SessionGone, nil
}

type fakeSandboxes struct {
	mu        sync.Mutex
	claims    map[string]*sandbox.SandboxClaim
	sandboxes map[string]*sandbox.Sandbox
	// pods holds each sandbox's pod by sandbox name; a sandbox without one
	// has no pod.
	pods      map[string]*corev1.Pod
	podReady  bool
	deleted   []string
	resumed   []string
	suspended []string
	stamped   []string
	// onSandboxOf runs under the lock on each SandboxOf: a hook for changes
	// that land between a sweep's listing and its per-claim reads.
	onSandboxOf func(*sandbox.SandboxClaim)
	// podOfFailures makes that many PodOf calls fail before they succeed;
	// sandboxOfFailures, suspendFailures and stampFailures do the same for
	// their calls.
	podOfFailures     int
	sandboxOfFailures int
	suspendFailures   int
	stampFailures     int
}

func (f *fakeSandboxes) EnsureClaim(_ context.Context, sessionID string) (*sandbox.SandboxClaim, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	name := sandbox.ClaimName(sessionID)
	if c, ok := f.claims[name]; ok {
		return c, nil
	}
	c := &sandbox.SandboxClaim{Name: name, SessionID: sessionID, SandboxName: "sb-" + name}
	f.claims[name] = c
	return c, nil
}

func (f *fakeSandboxes) GetClaim(_ context.Context, name string) (*sandbox.SandboxClaim, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	c, ok := f.claims[name]
	if !ok {
		return nil, apierrors.NewNotFound(corev1.Resource("sandboxclaims"), name)
	}
	return c, nil
}

func (f *fakeSandboxes) ListClaims(context.Context) ([]sandbox.SandboxClaim, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	var out []sandbox.SandboxClaim
	for _, c := range f.claims {
		out = append(out, *c)
	}
	return out, nil
}

func (f *fakeSandboxes) DeleteClaim(_ context.Context, name string) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	delete(f.claims, name)
	f.deleted = append(f.deleted, name)
	return nil
}

func (f *fakeSandboxes) StampLastActive(_ context.Context, name string, at time.Time) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if c, ok := f.claims[name]; ok {
		c.LastActive = at
	}
	f.stamped = append(f.stamped, name)
	return nil
}

func (f *fakeSandboxes) StampHandoff(ctx context.Context, name string, podUID types.UID, restarts int32) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.stampFailures > 0 {
		f.stampFailures--
		return errors.New("transient")
	}
	if c, ok := f.claims[name]; ok {
		c.HandoffRestarts = &restarts
		c.HandoffPodUID = podUID
	}
	return nil
}

func (f *fakeSandboxes) SandboxOf(_ context.Context, claim *sandbox.SandboxClaim) (*sandbox.Sandbox, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.onSandboxOf != nil {
		f.onSandboxOf(f.claims[claim.Name])
	}
	if f.sandboxOfFailures > 0 {
		f.sandboxOfFailures--
		return nil, errors.New("transient")
	}
	sb, ok := f.sandboxes[claim.SandboxName]
	if !ok {
		return nil, sandbox.ErrNotBound
	}
	return sb, nil
}

func (f *fakeSandboxes) PodOf(ctx context.Context, sb *sandbox.Sandbox) (*corev1.Pod, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.podOfFailures > 0 {
		f.podOfFailures--
		return nil, errors.New("transient")
	}
	pod, ok := f.pods[sb.Name]
	if !ok {
		return nil, apierrors.NewNotFound(corev1.Resource("pods"), sb.PodName)
	}
	return pod, nil
}

func (f *fakeSandboxes) Suspend(_ context.Context, name string) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.suspendFailures > 0 {
		f.suspendFailures--
		return errors.New("transient")
	}
	f.sandboxes[name].Suspended = true
	delete(f.pods, name)
	f.suspended = append(f.suspended, name)
	return nil
}

func (f *fakeSandboxes) Resume(_ context.Context, name string, _ time.Duration) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.sandboxes[name].Suspended = false
	f.resumed = append(f.resumed, name)
	return nil
}

func (f *fakeSandboxes) WaitForRunningPod(ctx context.Context, claimName string, _ time.Duration) (*sandbox.Sandbox, *corev1.Pod, error) {
	f.mu.Lock()
	ready := f.podReady
	f.mu.Unlock()
	if !ready {
		<-ctx.Done()
		return nil, nil, ctx.Err()
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	name := "sb-" + claimName
	sb, ok := f.sandboxes[name]
	if !ok {
		sb = &sandbox.Sandbox{Name: name, PodName: name}
	}
	pod, ok := f.pods[name]
	if !ok {
		pod = &corev1.Pod{ObjectMeta: metav1.ObjectMeta{Name: name, UID: types.UID("pod-" + claimName)}}
		f.pods[name] = pod
	}
	return sb, pod, nil
}

// runningPod gives the fake a pod for claim whose worker reports status.
func (f *fakeSandboxes) runningPod(claim *sandbox.SandboxClaim, status corev1.ContainerStatus) *corev1.Pod {
	f.mu.Lock()
	defer f.mu.Unlock()
	status.Name = sandbox.WorkerContainer
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{Name: claim.SandboxName, UID: types.UID("pod-" + claim.Name)},
		Status:     corev1.PodStatus{ContainerStatuses: []corev1.ContainerStatus{status}},
	}
	f.sandboxes[claim.SandboxName] = &sandbox.Sandbox{Name: claim.SandboxName, PodName: claim.SandboxName}
	f.pods[claim.SandboxName] = pod
	return pod
}

func exitedAt(at time.Time) corev1.ContainerState {
	return corev1.ContainerState{Terminated: &corev1.ContainerStateTerminated{FinishedAt: metav1.NewTime(at)}}
}

// collectWhenPending is deliverWhenPending with a delivery that always succeeds.
func collectWhenPending(t *testing.T, d *Dispatcher, podUID types.UID, handle <-chan error) {
	t.Helper()
	if err := deliverWhenPending(t, d, podUID, handle, func(workitem.SessionWorkItem) error { return nil }); err != nil {
		t.Fatal(err)
	}
}

// deliverWhenPending passes podUID's work item to deliver once Handle,
// running in the background, has put it. It returns the error deliver
// returned. It fails the test if handle returns first, if the pod had
// collected before, or if no work item appears within hangAfter.
func deliverWhenPending(t *testing.T, d *Dispatcher, podUID types.UID, handle <-chan error, deliver func(workitem.SessionWorkItem) error) error {
	t.Helper()
	deadline := time.Now().Add(hangAfter)
	for {
		status, err := d.workItems.Collect(podUID, deliver)
		if status == workitem.Collected {
			return err
		}
		if status != workitem.NotPending {
			t.Fatalf("pod %s had already collected its work item (Collect returned status %d, not NotPending)", podUID, status)
		}
		select {
		case err := <-handle:
			t.Fatalf("Handle returned %v before pod %s collected its work item", err, podUID)
		default:
		}
		if time.Now().After(deadline) {
			t.Fatalf("no work item appeared for pod %s within %v", podUID, hangAfter)
		}
		time.Sleep(time.Millisecond)
	}
}

func mustPut(t *testing.T, d *Dispatcher, podUID types.UID, item workitem.SessionWorkItem) {
	t.Helper()
	if _, ok := d.workItems.Put(podUID, item); !ok {
		t.Fatalf("Put for pod %s refused", podUID)
	}
}

// hangAfter is how long a test waits where the wait only guards against a
// hang. It is long on purpose: a step that takes milliseconds on an idle
// machine can take seconds on a busy CI runner under the race detector.
const hangAfter = 30 * time.Second

// expiringHandoffWait is the HandoffWait of tests that need Handle's wait to
// run out. Such a test must run inside synctest.Test. There, time moves only
// while every goroutine in the test is blocked. A stalled machine therefore
// cannot end Handle's wait before Handle blocks.
const expiringHandoffWait = 200 * time.Millisecond

// failedWithoutTimeout reports whether Handle failed for a reason other than
// HandoffWait running out. That is ErrHandoffTimeout once the work item is put,
// and context.DeadlineExceeded from these fakes before.
func failedWithoutTimeout(err error) bool {
	return err != nil && !errors.Is(err, ErrHandoffTimeout) && !errors.Is(err, context.DeadlineExceeded)
}

func newDispatcher(t *testing.T, sessions fakeSessions) (*Dispatcher, *fakeQueue, *fakeSandboxes) {
	t.Helper()
	queue := &fakeQueue{}
	sandboxes := &fakeSandboxes{claims: map[string]*sandbox.SandboxClaim{}, sandboxes: map[string]*sandbox.Sandbox{}, pods: map[string]*corev1.Pod{}, podReady: true}
	// Most tests wait on HandoffWait only to catch a hang, so it is hangAfter.
	// Tests that need it to run out set expiringHandoffWait.
	d := New(queue, sessions, sandboxes, workitem.NewStore(), slog.New(slog.NewTextHandler(io.Discard, nil)),
		Options{EnvironmentID: "env_1", BaseURL: "https://api", HandoffWait: hangAfter, PodPollInterval: time.Millisecond, SweepInterval: time.Hour, MaxConcurrentHandoffs: 2})
	d.Now = func() time.Time { return time.Date(2026, 9, 3, 12, 0, 0, 0, time.UTC) }
	return d, queue, sandboxes
}
