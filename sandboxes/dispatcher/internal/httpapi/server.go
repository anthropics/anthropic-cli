package httpapi

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/anthropics/anthropic-cli/sandboxes/dispatcher/internal/workitem"
)

// Authenticator resolves the pod a request came from, or refuses it. The
// upstream path verifies a projected token; the add-on path verifies an
// assignment nonce against the caller's network position.
type Authenticator interface {
	Authenticate(r *http.Request) (Identity, error)
}

// Server holds the collaborators of the pod-facing route.
type Server struct {
	Auth      Authenticator
	WorkItems *workitem.Store
	Logger    *slog.Logger
}

// Handler routes exactly GET /healthz and GET /work-item.
func (s *Server) Handler() http.Handler {
	if s.Auth == nil || s.WorkItems == nil || s.Logger == nil {
		panic("httpapi.Server: nil collaborator")
	}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) })
	mux.HandleFunc("GET /work-item", s.getWorkItem)
	return mux
}

// ListenAndServe serves handler over TLS until ctx ends.
func ListenAndServe(ctx context.Context, addr, certFile, keyFile string, handler http.Handler) error {
	srv := &http.Server{
		Addr:              addr,
		Handler:           handler,
		TLSConfig:         &tls.Config{MinVersion: tls.VersionTLS12},
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
	errc := make(chan error, 1)
	go func() { errc <- srv.ListenAndServeTLS(certFile, keyFile) }()
	select {
	case err := <-errc:
		return err
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return srv.Shutdown(shutdownCtx)
	}
}

// getWorkItem counts a work item as handed out only once its bytes have
// reached the connection: HEAD, which the GET pattern also matches, would consume
// it while sending nothing.
func (s *Server) getWorkItem(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "only GET collects a work item", http.StatusMethodNotAllowed)
		return
	}
	identity, ok := s.authenticate(w, r)
	if !ok {
		return
	}
	var delivered workitem.SessionWorkItem
	status, err := s.WorkItems.Collect(identity.PodUID, func(a workitem.SessionWorkItem) error {
		delivered = a
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(a); err != nil {
			return err
		}
		return http.NewResponseController(w).Flush()
	})
	switch {
	case status == workitem.NotPending:
		http.Error(w, "no work item for this pod yet", http.StatusNotFound)
	case status == workitem.AlreadyCollected:
		http.Error(w, "this pod's work item is closed", http.StatusGone)
	case err != nil:
		s.Logger.Warn("work item write failed; keeping it for the pod's retry", "pod_uid", identity.PodUID, "error", err)
	case status == workitem.CollectedAgain:
		s.Logger.Info("work item handed out again to the same pod", "pod_uid", identity.PodUID, "session_id", delivered.SessionID, "work_id", delivered.WorkID)
	default:
		s.Logger.Info("work item handed out", "pod_uid", identity.PodUID, "session_id", delivered.SessionID, "work_id", delivered.WorkID)
	}
}

// authenticate resolves the caller's pod through the configured
// authenticator, writing the refusal itself when the caller does not verify.
func (s *Server) authenticate(w http.ResponseWriter, r *http.Request) (Identity, bool) {
	identity, err := s.Auth.Authenticate(r)
	switch {
	case err == nil:
		return identity, true
	case errors.Is(err, ErrUnauthenticated):
		http.Error(w, err.Error(), http.StatusUnauthorized)
	case errors.Is(err, ErrNotSandbox), errors.Is(err, ErrNoPod):
		http.Error(w, err.Error(), http.StatusForbidden)
	default:
		s.Logger.Error("authenticating the caller", "error", err)
		http.Error(w, "could not verify the caller", http.StatusInternalServerError)
	}
	return Identity{}, false
}
