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

var organizationSpendLimitsIncreaseRequestsRetrieve = cli.Command{
	Name:    "retrieve",
	Usage:   "Retrieve a spend limit increase request.",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:        "spend-limit-increase-request-id",
			Usage:       "ID of the spend limit increase request.",
			Required:    true,
			PathParam:   "spend_limit_increase_request_id",
			DataAliases: []string{"id"},
		},
	},
	Action:          handleOrganizationSpendLimitsIncreaseRequestsRetrieve,
	HideHelpCommand: true,
}

var organizationSpendLimitsIncreaseRequestsList = cli.Command{
	Name:    "list",
	Usage:   "List spend limit increase requests, most recent first.",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[[]string]{
			Name:      "actor-id",
			Usage:     "Filter by requester, as `user_...` tagged IDs.",
			QueryPath: "actor_ids",
		},
		&requestflag.Flag[int64]{
			Name:      "limit",
			Default:   20,
			QueryPath: "limit",
		},
		&requestflag.Flag[string]{
			Name:      "page",
			Usage:     "Opaque cursor from a previous response's `next_page`.",
			QueryPath: "page",
		},
		&requestflag.Flag[[]string]{
			Name:      "status",
			Usage:     "Filter by status. Omit to return all.",
			QueryPath: "status",
		},
		&requestflag.Flag[int64]{
			Name:  "max-items",
			Usage: "The maximum number of items to return (use -1 for unlimited).",
		},
	},
	Action:          handleOrganizationSpendLimitsIncreaseRequestsList,
	HideHelpCommand: true,
}

var organizationSpendLimitsIncreaseRequestsApprove = cli.Command{
	Name:    "approve",
	Usage:   "Approve a pending spend limit increase request.",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:        "spend-limit-increase-request-id",
			Usage:       "ID of the spend limit increase request.",
			Required:    true,
			PathParam:   "spend_limit_increase_request_id",
			DataAliases: []string{"id"},
		},
		&requestflag.Flag[string]{
			Name:     "amount",
			Usage:    "New per-user spend limit as a non-negative integer decimal string (minor units).",
			Required: true,
			BodyPath: "amount",
		},
		&requestflag.Flag[*string]{
			Name:     "period",
			Usage:    `Allowed values: "daily", "monthly", "weekly".`,
			BodyPath: "period",
		},
		&requestflag.Flag[bool]{
			Name:     "suppress-notification",
			BodyPath: "suppress_notification",
		},
	},
	Action:          handleOrganizationSpendLimitsIncreaseRequestsApprove,
	HideHelpCommand: true,
}

var organizationSpendLimitsIncreaseRequestsDeny = cli.Command{
	Name:    "deny",
	Usage:   "Deny a pending spend limit increase request.",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:        "spend-limit-increase-request-id",
			Usage:       "ID of the spend limit increase request.",
			Required:    true,
			PathParam:   "spend_limit_increase_request_id",
			DataAliases: []string{"id"},
		},
		&requestflag.Flag[bool]{
			Name:     "suppress-notification",
			BodyPath: "suppress_notification",
		},
	},
	Action:          handleOrganizationSpendLimitsIncreaseRequestsDeny,
	HideHelpCommand: true,
}

func handleOrganizationSpendLimitsIncreaseRequestsRetrieve(ctx context.Context, cmd *cli.Command) error {
	client := anthropic.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("spend-limit-increase-request-id") && len(unusedArgs) > 0 {
		cmd.Set("spend-limit-increase-request-id", unusedArgs[0])
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
	_, err = client.Organization.SpendLimits.IncreaseRequests.Get(ctx, cmd.Value("spend-limit-increase-request-id").(string), options...)
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
		Title:          "organization:spend-limits:increase-requests retrieve",
		Transform:      transform,
	})
}

func handleOrganizationSpendLimitsIncreaseRequestsList(ctx context.Context, cmd *cli.Command) error {
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

	params := anthropic.OrganizationSpendLimitIncreaseRequestListParams{}

	format := "explore"
	explicitFormat := cmd.Root().IsSet("format")
	if explicitFormat {
		format = cmd.Root().String("format")
	}
	transform := cmd.Root().String("transform")
	if format == "raw" {
		var res []byte
		options = append(options, option.WithResponseBodyInto(&res))
		_, err = client.Organization.SpendLimits.IncreaseRequests.List(ctx, params, options...)
		if err != nil {
			return err
		}
		obj := gjson.ParseBytes(res)
		return ShowJSON(obj, ShowJSONOpts{
			ExplicitFormat: explicitFormat,
			Format:         format,
			RawOutput:      cmd.Root().Bool("raw-output"),
			Title:          "organization:spend-limits:increase-requests list",
			Transform:      transform,
		})
	} else {
		iter := client.Organization.SpendLimits.IncreaseRequests.ListAutoPaging(ctx, params, options...)
		maxItems := int64(-1)
		if cmd.IsSet("max-items") {
			maxItems = cmd.Value("max-items").(int64)
		}
		return ShowJSONIterator(iter, maxItems, ShowJSONOpts{
			ExplicitFormat: explicitFormat,
			Format:         format,
			RawOutput:      cmd.Root().Bool("raw-output"),
			Title:          "organization:spend-limits:increase-requests list",
			Transform:      transform,
		})
	}
}

func handleOrganizationSpendLimitsIncreaseRequestsApprove(ctx context.Context, cmd *cli.Command) error {
	client := anthropic.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("spend-limit-increase-request-id") && len(unusedArgs) > 0 {
		cmd.Set("spend-limit-increase-request-id", unusedArgs[0])
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

	params := anthropic.OrganizationSpendLimitIncreaseRequestApproveParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Organization.SpendLimits.IncreaseRequests.Approve(
		ctx,
		cmd.Value("spend-limit-increase-request-id").(string),
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
		Title:          "organization:spend-limits:increase-requests approve",
		Transform:      transform,
	})
}

func handleOrganizationSpendLimitsIncreaseRequestsDeny(ctx context.Context, cmd *cli.Command) error {
	client := anthropic.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("spend-limit-increase-request-id") && len(unusedArgs) > 0 {
		cmd.Set("spend-limit-increase-request-id", unusedArgs[0])
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

	params := anthropic.OrganizationSpendLimitIncreaseRequestDenyParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Organization.SpendLimits.IncreaseRequests.Deny(
		ctx,
		cmd.Value("spend-limit-increase-request-id").(string),
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
		Title:          "organization:spend-limits:increase-requests deny",
		Transform:      transform,
	})
}
