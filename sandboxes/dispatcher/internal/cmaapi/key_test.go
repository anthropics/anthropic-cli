package cmaapi

import (
	"bytes"
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/anthropics/anthropic-sdk-go/v2"
)

// everyCall returns each call of api that reaches the API, by name.
func everyCall(api *Client) map[string]func(context.Context) error {
	return map[string]func(context.Context) error{
		"VerifyEnvironment": api.VerifyEnvironment,
		"Poll": func(ctx context.Context) error {
			_, err := api.Poll(ctx, "w-1", time.Second)
			return err
		},
		"Ack": func(ctx context.Context) error { return api.Ack(ctx, "work_1") },
		"SessionState": func(ctx context.Context) error {
			_, err := api.SessionState(ctx, "session_1")
			return err
		},
	}
}

// A refusal can quote the request it refuses, and every request carries the
// environment key. The dispatcher logs the error of every call. That error
// keeps the refusal's words and holds a placeholder in place of the key.
func TestTheEnvironmentKeyIsInNoCallsError(t *testing.T) {
	const words, sentKey = "the refusal's words", "key-1"
	for shape, refuse := range map[string]func(w http.ResponseWriter, sent string){
		"the API's envelope": func(w http.ResponseWriter, sent string) {
			_, _ = w.Write([]byte(`{"type":"error","error":{"type":"permission_error","message":"` + words + " " + sent + `"}}`))
		},
		"a copy of the request": func(w http.ResponseWriter, sent string) {
			_, _ = w.Write([]byte(`{"message":"` + words + `","request":{"headers":{"authorization":"` + sent + `"}}}`))
		},
	} {
		// Given: an upstream that refuses every request, and puts the
		// request's Authorization header in its body and its request id.
		api := newAPI(t, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Request-Id", r.Header.Get("Authorization"))
			w.WriteHeader(http.StatusForbidden)
			refuse(w, r.Header.Get("Authorization"))
		})
		for name, call := range everyCall(api) {
			t.Run(shape+"/"+name, func(t *testing.T) {
				// When: the call is made and its error is logged.
				err := call(context.Background())
				var logged bytes.Buffer
				slog.New(slog.NewTextHandler(&logged, nil)).Warn("the call failed", "error", err)

				// Then: the line holds the status and the words, and not the key.
				if err == nil || !strings.Contains(logged.String(), "403 Forbidden") || !strings.Contains(logged.String(), words) {
					t.Fatalf("logged %q, want 403 Forbidden and %q", logged.String(), words)
				}
				if strings.Contains(logged.String(), sentKey) {
					t.Errorf("the log line holds the environment key: %s", logged.String())
				}
			})
		}
	}
}

// An error that holds no key is returned as the SDK made it, so a caller can
// still read its status.
func TestARefusalWithoutTheKeyIsUnchanged(t *testing.T) {
	api := newAPI(t, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"type":"error","error":{"type":"authentication_error","message":"no"}}`))
	})
	for name, call := range everyCall(api) {
		var apiErr *anthropic.Error
		if err := call(context.Background()); !errors.As(err, &apiErr) || apiErr.StatusCode != http.StatusUnauthorized {
			t.Errorf("%s: err = %v, want the SDK's error with status 401", name, err)
		}
	}
}
