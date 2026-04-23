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

var betaMemoryStoresMemoriesCreate = cli.Command{
	Name:    "create",
	Usage:   "CreateMemory",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:     "memory-store-id",
			Required: true,
		},
		&requestflag.Flag[any]{
			Name:     "content",
			Required: true,
			BodyPath: "content",
		},
		&requestflag.Flag[string]{
			Name:     "path",
			Required: true,
			BodyPath: "path",
		},
		&requestflag.Flag[string]{
			Name:      "view",
			Usage:     "MemoryView enum",
			QueryPath: "view",
		},
		&requestflag.Flag[[]string]{
			Name:       "beta",
			Usage:      "Optional header to specify the beta version(s) you want to use.",
			HeaderPath: "anthropic-beta",
		},
	},
	Action:          handleBetaMemoryStoresMemoriesCreate,
	HideHelpCommand: true,
}

var betaMemoryStoresMemoriesRetrieve = cli.Command{
	Name:    "retrieve",
	Usage:   "GetMemory",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:     "memory-store-id",
			Required: true,
		},
		&requestflag.Flag[string]{
			Name:     "memory-id",
			Required: true,
		},
		&requestflag.Flag[string]{
			Name:      "view",
			Usage:     "MemoryView enum",
			QueryPath: "view",
		},
		&requestflag.Flag[[]string]{
			Name:       "beta",
			Usage:      "Optional header to specify the beta version(s) you want to use.",
			HeaderPath: "anthropic-beta",
		},
	},
	Action:          handleBetaMemoryStoresMemoriesRetrieve,
	HideHelpCommand: true,
}

var betaMemoryStoresMemoriesUpdate = requestflag.WithInnerFlags(cli.Command{
	Name:    "update",
	Usage:   "UpdateMemory",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:     "memory-store-id",
			Required: true,
		},
		&requestflag.Flag[string]{
			Name:     "memory-id",
			Required: true,
		},
		&requestflag.Flag[string]{
			Name:      "view",
			Usage:     "MemoryView enum",
			QueryPath: "view",
		},
		&requestflag.Flag[any]{
			Name:     "content",
			BodyPath: "content",
		},
		&requestflag.Flag[any]{
			Name:     "path",
			BodyPath: "path",
		},
		&requestflag.Flag[map[string]any]{
			Name:     "precondition",
			BodyPath: "precondition",
		},
		&requestflag.Flag[[]string]{
			Name:       "beta",
			Usage:      "Optional header to specify the beta version(s) you want to use.",
			HeaderPath: "anthropic-beta",
		},
	},
	Action:          handleBetaMemoryStoresMemoriesUpdate,
	HideHelpCommand: true,
}, map[string][]requestflag.HasOuterFlag{
	"precondition": {
		&requestflag.InnerFlag[string]{
			Name:       "precondition.type",
			Usage:      `Allowed values: "content_sha256".`,
			InnerField: "type",
		},
		&requestflag.InnerFlag[string]{
			Name:       "precondition.content-sha256",
			InnerField: "content_sha256",
		},
	},
})

var betaMemoryStoresMemoriesList = cli.Command{
	Name:    "list",
	Usage:   "ListMemories",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:     "memory-store-id",
			Required: true,
		},
		&requestflag.Flag[int64]{
			Name:      "depth",
			Usage:     "Query parameter for depth",
			QueryPath: "depth",
		},
		&requestflag.Flag[int64]{
			Name:      "limit",
			Usage:     "Query parameter for limit",
			QueryPath: "limit",
		},
		&requestflag.Flag[string]{
			Name:      "order",
			Usage:     "Query parameter for order",
			QueryPath: "order",
		},
		&requestflag.Flag[string]{
			Name:      "order-by",
			Usage:     "Query parameter for order_by",
			QueryPath: "order_by",
		},
		&requestflag.Flag[string]{
			Name:      "page",
			Usage:     "Query parameter for page",
			QueryPath: "page",
		},
		&requestflag.Flag[string]{
			Name:      "path-prefix",
			Usage:     "Optional path prefix filter (raw string-prefix match; include a trailing slash for directory-scoped lists). This value appears in request URLs. Do not include secrets or personally identifiable information.",
			QueryPath: "path_prefix",
		},
		&requestflag.Flag[string]{
			Name:      "view",
			Usage:     "MemoryView enum",
			QueryPath: "view",
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
	Action:          handleBetaMemoryStoresMemoriesList,
	HideHelpCommand: true,
}

var betaMemoryStoresMemoriesDelete = cli.Command{
	Name:    "delete",
	Usage:   "DeleteMemory",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:     "memory-store-id",
			Required: true,
		},
		&requestflag.Flag[string]{
			Name:     "memory-id",
			Required: true,
		},
		&requestflag.Flag[string]{
			Name:      "expected-content-sha256",
			Usage:     "Query parameter for expected_content_sha256",
			QueryPath: "expected_content_sha256",
		},
		&requestflag.Flag[[]string]{
			Name:       "beta",
			Usage:      "Optional header to specify the beta version(s) you want to use.",
			HeaderPath: "anthropic-beta",
		},
	},
	Action:          handleBetaMemoryStoresMemoriesDelete,
	HideHelpCommand: true,
}

func handleBetaMemoryStoresMemoriesCreate(ctx context.Context, cmd *cli.Command) error {
	client := anthropic.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("memory-store-id") && len(unusedArgs) > 0 {
		cmd.Set("memory-store-id", unusedArgs[0])
		unusedArgs = unusedArgs[1:]
	}
	if len(unusedArgs) > 0 {
		return fmt.Errorf("Unexpected extra arguments: %v", unusedArgs)
	}

	params := anthropic.BetaMemoryStoreMemoryNewParams{}

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
	_, err = client.Beta.MemoryStores.Memories.New(
		ctx,
		cmd.Value("memory-store-id").(string),
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
		Title:          "beta:memory-stores:memories create",
		Transform:      transform,
	})
}

func handleBetaMemoryStoresMemoriesRetrieve(ctx context.Context, cmd *cli.Command) error {
	client := anthropic.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("memory-id") && len(unusedArgs) > 0 {
		cmd.Set("memory-id", unusedArgs[0])
		unusedArgs = unusedArgs[1:]
	}
	if len(unusedArgs) > 0 {
		return fmt.Errorf("Unexpected extra arguments: %v", unusedArgs)
	}

	params := anthropic.BetaMemoryStoreMemoryGetParams{
		MemoryStoreID: cmd.Value("memory-store-id").(string),
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
	_, err = client.Beta.MemoryStores.Memories.Get(
		ctx,
		cmd.Value("memory-id").(string),
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
		Title:          "beta:memory-stores:memories retrieve",
		Transform:      transform,
	})
}

func handleBetaMemoryStoresMemoriesUpdate(ctx context.Context, cmd *cli.Command) error {
	client := anthropic.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("memory-id") && len(unusedArgs) > 0 {
		cmd.Set("memory-id", unusedArgs[0])
		unusedArgs = unusedArgs[1:]
	}
	if len(unusedArgs) > 0 {
		return fmt.Errorf("Unexpected extra arguments: %v", unusedArgs)
	}

	params := anthropic.BetaMemoryStoreMemoryUpdateParams{
		MemoryStoreID: cmd.Value("memory-store-id").(string),
	}

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
	_, err = client.Beta.MemoryStores.Memories.Update(
		ctx,
		cmd.Value("memory-id").(string),
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
		Title:          "beta:memory-stores:memories update",
		Transform:      transform,
	})
}

func handleBetaMemoryStoresMemoriesList(ctx context.Context, cmd *cli.Command) error {
	client := anthropic.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("memory-store-id") && len(unusedArgs) > 0 {
		cmd.Set("memory-store-id", unusedArgs[0])
		unusedArgs = unusedArgs[1:]
	}
	if len(unusedArgs) > 0 {
		return fmt.Errorf("Unexpected extra arguments: %v", unusedArgs)
	}

	params := anthropic.BetaMemoryStoreMemoryListParams{}

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

	format := "explore"
	explicitFormat := cmd.Root().IsSet("format")
	if explicitFormat {
		format = cmd.Root().String("format")
	}
	transform := cmd.Root().String("transform")
	if format == "raw" {
		var res []byte
		options = append(options, option.WithResponseBodyInto(&res))
		_, err = client.Beta.MemoryStores.Memories.List(
			ctx,
			cmd.Value("memory-store-id").(string),
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
			Title:          "beta:memory-stores:memories list",
			Transform:      transform,
		})
	} else {
		iter := client.Beta.MemoryStores.Memories.ListAutoPaging(
			ctx,
			cmd.Value("memory-store-id").(string),
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
			Title:          "beta:memory-stores:memories list",
			Transform:      transform,
		})
	}
}

func handleBetaMemoryStoresMemoriesDelete(ctx context.Context, cmd *cli.Command) error {
	client := anthropic.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("memory-id") && len(unusedArgs) > 0 {
		cmd.Set("memory-id", unusedArgs[0])
		unusedArgs = unusedArgs[1:]
	}
	if len(unusedArgs) > 0 {
		return fmt.Errorf("Unexpected extra arguments: %v", unusedArgs)
	}

	params := anthropic.BetaMemoryStoreMemoryDeleteParams{
		MemoryStoreID: cmd.Value("memory-store-id").(string),
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
	_, err = client.Beta.MemoryStores.Memories.Delete(
		ctx,
		cmd.Value("memory-id").(string),
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
		Title:          "beta:memory-stores:memories delete",
		Transform:      transform,
	})
}
