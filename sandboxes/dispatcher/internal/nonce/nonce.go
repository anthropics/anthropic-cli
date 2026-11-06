// Package nonce authenticates a sandbox pod to the dispatcher without any
// Kubernetes-issued credential. The GKE Agent Sandbox add-on's admission
// policy bans every kubelet credential in a sandbox pod (projected tokens,
// ClusterTrustBundles, PodCertificates), so the dispatcher stamps a random
// per-sandbox secret into the pod and checks possession of that secret
// together with the pod's network position, which its own watch observes.
package nonce

import (
	"context"
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"net"
	"sync"

	"k8s.io/apimachinery/pkg/types"
)

var (
	// ErrUnknownNonce means the presented nonce matched no live sandbox.
	ErrUnknownNonce = errors.New("assignment nonce is not valid")
	// ErrNoRunningPod means the nonce's sandbox has no pod observed Running.
	ErrNoRunningPod = errors.New("the sandbox has no running pod")
	// ErrSourceMismatch means the caller's address is not the sandbox pod's.
	ErrSourceMismatch = errors.New("request did not come from the sandbox pod")
)

// nonceBytes is the length of a fresh nonce before base64url encoding.
const nonceBytes = 32

// PodLocator answers where a sandbox's pod is, from the dispatcher's own
// watch. A pod that is absent, terminating, or not yet Running is reported
// running=false; every ambiguity a caller cannot resolve is a rejection.
type PodLocator interface {
	SandboxPod(ctx context.Context, sandboxName string) (podUID types.UID, podIP string, running bool, err error)
}

// PodLocatorFunc adapts a function to a PodLocator, so the sandbox client can
// be wired as the locator after the registry it depends on is built.
type PodLocatorFunc func(ctx context.Context, sandboxName string) (types.UID, string, bool, error)

// SandboxPod calls f.
func (f PodLocatorFunc) SandboxPod(ctx context.Context, sandboxName string) (types.UID, string, bool, error) {
	return f(ctx, sandboxName)
}

// Registry maps each live sandbox's nonce to the sandbox, and authenticates a
// pod's request against the nonce plus the sandbox pod's network position. It
// is safe for concurrent use.
type Registry struct {
	locator PodLocator
	mu      sync.Mutex
	// bySandbox holds each sandbox's current nonce, keyed by sandbox name.
	bySandbox map[string]string
}

// New returns a Registry that reads pod positions through locator.
func New(locator PodLocator) *Registry {
	if locator == nil {
		panic("nonce.New: nil PodLocator")
	}
	return &Registry{locator: locator, bySandbox: map[string]string{}}
}

// Mint returns a fresh nonce for a new sandbox and registers it under
// sandboxName. The caller stamps the returned value into the sandbox's fetch
// container so the pod can present it.
func (r *Registry) Mint(sandboxName string) (string, error) {
	nonce, err := newNonce()
	if err != nil {
		return "", err
	}
	r.Set(sandboxName, nonce)
	return nonce, nil
}

// Set registers nonce as sandboxName's current nonce, replacing any earlier
// one.
func (r *Registry) Set(sandboxName, nonce string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.bySandbox[sandboxName] = nonce
}

// Known reports whether sandboxName currently has a registered nonce. A
// restarted dispatcher holds none, so its caller recreates or re-arms the
// sandbox rather than trust one whose nonce it cannot check.
func (r *Registry) Known(sandboxName string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	_, ok := r.bySandbox[sandboxName]
	return ok
}

// Clear forgets sandboxName's nonce once its sandbox is gone.
func (r *Registry) Clear(sandboxName string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.bySandbox, sandboxName)
}

// Authenticate resolves the pod a request came from, or refuses it. It
// accepts only when the presented nonce matches a live sandbox, that
// sandbox's pod is observed Running, and srcIP is that pod's address. srcIP
// is the request's TCP peer address; a caller-supplied forwarding header is
// never its source. Every read failure or missing pod is a refusal.
func (r *Registry) Authenticate(ctx context.Context, presentedNonce, srcIP string) (types.UID, error) {
	sandboxName, ok := r.match(presentedNonce)
	if !ok {
		return "", ErrUnknownNonce
	}
	podUID, podIP, running, err := r.locator.SandboxPod(ctx, sandboxName)
	if err != nil || !running {
		return "", ErrNoRunningPod
	}
	if !sameIP(srcIP, podIP) {
		return "", ErrSourceMismatch
	}
	return podUID, nil
}

// match returns the sandbox whose nonce equals presented, comparing every
// registered nonce in constant time so a hit is not distinguishable by
// timing. An empty nonce never matches.
func (r *Registry) match(presented string) (string, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	matched := ""
	found := false
	for sandboxName, nonce := range r.bySandbox {
		if presented != "" && subtle.ConstantTimeCompare([]byte(nonce), []byte(presented)) == 1 {
			matched, found = sandboxName, true
		}
	}
	return matched, found
}

func newNonce() (string, error) {
	b := make([]byte, nonceBytes)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

// sameIP reports whether two textual addresses are the same IP, comparing
// parsed addresses so equal IPs written differently still match. An
// unparseable address never matches.
func sameIP(a, b string) bool {
	ipA, ipB := net.ParseIP(a), net.ParseIP(b)
	return ipA != nil && ipB != nil && ipA.Equal(ipB)
}
