package httpapi

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	authnv1 "k8s.io/api/authentication/v1"
	"k8s.io/apimachinery/pkg/runtime"
	kubefake "k8s.io/client-go/kubernetes/fake"
	k8stesting "k8s.io/client-go/testing"

	"github.com/anthropics/anthropic-cli/sandboxes/dispatcher/internal/workitem"
)

// fakeReviewer maps tokens to outcomes.
type fakeReviewer map[string]any // Identity or error

func (f fakeReviewer) Review(_ context.Context, token string) (Identity, error) {
	switch v := f[token].(type) {
	case Identity:
		return v, nil
	case error:
		return Identity{}, v
	}
	return Identity{}, ErrUnauthenticated
}

func newServer(t *testing.T) (*httptest.Server, *workitem.Store) {
	t.Helper()
	return newServerLoggingTo(t, io.Discard)
}

func newServerLoggingTo(t *testing.T, logs io.Writer) (*httptest.Server, *workitem.Store) {
	t.Helper()
	store := workitem.NewStore()
	reviewer := fakeReviewer{
		"tok-a":        Identity{PodUID: "pod-a"},
		"tok-b":        Identity{PodUID: "pod-b"},
		"tok-badaud":   ErrUnauthenticated,
		"tok-other-sa": ErrNotSandbox,
		"tok-nopod":    ErrNoPod,
	}
	s := &Server{Auth: TokenAuth{Reviewer: reviewer}, WorkItems: store, Logger: slog.New(slog.NewTextHandler(logs, nil))}
	srv := httptest.NewServer(s.Handler())
	t.Cleanup(srv.Close)
	return srv, store
}

func call(t *testing.T, srv *httptest.Server, method, path, token, body string) (int, string) {
	t.Helper()
	req, _ := http.NewRequest(method, srv.URL+path, strings.NewReader(body))
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	resp, err := srv.Client().Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	out, _ := io.ReadAll(resp.Body)
	return resp.StatusCode, string(out)
}

func TestWorkItemAuth(t *testing.T) {
	srv, store := newServer(t)
	store.Put("pod-a", workitem.SessionWorkItem{SessionID: "session_a"})
	for _, tc := range []struct {
		token string
		want  int
	}{
		{"", http.StatusUnauthorized},
		{"tok-badaud", http.StatusUnauthorized},
		{"tok-other-sa", http.StatusForbidden},
		{"tok-nopod", http.StatusForbidden},
		{"tok-b", http.StatusNotFound},
	} {
		if got, _ := call(t, srv, http.MethodGet, "/work-item", tc.token, ""); got != tc.want {
			t.Errorf("token %q: status %d, want %d", tc.token, got, tc.want)
		}
	}
	if status, _ := store.Collect("pod-a", func(workitem.SessionWorkItem) error { return nil }); status != workitem.Collected {
		t.Error("refused requests must not consume the work item")
	}
}

func TestSamePodGetsItsWorkItemUntilClosed(t *testing.T) {
	var logs strings.Builder
	srv, store := newServerLoggingTo(t, &logs)
	store.Put("pod-a", workitem.SessionWorkItem{SessionID: "session_a", WorkID: "work_a", Secret: "st-secret", BaseURL: "https://api", EnvironmentID: "env_a"})
	status, first := call(t, srv, http.MethodGet, "/work-item", "tok-a", "")
	if status != http.StatusOK {
		t.Fatalf("status %d: %s", status, first)
	}
	if !strings.Contains(first, `"secret":"st-secret"`) {
		t.Errorf("body %s; the work item must carry the queue's payload under the secret key", first)
	}
	var got workitem.SessionWorkItem
	if err := json.Unmarshal([]byte(first), &got); err != nil || got.Secret != "st-secret" || got.WorkID != "work_a" {
		t.Errorf("body %s (%v)", first, err)
	}
	if status, again := call(t, srv, http.MethodGet, "/work-item", "tok-a", ""); status != http.StatusOK || again != first {
		t.Errorf("second GET from the same pod = %d %s; want 200 with the same body", status, again)
	}
	if status, _ := call(t, srv, http.MethodGet, "/work-item", "tok-b", ""); status != http.StatusNotFound {
		t.Errorf("another pod's status %d, want 404", status)
	}
	if !strings.Contains(logs.String(), "handed out again") || !strings.Contains(logs.String(), "session_a") {
		t.Errorf("the repeat answer was not logged with its session: %s", logs.String())
	}
	if strings.Contains(logs.String(), "st-secret") {
		t.Error("the work item's secret reached the log")
	}
	store.Close("pod-a")
	if status, _ := call(t, srv, http.MethodGet, "/work-item", "tok-a", ""); status != http.StatusGone {
		t.Errorf("status after Close %d, want 410", status)
	}
}

func TestRoutingIsExact(t *testing.T) {
	srv, store := newServer(t)
	store.Put("pod-a", workitem.SessionWorkItem{SessionID: "session_a"})
	for _, tc := range []struct{ method, path string }{
		{http.MethodPost, "/work-item"}, {http.MethodHead, "/work-item"}, {http.MethodGet, "/work-item/x"}, {http.MethodGet, "/"},
	} {
		if status, _ := call(t, srv, tc.method, tc.path, "tok-a", ""); status != http.StatusNotFound && status != http.StatusMethodNotAllowed {
			t.Errorf("%s %s: status %d", tc.method, tc.path, status)
		}
	}
	if status, _ := store.Collect("pod-a", func(workitem.SessionWorkItem) error { return nil }); status != workitem.Collected {
		t.Error("a request on the wrong method or path must not consume the work item")
	}
	if status, _ := call(t, srv, http.MethodGet, "/healthz", "", ""); status != http.StatusOK {
		t.Errorf("healthz status %d", status)
	}
}

func TestKubeTokenReviewer(t *testing.T) {
	kube := kubefake.NewSimpleClientset()
	kube.PrependReactor("create", "tokenreviews", func(action k8stesting.Action) (bool, runtime.Object, error) {
		review := action.(k8stesting.CreateAction).GetObject().(*authnv1.TokenReview)
		if len(review.Spec.Audiences) != 1 || review.Spec.Audiences[0] != "cma-dispatcher" {
			t.Errorf("audiences = %v", review.Spec.Audiences)
		}
		sandboxUser := authnv1.UserInfo{Username: "system:serviceaccount:ns:cma-sandbox", Extra: map[string]authnv1.ExtraValue{PodUIDExtraKey: {"pod-1"}}}
		switch review.Spec.Token {
		case "sandbox":
			review.Status = authnv1.TokenReviewStatus{Authenticated: true, Audiences: []string{"cma-dispatcher"}, User: sandboxUser}
		case "apiserver-only":
			review.Status = authnv1.TokenReviewStatus{Authenticated: true, User: sandboxUser}
		case "other-sa":
			review.Status = authnv1.TokenReviewStatus{Authenticated: true, Audiences: []string{"cma-dispatcher"}, User: authnv1.UserInfo{
				Username: "system:serviceaccount:ns:default", Extra: map[string]authnv1.ExtraValue{PodUIDExtraKey: {"pod-1"}}}}
		case "unbound":
			review.Status = authnv1.TokenReviewStatus{Authenticated: true, Audiences: []string{"cma-dispatcher"}, User: authnv1.UserInfo{Username: "system:serviceaccount:ns:cma-sandbox"}}
		default:
			review.Status = authnv1.TokenReviewStatus{Authenticated: false}
		}
		return true, review, nil
	})
	if _, err := NewKubeTokenReviewer(kube, "ns", "default", "cma-dispatcher"); err == nil {
		t.Error("the default service account must be refused")
	}
	r, err := NewKubeTokenReviewer(kube, "ns", "cma-sandbox", "cma-dispatcher")
	if err != nil {
		t.Fatal(err)
	}
	id, err := r.Review(context.Background(), "sandbox")
	if err != nil || id.PodUID != "pod-1" {
		t.Errorf("sandbox token: %+v, %v", id, err)
	}
	for token, want := range map[string]error{"wrong-audience": ErrUnauthenticated, "apiserver-only": ErrUnauthenticated, "other-sa": ErrNotSandbox, "unbound": ErrNoPod} {
		if _, err := r.Review(context.Background(), token); !errors.Is(err, want) {
			t.Errorf("%s: err = %v, want %v", token, err, want)
		}
	}
}
