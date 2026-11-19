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

var betaOrganizationWorkspacesRateLimitsUpdate = requestflag.WithInnerFlags(cli.Command{
	Name:    "update",
	Usage:   "Set or clear rate-limit overrides for one rate-limit group in a workspace.",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "workspace-id",
			Usage:     "The ID of the workspace.",
			Required:  true,
			PathParam: "workspace_id",
		},
		&requestflag.Flag[string]{
			Name:        "rate-limit-id",
			Usage:       "The `id` of the organization's rate-limit entry to set overrides for.",
			Required:    true,
			PathParam:   "rate_limit_id",
			DataAliases: []string{"id"},
		},
		&requestflag.Flag[[]map[string]any]{
			Name:     "limit",
			Usage:    "Limits to set or clear. Unspecified limit types are unchanged.",
			Required: true,
			BodyPath: "limits",
		},
	},
	Action:          handleBetaOrganizationWorkspacesRateLimitsUpdate,
	HideHelpCommand: true,
}, map[string][]requestflag.HasOuterFlag{
	"limit": {
		&requestflag.InnerFlag[string]{
			Name:       "limit.type",
			Usage:      "Rate limiter type. `input_tokens_per_minute` counts every input token, cache reads included. `input_tokens_per_minute_cache_aware` does not count cache reads, and on current models the organization's input-token limit usually has this type. A workspace limit of either type applies alongside the organization's limit, not in place of it. Because cache reads count toward an `input_tokens_per_minute` limit, it can reject requests that read mostly from the cache while the organization's limit still has room.",
			InnerField: "type",
		},
		&requestflag.InnerFlag[*int64]{
			Name:       "limit.value",
			Usage:      "New limit value. Set to null to clear the workspace override and fall back to the organization limit.",
			InnerField: "value",
		},
	},
})

var betaOrganizationWorkspacesRateLimitsList = cli.Command{
	Name:    "list",
	Usage:   "List a workspace's rate limits.",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "workspace-id",
			Usage:     "The ID of the workspace.",
			Required:  true,
			PathParam: "workspace_id",
		},
		&requestflag.Flag[string]{
			Name:      "group-type",
			Usage:     "Filter by group type.",
			QueryPath: "group_type",
		},
		&requestflag.Flag[bool]{
			Name:      "include-inherited",
			Usage:     "Also list the limiter values the workspace inherits from the organization, including groups with no workspace-level override.",
			Default:   false,
			QueryPath: "include_inherited",
		},
		&requestflag.Flag[int64]{
			Name:      "limit",
			Usage:     "Maximum number of items to return per page. Ranges from `1` to `1000`.\n\nWhen omitted, every remaining entry is returned in a single page and `next_page` is `null`.",
			QueryPath: "limit",
		},
		&requestflag.Flag[string]{
			Name:      "page",
			Usage:     "Opaque cursor from a previous response's `next_page`.",
			QueryPath: "page",
		},
		&requestflag.Flag[int64]{
			Name:  "max-items",
			Usage: "The maximum number of items to return (use -1 for unlimited).",
		},
	},
	Action:          handleBetaOrganizationWorkspacesRateLimitsList,
	HideHelpCommand: true,
}

func handleBetaOrganizationWorkspacesRateLimitsUpdate(ctx context.Context, cmd *cli.Command) error {
	client := anthropic.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("rate-limit-id") && len(unusedArgs) > 0 {
		cmd.Set("rate-limit-id", unusedArgs[0])
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

	params := anthropic.BetaOrganizationWorkspaceRateLimitUpdateParams{
		WorkspaceID: cmd.Value("workspace-id").(string),
	}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Beta.Organization.Workspaces.RateLimits.Update(
		ctx,
		cmd.Value("rate-limit-id").(string),
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
		Title:          "beta:organization:workspaces:rate-limits update",
		Transform:      transform,
	})
}

func handleBetaOrganizationWorkspacesRateLimitsList(ctx context.Context, cmd *cli.Command) error {
	client := anthropic.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("workspace-id") && len(unusedArgs) > 0 {
		cmd.Set("workspace-id", unusedArgs[0])
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

	params := anthropic.BetaOrganizationWorkspaceRateLimitListParams{}

	format := "explore"
	explicitFormat := cmd.Root().IsSet("format")
	if explicitFormat {
		format = cmd.Root().String("format")
	}
	transform := cmd.Root().String("transform")
	if format == "raw" {
		var res []byte
		options = append(options, option.WithResponseBodyInto(&res))
		_, err = client.Beta.Organization.Workspaces.RateLimits.List(
			ctx,
			cmd.Value("workspace-id").(string),
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
			Title:          "beta:organization:workspaces:rate-limits list",
			Transform:      transform,
		})
	} else {
		iter := client.Beta.Organization.Workspaces.RateLimits.ListAutoPaging(
			ctx,
			cmd.Value("workspace-id").(string),
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
			Title:          "beta:organization:workspaces:rate-limits list",
			Transform:      transform,
		})
	}
}
