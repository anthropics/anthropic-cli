// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cmd

import (
	"testing"

	"github.com/stainless-sdks/anthropic-cli/internal/mocktest"
	"github.com/stainless-sdks/anthropic-cli/internal/requestflag"
)

func TestMessagesCreate(t *testing.T) {
	mocktest.TestRunMockTestWithFlags(
		t,
		"messages", "create",
		"--max-tokens", "1024",
		"--message", "{content: [{text: x, type: text, cache_control: {type: ephemeral, ttl: 5m}, citations: [{cited_text: cited_text, document_index: 0, document_title: x, end_char_index: 0, start_char_index: 0, type: char_location}]}], role: user}",
		"--model", "claude-sonnet-4-5-20250929",
		"--metadata", "{user_id: 13803d75-b4b5-4c3e-b2a2-6f21399b021b}",
		"--service-tier", "auto",
		"--stop-sequence", "string",
		"--stream=false",
		"--system", "[{text: Today's date is 2024-06-01., type: text, cache_control: {type: ephemeral, ttl: 5m}, citations: [{cited_text: cited_text, document_index: 0, document_title: x, end_char_index: 0, start_char_index: 0, type: char_location}]}]",
		"--temperature", "1",
		"--thinking", "{budget_tokens: 1024, type: enabled}",
		"--tool-choice", "{type: auto, disable_parallel_tool_use: true}",
		"--tool", "{name: name, cache_control: {type: ephemeral, ttl: 5m}, description: Get the current weather in a given location, type: custom}",
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
		"--model", "claude-sonnet-4-5-20250929",
		"--metadata.user_id", "13803d75-b4b5-4c3e-b2a2-6f21399b021b",
		"--service-tier", "auto",
		"--stop-sequence", "string",
		"--stream=false",
		"--system", "[{text: Today's date is 2024-06-01., type: text, cache_control: {type: ephemeral, ttl: 5m}, citations: [{cited_text: cited_text, document_index: 0, document_title: x, end_char_index: 0, start_char_index: 0, type: char_location}]}]",
		"--temperature", "1",
		"--thinking", "{budget_tokens: 1024, type: enabled}",
		"--tool-choice", "{type: auto, disable_parallel_tool_use: true}",
		"--tool", "{name: name, cache_control: {type: ephemeral, ttl: 5m}, description: Get the current weather in a given location, type: custom}",
		"--top-k", "5",
		"--top-p", "0.7",
	)
}

func TestMessagesCountTokens(t *testing.T) {
	mocktest.TestRunMockTestWithFlags(
		t,
		"messages", "count-tokens",
		"--message", "{content: [{text: What is a quaternion?, type: text, cache_control: {type: ephemeral, ttl: 5m}, citations: [{cited_text: cited_text, document_index: 0, document_title: x, end_char_index: 0, start_char_index: 0, type: char_location}]}], role: user}",
		"--model", "claude-opus-4-5-20251101",
		"--system", "[{text: Today's date is 2024-06-01., type: text, cache_control: {type: ephemeral, ttl: 5m}, citations: [{cited_text: cited_text, document_index: 0, document_title: x, end_char_index: 0, start_char_index: 0, type: char_location}]}]",
		"--thinking", "{budget_tokens: 1024, type: enabled}",
		"--tool-choice", "{type: auto, disable_parallel_tool_use: true}",
		"--tool", "{name: name, cache_control: {type: ephemeral, ttl: 5m}, description: Get the current weather in a given location, type: custom}",
	)

	// Check that inner flags have been set up correctly
	requestflag.CheckInnerFlags(messagesCountTokens)

	// Alternative argument passing style using inner flags
	mocktest.TestRunMockTestWithFlags(
		t,
		"messages", "count-tokens",
		"--message.content", "[{text: What is a quaternion?, type: text, cache_control: {type: ephemeral, ttl: 5m}, citations: [{cited_text: cited_text, document_index: 0, document_title: x, end_char_index: 0, start_char_index: 0, type: char_location}]}]",
		"--message.role", "user",
		"--model", "claude-opus-4-5-20251101",
		"--system", "[{text: Today's date is 2024-06-01., type: text, cache_control: {type: ephemeral, ttl: 5m}, citations: [{cited_text: cited_text, document_index: 0, document_title: x, end_char_index: 0, start_char_index: 0, type: char_location}]}]",
		"--thinking", "{budget_tokens: 1024, type: enabled}",
		"--tool-choice", "{type: auto, disable_parallel_tool_use: true}",
		"--tool", "{name: name, cache_control: {type: ephemeral, ttl: 5m}, description: Get the current weather in a given location, type: custom}",
	)
}
