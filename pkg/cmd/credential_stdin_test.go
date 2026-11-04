package cmd

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/anthropics/anthropic-sdk-go/v2/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli/v3"
)

// stdinCredentialApp mirrors the root flags the credential helpers touch so
// tests exercise real urfave parsing without the package-level Command tree
// (which caches flag state across Runs).
func stdinCredentialApp(action func(*cli.Command) error) *cli.Command {
	return &cli.Command{
		Name: "ant",
		Flags: []cli.Flag{
			&cli.BoolFlag{Name: "api-key-stdin"},
			&cli.BoolFlag{Name: "auth-token-stdin"},
		},
		Commands: []*cli.Command{{Name: "probe", Action: func(_ context.Context, c *cli.Command) error { return action(c) }}},
	}
}

func TestReadStdinCredential(t *testing.T) {
	clearEnv(t, "ANTHROPIC_API_KEY")
	clearEnv(t, "ANTHROPIC_AUTH_TOKEN")

	type result struct{ target, secret string }
	runProbe := func(t *testing.T, stdin string, isTTY bool, argv ...string) (result, error) {
		t.Helper()
		var got result
		var readErr error
		app := stdinCredentialApp(func(c *cli.Command) error {
			got.target, got.secret, readErr = readStdinCredential(c, strings.NewReader(stdin), isTTY)
			return nil
		})
		require.NoError(t, app.Run(context.Background(), append([]string{"ant"}, argv...)))
		return got, readErr
	}

	t.Run("api key from stdin, whitespace trimmed", func(t *testing.T) {
		got, err := runProbe(t, "  sk-ant-from-stdin\n", false, "probe", "--api-key-stdin")
		require.NoError(t, err)
		assert.Equal(t, result{"api-key", "sk-ant-from-stdin"}, got)
	})
	t.Run("auth token from stdin", func(t *testing.T) {
		got, err := runProbe(t, "tok\n", false, "--auth-token-stdin", "probe")
		require.NoError(t, err)
		assert.Equal(t, result{"auth-token", "tok"}, got)
	})
	t.Run("no stdin flag: stdin untouched, nothing read", func(t *testing.T) {
		got, err := runProbe(t, "ignored", false, "probe")
		require.NoError(t, err)
		assert.Equal(t, result{}, got)
	})
	t.Run("both stdin flags", func(t *testing.T) {
		_, err := runProbe(t, "x", false, "probe", "--api-key-stdin", "--auth-token-stdin")
		assert.ErrorContains(t, err, "mutually exclusive")
	})
	t.Run("refuses to read from a terminal", func(t *testing.T) {
		_, err := runProbe(t, "", true, "probe", "--api-key-stdin")
		assert.ErrorContains(t, err, "pipe it in")
	})
	t.Run("empty stdin", func(t *testing.T) {
		_, err := runProbe(t, " \n", false, "probe", "--api-key-stdin")
		assert.ErrorContains(t, err, "stdin was empty")
	})
}

func resetStdinCredential() {
	stdinCredentialOnce, stdinCredentialErr, stdinConsumedByCredential = sync.Once{}, nil, ""
	stdinCredential.target, stdinCredential.value = "", ""
}

func withStdinCredential(t *testing.T, stdin string) {
	t.Helper()
	savedIn, savedTTY := credentialStdin, credentialStdinIsTTY
	credentialStdin = func() io.Reader { return strings.NewReader(stdin) }
	credentialStdinIsTTY = func() bool { return false }
	resetStdinCredential()
	t.Cleanup(func() {
		credentialStdin, credentialStdinIsTTY = savedIn, savedTTY
		resetStdinCredential()
	})
}

// TestDirectCredential: stdin beats the env var for its own credential and
// leaves the other credential's env var in place.
func TestDirectCredential(t *testing.T) {
	type result struct {
		apiKey, authToken string
		apiKeySet, tokSet bool
	}
	run := func(t *testing.T, stdin string, argv ...string) result {
		t.Helper()
		withStdinCredential(t, stdin)
		var got result
		app := stdinCredentialApp(func(c *cli.Command) error {
			if err := applyStdinCredential(c.Root()); err != nil {
				return err
			}
			got.apiKey, got.apiKeySet = directCredential("api-key")
			got.authToken, got.tokSet = directCredential("auth-token")
			return nil
		})
		require.NoError(t, app.Run(context.Background(), append([]string{"ant"}, argv...)))
		return got
	}
	clearEnv(t, "ANTHROPIC_API_KEY")
	clearEnv(t, "ANTHROPIC_AUTH_TOKEN")
	t.Run("nothing configured", func(t *testing.T) {
		assert.Equal(t, result{}, run(t, "", "probe"))
	})
	t.Run("env var", func(t *testing.T) {
		t.Setenv("ANTHROPIC_AUTH_TOKEN", "tok-env")
		assert.Equal(t, result{authToken: "tok-env", tokSet: true}, run(t, "", "probe"))
	})
	t.Run("stdin beats the env var for the same credential", func(t *testing.T) {
		t.Setenv("ANTHROPIC_API_KEY", "sk-env")
		assert.Equal(t, result{apiKey: "sk-stdin", apiKeySet: true}, run(t, "sk-stdin\n", "probe", "--api-key-stdin"))
	})
	t.Run("stdin for one credential leaves the other's env var in place", func(t *testing.T) {
		t.Setenv("ANTHROPIC_AUTH_TOKEN", "tok-env")
		assert.Equal(t, result{apiKey: "sk-stdin", apiKeySet: true, authToken: "tok-env", tokSet: true},
			run(t, "sk-stdin", "--api-key-stdin", "probe"))
	})
}

// TestRootTakesNoCredentialOnCommandLine: directCredential never reads these
// flags, so a root that declared them would accept a credential and ignore it.
func TestRootTakesNoCredentialOnCommandLine(t *testing.T) {
	for _, f := range Command.Flags {
		assert.NotContains(t, f.Names(), "api-key")
		assert.NotContains(t, f.Names(), "auth-token")
	}
}

// TestCredentialSourceLabel: diagnostics name the env var, or the stdin flag
// when the credential was piped in.
func TestCredentialSourceLabel(t *testing.T) {
	label := func(target string, argv ...string) (got string) {
		app := stdinCredentialApp(func(c *cli.Command) error {
			got = credentialSourceLabel(c.Root(), target)
			return nil
		})
		require.NoError(t, app.Run(context.Background(), append([]string{"ant"}, argv...)))
		return got
	}
	assert.Equal(t, "ANTHROPIC_API_KEY", label("api-key", "probe"))
	assert.Equal(t, "ANTHROPIC_AUTH_TOKEN", label("auth-token", "probe"))
	assert.Equal(t, "--auth-token-stdin", label("auth-token", "--auth-token-stdin", "probe"))
}

// TestAuthStatusStdinCredential: with a logged-in profile, a key piped via
// --api-key-stdin is reported as such — not as ANTHROPIC_API_KEY,
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
	withStdinCredential(t, "sk-ant-api03-FROMSTDIN\n")

	out, err := runStatus(t, "--api-key-stdin")
	require.NoError(t, err)
	assert.Contains(t, out, "(active) * --api-key-stdin")
	assert.NotContains(t, out, "ANTHROPIC_API_KEY")
	assert.NotContains(t, out, "is set in your environment")
	assert.NotContains(t, out, "FROMSTDIN")
}

// TestAuthStatusEnvOverrideWarning: with a logged-in profile, the "is set in
// your environment … unset it" advice is printed for a credential from its env
// var, and not for one piped in, where it would be false.
func TestAuthStatusEnvOverrideWarning(t *testing.T) {
	status := func(t *testing.T, argv ...string) string {
		t.Helper()
		activeProfileWork(t, os.Getenv("ANTHROPIC_CONFIG_DIR"))
		withStdinCredential(t, "secret-stdin\n")
		out, err := runStatusKeepingCredentialEnv(t, argv...)
		require.NoError(t, err)
		return out
	}

	for _, envVar := range []string{"ANTHROPIC_API_KEY", "ANTHROPIC_AUTH_TOKEN"} {
		t.Run(envVar+" warns", func(t *testing.T) {
			clearCredentialEnv(t)
			t.Setenv(envVar, "secret-env")

			out := status(t)

			assert.Contains(t, out, envVar+" is set in your environment and overrides the logged-in profile.")
			assert.Contains(t, out, "unset "+envVar)
		})
	}

	for _, flag := range []string{"--api-key-stdin", "--auth-token-stdin"} {
		t.Run(flag+" does not warn", func(t *testing.T) {
			clearCredentialEnv(t)

			out := status(t, flag)

			assert.NotContains(t, out, "is set in your environment")
		})
	}
}

// TestMultiAuthWarningNamesStdin: the multi-source notice names the stdin
// flag as the winner rather than the env var.
func TestMultiAuthWarningNamesStdin(t *testing.T) {
	resetWarnOnce(t)
	out := captureStderr(t, func() {
		warnIfMultipleAuthSources("--api-key-stdin", "", false, false, false, false, false, true, "", "")
	})
	assert.Contains(t, out, "using --api-key-stdin per precedence")
	assert.NotContains(t, out, "ANTHROPIC_API_KEY")
}

// TestAPIKeyStdinEndToEnd runs the real binary: the key piped on stdin is the
// one sent in the Authorization header, --debug does not print it, and a
// piped key does not get parsed as a request body. Subprocess for the same
// reason as TestWorkspaceIDHeader.
func TestAPIKeyStdinEndToEnd(t *testing.T) {
	var gotKey string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotKey, _ = strings.CutPrefix(r.Header.Get("Authorization"), "Bearer ")
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"data":[],"has_more":false,"first_id":null,"last_id":null}`)
	}))
	defer srv.Close()

	cmd := exec.Command("go", "run", "./cmd/ant", "--debug", "--api-key-stdin", "--base-url", srv.URL, "models", "list")
	cmd.Dir = "../.."
	cmd.Env = append(cmd.Environ(), "ANTHROPIC_CONFIG_DIR="+t.TempDir(), "ANTHROPIC_API_KEY=sk-from-env")
	cmd.Stdin = strings.NewReader("sk-from-stdin\n")
	out, err := cmd.CombinedOutput()
	require.NoError(t, err, "%s", out)
	assert.Equal(t, "sk-from-stdin", gotKey)
	assert.Contains(t, string(out), "Authorization: Bearer <REDACTED>")
	assert.NotContains(t, string(out), "sk-from-stdin")

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
