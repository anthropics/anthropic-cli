package nonce

import (
	"context"
	"errors"
	"testing"

	"k8s.io/apimachinery/pkg/types"
)

// fakeLocator answers SandboxPod from a fixed table, so the tests drive the
// pod's IP, Running state, and read errors.
type fakeLocator map[string]podState

type podState struct {
	uid     types.UID
	ip      string
	running bool
	err     error
}

func (f fakeLocator) SandboxPod(_ context.Context, sandboxName string) (types.UID, string, bool, error) {
	p := f[sandboxName]
	return p.uid, p.ip, p.running, p.err
}

func TestNonceMatchSucceeds(t *testing.T) {
	loc := fakeLocator{"sb-a": {uid: "pod-a", ip: "10.0.0.5", running: true}}
	r := New(loc)
	nonce, err := r.Mint("sb-a")
	if err != nil {
		t.Fatal(err)
	}
	uid, err := r.Authenticate(context.Background(), nonce, "10.0.0.5")
	if err != nil || uid != "pod-a" {
		t.Fatalf("Authenticate = %q, %v; want pod-a with the matching nonce and source IP", uid, err)
	}
}

func TestWrongNonceRejected(t *testing.T) {
	loc := fakeLocator{"sb-a": {uid: "pod-a", ip: "10.0.0.5", running: true}}
	r := New(loc)
	if _, err := r.Mint("sb-a"); err != nil {
		t.Fatal(err)
	}
	for _, presented := range []string{"", "not-the-nonce"} {
		if _, err := r.Authenticate(context.Background(), presented, "10.0.0.5"); !errors.Is(err, ErrUnknownNonce) {
			t.Errorf("nonce %q: err = %v, want ErrUnknownNonce", presented, err)
		}
	}
}

func TestCorrectNonceFromWrongSourceIPRejected(t *testing.T) {
	loc := fakeLocator{"sb-a": {uid: "pod-a", ip: "10.0.0.5", running: true}}
	r := New(loc)
	nonce, _ := r.Mint("sb-a")
	if _, err := r.Authenticate(context.Background(), nonce, "10.0.0.6"); !errors.Is(err, ErrSourceMismatch) {
		t.Errorf("err = %v, want ErrSourceMismatch: possession of the nonce is not enough from another pod's address", err)
	}
}

func TestPodNotRunningRejected(t *testing.T) {
	for name, state := range map[string]podState{
		"pending":    {uid: "pod-a", ip: "10.0.0.5", running: false},
		"read error": {err: errors.New("apiserver down")},
		"gone":       {},
	} {
		loc := fakeLocator{"sb-a": state}
		r := New(loc)
		nonce, _ := r.Mint("sb-a")
		if _, err := r.Authenticate(context.Background(), nonce, "10.0.0.5"); !errors.Is(err, ErrNoRunningPod) {
			t.Errorf("%s: err = %v, want ErrNoRunningPod (fail closed)", name, err)
		}
	}
}

// A pod deleted or moved to a new address fails the live position check on
// the next request, without any explicit invalidation call: Authenticate
// reads the sandbox's current pod each time.
func TestDeletionOrIPChangeInvalidates(t *testing.T) {
	loc := fakeLocator{"sb-a": {uid: "pod-a", ip: "10.0.0.5", running: true}}
	r := New(loc)
	nonce, _ := r.Mint("sb-a")
	if _, err := r.Authenticate(context.Background(), nonce, "10.0.0.5"); err != nil {
		t.Fatalf("baseline: %v", err)
	}
	loc["sb-a"] = podState{}
	if _, err := r.Authenticate(context.Background(), nonce, "10.0.0.5"); !errors.Is(err, ErrNoRunningPod) {
		t.Errorf("after deletion: err = %v, want ErrNoRunningPod", err)
	}
	loc["sb-a"] = podState{uid: "pod-b", ip: "10.0.0.9", running: true}
	if _, err := r.Authenticate(context.Background(), nonce, "10.0.0.5"); !errors.Is(err, ErrSourceMismatch) {
		t.Errorf("after IP change: err = %v, want ErrSourceMismatch", err)
	}
}

// After a restart the registry starts empty and refuses every request until
// the sandbox's nonce is set again from the durable Sandbox object.
func TestRestartFailsClosedUntilReregistered(t *testing.T) {
	loc := fakeLocator{"sb-a": {uid: "pod-a", ip: "10.0.0.5", running: true}}
	first := New(loc)
	nonce, _ := first.Mint("sb-a")

	restarted := New(loc)
	if _, err := restarted.Authenticate(context.Background(), nonce, "10.0.0.5"); !errors.Is(err, ErrUnknownNonce) {
		t.Errorf("before re-registration: err = %v, want ErrUnknownNonce", err)
	}
	restarted.Set("sb-a", nonce)
	if _, err := restarted.Authenticate(context.Background(), nonce, "10.0.0.5"); err != nil {
		t.Errorf("after re-registration from the Sandbox: %v", err)
	}
}

func TestClearForgetsTheNonce(t *testing.T) {
	loc := fakeLocator{"sb-a": {uid: "pod-a", ip: "10.0.0.5", running: true}}
	r := New(loc)
	nonce, _ := r.Mint("sb-a")
	r.Clear("sb-a")
	if _, err := r.Authenticate(context.Background(), nonce, "10.0.0.5"); !errors.Is(err, ErrUnknownNonce) {
		t.Errorf("err = %v, want ErrUnknownNonce after Clear", err)
	}
}

func TestMintGivesEachSandboxADistinctNonce(t *testing.T) {
	r := New(fakeLocator{})
	a, _ := r.Mint("sb-a")
	b, _ := r.Mint("sb-b")
	if a == "" || a == b {
		t.Errorf("nonces a=%q b=%q; want two distinct non-empty values", a, b)
	}
}
