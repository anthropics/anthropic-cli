// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cmd

import (
	"testing"

	"github.com/anthropics/anthropic-cli/internal/mocktest"
	"github.com/anthropics/anthropic-cli/internal/requestflag"
)

func TestBetaMemoryStoresMemoriesCreate(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:memory-stores:memories", "create",
			"--memory-store-id", "memory_store_id",
			"--content", "content",
			"--path", "xx",
			"--view", "basic",
			"--beta", "message-batches-2024-09-24",
		)
	})

	t.Run("piping data", func(t *testing.T) {
		// Test piping YAML data over stdin
		pipeData := []byte("" +
			"content: content\n" +
			"path: xx\n")
		mocktest.TestRunMockTestWithPipeAndFlags(
			t, pipeData,
			"--api-key", "string",
			"beta:memory-stores:memories", "create",
			"--memory-store-id", "memory_store_id",
			"--view", "basic",
			"--beta", "message-batches-2024-09-24",
		)
	})
}

func TestBetaMemoryStoresMemoriesRetrieve(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:memory-stores:memories", "retrieve",
			"--memory-store-id", "memory_store_id",
			"--memory-id", "memory_id",
			"--view", "basic",
			"--beta", "message-batches-2024-09-24",
		)
	})
}

func TestBetaMemoryStoresMemoriesUpdate(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:memory-stores:memories", "update",
			"--memory-store-id", "memory_store_id",
			"--memory-id", "memory_id",
			"--view", "basic",
			"--content", "content",
			"--path", "xx",
			"--precondition", "{type: content_sha256, content_sha256: content_sha256}",
			"--beta", "message-batches-2024-09-24",
		)
	})

	t.Run("inner flags", func(t *testing.T) {
		// Check that inner flags have been set up correctly
		requestflag.CheckInnerFlags(betaMemoryStoresMemoriesUpdate)

		// Alternative argument passing style using inner flags
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:memory-stores:memories", "update",
			"--memory-store-id", "memory_store_id",
			"--memory-id", "memory_id",
			"--view", "basic",
			"--content", "content",
			"--path", "xx",
			"--precondition.type", "content_sha256",
			"--precondition.content-sha256", "content_sha256",
			"--beta", "message-batches-2024-09-24",
		)
	})

	t.Run("piping data", func(t *testing.T) {
		// Test piping YAML data over stdin
		pipeData := []byte("" +
			"content: content\n" +
			"path: xx\n" +
			"precondition:\n" +
			"  type: content_sha256\n" +
			"  content_sha256: content_sha256\n")
		mocktest.TestRunMockTestWithPipeAndFlags(
			t, pipeData,
			"--api-key", "string",
			"beta:memory-stores:memories", "update",
			"--memory-store-id", "memory_store_id",
			"--memory-id", "memory_id",
			"--view", "basic",
			"--beta", "message-batches-2024-09-24",
		)
	})
}

func TestBetaMemoryStoresMemoriesList(t *testing.T) {
	t.Skip("buildURL drops path-level query params (SDK-4349)")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:memory-stores:memories", "list",
			"--max-items", "10",
			"--memory-store-id", "memory_store_id",
			"--depth", "0",
			"--limit", "0",
			"--order", "asc",
			"--order-by", "order_by",
			"--page", "page",
			"--path-prefix", "path_prefix",
			"--view", "basic",
			"--beta", "message-batches-2024-09-24",
		)
	})
}

func TestBetaMemoryStoresMemoriesDelete(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:memory-stores:memories", "delete",
			"--memory-store-id", "memory_store_id",
			"--memory-id", "memory_id",
			"--expected-content-sha256", "expected_content_sha256",
			"--beta", "message-batches-2024-09-24",
		)
	})
}
