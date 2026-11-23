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

var betaMemoryStoresMemoriesCreate = cli.Command{
	Name:    "create",
	Usage:   "Create a Memory",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "memory-store-id",
			Usage:     "The ID of the memory store to create the memory in (`memstore_...`).",
			Required:  true,
			PathParam: "memory_store_id",
		},
		&requestflag.Flag[*string]{
			Name:     "content",
			Usage:    "The UTF-8 text content of the new memory. Pass `\"\"` to create an empty memory.\n\nThe maximum size is 100 kB (102,400 bytes).",
			Required: true,
			BodyPath: "content",
		},
		&requestflag.Flag[string]{
			Name:     "path",
			Usage:    "The path of the new memory, such as `/projects/foo/notes.md`. Paths are case-sensitive.\n\nA path starts with `/` and has 1 to 20 segments. A segment can't be empty and can't be `.` or `..`. A path can't contain control characters, format characters, or the line and paragraph separators U+2028 and U+2029. It must be in Unicode Normalization Form C and at most 1,024 bytes long.",
			Required: true,
			BodyPath: "path",
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
	Action:          handleBetaMemoryStoresMemoriesCreate,
	HideHelpCommand: true,
}

var betaMemoryStoresMemoriesRetrieve = cli.Command{
	Name:    "retrieve",
	Usage:   "Retrieve a Memory",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "memory-store-id",
			Usage:     "The ID of the memory store that holds the memory (`memstore_...`).",
			Required:  true,
			PathParam: "memory_store_id",
		},
		&requestflag.Flag[string]{
			Name:      "memory-id",
			Usage:     "The ID of the memory to retrieve (`mem_...`).",
			Required:  true,
			PathParam: "memory_id",
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
	Action:          handleBetaMemoryStoresMemoriesRetrieve,
	HideHelpCommand: true,
}

var betaMemoryStoresMemoriesUpdate = requestflag.WithInnerFlags(cli.Command{
	Name:    "update",
	Usage:   "Update a Memory",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "memory-store-id",
			Usage:     "The ID of the memory store that holds the memory (`memstore_...`).",
			Required:  true,
			PathParam: "memory_store_id",
		},
		&requestflag.Flag[string]{
			Name:        "memory-id",
			Usage:       "The ID of the memory to update (`mem_...`).",
			Required:    true,
			PathParam:   "memory_id",
			DataAliases: []string{"id"},
		},
		&requestflag.Flag[string]{
			Name:      "view",
			Usage:     "Whether a `memory` or `memory_version` in the response includes its `content`. Defaults to `full` on retrieve operations and to `basic` on list, create, and update operations.\n\nWith `full`, list operations return at most 20 items per page.\n\n- `basic` - Return the object with `content` set to `null`. The `content_size_bytes` and `content_sha256` fields remain populated, so sync clients can diff without fetching content.\n- `full` - Return the object with `content` populated. On list endpoints, `view=full` caps `limit` at 20.",
			QueryPath: "view",
		},
		&requestflag.Flag[*string]{
			Name:     "content",
			Usage:    "New UTF-8 text content for the memory. Maximum 100 kB (102,400 bytes). Omit to leave the content unchanged (e.g., for a rename-only update).",
			BodyPath: "content",
		},
		&requestflag.Flag[*string]{
			Name:     "path",
			Usage:    "The new path of the memory, which renames it. Omit it to keep the current path.\n\nThe memory keeps its `id`. A path starts with `/` and has 1 to 20 segments. A segment can't be empty and can't be `.` or `..`. A path can't contain control characters, format characters, or the line and paragraph separators U+2028 and U+2029. It must be in Unicode Normalization Form C and at most 1,024 bytes long. Paths are case-sensitive.",
			BodyPath: "path",
		},
		&requestflag.Flag[map[string]any]{
			Name:     "precondition",
			Usage:    "Optional condition that must hold for an update to apply. When omitted, the update is unconditional. Asserts the current state of the memory being updated. When an update changes `path`, the precondition still refers to the memory's current content, not the destination path. Currently the only supported variant is `content_sha256`.",
			BodyPath: "precondition",
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
			Usage:      "Expected `content_sha256` of the stored memory (64 lowercase hexadecimal characters). Typically the `content_sha256` returned by a prior read or list call. Because the server applies no content normalization, clients can also compute this locally as the SHA-256 of the UTF-8 content bytes.",
			InnerField: "content_sha256",
		},
	},
})

var betaMemoryStoresMemoriesList = cli.Command{
	Name:    "list",
	Usage:   "List Memories",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "memory-store-id",
			Usage:     "The ID of the memory store to list memories from (`memstore_...`).",
			Required:  true,
			PathParam: "memory_store_id",
		},
		&requestflag.Flag[int64]{
			Name:      "depth",
			Usage:     "How far below `path_prefix` to list. Defaults to `0`, which lists every memory below the prefix.\n\nPass `1` to list only the direct children. Deeper memories are then grouped into `memory_prefix` items.",
			QueryPath: "depth",
		},
		&requestflag.Flag[int64]{
			Name:      "limit",
			Usage:     "The maximum number of items to return per page, from 1 to 100. Defaults to 20.\n\nWith `view=full`, a page holds at most 20 items. `memory` and `memory_prefix` items both count toward the limit.",
			QueryPath: "limit",
		},
		&requestflag.Flag[string]{
			Name:      "page",
			Usage:     "The `next_page` value from a previous response, to get the next page. Omit it to get the first page.",
			QueryPath: "page",
		},
		&requestflag.Flag[string]{
			Name:      "path-prefix",
			Usage:     "Return only memories whose path starts with this prefix, such as `/notes/`. This value appears in request URLs. Do not include secrets or personally identifiable information.\n\nOmit it to list the whole store. The prefix must end with `/` and matches whole segments, so `/notes/` matches `/notes/todo.md` but not `/notes-archive/todo.md`.",
			QueryPath: "path_prefix",
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
	Action:          handleBetaMemoryStoresMemoriesList,
	HideHelpCommand: true,
}

var betaMemoryStoresMemoriesDelete = cli.Command{
	Name:    "delete",
	Usage:   "Delete a Memory",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "memory-store-id",
			Usage:     "The ID of the memory store that holds the memory (`memstore_...`).",
			Required:  true,
			PathParam: "memory_store_id",
		},
		&requestflag.Flag[string]{
			Name:      "memory-id",
			Usage:     "The ID of the memory to delete (`mem_...`).",
			Required:  true,
			PathParam: "memory_id",
		},
		&requestflag.Flag[string]{
			Name:      "expected-content-sha256",
			Usage:     "Delete the memory only if its current `content_sha256` equals this value, given as 64 lowercase hexadecimal characters. Omit it to delete unconditionally.\n\nIf the hashes differ, the request fails with HTTP status 409 and nothing is deleted.",
			QueryPath: "expected_content_sha256",
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

	params := anthropic.BetaMemoryStoreMemoryNewParams{}

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

	params := anthropic.BetaMemoryStoreMemoryGetParams{
		MemoryStoreID: cmd.Value("memory-store-id").(string),
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

	params := anthropic.BetaMemoryStoreMemoryUpdateParams{
		MemoryStoreID: cmd.Value("memory-store-id").(string),
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

	params := anthropic.BetaMemoryStoreMemoryListParams{}

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

	params := anthropic.BetaMemoryStoreMemoryDeleteParams{
		MemoryStoreID: cmd.Value("memory-store-id").(string),
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
