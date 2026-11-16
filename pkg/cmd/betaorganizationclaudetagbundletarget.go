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

var betaOrganizationClaudeTagBundlesTargetsCreate = requestflag.WithInnerFlags(cli.Command{
	Name:    "create",
	Usage:   "Attach a Bundle to a Slack channel or a Slack workspace. A Bundle attached to a\nworkspace can be used in the workspace and in every channel in it. Attaching a\nBundle to a channel or workspace it is already attached to returns the existing\nBundle Target. The channel or workspace must be set up: attaching never sets a\nchannel up. A group direct message takes no Bundle and returns 400.",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "bundle-id",
			Usage:     "The Bundle's identifier, starting with `capp_`.",
			Required:  true,
			PathParam: "bundle_id",
		},
		&requestflag.Flag[map[string]any]{
			Name:     "target",
			Required: true,
			BodyRoot: true,
		},
		&requestflag.Flag[[]string]{
			Name:       "beta",
			Usage:      "This endpoint is in beta: requests must send `experimental-claude-tag-2026-10-21` in this header.",
			HeaderPath: "anthropic-beta",
		},
	},
	Action:          handleBetaOrganizationClaudeTagBundlesTargetsCreate,
	HideHelpCommand: true,
}, map[string][]requestflag.HasOuterFlag{
	"target": {
		&requestflag.InnerFlag[string]{
			Name:       "target.type",
			Usage:      `Allowed values: "slack_channel", "slack_workspace".`,
			InnerField: "type",
		},
		&requestflag.InnerFlag[string]{
			Name:       "target.slack-channel-id",
			Usage:      "The identifier of the Slack channel to attach the Bundle to, starting with `slch_`. It is the channel object's `id`, not Slack's `channel_id`.",
			InnerField: "slack_channel_id",
		},
		&requestflag.InnerFlag[string]{
			Name:       "target.slack-workspace-id",
			Usage:      "The identifier of the Slack workspace to attach the Bundle to, starting with `slws_`.",
			InnerField: "slack_workspace_id",
		},
	},
})

var betaOrganizationClaudeTagBundlesTargetsRetrieve = cli.Command{
	Name:    "retrieve",
	Usage:   "Get one Bundle Target. Returns 404 when `target_id` is the identifier of another\nBundle's Bundle Target, or the Bundle Target was detached.",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "bundle-id",
			Usage:     "The Bundle's identifier, starting with `capp_`.",
			Required:  true,
			PathParam: "bundle_id",
		},
		&requestflag.Flag[string]{
			Name:        "target-id",
			Usage:       "The Bundle Target's identifier, starting with `btgt_`. Attach Bundle and List Bundle Targets return it. It is never the identifier of the channel or workspace.",
			Required:    true,
			PathParam:   "target_id",
			DataAliases: []string{"id"},
		},
		&requestflag.Flag[[]string]{
			Name:       "beta",
			Usage:      "This endpoint is in beta: requests must send `experimental-claude-tag-2026-10-21` in this header.",
			HeaderPath: "anthropic-beta",
		},
	},
	Action:          handleBetaOrganizationClaudeTagBundlesTargetsRetrieve,
	HideHelpCommand: true,
}

var betaOrganizationClaudeTagBundlesTargetsList = cli.Command{
	Name:    "list",
	Usage:   "List the Slack channels and Slack workspaces a Bundle is attached to, newest\nfirst by `created_at`. Only channels and workspaces the organization owns are\nlisted; a Bundle that is attached to anything else has no Bundle Target for it.",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "bundle-id",
			Usage:     "The Bundle's identifier, starting with `capp_`.",
			Required:  true,
			PathParam: "bundle_id",
		},
		&requestflag.Flag[int64]{
			Name:      "limit",
			Usage:     "Number of items to return per page, between 1 and 100. Defaults to 20.",
			QueryPath: "limit",
		},
		&requestflag.Flag[string]{
			Name:      "page",
			Usage:     "Opaque cursor from a previous response's `next_page`.",
			QueryPath: "page",
		},
		&requestflag.Flag[[]string]{
			Name:       "beta",
			Usage:      "This endpoint is in beta: requests must send `experimental-claude-tag-2026-10-21` in this header.",
			HeaderPath: "anthropic-beta",
		},
		&requestflag.Flag[int64]{
			Name:  "max-items",
			Usage: "The maximum number of items to return (use -1 for unlimited).",
		},
	},
	Action:          handleBetaOrganizationClaudeTagBundlesTargetsList,
	HideHelpCommand: true,
}

var betaOrganizationClaudeTagBundlesTargetsDelete = cli.Command{
	Name:    "delete",
	Usage:   "Detach a Bundle from a Slack channel or a Slack workspace. A repeated detach\nreturns 404.",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "bundle-id",
			Usage:     "The Bundle's identifier, starting with `capp_`.",
			Required:  true,
			PathParam: "bundle_id",
		},
		&requestflag.Flag[string]{
			Name:        "target-id",
			Usage:       "The Bundle Target's identifier, starting with `btgt_`. Attach Bundle and List Bundle Targets return it. It is never the identifier of the channel or workspace.",
			Required:    true,
			PathParam:   "target_id",
			DataAliases: []string{"id"},
		},
		&requestflag.Flag[[]string]{
			Name:       "beta",
			Usage:      "This endpoint is in beta: requests must send `experimental-claude-tag-2026-10-21` in this header.",
			HeaderPath: "anthropic-beta",
		},
	},
	Action:          handleBetaOrganizationClaudeTagBundlesTargetsDelete,
	HideHelpCommand: true,
}

func handleBetaOrganizationClaudeTagBundlesTargetsCreate(ctx context.Context, cmd *cli.Command) error {
	client := anthropic.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("bundle-id") && len(unusedArgs) > 0 {
		cmd.Set("bundle-id", unusedArgs[0])
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

	params := anthropic.BetaOrganizationClaudeTagBundleTargetNewParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Beta.Organization.ClaudeTag.Bundles.Targets.New(
		ctx,
		cmd.Value("bundle-id").(string),
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
		Title:          "beta:organization:claude-tag:bundles:targets create",
		Transform:      transform,
	})
}

func handleBetaOrganizationClaudeTagBundlesTargetsRetrieve(ctx context.Context, cmd *cli.Command) error {
	client := anthropic.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("target-id") && len(unusedArgs) > 0 {
		cmd.Set("target-id", unusedArgs[0])
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

	params := anthropic.BetaOrganizationClaudeTagBundleTargetGetParams{
		BundleID: cmd.Value("bundle-id").(string),
	}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Beta.Organization.ClaudeTag.Bundles.Targets.Get(
		ctx,
		cmd.Value("target-id").(string),
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
		Title:          "beta:organization:claude-tag:bundles:targets retrieve",
		Transform:      transform,
	})
}

func handleBetaOrganizationClaudeTagBundlesTargetsList(ctx context.Context, cmd *cli.Command) error {
	client := anthropic.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("bundle-id") && len(unusedArgs) > 0 {
		cmd.Set("bundle-id", unusedArgs[0])
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

	params := anthropic.BetaOrganizationClaudeTagBundleTargetListParams{}

	format := "explore"
	explicitFormat := cmd.Root().IsSet("format")
	if explicitFormat {
		format = cmd.Root().String("format")
	}
	transform := cmd.Root().String("transform")
	if format == "raw" {
		var res []byte
		options = append(options, option.WithResponseBodyInto(&res))
		_, err = client.Beta.Organization.ClaudeTag.Bundles.Targets.List(
			ctx,
			cmd.Value("bundle-id").(string),
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
			Title:          "beta:organization:claude-tag:bundles:targets list",
			Transform:      transform,
		})
	} else {
		iter := client.Beta.Organization.ClaudeTag.Bundles.Targets.ListAutoPaging(
			ctx,
			cmd.Value("bundle-id").(string),
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
			Title:          "beta:organization:claude-tag:bundles:targets list",
			Transform:      transform,
		})
	}
}

func handleBetaOrganizationClaudeTagBundlesTargetsDelete(ctx context.Context, cmd *cli.Command) error {
	client := anthropic.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("target-id") && len(unusedArgs) > 0 {
		cmd.Set("target-id", unusedArgs[0])
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

	params := anthropic.BetaOrganizationClaudeTagBundleTargetDeleteParams{
		BundleID: cmd.Value("bundle-id").(string),
	}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Beta.Organization.ClaudeTag.Bundles.Targets.Delete(
		ctx,
		cmd.Value("target-id").(string),
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
		Title:          "beta:organization:claude-tag:bundles:targets delete",
		Transform:      transform,
	})
}
