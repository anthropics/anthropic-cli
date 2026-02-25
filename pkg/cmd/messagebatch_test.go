// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cmd

import (
	"testing"

	"github.com/anthropics/anthropic-cli/internal/mocktest"
	"github.com/anthropics/anthropic-cli/internal/requestflag"
)

func TestMessagesBatchesCreate(t *testing.T) {
	mocktest.TestRunMockTestWithFlags(
		t,
		"messages:batches", "create",
		"--request", "{custom_id: my-custom-id-1, params: {max_tokens: 1024, messages: [{content: [{text: x, type: text, cache_control: {type: ephemeral, ttl: 5m}, citations: [{cited_text: cited_text, document_index: 0, document_title: x, end_char_index: 0, start_char_index: 0, type: char_location}]}], role: user}], model: claude-opus-4-6, cache_control: {type: ephemeral, ttl: 5m}, container: container, inference_geo: inference_geo, metadata: {user_id: 13803d75-b4b5-4c3e-b2a2-6f21399b021b}, output_config: {effort: low, format: {schema: {foo: bar}, type: json_schema}}, service_tier: auto, stop_sequences: [string], stream: true, system: [{text: Today's date is 2024-06-01., type: text, cache_control: {type: ephemeral, ttl: 5m}, citations: [{cited_text: cited_text, document_index: 0, document_title: x, end_char_index: 0, start_char_index: 0, type: char_location}]}], temperature: 1, thinking: {budget_tokens: 1024, type: enabled}, tool_choice: {type: auto, disable_parallel_tool_use: true}, tools: [{input_schema: {type: object, properties: {location: bar, unit: bar}, required: [location]}, name: name, allowed_callers: [direct], cache_control: {type: ephemeral, ttl: 5m}, defer_loading: true, description: Get the current weather in a given location, eager_input_streaming: true, input_examples: [{foo: bar}], strict: true, type: custom}], top_k: 5, top_p: 0.7}}",
	)

	// Check that inner flags have been set up correctly
	requestflag.CheckInnerFlags(messagesBatchesCreate)

	// Alternative argument passing style using inner flags
	mocktest.TestRunMockTestWithFlags(
		t,
		"messages:batches", "create",
		"--request.custom-id", "my-custom-id-1",
		"--request.params", "{max_tokens: 1024, messages: [{content: [{text: x, type: text, cache_control: {type: ephemeral, ttl: 5m}, citations: [{cited_text: cited_text, document_index: 0, document_title: x, end_char_index: 0, start_char_index: 0, type: char_location}]}], role: user}], model: claude-opus-4-6, cache_control: {type: ephemeral, ttl: 5m}, container: container, inference_geo: inference_geo, metadata: {user_id: 13803d75-b4b5-4c3e-b2a2-6f21399b021b}, output_config: {effort: low, format: {schema: {foo: bar}, type: json_schema}}, service_tier: auto, stop_sequences: [string], stream: true, system: [{text: Today's date is 2024-06-01., type: text, cache_control: {type: ephemeral, ttl: 5m}, citations: [{cited_text: cited_text, document_index: 0, document_title: x, end_char_index: 0, start_char_index: 0, type: char_location}]}], temperature: 1, thinking: {budget_tokens: 1024, type: enabled}, tool_choice: {type: auto, disable_parallel_tool_use: true}, tools: [{input_schema: {type: object, properties: {location: bar, unit: bar}, required: [location]}, name: name, allowed_callers: [direct], cache_control: {type: ephemeral, ttl: 5m}, defer_loading: true, description: Get the current weather in a given location, eager_input_streaming: true, input_examples: [{foo: bar}], strict: true, type: custom}], top_k: 5, top_p: 0.7}",
	)
}

func TestMessagesBatchesRetrieve(t *testing.T) {
	mocktest.TestRunMockTestWithFlags(
		t,
		"messages:batches", "retrieve",
		"--message-batch-id", "message_batch_id",
	)
}

func TestMessagesBatchesList(t *testing.T) {
	mocktest.TestRunMockTestWithFlags(
		t,
		"messages:batches", "list",
		"--after-id", "after_id",
		"--before-id", "before_id",
		"--limit", "1",
	)
}

func TestMessagesBatchesDelete(t *testing.T) {
	mocktest.TestRunMockTestWithFlags(
		t,
		"messages:batches", "delete",
		"--message-batch-id", "message_batch_id",
	)
}

func TestMessagesBatchesCancel(t *testing.T) {
	mocktest.TestRunMockTestWithFlags(
		t,
		"messages:batches", "cancel",
		"--message-batch-id", "message_batch_id",
	)
}

func TestMessagesBatchesResults(t *testing.T) {
	t.Skip("Mock server doesn't support application/x-jsonl responses")
	mocktest.TestRunMockTestWithFlags(
		t,
		"messages:batches", "results",
		"--message-batch-id", "message_batch_id",
	)
}
