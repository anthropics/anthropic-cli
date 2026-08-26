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

var betaOrganizationServiceAccountsWorkspacesList = cli.Command{
	Name:    "list",
	Usage:   "**Requires an OAuth access token with the `org:admin` scope**, from\n`ant auth login --scope org:admin` or a workload identity federation rule; Admin\nAPI keys are not accepted. See\n[Manage WIF with the Admin API](/docs/en/manage-claude/wif-admin-api).",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "service-account-id",
			Usage:     "ID of the service account.",
			Required:  true,
			PathParam: "service_account_id",
		},
		&requestflag.Flag[int64]{
			Name:      "limit",
			Usage:     "Number of results per page.",
			Default:   20,
			QueryPath: "limit",
		},
		&requestflag.Flag[string]{
			Name:      "page",
			Usage:     "Opaque cursor from a previous response's `next_page`.",
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
	Action:          handleBetaOrganizationServiceAccountsWorkspacesList,
	HideHelpCommand: true,
}

var betaOrganizationServiceAccountsWorkspacesAdd = cli.Command{
	Name:    "add",
	Usage:   "**Requires an OAuth access token with the `org:admin` scope**, from\n`ant auth login --scope org:admin` or a workload identity federation rule; Admin\nAPI keys are not accepted. See\n[Manage WIF with the Admin API](/docs/en/manage-claude/wif-admin-api).",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "service-account-id",
			Usage:     "ID of the service account.",
			Required:  true,
			PathParam: "service_account_id",
		},
		&requestflag.Flag[string]{
			Name:     "workspace-id",
			Usage:    "Tagged workspace ID to add the service account to.",
			Required: true,
			BodyPath: "workspace_id",
		},
		&requestflag.Flag[string]{
			Name:     "workspace-role",
			Usage:    `Allowed values: "workspace_admin", "workspace_developer", "workspace_restricted_developer", "workspace_user".`,
			Required: true,
			BodyPath: "workspace_role",
		},
		&requestflag.Flag[[]string]{
			Name:       "beta",
			Usage:      "Optional header to specify the beta version(s) you want to use.",
			HeaderPath: "anthropic-beta",
		},
	},
	Action:          handleBetaOrganizationServiceAccountsWorkspacesAdd,
	HideHelpCommand: true,
}

var betaOrganizationServiceAccountsWorkspacesRemove = cli.Command{
	Name:    "remove",
	Usage:   "**Requires an OAuth access token with the `org:admin` scope**, from\n`ant auth login --scope org:admin` or a workload identity federation rule; Admin\nAPI keys are not accepted. See\n[Manage WIF with the Admin API](/docs/en/manage-claude/wif-admin-api).",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "service-account-id",
			Usage:     "ID of the service account.",
			Required:  true,
			PathParam: "service_account_id",
		},
		&requestflag.Flag[string]{
			Name:      "workspace-id",
			Usage:     "ID of the workspace.",
			Required:  true,
			PathParam: "workspace_id",
		},
		&requestflag.Flag[[]string]{
			Name:       "beta",
			Usage:      "Optional header to specify the beta version(s) you want to use.",
			HeaderPath: "anthropic-beta",
		},
	},
	Action:          handleBetaOrganizationServiceAccountsWorkspacesRemove,
	HideHelpCommand: true,
}

func handleBetaOrganizationServiceAccountsWorkspacesList(ctx context.Context, cmd *cli.Command) error {
	client := anthropic.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("service-account-id") && len(unusedArgs) > 0 {
		cmd.Set("service-account-id", unusedArgs[0])
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

	params := anthropic.BetaOrganizationServiceAccountWorkspaceListParams{}

	format := "explore"
	explicitFormat := cmd.Root().IsSet("format")
	if explicitFormat {
		format = cmd.Root().String("format")
	}
	transform := cmd.Root().String("transform")
	if format == "raw" {
		var res []byte
		options = append(options, option.WithResponseBodyInto(&res))
		_, err = client.Beta.Organization.ServiceAccounts.Workspaces.List(
			ctx,
			cmd.Value("service-account-id").(string),
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
			Title:          "beta:organization:service-accounts:workspaces list",
			Transform:      transform,
		})
	} else {
		iter := client.Beta.Organization.ServiceAccounts.Workspaces.ListAutoPaging(
			ctx,
			cmd.Value("service-account-id").(string),
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
			Title:          "beta:organization:service-accounts:workspaces list",
			Transform:      transform,
		})
	}
}

func handleBetaOrganizationServiceAccountsWorkspacesAdd(ctx context.Context, cmd *cli.Command) error {
	client := anthropic.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("service-account-id") && len(unusedArgs) > 0 {
		cmd.Set("service-account-id", unusedArgs[0])
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

	params := anthropic.BetaOrganizationServiceAccountWorkspaceAddParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Beta.Organization.ServiceAccounts.Workspaces.Add(
		ctx,
		cmd.Value("service-account-id").(string),
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
		Title:          "beta:organization:service-accounts:workspaces add",
		Transform:      transform,
	})
}

func handleBetaOrganizationServiceAccountsWorkspacesRemove(ctx context.Context, cmd *cli.Command) error {
	client := anthropic.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("workspace-id") && len(unusedArgs) > 0 {
		cmd.Set("workspace-id", unusedArgs[0])
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

	params := anthropic.BetaOrganizationServiceAccountWorkspaceRemoveParams{
		ServiceAccountID: cmd.Value("service-account-id").(string),
	}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Beta.Organization.ServiceAccounts.Workspaces.Remove(
		ctx,
		cmd.Value("workspace-id").(string),
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
		Title:          "beta:organization:service-accounts:workspaces remove",
		Transform:      transform,
	})
}
