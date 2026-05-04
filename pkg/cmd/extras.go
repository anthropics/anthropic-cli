package cmd

// Hand-written additions to the codegen-owned Command tree. Anything that
// must survive a regen lives here rather than in the generated files.

import (
	"github.com/urfave/cli/v3"
)

func init() {
	Command.Flags = append(Command.Flags,
		&cli.StringFlag{
			Name:    "profile",
			Usage:   "Named auth profile to use (default: active profile from active_config)",
			Sources: cli.EnvVars("ANTHROPIC_PROFILE"),
		},
		&cli.StringFlag{
			Name:    "identity-token",
			Usage:   "Signed OIDC JWT for federation (jwt-bearer grant). Mutually exclusive with --identity-token-file.",
			Sources: cli.EnvVars("ANTHROPIC_IDENTITY_TOKEN"),
		},
		&cli.StringFlag{
			Name:    "identity-token-file",
			Usage:   "Path to a file containing a signed OIDC JWT for federation. Re-read on every request, supporting rotating tokens (K8s projected SA, GitHub Actions OIDC).",
			Sources: cli.EnvVars("ANTHROPIC_IDENTITY_TOKEN_FILE"),
		},
		&cli.StringFlag{
			Name:    "federation-rule",
			Usage:   "Tagged ID of the OIDC federation rule (fdrl_...).",
			Sources: cli.EnvVars("ANTHROPIC_FEDERATION_RULE_ID"),
		},
		&cli.StringFlag{
			Name:    "organization-id",
			Usage:   "Anthropic organization UUID for federation token minting.",
			Sources: cli.EnvVars("ANTHROPIC_ORGANIZATION_ID"),
		},
		&cli.StringFlag{
			Name:    "service-account-id",
			Usage:   "Optional service-account tagged ID (svac_...) for target_type=SERVICE_ACCOUNT federation rules.",
			Sources: cli.EnvVars("ANTHROPIC_SERVICE_ACCOUNT_ID"),
		},
	)
}
