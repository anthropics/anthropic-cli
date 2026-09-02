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

var betaAgentsCreate = requestflag.WithInnerFlags(cli.Command{
	Name:    "create",
	Usage:   "Create Agent",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[map[string]any]{
			Name:     "model",
			Usage:    "Model identifier. Accepts the [model string](https://platform.claude.com/docs/en/about-claude/models/overview#latest-models-comparison), e.g. `claude-opus-5`, or a `model_config` object for additional configuration control",
			Required: true,
			BodyPath: "model",
		},
		&requestflag.Flag[string]{
			Name:     "name",
			Usage:    "Human-readable name for the agent.",
			Required: true,
			BodyPath: "name",
		},
		&requestflag.Flag[*string]{
			Name:     "description",
			Usage:    "Description of what the agent does.",
			BodyPath: "description",
		},
		&requestflag.Flag[[]map[string]any]{
			Name:     "mcp-server",
			Usage:    "MCP servers this agent connects to. Maximum 20. Names must be unique within the array. Every server must be referenced by an `mcp_toolset` in `tools`; unreferenced servers are rejected. See the [MCP connector guide](https://platform.claude.com/docs/en/managed-agents/mcp-connector).",
			BodyPath: "mcp_servers",
		},
		&requestflag.Flag[map[string]any]{
			Name:     "metadata",
			Usage:    "Arbitrary key-value metadata. Maximum 16 pairs, keys up to 64 chars, values up to 512 chars.",
			BodyPath: "metadata",
		},
		&requestflag.Flag[map[string]any]{
			Name:     "multiagent",
			Usage:    "A coordinator topology: the session's primary thread orchestrates work by spawning session threads, each running an agent drawn from the `agents` roster.",
			BodyPath: "multiagent",
		},
		&requestflag.Flag[[]map[string]any]{
			Name:     "skill",
			Usage:    "Skills available to the agent.",
			BodyPath: "skills",
		},
		&requestflag.Flag[*string]{
			Name:     "system",
			Usage:    "System prompt for the agent.",
			BodyPath: "system",
		},
		&requestflag.Flag[[]map[string]any]{
			Name:     "tool",
			Usage:    "Tool configurations available to the agent. Maximum of 128 tools across all toolsets allowed.",
			BodyPath: "tools",
		},
		&requestflag.Flag[[]string]{
			Name:       "beta",
			Usage:      "Optional header to specify the beta version(s) you want to use.",
			HeaderPath: "anthropic-beta",
		},
		&requestflag.Flag[string]{
			Name:       "workspace-id",
			HeaderPath: "anthropic-workspace-id",
		},
	},
	Action:          handleBetaAgentsCreate,
	HideHelpCommand: true,
}, map[string][]requestflag.HasOuterFlag{
	"model": {
		&requestflag.InnerFlag[string]{
			Name:       "model.id",
			Usage:      "The model that will power your agent.\n\nSee [models](https://docs.anthropic.com/en/docs/models-overview) for additional details and options.",
			InnerField: "id",
		},
		&requestflag.InnerFlag[any]{
			Name:       "model.effort",
			InnerField: "effort",
		},
		&requestflag.InnerFlag[*string]{
			Name:       "model.inference-geo",
			Usage:      "Geographic region for model inference. When unset, requests fall through to the workspace's default_inference_geo. On update, `model` is whole-object replacement — omitting inference_geo clears it.",
			InnerField: "inference_geo",
		},
		&requestflag.InnerFlag[*string]{
			Name:       "model.speed",
			Usage:      "Inference speed mode. `fast` provides significantly faster output token generation at premium pricing. Not all models support `fast`; invalid combinations are rejected at create time.",
			InnerField: "speed",
		},
	},
	"mcp-server": {
		&requestflag.InnerFlag[string]{
			Name:       "mcp-server.name",
			Usage:      "Unique name for this server, referenced by mcp_toolset configurations. 1-255 characters.",
			InnerField: "name",
		},
		&requestflag.InnerFlag[string]{
			Name:       "mcp-server.type",
			Usage:      `Allowed values: "url".`,
			InnerField: "type",
		},
		&requestflag.InnerFlag[string]{
			Name:       "mcp-server.url",
			Usage:      "Endpoint URL for the MCP server.",
			InnerField: "url",
		},
	},
	"multiagent": {
		&requestflag.InnerFlag[[]any]{
			Name:       "multiagent.agents",
			Usage:      "Agents the coordinator may spawn as session threads. 1–20 entries. Each entry is an agent ID string, a versioned `{\"type\":\"agent\",\"id\",\"version\"}` reference, or `{\"type\":\"self\"}` to allow recursive self-invocation. Entries must reference distinct agents (after resolving `self` and string forms); at most one `self`. Referenced agents must exist, must not be archived, and must not themselves have `multiagent` set (depth limit 1).",
			InnerField: "agents",
		},
		&requestflag.InnerFlag[string]{
			Name:       "multiagent.type",
			Usage:      `Allowed values: "coordinator".`,
			InnerField: "type",
		},
	},
	"skill": {
		&requestflag.InnerFlag[string]{
			Name:       "skill.skill-id",
			Usage:      `Identifier of the Anthropic skill (e.g., "xlsx").`,
			InnerField: "skill_id",
		},
		&requestflag.InnerFlag[string]{
			Name:       "skill.type",
			Usage:      `Allowed values: "anthropic", "custom".`,
			InnerField: "type",
		},
		&requestflag.InnerFlag[*string]{
			Name:       "skill.version",
			Usage:      "Version to pin. Defaults to latest if omitted.",
			InnerField: "version",
		},
	},
	"tool": {
		&requestflag.InnerFlag[string]{
			Name:       "tool.type",
			Usage:      `Allowed values: "agent_toolset_20260401", "mcp_toolset", "custom".`,
			InnerField: "type",
		},
		&requestflag.InnerFlag[any]{
			Name:       "tool.configs",
			InnerField: "configs",
		},
		&requestflag.InnerFlag[any]{
			Name:       "tool.default-config",
			InnerField: "default_config",
		},
		&requestflag.InnerFlag[string]{
			Name:       "tool.description",
			Usage:      "Description of what the tool does, shown to the agent to help it decide when to use the tool.",
			InnerField: "description",
		},
		&requestflag.InnerFlag[map[string]any]{
			Name:       "tool.input-schema",
			Usage:      "JSON Schema for custom tool input parameters.",
			InnerField: "input_schema",
		},
		&requestflag.InnerFlag[string]{
			Name:       "tool.mcp-server-name",
			Usage:      "Name of the MCP server. Must match a server name from the mcp_servers array. 1-255 characters.",
			InnerField: "mcp_server_name",
		},
		&requestflag.InnerFlag[string]{
			Name:       "tool.name",
			Usage:      "Unique name for the tool. 1-128 characters; letters, digits, underscores, and hyphens.",
			InnerField: "name",
		},
	},
})

var betaAgentsRetrieve = cli.Command{
	Name:    "retrieve",
	Usage:   "Get Agent",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "agent-id",
			Required:  true,
			PathParam: "agent_id",
		},
		&requestflag.Flag[int64]{
			Name:      "version",
			Usage:     "Agent version. Omit for the most recent version. Must be at least 1 if specified.",
			QueryPath: "version",
		},
		&requestflag.Flag[[]string]{
			Name:       "beta",
			Usage:      "Optional header to specify the beta version(s) you want to use.",
			HeaderPath: "anthropic-beta",
		},
		&requestflag.Flag[string]{
			Name:       "workspace-id",
			HeaderPath: "anthropic-workspace-id",
		},
	},
	Action:          handleBetaAgentsRetrieve,
	HideHelpCommand: true,
}

var betaAgentsUpdate = requestflag.WithInnerFlags(cli.Command{
	Name:    "update",
	Usage:   "Update Agent",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:        "agent-id",
			Required:    true,
			PathParam:   "agent_id",
			DataAliases: []string{"id"},
		},
		&requestflag.Flag[*string]{
			Name:     "description",
			Usage:    "Description. Omit to preserve; send empty string or null to clear.",
			BodyPath: "description",
		},
		&requestflag.Flag[any]{
			Name:     "mcp-server",
			Usage:    "MCP servers. Full replacement. Omit to preserve; send empty array or `null` to clear. Names must be unique. Maximum 20. Every server must be referenced by an `mcp_toolset` in the agent's resulting `tools`; unreferenced servers are rejected. See the [MCP connector guide](https://platform.claude.com/docs/en/managed-agents/mcp-connector).",
			BodyPath: "mcp_servers",
		},
		&requestflag.Flag[map[string]any]{
			Name:     "metadata",
			Usage:    "Metadata patch. Set a key to a string to upsert it, or to null to delete it. Omit the field to preserve. The stored bag is limited to 16 keys (up to 64 chars each) with values up to 512 chars.",
			BodyPath: "metadata",
		},
		&requestflag.Flag[map[string]any]{
			Name:     "model",
			Usage:    "Model identifier. Accepts the [model string](https://platform.claude.com/docs/en/about-claude/models/overview#latest-models-comparison), e.g. `claude-opus-5`, or a `model_config` object for additional configuration control. Omit to preserve. Cannot be cleared.",
			BodyPath: "model",
		},
		&requestflag.Flag[map[string]any]{
			Name:     "multiagent",
			Usage:    "A coordinator topology: the session's primary thread orchestrates work by spawning session threads, each running an agent drawn from the `agents` roster.",
			BodyPath: "multiagent",
		},
		&requestflag.Flag[string]{
			Name:     "name",
			Usage:    "Human-readable name. Must be non-empty. Omit to preserve. Cannot be cleared.",
			BodyPath: "name",
		},
		&requestflag.Flag[any]{
			Name:     "skill",
			Usage:    "Skills. Full replacement. Omit to preserve; send empty array or null to clear.",
			BodyPath: "skills",
		},
		&requestflag.Flag[*string]{
			Name:     "system",
			Usage:    "System prompt. Omit to preserve; send empty string or null to clear.",
			BodyPath: "system",
		},
		&requestflag.Flag[any]{
			Name:     "tool",
			Usage:    "Tool configurations available to the agent. Full replacement. Omit to preserve; send empty array or null to clear. Maximum of 128 tools across all toolsets allowed.",
			BodyPath: "tools",
		},
		&requestflag.Flag[int64]{
			Name:     "version",
			Usage:    "The agent's current version, used to prevent concurrent overwrites. Obtain this value from a create or retrieve response. Must be at least 1 if specified. When supplied, the request fails if it does not match the server's current version; omit to apply the update unconditionally.",
			BodyPath: "version",
		},
		&requestflag.Flag[[]string]{
			Name:       "beta",
			Usage:      "Optional header to specify the beta version(s) you want to use.",
			HeaderPath: "anthropic-beta",
		},
		&requestflag.Flag[string]{
			Name:       "workspace-id",
			HeaderPath: "anthropic-workspace-id",
		},
	},
	Action:          handleBetaAgentsUpdate,
	HideHelpCommand: true,
}, map[string][]requestflag.HasOuterFlag{
	"mcp-server": {
		&requestflag.InnerFlag[string]{
			Name:                  "mcp-server.name",
			Usage:                 "Unique name for this server, referenced by mcp_toolset configurations. 1-255 characters.",
			InnerField:            "name",
			OuterIsArrayOfObjects: true,
		},
		&requestflag.InnerFlag[string]{
			Name:                  "mcp-server.type",
			Usage:                 `Allowed values: "url".`,
			InnerField:            "type",
			OuterIsArrayOfObjects: true,
		},
		&requestflag.InnerFlag[string]{
			Name:                  "mcp-server.url",
			Usage:                 "Endpoint URL for the MCP server.",
			InnerField:            "url",
			OuterIsArrayOfObjects: true,
		},
	},
	"model": {
		&requestflag.InnerFlag[string]{
			Name:       "model.id",
			Usage:      "The model that will power your agent.\n\nSee [models](https://docs.anthropic.com/en/docs/models-overview) for additional details and options.",
			InnerField: "id",
		},
		&requestflag.InnerFlag[any]{
			Name:       "model.effort",
			InnerField: "effort",
		},
		&requestflag.InnerFlag[*string]{
			Name:       "model.inference-geo",
			Usage:      "Geographic region for model inference. When unset, requests fall through to the workspace's default_inference_geo. On update, `model` is whole-object replacement — omitting inference_geo clears it.",
			InnerField: "inference_geo",
		},
		&requestflag.InnerFlag[*string]{
			Name:       "model.speed",
			Usage:      "Inference speed mode. `fast` provides significantly faster output token generation at premium pricing. Not all models support `fast`; invalid combinations are rejected at create time.",
			InnerField: "speed",
		},
	},
	"multiagent": {
		&requestflag.InnerFlag[[]any]{
			Name:       "multiagent.agents",
			Usage:      "Agents the coordinator may spawn as session threads. 1–20 entries. Each entry is an agent ID string, a versioned `{\"type\":\"agent\",\"id\",\"version\"}` reference, or `{\"type\":\"self\"}` to allow recursive self-invocation. Entries must reference distinct agents (after resolving `self` and string forms); at most one `self`. Referenced agents must exist, must not be archived, and must not themselves have `multiagent` set (depth limit 1).",
			InnerField: "agents",
		},
		&requestflag.InnerFlag[string]{
			Name:       "multiagent.type",
			Usage:      `Allowed values: "coordinator".`,
			InnerField: "type",
		},
	},
	"skill": {
		&requestflag.InnerFlag[string]{
			Name:                  "skill.skill-id",
			Usage:                 `Identifier of the Anthropic skill (e.g., "xlsx").`,
			InnerField:            "skill_id",
			OuterIsArrayOfObjects: true,
		},
		&requestflag.InnerFlag[string]{
			Name:                  "skill.type",
			Usage:                 `Allowed values: "anthropic", "custom".`,
			InnerField:            "type",
			OuterIsArrayOfObjects: true,
		},
		&requestflag.InnerFlag[*string]{
			Name:                  "skill.version",
			Usage:                 "Version to pin. Defaults to latest if omitted.",
			InnerField:            "version",
			OuterIsArrayOfObjects: true,
		},
	},
	"tool": {
		&requestflag.InnerFlag[string]{
			Name:                  "tool.type",
			Usage:                 `Allowed values: "agent_toolset_20260401", "mcp_toolset", "custom".`,
			InnerField:            "type",
			OuterIsArrayOfObjects: true,
		},
		&requestflag.InnerFlag[any]{
			Name:                  "tool.configs",
			InnerField:            "configs",
			OuterIsArrayOfObjects: true,
		},
		&requestflag.InnerFlag[any]{
			Name:                  "tool.default-config",
			InnerField:            "default_config",
			OuterIsArrayOfObjects: true,
		},
		&requestflag.InnerFlag[string]{
			Name:                  "tool.description",
			Usage:                 "Description of what the tool does, shown to the agent to help it decide when to use the tool.",
			InnerField:            "description",
			OuterIsArrayOfObjects: true,
		},
		&requestflag.InnerFlag[map[string]any]{
			Name:                  "tool.input-schema",
			Usage:                 "JSON Schema for custom tool input parameters.",
			InnerField:            "input_schema",
			OuterIsArrayOfObjects: true,
		},
		&requestflag.InnerFlag[string]{
			Name:                  "tool.mcp-server-name",
			Usage:                 "Name of the MCP server. Must match a server name from the mcp_servers array. 1-255 characters.",
			InnerField:            "mcp_server_name",
			OuterIsArrayOfObjects: true,
		},
		&requestflag.InnerFlag[string]{
			Name:                  "tool.name",
			Usage:                 "Unique name for the tool. 1-128 characters; letters, digits, underscores, and hyphens.",
			InnerField:            "name",
			OuterIsArrayOfObjects: true,
		},
	},
})

var betaAgentsList = cli.Command{
	Name:    "list",
	Usage:   "List Agents",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[any]{
			Name:      "created-at-gte",
			Usage:     "Return agents created at or after this time (inclusive).",
			QueryPath: "created_at[gte]",
		},
		&requestflag.Flag[any]{
			Name:      "created-at-lte",
			Usage:     "Return agents created at or before this time (inclusive).",
			QueryPath: "created_at[lte]",
		},
		&requestflag.Flag[bool]{
			Name:      "include-archived",
			Usage:     "Include archived agents in results. Defaults to false.",
			QueryPath: "include_archived",
		},
		&requestflag.Flag[int64]{
			Name:      "limit",
			Usage:     "Maximum results per page. Default 20, maximum 100.",
			QueryPath: "limit",
		},
		&requestflag.Flag[string]{
			Name:      "page",
			Usage:     "Opaque pagination cursor from a previous response.",
			QueryPath: "page",
		},
		&requestflag.Flag[[]string]{
			Name:       "beta",
			Usage:      "Optional header to specify the beta version(s) you want to use.",
			HeaderPath: "anthropic-beta",
		},
		&requestflag.Flag[string]{
			Name:       "workspace-id",
			HeaderPath: "anthropic-workspace-id",
		},
		&requestflag.Flag[int64]{
			Name:  "max-items",
			Usage: "The maximum number of items to return (use -1 for unlimited).",
		},
	},
	Action:          handleBetaAgentsList,
	HideHelpCommand: true,
}

var betaAgentsArchive = cli.Command{
	Name:    "archive",
	Usage:   "Archive Agent",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "agent-id",
			Required:  true,
			PathParam: "agent_id",
		},
		&requestflag.Flag[[]string]{
			Name:       "beta",
			Usage:      "Optional header to specify the beta version(s) you want to use.",
			HeaderPath: "anthropic-beta",
		},
		&requestflag.Flag[string]{
			Name:       "workspace-id",
			HeaderPath: "anthropic-workspace-id",
		},
	},
	Action:          handleBetaAgentsArchive,
	HideHelpCommand: true,
}

func handleBetaAgentsCreate(ctx context.Context, cmd *cli.Command) error {
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

	params := anthropic.BetaAgentNewParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Beta.Agents.New(ctx, params, options...)
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
		Title:          "beta:agents create",
		Transform:      transform,
	})
}

func handleBetaAgentsRetrieve(ctx context.Context, cmd *cli.Command) error {
	client := anthropic.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("agent-id") && len(unusedArgs) > 0 {
		cmd.Set("agent-id", unusedArgs[0])
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

	params := anthropic.BetaAgentGetParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Beta.Agents.Get(
		ctx,
		cmd.Value("agent-id").(string),
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
		Title:          "beta:agents retrieve",
		Transform:      transform,
	})
}

func handleBetaAgentsUpdate(ctx context.Context, cmd *cli.Command) error {
	client := anthropic.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("agent-id") && len(unusedArgs) > 0 {
		cmd.Set("agent-id", unusedArgs[0])
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

	params := anthropic.BetaAgentUpdateParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Beta.Agents.Update(
		ctx,
		cmd.Value("agent-id").(string),
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
		Title:          "beta:agents update",
		Transform:      transform,
	})
}

func handleBetaAgentsList(ctx context.Context, cmd *cli.Command) error {
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

	params := anthropic.BetaAgentListParams{}

	format := "explore"
	explicitFormat := cmd.Root().IsSet("format")
	if explicitFormat {
		format = cmd.Root().String("format")
	}
	transform := cmd.Root().String("transform")
	if format == "raw" {
		var res []byte
		options = append(options, option.WithResponseBodyInto(&res))
		_, err = client.Beta.Agents.List(ctx, params, options...)
		if err != nil {
			return err
		}
		obj := gjson.ParseBytes(res)
		return ShowJSON(obj, ShowJSONOpts{
			ExplicitFormat: explicitFormat,
			Format:         format,
			RawOutput:      cmd.Root().Bool("raw-output"),
			Title:          "beta:agents list",
			Transform:      transform,
		})
	} else {
		iter := client.Beta.Agents.ListAutoPaging(ctx, params, options...)
		maxItems := int64(-1)
		if cmd.IsSet("max-items") {
			maxItems = cmd.Value("max-items").(int64)
		}
		return ShowJSONIterator(iter, maxItems, ShowJSONOpts{
			ExplicitFormat: explicitFormat,
			Format:         format,
			RawOutput:      cmd.Root().Bool("raw-output"),
			Title:          "beta:agents list",
			Transform:      transform,
		})
	}
}

func handleBetaAgentsArchive(ctx context.Context, cmd *cli.Command) error {
	client := anthropic.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("agent-id") && len(unusedArgs) > 0 {
		cmd.Set("agent-id", unusedArgs[0])
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

	params := anthropic.BetaAgentArchiveParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Beta.Agents.Archive(
		ctx,
		cmd.Value("agent-id").(string),
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
		Title:          "beta:agents archive",
		Transform:      transform,
	})
}
