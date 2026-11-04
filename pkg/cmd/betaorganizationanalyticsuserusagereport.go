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

var betaOrganizationAnalyticsUserUsageReportList = cli.Command{
	Name:    "list",
	Usage:   "[Legacy] Get per-user token usage across a date range.",
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
			Usage:     "Time-bucket granularity. When set, each row's `starting_at` and `ending_at` are populated and one actor may span several rows (one per time bucket with usage). The time bucket counts toward `limit`, so one page can return multiple rows for the same actor. `ending_at` is required when `bucket_width` is set, and with `bucket_width=\"1m\"` the range may span at most 24 hours. When omitted, each row aggregates the full `[starting_at, ending_at)` range.",
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
		&requestflag.Flag[bool]{
			Name:      "exclude-deleted-users",
			Usage:     "If true, omit rows for users who are deleted (`deleted: true`). A page may contain fewer than `limit` rows; use `has_more` and `next_page` to paginate as usual.",
			Default:   false,
			QueryPath: "exclude_deleted_users",
		},
		&requestflag.Flag[[]string]{
			Name:      "group-by",
			Usage:     "Break each actor's row out by the given dimensions. Accepts the same values as the bucketed `/usage_report` endpoint. `limit` bounds (actor × time bucket × dimension) rows — with dimensions or `bucket_width` present, one actor may span several rows.",
			QueryPath: "group_by",
		},
		&requestflag.Flag[[]string]{
			Name:      "inference-geo",
			Usage:     "Filter to specific inference regions. `not_available` matches rows where the region is unset. Use `group_by[]=inference_geo` to break out per-region values.",
			QueryPath: "inference_geos",
		},
		&requestflag.Flag[int64]{
			Name:      "limit",
			Usage:     "Number of rows per page (1-1000, default 20). One row per actor unless `group_by[]` or `bucket_width` splits an actor across rows; `cost_type`/`token_type` fan-out rows (cost endpoint only) are the exception — they do not count toward this limit, so `data` can exceed it.",
			Default:   20,
			QueryPath: "limit",
		},
		&requestflag.Flag[[]string]{
			Name:      "model",
			Usage:     "Models to include. Defaults to all models. Use `group_by[]=model` to break out per-model values.",
			QueryPath: "models",
		},
		&requestflag.Flag[string]{
			Name:      "order",
			Usage:     "Sort direction. Defaults to `desc`.",
			Default:   "desc",
			QueryPath: "order",
		},
		&requestflag.Flag[string]{
			Name:      "order-by",
			Usage:     "Metric to rank actors by. Defaults to `total_tokens`.",
			Default:   "total_tokens",
			QueryPath: "order_by",
		},
		&requestflag.Flag[string]{
			Name:      "page",
			Usage:     "Opaque cursor from a previous response's `next_page` field.",
			QueryPath: "page",
		},
		&requestflag.Flag[[]string]{
			Name:      "product",
			Usage:     "Product surfaces to include. Defaults to all products.",
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
	Action:          handleBetaOrganizationAnalyticsUserUsageReportList,
	HideHelpCommand: true,
}

func handleBetaOrganizationAnalyticsUserUsageReportList(ctx context.Context, cmd *cli.Command) error {
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

	params := anthropic.BetaOrganizationAnalyticsUserUsageReportListParams{}

	format := "explore"
	explicitFormat := cmd.Root().IsSet("format")
	if explicitFormat {
		format = cmd.Root().String("format")
	}
	transform := cmd.Root().String("transform")
	if format == "raw" {
		var res []byte
		options = append(options, option.WithResponseBodyInto(&res))
		_, err = client.Beta.Organization.Analytics.UserUsageReport.List(ctx, params, options...)
		if err != nil {
			return err
		}
		obj := gjson.ParseBytes(res)
		return ShowJSON(obj, ShowJSONOpts{
			ExplicitFormat: explicitFormat,
			Format:         format,
			RawOutput:      cmd.Root().Bool("raw-output"),
			Title:          "beta:organization:analytics:user-usage-report list",
			Transform:      transform,
		})
	} else {
		iter := client.Beta.Organization.Analytics.UserUsageReport.ListAutoPaging(ctx, params, options...)
		maxItems := int64(-1)
		if cmd.IsSet("max-items") {
			maxItems = cmd.Value("max-items").(int64)
		}
		return ShowJSONIterator(iter, maxItems, ShowJSONOpts{
			ExplicitFormat: explicitFormat,
			Format:         format,
			RawOutput:      cmd.Root().Bool("raw-output"),
			Title:          "beta:organization:analytics:user-usage-report list",
			Transform:      transform,
		})
	}
}
