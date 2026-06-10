// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cmd

import (
	"context"
	"fmt"

	"github.com/anthropics/anthropic-cli/internal/apiquery"
	"github.com/anthropics/anthropic-cli/internal/requestflag"
	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/tidwall/gjson"
	"github.com/urfave/cli/v3"
)

var betaDeploymentRunsRetrieve = cli.Command{
	Name:    "retrieve",
	Usage:   "Get Deployment Run",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "deployment-run-id",
			Required:  true,
			PathParam: "deployment_run_id",
		},
		&requestflag.Flag[[]string]{
			Name:       "beta",
			Usage:      "Optional header to specify the beta version(s) you want to use.",
			HeaderPath: "anthropic-beta",
		},
	},
	Action:          handleBetaDeploymentRunsRetrieve,
	HideHelpCommand: true,
}

var betaDeploymentRunsList = cli.Command{
	Name:    "list",
	Usage:   "List Deployment Runs",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[any]{
			Name:      "created-at-gt",
			Usage:     "Return runs created strictly after this time (exclusive).",
			QueryPath: "created_at[gt]",
		},
		&requestflag.Flag[any]{
			Name:      "created-at-gte",
			Usage:     "Return runs created at or after this time (inclusive).",
			QueryPath: "created_at[gte]",
		},
		&requestflag.Flag[any]{
			Name:      "created-at-lt",
			Usage:     "Return runs created strictly before this time (exclusive).",
			QueryPath: "created_at[lt]",
		},
		&requestflag.Flag[any]{
			Name:      "created-at-lte",
			Usage:     "Return runs created at or before this time (inclusive).",
			QueryPath: "created_at[lte]",
		},
		&requestflag.Flag[string]{
			Name:      "deployment-id",
			Usage:     "Filter to a specific deployment. Omit to list across all deployments in the workspace. Filtering by a non-existent deployment_id returns 200 with empty data.",
			QueryPath: "deployment_id",
		},
		&requestflag.Flag[bool]{
			Name:      "has-error",
			Usage:     "Filter: true for runs with non-null error, false for runs with non-null session_id. Omit for all.",
			QueryPath: "has_error",
		},
		&requestflag.Flag[int64]{
			Name:      "limit",
			Usage:     "Maximum results per page. Default 20, maximum 1000.",
			QueryPath: "limit",
		},
		&requestflag.Flag[string]{
			Name:      "page",
			Usage:     "Opaque pagination cursor. Pass next_page from the previous response. Invalid or expired cursors return 400.",
			QueryPath: "page",
		},
		&requestflag.Flag[string]{
			Name:      "trigger-type",
			Usage:     "What triggered a deployment run.",
			QueryPath: "trigger_type",
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
	Action:          handleBetaDeploymentRunsList,
	HideHelpCommand: true,
}

func handleBetaDeploymentRunsRetrieve(ctx context.Context, cmd *cli.Command) error {
	client := anthropic.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("deployment-run-id") && len(unusedArgs) > 0 {
		cmd.Set("deployment-run-id", unusedArgs[0])
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

	params := anthropic.BetaDeploymentRunGetParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Beta.DeploymentRuns.Get(
		ctx,
		cmd.Value("deployment-run-id").(string),
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
		Title:          "beta:deployment-runs retrieve",
		Transform:      transform,
	})
}

func handleBetaDeploymentRunsList(ctx context.Context, cmd *cli.Command) error {
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

	params := anthropic.BetaDeploymentRunListParams{}

	format := "explore"
	explicitFormat := cmd.Root().IsSet("format")
	if explicitFormat {
		format = cmd.Root().String("format")
	}
	transform := cmd.Root().String("transform")
	if format == "raw" {
		var res []byte
		options = append(options, option.WithResponseBodyInto(&res))
		_, err = client.Beta.DeploymentRuns.List(ctx, params, options...)
		if err != nil {
			return err
		}
		obj := gjson.ParseBytes(res)
		return ShowJSON(obj, ShowJSONOpts{
			ExplicitFormat: explicitFormat,
			Format:         format,
			RawOutput:      cmd.Root().Bool("raw-output"),
			Title:          "beta:deployment-runs list",
			Transform:      transform,
		})
	} else {
		iter := client.Beta.DeploymentRuns.ListAutoPaging(ctx, params, options...)
		maxItems := int64(-1)
		if cmd.IsSet("max-items") {
			maxItems = cmd.Value("max-items").(int64)
		}
		return ShowJSONIterator(iter, maxItems, ShowJSONOpts{
			ExplicitFormat: explicitFormat,
			Format:         format,
			RawOutput:      cmd.Root().Bool("raw-output"),
			Title:          "beta:deployment-runs list",
			Transform:      transform,
		})
	}
}
