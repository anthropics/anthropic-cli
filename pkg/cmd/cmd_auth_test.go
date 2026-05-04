package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/anthropics/anthropic-sdk-go/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli/v3"
)

// TestWriteCallbackPageEscapesHTML guards the loopback success/error page
// against markup injection: the OAuth server controls error/error_description
// and we echo them back into the browser response.
func TestWriteCallbackPageEscapesHTML(t *testing.T) {
	rec := httptest.NewRecorder()
	writeCallbackPage(rec, http.StatusBadRequest,
		"<script>alert(1)</script>",
		`oops "&' <img src=x onerror=alert(1)>`)
	body := rec.Body.String()
	assert.Equal(t, "text/html; charset=utf-8", rec.Header().Get("Content-Type"))
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.NotContains(t, body, "<script>")
	assert.NotContains(t, body, "<img")
	assert.Contains(t, body, "&lt;script&gt;alert(1)&lt;/script&gt;")
	assert.Contains(t, body, "&lt;img src=x onerror=alert(1)&gt;")
}

// TestLoginMode covers the three-way flag dispatch for `ant auth login`:
// default (browser + loopback), SSH-forward (--no-browser + fixed port),
// and headless (--no-browser alone, copy-paste via stdin).
func TestLoginMode(t *testing.T) {
	newCmd := func() *cli.Command {
		return &cli.Command{
			Name:   "login",
			Action: func(context.Context, *cli.Command) error { return nil },
			Flags: []cli.Flag{
				&cli.BoolFlag{Name: "no-browser"},
				&cli.IntFlag{Name: "callback-port", Value: 0},
			},
		}
	}

	cases := []struct {
		name string
		argv []string
		want loginModeEnum
	}{
		{"default no flags", nil, loginModeDefault},
		{"default with port but no --no-browser", []string{"--callback-port=9000"}, loginModeDefault},
		{"headless with --no-browser alone", []string{"--no-browser"}, loginModeHeadless},
		{"ssh-forward with --no-browser and port", []string{"--no-browser", "--callback-port=9000"}, loginModeSSHForward},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			cmd := newCmd()
			var got loginModeEnum
			cmd.Action = func(_ context.Context, c *cli.Command) error {
				got = loginMode(c)
				return nil
			}
			argv := append([]string{"login"}, tc.argv...)
			require.NoError(t, cmd.Run(context.Background(), argv))
			assert.Equal(t, tc.want, got)
		})
	}
}

// TestManualRedirectURIForAuthorize covers --no-browser redirect_uri
// derivation: valid console URLs map to
// <consoleURL>/oauth/code/callback?app=anthropic-cli; unparseable / hostless
// inputs return "". The ?app=anthropic-cli query param is part of the
// registered redirect_uri on the OAuth client and must be preserved verbatim
// — the authorize call is rejected without it.
func TestManualRedirectURIForAuthorize(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		{"https://platform.claude.com", "https://platform.claude.com/oauth/code/callback?app=anthropic-cli"},
		{"https://platform.claude.com/", "https://platform.claude.com/oauth/code/callback?app=anthropic-cli"},
		{"https://console.example.test", "https://console.example.test/oauth/code/callback?app=anthropic-cli"},
		// Unknown hosts are derived too — server-side redirect_uri validation
		// against the OAuth client's registered list is the rejection point.
		{"https://unknown.example.com", "https://unknown.example.com/oauth/code/callback?app=anthropic-cli"},
		{"not a url", ""},
		{"", ""},
	}
	for _, tc := range cases {
		assert.Equal(t, tc.want, manualRedirectURIForAuthorize(tc.in), "input=%q", tc.in)
	}
}

// TestReadManualCode covers the "<code>#<state>" parse: Console's
// /oauth/code/callback renders the code and state as a single pasteable
// string, so readManualCode must split on '#', verify the state when
// present, and accept a bare code for forward-compat.
func TestReadManualCode(t *testing.T) {
	const wantState = "expected-state"
	cases := []struct {
		name, in, wantCode, wantErr string
	}{
		{"code with matching state", "abc123#" + wantState + "\n", "abc123", ""},
		{"bare code (no state suffix)", "abc123\n", "abc123", ""},
		{"surrounding whitespace trimmed", "  abc123#" + wantState + "  \n", "abc123", ""},
		{"state mismatch rejected", "abc123#wrong\n", "", "pasted state does not match"},
		{"empty input rejected", "\n", "", "empty code pasted"},
		{"state only (no code) rejected", "#" + wantState + "\n", "", "empty code pasted"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := readManualCode(context.Background(), io.Discard, strings.NewReader(tc.in), wantState)
			if tc.wantErr != "" {
				require.ErrorContains(t, err, tc.wantErr)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.wantCode, got)
		})
	}
}

// clearEnv unsets key for the duration of the test with t.Setenv-style
// restore — t.Setenv alone would leave the var set to "", which os.LookupEnv
// reports as present.
func clearEnv(t *testing.T, key string) {
	t.Helper()
	t.Setenv(key, "")
	os.Unsetenv(key)
}

// run builds a one-off urfave/cli command and runs it with argv, so tests
// exercise the same flag parsing the binary uses. The wrapping app declares
// the global --profile flag (mirroring extras.go) so subcommands resolve it
// the same way production does.
func run(t *testing.T, defn *cli.Command, argv ...string) error {
	t.Helper()
	app := &cli.Command{
		Name:     "ant",
		Flags:    []cli.Flag{&cli.StringFlag{Name: "profile", Sources: cli.EnvVars("ANTHROPIC_PROFILE")}},
		Commands: []*cli.Command{defn},
	}
	return app.Run(context.Background(), append([]string{"ant"}, argv...))
}

// TestExchangeCodeRequestShape guards the wire format of the
// authorization_code → token exchange. The spec recommends JSON +
// anthropic-beta on every token POST, but doing so routes to api-go's
// jwt-bearer-only handler in prod and 400s. This test fails if someone
// re-applies that recommendation.
func TestExchangeCodeRequestShape(t *testing.T) {
	var got struct {
		path, contentType, beta string
		form                    url.Values
	}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got.path = r.URL.Path
		got.contentType = r.Header.Get("Content-Type")
		got.beta = r.Header.Get("Anthropic-Beta")
		_ = r.ParseForm()
		got.form = r.PostForm
		_ = json.NewEncoder(w).Encode(tokenResponse{
			AccessToken: "sk-ant-oat01-TEST", RefreshToken: "rt-TEST",
			ExpiresIn: 600, TokenType: "Bearer", Scope: "user:inference",
		})
	}))
	defer srv.Close()

	tok, err := exchangeCode(context.Background(), srv.URL, "client-x", "code-x", "verifier-x", "http://localhost:1/cb", "state-x", false)
	require.NoError(t, err)

	assert.Equal(t, "/v1/oauth/token", got.path)
	assert.Equal(t, "application/x-www-form-urlencoded", got.contentType)
	assert.Empty(t, got.beta, "must not send Anthropic-Beta on authorization_code grant (routes to api-go jwt-bearer handler)")
	assert.Equal(t, url.Values{
		"grant_type":    {"authorization_code"},
		"code":          {"code-x"},
		"code_verifier": {"verifier-x"},
		"client_id":     {"client-x"},
		"redirect_uri":  {"http://localhost:1/cb"},
		"state":         {"state-x"},
	}, got.form)
	// workspace_id (and the legacy workspace_uuid name) belong on
	// /oauth/authorize, not /oauth/token — the authorization carries the
	// workspace binding; the token endpoint just cashes in a pre-bound code.
	for _, k := range []string{"workspace_id", "workspace_uuid"} {
		_, has := got.form[k]
		assert.False(t, has, "%s must not be sent on /oauth/token; it belongs on /oauth/authorize", k)
	}
	assert.Equal(t, "sk-ant-oat01-TEST", tok.AccessToken)
	assert.Equal(t, "user:inference", tok.Scope)
}

// A 200 response with an empty access_token would otherwise be written to
// disk silently; the guard fails loud on that shape.
func TestExchangeCodeRejectsEmptyAccessToken(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Request-Id", "req_test_empty")
		_ = json.NewEncoder(w).Encode(tokenResponse{
			AccessToken: "", RefreshToken: "rt-TEST", ExpiresIn: 600, TokenType: "Bearer",
		})
	}))
	defer srv.Close()

	tok, err := exchangeCode(context.Background(), srv.URL, "client-x", "code-x", "verifier-x", "http://localhost:1/cb", "state-x", false)
	require.Error(t, err)
	assert.Nil(t, tok)
	assert.Contains(t, err.Error(), "empty access_token")
	assert.Contains(t, err.Error(), "req_test_empty")
}

func TestActiveProfileResolution(t *testing.T) {
	withFile := t.TempDir()
	require.NoError(t, os.WriteFile(config.ActiveConfigPath(withFile), []byte("from-file\n"), 0o600))

	for _, tc := range []struct {
		name, dir, flag, env, want string
	}{
		{"default when nothing set and no file", t.TempDir(), "", "", "default"},
		{"file when present", withFile, "", "", "from-file"},
		{"env beats file", withFile, "", "from-env", "from-env"},
		{"flag beats env", withFile, "from-flag", "from-env", "from-flag"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv("ANTHROPIC_CONFIG_DIR", tc.dir)
			if tc.env != "" {
				t.Setenv("ANTHROPIC_PROFILE", tc.env)
			} else {
				clearEnv(t, "ANTHROPIC_PROFILE")
			}
			cmd := &cli.Command{Flags: []cli.Flag{&cli.StringFlag{Name: "profile"}}}
			if tc.flag != "" {
				_ = cmd.Set("profile", tc.flag)
			}
			got, _ := activeProfile(cmd)
			assert.Equal(t, tc.want, got)
		})
	}
}

// TestLoadProfileConfigBypassesEnv guards the regression the old
// os.Setenv("ANTHROPIC_PROFILE", …) hack existed to work around: when the
// CLI asks for a specific profile (e.g. via --profile), loadProfileConfig
// must load that profile even if ANTHROPIC_PROFILE points elsewhere.
func TestLoadProfileConfigBypassesEnv(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("ANTHROPIC_CONFIG_DIR", dir)
	require.NoError(t, config.SaveProfile(dir, "target", &config.Config{
		AuthenticationInfo: &config.AuthenticationInfo{Type: config.AuthenticationTypeUserOAuth, UserOAuth: &config.UserOAuth{}},
		BaseURL:            "https://target.example",
	}))
	require.NoError(t, config.SaveProfile(dir, "decoy", &config.Config{
		AuthenticationInfo: &config.AuthenticationInfo{Type: config.AuthenticationTypeUserOAuth, UserOAuth: &config.UserOAuth{}},
		BaseURL:            "https://decoy.example",
	}))
	t.Setenv("ANTHROPIC_PROFILE", "decoy")

	cfg, err := loadProfileConfig(dir, "target")
	require.NoError(t, err)
	assert.Equal(t, "https://target.example", cfg.BaseURL,
		"loadProfileConfig must honour the explicit profile arg, not ANTHROPIC_PROFILE")

	// Empty profile name is rejected by the SDK's validateProfileName — never
	// reachable from CLI callers (activeProfile always returns at least
	// "default"), but the contract is "explicit name only".
	_, err = loadProfileConfig(dir, "")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "empty")
}

// TestResolveRequestedScope covers the --scope precedence: flag → profile's
// stored UserOAuth.Scope → oauthScope default.
func TestResolveRequestedScope(t *testing.T) {
	withScope := func(s string) *config.Config {
		return &config.Config{AuthenticationInfo: &config.AuthenticationInfo{
			Type: config.AuthenticationTypeUserOAuth, UserOAuth: &config.UserOAuth{Scope: s},
		}}
	}
	cases := []struct {
		name, flag string
		prev       *config.Config
		want       string
	}{
		{"flag wins over profile", "user:profile", withScope("user:inference"), "user:profile"},
		{"profile scope when no flag", "", withScope("user:profile"), "user:profile"},
		{"default when no flag and empty profile scope", "", withScope(""), oauthScope},
		{"default when no flag and nil profile", "", nil, oauthScope},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, resolveRequestedScope(tc.flag, tc.prev))
		})
	}
}

func TestResolveClientID(t *testing.T) {
	withClientID := func(id string) *config.Config {
		return &config.Config{AuthenticationInfo: &config.AuthenticationInfo{
			Type: config.AuthenticationTypeUserOAuth, UserOAuth: &config.UserOAuth{ClientID: id},
		}}
	}
	cases := []struct {
		name, flag string
		prev       *config.Config
		want       string
	}{
		{"flag wins over profile", "from-flag", withClientID("from-profile"), "from-flag"},
		{"profile when no flag", "", withClientID("from-profile"), "from-profile"},
		{"prod default when no flag and empty profile client_id", "", withClientID(""), oauthClientIDProd},
		{"prod default when no flag and nil profile", "", nil, oauthClientIDProd},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, resolveClientID(tc.flag, tc.prev))
		})
	}
}

func TestResolveConsoleURL(t *testing.T) {
	withConsoleURL := func(u string) *config.Config {
		return &config.Config{AuthenticationInfo: &config.AuthenticationInfo{
			Type: config.AuthenticationTypeUserOAuth, UserOAuth: &config.UserOAuth{ConsoleURL: u},
		}}
	}
	cases := []struct {
		name, flag string
		prev       *config.Config
		want       string
	}{
		{"flag wins over profile", "https://from-flag.example", withConsoleURL("https://from-profile.example"), "https://from-flag.example"},
		{"flag trailing slash trimmed", "https://from-flag.example/", nil, "https://from-flag.example"},
		{"profile when no flag", "", withConsoleURL("https://from-profile.example"), "https://from-profile.example"},
		{"profile trailing slash trimmed", "", withConsoleURL("https://from-profile.example/"), "https://from-profile.example"},
		{"default when no flag and nil profile", "", nil, defaultConsoleURL},
		{"default when no flag and empty profile console_url", "", withConsoleURL(""), defaultConsoleURL},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, resolveConsoleURL(tc.flag, tc.prev))
		})
	}
}

func TestResolveBaseURL(t *testing.T) {
	with := func(u string) *config.Config { return &config.Config{BaseURL: u} }
	cases := []struct {
		name, flag string
		prev       *config.Config
		want       string
	}{
		{"flag wins over profile", "https://from-flag.example", with("https://from-profile.example"), "https://from-flag.example"},
		{"flag trailing slash trimmed", "https://from-flag.example/", nil, "https://from-flag.example"},
		{"profile when no flag", "", with("https://from-profile.example"), "https://from-profile.example"},
		{"default when no flag and nil profile", "", nil, defaultBaseURL},
		{"default when no flag and empty profile base_url", "", with(""), defaultBaseURL},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, resolveBaseURL(tc.flag, tc.prev))
		})
	}
}

// captureStderr swaps os.Stderr while fn runs and returns whatever was written.
func captureStderr(t *testing.T, fn func()) string {
	t.Helper()
	r, w, err := os.Pipe()
	require.NoError(t, err)
	old := os.Stderr
	os.Stderr = w
	defer func() { os.Stderr = old }()
	fn()
	_ = w.Close()
	out, _ := io.ReadAll(r)
	return string(out)
}

// resetWarnOnce resets the package-level one-shot warning guards. Registered
// as a t.Cleanup so subtests don't bleed state regardless of order.
func resetWarnOnce(t *testing.T) {
	t.Helper()
	multiAuthWarnOnce = sync.Once{}
	clientIDDefaultedOnce = sync.Once{}
	t.Cleanup(func() {
		multiAuthWarnOnce = sync.Once{}
		clientIDDefaultedOnce = sync.Once{}
	})
}

// TestMultiAuthWarning covers the one-shot stderr notice when multiple
// credential sources are configured: it lists the sources, names the
// precedence winner, and only fires once.
func TestMultiAuthWarning(t *testing.T) {
	reset := func() { resetWarnOnce(t) }

	t.Run("api-key and explicit profile", func(t *testing.T) {
		reset()
		out := captureStderr(t, func() { warnIfMultipleAuthSources(true, false, true, false, false) })
		assert.Contains(t, out, "multiple auth sources configured")
		assert.Contains(t, out, "--api-key / ANTHROPIC_API_KEY")
		assert.Contains(t, out, "profile from --profile / ANTHROPIC_PROFILE")
		assert.Contains(t, out, "using --api-key / ANTHROPIC_API_KEY per precedence")
		assert.Contains(t, out, "ant auth status")
	})

	t.Run("federation beats implicit profile", func(t *testing.T) {
		reset()
		out := captureStderr(t, func() { warnIfMultipleAuthSources(false, false, false, true, true) })
		assert.Contains(t, out, "federation env")
		assert.Contains(t, out, "active profile (active_config)")
		assert.Contains(t, out, "using federation env per precedence")
	})

	t.Run("explicit profile beats federation", func(t *testing.T) {
		reset()
		out := captureStderr(t, func() { warnIfMultipleAuthSources(false, false, true, true, false) })
		assert.Contains(t, out, "using profile from --profile / ANTHROPIC_PROFILE per precedence")
	})

	t.Run("single source is silent", func(t *testing.T) {
		reset()
		out := captureStderr(t, func() { warnIfMultipleAuthSources(true, false, false, false, false) })
		assert.Empty(t, out)
	})

	t.Run("emits once", func(t *testing.T) {
		reset()
		first := captureStderr(t, func() { warnIfMultipleAuthSources(false, true, true, true, false) })
		second := captureStderr(t, func() { warnIfMultipleAuthSources(false, true, true, true, false) })
		assert.NotEmpty(t, first)
		assert.Empty(t, second)
	})
}

// TestCredentialPrecedenceFederationVsProfile verifies that federation env
// vars beat an implicitly-resolved profile (active_config) but lose to an
// explicitly-named one (ANTHROPIC_PROFILE), per the User Guide's 5-tier
// precedence.
func TestCredentialPrecedenceFederationVsProfile(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("ANTHROPIC_CONFIG_DIR", dir)
	clearEnv(t, "ANTHROPIC_API_KEY")
	clearEnv(t, "ANTHROPIC_AUTH_TOKEN")
	clearEnv(t, "ANTHROPIC_PROFILE")

	// Seed a usable user_oauth profile (config + creds) and make it active.
	require.NoError(t, config.SaveProfile(dir, "p", &config.Config{
		AuthenticationInfo: &config.AuthenticationInfo{
			Type: config.AuthenticationTypeUserOAuth, UserOAuth: &config.UserOAuth{},
		},
	}))
	require.NoError(t, config.WriteCredentials(config.ProfileCredentialsPath(dir, "p"),
		config.Credentials{AccessToken: "tok"}))
	require.NoError(t, config.SetActiveProfile(dir, "p"))

	root := &cli.Command{Flags: []cli.Flag{
		&cli.StringFlag{Name: "api-key"}, &cli.StringFlag{Name: "auth-token"},
		&cli.StringFlag{Name: "identity-token"}, &cli.StringFlag{Name: "identity-token-file"},
		&cli.StringFlag{Name: "federation-rule"}, &cli.StringFlag{Name: "organization-id"},
		&cli.StringFlag{Name: "service-account-id"},
		&cli.StringFlag{Name: "profile"},
	}}

	t.Run("implicit profile, no federation → profile usable", func(t *testing.T) {
		cfg, explicit := loadProfileIfUsable(root)
		assert.NotNil(t, cfg)
		assert.False(t, explicit)
	})

	t.Run("ANTHROPIC_PROFILE set → explicit", func(t *testing.T) {
		t.Setenv("ANTHROPIC_PROFILE", "p")
		cfg, explicit := loadProfileIfUsable(root)
		assert.NotNil(t, cfg)
		assert.True(t, explicit)
	})
}

// TestGlobalProfileFlag verifies that the global --profile flag (extras.go)
// reaches loadProfileIfUsable from any subcommand position: with
// active_config=a and `--profile b`, the loader returns b's config and
// reports it as explicitly named.
func TestGlobalProfileFlag(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("ANTHROPIC_CONFIG_DIR", dir)
	clearEnv(t, "ANTHROPIC_PROFILE")
	clearEnv(t, "ANTHROPIC_API_KEY")
	clearEnv(t, "ANTHROPIC_AUTH_TOKEN")

	for _, name := range []string{"a", "b"} {
		require.NoError(t, config.SaveProfile(dir, name, &config.Config{
			AuthenticationInfo: &config.AuthenticationInfo{
				Type: config.AuthenticationTypeUserOAuth, UserOAuth: &config.UserOAuth{ClientID: "client-" + name},
			},
		}))
		require.NoError(t, config.WriteCredentials(config.ProfileCredentialsPath(dir, name),
			config.Credentials{AccessToken: "tok-" + name}))
	}
	require.NoError(t, config.SetActiveProfile(dir, "a"))

	var cfg *config.Config
	var explicit bool
	probe := &cli.Command{
		Name: "probe",
		Action: func(_ context.Context, c *cli.Command) error {
			cfg, explicit = loadProfileIfUsable(c)
			return nil
		},
	}

	// No --profile: implicit resolution to active_config (= "a").
	require.NoError(t, run(t, probe, "probe"))
	require.NotNil(t, cfg)
	assert.False(t, explicit)
	assert.Equal(t, "client-a", cfg.AuthenticationInfo.UserOAuth.ClientID)

	// --profile b after the subcommand name: global flag reaches the action,
	// loadProfileIfUsable returns b's config, explicit=true.
	require.NoError(t, run(t, probe, "probe", "--profile", "b"))
	require.NotNil(t, cfg)
	assert.True(t, explicit)
	assert.Equal(t, "client-b", cfg.AuthenticationInfo.UserOAuth.ClientID)
}

// TestResolveOAuthOption_Federation verifies that fully-configured federation
// flags produce a single SDK request option (option.WithFederationTokenProvider
// → auth middleware), partial config errors with the missing inputs, and an
// empty Federation falls through with (nil, nil).
func TestResolveOAuthOption_Federation(t *testing.T) {
	t.Run("complete federation → one option", func(t *testing.T) {
		fed := federation{Assertion: "jwt", Rule: "fdrl_x", OrganizationID: "org-x", ServiceAccountID: "svac_x"}
		opts, err := resolveOAuthOption(fed)
		require.NoError(t, err)
		require.GreaterOrEqual(t, len(opts), 1, "federation path must produce the provider option; additional options may be added")
	})

	t.Run("partial federation → error listing missing", func(t *testing.T) {
		fed := federation{Rule: "fdrl_x"}
		_, err := resolveOAuthOption(fed)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "ANTHROPIC_IDENTITY_TOKEN")
		assert.Contains(t, err.Error(), "ANTHROPIC_ORGANIZATION_ID")
	})

	t.Run("no source configured → nil, nil", func(t *testing.T) {
		opts, err := resolveOAuthOption(federation{})
		require.NoError(t, err)
		assert.Nil(t, opts)
	})

	t.Run("assertion file path", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "id.jwt")
		require.NoError(t, os.WriteFile(path, []byte("from-file"), 0o600))
		fed := federation{AssertionFile: path, Rule: "fdrl_x", OrganizationID: "org-x"}
		opts, err := resolveOAuthOption(fed)
		require.NoError(t, err)
		require.GreaterOrEqual(t, len(opts), 1)
	})
}

// TestLoadProfileFillsClientIDDefault verifies that when a user_oauth profile
// config omits client_id (the bootstrap-only-if-set case), the CLI fills in
// oauthClientIDProd at request time so the SDK's refresh path has what it
// needs. A non-empty client_id is left untouched.
func TestLoadProfileFillsClientIDDefault(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("ANTHROPIC_CONFIG_DIR", dir)
	clearEnv(t, "ANTHROPIC_PROFILE")
	clearEnv(t, "ANTHROPIC_API_KEY")
	clearEnv(t, "ANTHROPIC_AUTH_TOKEN")

	seed := func(name, clientID string) {
		require.NoError(t, config.SaveProfile(dir, name, &config.Config{
			AuthenticationInfo: &config.AuthenticationInfo{
				Type: config.AuthenticationTypeUserOAuth, UserOAuth: &config.UserOAuth{ClientID: clientID},
			},
		}))
		require.NoError(t, config.WriteCredentials(config.ProfileCredentialsPath(dir, name),
			config.Credentials{AccessToken: "tok"}))
		require.NoError(t, config.SetActiveProfile(dir, name))
	}

	t.Run("empty client_id gets prod default and warns once", func(t *testing.T) {
		resetWarnOnce(t)
		seed("noclient", "")
		var cfg *config.Config
		out := captureStderr(t, func() { cfg, _ = loadProfileIfUsable(nil) })
		require.NotNil(t, cfg)
		assert.Equal(t, oauthClientIDProd, cfg.AuthenticationInfo.UserOAuth.ClientID)
		assert.Contains(t, out, "missing client_id")
		assert.Contains(t, out, "ant auth login")
		// Second call: still defaults, but no second warning.
		out2 := captureStderr(t, func() { _, _ = loadProfileIfUsable(nil) })
		assert.Empty(t, out2)
	})

	t.Run("explicit client_id preserved without warning", func(t *testing.T) {
		resetWarnOnce(t)
		seed("withclient", "custom-client")
		var cfg *config.Config
		out := captureStderr(t, func() { cfg, _ = loadProfileIfUsable(nil) })
		require.NotNil(t, cfg)
		assert.Equal(t, "custom-client", cfg.AuthenticationInfo.UserOAuth.ClientID)
		assert.Empty(t, out)
	})
}

func TestProfileSetDoesNotPersistDefaultCredentialsPath(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("ANTHROPIC_CONFIG_DIR", dir)
	clearEnv(t, "ANTHROPIC_PROFILE")
	require.NoError(t, config.SaveProfile(dir, "p", &config.Config{AuthenticationInfo: &config.AuthenticationInfo{Type: config.AuthenticationTypeUserOAuth, UserOAuth: &config.UserOAuth{}}}))
	require.NoError(t, config.SetActiveProfile(dir, "p"))

	require.NoError(t, run(t, &cli.Command{
		Name: "profile", Commands: []*cli.Command{{
			Name: "set", ArgsUsage: "<k> <v>", Action: profileSet,
		}},
	}, "profile", "set", "workspace_id", "wrkspc_01X"))

	raw, err := os.ReadFile(config.ProfilePath(dir, "p"))
	require.NoError(t, err)
	var m map[string]any
	require.NoError(t, json.Unmarshal(raw, &m))
	assert.Equal(t, "wrkspc_01X", m["workspace_id"])
	_, hasCredPath := m["credentials_path"]
	assert.False(t, hasCredPath, "profile set must not bake the resolved default credentials_path into the config")
}

// TestAuthLoginWritesSpecFiles runs the full authLogin flow against a mock
// token server: PKCE → callback → exchange → SDK writers. It asserts the
// three-file directory layout and field shapes the SDK readers expect.
// loginCmdDef returns an `auth login` command tree mirroring production flag
// shapes. --profile is global (declared by run()'s wrapping app), so it isn't
// listed here.
func loginCmdDef() *cli.Command {
	return &cli.Command{
		Name: "auth", Commands: []*cli.Command{{
			Name: "login", Action: authLogin,
			Flags: []cli.Flag{
				&cli.StringFlag{Name: "console-url"},
				&cli.StringFlag{Name: "base-url"},
				&cli.BoolFlag{Name: "no-browser"},
				&cli.IntFlag{Name: "callback-port", Value: 0},
				&cli.DurationFlag{Name: "timeout", Value: 30 * time.Second},
				&cli.StringFlag{Name: "client-id"},
				&cli.StringFlag{Name: "scope"},
				&cli.StringFlag{Name: "workspace-id"},
				&cli.BoolFlag{Name: "debug"},
			},
		}},
	}
}

// driveLogin runs `auth login --no-browser --callback-port 0` with extraArgs,
// scrapes the authorize URL from stderr, delivers the loopback callback, and
// waits for completion. Returns the parsed authorize URL so callers can assert
// on its query params.
func driveLogin(t *testing.T, tokenSrvURL string, extraArgs ...string) *url.URL {
	t.Helper()
	u, _, err := driveLoginErr(t, tokenSrvURL, extraArgs...)
	require.NoError(t, err)
	return u
}

// driveLoginWithArgs runs the full authLogin → callback → token-exchange flow
// with the given complete args (must include "auth", "login", --base-url, and
// any flags under test). It captures stderr, scrapes the printed authorize
// URL, fires the loopback callback with ?code=CODE, and returns the parsed
// authorize URL, the full stderr output, and authLogin's final error.
func driveLoginWithArgs(t *testing.T, args []string) (*url.URL, string, error) {
	t.Helper()

	r, w, perr := os.Pipe()
	require.NoError(t, perr)
	oldStderr := os.Stderr
	os.Stderr = w
	t.Cleanup(func() { os.Stderr = oldStderr })

	done := make(chan error, 1)
	go func() { done <- run(t, loginCmdDef(), args...) }()

	// Drain the pipe into a buffer in the background so writes never block,
	// and poll the buffer for the authorize URL. Reading via a single
	// io.Copy (rather than bufio.Scanner then io.Copy) means no bytes are
	// stranded in a scanner's internal buffer, so the returned stderr is
	// the complete output.
	var (
		out   bytes.Buffer
		outMu sync.Mutex
	)
	copyDone := make(chan struct{})
	go func() {
		buf := make([]byte, 4096)
		for {
			n, rerr := r.Read(buf)
			if n > 0 {
				outMu.Lock()
				out.Write(buf[:n])
				outMu.Unlock()
			}
			if rerr != nil {
				close(copyDone)
				return
			}
		}
	}()

	authorizeURLRe := regexp.MustCompile(`https://\S+/oauth/authorize\?\S+`)
	var authorizeURL string
	deadline := time.Now().Add(10 * time.Second)
	for authorizeURL == "" && time.Now().Before(deadline) {
		outMu.Lock()
		authorizeURL = authorizeURLRe.FindString(out.String())
		outMu.Unlock()
		if authorizeURL == "" {
			time.Sleep(5 * time.Millisecond)
		}
	}
	require.NotEmpty(t, authorizeURL, "authLogin should print an authorize URL")
	u, _ := url.Parse(authorizeURL)
	state := u.Query().Get("state")
	redirect := u.Query().Get("redirect_uri")
	// A request to the wrong path is stray local traffic, not the OAuth
	// redirect — it must 404 without consuming the auth attempt.
	ru, _ := url.Parse(redirect)
	stray, err := http.Get(ru.Scheme + "://" + ru.Host + "/wrong?code=CODE&state=" + url.QueryEscape(state))
	require.NoError(t, err)
	stray.Body.Close()
	require.Equal(t, http.StatusNotFound, stray.StatusCode)
	resp, err := http.Get(redirect + "?code=CODE&state=" + url.QueryEscape(state))
	require.NoError(t, err)
	resp.Body.Close()

	select {
	case err := <-done:
		os.Stderr = oldStderr
		_ = w.Close()
		<-copyDone
		outMu.Lock()
		defer outMu.Unlock()
		return u, out.String(), err
	case <-time.After(10 * time.Second):
		t.Fatal("authLogin did not complete")
		panic("unreachable")
	}
}

// driveLoginErr wraps driveLoginWithArgs with the common defaults (including
// --workspace-id wrkspc_test) and asserts the workspace_id query-param
// contract. Tests that need to omit --workspace-id should call
// driveLoginWithArgs directly.
func driveLoginErr(t *testing.T, tokenSrvURL string, extraArgs ...string) (*url.URL, string, error) {
	t.Helper()
	args := append([]string{"auth", "login", "--no-browser", "--callback-port", "0",
		"--base-url", tokenSrvURL, "--workspace-id", "wrkspc_test"}, extraArgs...)
	// If a later --workspace-id was passed in extraArgs it wins (urfave/cli
	// takes the last value for a repeated StringFlag), so compute the expected
	// query value from args rather than hardcoding.
	expectedWorkspaceID := "wrkspc_test"
	for i, a := range args {
		if a == "--workspace-id" && i+1 < len(args) {
			expectedWorkspaceID = args[i+1]
		}
	}
	u, out, err := driveLoginWithArgs(t, args)
	assert.Equal(t, expectedWorkspaceID, u.Query().Get("workspace_id"),
		"workspace_id must be sent on /oauth/authorize when --workspace-id is provided")
	assert.Empty(t, u.Query().Get("workspace_uuid"),
		"legacy workspace_uuid key must not be sent (backend uses workspace_id)")
	return u, out, err
}

// newTokenServer returns an httptest.Server that 200s every request with the
// given token response.
func newTokenServer(t *testing.T, resp tokenResponse) *httptest.Server {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(resp)
	}))
	t.Cleanup(srv.Close)
	return srv
}

func TestAuthLoginWritesSpecFiles(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("ANTHROPIC_CONFIG_DIR", dir)
	clearEnv(t, "ANTHROPIC_PROFILE")

	tokenSrv := newTokenServer(t, tokenResponse{
		AccessToken: "sk-ant-oat01-E2E", RefreshToken: "rt-E2E",
		ExpiresIn: 600, Scope: "user:inference",
		Organization: tokenOrganization{UUID: "org-E2E", Name: "E2E Org"},
		Account:      tokenAccount{EmailAddress: "e2e@example.com"},
		Workspace:    tokenWorkspace{ID: "wrkspc_test", Name: "Test Workspace"},
	})

	u := driveLogin(t, tokenSrv.URL, "--profile", "e2e", "--scope", "user:profile",
		"--client-id", "test-client")
	assert.Equal(t, "user:profile", u.Query().Get("scope"), "--scope must drive the authorize URL scope param")

	// configs/e2e.json — client_id/scope/base_url persisted because all three
	// flags were explicitly set on this bootstrap.
	var cfg map[string]any
	require.NoError(t, json.Unmarshal(mustRead(t, config.ProfilePath(dir, "e2e")), &cfg))
	assert.Equal(t, "1.0", cfg["version"])
	assert.Equal(t, "wrkspc_test", cfg["workspace_id"])
	assert.True(t, strings.HasPrefix(cfg["workspace_id"].(string), "wrkspc_"),
		"config workspace_id must be the tagged form")
	assert.Equal(t, tokenSrv.URL, cfg["base_url"])
	auth, ok := cfg["authentication"].(map[string]any)
	require.True(t, ok, "authentication block present")
	assert.Equal(t, "user_oauth", auth["type"])
	assert.Equal(t, "test-client", auth["client_id"])
	assert.Equal(t, "user:inference", auth["scope"], "persisted scope is the granted value")

	// credentials/e2e.json
	var creds map[string]any
	require.NoError(t, json.Unmarshal(mustRead(t, config.ProfileCredentialsPath(dir, "e2e")), &creds))
	assert.Equal(t, "1.0", creds["version"])
	assert.Equal(t, "oauth_token", creds["type"])
	assert.Equal(t, "sk-ant-oat01-E2E", creds["access_token"])
	assert.Equal(t, "rt-E2E", creds["refresh_token"])
	_, isFloat := creds["expires_at"].(float64) // JSON numbers decode to float64
	assert.True(t, isFloat, "expires_at must be a unix-seconds number, not an RFC3339 string")
	assert.Equal(t, "user:inference", creds["scope"], "granted scope persisted to credentials")
	assert.Equal(t, "org-E2E", creds["organization_uuid"])
	assert.Equal(t, "E2E Org", creds["organization_name"])
	assert.Equal(t, "e2e@example.com", creds["account_email"])
	// Workspace fields persisted from token response in tagged form.
	assert.Equal(t, "wrkspc_test", creds["workspace_id"])
	assert.True(t, strings.HasPrefix(creds["workspace_id"].(string), "wrkspc_"),
		"creds workspace_id must be the tagged form")
	assert.Equal(t, "Test Workspace", creds["workspace_name"])
	_, hasLegacy := creds["workspace_uuid"]
	assert.False(t, hasLegacy, "creds must not write the legacy workspace_uuid key")

	// active_config — written because --profile was given and none existed.
	assert.Equal(t, "e2e", strings.TrimSpace(string(mustRead(t, config.ActiveConfigPath(dir)))))
}

// driveDefaultLogin runs authLogin in default (browser) mode with openBrowser
// and stdin stubbed. The opened channel receives the loopback authorize URL
// when openBrowser is called; pasted (if non-empty) is fed to stdin. Returns
// the form params the token server received on the exchange POST, the full
// stderr, and authLogin's error. The caller is responsible for delivering the
// loopback callback (via opened) when testing the callback-wins path.
func driveDefaultLogin(t *testing.T, opened chan<- string, pasted string) (url.Values, string, error) {
	t.Helper()
	dir := t.TempDir()
	t.Setenv("ANTHROPIC_CONFIG_DIR", dir)
	clearEnv(t, "ANTHROPIC_PROFILE")

	origOpenBrowser := openBrowser
	openBrowser = func(u string) error {
		select {
		case opened <- u:
		default:
		}
		return nil
	}
	t.Cleanup(func() { openBrowser = origOpenBrowser })

	stdinR, stdinW, _ := os.Pipe()
	origStdin := stdin
	stdin = stdinR
	t.Cleanup(func() { stdin = origStdin; stdinR.Close(); stdinW.Close() })
	if pasted != "" {
		go func() { stdinW.WriteString(pasted + "\n") }()
	}

	var exchangeForm url.Values
	tokenSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = r.ParseForm()
		exchangeForm = r.PostForm
		_ = json.NewEncoder(w).Encode(tokenResponse{
			AccessToken: "sk-ant-oat01-BROWSER", RefreshToken: "rt-BROWSER",
			ExpiresIn: 600, Scope: "user:inference",
			Workspace: tokenWorkspace{ID: "wrkspc_test"},
		})
	}))
	t.Cleanup(tokenSrv.Close)

	r, w, _ := os.Pipe()
	oldStderr := os.Stderr
	os.Stderr = w
	t.Cleanup(func() { os.Stderr = oldStderr })
	var out bytes.Buffer
	copyDone := make(chan struct{})
	go func() { io.Copy(&out, r); close(copyDone) }()

	done := make(chan error, 1)
	go func() {
		done <- run(t, loginCmdDef(), "auth", "login",
			"--base-url", tokenSrv.URL, "--console-url", "https://console.test")
	}()

	select {
	case err := <-done:
		os.Stderr = oldStderr
		_ = w.Close()
		<-copyDone
		return exchangeForm, out.String(), err
	case <-time.After(10 * time.Second):
		t.Fatal("authLogin did not complete")
		panic("unreachable")
	}
}

// TestAuthLoginDefaultPrintsURL covers the default (browser) login path:
// the authorize URL printed to stderr uses the manual redirect_uri (so a
// remote user can paste the resulting code), while openBrowser is called with
// the loopback URL (so a same-machine browser auto-completes via callback).
func TestAuthLoginDefaultPrintsURL(t *testing.T) {
	opened := make(chan string, 1)

	// Callback-wins path: nothing on stdin; once openBrowser fires, deliver
	// the loopback callback the same way a local browser would.
	formCh := make(chan url.Values, 1)
	stderrCh := make(chan string, 1)
	go func() {
		form, stderr, err := driveDefaultLogin(t, opened, "")
		require.NoError(t, err)
		formCh <- form
		stderrCh <- stderr
	}()

	loopbackURL := <-opened
	lu, _ := url.Parse(loopbackURL)
	redirect := lu.Query().Get("redirect_uri")
	require.True(t, strings.HasPrefix(redirect, "http://127.0.0.1:") || strings.HasPrefix(redirect, "http://localhost:"),
		"openBrowser must receive the loopback redirect_uri, got %q", redirect)
	resp, err := http.Get(redirect + "?code=CB-CODE&state=" + url.QueryEscape(lu.Query().Get("state")))
	require.NoError(t, err)
	resp.Body.Close()

	form := <-formCh
	stderr := <-stderrCh
	assert.Equal(t, redirect, form.Get("redirect_uri"),
		"callback-wins path must exchange with the loopback redirect_uri")
	assert.Equal(t, "CB-CODE", form.Get("code"))

	// Stderr printed the manual URL, not the loopback one.
	printed := regexp.MustCompile(`https://\S+/oauth/authorize\?\S+`).FindString(stderr)
	require.NotEmpty(t, printed, "default mode must print an authorize URL")
	pu, _ := url.Parse(printed)
	assert.Equal(t, "https://console.test/oauth/code/callback?app=anthropic-cli",
		pu.Query().Get("redirect_uri"),
		"printed URL must use the manual redirect_uri so the displayed code can be pasted")
	assert.NotEqual(t, loopbackURL, printed)
	assert.Contains(t, stderr, "paste the code here")
	assert.NotContains(t, stderr, "--no-browser",
		"default mode no longer needs to hint at --no-browser")
}

// TestAuthLoginDefaultAcceptsPastedCode covers the stdin-wins side of the
// default-mode race: a code pasted on stdin completes login without any
// loopback callback, and the token exchange sends the manual redirect_uri.
func TestAuthLoginDefaultAcceptsPastedCode(t *testing.T) {
	opened := make(chan string, 1)
	form, stderr, err := driveDefaultLogin(t, opened, "PASTED-CODE")
	require.NoError(t, err)

	select {
	case <-opened:
	default:
		t.Fatal("openBrowser should still be called even when stdin wins")
	}

	assert.Equal(t, "PASTED-CODE", form.Get("code"))
	assert.Equal(t, "https://console.test/oauth/code/callback?app=anthropic-cli",
		form.Get("redirect_uri"),
		"stdin-wins path must exchange with the manual redirect_uri")
	assert.Contains(t, stderr, "Logged in")
}

// TestAuthLoginActivatesExplicitProfile covers the active_config write
// behavior: an explicit --profile always retargets active_config (and the
// success output says so); resolving via active_config / "default" leaves an
// existing active_config alone.
func TestAuthLoginActivatesExplicitProfile(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("ANTHROPIC_CONFIG_DIR", dir)
	clearEnv(t, "ANTHROPIC_PROFILE")

	srv := newTokenServer(t, tokenResponse{AccessToken: "tok", RefreshToken: "rt", ExpiresIn: 600})

	// First login --profile a → active_config=a (set from absent).
	driveLogin(t, srv.URL, "--profile", "a")
	require.Equal(t, "a", strings.TrimSpace(string(mustRead(t, config.ActiveConfigPath(dir)))))

	// Login --profile b → active_config retargeted to b.
	driveLogin(t, srv.URL, "--profile", "b")
	assert.Equal(t, "b", strings.TrimSpace(string(mustRead(t, config.ActiveConfigPath(dir)))),
		"explicit --profile must retarget active_config")

	// Login without --profile (resolves to b via active_config) → still b.
	driveLogin(t, srv.URL)
	assert.Equal(t, "b", strings.TrimSpace(string(mustRead(t, config.ActiveConfigPath(dir)))),
		"login without --profile must not retarget an existing active_config")

	// Login --profile b again (already active) → no change, but no error.
	driveLogin(t, srv.URL, "--profile", "b")
	assert.Equal(t, "b", strings.TrimSpace(string(mustRead(t, config.ActiveConfigPath(dir)))))
}

// TestAuthLoginHonorsActiveConfig is a regression test for `auth login`
// writing to "default" instead of the active_config profile when no
// --profile/ANTHROPIC_PROFILE is set. Before the fix, the login flag had
// Value:"default" and authLogin read c.String("profile") directly, bypassing
// the active_config tier of resolution that logout/status/print-credentials honor.
func TestAuthLoginHonorsActiveConfig(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("ANTHROPIC_CONFIG_DIR", dir)
	clearEnv(t, "ANTHROPIC_PROFILE")
	require.NoError(t, config.SetActiveProfile(dir, "personal"))

	tokenSrv := newTokenServer(t, tokenResponse{
		AccessToken: "sk-ant-oat01-PERSONAL", RefreshToken: "rt-P",
		ExpiresIn: 600, Scope: "user:inference",
	})

	driveLogin(t, tokenSrv.URL) // no --profile

	assert.FileExists(t, config.ProfilePath(dir, "personal"),
		"login without --profile must write to the active_config profile")
	assert.FileExists(t, config.ProfileCredentialsPath(dir, "personal"))
	assert.NoFileExists(t, config.ProfilePath(dir, "default"),
		"login without --profile must NOT write to \"default\" when active_config points elsewhere")
	assert.Equal(t, "personal",
		strings.TrimSpace(string(mustRead(t, config.ActiveConfigPath(dir)))),
		"existing active_config must be left untouched")
}

// TestAuthLoginCapturesOrganization covers organization handling in
// authLogin: the token response's organization.uuid is written to
// cfg.OrganizationID; on a fresh profile (no prev org) the authorize URL
// has no orgUUID hint.
func TestAuthLoginCapturesOrganization(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("ANTHROPIC_CONFIG_DIR", dir)
	clearEnv(t, "ANTHROPIC_PROFILE")

	srv := newTokenServer(t, tokenResponse{
		AccessToken: "tok", RefreshToken: "rt", ExpiresIn: 600,
		Organization: tokenOrganization{UUID: "org-A", Name: "Org A"},
		Account:      tokenAccount{EmailAddress: "user@example.com"},
	})

	u := driveLogin(t, srv.URL, "--profile", "fresh")
	assert.Empty(t, u.Query().Get("orgUUID"),
		"no org hint on a profile with no prior organization_id")

	var cfg map[string]any
	require.NoError(t, json.Unmarshal(mustRead(t, config.ProfilePath(dir, "fresh")), &cfg))
	assert.Equal(t, "org-A", cfg["organization_id"],
		"organization_id must be captured from the token response")
}

// TestAuthLoginOrgHintAndMismatch covers re-login on a profile with a
// stored organization_id: the authorize URL carries an orgUUID hint, a
// matching response succeeds quietly, and a mismatching response warns on
// stderr but still completes (Console enforces ?orgUUID upstream, so the
// post-exchange check is a backstop only).
func TestAuthLoginOrgHintAndMismatch(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("ANTHROPIC_CONFIG_DIR", dir)
	clearEnv(t, "ANTHROPIC_PROFILE")

	// Seed a profile already pinned to org-A.
	require.NoError(t, config.SaveProfile(dir, "pinned", &config.Config{
		AuthenticationInfo: &config.AuthenticationInfo{
			Type: config.AuthenticationTypeUserOAuth, UserOAuth: &config.UserOAuth{},
		},
		OrganizationID: "org-A",
	}))

	t.Run("hint and match succeeds", func(t *testing.T) {
		srv := newTokenServer(t, tokenResponse{
			AccessToken: "tok", RefreshToken: "rt", ExpiresIn: 600,
			Organization: tokenOrganization{UUID: "org-A", Name: "Org A"},
		})
		u := driveLogin(t, srv.URL, "--profile", "pinned")
		assert.Equal(t, "org-A", u.Query().Get("orgUUID"),
			"authorize URL must hint the profile's stored org")
		assert.FileExists(t, config.ProfileCredentialsPath(dir, "pinned"))
	})

	t.Run("mismatch warns and proceeds", func(t *testing.T) {
		// Re-seed without credentials so the writes-creds assertion is meaningful.
		require.NoError(t, os.Remove(config.ProfileCredentialsPath(dir, "pinned")))
		srv := newTokenServer(t, tokenResponse{
			AccessToken: "tok", RefreshToken: "rt", ExpiresIn: 600,
			Organization: tokenOrganization{UUID: "org-B", Name: "Org B"},
		})
		_, stderr, err := driveLoginErr(t, srv.URL, "--profile", "pinned")
		require.NoError(t, err, "org mismatch is a warning, not an error")
		assert.Contains(t, stderr, "org-B")
		assert.Contains(t, stderr, "org-A")
		assert.Contains(t, stderr, "organization_id")
		assert.FileExists(t, config.ProfileCredentialsPath(dir, "pinned"),
			"login proceeds and writes credentials despite mismatch warning")
		var cfg map[string]any
		require.NoError(t, json.Unmarshal(mustRead(t, config.ProfilePath(dir, "pinned")), &cfg))
		assert.Equal(t, "org-A", cfg["organization_id"],
			"existing profile config is never rewritten on re-login")
	})
}

func runPrintCredentials(t *testing.T, args ...string) (string, error) {
	t.Helper()
	return captureStdout(t, func() error {
		return run(t, &cli.Command{Name: "auth", Commands: []*cli.Command{{
			Name: "print-credentials", Action: authPrintCredentials,
			Flags: []cli.Flag{&cli.BoolFlag{Name: "access-token"}, &cli.BoolFlag{Name: "env"}},
		}}}, append([]string{"auth", "print-credentials"}, args...)...)
	})
}

func TestAuthPrintCredentials(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("ANTHROPIC_CONFIG_DIR", dir)
	clearEnv(t, "ANTHROPIC_PROFILE")
	require.NoError(t, config.SaveProfile(dir, "default", &config.Config{
		AuthenticationInfo: &config.AuthenticationInfo{Type: config.AuthenticationTypeUserOAuth, UserOAuth: &config.UserOAuth{}},
	}))
	require.NoError(t, config.SetActiveProfile(dir, "default"))
	exp := time.Now().Add(time.Hour)
	credsPath := config.ProfileCredentialsPath(dir, "default")
	require.NoError(t, config.WriteCredentials(credsPath, config.Credentials{
		AccessToken: "sk-ant-oat01-FRESH", RefreshToken: "rt-FRESH", ExpiresAt: &exp,
		Scope: "user:inference", OrganizationUUID: "org-X", OrganizationName: "X", AccountEmail: "x@example.com",
	}))

	t.Run("default emits the on-disk JSON shape", func(t *testing.T) {
		out, err := runPrintCredentials(t)
		require.NoError(t, err)
		var got map[string]any
		require.NoError(t, json.Unmarshal([]byte(out), &got))
		assert.Equal(t, "sk-ant-oat01-FRESH", got["access_token"])
		assert.Equal(t, "rt-FRESH", got["refresh_token"])
		assert.Equal(t, "org-X", got["organization_uuid"])
		assert.Equal(t, "oauth_token", got["type"])
		// stdout JSON exactly equals what's on disk (modulo indent), so the
		// "json object equivalent to the file" contract holds.
		var disk map[string]any
		require.NoError(t, json.Unmarshal(mustRead(t, credsPath), &disk))
		assert.Equal(t, disk, got)
	})

	t.Run("--access-token emits only the bare token", func(t *testing.T) {
		out, err := runPrintCredentials(t, "--access-token")
		require.NoError(t, err)
		assert.Equal(t, "sk-ant-oat01-FRESH\n", out)
	})

	t.Run("--env emits .env KEY=value lines", func(t *testing.T) {
		out, err := runPrintCredentials(t, "--env")
		require.NoError(t, err)
		assert.Equal(t, "ANTHROPIC_AUTH_TOKEN=sk-ant-oat01-FRESH\n", out)
	})

	t.Run("--access-token and --env are mutually exclusive", func(t *testing.T) {
		_, err := runPrintCredentials(t, "--access-token", "--env")
		require.ErrorContains(t, err, "mutually exclusive")
	})
}

func TestAuthPrintCredentials_EnvIncludesBaseURL(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("ANTHROPIC_CONFIG_DIR", dir)
	clearEnv(t, "ANTHROPIC_PROFILE")
	require.NoError(t, config.SaveProfile(dir, "default", &config.Config{
		AuthenticationInfo: &config.AuthenticationInfo{Type: config.AuthenticationTypeUserOAuth, UserOAuth: &config.UserOAuth{}},
		BaseURL:            "https://staging.example",
	}))
	require.NoError(t, config.SetActiveProfile(dir, "default"))
	exp := time.Now().Add(time.Hour)
	require.NoError(t, config.WriteCredentials(config.ProfileCredentialsPath(dir, "default"),
		config.Credentials{AccessToken: "sk-ant-oat01-ENV", ExpiresAt: &exp}))

	out, err := runPrintCredentials(t, "--env")
	require.NoError(t, err)
	assert.Equal(t,
		"ANTHROPIC_AUTH_TOKEN=sk-ant-oat01-ENV\n"+
			"ANTHROPIC_BASE_URL=https://staging.example\n",
		out)
}

func TestAuthPrintCredentials_RefreshesExpired(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("ANTHROPIC_CONFIG_DIR", dir)
	clearEnv(t, "ANTHROPIC_PROFILE")

	var gotReq map[string]string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/v1/oauth/token", r.URL.Path)
		require.Equal(t, "application/json", r.Header.Get("Content-Type"))
		require.Equal(t, betaUserOAuth, r.Header.Get("anthropic-beta"))
		require.NoError(t, json.NewDecoder(r.Body).Decode(&gotReq))
		_ = json.NewEncoder(w).Encode(tokenResponse{
			AccessToken: "sk-ant-oat01-REFRESHED", RefreshToken: "rt-NEW", ExpiresIn: 600,
		})
	}))
	t.Cleanup(srv.Close)

	require.NoError(t, config.SaveProfile(dir, "default", &config.Config{
		AuthenticationInfo: &config.AuthenticationInfo{
			Type: config.AuthenticationTypeUserOAuth, UserOAuth: &config.UserOAuth{ClientID: "cli-client"},
		},
		BaseURL: srv.URL,
	}))
	require.NoError(t, config.SetActiveProfile(dir, "default"))
	past := time.Now().Add(-time.Hour)
	credsPath := config.ProfileCredentialsPath(dir, "default")
	require.NoError(t, config.WriteCredentials(credsPath, config.Credentials{
		AccessToken: "sk-ant-oat01-STALE", RefreshToken: "rt-OLD", ExpiresAt: &past,
		OrganizationUUID: "org-KEEP",
	}))

	out, err := runPrintCredentials(t, "--access-token")
	require.NoError(t, err)
	assert.Equal(t, "sk-ant-oat01-REFRESHED\n", out)
	assert.Equal(t, "refresh_token", gotReq["grant_type"])
	assert.Equal(t, "rt-OLD", gotReq["refresh_token"])
	assert.Equal(t, "cli-client", gotReq["client_id"])

	// Refreshed credentials are persisted; org metadata not in the refresh
	// response is preserved from the prior file.
	var disk map[string]any
	require.NoError(t, json.Unmarshal(mustRead(t, credsPath), &disk))
	assert.Equal(t, "sk-ant-oat01-REFRESHED", disk["access_token"])
	assert.Equal(t, "rt-NEW", disk["refresh_token"])
	assert.Equal(t, "org-KEEP", disk["organization_uuid"])
}

func TestAuthPrintCredentials_RefreshFailureWarnsAndPrintsStale(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("ANTHROPIC_CONFIG_DIR", dir)
	clearEnv(t, "ANTHROPIC_PROFILE")

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, `{"error":"server_error"}`, http.StatusBadGateway)
	}))
	t.Cleanup(srv.Close)

	require.NoError(t, config.SaveProfile(dir, "default", &config.Config{
		AuthenticationInfo: &config.AuthenticationInfo{
			Type: config.AuthenticationTypeUserOAuth, UserOAuth: &config.UserOAuth{ClientID: "c"},
		},
		BaseURL: srv.URL,
	}))
	require.NoError(t, config.SetActiveProfile(dir, "default"))
	soon := time.Now().Add(30 * time.Second)
	require.NoError(t, config.WriteCredentials(config.ProfileCredentialsPath(dir, "default"),
		config.Credentials{AccessToken: "sk-ant-oat01-STILLOK", RefreshToken: "rt", ExpiresAt: &soon}))

	out, err := runPrintCredentials(t, "--access-token")
	require.NoError(t, err, "refresh failure should warn-and-proceed, not error")
	assert.Equal(t, "sk-ant-oat01-STILLOK\n", out)
}

// TestAuthLoginBootstrapOnly covers the principle that `auth login` writes
// configs/<profile>.json only when the profile doesn't already exist;
// re-login on an existing profile produces credentials only.
func TestAuthLoginBootstrapOnly(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("ANTHROPIC_CONFIG_DIR", dir)
	clearEnv(t, "ANTHROPIC_PROFILE")

	srv := newTokenServer(t, tokenResponse{
		AccessToken: "tok", RefreshToken: "rt", ExpiresIn: 600,
		Scope:        "user:inference",
		Organization: tokenOrganization{UUID: "org-BOOT", Name: "Boot"},
	})

	t.Run("fresh profile omits unset override flags", func(t *testing.T) {
		driveLogin(t, srv.URL, "--profile", "fresh")
		var cfg map[string]any
		require.NoError(t, json.Unmarshal(mustRead(t, config.ProfilePath(dir, "fresh")), &cfg))
		assert.Equal(t, "org-BOOT", cfg["organization_id"])
		auth, ok := cfg["authentication"].(map[string]any)
		require.True(t, ok)
		assert.Equal(t, "user_oauth", auth["type"])
		// client_id is always written: the profile file is shared with SDK
		// consumers whose refresh path requires it, and the refresh_token is
		// bound to this client. With no flag, resolveClientID → prod default.
		assert.Equal(t, oauthClientIDProd, auth["client_id"])
		_, hasScope := auth["scope"]
		assert.False(t, hasScope, "scope must be omitted when --scope not set")
	})

	t.Run("fresh profile persists explicit override flags", func(t *testing.T) {
		driveLogin(t, srv.URL, "--profile", "withflags",
			"--client-id", "custom-client", "--scope", "user:profile")
		var cfg map[string]any
		require.NoError(t, json.Unmarshal(mustRead(t, config.ProfilePath(dir, "withflags")), &cfg))
		auth := cfg["authentication"].(map[string]any)
		assert.Equal(t, "custom-client", auth["client_id"])
		assert.Equal(t, "user:inference", auth["scope"], "persisted scope is the granted value")
	})

	t.Run("existing profile config untouched by re-login", func(t *testing.T) {
		require.NoError(t, config.SaveProfile(dir, "keep", &config.Config{
			AuthenticationInfo: &config.AuthenticationInfo{Type: config.AuthenticationTypeUserOAuth, UserOAuth: &config.UserOAuth{}},
			WorkspaceID:        "wrkspc_test",
			OrganizationID:     "org-BOOT",
		}))
		before := mustRead(t, config.ProfilePath(dir, "keep"))

		driveLogin(t, srv.URL, "--profile", "keep",
			"--client-id", "ignored-client", "--scope", "user:profile")

		after := mustRead(t, config.ProfilePath(dir, "keep"))
		assert.Equal(t, string(before), string(after),
			"re-login on an existing profile must not rewrite configs/<profile>.json")
		assert.FileExists(t, config.ProfileCredentialsPath(dir, "keep"),
			"credentials must still be written")
	})

	t.Run("re-login without --workspace-id preserves stored value", func(t *testing.T) {
		require.NoError(t, config.SaveProfile(dir, "preserve", &config.Config{
			AuthenticationInfo: &config.AuthenticationInfo{Type: config.AuthenticationTypeUserOAuth, UserOAuth: &config.UserOAuth{}},
			WorkspaceID:        "wrkspc_alpha",
			OrganizationID:     "org-BOOT",
		}))
		clearEnv(t, "ANTHROPIC_WORKSPACE_ID")
		before := mustRead(t, config.ProfilePath(dir, "preserve"))

		u, out, err := driveLoginWithArgs(t, []string{"auth", "login", "--no-browser",
			"--callback-port", "0", "--base-url", srv.URL, "--profile", "preserve"})
		require.NoError(t, err)
		assert.Equal(t, "wrkspc_alpha", u.Query().Get("workspace_id"),
			"no --workspace-id flag → stored profile value is sent on /oauth/authorize")
		// prev.WorkspaceID == effectiveWorkspaceID → no drift messaging.
		assert.NotContains(t, out, "Token bound to workspace",
			"no drift messaging when token binds to the profile's stored workspace")

		assert.Equal(t, string(before), string(mustRead(t, config.ProfilePath(dir, "preserve"))),
			"re-login must not rewrite configs/<profile>.json")
	})

	t.Run("re-login on profile without workspace_id prints set hint", func(t *testing.T) {
		require.NoError(t, config.SaveProfile(dir, "nows", &config.Config{
			AuthenticationInfo: &config.AuthenticationInfo{Type: config.AuthenticationTypeUserOAuth, UserOAuth: &config.UserOAuth{}},
			OrganizationID:     "org-BOOT",
		}))
		before := mustRead(t, config.ProfilePath(dir, "nows"))

		wsSrv := newTokenServer(t, tokenResponse{
			AccessToken: "tok", RefreshToken: "rt", ExpiresIn: 600,
			Organization: tokenOrganization{UUID: "org-BOOT", Name: "Boot"},
			Workspace:    tokenWorkspace{ID: "wrkspc_hint", Name: "Hint Workspace"},
		})
		_, out, err := driveLoginWithArgs(t, []string{"auth", "login", "--no-browser",
			"--callback-port", "0", "--base-url", wsSrv.URL, "--profile", "nows"})
		require.NoError(t, err)

		assert.Contains(t, out, "→ Token bound to workspace \"Hint Workspace\" (wrkspc_hint)")
		assert.Contains(t, out, "ant profile set workspace_id wrkspc_hint --profile nows")
		assert.Equal(t, string(before), string(mustRead(t, config.ProfilePath(dir, "nows"))),
			"re-login must not rewrite configs/<profile>.json")
	})

	t.Run("re-login with workspace mismatch warns instead of retargeting", func(t *testing.T) {
		require.NoError(t, config.SaveProfile(dir, "mismatch", &config.Config{
			AuthenticationInfo: &config.AuthenticationInfo{Type: config.AuthenticationTypeUserOAuth, UserOAuth: &config.UserOAuth{}},
			WorkspaceID:        "wrkspc_OLD",
			OrganizationID:     "org-BOOT",
		}))
		before := mustRead(t, config.ProfilePath(dir, "mismatch"))

		wsSrv := newTokenServer(t, tokenResponse{
			AccessToken: "tok", RefreshToken: "rt", ExpiresIn: 600,
			Organization: tokenOrganization{UUID: "org-BOOT", Name: "Boot"},
			Workspace:    tokenWorkspace{ID: "wrkspc_NEW", Name: "New Workspace"},
		})
		_, out, err := driveLoginErr(t, wsSrv.URL, "--profile", "mismatch",
			"--workspace-id", "wrkspc_NEW")
		require.NoError(t, err)

		assert.Contains(t, out, "⚠ Token bound to workspace \"New Workspace\" (wrkspc_NEW)")
		assert.Contains(t, out, "but profile \"mismatch\" targets wrkspc_OLD")
		assert.Contains(t, out, "ant profile set workspace_id wrkspc_NEW --profile mismatch")
		assert.Contains(t, out, "ant auth login --profile <new> --workspace-id wrkspc_NEW")
		assert.Equal(t, string(before), string(mustRead(t, config.ProfilePath(dir, "mismatch"))),
			"re-login must not rewrite configs/<profile>.json (warn-only, no retarget)")
	})
}

// TestAuthLoginAcceptsWorkspaceFromTokenResponse: with no --workspace-id and
// no stored value, Console picker resolves it; we model that by having the
// mock token server return a Workspace block, and assert the CLI persists it.
func TestAuthLoginAcceptsWorkspaceFromTokenResponse(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("ANTHROPIC_CONFIG_DIR", dir)
	clearEnv(t, "ANTHROPIC_PROFILE")
	clearEnv(t, "ANTHROPIC_WORKSPACE_ID")

	srv := newTokenServer(t, tokenResponse{
		AccessToken: "tok", RefreshToken: "rt", ExpiresIn: 600,
		Workspace: tokenWorkspace{ID: "wrkspc_picker", Name: "Picker Workspace"},
	})

	u, _, err := driveLoginWithArgs(t, []string{"auth", "login", "--no-browser",
		"--callback-port", "0", "--base-url", srv.URL, "--profile", "fresh"})
	require.NoError(t, err)
	assert.Empty(t, u.Query().Get("workspace_id"),
		"workspace_id must be omitted so Console renders the picker")

	// Picker-resolved workspace must land in the profile config and credentials.
	var cfg map[string]any
	require.NoError(t, json.Unmarshal(mustRead(t, config.ProfilePath(dir, "fresh")), &cfg))
	assert.Equal(t, "wrkspc_picker", cfg["workspace_id"], "token response workspace persisted to profile")
	assert.True(t, strings.HasPrefix(cfg["workspace_id"].(string), "wrkspc_"),
		"config workspace_id must be the tagged form")
	var creds map[string]any
	require.NoError(t, json.Unmarshal(mustRead(t, config.ProfileCredentialsPath(dir, "fresh")), &creds))
	assert.Equal(t, "wrkspc_picker", creds["workspace_id"])
	assert.Equal(t, "Picker Workspace", creds["workspace_name"])
}

// TestAuthLoginRequiresWorkspaceFromSomewhere: with no flag, no stored value,
// AND a token response that omits the workspace block (older backend or a
// federation token) we error rather than write a profile with empty workspace.
func TestAuthLoginRequiresWorkspaceFromSomewhere(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("ANTHROPIC_CONFIG_DIR", dir)
	clearEnv(t, "ANTHROPIC_PROFILE")
	clearEnv(t, "ANTHROPIC_WORKSPACE_ID")

	srv := newTokenServer(t, tokenResponse{
		AccessToken: "tok", RefreshToken: "rt", ExpiresIn: 600,
		// no Workspace block
	})

	_, _, err := driveLoginWithArgs(t, []string{"auth", "login", "--no-browser",
		"--callback-port", "0", "--base-url", srv.URL, "--profile", "fresh"})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "no workspace bound")
	assert.NoFileExists(t, config.ProfilePath(dir, "fresh"),
		"nothing written when no workspace can be resolved")
}

// TestResolveWorkspaceIDPrecedence pins the lookup chain
// flag → prev.WorkspaceID → ANTHROPIC_WORKSPACE_ID → "" and that the
// env var is fill-missing only (does not override the profile).
func TestResolveWorkspaceIDPrecedence(t *testing.T) {
	prev := &config.Config{WorkspaceID: "wrkspc_prev"}
	for _, tc := range []struct {
		name, flag, env string
		prev            *config.Config
		want            string
	}{
		{"flag wins over prev and env", "wrkspc_flag", "wrkspc_env", prev, "wrkspc_flag"},
		{"prev wins over env", "", "wrkspc_env", prev, "wrkspc_prev"},
		{"env fills only when prev empty", "", "wrkspc_env",
			&config.Config{WorkspaceID: ""}, "wrkspc_env"},
		{"env fills when no prev", "", "wrkspc_env", nil, "wrkspc_env"},
		{"empty when nothing set", "", "", nil, ""},
		{"flag is trimmed", "  wrkspc_flag  ", "", nil, "wrkspc_flag"},
		{"env is trimmed", "", "  wrkspc_env  ", nil, "wrkspc_env"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv("ANTHROPIC_WORKSPACE_ID", tc.env)
			assert.Equal(t, tc.want, resolveWorkspaceID(tc.flag, tc.prev))
		})
	}
}

// TestAuthLogoutAllScopesDeletion guards against `--all` removing anything
// outside the configs/credentials/active_config set. ANTHROPIC_CONFIG_DIR is
// user-controlled, so RemoveAll(dir) on a misconfigured value would be
// catastrophic.
func TestAuthLogoutAllScopesDeletion(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("ANTHROPIC_CONFIG_DIR", dir)
	require.NoError(t, config.SaveProfile(dir, "p", &config.Config{AuthenticationInfo: &config.AuthenticationInfo{Type: config.AuthenticationTypeUserOAuth, UserOAuth: &config.UserOAuth{}}}))
	require.NoError(t, config.SetActiveProfile(dir, "p"))
	exp := time.Now().Add(time.Hour)
	require.NoError(t, config.WriteCredentials(config.ProfileCredentialsPath(dir, "p"),
		config.Credentials{AccessToken: "x", ExpiresAt: &exp}))
	// Unrelated content that happens to live under the configured dir.
	bystander := filepath.Join(dir, "unrelated.txt")
	require.NoError(t, os.WriteFile(bystander, []byte("keep me"), 0o644))

	require.NoError(t, run(t, &cli.Command{
		Name: "auth", Commands: []*cli.Command{{
			Name: "logout", Action: authLogout,
			Flags: []cli.Flag{&cli.BoolFlag{Name: "all"}},
		}},
	}, "auth", "logout", "--all"))

	assert.NoFileExists(t, config.ProfilePath(dir, "p"))
	assert.NoFileExists(t, config.ProfileCredentialsPath(dir, "p"))
	assert.NoFileExists(t, config.ActiveConfigPath(dir))
	assert.FileExists(t, bystander, "logout --all must not remove files it didn't create")
	assert.DirExists(t, dir, "logout --all must not remove the config dir itself")
}

// TestAuthLogoutAllWithoutActiveConfig covers the case where the user never
// had an active_config pointer (e.g. only ever wrote profiles by hand). The
// --all path uses os.RemoveAll, so a missing active_config must be a no-op
// rather than an error.
func TestAuthLogoutAllWithoutActiveConfig(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("ANTHROPIC_CONFIG_DIR", dir)
	require.NoError(t, config.SaveProfile(dir, "foo", &config.Config{AuthenticationInfo: &config.AuthenticationInfo{Type: config.AuthenticationTypeUserOAuth, UserOAuth: &config.UserOAuth{}}}))
	exp := time.Now().Add(time.Hour)
	require.NoError(t, config.WriteCredentials(config.ProfileCredentialsPath(dir, "foo"),
		config.Credentials{AccessToken: "x", ExpiresAt: &exp}))
	// Precondition: active_config does not exist.
	require.NoFileExists(t, config.ActiveConfigPath(dir))

	require.NoError(t, run(t, &cli.Command{
		Name: "auth", Commands: []*cli.Command{{
			Name: "logout", Action: authLogout,
			Flags: []cli.Flag{&cli.BoolFlag{Name: "all"}},
		}},
	}, "auth", "logout", "--all"))

	assert.NoFileExists(t, config.ProfilePath(dir, "foo"))
	assert.NoFileExists(t, config.ProfileCredentialsPath(dir, "foo"))
	assert.NoFileExists(t, config.ActiveConfigPath(dir))
}

// runStatus runs `auth status` against a synthetic root that carries the
// global flags authStatus reads via c.Root() (api-key/auth-token/base-url/
// organization-id/federation inputs). None are set in these tests; the flags
// exist so root.String/root.IsSet resolve to zero values rather than depend
// on undefined-flag behaviour.
func runStatus(t *testing.T) (string, error) {
	t.Helper()
	root := &cli.Command{
		Name: "ant",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "profile", Sources: cli.EnvVars("ANTHROPIC_PROFILE")},
			&cli.StringFlag{Name: "api-key"},
			&cli.StringFlag{Name: "auth-token"},
			&cli.StringFlag{Name: "base-url"},
			&cli.StringFlag{Name: "organization-id"},
			&cli.StringFlag{Name: "identity-token"},
			&cli.StringFlag{Name: "identity-token-file"},
			&cli.StringFlag{Name: "federation-rule"},
			&cli.StringFlag{Name: "service-account-id"},
		},
		Commands: []*cli.Command{{
			Name: "auth", Commands: []*cli.Command{{
				Name: "status", Action: authStatus,
			}},
		}},
	}
	return captureStdout(t, func() error {
		return root.Run(context.Background(), []string{"ant", "auth", "status"})
	})
}

// TestAuthStatusProfileNotConfigured covers a fresh config dir with no
// configs/<profile>.json at all (typical first run before `ant auth login`).
// Status should succeed and surface an actionable hint, distinct from the
// "config present, credentials missing" state.
func TestAuthStatusProfileNotConfigured(t *testing.T) {
	t.Setenv("ANTHROPIC_CONFIG_DIR", t.TempDir())
	clearEnv(t, "ANTHROPIC_PROFILE")
	clearEnv(t, "ANTHROPIC_BASE_URL")

	out, err := runStatus(t)
	require.NoError(t, err)
	assert.Contains(t, out, "Active profile:  default")
	assert.Contains(t, out, `profile "default" not configured`)
	assert.Contains(t, out, "ant auth login")
}

// TestAuthStatusCredentialsMissing covers a profile config that exists on
// disk but whose credentials file is absent — the user logged out, or wrote
// the config by hand. Status should report the missing credentials cleanly
// rather than panic or surface a raw fs error.
func TestAuthStatusCredentialsMissing(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("ANTHROPIC_CONFIG_DIR", dir)
	clearEnv(t, "ANTHROPIC_PROFILE")
	clearEnv(t, "ANTHROPIC_BASE_URL")
	require.NoError(t, config.SaveProfile(dir, "default", &config.Config{
		AuthenticationInfo: &config.AuthenticationInfo{Type: config.AuthenticationTypeUserOAuth, UserOAuth: &config.UserOAuth{}},
	}))
	// No credentials/default.json written.
	require.NoFileExists(t, config.ProfileCredentialsPath(dir, "default"))

	out, err := runStatus(t)
	require.NoError(t, err)
	assert.Contains(t, out, "Active profile:  default")
	assert.Contains(t, out, "configured but not logged in — run `ant auth login`")
}

// TestAuthStatusExpiredToken covers a profile whose stored access_token has
// already expired. Status should still treat the profile credential as
// present (it may be refreshable) and surface the expiry as "expired … ago".
func TestAuthStatusExpiredToken(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("ANTHROPIC_CONFIG_DIR", dir)
	clearEnv(t, "ANTHROPIC_PROFILE")
	clearEnv(t, "ANTHROPIC_BASE_URL")
	require.NoError(t, config.SaveProfile(dir, "default", &config.Config{
		AuthenticationInfo: &config.AuthenticationInfo{Type: config.AuthenticationTypeUserOAuth, UserOAuth: &config.UserOAuth{}},
	}))
	past := time.Now().Add(-time.Hour)
	require.NoError(t, config.WriteCredentials(config.ProfileCredentialsPath(dir, "default"),
		config.Credentials{AccessToken: "sk-ant-oat01-EXPIRED", ExpiresAt: &past}))

	out, err := runStatus(t)
	require.NoError(t, err)
	assert.Contains(t, out, "Profile (user_oauth)")
	assert.Contains(t, out, "expired")
	assert.Contains(t, out, "ago")
}

// TestAuthStatusLoggedInLine covers the "Logged in to <org> as <email>"
// headline when the credentials file carries OrganizationName and
// AccountEmail (written on every login). Older creds without those fields
// should omit the line entirely.
func TestAuthStatusLoggedInLine(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("ANTHROPIC_CONFIG_DIR", dir)
	clearEnv(t, "ANTHROPIC_PROFILE")
	clearEnv(t, "ANTHROPIC_BASE_URL")
	require.NoError(t, config.SaveProfile(dir, "default", &config.Config{
		AuthenticationInfo: &config.AuthenticationInfo{Type: config.AuthenticationTypeUserOAuth, UserOAuth: &config.UserOAuth{}},
	}))
	exp := time.Now().Add(time.Hour)

	t.Run("full creds", func(t *testing.T) {
		require.NoError(t, config.WriteCredentials(config.ProfileCredentialsPath(dir, "default"),
			config.Credentials{
				AccessToken: "sk-ant-oat01-X", ExpiresAt: &exp,
				OrganizationName: "Acme Inc", AccountEmail: "user@example.com",
			}))
		out, err := runStatus(t)
		require.NoError(t, err)
		assert.Contains(t, out, "Logged in to Acme Inc as user@example.com")
	})

	t.Run("legacy creds without org/email", func(t *testing.T) {
		require.NoError(t, config.WriteCredentials(config.ProfileCredentialsPath(dir, "default"),
			config.Credentials{AccessToken: "sk-ant-oat01-X", ExpiresAt: &exp}))
		out, err := runStatus(t)
		require.NoError(t, err)
		assert.NotContains(t, out, "Logged in to")
		assert.NotContains(t, out, "Logged in as")
	})
}

func mustRead(t *testing.T, path string) []byte {
	t.Helper()
	b, err := os.ReadFile(path)
	require.NoError(t, err)
	return b
}
