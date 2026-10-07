package cmd

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/anthropics/anthropic-sdk-go/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOAuthTokenStaysOnTheConsoleDeployment(t *testing.T) {
	t.Run("default console rejects a different token host", func(t *testing.T) {
		hits := 0
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			hits++
			w.WriteHeader(http.StatusOK)
		}))
		t.Cleanup(srv.Close)

		_, err := oauthTokenBase(defaultConsoleURL, srv.URL)
		require.Error(t, err)
		require.ErrorContains(t, err, "refusing to send OAuth credentials")
		assert.Equal(t, 0, hits)

		got, err := oauthTokenBase(defaultConsoleURL, "")
		require.NoError(t, err)
		assert.Equal(t, defaultBaseURL, got)
	})

	t.Run("print-credentials does not post the refresh token to the profile base URL", func(t *testing.T) {
		dir := t.TempDir()
		t.Setenv("ANTHROPIC_CONFIG_DIR", dir)
		clearEnv(t, "ANTHROPIC_PROFILE")

		hits := 0
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			hits++
			w.WriteHeader(http.StatusOK)
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
		require.NoError(t, config.WriteCredentials(config.ProfileCredentialsPath(dir, "default"), config.Credentials{
			AccessToken: "sk-ant-oat01-STALE", RefreshToken: "rt-PROOF", ExpiresAt: &past,
		}))

		out, err := runPrintCredentials(t, "--access-token")
		require.NoError(t, err)
		assert.Equal(t, "sk-ant-oat01-STALE\n", out)
		assert.Equal(t, 0, hits)
	})
}
