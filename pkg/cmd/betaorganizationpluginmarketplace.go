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

var betaOrganizationPluginMarketplacesRetrieve = cli.Command{
	Name:    "retrieve",
	Usage:   "Retrieve a plugin marketplace by ID.",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "marketplace-id",
			Usage:     "ID of the plugin marketplace (prefixed `marketplace_`).",
			Required:  true,
			PathParam: "marketplace_id",
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
	Action:          handleBetaOrganizationPluginMarketplacesRetrieve,
	HideHelpCommand: true,
}

var betaOrganizationPluginMarketplacesUpdate = cli.Command{
	Name:    "update",
	Usage:   "Set the default installation setting of one of the organization's own plugin\nmarketplaces. Every Plugin in it without a setting of its own gets this default\nas its organization-wide setting, including Plugins added later.",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "marketplace-id",
			Usage:     "ID of the plugin marketplace (prefixed `marketplace_`).",
			Required:  true,
			PathParam: "marketplace_id",
		},
		&requestflag.Flag[string]{
			Name:     "default-installation-preference",
			Usage:    "The organization-wide installation setting every Plugin in the marketplace without one of its own gets: one of `required`, `auto_install`, `available`, `not_available`. Once set it can be changed but not removed.\n\nAllowed values: \"auto_install\", \"available\", \"not_available\", \"required\".",
			Required: true,
			BodyPath: "default_installation_preference",
		},
		&requestflag.Flag[[]string]{
			Name:       "beta",
			Usage:      "This endpoint is in beta: requests must send `ce-plugins-2026-09-01` in this header.",
			HeaderPath: "anthropic-beta",
		},
	},
	Action:          handleBetaOrganizationPluginMarketplacesUpdate,
	HideHelpCommand: true,
}

var betaOrganizationPluginMarketplacesList = cli.Command{
	Name:    "list",
	Usage:   "List the plugin marketplaces Plugins live in, newest first: the organization's\nown and its members' personal ones.",
	Suggest: true,
	Flags: []cli.Flag{
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
			Name:      "owner-type",
			Usage:     "`organization` for the organization's plugin marketplaces, `user` for members' personal plugin marketplaces.\n\nAllowed values: \"organization\", \"user\".",
			QueryPath: "owner_type",
		},
		&requestflag.Flag[string]{
			Name:      "page",
			Usage:     "Optionally set to the `next_page` token from the previous response.",
			QueryPath: "page",
		},
		&requestflag.Flag[string]{
			Name:      "source",
			Usage:     "Only plugin marketplaces with this `source`: `manual` for those whose Plugins are uploaded; `github`, `gitlab` or `public_git` for those synchronized from a Git repository. `directory` (Anthropic's catalog) is never listed here.\n\nAllowed values: \"directory\", \"github\", \"gitlab\", \"manual\", \"public_git\".",
			QueryPath: "source",
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
	Action:          handleBetaOrganizationPluginMarketplacesList,
	HideHelpCommand: true,
}

var betaOrganizationPluginMarketplacesValidateArchive = cli.Command{
	Name:    "validate-archive",
	Usage:   "Check whether a plugin marketplace, uploaded as a `.zip` of the marketplace\ndirectory, would synchronize into claude.ai, without connecting or storing it.",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "archive",
			Usage:     "A .zip of the marketplace directory (its contents at the root, or wrapped in one folder as a Git host's download produces), sent as a file part with a filename; DEFLATE- or STORE-compressed, at most 32 MB. A part sent without a filename, a second archive part, or any other form field is a 400; a larger archive is a 413.",
			Required:  true,
			BodyPath:  "archive",
			FileInput: true,
		},
		&requestflag.Flag[[]string]{
			Name:       "beta",
			Usage:      "This endpoint is in beta: requests must send `ce-plugins-2026-09-01` in this header.",
			HeaderPath: "anthropic-beta",
		},
	},
	Action:          handleBetaOrganizationPluginMarketplacesValidateArchive,
	HideHelpCommand: true,
}

var betaOrganizationPluginMarketplacesValidateRepository = cli.Command{
	Name:    "validate-repository",
	Usage:   "Check whether a plugin marketplace held in a public GitHub repository would\nsynchronize into claude.ai, without connecting or storing it.",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:     "repository-url",
			Usage:    "The `https://` URL of a public repository on github.com that holds the marketplace. A URL longer than 2,048 characters, on any other host or scheme, with a port other than 443 or 80, with credentials in it, or that does not name exactly `owner/repo` (a trailing `.git` or `/` is accepted) is a 400. A private or missing repository is reported in the response, not refused.",
			Required: true,
			BodyPath: "repository_url",
		},
		&requestflag.Flag[*string]{
			Name:     "ref",
			Usage:    "The branch to validate the tip of, or the full 40-character SHA of the commit to validate. When omitted, the branch a synchronization would read (usually the repository's default branch); if that is not the default branch, the report's `ref` says which branch was read. An empty string, a value longer than 255 characters, or one that contains `..` or a character other than ASCII letters, digits, `.`, `_`, `-`, `+`, and `/`, is a 400. Any other value that is not a 40-character SHA is read as a branch name; if the repository has no such branch, the response is a report with `valid: false`, not a 400.",
			BodyPath: "ref",
		},
		&requestflag.Flag[[]string]{
			Name:       "beta",
			Usage:      "This endpoint is in beta: requests must send `ce-plugins-2026-09-01` in this header.",
			HeaderPath: "anthropic-beta",
		},
	},
	Action:          handleBetaOrganizationPluginMarketplacesValidateRepository,
	HideHelpCommand: true,
}

func handleBetaOrganizationPluginMarketplacesRetrieve(ctx context.Context, cmd *cli.Command) error {
	client := anthropic.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("marketplace-id") && len(unusedArgs) > 0 {
		cmd.Set("marketplace-id", unusedArgs[0])
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

	params := anthropic.BetaOrganizationPluginMarketplaceGetParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Beta.Organization.PluginMarketplaces.Get(
		ctx,
		cmd.Value("marketplace-id").(string),
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
		Title:          "beta:organization:plugin-marketplaces retrieve",
		Transform:      transform,
	})
}

func handleBetaOrganizationPluginMarketplacesUpdate(ctx context.Context, cmd *cli.Command) error {
	client := anthropic.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("marketplace-id") && len(unusedArgs) > 0 {
		cmd.Set("marketplace-id", unusedArgs[0])
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

	params := anthropic.BetaOrganizationPluginMarketplaceUpdateParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Beta.Organization.PluginMarketplaces.Update(
		ctx,
		cmd.Value("marketplace-id").(string),
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
		Title:          "beta:organization:plugin-marketplaces update",
		Transform:      transform,
	})
}

func handleBetaOrganizationPluginMarketplacesList(ctx context.Context, cmd *cli.Command) error {
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

	params := anthropic.BetaOrganizationPluginMarketplaceListParams{}

	format := "explore"
	explicitFormat := cmd.Root().IsSet("format")
	if explicitFormat {
		format = cmd.Root().String("format")
	}
	transform := cmd.Root().String("transform")
	if format == "raw" {
		var res []byte
		options = append(options, option.WithResponseBodyInto(&res))
		_, err = client.Beta.Organization.PluginMarketplaces.List(ctx, params, options...)
		if err != nil {
			return err
		}
		obj := gjson.ParseBytes(res)
		return ShowJSON(obj, ShowJSONOpts{
			ExplicitFormat: explicitFormat,
			Format:         format,
			RawOutput:      cmd.Root().Bool("raw-output"),
			Title:          "beta:organization:plugin-marketplaces list",
			Transform:      transform,
		})
	} else {
		iter := client.Beta.Organization.PluginMarketplaces.ListAutoPaging(ctx, params, options...)
		maxItems := int64(-1)
		if cmd.IsSet("max-items") {
			maxItems = cmd.Value("max-items").(int64)
		}
		return ShowJSONIterator(iter, maxItems, ShowJSONOpts{
			ExplicitFormat: explicitFormat,
			Format:         format,
			RawOutput:      cmd.Root().Bool("raw-output"),
			Title:          "beta:organization:plugin-marketplaces list",
			Transform:      transform,
		})
	}
}

func handleBetaOrganizationPluginMarketplacesValidateArchive(ctx context.Context, cmd *cli.Command) error {
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

	params := anthropic.BetaOrganizationPluginMarketplaceValidateArchiveParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Beta.Organization.PluginMarketplaces.ValidateArchive(ctx, params, options...)
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
		Title:          "beta:organization:plugin-marketplaces validate-archive",
		Transform:      transform,
	})
}

func handleBetaOrganizationPluginMarketplacesValidateRepository(ctx context.Context, cmd *cli.Command) error {
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

	params := anthropic.BetaOrganizationPluginMarketplaceValidateRepositoryParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Beta.Organization.PluginMarketplaces.ValidateRepository(ctx, params, options...)
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
		Title:          "beta:organization:plugin-marketplaces validate-repository",
		Transform:      transform,
	})
}
