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

var betaOrganizationPluginsCreate = cli.Command{
	Name:    "create",
	Usage:   "Create an organization-owned Plugin and its first version by uploading the\nversion's files.",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[[]string]{
			Name:      "file",
			Usage:     "The version's files: one part per file, the part's filename being the file's path within the Plugin (for example `skills/review-pr/SKILL.md`), or a single `.zip` or `.plugin` archive holding them all. On the wire each part is named `files[]`, and a part named plain `files` is not read; with cURL, `-F 'files[]=@SKILL.md;filename=skills/review-pr/SKILL.md'`. The files must include the manifest, `.claude-plugin/plugin.json`.",
			Required:  true,
			BodyPath:  "files",
			FileInput: true,
		},
		&requestflag.Flag[string]{
			Name:     "marketplace-id",
			Usage:    "ID of the organization-owned plugin marketplace to create the Plugin in (prefixed `marketplace_`). It must be a `manual` marketplace, one whose Plugins are uploaded rather than synchronized from a repository. When omitted, the Plugin is created in the organization's library marketplace, an organization-owned `manual` marketplace created on first use.",
			BodyPath: "marketplace_id",
		},
		&requestflag.Flag[string]{
			Name:     "release-notes",
			Usage:    "Release notes stored with the version and shown in its version history in claude.ai; up to 5,000 characters.",
			BodyPath: "release_notes",
		},
		&requestflag.Flag[[]string]{
			Name:       "beta",
			Usage:      "This endpoint is in beta: requests must send `ce-plugins-2026-09-01` in this header.",
			HeaderPath: "anthropic-beta",
		},
	},
	Action:          handleBetaOrganizationPluginsCreate,
	HideHelpCommand: true,
}

var betaOrganizationPluginsRetrieve = cli.Command{
	Name:    "retrieve",
	Usage:   "Retrieve a Plugin by ID.",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "plugin-id",
			Usage:     "ID of the Plugin (prefixed `plugin_`).",
			Required:  true,
			PathParam: "plugin_id",
		},
		&requestflag.Flag[string]{
			Name:      "organization-id",
			Usage:     "For a `read:org_audit` or `read:compliance_org_data` key created for all of a parent organization's linked organizations: a child organization of that parent to read instead of the organization the key was created in, given as the organization's UUID or its `org_`-prefixed ID. A value that is neither returns a 400; an organization that is not a child of the key's parent, or where the Plugins API is not available, returns a 404. Any other key may pass only its own organization's ID here; another organization returns a 404.",
			QueryPath: "organization_id",
		},
		&requestflag.Flag[[]string]{
			Name:       "beta",
			Usage:      "This endpoint is in beta: requests must send `ce-plugins-2026-09-01` in this header.",
			HeaderPath: "anthropic-beta",
		},
	},
	Action:          handleBetaOrganizationPluginsRetrieve,
	HideHelpCommand: true,
}

var betaOrganizationPluginsUpdate = cli.Command{
	Name:    "update",
	Usage:   "Change which stored version of an organization-owned Plugin is served to\nmembers, for example to roll back to an earlier one. This pins the served\nversion: later uploads are stored but no longer change what is served, and\npinning cannot currently be undone, here or in claude.ai.",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "plugin-id",
			Usage:     "ID of the Plugin (prefixed `plugin_`).",
			Required:  true,
			PathParam: "plugin_id",
		},
		&requestflag.Flag[string]{
			Name:     "served-version-id",
			Usage:    "Serve this version of the Plugin (prefixed `pluginver_`) and pin the served version to it; `latest` is not accepted.",
			Required: true,
			BodyPath: "served_version_id",
		},
		&requestflag.Flag[[]string]{
			Name:       "beta",
			Usage:      "This endpoint is in beta: requests must send `ce-plugins-2026-09-01` in this header.",
			HeaderPath: "anthropic-beta",
		},
	},
	Action:          handleBetaOrganizationPluginsUpdate,
	HideHelpCommand: true,
}

var betaOrganizationPluginsList = cli.Command{
	Name:    "list",
	Usage:   "List the Plugins created under the organization, newest first: those in the\norganization's own plugin marketplaces and those in members' personal plugin\nmarketplaces.",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[any]{
			Name:      "created-at-gt",
			Usage:     "RFC 3339 timestamp bound; combine [gte], [gt], [lte], [lt].",
			QueryPath: "created_at[gt]",
		},
		&requestflag.Flag[any]{
			Name:      "created-at-gte",
			Usage:     "RFC 3339 timestamp bound; combine [gte], [gt], [lte], [lt].",
			QueryPath: "created_at[gte]",
		},
		&requestflag.Flag[any]{
			Name:      "created-at-lt",
			Usage:     "RFC 3339 timestamp bound; combine [gte], [gt], [lte], [lt].",
			QueryPath: "created_at[lt]",
		},
		&requestflag.Flag[any]{
			Name:      "created-at-lte",
			Usage:     "RFC 3339 timestamp bound; combine [gte], [gt], [lte], [lt].",
			QueryPath: "created_at[lte]",
		},
		&requestflag.Flag[int64]{
			Name:      "limit",
			Usage:     "Number of items to return per page.\n\nDefaults to `20`. Ranges from `1` to `100`.",
			Default:   20,
			QueryPath: "limit",
		},
		&requestflag.Flag[string]{
			Name:      "marketplace-id",
			Usage:     "Only Plugins in this plugin marketplace (prefixed `marketplace_`).",
			QueryPath: "marketplace_id",
		},
		&requestflag.Flag[string]{
			Name:      "organization-id",
			Usage:     "For a `read:org_audit` or `read:compliance_org_data` key created for all of a parent organization's linked organizations: a child organization of that parent to read instead of the organization the key was created in, given as the organization's UUID or its `org_`-prefixed ID. A value that is neither returns a 400; an organization that is not a child of the key's parent, or where the Plugins API is not available, returns a 404. Any other key may pass only its own organization's ID here; another organization returns a 404.",
			QueryPath: "organization_id",
		},
		&requestflag.Flag[string]{
			Name:      "owner-type",
			Usage:     "`organization` for Plugins in the organization's plugin marketplaces, `user` for Plugins in members' personal plugin marketplaces.",
			QueryPath: "owner_type",
		},
		&requestflag.Flag[string]{
			Name:      "owner-user-id",
			Usage:     "Only Plugins in this member's personal plugin marketplaces (prefixed `user_`); a removed member's ID is accepted.",
			QueryPath: "owner_user_id",
		},
		&requestflag.Flag[string]{
			Name:      "page",
			Usage:     "Optionally set to the `next_page` token from the previous response.",
			QueryPath: "page",
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
	Action:          handleBetaOrganizationPluginsList,
	HideHelpCommand: true,
}

var betaOrganizationPluginsDelete = cli.Command{
	Name:    "delete",
	Usage:   "Permanently delete a Plugin and every version it holds, exactly as when an\nadministrator deletes it in claude.ai. The Plugin may belong to the organization\nor to a member, including a member who has since left the organization. Do not\ndelete a Plugin only because your pipeline did not create it: a member or an\nadministrator might have added it in claude.ai.",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "plugin-id",
			Usage:     "ID of the Plugin (prefixed `plugin_`).",
			Required:  true,
			PathParam: "plugin_id",
		},
		&requestflag.Flag[[]string]{
			Name:       "beta",
			Usage:      "This endpoint is in beta: requests must send `ce-plugins-2026-09-01` in this header.",
			HeaderPath: "anthropic-beta",
		},
	},
	Action:          handleBetaOrganizationPluginsDelete,
	HideHelpCommand: true,
}

func handleBetaOrganizationPluginsCreate(ctx context.Context, cmd *cli.Command) error {
	client := anthropic.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()

	if len(unusedArgs) > 0 {
		return fmt.Errorf("Unexpected extra arguments: %v", unusedArgs)
	}

	options, err := flagOptions(
		cmd,
		apiquery.NestedQueryFormatBrackets,
		apiquery.ArrayQueryFormatBrackets,
		MultipartFormEncoded,
		false,
	)
	if err != nil {
		return err
	}

	params := anthropic.BetaOrganizationPluginNewParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Beta.Organization.Plugins.New(ctx, params, options...)
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
		Title:          "beta:organization:plugins create",
		Transform:      transform,
	})
}

func handleBetaOrganizationPluginsRetrieve(ctx context.Context, cmd *cli.Command) error {
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

	params := anthropic.BetaOrganizationPluginGetParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Beta.Organization.Plugins.Get(
		ctx,
		cmd.Value("plugin-id").(string),
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
		Title:          "beta:organization:plugins retrieve",
		Transform:      transform,
	})
}

func handleBetaOrganizationPluginsUpdate(ctx context.Context, cmd *cli.Command) error {
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
		ApplicationJSON,
		false,
	)
	if err != nil {
		return err
	}

	params := anthropic.BetaOrganizationPluginUpdateParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Beta.Organization.Plugins.Update(
		ctx,
		cmd.Value("plugin-id").(string),
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
		Title:          "beta:organization:plugins update",
		Transform:      transform,
	})
}

func handleBetaOrganizationPluginsList(ctx context.Context, cmd *cli.Command) error {
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

	params := anthropic.BetaOrganizationPluginListParams{}

	format := "explore"
	explicitFormat := cmd.Root().IsSet("format")
	if explicitFormat {
		format = cmd.Root().String("format")
	}
	transform := cmd.Root().String("transform")
	if format == "raw" {
		var res []byte
		options = append(options, option.WithResponseBodyInto(&res))
		_, err = client.Beta.Organization.Plugins.List(ctx, params, options...)
		if err != nil {
			return err
		}
		obj := gjson.ParseBytes(res)
		return ShowJSON(obj, ShowJSONOpts{
			ExplicitFormat: explicitFormat,
			Format:         format,
			RawOutput:      cmd.Root().Bool("raw-output"),
			Title:          "beta:organization:plugins list",
			Transform:      transform,
		})
	} else {
		iter := client.Beta.Organization.Plugins.ListAutoPaging(ctx, params, options...)
		maxItems := int64(-1)
		if cmd.IsSet("max-items") {
			maxItems = cmd.Value("max-items").(int64)
		}
		return ShowJSONIterator(iter, maxItems, ShowJSONOpts{
			ExplicitFormat: explicitFormat,
			Format:         format,
			RawOutput:      cmd.Root().Bool("raw-output"),
			Title:          "beta:organization:plugins list",
			Transform:      transform,
		})
	}
}

func handleBetaOrganizationPluginsDelete(ctx context.Context, cmd *cli.Command) error {
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

	params := anthropic.BetaOrganizationPluginDeleteParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Beta.Organization.Plugins.Delete(
		ctx,
		cmd.Value("plugin-id").(string),
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
		Title:          "beta:organization:plugins delete",
		Transform:      transform,
	})
}
