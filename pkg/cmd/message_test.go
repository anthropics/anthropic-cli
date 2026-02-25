// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cmd

import (
	"testing"

	"github.com/anthropics/anthropic-cli/internal/mocktest"
	"github.com/anthropics/anthropic-cli/internal/requestflag"
)

func TestMessagesCreate(t *testing.T) {
	mocktest.TestRunMockTestWithFlags(
		t,
		"messages", "create",
		"--max-tokens", "1024",
		"--message", "{content: [{text: x, type: text, cache_control: {type: ephemeral, ttl: 5m}, citations: [{cited_text: cited_text, document_index: 0, document_title: x, end_char_index: 0, start_char_index: 0, type: char_location}]}], role: user}",
		"--model", "claude-opus-4-6",
		"--cache-control", "{type: ephemeral, ttl: 5m}",
		"--container", "container",
		"--inference-geo", "inference_geo",
		"--metadata", "{user_id: 13803d75-b4b5-4c3e-b2a2-6f21399b021b}",
		"--output-config", "{effort: low, format: {schema: {foo: bar}, type: json_schema}}",
		"--service-tier", "auto",
		"--stop-sequence", "string",
		"--stream=false",
		"--system", "{text: Today's date is 2024-06-01., type: text, cache_control: {type: ephemeral, ttl: 5m}, citations: [{cited_text: cited_text, document_index: 0, document_title: x, end_char_index: 0, start_char_index: 0, type: char_location}]}",
		"--temperature", "1",
		"--thinking", "{budget_tokens: 1024, type: enabled}",
		"--tool-choice", "{type: auto, disable_parallel_tool_use: true}",
		"--tool", "{input_schema: {type: object, properties: {location: bar, unit: bar}, required: [location]}, name: name, allowed_callers: [direct], cache_control: {type: ephemeral, ttl: 5m}, defer_loading: true, description: Get the current weather in a given location, eager_input_streaming: true, input_examples: [{foo: bar}], strict: true, type: custom}",
		"--top-k", "5",
		"--top-p", "0.7",
	)

	// Check that inner flags have been set up correctly
	requestflag.CheckInnerFlags(messagesCreate)

	// Alternative argument passing style using inner flags
	mocktest.TestRunMockTestWithFlags(
		t,
		"messages", "create",
		"--max-tokens", "1024",
		"--message.content", "[{text: x, type: text, cache_control: {type: ephemeral, ttl: 5m}, citations: [{cited_text: cited_text, document_index: 0, document_title: x, end_char_index: 0, start_char_index: 0, type: char_location}]}]",
		"--message.role", "user",
		"--model", "claude-opus-4-6",
		"--cache-control.type", "ephemeral",
		"--cache-control.ttl", "5m",
		"--container", "container",
		"--inference-geo", "inference_geo",
		"--metadata.user-id", "13803d75-b4b5-4c3e-b2a2-6f21399b021b",
		"--output-config.effort", "low",
		"--output-config.format", "{schema: {foo: bar}, type: json_schema}",
		"--service-tier", "auto",
		"--stop-sequence", "string",
		"--stream=false",
		"--system", "{text: Today's date is 2024-06-01., type: text, cache_control: {type: ephemeral, ttl: 5m}, citations: [{cited_text: cited_text, document_index: 0, document_title: x, end_char_index: 0, start_char_index: 0, type: char_location}]}",
		"--temperature", "1",
		"--thinking", "{budget_tokens: 1024, type: enabled}",
		"--tool-choice", "{type: auto, disable_parallel_tool_use: true}",
		"--tool", "{input_schema: {type: object, properties: {location: bar, unit: bar}, required: [location]}, name: name, allowed_callers: [direct], cache_control: {type: ephemeral, ttl: 5m}, defer_loading: true, description: Get the current weather in a given location, eager_input_streaming: true, input_examples: [{foo: bar}], strict: true, type: custom}",
		"--top-k", "5",
		"--top-p", "0.7",
	)
}

func TestMessagesCountTokens(t *testing.T) {
	mocktest.TestRunMockTestWithFlags(
		t,
		"messages", "count-tokens",
		"--message", "{content: [{text: x, type: text, cache_control: {type: ephemeral, ttl: 5m}, citations: [{cited_text: cited_text, document_index: 0, document_title: x, end_char_index: 0, start_char_index: 0, type: char_location}]}], role: user}",
		"--model", "claude-opus-4-6",
		"--cache-control", "{type: ephemeral, ttl: 5m}",
		"--output-config", "{effort: low, format: {schema: {foo: bar}, type: json_schema}}",
		"--system", "[{text: Today's date is 2024-06-01., type: text, cache_control: {type: ephemeral, ttl: 5m}, citations: [{cited_text: cited_text, document_index: 0, document_title: x, end_char_index: 0, start_char_index: 0, type: char_location}]}]",
		"--thinking", "{budget_tokens: 1024, type: enabled}",
		"--tool-choice", "{type: auto, disable_parallel_tool_use: true}",
		"--tool", "{input_schema: {type: object, properties: {location: bar, unit: bar}, required: [location]}, name: name, allowed_callers: [direct], cache_control: {type: ephemeral, ttl: 5m}, defer_loading: true, description: Get the current weather in a given location, eager_input_streaming: true, input_examples: [{foo: bar}], strict: true, type: custom}",
	)

	// Check that inner flags have been set up correctly
	requestflag.CheckInnerFlags(messagesCountTokens)

	// Alternative argument passing style using inner flags
	mocktest.TestRunMockTestWithFlags(
		t,
		"messages", "count-tokens",
		"--message.content", "[{text: x, type: text, cache_control: {type: ephemeral, ttl: 5m}, citations: [{cited_text: cited_text, document_index: 0, document_title: x, end_char_index: 0, start_char_index: 0, type: char_location}]}]",
		"--message.role", "user",
		"--model", "claude-opus-4-6",
		"--cache-control.type", "ephemeral",
		"--cache-control.ttl", "5m",
		"--output-config.effort", "low",
		"--output-config.format", "{schema: {foo: bar}, type: json_schema}",
		"--system", "[{text: Today's date is 2024-06-01., type: text, cache_control: {type: ephemeral, ttl: 5m}, citations: [{cited_text: cited_text, document_index: 0, document_title: x, end_char_index: 0, start_char_index: 0, type: char_location}]}]",
		"--thinking", "{budget_tokens: 1024, type: enabled}",
		"--tool-choice", "{type: auto, disable_parallel_tool_use: true}",
		"--tool", "{input_schema: {type: object, properties: {location: bar, unit: bar}, required: [location]}, name: name, allowed_callers: [direct], cache_control: {type: ephemeral, ttl: 5m}, defer_loading: true, description: Get the current weather in a given location, eager_input_streaming: true, input_examples: [{foo: bar}], strict: true, type: custom}",
	)
}
