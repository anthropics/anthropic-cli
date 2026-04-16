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

var messagesCreate = requestflag.WithInnerFlags(cli.Command{
	Name:    "create",
	Usage:   "Send a structured list of input messages with text and/or image content, and the\nmodel will generate the next message in the conversation.",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[int64]{
			Name:     "max-tokens",
			Usage:    "The maximum number of tokens to generate before stopping.\n\nNote that our models may stop _before_ reaching this maximum. This parameter only specifies the absolute maximum number of tokens to generate.\n\nDifferent models have different maximum values for this parameter.  See [models](https://docs.claude.com/en/docs/models-overview) for details.",
			Required: true,
			BodyPath: "max_tokens",
		},
		&requestflag.Flag[[]map[string]any]{
			Name:     "message",
			Usage:    "Input messages.\n\nOur models are trained to operate on alternating `user` and `assistant` conversational turns. When creating a new `Message`, you specify the prior conversational turns with the `messages` parameter, and the model then generates the next `Message` in the conversation. Consecutive `user` or `assistant` turns in your request will be combined into a single turn.\n\nEach input message must be an object with a `role` and `content`. You can specify a single `user`-role message, or you can include multiple `user` and `assistant` messages.\n\nIf the final message uses the `assistant` role, the response content will continue immediately from the content in that message. This can be used to constrain part of the model's response.\n\nExample with a single `user` message:\n\n```json\n[{\"role\": \"user\", \"content\": \"Hello, Claude\"}]\n```\n\nExample with multiple conversational turns:\n\n```json\n[\n  {\"role\": \"user\", \"content\": \"Hello there.\"},\n  {\"role\": \"assistant\", \"content\": \"Hi, I'm Claude. How can I help you?\"},\n  {\"role\": \"user\", \"content\": \"Can you explain LLMs in plain English?\"},\n]\n```\n\nExample with a partially-filled response from Claude:\n\n```json\n[\n  {\"role\": \"user\", \"content\": \"What's the Greek name for Sun? (A) Sol (B) Helios (C) Sun\"},\n  {\"role\": \"assistant\", \"content\": \"The best answer is (\"},\n]\n```\n\nEach input message `content` may be either a single `string` or an array of content blocks, where each block has a specific `type`. Using a `string` for `content` is shorthand for an array of one content block of type `\"text\"`. The following input messages are equivalent:\n\n```json\n{\"role\": \"user\", \"content\": \"Hello, Claude\"}\n```\n\n```json\n{\"role\": \"user\", \"content\": [{\"type\": \"text\", \"text\": \"Hello, Claude\"}]}\n```\n\nSee [input examples](https://docs.claude.com/en/api/messages-examples).\n\nNote that if you want to include a [system prompt](https://docs.claude.com/en/docs/system-prompts), you can use the top-level `system` parameter — there is no `\"system\"` role for input messages in the Messages API.\n\nThere is a limit of 100,000 messages in a single request.",
			Required: true,
			BodyPath: "messages",
		},
		&requestflag.Flag[string]{
			Name:     "model",
			Usage:    "The model that will complete your prompt.\\n\\nSee [models](https://docs.anthropic.com/en/docs/models-overview) for additional details and options.",
			Required: true,
			BodyPath: "model",
		},
		&requestflag.Flag[map[string]any]{
			Name:     "cache-control",
			BodyPath: "cache_control",
		},
		&requestflag.Flag[any]{
			Name:     "container",
			Usage:    "Container identifier for reuse across requests.",
			BodyPath: "container",
		},
		&requestflag.Flag[any]{
			Name:     "inference-geo",
			Usage:    "Specifies the geographic region for inference processing. If not specified, the workspace's `default_inference_geo` is used.",
			BodyPath: "inference_geo",
		},
		&requestflag.Flag[map[string]any]{
			Name:     "metadata",
			BodyPath: "metadata",
		},
		&requestflag.Flag[map[string]any]{
			Name:     "output-config",
			BodyPath: "output_config",
		},
		&requestflag.Flag[string]{
			Name:     "service-tier",
			Usage:    "Determines whether to use priority capacity (if available) or standard capacity for this request.\n\nAnthropic offers different levels of service for your API requests. See [service-tiers](https://docs.claude.com/en/api/service-tiers) for details.",
			BodyPath: "service_tier",
		},
		&requestflag.Flag[[]string]{
			Name:     "stop-sequence",
			Usage:    "Custom text sequences that will cause the model to stop generating.\n\nOur models will normally stop when they have naturally completed their turn, which will result in a response `stop_reason` of `\"end_turn\"`.\n\nIf you want the model to stop generating when it encounters custom strings of text, you can use the `stop_sequences` parameter. If the model encounters one of the custom sequences, the response `stop_reason` value will be `\"stop_sequence\"` and the response `stop_sequence` value will contain the matched stop sequence.",
			BodyPath: "stop_sequences",
		},
		&requestflag.Flag[bool]{
			Name:     "stream",
			Usage:    "Whether to incrementally stream the response using server-sent events.\n\nSee [streaming](https://docs.claude.com/en/api/messages-streaming) for details.",
			BodyPath: "stream",
		},
		&requestflag.Flag[any]{
			Name:     "system",
			Usage:    "System prompt.\n\nA system prompt is a way of providing context and instructions to Claude, such as specifying a particular goal or role. See our [guide to system prompts](https://docs.claude.com/en/docs/system-prompts).",
			BodyPath: "system",
		},
		&requestflag.Flag[float64]{
			Name:     "temperature",
			Usage:    "Amount of randomness injected into the response.\n\nDefaults to `1.0`. Ranges from `0.0` to `1.0`. Use `temperature` closer to `0.0` for analytical / multiple choice, and closer to `1.0` for creative and generative tasks.\n\nNote that even with `temperature` of `0.0`, the results will not be fully deterministic.",
			BodyPath: "temperature",
		},
		&requestflag.Flag[map[string]any]{
			Name:     "thinking",
			Usage:    "Configuration for enabling Claude's extended thinking.\n\nWhen enabled, responses include `thinking` content blocks showing Claude's thinking process before the final answer. Requires a minimum budget of 1,024 tokens and counts towards your `max_tokens` limit.\n\nSee [extended thinking](https://docs.claude.com/en/docs/build-with-claude/extended-thinking) for details.",
			BodyPath: "thinking",
		},
		&requestflag.Flag[map[string]any]{
			Name:     "tool-choice",
			Usage:    "How the model should use the provided tools. The model can use a specific tool, any available tool, decide by itself, or not use tools at all.",
			BodyPath: "tool_choice",
		},
		&requestflag.Flag[[]map[string]any]{
			Name:     "tool",
			Usage:    "Definitions of tools that the model may use.\n\nIf you include `tools` in your API request, the model may return `tool_use` content blocks that represent the model's use of those tools. You can then run those tools using the tool input generated by the model and then optionally return results back to the model using `tool_result` content blocks.\n\nThere are two types of tools: **client tools** and **server tools**. The behavior described below applies to client tools. For [server tools](https://docs.claude.com/en/docs/agents-and-tools/tool-use/overview\\#server-tools), see their individual documentation as each has its own behavior (e.g., the [web search tool](https://docs.claude.com/en/docs/agents-and-tools/tool-use/web-search-tool)).\n\nEach tool definition includes:\n\n* `name`: Name of the tool.\n* `description`: Optional, but strongly-recommended description of the tool.\n* `input_schema`: [JSON schema](https://json-schema.org/draft/2020-12) for the tool `input` shape that the model will produce in `tool_use` output content blocks.\n\nFor example, if you defined `tools` as:\n\n```json\n[\n  {\n    \"name\": \"get_stock_price\",\n    \"description\": \"Get the current stock price for a given ticker symbol.\",\n    \"input_schema\": {\n      \"type\": \"object\",\n      \"properties\": {\n        \"ticker\": {\n          \"type\": \"string\",\n          \"description\": \"The stock ticker symbol, e.g. AAPL for Apple Inc.\"\n        }\n      },\n      \"required\": [\"ticker\"]\n    }\n  }\n]\n```\n\nAnd then asked the model \"What's the S&P 500 at today?\", the model might produce `tool_use` content blocks in the response like this:\n\n```json\n[\n  {\n    \"type\": \"tool_use\",\n    \"id\": \"toolu_01D7FLrfh4GYq7yT1ULFeyMV\",\n    \"name\": \"get_stock_price\",\n    \"input\": { \"ticker\": \"^GSPC\" }\n  }\n]\n```\n\nYou might then run your `get_stock_price` tool with `{\"ticker\": \"^GSPC\"}` as an input, and return the following back to the model in a subsequent `user` message:\n\n```json\n[\n  {\n    \"type\": \"tool_result\",\n    \"tool_use_id\": \"toolu_01D7FLrfh4GYq7yT1ULFeyMV\",\n    \"content\": \"259.75 USD\"\n  }\n]\n```\n\nTools can be used for workflows that include running client-side tools and functions, or more generally whenever you want the model to produce a particular JSON structure of output.\n\nSee our [guide](https://docs.claude.com/en/docs/tool-use) for more details.",
			BodyPath: "tools",
		},
		&requestflag.Flag[int64]{
			Name:     "top-k",
			Usage:    "Only sample from the top K options for each subsequent token.\n\nUsed to remove \"long tail\" low probability responses. [Learn more technical details here](https://towardsdatascience.com/how-to-sample-from-language-models-682bceb97277).\n\nRecommended for advanced use cases only. You usually only need to use `temperature`.",
			BodyPath: "top_k",
		},
		&requestflag.Flag[float64]{
			Name:     "top-p",
			Usage:    "Use nucleus sampling.\n\nIn nucleus sampling, we compute the cumulative distribution over all the options for each subsequent token in decreasing probability order and cut it off once it reaches a particular probability specified by `top_p`. You should either alter `temperature` or `top_p`, but not both.\n\nRecommended for advanced use cases only. You usually only need to use `temperature`.",
			BodyPath: "top_p",
		},
		&requestflag.Flag[int64]{
			Name:  "max-items",
			Usage: "The maximum number of items to return (use -1 for unlimited).",
		},
	},
	Action:          handleMessagesCreate,
	HideHelpCommand: true,
}, map[string][]requestflag.HasOuterFlag{
	"message": {
		&requestflag.InnerFlag[[]map[string]any]{
			Name:       "message.content",
			InnerField: "content",
		},
		&requestflag.InnerFlag[string]{
			Name:       "message.role",
			Usage:      `Allowed values: "user", "assistant".`,
			InnerField: "role",
		},
	},
	"cache-control": {
		&requestflag.InnerFlag[string]{
			Name:       "cache-control.type",
			Usage:      `Allowed values: "ephemeral".`,
			InnerField: "type",
		},
		&requestflag.InnerFlag[string]{
			Name:       "cache-control.ttl",
			Usage:      "The time-to-live for the cache control breakpoint.\n\nThis may be one the following values:\n- `5m`: 5 minutes\n- `1h`: 1 hour\n\nDefaults to `5m`.",
			InnerField: "ttl",
		},
	},
	"metadata": {
		&requestflag.InnerFlag[any]{
			Name:       "metadata.user-id",
			Usage:      "An external identifier for the user who is associated with the request.\n\nThis should be a uuid, hash value, or other opaque identifier. Anthropic may use this id to help detect abuse. Do not include any identifying information such as name, email address, or phone number.",
			InnerField: "user_id",
		},
	},
	"output-config": {
		&requestflag.InnerFlag[any]{
			Name:       "output-config.effort",
			Usage:      "All possible effort levels.",
			InnerField: "effort",
		},
		&requestflag.InnerFlag[map[string]any]{
			Name:       "output-config.format",
			InnerField: "format",
		},
	},
})

var messagesCountTokens = requestflag.WithInnerFlags(cli.Command{
	Name:    "count-tokens",
	Usage:   "Count the number of tokens in a Message.",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[[]map[string]any]{
			Name:     "message",
			Usage:    "Input messages.\n\nOur models are trained to operate on alternating `user` and `assistant` conversational turns. When creating a new `Message`, you specify the prior conversational turns with the `messages` parameter, and the model then generates the next `Message` in the conversation. Consecutive `user` or `assistant` turns in your request will be combined into a single turn.\n\nEach input message must be an object with a `role` and `content`. You can specify a single `user`-role message, or you can include multiple `user` and `assistant` messages.\n\nIf the final message uses the `assistant` role, the response content will continue immediately from the content in that message. This can be used to constrain part of the model's response.\n\nExample with a single `user` message:\n\n```json\n[{\"role\": \"user\", \"content\": \"Hello, Claude\"}]\n```\n\nExample with multiple conversational turns:\n\n```json\n[\n  {\"role\": \"user\", \"content\": \"Hello there.\"},\n  {\"role\": \"assistant\", \"content\": \"Hi, I'm Claude. How can I help you?\"},\n  {\"role\": \"user\", \"content\": \"Can you explain LLMs in plain English?\"},\n]\n```\n\nExample with a partially-filled response from Claude:\n\n```json\n[\n  {\"role\": \"user\", \"content\": \"What's the Greek name for Sun? (A) Sol (B) Helios (C) Sun\"},\n  {\"role\": \"assistant\", \"content\": \"The best answer is (\"},\n]\n```\n\nEach input message `content` may be either a single `string` or an array of content blocks, where each block has a specific `type`. Using a `string` for `content` is shorthand for an array of one content block of type `\"text\"`. The following input messages are equivalent:\n\n```json\n{\"role\": \"user\", \"content\": \"Hello, Claude\"}\n```\n\n```json\n{\"role\": \"user\", \"content\": [{\"type\": \"text\", \"text\": \"Hello, Claude\"}]}\n```\n\nSee [input examples](https://docs.claude.com/en/api/messages-examples).\n\nNote that if you want to include a [system prompt](https://docs.claude.com/en/docs/system-prompts), you can use the top-level `system` parameter — there is no `\"system\"` role for input messages in the Messages API.\n\nThere is a limit of 100,000 messages in a single request.",
			Required: true,
			BodyPath: "messages",
		},
		&requestflag.Flag[string]{
			Name:     "model",
			Usage:    "The model that will complete your prompt.\\n\\nSee [models](https://docs.anthropic.com/en/docs/models-overview) for additional details and options.",
			Required: true,
			BodyPath: "model",
		},
		&requestflag.Flag[map[string]any]{
			Name:     "cache-control",
			BodyPath: "cache_control",
		},
		&requestflag.Flag[map[string]any]{
			Name:     "output-config",
			BodyPath: "output_config",
		},
		&requestflag.Flag[any]{
			Name:     "system",
			Usage:    "System prompt.\n\nA system prompt is a way of providing context and instructions to Claude, such as specifying a particular goal or role. See our [guide to system prompts](https://docs.claude.com/en/docs/system-prompts).",
			BodyPath: "system",
		},
		&requestflag.Flag[map[string]any]{
			Name:     "thinking",
			Usage:    "Configuration for enabling Claude's extended thinking.\n\nWhen enabled, responses include `thinking` content blocks showing Claude's thinking process before the final answer. Requires a minimum budget of 1,024 tokens and counts towards your `max_tokens` limit.\n\nSee [extended thinking](https://docs.claude.com/en/docs/build-with-claude/extended-thinking) for details.",
			BodyPath: "thinking",
		},
		&requestflag.Flag[map[string]any]{
			Name:     "tool-choice",
			Usage:    "How the model should use the provided tools. The model can use a specific tool, any available tool, decide by itself, or not use tools at all.",
			BodyPath: "tool_choice",
		},
		&requestflag.Flag[[]map[string]any]{
			Name:     "tool",
			Usage:    "Definitions of tools that the model may use.\n\nIf you include `tools` in your API request, the model may return `tool_use` content blocks that represent the model's use of those tools. You can then run those tools using the tool input generated by the model and then optionally return results back to the model using `tool_result` content blocks.\n\nThere are two types of tools: **client tools** and **server tools**. The behavior described below applies to client tools. For [server tools](https://docs.claude.com/en/docs/agents-and-tools/tool-use/overview\\#server-tools), see their individual documentation as each has its own behavior (e.g., the [web search tool](https://docs.claude.com/en/docs/agents-and-tools/tool-use/web-search-tool)).\n\nEach tool definition includes:\n\n* `name`: Name of the tool.\n* `description`: Optional, but strongly-recommended description of the tool.\n* `input_schema`: [JSON schema](https://json-schema.org/draft/2020-12) for the tool `input` shape that the model will produce in `tool_use` output content blocks.\n\nFor example, if you defined `tools` as:\n\n```json\n[\n  {\n    \"name\": \"get_stock_price\",\n    \"description\": \"Get the current stock price for a given ticker symbol.\",\n    \"input_schema\": {\n      \"type\": \"object\",\n      \"properties\": {\n        \"ticker\": {\n          \"type\": \"string\",\n          \"description\": \"The stock ticker symbol, e.g. AAPL for Apple Inc.\"\n        }\n      },\n      \"required\": [\"ticker\"]\n    }\n  }\n]\n```\n\nAnd then asked the model \"What's the S&P 500 at today?\", the model might produce `tool_use` content blocks in the response like this:\n\n```json\n[\n  {\n    \"type\": \"tool_use\",\n    \"id\": \"toolu_01D7FLrfh4GYq7yT1ULFeyMV\",\n    \"name\": \"get_stock_price\",\n    \"input\": { \"ticker\": \"^GSPC\" }\n  }\n]\n```\n\nYou might then run your `get_stock_price` tool with `{\"ticker\": \"^GSPC\"}` as an input, and return the following back to the model in a subsequent `user` message:\n\n```json\n[\n  {\n    \"type\": \"tool_result\",\n    \"tool_use_id\": \"toolu_01D7FLrfh4GYq7yT1ULFeyMV\",\n    \"content\": \"259.75 USD\"\n  }\n]\n```\n\nTools can be used for workflows that include running client-side tools and functions, or more generally whenever you want the model to produce a particular JSON structure of output.\n\nSee our [guide](https://docs.claude.com/en/docs/tool-use) for more details.",
			BodyPath: "tools",
		},
	},
	Action:          handleMessagesCountTokens,
	HideHelpCommand: true,
}, map[string][]requestflag.HasOuterFlag{
	"message": {
		&requestflag.InnerFlag[[]map[string]any]{
			Name:       "message.content",
			InnerField: "content",
		},
		&requestflag.InnerFlag[string]{
			Name:       "message.role",
			Usage:      `Allowed values: "user", "assistant".`,
			InnerField: "role",
		},
	},
	"cache-control": {
		&requestflag.InnerFlag[string]{
			Name:       "cache-control.type",
			Usage:      `Allowed values: "ephemeral".`,
			InnerField: "type",
		},
		&requestflag.InnerFlag[string]{
			Name:       "cache-control.ttl",
			Usage:      "The time-to-live for the cache control breakpoint.\n\nThis may be one the following values:\n- `5m`: 5 minutes\n- `1h`: 1 hour\n\nDefaults to `5m`.",
			InnerField: "ttl",
		},
	},
	"output-config": {
		&requestflag.InnerFlag[any]{
			Name:       "output-config.effort",
			Usage:      "All possible effort levels.",
			InnerField: "effort",
		},
		&requestflag.InnerFlag[map[string]any]{
			Name:       "output-config.format",
			InnerField: "format",
		},
	},
})

func handleMessagesCreate(ctx context.Context, cmd *cli.Command) error {
	client := anthropic.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()

	if len(unusedArgs) > 0 {
		return fmt.Errorf("Unexpected extra arguments: %v", unusedArgs)
	}

	params := anthropic.MessageNewParams{}

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

	format := cmd.Root().String("format")
	explicitFormat := cmd.Root().IsSet("format")
	transform := cmd.Root().String("transform")
	if cmd.Bool("stream") {
		stream := client.Messages.NewStreaming(ctx, params, options...)
		maxItems := int64(-1)
		if cmd.IsSet("max-items") {
			maxItems = cmd.Value("max-items").(int64)
		}
		return ShowJSONIterator(stream, maxItems, ShowJSONOpts{
			ExplicitFormat: explicitFormat,
			Format:         format,
			RawOutput:      cmd.Root().Bool("raw-output"),
			Title:          "messages create",
			Transform:      transform,
		})
	} else {
		var res []byte
		options = append(options, option.WithResponseBodyInto(&res))
		_, err = client.Messages.New(ctx, params, options...)
		if err != nil {
			return err
		}

		obj := gjson.ParseBytes(res)
		return ShowJSON(obj, ShowJSONOpts{
			ExplicitFormat: explicitFormat,
			Format:         format,
			RawOutput:      cmd.Root().Bool("raw-output"),
			Title:          "messages create",
			Transform:      transform,
		})
	}
}

func handleMessagesCountTokens(ctx context.Context, cmd *cli.Command) error {
	client := anthropic.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()

	if len(unusedArgs) > 0 {
		return fmt.Errorf("Unexpected extra arguments: %v", unusedArgs)
	}

	params := anthropic.MessageCountTokensParams{}

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
	_, err = client.Messages.CountTokens(ctx, params, options...)
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
		Title:          "messages count-tokens",
		Transform:      transform,
	})
}
