package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/anthropics/anthropic-cli/sandboxes/dispatcher/internal/version"
)

// runMainEnv set to 1 makes the test binary run main with its own arguments
// instead of running tests.
const runMainEnv = "DISPATCHER_TEST_RUN_MAIN"

func TestMain(m *testing.M) {
	if os.Getenv(runMainEnv) == "1" {
		main()
		os.Exit(0)
	}
	os.Exit(m.Run())
}

// runMain runs main in a child process and returns its exit status, stdout and
// stderr. The child sees env plus PATH, a scratch HOME and a scratch
// GOCOVERDIR, which keeps a coverage warning off stderr. It fails after 30s.
func runMain(t *testing.T, env []string, args ...string) (int, string, string) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, os.Args[0], args...)
	cmd.WaitDelay = 5 * time.Second
	cmd.Env = append([]string{
		runMainEnv + "=1",
		"PATH=" + os.Getenv("PATH"),
		"HOME=" + t.TempDir(),
		"GOCOVERDIR=" + t.TempDir(),
	}, env...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &stdout, &stderr
	err := cmd.Run()
	if ctx.Err() != nil {
		t.Fatalf("%q: no exit within 30s; stderr %q", args, stderr.String())
	}
	if exit, ok := err.(*exec.ExitError); ok {
		return exit.ExitCode(), stdout.String(), stderr.String()
	}
	if err != nil {
		t.Fatalf("running %q: %v", args, err)
	}
	return 0, stdout.String(), stderr.String()
}

// The child has no environment key. A version query that reached cfg.Load
// would fail instead of printing.
func TestVersionFlag(t *testing.T) {
	for _, args := range [][]string{
		{"-version"},
		{"--version"},
		{"serve", "-version"},
		{"preflight", "-version"},
	} {
		code, stdout, stderr := runMain(t, nil, args...)
		if code != 0 {
			t.Errorf("%q: exit status %d, want 0; stderr %q", args, code, stderr)
		}
		if want := "cma-dispatcher " + version.String() + "\n"; stdout != want {
			t.Errorf("%q: stdout = %q, want %q", args, stdout, want)
		}
		if stderr != "" {
			t.Errorf("%q: stderr = %q, want nothing", args, stderr)
		}
	}
}

func TestServeLogsVersionAtStart(t *testing.T) {
	api := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"type":"error","error":{"type":"authentication_error","message":"invalid key"}}`))
	}))
	defer api.Close()
	env := []string{
		"ANTHROPIC_ENVIRONMENT_KEY=test-key",
		"ANTHROPIC_ENVIRONMENT_ID=env_test",
		"ANTHROPIC_BASE_URL=" + api.URL,
	}

	code, _, stderr := runMain(t, env, "serve", "-namespace", "test", "-log-format", "json")

	if code != 1 {
		t.Fatalf("exit status %d, want 1 from the refused key; stderr %q", code, stderr)
	}
	first, _, _ := strings.Cut(stderr, "\n")
	var line struct {
		Msg     string `json:"msg"`
		Version string `json:"version"`
	}
	if err := json.Unmarshal([]byte(first), &line); err != nil {
		t.Fatalf("first stderr line %q is not a JSON log line: %v", first, err)
	}
	if line.Msg != "cma-dispatcher starting" || line.Version != version.String() {
		t.Errorf("first log line = %q, want msg %q with version %q", first, "cma-dispatcher starting", version.String())
	}
}

// serve and preflight each log the http2debug warning once, at WARN, when
// GODEBUG turns the logging on. With GODEBUG empty they log none.
// serve then stops on the refused key, and preflight on the missing cluster.
func TestSubcommandsWarnOfHTTP2DebugLogging(t *testing.T) {
	api := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"type":"error","error":{"type":"authentication_error","message":"invalid key"}}`))
	}))
	defer api.Close()
	for _, godebug := range []struct {
		value string
		want  int
	}{
		{"http2debug=1", 1},
		{"", 0},
	} {
		for _, c := range []struct{ cmd, stop string }{
			{"serve", "work stats"},
			{"preflight", "connecting to the cluster"},
		} {
			env := []string{
				"ANTHROPIC_ENVIRONMENT_KEY=test-key",
				"ANTHROPIC_ENVIRONMENT_ID=env_test",
				"ANTHROPIC_BASE_URL=" + api.URL,
				"GODEBUG=" + godebug.value,
			}

			code, _, stderr := runMain(t, env, c.cmd, "-namespace", "test", "-log-format", "json")

			if code != 1 || !strings.Contains(stderr, c.stop) {
				t.Fatalf("%s with GODEBUG=%q: exit status %d, want 1 on %q; stderr %q", c.cmd, godebug.value, code, c.stop, stderr)
			}
			warnings := 0
			for _, l := range strings.Split(stderr, "\n") {
				if strings.Contains(l, `"level":"WARN"`) && strings.Contains(l, "remove http2debug from GODEBUG") {
					warnings++
				}
			}
			if warnings != godebug.want {
				t.Errorf("%s with GODEBUG=%q: %d http2debug warnings, want %d; stderr %q", c.cmd, godebug.value, warnings, godebug.want, stderr)
			}
		}
	}
}

// noCluster writes a kubeconfig whose cluster refuses every connection and
// returns the KUBECONFIG variable that names it. preflight needs one to get as
// far as its check of the environment key.
func noCluster(t *testing.T) string {
	t.Helper()
	closed := httptest.NewServer(http.NotFoundHandler())
	closed.Close()
	path := filepath.Join(t.TempDir(), "kubeconfig")
	text := "apiVersion: v1\nkind: Config\ncurrent-context: c\n" +
		"clusters: [{name: c, cluster: {server: \"" + closed.URL + "\"}}]\n" +
		"contexts: [{name: c, context: {cluster: c, user: u}}]\n" +
		"users: [{name: u, user: {}}]\n"
	if err := os.WriteFile(path, []byte(text), 0o600); err != nil {
		t.Fatal(err)
	}
	return "KUBECONFIG=" + path
}

// What serve and preflight print reaches a pod's log or a terminal. For each
// value of ANTHROPIC_BASE_URL below, neither prints a piece of its user name
// or password, in either case: Go lower-cases what it reads as a scheme.
func TestNoSubcommandPrintsTheBaseURLsUserOrPassword(t *testing.T) {
	api := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"type":"error","error":{"type":"authentication_error","message":"invalid key"}}`))
	}))
	t.Cleanup(api.Close)
	closed := httptest.NewServer(http.NotFoundHandler())
	closed.Close()
	host, closedHost := strings.TrimPrefix(api.URL, "http://"), strings.TrimPrefix(closed.URL, "http://")
	const refused = "cma-dispatcher: ANTHROPIC_BASE_URL is not a URL that names a host (its value is not shown)"
	kubeconfig := noCluster(t)

	for _, tc := range []struct{ name, value, words string }{
		{"a space in the password", "http://Qz7Kx9Wv:Jq4Z p8Xy@" + host, refused},
		{"a bad percent escape", "http://Qz7Kx9Wv:Jq4Z%zzp8Xy@" + host, refused},
		{"a password where the port goes", "http://Qz7Kx9Wv:Jq4Zp8Xy", refused},
		{"a slash in the password", "http://Qz7Kx9Wv:Jq4Z/p8Xy@" + host, refused},
		{"no scheme", "Qz7Kx9Wv:Jq4Zp8Xy@" + host, refused},
		{"no host", "http://Qz7Kx9Wv:Jq4Zp8Xy@", refused},
		{"a refused key", "http://Qz7Kx9Wv:Jq4Zp8Xy@" + host, `GET "http://` + host + `/v1/environments/env_test/work/stats`},
		{"no answer", "http://Qz7Kx9Wv:Jq4Zp8Xy@" + closedHost, `Get "http://` + closedHost + `/v1/environments/env_test/work/stats`},
	} {
		for _, cmd := range []string{"serve", "preflight"} {
			for _, format := range []string{"text", "json"} {
				t.Run(tc.name+"/"+cmd+"/"+format, func(t *testing.T) {
					t.Parallel()
					env := []string{
						"ANTHROPIC_ENVIRONMENT_KEY=test-key",
						"ANTHROPIC_ENVIRONMENT_ID=env_test",
						"ANTHROPIC_BASE_URL=" + tc.value,
						kubeconfig,
					}

					code, stdout, stderr := runMain(t, env, cmd, "-namespace", "test", "-log-format", format, "-skip-network-probe")

					printed := stdout + stderr
					if code != 1 || !strings.Contains(printed, tc.words) {
						t.Fatalf("exit status %d, want 1 and %q; stdout %q, stderr %q", code, tc.words, stdout, stderr)
					}
					for _, piece := range []string{"qz7kx9wv", "jq4z", "p8xy"} {
						if strings.Contains(strings.ToLower(printed), piece) {
							t.Errorf("%q of the value is printed; stdout %q, stderr %q", piece, stdout, stderr)
						}
					}
				})
			}
		}
	}
}
