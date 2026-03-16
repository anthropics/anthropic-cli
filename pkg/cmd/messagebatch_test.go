// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cmd

import (
	"testing"

	"github.com/anthropics/anthropic-cli/internal/mocktest"
	"github.com/anthropics/anthropic-cli/internal/requestflag"
)

func TestMessagesBatchesCreate(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t, "messages:batches", "create",
			"--api-key", "string",
			"--request", "{custom_id: my-custom-id-1, params: {max_tokens: 1024, messages: [{content: [{text: x, type: text, cache_control: {type: ephemeral, ttl: 5m}, citations: [{cited_text: cited_text, document_index: 0, document_title: x, end_char_index: 0, start_char_index: 0, type: char_location}]}], role: user}], model: claude-opus-4-6, cache_control: {type: ephemeral, ttl: 5m}, container: container, inference_geo: inference_geo, metadata: {user_id: 13803d75-b4b5-4c3e-b2a2-6f21399b021b}, output_config: {effort: low, format: {schema: {foo: bar}, type: json_schema}}, service_tier: auto, stop_sequences: [string], stream: true, system: [{text: Today's date is 2024-06-01., type: text, cache_control: {type: ephemeral, ttl: 5m}, citations: [{cited_text: cited_text, document_index: 0, document_title: x, end_char_index: 0, start_char_index: 0, type: char_location}]}], temperature: 1, thinking: {type: adaptive, display: summarized}, tool_choice: {type: auto, disable_parallel_tool_use: true}, tools: [{input_schema: {type: object, properties: {location: bar, unit: bar}, required: [location]}, name: name, allowed_callers: [direct], cache_control: {type: ephemeral, ttl: 5m}, defer_loading: true, description: Get the current weather in a given location, eager_input_streaming: true, input_examples: [{foo: bar}], strict: true, type: custom}], top_k: 5, top_p: 0.7}}",
		)
	})

	t.Run("inner flags", func(t *testing.T) {
		// Check that inner flags have been set up correctly
		requestflag.CheckInnerFlags(messagesBatchesCreate)

		// Alternative argument passing style using inner flags
		mocktest.TestRunMockTestWithFlags(
			t, "messages:batches", "create",
			"--api-key", "string",
			"--request.custom-id", "my-custom-id-1",
			"--request.params", "{max_tokens: 1024, messages: [{content: [{text: x, type: text, cache_control: {type: ephemeral, ttl: 5m}, citations: [{cited_text: cited_text, document_index: 0, document_title: x, end_char_index: 0, start_char_index: 0, type: char_location}]}], role: user}], model: claude-opus-4-6, cache_control: {type: ephemeral, ttl: 5m}, container: container, inference_geo: inference_geo, metadata: {user_id: 13803d75-b4b5-4c3e-b2a2-6f21399b021b}, output_config: {effort: low, format: {schema: {foo: bar}, type: json_schema}}, service_tier: auto, stop_sequences: [string], stream: true, system: [{text: Today's date is 2024-06-01., type: text, cache_control: {type: ephemeral, ttl: 5m}, citations: [{cited_text: cited_text, document_index: 0, document_title: x, end_char_index: 0, start_char_index: 0, type: char_location}]}], temperature: 1, thinking: {type: adaptive, display: summarized}, tool_choice: {type: auto, disable_parallel_tool_use: true}, tools: [{input_schema: {type: object, properties: {location: bar, unit: bar}, required: [location]}, name: name, allowed_callers: [direct], cache_control: {type: ephemeral, ttl: 5m}, defer_loading: true, description: Get the current weather in a given location, eager_input_streaming: true, input_examples: [{foo: bar}], strict: true, type: custom}], top_k: 5, top_p: 0.7}",
		)
	})

	t.Run("piping data", func(t *testing.T) {
		// Test piping YAML data over stdin
		pipeData := []byte("" +
			"requests:\n" +
			"  - custom_id: my-custom-id-1\n" +
			"    params:\n" +
			"      max_tokens: 1024\n" +
			"      messages:\n" +
			"        - content:\n" +
			"            - text: x\n" +
			"              type: text\n" +
			"              cache_control:\n" +
			"                type: ephemeral\n" +
			"                ttl: 5m\n" +
			"              citations:\n" +
			"                - cited_text: cited_text\n" +
			"                  document_index: 0\n" +
			"                  document_title: x\n" +
			"                  end_char_index: 0\n" +
			"                  start_char_index: 0\n" +
			"                  type: char_location\n" +
			"          role: user\n" +
			"      model: claude-opus-4-6\n" +
			"      cache_control:\n" +
			"        type: ephemeral\n" +
			"        ttl: 5m\n" +
			"      container: container\n" +
			"      inference_geo: inference_geo\n" +
			"      metadata:\n" +
			"        user_id: 13803d75-b4b5-4c3e-b2a2-6f21399b021b\n" +
			"      output_config:\n" +
			"        effort: low\n" +
			"        format:\n" +
			"          schema:\n" +
			"            foo: bar\n" +
			"          type: json_schema\n" +
			"      service_tier: auto\n" +
			"      stop_sequences:\n" +
			"        - string\n" +
			"      stream: true\n" +
			"      system:\n" +
			"        - text: Today's date is 2024-06-01.\n" +
			"          type: text\n" +
			"          cache_control:\n" +
			"            type: ephemeral\n" +
			"            ttl: 5m\n" +
			"          citations:\n" +
			"            - cited_text: cited_text\n" +
			"              document_index: 0\n" +
			"              document_title: x\n" +
			"              end_char_index: 0\n" +
			"              start_char_index: 0\n" +
			"              type: char_location\n" +
			"      temperature: 1\n" +
			"      thinking:\n" +
			"        type: adaptive\n" +
			"        display: summarized\n" +
			"      tool_choice:\n" +
			"        type: auto\n" +
			"        disable_parallel_tool_use: true\n" +
			"      tools:\n" +
			"        - input_schema:\n" +
			"            type: object\n" +
			"            properties:\n" +
			"              location: bar\n" +
			"              unit: bar\n" +
			"            required:\n" +
			"              - location\n" +
			"          name: name\n" +
			"          allowed_callers:\n" +
			"            - direct\n" +
			"          cache_control:\n" +
			"            type: ephemeral\n" +
			"            ttl: 5m\n" +
			"          defer_loading: true\n" +
			"          description: Get the current weather in a given location\n" +
			"          eager_input_streaming: true\n" +
			"          input_examples:\n" +
			"            - foo: bar\n" +
			"          strict: true\n" +
			"          type: custom\n" +
			"      top_k: 5\n" +
			"      top_p: 0.7\n")
		mocktest.TestRunMockTestWithPipeAndFlags(
			t, pipeData, "messages:batches", "create",
			"--api-key", "string",
		)
	})
}

func TestMessagesBatchesRetrieve(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t, "messages:batches", "retrieve",
			"--api-key", "string",
			"--message-batch-id", "message_batch_id",
		)
	})
}

func TestMessagesBatchesList(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t, "messages:batches", "list",
			"--api-key", "string",
			"--max-items", "10",
			"--after-id", "after_id",
			"--before-id", "before_id",
			"--limit", "1",
		)
	})
}

func TestMessagesBatchesDelete(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t, "messages:batches", "delete",
			"--api-key", "string",
			"--message-batch-id", "message_batch_id",
		)
	})
}

func TestMessagesBatchesCancel(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t, "messages:batches", "cancel",
			"--api-key", "string",
			"--message-batch-id", "message_batch_id",
		)
	})
}

func TestMessagesBatchesResults(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t, "messages:batches", "results",
			"--api-key", "string",
			"--max-items", "10",
			"--message-batch-id", "message_batch_id",
		)
	})
}
