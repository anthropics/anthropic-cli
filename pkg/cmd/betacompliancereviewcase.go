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

var betaComplianceReviewCasesRetrieve = cli.Command{
	Name:    "retrieve",
	Usage:   "Read one review case.",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:        "case-id",
			Usage:       "The ID of the case (`rvc_...`).",
			Required:    true,
			PathParam:   "case_id",
			DataAliases: []string{"id"},
		},
		&requestflag.Flag[[]string]{
			Name:       "beta",
			Usage:      "Optional header to specify the beta version(s) you want to use.",
			HeaderPath: "anthropic-beta",
		},
	},
	Action:          handleBetaComplianceReviewCasesRetrieve,
	HideHelpCommand: true,
}

var betaComplianceReviewCasesList = cli.Command{
	Name:    "list",
	Usage:   "List the review cases of every organization that the API key can reach, newest\nfirst by default. Both open and closed cases are returned unless `status` is\nset.",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[any]{
			Name:      "created-at-gt",
			Usage:     "Return only cases whose `created_at` is after this time.",
			QueryPath: "created_at[gt]",
		},
		&requestflag.Flag[any]{
			Name:      "created-at-gte",
			Usage:     "Return only cases whose `created_at` is at or after this time.",
			QueryPath: "created_at[gte]",
		},
		&requestflag.Flag[any]{
			Name:      "created-at-lt",
			Usage:     "Return only cases whose `created_at` is before this time.",
			QueryPath: "created_at[lt]",
		},
		&requestflag.Flag[any]{
			Name:      "created-at-lte",
			Usage:     "Return only cases whose `created_at` is at or before this time.",
			QueryPath: "created_at[lte]",
		},
		&requestflag.Flag[int64]{
			Name:      "limit",
			Usage:     "Number of cases to return per page, between 1 and 100. Defaults to 20.",
			QueryPath: "limit",
		},
		&requestflag.Flag[string]{
			Name:      "order",
			Usage:     "Sort direction, by `order_by` and then by `id`. Defaults to `desc`.\n\n- `asc` - Oldest first.\n- `desc` - Newest first.",
			QueryPath: "order",
		},
		&requestflag.Flag[string]{
			Name:      "order-by",
			Usage:     "The field to sort by. Defaults to `created_at`.\n\n- `created_at` - Sort by when Anthropic sent the case.\n- `updated_at` - Sort by when the case last changed.",
			QueryPath: "order_by",
		},
		&requestflag.Flag[string]{
			Name:      "page",
			Usage:     "Opaque cursor from a previous response's `next_page`. Leave it out to get the first page. A cursor works only with the filters and sort order of the request that returned it.",
			QueryPath: "page",
		},
		&requestflag.Flag[string]{
			Name:      "status",
			Usage:     "Return only cases with this status. Leave it out to get both open and closed cases.\n\n- `open` - Cases that Anthropic is waiting for a response to.\n- `closed` - Cases that Anthropic has stopped waiting on.",
			QueryPath: "status",
		},
		&requestflag.Flag[any]{
			Name:      "updated-at-gt",
			Usage:     "Return only cases whose `updated_at` is after this time.",
			QueryPath: "updated_at[gt]",
		},
		&requestflag.Flag[any]{
			Name:      "updated-at-gte",
			Usage:     "Return cases whose `updated_at` is at or after this time. Cases slightly older than this time are also returned, so that a sync misses no change.",
			QueryPath: "updated_at[gte]",
		},
		&requestflag.Flag[any]{
			Name:      "updated-at-lt",
			Usage:     "Return only cases whose `updated_at` is before this time.",
			QueryPath: "updated_at[lt]",
		},
		&requestflag.Flag[any]{
			Name:      "updated-at-lte",
			Usage:     "Return only cases whose `updated_at` is at or before this time.",
			QueryPath: "updated_at[lte]",
		},
		&requestflag.Flag[string]{
			Name:      "user-profile-id",
			Usage:     "Return only cases where a message's `user_profile_id` is this ID, including cases that also have messages of other user profiles. A case stops matching once its messages' `user_profile_id` is null. An ID that matches no reachable case returns no cases. A request that sends the `anthropic-user-profile-id` header without this parameter returns a 400 error.",
			QueryPath: "user_profile_id",
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
	Action:          handleBetaComplianceReviewCasesList,
	HideHelpCommand: true,
}

var betaComplianceReviewCasesRespond = requestflag.WithInnerFlags(cli.Command{
	Name:    "respond",
	Usage:   "Send your conclusion on a review case. A response closes the case, and a closed\ncase can't be reopened.",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:        "case-id",
			Usage:       "The ID of the case (`rvc_...`).",
			Required:    true,
			PathParam:   "case_id",
			DataAliases: []string{"id"},
		},
		&requestflag.Flag[string]{
			Name:     "determination",
			Usage:    "Your conclusion on the case as a whole. Anthropic doesn't check it against `verdicts`.\n\n- `violation` - The messages violated a policy.\n- `no_violation` - The messages violated no policy.",
			Required: true,
			BodyPath: "determination",
		},
		&requestflag.Flag[[]string]{
			Name:     "policy-category",
			Usage:    "The policy categories that you found were violated. Send one to four with `violation`, and none with `no_violation`. They need not match the case's own. Each value is one of `cyber`, `bio`, `chem`, `radnuc` (radiological or nuclear weapons), `cgw` (conventional weapons), `mgw` (military-grade weapons systems), `io` (influence operations), `fraud`, `violent_extremism`, `surveillance`, `human_exploitation`, `self_harm`, `adult_content`, `reasoning_extraction`, `frontier_llm` (use of the model to accelerate frontier AI research and development in restricted ways), `compromised_key` or `other`, and can appear once.",
			BodyPath: "policy_categories",
		},
		&requestflag.Flag[[]map[string]any]{
			Name:     "verdict",
			Usage:    "Your conclusions on individual messages of the case. They need not cover every message, and each `item_id` can appear once.",
			BodyPath: "verdicts",
		},
		&requestflag.Flag[[]string]{
			Name:       "beta",
			Usage:      "Optional header to specify the beta version(s) you want to use.",
			HeaderPath: "anthropic-beta",
		},
	},
	Action:          handleBetaComplianceReviewCasesRespond,
	HideHelpCommand: true,
}, map[string][]requestflag.HasOuterFlag{
	"verdict": {
		&requestflag.InnerFlag[string]{
			Name:       "verdict.determination",
			Usage:      "Your conclusion on this message.\n\n- `violation` - The messages violated a policy.\n- `no_violation` - The messages violated no policy.",
			InnerField: "determination",
		},
		&requestflag.InnerFlag[string]{
			Name:       "verdict.item-id",
			Usage:      "The `item_id` of one of the case's `transcripts`.",
			InnerField: "item_id",
		},
		&requestflag.InnerFlag[*string]{
			Name:       "verdict.policy-category",
			Usage:      "The policy category that this message violated, from the values that the request's `policy_categories` accepts. Optional with `violation`. Leave it out with `no_violation`.\n\nAllowed values: \"cyber\", \"bio\", \"chem\", \"radnuc\", \"cgw\", \"mgw\", \"io\", \"fraud\", \"violent_extremism\", \"surveillance\", \"human_exploitation\", \"self_harm\", \"adult_content\", \"reasoning_extraction\", \"frontier_llm\", \"compromised_key\", \"other\".",
			InnerField: "policy_category",
		},
	},
})

func handleBetaComplianceReviewCasesRetrieve(ctx context.Context, cmd *cli.Command) error {
	client := anthropic.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("case-id") && len(unusedArgs) > 0 {
		cmd.Set("case-id", unusedArgs[0])
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

	params := anthropic.BetaComplianceReviewCaseGetParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Beta.Compliance.Review.Cases.Get(
		ctx,
		cmd.Value("case-id").(string),
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
		Title:          "beta:compliance:review:cases retrieve",
		Transform:      transform,
	})
}

func handleBetaComplianceReviewCasesList(ctx context.Context, cmd *cli.Command) error {
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

	params := anthropic.BetaComplianceReviewCaseListParams{}

	format := "explore"
	explicitFormat := cmd.Root().IsSet("format")
	if explicitFormat {
		format = cmd.Root().String("format")
	}
	transform := cmd.Root().String("transform")
	if format == "raw" {
		var res []byte
		options = append(options, option.WithResponseBodyInto(&res))
		_, err = client.Beta.Compliance.Review.Cases.List(ctx, params, options...)
		if err != nil {
			return err
		}
		obj := gjson.ParseBytes(res)
		return ShowJSON(obj, ShowJSONOpts{
			ExplicitFormat: explicitFormat,
			Format:         format,
			RawOutput:      cmd.Root().Bool("raw-output"),
			Title:          "beta:compliance:review:cases list",
			Transform:      transform,
		})
	} else {
		iter := client.Beta.Compliance.Review.Cases.ListAutoPaging(ctx, params, options...)
		maxItems := int64(-1)
		if cmd.IsSet("max-items") {
			maxItems = cmd.Value("max-items").(int64)
		}
		return ShowJSONIterator(iter, maxItems, ShowJSONOpts{
			ExplicitFormat: explicitFormat,
			Format:         format,
			RawOutput:      cmd.Root().Bool("raw-output"),
			Title:          "beta:compliance:review:cases list",
			Transform:      transform,
		})
	}
}

func handleBetaComplianceReviewCasesRespond(ctx context.Context, cmd *cli.Command) error {
	client := anthropic.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("case-id") && len(unusedArgs) > 0 {
		cmd.Set("case-id", unusedArgs[0])
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

	params := anthropic.BetaComplianceReviewCaseRespondParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Beta.Compliance.Review.Cases.Respond(
		ctx,
		cmd.Value("case-id").(string),
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
		Title:          "beta:compliance:review:cases respond",
		Transform:      transform,
	})
}
