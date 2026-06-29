package cmd

// Hand-written additions to the codegen-owned Command tree. Anything that
// must survive a regen lives here rather than in the generated files.

import (
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/urfave/cli/v3"
)

// extraClientFlags are the global flags that map onto SDK request options,
// shared by every client construction site (getDefaultRequestOptions,
// newWorkerClient).
type extraClientFlags struct {
	BaseURL string
	// Sent as the anthropic-workspace-id header. Required on Claude Platform
	// on AWS; ignored by the first-party API.
	WorkspaceID string
}

func extraClientFlagsFromCmd(cmd *cli.Command) extraClientFlags {
	return extraClientFlags{
		BaseURL:     cmd.String("base-url"),
		WorkspaceID: cmd.String("workspace-id"),
	}
}

func (o extraClientFlags) requestOptions() []option.RequestOption {
	var opts []option.RequestOption
	if o.WorkspaceID != "" {
		opts = append(opts, option.WithHeader("anthropic-workspace-id", o.WorkspaceID))
	}
	if o.BaseURL != "" {
		opts = append(opts, option.WithBaseURL(o.BaseURL))
	}
	return opts
}

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
		&cli.StringFlag{
			Name:    "workspace-id",
			Usage:   "Workspace (wrkspc_...) to send as the anthropic-workspace-id header; only needed on Claude Platform on AWS.",
			Sources: cli.EnvVars("ANTHROPIC_WORKSPACE_ID"),
		},
	)
}
