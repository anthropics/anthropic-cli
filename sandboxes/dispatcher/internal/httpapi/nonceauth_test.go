package httpapi

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"k8s.io/apimachinery/pkg/types"

	"github.com/anthropics/anthropic-cli/sandboxes/dispatcher/internal/nonce"
)

// fakeRegistry authenticates by comparing the nonce and source IP to a fixed
// expectation, returning the nonce package's own errors.
type fakeRegistry struct {
	nonce, srcIP string
	podUID       types.UID
}

func (f fakeRegistry) Authenticate(_ context.Context, presented, srcIP string) (types.UID, error) {
	if presented != f.nonce {
		return "", nonce.ErrUnknownNonce
	}
	if srcIP != f.srcIP {
		return "", nonce.ErrSourceMismatch
	}
	return f.podUID, nil
}

func requestFrom(t *testing.T, header, remoteAddr string) *http.Request {
	t.Helper()
	r := httptest.NewRequest(http.MethodGet, "/work-item", nil)
	if header != "" {
		r.Header.Set(NonceHeader, header)
	}
	r.RemoteAddr = remoteAddr
	return r
}

func TestNonceAuthMapsRegistryOutcomesToStatuses(t *testing.T) {
	auth := NonceAuth{Registry: fakeRegistry{nonce: "good", srcIP: "10.0.0.5", podUID: "pod-a"}}
	for name, tc := range map[string]struct {
		header, remoteAddr string
		wantErr            error
		wantUID            types.UID
	}{
		"match":          {"good", "10.0.0.5:5000", nil, "pod-a"},
		"no header":      {"", "10.0.0.5:5000", ErrUnauthenticated, ""},
		"wrong nonce":    {"bad", "10.0.0.5:5000", ErrUnauthenticated, ""},
		"wrong source":   {"good", "10.0.0.6:5000", ErrNoPod, ""},
		"no source port": {"good", "10.0.0.5", ErrNoPod, ""},
	} {
		id, err := auth.Authenticate(requestFrom(t, tc.header, tc.remoteAddr))
		if tc.wantErr == nil {
			if err != nil || id.PodUID != tc.wantUID {
				t.Errorf("%s: id=%v err=%v; want pod %q", name, id, err, tc.wantUID)
			}
			continue
		}
		if !errors.Is(err, tc.wantErr) {
			t.Errorf("%s: err = %v, want %v", name, err, tc.wantErr)
		}
	}
}

// The header, not any forwarding header, carries the identity; the source is
// the TCP peer alone, so a spoofed X-Forwarded-For changes nothing.
func TestNonceAuthIgnoresForwardingHeaders(t *testing.T) {
	auth := NonceAuth{Registry: fakeRegistry{nonce: "good", srcIP: "10.0.0.5", podUID: "pod-a"}}
	r := requestFrom(t, "good", "10.0.0.6:5000")
	r.Header.Set("X-Forwarded-For", "10.0.0.5")
	if _, err := auth.Authenticate(r); !errors.Is(err, ErrNoPod) {
		t.Errorf("err = %v; a forged X-Forwarded-For must not stand in for the peer address", err)
	}
}
