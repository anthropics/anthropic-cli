// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cmd

import (
	"testing"

	"github.com/anthropics/anthropic-cli/internal/mocktest"
)

func TestBetaVaultsCreate(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:vaults", "create",
			"--display-name", "Example vault",
			"--metadata", "{environment: production}",
			"--beta", "message-batches-2024-09-24",
		)
	})

	t.Run("piping data", func(t *testing.T) {
		// Test piping YAML data over stdin
		pipeData := []byte("" +
			"display_name: Example vault\n" +
			"metadata:\n" +
			"  environment: production\n")
		mocktest.TestRunMockTestWithPipeAndFlags(
			t, pipeData,
			"--api-key", "string",
			"beta:vaults", "create",
			"--beta", "message-batches-2024-09-24",
		)
	})
}

func TestBetaVaultsRetrieve(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:vaults", "retrieve",
			"--vault-id", "vlt_011CZkZDLs7fYzm1hXNPeRjv",
			"--beta", "message-batches-2024-09-24",
		)
	})
}

func TestBetaVaultsUpdate(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:vaults", "update",
			"--vault-id", "vlt_011CZkZDLs7fYzm1hXNPeRjv",
			"--display-name", "Example vault",
			"--metadata", "{environment: production}",
			"--beta", "message-batches-2024-09-24",
		)
	})

	t.Run("piping data", func(t *testing.T) {
		// Test piping YAML data over stdin
		pipeData := []byte("" +
			"display_name: Example vault\n" +
			"metadata:\n" +
			"  environment: production\n")
		mocktest.TestRunMockTestWithPipeAndFlags(
			t, pipeData,
			"--api-key", "string",
			"beta:vaults", "update",
			"--vault-id", "vlt_011CZkZDLs7fYzm1hXNPeRjv",
			"--beta", "message-batches-2024-09-24",
		)
	})
}

func TestBetaVaultsList(t *testing.T) {
	t.Skip("buildURL drops path-level query params (SDK-4349)")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:vaults", "list",
			"--max-items", "10",
			"--include-archived=true",
			"--limit", "0",
			"--page", "page",
			"--beta", "message-batches-2024-09-24",
		)
	})
}

func TestBetaVaultsDelete(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:vaults", "delete",
			"--vault-id", "vlt_011CZkZDLs7fYzm1hXNPeRjv",
			"--beta", "message-batches-2024-09-24",
		)
	})
}

func TestBetaVaultsArchive(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:vaults", "archive",
			"--vault-id", "vlt_011CZkZDLs7fYzm1hXNPeRjv",
			"--beta", "message-batches-2024-09-24",
		)
	})
}
