package cmd

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/anthropics/anthropic-cli/v2/internal/declarative/core"
	"github.com/anthropics/anthropic-sdk-go/v2"
	"github.com/anthropics/anthropic-sdk-go/v2/config"
	"github.com/anthropics/anthropic-sdk-go/v2/option"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli/v3"
)

const fakeBedrockAPIKey = "fake-bedrock-api-key"

func clearBedrockEnv(t *testing.T) {
	t.Helper()
	for _, k := range []string{
		"ANTHROPIC_BEDROCK_REGION", "AWS_BEARER_TOKEN_BEDROCK", "ANTHROPIC_BEDROCK_MANTLE_BASE_URL",
	} {
		clearEnv(t, k)
	}
}

func setupBedrockTest(t *testing.T) {
	t.Helper()
	setupFoundryTest(t)
	clearBedrockEnv(t)
}

func bedrockTestRoot(action cli.ActionFunc) *cli.Command {
	root := foundryTestRoot(action)
	root.Flags = append(root.Flags,
		&cli.StringFlag{Name: "bedrock-region", Sources: cli.EnvVars("ANTHROPIC_BEDROCK_REGION")},
		&cli.BoolFlag{Name: "bedrock-api-key-stdin"},
		&cli.StringFlag{Name: "bedrock-base-url", Sources: cli.EnvVars("ANTHROPIC_BEDROCK_MANTLE_BASE_URL")},
	)
	return root
}

// runBedrockLeaf runs a leaf command through getDefaultRequestOptions and
// returns the last request it sends, captured at the transport, after
// signing. The transport comes first and answers with a token, so
// federation's exchange stays here too.
func runBedrockLeaf(t *testing.T, argv ...string) *http.Request {
	t.Helper()
	var captured *http.Request
	root := bedrockTestRoot(func(ctx context.Context, c *cli.Command) error {
		opts := append([]option.RequestOption{
			option.WithHTTPClient(&http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
				captured = req
				return &http.Response{
					StatusCode: http.StatusOK,
					Header:     http.Header{"Content-Type": []string{"application/json"}},
					Body:       io.NopCloser(strings.NewReader(`{"access_token":"fake-federation-token","token_type":"Bearer","expires_in":3600}`)),
					Request:    req,
				}, nil
			})}),
			option.WithMaxRetries(0),
		}, getDefaultRequestOptions(c)...)
		client := anthropic.NewClient(opts...)
		return client.Post(ctx, "/v1/messages", nil, nil)
	})
	require.NoError(t, root.Run(context.Background(), append(append([]string{"ant"}, argv...), "leaf")))
	require.NotNil(t, captured, "no request captured")
	return captured
}

// assertBedrockSigV4Request checks the request went to the regional Bedrock
// endpoint, signed for the bedrock-mantle service with the fake key, with no
// Anthropic credential and no workspace header.
func assertBedrockSigV4Request(t *testing.T, req *http.Request, region string) {
	t.Helper()
	assert.Equal(t, "https://bedrock-mantle."+region+".api.aws/anthropic/v1/messages", req.URL.String())
	assert.Empty(t, req.Header.Values("anthropic-workspace-id"))
	assert.NotEmpty(t, req.Header.Get("X-Amz-Date"))
	assert.Empty(t, req.Header.Get("x-api-key"))

	m := sigV4Authorization.FindStringSubmatch(req.Header.Get("Authorization"))
	require.NotNil(t, m, "Authorization %q", req.Header.Get("Authorization"))
	assert.Equal(t, fakeAWSAccessKey, m[1])
	assert.Equal(t, region, m[2])
	assert.Equal(t, "bedrock-mantle", m[3])
	signed := strings.Split(m[4], ";")
	assert.Contains(t, signed, "host")
	assert.Contains(t, signed, "x-amz-date")
}

// assertBedrockAPIKey checks the key is the only credential on the request.
func assertBedrockAPIKey(t *testing.T, req *http.Request) {
	t.Helper()
	assert.Equal(t, []string{"Bearer " + fakeBedrockAPIKey}, req.Header.Values("Authorization"))
	assert.Empty(t, req.Header.Values("x-api-key"))
}

// TestBedrockTier pins what the Bedrock credential tier sends: the regional
// endpoint URL, and either a SigV4 signature or the Bedrock API key.
func TestBedrockTier(t *testing.T) {
	t.Run("--bedrock-region signs with the AWS credential chain", func(t *testing.T) {
		setupBedrockTest(t)
		useFakeAWSCredentials(t, "")
		req := runBedrockLeaf(t, "--bedrock-region", "us-east-1")
		assertBedrockSigV4Request(t, req, "us-east-1")
		assert.Empty(t, req.Header.Get("X-Amz-Security-Token"))
	})

	t.Run("ANTHROPIC_BEDROCK_REGION selects the tier", func(t *testing.T) {
		setupBedrockTest(t)
		useFakeAWSCredentials(t, "")
		t.Setenv("ANTHROPIC_BEDROCK_REGION", "ap-northeast-1")
		assertBedrockSigV4Request(t, runBedrockLeaf(t), "ap-northeast-1")
	})

	t.Run("a piped key selects the tier; the region comes from AWS_DEFAULT_REGION", func(t *testing.T) {
		setupBedrockTest(t)
		t.Setenv("AWS_DEFAULT_REGION", "eu-west-1")
		withStdinCredential(t, fakeBedrockAPIKey+"\n")
		req := runBedrockLeaf(t, "--bedrock-api-key-stdin")
		assert.Equal(t, "https://bedrock-mantle.eu-west-1.api.aws/anthropic/v1/messages", req.URL.String())
		assertBedrockAPIKey(t, req)
	})

	t.Run("temporary credentials send the session token", func(t *testing.T) {
		setupBedrockTest(t)
		useFakeAWSCredentials(t, fakeAWSSessionToken)
		req := runBedrockLeaf(t, "--bedrock-region", "us-east-1")
		assertBedrockSigV4Request(t, req, "us-east-1")
		assert.Equal(t, fakeAWSSessionToken, req.Header.Get("X-Amz-Security-Token"))
	})

	t.Run("AWS_BEARER_TOKEN_BEDROCK sends the key instead of signing", func(t *testing.T) {
		setupBedrockTest(t)
		useFakeAWSCredentials(t, "")
		t.Setenv("AWS_BEARER_TOKEN_BEDROCK", fakeBedrockAPIKey)
		req := runBedrockLeaf(t, "--bedrock-region", "us-east-1")
		assert.Equal(t, "bedrock-mantle.us-east-1.api.aws", req.URL.Host)
		assertBedrockAPIKey(t, req)
	})

	t.Run("ANTHROPIC_CUSTOM_HEADERS is not sent", func(t *testing.T) {
		setupBedrockTest(t)
		useFakeAWSCredentials(t, "")
		t.Setenv("ANTHROPIC_CUSTOM_HEADERS", "X-Custom: one\nAuthorization: Bearer from-headers")
		req := runBedrockLeaf(t, "--bedrock-region", "us-east-1")
		assertBedrockSigV4Request(t, req, "us-east-1")
		assert.Empty(t, req.Header.Get("X-Custom"))
	})

	t.Run("--bedrock-base-url replaces the endpoint", func(t *testing.T) {
		setupBedrockTest(t)
		useFakeAWSCredentials(t, "")
		req := runBedrockLeaf(t, "--bedrock-region", "us-east-1", "--bedrock-base-url", "https://endpoint.invalid")
		assert.Equal(t, "endpoint.invalid", req.URL.Host)
		assert.True(t, strings.HasPrefix(req.Header.Get("Authorization"), "AWS4-HMAC-SHA256 "))
	})

	t.Run("--base-url never receives AWS credentials", func(t *testing.T) {
		setupBedrockTest(t)
		useFakeAWSCredentials(t, "")
		t.Setenv("ANTHROPIC_BASE_URL", "https://anthropic-proxy.invalid")
		var req *http.Request
		stderr := captureStderr(t, func() {
			req = runBedrockLeaf(t, "--bedrock-region", "us-east-1")
		})
		assertBedrockSigV4Request(t, req, "us-east-1")
		assert.Contains(t, stderr, "--base-url / ANTHROPIC_BASE_URL is not used with Amazon Bedrock credentials")
		assert.Contains(t, stderr, "https://bedrock-mantle.us-east-1.api.aws/anthropic")
	})

	t.Run("ambient AWS and Bedrock settings don't select the tier", func(t *testing.T) {
		setupBedrockTest(t)
		useFakeAWSCredentials(t, "")
		t.Setenv("AWS_REGION", "us-east-1")
		t.Setenv("AWS_PROFILE", "default")
		t.Setenv("AWS_BEARER_TOKEN_BEDROCK", fakeBedrockAPIKey)
		t.Setenv("ANTHROPIC_BEDROCK_MANTLE_BASE_URL", "https://endpoint.invalid")
		t.Setenv("ANTHROPIC_API_KEY", "fake-first-party-key")
		req := runBedrockLeaf(t)
		assert.Equal(t, "api.anthropic.com", req.URL.Host)
		assert.Equal(t, []string{"Bearer fake-first-party-key"}, req.Header.Values("Authorization"))
		assert.Empty(t, req.Header.Values("x-api-key"))
	})
}

// TestBedrockSelection pins what selects the Bedrock tier: a region or a
// piped key the user set for this CLI, not AWS or Bedrock settings already on
// the machine.
func TestBedrockSelection(t *testing.T) {
	for _, tc := range []struct {
		name  string
		env   map[string]string
		argv  []string
		stdin string
		want  bool
	}{
		{name: "--bedrock-region", argv: []string{"--bedrock-region", "us-east-1"}, want: true},
		{name: "ANTHROPIC_BEDROCK_REGION", env: map[string]string{"ANTHROPIC_BEDROCK_REGION": "us-east-1"}, want: true},
		{name: "--bedrock-api-key-stdin", argv: []string{"--bedrock-api-key-stdin"}, stdin: fakeBedrockAPIKey + "\n", want: true},
		{name: "AWS_BEARER_TOKEN_BEDROCK", env: map[string]string{"AWS_BEARER_TOKEN_BEDROCK": fakeBedrockAPIKey}},
		{name: "AWS_REGION", env: map[string]string{"AWS_REGION": "us-east-1"}},
		{name: "AWS_PROFILE", env: map[string]string{"AWS_PROFILE": "default"}},
		{name: "--bedrock-base-url", argv: []string{"--bedrock-base-url", "https://endpoint.invalid"}},
		{name: "ANTHROPIC_BEDROCK_MANTLE_BASE_URL", env: map[string]string{"ANTHROPIC_BEDROCK_MANTLE_BASE_URL": "https://endpoint.invalid"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			setupBedrockTest(t)
			for k, v := range tc.env {
				t.Setenv(k, v)
			}
			if tc.stdin != "" {
				withStdinCredential(t, tc.stdin)
			}
			var got bool
			root := bedrockTestRoot(func(ctx context.Context, c *cli.Command) error {
				require.NoError(t, applyStdinCredential(c.Root()))
				got = bedrockFromRoot(c.Root()).AnySet()
				return nil
			})
			require.NoError(t, root.Run(context.Background(), append(append([]string{"ant"}, tc.argv...), "leaf")))
			assert.Equal(t, tc.want, got)
		})
	}
}

// TestBedrockPrecedence pins where the Bedrock tier sits: below every other
// tier the user addressed, above the implicit profile. Each case checks the
// request, the multi-source notice, `auth status` and bedrockInUse agree.
func TestBedrockPrecedence(t *testing.T) {
	bedrockArgs := []string{"--bedrock-region", "us-east-1"}
	setup := func(t *testing.T) {
		setupBedrockTest(t)
		useFakeAWSCredentials(t, "")
	}

	for _, tc := range []struct {
		name                             string
		other                            func(*testing.T)
		argv                             []string
		wantHost, wantActive, wantWinner string
		wantAuthorization                string
		wantBedrock                      bool
	}{
		{
			name:              "an API key",
			other:             func(t *testing.T) { t.Setenv("ANTHROPIC_API_KEY", "fake-first-party-key") },
			wantHost:          "api.anthropic.com",
			wantActive:        "ANTHROPIC_API_KEY",
			wantWinner:        "ANTHROPIC_API_KEY",
			wantAuthorization: "Bearer fake-first-party-key",
		},
		{
			name:              "an auth token",
			other:             func(t *testing.T) { t.Setenv("ANTHROPIC_AUTH_TOKEN", "fake-first-party-token") },
			wantHost:          "api.anthropic.com",
			wantActive:        "ANTHROPIC_AUTH_TOKEN",
			wantWinner:        "ANTHROPIC_AUTH_TOKEN",
			wantAuthorization: "Bearer fake-first-party-token",
		},
		{
			name:              "a named profile",
			other:             writeAWSTestProfile,
			argv:              []string{"--profile", "p"},
			wantHost:          "profile.invalid",
			wantActive:        "Profile (user_oauth) [via --profile/ANTHROPIC_PROFILE]",
			wantWinner:        "profile from --profile / ANTHROPIC_PROFILE",
			wantAuthorization: "Bearer profile-token",
		},
		{
			name:              "federation",
			argv:              []string{"--identity-token", "fake-identity-token", "--federation-rule", "fdrl_01", "--organization-id", "org_01"},
			wantHost:          "api.anthropic.com",
			wantActive:        "Federation (jwt-bearer)",
			wantWinner:        "federation env",
			wantAuthorization: "Bearer fake-federation-token",
		},
		{
			name:              "Google Cloud",
			other:             useFakeGoogleCredentials,
			argv:              []string{"--google-cloud-project", "my-project", "--google-cloud-workspace-id", "wrkspc_gc"},
			wantHost:          "claude.googleapis.com",
			wantActive:        "Google Cloud",
			wantWinner:        "Google Cloud (--google-cloud-* / ANTHROPIC_GOOGLE_CLOUD_*)",
			wantAuthorization: "Bearer " + fakeGoogleToken,
		},
		{
			name:       "AWS",
			argv:       []string{"--aws-region", "us-east-1", "--aws-workspace-id", "wrkspc_01"},
			wantHost:   "aws-external-anthropic.us-east-1.api.aws",
			wantActive: "AWS",
			wantWinner: "AWS (--aws-* / ANTHROPIC_AWS_*)",
		},
		{
			name:       "Foundry",
			other:      func(t *testing.T) { t.Setenv("ANTHROPIC_FOUNDRY_API_KEY", fakeFoundryAPIKey) },
			argv:       []string{"--foundry-resource", "my-resource"},
			wantHost:   "my-resource.services.ai.azure.com",
			wantActive: "Foundry (API key, ANTHROPIC_FOUNDRY_API_KEY)",
			wantWinner: "Foundry (--foundry-* / ANTHROPIC_FOUNDRY_*)",
		},
		{
			name:        "the implicit profile",
			other:       writeAWSTestProfile,
			wantHost:    "bedrock-mantle.us-east-1.api.aws",
			wantActive:  "Bedrock",
			wantWinner:  "Bedrock (--bedrock-* / ANTHROPIC_BEDROCK_REGION)",
			wantBedrock: true,
		},
	} {
		t.Run("Bedrock and "+tc.name, func(t *testing.T) {
			setup(t)
			if tc.other != nil {
				tc.other(t)
			}
			argv := append(bedrockArgs, tc.argv...)

			var req *http.Request
			stderr := captureStderr(t, func() { req = runBedrockLeaf(t, argv...) })
			status, err := runStatusKeepingCredentialEnv(t, argv...)
			require.NoError(t, err)
			var inUse bool
			root := bedrockTestRoot(func(ctx context.Context, c *cli.Command) error {
				inUse = bedrockInUse(c)
				return nil
			})
			require.NoError(t, root.Run(context.Background(), append(append([]string{"ant"}, argv...), "leaf")))

			assert.Equal(t, tc.wantHost, req.URL.Host)
			if tc.wantBedrock {
				assertBedrockSigV4Request(t, req, "us-east-1")
			}
			if tc.wantAuthorization != "" {
				assert.Equal(t, []string{tc.wantAuthorization}, req.Header.Values("Authorization"))
			}
			assert.Contains(t, stderr, "Bedrock (--bedrock-* / ANTHROPIC_BEDROCK_REGION)")
			assert.Contains(t, stderr, "using "+tc.wantWinner+" per precedence")
			assert.Contains(t, status, fmt.Sprintf("(active) * %-46s ", tc.wantActive))
			assert.Equal(t, tc.wantBedrock, inUse)
		})
	}

	t.Run("a region and a key are one source", func(t *testing.T) {
		setupBedrockTest(t)
		t.Setenv("AWS_BEARER_TOKEN_BEDROCK", fakeBedrockAPIKey)

		stderr := captureStderr(t, func() { runBedrockLeaf(t, bedrockArgs...) })

		assert.Empty(t, stderr)
	})

	t.Run("an AWS API key goes to Claude Platform on AWS, not Bedrock", func(t *testing.T) {
		setupBedrockTest(t)
		t.Setenv("AWS_REGION", "us-east-1")
		t.Setenv("ANTHROPIC_AWS_API_KEY", fakeAWSAPIKey)
		t.Setenv("ANTHROPIC_AWS_WORKSPACE_ID", "wrkspc_01")
		req := runBedrockLeaf(t, bedrockArgs...)
		assert.Equal(t, "aws-external-anthropic.us-east-1.api.aws", req.URL.Host)
		assert.Equal(t, []string{"Bearer " + fakeAWSAPIKey}, req.Header.Values("Authorization"))
	})
}

// TestBedrockCredentialsError guards the message a user sees when the
// Bedrock client can't be built: it must say where credentials come from and
// which flags address the endpoint.
func TestBedrockCredentialsError(t *testing.T) {
	setupBedrockTest(t)
	_, err := bedrockPlatform{APIKey: fakeBedrockAPIKey}.requestOptions(context.Background())
	require.Error(t, err)
	for _, want := range []string{"credential chain", "AWS_BEARER_TOKEN_BEDROCK", "--bedrock-region", "--bedrock-base-url"} {
		assert.Contains(t, err.Error(), want)
	}

	_, err = bedrockPlatform{Region: "us-east-1", BaseURL: "endpoint.invalid"}.requestOptions(context.Background())
	assert.ErrorContains(t, err, "--bedrock-base-url / ANTHROPIC_BEDROCK_MANTLE_BASE_URL")
}

// TestAuthStatusBedrock pins that `auth status` reports the tier the client
// would use, where requests go, and no part of a key.
func TestAuthStatusBedrock(t *testing.T) {
	t.Run("the AWS credential chain", func(t *testing.T) {
		setupBedrockTest(t)

		out, err := runStatus(t, "--bedrock-region", "us-east-1")

		require.NoError(t, err)
		assert.Regexp(t, `\(active\) \* Bedrock\s+AWS credential chain`, out)
		assert.Regexp(t, `region:\s+us-east-1`, out)
		assert.Regexp(t, `\(active\) \* Bedrock endpoint\s+https://bedrock-mantle\.us-east-1\.api\.aws/anthropic`, out)
	})

	t.Run("a piped key, with the region from AWS_REGION", func(t *testing.T) {
		setupBedrockTest(t)
		t.Setenv("AWS_REGION", "eu-west-1")
		withStdinCredential(t, fakeBedrockAPIKey+"\n")

		out, err := runStatus(t, "--bedrock-api-key-stdin")

		require.NoError(t, err)
		assert.Regexp(t, `\(active\) \* Bedrock \(API key, --bedrock-api-key-stdin\)`, out)
		assert.Regexp(t, `region:\s+\(from AWS_REGION or AWS_DEFAULT_REGION\)`, out)
		assert.Regexp(t, `\(active\) \* Bedrock endpoint\s+https://bedrock-mantle\.eu-west-1\.api\.aws/anthropic`, out)
	})

	t.Run("an API key is shown only by its length", func(t *testing.T) {
		setupBedrockTest(t)
		t.Setenv("AWS_BEARER_TOKEN_BEDROCK", "fake-bedrock-key-0123456789abcdef")

		out, err := runStatus(t, "--bedrock-region", "us-east-1")

		require.NoError(t, err)
		assert.Regexp(t, `\(active\) \* Bedrock \(API key, AWS_BEARER_TOKEN_BEDROCK\)\s+\(set, 33 chars\)`, out)
		assert.NotContains(t, out, "fake-bed")
		assert.NotContains(t, out, "0123456789abcdef")
	})

	t.Run("a generic base URL is listed but not active", func(t *testing.T) {
		setupBedrockTest(t)

		out, err := runStatus(t, "--base-url", "https://anthropic-proxy.invalid", "--bedrock-region", "us-east-1")

		require.NoError(t, err)
		assert.Regexp(t, `\(active\) \* Bedrock endpoint\s+https://bedrock-mantle\.us-east-1\.api\.aws/anthropic`, out)
		assert.Regexp(t, `\n\s+\* --base-url flag\s+https://anthropic-proxy.invalid`, out)
	})

	t.Run("the Bedrock base URL is the active one when set", func(t *testing.T) {
		setupBedrockTest(t)

		out, err := runStatus(t, "--bedrock-region", "us-east-1", "--bedrock-base-url", "https://endpoint.invalid")

		require.NoError(t, err)
		assert.Regexp(t, `\(active\) \* --bedrock-base-url / ANTHROPIC_BEDROCK_MANTLE_BASE_URL\s+https://endpoint.invalid`, out)
	})
}

// TestBedrockInUse pins what `ant apply` records as the host, how it names
// the credentials, and whether it treats the active profile as the one in use.
func TestBedrockInUse(t *testing.T) {
	var credentials string
	probe := func(t *testing.T, argv ...string) (inUse bool, baseURL string, profile *config.Config) {
		t.Helper()
		root := bedrockTestRoot(func(ctx context.Context, c *cli.Command) error {
			inUse, baseURL, profile = bedrockInUse(c), currentBaseURL(c), profileInUse(c)
			credentials = describeOrigin(c, core.Origin{BaseURL: baseURL}).Credentials
			return nil
		})
		require.NoError(t, root.Run(context.Background(), append(append([]string{"ant"}, argv...), "leaf")))
		return inUse, baseURL, profile
	}

	t.Run("regional endpoint", func(t *testing.T) {
		setupBedrockTest(t)
		inUse, baseURL, _ := probe(t, "--bedrock-region", "us-east-1")
		assert.True(t, inUse)
		assert.Equal(t, "https://bedrock-mantle.us-east-1.api.aws/anthropic", baseURL)
		assert.Equal(t, "Bedrock (AWS credential chain)", credentials)
	})

	t.Run("--bedrock-base-url, API key", func(t *testing.T) {
		setupBedrockTest(t)
		t.Setenv("AWS_BEARER_TOKEN_BEDROCK", fakeBedrockAPIKey)
		t.Setenv("ANTHROPIC_BASE_URL", "https://anthropic-proxy.invalid")
		_, baseURL, _ := probe(t, "--bedrock-region", "us-east-1", "--bedrock-base-url", "https://endpoint.invalid")
		assert.Equal(t, "https://endpoint.invalid", baseURL)
		assert.Equal(t, "Bedrock (API key)", credentials)
	})

	t.Run("not in use when Foundry wins", func(t *testing.T) {
		setupBedrockTest(t)
		t.Setenv("ANTHROPIC_FOUNDRY_API_KEY", fakeFoundryAPIKey)
		inUse, baseURL, _ := probe(t, "--foundry-resource", "my-resource", "--bedrock-region", "us-east-1")
		assert.False(t, inUse)
		assert.Equal(t, "https://my-resource.services.ai.azure.com/anthropic", baseURL)
		assert.Equal(t, "Foundry (API key)", credentials)
	})

	t.Run("not in use when an API key wins", func(t *testing.T) {
		setupBedrockTest(t)
		t.Setenv("ANTHROPIC_API_KEY", "fake-first-party-key")
		inUse, baseURL, _ := probe(t, "--bedrock-region", "us-east-1")
		assert.False(t, inUse)
		assert.Equal(t, defaultBaseURL, baseURL)
	})

	t.Run("the implicit profile is not in use", func(t *testing.T) {
		setupBedrockTest(t)
		writeAWSTestProfile(t)
		inUse, _, profile := probe(t, "--bedrock-region", "us-east-1")
		assert.True(t, inUse)
		assert.Nil(t, profile)
	})
}
