// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cmd

import (
	"testing"

	"github.com/anthropics/anthropic-cli/internal/mocktest"
	"github.com/anthropics/anthropic-cli/internal/requestflag"
)

func TestCompletionsCreate(t *testing.T) {
	mocktest.TestRunMockTestWithFlags(
		t,
		"completions", "create",
		"--max-tokens-to-sample", "256",
		"--model", "claude-opus-4-6",
		"--prompt", "\n\nHuman: Hello, world!\n\nAssistant:",
		"--metadata", "{user_id: 13803d75-b4b5-4c3e-b2a2-6f21399b021b}",
		"--stop-sequence", "string",
		"--stream=false",
		"--temperature", "1",
		"--top-k", "5",
		"--top-p", "0.7",
		"--beta", "message-batches-2024-09-24",
	)

	// Check that inner flags have been set up correctly
	requestflag.CheckInnerFlags(completionsCreate)

	// Alternative argument passing style using inner flags
	mocktest.TestRunMockTestWithFlags(
		t,
		"completions", "create",
		"--max-tokens-to-sample", "256",
		"--model", "claude-opus-4-6",
		"--prompt", "\n\nHuman: Hello, world!\n\nAssistant:",
		"--metadata.user-id", "13803d75-b4b5-4c3e-b2a2-6f21399b021b",
		"--stop-sequence", "string",
		"--stream=false",
		"--temperature", "1",
		"--top-k", "5",
		"--top-p", "0.7",
		"--beta", "message-batches-2024-09-24",
	)
}
