package httpapi

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"

	"k8s.io/apimachinery/pkg/types"

	"github.com/anthropics/anthropic-cli/sandboxes/dispatcher/internal/nonce"
)

// NonceHeader carries the sandbox pod's assignment nonce on the add-on path.
const NonceHeader = "X-CMA-Assignment-Nonce"

// NonceRegistry resolves the pod a nonce-bearing request came from, checking
// the nonce against the sandbox pod's observed network position.
type NonceRegistry interface {
	Authenticate(ctx context.Context, presentedNonce, srcIP string) (types.UID, error)
}

// NonceAuth authenticates a caller by the assignment nonce it presents plus
// the address it connects from.
type NonceAuth struct {
	Registry NonceRegistry
}

// Authenticate reads the nonce header and the request's TCP peer address, and
// resolves the sandbox pod they identify. The source address is the peer
// address alone; a forwarding header a caller could set is never read.
func (a NonceAuth) Authenticate(r *http.Request) (Identity, error) {
	presented := r.Header.Get(NonceHeader)
	if presented == "" {
		return Identity{}, fmt.Errorf("%w: an assignment nonce is required", ErrUnauthenticated)
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return Identity{}, fmt.Errorf("%w: the request has no resolvable source address", ErrNoPod)
	}
	podUID, err := a.Registry.Authenticate(r.Context(), presented, host)
	switch {
	case err == nil:
		return Identity{PodUID: podUID}, nil
	case errors.Is(err, nonce.ErrUnknownNonce):
		return Identity{}, fmt.Errorf("%w: %v", ErrUnauthenticated, err)
	default:
		return Identity{}, fmt.Errorf("%w: %v", ErrNoPod, err)
	}
}
