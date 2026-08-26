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

var betaOrganizationFederationRulesWorkspacesList = cli.Command{
	Name:    "list",
	Usage:   "**Requires an OAuth access token with the `org:admin` scope**, from\n`ant auth login --scope org:admin` or a workload identity federation rule; Admin\nAPI keys are not accepted. See\n[Manage WIF with the Admin API](/docs/en/manage-claude/wif-admin-api).",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "federation-rule-id",
			Usage:     "ID of the federation rule.",
			Required:  true,
			PathParam: "federation_rule_id",
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
	Action:          handleBetaOrganizationFederationRulesWorkspacesList,
	HideHelpCommand: true,
}

var betaOrganizationFederationRulesWorkspacesAdd = cli.Command{
	Name:    "add",
	Usage:   "**Requires an OAuth access token with the `org:admin` scope**, from\n`ant auth login --scope org:admin` or a workload identity federation rule; Admin\nAPI keys are not accepted. See\n[Manage WIF with the Admin API](/docs/en/manage-claude/wif-admin-api).",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "federation-rule-id",
			Usage:     "ID of the federation rule.",
			Required:  true,
			PathParam: "federation_rule_id",
		},
		&requestflag.Flag[string]{
			Name:     "workspace-id",
			Usage:    "Tagged ID of the workspace to enable this rule for.",
			Required: true,
			BodyPath: "workspace_id",
		},
		&requestflag.Flag[[]string]{
			Name:       "beta",
			Usage:      "Optional header to specify the beta version(s) you want to use.",
			HeaderPath: "anthropic-beta",
		},
	},
	Action:          handleBetaOrganizationFederationRulesWorkspacesAdd,
	HideHelpCommand: true,
}

var betaOrganizationFederationRulesWorkspacesRemove = cli.Command{
	Name:    "remove",
	Usage:   "**Requires an OAuth access token with the `org:admin` scope**, from\n`ant auth login --scope org:admin` or a workload identity federation rule; Admin\nAPI keys are not accepted. See\n[Manage WIF with the Admin API](/docs/en/manage-claude/wif-admin-api).",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "federation-rule-id",
			Usage:     "ID of the federation rule.",
			Required:  true,
			PathParam: "federation_rule_id",
		},
		&requestflag.Flag[string]{
			Name:      "workspace-id",
			Usage:     "ID of the workspace to disable for.",
			Required:  true,
			PathParam: "workspace_id",
		},
		&requestflag.Flag[[]string]{
			Name:       "beta",
			Usage:      "Optional header to specify the beta version(s) you want to use.",
			HeaderPath: "anthropic-beta",
		},
	},
	Action:          handleBetaOrganizationFederationRulesWorkspacesRemove,
	HideHelpCommand: true,
}

func handleBetaOrganizationFederationRulesWorkspacesList(ctx context.Context, cmd *cli.Command) error {
	client := anthropic.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("federation-rule-id") && len(unusedArgs) > 0 {
		cmd.Set("federation-rule-id", unusedArgs[0])
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

	params := anthropic.BetaOrganizationFederationRuleWorkspaceListParams{}

	format := "explore"
	explicitFormat := cmd.Root().IsSet("format")
	if explicitFormat {
		format = cmd.Root().String("format")
	}
	transform := cmd.Root().String("transform")
	if format == "raw" {
		var res []byte
		options = append(options, option.WithResponseBodyInto(&res))
		_, err = client.Beta.Organization.Federation.Rules.Workspaces.List(
			ctx,
			cmd.Value("federation-rule-id").(string),
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
			Title:          "beta:organization:federation:rules:workspaces list",
			Transform:      transform,
		})
	} else {
		iter := client.Beta.Organization.Federation.Rules.Workspaces.ListAutoPaging(
			ctx,
			cmd.Value("federation-rule-id").(string),
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
			Title:          "beta:organization:federation:rules:workspaces list",
			Transform:      transform,
		})
	}
}

func handleBetaOrganizationFederationRulesWorkspacesAdd(ctx context.Context, cmd *cli.Command) error {
	client := anthropic.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("federation-rule-id") && len(unusedArgs) > 0 {
		cmd.Set("federation-rule-id", unusedArgs[0])
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

	params := anthropic.BetaOrganizationFederationRuleWorkspaceAddParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Beta.Organization.Federation.Rules.Workspaces.Add(
		ctx,
		cmd.Value("federation-rule-id").(string),
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
		Title:          "beta:organization:federation:rules:workspaces add",
		Transform:      transform,
	})
}

func handleBetaOrganizationFederationRulesWorkspacesRemove(ctx context.Context, cmd *cli.Command) error {
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

	params := anthropic.BetaOrganizationFederationRuleWorkspaceRemoveParams{
		FederationRuleID: cmd.Value("federation-rule-id").(string),
	}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Beta.Organization.Federation.Rules.Workspaces.Remove(
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
		Title:          "beta:organization:federation:rules:workspaces remove",
		Transform:      transform,
	})
}
