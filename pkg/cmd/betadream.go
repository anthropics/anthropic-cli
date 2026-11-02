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

var betaDreamsCreate = requestflag.WithInnerFlags(cli.Command{
	Name:    "create",
	Usage:   "Start an asynchronous job that uses past sessions to produce a reorganized\nversion of a memory store and get back the dream to poll for the result.",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[[]map[string]any]{
			Name:     "input",
			Usage:    "The memory store and sessions for the dream to read, as exactly one `memory_store` entry and exactly one `sessions` entry.",
			Required: true,
			BodyPath: "inputs",
		},
		&requestflag.Flag[any]{
			Name:     "model",
			Usage:    "The model that runs a dream, given as a model ID or as an object with `id` and `speed`.\n\nIn the object form, `speed` can only be `standard`.\n\nThe [Availability section of the Dreams guide](https://platform.claude.com/docs/en/managed-agents/dreams#feature-availability) lists the supported models.",
			Required: true,
			BodyPath: "model",
		},
		&requestflag.Flag[*string]{
			Name:     "instructions",
			Usage:    "Guidance that steers how the dream reads the sessions and organizes the output memory store, from 1 to 4,096 characters.\n\nSee the [Dreams guide](https://platform.claude.com/docs/en/managed-agents/dreams#steer-with-instructions) for what kinds of instructions work well.",
			BodyPath: "instructions",
		},
		&requestflag.Flag[map[string]any]{
			Name:     "output-behavior",
			Usage:    "Which memory store a dream writes its result to. Defaults to `create_new` when left out of a create request.",
			BodyPath: "output_behavior",
		},
		&requestflag.Flag[[]string]{
			Name:       "beta",
			Usage:      "This endpoint is in beta: requests must send `dreaming-2026-04-21` in this header.",
			HeaderPath: "anthropic-beta",
		},
		&requestflag.Flag[string]{
			Name:       "workspace-id",
			Usage:      "Optional header to select the Workspace for this request. The value is a Workspace ID (for example, `wrkspc_011CZkZaBF1tNoB5wlCeusgy`).\n\nOnly needed for credentials that can act on more than one Workspace. A credential that belongs to a specific Workspace may omit it; if sent, it must match that Workspace.",
			HeaderPath: "anthropic-workspace-id",
		},
	},
	Action:          handleBetaDreamsCreate,
	HideHelpCommand: true,
}, map[string][]requestflag.HasOuterFlag{
	"input": {
		&requestflag.InnerFlag[string]{
			Name:       "input.type",
			Usage:      `Allowed values: "memory_store", "sessions".`,
			InnerField: "type",
		},
		&requestflag.InnerFlag[string]{
			Name:       "input.memory-store-id",
			Usage:      "The ID of the memory store for the dream to read (`memstore_...`).\n\nThe memory store must be in the same workspace as the dream and must not be archived.",
			InnerField: "memory_store_id",
		},
		&requestflag.InnerFlag[[]string]{
			Name:       "input.session-ids",
			Usage:      "The IDs of the sessions whose transcripts the dream reads (`sesn_...`).\n\nGive 1 to 100 IDs, with no duplicates. Each session must be in the same workspace as the dream. Responses list the IDs in sorted order.",
			InnerField: "session_ids",
		},
	},
	"model": {
		&requestflag.InnerFlag[string]{
			Name:       "model.id",
			Usage:      "The ID of the model to run the dream with.\n\nThe ID can be 1 to 256 characters long.\n\nThe [Availability section of the Dreams guide](https://platform.claude.com/docs/en/managed-agents/dreams#feature-availability) lists the supported models.",
			InnerField: "id",
		},
		&requestflag.InnerFlag[*string]{
			Name:       "model.speed",
			Usage:      "How fast the model generates output for the dream. Defaults to `standard`.\n\nDreams accept only `standard`.",
			InnerField: "speed",
		},
	},
	"output-behavior": {
		&requestflag.InnerFlag[string]{
			Name:       "output-behavior.type",
			Usage:      `Allowed values: "create_new", "update_existing".`,
			InnerField: "type",
		},
		&requestflag.InnerFlag[string]{
			Name:       "output-behavior.memory-store-id",
			Usage:      "The ID of the memory store for the dream to write its result to (`memstore_...`). It must be the memory store in the `memory_store` entry of `inputs`.",
			InnerField: "memory_store_id",
		},
	},
})

var betaDreamsRetrieve = cli.Command{
	Name:    "retrieve",
	Usage:   "Get a dream by ID to check its status, output memory store, and token usage.",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:        "dream-id",
			Usage:       "The ID of the dream to get (`drm_...`).",
			Required:    true,
			PathParam:   "dream_id",
			DataAliases: []string{"id"},
		},
		&requestflag.Flag[[]string]{
			Name:       "beta",
			Usage:      "This endpoint is in beta: requests must send `dreaming-2026-04-21` in this header.",
			HeaderPath: "anthropic-beta",
		},
		&requestflag.Flag[string]{
			Name:       "workspace-id",
			Usage:      "Optional header to select the Workspace for this request. The value is a Workspace ID (for example, `wrkspc_011CZkZaBF1tNoB5wlCeusgy`).\n\nOnly needed for credentials that can act on more than one Workspace. A credential that belongs to a specific Workspace may omit it; if sent, it must match that Workspace.",
			HeaderPath: "anthropic-workspace-id",
		},
	},
	Action:          handleBetaDreamsRetrieve,
	HideHelpCommand: true,
}

var betaDreamsList = cli.Command{
	Name:    "list",
	Usage:   "List the dreams in the workspace, newest first.",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[any]{
			Name:      "created-at-gt",
			Usage:     "Return only dreams created after this time (exclusive), in RFC 3339.",
			QueryPath: "created_at[gt]",
		},
		&requestflag.Flag[any]{
			Name:      "created-at-lt",
			Usage:     "Return only dreams created before this time (exclusive), in RFC 3339.",
			QueryPath: "created_at[lt]",
		},
		&requestflag.Flag[bool]{
			Name:      "include-archived",
			Usage:     "Whether to include archived dreams. Defaults to `false`.",
			QueryPath: "include_archived",
		},
		&requestflag.Flag[int64]{
			Name:      "limit",
			Usage:     "The maximum number of dreams to return, from 1 to 100. Defaults to 20.",
			QueryPath: "limit",
		},
		&requestflag.Flag[string]{
			Name:      "page",
			Usage:     "The cursor for the page to return, taken from `next_page` in a previous response.\n\nLeave it out to get the first page.",
			QueryPath: "page",
		},
		&requestflag.Flag[[]string]{
			Name:      "status",
			Usage:     "Return only dreams that have one of these statuses.\n\nRepeat the parameter to give more than one status. Leave it out to return dreams of every status.",
			QueryPath: "statuses",
		},
		&requestflag.Flag[[]string]{
			Name:       "beta",
			Usage:      "This endpoint is in beta: requests must send `dreaming-2026-04-21` in this header.",
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
	Action:          handleBetaDreamsList,
	HideHelpCommand: true,
}

var betaDreamsArchive = cli.Command{
	Name:    "archive",
	Usage:   "Hide a `completed`, `failed`, or `canceled` dream from the default list of\ndreams.",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:        "dream-id",
			Usage:       "The ID of the dream to archive (`drm_...`).",
			Required:    true,
			PathParam:   "dream_id",
			DataAliases: []string{"id"},
		},
		&requestflag.Flag[[]string]{
			Name:       "beta",
			Usage:      "This endpoint is in beta: requests must send `dreaming-2026-04-21` in this header.",
			HeaderPath: "anthropic-beta",
		},
		&requestflag.Flag[string]{
			Name:       "workspace-id",
			Usage:      "Optional header to select the Workspace for this request. The value is a Workspace ID (for example, `wrkspc_011CZkZaBF1tNoB5wlCeusgy`).\n\nOnly needed for credentials that can act on more than one Workspace. A credential that belongs to a specific Workspace may omit it; if sent, it must match that Workspace.",
			HeaderPath: "anthropic-workspace-id",
		},
	},
	Action:          handleBetaDreamsArchive,
	HideHelpCommand: true,
}

var betaDreamsCancel = cli.Command{
	Name:    "cancel",
	Usage:   "Stop a `pending` or `running` dream.",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:        "dream-id",
			Usage:       "The ID of the dream to cancel (`drm_...`).",
			Required:    true,
			PathParam:   "dream_id",
			DataAliases: []string{"id"},
		},
		&requestflag.Flag[[]string]{
			Name:       "beta",
			Usage:      "This endpoint is in beta: requests must send `dreaming-2026-04-21` in this header.",
			HeaderPath: "anthropic-beta",
		},
		&requestflag.Flag[string]{
			Name:       "workspace-id",
			Usage:      "Optional header to select the Workspace for this request. The value is a Workspace ID (for example, `wrkspc_011CZkZaBF1tNoB5wlCeusgy`).\n\nOnly needed for credentials that can act on more than one Workspace. A credential that belongs to a specific Workspace may omit it; if sent, it must match that Workspace.",
			HeaderPath: "anthropic-workspace-id",
		},
	},
	Action:          handleBetaDreamsCancel,
	HideHelpCommand: true,
}

func handleBetaDreamsCreate(ctx context.Context, cmd *cli.Command) error {
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

	params := anthropic.BetaDreamNewParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Beta.Dreams.New(ctx, params, options...)
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
		Title:          "beta:dreams create",
		Transform:      transform,
	})
}

func handleBetaDreamsRetrieve(ctx context.Context, cmd *cli.Command) error {
	client := anthropic.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("dream-id") && len(unusedArgs) > 0 {
		cmd.Set("dream-id", unusedArgs[0])
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

	params := anthropic.BetaDreamGetParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Beta.Dreams.Get(
		ctx,
		cmd.Value("dream-id").(string),
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
		Title:          "beta:dreams retrieve",
		Transform:      transform,
	})
}

func handleBetaDreamsList(ctx context.Context, cmd *cli.Command) error {
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

	params := anthropic.BetaDreamListParams{}

	format := "explore"
	explicitFormat := cmd.Root().IsSet("format")
	if explicitFormat {
		format = cmd.Root().String("format")
	}
	transform := cmd.Root().String("transform")
	if format == "raw" {
		var res []byte
		options = append(options, option.WithResponseBodyInto(&res))
		_, err = client.Beta.Dreams.List(ctx, params, options...)
		if err != nil {
			return err
		}
		obj := gjson.ParseBytes(res)
		return ShowJSON(obj, ShowJSONOpts{
			ExplicitFormat: explicitFormat,
			Format:         format,
			RawOutput:      cmd.Root().Bool("raw-output"),
			Title:          "beta:dreams list",
			Transform:      transform,
		})
	} else {
		iter := client.Beta.Dreams.ListAutoPaging(ctx, params, options...)
		maxItems := int64(-1)
		if cmd.IsSet("max-items") {
			maxItems = cmd.Value("max-items").(int64)
		}
		return ShowJSONIterator(iter, maxItems, ShowJSONOpts{
			ExplicitFormat: explicitFormat,
			Format:         format,
			RawOutput:      cmd.Root().Bool("raw-output"),
			Title:          "beta:dreams list",
			Transform:      transform,
		})
	}
}

func handleBetaDreamsArchive(ctx context.Context, cmd *cli.Command) error {
	client := anthropic.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("dream-id") && len(unusedArgs) > 0 {
		cmd.Set("dream-id", unusedArgs[0])
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

	params := anthropic.BetaDreamArchiveParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Beta.Dreams.Archive(
		ctx,
		cmd.Value("dream-id").(string),
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
		Title:          "beta:dreams archive",
		Transform:      transform,
	})
}

func handleBetaDreamsCancel(ctx context.Context, cmd *cli.Command) error {
	client := anthropic.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("dream-id") && len(unusedArgs) > 0 {
		cmd.Set("dream-id", unusedArgs[0])
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

	params := anthropic.BetaDreamCancelParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Beta.Dreams.Cancel(
		ctx,
		cmd.Value("dream-id").(string),
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
		Title:          "beta:dreams cancel",
		Transform:      transform,
	})
}
