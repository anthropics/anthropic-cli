// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cmd

import (
	"testing"

	"github.com/anthropics/anthropic-cli/internal/mocktest"
	"github.com/anthropics/anthropic-cli/internal/requestflag"
)

func TestCompletionsCreate(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"completions", "create",
			"--max-items", "10",
			"--max-tokens-to-sample", "256",
			"--model", "claude-mythos-preview",
			"--prompt", "\n\nHuman: Hello, world!\n\nAssistant:",
			"--metadata", "{user_id: 13803d75-b4b5-4c3e-b2a2-6f21399b021b}",
			"--stop-sequence", "string",
			"--stream=false",
			"--temperature", "1",
			"--top-k", "5",
			"--top-p", "0.7",
			"--beta", "message-batches-2024-09-24",
		)
	})

	t.Run("inner flags", func(t *testing.T) {
		// Check that inner flags have been set up correctly
		requestflag.CheckInnerFlags(completionsCreate)

		// Alternative argument passing style using inner flags
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"completions", "create",
			"--max-items", "10",
			"--max-tokens-to-sample", "256",
			"--model", "claude-mythos-preview",
			"--prompt", "\n\nHuman: Hello, world!\n\nAssistant:",
			"--metadata.user-id", "13803d75-b4b5-4c3e-b2a2-6f21399b021b",
			"--stop-sequence", "string",
			"--stream=false",
			"--temperature", "1",
			"--top-k", "5",
			"--top-p", "0.7",
			"--beta", "message-batches-2024-09-24",
		)
	})

	t.Run("piping data", func(t *testing.T) {
		// Test piping YAML data over stdin
		pipeData := []byte("" +
			"max_tokens_to_sample: 256\n" +
			"model: claude-mythos-preview\n" +
			"prompt: |-\n" +
			"\n" +
			"\n" +
			"  Human: Hello, world!\n" +
			"\n" +
			"  Assistant:\n" +
			"metadata:\n" +
			"  user_id: 13803d75-b4b5-4c3e-b2a2-6f21399b021b\n" +
			"stop_sequences:\n" +
			"  - string\n" +
			"stream: false\n" +
			"temperature: 1\n" +
			"top_k: 5\n" +
			"top_p: 0.7\n")
		mocktest.TestRunMockTestWithPipeAndFlags(
			t, pipeData,
			"--api-key", "string",
			"completions", "create",
			"--max-items", "10",
			"--beta", "message-batches-2024-09-24",
		)
	})
}
