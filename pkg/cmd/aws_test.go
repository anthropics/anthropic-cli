package cmd

import (
	"context"
	"io"
	"net/http"
	"regexp"
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

const (
	fakeAWSAccessKey    = "AKIAFAKEFAKEFAKE0000"
	fakeAWSSessionToken = "fake-aws-session-token"
	fakeAWSAPIKey       = "fake-aws-api-key"
)

// clearAWSEnv unsets every variable the AWS tier and the AWS credential chain
// read, so a developer machine or CI host with AWS configured can't leak into
// the tests.
func clearAWSEnv(t *testing.T) {
	t.Helper()
	for _, k := range []string{
		"AWS_REGION", "AWS_DEFAULT_REGION",
		"ANTHROPIC_AWS_WORKSPACE_ID", "ANTHROPIC_AWS_API_KEY", "ANTHROPIC_AWS_BASE_URL",
		"AWS_ACCESS_KEY_ID", "AWS_SECRET_ACCESS_KEY", "AWS_SESSION_TOKEN",
		"AWS_PROFILE", "AWS_DEFAULT_PROFILE",
	} {
		clearEnv(t, k)
	}
	t.Setenv("AWS_EC2_METADATA_DISABLED", "true")
}

// useFakeAWSCredentials puts static AWS credentials in the environment, the
// first source the default credential chain reads. Signing is local, so no
// request reaches AWS.
func useFakeAWSCredentials(t *testing.T, sessionToken string) {
	t.Helper()
	t.Setenv("AWS_ACCESS_KEY_ID", fakeAWSAccessKey)
	t.Setenv("AWS_SECRET_ACCESS_KEY", "fake-aws-secret-key")
	if sessionToken != "" {
		t.Setenv("AWS_SESSION_TOKEN", sessionToken)
	}
}

func setupAWSTest(t *testing.T) {
	t.Helper()
	clearCredentialEnv(t)
	clearGoogleCloudEnv(t)
	clearAWSEnv(t)
	resetWarnOnce(t)
	resetStdinCredential()
	t.Cleanup(resetStdinCredential)
}

// writeAWSTestProfile makes "p" the active profile, holding a static token
// for profile.invalid.
func writeAWSTestProfile(t *testing.T) {
	t.Helper()
	dir := t.TempDir()
	t.Setenv("ANTHROPIC_CONFIG_DIR", dir)
	credentials := config.ProfileCredentialsPath(dir, "p")
	require.NoError(t, config.SaveProfile(dir, "p", &config.Config{
		BaseURL: "https://profile.invalid",
		AuthenticationInfo: &config.AuthenticationInfo{
			Type:            config.AuthenticationTypeUserOAuth,
			CredentialsPath: credentials,
			UserOAuth:       &config.UserOAuth{},
		},
	}))
	require.NoError(t, config.WriteCredentials(credentials, config.Credentials{AccessToken: "profile-token"}))
	require.NoError(t, config.SetActiveProfile(dir, "p"))
}

// awsTestRoot is googleCloudTestRoot plus the AWS flags and `auth status`.
func awsTestRoot(action cli.ActionFunc) *cli.Command {
	root := googleCloudTestRoot(action)
	root.Flags = append(root.Flags,
		&cli.StringFlag{Name: "aws-region"},
		&cli.StringFlag{Name: "aws-workspace-id", Sources: cli.EnvVars("ANTHROPIC_AWS_WORKSPACE_ID")},
		&cli.BoolFlag{Name: "aws-api-key-stdin"},
		&cli.StringFlag{Name: "aws-base-url", Sources: cli.EnvVars("ANTHROPIC_AWS_BASE_URL")},
	)
	root.Commands = append(root.Commands, &cli.Command{
		Name: "auth", Commands: []*cli.Command{{Name: "status", Action: authStatus}},
	})
	return root
}

// runAWSLeaf runs a leaf command through getDefaultRequestOptions and returns
// the request it sends, captured at the transport, after signing.
func runAWSLeaf(t *testing.T, argv ...string) *http.Request {
	t.Helper()
	var captured *http.Request
	root := awsTestRoot(func(ctx context.Context, c *cli.Command) error {
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
	})
	require.NoError(t, root.Run(context.Background(), append(append([]string{"ant"}, argv...), "leaf")))
	require.NotNil(t, captured, "no request captured")
	return captured
}

var sigV4Authorization = regexp.MustCompile(`^AWS4-HMAC-SHA256 Credential=([^/]+)/\d{8}/([^/]+)/([^/]+)/aws4_request, SignedHeaders=([^,]+), Signature=[0-9a-f]{64}$`)

// assertSigV4Request checks the request went to the regional gateway, signed
// for the aws-external-anthropic service with the fake key, with the
// workspace header covered by the signature and no Anthropic credential.
func assertSigV4Request(t *testing.T, req *http.Request, region, workspaceID string) {
	t.Helper()
	assert.Equal(t, "https", req.URL.Scheme)
	assert.Equal(t, "aws-external-anthropic."+region+".api.aws", req.URL.Host)
	assert.Equal(t, "/v1/models", req.URL.Path)
	assert.Equal(t, []string{workspaceID}, req.Header.Values("anthropic-workspace-id"))
	assert.NotEmpty(t, req.Header.Get("X-Amz-Date"))
	assert.Empty(t, req.Header.Get("x-api-key"))

	m := sigV4Authorization.FindStringSubmatch(req.Header.Get("Authorization"))
	require.NotNil(t, m, "Authorization %q", req.Header.Get("Authorization"))
	assert.Equal(t, fakeAWSAccessKey, m[1])
	assert.Equal(t, region, m[2])
	assert.Equal(t, "aws-external-anthropic", m[3])
	signed := strings.Split(m[4], ";")
	assert.Contains(t, signed, "host")
	assert.Contains(t, signed, "x-amz-date")
	assert.Contains(t, signed, "anthropic-workspace-id")
}

// TestAWSTier pins when the AWS credential tier is selected and what it
// sends: the regional gateway URL, the workspace header, and either a SigV4
// signature or the AWS API key.
func TestAWSTier(t *testing.T) {
	t.Run("flags select the tier and sign with the AWS credential chain", func(t *testing.T) {
		setupAWSTest(t)
		useFakeAWSCredentials(t, "")
		req := runAWSLeaf(t, "--aws-region", "us-east-1", "--aws-workspace-id", "wrkspc_01")
		assertSigV4Request(t, req, "us-east-1", "wrkspc_01")
		assert.Empty(t, req.Header.Get("X-Amz-Security-Token"))
	})

	t.Run("ANTHROPIC_AWS_WORKSPACE_ID selects the tier; the region comes from AWS_DEFAULT_REGION", func(t *testing.T) {
		setupAWSTest(t)
		useFakeAWSCredentials(t, "")
		t.Setenv("AWS_DEFAULT_REGION", "eu-west-1")
		t.Setenv("ANTHROPIC_AWS_WORKSPACE_ID", "wrkspc_02")
		assertSigV4Request(t, runAWSLeaf(t), "eu-west-1", "wrkspc_02")
	})

	t.Run("temporary credentials send the session token", func(t *testing.T) {
		setupAWSTest(t)
		useFakeAWSCredentials(t, fakeAWSSessionToken)
		req := runAWSLeaf(t, "--aws-region", "us-east-1", "--aws-workspace-id", "wrkspc_01")
		assertSigV4Request(t, req, "us-east-1", "wrkspc_01")
		assert.Equal(t, fakeAWSSessionToken, req.Header.Get("X-Amz-Security-Token"))
	})

	t.Run("ANTHROPIC_AWS_API_KEY sends the key instead of signing", func(t *testing.T) {
		setupAWSTest(t)
		useFakeAWSCredentials(t, "")
		t.Setenv("ANTHROPIC_AWS_API_KEY", fakeAWSAPIKey)
		req := runAWSLeaf(t, "--aws-region", "us-east-1", "--aws-workspace-id", "wrkspc_01")
		assert.Equal(t, "aws-external-anthropic.us-east-1.api.aws", req.URL.Host)
		assert.Equal(t, "Bearer "+fakeAWSAPIKey, req.Header.Get("Authorization"))
		assert.Equal(t, "wrkspc_01", req.Header.Get("anthropic-workspace-id"))
		assert.Empty(t, req.Header.Get("x-api-key"))
	})

	t.Run("--aws-api-key-stdin sends the piped key", func(t *testing.T) {
		setupAWSTest(t)
		withStdinCredential(t, fakeAWSAPIKey+"\n")
		req := runAWSLeaf(t, "--aws-api-key-stdin", "--aws-region", "us-east-1", "--aws-workspace-id", "wrkspc_01")
		assert.Equal(t, "Bearer "+fakeAWSAPIKey, req.Header.Get("Authorization"))
		assert.Empty(t, req.Header.Get("x-api-key"))
	})

	t.Run("--workspace-id doesn't override the AWS workspace", func(t *testing.T) {
		setupAWSTest(t)
		useFakeAWSCredentials(t, "")
		req := runAWSLeaf(t, "--aws-region", "us-east-1", "--workspace-id", "wrkspc_global", "--aws-workspace-id", "wrkspc_aws")
		assertSigV4Request(t, req, "us-east-1", "wrkspc_aws")
	})

	t.Run("ANTHROPIC_CUSTOM_HEADERS is not sent", func(t *testing.T) {
		setupAWSTest(t)
		useFakeAWSCredentials(t, "")
		t.Setenv("ANTHROPIC_CUSTOM_HEADERS", "X-Custom: one\nAuthorization: Bearer from-headers")
		req := runAWSLeaf(t, "--aws-region", "us-east-1", "--aws-workspace-id", "wrkspc_01")
		assertSigV4Request(t, req, "us-east-1", "wrkspc_01")
		assert.Empty(t, req.Header.Get("X-Custom"))
	})

	t.Run("--aws-base-url replaces the gateway", func(t *testing.T) {
		setupAWSTest(t)
		useFakeAWSCredentials(t, "")
		req := runAWSLeaf(t, "--aws-region", "us-east-1", "--aws-workspace-id", "wrkspc_01", "--aws-base-url", "https://gateway.invalid")
		assert.Equal(t, "gateway.invalid", req.URL.Host)
		assert.True(t, strings.HasPrefix(req.Header.Get("Authorization"), "AWS4-HMAC-SHA256 "))
	})

	t.Run("--base-url never receives AWS credentials", func(t *testing.T) {
		setupAWSTest(t)
		useFakeAWSCredentials(t, "")
		t.Setenv("ANTHROPIC_BASE_URL", "https://anthropic-proxy.invalid")
		var req *http.Request
		stderr := captureStderr(t, func() {
			req = runAWSLeaf(t, "--aws-region", "us-east-1", "--aws-workspace-id", "wrkspc_01")
		})
		assertSigV4Request(t, req, "us-east-1", "wrkspc_01")
		assert.Contains(t, stderr, "--base-url / ANTHROPIC_BASE_URL is not used with AWS credentials")
		assert.Contains(t, stderr, "https://aws-external-anthropic.us-east-1.api.aws")
	})

	t.Run("ambient AWS settings don't select the tier", func(t *testing.T) {
		setupAWSTest(t)
		useFakeAWSCredentials(t, "")
		t.Setenv("AWS_REGION", "us-east-1")
		t.Setenv("AWS_PROFILE", "default")
		t.Setenv("ANTHROPIC_AWS_API_KEY", fakeAWSAPIKey)
		t.Setenv("ANTHROPIC_AWS_BASE_URL", "https://gateway.invalid")
		t.Setenv("ANTHROPIC_API_KEY", "fake-first-party-key")
		req := runAWSLeaf(t, "--workspace-id", "wrkspc_01")
		assert.Equal(t, "api.anthropic.com", req.URL.Host)
		assert.Equal(t, "Bearer fake-first-party-key", req.Header.Get("Authorization"))
		assert.Empty(t, req.Header.Get("x-api-key"))
	})
}

// TestAWSSelection pins what selects the AWS tier: an AWS region, workspace or
// API key the user set for this CLI, not AWS settings already on the machine.
func TestAWSSelection(t *testing.T) {
	for _, tc := range []struct {
		name  string
		env   map[string]string
		argv  []string
		stdin string
		want  bool
	}{
		{name: "--aws-region", argv: []string{"--aws-region", "us-east-1"}, want: true},
		{name: "--aws-workspace-id", argv: []string{"--aws-workspace-id", "wrkspc_01"}, want: true},
		{name: "ANTHROPIC_AWS_WORKSPACE_ID", env: map[string]string{"ANTHROPIC_AWS_WORKSPACE_ID": "wrkspc_01"}, want: true},
		{name: "ANTHROPIC_AWS_API_KEY", env: map[string]string{"ANTHROPIC_AWS_API_KEY": fakeAWSAPIKey}, want: true},
		{name: "--aws-api-key-stdin", argv: []string{"--aws-api-key-stdin"}, stdin: fakeAWSAPIKey + "\n", want: true},
		{name: "AWS_REGION", env: map[string]string{"AWS_REGION": "us-east-1"}},
		{name: "AWS_PROFILE", env: map[string]string{"AWS_PROFILE": "default"}},
		{name: "--aws-base-url", argv: []string{"--aws-base-url", "https://gateway.invalid"}},
		{name: "--workspace-id", argv: []string{"--workspace-id", "wrkspc_01"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			setupAWSTest(t)
			for k, v := range tc.env {
				t.Setenv(k, v)
			}
			if tc.stdin != "" {
				withStdinCredential(t, tc.stdin)
			}
			var got bool
			root := awsTestRoot(func(ctx context.Context, c *cli.Command) error {
				require.NoError(t, applyStdinCredential(c.Root()))
				got = awsFromRoot(c.Root()).AnySet()
				return nil
			})
			require.NoError(t, root.Run(context.Background(), append(append([]string{"ant"}, tc.argv...), "leaf")))
			assert.Equal(t, tc.want, got)
		})
	}
}

// TestAWSPrecedence pins where the AWS tier sits: below an API key, an auth
// token, a named profile, federation and Google Cloud; above the implicit
// profile. Same place as the Google Cloud tier.
func TestAWSPrecedence(t *testing.T) {
	awsArgs := []string{"--aws-region", "us-east-1", "--aws-workspace-id", "wrkspc_01"}

	t.Run("an API key beats AWS, with a notice", func(t *testing.T) {
		setupAWSTest(t)
		useFakeAWSCredentials(t, "")
		t.Setenv("ANTHROPIC_API_KEY", "fake-first-party-key")
		var req *http.Request
		stderr := captureStderr(t, func() { req = runAWSLeaf(t, awsArgs...) })
		assert.Equal(t, "Bearer fake-first-party-key", req.Header.Get("Authorization"))
		assert.Empty(t, req.Header.Get("x-api-key"))
		assert.Contains(t, stderr, "AWS (--aws-* / ANTHROPIC_AWS_*)")
		assert.Contains(t, stderr, "using ANTHROPIC_API_KEY per precedence")
	})

	t.Run("an AWS API key alone is not a second source", func(t *testing.T) {
		setupAWSTest(t)
		t.Setenv("ANTHROPIC_AWS_API_KEY", fakeAWSAPIKey)
		stderr := captureStderr(t, func() { runAWSLeaf(t, awsArgs...) })
		assert.Empty(t, stderr)
	})

	t.Run("Google Cloud beats AWS", func(t *testing.T) {
		setupAWSTest(t)
		useFakeGoogleCredentials(t)
		useFakeAWSCredentials(t, "")
		req := runAWSLeaf(t, append(awsArgs, "--google-cloud-project", "my-project", "--google-cloud-workspace-id", "wrkspc_gc")...)
		assertGoogleCloudRequest(t, req, "my-project", "global", "wrkspc_gc")
	})

	t.Run("AWS beats the implicit profile", func(t *testing.T) {
		setupAWSTest(t)
		useFakeAWSCredentials(t, "")
		writeAWSTestProfile(t)
		assertSigV4Request(t, runAWSLeaf(t, awsArgs...), "us-east-1", "wrkspc_01")
	})

	t.Run("a named profile beats AWS", func(t *testing.T) {
		setupAWSTest(t)
		useFakeAWSCredentials(t, "")
		writeAWSTestProfile(t)
		req := runAWSLeaf(t, append([]string{"--profile", "p"}, awsArgs...)...)
		assert.Equal(t, "profile.invalid", req.URL.Host)
		assert.Equal(t, "Bearer profile-token", req.Header.Get("Authorization"))
	})
}

// TestAWSCredentialsError guards the message a user sees when the AWS client
// can't be built: it must say where credentials come from and which flags
// address the gateway.
func TestAWSCredentialsError(t *testing.T) {
	setupAWSTest(t)
	_, err := awsPlatform{WorkspaceID: "wrkspc_01", APIKey: fakeAWSAPIKey}.requestOptions(context.Background())
	require.Error(t, err)
	for _, want := range []string{"credential chain", "ANTHROPIC_AWS_API_KEY", "--aws-region", "--aws-workspace-id"} {
		assert.Contains(t, err.Error(), want)
	}

	_, err = awsPlatform{Region: "us-east-1", WorkspaceID: "wrkspc_01", BaseURL: "gateway.invalid"}.requestOptions(context.Background())
	assert.ErrorContains(t, err, "--aws-base-url / ANTHROPIC_AWS_BASE_URL")
}

// TestAuthStatusAWS pins that `auth status` reports the tier the client would
// use, where requests go, and never a key.
func TestAuthStatusAWS(t *testing.T) {
	setupAWSTest(t)

	out, err := runStatus(t, "--aws-region", "us-east-1", "--aws-workspace-id", "wrkspc_01")
	require.NoError(t, err)
	assert.Regexp(t, `\(active\) \* AWS\s+AWS credential chain`, out)
	assert.Regexp(t, `region:\s+us-east-1`, out)
	assert.Regexp(t, `workspace_id:\s+wrkspc_01`, out)
	assert.Regexp(t, `\(active\) \* AWS gateway\s+https://aws-external-anthropic\.us-east-1\.api\.aws`, out)

	t.Run("the region comes from AWS_REGION", func(t *testing.T) {
		t.Setenv("AWS_REGION", "eu-west-1")
		out, err := runStatus(t, "--aws-workspace-id", "wrkspc_01")
		require.NoError(t, err)
		assert.Regexp(t, `region:\s+\(from AWS_REGION or AWS_DEFAULT_REGION\)`, out)
		assert.Regexp(t, `\(active\) \* AWS gateway\s+https://aws-external-anthropic\.eu-west-1\.api\.aws`, out)
	})

	t.Run("a generic base URL is listed but not active", func(t *testing.T) {
		out, err := runStatus(t, "--base-url", "https://anthropic-proxy.invalid", "--aws-region", "us-east-1", "--aws-workspace-id", "wrkspc_01")
		require.NoError(t, err)
		assert.Regexp(t, `\(active\) \* AWS gateway\s+https://aws-external-anthropic\.us-east-1\.api\.aws`, out)
		assert.Regexp(t, `\n\s+\* --base-url flag\s+https://anthropic-proxy.invalid`, out)
	})

	t.Run("the AWS base URL is the active one when set", func(t *testing.T) {
		out, err := runStatus(t, "--aws-workspace-id", "wrkspc_01", "--aws-base-url", "https://gateway.invalid")
		require.NoError(t, err)
		assert.Regexp(t, `\(active\) \* --aws-base-url / ANTHROPIC_AWS_BASE_URL\s+https://gateway.invalid`, out)
	})

	t.Run("an AWS API key is redacted", func(t *testing.T) {
		t.Setenv("ANTHROPIC_AWS_API_KEY", "fake-aws-api-key-0123456789abcdef")
		out, err := runStatus(t, "--aws-region", "us-east-1", "--aws-workspace-id", "wrkspc_01")
		require.NoError(t, err)
		assert.Regexp(t, `\(active\) \* AWS \(API key, ANTHROPIC_AWS_API_KEY\)\s+\(set, 33 chars\)`, out)
		assert.NotContains(t, out, "fake-aws")
		assert.NotContains(t, out, "0123456789abcdef")
	})

	t.Run("an AWS API key selects the tier", func(t *testing.T) {
		t.Setenv("ANTHROPIC_AWS_API_KEY", fakeAWSAPIKey)
		t.Setenv("AWS_REGION", "us-east-1")
		out, err := runStatus(t)
		require.NoError(t, err)
		assert.Regexp(t, `\(active\) \* AWS \(API key, ANTHROPIC_AWS_API_KEY\)`, out)
		assert.Regexp(t, `\(active\) \* AWS gateway\s+https://aws-external-anthropic\.us-east-1\.api\.aws`, out)
	})

	t.Run("the AWS workspace is the active one", func(t *testing.T) {
		out, err := runStatus(t, "--aws-region", "us-east-1", "--aws-workspace-id", "wrkspc_aws", "--workspace-id", "wrkspc_global")
		require.NoError(t, err)
		assert.Regexp(t, `\(active\) \* AWS workspace\s+wrkspc_aws`, out)
		assert.Regexp(t, `\n\s+\* --workspace-id / ANTHROPIC_WORKSPACE_ID\s+wrkspc_global`, out)
	})

	t.Run("an API key wins over AWS", func(t *testing.T) {
		out, err := captureStdout(t, func() error {
			t.Setenv("ANTHROPIC_API_KEY", "fake-first-party-key")
			return awsTestRoot(nil).Run(context.Background(), []string{"ant", "--aws-region", "us-east-1", "auth", "status"})
		})
		require.NoError(t, err)
		assert.Regexp(t, `\(active\) \* ANTHROPIC_API_KEY`, out)
		assert.Regexp(t, `\n\s+\* AWS\s`, out)
	})
}

// TestAWSInUse pins what `ant apply` records as the host, how it names the
// credentials, and whether it treats the active profile as the one in use.
func TestAWSInUse(t *testing.T) {
	var credentials string
	probe := func(t *testing.T, argv ...string) (inUse bool, baseURL string, profile *config.Config) {
		t.Helper()
		root := awsTestRoot(func(ctx context.Context, c *cli.Command) error {
			inUse, baseURL, profile = awsInUse(c), currentBaseURL(c), profileInUse(c)
			credentials = describeOrigin(c, core.Origin{BaseURL: baseURL}).Credentials
			return nil
		})
		require.NoError(t, root.Run(context.Background(), append(append([]string{"ant"}, argv...), "leaf")))
		return inUse, baseURL, profile
	}

	t.Run("regional gateway", func(t *testing.T) {
		setupAWSTest(t)
		inUse, baseURL, _ := probe(t, "--aws-region", "us-east-1")
		assert.True(t, inUse)
		assert.Equal(t, "https://aws-external-anthropic.us-east-1.api.aws", baseURL)
		assert.Equal(t, "AWS (AWS credential chain)", credentials)
	})

	t.Run("--aws-base-url, API key", func(t *testing.T) {
		setupAWSTest(t)
		t.Setenv("ANTHROPIC_AWS_API_KEY", fakeAWSAPIKey)
		t.Setenv("ANTHROPIC_BASE_URL", "https://anthropic-proxy.invalid")
		_, baseURL, _ := probe(t, "--aws-region", "us-east-1", "--aws-base-url", "https://gateway.invalid")
		assert.Equal(t, "https://gateway.invalid", baseURL)
		assert.Equal(t, "AWS (API key)", credentials)
	})

	t.Run("not in use when an API key wins", func(t *testing.T) {
		setupAWSTest(t)
		t.Setenv("ANTHROPIC_API_KEY", "fake-first-party-key")
		inUse, baseURL, _ := probe(t, "--aws-region", "us-east-1")
		assert.False(t, inUse)
		assert.Equal(t, defaultBaseURL, baseURL)
	})

	t.Run("the implicit profile is not in use", func(t *testing.T) {
		setupAWSTest(t)
		writeAWSTestProfile(t)
		inUse, _, profile := probe(t, "--aws-region", "us-east-1")
		assert.True(t, inUse)
		assert.Nil(t, profile)
	})
}
