package cmd

import (
	"context"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/config"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli/v3"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) { return f(req) }

// TestBaseURLEnvRealTree drives the production Command tree end to end: a
// probe subcommand runs getDefaultRequestOptions and issues a request through
// the real SDK client, asserting which host it reaches. Guards the v1.5.0
// regression where ANTHROPIC_BASE_URL was silently ignored and the env
// credential was sent to the SDK default endpoint instead of the configured
// deployment. Single Run: urfave applies a flag's env sources once per
// process, so the global tree cannot be re-Run per subtest.
func TestBaseURLEnvRealTree(t *testing.T) {
	clearCredentialEnv(t)
	t.Setenv("ANTHROPIC_BASE_URL", "https://env.invalid")
	t.Setenv("ANTHROPIC_API_KEY", "test-fake-api-key-not-real")

	var captured *http.Request
	probe := &cli.Command{
		Name: "probe",
		Action: func(ctx context.Context, c *cli.Command) error {
			opts := append(getDefaultRequestOptions(c),
				option.WithHTTPClient(&http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
					captured = req
					return &http.Response{
						StatusCode: http.StatusOK,
						Header:     http.Header{"Content-Type": []string{"application/json"}},
						Body:       io.NopCloser(strings.NewReader("{}")),
						Request:    req,
					}, nil
				})}),
				option.WithMaxRetries(0),
			)
			client := anthropic.NewClient(opts...)
			return client.Get(ctx, "/v1/models", nil, nil)
		},
	}
	saved := Command.Commands
	Command.Commands = append(append([]*cli.Command{}, saved...), probe)
	t.Cleanup(func() { Command.Commands = saved })

	require.NoError(t, Command.Run(context.Background(), []string{"ant", "probe"}))
	require.NotNil(t, captured, "no request captured")
	assert.Equal(t, "env.invalid", captured.URL.Host, "request host")
	assert.Equal(t, "test-fake-api-key-not-real", captured.Header.Get("x-api-key"))
}

// runShadowedLeaf runs a fresh tree whose root mirrors the global credential /
// workspace flags and whose leaf declares local --workspace-id and
// --service-account-id (as generated admin subcommands do for path params),
// returning the request getDefaultRequestOptions produced from the leaf.
func runShadowedLeaf(t *testing.T, argv ...string) *http.Request {
	t.Helper()
	var captured *http.Request
	leaf := &cli.Command{
		Name: "leaf",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "workspace-id"},
			&cli.StringFlag{Name: "service-account-id"},
		},
		Action: func(ctx context.Context, c *cli.Command) error {
			opts := append(getDefaultRequestOptions(c),
				option.WithHTTPClient(&http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
					captured = req
					return &http.Response{
						StatusCode: http.StatusOK,
						Header:     http.Header{"Content-Type": []string{"application/json"}},
						Body:       io.NopCloser(strings.NewReader("{}")),
						Request:    req,
					}, nil
				})}),
				option.WithMaxRetries(0),
			)
			client := anthropic.NewClient(opts...)
			return client.Get(ctx, "/v1/models", nil, nil)
		},
	}
	root := &cli.Command{
		Name: "ant",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "api-key", Sources: cli.EnvVars("ANTHROPIC_API_KEY")},
			&cli.StringFlag{Name: "auth-token"},
			&cli.BoolFlag{Name: "api-key-stdin"},
			&cli.BoolFlag{Name: "auth-token-stdin"},
			&cli.StringFlag{Name: "webhook-key"},
			&cli.StringFlag{Name: "base-url"},
			&cli.StringFlag{Name: "profile", Sources: cli.EnvVars("ANTHROPIC_PROFILE")},
			&cli.StringFlag{Name: "identity-token"},
			&cli.StringFlag{Name: "identity-token-file"},
			&cli.StringFlag{Name: "federation-rule"},
			&cli.StringFlag{Name: "organization-id"},
			&cli.StringFlag{Name: "service-account-id", Sources: cli.EnvVars("ANTHROPIC_SERVICE_ACCOUNT_ID")},
			&cli.StringFlag{Name: "workspace-id", Sources: cli.EnvVars("ANTHROPIC_WORKSPACE_ID")},
		},
		Commands: []*cli.Command{leaf},
	}
	require.NoError(t, root.Run(context.Background(), append([]string{"ant"}, argv...)))
	require.NotNil(t, captured, "no request captured")
	return captured
}

func clearCredentialEnv(t *testing.T) {
	t.Helper()
	t.Setenv("ANTHROPIC_CONFIG_DIR", t.TempDir())
	for _, k := range []string{
		"ANTHROPIC_BASE_URL", "ANTHROPIC_API_KEY", "ANTHROPIC_AUTH_TOKEN",
		"ANTHROPIC_PROFILE", "ANTHROPIC_IDENTITY_TOKEN", "ANTHROPIC_IDENTITY_TOKEN_FILE",
		"ANTHROPIC_FEDERATION_RULE_ID", "ANTHROPIC_ORGANIZATION_ID",
		"ANTHROPIC_SERVICE_ACCOUNT_ID", "ANTHROPIC_WORKSPACE_ID",
	} {
		clearEnv(t, k)
	}
}

// TestWorkspaceHeaderIgnoresLocalFlag: the anthropic-workspace-id header
// comes from the global --workspace-id / ANTHROPIC_WORKSPACE_ID only; a
// subcommand's same-named local flag (a path/body param) must neither be sent
// as the header nor hide the global value.
func TestWorkspaceHeaderIgnoresLocalFlag(t *testing.T) {
	t.Run("env reaches header past local flag", func(t *testing.T) {
		clearCredentialEnv(t)
		t.Setenv("ANTHROPIC_API_KEY", "k")
		t.Setenv("ANTHROPIC_WORKSPACE_ID", "wrkspc_env")
		req := runShadowedLeaf(t, "leaf", "--workspace-id", "wrkspc_path")
		assert.Equal(t, "wrkspc_env", req.Header.Get("anthropic-workspace-id"))
	})

	t.Run("global flag reaches header past local flag", func(t *testing.T) {
		clearCredentialEnv(t)
		t.Setenv("ANTHROPIC_API_KEY", "k")
		req := runShadowedLeaf(t, "--workspace-id", "wrkspc_global", "leaf", "--workspace-id", "wrkspc_path")
		assert.Equal(t, "wrkspc_global", req.Header.Get("anthropic-workspace-id"))
	})

	t.Run("local flag alone sends no header", func(t *testing.T) {
		clearCredentialEnv(t)
		t.Setenv("ANTHROPIC_API_KEY", "k")
		req := runShadowedLeaf(t, "leaf", "--workspace-id", "wrkspc_path")
		assert.Empty(t, req.Header.Get("anthropic-workspace-id"))
	})
}

// TestLocalServiceAccountFlagDoesNotSelectFederation: with only an implicit
// active_config profile configured, a subcommand's local --service-account-id
// must not make the credential chain take the federation tier — the profile's
// token is used and no federation source is reported.
func TestLocalServiceAccountFlagDoesNotSelectFederation(t *testing.T) {
	clearCredentialEnv(t)
	resetWarnOnce(t)
	dir := os.Getenv("ANTHROPIC_CONFIG_DIR")
	require.NoError(t, config.SaveProfile(dir, "p", &config.Config{
		AuthenticationInfo: &config.AuthenticationInfo{
			Type: config.AuthenticationTypeUserOAuth, UserOAuth: &config.UserOAuth{ClientID: "c"},
		},
	}))
	require.NoError(t, config.WriteCredentials(config.ProfileCredentialsPath(dir, "p"),
		config.Credentials{AccessToken: "tok", RefreshToken: "rt"}))
	require.NoError(t, config.SetActiveProfile(dir, "p"))

	var req *http.Request
	stderr := captureStderr(t, func() {
		req = runShadowedLeaf(t, "leaf", "--service-account-id", "svac_x")
	})
	assert.Equal(t, "Bearer tok", req.Header.Get("Authorization"), "implicit profile must win; a local --service-account-id is not federation config")
	assert.NotContains(t, stderr, "federation")
}
