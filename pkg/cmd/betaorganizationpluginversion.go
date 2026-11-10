package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/anthropics/anthropic-cli/v2/internal/apiquery"
	"github.com/anthropics/anthropic-cli/v2/internal/requestflag"
	"github.com/anthropics/anthropic-cli/v2/internal/termsafe"
	"github.com/anthropics/anthropic-sdk-go/v2"
	"github.com/anthropics/anthropic-sdk-go/v2/option"
	"github.com/tidwall/gjson"
	"github.com/urfave/cli/v3"
)

var betaOrganizationPluginsVersionsCreate = cli.Command{
	Name:    "create",
	Usage:   "Add a version to an organization-owned Plugin by uploading the new version's\nfiles; it becomes the version served to members unless the Plugin's served\nversion has been pinned.",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "plugin-id",
			Usage:     "ID of the Plugin (prefixed `plugin_`).",
			Required:  true,
			PathParam: "plugin_id",
		},
		&requestflag.Flag[[]string]{
			Name:      "file",
			Usage:     "The version's files: one part per file, the part's filename being the file's path within the Plugin (for example `skills/review-pr/SKILL.md`), or a single `.zip` or `.plugin` archive holding them all. On the wire each part is named `files[]`, and a part named plain `files` is not read; with cURL, `-F 'files[]=@SKILL.md;filename=skills/review-pr/SKILL.md'`. The files must include the manifest, `.claude-plugin/plugin.json`.",
			Required:  true,
			BodyPath:  "files",
			FileInput: true,
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
	Action:          handleBetaOrganizationPluginsVersionsCreate,
	HideHelpCommand: true,
}

var betaOrganizationPluginsVersionsRetrieve = cli.Command{
	Name:    "retrieve",
	Usage:   "Retrieve one version of a Plugin by its ID, or the Plugin's newest version.",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "plugin-id",
			Usage:     "ID of the Plugin (prefixed `plugin_`).",
			Required:  true,
			PathParam: "plugin_id",
		},
		&requestflag.Flag[string]{
			Name:      "version",
			Usage:     "ID of the Plugin Version (prefixed `pluginver_`), or `latest` for the newest one.",
			Required:  true,
			PathParam: "version",
		},
		&requestflag.Flag[string]{
			Name:      "organization-id",
			Usage:     "For a `read:org_audit` or `read:compliance_org_data` key created for all of a parent organization's linked organizations: a child organization of that parent to read instead of the organization the key was created in, given as the organization's UUID or its `org_`-prefixed ID. A value that is neither returns a 400; an organization that is not a child of the key's parent, or where the Plugins API is not available, returns a 404. Any other key may pass only its own organization's ID here; another organization returns a 404. A key with no organization of its own must always pass `organization_id`.",
			QueryPath: "organization_id",
		},
		&requestflag.Flag[[]string]{
			Name:       "beta",
			Usage:      "This endpoint is in beta: requests must send `ce-plugins-2026-09-01` in this header.",
			HeaderPath: "anthropic-beta",
		},
	},
	Action:          handleBetaOrganizationPluginsVersionsRetrieve,
	HideHelpCommand: true,
}

var betaOrganizationPluginsVersionsList = cli.Command{
	Name:    "list",
	Usage:   "List a Plugin's versions, newest first.",
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
			Usage:     "Number of items to return per page.\n\nDefaults to `20`. Ranges from `1` to `1000`.",
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
	Action:          handleBetaOrganizationPluginsVersionsList,
	HideHelpCommand: true,
}

var betaOrganizationPluginsVersionsDownload = cli.Command{
	Name:    "download",
	Usage:   "Download one version's `.zip` archive, exactly as stored. Each download of a\nPlugin from a member's personal plugin marketplace is recorded on the Compliance\nAPI activity feed.",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "plugin-id",
			Usage:     "ID of the Plugin (prefixed `plugin_`).",
			Required:  true,
			PathParam: "plugin_id",
		},
		&requestflag.Flag[string]{
			Name:      "version",
			Usage:     "ID of the Plugin Version (prefixed `pluginver_`). `latest` is not accepted here.",
			Required:  true,
			PathParam: "version",
		},
		&requestflag.Flag[string]{
			Name:      "organization-id",
			Usage:     "For a `read:org_audit` or `read:compliance_org_data` key created for all of a parent organization's linked organizations: a child organization of that parent to read instead of the organization the key was created in, given as the organization's UUID or its `org_`-prefixed ID. A value that is neither returns a 400; an organization that is not a child of the key's parent, or where the Plugins API is not available, returns a 404. Any other key may pass only its own organization's ID here; another organization returns a 404. A key with no organization of its own must always pass `organization_id`.",
			QueryPath: "organization_id",
		},
		&requestflag.Flag[[]string]{
			Name:       "beta",
			Usage:      "This endpoint is in beta: requests must send `ce-plugins-2026-09-01` in this header.",
			HeaderPath: "anthropic-beta",
		},
		&requestflag.Flag[string]{
			Name:    "output",
			Aliases: []string{"o"},
			Usage:   "The file where the response contents will be stored. Use the value '-' to force output to stdout.",
		},
	},
	Action:          handleBetaOrganizationPluginsVersionsDownload,
	HideHelpCommand: true,
}

func handleBetaOrganizationPluginsVersionsCreate(ctx context.Context, cmd *cli.Command) error {
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
		MultipartFormEncoded,
		false,
	)
	if err != nil {
		return err
	}

	params := anthropic.BetaOrganizationPluginVersionNewParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Beta.Organization.Plugins.Versions.New(
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
		Title:          "beta:organization:plugins:versions create",
		Transform:      transform,
	})
}

func handleBetaOrganizationPluginsVersionsRetrieve(ctx context.Context, cmd *cli.Command) error {
	client := anthropic.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("version") && len(unusedArgs) > 0 {
		cmd.Set("version", unusedArgs[0])
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

	params := anthropic.BetaOrganizationPluginVersionGetParams{
		PluginID: cmd.Value("plugin-id").(string),
	}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Beta.Organization.Plugins.Versions.Get(
		ctx,
		cmd.Value("version").(string),
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
		Title:          "beta:organization:plugins:versions retrieve",
		Transform:      transform,
	})
}

func handleBetaOrganizationPluginsVersionsList(ctx context.Context, cmd *cli.Command) error {
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

	params := anthropic.BetaOrganizationPluginVersionListParams{}

	format := "explore"
	explicitFormat := cmd.Root().IsSet("format")
	if explicitFormat {
		format = cmd.Root().String("format")
	}
	transform := cmd.Root().String("transform")
	if format == "raw" {
		var res []byte
		options = append(options, option.WithResponseBodyInto(&res))
		_, err = client.Beta.Organization.Plugins.Versions.List(
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
			Title:          "beta:organization:plugins:versions list",
			Transform:      transform,
		})
	} else {
		iter := client.Beta.Organization.Plugins.Versions.ListAutoPaging(
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
			Title:          "beta:organization:plugins:versions list",
			Transform:      transform,
		})
	}
}

func handleBetaOrganizationPluginsVersionsDownload(ctx context.Context, cmd *cli.Command) error {
	client := anthropic.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("version") && len(unusedArgs) > 0 {
		cmd.Set("version", unusedArgs[0])
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

	params := anthropic.BetaOrganizationPluginVersionDownloadParams{
		PluginID: cmd.Value("plugin-id").(string),
	}

	response, err := client.Beta.Organization.Plugins.Versions.Download(
		ctx,
		cmd.Value("version").(string),
		params,
		options...,
	)
	if err != nil {
		return err
	}
	message, err := writeBinaryResponse(response, os.Stdout, cmd.String("output"))
	if message != "" {
		fmt.Fprintln(termsafe.Stdout, message)
	}
	return err
}
