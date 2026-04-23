// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cmd

import (
	"testing"

	"github.com/anthropics/anthropic-cli/internal/mocktest"
	"github.com/anthropics/anthropic-cli/internal/requestflag"
)

func TestBetaMessagesCreate(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:messages", "create",
			"--max-items", "10",
			"--max-tokens", "1024",
			"--message", "{content: [{text: x, type: text, cache_control: {type: ephemeral, ttl: 5m}, citations: [{cited_text: cited_text, document_index: 0, document_title: x, end_char_index: 0, start_char_index: 0, type: char_location}]}], role: user}",
			"--model", "claude-opus-4-6",
			"--cache-control", "{type: ephemeral, ttl: 5m}",
			"--container", "{id: id, skills: [{skill_id: pdf, type: anthropic, version: latest}]}",
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
			"--system", "[{text: Today's date is 2024-06-01., type: text, cache_control: {type: ephemeral, ttl: 5m}, citations: [{cited_text: cited_text, document_index: 0, document_title: x, end_char_index: 0, start_char_index: 0, type: char_location}]}]",
			"--temperature", "1",
			"--thinking", "{type: adaptive, display: summarized}",
			"--tool-choice", "{type: auto, disable_parallel_tool_use: true}",
			"--tool", "{input_schema: {type: object, properties: {location: bar, unit: bar}, required: [location]}, name: name, allowed_callers: [direct], cache_control: {type: ephemeral, ttl: 5m}, defer_loading: true, description: Get the current weather in a given location, eager_input_streaming: true, input_examples: [{foo: bar}], strict: true, type: custom}",
			"--top-k", "5",
			"--top-p", "0.7",
			"--beta", "message-batches-2024-09-24",
		)
	})

	t.Run("inner flags", func(t *testing.T) {
		// Check that inner flags have been set up correctly
		requestflag.CheckInnerFlags(betaMessagesCreate)

		// Alternative argument passing style using inner flags
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:messages", "create",
			"--max-items", "10",
			"--max-tokens", "1024",
			"--message.content", "[{text: x, type: text, cache_control: {type: ephemeral, ttl: 5m}, citations: [{cited_text: cited_text, document_index: 0, document_title: x, end_char_index: 0, start_char_index: 0, type: char_location}]}]",
			"--message.role", "user",
			"--model", "claude-opus-4-6",
			"--cache-control.type", "ephemeral",
			"--cache-control.ttl", "5m",
			"--container", "{id: id, skills: [{skill_id: pdf, type: anthropic, version: latest}]}",
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
			"--system", "[{text: Today's date is 2024-06-01., type: text, cache_control: {type: ephemeral, ttl: 5m}, citations: [{cited_text: cited_text, document_index: 0, document_title: x, end_char_index: 0, start_char_index: 0, type: char_location}]}]",
			"--temperature", "1",
			"--thinking", "{type: adaptive, display: summarized}",
			"--tool-choice", "{type: auto, disable_parallel_tool_use: true}",
			"--tool", "{input_schema: {type: object, properties: {location: bar, unit: bar}, required: [location]}, name: name, allowed_callers: [direct], cache_control: {type: ephemeral, ttl: 5m}, defer_loading: true, description: Get the current weather in a given location, eager_input_streaming: true, input_examples: [{foo: bar}], strict: true, type: custom}",
			"--top-k", "5",
			"--top-p", "0.7",
			"--beta", "message-batches-2024-09-24",
		)
	})

	t.Run("piping data", func(t *testing.T) {
		// Test piping YAML data over stdin
		pipeData := []byte("" +
			"max_tokens: 1024\n" +
			"messages:\n" +
			"  - content:\n" +
			"      - text: x\n" +
			"        type: text\n" +
			"        cache_control:\n" +
			"          type: ephemeral\n" +
			"          ttl: 5m\n" +
			"        citations:\n" +
			"          - cited_text: cited_text\n" +
			"            document_index: 0\n" +
			"            document_title: x\n" +
			"            end_char_index: 0\n" +
			"            start_char_index: 0\n" +
			"            type: char_location\n" +
			"    role: user\n" +
			"model: claude-opus-4-6\n" +
			"cache_control:\n" +
			"  type: ephemeral\n" +
			"  ttl: 5m\n" +
			"container:\n" +
			"  id: id\n" +
			"  skills:\n" +
			"    - skill_id: pdf\n" +
			"      type: anthropic\n" +
			"      version: latest\n" +
			"context_management:\n" +
			"  edits:\n" +
			"    - type: clear_tool_uses_20250919\n" +
			"      clear_at_least:\n" +
			"        type: input_tokens\n" +
			"        value: 0\n" +
			"      clear_tool_inputs: true\n" +
			"      exclude_tools:\n" +
			"        - string\n" +
			"      keep:\n" +
			"        type: tool_uses\n" +
			"        value: 0\n" +
			"      trigger:\n" +
			"        type: input_tokens\n" +
			"        value: 1\n" +
			"inference_geo: inference_geo\n" +
			"mcp_servers:\n" +
			"  - name: name\n" +
			"    type: url\n" +
			"    url: url\n" +
			"    authorization_token: authorization_token\n" +
			"    tool_configuration:\n" +
			"      allowed_tools:\n" +
			"        - string\n" +
			"      enabled: true\n" +
			"metadata:\n" +
			"  user_id: 13803d75-b4b5-4c3e-b2a2-6f21399b021b\n" +
			"output_config:\n" +
			"  effort: low\n" +
			"  format:\n" +
			"    schema:\n" +
			"      foo: bar\n" +
			"    type: json_schema\n" +
			"output_format:\n" +
			"  schema:\n" +
			"    foo: bar\n" +
			"  type: json_schema\n" +
			"service_tier: auto\n" +
			"speed: standard\n" +
			"stop_sequences:\n" +
			"  - string\n" +
			"stream: false\n" +
			"system:\n" +
			"  - text: Today's date is 2024-06-01.\n" +
			"    type: text\n" +
			"    cache_control:\n" +
			"      type: ephemeral\n" +
			"      ttl: 5m\n" +
			"    citations:\n" +
			"      - cited_text: cited_text\n" +
			"        document_index: 0\n" +
			"        document_title: x\n" +
			"        end_char_index: 0\n" +
			"        start_char_index: 0\n" +
			"        type: char_location\n" +
			"temperature: 1\n" +
			"thinking:\n" +
			"  type: adaptive\n" +
			"  display: summarized\n" +
			"tool_choice:\n" +
			"  type: auto\n" +
			"  disable_parallel_tool_use: true\n" +
			"tools:\n" +
			"  - input_schema:\n" +
			"      type: object\n" +
			"      properties:\n" +
			"        location: bar\n" +
			"        unit: bar\n" +
			"      required:\n" +
			"        - location\n" +
			"    name: name\n" +
			"    allowed_callers:\n" +
			"      - direct\n" +
			"    cache_control:\n" +
			"      type: ephemeral\n" +
			"      ttl: 5m\n" +
			"    defer_loading: true\n" +
			"    description: Get the current weather in a given location\n" +
			"    eager_input_streaming: true\n" +
			"    input_examples:\n" +
			"      - foo: bar\n" +
			"    strict: true\n" +
			"    type: custom\n" +
			"top_k: 5\n" +
			"top_p: 0.7\n")
		mocktest.TestRunMockTestWithPipeAndFlags(
			t, pipeData,
			"--api-key", "string",
			"beta:messages", "create",
			"--max-items", "10",
			"--beta", "message-batches-2024-09-24",
		)
	})
}

func TestBetaMessagesCountTokens(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
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
			"--thinking", "{type: adaptive, display: summarized}",
			"--tool-choice", "{type: auto, disable_parallel_tool_use: true}",
			"--tool", "{input_schema: {type: object, properties: {location: bar, unit: bar}, required: [location]}, name: name, allowed_callers: [direct], cache_control: {type: ephemeral, ttl: 5m}, defer_loading: true, description: Get the current weather in a given location, eager_input_streaming: true, input_examples: [{foo: bar}], strict: true, type: custom}",
			"--beta", "message-batches-2024-09-24",
		)
	})

	t.Run("inner flags", func(t *testing.T) {
		// Check that inner flags have been set up correctly
		requestflag.CheckInnerFlags(betaMessagesCountTokens)

		// Alternative argument passing style using inner flags
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
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
			"--thinking", "{type: adaptive, display: summarized}",
			"--tool-choice", "{type: auto, disable_parallel_tool_use: true}",
			"--tool", "{input_schema: {type: object, properties: {location: bar, unit: bar}, required: [location]}, name: name, allowed_callers: [direct], cache_control: {type: ephemeral, ttl: 5m}, defer_loading: true, description: Get the current weather in a given location, eager_input_streaming: true, input_examples: [{foo: bar}], strict: true, type: custom}",
			"--beta", "message-batches-2024-09-24",
		)
	})

	t.Run("piping data", func(t *testing.T) {
		// Test piping YAML data over stdin
		pipeData := []byte("" +
			"messages:\n" +
			"  - content:\n" +
			"      - text: x\n" +
			"        type: text\n" +
			"        cache_control:\n" +
			"          type: ephemeral\n" +
			"          ttl: 5m\n" +
			"        citations:\n" +
			"          - cited_text: cited_text\n" +
			"            document_index: 0\n" +
			"            document_title: x\n" +
			"            end_char_index: 0\n" +
			"            start_char_index: 0\n" +
			"            type: char_location\n" +
			"    role: user\n" +
			"model: claude-opus-4-6\n" +
			"cache_control:\n" +
			"  type: ephemeral\n" +
			"  ttl: 5m\n" +
			"context_management:\n" +
			"  edits:\n" +
			"    - type: clear_tool_uses_20250919\n" +
			"      clear_at_least:\n" +
			"        type: input_tokens\n" +
			"        value: 0\n" +
			"      clear_tool_inputs: true\n" +
			"      exclude_tools:\n" +
			"        - string\n" +
			"      keep:\n" +
			"        type: tool_uses\n" +
			"        value: 0\n" +
			"      trigger:\n" +
			"        type: input_tokens\n" +
			"        value: 1\n" +
			"mcp_servers:\n" +
			"  - name: name\n" +
			"    type: url\n" +
			"    url: url\n" +
			"    authorization_token: authorization_token\n" +
			"    tool_configuration:\n" +
			"      allowed_tools:\n" +
			"        - string\n" +
			"      enabled: true\n" +
			"output_config:\n" +
			"  effort: low\n" +
			"  format:\n" +
			"    schema:\n" +
			"      foo: bar\n" +
			"    type: json_schema\n" +
			"output_format:\n" +
			"  schema:\n" +
			"    foo: bar\n" +
			"  type: json_schema\n" +
			"speed: standard\n" +
			"system:\n" +
			"  - text: Today's date is 2024-06-01.\n" +
			"    type: text\n" +
			"    cache_control:\n" +
			"      type: ephemeral\n" +
			"      ttl: 5m\n" +
			"    citations:\n" +
			"      - cited_text: cited_text\n" +
			"        document_index: 0\n" +
			"        document_title: x\n" +
			"        end_char_index: 0\n" +
			"        start_char_index: 0\n" +
			"        type: char_location\n" +
			"thinking:\n" +
			"  type: adaptive\n" +
			"  display: summarized\n" +
			"tool_choice:\n" +
			"  type: auto\n" +
			"  disable_parallel_tool_use: true\n" +
			"tools:\n" +
			"  - input_schema:\n" +
			"      type: object\n" +
			"      properties:\n" +
			"        location: bar\n" +
			"        unit: bar\n" +
			"      required:\n" +
			"        - location\n" +
			"    name: name\n" +
			"    allowed_callers:\n" +
			"      - direct\n" +
			"    cache_control:\n" +
			"      type: ephemeral\n" +
			"      ttl: 5m\n" +
			"    defer_loading: true\n" +
			"    description: Get the current weather in a given location\n" +
			"    eager_input_streaming: true\n" +
			"    input_examples:\n" +
			"      - foo: bar\n" +
			"    strict: true\n" +
			"    type: custom\n")
		mocktest.TestRunMockTestWithPipeAndFlags(
			t, pipeData,
			"--api-key", "string",
			"beta:messages", "count-tokens",
			"--beta", "message-batches-2024-09-24",
		)
	})
}
