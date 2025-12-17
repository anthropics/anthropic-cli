// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/stainless-sdks/anthropic-cli/internal/apiquery"
	"github.com/stainless-sdks/anthropic-cli/internal/requestflag"
	"github.com/tidwall/gjson"
	"github.com/urfave/cli/v3"
)

var messagesBatchesCreate = cli.Command{
	Name:  "create",
	Usage: "Send a batch of Message creation requests.",
	Flags: []cli.Flag{
		&requestflag.Flag[[]any]{
			Name:     "request",
			Usage:    "List of requests for prompt completion. Each is an individual request to create a Message.",
			BodyPath: "requests",
		},
	},
	Action:          handleMessagesBatchesCreate,
	HideHelpCommand: true,
}

var messagesBatchesRetrieve = cli.Command{
	Name:  "retrieve",
	Usage: "This endpoint is idempotent and can be used to poll for Message Batch\ncompletion. To access the results of a Message Batch, make a request to the\n`results_url` field in the response.",
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:  "message-batch-id",
			Usage: "ID of the Message Batch.",
		},
	},
	Action:          handleMessagesBatchesRetrieve,
	HideHelpCommand: true,
}

var messagesBatchesList = cli.Command{
	Name:  "list",
	Usage: "List all Message Batches within a Workspace. Most recently created batches are\nreturned first.",
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "after-id",
			Usage:     "ID of the object to use as a cursor for pagination. When provided, returns the page of results immediately after this object.",
			QueryPath: "after_id",
		},
		&requestflag.Flag[string]{
			Name:      "before-id",
			Usage:     "ID of the object to use as a cursor for pagination. When provided, returns the page of results immediately before this object.",
			QueryPath: "before_id",
		},
		&requestflag.Flag[int64]{
			Name:      "limit",
			Usage:     "Number of items to return per page.\n\nDefaults to `20`. Ranges from `1` to `1000`.",
			Default:   20,
			QueryPath: "limit",
		},
	},
	Action:          handleMessagesBatchesList,
	HideHelpCommand: true,
}

var messagesBatchesDelete = cli.Command{
	Name:  "delete",
	Usage: "Delete a Message Batch.",
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:  "message-batch-id",
			Usage: "ID of the Message Batch.",
		},
	},
	Action:          handleMessagesBatchesDelete,
	HideHelpCommand: true,
}

var messagesBatchesCancel = cli.Command{
	Name:  "cancel",
	Usage: "Batches may be canceled any time before processing ends. Once cancellation is\ninitiated, the batch enters a `canceling` state, at which time the system may\ncomplete any in-progress, non-interruptible requests before finalizing\ncancellation.",
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:  "message-batch-id",
			Usage: "ID of the Message Batch.",
		},
	},
	Action:          handleMessagesBatchesCancel,
	HideHelpCommand: true,
}

var messagesBatchesResults = cli.Command{
	Name:  "results",
	Usage: "Streams the results of a Message Batch as a `.jsonl` file.",
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:  "message-batch-id",
			Usage: "ID of the Message Batch.",
		},
	},
	Action:          handleMessagesBatchesResults,
	HideHelpCommand: true,
}

func handleMessagesBatchesCreate(ctx context.Context, cmd *cli.Command) error {
	client := anthropic.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()

	if len(unusedArgs) > 0 {
		return fmt.Errorf("Unexpected extra arguments: %v", unusedArgs)
	}

	params := anthropic.MessageBatchNewParams{}

	options, err := flagOptions(
		cmd,
		apiquery.NestedQueryFormatBrackets,
		apiquery.ArrayQueryFormatComma,
		ApplicationJSON,
		false,
	)
	if err != nil {
		return err
	}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Messages.Batches.New(ctx, params, options...)
	if err != nil {
		return err
	}

	obj := gjson.ParseBytes(res)
	format := cmd.Root().String("format")
	transform := cmd.Root().String("transform")
	return ShowJSON(os.Stdout, "messages:batches create", obj, format, transform)
}

func handleMessagesBatchesRetrieve(ctx context.Context, cmd *cli.Command) error {
	client := anthropic.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("message-batch-id") && len(unusedArgs) > 0 {
		cmd.Set("message-batch-id", unusedArgs[0])
		unusedArgs = unusedArgs[1:]
	}
	if len(unusedArgs) > 0 {
		return fmt.Errorf("Unexpected extra arguments: %v", unusedArgs)
	}

	options, err := flagOptions(
		cmd,
		apiquery.NestedQueryFormatBrackets,
		apiquery.ArrayQueryFormatComma,
		EmptyBody,
		false,
	)
	if err != nil {
		return err
	}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Messages.Batches.Get(ctx, cmd.Value("message-batch-id").(string), options...)
	if err != nil {
		return err
	}

	obj := gjson.ParseBytes(res)
	format := cmd.Root().String("format")
	transform := cmd.Root().String("transform")
	return ShowJSON(os.Stdout, "messages:batches retrieve", obj, format, transform)
}

func handleMessagesBatchesList(ctx context.Context, cmd *cli.Command) error {
	client := anthropic.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()

	if len(unusedArgs) > 0 {
		return fmt.Errorf("Unexpected extra arguments: %v", unusedArgs)
	}

	params := anthropic.MessageBatchListParams{}

	options, err := flagOptions(
		cmd,
		apiquery.NestedQueryFormatBrackets,
		apiquery.ArrayQueryFormatComma,
		EmptyBody,
		false,
	)
	if err != nil {
		return err
	}

	format := cmd.Root().String("format")
	transform := cmd.Root().String("transform")
	if format == "raw" {
		var res []byte
		options = append(options, option.WithResponseBodyInto(&res))
		_, err = client.Messages.Batches.List(ctx, params, options...)
		if err != nil {
			return err
		}
		obj := gjson.ParseBytes(res)
		return ShowJSON(os.Stdout, "messages:batches list", obj, format, transform)
	} else {
		iter := client.Messages.Batches.ListAutoPaging(ctx, params, options...)
		return streamOutput("messages:batches list", func(w *os.File) error {
			for iter.Next() {
				item := iter.Current()
				obj := gjson.Parse(item.RawJSON())
				if err := ShowJSON(w, "messages:batches list", obj, format, transform); err != nil {
					return err
				}
			}
			return iter.Err()
		})
	}
}

func handleMessagesBatchesDelete(ctx context.Context, cmd *cli.Command) error {
	client := anthropic.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("message-batch-id") && len(unusedArgs) > 0 {
		cmd.Set("message-batch-id", unusedArgs[0])
		unusedArgs = unusedArgs[1:]
	}
	if len(unusedArgs) > 0 {
		return fmt.Errorf("Unexpected extra arguments: %v", unusedArgs)
	}

	options, err := flagOptions(
		cmd,
		apiquery.NestedQueryFormatBrackets,
		apiquery.ArrayQueryFormatComma,
		EmptyBody,
		false,
	)
	if err != nil {
		return err
	}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Messages.Batches.Delete(ctx, cmd.Value("message-batch-id").(string), options...)
	if err != nil {
		return err
	}

	obj := gjson.ParseBytes(res)
	format := cmd.Root().String("format")
	transform := cmd.Root().String("transform")
	return ShowJSON(os.Stdout, "messages:batches delete", obj, format, transform)
}

func handleMessagesBatchesCancel(ctx context.Context, cmd *cli.Command) error {
	client := anthropic.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("message-batch-id") && len(unusedArgs) > 0 {
		cmd.Set("message-batch-id", unusedArgs[0])
		unusedArgs = unusedArgs[1:]
	}
	if len(unusedArgs) > 0 {
		return fmt.Errorf("Unexpected extra arguments: %v", unusedArgs)
	}

	options, err := flagOptions(
		cmd,
		apiquery.NestedQueryFormatBrackets,
		apiquery.ArrayQueryFormatComma,
		EmptyBody,
		false,
	)
	if err != nil {
		return err
	}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Messages.Batches.Cancel(ctx, cmd.Value("message-batch-id").(string), options...)
	if err != nil {
		return err
	}

	obj := gjson.ParseBytes(res)
	format := cmd.Root().String("format")
	transform := cmd.Root().String("transform")
	return ShowJSON(os.Stdout, "messages:batches cancel", obj, format, transform)
}

func handleMessagesBatchesResults(ctx context.Context, cmd *cli.Command) error {
	client := anthropic.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("message-batch-id") && len(unusedArgs) > 0 {
		cmd.Set("message-batch-id", unusedArgs[0])
		unusedArgs = unusedArgs[1:]
	}
	if len(unusedArgs) > 0 {
		return fmt.Errorf("Unexpected extra arguments: %v", unusedArgs)
	}

	options, err := flagOptions(
		cmd,
		apiquery.NestedQueryFormatBrackets,
		apiquery.ArrayQueryFormatComma,
		EmptyBody,
		false,
	)
	if err != nil {
		return err
	}

	format := cmd.Root().String("format")
	transform := cmd.Root().String("transform")
	stream := client.Messages.Batches.ResultsStreaming(ctx, cmd.Value("message-batch-id").(string), options...)
	for stream.Next() {
		response := stream.Current()
		jsonData, err := json.Marshal(response)
		if err != nil {
			return err
		}
		obj := gjson.ParseBytes(jsonData)
		ShowJSON(os.Stdout, "messages:batches results", obj, format, transform)
	}
	return stream.Err()
}
