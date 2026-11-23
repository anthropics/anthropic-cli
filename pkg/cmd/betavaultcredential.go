package cmd

import (
	"context"
	"fmt"

	"github.com/anthropics/anthropic-cli/v2/internal/apiquery"
	"github.com/anthropics/anthropic-cli/v2/internal/requestflag"
	"github.com/anthropics/anthropic-sdk-go/v2"
	"github.com/anthropics/anthropic-sdk-go/v2/option"
	"github.com/tidwall/gjson"
	"github.com/urfave/cli/v3"
)

var betaVaultsCredentialsCreate = requestflag.WithInnerFlags(cli.Command{
	Name:    "create",
	Usage:   "Create a Credential",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "vault-id",
			Usage:     "Identifier of the vault to create the credential in.",
			Required:  true,
			PathParam: "vault_id",
		},
		&requestflag.Flag[map[string]any]{
			Name:     "auth",
			Usage:    "Authentication configuration for the credential.",
			Required: true,
			BodyPath: "auth",
		},
		&requestflag.Flag[*string]{
			Name:     "display-name",
			Usage:    "Human-readable name for the credential. Up to 255 characters.",
			BodyPath: "display_name",
		},
		&requestflag.Flag[map[string]any]{
			Name:     "metadata",
			Usage:    "Arbitrary key-value metadata to attach to the credential. Maximum 16 pairs, keys up to 64 chars, values up to 512 chars.",
			BodyPath: "metadata",
		},
		&requestflag.Flag[[]string]{
			Name:       "beta",
			Usage:      "This endpoint is in beta: requests must send `managed-agents-2026-04-01` in this header.",
			HeaderPath: "anthropic-beta",
		},
		&requestflag.Flag[string]{
			Name:       "workspace-id",
			Usage:      "Optional header to select the Workspace for this request. The value is a Workspace ID (for example, `wrkspc_011CZkZaBF1tNoB5wlCeusgy`).\n\nOnly needed for credentials that can act on more than one Workspace. A credential that belongs to a specific Workspace may omit it; if sent, it must match that Workspace.",
			HeaderPath: "anthropic-workspace-id",
		},
	},
	Action:          handleBetaVaultsCredentialsCreate,
	HideHelpCommand: true,
}, map[string][]requestflag.HasOuterFlag{
	"auth": {
		&requestflag.InnerFlag[string]{
			Name:       "auth.type",
			Usage:      `Allowed values: "mcp_oauth", "static_bearer", "environment_variable".`,
			InnerField: "type",
		},
		&requestflag.InnerFlag[string]{
			Name:       "auth.token",
			Usage:      "The static bearer token to send to the MCP server, such as an API key or personal access token.",
			InnerField: "token",
		},
		&requestflag.InnerFlag[string]{
			Name:       "auth.access-token",
			Usage:      "The OAuth access token to send to the MCP server.",
			InnerField: "access_token",
		},
		&requestflag.InnerFlag[any]{
			Name:       "auth.expires-at",
			Usage:      "When `access_token` expires, in RFC 3339.\n\nMust be in the future.",
			InnerField: "expires_at",
		},
		&requestflag.InnerFlag[map[string]any]{
			Name:       "auth.injection-location",
			Usage:      "Where in the outbound request the secret value may be substituted.",
			InnerField: "injection_location",
		},
		&requestflag.InnerFlag[string]{
			Name:       "auth.mcp-server-url",
			Usage:      "URL of the MCP server this credential authenticates against.",
			InnerField: "mcp_server_url",
		},
		&requestflag.InnerFlag[map[string]any]{
			Name:       "auth.networking",
			Usage:      "Outbound hosts the secret value is substituted on.",
			InnerField: "networking",
		},
		&requestflag.InnerFlag[map[string]any]{
			Name:       "auth.refresh",
			Usage:      "Refresh token configuration, if the credential supports token refresh.",
			InnerField: "refresh",
		},
		&requestflag.InnerFlag[string]{
			Name:       "auth.secret-name",
			Usage:      "Name of the environment variable. Immutable after create.",
			InnerField: "secret_name",
		},
		&requestflag.InnerFlag[string]{
			Name:       "auth.secret-value",
			Usage:      "The secret to substitute into the agent's outbound requests. It's write-only and never returned in responses.",
			InnerField: "secret_value",
		},
	},
})

var betaVaultsCredentialsRetrieve = cli.Command{
	Name:    "retrieve",
	Usage:   "Retrieve a Credential",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "vault-id",
			Usage:     "Identifier of the vault containing the credential.",
			Required:  true,
			PathParam: "vault_id",
		},
		&requestflag.Flag[string]{
			Name:      "credential-id",
			Usage:     "Unique identifier of the credential to retrieve.",
			Required:  true,
			PathParam: "credential_id",
		},
		&requestflag.Flag[[]string]{
			Name:       "beta",
			Usage:      "This endpoint is in beta: requests must send `managed-agents-2026-04-01` in this header.",
			HeaderPath: "anthropic-beta",
		},
		&requestflag.Flag[string]{
			Name:       "workspace-id",
			Usage:      "Optional header to select the Workspace for this request. The value is a Workspace ID (for example, `wrkspc_011CZkZaBF1tNoB5wlCeusgy`).\n\nOnly needed for credentials that can act on more than one Workspace. A credential that belongs to a specific Workspace may omit it; if sent, it must match that Workspace.",
			HeaderPath: "anthropic-workspace-id",
		},
	},
	Action:          handleBetaVaultsCredentialsRetrieve,
	HideHelpCommand: true,
}

var betaVaultsCredentialsUpdate = requestflag.WithInnerFlags(cli.Command{
	Name:    "update",
	Usage:   "Update a Credential",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "vault-id",
			Usage:     "Identifier of the vault containing the credential.",
			Required:  true,
			PathParam: "vault_id",
		},
		&requestflag.Flag[string]{
			Name:        "credential-id",
			Usage:       "Unique identifier of the credential to update.",
			Required:    true,
			PathParam:   "credential_id",
			DataAliases: []string{"id"},
		},
		&requestflag.Flag[map[string]any]{
			Name:     "auth",
			Usage:    "New authentication settings for the credential. The `type` must match the credential's current type, which can't be changed.",
			BodyPath: "auth",
		},
		&requestflag.Flag[*string]{
			Name:     "display-name",
			Usage:    "Updated human-readable name for the credential. 1-255 characters.",
			BodyPath: "display_name",
		},
		&requestflag.Flag[map[string]any]{
			Name:     "metadata",
			Usage:    "Metadata patch. Set a key to a string to upsert it, or to null to delete it. Omitted keys are preserved.",
			BodyPath: "metadata",
		},
		&requestflag.Flag[[]string]{
			Name:       "beta",
			Usage:      "This endpoint is in beta: requests must send `managed-agents-2026-04-01` in this header.",
			HeaderPath: "anthropic-beta",
		},
		&requestflag.Flag[string]{
			Name:       "workspace-id",
			Usage:      "Optional header to select the Workspace for this request. The value is a Workspace ID (for example, `wrkspc_011CZkZaBF1tNoB5wlCeusgy`).\n\nOnly needed for credentials that can act on more than one Workspace. A credential that belongs to a specific Workspace may omit it; if sent, it must match that Workspace.",
			HeaderPath: "anthropic-workspace-id",
		},
	},
	Action:          handleBetaVaultsCredentialsUpdate,
	HideHelpCommand: true,
}, map[string][]requestflag.HasOuterFlag{
	"auth": {
		&requestflag.InnerFlag[string]{
			Name:       "auth.type",
			Usage:      `Allowed values: "mcp_oauth", "static_bearer", "environment_variable".`,
			InnerField: "type",
		},
		&requestflag.InnerFlag[*string]{
			Name:       "auth.token",
			Usage:      "The new static bearer token to replace the stored one.",
			InnerField: "token",
		},
		&requestflag.InnerFlag[*string]{
			Name:       "auth.access-token",
			Usage:      "The new OAuth access token to replace the stored one.",
			InnerField: "access_token",
		},
		&requestflag.InnerFlag[any]{
			Name:       "auth.expires-at",
			Usage:      "When the credential's access token expires, in RFC 3339. Leave it out to keep the stored expiry.\n\nIf you send a new `access_token`, send its `expires_at` too, or the new token keeps the old expiry.\n\nMust be in the future.",
			InnerField: "expires_at",
		},
		&requestflag.InnerFlag[map[string]any]{
			Name:       "auth.injection-location",
			Usage:      "Updated injection location.",
			InnerField: "injection_location",
		},
		&requestflag.InnerFlag[map[string]any]{
			Name:       "auth.networking",
			Usage:      "The new networking scope, which replaces the stored one as a whole.",
			InnerField: "networking",
		},
		&requestflag.InnerFlag[map[string]any]{
			Name:       "auth.refresh",
			Usage:      "Updated refresh token configuration.",
			InnerField: "refresh",
		},
		&requestflag.InnerFlag[*string]{
			Name:       "auth.secret-value",
			Usage:      "The new secret value to replace the stored one.",
			InnerField: "secret_value",
		},
	},
})

var betaVaultsCredentialsList = cli.Command{
	Name:    "list",
	Usage:   "List Credentials",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "vault-id",
			Usage:     "Identifier of the vault to list credentials for.",
			Required:  true,
			PathParam: "vault_id",
		},
		&requestflag.Flag[bool]{
			Name:      "include-archived",
			Usage:     "Whether to include archived credentials in the results. Defaults to `false`.",
			QueryPath: "include_archived",
		},
		&requestflag.Flag[int64]{
			Name:      "limit",
			Usage:     "Maximum number of credentials to return per page. Defaults to 20, maximum 100.",
			QueryPath: "limit",
		},
		&requestflag.Flag[string]{
			Name:      "page",
			Usage:     "The `next_page` value from a previous [List Credentials](https://platform.claude.com/docs/en/api/beta/vaults/credentials/list) response, to get the next page. Leave it out to get the first page.",
			QueryPath: "page",
		},
		&requestflag.Flag[[]string]{
			Name:       "beta",
			Usage:      "This endpoint is in beta: requests must send `managed-agents-2026-04-01` in this header.",
			HeaderPath: "anthropic-beta",
		},
		&requestflag.Flag[string]{
			Name:       "workspace-id",
			Usage:      "Optional header to select the Workspace for this request. The value is a Workspace ID (for example, `wrkspc_011CZkZaBF1tNoB5wlCeusgy`).\n\nOnly needed for credentials that can act on more than one Workspace. A credential that belongs to a specific Workspace may omit it; if sent, it must match that Workspace.",
			HeaderPath: "anthropic-workspace-id",
		},
		&requestflag.Flag[int64]{
			Name:  "max-items",
			Usage: "The maximum number of items to return (use -1 for unlimited).",
		},
	},
	Action:          handleBetaVaultsCredentialsList,
	HideHelpCommand: true,
}

var betaVaultsCredentialsDelete = cli.Command{
	Name:    "delete",
	Usage:   "Delete a Credential",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "vault-id",
			Usage:     "Identifier of the vault containing the credential.",
			Required:  true,
			PathParam: "vault_id",
		},
		&requestflag.Flag[string]{
			Name:      "credential-id",
			Usage:     "Unique identifier of the credential to delete.",
			Required:  true,
			PathParam: "credential_id",
		},
		&requestflag.Flag[[]string]{
			Name:       "beta",
			Usage:      "This endpoint is in beta: requests must send `managed-agents-2026-04-01` in this header.",
			HeaderPath: "anthropic-beta",
		},
		&requestflag.Flag[string]{
			Name:       "workspace-id",
			Usage:      "Optional header to select the Workspace for this request. The value is a Workspace ID (for example, `wrkspc_011CZkZaBF1tNoB5wlCeusgy`).\n\nOnly needed for credentials that can act on more than one Workspace. A credential that belongs to a specific Workspace may omit it; if sent, it must match that Workspace.",
			HeaderPath: "anthropic-workspace-id",
		},
	},
	Action:          handleBetaVaultsCredentialsDelete,
	HideHelpCommand: true,
}

var betaVaultsCredentialsArchive = cli.Command{
	Name:    "archive",
	Usage:   "Archive a Credential",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "vault-id",
			Usage:     "Identifier of the vault containing the credential.",
			Required:  true,
			PathParam: "vault_id",
		},
		&requestflag.Flag[string]{
			Name:      "credential-id",
			Usage:     "Unique identifier of the credential to archive.",
			Required:  true,
			PathParam: "credential_id",
		},
		&requestflag.Flag[[]string]{
			Name:       "beta",
			Usage:      "This endpoint is in beta: requests must send `managed-agents-2026-04-01` in this header.",
			HeaderPath: "anthropic-beta",
		},
		&requestflag.Flag[string]{
			Name:       "workspace-id",
			Usage:      "Optional header to select the Workspace for this request. The value is a Workspace ID (for example, `wrkspc_011CZkZaBF1tNoB5wlCeusgy`).\n\nOnly needed for credentials that can act on more than one Workspace. A credential that belongs to a specific Workspace may omit it; if sent, it must match that Workspace.",
			HeaderPath: "anthropic-workspace-id",
		},
	},
	Action:          handleBetaVaultsCredentialsArchive,
	HideHelpCommand: true,
}

var betaVaultsCredentialsMCPOAuthValidate = cli.Command{
	Name:    "mcp-oauth-validate",
	Usage:   "Validate an MCP OAuth Credential",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "vault-id",
			Usage:     "Identifier of the vault containing the credential.",
			Required:  true,
			PathParam: "vault_id",
		},
		&requestflag.Flag[string]{
			Name:      "credential-id",
			Usage:     "Unique identifier of the credential to validate.",
			Required:  true,
			PathParam: "credential_id",
		},
		&requestflag.Flag[[]string]{
			Name:       "beta",
			Usage:      "This endpoint is in beta: requests must send `managed-agents-2026-04-01` in this header.",
			HeaderPath: "anthropic-beta",
		},
		&requestflag.Flag[string]{
			Name:       "workspace-id",
			Usage:      "Optional header to select the Workspace for this request. The value is a Workspace ID (for example, `wrkspc_011CZkZaBF1tNoB5wlCeusgy`).\n\nOnly needed for credentials that can act on more than one Workspace. A credential that belongs to a specific Workspace may omit it; if sent, it must match that Workspace.",
			HeaderPath: "anthropic-workspace-id",
		},
	},
	Action:          handleBetaVaultsCredentialsMCPOAuthValidate,
	HideHelpCommand: true,
}

func handleBetaVaultsCredentialsCreate(ctx context.Context, cmd *cli.Command) error {
	client := anthropic.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("vault-id") && len(unusedArgs) > 0 {
		cmd.Set("vault-id", unusedArgs[0])
		unusedArgs = unusedArgs[1:]
	}
	if len(unusedArgs) > 0 {
		return fmt.Errorf("Unexpected extra arguments: %v", unusedArgs)
	}

	options, err := flagOptions(
		cmd,
		apiquery.NestedQueryFormatBrackets,
		apiquery.ArrayQueryFormatBrackets,
		ApplicationJSON,
		false,
	)
	if err != nil {
		return err
	}

	params := anthropic.BetaVaultCredentialNewParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Beta.Vaults.Credentials.New(
		ctx,
		cmd.Value("vault-id").(string),
		params,
		options...,
	)
	if err != nil {
		return err
	}

	obj := gjson.ParseBytes(res)
	format := cmd.Root().String("format")
	explicitFormat := cmd.Root().IsSet("format")
	transform := cmd.Root().String("transform")
	return ShowJSON(obj, ShowJSONOpts{
		ExplicitFormat: explicitFormat,
		Format:         format,
		RawOutput:      cmd.Root().Bool("raw-output"),
		Title:          "beta:vaults:credentials create",
		Transform:      transform,
	})
}

func handleBetaVaultsCredentialsRetrieve(ctx context.Context, cmd *cli.Command) error {
	client := anthropic.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("credential-id") && len(unusedArgs) > 0 {
		cmd.Set("credential-id", unusedArgs[0])
		unusedArgs = unusedArgs[1:]
	}
	if len(unusedArgs) > 0 {
		return fmt.Errorf("Unexpected extra arguments: %v", unusedArgs)
	}

	options, err := flagOptions(
		cmd,
		apiquery.NestedQueryFormatBrackets,
		apiquery.ArrayQueryFormatBrackets,
		EmptyBody,
		false,
	)
	if err != nil {
		return err
	}

	params := anthropic.BetaVaultCredentialGetParams{
		VaultID: cmd.Value("vault-id").(string),
	}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Beta.Vaults.Credentials.Get(
		ctx,
		cmd.Value("credential-id").(string),
		params,
		options...,
	)
	if err != nil {
		return err
	}

	obj := gjson.ParseBytes(res)
	format := "explore"
	explicitFormat := cmd.Root().IsSet("format")
	if explicitFormat {
		format = cmd.Root().String("format")
	}
	transform := cmd.Root().String("transform")
	return ShowJSON(obj, ShowJSONOpts{
		ExplicitFormat: explicitFormat,
		Format:         format,
		RawOutput:      cmd.Root().Bool("raw-output"),
		Title:          "beta:vaults:credentials retrieve",
		Transform:      transform,
	})
}

func handleBetaVaultsCredentialsUpdate(ctx context.Context, cmd *cli.Command) error {
	client := anthropic.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("credential-id") && len(unusedArgs) > 0 {
		cmd.Set("credential-id", unusedArgs[0])
		unusedArgs = unusedArgs[1:]
	}
	if len(unusedArgs) > 0 {
		return fmt.Errorf("Unexpected extra arguments: %v", unusedArgs)
	}

	options, err := flagOptions(
		cmd,
		apiquery.NestedQueryFormatBrackets,
		apiquery.ArrayQueryFormatBrackets,
		ApplicationJSON,
		false,
	)
	if err != nil {
		return err
	}

	params := anthropic.BetaVaultCredentialUpdateParams{
		VaultID: cmd.Value("vault-id").(string),
	}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Beta.Vaults.Credentials.Update(
		ctx,
		cmd.Value("credential-id").(string),
		params,
		options...,
	)
	if err != nil {
		return err
	}

	obj := gjson.ParseBytes(res)
	format := cmd.Root().String("format")
	explicitFormat := cmd.Root().IsSet("format")
	transform := cmd.Root().String("transform")
	return ShowJSON(obj, ShowJSONOpts{
		ExplicitFormat: explicitFormat,
		Format:         format,
		RawOutput:      cmd.Root().Bool("raw-output"),
		Title:          "beta:vaults:credentials update",
		Transform:      transform,
	})
}

func handleBetaVaultsCredentialsList(ctx context.Context, cmd *cli.Command) error {
	client := anthropic.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("vault-id") && len(unusedArgs) > 0 {
		cmd.Set("vault-id", unusedArgs[0])
		unusedArgs = unusedArgs[1:]
	}
	if len(unusedArgs) > 0 {
		return fmt.Errorf("Unexpected extra arguments: %v", unusedArgs)
	}

	options, err := flagOptions(
		cmd,
		apiquery.NestedQueryFormatBrackets,
		apiquery.ArrayQueryFormatBrackets,
		EmptyBody,
		false,
	)
	if err != nil {
		return err
	}

	params := anthropic.BetaVaultCredentialListParams{}

	format := "explore"
	explicitFormat := cmd.Root().IsSet("format")
	if explicitFormat {
		format = cmd.Root().String("format")
	}
	transform := cmd.Root().String("transform")
	if format == "raw" {
		var res []byte
		options = append(options, option.WithResponseBodyInto(&res))
		_, err = client.Beta.Vaults.Credentials.List(
			ctx,
			cmd.Value("vault-id").(string),
			params,
			options...,
		)
		if err != nil {
			return err
		}
		obj := gjson.ParseBytes(res)
		return ShowJSON(obj, ShowJSONOpts{
			ExplicitFormat: explicitFormat,
			Format:         format,
			RawOutput:      cmd.Root().Bool("raw-output"),
			Title:          "beta:vaults:credentials list",
			Transform:      transform,
		})
	} else {
		iter := client.Beta.Vaults.Credentials.ListAutoPaging(
			ctx,
			cmd.Value("vault-id").(string),
			params,
			options...,
		)
		maxItems := int64(-1)
		if cmd.IsSet("max-items") {
			maxItems = cmd.Value("max-items").(int64)
		}
		return ShowJSONIterator(iter, maxItems, ShowJSONOpts{
			ExplicitFormat: explicitFormat,
			Format:         format,
			RawOutput:      cmd.Root().Bool("raw-output"),
			Title:          "beta:vaults:credentials list",
			Transform:      transform,
		})
	}
}

func handleBetaVaultsCredentialsDelete(ctx context.Context, cmd *cli.Command) error {
	client := anthropic.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("credential-id") && len(unusedArgs) > 0 {
		cmd.Set("credential-id", unusedArgs[0])
		unusedArgs = unusedArgs[1:]
	}
	if len(unusedArgs) > 0 {
		return fmt.Errorf("Unexpected extra arguments: %v", unusedArgs)
	}

	options, err := flagOptions(
		cmd,
		apiquery.NestedQueryFormatBrackets,
		apiquery.ArrayQueryFormatBrackets,
		EmptyBody,
		false,
	)
	if err != nil {
		return err
	}

	params := anthropic.BetaVaultCredentialDeleteParams{
		VaultID: cmd.Value("vault-id").(string),
	}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Beta.Vaults.Credentials.Delete(
		ctx,
		cmd.Value("credential-id").(string),
		params,
		options...,
	)
	if err != nil {
		return err
	}

	obj := gjson.ParseBytes(res)
	format := cmd.Root().String("format")
	explicitFormat := cmd.Root().IsSet("format")
	transform := cmd.Root().String("transform")
	return ShowJSON(obj, ShowJSONOpts{
		ExplicitFormat: explicitFormat,
		Format:         format,
		RawOutput:      cmd.Root().Bool("raw-output"),
		Title:          "beta:vaults:credentials delete",
		Transform:      transform,
	})
}

func handleBetaVaultsCredentialsArchive(ctx context.Context, cmd *cli.Command) error {
	client := anthropic.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("credential-id") && len(unusedArgs) > 0 {
		cmd.Set("credential-id", unusedArgs[0])
		unusedArgs = unusedArgs[1:]
	}
	if len(unusedArgs) > 0 {
		return fmt.Errorf("Unexpected extra arguments: %v", unusedArgs)
	}

	options, err := flagOptions(
		cmd,
		apiquery.NestedQueryFormatBrackets,
		apiquery.ArrayQueryFormatBrackets,
		EmptyBody,
		false,
	)
	if err != nil {
		return err
	}

	params := anthropic.BetaVaultCredentialArchiveParams{
		VaultID: cmd.Value("vault-id").(string),
	}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Beta.Vaults.Credentials.Archive(
		ctx,
		cmd.Value("credential-id").(string),
		params,
		options...,
	)
	if err != nil {
		return err
	}

	obj := gjson.ParseBytes(res)
	format := cmd.Root().String("format")
	explicitFormat := cmd.Root().IsSet("format")
	transform := cmd.Root().String("transform")
	return ShowJSON(obj, ShowJSONOpts{
		ExplicitFormat: explicitFormat,
		Format:         format,
		RawOutput:      cmd.Root().Bool("raw-output"),
		Title:          "beta:vaults:credentials archive",
		Transform:      transform,
	})
}

func handleBetaVaultsCredentialsMCPOAuthValidate(ctx context.Context, cmd *cli.Command) error {
	client := anthropic.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("credential-id") && len(unusedArgs) > 0 {
		cmd.Set("credential-id", unusedArgs[0])
		unusedArgs = unusedArgs[1:]
	}
	if len(unusedArgs) > 0 {
		return fmt.Errorf("Unexpected extra arguments: %v", unusedArgs)
	}

	options, err := flagOptions(
		cmd,
		apiquery.NestedQueryFormatBrackets,
		apiquery.ArrayQueryFormatBrackets,
		EmptyBody,
		false,
	)
	if err != nil {
		return err
	}

	params := anthropic.BetaVaultCredentialMCPOAuthValidateParams{
		VaultID: cmd.Value("vault-id").(string),
	}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Beta.Vaults.Credentials.MCPOAuthValidate(
		ctx,
		cmd.Value("credential-id").(string),
		params,
		options...,
	)
	if err != nil {
		return err
	}

	obj := gjson.ParseBytes(res)
	format := cmd.Root().String("format")
	explicitFormat := cmd.Root().IsSet("format")
	transform := cmd.Root().String("transform")
	return ShowJSON(obj, ShowJSONOpts{
		ExplicitFormat: explicitFormat,
		Format:         format,
		RawOutput:      cmd.Root().Bool("raw-output"),
		Title:          "beta:vaults:credentials mcp-oauth-validate",
		Transform:      transform,
	})
}
