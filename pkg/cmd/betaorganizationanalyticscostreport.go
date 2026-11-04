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

var betaOrganizationAnalyticsCostReportList = cli.Command{
	Name:    "list",
	Usage:   "[Legacy] Get cost in USD over time across a date range.",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[any]{
			Name:      "starting-at",
			Usage:     "Start of range, inclusive. RFC 3339 tz-aware. Must be within the last 365 days and no earlier than 2026-01-01T00:00:00Z.",
			Required:  true,
			QueryPath: "starting_at",
		},
		&requestflag.Flag[string]{
			Name:      "bucket-width",
			Usage:     "Time bucket granularity.",
			Default:   "1d",
			QueryPath: "bucket_width",
		},
		&requestflag.Flag[[]string]{
			Name:      "claude-tag-category",
			Usage:     "Filter to Claude Tag (Claude in Slack) usage in specific spend categories. Usage with no category never matches. Use `group_by[]=claude_tag_category` to break out per-category values.",
			QueryPath: "claude_tag_categories",
		},
		&requestflag.Flag[[]string]{
			Name:      "claude-tag-user-id",
			Usage:     "Filter to Claude Tag (Claude in Slack) usage attributed to specific Slack users, by Slack user ID (for example `U0123ABCDEF`), not claude.ai user ID. Usage that is not Claude Tag, and Claude Tag usage not attributed to a single user, never matches. Use `group_by[]=claude_tag_user_id` to break out per-user values.",
			QueryPath: "claude_tag_user_ids",
		},
		&requestflag.Flag[[]string]{
			Name:      "context-window",
			Usage:     "Filter to specific context-window pricing tiers. Use `group_by[]=context_window` to break out per-tier values.",
			QueryPath: "context_windows",
		},
		&requestflag.Flag[any]{
			Name:      "ending-at",
			Usage:     "End of range, exclusive. When omitted, defaults to the earlier of now and `starting_at` + 31 days. The range may span at most 31 days.",
			QueryPath: "ending_at",
		},
		&requestflag.Flag[[]string]{
			Name:      "group-by",
			Usage:     "Dimensions to break each time bucket out by. Defaults to no grouping (one total per bucket). Each bucket reports at most its top 100 groups; a group beyond that cap has no row in that bucket (there is no remainder row), so grouped buckets are not exhaustive when a dimension has more than 100 distinct values.",
			QueryPath: "group_by",
		},
		&requestflag.Flag[[]string]{
			Name:      "inference-geo",
			Usage:     "Filter to specific inference regions. `not_available` matches rows where the region is unset. Use `group_by[]=inference_geo` to break out per-region values.",
			QueryPath: "inference_geos",
		},
		&requestflag.Flag[int64]{
			Name:      "limit",
			Usage:     "Maximum number of time buckets per page. Defaults and caps vary by `bucket_width` (`1d`: default 7, max 31; `1h`: default 24, max 168; `1m`: default 60, max 256).",
			QueryPath: "limit",
		},
		&requestflag.Flag[[]string]{
			Name:      "model",
			Usage:     "Models to include. Defaults to all models. Use `group_by[]=model` to break out per-model values.",
			QueryPath: "models",
		},
		&requestflag.Flag[string]{
			Name:      "page",
			Usage:     "Opaque cursor from a previous response's `next_page` field.",
			QueryPath: "page",
		},
		&requestflag.Flag[[]string]{
			Name:      "product",
			Usage:     "Product surfaces to include. Defaults to all products. Use `group_by[]=product` to break out per-product values.",
			QueryPath: "products",
		},
		&requestflag.Flag[[]string]{
			Name:      "rbac-group-id",
			Usage:     "Filter to usage attributed to specific RBAC groups. Accepts tagged RBAC group IDs (`rbac_group_...`) or bare group UUIDs. A row matches when the user belonged to any of the listed groups on the (UTC) day the usage occurred; usage with no group attribution never matches.",
			QueryPath: "rbac_group_ids",
		},
		&requestflag.Flag[[]string]{
			Name:      "slack-channel-id",
			Usage:     "Filter to usage originating from specific Slack channels. Use `group_by[]=slack_channel_id` to break out per-channel values.",
			QueryPath: "slack_channel_ids",
		},
		&requestflag.Flag[[]string]{
			Name:      "speed",
			Usage:     "Filter to fast or standard inference mode. Use `group_by[]=speed` to break out per-mode values.",
			QueryPath: "speeds",
		},
		&requestflag.Flag[[]string]{
			Name:      "user-id",
			Usage:     "Filter to specific users by tagged user ID.",
			QueryPath: "user_ids",
		},
		&requestflag.Flag[int64]{
			Name:  "max-items",
			Usage: "The maximum number of items to return (use -1 for unlimited).",
		},
	},
	Action:          handleBetaOrganizationAnalyticsCostReportList,
	HideHelpCommand: true,
}

func handleBetaOrganizationAnalyticsCostReportList(ctx context.Context, cmd *cli.Command) error {
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

	params := anthropic.BetaOrganizationAnalyticsCostReportListParams{}

	format := "explore"
	explicitFormat := cmd.Root().IsSet("format")
	if explicitFormat {
		format = cmd.Root().String("format")
	}
	transform := cmd.Root().String("transform")
	if format == "raw" {
		var res []byte
		options = append(options, option.WithResponseBodyInto(&res))
		_, err = client.Beta.Organization.Analytics.CostReport.List(ctx, params, options...)
		if err != nil {
			return err
		}
		obj := gjson.ParseBytes(res)
		return ShowJSON(obj, ShowJSONOpts{
			ExplicitFormat: explicitFormat,
			Format:         format,
			RawOutput:      cmd.Root().Bool("raw-output"),
			Title:          "beta:organization:analytics:cost-report list",
			Transform:      transform,
		})
	} else {
		iter := client.Beta.Organization.Analytics.CostReport.ListAutoPaging(ctx, params, options...)
		maxItems := int64(-1)
		if cmd.IsSet("max-items") {
			maxItems = cmd.Value("max-items").(int64)
		}
		return ShowJSONIterator(iter, maxItems, ShowJSONOpts{
			ExplicitFormat: explicitFormat,
			Format:         format,
			RawOutput:      cmd.Root().Bool("raw-output"),
			Title:          "beta:organization:analytics:cost-report list",
			Transform:      transform,
		})
	}
}
