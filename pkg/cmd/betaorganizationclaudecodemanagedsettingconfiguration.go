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

var betaOrganizationClaudeCodeManagedSettingsConfigurationRetrieve = cli.Command{
	Name:    "retrieve",
	Usage:   "Retrieve the organization's Claude Code managed settings mode and RBAC Group\norder.",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[[]string]{
			Name:       "beta",
			Usage:      "This endpoint is in beta: requests must send `ce-claude-code-managed-settings-2026-10-06` in this header.",
			HeaderPath: "anthropic-beta",
		},
	},
	Action:          handleBetaOrganizationClaudeCodeManagedSettingsConfigurationRetrieve,
	HideHelpCommand: true,
}

var betaOrganizationClaudeCodeManagedSettingsConfigurationUpdate = requestflag.WithInnerFlags(cli.Command{
	Name:    "update",
	Usage:   "Change the organization's Claude Code managed settings mode, the RBAC Group\norder, or both, in one atomic write.",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[map[string]any]{
			Name:     "mode",
			BodyPath: "mode",
		},
		&requestflag.Flag[[]string]{
			Name:     "rbac-group-order",
			Usage:    "Every RBAC Group that has settings, each listed once, in the new order, highest precedence first. Requires `rbac_group_order_version`.",
			BodyPath: "rbac_group_order",
		},
		&requestflag.Flag[string]{
			Name:     "rbac-group-order-version",
			Usage:    "The `rbac_group_order_version` from the read this change is based on. Requires `rbac_group_order`.",
			BodyPath: "rbac_group_order_version",
		},
		&requestflag.Flag[[]string]{
			Name:       "beta",
			Usage:      "This endpoint is in beta: requests must send `ce-claude-code-managed-settings-2026-10-06` in this header.",
			HeaderPath: "anthropic-beta",
		},
	},
	Action:          handleBetaOrganizationClaudeCodeManagedSettingsConfigurationUpdate,
	HideHelpCommand: true,
}, map[string][]requestflag.HasOuterFlag{
	"mode": {
		&requestflag.InnerFlag[string]{
			Name:       "mode.type",
			Usage:      `Allowed values: "organization", "rbac_groups".`,
			InnerField: "type",
		},
	},
})

func handleBetaOrganizationClaudeCodeManagedSettingsConfigurationRetrieve(ctx context.Context, cmd *cli.Command) error {
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

	params := anthropic.BetaOrganizationClaudeCodeManagedSettingConfigurationGetParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Beta.Organization.ClaudeCode.ManagedSettings.Configuration.Get(ctx, params, options...)
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
		Title:          "beta:organization:claude-code:managed-settings:configuration retrieve",
		Transform:      transform,
	})
}

func handleBetaOrganizationClaudeCodeManagedSettingsConfigurationUpdate(ctx context.Context, cmd *cli.Command) error {
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

	params := anthropic.BetaOrganizationClaudeCodeManagedSettingConfigurationUpdateParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Beta.Organization.ClaudeCode.ManagedSettings.Configuration.Update(ctx, params, options...)
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
		Title:          "beta:organization:claude-code:managed-settings:configuration update",
		Transform:      transform,
	})
}
