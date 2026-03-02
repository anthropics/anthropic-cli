// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cmd

import (
	"testing"

	"github.com/anthropics/anthropic-cli/internal/mocktest"
	"github.com/anthropics/anthropic-cli/internal/requestflag"
)

func TestBetaMessagesBatchesCreate(t *testing.T) {
	t.Skip("prism validates based on the non-beta endpoint")
	mocktest.TestRunMockTestWithFlags(
		t,
		"beta:messages:batches", "create",
		"--api-key", "string",
		"--request", "{custom_id: my-custom-id-1, params: {max_tokens: 1024, messages: [{content: [{text: x, type: text, cache_control: {type: ephemeral, ttl: 5m}, citations: [{cited_text: cited_text, document_index: 0, document_title: x, end_char_index: 0, start_char_index: 0, type: char_location}]}], role: user}], model: claude-opus-4-6, cache_control: {type: ephemeral, ttl: 5m}, container: {id: id, skills: [{skill_id: x, type: anthropic, version: x}]}, context_management: {edits: [{type: clear_tool_uses_20250919, clear_at_least: {type: input_tokens, value: 0}, clear_tool_inputs: true, exclude_tools: [string], keep: {type: tool_uses, value: 0}, trigger: {type: input_tokens, value: 1}}]}, inference_geo: inference_geo, mcp_servers: [{name: name, type: url, url: url, authorization_token: authorization_token, tool_configuration: {allowed_tools: [string], enabled: true}}], metadata: {user_id: 13803d75-b4b5-4c3e-b2a2-6f21399b021b}, output_config: {effort: low, format: {schema: {foo: bar}, type: json_schema}}, output_format: {schema: {foo: bar}, type: json_schema}, service_tier: auto, speed: standard, stop_sequences: [string], stream: true, system: [{text: Today's date is 2024-06-01., type: text, cache_control: {type: ephemeral, ttl: 5m}, citations: [{cited_text: cited_text, document_index: 0, document_title: x, end_char_index: 0, start_char_index: 0, type: char_location}]}], temperature: 1, thinking: {budget_tokens: 1024, type: enabled}, tool_choice: {type: auto, disable_parallel_tool_use: true}, tools: [{input_schema: {type: object, properties: {location: bar, unit: bar}, required: [location]}, name: name, allowed_callers: [direct], cache_control: {type: ephemeral, ttl: 5m}, defer_loading: true, description: Get the current weather in a given location, eager_input_streaming: true, input_examples: [{foo: bar}], strict: true, type: custom}], top_k: 5, top_p: 0.7}}",
		"--beta", "message-batches-2024-09-24",
	)

	// Check that inner flags have been set up correctly
	requestflag.CheckInnerFlags(betaMessagesBatchesCreate)

	// Alternative argument passing style using inner flags
	mocktest.TestRunMockTestWithFlags(
		t,
		"beta:messages:batches", "create",
		"--request.custom-id", "my-custom-id-1",
		"--request.params", "{max_tokens: 1024, messages: [{content: [{text: x, type: text, cache_control: {type: ephemeral, ttl: 5m}, citations: [{cited_text: cited_text, document_index: 0, document_title: x, end_char_index: 0, start_char_index: 0, type: char_location}]}], role: user}], model: claude-opus-4-6, cache_control: {type: ephemeral, ttl: 5m}, container: {id: id, skills: [{skill_id: x, type: anthropic, version: x}]}, context_management: {edits: [{type: clear_tool_uses_20250919, clear_at_least: {type: input_tokens, value: 0}, clear_tool_inputs: true, exclude_tools: [string], keep: {type: tool_uses, value: 0}, trigger: {type: input_tokens, value: 1}}]}, inference_geo: inference_geo, mcp_servers: [{name: name, type: url, url: url, authorization_token: authorization_token, tool_configuration: {allowed_tools: [string], enabled: true}}], metadata: {user_id: 13803d75-b4b5-4c3e-b2a2-6f21399b021b}, output_config: {effort: low, format: {schema: {foo: bar}, type: json_schema}}, output_format: {schema: {foo: bar}, type: json_schema}, service_tier: auto, speed: standard, stop_sequences: [string], stream: true, system: [{text: Today's date is 2024-06-01., type: text, cache_control: {type: ephemeral, ttl: 5m}, citations: [{cited_text: cited_text, document_index: 0, document_title: x, end_char_index: 0, start_char_index: 0, type: char_location}]}], temperature: 1, thinking: {budget_tokens: 1024, type: enabled}, tool_choice: {type: auto, disable_parallel_tool_use: true}, tools: [{input_schema: {type: object, properties: {location: bar, unit: bar}, required: [location]}, name: name, allowed_callers: [direct], cache_control: {type: ephemeral, ttl: 5m}, defer_loading: true, description: Get the current weather in a given location, eager_input_streaming: true, input_examples: [{foo: bar}], strict: true, type: custom}], top_k: 5, top_p: 0.7}",
		"--beta", "message-batches-2024-09-24",
	)
}

func TestBetaMessagesBatchesRetrieve(t *testing.T) {
	mocktest.TestRunMockTestWithFlags(
		t,
		"beta:messages:batches", "retrieve",
		"--api-key", "string",
		"--message-batch-id", "message_batch_id",
		"--beta", "message-batches-2024-09-24",
	)
}

func TestBetaMessagesBatchesList(t *testing.T) {
	mocktest.TestRunMockTestWithFlags(
		t,
		"beta:messages:batches", "list",
		"--api-key", "string",
		"--after-id", "after_id",
		"--before-id", "before_id",
		"--limit", "1",
		"--beta", "message-batches-2024-09-24",
	)
}

func TestBetaMessagesBatchesDelete(t *testing.T) {
	mocktest.TestRunMockTestWithFlags(
		t,
		"beta:messages:batches", "delete",
		"--api-key", "string",
		"--message-batch-id", "message_batch_id",
		"--beta", "message-batches-2024-09-24",
	)
}

func TestBetaMessagesBatchesCancel(t *testing.T) {
	mocktest.TestRunMockTestWithFlags(
		t,
		"beta:messages:batches", "cancel",
		"--api-key", "string",
		"--message-batch-id", "message_batch_id",
		"--beta", "message-batches-2024-09-24",
	)
}

func TestBetaMessagesBatchesResults(t *testing.T) {
	t.Skip("Mock server doesn't support application/x-jsonl responses")
	mocktest.TestRunMockTestWithFlags(
		t,
		"beta:messages:batches", "results",
		"--api-key", "string",
		"--message-batch-id", "message_batch_id",
		"--beta", "message-batches-2024-09-24",
	)
}
