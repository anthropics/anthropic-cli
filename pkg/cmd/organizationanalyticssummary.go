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

var organizationAnalyticsSummariesList = cli.Command{
	Name:    "list",
	Usage:   "Get organization-wide activity summaries for a date range.",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[any]{
			Name:      "starting-date",
			Usage:     "UTC date in YYYY-MM-DD format. Start of the date range (inclusive). Data is typically available with a 1-day lag (varies by query; the error for a too-recent date names the latest available day) and may be revised by a few percent over the following days. No earlier than 2026-01-01.",
			Required:  true,
			QueryPath: "starting_date",
		},
		&requestflag.Flag[any]{
			Name:      "ending-date",
			Usage:     "UTC date in YYYY-MM-DD format. End of the date range (exclusive). Data is typically available with a 1-day lag, so this can be at most today — which is also the default when omitted, making the last entry cover the most recent available day. Data may be revised by a few percent over the following days. The range may span at most 366 days.",
			QueryPath: "ending_date",
		},
		&requestflag.Flag[[]string]{
			Name:      "filter",
			Usage:     "Filters as `dimension:value`. Only `rbac_group_id` is supported (e.g. `filter[]=rbac_group_id:{id}`); repeat the param to OR across groups. Scopes the whole day series to members of the matching group(s), re-aggregated from member-level activity — org-wide seat/invite fields and the adoption rates derived from them are null on scoped rows. `rbac_group_id` accepts the tagged id (`rbac_group_...`, as emitted in responses and by the spend-limits API) or a bare group UUID, and matches users who held the group at any point during each UTC day (time-of-usage attribution). At most 100 entries.",
			QueryPath: "filter",
		},
		&requestflag.Flag[int64]{
			Name:      "limit",
			Usage:     "Number of results per page (1-1000, default 100). The day series (at most 366 entries) is currently returned in full in a single page, so `limit` does not yet shorten it.",
			QueryPath: "limit",
		},
		&requestflag.Flag[string]{
			Name:      "page",
			Usage:     "Opaque cursor from a previous response's `next_page` field. `next_page` is currently always null, so there is never a cursor to send.",
			QueryPath: "page",
		},
		&requestflag.Flag[int64]{
			Name:  "max-items",
			Usage: "The maximum number of items to return (use -1 for unlimited).",
		},
	},
	Action:          handleOrganizationAnalyticsSummariesList,
	HideHelpCommand: true,
}

func handleOrganizationAnalyticsSummariesList(ctx context.Context, cmd *cli.Command) error {
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

	params := anthropic.OrganizationAnalyticsSummaryListParams{}

	format := "explore"
	explicitFormat := cmd.Root().IsSet("format")
	if explicitFormat {
		format = cmd.Root().String("format")
	}
	transform := cmd.Root().String("transform")
	if format == "raw" {
		var res []byte
		options = append(options, option.WithResponseBodyInto(&res))
		_, err = client.Organization.Analytics.Summaries.List(ctx, params, options...)
		if err != nil {
			return err
		}
		obj := gjson.ParseBytes(res)
		return ShowJSON(obj, ShowJSONOpts{
			ExplicitFormat: explicitFormat,
			Format:         format,
			RawOutput:      cmd.Root().Bool("raw-output"),
			Title:          "organization:analytics:summaries list",
			Transform:      transform,
		})
	} else {
		iter := client.Organization.Analytics.Summaries.ListAutoPaging(ctx, params, options...)
		maxItems := int64(-1)
		if cmd.IsSet("max-items") {
			maxItems = cmd.Value("max-items").(int64)
		}
		return ShowJSONIterator(iter, maxItems, ShowJSONOpts{
			ExplicitFormat: explicitFormat,
			Format:         format,
			RawOutput:      cmd.Root().Bool("raw-output"),
			Title:          "organization:analytics:summaries list",
			Transform:      transform,
		})
	}
}
