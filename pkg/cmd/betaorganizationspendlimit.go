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

var betaOrganizationSpendLimitsRetrieve = cli.Command{
	Name:    "retrieve",
	Usage:   "Retrieve a spend limit by ID.",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "spend-limit-id",
			Usage:     "ID of the Spend Limit.",
			Required:  true,
			PathParam: "spend_limit_id",
		},
	},
	Action:          handleBetaOrganizationSpendLimitsRetrieve,
	HideHelpCommand: true,
}

var betaOrganizationSpendLimitsList = cli.Command{
	Name:    "list",
	Usage:   "List the organization's spend limits.",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[int64]{
			Name:      "limit",
			Usage:     "Maximum number of limits per page. Defaults to `20`.",
			Default:   20,
			QueryPath: "limit",
		},
		&requestflag.Flag[string]{
			Name:      "page",
			Usage:     "Opaque cursor from a previous response's `next_page` field.",
			QueryPath: "page",
		},
		&requestflag.Flag[[]string]{
			Name:      "scope-type",
			Usage:     "Return only limits with these scope types. Omit for all. A scope type that doesn't apply to the organization gets a 400 error.",
			QueryPath: "scope_type",
		},
		&requestflag.Flag[string]{
			Name:      "workspace-id",
			Usage:     "Return only the limits whose `scope.workspace_id` is this workspace (`wrkspc_...`). Claude Console organizations only; other organizations get a 400 error.",
			QueryPath: "workspace_id",
		},
		&requestflag.Flag[[]string]{
			Name:       "beta",
			Usage:      "This endpoint is in beta: requests must send `spend-limit-reads-2026-09-26` in this header.",
			HeaderPath: "anthropic-beta",
		},
		&requestflag.Flag[int64]{
			Name:  "max-items",
			Usage: "The maximum number of items to return (use -1 for unlimited).",
		},
	},
	Action:          handleBetaOrganizationSpendLimitsList,
	HideHelpCommand: true,
}

var betaOrganizationSpendLimitsDelete = cli.Command{
	Name:    "delete",
	Usage:   "Delete a spend limit.",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "spend-limit-id",
			Usage:     "ID of the Spend Limit.",
			Required:  true,
			PathParam: "spend_limit_id",
		},
	},
	Action:          handleBetaOrganizationSpendLimitsDelete,
	HideHelpCommand: true,
}

var betaOrganizationSpendLimitsSet = requestflag.WithInnerFlags(cli.Command{
	Name:    "set",
	Usage:   "Set a spend limit.",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[*string]{
			Name:     "amount",
			Usage:    "Limit amount as a non-negative integer decimal string in the minor unit of the organization's billing currency (cents for USD): \"50000\" is $500.00. `null` sets an explicit no-limit override for this scope and `period` only — each period resolves independently, so caps for other periods still apply.",
			Required: true,
			BodyPath: "amount",
		},
		&requestflag.Flag[map[string]any]{
			Name:     "scope",
			Usage:    "What the limit applies to. Claude Enterprise organizations set `user` limits. Claude Console organizations set `organization` and `workspace` limits. Any other combination returns 400. Setting `organization` and `workspace` limits through the API is in an early access preview. To request access, contact your Anthropic account team.",
			Required: true,
			BodyPath: "scope",
		},
		&requestflag.Flag[string]{
			Name:     "period",
			Usage:    `Allowed values: "daily", "monthly", "weekly".`,
			BodyPath: "period",
		},
		&requestflag.Flag[[]string]{
			Name:       "beta",
			Usage:      "Optional header to specify the beta version(s) you want to use.",
			HeaderPath: "anthropic-beta",
		},
	},
	Action:          handleBetaOrganizationSpendLimitsSet,
	HideHelpCommand: true,
}, map[string][]requestflag.HasOuterFlag{
	"scope": {
		&requestflag.InnerFlag[string]{
			Name:       "scope.type",
			Usage:      "Scope type. Always `user` for this scope.\n\nAllowed values: \"user\", \"organization\", \"workspace\".",
			InnerField: "type",
		},
		&requestflag.InnerFlag[string]{
			Name:       "scope.user-id",
			Usage:      "Tagged ID of the member the spend limit applies to.",
			InnerField: "user_id",
		},
		&requestflag.InnerFlag[string]{
			Name:       "scope.workspace-id",
			Usage:      "Tagged ID of the workspace the spend limit applies to.",
			InnerField: "workspace_id",
		},
	},
})

func handleBetaOrganizationSpendLimitsRetrieve(ctx context.Context, cmd *cli.Command) error {
	client := anthropic.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("spend-limit-id") && len(unusedArgs) > 0 {
		cmd.Set("spend-limit-id", unusedArgs[0])
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

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Beta.Organization.SpendLimits.Get(ctx, cmd.Value("spend-limit-id").(string), options...)
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
		Title:          "beta:organization:spend-limits retrieve",
		Transform:      transform,
	})
}

func handleBetaOrganizationSpendLimitsList(ctx context.Context, cmd *cli.Command) error {
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

	params := anthropic.BetaOrganizationSpendLimitListParams{}

	format := "explore"
	explicitFormat := cmd.Root().IsSet("format")
	if explicitFormat {
		format = cmd.Root().String("format")
	}
	transform := cmd.Root().String("transform")
	if format == "raw" {
		var res []byte
		options = append(options, option.WithResponseBodyInto(&res))
		_, err = client.Beta.Organization.SpendLimits.List(ctx, params, options...)
		if err != nil {
			return err
		}
		obj := gjson.ParseBytes(res)
		return ShowJSON(obj, ShowJSONOpts{
			ExplicitFormat: explicitFormat,
			Format:         format,
			RawOutput:      cmd.Root().Bool("raw-output"),
			Title:          "beta:organization:spend-limits list",
			Transform:      transform,
		})
	} else {
		iter := client.Beta.Organization.SpendLimits.ListAutoPaging(ctx, params, options...)
		maxItems := int64(-1)
		if cmd.IsSet("max-items") {
			maxItems = cmd.Value("max-items").(int64)
		}
		return ShowJSONIterator(iter, maxItems, ShowJSONOpts{
			ExplicitFormat: explicitFormat,
			Format:         format,
			RawOutput:      cmd.Root().Bool("raw-output"),
			Title:          "beta:organization:spend-limits list",
			Transform:      transform,
		})
	}
}

func handleBetaOrganizationSpendLimitsDelete(ctx context.Context, cmd *cli.Command) error {
	client := anthropic.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("spend-limit-id") && len(unusedArgs) > 0 {
		cmd.Set("spend-limit-id", unusedArgs[0])
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

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Beta.Organization.SpendLimits.Delete(ctx, cmd.Value("spend-limit-id").(string), options...)
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
		Title:          "beta:organization:spend-limits delete",
		Transform:      transform,
	})
}

func handleBetaOrganizationSpendLimitsSet(ctx context.Context, cmd *cli.Command) error {
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

	params := anthropic.BetaOrganizationSpendLimitSetParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Beta.Organization.SpendLimits.Set(ctx, params, options...)
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
		Title:          "beta:organization:spend-limits set",
		Transform:      transform,
	})
}
