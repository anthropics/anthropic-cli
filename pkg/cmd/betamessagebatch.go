// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cmd

import (
	"context"
	"fmt"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/stainless-sdks/anthropic-cli/internal/apiquery"
	"github.com/stainless-sdks/anthropic-cli/internal/requestflag"
	"github.com/tidwall/gjson"
	"github.com/urfave/cli/v3"
)

var betaMessagesBatchesCreate = cli.Command{
	Name:  "create",
	Usage: "Send a batch of Message creation requests.",
	Flags: []cli.Flag{
		&requestflag.YAMLSliceFlag{
			Name:  "request",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: requestflag.RequestConfig{
				BodyPath: "requests",
			},
		},
		&requestflag.YAMLSliceFlag{
			Name:  "beta",
			Usage: "Optional header to specify the beta version(s) you want to use.",
			Config: requestflag.RequestConfig{
				HeaderPath: "anthropic-beta",
			},
		},
	},
	Action:          handleBetaMessagesBatchesCreate,
	HideHelpCommand: true,
}

var betaMessagesBatchesRetrieve = cli.Command{
	Name:  "retrieve",
	Usage: "This endpoint is idempotent and can be used to poll for Message Batch\ncompletion. To access the results of a Message Batch, make a request to the\n`results_url` field in the response.",
	Flags: []cli.Flag{
		&requestflag.StringFlag{
			Name:  "message-batch-id",
			Usage: "ID of the Message Batch.",
		},
		&requestflag.YAMLSliceFlag{
			Name:  "beta",
			Usage: "Optional header to specify the beta version(s) you want to use.",
			Config: requestflag.RequestConfig{
				HeaderPath: "anthropic-beta",
			},
		},
	},
	Action:          handleBetaMessagesBatchesRetrieve,
	HideHelpCommand: true,
}

var betaMessagesBatchesList = cli.Command{
	Name:  "list",
	Usage: "List all Message Batches within a Workspace. Most recently created batches are\nreturned first.",
	Flags: []cli.Flag{
		&requestflag.StringFlag{
			Name:  "after-id",
			Usage: "ID of the object to use as a cursor for pagination. When provided, returns the page of results immediately after this object.",
			Config: requestflag.RequestConfig{
				QueryPath: "after_id",
			},
		},
		&requestflag.StringFlag{
			Name:  "before-id",
			Usage: "ID of the object to use as a cursor for pagination. When provided, returns the page of results immediately before this object.",
			Config: requestflag.RequestConfig{
				QueryPath: "before_id",
			},
		},
		&requestflag.IntFlag{
			Name:  "limit",
			Usage: "Number of items to return per page.\n\nDefaults to `20`. Ranges from `1` to `1000`.",
			Value: requestflag.Value[int64](20),
			Config: requestflag.RequestConfig{
				QueryPath: "limit",
			},
		},
		&requestflag.YAMLSliceFlag{
			Name:  "beta",
			Usage: "Optional header to specify the beta version(s) you want to use.",
			Config: requestflag.RequestConfig{
				HeaderPath: "anthropic-beta",
			},
		},
	},
	Action:          handleBetaMessagesBatchesList,
	HideHelpCommand: true,
}

var betaMessagesBatchesDelete = cli.Command{
	Name:  "delete",
	Usage: "Delete a Message Batch.",
	Flags: []cli.Flag{
		&requestflag.StringFlag{
			Name:  "message-batch-id",
			Usage: "ID of the Message Batch.",
		},
		&requestflag.YAMLSliceFlag{
			Name:  "beta",
			Usage: "Optional header to specify the beta version(s) you want to use.",
			Config: requestflag.RequestConfig{
				HeaderPath: "anthropic-beta",
			},
		},
	},
	Action:          handleBetaMessagesBatchesDelete,
	HideHelpCommand: true,
}

var betaMessagesBatchesCancel = cli.Command{
	Name:  "cancel",
	Usage: "Batches may be canceled any time before processing ends. Once cancellation is\ninitiated, the batch enters a `canceling` state, at which time the system may\ncomplete any in-progress, non-interruptible requests before finalizing\ncancellation.",
	Flags: []cli.Flag{
		&requestflag.StringFlag{
			Name:  "message-batch-id",
			Usage: "ID of the Message Batch.",
		},
		&requestflag.YAMLSliceFlag{
			Name:  "beta",
			Usage: "Optional header to specify the beta version(s) you want to use.",
			Config: requestflag.RequestConfig{
				HeaderPath: "anthropic-beta",
			},
		},
	},
	Action:          handleBetaMessagesBatchesCancel,
	HideHelpCommand: true,
}

var betaMessagesBatchesResults = cli.Command{
	Name:  "results",
	Usage: "Streams the results of a Message Batch as a `.jsonl` file.",
	Flags: []cli.Flag{
		&requestflag.StringFlag{
			Name:  "message-batch-id",
			Usage: "ID of the Message Batch.",
		},
		&requestflag.YAMLSliceFlag{
			Name:  "beta",
			Usage: "Optional header to specify the beta version(s) you want to use.",
			Config: requestflag.RequestConfig{
				HeaderPath: "anthropic-beta",
			},
		},
	},
	Action:          handleBetaMessagesBatchesResults,
	HideHelpCommand: true,
}

func handleBetaMessagesBatchesCreate(ctx context.Context, cmd *cli.Command) error {
	client := anthropic.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if len(unusedArgs) > 0 {
		return fmt.Errorf("Unexpected extra arguments: %v", unusedArgs)
	}
	params := anthropic.BetaMessageBatchNewParams{}

	options, err := flagOptions(
		cmd,
		apiquery.NestedQueryFormatBrackets,
		apiquery.ArrayQueryFormatComma,
		ApplicationJSON,
	)
	if err != nil {
		return err
	}
	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Beta.Messages.Batches.New(
		ctx,
		params,
		options...,
	)
	if err != nil {
		return err
	}

	json := gjson.Parse(string(res))
	format := cmd.Root().String("format")
	transform := cmd.Root().String("transform")
	return ShowJSON("beta:messages:batches create", json, format, transform)
}

func handleBetaMessagesBatchesRetrieve(ctx context.Context, cmd *cli.Command) error {
	client := anthropic.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("message-batch-id") && len(unusedArgs) > 0 {
		cmd.Set("message-batch-id", unusedArgs[0])
		unusedArgs = unusedArgs[1:]
	}
	if len(unusedArgs) > 0 {
		return fmt.Errorf("Unexpected extra arguments: %v", unusedArgs)
	}
	params := anthropic.BetaMessageBatchGetParams{}

	options, err := flagOptions(
		cmd,
		apiquery.NestedQueryFormatBrackets,
		apiquery.ArrayQueryFormatComma,
		ApplicationJSON,
	)
	if err != nil {
		return err
	}
	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Beta.Messages.Batches.Get(
		ctx,
		requestflag.CommandRequestValue[string](cmd, "message-batch-id"),
		params,
		options...,
	)
	if err != nil {
		return err
	}

	json := gjson.Parse(string(res))
	format := cmd.Root().String("format")
	transform := cmd.Root().String("transform")
	return ShowJSON("beta:messages:batches retrieve", json, format, transform)
}

func handleBetaMessagesBatchesList(ctx context.Context, cmd *cli.Command) error {
	client := anthropic.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if len(unusedArgs) > 0 {
		return fmt.Errorf("Unexpected extra arguments: %v", unusedArgs)
	}
	params := anthropic.BetaMessageBatchListParams{}

	options, err := flagOptions(
		cmd,
		apiquery.NestedQueryFormatBrackets,
		apiquery.ArrayQueryFormatComma,
		ApplicationJSON,
	)
	if err != nil {
		return err
	}
	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Beta.Messages.Batches.List(
		ctx,
		params,
		options...,
	)
	if err != nil {
		return err
	}

	json := gjson.Parse(string(res))
	format := cmd.Root().String("format")
	transform := cmd.Root().String("transform")
	return ShowJSON("beta:messages:batches list", json, format, transform)
}

func handleBetaMessagesBatchesDelete(ctx context.Context, cmd *cli.Command) error {
	client := anthropic.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("message-batch-id") && len(unusedArgs) > 0 {
		cmd.Set("message-batch-id", unusedArgs[0])
		unusedArgs = unusedArgs[1:]
	}
	if len(unusedArgs) > 0 {
		return fmt.Errorf("Unexpected extra arguments: %v", unusedArgs)
	}
	params := anthropic.BetaMessageBatchDeleteParams{}

	options, err := flagOptions(
		cmd,
		apiquery.NestedQueryFormatBrackets,
		apiquery.ArrayQueryFormatComma,
		ApplicationJSON,
	)
	if err != nil {
		return err
	}
	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Beta.Messages.Batches.Delete(
		ctx,
		requestflag.CommandRequestValue[string](cmd, "message-batch-id"),
		params,
		options...,
	)
	if err != nil {
		return err
	}

	json := gjson.Parse(string(res))
	format := cmd.Root().String("format")
	transform := cmd.Root().String("transform")
	return ShowJSON("beta:messages:batches delete", json, format, transform)
}

func handleBetaMessagesBatchesCancel(ctx context.Context, cmd *cli.Command) error {
	client := anthropic.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("message-batch-id") && len(unusedArgs) > 0 {
		cmd.Set("message-batch-id", unusedArgs[0])
		unusedArgs = unusedArgs[1:]
	}
	if len(unusedArgs) > 0 {
		return fmt.Errorf("Unexpected extra arguments: %v", unusedArgs)
	}
	params := anthropic.BetaMessageBatchCancelParams{}

	options, err := flagOptions(
		cmd,
		apiquery.NestedQueryFormatBrackets,
		apiquery.ArrayQueryFormatComma,
		ApplicationJSON,
	)
	if err != nil {
		return err
	}
	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Beta.Messages.Batches.Cancel(
		ctx,
		requestflag.CommandRequestValue[string](cmd, "message-batch-id"),
		params,
		options...,
	)
	if err != nil {
		return err
	}

	json := gjson.Parse(string(res))
	format := cmd.Root().String("format")
	transform := cmd.Root().String("transform")
	return ShowJSON("beta:messages:batches cancel", json, format, transform)
}

func handleBetaMessagesBatchesResults(ctx context.Context, cmd *cli.Command) error {
	client := anthropic.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("message-batch-id") && len(unusedArgs) > 0 {
		cmd.Set("message-batch-id", unusedArgs[0])
		unusedArgs = unusedArgs[1:]
	}
	if len(unusedArgs) > 0 {
		return fmt.Errorf("Unexpected extra arguments: %v", unusedArgs)
	}
	params := anthropic.BetaMessageBatchResultsParams{}

	options, err := flagOptions(
		cmd,
		apiquery.NestedQueryFormatBrackets,
		apiquery.ArrayQueryFormatComma,
		ApplicationJSON,
	)
	if err != nil {
		return err
	}
	stream := client.Beta.Messages.Batches.ResultsStreaming(
		ctx,
		requestflag.CommandRequestValue[string](cmd, "message-batch-id"),
		params,
		options...,
	)
	for stream.Next() {
		fmt.Printf("%s\n", stream.Current().RawJSON())
	}
	return stream.Err()
}
