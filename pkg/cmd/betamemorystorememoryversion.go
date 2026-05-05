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

var betaMemoryStoresMemoryVersionsRetrieve = cli.Command{
	Name:    "retrieve",
	Usage:   "Retrieve a memory version",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "memory-store-id",
			Required:  true,
			PathParam: "memory_store_id",
		},
		&requestflag.Flag[string]{
			Name:      "memory-version-id",
			Required:  true,
			PathParam: "memory_version_id",
		},
		&requestflag.Flag[string]{
			Name:      "view",
			Usage:     "Selects which projection of a `memory` or `memory_version` the server returns. `basic` returns the object with `content` set to `null`; `full` populates `content`. When omitted, the default is endpoint-specific: retrieve operations default to `full`; list, create, and update operations default to `basic`. Listing with `view=full` caps `limit` at 20.",
			QueryPath: "view",
		},
		&requestflag.Flag[[]string]{
			Name:       "beta",
			Usage:      "Optional header to specify the beta version(s) you want to use.",
			HeaderPath: "anthropic-beta",
		},
	},
	Action:          handleBetaMemoryStoresMemoryVersionsRetrieve,
	HideHelpCommand: true,
}

var betaMemoryStoresMemoryVersionsList = cli.Command{
	Name:    "list",
	Usage:   "List memory versions",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "memory-store-id",
			Required:  true,
			PathParam: "memory_store_id",
		},
		&requestflag.Flag[string]{
			Name:      "api-key-id",
			Usage:     "Query parameter for api_key_id",
			QueryPath: "api_key_id",
		},
		&requestflag.Flag[any]{
			Name:      "created-at-gte",
			Usage:     "Return versions created at or after this time (inclusive).",
			QueryPath: "created_at[gte]",
		},
		&requestflag.Flag[any]{
			Name:      "created-at-lte",
			Usage:     "Return versions created at or before this time (inclusive).",
			QueryPath: "created_at[lte]",
		},
		&requestflag.Flag[int64]{
			Name:      "limit",
			Usage:     "Query parameter for limit",
			QueryPath: "limit",
		},
		&requestflag.Flag[string]{
			Name:      "memory-id",
			Usage:     "Query parameter for memory_id",
			QueryPath: "memory_id",
		},
		&requestflag.Flag[string]{
			Name:      "operation",
			Usage:     "The kind of mutation a `memory_version` records. Every non-no-op mutation to a memory appends exactly one version row with one of these values.",
			QueryPath: "operation",
		},
		&requestflag.Flag[string]{
			Name:      "page",
			Usage:     "Query parameter for page",
			QueryPath: "page",
		},
		&requestflag.Flag[string]{
			Name:      "session-id",
			Usage:     "Query parameter for session_id",
			QueryPath: "session_id",
		},
		&requestflag.Flag[string]{
			Name:      "view",
			Usage:     "Selects which projection of a `memory` or `memory_version` the server returns. `basic` returns the object with `content` set to `null`; `full` populates `content`. When omitted, the default is endpoint-specific: retrieve operations default to `full`; list, create, and update operations default to `basic`. Listing with `view=full` caps `limit` at 20.",
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
	Action:          handleBetaMemoryStoresMemoryVersionsList,
	HideHelpCommand: true,
}

var betaMemoryStoresMemoryVersionsRedact = cli.Command{
	Name:    "redact",
	Usage:   "Redact a memory version",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "memory-store-id",
			Required:  true,
			PathParam: "memory_store_id",
		},
		&requestflag.Flag[string]{
			Name:      "memory-version-id",
			Required:  true,
			PathParam: "memory_version_id",
		},
		&requestflag.Flag[[]string]{
			Name:       "beta",
			Usage:      "Optional header to specify the beta version(s) you want to use.",
			HeaderPath: "anthropic-beta",
		},
	},
	Action:          handleBetaMemoryStoresMemoryVersionsRedact,
	HideHelpCommand: true,
}

func handleBetaMemoryStoresMemoryVersionsRetrieve(ctx context.Context, cmd *cli.Command) error {
	client := anthropic.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("memory-version-id") && len(unusedArgs) > 0 {
		cmd.Set("memory-version-id", unusedArgs[0])
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

	params := anthropic.BetaMemoryStoreMemoryVersionGetParams{
		MemoryStoreID: cmd.Value("memory-store-id").(string),
	}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Beta.MemoryStores.MemoryVersions.Get(
		ctx,
		cmd.Value("memory-version-id").(string),
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
		Title:          "beta:memory-stores:memory-versions retrieve",
		Transform:      transform,
	})
}

func handleBetaMemoryStoresMemoryVersionsList(ctx context.Context, cmd *cli.Command) error {
	client := anthropic.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("memory-store-id") && len(unusedArgs) > 0 {
		cmd.Set("memory-store-id", unusedArgs[0])
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

	params := anthropic.BetaMemoryStoreMemoryVersionListParams{}

	format := "explore"
	explicitFormat := cmd.Root().IsSet("format")
	if explicitFormat {
		format = cmd.Root().String("format")
	}
	transform := cmd.Root().String("transform")
	if format == "raw" {
		var res []byte
		options = append(options, option.WithResponseBodyInto(&res))
		_, err = client.Beta.MemoryStores.MemoryVersions.List(
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
			Title:          "beta:memory-stores:memory-versions list",
			Transform:      transform,
		})
	} else {
		iter := client.Beta.MemoryStores.MemoryVersions.ListAutoPaging(
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
			Title:          "beta:memory-stores:memory-versions list",
			Transform:      transform,
		})
	}
}

func handleBetaMemoryStoresMemoryVersionsRedact(ctx context.Context, cmd *cli.Command) error {
	client := anthropic.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("memory-version-id") && len(unusedArgs) > 0 {
		cmd.Set("memory-version-id", unusedArgs[0])
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

	params := anthropic.BetaMemoryStoreMemoryVersionRedactParams{
		MemoryStoreID: cmd.Value("memory-store-id").(string),
	}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Beta.MemoryStores.MemoryVersions.Redact(
		ctx,
		cmd.Value("memory-version-id").(string),
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
		Title:          "beta:memory-stores:memory-versions redact",
		Transform:      transform,
	})
}
