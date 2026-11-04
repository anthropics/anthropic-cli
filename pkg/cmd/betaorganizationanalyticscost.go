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

var betaOrganizationAnalyticsCostList = cli.Command{
	Name:    "list",
	Usage:   "Get cost in USD across a date range, broken out by dimension.",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[any]{
			Name:      "starting-at",
			Usage:     "Start of range, inclusive. RFC 3339 tz-aware. Must be within the last 365 days and no earlier than 2026-01-01T00:00:00Z.",
			Required:  true,
			QueryPath: "starting_at",
		},
		&requestflag.Flag[[]string]{
			Name:      "actor-id",
			Usage:     "Filter to usage attributed to specific identities, by tagged ID: a `user_...`, `svac_...` (service account), or `apikey_...` ID. An `apikey_` entry matches only usage attributed to the key itself, not usage on that key attributed to a user or service account (use `api_key_ids[]` for every request a key authenticated). Federated identities have no tagged ID; select them with `actor_types[]`. Use `group_by[]=actor` to break out per-identity values.",
			QueryPath: "actor_ids",
		},
		&requestflag.Flag[[]string]{
			Name:      "actor-type",
			Usage:     "Filter to usage attributed to specific kinds of identity. Use `group_by[]=actor_type` to break out per-kind values.",
			QueryPath: "actor_types",
		},
		&requestflag.Flag[[]string]{
			Name:      "agent-id",
			Usage:     "Filter to usage incurred by sessions of specific Claude Managed Agents agents. Use `group_by[]=agent_id` to break out per-agent values. Claude Platform organizations only.",
			QueryPath: "agent_ids",
		},
		&requestflag.Flag[[]string]{
			Name:      "api-key-id",
			Usage:     "Filter to usage authenticated by specific API keys, by tagged API key ID (`apikey_...`), whichever identity the usage is attributed to. Include an empty entry (`api_key_ids[]=`) to also match usage that no API key authenticated. Use `group_by[]=api_key_id` to break out per-key values; usage without a key has a null `api_key_id`. Claude Platform organizations only.",
			QueryPath: "api_key_ids",
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
		&requestflag.Flag[[]string]{
			Name:      "deployment-id",
			Usage:     "Filter to usage incurred by Claude Managed Agents sessions that specific deployments started. Use `group_by[]=deployment_id` to break out per-deployment values. Claude Platform organizations only.",
			QueryPath: "deployment_ids",
		},
		&requestflag.Flag[any]{
			Name:      "ending-at",
			Usage:     "End of range, exclusive. Required when `time_bucket_width` is set. Without `time_bucket_width`, defaults to the earlier of now and `starting_at` + 31 days. The range may span at most 31 days.",
			QueryPath: "ending_at",
		},
		&requestflag.Flag[[]string]{
			Name:      "environment-id",
			Usage:     "Filter to usage incurred by Claude Managed Agents sessions that ran in specific environments. Use `group_by[]=environment_id` to break out per-environment values. Claude Platform organizations only.",
			QueryPath: "environment_ids",
		},
		&requestflag.Flag[bool]{
			Name:      "exclude-deleted-users",
			Usage:     "If true, omit rows whose `actor` is a `user_actor` with `deleted: true`. Requires `group_by[]=actor`. Rows for other actor types, and rows whose `actor` is null, are always kept. A page may contain fewer than `limit` rows; use `has_more` and `next_page` to paginate as usual.",
			Default:   false,
			QueryPath: "exclude_deleted_users",
		},
		&requestflag.Flag[[]string]{
			Name:      "group-by",
			Usage:     "Dimensions to break cost out by (at most 5, not counting `cost_type` and `token_type`): one row per combination of their values with cost in the range (and per time bucket, when `time_bucket_width` is set), every combination returned across pages rather than a capped top set. `cost_type` and `token_type` do not count toward `limit`: `cost_type` returns one row per cost component (tokens, web search, code execution, and `session_usage`, Claude Managed Agents session time); `token_type` returns one row per token type, each with `cost_type: \"tokens\"`; combining both returns the per-token-type rows plus the web-search, code-execution and session-usage rows. When omitted, a single row aggregates the whole range.",
			QueryPath: "group_by",
		},
		&requestflag.Flag[[]string]{
			Name:      "inference-geo",
			Usage:     "Filter to specific inference regions. `not_available` matches rows where the region is unset. Use `group_by[]=inference_geo` to break out per-region values.",
			QueryPath: "inference_geos",
		},
		&requestflag.Flag[int64]{
			Name:      "limit",
			Usage:     "Number of rows per page (1-1000, default 20). `cost_type`/`token_type` fan-out rows (cost endpoint only) do not count toward this limit, so `data` can exceed it.",
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
			Usage:     `Allowed values: "asc", "desc".`,
			QueryPath: "order",
		},
		&requestflag.Flag[string]{
			Name:      "order-by",
			Usage:     `Allowed values: "amount", "list_amount", "time".`,
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
			Name:      "service-tier",
			Usage:     "Filter to specific service tiers. Use `group_by[]=service_tier` to break out per-tier values. Claude Platform organizations only; matches nothing for a Claude Enterprise organization.",
			QueryPath: "service_tiers",
		},
		&requestflag.Flag[[]string]{
			Name:      "session-id",
			Usage:     "Filter to usage incurred by specific Claude Managed Agents sessions, by tagged session ID (`sesn_...`). Use `group_by[]=session_id` to break out per-session values. Claude Platform organizations only.",
			QueryPath: "session_ids",
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
		&requestflag.Flag[string]{
			Name:      "time-bucket-width",
			Usage:     `Allowed values: "1d", "1h", "1m".`,
			QueryPath: "time_bucket_width",
		},
		&requestflag.Flag[[]string]{
			Name:      "user-profile-id",
			Usage:     "Filter to usage attributed to specific user profiles, by tagged user profile ID (`uprof_...`). Use `group_by[]=user_profile_id` to break out per-profile values. Claude Platform organizations only.",
			QueryPath: "user_profile_ids",
		},
		&requestflag.Flag[[]string]{
			Name:      "workspace-id",
			Usage:     "Filter to usage in specific workspaces, by tagged workspace ID (`wrkspc_...`). Use `group_by[]=workspace_id` to break out per-workspace values. Claude Platform organizations only; matches nothing for a Claude Enterprise organization.",
			QueryPath: "workspace_ids",
		},
		&requestflag.Flag[[]string]{
			Name:       "beta",
			Usage:      "This endpoint is in beta: requests must send `analytics-usage-cost-2026-08-24` in this header.",
			HeaderPath: "anthropic-beta",
		},
		&requestflag.Flag[int64]{
			Name:  "max-items",
			Usage: "The maximum number of items to return (use -1 for unlimited).",
		},
	},
	Action:          handleBetaOrganizationAnalyticsCostList,
	HideHelpCommand: true,
}

func handleBetaOrganizationAnalyticsCostList(ctx context.Context, cmd *cli.Command) error {
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

	params := anthropic.BetaOrganizationAnalyticsCostListParams{}

	format := "explore"
	explicitFormat := cmd.Root().IsSet("format")
	if explicitFormat {
		format = cmd.Root().String("format")
	}
	transform := cmd.Root().String("transform")
	if format == "raw" {
		var res []byte
		options = append(options, option.WithResponseBodyInto(&res))
		_, err = client.Beta.Organization.Analytics.Cost.List(ctx, params, options...)
		if err != nil {
			return err
		}
		obj := gjson.ParseBytes(res)
		return ShowJSON(obj, ShowJSONOpts{
			ExplicitFormat: explicitFormat,
			Format:         format,
			RawOutput:      cmd.Root().Bool("raw-output"),
			Title:          "beta:organization:analytics:cost list",
			Transform:      transform,
		})
	} else {
		iter := client.Beta.Organization.Analytics.Cost.ListAutoPaging(ctx, params, options...)
		maxItems := int64(-1)
		if cmd.IsSet("max-items") {
			maxItems = cmd.Value("max-items").(int64)
		}
		return ShowJSONIterator(iter, maxItems, ShowJSONOpts{
			ExplicitFormat: explicitFormat,
			Format:         format,
			RawOutput:      cmd.Root().Bool("raw-output"),
			Title:          "beta:organization:analytics:cost list",
			Transform:      transform,
		})
	}
}
