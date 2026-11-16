package cmd

import (
	"testing"

	"github.com/anthropics/anthropic-cli/v2/internal/mocktest"
)

func TestBetaMemoryStoresMemoryVersionsRetrieve(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "string")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--base-url", mocktest.ServerWithQuery(t, "beta=true"),
			"beta:memory-stores:memory-versions", "retrieve",
			"--memory-store-id", "memory_store_id",
			"--memory-version-id", "memory_version_id",
			"--view", "basic",
			"--beta", "computer-use-2024-10-22",
			"--workspace-id", "wrkspc_011CZkZaBF1tNoB5wlCeusgy",
		)
	})
}

func TestBetaMemoryStoresMemoryVersionsList(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "string")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--base-url", mocktest.ServerWithQuery(t, "beta=true"),
			"beta:memory-stores:memory-versions", "list",
			"--max-items", "10",
			"--memory-store-id", "memory_store_id",
			"--api-key-id", "api_key_id",
			"--created-at-gte", "'2019-12-27T18:11:19.117Z'",
			"--created-at-lte", "'2019-12-27T18:11:19.117Z'",
			"--limit", "0",
			"--memory-id", "memory_id",
			"--operation", "created",
			"--page", "page",
			"--service-account-id", "service_account_id",
			"--session-id", "session_id",
			"--view", "basic",
			"--beta", "computer-use-2024-10-22",
			"--workspace-id", "wrkspc_011CZkZaBF1tNoB5wlCeusgy",
		)
	})
}

func TestBetaMemoryStoresMemoryVersionsRedact(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "string")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--base-url", mocktest.ServerWithQuery(t, "beta=true"),
			"beta:memory-stores:memory-versions", "redact",
			"--memory-store-id", "memory_store_id",
			"--memory-version-id", "memory_version_id",
			"--beta", "computer-use-2024-10-22",
			"--workspace-id", "wrkspc_011CZkZaBF1tNoB5wlCeusgy",
		)
	})
}
