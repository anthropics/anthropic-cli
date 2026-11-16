package cmd

import (
	"testing"

	"github.com/anthropics/anthropic-cli/v2/internal/mocktest"
	"github.com/anthropics/anthropic-cli/v2/internal/requestflag"
)

func TestBetaMessagesCreate(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "string")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--base-url", mocktest.ServerWithQuery(t, "beta=true"),
			"beta:messages", "create",
			"--max-items", "10",
			"--max-tokens", "1024",
			"--message", "{content: [{text: x, type: text, cache_control: {type: ephemeral, ttl: 5m}, citations: [{cited_text: The grass is green. The sky is blue., document_index: 0, document_title: x, end_char_index: 0, start_char_index: 0, type: char_location}]}], role: user}",
			"--model", "claude-opus-5-5",
			"--cache-control", "{type: ephemeral, ttl: 5m}",
			"--compaction", "{type: summarize, instructions: instructions}",
			"--container", "{id: id, skills: [{skill_id: pdf, type: anthropic, version: latest}]}",
			"--context-management", "{edits: [{type: clear_tool_uses_20250919, clear_at_least: {type: input_tokens, value: 0}, clear_tool_inputs: true, exclude_tools: [string], keep: {type: tool_uses, value: 0}, trigger: {type: input_tokens, value: 1}}]}",
			"--diagnostics", "{previous_message_id: previous_message_id}",
			"--fallback-credit-token", "x",
			"--fallbacks", "default",
			"--inference-geo", "inference_geo",
			"--mcp-server", "{name: name, type: url, url: url, authorization_token: authorization_token, tool_configuration: {allowed_tools: [string], enabled: true}}",
			"--metadata", "{user_id: 13803d75-b4b5-4c3e-b2a2-6f21399b021b}",
			"--output-config", "{effort: low, format: {schema: {foo: bar}, type: json_schema}, task_budget: {total: 1024, type: tokens, remaining: 0}}",
			"--output-format", "{schema: {foo: bar}, type: json_schema}",
			"--service-tier", "auto",
			"--speed", "standard",
			"--stop-sequence", "string",
			"--stream=false",
			"--system", "[{text: Today's date is 2024-06-01., type: text, cache_control: {type: ephemeral, ttl: 5m}, citations: [{cited_text: The grass is green. The sky is blue., document_index: 0, document_title: x, end_char_index: 0, start_char_index: 0, type: char_location}]}]",
			"--thinking", "{type: adaptive, block_binding: {prefix_mismatch_behavior: error}, display: summarized}",
			"--thread", "{type: create}",
			"--tool-choice", "{type: auto, disable_parallel_tool_use: true}",
			"--tool", "{input_schema: {type: object, properties: {location: bar, unit: bar}, required: [location]}, name: name, allowed_callers: [direct], cache_control: {type: ephemeral, ttl: 5m}, defer_loading: true, description: Get the current weather in a given location, eager_input_streaming: true, input_examples: [{foo: bar}], strict: true, type: custom}",
			"--beta", "computer-use-2024-10-22",
			"--user-profile-id", "anthropic-user-profile-id",
			"--workspace-id", "wrkspc_011CZkZaBF1tNoB5wlCeusgy",
		)
	})

	t.Run("inner flags", func(t *testing.T) {
		// Check that inner flags have been set up correctly
		requestflag.CheckInnerFlags(betaMessagesCreate)

		// Alternative argument passing style using inner flags
		mocktest.TestRunMockTestWithFlags(
			t,
			"--base-url", mocktest.ServerWithQuery(t, "beta=true"),
			"beta:messages", "create",
			"--max-items", "10",
			"--max-tokens", "1024",
			"--message", "{content: [{text: x, type: text, cache_control: {type: ephemeral, ttl: 5m}, citations: [{cited_text: The grass is green. The sky is blue., document_index: 0, document_title: x, end_char_index: 0, start_char_index: 0, type: char_location}]}], role: user}",
			"--model", "claude-opus-5-5",
			"--cache-control.type", "ephemeral",
			"--cache-control.ttl", "5m",
			"--compaction.type", "summarize",
			"--compaction.instructions", "instructions",
			"--container", "{id: id, skills: [{skill_id: pdf, type: anthropic, version: latest}]}",
			"--context-management.edits", "[{type: clear_tool_uses_20250919, clear_at_least: {type: input_tokens, value: 0}, clear_tool_inputs: true, exclude_tools: [string], keep: {type: tool_uses, value: 0}, trigger: {type: input_tokens, value: 1}}]",
			"--diagnostics.previous-message-id", "previous_message_id",
			"--fallback-credit-token", "x",
			"--fallbacks", "default",
			"--inference-geo", "inference_geo",
			"--mcp-server.name", "name",
			"--mcp-server.type", "url",
			"--mcp-server.url", "url",
			"--mcp-server.authorization-token", "authorization_token",
			"--mcp-server.tool-configuration", "{allowed_tools: [string], enabled: true}",
			"--metadata.user-id", "13803d75-b4b5-4c3e-b2a2-6f21399b021b",
			"--output-config.effort", "low",
			"--output-config.format", "{schema: {foo: bar}, type: json_schema}",
			"--output-config.task-budget", "{total: 1024, type: tokens, remaining: 0}",
			"--output-format.schema", "{foo: bar}",
			"--output-format.type", "json_schema",
			"--service-tier", "auto",
			"--speed", "standard",
			"--stop-sequence", "string",
			"--stream=false",
			"--system", "[{text: Today's date is 2024-06-01., type: text, cache_control: {type: ephemeral, ttl: 5m}, citations: [{cited_text: The grass is green. The sky is blue., document_index: 0, document_title: x, end_char_index: 0, start_char_index: 0, type: char_location}]}]",
			"--thinking", "{type: adaptive, block_binding: {prefix_mismatch_behavior: error}, display: summarized}",
			"--thread", "{type: create}",
			"--tool-choice", "{type: auto, disable_parallel_tool_use: true}",
			"--tool", "{input_schema: {type: object, properties: {location: bar, unit: bar}, required: [location]}, name: name, allowed_callers: [direct], cache_control: {type: ephemeral, ttl: 5m}, defer_loading: true, description: Get the current weather in a given location, eager_input_streaming: true, input_examples: [{foo: bar}], strict: true, type: custom}",
			"--beta", "computer-use-2024-10-22",
			"--user-profile-id", "anthropic-user-profile-id",
			"--workspace-id", "wrkspc_011CZkZaBF1tNoB5wlCeusgy",
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
			"          - cited_text: The grass is green. The sky is blue.\n" +
			"            document_index: 0\n" +
			"            document_title: x\n" +
			"            end_char_index: 0\n" +
			"            start_char_index: 0\n" +
			"            type: char_location\n" +
			"    role: user\n" +
			"model: claude-opus-5-5\n" +
			"cache_control:\n" +
			"  type: ephemeral\n" +
			"  ttl: 5m\n" +
			"compaction:\n" +
			"  type: summarize\n" +
			"  instructions: instructions\n" +
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
			"diagnostics:\n" +
			"  previous_message_id: previous_message_id\n" +
			"fallback_credit_token: x\n" +
			"fallbacks: default\n" +
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
			"  task_budget:\n" +
			"    total: 1024\n" +
			"    type: tokens\n" +
			"    remaining: 0\n" +
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
			"      - cited_text: The grass is green. The sky is blue.\n" +
			"        document_index: 0\n" +
			"        document_title: x\n" +
			"        end_char_index: 0\n" +
			"        start_char_index: 0\n" +
			"        type: char_location\n" +
			"thinking:\n" +
			"  type: adaptive\n" +
			"  block_binding:\n" +
			"    prefix_mismatch_behavior: error\n" +
			"  display: summarized\n" +
			"thread:\n" +
			"  type: create\n" +
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
			"--base-url", mocktest.ServerWithQuery(t, "beta=true"),
			"beta:messages", "create",
			"--max-items", "10",
			"--beta", "computer-use-2024-10-22",
			"--user-profile-id", "anthropic-user-profile-id",
			"--workspace-id", "wrkspc_011CZkZaBF1tNoB5wlCeusgy",
		)
	})
}

func TestBetaMessagesCountTokens(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "string")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--base-url", mocktest.ServerWithQuery(t, "beta=true"),
			"beta:messages", "count-tokens",
			"--message", "{content: [{text: x, type: text, cache_control: {type: ephemeral, ttl: 5m}, citations: [{cited_text: The grass is green. The sky is blue., document_index: 0, document_title: x, end_char_index: 0, start_char_index: 0, type: char_location}]}], role: user}",
			"--model", "claude-opus-5-5",
			"--cache-control", "{type: ephemeral, ttl: 5m}",
			"--compaction", "{type: summarize, instructions: instructions}",
			"--context-management", "{edits: [{type: clear_tool_uses_20250919, clear_at_least: {type: input_tokens, value: 0}, clear_tool_inputs: true, exclude_tools: [string], keep: {type: tool_uses, value: 0}, trigger: {type: input_tokens, value: 1}}]}",
			"--mcp-server", "{name: name, type: url, url: url, authorization_token: authorization_token, tool_configuration: {allowed_tools: [string], enabled: true}}",
			"--output-config", "{effort: low, format: {schema: {foo: bar}, type: json_schema}, task_budget: {total: 1024, type: tokens, remaining: 0}}",
			"--output-format", "{schema: {foo: bar}, type: json_schema}",
			"--speed", "standard",
			"--system", "[{text: Today's date is 2024-06-01., type: text, cache_control: {type: ephemeral, ttl: 5m}, citations: [{cited_text: The grass is green. The sky is blue., document_index: 0, document_title: x, end_char_index: 0, start_char_index: 0, type: char_location}]}]",
			"--thinking", "{type: adaptive, block_binding: {prefix_mismatch_behavior: error}, display: summarized}",
			"--tool-choice", "{type: auto, disable_parallel_tool_use: true}",
			"--tool", "{input_schema: {type: object, properties: {location: bar, unit: bar}, required: [location]}, name: name, allowed_callers: [direct], cache_control: {type: ephemeral, ttl: 5m}, defer_loading: true, description: Get the current weather in a given location, eager_input_streaming: true, input_examples: [{foo: bar}], strict: true, type: custom}",
			"--beta", "computer-use-2024-10-22",
			"--user-profile-id", "anthropic-user-profile-id",
			"--workspace-id", "wrkspc_011CZkZaBF1tNoB5wlCeusgy",
		)
	})

	t.Run("inner flags", func(t *testing.T) {
		// Check that inner flags have been set up correctly
		requestflag.CheckInnerFlags(betaMessagesCountTokens)

		// Alternative argument passing style using inner flags
		mocktest.TestRunMockTestWithFlags(
			t,
			"--base-url", mocktest.ServerWithQuery(t, "beta=true"),
			"beta:messages", "count-tokens",
			"--message", "{content: [{text: x, type: text, cache_control: {type: ephemeral, ttl: 5m}, citations: [{cited_text: The grass is green. The sky is blue., document_index: 0, document_title: x, end_char_index: 0, start_char_index: 0, type: char_location}]}], role: user}",
			"--model", "claude-opus-5-5",
			"--cache-control.type", "ephemeral",
			"--cache-control.ttl", "5m",
			"--compaction.type", "summarize",
			"--compaction.instructions", "instructions",
			"--context-management.edits", "[{type: clear_tool_uses_20250919, clear_at_least: {type: input_tokens, value: 0}, clear_tool_inputs: true, exclude_tools: [string], keep: {type: tool_uses, value: 0}, trigger: {type: input_tokens, value: 1}}]",
			"--mcp-server.name", "name",
			"--mcp-server.type", "url",
			"--mcp-server.url", "url",
			"--mcp-server.authorization-token", "authorization_token",
			"--mcp-server.tool-configuration", "{allowed_tools: [string], enabled: true}",
			"--output-config.effort", "low",
			"--output-config.format", "{schema: {foo: bar}, type: json_schema}",
			"--output-config.task-budget", "{total: 1024, type: tokens, remaining: 0}",
			"--output-format.schema", "{foo: bar}",
			"--output-format.type", "json_schema",
			"--speed", "standard",
			"--system", "[{text: Today's date is 2024-06-01., type: text, cache_control: {type: ephemeral, ttl: 5m}, citations: [{cited_text: The grass is green. The sky is blue., document_index: 0, document_title: x, end_char_index: 0, start_char_index: 0, type: char_location}]}]",
			"--thinking", "{type: adaptive, block_binding: {prefix_mismatch_behavior: error}, display: summarized}",
			"--tool-choice", "{type: auto, disable_parallel_tool_use: true}",
			"--tool", "{input_schema: {type: object, properties: {location: bar, unit: bar}, required: [location]}, name: name, allowed_callers: [direct], cache_control: {type: ephemeral, ttl: 5m}, defer_loading: true, description: Get the current weather in a given location, eager_input_streaming: true, input_examples: [{foo: bar}], strict: true, type: custom}",
			"--beta", "computer-use-2024-10-22",
			"--user-profile-id", "anthropic-user-profile-id",
			"--workspace-id", "wrkspc_011CZkZaBF1tNoB5wlCeusgy",
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
			"          - cited_text: The grass is green. The sky is blue.\n" +
			"            document_index: 0\n" +
			"            document_title: x\n" +
			"            end_char_index: 0\n" +
			"            start_char_index: 0\n" +
			"            type: char_location\n" +
			"    role: user\n" +
			"model: claude-opus-5-5\n" +
			"cache_control:\n" +
			"  type: ephemeral\n" +
			"  ttl: 5m\n" +
			"compaction:\n" +
			"  type: summarize\n" +
			"  instructions: instructions\n" +
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
			"  task_budget:\n" +
			"    total: 1024\n" +
			"    type: tokens\n" +
			"    remaining: 0\n" +
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
			"      - cited_text: The grass is green. The sky is blue.\n" +
			"        document_index: 0\n" +
			"        document_title: x\n" +
			"        end_char_index: 0\n" +
			"        start_char_index: 0\n" +
			"        type: char_location\n" +
			"thinking:\n" +
			"  type: adaptive\n" +
			"  block_binding:\n" +
			"    prefix_mismatch_behavior: error\n" +
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
			"--base-url", mocktest.ServerWithQuery(t, "beta=true"),
			"beta:messages", "count-tokens",
			"--beta", "computer-use-2024-10-22",
			"--user-profile-id", "anthropic-user-profile-id",
			"--workspace-id", "wrkspc_011CZkZaBF1tNoB5wlCeusgy",
		)
	})
}
