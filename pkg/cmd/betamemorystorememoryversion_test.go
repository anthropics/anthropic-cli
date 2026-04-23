// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cmd

import (
	"testing"

	"github.com/anthropics/anthropic-cli/internal/mocktest"
)

func TestBetaMemoryStoresMemoryVersionsRetrieve(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:memory-stores:memory-versions", "retrieve",
			"--memory-store-id", "memory_store_id",
			"--memory-version-id", "memory_version_id",
			"--view", "basic",
			"--beta", "message-batches-2024-09-24",
		)
	})
}

func TestBetaMemoryStoresMemoryVersionsList(t *testing.T) {
	t.Skip("buildURL drops path-level query params (SDK-4349)")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
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
			"--session-id", "session_id",
			"--view", "basic",
			"--beta", "message-batches-2024-09-24",
		)
	})
}

func TestBetaMemoryStoresMemoryVersionsRedact(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:memory-stores:memory-versions", "redact",
			"--memory-store-id", "memory_store_id",
			"--memory-version-id", "memory_version_id",
			"--beta", "message-batches-2024-09-24",
		)
	})
}
