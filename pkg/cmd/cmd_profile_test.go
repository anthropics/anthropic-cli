package cmd

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/anthropics/anthropic-sdk-go/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli/v3"
)

// captureStdout swaps os.Stdout for a pipe while fn runs and returns whatever
// fn wrote. Profile commands print to stdout directly; this lets tests assert
// on that output without restructuring the handlers.
func captureStdout(t *testing.T, fn func() error) (string, error) {
	t.Helper()
	r, w, _ := os.Pipe()
	old := os.Stdout
	os.Stdout = w
	defer func() { os.Stdout = old }()

	runErr := fn()
	_ = w.Close()
	out, _ := io.ReadAll(r)
	return string(out), runErr
}

func profileCmd() *cli.Command {
	return &cli.Command{
		Name: "profile",
		Commands: []*cli.Command{
			{Name: "activate", ArgsUsage: "<name>", Action: profileActivate},
			{Name: "list", Action: profileList},
			{Name: "get", ArgsUsage: "[<key>]", Action: profileGet},
			{Name: "set", ArgsUsage: "<k> <v>", Action: profileSet},
		},
	}
}

func TestProfileActivate(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("ANTHROPIC_CONFIG_DIR", dir)
	require.NoError(t, config.SaveProfile(dir, "work", &config.Config{AuthenticationInfo: &config.AuthenticationInfo{Type: config.AuthenticationTypeUserOAuth, UserOAuth: &config.UserOAuth{}}}))

	t.Run("writes active_config", func(t *testing.T) {
		require.NoError(t, run(t, profileCmd(), "profile", "activate", "work"))
		assert.Equal(t, "work", strings.TrimSpace(string(mustRead(t, config.ActiveConfigPath(dir)))))
	})

	t.Run("rejects nonexistent profile with login hint", func(t *testing.T) {
		err := run(t, profileCmd(), "profile", "activate", "nope")
		require.Error(t, err)
		assert.Contains(t, err.Error(), `profile "nope" doesn't exist`)
		assert.Contains(t, err.Error(), "ant auth login --profile nope")
		// Must not have clobbered active_config on a failed activate.
		assert.Equal(t, "work", strings.TrimSpace(string(mustRead(t, config.ActiveConfigPath(dir)))))
	})

	t.Run("requires exactly one arg", func(t *testing.T) {
		require.Error(t, run(t, profileCmd(), "profile", "activate"))
		require.Error(t, run(t, profileCmd(), "profile", "activate", "a", "b"))
	})
}

func TestProfileList(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("ANTHROPIC_CONFIG_DIR", dir)
	clearEnv(t, "ANTHROPIC_PROFILE")

	require.NoError(t, config.SaveProfile(dir, "default", &config.Config{AuthenticationInfo: &config.AuthenticationInfo{Type: config.AuthenticationTypeUserOAuth, UserOAuth: &config.UserOAuth{}}}))
	require.NoError(t, config.SaveProfile(dir, "staging", &config.Config{AuthenticationInfo: &config.AuthenticationInfo{Type: config.AuthenticationTypeUserOAuth, UserOAuth: &config.UserOAuth{}}}))
	require.NoError(t, config.SetActiveProfile(dir, "staging"))
	// Non-.json file in the profiles dir must be ignored.
	require.NoError(t, os.WriteFile(filepath.Join(config.ProfilesDir(dir), "README"), []byte("x"), 0o644))

	out, err := captureStdout(t, func() error { return run(t, profileCmd(), "profile", "list") })
	require.NoError(t, err)
	assert.Equal(t, "  default\n* staging\n", out)
}

// TestProfileListEmpty doubles as the missing-dir case: a fresh tempdir has
// no configs/ subdirectory at all. config.ListProfiles handles ErrNotExist
// internally (returns nil, nil), so profileList renders the friendly hint
// rather than erroring.
func TestProfileListEmpty(t *testing.T) {
	t.Setenv("ANTHROPIC_CONFIG_DIR", t.TempDir())
	out, err := captureStdout(t, func() error { return run(t, profileCmd(), "profile", "list") })
	require.NoError(t, err)
	// "(no profiles ...)" goes to stderr; stdout is empty.
	assert.Empty(t, out)
}

func TestProfileGet(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("ANTHROPIC_CONFIG_DIR", dir)
	clearEnv(t, "ANTHROPIC_PROFILE")
	require.NoError(t, config.SaveProfile(dir, "default", &config.Config{
		AuthenticationInfo: &config.AuthenticationInfo{Type: config.AuthenticationTypeUserOAuth, UserOAuth: &config.UserOAuth{}},
		BaseURL:            "https://x",
		WorkspaceID:        "wrkspc_01G",
	}))
	require.NoError(t, config.SetActiveProfile(dir, "default"))

	t.Run("no key dumps full config as JSON", func(t *testing.T) {
		out, err := captureStdout(t, func() error { return run(t, profileCmd(), "profile", "get") })
		require.NoError(t, err)
		var m map[string]any
		require.NoError(t, json.Unmarshal([]byte(out), &m))
		assert.Equal(t, "wrkspc_01G", m["workspace_id"])
		auth, ok := m["authentication"].(map[string]any)
		require.True(t, ok, "authentication block present")
		assert.Equal(t, "user_oauth", auth["type"])
	})

	t.Run("single key prints value", func(t *testing.T) {
		out, err := captureStdout(t, func() error { return run(t, profileCmd(), "profile", "get", "base_url") })
		require.NoError(t, err)
		assert.Equal(t, "https://x\n", out)
	})

	t.Run("nested key prints pretty JSON", func(t *testing.T) {
		out, err := captureStdout(t, func() error { return run(t, profileCmd(), "profile", "get", "authentication") })
		require.NoError(t, err)
		var m map[string]any
		require.NoError(t, json.Unmarshal([]byte(out), &m), "authentication block must be valid JSON, got %q", out)
		assert.Equal(t, "user_oauth", m["type"])
	})

	t.Run("known but unset key prints empty and exits 0", func(t *testing.T) {
		out, err := captureStdout(t, func() error { return run(t, profileCmd(), "profile", "get", "scope") })
		require.NoError(t, err)
		assert.Equal(t, "\n", out)
	})

	t.Run("unknown key errors", func(t *testing.T) {
		_, err := captureStdout(t, func() error { return run(t, profileCmd(), "profile", "get", "nope") })
		require.Error(t, err)
		assert.Contains(t, err.Error(), `no such field "nope"`)
	})
}

func TestProfileSetValidation(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("ANTHROPIC_CONFIG_DIR", dir)
	clearEnv(t, "ANTHROPIC_PROFILE")
	require.NoError(t, config.SaveProfile(dir, "default", &config.Config{AuthenticationInfo: &config.AuthenticationInfo{Type: config.AuthenticationTypeUserOAuth, UserOAuth: &config.UserOAuth{}}}))
	require.NoError(t, config.SetActiveProfile(dir, "default"))

	t.Run("rejects unsupported key", func(t *testing.T) {
		err := run(t, profileCmd(), "profile", "set", "nope", "x")
		require.Error(t, err)
		assert.Contains(t, err.Error(), `unsupported key "nope"`)
	})

	t.Run("requires exactly two args", func(t *testing.T) {
		require.Error(t, run(t, profileCmd(), "profile", "set", "workspace_id"))
	})

	t.Run("each supported key persists", func(t *testing.T) {
		for _, tc := range []struct{ key, val string }{
			{"workspace_id", "wrkspc_01S"},
			{"base_url", "https://s"},
			{"organization_id", "00000000-0000-0000-0000-000000000000"},
			{"scope", "user:inference"},
			{"client_id", "client-x"},
			{"console_url", "https://console.example/"},
		} {
			require.NoError(t, run(t, profileCmd(), "profile", "set", tc.key, tc.val))
		}
		var m map[string]any
		require.NoError(t, json.Unmarshal(mustRead(t, config.ProfilePath(dir, "default")), &m))
		assert.Equal(t, "wrkspc_01S", m["workspace_id"])
		assert.Equal(t, "https://s", m["base_url"])
		assert.Equal(t, "00000000-0000-0000-0000-000000000000", m["organization_id"])
		auth, ok := m["authentication"].(map[string]any)
		require.True(t, ok, "authentication block present")
		assert.Equal(t, "user_oauth", auth["type"])
		assert.Equal(t, "user:inference", auth["scope"])
		assert.Equal(t, "client-x", auth["client_id"])
		assert.Equal(t, "https://console.example", auth["console_url"], "trailing slash trimmed on save")
	})
}
