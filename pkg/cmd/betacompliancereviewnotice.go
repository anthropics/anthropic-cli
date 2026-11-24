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

var betaComplianceReviewNoticesList = cli.Command{
	Name:    "list",
	Usage:   "List Anthropic's notices to the organizations that the API key can reach, newest\nfirst by default. A notice is returned only if the API key can read every case\nin its `case_ids`.",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "case-id",
			Usage:     "Return only notices whose `case_ids` include this case (`rvc_...`). A case that the API key can't read returns an empty list.",
			QueryPath: "case_id",
		},
		&requestflag.Flag[any]{
			Name:      "created-at-gt",
			Usage:     "Return only notices whose `created_at` is after this time.",
			QueryPath: "created_at[gt]",
		},
		&requestflag.Flag[any]{
			Name:      "created-at-gte",
			Usage:     "Return notices whose `created_at` is at or after this time. Notices slightly older than this time are also returned, so that a poll misses no notice.",
			QueryPath: "created_at[gte]",
		},
		&requestflag.Flag[any]{
			Name:      "created-at-lt",
			Usage:     "Return only notices whose `created_at` is before this time.",
			QueryPath: "created_at[lt]",
		},
		&requestflag.Flag[any]{
			Name:      "created-at-lte",
			Usage:     "Return only notices whose `created_at` is at or before this time.",
			QueryPath: "created_at[lte]",
		},
		&requestflag.Flag[int64]{
			Name:      "limit",
			Usage:     "Number of notices to return per page, between 1 and 100. Defaults to 20.",
			QueryPath: "limit",
		},
		&requestflag.Flag[string]{
			Name:      "order",
			Usage:     "Sort direction, by `created_at` and then by `id`. Defaults to `desc`.\n\n- `asc` - Oldest first.\n- `desc` - Newest first.",
			QueryPath: "order",
		},
		&requestflag.Flag[string]{
			Name:      "page",
			Usage:     "Opaque cursor from a previous response's `next_page`. Leave it out to get the first page. A cursor works only with the filters and sort order of the request that returned it.",
			QueryPath: "page",
		},
		&requestflag.Flag[[]string]{
			Name:       "beta",
			Usage:      "Optional header to specify the beta version(s) you want to use.",
			HeaderPath: "anthropic-beta",
		},
		&requestflag.Flag[int64]{
			Name:  "max-items",
			Usage: "The maximum number of items to return (use -1 for unlimited).",
		},
	},
	Action:          handleBetaComplianceReviewNoticesList,
	HideHelpCommand: true,
}

func handleBetaComplianceReviewNoticesList(ctx context.Context, cmd *cli.Command) error {
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

	params := anthropic.BetaComplianceReviewNoticeListParams{}

	format := "explore"
	explicitFormat := cmd.Root().IsSet("format")
	if explicitFormat {
		format = cmd.Root().String("format")
	}
	transform := cmd.Root().String("transform")
	if format == "raw" {
		var res []byte
		options = append(options, option.WithResponseBodyInto(&res))
		_, err = client.Beta.Compliance.Review.Notices.List(ctx, params, options...)
		if err != nil {
			return err
		}
		obj := gjson.ParseBytes(res)
		return ShowJSON(obj, ShowJSONOpts{
			ExplicitFormat: explicitFormat,
			Format:         format,
			RawOutput:      cmd.Root().Bool("raw-output"),
			Title:          "beta:compliance:review:notices list",
			Transform:      transform,
		})
	} else {
		iter := client.Beta.Compliance.Review.Notices.ListAutoPaging(ctx, params, options...)
		maxItems := int64(-1)
		if cmd.IsSet("max-items") {
			maxItems = cmd.Value("max-items").(int64)
		}
		return ShowJSONIterator(iter, maxItems, ShowJSONOpts{
			ExplicitFormat: explicitFormat,
			Format:         format,
			RawOutput:      cmd.Root().Bool("raw-output"),
			Title:          "beta:compliance:review:notices list",
			Transform:      transform,
		})
	}
}
