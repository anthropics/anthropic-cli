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

var betaMemoryStoresMemoryVersionsRetrieve = cli.Command{
	Name:    "retrieve",
	Usage:   "Retrieve a Memory Version",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "memory-store-id",
			Usage:     "The ID of the memory store that holds the version (`memstore_...`).",
			Required:  true,
			PathParam: "memory_store_id",
		},
		&requestflag.Flag[string]{
			Name:      "memory-version-id",
			Usage:     "The ID of the memory version to retrieve (`memver_...`).",
			Required:  true,
			PathParam: "memory_version_id",
		},
		&requestflag.Flag[string]{
			Name:      "view",
			Usage:     "Whether a `memory` or `memory_version` in the response includes its `content`. Defaults to `full` on retrieve operations and to `basic` on list, create, and update operations.\n\nWith `full`, list operations return at most 20 items per page.\n\n- `basic` - Return the object with `content` set to `null`. The `content_size_bytes` and `content_sha256` fields remain populated, so sync clients can diff without fetching content.\n- `full` - Return the object with `content` populated. On list endpoints, `view=full` caps `limit` at 20.",
			QueryPath: "view",
		},
		&requestflag.Flag[[]string]{
			Name:       "beta",
			Usage:      "This endpoint is in beta: requests must send `agent-memory-2026-07-22` or `managed-agents-2026-04-01` in this header.",
			HeaderPath: "anthropic-beta",
		},
		&requestflag.Flag[string]{
			Name:       "workspace-id",
			Usage:      "Optional header to select the Workspace for this request. The value is a Workspace ID (for example, `wrkspc_011CZkZaBF1tNoB5wlCeusgy`).\n\nOnly needed for credentials that can act on more than one Workspace. A credential that belongs to a specific Workspace may omit it; if sent, it must match that Workspace.",
			HeaderPath: "anthropic-workspace-id",
		},
	},
	Action:          handleBetaMemoryStoresMemoryVersionsRetrieve,
	HideHelpCommand: true,
}

var betaMemoryStoresMemoryVersionsList = cli.Command{
	Name:    "list",
	Usage:   "List Memory Versions",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "memory-store-id",
			Usage:     "The ID of the memory store whose version history to list (`memstore_...`).",
			Required:  true,
			PathParam: "memory_store_id",
		},
		&requestflag.Flag[string]{
			Name:      "api-key-id",
			Usage:     "Return only versions written with the API key that has this ID.",
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
			Usage:     "The maximum number of versions to return per page, from 1 to 100. Defaults to 20.\n\nWith `view=full`, a page holds at most 20 versions.",
			QueryPath: "limit",
		},
		&requestflag.Flag[string]{
			Name:      "memory-id",
			Usage:     "Return only versions of the memory with this ID (`mem_...`).\n\nThe filter still works after the memory is deleted. While that memory's versions are retained, the results include the version whose `operation` is `deleted`.",
			QueryPath: "memory_id",
		},
		&requestflag.Flag[string]{
			Name:      "operation",
			Usage:     "The kind of change that a memory version records.\n\nEvery change to a memory adds one version. If you update a memory and nothing changes, the memory doesn't get a new version.\n\n- `created` - The memory was created. The first version in any memory's lineage.\n- `modified` - The memory's `content`, `path`, or both were changed via update. Writes the agent makes through the filesystem mount also appear as `modified`.\n- `deleted` - The memory was deleted. The `content`, `content_size_bytes`, and `content_sha256` fields are `null` on this version. The preceding version, while it is retained, records the deleted content's size and hash.",
			QueryPath: "operation",
		},
		&requestflag.Flag[string]{
			Name:      "page",
			Usage:     "The `next_page` value from a previous response, to get the next page. Omit it to get the first page.",
			QueryPath: "page",
		},
		&requestflag.Flag[string]{
			Name:      "service-account-id",
			Usage:     "Return only versions written by the service account with this ID (`svac_...`).",
			QueryPath: "service_account_id",
		},
		&requestflag.Flag[string]{
			Name:      "session-id",
			Usage:     "Return only versions written by the session with this ID.",
			QueryPath: "session_id",
		},
		&requestflag.Flag[string]{
			Name:      "view",
			Usage:     "Whether a `memory` or `memory_version` in the response includes its `content`. Defaults to `full` on retrieve operations and to `basic` on list, create, and update operations.\n\nWith `full`, list operations return at most 20 items per page.\n\n- `basic` - Return the object with `content` set to `null`. The `content_size_bytes` and `content_sha256` fields remain populated, so sync clients can diff without fetching content.\n- `full` - Return the object with `content` populated. On list endpoints, `view=full` caps `limit` at 20.",
			QueryPath: "view",
		},
		&requestflag.Flag[[]string]{
			Name:       "beta",
			Usage:      "This endpoint is in beta: requests must send `agent-memory-2026-07-22` or `managed-agents-2026-04-01` in this header.",
			HeaderPath: "anthropic-beta",
		},
		&requestflag.Flag[string]{
			Name:       "workspace-id",
			Usage:      "Optional header to select the Workspace for this request. The value is a Workspace ID (for example, `wrkspc_011CZkZaBF1tNoB5wlCeusgy`).\n\nOnly needed for credentials that can act on more than one Workspace. A credential that belongs to a specific Workspace may omit it; if sent, it must match that Workspace.",
			HeaderPath: "anthropic-workspace-id",
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
	Usage:   "Redact a Memory Version",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "memory-store-id",
			Usage:     "The ID of the memory store that holds the version (`memstore_...`).",
			Required:  true,
			PathParam: "memory_store_id",
		},
		&requestflag.Flag[string]{
			Name:      "memory-version-id",
			Usage:     "The ID of the memory version to redact (`memver_...`).",
			Required:  true,
			PathParam: "memory_version_id",
		},
		&requestflag.Flag[[]string]{
			Name:       "beta",
			Usage:      "This endpoint is in beta: requests must send `agent-memory-2026-07-22` or `managed-agents-2026-04-01` in this header.",
			HeaderPath: "anthropic-beta",
		},
		&requestflag.Flag[string]{
			Name:       "workspace-id",
			Usage:      "Optional header to select the Workspace for this request. The value is a Workspace ID (for example, `wrkspc_011CZkZaBF1tNoB5wlCeusgy`).\n\nOnly needed for credentials that can act on more than one Workspace. A credential that belongs to a specific Workspace may omit it; if sent, it must match that Workspace.",
			HeaderPath: "anthropic-workspace-id",
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
