package cmd

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli/v3"

	"github.com/anthropics/anthropic-sdk-go/aws"
	"github.com/anthropics/anthropic-sdk-go/option"
)

// awsFlags returns the AWS tier flags (mirroring extras.go) plus the
// first-party flags the conflict check reads, so tests exercise the same
// flag layer production uses.
func awsFlags() []cli.Flag {
	return []cli.Flag{
		&cli.BoolFlag{Name: "aws", Sources: cli.EnvVars("ANTHROPIC_USE_AWS")},
		&cli.StringFlag{Name: "aws-region", Sources: cli.EnvVars("AWS_REGION", "AWS_DEFAULT_REGION")},
		&cli.StringFlag{Name: "aws-workspace-id", Sources: cli.EnvVars("ANTHROPIC_AWS_WORKSPACE_ID")},
		&cli.StringFlag{Name: "aws-api-key", Sources: cli.EnvVars("ANTHROPIC_AWS_API_KEY")},
		&cli.StringFlag{Name: "base-url"},
		// First-party creds the AWS conflict notice inspects.
		&cli.StringFlag{Name: "api-key"},
		&cli.StringFlag{Name: "auth-token"},
		&cli.StringFlag{Name: "profile", Sources: cli.EnvVars("ANTHROPIC_PROFILE")},
		&cli.StringFlag{Name: "identity-token"},
		&cli.StringFlag{Name: "identity-token-file"},
		&cli.StringFlag{Name: "federation-rule"},
		&cli.StringFlag{Name: "organization-id"},
		&cli.StringFlag{Name: "service-account-id"},
	}
}

// clearAWSEnv unsets every env var the AWS flag Sources read, so ambient
// values on the test host (a developer laptop or EC2 instance commonly has
// AWS_REGION set) cannot leak into flag resolution.
func clearAWSEnv(t *testing.T) {
	t.Helper()
	for _, k := range []string{
		"ANTHROPIC_USE_AWS", "AWS_REGION", "AWS_DEFAULT_REGION",
		"ANTHROPIC_AWS_WORKSPACE_ID", "ANTHROPIC_AWS_API_KEY", "ANTHROPIC_AWS_BASE_URL",
		"ANTHROPIC_PROFILE", "ANTHROPIC_API_KEY", "ANTHROPIC_AUTH_TOKEN",
	} {
		clearEnv(t, k)
	}
}

// newAWSCmd builds a *cli.Command carrying awsFlags() with the given flag
// values pre-set, so helpers that read cmd.String/cmd.IsSet behave exactly as
// they do after real flag parsing.
func newAWSCmd(t *testing.T, set map[string]string) *cli.Command {
	t.Helper()
	cmd := &cli.Command{Flags: awsFlags()}
	for k, v := range set {
		require.NoError(t, cmd.Set(k, v))
	}
	return cmd
}

// TestBuildAWSConfig asserts buildAWSConfig copies the four flag values into
// the matching aws.ClientConfig fields — including the empty-flag case (empty
// in ⇒ empty out; the SDK does the env/regional fallback, not us).
func TestBuildAWSConfig(t *testing.T) {
	for _, tc := range []struct {
		name string
		set  map[string]string
		want aws.ClientConfig
	}{
		{
			name: "all set",
			set: map[string]string{
				"aws-workspace-id": "wrkspc_123",
				"aws-region":       "us-west-2",
				"aws-api-key":      "fake-aws-gateway-key",
				"base-url":         "https://staging.example.com",
			},
			want: aws.ClientConfig{
				WorkspaceID: "wrkspc_123",
				AWSRegion:   "us-west-2",
				APIKey:      "fake-aws-gateway-key",
				BaseURL:     "https://staging.example.com",
			},
		},
		{
			name: "sigv4 mode (no api key)",
			set: map[string]string{
				"aws-workspace-id": "wrkspc_456",
				"aws-region":       "eu-central-1",
			},
			want: aws.ClientConfig{
				WorkspaceID: "wrkspc_456",
				AWSRegion:   "eu-central-1",
			},
		},
		{
			name: "all empty — SDK falls back",
			set:  map[string]string{},
			want: aws.ClientConfig{},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			clearAWSEnv(t)
			cmd := newAWSCmd(t, tc.set)
			got := buildAWSConfig(cmd)
			assert.Equal(t, tc.want, got)
		})
	}
}

// TestAWSConflictWarning mirrors TestMultiAuthWarning: the AWS conflict notice
// fires when --aws is paired with a first-party credential (which --aws
// overrides), and is silent when only --aws-api-key (mode selector, not a
// cross-tier conflict) is set alongside --aws.
func TestAWSConflictWarning(t *testing.T) {
	t.Run("warns on first-party api-key", func(t *testing.T) {
		clearAWSEnv(t)
		resetWarnOnce(t)
		cmd := newAWSCmd(t, map[string]string{"aws": "true", "api-key": "fake-first-party-key"})
		out := captureStderr(t, func() { warnIfAWSConflict(cmd) })
		assert.Contains(t, out, "--aws is active")
		assert.Contains(t, out, "--api-key / ANTHROPIC_API_KEY")
	})

	t.Run("warns on auth-token", func(t *testing.T) {
		clearAWSEnv(t)
		resetWarnOnce(t)
		cmd := newAWSCmd(t, map[string]string{"aws": "true", "auth-token": "tok"})
		out := captureStderr(t, func() { warnIfAWSConflict(cmd) })
		assert.Contains(t, out, "--auth-token / ANTHROPIC_AUTH_TOKEN")
	})

	t.Run("warns on explicit profile", func(t *testing.T) {
		clearAWSEnv(t)
		resetWarnOnce(t)
		cmd := newAWSCmd(t, map[string]string{"aws": "true", "profile": "work"})
		out := captureStderr(t, func() { warnIfAWSConflict(cmd) })
		assert.Contains(t, out, "profile from --profile / ANTHROPIC_PROFILE")
	})

	t.Run("warns on federation env", func(t *testing.T) {
		clearAWSEnv(t)
		resetWarnOnce(t)
		cmd := newAWSCmd(t, map[string]string{"aws": "true", "federation-rule": "fdrl_x"})
		out := captureStderr(t, func() { warnIfAWSConflict(cmd) })
		assert.Contains(t, out, "federation env")
	})

	t.Run("silent for --aws-api-key alone (mode selector, must not warn)", func(t *testing.T) {
		clearAWSEnv(t)
		resetWarnOnce(t)
		cmd := newAWSCmd(t, map[string]string{"aws": "true", "aws-api-key": "fake-aws-gateway-key"})
		out := captureStderr(t, func() { warnIfAWSConflict(cmd) })
		assert.Empty(t, out, "ANTHROPIC_AWS_API_KEY selects API-key mode; it is not a cross-tier conflict")
	})

	t.Run("silent when no first-party cred set", func(t *testing.T) {
		clearAWSEnv(t)
		resetWarnOnce(t)
		cmd := newAWSCmd(t, map[string]string{"aws": "true", "aws-region": "us-west-2", "aws-workspace-id": "wrkspc_1"})
		out := captureStderr(t, func() { warnIfAWSConflict(cmd) })
		assert.Empty(t, out)
	})

	t.Run("emits once", func(t *testing.T) {
		clearAWSEnv(t)
		resetWarnOnce(t)
		cmd := newAWSCmd(t, map[string]string{"aws": "true", "api-key": "fake-first-party-key"})
		first := captureStderr(t, func() { warnIfAWSConflict(cmd) })
		second := captureStderr(t, func() { warnIfAWSConflict(cmd) })
		assert.NotEmpty(t, first)
		assert.Empty(t, second)
	})
}

// runAWSStatus runs `auth status` against a root carrying the AWS tier flags
// (plus the first-party flags authStatus reads via c.Root()), passing the given
// global flag values as argv before the subcommand — mirroring how the binary
// parses them. (Values must come through parsing, not pre-Set: Run re-parses
// argv and resets flags.) Mirrors runStatus but for the AWS path.
func runAWSStatus(t *testing.T, set map[string]string) (string, error) {
	t.Helper()
	root := &cli.Command{
		Name:  "ant",
		Flags: awsFlags(),
		Commands: []*cli.Command{{
			Name: "auth", Commands: []*cli.Command{{
				Name: "status", Action: authStatus,
			}},
		}},
	}
	argv := []string{"ant"}
	for k, v := range set {
		if k == "aws" {
			argv = append(argv, "--aws")
			continue
		}
		argv = append(argv, "--"+k, v)
	}
	argv = append(argv, "auth", "status")
	return captureStdout(t, func() error {
		return root.Run(t.Context(), argv)
	})
}

// TestAWSPrecedenceInvariant guards the cross-site invariant the plan requires
// (ANT_AWS_PLAN.md "Three-site precedence invariant"): AWS is tier-0 and must
// win across all three sites that mirror each other — (a) the request-options
// short-circuit in getDefaultRequestOptions, (b) warnIfMultipleAuthSources'
// ordering, and (c) credWinner in authStatus. The test pairs --aws with the
// TOP first-party tier (--api-key, tier 1) so AWS is genuinely beating the
// strongest competing source, and asserts AWS wins at the two observable sites
// in lockstep. API-key mode keeps aws.NewClient network-free; resetWarnOnce
// isolates the shared one-shot guard.
func TestAWSPrecedenceInvariant(t *testing.T) {
	awsAndTopFirstParty := map[string]string{
		"aws":              "true",
		"aws-region":       "us-west-2",
		"aws-workspace-id": "wrkspc_abc",
		"aws-api-key":      "fake-aws-gateway-key",
		"api-key":          "fake-first-party-secret", // tier-1 first-party; must be overridden
	}

	t.Run("site (a): request options short-circuit before first-party switch", func(t *testing.T) {
		clearAWSEnv(t)
		resetWarnOnce(t)
		cmd := newAWSCmd(t, awsAndTopFirstParty)
		var opts []option.RequestOption
		out := captureStderr(t, func() {
			opts = getDefaultRequestOptions(cmd)
		})
		// AWS branch ran (its conflict notice), NOT the first-party switch
		// (which would emit the "multiple auth sources configured" notice).
		assert.Contains(t, out, "--aws is active")
		assert.NotContains(t, out, "multiple auth sources configured")
		// Prove the AWS backend opts actually flowed through — not merely that
		// the slice is non-empty (the base opts alone make it non-empty, so
		// NotEmpty could never fail). The AWS branch returns base opts +
		// awsClient.Options. Derive the base count empirically from a no-aws,
		// no-credential invocation (self-calibrating — survives base-opt
		// changes, no magic number), then build the AWS client independently
		// (network-free in API-key mode) and assert the totals reconcile.
		baseCmd := newAWSCmd(t, nil) // no flags set ⇒ switch appends nothing
		baseCount := len(getDefaultRequestOptions(baseCmd))
		awsClient, err := aws.NewClient(context.Background(), buildAWSConfig(cmd))
		require.NoError(t, err)
		require.NotEmpty(t, awsClient.Options, "sanity: AWS client should contribute options in API-key mode")
		assert.Equal(t, baseCount+len(awsClient.Options), len(opts),
			"getDefaultRequestOptions must return base opts + the AWS backend opts")
	})

	t.Run("site (c): auth status reports AWS winner, not the first-party api-key", func(t *testing.T) {
		clearAWSEnv(t)
		out, err := runAWSStatus(t, awsAndTopFirstParty)
		require.NoError(t, err)
		// AWS is the winner; the tier-1 first-party api-key is neither shown as
		// the credential winner nor leaked.
		assert.Contains(t, out, "Claude Platform on AWS")
		assert.Contains(t, out, "AWS gateway (API key, x-api-key)")
		assert.NotContains(t, out, "fake-first-party-secret")
		assert.NotContains(t, out, "--api-key / ANTHROPIC_API_KEY")
	})
}

// TestAuthStatusAWSRow asserts that when Bool("aws") is set, auth status shows
// AWS as the active backend with the correct mode label, and does NOT render
// the first-party profile/key/federation winner rows.
func TestAuthStatusAWSRow(t *testing.T) {
	t.Run("API-key mode", func(t *testing.T) {
		clearAWSEnv(t)
		out, err := runAWSStatus(t, map[string]string{
			"aws":              "true",
			"aws-region":       "us-west-2",
			"aws-workspace-id": "wrkspc_abc",
			"aws-api-key":      "fake-aws-gateway-key-1234567890",
		})
		require.NoError(t, err)
		assert.Contains(t, out, "Claude Platform on AWS")
		assert.Contains(t, out, "AWS gateway (API key, x-api-key)")
		assert.Contains(t, out, "us-west-2")
		assert.Contains(t, out, "wrkspc_abc")
		// Regional base URL derived from region.
		assert.Contains(t, out, "https://aws-external-anthropic.us-west-2.api.aws")
		// Secret is redacted, not printed in full.
		assert.NotContains(t, out, "fake-aws-gateway-key-1234567890")
		// First-party winner rows must not appear.
		assert.NotContains(t, out, "Active profile:")
		assert.NotContains(t, out, "no credential configured")
	})

	t.Run("SigV4 mode", func(t *testing.T) {
		clearAWSEnv(t)
		out, err := runAWSStatus(t, map[string]string{
			"aws":              "true",
			"aws-region":       "eu-central-1",
			"aws-workspace-id": "wrkspc_xyz",
		})
		require.NoError(t, err)
		assert.Contains(t, out, "AWS gateway (SigV4)")
		assert.Contains(t, out, "AWS credential chain")
		assert.Contains(t, out, "eu-central-1")
		assert.NotContains(t, out, "AWS gateway (API key")
	})

	t.Run("base-url override wins over regional derivation", func(t *testing.T) {
		clearAWSEnv(t)
		out, err := runAWSStatus(t, map[string]string{
			"aws":              "true",
			"aws-region":       "us-west-2",
			"aws-workspace-id": "wrkspc_abc",
			"base-url":         "https://staging.example.com",
		})
		require.NoError(t, err)
		assert.Contains(t, out, "https://staging.example.com")
		assert.NotContains(t, out, "https://aws-external-anthropic.us-west-2.api.aws")
	})
}
