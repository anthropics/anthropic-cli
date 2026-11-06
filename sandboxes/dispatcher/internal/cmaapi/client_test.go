package cmaapi

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"slices"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

func newAPI(t *testing.T, handler http.HandlerFunc) *Client {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)
	return New(srv.URL, "env_1", "key-1", 5*time.Second)
}

func TestStalledUpstreamTimesOut(t *testing.T) {
	// Given: an upstream that answers each request before stallFrom with a
	// retryable 503 and a 1ms retry hint, and never answers from stallFrom on.
	var requests, stallFrom atomic.Int32
	arrived := make(chan struct{}, 1)
	release := make(chan struct{})
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if requests.Add(1) < stallFrom.Load() {
			w.Header().Set("Retry-After-Ms", "1")
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		select {
		case arrived <- struct{}{}:
		default: // The SDK may retry a stalled attempt. Only the first stall is signalled.
		}
		select {
		case <-r.Context().Done():
		case <-release:
		}
	}))
	t.Cleanup(srv.Close)
	t.Cleanup(func() { close(release) }) // Runs before srv.Close. A handler still blocked on release cannot then hang Close.
	api := New(srv.URL, "env_1", "key-1", 200*time.Millisecond)
	poll := func(ctx context.Context) error {
		_, err := api.Poll(ctx, "w-1", time.Second)
		return err
	}

	for _, tc := range []struct {
		name      string
		stallFrom int32
		call      func(context.Context) error
	}{
		{"verify stalls at once", 1, api.VerifyEnvironment},
		{"poll stalls at once", 1, poll},
		{"verify stalls on the last attempt", 3, api.VerifyEnvironment},
	} {
		// When: the call runs against the upstream with no deadline of its own.
		requests.Store(0)
		stallFrom.Store(tc.stallFrom)
		select {
		case <-arrived:
		default:
		}
		errc := make(chan error, 1)
		go func() { errc <- tc.call(context.Background()) }()
		select {
		case <-arrived:
		case <-time.After(5 * time.Second):
			t.Fatalf("%s: no request stalled within 5s", tc.name)
		}
		// Then: the call fails with the deadline error instead of hanging. The
		// 5s bound has room for the SDK to retry the stalled attempt twice,
		// with its backoff.
		select {
		case err := <-errc:
			if !errors.Is(err, context.DeadlineExceeded) {
				t.Errorf("%s: err = %v, want context.DeadlineExceeded", tc.name, err)
			}
		case <-time.After(5 * time.Second):
			t.Fatalf("%s: still blocked 5s after the 200ms request timeout", tc.name)
		}
	}
}

func requireEnvKey(t *testing.T, r *http.Request) {
	t.Helper()
	if r.Header.Get("Authorization") != "Bearer key-1" {
		t.Errorf("%s %s: Authorization = %q", r.Method, r.URL.Path, r.Header.Get("Authorization"))
	}
	if r.Header.Get("X-Api-Key") != "" {
		t.Errorf("%s %s: an API key header must not be sent", r.Method, r.URL.Path)
	}
}

func TestSessionState(t *testing.T) {
	for name, tc := range map[string]struct {
		status int
		body   string
		want   SessionState
		errs   bool
	}{
		"running":    {200, `{"id":"session_1","type":"session","status":"running"}`, SessionAlive, false},
		"terminated": {200, `{"id":"session_1","type":"session","status":"terminated"}`, SessionEnded, false},
		"gone":       {404, `{"type":"error","error":{"type":"not_found_error","message":"Session not found: session_1"}}`, SessionGone, false},
		"forbidden":  {403, `{"type":"error","error":{"type":"permission_error","message":"no"}}`, SessionAlive, true},
		"outage":     {503, `{"type":"error","error":{"type":"api_error","message":"no"}}`, SessionAlive, true},
	} {
		api := newAPI(t, func(w http.ResponseWriter, r *http.Request) {
			requireEnvKey(t, r)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(tc.status)
			_, _ = w.Write([]byte(tc.body))
		})
		got, err := api.SessionState(context.Background(), "session_1")
		if (err != nil) != tc.errs || got != tc.want {
			t.Errorf("%s: state %v err %v, want %v (error %v)", name, got, err, tc.want, tc.errs)
		}
	}
}

func TestPollAndAck(t *testing.T) {
	// Poll forwards the secret unread. Only the sandbox's worker decodes it.
	secret := "a format this dispatcher has never seen"
	var acked string
	api := newAPI(t, func(w http.ResponseWriter, r *http.Request) {
		requireEnvKey(t, r)
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/v1/environments/env_1/work/poll":
			if q := r.URL.Query(); q.Get("reclaim_older_than_ms") != "120000" || q.Get("block_ms") != "999" {
				t.Errorf("poll query = %v", q)
			}
			if r.Header.Get("Anthropic-Worker-Id") != "w-1" {
				t.Errorf("worker id header = %q", r.Header.Get("Anthropic-Worker-Id"))
			}
			_ = json.NewEncoder(w).Encode(map[string]any{"id": "work_1", "type": "work", "environment_id": "env_1",
				"data": map[string]any{"type": "session", "id": "session_1"}, "secret": secret, "state": "starting"})
		case r.Method == http.MethodPost && r.URL.Path == "/v1/environments/env_1/work/work_1/ack":
			acked = "work_1"
			_ = json.NewEncoder(w).Encode(map[string]any{"id": "work_1", "type": "work", "state": "active"})
		case r.URL.Path == "/v1/environments/env_1/work/stats":
			_ = json.NewEncoder(w).Encode(map[string]any{"type": "work_queue_stats", "depth": 0, "pending": 0, "workers_polling": 1})
		default:
			t.Errorf("unexpected %s %s", r.Method, r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
		}
	})
	if err := api.VerifyEnvironment(context.Background()); err != nil {
		t.Fatal(err)
	}
	item, err := api.Poll(context.Background(), "w-1", 120*time.Second)
	if err != nil {
		t.Fatal(err)
	}
	if item == nil || item.ID != "work_1" || item.SessionID != "session_1" || item.Secret != secret {
		t.Fatalf("item = %+v; want the queue's secret carried verbatim", item)
	}
	if err := api.Ack(context.Background(), item.ID); err != nil || acked != "work_1" {
		t.Errorf("ack: %v (acked %q)", err, acked)
	}
}

func TestPollEmpty(t *testing.T) {
	api := newAPI(t, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`null`))
	})
	item, err := api.Poll(context.Background(), "w-1", time.Second)
	if err != nil || item != nil {
		t.Errorf("empty poll: %+v, %v", item, err)
	}
}

// A proxy in front of the API can answer 404 itself, for example while it has
// no route to the server, and the API answers an organization without the beta
// with not_found_error and "not found". Neither says anything about the session,
// and the caller deletes the sandbox of a session that is gone.
func TestA404ThatIsNotTheAPIsOwnLeavesTheSessionAlive(t *testing.T) {
	for name, page := range map[string]struct{ contentType, body string }{
		"an HTML page":                {"text/html", "<html><body><h1>404 Not Found</h1></body></html>"},
		"plain text":                  {"text/plain; charset=utf-8", "404 page not found\n"},
		"no body":                     {"", ""},
		"JSON of another shape":       {"application/json", `{"message":"no Route matched with those values"}`},
		"an envelope of another type": {"application/json", `{"type":"error","error":{"type":"api_error","message":"no"}}`},
		"the beta gate's answer":      {"application/json", `{"type":"error","error":{"type":"not_found_error","message":"not found"}}`},
		"another session's id":        {"application/json", `{"type":"error","error":{"type":"not_found_error","message":"Session not found: session_2"}}`},
	} {
		t.Run(name, func(t *testing.T) {
			// Given: an upstream that answers every request 404 with page.
			api := newAPI(t, func(w http.ResponseWriter, _ *http.Request) {
				if page.contentType != "" {
					w.Header().Set("Content-Type", page.contentType)
				}
				w.WriteHeader(http.StatusNotFound)
				_, _ = w.Write([]byte(page.body))
			})

			// When: the session's state is asked for.
			got, err := api.SessionState(context.Background(), "session_1")

			// Then: the call fails, and the session does not read as gone.
			if err == nil || got != SessionAlive {
				t.Fatalf("state %v err %v, want SessionAlive and an error", got, err)
			}
			if !strings.Contains(err.Error(), "404 Not Found") {
				t.Errorf("the error lacks the status: %v", err)
			}
		})
	}
}

// withUserinfo returns address, an http URL, with a made-up user name and
// password in it.
func withUserinfo(address string) string {
	return strings.Replace(address, "http://", "http://Qz7Kx9Wv:Jq4Z%20p8Xy@", 1)
}

// The dispatcher logs the error of every call. The SDK's error for a refusal
// quotes the request's URL as it stands, and Go's error for a request that got
// no answer masks a password and leaves the user name.
func TestTheBaseURLsUserAndPasswordAreInNoCallsError(t *testing.T) {
	refusing := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"type":"error","error":{"type":"authentication_error","message":"no"}}`))
	}))
	t.Cleanup(refusing.Close)
	closed := httptest.NewServer(http.NotFoundHandler())
	closed.Close()

	for _, upstream := range []struct{ name, address, words string }{
		{"a refusal", refusing.URL, "401 Unauthorized"},
		{"no answer", closed.URL, "connection refused"},
	} {
		api := New(withUserinfo(upstream.address), "env_1", "key-1", 5*time.Second)
		for name, call := range everyCall(api) {
			t.Run(upstream.name+"/"+name, func(t *testing.T) {
				t.Parallel()
				err := call(context.Background())
				var logged bytes.Buffer
				slog.New(slog.NewTextHandler(&logged, nil)).Warn("the call failed", "error", err)

				host := strings.TrimPrefix(upstream.address, "http://")
				if err == nil || !strings.Contains(logged.String(), upstream.words) || !strings.Contains(logged.String(), "http://"+host+"/v1/") {
					t.Fatalf("logged %q, want %q and the address http://%s", logged.String(), upstream.words, host)
				}
				for _, piece := range []string{"qz7kx9wv", "jq4z", "p8xy"} {
					if strings.Contains(strings.ToLower(logged.String()), piece) {
						t.Errorf("the log line holds %q of the base URL: %s", piece, logged.String())
					}
				}
			})
		}
	}
}

// Go sends a URL's user name and password only on a request that has no
// Authorization header, and every request here carries the environment key in
// it. An upstream therefore sees one request, whether the base URL holds them
// or not.
func TestTheBaseURLsUserAndPasswordChangeNoRequest(t *testing.T) {
	seen := make(chan string, 1)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var headers []string
		for name, values := range r.Header {
			headers = append(headers, name+": "+strings.Join(values, ","))
		}
		slices.Sort(headers)
		seen <- r.Method + " " + r.RequestURI + " host " + r.Host + "\n" + strings.Join(headers, "\n")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"type":"error","error":{"type":"not_found_error","message":"Session not found: session_1"}}`))
	}))
	t.Cleanup(srv.Close)

	plain := everyCall(New(srv.URL, "env_1", "key-1", 5*time.Second))
	held := New(withUserinfo(srv.URL), "env_1", "key-1", 5*time.Second)
	for name, call := range everyCall(held) {
		_ = plain[name](context.Background())
		want := <-seen
		_ = call(context.Background())
		if got := <-seen; got != want {
			t.Errorf("%s: the upstream saw %q, want the request of a base URL without them, %q", name, got, want)
		}
		if !strings.Contains(want, "Authorization: Bearer key-1") {
			t.Errorf("%s: the upstream saw %q, want the environment key", name, want)
		}
	}
	// A caller can still read the status of a refusal.
	state, err := held.SessionState(context.Background(), "session_1")
	<-seen
	if err != nil || state != SessionGone {
		t.Errorf("SessionState = %v, %v, want SessionGone from the 404", state, err)
	}
}
