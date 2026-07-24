// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cmd

import (
	"testing"

	"github.com/anthropics/anthropic-cli/internal/mocktest"
)

func TestBetaDreamsCreate(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:dreams", "create",
			"--input", "{memory_store_id: x, type: memory_store}",
			"--model", "string",
			"--instructions", "x",
			"--beta", "message-batches-2024-09-24",
		)
	})

	t.Run("piping data", func(t *testing.T) {
		// Test piping YAML data over stdin
		pipeData := []byte("" +
			"inputs:\n" +
			"  - memory_store_id: x\n" +
			"    type: memory_store\n" +
			"model: string\n" +
			"instructions: x\n")
		mocktest.TestRunMockTestWithPipeAndFlags(
			t, pipeData,
			"--api-key", "string",
			"beta:dreams", "create",
			"--beta", "message-batches-2024-09-24",
		)
	})
}

func TestBetaDreamsRetrieve(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:dreams", "retrieve",
			"--dream-id", "dream_id",
			"--beta", "message-batches-2024-09-24",
		)
	})
}

func TestBetaDreamsList(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:dreams", "list",
			"--max-items", "10",
			"--created-at-gt", "'2019-12-27T18:11:19.117Z'",
			"--created-at-lt", "'2019-12-27T18:11:19.117Z'",
			"--include-archived=true",
			"--limit", "0",
			"--page", "page",
			"--status", "pending",
			"--beta", "message-batches-2024-09-24",
		)
	})
}

func TestBetaDreamsArchive(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:dreams", "archive",
			"--dream-id", "dream_id",
			"--beta", "message-batches-2024-09-24",
		)
	})
}

func TestBetaDreamsCancel(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:dreams", "cancel",
			"--dream-id", "dream_id",
			"--beta", "message-batches-2024-09-24",
		)
	})
}
