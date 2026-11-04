package cmd

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os/exec"
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
	fakeFoundryAPIKey    = "fake-foundry-api-key"
	fakeFoundryAuthToken = "fake-foundry-entra-id-token"
)

func clearFoundryEnv(t *testing.T) {
	t.Helper()
	for _, k := range []string{
		"ANTHROPIC_FOUNDRY_RESOURCE", "ANTHROPIC_FOUNDRY_API_KEY",
		"ANTHROPIC_FOUNDRY_AUTH_TOKEN", "ANTHROPIC_FOUNDRY_BASE_URL",
	} {
		clearEnv(t, k)
	}
}

func setupFoundryTest(t *testing.T) {
	t.Helper()
	setupAWSTest(t)
	clearFoundryEnv(t)
}

func foundryTestRoot(action cli.ActionFunc) *cli.Command {
	root := awsTestRoot(action)
	root.Flags = append(root.Flags,
		&cli.StringFlag{Name: "foundry-resource"},
		&cli.BoolFlag{Name: "foundry-api-key-stdin"},
		&cli.BoolFlag{Name: "foundry-auth-token-stdin"},
		&cli.StringFlag{Name: "foundry-base-url", Sources: cli.EnvVars("ANTHROPIC_FOUNDRY_BASE_URL")},
	)
	return root
}

// runFoundryLeaf runs a leaf command through getDefaultRequestOptions and
// returns the last request it sends, captured at the transport. The transport
// comes first and answers with a token, so federation's exchange stays here too.
func runFoundryLeaf(t *testing.T, argv ...string) *http.Request {
	t.Helper()
	var captured *http.Request
	root := foundryTestRoot(func(ctx context.Context, c *cli.Command) error {
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

func assertFoundryURL(t *testing.T, req *http.Request, resource string) {
	t.Helper()
	assert.Equal(t, "https://"+resource+".services.ai.azure.com/anthropic/v1/messages", req.URL.String())
}

// assertFoundryAPIKey checks the key is the only credential on the request.
func assertFoundryAPIKey(t *testing.T, req *http.Request) {
	t.Helper()
	assert.Equal(t, []string{fakeFoundryAPIKey}, req.Header.Values("x-api-key"))
	assert.Empty(t, req.Header.Values("Authorization"))
}

// assertFoundryAuthToken checks the token is the only credential on the request.
func assertFoundryAuthToken(t *testing.T, req *http.Request) {
	t.Helper()
	assert.Equal(t, []string{"Bearer " + fakeFoundryAuthToken}, req.Header.Values("Authorization"))
	assert.Empty(t, req.Header.Values("x-api-key"))
}

// TestFoundryTier pins what the Foundry credential tier sends: the resource's
// URL, and either the API key or the Entra ID token.
func TestFoundryTier(t *testing.T) {
	t.Run("ANTHROPIC_FOUNDRY_API_KEY is sent as x-api-key", func(t *testing.T) {
		setupFoundryTest(t)
		t.Setenv("ANTHROPIC_FOUNDRY_API_KEY", fakeFoundryAPIKey)

		req := runFoundryLeaf(t, "--foundry-resource", "my-resource")

		assertFoundryURL(t, req, "my-resource")
		assertFoundryAPIKey(t, req)
	})

	t.Run("ANTHROPIC_FOUNDRY_AUTH_TOKEN is sent as a bearer token", func(t *testing.T) {
		setupFoundryTest(t)
		t.Setenv("ANTHROPIC_FOUNDRY_AUTH_TOKEN", fakeFoundryAuthToken)

		req := runFoundryLeaf(t, "--foundry-resource", "my-resource")

		assertFoundryURL(t, req, "my-resource")
		assertFoundryAuthToken(t, req)
	})

	t.Run("ANTHROPIC_FOUNDRY_RESOURCE names the resource", func(t *testing.T) {
		setupFoundryTest(t)
		t.Setenv("ANTHROPIC_FOUNDRY_RESOURCE", "my-resource")
		t.Setenv("ANTHROPIC_FOUNDRY_API_KEY", fakeFoundryAPIKey)

		req := runFoundryLeaf(t)

		assertFoundryURL(t, req, "my-resource")
		assertFoundryAPIKey(t, req)
	})

	t.Run("--foundry-api-key-stdin sends the piped key", func(t *testing.T) {
		setupFoundryTest(t)
		withStdinCredential(t, fakeFoundryAPIKey+"\n")

		req := runFoundryLeaf(t, "--foundry-api-key-stdin", "--foundry-resource", "my-resource")

		assertFoundryAPIKey(t, req)
	})

	t.Run("--foundry-auth-token-stdin sends the piped token", func(t *testing.T) {
		setupFoundryTest(t)
		withStdinCredential(t, fakeFoundryAuthToken+"\n")

		req := runFoundryLeaf(t, "--foundry-auth-token-stdin", "--foundry-resource", "my-resource")

		assertFoundryAuthToken(t, req)
	})

	t.Run("--foundry-base-url takes the place of the resource", func(t *testing.T) {
		setupFoundryTest(t)
		t.Setenv("ANTHROPIC_FOUNDRY_API_KEY", fakeFoundryAPIKey)

		req := runFoundryLeaf(t, "--foundry-base-url", "https://gateway.invalid/anthropic")

		assert.Equal(t, "https://gateway.invalid/anthropic/v1/messages", req.URL.String())
		assertFoundryAPIKey(t, req)
	})

	t.Run("ANTHROPIC_CUSTOM_HEADERS is not sent", func(t *testing.T) {
		setupFoundryTest(t)
		t.Setenv("ANTHROPIC_FOUNDRY_AUTH_TOKEN", fakeFoundryAuthToken)
		t.Setenv("ANTHROPIC_CUSTOM_HEADERS", "X-Custom: one\nAuthorization: Bearer from-headers")

		req := runFoundryLeaf(t, "--foundry-resource", "my-resource")

		assertFoundryAuthToken(t, req)
		assert.Empty(t, req.Header.Get("X-Custom"))
	})

	t.Run("--base-url never receives Foundry credentials", func(t *testing.T) {
		setupFoundryTest(t)
		t.Setenv("ANTHROPIC_FOUNDRY_API_KEY", fakeFoundryAPIKey)
		t.Setenv("ANTHROPIC_BASE_URL", "https://anthropic-proxy.invalid")

		var req *http.Request
		stderr := captureStderr(t, func() { req = runFoundryLeaf(t, "--foundry-resource", "my-resource") })

		assertFoundryURL(t, req, "my-resource")
		assert.Contains(t, stderr, "--base-url / ANTHROPIC_BASE_URL is not used with Foundry credentials")
		assert.Contains(t, stderr, "https://my-resource.services.ai.azure.com/anthropic")
	})
}

// TestFoundryKeyAndToken pins which credential is used when both are set: the
// piped one, else the token, with a note.
func TestFoundryKeyAndToken(t *testing.T) {
	for _, tc := range []struct {
		name                     string
		env                      map[string]string
		stdinFlag, piped         string
		assertCredential         func(*testing.T, *http.Request)
		wantStderr               string
		wantActive, wantInactive string
	}{
		{
			name:             "both in the environment",
			env:              map[string]string{"ANTHROPIC_FOUNDRY_API_KEY": fakeFoundryAPIKey, "ANTHROPIC_FOUNDRY_AUTH_TOKEN": fakeFoundryAuthToken},
			assertCredential: assertFoundryAuthToken,
			wantStderr:       "Note: ANTHROPIC_FOUNDRY_API_KEY and ANTHROPIC_FOUNDRY_AUTH_TOKEN are both set in your environment; using ANTHROPIC_FOUNDRY_AUTH_TOKEN. Run `ant auth status` for details.\n",
			wantActive:       "Foundry (Entra ID token, ANTHROPIC_FOUNDRY_AUTH_TOKEN)",
			wantInactive:     "Foundry (API key, ANTHROPIC_FOUNDRY_API_KEY)",
		},
		{
			name:             "a piped key, a token in the environment",
			env:              map[string]string{"ANTHROPIC_FOUNDRY_AUTH_TOKEN": fakeFoundryAuthToken},
			stdinFlag:        "--foundry-api-key-stdin",
			piped:            fakeFoundryAPIKey,
			assertCredential: assertFoundryAPIKey,
			wantActive:       "Foundry (API key, --foundry-api-key-stdin)",
			wantInactive:     "Foundry (Entra ID token, ANTHROPIC_FOUNDRY_AUTH_TOKEN)",
		},
		{
			name:             "a piped token, a key in the environment",
			env:              map[string]string{"ANTHROPIC_FOUNDRY_API_KEY": fakeFoundryAPIKey},
			stdinFlag:        "--foundry-auth-token-stdin",
			piped:            fakeFoundryAuthToken,
			assertCredential: assertFoundryAuthToken,
			wantActive:       "Foundry (Entra ID token, --foundry-auth-token-stdin)",
			wantInactive:     "Foundry (API key, ANTHROPIC_FOUNDRY_API_KEY)",
		},
	} {
		setup := func(t *testing.T) []string {
			setupFoundryTest(t)
			for k, v := range tc.env {
				t.Setenv(k, v)
			}
			argv := []string{"--foundry-resource", "my-resource"}
			if tc.stdinFlag != "" {
				withStdinCredential(t, tc.piped+"\n")
				argv = append(argv, tc.stdinFlag)
			}
			return argv
		}

		t.Run(tc.name+": request", func(t *testing.T) {
			argv := setup(t)

			var req *http.Request
			stderr := captureStderr(t, func() {
				runFoundryLeaf(t, argv...)
				req = runFoundryLeaf(t, argv...)
			})

			tc.assertCredential(t, req)
			assert.Equal(t, tc.wantStderr, stderr)
		})

		t.Run(tc.name+": auth status", func(t *testing.T) {
			out, err := runStatus(t, setup(t)...)

			require.NoError(t, err)
			assert.Regexp(t, `\(active\) \* `+regexp.QuoteMeta(tc.wantActive), out)
			assert.Regexp(t, `\n\s+\* `+regexp.QuoteMeta(tc.wantInactive), out)
		})
	}
}

// TestFoundryEndToEnd runs the real binary, so it covers the real flags'
// environment variables, what --debug and `auth status` print, and the errors
// that exit the process. Subprocess for the same reason as TestWorkspaceIDHeader.
func TestFoundryEndToEnd(t *testing.T) {
	setupFoundryTest(t)
	bin := buildAnt(t)
	ant := func(env []string, args ...string) (string, error) {
		cmd := exec.Command(bin, args...)
		cmd.Env = append(cmd.Environ(), env...)
		out, err := cmd.CombinedOutput()
		return string(out), err
	}
	messagesCreate := []string{"messages", "create", "--max-tokens", "1", "--model", "my-deployment", "--message", "{role: user, content: hi}"}
	const (
		resource = "ANTHROPIC_FOUNDRY_RESOURCE=my-resource"
		baseURL  = "ANTHROPIC_FOUNDRY_BASE_URL=https://gateway.invalid/anthropic"
		apiKey   = "ANTHROPIC_FOUNDRY_API_KEY=" + fakeFoundryAPIKey
	)

	for _, tc := range []struct {
		name, env, secret, wantDebug string
		assertCredential             func(*testing.T, *http.Request)
		alsoSet                      []string
	}{
		{"an API key", "ANTHROPIC_FOUNDRY_API_KEY", fakeFoundryAPIKey, "X-Api-Key: <REDACTED>", assertFoundryAPIKey, nil},
		{"an Entra ID token", "ANTHROPIC_FOUNDRY_AUTH_TOKEN", fakeFoundryAuthToken, "Authorization: Bearer <REDACTED>", assertFoundryAuthToken, nil},
		{"an Entra ID token with a key also set", "ANTHROPIC_FOUNDRY_AUTH_TOKEN", fakeFoundryAuthToken, "Authorization: Bearer <REDACTED>", assertFoundryAuthToken, []string{apiKey}},
	} {
		t.Run(tc.name+" reaches ANTHROPIC_FOUNDRY_BASE_URL and --debug redacts it", func(t *testing.T) {
			var got *http.Request
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				got = r
				w.Header().Set("Content-Type", "application/json")
				fmt.Fprint(w, `{}`)
			}))
			defer srv.Close()

			out, err := ant(append([]string{"ANTHROPIC_FOUNDRY_BASE_URL=" + srv.URL + "/anthropic", tc.env + "=" + tc.secret}, tc.alsoSet...),
				append([]string{"--debug"}, messagesCreate...)...)

			require.NoError(t, err, out)
			require.NotNil(t, got, "no request received")
			assert.Equal(t, "/anthropic/v1/messages", got.URL.Path)
			tc.assertCredential(t, got)
			assert.Contains(t, out, tc.wantDebug)
			assert.NotContains(t, out, tc.secret)
		})

		t.Run("auth status shows only the length of "+tc.name, func(t *testing.T) {
			out, err := ant(append([]string{tc.env + "=" + tc.secret}, tc.alsoSet...), "auth", "status")

			require.NoError(t, err, out)
			assert.Contains(t, out, fmt.Sprintf("(set, %d chars)", len(tc.secret)))
			assert.NotContains(t, out, tc.secret[:8])
		})
	}

	t.Run("ANTHROPIC_FOUNDRY_RESOURCE names the resource", func(t *testing.T) {
		out, err := ant([]string{resource, apiKey}, "auth", "status")

		require.NoError(t, err, out)
		assert.Contains(t, out, "https://my-resource.services.ai.azure.com/anthropic")
	})

	t.Run("ANTHROPIC_FOUNDRY_RESOURCE alone leaves the login profile in use", func(t *testing.T) {
		writeAWSTestProfile(t)

		out, err := ant([]string{resource}, "auth", "status")

		require.NoError(t, err, out)
		assert.Regexp(t, `\(active\) \* Profile \(user_oauth\) \[via active_config\]`, out)
		assert.NotContains(t, out, "* Foundry")
	})

	for _, tc := range []struct {
		name             string
		env, args        []string
		wantReason, hint string
	}{
		{"a resource and a base URL", []string{resource, baseURL, apiKey}, nil, "mutually exclusive", "--foundry-base-url / ANTHROPIC_FOUNDRY_BASE_URL"},
		{"a piped key and a piped token", nil, []string{"--foundry-api-key-stdin", "--foundry-auth-token-stdin"}, "mutually exclusive", "--foundry-api-key-stdin and --foundry-auth-token-stdin"},
		{"no credential", nil, []string{"--foundry-resource", "my-resource"}, "no credentials found", "ANTHROPIC_FOUNDRY_API_KEY / --foundry-api-key-stdin"},
		{"no endpoint", []string{apiKey}, nil, "no endpoint found", "--foundry-resource / ANTHROPIC_FOUNDRY_RESOURCE"},
		{"a base URL without a scheme", []string{"ANTHROPIC_FOUNDRY_BASE_URL=gateway.invalid", apiKey}, nil, "missing a scheme", "--foundry-base-url / ANTHROPIC_FOUNDRY_BASE_URL"},
	} {
		t.Run(tc.name+" is an error that names what to set", func(t *testing.T) {
			out, err := ant(tc.env, append(tc.args, messagesCreate...)...)

			require.Error(t, err, out)
			assert.Contains(t, out, tc.wantReason)
			assert.Contains(t, out, tc.hint)
			assert.NotContains(t, out, fakeFoundryAPIKey)
			assert.NotContains(t, out, fakeFoundryAuthToken)
		})
	}
}

// TestFoundrySelection pins what selects the Foundry tier: the resource flag,
// an API key or a token, but not the resource's environment variable or the
// base URL.
func TestFoundrySelection(t *testing.T) {
	for _, tc := range []struct {
		name string
		env  map[string]string
		argv []string
		want bool
	}{
		{name: "--foundry-resource", argv: []string{"--foundry-resource", "my-resource"}, want: true},
		{name: "ANTHROPIC_FOUNDRY_RESOURCE", env: map[string]string{"ANTHROPIC_FOUNDRY_RESOURCE": "my-resource"}},
		{name: "ANTHROPIC_FOUNDRY_API_KEY", env: map[string]string{"ANTHROPIC_FOUNDRY_API_KEY": fakeFoundryAPIKey}, want: true},
		{name: "--foundry-api-key-stdin", argv: []string{"--foundry-api-key-stdin"}, want: true},
		{name: "ANTHROPIC_FOUNDRY_AUTH_TOKEN", env: map[string]string{"ANTHROPIC_FOUNDRY_AUTH_TOKEN": fakeFoundryAuthToken}, want: true},
		{name: "--foundry-auth-token-stdin", argv: []string{"--foundry-auth-token-stdin"}, want: true},
		{name: "--foundry-base-url", argv: []string{"--foundry-base-url", "https://gateway.invalid/anthropic"}},
		{name: "ANTHROPIC_FOUNDRY_BASE_URL", env: map[string]string{"ANTHROPIC_FOUNDRY_BASE_URL": "https://gateway.invalid/anthropic"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			setupFoundryTest(t)
			withStdinCredential(t, "piped-credential\n")
			for k, v := range tc.env {
				t.Setenv(k, v)
			}

			out, err := runStatus(t, tc.argv...)

			require.NoError(t, err)
			assert.Equal(t, tc.want, strings.Contains(out, "(active) * Foundry"), out)
		})
	}
}

// TestFoundryPrecedence pins where the Foundry tier sits: below an API key, an
// auth token, a named profile, federation, Google Cloud and AWS; above the
// implicit profile. Each credential goes only to its own host.
func TestFoundryPrecedence(t *testing.T) {
	foundryArgs := []string{"--foundry-resource", "my-resource"}
	setup := func(t *testing.T) {
		setupFoundryTest(t)
		t.Setenv("ANTHROPIC_FOUNDRY_API_KEY", fakeFoundryAPIKey)
	}

	for _, tc := range []struct {
		name                             string
		other                            func(*testing.T)
		argv                             []string
		wantHost, wantActive, wantWinner string
		wantAuthorization                string
		wantFoundry                      bool
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
			other:      func(t *testing.T) { useFakeAWSCredentials(t, "") },
			argv:       []string{"--aws-region", "us-east-1", "--aws-workspace-id", "wrkspc_01"},
			wantHost:   "aws-external-anthropic.us-east-1.api.aws",
			wantActive: "AWS",
			wantWinner: "AWS (--aws-* / ANTHROPIC_AWS_*)",
		},
		{
			name:        "the implicit profile",
			other:       writeAWSTestProfile,
			wantHost:    "my-resource.services.ai.azure.com",
			wantActive:  "Foundry (API key, ANTHROPIC_FOUNDRY_API_KEY)",
			wantWinner:  "Foundry (--foundry-* / ANTHROPIC_FOUNDRY_*)",
			wantFoundry: true,
		},
	} {
		t.Run("Foundry and "+tc.name, func(t *testing.T) {
			setup(t)
			if tc.other != nil {
				tc.other(t)
			}
			argv := append(foundryArgs, tc.argv...)

			var req *http.Request
			stderr := captureStderr(t, func() { req = runFoundryLeaf(t, argv...) })
			status, err := runStatusKeepingCredentialEnv(t, argv...)
			require.NoError(t, err)
			var inUse bool
			root := foundryTestRoot(func(ctx context.Context, c *cli.Command) error {
				inUse = foundryInUse(c)
				return nil
			})
			require.NoError(t, root.Run(context.Background(), append(append([]string{"ant"}, argv...), "leaf")))

			assert.Equal(t, tc.wantHost, req.URL.Host)
			if tc.wantFoundry {
				assertFoundryAPIKey(t, req)
			} else {
				assert.Empty(t, req.Header.Values("x-api-key"))
			}
			if tc.wantAuthorization != "" {
				assert.Equal(t, []string{tc.wantAuthorization}, req.Header.Values("Authorization"))
			}
			assert.Contains(t, stderr, "Foundry (--foundry-* / ANTHROPIC_FOUNDRY_*)")
			assert.Contains(t, stderr, "using "+tc.wantWinner+" per precedence")
			assert.Contains(t, status, fmt.Sprintf("(active) * %-46s ", tc.wantActive))
			assert.Equal(t, tc.wantFoundry, inUse)
		})
	}

	t.Run("a resource and a key are one source", func(t *testing.T) {
		setup(t)

		stderr := captureStderr(t, func() { runFoundryLeaf(t, foundryArgs...) })

		assert.Empty(t, stderr)
	})

	t.Run("ANTHROPIC_FOUNDRY_RESOURCE alone leaves the implicit profile in use", func(t *testing.T) {
		setupFoundryTest(t)
		writeAWSTestProfile(t)
		t.Setenv("ANTHROPIC_FOUNDRY_RESOURCE", "my-resource")

		var req *http.Request
		stderr := captureStderr(t, func() { req = runFoundryLeaf(t) })

		assert.Equal(t, "profile.invalid", req.URL.Host)
		assert.Equal(t, []string{"Bearer profile-token"}, req.Header.Values("Authorization"))
		assert.Empty(t, stderr)
	})
}

// TestAuthStatusFoundry pins that `auth status` reports the tier the client
// would use, where requests go, and no part of a key or token.
func TestAuthStatusFoundry(t *testing.T) {
	t.Run("an API key", func(t *testing.T) {
		setupFoundryTest(t)
		t.Setenv("ANTHROPIC_FOUNDRY_API_KEY", "fake-foundry-key-0123456789abcdef")

		out, err := runStatus(t, "--foundry-resource", "my-resource")

		require.NoError(t, err)
		assert.Regexp(t, `\(active\) \* Foundry \(API key, ANTHROPIC_FOUNDRY_API_KEY\)\s+\(set, 33 chars\)`, out)
		assert.Regexp(t, `resource:\s+my-resource`, out)
		assert.Regexp(t, `\(active\) \* Foundry resource\s+https://my-resource\.services\.ai\.azure\.com/anthropic`, out)
		assert.NotContains(t, out, "fake-fou")
		assert.NotContains(t, out, "0123456789abcdef")
	})

	t.Run("an Entra ID token", func(t *testing.T) {
		setupFoundryTest(t)
		t.Setenv("ANTHROPIC_FOUNDRY_AUTH_TOKEN", "fake-entra-token-0123456789abcdef")

		out, err := runStatus(t, "--foundry-resource", "my-resource")

		require.NoError(t, err)
		assert.Regexp(t, `\(active\) \* Foundry \(Entra ID token, ANTHROPIC_FOUNDRY_AUTH_TOKEN\)\s+\(set, 33 chars\)`, out)
		assert.NotContains(t, out, "fake-ent")
		assert.NotContains(t, out, "0123456789abcdef")
	})

	t.Run("ANTHROPIC_FOUNDRY_RESOURCE names the resource", func(t *testing.T) {
		setupFoundryTest(t)
		t.Setenv("ANTHROPIC_FOUNDRY_RESOURCE", "my-resource")
		t.Setenv("ANTHROPIC_FOUNDRY_API_KEY", fakeFoundryAPIKey)

		out, err := runStatus(t)

		require.NoError(t, err)
		assert.Regexp(t, `resource:\s+my-resource`, out)
		assert.Regexp(t, `\(active\) \* Foundry resource\s+https://my-resource\.services\.ai\.azure\.com/anthropic`, out)
	})

	t.Run("ANTHROPIC_FOUNDRY_RESOURCE alone adds no row", func(t *testing.T) {
		setupFoundryTest(t)
		writeAWSTestProfile(t)
		t.Setenv("ANTHROPIC_FOUNDRY_RESOURCE", "my-resource")

		out, err := runStatus(t)

		require.NoError(t, err)
		assert.Regexp(t, `\(active\) \* Profile \(user_oauth\) \[via active_config\]`, out)
		assert.NotContains(t, out, "* Foundry")
	})

	t.Run("a piped key names its flag", func(t *testing.T) {
		setupFoundryTest(t)
		withStdinCredential(t, fakeFoundryAPIKey+"\n")

		out, err := runStatus(t, "--foundry-api-key-stdin", "--foundry-resource", "my-resource")

		require.NoError(t, err)
		assert.Regexp(t, `\(active\) \* Foundry \(API key, --foundry-api-key-stdin\)`, out)
	})

	t.Run("no credential", func(t *testing.T) {
		setupFoundryTest(t)

		out, err := runStatus(t, "--foundry-resource", "my-resource")

		require.NoError(t, err)
		assert.Regexp(t, `\(active\) \* Foundry\s+\(no API key or Entra ID token set\)`, out)
	})

	t.Run("no resource", func(t *testing.T) {
		setupFoundryTest(t)
		t.Setenv("ANTHROPIC_FOUNDRY_API_KEY", fakeFoundryAPIKey)

		out, err := runStatus(t)

		require.NoError(t, err)
		assert.Regexp(t, `resource:\s+\(not set\)`, out)
		assert.Regexp(t, `\(active\) \* Foundry resource\s+\(set --foundry-resource\)`, out)
	})

	t.Run("a generic base URL is listed but not active", func(t *testing.T) {
		setupFoundryTest(t)

		out, err := runStatus(t, "--base-url", "https://anthropic-proxy.invalid", "--foundry-resource", "my-resource")

		require.NoError(t, err)
		assert.Regexp(t, `\(active\) \* Foundry resource\s+https://my-resource\.services\.ai\.azure\.com/anthropic`, out)
		assert.Regexp(t, `\n\s+\* --base-url flag\s+https://anthropic-proxy.invalid`, out)
	})

	t.Run("the Foundry base URL is the active one when set", func(t *testing.T) {
		setupFoundryTest(t)
		t.Setenv("ANTHROPIC_FOUNDRY_API_KEY", fakeFoundryAPIKey)

		out, err := runStatus(t, "--foundry-base-url", "https://gateway.invalid/anthropic")

		require.NoError(t, err)
		assert.Regexp(t, `\(active\) \* --foundry-base-url / ANTHROPIC_FOUNDRY_BASE_URL\s+https://gateway.invalid/anthropic`, out)
	})

	t.Run("an API key wins over Foundry", func(t *testing.T) {
		setupFoundryTest(t)
		t.Setenv("ANTHROPIC_API_KEY", "fake-first-party-key")

		out, err := runStatusKeepingCredentialEnv(t, "--foundry-resource", "my-resource")

		require.NoError(t, err)
		assert.Regexp(t, `\(active\) \* ANTHROPIC_API_KEY`, out)
		assert.Regexp(t, `\n\s+\* Foundry\s`, out)
	})

	t.Run("Foundry wins over the implicit profile", func(t *testing.T) {
		setupFoundryTest(t)
		writeAWSTestProfile(t)
		t.Setenv("ANTHROPIC_FOUNDRY_API_KEY", fakeFoundryAPIKey)

		out, err := runStatus(t, "--foundry-resource", "my-resource")

		require.NoError(t, err)
		assert.Regexp(t, `\(active\) \* Foundry \(API key`, out)
		assert.Regexp(t, `\n\s+\* Profile \(user_oauth\) \[via active_config\]`, out)
	})
}

// TestFoundryInUse pins what `ant apply` records as the host, how it names the
// credentials, and whether it treats the active profile as the one in use.
func TestFoundryInUse(t *testing.T) {
	var credentials string
	probe := func(t *testing.T, argv ...string) (inUse bool, baseURL string, profile *config.Config) {
		t.Helper()
		root := foundryTestRoot(func(ctx context.Context, c *cli.Command) error {
			require.NoError(t, applyStdinCredential(c.Root()))
			inUse, baseURL, profile = foundryInUse(c), currentBaseURL(c), profileInUse(c)
			credentials = describeOrigin(c, core.Origin{BaseURL: baseURL}).Credentials
			return nil
		})
		require.NoError(t, root.Run(context.Background(), append(append([]string{"ant"}, argv...), "leaf")))
		return inUse, baseURL, profile
	}

	t.Run("resource, API key", func(t *testing.T) {
		setupFoundryTest(t)
		t.Setenv("ANTHROPIC_FOUNDRY_API_KEY", fakeFoundryAPIKey)

		inUse, baseURL, _ := probe(t, "--foundry-resource", "my-resource")

		assert.True(t, inUse)
		assert.Equal(t, "https://my-resource.services.ai.azure.com/anthropic", baseURL)
		assert.Equal(t, "Foundry (API key)", credentials)
	})

	t.Run("--foundry-base-url, Entra ID token", func(t *testing.T) {
		setupFoundryTest(t)
		t.Setenv("ANTHROPIC_FOUNDRY_AUTH_TOKEN", fakeFoundryAuthToken)
		t.Setenv("ANTHROPIC_BASE_URL", "https://anthropic-proxy.invalid")

		_, baseURL, _ := probe(t, "--foundry-base-url", "https://gateway.invalid/anthropic")

		assert.Equal(t, "https://gateway.invalid/anthropic", baseURL)
		assert.Equal(t, "Foundry (Entra ID token)", credentials)
	})

	t.Run("a key and a token", func(t *testing.T) {
		setupFoundryTest(t)
		t.Setenv("ANTHROPIC_FOUNDRY_API_KEY", fakeFoundryAPIKey)
		t.Setenv("ANTHROPIC_FOUNDRY_AUTH_TOKEN", fakeFoundryAuthToken)

		probe(t, "--foundry-resource", "my-resource")

		assert.Equal(t, "Foundry (Entra ID token)", credentials)
	})

	t.Run("a piped key and a token", func(t *testing.T) {
		setupFoundryTest(t)
		withStdinCredential(t, fakeFoundryAPIKey+"\n")
		t.Setenv("ANTHROPIC_FOUNDRY_AUTH_TOKEN", fakeFoundryAuthToken)

		probe(t, "--foundry-resource", "my-resource", "--foundry-api-key-stdin")

		assert.Equal(t, "Foundry (API key)", credentials)
	})

	t.Run("not in use with ANTHROPIC_FOUNDRY_RESOURCE alone", func(t *testing.T) {
		setupFoundryTest(t)
		writeAWSTestProfile(t)
		t.Setenv("ANTHROPIC_FOUNDRY_RESOURCE", "my-resource")

		inUse, baseURL, profile := probe(t)

		assert.False(t, inUse)
		assert.Equal(t, "https://profile.invalid", baseURL)
		assert.NotNil(t, profile)
	})

	t.Run("not in use when an API key wins", func(t *testing.T) {
		setupFoundryTest(t)
		t.Setenv("ANTHROPIC_API_KEY", "fake-first-party-key")

		inUse, baseURL, _ := probe(t, "--foundry-resource", "my-resource")

		assert.False(t, inUse)
		assert.Equal(t, defaultBaseURL, baseURL)
	})

	t.Run("not in use when AWS wins", func(t *testing.T) {
		setupFoundryTest(t)

		inUse, baseURL, _ := probe(t, "--foundry-resource", "my-resource", "--aws-region", "us-east-1")

		assert.False(t, inUse)
		assert.Equal(t, "https://aws-external-anthropic.us-east-1.api.aws", baseURL)
	})

	t.Run("the implicit profile is not in use", func(t *testing.T) {
		setupFoundryTest(t)
		writeAWSTestProfile(t)

		inUse, _, profile := probe(t, "--foundry-resource", "my-resource")

		assert.True(t, inUse)
		assert.Nil(t, profile)
	})
}
