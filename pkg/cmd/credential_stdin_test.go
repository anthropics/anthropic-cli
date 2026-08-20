package cmd

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os/exec"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/anthropics/anthropic-cli/internal/requestflag"
	"github.com/anthropics/anthropic-sdk-go/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli/v3"
)

func withCommandLine(t *testing.T, args ...string) {
	t.Helper()
	saved := commandLineArgs
	commandLineArgs = func() []string { return args }
	argvCredentialWarnOnce = sync.Once{}
	t.Cleanup(func() { commandLineArgs = saved; argvCredentialWarnOnce = sync.Once{} })
}

// TestCredentialFlagsOnCommandLine pins what counts as "a credential on the
// command line": both dash styles and the = form, nothing after "--", and
// never the *-stdin flags or unrelated flags that merely share a prefix.
func TestCredentialFlagsOnCommandLine(t *testing.T) {
	for _, tc := range []struct {
		name string
		argv []string
		want []string
	}{
		{"space form", []string{"models", "list", "--api-key", "sk"}, []string{"api-key"}},
		{"equals form", []string{"--auth-token=tok", "models", "list"}, []string{"auth-token"}},
		{"single dash", []string{"-api-key", "sk", "models", "list"}, []string{"api-key"}},
		{"both", []string{"--api-key", "a", "--auth-token", "b"}, []string{"api-key", "auth-token"}},
		{"stdin flag is not the deprecated flag", []string{"--api-key-stdin", "models", "list"}, nil},
		{"lookalike", []string{"beta:memory-stores:memory-versions", "list", "--api-key-id", "x"}, nil},
		{"after terminator", []string{"models", "list", "--", "--api-key", "sk"}, nil},
		{"env only", []string{"models", "list"}, nil},
	} {
		t.Run(tc.name, func(t *testing.T) {
			withCommandLine(t, tc.argv...)
			assert.Equal(t, tc.want, credentialFlagsOnCommandLine())
		})
	}
}

// TestArgvCredentialWarning: the deprecation notice names the flag and the
// replacements, never the secret, and fires once per process.
func TestArgvCredentialWarning(t *testing.T) {
	withCommandLine(t, "--api-key", "sk-ant-SECRET", "models", "list")
	var b strings.Builder
	warnIfCredentialOnCommandLine(&b)
	warnIfCredentialOnCommandLine(&b)
	out := b.String()
	assert.Equal(t, 1, strings.Count(out, "Warning:"), "emits once")
	assert.Contains(t, out, "--api-key is deprecated")
	assert.Contains(t, out, "ANTHROPIC_API_KEY, --api-key-stdin")
	assert.NotContains(t, out, "sk-ant-SECRET")

	withCommandLine(t, "models", "list")
	b.Reset()
	warnIfCredentialOnCommandLine(&b)
	assert.Empty(t, b.String(), "silent when the credential came from env or a profile")
}

// stdinCredentialApp mirrors the root flags readStdinCredential touches so
// tests exercise real urfave parsing without the package-level Command tree
// (which caches flag state across Runs).
func stdinCredentialApp(action func(*cli.Command) error) *cli.Command {
	return &cli.Command{
		Name: "ant",
		Flags: []cli.Flag{
			&requestflag.Flag[string]{Name: "api-key", Sources: cli.EnvVars("ANTHROPIC_API_KEY")},
			&requestflag.Flag[string]{Name: "auth-token", Sources: cli.EnvVars("ANTHROPIC_AUTH_TOKEN")},
			&cli.BoolFlag{Name: "api-key-stdin"},
			&cli.BoolFlag{Name: "auth-token-stdin"},
		},
		Commands: []*cli.Command{{Name: "probe", Action: func(_ context.Context, c *cli.Command) error { return action(c) }}},
	}
}

func TestReadStdinCredential(t *testing.T) {
	clearEnv(t, "ANTHROPIC_API_KEY")
	clearEnv(t, "ANTHROPIC_AUTH_TOKEN")

	type result struct{ apiKey, authToken string }
	runProbe := func(t *testing.T, stdin string, isTTY bool, onCommandLine []string, argv ...string) (result, error) {
		t.Helper()
		var got result
		var readErr error
		app := stdinCredentialApp(func(c *cli.Command) error {
			readErr = readStdinCredential(c, strings.NewReader(stdin), isTTY, onCommandLine)
			got = result{c.String("api-key"), c.String("auth-token")}
			return nil
		})
		require.NoError(t, app.Run(context.Background(), append([]string{"ant"}, argv...)))
		return got, readErr
	}

	t.Run("api key from stdin lands on --api-key, whitespace trimmed", func(t *testing.T) {
		got, err := runProbe(t, "  sk-ant-from-stdin\n", false, nil, "probe", "--api-key-stdin")
		require.NoError(t, err)
		assert.Equal(t, result{apiKey: "sk-ant-from-stdin"}, got)
	})
	t.Run("auth token from stdin lands on --auth-token", func(t *testing.T) {
		got, err := runProbe(t, "tok\n", false, nil, "--auth-token-stdin", "probe")
		require.NoError(t, err)
		assert.Equal(t, result{authToken: "tok"}, got)
	})
	t.Run("stdin beats the env var for the same credential", func(t *testing.T) {
		t.Setenv("ANTHROPIC_API_KEY", "sk-from-env")
		got, err := runProbe(t, "sk-from-stdin", false, nil, "probe", "--api-key-stdin")
		require.NoError(t, err)
		assert.Equal(t, "sk-from-stdin", got.apiKey)
	})
	t.Run("no stdin flag: stdin untouched, nothing set", func(t *testing.T) {
		got, err := runProbe(t, "ignored", false, nil, "probe")
		require.NoError(t, err)
		assert.Equal(t, result{}, got)
	})
	t.Run("both stdin flags", func(t *testing.T) {
		_, err := runProbe(t, "x", false, nil, "probe", "--api-key-stdin", "--auth-token-stdin")
		assert.ErrorContains(t, err, "mutually exclusive")
	})
	t.Run("stdin flag plus the same credential on argv", func(t *testing.T) {
		_, err := runProbe(t, "x", false, []string{"api-key"}, "probe", "--api-key-stdin", "--api-key", "y")
		assert.ErrorContains(t, err, "--api-key-stdin cannot be combined with --api-key")
	})
	t.Run("refuses to read from a terminal", func(t *testing.T) {
		_, err := runProbe(t, "", true, nil, "probe", "--api-key-stdin")
		assert.ErrorContains(t, err, "pipe it in")
	})
	t.Run("empty stdin", func(t *testing.T) {
		_, err := runProbe(t, " \n", false, nil, "probe", "--api-key-stdin")
		assert.ErrorContains(t, err, "stdin was empty")
	})
}

func withStdinCredential(t *testing.T, stdin string) {
	t.Helper()
	savedIn, savedTTY := credentialStdin, credentialStdinIsTTY
	credentialStdin = func() io.Reader { return strings.NewReader(stdin) }
	credentialStdinIsTTY = func() bool { return false }
	stdinCredentialOnce, stdinCredentialErr, stdinConsumedByCredential = sync.Once{}, nil, ""
	t.Cleanup(func() {
		credentialStdin, credentialStdinIsTTY = savedIn, savedTTY
		stdinCredentialOnce, stdinCredentialErr, stdinConsumedByCredential = sync.Once{}, nil, ""
	})
}

// TestAuthStatusStdinCredential: with a logged-in profile, a key piped via
// --api-key-stdin is reported as such — not as --api-key / ANTHROPIC_API_KEY,
// and without the "is set in your environment … unset it" advice, which would
// be false for a per-invocation stdin credential.
func TestAuthStatusStdinCredential(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("ANTHROPIC_CONFIG_DIR", dir)
	clearEnv(t, "ANTHROPIC_PROFILE")
	clearEnv(t, "ANTHROPIC_BASE_URL")
	require.NoError(t, config.SaveProfile(dir, "default", &config.Config{
		AuthenticationInfo: &config.AuthenticationInfo{Type: config.AuthenticationTypeUserOAuth, UserOAuth: &config.UserOAuth{}},
	}))
	exp := time.Now().Add(time.Hour)
	require.NoError(t, config.WriteCredentials(config.ProfileCredentialsPath(dir, "default"),
		config.Credentials{AccessToken: "sk-ant-oat01-X", ExpiresAt: &exp}))
	withCommandLine(t, "--api-key-stdin", "auth", "status")
	withStdinCredential(t, "sk-ant-api03-FROMSTDIN\n")

	out, err := runStatus(t, "--api-key-stdin")
	require.NoError(t, err)
	assert.Contains(t, out, "(active) * --api-key-stdin")
	assert.NotContains(t, out, "--api-key / ANTHROPIC_API_KEY")
	assert.NotContains(t, out, "is set in your environment")
	assert.NotContains(t, out, "FROMSTDIN")
}

// TestMultiAuthWarningNamesStdin: the multi-source notice names the stdin
// flag as the winner rather than the deprecated --api-key.
func TestMultiAuthWarningNamesStdin(t *testing.T) {
	resetWarnOnce(t)
	out := captureStderr(t, func() { warnIfMultipleAuthSources("--api-key-stdin", "", false, false, true) })
	assert.Contains(t, out, "using --api-key-stdin per precedence")
	assert.NotContains(t, out, "ANTHROPIC_API_KEY")
}

// TestAPIKeyStdinEndToEnd runs the real binary: the key piped on stdin is the
// one sent as x-api-key, no deprecation warning is printed, and a piped key
// does not get parsed as a request body. Subprocess for the same reason as
// TestWorkspaceIDHeader.
func TestAPIKeyStdinEndToEnd(t *testing.T) {
	var gotKey string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotKey = r.Header.Get("x-api-key")
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"data":[],"has_more":false,"first_id":null,"last_id":null}`)
	}))
	defer srv.Close()

	cmd := exec.Command("go", "run", "./cmd/ant", "--api-key-stdin", "--base-url", srv.URL, "models", "list")
	cmd.Dir = "../.."
	cmd.Env = append(cmd.Environ(), "ANTHROPIC_CONFIG_DIR="+t.TempDir(), "ANTHROPIC_API_KEY=sk-from-env")
	cmd.Stdin = strings.NewReader("sk-from-stdin\n")
	out, err := cmd.CombinedOutput()
	require.NoError(t, err, "%s", out)
	assert.Equal(t, "sk-from-stdin", gotKey)
	assert.NotContains(t, string(out), "deprecated")

	cmd = exec.Command("go", "run", "./cmd/ant", "--api-key", "sk-from-argv", "--base-url", srv.URL, "models", "list")
	cmd.Dir = "../.."
	cmd.Env = append(cmd.Environ(), "ANTHROPIC_CONFIG_DIR="+t.TempDir())
	out, err = cmd.CombinedOutput()
	require.NoError(t, err, "%s", out)
	assert.Equal(t, "sk-from-argv", gotKey, "deprecated flag still works")
	assert.Contains(t, string(out), "--api-key is deprecated")
	assert.NotContains(t, string(out), "sk-from-argv")

	// A second stdin consumer ("-" file param) must fail loudly rather than
	// upload whatever is left of stdin after the key was read.
	cmd = exec.Command("go", "run", "./cmd/ant", "--api-key-stdin", "--base-url", srv.URL, "beta:files", "upload", "--file", "-")
	cmd.Dir = "../.."
	cmd.Env = append(cmd.Environ(), "ANTHROPIC_CONFIG_DIR="+t.TempDir())
	cmd.Stdin = strings.NewReader("sk-from-stdin\n")
	out, err = cmd.CombinedOutput()
	require.Error(t, err, "%s", out)
	assert.Contains(t, string(out), "stdin was already consumed by --api-key-stdin")
}
