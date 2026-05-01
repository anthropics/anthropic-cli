// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cmd

import (
	"testing"

	"github.com/anthropics/anthropic-cli/internal/mocktest"
)

func TestBetaMemoryStoresCreate(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:memory-stores", "create",
			"--name", "x",
			"--description", "description",
			"--metadata", "{foo: string}",
			"--beta", "message-batches-2024-09-24",
		)
	})

	t.Run("piping data", func(t *testing.T) {
		// Test piping YAML data over stdin
		pipeData := []byte("" +
			"name: x\n" +
			"description: description\n" +
			"metadata:\n" +
			"  foo: string\n")
		mocktest.TestRunMockTestWithPipeAndFlags(
			t, pipeData,
			"--api-key", "string",
			"beta:memory-stores", "create",
			"--beta", "message-batches-2024-09-24",
		)
	})
}

func TestBetaMemoryStoresRetrieve(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:memory-stores", "retrieve",
			"--memory-store-id", "memory_store_id",
			"--beta", "message-batches-2024-09-24",
		)
	})
}

func TestBetaMemoryStoresUpdate(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:memory-stores", "update",
			"--memory-store-id", "memory_store_id",
			"--description", "description",
			"--metadata", "{foo: string}",
			"--name", "x",
			"--beta", "message-batches-2024-09-24",
		)
	})

	t.Run("piping data", func(t *testing.T) {
		// Test piping YAML data over stdin
		pipeData := []byte("" +
			"description: description\n" +
			"metadata:\n" +
			"  foo: string\n" +
			"name: x\n")
		mocktest.TestRunMockTestWithPipeAndFlags(
			t, pipeData,
			"--api-key", "string",
			"beta:memory-stores", "update",
			"--memory-store-id", "memory_store_id",
			"--beta", "message-batches-2024-09-24",
		)
	})
}

func TestBetaMemoryStoresList(t *testing.T) {
	t.Skip("buildURL drops path-level query params (SDK-4349)")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:memory-stores", "list",
			"--max-items", "10",
			"--created-at-gte", "'2019-12-27T18:11:19.117Z'",
			"--created-at-lte", "'2019-12-27T18:11:19.117Z'",
			"--include-archived=true",
			"--limit", "0",
			"--page", "page",
			"--beta", "message-batches-2024-09-24",
		)
	})
}

func TestBetaMemoryStoresDelete(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:memory-stores", "delete",
			"--memory-store-id", "memory_store_id",
			"--beta", "message-batches-2024-09-24",
		)
	})
}

func TestBetaMemoryStoresArchive(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:memory-stores", "archive",
			"--memory-store-id", "memory_store_id",
			"--beta", "message-batches-2024-09-24",
		)
	})
}
