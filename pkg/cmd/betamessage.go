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

var betaMessagesCreate = requestflag.WithInnerFlags(cli.Command{
	Name:    "create",
	Usage:   "Send a structured list of input messages with text and/or image content, and the\nmodel will generate the next message in the conversation.",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[int64]{
			Name:     "max-tokens",
			Usage:    "The maximum number of tokens to generate before stopping.\n\nNote that our models may stop _before_ reaching this maximum. This parameter only specifies the absolute maximum number of tokens to generate.\n\nSet to `0` to populate the [prompt cache](https://platform.claude.com/docs/en/build-with-claude/prompt-caching#pre-warming-the-cache) without generating a response.\n\nDifferent models have different maximum values for this parameter.  See [models](https://platform.claude.com/docs/en/about-claude/models/overview) for details.",
			Required: true,
			BodyPath: "max_tokens",
		},
		&requestflag.Flag[[]map[string]any]{
			Name:     "message",
			Usage:    "Input messages.\n\nOur models are trained to operate on alternating `user` and `assistant` conversational turns. When creating a new `Message`, you specify the prior conversational turns with the `messages` parameter, and the model then generates the next `Message` in the conversation. Consecutive `user` or `assistant` turns in your request will be combined into a single turn.\n\nEach input message must be an object with a `role` and `content`. You can specify a single `user`-role message, or you can include multiple `user` and `assistant` messages.\n\nIf the final message uses the `assistant` role, the response content will continue immediately from the content in that message. This can be used to constrain part of the model's response.\n\nExample with a single `user` message:\n\n```json\n[{\"role\": \"user\", \"content\": \"Hello, Claude\"}]\n```\n\nExample with multiple conversational turns:\n\n```json\n[\n  {\"role\": \"user\", \"content\": \"Hello there.\"},\n  {\"role\": \"assistant\", \"content\": \"Hi, I'm Claude. How can I help you?\"},\n  {\"role\": \"user\", \"content\": \"Can you explain LLMs in plain English?\"},\n]\n```\n\nExample with a partially-filled response from Claude:\n\n```json\n[\n  {\"role\": \"user\", \"content\": \"What's the Greek name for Sun? (A) Sol (B) Helios (C) Sun\"},\n  {\"role\": \"assistant\", \"content\": \"The best answer is (\"},\n]\n```\n\nEach input message `content` may be either a single `string` or an array of content blocks, where each block has a specific `type`. Using a `string` for `content` is shorthand for an array of one content block of type `\"text\"`. The following input messages are equivalent:\n\n```json\n{\"role\": \"user\", \"content\": \"Hello, Claude\"}\n```\n\n```json\n{\"role\": \"user\", \"content\": [{\"type\": \"text\", \"text\": \"Hello, Claude\"}]}\n```\n\nSee [input examples](https://platform.claude.com/docs/en/build-with-claude/working-with-messages).\n\nNote that if you want to include a [system prompt](https://platform.claude.com/docs/en/build-with-claude/prompt-engineering/claude-prompting-best-practices#give-claude-a-role), you can use the top-level `system` parameter — there is no `\"system\"` role for input messages in the Messages API.\n\nThere is a limit of 100,000 messages in a single request.",
			Required: true,
			BodyPath: "messages",
		},
		&requestflag.Flag[string]{
			Name:     "model",
			Usage:    "The model that will complete your prompt.\n\nSee [models](https://docs.anthropic.com/en/docs/models-overview) for additional details and options.",
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
		&requestflag.Flag[map[string]any]{
			Name:     "context-management",
			BodyPath: "context_management",
		},
		&requestflag.Flag[map[string]any]{
			Name:     "diagnostics",
			Usage:    "Request-level diagnostics. Currently carries the previous response\nid for prompt-cache divergence reporting.",
			BodyPath: "diagnostics",
		},
		&requestflag.Flag[any]{
			Name:     "fallback-credit-token",
			Usage:    "The `fallback_credit_token` from a prior refusal's `stop_details`.\n\nWhen a preceding request was refused and returned a `fallback_credit_token`,\npass that code here on the retry to have the retry's cache-creation tokens\nfor the prefix that was warm on the refused model billed at the cache-read\nrate. Must be redeemed by the same organization and workspace, with the same\nrequest body (optionally extended by one appended `assistant` message whose\ncontent is the partial text — with any trailing whitespace stripped from\nthe final text block — and paired server-tool blocks streamed before the\nrefusal; the appended-assistant form is not available for requests with\n`output_format` set or forced `tool_choice`), on an eligible fallback\nmodel, on the same platform,\nand within 5 minutes of the refusal; a mismatch is a 400. A token minted\nmid-server-tool-loop whose partial content was continuable may only be\nredeemed with the appended-assistant form — if an exact-body retry is\nrejected with a 400 saying the token must be redeemed by continuing the\npartial response, retry with the appended-assistant form instead.\n\nWhen the appended-assistant form is used on a model that otherwise disallows\nassistant-turn prefill, this token also authorizes that one prefill.",
			BodyPath: "fallback_credit_token",
		},
		&requestflag.Flag[any]{
			Name:     "fallbacks",
			Usage:    `Opt-in server-side retry on one or more substitute models when the requested model declines for policy reasons. Tried in order: if the first entry also declines, the second is tried, and so on. The string "default" requests the requested model's server-defined default fallback configuration.`,
			BodyPath: "fallbacks",
		},
		&requestflag.Flag[*string]{
			Name:     "inference-geo",
			Usage:    "Specifies the geographic region for inference processing. If not specified, the workspace's `default_inference_geo` is used.",
			BodyPath: "inference_geo",
		},
		&requestflag.Flag[[]map[string]any]{
			Name:     "mcp-server",
			Usage:    "MCP servers to be utilized in this request",
			BodyPath: "mcp_servers",
		},
		&requestflag.Flag[map[string]any]{
			Name:     "metadata",
			BodyPath: "metadata",
		},
		&requestflag.Flag[map[string]any]{
			Name:     "output-config",
			BodyPath: "output_config",
		},
		&requestflag.Flag[map[string]any]{
			Name:     "output-format",
			BodyPath: "output_format",
		},
		&requestflag.Flag[string]{
			Name:     "service-tier",
			Usage:    "Determines whether to use priority capacity (if available) or standard capacity for this request.\n\nAnthropic offers different levels of service for your API requests. See [service-tiers](https://platform.claude.com/docs/en/api/service-tiers) for details.",
			BodyPath: "service_tier",
		},
		&requestflag.Flag[*string]{
			Name:     "speed",
			Usage:    "Inference speed mode. `fast` provides significantly faster output token generation at premium pricing. Not all models support `fast`; invalid combinations are rejected at create time.",
			BodyPath: "speed",
		},
		&requestflag.Flag[[]string]{
			Name:     "stop-sequence",
			Usage:    "Custom text sequences that will cause the model to stop generating.\n\nOur models will normally stop when they have naturally completed their turn, which will result in a response `stop_reason` of `\"end_turn\"`.\n\nIf you want the model to stop generating when it encounters custom strings of text, you can use the `stop_sequences` parameter. If the model encounters one of the custom sequences, the response `stop_reason` value will be `\"stop_sequence\"` and the response `stop_sequence` value will contain the matched stop sequence.",
			BodyPath: "stop_sequences",
		},
		&requestflag.Flag[bool]{
			Name:     "stream",
			Usage:    "Whether to incrementally stream the response using server-sent events.\n\nSee [streaming](https://platform.claude.com/docs/en/build-with-claude/streaming) for details.",
			BodyPath: "stream",
		},
		&requestflag.Flag[any]{
			Name:     "system",
			Usage:    "System prompt.\n\nA system prompt is a way of providing context and instructions to Claude, such as specifying a particular goal or role. See our [guide to system prompts](https://platform.claude.com/docs/en/build-with-claude/prompt-engineering/claude-prompting-best-practices#give-claude-a-role).",
			BodyPath: "system",
		},
		&requestflag.Flag[float64]{
			Name:     "temperature",
			Usage:    "Amount of randomness injected into the response.\n\nDefaults to `1.0`. Ranges from `0.0` to `1.0`. Use `temperature` closer to `0.0` for analytical / multiple choice, and closer to `1.0` for creative and generative tasks.\n\nNote that even with `temperature` of `0.0`, the results will not be fully deterministic.",
			BodyPath: "temperature",
		},
		&requestflag.Flag[map[string]any]{
			Name:     "thinking",
			Usage:    "Configuration for enabling Claude's extended thinking.\n\nWhen enabled, responses include `thinking` content blocks showing Claude's thinking process before the final answer. Requires a minimum budget of 1,024 tokens and counts towards your `max_tokens` limit.\n\nSee [extended thinking](https://platform.claude.com/docs/en/build-with-claude/extended-thinking) for details.",
			BodyPath: "thinking",
		},
		&requestflag.Flag[map[string]any]{
			Name:     "tool-choice",
			Usage:    "How the model should use the provided tools. The model can use a specific tool, any available tool, decide by itself, or not use tools at all.",
			BodyPath: "tool_choice",
		},
		&requestflag.Flag[[]map[string]any]{
			Name:     "tool",
			Usage:    "Definitions of tools that the model may use.\n\nIf you include `tools` in your API request, the model may return `tool_use` content blocks that represent the model's use of those tools. You can then run those tools using the tool input generated by the model and then optionally return results back to the model using `tool_result` content blocks.\n\nThere are two types of tools: **client tools** and **server tools**. The behavior described below applies to client tools. For [server tools](https://platform.claude.com/docs/en/agents-and-tools/tool-use/server-tools), see their individual documentation as each has its own behavior (e.g., the [web search tool](https://platform.claude.com/docs/en/agents-and-tools/tool-use/web-search-tool)).\n\nEach tool definition includes:\n\n* `name`: Name of the tool.\n* `description`: Optional, but strongly-recommended description of the tool.\n* `input_schema`: [JSON schema](https://json-schema.org/draft/2020-12) for the tool `input` shape that the model will produce in `tool_use` output content blocks.\n\nFor example, if you defined `tools` as:\n\n```json\n[\n  {\n    \"name\": \"get_stock_price\",\n    \"description\": \"Get the current stock price for a given ticker symbol.\",\n    \"input_schema\": {\n      \"type\": \"object\",\n      \"properties\": {\n        \"ticker\": {\n          \"type\": \"string\",\n          \"description\": \"The stock ticker symbol, e.g. AAPL for Apple Inc.\"\n        }\n      },\n      \"required\": [\"ticker\"]\n    }\n  }\n]\n```\n\nAnd then asked the model \"What's the S&P 500 at today?\", the model might produce `tool_use` content blocks in the response like this:\n\n```json\n[\n  {\n    \"type\": \"tool_use\",\n    \"id\": \"toolu_01D7FLrfh4GYq7yT1ULFeyMV\",\n    \"name\": \"get_stock_price\",\n    \"input\": { \"ticker\": \"^GSPC\" }\n  }\n]\n```\n\nYou might then run your `get_stock_price` tool with `{\"ticker\": \"^GSPC\"}` as an input, and return the following back to the model in a subsequent `user` message:\n\n```json\n[\n  {\n    \"type\": \"tool_result\",\n    \"tool_use_id\": \"toolu_01D7FLrfh4GYq7yT1ULFeyMV\",\n    \"content\": \"259.75 USD\"\n  }\n]\n```\n\nTools can be used for workflows that include running client-side tools and functions, or more generally whenever you want the model to produce a particular JSON structure of output.\n\nSee our [guide](https://platform.claude.com/docs/en/agents-and-tools/tool-use/overview) for more details.",
			BodyPath: "tools",
		},
		&requestflag.Flag[int64]{
			Name:     "top-k",
			Usage:    "Only sample from the top K options for each subsequent token.\n\nUsed to remove \"long tail\" low probability responses. [Learn more technical details here](https://towardsdatascience.com/how-to-sample-from-language-models-682bceb97277).\n\nRecommended for advanced use cases only.",
			BodyPath: "top_k",
		},
		&requestflag.Flag[float64]{
			Name:     "top-p",
			Usage:    "Use nucleus sampling.\n\nIn nucleus sampling, we compute the cumulative distribution over all the options for each subsequent token in decreasing probability order and cut it off once it reaches a particular probability specified by `top_p`.\n\nRecommended for advanced use cases only.",
			BodyPath: "top_p",
		},
		&requestflag.Flag[[]string]{
			Name:       "beta",
			Usage:      "Optional header to specify the beta version(s) you want to use.",
			HeaderPath: "anthropic-beta",
		},
		&requestflag.Flag[string]{
			Name:       "user-profile-id",
			Usage:      "The user profile ID to attribute this request to. Use when acting on behalf of a party other than your organization. Requires the `user-profiles` beta header.",
			HeaderPath: "anthropic-user-profile-id",
		},
		&requestflag.Flag[int64]{
			Name:  "max-items",
			Usage: "The maximum number of items to return (use -1 for unlimited).",
		},
	},
	Action:          handleBetaMessagesCreate,
	HideHelpCommand: true,
}, map[string][]requestflag.HasOuterFlag{
	"message": {
		&requestflag.InnerFlag[[]map[string]any]{
			Name:       "message.content",
			InnerField: "content",
		},
		&requestflag.InnerFlag[string]{
			Name:       "message.role",
			Usage:      `Allowed values: "user", "assistant", "system".`,
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
			Usage:      "The time-to-live for the cache control breakpoint.\n\nThis may be one the following values:\n- `5m`: 5 minutes\n- `1h`: 1 hour\n\nDefaults to `5m`. See [prompt caching pricing](https://platform.claude.com/docs/en/build-with-claude/prompt-caching) for details.",
			InnerField: "ttl",
		},
	},
	"container": {
		&requestflag.InnerFlag[*string]{
			Name:       "container.id",
			Usage:      "Container id",
			InnerField: "id",
		},
		&requestflag.InnerFlag[any]{
			Name:       "container.skills",
			InnerField: "skills",
		},
	},
	"context-management": {
		&requestflag.InnerFlag[[]map[string]any]{
			Name:       "context-management.edits",
			Usage:      "List of context management edits to apply",
			InnerField: "edits",
		},
	},
	"diagnostics": {
		&requestflag.InnerFlag[*string]{
			Name:       "diagnostics.previous-message-id",
			Usage:      "The `id` (`msg_...`) from this client's previous /v1/messages response. The server compares that request's prompt fingerprint against this one and returns `diagnostics.cache_miss_reason` when the prompt-cache prefix could not be reused. Pass `null` on the first turn to opt in without a prior message to compare.",
			InnerField: "previous_message_id",
		},
	},
	"fallback-credit-token": {
		&requestflag.InnerFlag[string]{
			Name:       "fallback-credit-token.token",
			Usage:      "The opaque `fallback_credit_token` from a prior refusal's `stop_details` — the same string the bare-string form carries.",
			InnerField: "token",
		},
		&requestflag.InnerFlag[string]{
			Name:       "fallback-credit-token.mode",
			Usage:      "How a failing token affects the retry. `strict` (the default, and the bare-string behavior): a failing redemption is a 400 and the retry is not served. `best_effort`: the retry is served either way — a token-layer failure no longer rejects the request; the retry proceeds at normal price and the outcome is reported on the response's `usage.fallback_credit`. Two failures stay hard in both modes: a malformed token, and combining `fallback_credit_token` with `fallbacks`.",
			InnerField: "mode",
		},
	},
	"mcp-server": {
		&requestflag.InnerFlag[string]{
			Name:       "mcp-server.name",
			InnerField: "name",
		},
		&requestflag.InnerFlag[string]{
			Name:       "mcp-server.type",
			Usage:      `Allowed values: "url".`,
			InnerField: "type",
		},
		&requestflag.InnerFlag[string]{
			Name:       "mcp-server.url",
			InnerField: "url",
		},
		&requestflag.InnerFlag[*string]{
			Name:       "mcp-server.authorization-token",
			InnerField: "authorization_token",
		},
		&requestflag.InnerFlag[map[string]any]{
			Name:       "mcp-server.tool-configuration",
			InnerField: "tool_configuration",
		},
	},
	"metadata": {
		&requestflag.InnerFlag[*string]{
			Name:       "metadata.user-id",
			Usage:      "An external identifier for the user who is associated with the request.\n\nThis should be a uuid, hash value, or other opaque identifier. Anthropic may use this id to help detect abuse. Do not include any identifying information such as name, email address, or phone number.",
			InnerField: "user_id",
		},
	},
	"output-config": {
		&requestflag.InnerFlag[*string]{
			Name:       "output-config.effort",
			Usage:      "All possible effort levels.",
			InnerField: "effort",
		},
		&requestflag.InnerFlag[map[string]any]{
			Name:       "output-config.format",
			InnerField: "format",
		},
		&requestflag.InnerFlag[map[string]any]{
			Name:       "output-config.task-budget",
			Usage:      "User-configurable total token budget across contexts.",
			InnerField: "task_budget",
		},
	},
	"output-format": {
		&requestflag.InnerFlag[map[string]any]{
			Name:       "output-format.schema",
			Usage:      "The JSON schema of the format",
			InnerField: "schema",
		},
		&requestflag.InnerFlag[string]{
			Name:       "output-format.type",
			Usage:      `Allowed values: "json_schema".`,
			InnerField: "type",
		},
	},
	"thinking": {
		&requestflag.InnerFlag[string]{
			Name:       "thinking.type",
			Usage:      `Allowed values: "enabled", "disabled", "adaptive".`,
			InnerField: "type",
		},
		&requestflag.InnerFlag[int64]{
			Name:       "thinking.budget-tokens",
			Usage:      "Determines how many tokens Claude can use for its internal reasoning process. Larger budgets can enable more thorough analysis for complex problems, improving response quality.\n\nMust be ≥1024 and less than `max_tokens`.\n\nSee [extended thinking](https://platform.claude.com/docs/en/build-with-claude/extended-thinking) for details.",
			InnerField: "budget_tokens",
		},
		&requestflag.InnerFlag[*string]{
			Name:       "thinking.display",
			Usage:      "Controls how thinking content appears in the response. When set to `summarized`, thinking is returned normally. When set to `omitted`, thinking content is redacted but a signature is returned for multi-turn continuity. Defaults to `summarized`.",
			InnerField: "display",
		},
	},
	"tool-choice": {
		&requestflag.InnerFlag[string]{
			Name:       "tool-choice.type",
			Usage:      `Allowed values: "auto", "any", "tool", "none".`,
			InnerField: "type",
		},
		&requestflag.InnerFlag[bool]{
			Name:       "tool-choice.disable-parallel-tool-use",
			Usage:      "Whether to disable parallel tool use.\n\nDefaults to `false`. If set to `true`, the model will output at most one tool use.",
			InnerField: "disable_parallel_tool_use",
		},
		&requestflag.InnerFlag[string]{
			Name:       "tool-choice.name",
			Usage:      "The name of the tool to use.",
			InnerField: "name",
		},
	},
	"tool": {
		&requestflag.InnerFlag[any]{
			Name:       "tool.allowed-callers",
			InnerField: "allowed_callers",
		},
		&requestflag.InnerFlag[any]{
			Name:       "tool.allowed-domains",
			InnerField: "allowed_domains",
		},
		&requestflag.InnerFlag[any]{
			Name:       "tool.blocked-domains",
			InnerField: "blocked_domains",
		},
		&requestflag.InnerFlag[map[string]any]{
			Name:       "tool.cache-control",
			InnerField: "cache_control",
		},
		&requestflag.InnerFlag[map[string]any]{
			Name:       "tool.caching",
			InnerField: "caching",
		},
		&requestflag.InnerFlag[map[string]any]{
			Name:       "tool.citations",
			InnerField: "citations",
		},
		&requestflag.InnerFlag[any]{
			Name:       "tool.configs",
			InnerField: "configs",
		},
		&requestflag.InnerFlag[map[string]any]{
			Name:       "tool.default-config",
			Usage:      "Default configuration for tools in an MCP toolset.",
			InnerField: "default_config",
		},
		&requestflag.InnerFlag[bool]{
			Name:       "tool.defer-loading",
			Usage:      "If true, tool will not be included in initial system prompt. Only loaded when returned via tool_reference from tool search.",
			InnerField: "defer_loading",
		},
		&requestflag.InnerFlag[string]{
			Name:       "tool.description",
			Usage:      "Description of what this tool does.\n\nTool descriptions should be as detailed as possible. The more information that the model has about what the tool is and how to use it, the better it will perform. You can use natural language descriptions to reinforce important aspects of the tool input JSON schema.",
			InnerField: "description",
		},
		&requestflag.InnerFlag[int64]{
			Name:       "tool.display-height-px",
			Usage:      "The height of the display in pixels.",
			InnerField: "display_height_px",
		},
		&requestflag.InnerFlag[*int64]{
			Name:       "tool.display-number",
			Usage:      "The X11 display number (e.g. 0, 1) for the display.",
			InnerField: "display_number",
		},
		&requestflag.InnerFlag[int64]{
			Name:       "tool.display-width-px",
			Usage:      "The width of the display in pixels.",
			InnerField: "display_width_px",
		},
		&requestflag.InnerFlag[*bool]{
			Name:       "tool.eager-input-streaming",
			Usage:      "Enable eager input streaming for this tool. When true, tool input parameters will be streamed incrementally as they are generated, and types will be inferred on-the-fly rather than buffering the full JSON output. When false, streaming is disabled for this tool even if the fine-grained-tool-streaming beta is active. When null (default), uses the default behavior based on beta headers.",
			InnerField: "eager_input_streaming",
		},
		&requestflag.InnerFlag[bool]{
			Name:       "tool.enable-zoom",
			Usage:      "Whether to enable an action to take a zoomed-in screenshot of the screen.",
			InnerField: "enable_zoom",
		},
		&requestflag.InnerFlag[any]{
			Name:       "tool.input-examples",
			InnerField: "input_examples",
		},
		&requestflag.InnerFlag[any]{
			Name:       "tool.input-schema",
			InnerField: "input_schema",
		},
		&requestflag.InnerFlag[*int64]{
			Name:       "tool.max-characters",
			Usage:      "Maximum number of characters to display when viewing a file. If not specified, defaults to displaying the full file.",
			InnerField: "max_characters",
		},
		&requestflag.InnerFlag[*int64]{
			Name:       "tool.max-content-tokens",
			Usage:      "Maximum number of tokens used by including web page text content in the context. The limit is approximate and does not apply to binary content such as PDFs.",
			InnerField: "max_content_tokens",
		},
		&requestflag.InnerFlag[*int64]{
			Name:       "tool.max-tokens",
			Usage:      "Bounds the advisor's total output (thinking + text) per call. When the advisor hits this cap, the returned advisor_result or advisor_redacted_result block carries stop_reason='max_tokens', and a truncation note is appended to the advice text the worker model sees (inside the encrypted blob in redacted mode). When set, the server also emits a remaining-tokens budget block in the advisor's prompt so the advisor self-shapes toward the cap. When omitted, the advisor model's default output cap applies and no budget block is emitted.",
			InnerField: "max_tokens",
		},
		&requestflag.InnerFlag[*int64]{
			Name:       "tool.max-uses",
			Usage:      "Maximum number of times the tool can be used in the API request.",
			InnerField: "max_uses",
		},
		&requestflag.InnerFlag[string]{
			Name:       "tool.mcp-server-name",
			Usage:      "Name of the MCP server to configure tools for",
			InnerField: "mcp_server_name",
		},
		&requestflag.InnerFlag[string]{
			Name:       "tool.model",
			Usage:      "The model that will complete your prompt.\n\nSee [models](https://docs.anthropic.com/en/docs/models-overview) for additional details and options.",
			InnerField: "model",
		},
		&requestflag.InnerFlag[string]{
			Name:       "tool.name",
			Usage:      "Name of the tool.\n\nThis is how the tool will be called by the model and in `tool_use` blocks.",
			InnerField: "name",
		},
		&requestflag.InnerFlag[string]{
			Name:       "tool.response-inclusion",
			Usage:      "How this tool's result blocks appear in the API response when the result was consumed by a completed code_execution call in the same turn. 'full' returns the complete content (default). 'excluded' drops the nested server_tool_use and result block pair entirely. Results from direct calls, or from code_execution calls that paused before completing, are always returned in full so they can be sent back on the next turn.",
			InnerField: "response_inclusion",
		},
		&requestflag.InnerFlag[bool]{
			Name:       "tool.strict",
			Usage:      "When true, guarantees schema validation on tool names and inputs",
			InnerField: "strict",
		},
		&requestflag.InnerFlag[*string]{
			Name:       "tool.type",
			Usage:      `Allowed values: "custom", "bash_20241022", "bash_20250124", "code_execution_20250522", "code_execution_20250825", "code_execution_20260120", "code_execution_20260521", "browser_toolset_20260801", "computer_20241022", "memory_20250818", "computer_20250124", "text_editor_20241022", "computer_20251124", "computer_toolset_20260801", "text_editor_20250124", "text_editor_20250429", "text_editor_20250728", "web_search_20250305", "web_fetch_20250910", "web_search_20260209", "web_fetch_20260209", "web_fetch_20260309", "web_search_20260318", "web_fetch_20260318", "advisor_20260301", "tool_search_tool_bm25_20251119", "tool_search_tool_bm25", "tool_search_tool_regex_20251119", "tool_search_tool_regex", "mcp_toolset".`,
			InnerField: "type",
		},
		&requestflag.InnerFlag[bool]{
			Name:       "tool.use-cache",
			Usage:      "Whether to use cached content. Set to false to bypass the cache and fetch fresh content. Only set to false when the user explicitly requests fresh content or when fetching rapidly-changing sources.",
			InnerField: "use_cache",
		},
		&requestflag.InnerFlag[map[string]any]{
			Name:       "tool.user-location",
			InnerField: "user_location",
		},
	},
})

var betaMessagesCountTokens = requestflag.WithInnerFlags(cli.Command{
	Name:    "count-tokens",
	Usage:   "Count the number of tokens in a Message.",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[[]map[string]any]{
			Name:     "message",
			Usage:    "Input messages.\n\nOur models are trained to operate on alternating `user` and `assistant` conversational turns. When creating a new `Message`, you specify the prior conversational turns with the `messages` parameter, and the model then generates the next `Message` in the conversation. Consecutive `user` or `assistant` turns in your request will be combined into a single turn.\n\nEach input message must be an object with a `role` and `content`. You can specify a single `user`-role message, or you can include multiple `user` and `assistant` messages.\n\nIf the final message uses the `assistant` role, the response content will continue immediately from the content in that message. This can be used to constrain part of the model's response.\n\nExample with a single `user` message:\n\n```json\n[{\"role\": \"user\", \"content\": \"Hello, Claude\"}]\n```\n\nExample with multiple conversational turns:\n\n```json\n[\n  {\"role\": \"user\", \"content\": \"Hello there.\"},\n  {\"role\": \"assistant\", \"content\": \"Hi, I'm Claude. How can I help you?\"},\n  {\"role\": \"user\", \"content\": \"Can you explain LLMs in plain English?\"},\n]\n```\n\nExample with a partially-filled response from Claude:\n\n```json\n[\n  {\"role\": \"user\", \"content\": \"What's the Greek name for Sun? (A) Sol (B) Helios (C) Sun\"},\n  {\"role\": \"assistant\", \"content\": \"The best answer is (\"},\n]\n```\n\nEach input message `content` may be either a single `string` or an array of content blocks, where each block has a specific `type`. Using a `string` for `content` is shorthand for an array of one content block of type `\"text\"`. The following input messages are equivalent:\n\n```json\n{\"role\": \"user\", \"content\": \"Hello, Claude\"}\n```\n\n```json\n{\"role\": \"user\", \"content\": [{\"type\": \"text\", \"text\": \"Hello, Claude\"}]}\n```\n\nSee [input examples](https://platform.claude.com/docs/en/build-with-claude/working-with-messages).\n\nNote that if you want to include a [system prompt](https://platform.claude.com/docs/en/build-with-claude/prompt-engineering/claude-prompting-best-practices#give-claude-a-role), you can use the top-level `system` parameter — there is no `\"system\"` role for input messages in the Messages API.\n\nThere is a limit of 100,000 messages in a single request.",
			Required: true,
			BodyPath: "messages",
		},
		&requestflag.Flag[string]{
			Name:     "model",
			Usage:    "The model that will complete your prompt.\n\nSee [models](https://docs.anthropic.com/en/docs/models-overview) for additional details and options.",
			Required: true,
			BodyPath: "model",
		},
		&requestflag.Flag[map[string]any]{
			Name:     "cache-control",
			BodyPath: "cache_control",
		},
		&requestflag.Flag[map[string]any]{
			Name:     "context-management",
			BodyPath: "context_management",
		},
		&requestflag.Flag[[]map[string]any]{
			Name:     "mcp-server",
			Usage:    "MCP servers to be utilized in this request",
			BodyPath: "mcp_servers",
		},
		&requestflag.Flag[map[string]any]{
			Name:     "output-config",
			BodyPath: "output_config",
		},
		&requestflag.Flag[map[string]any]{
			Name:     "output-format",
			BodyPath: "output_format",
		},
		&requestflag.Flag[*string]{
			Name:     "speed",
			Usage:    "Inference speed mode. `fast` provides significantly faster output token generation at premium pricing. Not all models support `fast`; invalid combinations are rejected at create time.",
			BodyPath: "speed",
		},
		&requestflag.Flag[any]{
			Name:     "system",
			Usage:    "System prompt.\n\nA system prompt is a way of providing context and instructions to Claude, such as specifying a particular goal or role. See our [guide to system prompts](https://platform.claude.com/docs/en/build-with-claude/prompt-engineering/claude-prompting-best-practices#give-claude-a-role).",
			BodyPath: "system",
		},
		&requestflag.Flag[map[string]any]{
			Name:     "thinking",
			Usage:    "Configuration for enabling Claude's extended thinking.\n\nWhen enabled, responses include `thinking` content blocks showing Claude's thinking process before the final answer. Requires a minimum budget of 1,024 tokens and counts towards your `max_tokens` limit.\n\nSee [extended thinking](https://platform.claude.com/docs/en/build-with-claude/extended-thinking) for details.",
			BodyPath: "thinking",
		},
		&requestflag.Flag[map[string]any]{
			Name:     "tool-choice",
			Usage:    "How the model should use the provided tools. The model can use a specific tool, any available tool, decide by itself, or not use tools at all.",
			BodyPath: "tool_choice",
		},
		&requestflag.Flag[[]map[string]any]{
			Name:     "tool",
			Usage:    "Definitions of tools that the model may use.\n\nIf you include `tools` in your API request, the model may return `tool_use` content blocks that represent the model's use of those tools. You can then run those tools using the tool input generated by the model and then optionally return results back to the model using `tool_result` content blocks.\n\nThere are two types of tools: **client tools** and **server tools**. The behavior described below applies to client tools. For [server tools](https://platform.claude.com/docs/en/agents-and-tools/tool-use/server-tools), see their individual documentation as each has its own behavior (e.g., the [web search tool](https://platform.claude.com/docs/en/agents-and-tools/tool-use/web-search-tool)).\n\nEach tool definition includes:\n\n* `name`: Name of the tool.\n* `description`: Optional, but strongly-recommended description of the tool.\n* `input_schema`: [JSON schema](https://json-schema.org/draft/2020-12) for the tool `input` shape that the model will produce in `tool_use` output content blocks.\n\nFor example, if you defined `tools` as:\n\n```json\n[\n  {\n    \"name\": \"get_stock_price\",\n    \"description\": \"Get the current stock price for a given ticker symbol.\",\n    \"input_schema\": {\n      \"type\": \"object\",\n      \"properties\": {\n        \"ticker\": {\n          \"type\": \"string\",\n          \"description\": \"The stock ticker symbol, e.g. AAPL for Apple Inc.\"\n        }\n      },\n      \"required\": [\"ticker\"]\n    }\n  }\n]\n```\n\nAnd then asked the model \"What's the S&P 500 at today?\", the model might produce `tool_use` content blocks in the response like this:\n\n```json\n[\n  {\n    \"type\": \"tool_use\",\n    \"id\": \"toolu_01D7FLrfh4GYq7yT1ULFeyMV\",\n    \"name\": \"get_stock_price\",\n    \"input\": { \"ticker\": \"^GSPC\" }\n  }\n]\n```\n\nYou might then run your `get_stock_price` tool with `{\"ticker\": \"^GSPC\"}` as an input, and return the following back to the model in a subsequent `user` message:\n\n```json\n[\n  {\n    \"type\": \"tool_result\",\n    \"tool_use_id\": \"toolu_01D7FLrfh4GYq7yT1ULFeyMV\",\n    \"content\": \"259.75 USD\"\n  }\n]\n```\n\nTools can be used for workflows that include running client-side tools and functions, or more generally whenever you want the model to produce a particular JSON structure of output.\n\nSee our [guide](https://platform.claude.com/docs/en/agents-and-tools/tool-use/overview) for more details.",
			BodyPath: "tools",
		},
		&requestflag.Flag[[]string]{
			Name:       "beta",
			Usage:      "Optional header to specify the beta version(s) you want to use.",
			HeaderPath: "anthropic-beta",
		},
		&requestflag.Flag[string]{
			Name:       "user-profile-id",
			Usage:      "The user profile ID to attribute this request to. Use when acting on behalf of a party other than your organization. Requires the `user-profiles` beta header.",
			HeaderPath: "anthropic-user-profile-id",
		},
	},
	Action:          handleBetaMessagesCountTokens,
	HideHelpCommand: true,
}, map[string][]requestflag.HasOuterFlag{
	"message": {
		&requestflag.InnerFlag[[]map[string]any]{
			Name:       "message.content",
			InnerField: "content",
		},
		&requestflag.InnerFlag[string]{
			Name:       "message.role",
			Usage:      `Allowed values: "user", "assistant", "system".`,
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
			Usage:      "The time-to-live for the cache control breakpoint.\n\nThis may be one the following values:\n- `5m`: 5 minutes\n- `1h`: 1 hour\n\nDefaults to `5m`. See [prompt caching pricing](https://platform.claude.com/docs/en/build-with-claude/prompt-caching) for details.",
			InnerField: "ttl",
		},
	},
	"context-management": {
		&requestflag.InnerFlag[[]map[string]any]{
			Name:       "context-management.edits",
			Usage:      "List of context management edits to apply",
			InnerField: "edits",
		},
	},
	"mcp-server": {
		&requestflag.InnerFlag[string]{
			Name:       "mcp-server.name",
			InnerField: "name",
		},
		&requestflag.InnerFlag[string]{
			Name:       "mcp-server.type",
			Usage:      `Allowed values: "url".`,
			InnerField: "type",
		},
		&requestflag.InnerFlag[string]{
			Name:       "mcp-server.url",
			InnerField: "url",
		},
		&requestflag.InnerFlag[*string]{
			Name:       "mcp-server.authorization-token",
			InnerField: "authorization_token",
		},
		&requestflag.InnerFlag[map[string]any]{
			Name:       "mcp-server.tool-configuration",
			InnerField: "tool_configuration",
		},
	},
	"output-config": {
		&requestflag.InnerFlag[*string]{
			Name:       "output-config.effort",
			Usage:      "All possible effort levels.",
			InnerField: "effort",
		},
		&requestflag.InnerFlag[map[string]any]{
			Name:       "output-config.format",
			InnerField: "format",
		},
		&requestflag.InnerFlag[map[string]any]{
			Name:       "output-config.task-budget",
			Usage:      "User-configurable total token budget across contexts.",
			InnerField: "task_budget",
		},
	},
	"output-format": {
		&requestflag.InnerFlag[map[string]any]{
			Name:       "output-format.schema",
			Usage:      "The JSON schema of the format",
			InnerField: "schema",
		},
		&requestflag.InnerFlag[string]{
			Name:       "output-format.type",
			Usage:      `Allowed values: "json_schema".`,
			InnerField: "type",
		},
	},
	"thinking": {
		&requestflag.InnerFlag[string]{
			Name:       "thinking.type",
			Usage:      `Allowed values: "enabled", "disabled", "adaptive".`,
			InnerField: "type",
		},
		&requestflag.InnerFlag[int64]{
			Name:       "thinking.budget-tokens",
			Usage:      "Determines how many tokens Claude can use for its internal reasoning process. Larger budgets can enable more thorough analysis for complex problems, improving response quality.\n\nMust be ≥1024 and less than `max_tokens`.\n\nSee [extended thinking](https://platform.claude.com/docs/en/build-with-claude/extended-thinking) for details.",
			InnerField: "budget_tokens",
		},
		&requestflag.InnerFlag[*string]{
			Name:       "thinking.display",
			Usage:      "Controls how thinking content appears in the response. When set to `summarized`, thinking is returned normally. When set to `omitted`, thinking content is redacted but a signature is returned for multi-turn continuity. Defaults to `summarized`.",
			InnerField: "display",
		},
	},
	"tool-choice": {
		&requestflag.InnerFlag[string]{
			Name:       "tool-choice.type",
			Usage:      `Allowed values: "auto", "any", "tool", "none".`,
			InnerField: "type",
		},
		&requestflag.InnerFlag[bool]{
			Name:       "tool-choice.disable-parallel-tool-use",
			Usage:      "Whether to disable parallel tool use.\n\nDefaults to `false`. If set to `true`, the model will output at most one tool use.",
			InnerField: "disable_parallel_tool_use",
		},
		&requestflag.InnerFlag[string]{
			Name:       "tool-choice.name",
			Usage:      "The name of the tool to use.",
			InnerField: "name",
		},
	},
	"tool": {
		&requestflag.InnerFlag[any]{
			Name:       "tool.allowed-callers",
			InnerField: "allowed_callers",
		},
		&requestflag.InnerFlag[any]{
			Name:       "tool.allowed-domains",
			InnerField: "allowed_domains",
		},
		&requestflag.InnerFlag[any]{
			Name:       "tool.blocked-domains",
			InnerField: "blocked_domains",
		},
		&requestflag.InnerFlag[map[string]any]{
			Name:       "tool.cache-control",
			InnerField: "cache_control",
		},
		&requestflag.InnerFlag[map[string]any]{
			Name:       "tool.caching",
			InnerField: "caching",
		},
		&requestflag.InnerFlag[map[string]any]{
			Name:       "tool.citations",
			InnerField: "citations",
		},
		&requestflag.InnerFlag[any]{
			Name:       "tool.configs",
			InnerField: "configs",
		},
		&requestflag.InnerFlag[map[string]any]{
			Name:       "tool.default-config",
			Usage:      "Default configuration for tools in an MCP toolset.",
			InnerField: "default_config",
		},
		&requestflag.InnerFlag[bool]{
			Name:       "tool.defer-loading",
			Usage:      "If true, tool will not be included in initial system prompt. Only loaded when returned via tool_reference from tool search.",
			InnerField: "defer_loading",
		},
		&requestflag.InnerFlag[string]{
			Name:       "tool.description",
			Usage:      "Description of what this tool does.\n\nTool descriptions should be as detailed as possible. The more information that the model has about what the tool is and how to use it, the better it will perform. You can use natural language descriptions to reinforce important aspects of the tool input JSON schema.",
			InnerField: "description",
		},
		&requestflag.InnerFlag[int64]{
			Name:       "tool.display-height-px",
			Usage:      "The height of the display in pixels.",
			InnerField: "display_height_px",
		},
		&requestflag.InnerFlag[*int64]{
			Name:       "tool.display-number",
			Usage:      "The X11 display number (e.g. 0, 1) for the display.",
			InnerField: "display_number",
		},
		&requestflag.InnerFlag[int64]{
			Name:       "tool.display-width-px",
			Usage:      "The width of the display in pixels.",
			InnerField: "display_width_px",
		},
		&requestflag.InnerFlag[*bool]{
			Name:       "tool.eager-input-streaming",
			Usage:      "Enable eager input streaming for this tool. When true, tool input parameters will be streamed incrementally as they are generated, and types will be inferred on-the-fly rather than buffering the full JSON output. When false, streaming is disabled for this tool even if the fine-grained-tool-streaming beta is active. When null (default), uses the default behavior based on beta headers.",
			InnerField: "eager_input_streaming",
		},
		&requestflag.InnerFlag[bool]{
			Name:       "tool.enable-zoom",
			Usage:      "Whether to enable an action to take a zoomed-in screenshot of the screen.",
			InnerField: "enable_zoom",
		},
		&requestflag.InnerFlag[any]{
			Name:       "tool.input-examples",
			InnerField: "input_examples",
		},
		&requestflag.InnerFlag[any]{
			Name:       "tool.input-schema",
			InnerField: "input_schema",
		},
		&requestflag.InnerFlag[*int64]{
			Name:       "tool.max-characters",
			Usage:      "Maximum number of characters to display when viewing a file. If not specified, defaults to displaying the full file.",
			InnerField: "max_characters",
		},
		&requestflag.InnerFlag[*int64]{
			Name:       "tool.max-content-tokens",
			Usage:      "Maximum number of tokens used by including web page text content in the context. The limit is approximate and does not apply to binary content such as PDFs.",
			InnerField: "max_content_tokens",
		},
		&requestflag.InnerFlag[*int64]{
			Name:       "tool.max-tokens",
			Usage:      "Bounds the advisor's total output (thinking + text) per call. When the advisor hits this cap, the returned advisor_result or advisor_redacted_result block carries stop_reason='max_tokens', and a truncation note is appended to the advice text the worker model sees (inside the encrypted blob in redacted mode). When set, the server also emits a remaining-tokens budget block in the advisor's prompt so the advisor self-shapes toward the cap. When omitted, the advisor model's default output cap applies and no budget block is emitted.",
			InnerField: "max_tokens",
		},
		&requestflag.InnerFlag[*int64]{
			Name:       "tool.max-uses",
			Usage:      "Maximum number of times the tool can be used in the API request.",
			InnerField: "max_uses",
		},
		&requestflag.InnerFlag[string]{
			Name:       "tool.mcp-server-name",
			Usage:      "Name of the MCP server to configure tools for",
			InnerField: "mcp_server_name",
		},
		&requestflag.InnerFlag[string]{
			Name:       "tool.model",
			Usage:      "The model that will complete your prompt.\n\nSee [models](https://docs.anthropic.com/en/docs/models-overview) for additional details and options.",
			InnerField: "model",
		},
		&requestflag.InnerFlag[string]{
			Name:       "tool.name",
			Usage:      "Name of the tool.\n\nThis is how the tool will be called by the model and in `tool_use` blocks.",
			InnerField: "name",
		},
		&requestflag.InnerFlag[string]{
			Name:       "tool.response-inclusion",
			Usage:      "How this tool's result blocks appear in the API response when the result was consumed by a completed code_execution call in the same turn. 'full' returns the complete content (default). 'excluded' drops the nested server_tool_use and result block pair entirely. Results from direct calls, or from code_execution calls that paused before completing, are always returned in full so they can be sent back on the next turn.",
			InnerField: "response_inclusion",
		},
		&requestflag.InnerFlag[bool]{
			Name:       "tool.strict",
			Usage:      "When true, guarantees schema validation on tool names and inputs",
			InnerField: "strict",
		},
		&requestflag.InnerFlag[*string]{
			Name:       "tool.type",
			Usage:      `Allowed values: "custom", "bash_20241022", "bash_20250124", "code_execution_20250522", "code_execution_20250825", "code_execution_20260120", "code_execution_20260521", "browser_toolset_20260801", "computer_20241022", "memory_20250818", "computer_20250124", "text_editor_20241022", "computer_20251124", "computer_toolset_20260801", "text_editor_20250124", "text_editor_20250429", "text_editor_20250728", "web_search_20250305", "web_fetch_20250910", "web_search_20260209", "web_fetch_20260209", "web_fetch_20260309", "web_search_20260318", "web_fetch_20260318", "advisor_20260301", "tool_search_tool_bm25_20251119", "tool_search_tool_bm25", "tool_search_tool_regex_20251119", "tool_search_tool_regex", "mcp_toolset".`,
			InnerField: "type",
		},
		&requestflag.InnerFlag[bool]{
			Name:       "tool.use-cache",
			Usage:      "Whether to use cached content. Set to false to bypass the cache and fetch fresh content. Only set to false when the user explicitly requests fresh content or when fetching rapidly-changing sources.",
			InnerField: "use_cache",
		},
		&requestflag.InnerFlag[map[string]any]{
			Name:       "tool.user-location",
			InnerField: "user_location",
		},
	},
})

func handleBetaMessagesCreate(ctx context.Context, cmd *cli.Command) error {
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

	params := anthropic.BetaMessageNewParams{}

	format := cmd.Root().String("format")
	explicitFormat := cmd.Root().IsSet("format")
	transform := cmd.Root().String("transform")
	if cmd.Bool("stream") {
		stream := client.Beta.Messages.NewStreaming(ctx, params, options...)
		maxItems := int64(-1)
		if cmd.IsSet("max-items") {
			maxItems = cmd.Value("max-items").(int64)
		}
		return ShowJSONIterator(stream, maxItems, ShowJSONOpts{
			ExplicitFormat: explicitFormat,
			Format:         format,
			RawOutput:      cmd.Root().Bool("raw-output"),
			Title:          "beta:messages create",
			Transform:      transform,
		})
	} else {
		var res []byte
		options = append(options, option.WithResponseBodyInto(&res))
		_, err = client.Beta.Messages.New(ctx, params, options...)
		if err != nil {
			return err
		}

		obj := gjson.ParseBytes(res)
		return ShowJSON(obj, ShowJSONOpts{
			ExplicitFormat: explicitFormat,
			Format:         format,
			RawOutput:      cmd.Root().Bool("raw-output"),
			Title:          "beta:messages create",
			Transform:      transform,
		})
	}
}

func handleBetaMessagesCountTokens(ctx context.Context, cmd *cli.Command) error {
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

	params := anthropic.BetaMessageCountTokensParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Beta.Messages.CountTokens(ctx, params, options...)
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
		Title:          "beta:messages count-tokens",
		Transform:      transform,
	})
}
