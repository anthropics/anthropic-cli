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

var betaEnvironmentsWorkRetrieve = cli.Command{
	Name:    "retrieve",
	Usage:   "Note: these endpoints are called automatically by the pre-built environment\nworker provided in the SDKs and CLI, for orchestrating sessions with self-hosted\nsandbox environments. They are included here as a reference; you do not need to\ninvoke them directly.",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "environment-id",
			Required:  true,
			PathParam: "environment_id",
		},
		&requestflag.Flag[string]{
			Name:      "work-id",
			Required:  true,
			PathParam: "work_id",
		},
		&requestflag.Flag[[]string]{
			Name:       "beta",
			Usage:      "Optional header to specify the beta version(s) you want to use.",
			HeaderPath: "anthropic-beta",
		},
	},
	Action:          handleBetaEnvironmentsWorkRetrieve,
	HideHelpCommand: true,
}

var betaEnvironmentsWorkUpdate = cli.Command{
	Name:    "update",
	Usage:   "Note: these endpoints are called automatically by the pre-built environment\nworker provided in the SDKs and CLI, for orchestrating sessions with self-hosted\nsandbox environments. They are included here as a reference; you do not need to\ninvoke them directly.",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "environment-id",
			Required:  true,
			PathParam: "environment_id",
		},
		&requestflag.Flag[string]{
			Name:      "work-id",
			Required:  true,
			PathParam: "work_id",
		},
		&requestflag.Flag[map[string]any]{
			Name:     "metadata",
			Usage:    "Metadata patch. Set a key to a string to upsert it, or to null to delete it. Omit the field to preserve existing metadata.",
			Required: true,
			BodyPath: "metadata",
		},
		&requestflag.Flag[[]string]{
			Name:       "beta",
			Usage:      "Optional header to specify the beta version(s) you want to use.",
			HeaderPath: "anthropic-beta",
		},
	},
	Action:          handleBetaEnvironmentsWorkUpdate,
	HideHelpCommand: true,
}

var betaEnvironmentsWorkList = cli.Command{
	Name:    "list",
	Usage:   "Note: these endpoints are called automatically by the pre-built environment\nworker provided in the SDKs and CLI, for orchestrating sessions with self-hosted\nsandbox environments. They are included here as a reference; you do not need to\ninvoke them directly.",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "environment-id",
			Required:  true,
			PathParam: "environment_id",
		},
		&requestflag.Flag[int64]{
			Name:      "limit",
			Usage:     "Maximum number of work items to return",
			Default:   20,
			QueryPath: "limit",
		},
		&requestflag.Flag[*string]{
			Name:      "page",
			Usage:     "Opaque cursor from previous response for pagination",
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
	Action:          handleBetaEnvironmentsWorkList,
	HideHelpCommand: true,
}

var betaEnvironmentsWorkAck = cli.Command{
	Name:    "ack",
	Usage:   "Note: these endpoints are called automatically by the pre-built environment\nworker provided in the SDKs and CLI, for orchestrating sessions with self-hosted\nsandbox environments. They are included here as a reference; you do not need to\ninvoke them directly.",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "environment-id",
			Required:  true,
			PathParam: "environment_id",
		},
		&requestflag.Flag[string]{
			Name:      "work-id",
			Required:  true,
			PathParam: "work_id",
		},
		&requestflag.Flag[[]string]{
			Name:       "beta",
			Usage:      "Optional header to specify the beta version(s) you want to use.",
			HeaderPath: "anthropic-beta",
		},
	},
	Action:          handleBetaEnvironmentsWorkAck,
	HideHelpCommand: true,
}

var betaEnvironmentsWorkHeartbeat = cli.Command{
	Name:    "heartbeat",
	Usage:   "Note: these endpoints are called automatically by the pre-built environment\nworker provided in the SDKs and CLI, for orchestrating sessions with self-hosted\nsandbox environments. They are included here as a reference; you do not need to\ninvoke them directly.",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "environment-id",
			Required:  true,
			PathParam: "environment_id",
		},
		&requestflag.Flag[string]{
			Name:      "work-id",
			Required:  true,
			PathParam: "work_id",
		},
		&requestflag.Flag[*int64]{
			Name:      "desired-ttl-seconds",
			Usage:     "Desired TTL in seconds",
			QueryPath: "desired_ttl_seconds",
		},
		&requestflag.Flag[*string]{
			Name:      "expected-last-heartbeat",
			Usage:     "Expected last_heartbeat for conditional update (optimistic concurrency). Use literal 'NO_HEARTBEAT' to claim an unclaimed lease (first heartbeat). For subsequent heartbeats, echo the server's previous last_heartbeat value exactly. Returns 412 Precondition Failed if the actual value doesn't match.",
			QueryPath: "expected_last_heartbeat",
		},
		&requestflag.Flag[[]string]{
			Name:       "beta",
			Usage:      "Optional header to specify the beta version(s) you want to use.",
			HeaderPath: "anthropic-beta",
		},
	},
	Action:          handleBetaEnvironmentsWorkHeartbeat,
	HideHelpCommand: true,
}

var betaEnvironmentsWorkPoll = cli.Command{
	Name:    "poll",
	Usage:   "Note: these endpoints are called automatically by the pre-built environment\nworker provided in the SDKs and CLI, for orchestrating sessions with self-hosted\nsandbox environments. They are included here as a reference; you do not need to\ninvoke them directly.",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "environment-id",
			Required:  true,
			PathParam: "environment_id",
		},
		&requestflag.Flag[*int64]{
			Name:      "block-ms",
			Usage:     "How long to wait for work to arrive before returning. Must be 1-999 in milliseconds. Defaults to non-blocking (returns immediately if no work is available).",
			QueryPath: "block_ms",
		},
		&requestflag.Flag[*int64]{
			Name:      "reclaim-older-than-ms",
			Usage:     "Reclaim unacknowledged work items older than this many milliseconds. If omitted, uses the default (5000ms).",
			QueryPath: "reclaim_older_than_ms",
		},
		&requestflag.Flag[[]string]{
			Name:       "beta",
			Usage:      "Optional header to specify the beta version(s) you want to use.",
			HeaderPath: "anthropic-beta",
		},
		&requestflag.Flag[string]{
			Name:       "anthropic-worker-id",
			Usage:      "Unique identifier for the specific worker polling, used to track aggregated environment-level work metrics in Console",
			HeaderPath: "Anthropic-Worker-ID",
		},
	},
	Action:          handleBetaEnvironmentsWorkPoll,
	HideHelpCommand: true,
}

var betaEnvironmentsWorkStats = cli.Command{
	Name:    "stats",
	Usage:   "Get statistics about the work queue for an environment.",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "environment-id",
			Required:  true,
			PathParam: "environment_id",
		},
		&requestflag.Flag[[]string]{
			Name:       "beta",
			Usage:      "Optional header to specify the beta version(s) you want to use.",
			HeaderPath: "anthropic-beta",
		},
	},
	Action:          handleBetaEnvironmentsWorkStats,
	HideHelpCommand: true,
}

var betaEnvironmentsWorkStop = cli.Command{
	Name:    "stop",
	Usage:   "Note: these endpoints are called automatically by the pre-built environment\nworker provided in the SDKs and CLI, for orchestrating sessions with self-hosted\nsandbox environments. They are included here as a reference; you do not need to\ninvoke them directly.",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "environment-id",
			Required:  true,
			PathParam: "environment_id",
		},
		&requestflag.Flag[string]{
			Name:      "work-id",
			Required:  true,
			PathParam: "work_id",
		},
		&requestflag.Flag[bool]{
			Name:     "force",
			Usage:    "If true, immediately stop work without graceful shutdown",
			Default:  false,
			BodyPath: "force",
		},
		&requestflag.Flag[[]string]{
			Name:       "beta",
			Usage:      "Optional header to specify the beta version(s) you want to use.",
			HeaderPath: "anthropic-beta",
		},
	},
	Action:          handleBetaEnvironmentsWorkStop,
	HideHelpCommand: true,
}

func handleBetaEnvironmentsWorkRetrieve(ctx context.Context, cmd *cli.Command) error {
	client := anthropic.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("work-id") && len(unusedArgs) > 0 {
		cmd.Set("work-id", unusedArgs[0])
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

	params := anthropic.BetaEnvironmentWorkGetParams{
		EnvironmentID: cmd.Value("environment-id").(string),
	}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Beta.Environments.Work.Get(
		ctx,
		cmd.Value("work-id").(string),
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
		Title:          "beta:environments:work retrieve",
		Transform:      transform,
	})
}

func handleBetaEnvironmentsWorkUpdate(ctx context.Context, cmd *cli.Command) error {
	client := anthropic.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("work-id") && len(unusedArgs) > 0 {
		cmd.Set("work-id", unusedArgs[0])
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

	params := anthropic.BetaEnvironmentWorkUpdateParams{
		EnvironmentID: cmd.Value("environment-id").(string),
	}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Beta.Environments.Work.Update(
		ctx,
		cmd.Value("work-id").(string),
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
		Title:          "beta:environments:work update",
		Transform:      transform,
	})
}

func handleBetaEnvironmentsWorkList(ctx context.Context, cmd *cli.Command) error {
	client := anthropic.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("environment-id") && len(unusedArgs) > 0 {
		cmd.Set("environment-id", unusedArgs[0])
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

	params := anthropic.BetaEnvironmentWorkListParams{}

	format := "explore"
	explicitFormat := cmd.Root().IsSet("format")
	if explicitFormat {
		format = cmd.Root().String("format")
	}
	transform := cmd.Root().String("transform")
	if format == "raw" {
		var res []byte
		options = append(options, option.WithResponseBodyInto(&res))
		_, err = client.Beta.Environments.Work.List(
			ctx,
			cmd.Value("environment-id").(string),
			params,
			options...,
		)
		if err != nil {
			return err
		}
		obj := gjson.ParseBytes(res)
		return ShowJSON(obj, ShowJSONOpts{
			ExplicitFormat: explicitFormat,
			Format:         format,
			RawOutput:      cmd.Root().Bool("raw-output"),
			Title:          "beta:environments:work list",
			Transform:      transform,
		})
	} else {
		iter := client.Beta.Environments.Work.ListAutoPaging(
			ctx,
			cmd.Value("environment-id").(string),
			params,
			options...,
		)
		maxItems := int64(-1)
		if cmd.IsSet("max-items") {
			maxItems = cmd.Value("max-items").(int64)
		}
		return ShowJSONIterator(iter, maxItems, ShowJSONOpts{
			ExplicitFormat: explicitFormat,
			Format:         format,
			RawOutput:      cmd.Root().Bool("raw-output"),
			Title:          "beta:environments:work list",
			Transform:      transform,
		})
	}
}

func handleBetaEnvironmentsWorkAck(ctx context.Context, cmd *cli.Command) error {
	client := anthropic.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("work-id") && len(unusedArgs) > 0 {
		cmd.Set("work-id", unusedArgs[0])
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

	params := anthropic.BetaEnvironmentWorkAckParams{
		EnvironmentID: cmd.Value("environment-id").(string),
	}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Beta.Environments.Work.Ack(
		ctx,
		cmd.Value("work-id").(string),
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
		Title:          "beta:environments:work ack",
		Transform:      transform,
	})
}

func handleBetaEnvironmentsWorkHeartbeat(ctx context.Context, cmd *cli.Command) error {
	client := anthropic.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("work-id") && len(unusedArgs) > 0 {
		cmd.Set("work-id", unusedArgs[0])
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

	params := anthropic.BetaEnvironmentWorkHeartbeatParams{
		EnvironmentID: cmd.Value("environment-id").(string),
	}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Beta.Environments.Work.Heartbeat(
		ctx,
		cmd.Value("work-id").(string),
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
		Title:          "beta:environments:work heartbeat",
		Transform:      transform,
	})
}

func handleBetaEnvironmentsWorkPoll(ctx context.Context, cmd *cli.Command) error {
	client := anthropic.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("environment-id") && len(unusedArgs) > 0 {
		cmd.Set("environment-id", unusedArgs[0])
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

	params := anthropic.BetaEnvironmentWorkPollParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Beta.Environments.Work.Poll(
		ctx,
		cmd.Value("environment-id").(string),
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
		Title:          "beta:environments:work poll",
		Transform:      transform,
	})
}

func handleBetaEnvironmentsWorkStats(ctx context.Context, cmd *cli.Command) error {
	client := anthropic.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("environment-id") && len(unusedArgs) > 0 {
		cmd.Set("environment-id", unusedArgs[0])
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

	params := anthropic.BetaEnvironmentWorkStatsParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Beta.Environments.Work.Stats(
		ctx,
		cmd.Value("environment-id").(string),
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
		Title:          "beta:environments:work stats",
		Transform:      transform,
	})
}

func handleBetaEnvironmentsWorkStop(ctx context.Context, cmd *cli.Command) error {
	client := anthropic.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("work-id") && len(unusedArgs) > 0 {
		cmd.Set("work-id", unusedArgs[0])
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

	params := anthropic.BetaEnvironmentWorkStopParams{
		EnvironmentID: cmd.Value("environment-id").(string),
	}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Beta.Environments.Work.Stop(
		ctx,
		cmd.Value("work-id").(string),
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
		Title:          "beta:environments:work stop",
		Transform:      transform,
	})
}
