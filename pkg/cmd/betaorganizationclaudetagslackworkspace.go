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

var betaOrganizationClaudeTagSlackWorkspacesRetrieve = cli.Command{
	Name:    "retrieve",
	Usage:   "Get one Slack workspace and its settings. By default `settings` holds what the\nworkspace itself stores: a setting that is null there is not set on the\nworkspace and is inherited. With `include_inherited=true`, `settings` holds what\napplies in the workspace once the organization's Slack defaults are taken into\naccount, and the request fails with 503 rather than resolve against a settings\nlevel that could not be read. Either way `sources` names the settings level each\nvalue comes from.",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "id",
			Usage:     "The workspace's identifier, starting with `slws_`.",
			Required:  true,
			PathParam: "id",
		},
		&requestflag.Flag[bool]{
			Name:      "include-inherited",
			Usage:     "When `true`, `settings` holds what applies in the workspace once the organization's Slack defaults are taken into account, and every member of `sources` names the level a value comes from. When `false` or left out, `settings` holds what the workspace itself stores.",
			QueryPath: "include_inherited",
		},
		&requestflag.Flag[[]string]{
			Name:       "beta",
			Usage:      "This endpoint is in beta: requests must send `experimental-claude-tag-2026-10-21` in this header.",
			HeaderPath: "anthropic-beta",
		},
	},
	Action:          handleBetaOrganizationClaudeTagSlackWorkspacesRetrieve,
	HideHelpCommand: true,
}

var betaOrganizationClaudeTagSlackWorkspacesUpdate = requestflag.WithInnerFlags(cli.Command{
	Name:    "update",
	Usage:   "Change the settings a Slack workspace stores. The update is partial: a setting\nleft out is unchanged, and a setting sent as null is cleared, so the workspace\ninherits it. If any value is refused, no setting in the request is stored.",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "id",
			Usage:     "The workspace's identifier, starting with `slws_`.",
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
	Action:          handleBetaOrganizationClaudeTagSlackWorkspacesUpdate,
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
			Usage:      "The environment Claude's sessions run in: the identifier of an active environment of the organization, starting with `env_`. Null clears it.",
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

var betaOrganizationClaudeTagSlackWorkspacesList = cli.Command{
	Name:    "list",
	Usage:   "List the set-up Slack workspaces the organization owns, newest first by\n`created_at`. Pass `slack_workspace_id` to look one workspace up by its Slack\nidentifier: `data` then holds that workspace, or is empty when it is not set up.\nA workspace of a Slack Enterprise Grid the organization has connected returns\n404 until it is set up. With `include_inherited=true` each workspace's\n`settings` holds what applies in it, and the request fails with 503 rather than\nresolve against a settings level that could not be read. A `page` cursor resumes\nonly a request with the same `include_inherited`.",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[bool]{
			Name:      "include-inherited",
			Usage:     "When `true`, `settings` holds what applies in each workspace once the organization's Slack defaults are taken into account, and every member of `sources` names the level a value comes from. When `false` or left out, `settings` holds what each workspace itself stores.",
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
			Name:      "slack-workspace-id",
			Usage:     "Slack's identifier for one workspace, as in the workspace object's `slack.workspace_id`, not its `id`. With it `data` holds that workspace alone, or nothing when it is not set up.",
			QueryPath: "slack_workspace_id",
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
	Action:          handleBetaOrganizationClaudeTagSlackWorkspacesList,
	HideHelpCommand: true,
}

func handleBetaOrganizationClaudeTagSlackWorkspacesRetrieve(ctx context.Context, cmd *cli.Command) error {
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

	params := anthropic.BetaOrganizationClaudeTagSlackWorkspaceGetParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Beta.Organization.ClaudeTag.SlackWorkspaces.Get(
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
		Title:          "beta:organization:claude-tag:slack-workspaces retrieve",
		Transform:      transform,
	})
}

func handleBetaOrganizationClaudeTagSlackWorkspacesUpdate(ctx context.Context, cmd *cli.Command) error {
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

	params := anthropic.BetaOrganizationClaudeTagSlackWorkspaceUpdateParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Beta.Organization.ClaudeTag.SlackWorkspaces.Update(
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
		Title:          "beta:organization:claude-tag:slack-workspaces update",
		Transform:      transform,
	})
}

func handleBetaOrganizationClaudeTagSlackWorkspacesList(ctx context.Context, cmd *cli.Command) error {
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

	params := anthropic.BetaOrganizationClaudeTagSlackWorkspaceListParams{}

	format := "explore"
	explicitFormat := cmd.Root().IsSet("format")
	if explicitFormat {
		format = cmd.Root().String("format")
	}
	transform := cmd.Root().String("transform")
	if format == "raw" {
		var res []byte
		options = append(options, option.WithResponseBodyInto(&res))
		_, err = client.Beta.Organization.ClaudeTag.SlackWorkspaces.List(ctx, params, options...)
		if err != nil {
			return err
		}
		obj := gjson.ParseBytes(res)
		return ShowJSON(obj, ShowJSONOpts{
			ExplicitFormat: explicitFormat,
			Format:         format,
			RawOutput:      cmd.Root().Bool("raw-output"),
			Title:          "beta:organization:claude-tag:slack-workspaces list",
			Transform:      transform,
		})
	} else {
		iter := client.Beta.Organization.ClaudeTag.SlackWorkspaces.ListAutoPaging(ctx, params, options...)
		maxItems := int64(-1)
		if cmd.IsSet("max-items") {
			maxItems = cmd.Value("max-items").(int64)
		}
		return ShowJSONIterator(iter, maxItems, ShowJSONOpts{
			ExplicitFormat: explicitFormat,
			Format:         format,
			RawOutput:      cmd.Root().Bool("raw-output"),
			Title:          "beta:organization:claude-tag:slack-workspaces list",
			Transform:      transform,
		})
	}
}
