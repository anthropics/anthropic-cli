// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cmd

import (
	"context"
	"fmt"
	"os"

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
			Usage:    "Model identifier. Accepts the [model string](https://platform.claude.com/docs/en/about-claude/models/overview#latest-models-comparison), e.g. `claude-opus-4-6`, or a `model_config` object for additional configuration control",
			Required: true,
			BodyPath: "model",
		},
		&requestflag.Flag[string]{
			Name:     "name",
			Usage:    "Human-readable name for the agent. 1-256 characters.",
			Required: true,
			BodyPath: "name",
		},
		&requestflag.Flag[any]{
			Name:     "description",
			Usage:    "Description of what the agent does. Up to 2048 characters.",
			BodyPath: "description",
		},
		&requestflag.Flag[[]map[string]any]{
			Name:     "mcp-server",
			Usage:    "MCP servers this agent connects to. Maximum 20. Names must be unique within the array.",
			BodyPath: "mcp_servers",
		},
		&requestflag.Flag[map[string]any]{
			Name:     "metadata",
			Usage:    "Arbitrary key-value metadata. Maximum 16 pairs, keys up to 64 chars, values up to 512 chars.",
			BodyPath: "metadata",
		},
		&requestflag.Flag[[]map[string]any]{
			Name:     "skill",
			Usage:    "Skills available to the agent. Maximum 20.",
			BodyPath: "skills",
		},
		&requestflag.Flag[any]{
			Name:     "system",
			Usage:    "System prompt for the agent. Up to 100,000 characters.",
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
	},
	Action:          handleBetaAgentsCreate,
	HideHelpCommand: true,
}, map[string][]requestflag.HasOuterFlag{
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
})

var betaAgentsRetrieve = cli.Command{
	Name:    "retrieve",
	Usage:   "Get Agent",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:     "agent-id",
			Required: true,
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
			Name:     "agent-id",
			Required: true,
		},
		&requestflag.Flag[int64]{
			Name:     "version",
			Usage:    "The agent's current version, used to prevent concurrent overwrites. Obtain this value from a create or retrieve response. The request fails if this does not match the server's current version.",
			Required: true,
			BodyPath: "version",
		},
		&requestflag.Flag[any]{
			Name:     "description",
			Usage:    "Description. Up to 2048 characters. Omit to preserve; send empty string or null to clear.",
			BodyPath: "description",
		},
		&requestflag.Flag[any]{
			Name:     "mcp-server",
			Usage:    "MCP servers. Full replacement. Omit to preserve; send empty array or null to clear. Names must be unique. Maximum 20.",
			BodyPath: "mcp_servers",
		},
		&requestflag.Flag[map[string]any]{
			Name:     "metadata",
			Usage:    "Metadata patch. Set a key to a string to upsert it, or to null to delete it. Omit the field to preserve. The stored bag is limited to 16 keys (up to 64 chars each) with values up to 512 chars.",
			BodyPath: "metadata",
		},
		&requestflag.Flag[map[string]any]{
			Name:     "model",
			Usage:    "Model identifier. Accepts the [model string](https://platform.claude.com/docs/en/about-claude/models/overview#latest-models-comparison), e.g. `claude-opus-4-6`, or a `model_config` object for additional configuration control. Omit to preserve. Cannot be cleared.",
			BodyPath: "model",
		},
		&requestflag.Flag[string]{
			Name:     "name",
			Usage:    "Human-readable name. 1-256 characters. Omit to preserve. Cannot be cleared.",
			BodyPath: "name",
		},
		&requestflag.Flag[any]{
			Name:     "skill",
			Usage:    "Skills. Full replacement. Omit to preserve; send empty array or null to clear. Maximum 20.",
			BodyPath: "skills",
		},
		&requestflag.Flag[any]{
			Name:     "system",
			Usage:    "System prompt. Up to 100,000 characters. Omit to preserve; send empty string or null to clear.",
			BodyPath: "system",
		},
		&requestflag.Flag[any]{
			Name:     "tool",
			Usage:    "Tool configurations available to the agent. Full replacement. Omit to preserve; send empty array or null to clear. Maximum of 128 tools across all toolsets allowed.",
			BodyPath: "tools",
		},
		&requestflag.Flag[[]string]{
			Name:       "beta",
			Usage:      "Optional header to specify the beta version(s) you want to use.",
			HeaderPath: "anthropic-beta",
		},
	},
	Action:          handleBetaAgentsUpdate,
	HideHelpCommand: true,
}, map[string][]requestflag.HasOuterFlag{
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
			Name:     "agent-id",
			Required: true,
		},
		&requestflag.Flag[[]string]{
			Name:       "beta",
			Usage:      "Optional header to specify the beta version(s) you want to use.",
			HeaderPath: "anthropic-beta",
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

	params := anthropic.BetaAgentNewParams{}

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
	_, err = client.Beta.Agents.New(ctx, params, options...)
	if err != nil {
		return err
	}

	obj := gjson.ParseBytes(res)
	format := cmd.Root().String("format")
	explicitFormat := cmd.Root().IsSet("format")
	transform := cmd.Root().String("transform")
	return ShowJSON(os.Stdout, os.Stderr, "beta:agents create", obj, format, explicitFormat, transform)
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

	params := anthropic.BetaAgentGetParams{}

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
	format := cmd.Root().String("format")
	explicitFormat := cmd.Root().IsSet("format")
	transform := cmd.Root().String("transform")
	return ShowJSON(os.Stdout, os.Stderr, "beta:agents retrieve", obj, format, explicitFormat, transform)
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

	params := anthropic.BetaAgentUpdateParams{}

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
	return ShowJSON(os.Stdout, os.Stderr, "beta:agents update", obj, format, explicitFormat, transform)
}

func handleBetaAgentsList(ctx context.Context, cmd *cli.Command) error {
	client := anthropic.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()

	if len(unusedArgs) > 0 {
		return fmt.Errorf("Unexpected extra arguments: %v", unusedArgs)
	}

	params := anthropic.BetaAgentListParams{}

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
	explicitFormat := cmd.Root().IsSet("format")
	transform := cmd.Root().String("transform")
	if format == "raw" {
		var res []byte
		options = append(options, option.WithResponseBodyInto(&res))
		_, err = client.Beta.Agents.List(ctx, params, options...)
		if err != nil {
			return err
		}
		obj := gjson.ParseBytes(res)
		return ShowJSON(os.Stdout, os.Stderr, "beta:agents list", obj, format, explicitFormat, transform)
	} else {
		iter := client.Beta.Agents.ListAutoPaging(ctx, params, options...)
		maxItems := int64(-1)
		if cmd.IsSet("max-items") {
			maxItems = cmd.Value("max-items").(int64)
		}
		return ShowJSONIterator(os.Stdout, os.Stderr, "beta:agents list", iter, format, explicitFormat, transform, maxItems)
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

	params := anthropic.BetaAgentArchiveParams{}

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
	return ShowJSON(os.Stdout, os.Stderr, "beta:agents archive", obj, format, explicitFormat, transform)
}
