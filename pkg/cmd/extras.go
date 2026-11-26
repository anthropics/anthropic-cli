package cmd

// Hand-written additions to the codegen-owned Command tree. Anything that
// must survive a regen lives here rather than in the generated files.

import (
	"os"
	"strings"

	"github.com/anthropics/anthropic-sdk-go/v2/option"
	"github.com/urfave/cli/v3"
)

// extraClientFlags are the global flags that map onto SDK request options,
// shared by every client construction site (getDefaultRequestOptions,
// newWorkerClient).
type extraClientFlags struct {
	BaseURL string
	// Sent as the anthropic-workspace-id header. Picks the workspace for API keys
	// that have access to multiple workspaces.
	WorkspaceID string
}

// extraClientFlagsFromCmd reads the client-level flags. urfave/cli resolves a
// flag name leaf-first, so a subcommand that redeclares a name shadows the
// root flag of the same name.
func extraClientFlagsFromCmd(cmd *cli.Command) extraClientFlags {
	return extraClientFlags{
		// A subcommand's --base-url (e.g. beta:worker) means the same thing as the
		// root's and is meant to override it, so the leaf-first lookup is right.
		BaseURL: cmd.String("base-url"),
		// A subcommand's --workspace-id is a different parameter: on the admin
		// commands it is the workspace being operated on (`/workspaces/{workspace_id}`
		// or a body field), not the `anthropic-workspace-id` header that scopes the
		// credential. Only the root flag / ANTHROPIC_WORKSPACE_ID feeds the header.
		WorkspaceID: cmd.Root().String("workspace-id"),
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

// customHeadersFromEnv reads ANTHROPIC_CUSTOM_HEADERS, which
// WithoutEnvironmentDefaults keeps the SDK from reading.
func customHeadersFromEnv() []option.RequestOption {
	var opts []option.RequestOption
	for _, line := range strings.Split(os.Getenv("ANTHROPIC_CUSTOM_HEADERS"), "\n") {
		if name, value, ok := strings.Cut(line, ":"); ok {
			opts = append(opts, option.WithHeader(strings.TrimSpace(name), strings.TrimSpace(value)))
		}
	}
	return opts
}

// init runs after cmd.go's, which builds Command — package-level init
// functions run in filename order, and "extras.go" sorts after "cmd.go".
// Hand-written commands are registered here for that reason.
func init() {
	Command.Commands = append(Command.Commands, &applyCommand)

	// cmd.go leaves both hooks unset, so nothing generated is replaced here.
	Command.Before = otelBefore
	Command.After = otelAfter

	// cmd.go declares --base-url without an env source and the SDK's own
	// env handling is disabled (WithoutEnvironmentDefaults), so attach
	// ANTHROPIC_BASE_URL here. TestBaseURLEnvRealTree fails if a regen
	// reshapes the flag and this patch stops applying.
	for _, f := range Command.Flags {
		if sf, ok := f.(*cli.StringFlag); ok && sf.Name == "base-url" {
			sf.Sources = cli.EnvVars("ANTHROPIC_BASE_URL")
		}
	}
	Command.Flags = append(Command.Flags,
		&cli.BoolFlag{
			Name:  "api-key-stdin",
			Usage: "Read the API key from standard input (e.g. op read op://vault/anthropic/key | ant --api-key-stdin models list). Consumes stdin, so the request body must come from flags.",
		},
		&cli.BoolFlag{
			Name:  "auth-token-stdin",
			Usage: "Read the bearer auth token from standard input. Consumes stdin, so the request body must come from flags.",
		},
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
			Usage:   "Workspace ID (wrkspc_...) to send as the anthropic-workspace-id header. Only needed for API keys that aren't scoped to a single workspace. On Claude Platform on AWS, use --aws-workspace-id.",
			Sources: cli.EnvVars("ANTHROPIC_WORKSPACE_ID"),
		},
		&cli.StringFlag{
			Name:    "google-cloud-project",
			Usage:   "Google Cloud project ID. Setting this, --google-cloud-location or --google-cloud-workspace-id sends requests to Claude Models on Google Cloud (Anthropic managed), signed in with Google Application Default Credentials. Defaults to the credentials' project.",
			Sources: cli.EnvVars("ANTHROPIC_GOOGLE_CLOUD_PROJECT"),
		},
		&cli.StringFlag{
			Name:    "google-cloud-location",
			Usage:   "Google Cloud location for Claude Models on Google Cloud (Anthropic managed). Defaults to global.",
			Sources: cli.EnvVars("ANTHROPIC_GOOGLE_CLOUD_LOCATION"),
		},
		&cli.StringFlag{
			Name:    "google-cloud-workspace-id",
			Usage:   "Workspace ID (wrkspc_...) to address on Claude Models on Google Cloud (Anthropic managed). It is part of the gateway URL there, so --workspace-id is not needed.",
			Sources: cli.EnvVars("ANTHROPIC_GOOGLE_CLOUD_WORKSPACE_ID"),
		},
		&cli.StringFlag{
			Name:    "google-cloud-base-url",
			Usage:   "Override the gateway URL for Claude Models on Google Cloud (Anthropic managed), once another --google-cloud-* flag has selected it. Google credentials are sent only here or to the gateway, never to --base-url.",
			Sources: cli.EnvVars("ANTHROPIC_GOOGLE_CLOUD_BASE_URL"),
		},
		&cli.StringFlag{
			Name:  "aws-region",
			Usage: "AWS region for Claude Platform on AWS. Setting this, --aws-workspace-id or an AWS API key sends requests to Claude Platform on AWS. Without an API key they are signed with AWS credentials from the default credential chain. Defaults to AWS_REGION, then AWS_DEFAULT_REGION.",
		},
		&cli.StringFlag{
			Name:    "aws-workspace-id",
			Usage:   "Workspace ID (wrkspc_...) to address on Claude Platform on AWS. It is sent as the anthropic-workspace-id header there, so --workspace-id is not needed.",
			Sources: cli.EnvVars("ANTHROPIC_AWS_WORKSPACE_ID"),
		},
		&cli.BoolFlag{
			Name:  "aws-api-key-stdin",
			Usage: "Read an API key for Claude Platform on AWS from standard input, instead of signing with AWS credentials (or set ANTHROPIC_AWS_API_KEY). Consumes stdin, so the request body must come from flags.",
		},
		&cli.StringFlag{
			Name:    "aws-base-url",
			Usage:   "Override the gateway URL for Claude Platform on AWS, once --aws-region, --aws-workspace-id or an AWS API key has selected it. AWS credentials are sent only here or to the gateway, never to --base-url.",
			Sources: cli.EnvVars("ANTHROPIC_AWS_BASE_URL"),
		},
		&cli.StringFlag{
			Name:  "foundry-resource",
			Usage: "Microsoft Foundry resource name, as in https://{resource}.services.ai.azure.com. Setting this, a Foundry API key or a Microsoft Entra ID token sends requests to Claude in Microsoft Foundry. Defaults to ANTHROPIC_FOUNDRY_RESOURCE, which doesn't do that on its own.",
		},
		&cli.BoolFlag{
			Name:  "foundry-api-key-stdin",
			Usage: "Read an API key for Claude in Microsoft Foundry from standard input (or set ANTHROPIC_FOUNDRY_API_KEY). Consumes stdin, so the request body must come from flags.",
		},
		&cli.BoolFlag{
			Name:  "foundry-auth-token-stdin",
			Usage: "Read a Microsoft Entra ID access token for Claude in Microsoft Foundry from standard input, instead of an API key (or set ANTHROPIC_FOUNDRY_AUTH_TOKEN). Consumes stdin, so the request body must come from flags.",
		},
		&cli.StringFlag{
			Name:    "foundry-base-url",
			Usage:   "Full base URL for Claude in Microsoft Foundry, including the /anthropic path, instead of --foundry-resource. Foundry credentials are sent only here or to the resource, never to --base-url.",
			Sources: cli.EnvVars("ANTHROPIC_FOUNDRY_BASE_URL"),
		},
		&cli.StringFlag{
			Name:    "bedrock-region",
			Usage:   "AWS region for Claude in Amazon Bedrock. Setting this or --bedrock-api-key-stdin sends requests to the Bedrock Messages API endpoint, where model IDs carry the anthropic. prefix. Without an API key they are signed with AWS credentials from the default credential chain. Defaults to AWS_REGION, then AWS_DEFAULT_REGION.",
			Sources: cli.EnvVars("ANTHROPIC_BEDROCK_REGION"),
		},
		&cli.BoolFlag{
			Name:  "bedrock-api-key-stdin",
			Usage: "Read a Bedrock API key for Claude in Amazon Bedrock from standard input, instead of signing with AWS credentials (or set AWS_BEARER_TOKEN_BEDROCK once --bedrock-region has selected Bedrock). Consumes stdin, so the request body must come from flags.",
		},
		&cli.StringFlag{
			Name:    "bedrock-base-url",
			Usage:   "Override the endpoint URL for Claude in Amazon Bedrock, once --bedrock-region or --bedrock-api-key-stdin has selected it. AWS credentials are sent only here or to the endpoint, never to --base-url.",
			Sources: cli.EnvVars("ANTHROPIC_BEDROCK_MANTLE_BASE_URL"),
		},
	)
}
