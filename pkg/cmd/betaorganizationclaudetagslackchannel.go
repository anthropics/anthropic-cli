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

var betaOrganizationClaudeTagSlackChannelsCreate = requestflag.WithInnerFlags(cli.Command{
	Name:    "create",
	Usage:   "Set a Slack channel up for Claude Tag and return it with its identifier. If the\nchannel is already set up, whether in claude.ai, by Claude's first use in it or\nby an earlier request, the existing channel is returned, so the request can be\nrepeated. `settings`, when sent, is applied as an update applies it. Before a\nchannel is set up, Claude's Slack app confirms with Slack that the channel is in\nthe workspace, is not archived and is visible to the app; a channel that fails\nthe confirmation returns 404. Setting the channel up and applying `settings` are\ntwo steps: when the second is refused, the response is the refusal, the channel\nstays set up with every setting unset, and the error's\n`details.slack_channel_id` holds the channel's identifier.",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[map[string]any]{
			Name:     "slack",
			Usage:    "The identifiers Slack shows for a channel. They are Slack's own, not identifiers of this API.",
			Required: true,
			BodyPath: "slack",
		},
		&requestflag.Flag[map[string]any]{
			Name:     "settings",
			Usage:    "The settings in a request that sets a Slack channel up. A setting left out is not set. A setting sent as null is cleared.",
			BodyPath: "settings",
		},
		&requestflag.Flag[[]string]{
			Name:       "beta",
			Usage:      "This endpoint is in beta: requests must send `experimental-claude-tag-2026-10-21` in this header.",
			HeaderPath: "anthropic-beta",
		},
	},
	Action:          handleBetaOrganizationClaudeTagSlackChannelsCreate,
	HideHelpCommand: true,
}, map[string][]requestflag.HasOuterFlag{
	"slack": {
		&requestflag.InnerFlag[string]{
			Name:       "slack.channel-id",
			Usage:      "Slack's identifier for the channel, starting with `C` or `G`, such as `C0456EFGH`.",
			InnerField: "channel_id",
		},
		&requestflag.InnerFlag[string]{
			Name:       "slack.workspace-id",
			Usage:      "Slack's identifier for the workspace the channel is in, starting with `T`, such as `T0123ABCD`.",
			InnerField: "workspace_id",
		},
	},
	"settings": {
		&requestflag.InnerFlag[*bool]{
			Name:       "settings.enabled",
			Usage:      "Whether Claude is on. `true` on a channel turns Claude on in it even when its workspace stores `false`. Null clears it.",
			InnerField: "enabled",
		},
		&requestflag.InnerFlag[*string]{
			Name:       "settings.environment-id",
			Usage:      "The environment Claude's sessions run in: the identifier of an active environment of the organization, starting with `env_`, or of a self-hosted runner pool of the organization that every member may use, starting with `ccpool_`. Null clears it.",
			InnerField: "environment_id",
		},
		&requestflag.InnerFlag[map[string]any]{
			Name:       "settings.guest-policy",
			Usage:      "What Claude does in a channel while a Slack guest is in it, for everyone in the channel. It changes nothing in a channel without Slack guests. More types can be added.",
			InnerField: "guest_policy",
		},
		&requestflag.InnerFlag[*string]{
			Name:       "settings.instructions",
			Usage:      "Instructions for Claude, up to 8,192 characters after leading and trailing white space is trimmed. Null or blank text clears them.",
			InnerField: "instructions",
		},
		&requestflag.InnerFlag[*string]{
			Name:       "settings.model",
			Usage:      "The model Claude uses: the identifier of a model, or the name of a model line, that the organization is entitled to for Claude Tag. Null clears it.",
			InnerField: "model",
		},
		&requestflag.InnerFlag[map[string]any]{
			Name:       "settings.reply-mode",
			Usage:      "How Claude takes part in a Slack channel. A channel alone stores one. More types can be added.",
			InnerField: "reply_mode",
		},
	},
})

var betaOrganizationClaudeTagSlackChannelsRetrieve = cli.Command{
	Name:    "retrieve",
	Usage:   "Get one Slack channel and its settings. By default `settings` holds what the\nchannel itself stores: a setting that is null there is not set on the channel\nand is inherited. With `include_inherited=true`, `settings` holds what applies\nin the channel once the settings levels above it are taken into account, and the\nrequest fails with 503 rather than resolve against a settings level that could\nnot be read. Either way `sources` names the settings level each value comes\nfrom, so a caller tells the channel's own value from an inherited one by its\nsource, not by whether it is null.",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "id",
			Usage:     "The channel's identifier, starting with `slch_`.",
			Required:  true,
			PathParam: "id",
		},
		&requestflag.Flag[bool]{
			Name:      "include-inherited",
			Usage:     "When `true`, `settings` holds what applies in the channel once the settings levels above it are taken into account, and every member of `sources` names the level a value comes from. When `false` or left out, `settings` holds what the channel itself stores.",
			QueryPath: "include_inherited",
		},
		&requestflag.Flag[[]string]{
			Name:       "beta",
			Usage:      "This endpoint is in beta: requests must send `experimental-claude-tag-2026-10-21` in this header.",
			HeaderPath: "anthropic-beta",
		},
	},
	Action:          handleBetaOrganizationClaudeTagSlackChannelsRetrieve,
	HideHelpCommand: true,
}

var betaOrganizationClaudeTagSlackChannelsUpdate = requestflag.WithInnerFlags(cli.Command{
	Name:    "update",
	Usage:   "Change the settings a Slack channel stores. The update is partial: a setting\nleft out is unchanged, and a setting sent as null is cleared, so the channel\ninherits it. If any value is refused, no setting in the request is stored. A\ngroup direct message takes no settings and returns 400.",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "id",
			Usage:     "The channel's identifier, starting with `slch_`.",
			Required:  true,
			PathParam: "id",
		},
		&requestflag.Flag[map[string]any]{
			Name:     "settings",
			Usage:    "The settings to change. A setting left out is unchanged. A setting sent as null is cleared, so it is inherited.",
			Required: true,
			BodyPath: "settings",
		},
		&requestflag.Flag[[]string]{
			Name:       "beta",
			Usage:      "This endpoint is in beta: requests must send `experimental-claude-tag-2026-10-21` in this header.",
			HeaderPath: "anthropic-beta",
		},
	},
	Action:          handleBetaOrganizationClaudeTagSlackChannelsUpdate,
	HideHelpCommand: true,
}, map[string][]requestflag.HasOuterFlag{
	"settings": {
		&requestflag.InnerFlag[*bool]{
			Name:       "settings.enabled",
			Usage:      "Whether Claude is on. `true` on a channel turns Claude on in it even when its workspace stores `false`. Null clears it.",
			InnerField: "enabled",
		},
		&requestflag.InnerFlag[*string]{
			Name:       "settings.environment-id",
			Usage:      "The environment Claude's sessions run in: the identifier of an active environment of the organization, starting with `env_`, or of a self-hosted runner pool of the organization that every member may use, starting with `ccpool_`. Null clears it.",
			InnerField: "environment_id",
		},
		&requestflag.InnerFlag[map[string]any]{
			Name:       "settings.guest-policy",
			Usage:      "What Claude does in a channel while a Slack guest is in it, for everyone in the channel. It changes nothing in a channel without Slack guests. More types can be added.",
			InnerField: "guest_policy",
		},
		&requestflag.InnerFlag[*string]{
			Name:       "settings.instructions",
			Usage:      "Instructions for Claude, up to 8,192 characters after leading and trailing white space is trimmed. Null or blank text clears them.",
			InnerField: "instructions",
		},
		&requestflag.InnerFlag[*string]{
			Name:       "settings.model",
			Usage:      "The model Claude uses: the identifier of a model, or the name of a model line, that the organization is entitled to for Claude Tag. Null clears it.",
			InnerField: "model",
		},
		&requestflag.InnerFlag[map[string]any]{
			Name:       "settings.reply-mode",
			Usage:      "How Claude takes part in a Slack channel. A channel alone stores one. More types can be added.",
			InnerField: "reply_mode",
		},
	},
})

var betaOrganizationClaudeTagSlackChannelsList = cli.Command{
	Name:    "list",
	Usage:   "List the set-up Slack channels of one Slack workspace, newest first by\n`created_at`. A channel that is not set up is not listed: Slack holds the full\nlist of a workspace's channels. Pass `slack_channel_id` as well to look one\nchannel up by its Slack identifiers: `data` then holds that channel, or is empty\nwhen it is not set up, and `limit` and `page` are ignored. A workspace that the\norganization owns only through a Slack Enterprise Grid returns 404 until it is\nset up. With `include_inherited=true` each channel's `settings` holds what\napplies in it, and the request fails with 503 rather than resolve against a\nsettings level that could not be read. A `page` cursor resumes only a request\nwith the same `include_inherited`.",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "slack-workspace-id",
			Usage:     "The Slack workspace whose channels to list: Slack's identifier for it, as in the channel object's `slack.workspace_id`, not the workspace object's `id`.",
			Required:  true,
			QueryPath: "slack_workspace_id",
		},
		&requestflag.Flag[bool]{
			Name:      "include-inherited",
			Usage:     "When `true`, `settings` holds what applies in each channel once the settings levels above it are taken into account, and every member of `sources` names the level a value comes from. When `false` or left out, `settings` holds what each channel itself stores.",
			QueryPath: "include_inherited",
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
		&requestflag.Flag[string]{
			Name:      "slack-channel-id",
			Usage:     "Slack's identifier for one channel of that workspace, as in the channel object's `slack.channel_id`. It is not the channel object's `id`. With it `data` holds that channel alone, or nothing when the channel is not set up.",
			QueryPath: "slack_channel_id",
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
	Action:          handleBetaOrganizationClaudeTagSlackChannelsList,
	HideHelpCommand: true,
}

var betaOrganizationClaudeTagSlackChannelsDelete = cli.Command{
	Name:    "delete",
	Usage:   "Remove a Slack channel from Claude Tag. This deletes the channel's stored\nsettings and its Bundle Targets at once, and then, in the background, the\nchannel's sessions and their transcripts, its memory, its routines and the\nartifacts published from it. None of it can be restored. The channel in Slack\nand its messages are unchanged. The request is refused with 400 and\n`error.details.error_code` `slack_channel_enabled` while `enabled` resolves to\n`true` for the channel, as a read with `include_inherited=true` returns it: set\n`enabled` to `false` on the channel first. This guards against removing a\nchannel in use by mistake; it is not a way to keep Claude out, because the\nremoval deletes the channel's own `enabled: false` with its other settings,\nafter which the channel inherits again and Claude's next use of it can set it up\nanew with a new identifier. The identifier returns 404 as soon as the request\nreturns, so a repeated removal returns 404. A group direct message returns 400,\nand so does a channel that holds configuration this API does not manage.",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "id",
			Usage:     "The channel's identifier, starting with `slch_`.",
			Required:  true,
			PathParam: "id",
		},
		&requestflag.Flag[[]string]{
			Name:       "beta",
			Usage:      "This endpoint is in beta: requests must send `experimental-claude-tag-2026-10-21` in this header.",
			HeaderPath: "anthropic-beta",
		},
	},
	Action:          handleBetaOrganizationClaudeTagSlackChannelsDelete,
	HideHelpCommand: true,
}

func handleBetaOrganizationClaudeTagSlackChannelsCreate(ctx context.Context, cmd *cli.Command) error {
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

	params := anthropic.BetaOrganizationClaudeTagSlackChannelNewParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Beta.Organization.ClaudeTag.SlackChannels.New(ctx, params, options...)
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
		Title:          "beta:organization:claude-tag:slack-channels create",
		Transform:      transform,
	})
}

func handleBetaOrganizationClaudeTagSlackChannelsRetrieve(ctx context.Context, cmd *cli.Command) error {
	client := anthropic.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("id") && len(unusedArgs) > 0 {
		cmd.Set("id", unusedArgs[0])
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

	params := anthropic.BetaOrganizationClaudeTagSlackChannelGetParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Beta.Organization.ClaudeTag.SlackChannels.Get(
		ctx,
		cmd.Value("id").(string),
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
		Title:          "beta:organization:claude-tag:slack-channels retrieve",
		Transform:      transform,
	})
}

func handleBetaOrganizationClaudeTagSlackChannelsUpdate(ctx context.Context, cmd *cli.Command) error {
	client := anthropic.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("id") && len(unusedArgs) > 0 {
		cmd.Set("id", unusedArgs[0])
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

	params := anthropic.BetaOrganizationClaudeTagSlackChannelUpdateParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Beta.Organization.ClaudeTag.SlackChannels.Update(
		ctx,
		cmd.Value("id").(string),
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
		Title:          "beta:organization:claude-tag:slack-channels update",
		Transform:      transform,
	})
}

func handleBetaOrganizationClaudeTagSlackChannelsList(ctx context.Context, cmd *cli.Command) error {
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

	params := anthropic.BetaOrganizationClaudeTagSlackChannelListParams{}

	format := "explore"
	explicitFormat := cmd.Root().IsSet("format")
	if explicitFormat {
		format = cmd.Root().String("format")
	}
	transform := cmd.Root().String("transform")
	if format == "raw" {
		var res []byte
		options = append(options, option.WithResponseBodyInto(&res))
		_, err = client.Beta.Organization.ClaudeTag.SlackChannels.List(ctx, params, options...)
		if err != nil {
			return err
		}
		obj := gjson.ParseBytes(res)
		return ShowJSON(obj, ShowJSONOpts{
			ExplicitFormat: explicitFormat,
			Format:         format,
			RawOutput:      cmd.Root().Bool("raw-output"),
			Title:          "beta:organization:claude-tag:slack-channels list",
			Transform:      transform,
		})
	} else {
		iter := client.Beta.Organization.ClaudeTag.SlackChannels.ListAutoPaging(ctx, params, options...)
		maxItems := int64(-1)
		if cmd.IsSet("max-items") {
			maxItems = cmd.Value("max-items").(int64)
		}
		return ShowJSONIterator(iter, maxItems, ShowJSONOpts{
			ExplicitFormat: explicitFormat,
			Format:         format,
			RawOutput:      cmd.Root().Bool("raw-output"),
			Title:          "beta:organization:claude-tag:slack-channels list",
			Transform:      transform,
		})
	}
}

func handleBetaOrganizationClaudeTagSlackChannelsDelete(ctx context.Context, cmd *cli.Command) error {
	client := anthropic.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("id") && len(unusedArgs) > 0 {
		cmd.Set("id", unusedArgs[0])
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

	params := anthropic.BetaOrganizationClaudeTagSlackChannelDeleteParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Beta.Organization.ClaudeTag.SlackChannels.Delete(
		ctx,
		cmd.Value("id").(string),
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
		Title:          "beta:organization:claude-tag:slack-channels delete",
		Transform:      transform,
	})
}
