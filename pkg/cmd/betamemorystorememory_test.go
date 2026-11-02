package cmd

import (
	"testing"

	"github.com/anthropics/anthropic-cli/v2/internal/mocktest"
	"github.com/anthropics/anthropic-cli/v2/internal/requestflag"
)

func TestBetaMemoryStoresMemoriesCreate(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "string")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--base-url", mocktest.ServerWithQuery(t, "beta=true"),
			"beta:memory-stores:memories", "create",
			"--memory-store-id", "memory_store_id",
			"--content", "content",
			"--path", "/P",
			"--view", "basic",
			"--beta", "message-batches-2024-09-24",
			"--workspace-id", "wrkspc_011CZkZaBF1tNoB5wlCeusgy",
		)
	})

	t.Run("piping data", func(t *testing.T) {
		// Test piping YAML data over stdin
		pipeData := []byte("" +
			"content: content\n" +
			"path: /P\n")
		mocktest.TestRunMockTestWithPipeAndFlags(
			t, pipeData,
			"--base-url", mocktest.ServerWithQuery(t, "beta=true"),
			"beta:memory-stores:memories", "create",
			"--memory-store-id", "memory_store_id",
			"--view", "basic",
			"--beta", "message-batches-2024-09-24",
			"--workspace-id", "wrkspc_011CZkZaBF1tNoB5wlCeusgy",
		)
	})
}

func TestBetaMemoryStoresMemoriesRetrieve(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "string")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--base-url", mocktest.ServerWithQuery(t, "beta=true"),
			"beta:memory-stores:memories", "retrieve",
			"--memory-store-id", "memory_store_id",
			"--memory-id", "memory_id",
			"--view", "basic",
			"--beta", "message-batches-2024-09-24",
			"--workspace-id", "wrkspc_011CZkZaBF1tNoB5wlCeusgy",
		)
	})
}

func TestBetaMemoryStoresMemoriesUpdate(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "string")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--base-url", mocktest.ServerWithQuery(t, "beta=true"),
			"beta:memory-stores:memories", "update",
			"--memory-store-id", "memory_store_id",
			"--memory-id", "memory_id",
			"--view", "basic",
			"--content", "content",
			"--path", "/P",
			"--precondition", "{type: content_sha256, content_sha256: 210b9798eb53baa4e69d31c1071cf03d212b8ad0ca30cf321e0ea82e120aac26}",
			"--beta", "message-batches-2024-09-24",
			"--workspace-id", "wrkspc_011CZkZaBF1tNoB5wlCeusgy",
		)
	})

	t.Run("inner flags", func(t *testing.T) {
		// Check that inner flags have been set up correctly
		requestflag.CheckInnerFlags(betaMemoryStoresMemoriesUpdate)

		// Alternative argument passing style using inner flags
		mocktest.TestRunMockTestWithFlags(
			t,
			"--base-url", mocktest.ServerWithQuery(t, "beta=true"),
			"beta:memory-stores:memories", "update",
			"--memory-store-id", "memory_store_id",
			"--memory-id", "memory_id",
			"--view", "basic",
			"--content", "content",
			"--path", "/P",
			"--precondition.type", "content_sha256",
			"--precondition.content-sha256", "210b9798eb53baa4e69d31c1071cf03d212b8ad0ca30cf321e0ea82e120aac26",
			"--beta", "message-batches-2024-09-24",
			"--workspace-id", "wrkspc_011CZkZaBF1tNoB5wlCeusgy",
		)
	})

	t.Run("piping data", func(t *testing.T) {
		// Test piping YAML data over stdin
		pipeData := []byte("" +
			"content: content\n" +
			"path: /P\n" +
			"precondition:\n" +
			"  type: content_sha256\n" +
			"  content_sha256: 210b9798eb53baa4e69d31c1071cf03d212b8ad0ca30cf321e0ea82e120aac26\n")
		mocktest.TestRunMockTestWithPipeAndFlags(
			t, pipeData,
			"--base-url", mocktest.ServerWithQuery(t, "beta=true"),
			"beta:memory-stores:memories", "update",
			"--memory-store-id", "memory_store_id",
			"--memory-id", "memory_id",
			"--view", "basic",
			"--beta", "message-batches-2024-09-24",
			"--workspace-id", "wrkspc_011CZkZaBF1tNoB5wlCeusgy",
		)
	})
}

func TestBetaMemoryStoresMemoriesList(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "string")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--base-url", mocktest.ServerWithQuery(t, "beta=true"),
			"beta:memory-stores:memories", "list",
			"--max-items", "10",
			"--memory-store-id", "memory_store_id",
			"--depth", "0",
			"--limit", "0",
			"--page", "page",
			"--path-prefix", "path_prefix",
			"--view", "basic",
			"--beta", "message-batches-2024-09-24",
			"--workspace-id", "wrkspc_011CZkZaBF1tNoB5wlCeusgy",
		)
	})
}

func TestBetaMemoryStoresMemoriesDelete(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "string")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--base-url", mocktest.ServerWithQuery(t, "beta=true"),
			"beta:memory-stores:memories", "delete",
			"--memory-store-id", "memory_store_id",
			"--memory-id", "memory_id",
			"--expected-content-sha256", "expected_content_sha256",
			"--beta", "message-batches-2024-09-24",
			"--workspace-id", "wrkspc_011CZkZaBF1tNoB5wlCeusgy",
		)
	})
}
