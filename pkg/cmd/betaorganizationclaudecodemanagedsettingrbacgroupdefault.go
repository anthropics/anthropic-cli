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

var betaOrganizationClaudeCodeManagedSettingsRBACGroupDefaultRetrieve = cli.Command{
	Name:    "retrieve",
	Usage:   "Retrieve the default Claude Code managed settings, which members who belong to\nno RBAC Group in `rbac_group_order` receive in `rbac_groups` mode.",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[[]string]{
			Name:       "beta",
			Usage:      "This endpoint is in beta: requests must send `ce-claude-code-managed-settings-2026-10-06` in this header.",
			HeaderPath: "anthropic-beta",
		},
	},
	Action:          handleBetaOrganizationClaudeCodeManagedSettingsRBACGroupDefaultRetrieve,
	HideHelpCommand: true,
}

var betaOrganizationClaudeCodeManagedSettingsRBACGroupDefaultUpdate = cli.Command{
	Name:    "update",
	Usage:   "Set or replace the default Claude Code managed settings, or clear them with\n`null`.",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[map[string]any]{
			Name:     "settings",
			Usage:    "The whole settings document, a JSON object of Claude Code managed settings. It replaces any stored document, or `null` clears the default settings. At most 1 MiB as the server stores it, which can be larger than the request body.",
			Required: true,
			BodyPath: "settings",
		},
		&requestflag.Flag[[]string]{
			Name:       "beta",
			Usage:      "This endpoint is in beta: requests must send `ce-claude-code-managed-settings-2026-10-06` in this header.",
			HeaderPath: "anthropic-beta",
		},
	},
	Action:          handleBetaOrganizationClaudeCodeManagedSettingsRBACGroupDefaultUpdate,
	HideHelpCommand: true,
}

func handleBetaOrganizationClaudeCodeManagedSettingsRBACGroupDefaultRetrieve(ctx context.Context, cmd *cli.Command) error {
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

	params := anthropic.BetaOrganizationClaudeCodeManagedSettingRBACGroupDefaultGetParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Beta.Organization.ClaudeCode.ManagedSettings.RBACGroupDefault.Get(ctx, params, options...)
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
		Title:          "beta:organization:claude-code:managed-settings:rbac-group-default retrieve",
		Transform:      transform,
	})
}

func handleBetaOrganizationClaudeCodeManagedSettingsRBACGroupDefaultUpdate(ctx context.Context, cmd *cli.Command) error {
	client := anthropic.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()

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

	params := anthropic.BetaOrganizationClaudeCodeManagedSettingRBACGroupDefaultUpdateParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Beta.Organization.ClaudeCode.ManagedSettings.RBACGroupDefault.Update(ctx, params, options...)
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
		Title:          "beta:organization:claude-code:managed-settings:rbac-group-default update",
		Transform:      transform,
	})
}
