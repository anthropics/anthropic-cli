package cmd

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/anthropics/anthropic-sdk-go"
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
	t.Setenv("ANTHROPIC_CONFIG_DIR", t.TempDir())
	for _, k := range []string{
		"ANTHROPIC_BASE_URL", "ANTHROPIC_API_KEY", "ANTHROPIC_AUTH_TOKEN",
		"ANTHROPIC_PROFILE", "ANTHROPIC_IDENTITY_TOKEN", "ANTHROPIC_IDENTITY_TOKEN_FILE",
		"ANTHROPIC_FEDERATION_RULE_ID", "ANTHROPIC_ORGANIZATION_ID",
		"ANTHROPIC_SERVICE_ACCOUNT_ID", "ANTHROPIC_WORKSPACE_ID",
	} {
		clearEnv(t, k)
	}
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
