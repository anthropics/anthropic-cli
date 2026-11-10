package cmd

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/anthropics/anthropic-sdk-go/v2/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli/v3"
)

const (
	logoutAccessToken  = "sk-ant-oat01-LOGOUT"
	logoutRefreshToken = "rt-LOGOUT"
)

// logoutDir points the CLI at an empty config dir and clears the variables
// that pick a profile or fill in a profile's base URL.
func logoutDir(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	t.Setenv("ANTHROPIC_CONFIG_DIR", dir)
	clearEnv(t, "ANTHROPIC_PROFILE")
	clearEnv(t, "ANTHROPIC_BASE_URL")
	return dir
}

// runLogout runs `ant auth logout` and returns what it wrote to stderr.
func runLogout(t *testing.T, args ...string) string {
	t.Helper()
	return captureStderr(t, func() {
		require.NoError(t, run(t, &cli.Command{
			Name: "auth", Commands: []*cli.Command{{
				Name: "logout", Action: authLogout,
				Flags: []cli.Flag{&cli.BoolFlag{Name: "all"}},
			}},
		}, append([]string{"auth", "logout"}, args...)...))
	})
}

func seedLogin(t *testing.T, dir, profile, baseURL, clientID string, creds config.Credentials) {
	t.Helper()
	require.NoError(t, config.SaveProfile(dir, profile, &config.Config{
		AuthenticationInfo: config.NewUserOAuthAuthentication(clientID),
		BaseURL:            baseURL,
	}))
	require.NoError(t, config.WriteCredentials(config.ProfileCredentialsPath(dir, profile), creds))
}

type revokeRequest struct {
	method, path, contentType, authorization string
	body                                     map[string]string
}

// revokeServer records the requests it gets and answers them with status and
// the request body, so that printing a response body would print a token.
type revokeServer struct {
	*httptest.Server
	mu       sync.Mutex
	requests []revokeRequest
}

func newRevokeServer(t *testing.T, status int) *revokeServer {
	t.Helper()
	s := &revokeServer{}
	s.Server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		raw, _ := io.ReadAll(r.Body)
		got := revokeRequest{
			method: r.Method, path: r.URL.Path,
			contentType: r.Header.Get("Content-Type"), authorization: r.Header.Get("Authorization"),
		}
		_ = json.Unmarshal(raw, &got.body)
		s.mu.Lock()
		s.requests = append(s.requests, got)
		s.mu.Unlock()
		w.WriteHeader(status)
		_, _ = w.Write(raw)
	}))
	t.Cleanup(s.Close)
	return s
}

func (s *revokeServer) got() []revokeRequest {
	s.mu.Lock()
	defer s.mu.Unlock()
	return slices.Clone(s.requests)
}

func wantRevoke(token, hint, clientID string) revokeRequest {
	return revokeRequest{
		method: "POST", path: "/v1/oauth/token/revoke", contentType: "application/json",
		body: map[string]string{"token": token, "token_type_hint": hint, "client_id": clientID},
	}
}

func TestAuthLogoutRevokesTokens(t *testing.T) {
	dir := logoutDir(t)
	srv := newRevokeServer(t, http.StatusOK)
	seedLogin(t, dir, "work", srv.URL, "client-work", config.Credentials{
		AccessToken: logoutAccessToken, RefreshToken: logoutRefreshToken,
	})

	stderr := runLogout(t, "--profile", "work")

	assert.Equal(t, []revokeRequest{
		{
			method: "POST", path: "/v1/oauth/token/revoke", contentType: "application/json",
			body: map[string]string{"token": logoutRefreshToken, "token_type_hint": "refresh_token", "client_id": "client-work"},
		},
		{
			method: "POST", path: "/v1/oauth/token/revoke", contentType: "application/json",
			body: map[string]string{"token": logoutAccessToken, "token_type_hint": "access_token", "client_id": "client-work"},
		},
	}, srv.got())
	assert.Equal(t, "✓ Logged out of profile \"work\".\n", stderr)
	assert.NoFileExists(t, config.ProfileCredentialsPath(dir, "work"))
	assert.FileExists(t, config.ProfilePath(dir, "work"), "logout keeps the profile config")
}

// TestAuthLogoutRevokeTarget: the tokens go where a refresh for the profile
// would send its refresh token, under the client a refresh would name.
func TestAuthLogoutRevokeTarget(t *testing.T) {
	creds := config.Credentials{RefreshToken: logoutRefreshToken}

	t.Run("the profile's base_url beats ANTHROPIC_BASE_URL", func(t *testing.T) {
		dir := logoutDir(t)
		srv, env := newRevokeServer(t, http.StatusOK), newRevokeServer(t, http.StatusOK)
		t.Setenv("ANTHROPIC_BASE_URL", env.URL)
		seedLogin(t, dir, "default", srv.URL+"/", "client-x", creds)

		runLogout(t)

		assert.Equal(t, []revokeRequest{wantRevoke(logoutRefreshToken, "refresh_token", "client-x")}, srv.got())
		assert.Empty(t, env.got())
	})

	t.Run("a profile with no base_url follows ANTHROPIC_BASE_URL", func(t *testing.T) {
		dir := logoutDir(t)
		env := newRevokeServer(t, http.StatusOK)
		t.Setenv("ANTHROPIC_BASE_URL", env.URL)
		seedLogin(t, dir, "default", "", "client-x", creds)

		runLogout(t)

		assert.Equal(t, []revokeRequest{wantRevoke(logoutRefreshToken, "refresh_token", "client-x")}, env.got())
	})

	t.Run("a login with no client_id names the default client", func(t *testing.T) {
		dir := logoutDir(t)
		srv := newRevokeServer(t, http.StatusOK)
		seedLogin(t, dir, "default", srv.URL, "", creds)

		runLogout(t)

		assert.Equal(t, []revokeRequest{wantRevoke(logoutRefreshToken, "refresh_token", oauthClientIDProd)}, srv.got())
	})
}

// TestAuthLogoutRevokeFailureStillLogsOut: whatever stops the revoke, logout
// removes the credentials, says so in one warning, and prints no token.
func TestAuthLogoutRevokeFailureStillLogsOut(t *testing.T) {
	creds := config.Credentials{AccessToken: logoutAccessToken, RefreshToken: logoutRefreshToken}

	for name, tc := range map[string]struct {
		// seed stores the login and returns a count of the requests seen, which
		// the server that never answers leaves at 0.
		seed         func(t *testing.T, dir string) func() int
		wantRequests int
		wantWarning  string
	}{
		"server error": {
			seed: func(t *testing.T, dir string) func() int {
				srv := newRevokeServer(t, http.StatusInternalServerError)
				seedLogin(t, dir, "default", srv.URL, "client-x", creds)
				return func() int { return len(srv.got()) }
			},
			wantRequests: 2,
			wantWarning:  "revoke endpoint returned HTTP 500\n",
		},
		"connection refused": {
			seed: func(t *testing.T, dir string) func() int {
				srv := newRevokeServer(t, http.StatusOK)
				srv.Close()
				seedLogin(t, dir, "default", srv.URL, "client-x", creds)
				return func() int { return len(srv.got()) }
			},
			wantWarning: "dial tcp",
		},
		"no answer before the timeout": {
			seed: func(t *testing.T, dir string) func() int {
				orig := revokeTimeout
				revokeTimeout = 50 * time.Millisecond
				t.Cleanup(func() { revokeTimeout = orig })
				release := make(chan struct{})
				srv := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) { <-release }))
				t.Cleanup(srv.Close)
				t.Cleanup(func() { close(release) })
				seedLogin(t, dir, "default", srv.URL, "client-x", creds)
				return func() int { return 0 }
			},
			wantWarning: "context deadline exceeded",
		},
		"cleartext base URL": {
			seed: func(t *testing.T, dir string) func() int {
				sent := 0
				orig := http.DefaultTransport
				http.DefaultTransport = roundTripFunc(func(*http.Request) (*http.Response, error) {
					sent++
					return nil, errors.New("no network in this test")
				})
				t.Cleanup(func() { http.DefaultTransport = orig })
				seedLogin(t, dir, "default", "http://gateway.example.com", "client-x", creds)
				return func() int { return sent }
			},
			wantWarning: "non-https token endpoint",
		},
		"unparsable credentials": {
			seed: func(t *testing.T, dir string) func() int {
				srv := newRevokeServer(t, http.StatusOK)
				seedLogin(t, dir, "default", srv.URL, "client-x", creds)
				require.NoError(t, os.WriteFile(config.ProfileCredentialsPath(dir, "default"),
					[]byte(`{"access_token": `+logoutAccessToken), 0o600))
				return func() int { return len(srv.got()) }
			},
			wantWarning: "could not be parsed",
		},
		"no profile config": {
			seed: func(t *testing.T, dir string) func() int {
				srv := newRevokeServer(t, http.StatusOK)
				seedLogin(t, dir, "default", srv.URL, "client-x", creds)
				require.NoError(t, os.Remove(config.ProfilePath(dir, "default")))
				return func() int { return len(srv.got()) }
			},
			wantWarning: "failed to read config file",
		},
		"profile config with no authentication block": {
			seed: func(t *testing.T, dir string) func() int {
				srv := newRevokeServer(t, http.StatusOK)
				seedLogin(t, dir, "default", srv.URL, "client-x", creds)
				require.NoError(t, os.WriteFile(config.ProfilePath(dir, "default"),
					[]byte(`{"base_url": "`+srv.URL+`"}`), 0o644))
				return func() int { return len(srv.got()) }
			},
			wantWarning: "missing authentication",
		},
	} {
		t.Run(name, func(t *testing.T) {
			dir := logoutDir(t)
			requests := tc.seed(t, dir)

			start := time.Now()
			stderr := runLogout(t)

			assert.Less(t, time.Since(start), 5*time.Second)
			assert.Equal(t, tc.wantRequests, requests())
			assert.NoFileExists(t, config.ProfileCredentialsPath(dir, "default"))
			assert.Equal(t, 1, strings.Count(stderr, "warning:"), stderr)
			assert.Contains(t, stderr, `could not revoke the tokens for profile "default"`)
			assert.Contains(t, stderr, tc.wantWarning)
			assert.Contains(t, stderr, `✓ Logged out of profile "default".`)
			assert.NotContains(t, stderr, logoutAccessToken)
			assert.NotContains(t, stderr, logoutRefreshToken)
		})
	}
}

func TestAuthLogoutRevokeWarningEscapesControlCharacters(t *testing.T) {
	dir := logoutDir(t)
	require.NoError(t, os.MkdirAll(config.ProfileCredentialsPath(dir, "a\x1bb"), 0o700))

	stderr := runLogout(t, "--profile", "a\x1bb")

	assert.Contains(t, stderr, `a\x1bb.json`)
	assert.NotContains(t, stderr, "\x1b")
}

// TestAuthLogoutRevokesOnlyLogins: nothing is sent, and nothing is warned
// about, when the profile holds no token that a login issued.
func TestAuthLogoutRevokesOnlyLogins(t *testing.T) {
	for name, tc := range map[string]struct {
		seed func(t *testing.T, dir, baseURL string)
		want string
	}{
		"not logged in": {
			seed: func(t *testing.T, dir, baseURL string) {
				seedLogin(t, dir, "default", baseURL, "client-x", config.Credentials{})
				require.NoError(t, os.Remove(config.ProfileCredentialsPath(dir, "default")))
			},
			want: "Not logged in (profile \"default\").\n",
		},
		"no profile at all": {
			seed: func(*testing.T, string, string) {},
			want: "Not logged in (profile \"default\").\n",
		},
		"credentials file with no tokens": {
			seed: func(t *testing.T, dir, baseURL string) {
				seedLogin(t, dir, "default", baseURL, "client-x", config.Credentials{})
			},
			want: "✓ Logged out of profile \"default\".\n",
		},
		"hand-written bearer token": {
			seed: func(t *testing.T, dir, baseURL string) {
				seedLogin(t, dir, "default", baseURL, "", config.Credentials{AccessToken: logoutAccessToken})
			},
			want: "✓ Logged out of profile \"default\".\n",
		},
		"cached federation token": {
			seed: func(t *testing.T, dir, baseURL string) {
				require.NoError(t, config.SaveProfile(dir, "default", &config.Config{
					AuthenticationInfo: config.NewOIDCFederationAuthentication(config.OIDCFederation{FederationRuleID: "fdrl_x"}),
					BaseURL:            baseURL,
				}))
				require.NoError(t, config.WriteCredentials(config.ProfileCredentialsPath(dir, "default"),
					config.Credentials{AccessToken: logoutAccessToken}))
			},
			want: "✓ Logged out of profile \"default\".\n",
		},
	} {
		t.Run(name, func(t *testing.T) {
			dir := logoutDir(t)
			srv := newRevokeServer(t, http.StatusOK)
			tc.seed(t, dir, srv.URL)

			stderr := runLogout(t)

			assert.Equal(t, tc.want, stderr)
			assert.Empty(t, srv.got())
			assert.NoFileExists(t, config.ProfileCredentialsPath(dir, "default"))
		})
	}
}

func TestAuthLogoutAllRevokesEachProfile(t *testing.T) {
	dir := logoutDir(t)
	a, b, down := newRevokeServer(t, http.StatusOK), newRevokeServer(t, http.StatusOK), newRevokeServer(t, http.StatusBadGateway)
	seedLogin(t, dir, "a", a.URL, "client-a", config.Credentials{AccessToken: "sk-ant-oat01-A", RefreshToken: "rt-A"})
	seedLogin(t, dir, "b", b.URL, "client-b", config.Credentials{AccessToken: "sk-ant-oat01-B", RefreshToken: "rt-B"})
	seedLogin(t, dir, "down", down.URL, "client-down", config.Credentials{AccessToken: "sk-ant-oat01-DOWN", RefreshToken: "rt-DOWN"})
	require.NoError(t, config.SetActiveProfile(dir, "a"))

	stderr := runLogout(t, "--all")

	assert.Equal(t, []revokeRequest{
		wantRevoke("rt-A", "refresh_token", "client-a"), wantRevoke("sk-ant-oat01-A", "access_token", "client-a"),
	}, a.got())
	assert.Equal(t, []revokeRequest{
		wantRevoke("rt-B", "refresh_token", "client-b"), wantRevoke("sk-ant-oat01-B", "access_token", "client-b"),
	}, b.got())
	assert.Len(t, down.got(), 2)
	assert.Equal(t, 1, strings.Count(stderr, "warning:"), stderr)
	assert.Contains(t, stderr, `could not revoke the tokens for profile "down"`)
	assert.NotContains(t, stderr, "DOWN")
	assert.Contains(t, stderr, "✓ Removed all profiles and credentials.")
	assert.NoDirExists(t, config.ProfilesDir(dir))
	assert.NoDirExists(t, config.CredentialsDir(dir))
	assert.NoFileExists(t, config.ActiveConfigPath(dir))
}

func TestAuthLogoutAllGetsPastAProfileItCannotLoad(t *testing.T) {
	dir := logoutDir(t)
	srv := newRevokeServer(t, http.StatusOK)
	seedLogin(t, dir, "a-bare", srv.URL, "client-bare", config.Credentials{AccessToken: "sk-ant-oat01-BARE", RefreshToken: "rt-BARE"})
	require.NoError(t, os.WriteFile(config.ProfilePath(dir, "a-bare"), []byte(`{"base_url": "`+srv.URL+`"}`), 0o644))
	seedLogin(t, dir, "b-login", srv.URL, "client-login", config.Credentials{RefreshToken: logoutRefreshToken})

	stderr := runLogout(t, "--all")

	assert.Equal(t, []revokeRequest{wantRevoke(logoutRefreshToken, "refresh_token", "client-login")}, srv.got())
	assert.Equal(t, 1, strings.Count(stderr, "warning:"), stderr)
	assert.Contains(t, stderr, `could not revoke the tokens for profile "a-bare"`)
	assert.NotContains(t, stderr, "BARE")
	assert.NoDirExists(t, config.ProfilesDir(dir))
	assert.NoDirExists(t, config.CredentialsDir(dir))
}

func TestAuthLogoutAllWarnsForCredentialsWithNoProfile(t *testing.T) {
	dir := logoutDir(t)
	srv := newRevokeServer(t, http.StatusOK)
	seedLogin(t, dir, "login", srv.URL, "client-login", config.Credentials{RefreshToken: logoutRefreshToken})
	seedLogin(t, dir, "orphan", srv.URL, "client-orphan", config.Credentials{AccessToken: "sk-ant-oat01-ORPHAN", RefreshToken: "rt-ORPHAN"})
	require.NoError(t, os.Remove(config.ProfilePath(dir, "orphan")))
	require.NoError(t, os.WriteFile(filepath.Join(config.CredentialsDir(dir), "login"), nil, 0o600))
	require.NoError(t, os.Mkdir(filepath.Join(config.CredentialsDir(dir), "dir.json"), 0o700))

	stderr := runLogout(t, "--all")

	assert.Equal(t, []revokeRequest{wantRevoke(logoutRefreshToken, "refresh_token", "client-login")}, srv.got())
	assert.Equal(t, 1, strings.Count(stderr, "warning:"), stderr)
	assert.Contains(t, stderr, `could not revoke the tokens for profile "orphan"`)
	assert.NotContains(t, stderr, "ORPHAN")
	assert.NoDirExists(t, config.CredentialsDir(dir))
}

func TestAuthLogoutAllWarnsWhenCredentialsCannotBeListed(t *testing.T) {
	dir := logoutDir(t)
	require.NoError(t, os.WriteFile(config.CredentialsDir(dir), nil, 0o600))

	stderr := runLogout(t, "--all")

	assert.Contains(t, stderr, "warning: could not list the credentials")
	assert.Contains(t, stderr, "✓ Removed all profiles and credentials.")
	assert.NoFileExists(t, config.CredentialsDir(dir))
}

func TestAuthLogoutAllWithNoCredentialsDirDoesNotWarn(t *testing.T) {
	logoutDir(t)

	stderr := runLogout(t, "--all")

	assert.Equal(t, "✓ Removed all profiles and credentials.\n", stderr)
}
