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

var betaDeploymentsCreate = requestflag.WithInnerFlags(cli.Command{
	Name:    "create",
	Usage:   "Create Deployment",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[any]{
			Name:     "agent",
			Usage:    "Agent to deploy. Accepts the `agent` ID string, which pins the latest version, or an `agent` object with both id and version specified. The agent must exist and not be archived.",
			Required: true,
			BodyPath: "agent",
		},
		&requestflag.Flag[string]{
			Name:     "environment-id",
			Usage:    "ID of the `environment` defining the container configuration for sessions created from this deployment.",
			Required: true,
			BodyPath: "environment_id",
		},
		&requestflag.Flag[[]map[string]any]{
			Name:     "initial-event",
			Usage:    "Events to send to each session immediately after creation. At least 1, maximum 50.",
			Required: true,
			BodyPath: "initial_events",
		},
		&requestflag.Flag[string]{
			Name:     "name",
			Usage:    "Human-readable name for the deployment.",
			Required: true,
			BodyPath: "name",
		},
		&requestflag.Flag[*string]{
			Name:     "description",
			Usage:    "Description of what the deployment does.",
			BodyPath: "description",
		},
		&requestflag.Flag[map[string]any]{
			Name:     "metadata",
			Usage:    "Arbitrary key-value metadata. Maximum 16 pairs, keys up to 64 chars, values up to 512 chars.",
			BodyPath: "metadata",
		},
		&requestflag.Flag[[]map[string]any]{
			Name:     "resource",
			Usage:    "Resources (e.g. repositories, files) to mount into each session's container. Maximum 500.",
			BodyPath: "resources",
		},
		&requestflag.Flag[map[string]any]{
			Name:     "schedule",
			Usage:    "5-field POSIX cron schedule. Literal wall-clock matching in the configured timezone.",
			BodyPath: "schedule",
		},
		&requestflag.Flag[[]string]{
			Name:     "vault-id",
			Usage:    "Vault IDs for stored credentials the agent can use during sessions created from this deployment. Maximum 50.",
			BodyPath: "vault_ids",
		},
		&requestflag.Flag[[]string]{
			Name:       "beta",
			Usage:      "Optional header to specify the beta version(s) you want to use.",
			HeaderPath: "anthropic-beta",
		},
	},
	Action:          handleBetaDeploymentsCreate,
	HideHelpCommand: true,
}, map[string][]requestflag.HasOuterFlag{
	"schedule": {
		&requestflag.InnerFlag[string]{
			Name:       "schedule.expression",
			Usage:      `5-field POSIX cron expression: minute hour day-of-month month day-of-week (e.g., "0 9 * * 1-5" for weekdays at 9am). Day-of-week is 0-7 where 0 and 7 both mean Sunday. Extended cron syntax - seconds or year fields, and the special characters L, W, #, and ? - is not supported, nor are predefined shortcuts (@daily).`,
			InnerField: "expression",
		},
		&requestflag.InnerFlag[string]{
			Name:       "schedule.timezone",
			Usage:      `Required. IANA timezone identifier (e.g., "America/Los_Angeles", "UTC"). Validated against the IANA timezone database.`,
			InnerField: "timezone",
		},
		&requestflag.InnerFlag[string]{
			Name:       "schedule.type",
			Usage:      `Allowed values: "cron".`,
			InnerField: "type",
		},
	},
})

var betaDeploymentsRetrieve = cli.Command{
	Name:    "retrieve",
	Usage:   "Get Deployment",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "deployment-id",
			Required:  true,
			PathParam: "deployment_id",
		},
		&requestflag.Flag[[]string]{
			Name:       "beta",
			Usage:      "Optional header to specify the beta version(s) you want to use.",
			HeaderPath: "anthropic-beta",
		},
	},
	Action:          handleBetaDeploymentsRetrieve,
	HideHelpCommand: true,
}

var betaDeploymentsUpdate = requestflag.WithInnerFlags(cli.Command{
	Name:    "update",
	Usage:   "Update Deployment",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:        "deployment-id",
			Required:    true,
			PathParam:   "deployment_id",
			DataAliases: []string{"id"},
		},
		&requestflag.Flag[any]{
			Name:     "agent",
			Usage:    "Agent to deploy. Accepts the `agent` ID string, which re-pins to the latest version, or an `agent` object with both id and version specified. Omit to preserve. Cannot be cleared.",
			BodyPath: "agent",
		},
		&requestflag.Flag[*string]{
			Name:     "description",
			Usage:    "Description. Omit to preserve; send empty string or null to clear.",
			BodyPath: "description",
		},
		&requestflag.Flag[string]{
			Name:     "environment-id",
			Usage:    "ID of the `environment` where sessions run. Omit to preserve. Cannot be cleared.",
			BodyPath: "environment_id",
		},
		&requestflag.Flag[[]map[string]any]{
			Name:     "initial-event",
			Usage:    "Initial events. Full replacement. Omit to preserve. Cannot be cleared. At least 1, maximum 50.",
			BodyPath: "initial_events",
		},
		&requestflag.Flag[map[string]any]{
			Name:     "metadata",
			Usage:    "Metadata patch. Set a key to a string to upsert it, or to null to delete it. Omit the field to preserve. The stored bag is limited to 16 keys (up to 64 chars each) with values up to 512 chars.",
			BodyPath: "metadata",
		},
		&requestflag.Flag[string]{
			Name:     "name",
			Usage:    "Human-readable name. Must be non-empty. Omit to preserve. Cannot be cleared.",
			BodyPath: "name",
		},
		&requestflag.Flag[any]{
			Name:     "resource",
			Usage:    "Session resources. Full replacement. Omit to preserve; send empty array or null to clear. Maximum 500.",
			BodyPath: "resources",
		},
		&requestflag.Flag[map[string]any]{
			Name:     "schedule",
			Usage:    "5-field POSIX cron schedule. Literal wall-clock matching in the configured timezone.",
			BodyPath: "schedule",
		},
		&requestflag.Flag[any]{
			Name:     "vault-id",
			Usage:    "Vault IDs. Full replacement. Omit to preserve; send empty array or null to clear. Maximum 50.",
			BodyPath: "vault_ids",
		},
		&requestflag.Flag[[]string]{
			Name:       "beta",
			Usage:      "Optional header to specify the beta version(s) you want to use.",
			HeaderPath: "anthropic-beta",
		},
	},
	Action:          handleBetaDeploymentsUpdate,
	HideHelpCommand: true,
}, map[string][]requestflag.HasOuterFlag{
	"schedule": {
		&requestflag.InnerFlag[string]{
			Name:       "schedule.expression",
			Usage:      `5-field POSIX cron expression: minute hour day-of-month month day-of-week (e.g., "0 9 * * 1-5" for weekdays at 9am). Day-of-week is 0-7 where 0 and 7 both mean Sunday. Extended cron syntax - seconds or year fields, and the special characters L, W, #, and ? - is not supported, nor are predefined shortcuts (@daily).`,
			InnerField: "expression",
		},
		&requestflag.InnerFlag[string]{
			Name:       "schedule.timezone",
			Usage:      `Required. IANA timezone identifier (e.g., "America/Los_Angeles", "UTC"). Validated against the IANA timezone database.`,
			InnerField: "timezone",
		},
		&requestflag.InnerFlag[string]{
			Name:       "schedule.type",
			Usage:      `Allowed values: "cron".`,
			InnerField: "type",
		},
	},
})

var betaDeploymentsList = cli.Command{
	Name:    "list",
	Usage:   "List Deployments",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "agent-id",
			Usage:     "Filter by agent ID.",
			QueryPath: "agent_id",
		},
		&requestflag.Flag[any]{
			Name:      "created-at-gte",
			Usage:     "Return deployments created at or after this time (inclusive).",
			QueryPath: "created_at[gte]",
		},
		&requestflag.Flag[any]{
			Name:      "created-at-lte",
			Usage:     "Return deployments created at or before this time (inclusive).",
			QueryPath: "created_at[lte]",
		},
		&requestflag.Flag[bool]{
			Name:      "include-archived",
			Usage:     "When true, includes archived deployments. Default: false (exclude archived).",
			QueryPath: "include_archived",
		},
		&requestflag.Flag[int64]{
			Name:      "limit",
			Usage:     "Maximum results per page. Default 20, maximum 100.",
			QueryPath: "limit",
		},
		&requestflag.Flag[string]{
			Name:      "page",
			Usage:     "Opaque pagination cursor.",
			QueryPath: "page",
		},
		&requestflag.Flag[string]{
			Name:      "status",
			Usage:     "Lifecycle status of a deployment.",
			QueryPath: "status",
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
	Action:          handleBetaDeploymentsList,
	HideHelpCommand: true,
}

var betaDeploymentsArchive = cli.Command{
	Name:    "archive",
	Usage:   "Archive Deployment",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "deployment-id",
			Required:  true,
			PathParam: "deployment_id",
		},
		&requestflag.Flag[[]string]{
			Name:       "beta",
			Usage:      "Optional header to specify the beta version(s) you want to use.",
			HeaderPath: "anthropic-beta",
		},
	},
	Action:          handleBetaDeploymentsArchive,
	HideHelpCommand: true,
}

var betaDeploymentsPause = cli.Command{
	Name:    "pause",
	Usage:   "Pause Deployment",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "deployment-id",
			Required:  true,
			PathParam: "deployment_id",
		},
		&requestflag.Flag[[]string]{
			Name:       "beta",
			Usage:      "Optional header to specify the beta version(s) you want to use.",
			HeaderPath: "anthropic-beta",
		},
	},
	Action:          handleBetaDeploymentsPause,
	HideHelpCommand: true,
}

var betaDeploymentsRun = cli.Command{
	Name:    "run",
	Usage:   "Run Deployment Now",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "deployment-id",
			Required:  true,
			PathParam: "deployment_id",
		},
		&requestflag.Flag[[]string]{
			Name:       "beta",
			Usage:      "Optional header to specify the beta version(s) you want to use.",
			HeaderPath: "anthropic-beta",
		},
	},
	Action:          handleBetaDeploymentsRun,
	HideHelpCommand: true,
}

var betaDeploymentsUnpause = cli.Command{
	Name:    "unpause",
	Usage:   "Unpause Deployment",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "deployment-id",
			Required:  true,
			PathParam: "deployment_id",
		},
		&requestflag.Flag[[]string]{
			Name:       "beta",
			Usage:      "Optional header to specify the beta version(s) you want to use.",
			HeaderPath: "anthropic-beta",
		},
	},
	Action:          handleBetaDeploymentsUnpause,
	HideHelpCommand: true,
}

func handleBetaDeploymentsCreate(ctx context.Context, cmd *cli.Command) error {
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

	params := anthropic.BetaDeploymentNewParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Beta.Deployments.New(ctx, params, options...)
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
		Title:          "beta:deployments create",
		Transform:      transform,
	})
}

func handleBetaDeploymentsRetrieve(ctx context.Context, cmd *cli.Command) error {
	client := anthropic.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("deployment-id") && len(unusedArgs) > 0 {
		cmd.Set("deployment-id", unusedArgs[0])
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

	params := anthropic.BetaDeploymentGetParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Beta.Deployments.Get(
		ctx,
		cmd.Value("deployment-id").(string),
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
		Title:          "beta:deployments retrieve",
		Transform:      transform,
	})
}

func handleBetaDeploymentsUpdate(ctx context.Context, cmd *cli.Command) error {
	client := anthropic.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("deployment-id") && len(unusedArgs) > 0 {
		cmd.Set("deployment-id", unusedArgs[0])
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

	params := anthropic.BetaDeploymentUpdateParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Beta.Deployments.Update(
		ctx,
		cmd.Value("deployment-id").(string),
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
		Title:          "beta:deployments update",
		Transform:      transform,
	})
}

func handleBetaDeploymentsList(ctx context.Context, cmd *cli.Command) error {
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

	params := anthropic.BetaDeploymentListParams{}

	format := "explore"
	explicitFormat := cmd.Root().IsSet("format")
	if explicitFormat {
		format = cmd.Root().String("format")
	}
	transform := cmd.Root().String("transform")
	if format == "raw" {
		var res []byte
		options = append(options, option.WithResponseBodyInto(&res))
		_, err = client.Beta.Deployments.List(ctx, params, options...)
		if err != nil {
			return err
		}
		obj := gjson.ParseBytes(res)
		return ShowJSON(obj, ShowJSONOpts{
			ExplicitFormat: explicitFormat,
			Format:         format,
			RawOutput:      cmd.Root().Bool("raw-output"),
			Title:          "beta:deployments list",
			Transform:      transform,
		})
	} else {
		iter := client.Beta.Deployments.ListAutoPaging(ctx, params, options...)
		maxItems := int64(-1)
		if cmd.IsSet("max-items") {
			maxItems = cmd.Value("max-items").(int64)
		}
		return ShowJSONIterator(iter, maxItems, ShowJSONOpts{
			ExplicitFormat: explicitFormat,
			Format:         format,
			RawOutput:      cmd.Root().Bool("raw-output"),
			Title:          "beta:deployments list",
			Transform:      transform,
		})
	}
}

func handleBetaDeploymentsArchive(ctx context.Context, cmd *cli.Command) error {
	client := anthropic.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("deployment-id") && len(unusedArgs) > 0 {
		cmd.Set("deployment-id", unusedArgs[0])
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

	params := anthropic.BetaDeploymentArchiveParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Beta.Deployments.Archive(
		ctx,
		cmd.Value("deployment-id").(string),
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
		Title:          "beta:deployments archive",
		Transform:      transform,
	})
}

func handleBetaDeploymentsPause(ctx context.Context, cmd *cli.Command) error {
	client := anthropic.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("deployment-id") && len(unusedArgs) > 0 {
		cmd.Set("deployment-id", unusedArgs[0])
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

	params := anthropic.BetaDeploymentPauseParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Beta.Deployments.Pause(
		ctx,
		cmd.Value("deployment-id").(string),
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
		Title:          "beta:deployments pause",
		Transform:      transform,
	})
}

func handleBetaDeploymentsRun(ctx context.Context, cmd *cli.Command) error {
	client := anthropic.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("deployment-id") && len(unusedArgs) > 0 {
		cmd.Set("deployment-id", unusedArgs[0])
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

	params := anthropic.BetaDeploymentRunParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Beta.Deployments.Run(
		ctx,
		cmd.Value("deployment-id").(string),
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
		Title:          "beta:deployments run",
		Transform:      transform,
	})
}

func handleBetaDeploymentsUnpause(ctx context.Context, cmd *cli.Command) error {
	client := anthropic.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("deployment-id") && len(unusedArgs) > 0 {
		cmd.Set("deployment-id", unusedArgs[0])
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

	params := anthropic.BetaDeploymentUnpauseParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Beta.Deployments.Unpause(
		ctx,
		cmd.Value("deployment-id").(string),
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
		Title:          "beta:deployments unpause",
		Transform:      transform,
	})
}
