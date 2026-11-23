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

var betaOrganizationPluginsInstallationSettingsList = cli.Command{
	Name:    "list",
	Usage:   "List an organization-owned Plugin's installation settings, which say which\nmembers it is for, most recently created first.",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "plugin-id",
			Usage:     "ID of the Plugin (prefixed `plugin_`).",
			Required:  true,
			PathParam: "plugin_id",
		},
		&requestflag.Flag[int64]{
			Name:      "limit",
			Usage:     "Number of items to return per page.\n\nDefaults to `20`. Ranges from `1` to `100`.",
			Default:   20,
			QueryPath: "limit",
		},
		&requestflag.Flag[string]{
			Name:      "organization-id",
			Usage:     "For a `read:org_audit` or `read:compliance_org_data` key created for all of a parent organization's linked organizations: a child organization of that parent to read instead of the organization the key was created in, given as the organization's UUID or its `org_`-prefixed ID. A value that is neither returns a 400; an organization that is not a child of the key's parent, or where the Plugins API is not available, returns a 404. Any other key may pass only its own organization's ID here; another organization returns a 404. A key with no organization of its own must always pass `organization_id`.",
			QueryPath: "organization_id",
		},
		&requestflag.Flag[string]{
			Name:      "page",
			Usage:     "Optionally set to the `next_page` token from the previous response.",
			QueryPath: "page",
		},
		&requestflag.Flag[string]{
			Name:      "target-type",
			Usage:     "Only settings for this kind of target: `organization` (the organization-wide setting) or `rbac_group` (an RBAC Group's).\n\nAllowed values: \"organization\", \"rbac_group\".",
			QueryPath: "target_type",
		},
		&requestflag.Flag[[]string]{
			Name:       "beta",
			Usage:      "This endpoint is in beta: requests must send `ce-plugins-2026-09-01` in this header.",
			HeaderPath: "anthropic-beta",
		},
		&requestflag.Flag[int64]{
			Name:  "max-items",
			Usage: "The maximum number of items to return (use -1 for unlimited).",
		},
	},
	Action:          handleBetaOrganizationPluginsInstallationSettingsList,
	HideHelpCommand: true,
}

var betaOrganizationPluginsInstallationSettingsRemove = cli.Command{
	Name:    "remove",
	Usage:   "Remove an organization-owned Plugin's own installation setting for the whole\norganization or for one RBAC Group.",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "plugin-id",
			Usage:     "ID of the Plugin (prefixed `plugin_`).",
			Required:  true,
			PathParam: "plugin_id",
		},
		&requestflag.Flag[string]{
			Name:      "target",
			Usage:     "The target whose own setting is removed: the literal `organization` for the Plugin's organization-wide setting, or an RBAC Group's ID (prefixed `rbac_group_`) for that group's own setting. Removing the `organization` setting returns the Plugin to its marketplace's default.",
			Required:  true,
			PathParam: "target",
		},
		&requestflag.Flag[[]string]{
			Name:       "beta",
			Usage:      "This endpoint is in beta: requests must send `ce-plugins-2026-09-01` in this header.",
			HeaderPath: "anthropic-beta",
		},
	},
	Action:          handleBetaOrganizationPluginsInstallationSettingsRemove,
	HideHelpCommand: true,
}

var betaOrganizationPluginsInstallationSettingsSet = cli.Command{
	Name:    "set",
	Usage:   "Set or change an organization-owned Plugin's installation setting for the whole\norganization or for one RBAC Group.",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "plugin-id",
			Usage:     "ID of the Plugin (prefixed `plugin_`).",
			Required:  true,
			PathParam: "plugin_id",
		},
		&requestflag.Flag[string]{
			Name:      "target",
			Usage:     "The target whose setting is written: the literal `organization` for the Plugin's organization-wide setting, or an RBAC Group's ID (prefixed `rbac_group_`) for that group's own setting. Writing the `organization` target stops the Plugin from inheriting its marketplace's default, even when the value written equals that default.",
			Required:  true,
			PathParam: "target",
		},
		&requestflag.Flag[string]{
			Name:     "installation-preference",
			Usage:    "The installation setting the target is to hold for this Plugin: one of `required`, `auto_install`, `available`, `not_available`.\n\nAllowed values: \"auto_install\", \"available\", \"not_available\", \"required\".",
			Required: true,
			BodyPath: "installation_preference",
		},
		&requestflag.Flag[[]string]{
			Name:       "beta",
			Usage:      "This endpoint is in beta: requests must send `ce-plugins-2026-09-01` in this header.",
			HeaderPath: "anthropic-beta",
		},
	},
	Action:          handleBetaOrganizationPluginsInstallationSettingsSet,
	HideHelpCommand: true,
}

func handleBetaOrganizationPluginsInstallationSettingsList(ctx context.Context, cmd *cli.Command) error {
	client := anthropic.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("plugin-id") && len(unusedArgs) > 0 {
		cmd.Set("plugin-id", unusedArgs[0])
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

	params := anthropic.BetaOrganizationPluginInstallationSettingListParams{}

	format := "explore"
	explicitFormat := cmd.Root().IsSet("format")
	if explicitFormat {
		format = cmd.Root().String("format")
	}
	transform := cmd.Root().String("transform")
	if format == "raw" {
		var res []byte
		options = append(options, option.WithResponseBodyInto(&res))
		_, err = client.Beta.Organization.Plugins.InstallationSettings.List(
			ctx,
			cmd.Value("plugin-id").(string),
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
			Title:          "beta:organization:plugins:installation-settings list",
			Transform:      transform,
		})
	} else {
		iter := client.Beta.Organization.Plugins.InstallationSettings.ListAutoPaging(
			ctx,
			cmd.Value("plugin-id").(string),
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
			Title:          "beta:organization:plugins:installation-settings list",
			Transform:      transform,
		})
	}
}

func handleBetaOrganizationPluginsInstallationSettingsRemove(ctx context.Context, cmd *cli.Command) error {
	client := anthropic.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("target") && len(unusedArgs) > 0 {
		cmd.Set("target", unusedArgs[0])
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

	params := anthropic.BetaOrganizationPluginInstallationSettingRemoveParams{
		PluginID: cmd.Value("plugin-id").(string),
	}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Beta.Organization.Plugins.InstallationSettings.Remove(
		ctx,
		cmd.Value("target").(string),
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
		Title:          "beta:organization:plugins:installation-settings remove",
		Transform:      transform,
	})
}

func handleBetaOrganizationPluginsInstallationSettingsSet(ctx context.Context, cmd *cli.Command) error {
	client := anthropic.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("target") && len(unusedArgs) > 0 {
		cmd.Set("target", unusedArgs[0])
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

	params := anthropic.BetaOrganizationPluginInstallationSettingSetParams{
		PluginID: cmd.Value("plugin-id").(string),
	}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Beta.Organization.Plugins.InstallationSettings.Set(
		ctx,
		cmd.Value("target").(string),
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
		Title:          "beta:organization:plugins:installation-settings set",
		Transform:      transform,
	})
}
