package cmd

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
)

// toolCall is one agent.tool_use the fake session emits. before, when set,
// runs once the previous call's result has been posted.
type toolCall struct {
	id     string
	name   string
	input  map[string]any
	before func()
}

type toolResult struct {
	isError bool
	text    string
}

// workerRunSession is a fake API server for one session with memory
// stores attached: the session lookup, the stores' listings, the lease
// heartbeat and stop, and an event stream that emits calls one at a time,
// each after the previous result was posted, then terminates the session.
type workerRunSession struct {
	t      *testing.T
	stores []map[string]any
	calls  []toolCall

	mu      sync.Mutex
	results map[string]toolResult
	// uploads is every request body the memory-store endpoints received.
	uploads []string
}

func (s *workerRunSession) result(id string) (toolResult, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	res, ok := s.results[id]
	return res, ok
}

func memoryJSON(storeID, path, content string, full bool) map[string]any {
	sum := sha256.Sum256([]byte(content))
	item := map[string]any{
		"type": "memory", "id": "mem_" + hex.EncodeToString([]byte(path)), "path": path,
		"content_sha256":     hex.EncodeToString(sum[:]),
		"content_size_bytes": len(content),
		"created_at":         "2026-05-11T12:00:00Z", "updated_at": "2026-05-11T12:00:00Z",
		"memory_store_id": storeID, "memory_version_id": "memver_1",
	}
	if full {
		item["content"] = content
	}
	return item
}

func (s *workerRunSession) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	path := r.URL.Path
	switch {
	case strings.HasPrefix(path, "/v1/memory_stores/"):
		storeID := strings.Split(strings.TrimPrefix(path, "/v1/memory_stores/"), "/")[0]
		if r.Method != http.MethodGet {
			body, _ := io.ReadAll(r.Body)
			s.mu.Lock()
			s.uploads = append(s.uploads, string(body))
			s.mu.Unlock()
			var in struct{ Path, Content string }
			_ = json.Unmarshal(body, &in)
			_ = json.NewEncoder(w).Encode(memoryJSON(storeID, in.Path, in.Content, true))
			return
		}
		full := r.URL.Query().Get("view") == "full"
		_ = json.NewEncoder(w).Encode(map[string]any{
			"data":      []any{memoryJSON(storeID, "/note.md", "remembered", full)},
			"next_page": nil,
		})
	case strings.HasSuffix(path, "/heartbeat"):
		_, _ = io.WriteString(w, `{"last_heartbeat":"2026-05-11T12:00:00Z","lease_extended":true,"state":"active","ttl_seconds":30,"type":"work_heartbeat"}`)
	case strings.HasSuffix(path, "/stop"):
		w.WriteHeader(http.StatusNoContent)
	case strings.HasSuffix(path, "/events/stream"):
		s.stream(w, r)
	case strings.HasSuffix(path, "/events") && r.Method == http.MethodPost:
		body, _ := io.ReadAll(r.Body)
		s.mu.Lock()
		for _, ev := range gjson.GetBytes(body, "events").Array() {
			s.results[ev.Get("tool_use_id").String()] = toolResult{
				isError: ev.Get("is_error").Bool(),
				text:    ev.Get("content.0.text").String(),
			}
		}
		s.mu.Unlock()
		_, _ = io.WriteString(w, `{"type":"send_session_events"}`)
	case strings.HasSuffix(path, "/events"):
		_, _ = io.WriteString(w, `{"data":[],"first_id":null,"has_more":false,"last_id":null}`)
	case strings.HasPrefix(path, "/v1/sessions/"):
		_ = json.NewEncoder(w).Encode(map[string]any{
			"id":        "ses_1",
			"agent":     map[string]any{"skills": []any{}},
			"resources": s.stores,
		})
	default:
		s.t.Errorf("unexpected request: %s %s", r.Method, path)
		http.Error(w, "unexpected", http.StatusNotImplemented)
	}
}

func (s *workerRunSession) stream(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	flusher := w.(http.Flusher)
	w.WriteHeader(http.StatusOK)
	flusher.Flush()
	for _, call := range s.calls {
		if call.before != nil {
			call.before()
		}
		data, _ := json.Marshal(map[string]any{
			"type": "agent.tool_use", "id": call.id, "name": call.name,
			"input": call.input, "processed_at": "2026-05-11T12:00:00Z",
		})
		fmt.Fprintf(w, "event: agent.tool_use\ndata: %s\n\n", data)
		flusher.Flush()
		for {
			if _, ok := s.result(call.id); ok {
				break
			}
			select {
			case <-r.Context().Done():
				return
			case <-time.After(2 * time.Millisecond):
			}
		}
	}
	_, _ = io.WriteString(w, "event: session.status_terminated\n"+
		`data: {"type":"session.status_terminated","id":"evt_term","processed_at":"2026-05-11T12:00:00Z"}`+"\n\n")
	flusher.Flush()
}

// The file tools `beta:worker run` serves must reach the session's attached
// memory stores at their mount paths, refuse writes into a read-only store,
// and still refuse everything outside the workdir and the stores. Symlinks
// left in a store must not carry outside files into the end-of-session sync.
func TestWorkerRunFileToolsReachMemoryStores(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("memory store mount paths are POSIX absolute paths")
	}
	clearWorkerEnv(t)

	workdir := t.TempDir()
	outside := t.TempDir()
	notes := filepath.Join(outside, "mnt", "memory", "notes")
	reference := filepath.Join(outside, "mnt", "memory", "reference")
	secret := filepath.Join(outside, "secret.txt")
	require.NoError(t, os.WriteFile(secret, []byte("top secret"), 0o600))
	link := filepath.Join(notes, "link.md")
	linkDir := filepath.Join(notes, "linked")

	session := &workerRunSession{
		t:       t,
		results: map[string]toolResult{},
		stores: []map[string]any{
			{"type": "memory_store", "memory_store_id": "memstore_notes", "name": "notes", "mount_path": notes, "access": "read_write"},
			{"type": "memory_store", "memory_store_id": "memstore_reference", "name": "reference", "mount_path": reference, "access": "read_only"},
		},
		calls: []toolCall{
			{id: "read_store", name: "read", input: map[string]any{"file_path": filepath.Join(notes, "note.md")}},
			{id: "write_store", name: "write", input: map[string]any{"file_path": filepath.Join(notes, "new.md"), "content": "learned"}},
			{id: "write_read_only_store", name: "write", input: map[string]any{"file_path": filepath.Join(reference, "new.md"), "content": "learned"}},
			{id: "read_outside", name: "read", input: map[string]any{"file_path": secret}},
			{id: "read_dot_dot", name: "read", input: map[string]any{"file_path": "../" + filepath.Base(outside) + "/secret.txt"}},
			{
				id: "read_symlink_out_of_store", name: "read", input: map[string]any{"file_path": link},
				before: func() {
					assert.NoError(t, os.Symlink(secret, link))
					assert.NoError(t, os.Symlink(outside, linkDir))
				},
			},
			{id: "read_through_symlinked_dir", name: "read", input: map[string]any{"file_path": filepath.Join(linkDir, "secret.txt")}},
			{id: "glob_store", name: "glob", input: map[string]any{"pattern": "*.md", "path": notes}},
		},
	}
	server := httptest.NewServer(session)
	defer server.Close()

	token, err := json.Marshal(map[string]string{"sessions_token": "tok"})
	require.NoError(t, err)

	done := make(chan error, 1)
	go func() {
		done <- run(t, workerRunCommandDef(), "run",
			"--session-id", "ses_1", "--work-id", "work_1", "--environment-id", "env_1",
			"--base-url", server.URL,
			"--work-secret-file", workSecretFile(t, base64.RawURLEncoding.EncodeToString(token)),
			"--workdir", workdir,
		)
	}()
	select {
	case err := <-done:
		require.NoError(t, err)
	case <-time.After(30 * time.Second):
		t.Fatal("beta:worker run did not finish")
	}

	want := func(id string, isError bool, contains string) {
		t.Helper()
		got, ok := session.result(id)
		require.True(t, ok, "no result posted for %s", id)
		assert.Equal(t, isError, got.isError, "%s: %s", id, got.text)
		assert.Contains(t, got.text, contains, id)
	}
	want("read_store", false, "remembered")
	want("write_store", false, "")
	want("write_read_only_store", true, "read-only")
	want("read_outside", true, "outside the session's working directory")
	want("read_dot_dot", true, "outside the session's working directory")
	want("read_symlink_out_of_store", true, "outside the session's working directory")
	want("read_through_symlinked_dir", true, "outside the session's working directory")
	want("glob_store", false, "note.md")

	// The symlinks were still in the store when the final sync and the store
	// cleanup ran: the sync uploaded the agent's own file and nothing else,
	// and the cleanup removed the links without touching their targets.
	session.mu.Lock()
	defer session.mu.Unlock()
	for id, res := range session.results {
		assert.NotContains(t, res.text, "top secret", id)
	}
	require.Len(t, session.uploads, 1)
	assert.Contains(t, session.uploads[0], "learned")
	assert.NotContains(t, session.uploads[0], "top secret")
	got, err := os.ReadFile(secret)
	require.NoError(t, err)
	assert.Equal(t, "top secret", string(got))
	assert.NoDirExists(t, notes)
}
