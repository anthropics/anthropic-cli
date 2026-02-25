// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cmd

import (
	"testing"

	"github.com/anthropics/anthropic-cli/internal/mocktest"
	"github.com/anthropics/anthropic-cli/internal/requestflag"
)

func TestBetaMessagesCreate(t *testing.T) {
	t.Skip("prism validates based on the non-beta endpoint")
	mocktest.TestRunMockTestWithFlags(
		t,
		"beta:messages", "create",
		"--max-tokens", "1024",
		"--message", "{content: [{text: x, type: text, cache_control: {type: ephemeral, ttl: 5m}, citations: [{cited_text: cited_text, document_index: 0, document_title: x, end_char_index: 0, start_char_index: 0, type: char_location}]}], role: user}",
		"--model", "claude-opus-4-6",
		"--cache-control", "{type: ephemeral, ttl: 5m}",
		"--container", "{id: id, skills: [{skill_id: x, type: anthropic, version: x}]}",
		"--context-management", "{edits: [{type: clear_tool_uses_20250919, clear_at_least: {type: input_tokens, value: 0}, clear_tool_inputs: true, exclude_tools: [string], keep: {type: tool_uses, value: 0}, trigger: {type: input_tokens, value: 1}}]}",
		"--inference-geo", "inference_geo",
		"--mcp-server", "{name: name, type: url, url: url, authorization_token: authorization_token, tool_configuration: {allowed_tools: [string], enabled: true}}",
		"--metadata", "{user_id: 13803d75-b4b5-4c3e-b2a2-6f21399b021b}",
		"--output-config", "{effort: low, format: {schema: {foo: bar}, type: json_schema}}",
		"--output-format", "{schema: {foo: bar}, type: json_schema}",
		"--service-tier", "auto",
		"--speed", "standard",
		"--stop-sequence", "string",
		"--stream=false",
		"--system", "{text: Today's date is 2024-06-01., type: text, cache_control: {type: ephemeral, ttl: 5m}, citations: [{cited_text: cited_text, document_index: 0, document_title: x, end_char_index: 0, start_char_index: 0, type: char_location}]}",
		"--temperature", "1",
		"--thinking", "{budget_tokens: 1024, type: enabled}",
		"--tool-choice", "{type: auto, disable_parallel_tool_use: true}",
		"--tool", "{input_schema: {type: object, properties: {location: bar, unit: bar}, required: [location]}, name: name, allowed_callers: [direct], cache_control: {type: ephemeral, ttl: 5m}, defer_loading: true, description: Get the current weather in a given location, eager_input_streaming: true, input_examples: [{foo: bar}], strict: true, type: custom}",
		"--top-k", "5",
		"--top-p", "0.7",
		"--beta", "message-batches-2024-09-24",
	)

	// Check that inner flags have been set up correctly
	requestflag.CheckInnerFlags(betaMessagesCreate)

	// Alternative argument passing style using inner flags
	mocktest.TestRunMockTestWithFlags(
		t,
		"beta:messages", "create",
		"--max-tokens", "1024",
		"--message.content", "[{text: x, type: text, cache_control: {type: ephemeral, ttl: 5m}, citations: [{cited_text: cited_text, document_index: 0, document_title: x, end_char_index: 0, start_char_index: 0, type: char_location}]}]",
		"--message.role", "user",
		"--model", "claude-opus-4-6",
		"--cache-control.type", "ephemeral",
		"--cache-control.ttl", "5m",
		"--container", "{id: id, skills: [{skill_id: x, type: anthropic, version: x}]}",
		"--context-management.edits", "[{type: clear_tool_uses_20250919, clear_at_least: {type: input_tokens, value: 0}, clear_tool_inputs: true, exclude_tools: [string], keep: {type: tool_uses, value: 0}, trigger: {type: input_tokens, value: 1}}]",
		"--inference-geo", "inference_geo",
		"--mcp-server.name", "name",
		"--mcp-server.type", "url",
		"--mcp-server.url", "url",
		"--mcp-server.authorization-token", "authorization_token",
		"--mcp-server.tool-configuration", "{allowed_tools: [string], enabled: true}",
		"--metadata.user-id", "13803d75-b4b5-4c3e-b2a2-6f21399b021b",
		"--output-config.effort", "low",
		"--output-config.format", "{schema: {foo: bar}, type: json_schema}",
		"--output-format.schema", "{foo: bar}",
		"--output-format.type", "json_schema",
		"--service-tier", "auto",
		"--speed", "standard",
		"--stop-sequence", "string",
		"--stream=false",
		"--system", "{text: Today's date is 2024-06-01., type: text, cache_control: {type: ephemeral, ttl: 5m}, citations: [{cited_text: cited_text, document_index: 0, document_title: x, end_char_index: 0, start_char_index: 0, type: char_location}]}",
		"--temperature", "1",
		"--thinking", "{budget_tokens: 1024, type: enabled}",
		"--tool-choice", "{type: auto, disable_parallel_tool_use: true}",
		"--tool", "{input_schema: {type: object, properties: {location: bar, unit: bar}, required: [location]}, name: name, allowed_callers: [direct], cache_control: {type: ephemeral, ttl: 5m}, defer_loading: true, description: Get the current weather in a given location, eager_input_streaming: true, input_examples: [{foo: bar}], strict: true, type: custom}",
		"--top-k", "5",
		"--top-p", "0.7",
		"--beta", "message-batches-2024-09-24",
	)
}

func TestBetaMessagesCountTokens(t *testing.T) {
	t.Skip("prism validates based on the non-beta endpoint")
	mocktest.TestRunMockTestWithFlags(
		t,
		"beta:messages", "count-tokens",
		"--message", "{content: [{text: x, type: text, cache_control: {type: ephemeral, ttl: 5m}, citations: [{cited_text: cited_text, document_index: 0, document_title: x, end_char_index: 0, start_char_index: 0, type: char_location}]}], role: user}",
		"--model", "claude-opus-4-6",
		"--cache-control", "{type: ephemeral, ttl: 5m}",
		"--context-management", "{edits: [{type: clear_tool_uses_20250919, clear_at_least: {type: input_tokens, value: 0}, clear_tool_inputs: true, exclude_tools: [string], keep: {type: tool_uses, value: 0}, trigger: {type: input_tokens, value: 1}}]}",
		"--mcp-server", "{name: name, type: url, url: url, authorization_token: authorization_token, tool_configuration: {allowed_tools: [string], enabled: true}}",
		"--output-config", "{effort: low, format: {schema: {foo: bar}, type: json_schema}}",
		"--output-format", "{schema: {foo: bar}, type: json_schema}",
		"--speed", "standard",
		"--system", "[{text: Today's date is 2024-06-01., type: text, cache_control: {type: ephemeral, ttl: 5m}, citations: [{cited_text: cited_text, document_index: 0, document_title: x, end_char_index: 0, start_char_index: 0, type: char_location}]}]",
		"--thinking", "{budget_tokens: 1024, type: enabled}",
		"--tool-choice", "{type: auto, disable_parallel_tool_use: true}",
		"--tool", "{input_schema: {type: object, properties: {location: bar, unit: bar}, required: [location]}, name: name, allowed_callers: [direct], cache_control: {type: ephemeral, ttl: 5m}, defer_loading: true, description: Get the current weather in a given location, eager_input_streaming: true, input_examples: [{foo: bar}], strict: true, type: custom}",
		"--beta", "message-batches-2024-09-24",
	)

	// Check that inner flags have been set up correctly
	requestflag.CheckInnerFlags(betaMessagesCountTokens)

	// Alternative argument passing style using inner flags
	mocktest.TestRunMockTestWithFlags(
		t,
		"beta:messages", "count-tokens",
		"--message.content", "[{text: x, type: text, cache_control: {type: ephemeral, ttl: 5m}, citations: [{cited_text: cited_text, document_index: 0, document_title: x, end_char_index: 0, start_char_index: 0, type: char_location}]}]",
		"--message.role", "user",
		"--model", "claude-opus-4-6",
		"--cache-control.type", "ephemeral",
		"--cache-control.ttl", "5m",
		"--context-management.edits", "[{type: clear_tool_uses_20250919, clear_at_least: {type: input_tokens, value: 0}, clear_tool_inputs: true, exclude_tools: [string], keep: {type: tool_uses, value: 0}, trigger: {type: input_tokens, value: 1}}]",
		"--mcp-server.name", "name",
		"--mcp-server.type", "url",
		"--mcp-server.url", "url",
		"--mcp-server.authorization-token", "authorization_token",
		"--mcp-server.tool-configuration", "{allowed_tools: [string], enabled: true}",
		"--output-config.effort", "low",
		"--output-config.format", "{schema: {foo: bar}, type: json_schema}",
		"--output-format.schema", "{foo: bar}",
		"--output-format.type", "json_schema",
		"--speed", "standard",
		"--system", "[{text: Today's date is 2024-06-01., type: text, cache_control: {type: ephemeral, ttl: 5m}, citations: [{cited_text: cited_text, document_index: 0, document_title: x, end_char_index: 0, start_char_index: 0, type: char_location}]}]",
		"--thinking", "{budget_tokens: 1024, type: enabled}",
		"--tool-choice", "{type: auto, disable_parallel_tool_use: true}",
		"--tool", "{input_schema: {type: object, properties: {location: bar, unit: bar}, required: [location]}, name: name, allowed_callers: [direct], cache_control: {type: ephemeral, ttl: 5m}, defer_loading: true, description: Get the current weather in a given location, eager_input_streaming: true, input_examples: [{foo: bar}], strict: true, type: custom}",
		"--beta", "message-batches-2024-09-24",
	)
}
