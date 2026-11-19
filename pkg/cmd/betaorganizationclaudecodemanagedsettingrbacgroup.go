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

var betaOrganizationClaudeCodeManagedSettingsRBACGroupsRetrieve = cli.Command{
	Name:    "retrieve",
	Usage:   "Retrieve an RBAC Group's Claude Code managed settings.",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:        "rbac-group-id",
			Usage:       "The RBAC Group's ID, prefixed `rbac_group_`.",
			Required:    true,
			PathParam:   "rbac_group_id",
			DataAliases: []string{"id"},
		},
		&requestflag.Flag[[]string]{
			Name:       "beta",
			Usage:      "This endpoint is in beta: requests must send `ce-claude-code-managed-settings-2026-10-06` in this header.",
			HeaderPath: "anthropic-beta",
		},
	},
	Action:          handleBetaOrganizationClaudeCodeManagedSettingsRBACGroupsRetrieve,
	HideHelpCommand: true,
}

var betaOrganizationClaudeCodeManagedSettingsRBACGroupsList = cli.Command{
	Name:    "list",
	Usage:   "List the RBAC Groups that have Claude Code managed settings, highest precedence\nfirst.",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[int64]{
			Name:      "limit",
			Usage:     "Number of items to return per page.\n\nDefaults to `20`. Ranges from `1` to `100`.",
			QueryPath: "limit",
		},
		&requestflag.Flag[string]{
			Name:      "page",
			Usage:     "Optionally set to the `next_page` token from the previous response.",
			QueryPath: "page",
		},
		&requestflag.Flag[[]string]{
			Name:       "beta",
			Usage:      "This endpoint is in beta: requests must send `ce-claude-code-managed-settings-2026-10-06` in this header.",
			HeaderPath: "anthropic-beta",
		},
		&requestflag.Flag[int64]{
			Name:  "max-items",
			Usage: "The maximum number of items to return (use -1 for unlimited).",
		},
	},
	Action:          handleBetaOrganizationClaudeCodeManagedSettingsRBACGroupsList,
	HideHelpCommand: true,
}

var betaOrganizationClaudeCodeManagedSettingsRBACGroupsDelete = cli.Command{
	Name:    "delete",
	Usage:   "Delete an RBAC Group's Claude Code managed settings and remove the group from\n`rbac_group_order`.",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:        "rbac-group-id",
			Usage:       "The RBAC Group's ID, prefixed `rbac_group_`.",
			Required:    true,
			PathParam:   "rbac_group_id",
			DataAliases: []string{"id"},
		},
		&requestflag.Flag[[]string]{
			Name:       "beta",
			Usage:      "This endpoint is in beta: requests must send `ce-claude-code-managed-settings-2026-10-06` in this header.",
			HeaderPath: "anthropic-beta",
		},
	},
	Action:          handleBetaOrganizationClaudeCodeManagedSettingsRBACGroupsDelete,
	HideHelpCommand: true,
}

var betaOrganizationClaudeCodeManagedSettingsRBACGroupsSet = cli.Command{
	Name:    "set",
	Usage:   "Create or replace an RBAC Group's Claude Code managed settings.",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:        "rbac-group-id",
			Usage:       "The RBAC Group's ID, prefixed `rbac_group_`.",
			Required:    true,
			PathParam:   "rbac_group_id",
			DataAliases: []string{"id"},
		},
		&requestflag.Flag[map[string]any]{
			Name:     "settings",
			Usage:    "The whole settings document, a JSON object of Claude Code managed settings. It replaces any stored document. At most 1 MiB as the server stores it, which can be larger than the request body.",
			Required: true,
			BodyPath: "settings",
		},
		&requestflag.Flag[[]string]{
			Name:       "beta",
			Usage:      "This endpoint is in beta: requests must send `ce-claude-code-managed-settings-2026-10-06` in this header.",
			HeaderPath: "anthropic-beta",
		},
	},
	Action:          handleBetaOrganizationClaudeCodeManagedSettingsRBACGroupsSet,
	HideHelpCommand: true,
}

func handleBetaOrganizationClaudeCodeManagedSettingsRBACGroupsRetrieve(ctx context.Context, cmd *cli.Command) error {
	client := anthropic.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("rbac-group-id") && len(unusedArgs) > 0 {
		cmd.Set("rbac-group-id", unusedArgs[0])
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

	params := anthropic.BetaOrganizationClaudeCodeManagedSettingRBACGroupGetParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Beta.Organization.ClaudeCode.ManagedSettings.RBACGroups.Get(
		ctx,
		cmd.Value("rbac-group-id").(string),
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
		Title:          "beta:organization:claude-code:managed-settings:rbac-groups retrieve",
		Transform:      transform,
	})
}

func handleBetaOrganizationClaudeCodeManagedSettingsRBACGroupsList(ctx context.Context, cmd *cli.Command) error {
	client := anthropic.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()

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

	params := anthropic.BetaOrganizationClaudeCodeManagedSettingRBACGroupListParams{}

	format := "explore"
	explicitFormat := cmd.Root().IsSet("format")
	if explicitFormat {
		format = cmd.Root().String("format")
	}
	transform := cmd.Root().String("transform")
	if format == "raw" {
		var res []byte
		options = append(options, option.WithResponseBodyInto(&res))
		_, err = client.Beta.Organization.ClaudeCode.ManagedSettings.RBACGroups.List(ctx, params, options...)
		if err != nil {
			return err
		}
		obj := gjson.ParseBytes(res)
		return ShowJSON(obj, ShowJSONOpts{
			ExplicitFormat: explicitFormat,
			Format:         format,
			RawOutput:      cmd.Root().Bool("raw-output"),
			Title:          "beta:organization:claude-code:managed-settings:rbac-groups list",
			Transform:      transform,
		})
	} else {
		iter := client.Beta.Organization.ClaudeCode.ManagedSettings.RBACGroups.ListAutoPaging(ctx, params, options...)
		maxItems := int64(-1)
		if cmd.IsSet("max-items") {
			maxItems = cmd.Value("max-items").(int64)
		}
		return ShowJSONIterator(iter, maxItems, ShowJSONOpts{
			ExplicitFormat: explicitFormat,
			Format:         format,
			RawOutput:      cmd.Root().Bool("raw-output"),
			Title:          "beta:organization:claude-code:managed-settings:rbac-groups list",
			Transform:      transform,
		})
	}
}

func handleBetaOrganizationClaudeCodeManagedSettingsRBACGroupsDelete(ctx context.Context, cmd *cli.Command) error {
	client := anthropic.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("rbac-group-id") && len(unusedArgs) > 0 {
		cmd.Set("rbac-group-id", unusedArgs[0])
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

	params := anthropic.BetaOrganizationClaudeCodeManagedSettingRBACGroupDeleteParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Beta.Organization.ClaudeCode.ManagedSettings.RBACGroups.Delete(
		ctx,
		cmd.Value("rbac-group-id").(string),
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
		Title:          "beta:organization:claude-code:managed-settings:rbac-groups delete",
		Transform:      transform,
	})
}

func handleBetaOrganizationClaudeCodeManagedSettingsRBACGroupsSet(ctx context.Context, cmd *cli.Command) error {
	client := anthropic.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("rbac-group-id") && len(unusedArgs) > 0 {
		cmd.Set("rbac-group-id", unusedArgs[0])
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

	params := anthropic.BetaOrganizationClaudeCodeManagedSettingRBACGroupSetParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Beta.Organization.ClaudeCode.ManagedSettings.RBACGroups.Set(
		ctx,
		cmd.Value("rbac-group-id").(string),
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
		Title:          "beta:organization:claude-code:managed-settings:rbac-groups set",
		Transform:      transform,
	})
}
