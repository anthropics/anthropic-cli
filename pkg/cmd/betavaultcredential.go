// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cmd

import (
	"context"
	"fmt"

	"github.com/anthropics/anthropic-cli/internal/apiquery"
	"github.com/anthropics/anthropic-cli/internal/requestflag"
	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/tidwall/gjson"
	"github.com/urfave/cli/v3"
)

var betaVaultsCredentialsCreate = requestflag.WithInnerFlags(cli.Command{
	Name:    "create",
	Usage:   "Create Credential",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "vault-id",
			Required:  true,
			PathParam: "vault_id",
		},
		&requestflag.Flag[map[string]any]{
			Name:     "auth",
			Usage:    "Authentication details for creating a credential.",
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
			Usage:      "Optional header to specify the beta version(s) you want to use.",
			HeaderPath: "anthropic-beta",
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
			Usage:      "Static bearer token value.",
			InnerField: "token",
		},
		&requestflag.InnerFlag[string]{
			Name:       "auth.access-token",
			Usage:      "OAuth access token.",
			InnerField: "access_token",
		},
		&requestflag.InnerFlag[any]{
			Name:       "auth.expires-at",
			Usage:      "A timestamp in RFC 3339 format",
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
			Usage:      "Substitute the secret on any host the session's Environment network policy permits egress to. The Environment's network policy is the only boundary on where the secret can reach.",
			InnerField: "networking",
		},
		&requestflag.InnerFlag[map[string]any]{
			Name:       "auth.refresh",
			Usage:      "OAuth refresh token parameters for creating a credential with refresh support.",
			InnerField: "refresh",
		},
		&requestflag.InnerFlag[string]{
			Name:       "auth.secret-name",
			Usage:      "Name of the environment variable. Immutable after create.",
			InnerField: "secret_name",
		},
		&requestflag.InnerFlag[string]{
			Name:       "auth.secret-value",
			Usage:      "Secret value. Write-only; never returned in responses.",
			InnerField: "secret_value",
		},
	},
})

var betaVaultsCredentialsRetrieve = cli.Command{
	Name:    "retrieve",
	Usage:   "Get Credential",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "vault-id",
			Required:  true,
			PathParam: "vault_id",
		},
		&requestflag.Flag[string]{
			Name:      "credential-id",
			Required:  true,
			PathParam: "credential_id",
		},
		&requestflag.Flag[[]string]{
			Name:       "beta",
			Usage:      "Optional header to specify the beta version(s) you want to use.",
			HeaderPath: "anthropic-beta",
		},
	},
	Action:          handleBetaVaultsCredentialsRetrieve,
	HideHelpCommand: true,
}

var betaVaultsCredentialsUpdate = requestflag.WithInnerFlags(cli.Command{
	Name:    "update",
	Usage:   "Update Credential",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "vault-id",
			Required:  true,
			PathParam: "vault_id",
		},
		&requestflag.Flag[string]{
			Name:        "credential-id",
			Required:    true,
			PathParam:   "credential_id",
			DataAliases: []string{"id"},
		},
		&requestflag.Flag[map[string]any]{
			Name:     "auth",
			Usage:    "Updated authentication details for a credential.",
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
			Usage:      "Optional header to specify the beta version(s) you want to use.",
			HeaderPath: "anthropic-beta",
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
			Usage:      "Updated static bearer token value.",
			InnerField: "token",
		},
		&requestflag.InnerFlag[*string]{
			Name:       "auth.access-token",
			Usage:      "Updated OAuth access token.",
			InnerField: "access_token",
		},
		&requestflag.InnerFlag[any]{
			Name:       "auth.expires-at",
			Usage:      "A timestamp in RFC 3339 format",
			InnerField: "expires_at",
		},
		&requestflag.InnerFlag[map[string]any]{
			Name:       "auth.injection-location",
			Usage:      "Updated injection location.",
			InnerField: "injection_location",
		},
		&requestflag.InnerFlag[map[string]any]{
			Name:       "auth.networking",
			Usage:      "Substitute the secret on any host the session's Environment network policy permits egress to. The Environment's network policy is the only boundary on where the secret can reach.",
			InnerField: "networking",
		},
		&requestflag.InnerFlag[map[string]any]{
			Name:       "auth.refresh",
			Usage:      "Parameters for updating OAuth refresh token configuration.",
			InnerField: "refresh",
		},
		&requestflag.InnerFlag[*string]{
			Name:       "auth.secret-value",
			Usage:      "Updated secret value.",
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
			Required:  true,
			PathParam: "vault_id",
		},
		&requestflag.Flag[bool]{
			Name:      "include-archived",
			Usage:     "Whether to include archived credentials in the results.",
			QueryPath: "include_archived",
		},
		&requestflag.Flag[int64]{
			Name:      "limit",
			Usage:     "Maximum number of credentials to return per page. Defaults to 20, maximum 100.",
			QueryPath: "limit",
		},
		&requestflag.Flag[string]{
			Name:      "page",
			Usage:     "Opaque pagination token from a previous `list_credentials` response.",
			QueryPath: "page",
		},
		&requestflag.Flag[[]string]{
			Name:       "beta",
			Usage:      "Optional header to specify the beta version(s) you want to use.",
			HeaderPath: "anthropic-beta",
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
	Usage:   "Delete Credential",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "vault-id",
			Required:  true,
			PathParam: "vault_id",
		},
		&requestflag.Flag[string]{
			Name:      "credential-id",
			Required:  true,
			PathParam: "credential_id",
		},
		&requestflag.Flag[[]string]{
			Name:       "beta",
			Usage:      "Optional header to specify the beta version(s) you want to use.",
			HeaderPath: "anthropic-beta",
		},
	},
	Action:          handleBetaVaultsCredentialsDelete,
	HideHelpCommand: true,
}

var betaVaultsCredentialsArchive = cli.Command{
	Name:    "archive",
	Usage:   "Archive Credential",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "vault-id",
			Required:  true,
			PathParam: "vault_id",
		},
		&requestflag.Flag[string]{
			Name:      "credential-id",
			Required:  true,
			PathParam: "credential_id",
		},
		&requestflag.Flag[[]string]{
			Name:       "beta",
			Usage:      "Optional header to specify the beta version(s) you want to use.",
			HeaderPath: "anthropic-beta",
		},
	},
	Action:          handleBetaVaultsCredentialsArchive,
	HideHelpCommand: true,
}

var betaVaultsCredentialsMCPOAuthValidate = cli.Command{
	Name:    "mcp-oauth-validate",
	Usage:   "Validate Credential",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "vault-id",
			Required:  true,
			PathParam: "vault_id",
		},
		&requestflag.Flag[string]{
			Name:      "credential-id",
			Required:  true,
			PathParam: "credential_id",
		},
		&requestflag.Flag[[]string]{
			Name:       "beta",
			Usage:      "Optional header to specify the beta version(s) you want to use.",
			HeaderPath: "anthropic-beta",
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
