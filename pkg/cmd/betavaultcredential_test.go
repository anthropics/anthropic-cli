// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cmd

import (
	"testing"

	"github.com/anthropics/anthropic-cli/internal/mocktest"
)

func TestBetaVaultsCredentialsCreate(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:vaults:credentials", "create",
			"--vault-id", "vlt_011CZkZDLs7fYzm1hXNPeRjv",
			"--auth", "{token: bearer_exampletoken, mcp_server_url: https://example-server.modelcontextprotocol.io/sse, type: static_bearer}",
			"--display-name", "Example credential",
			"--metadata", "{environment: production}",
			"--beta", "message-batches-2024-09-24",
		)
	})

	t.Run("piping data", func(t *testing.T) {
		// Test piping YAML data over stdin
		pipeData := []byte("" +
			"auth:\n" +
			"  token: bearer_exampletoken\n" +
			"  mcp_server_url: https://example-server.modelcontextprotocol.io/sse\n" +
			"  type: static_bearer\n" +
			"display_name: Example credential\n" +
			"metadata:\n" +
			"  environment: production\n")
		mocktest.TestRunMockTestWithPipeAndFlags(
			t, pipeData,
			"--api-key", "string",
			"beta:vaults:credentials", "create",
			"--vault-id", "vlt_011CZkZDLs7fYzm1hXNPeRjv",
			"--beta", "message-batches-2024-09-24",
		)
	})
}

func TestBetaVaultsCredentialsRetrieve(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:vaults:credentials", "retrieve",
			"--vault-id", "vlt_011CZkZDLs7fYzm1hXNPeRjv",
			"--credential-id", "vcrd_011CZkZEMt8gZan2iYOQfSkw",
			"--beta", "message-batches-2024-09-24",
		)
	})
}

func TestBetaVaultsCredentialsUpdate(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:vaults:credentials", "update",
			"--vault-id", "vlt_011CZkZDLs7fYzm1hXNPeRjv",
			"--credential-id", "vcrd_011CZkZEMt8gZan2iYOQfSkw",
			"--auth", "{type: mcp_oauth, access_token: x, expires_at: '2019-12-27T18:11:19.117Z', refresh: {refresh_token: x, scope: scope, token_endpoint_auth: {type: client_secret_basic, client_secret: x}}}",
			"--display-name", "Example credential",
			"--metadata", "{environment: production}",
			"--beta", "message-batches-2024-09-24",
		)
	})

	t.Run("piping data", func(t *testing.T) {
		// Test piping YAML data over stdin
		pipeData := []byte("" +
			"auth:\n" +
			"  type: mcp_oauth\n" +
			"  access_token: x\n" +
			"  expires_at: '2019-12-27T18:11:19.117Z'\n" +
			"  refresh:\n" +
			"    refresh_token: x\n" +
			"    scope: scope\n" +
			"    token_endpoint_auth:\n" +
			"      type: client_secret_basic\n" +
			"      client_secret: x\n" +
			"display_name: Example credential\n" +
			"metadata:\n" +
			"  environment: production\n")
		mocktest.TestRunMockTestWithPipeAndFlags(
			t, pipeData,
			"--api-key", "string",
			"beta:vaults:credentials", "update",
			"--vault-id", "vlt_011CZkZDLs7fYzm1hXNPeRjv",
			"--credential-id", "vcrd_011CZkZEMt8gZan2iYOQfSkw",
			"--beta", "message-batches-2024-09-24",
		)
	})
}

func TestBetaVaultsCredentialsList(t *testing.T) {
	t.Skip("buildURL drops path-level query params (SDK-4349)")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:vaults:credentials", "list",
			"--max-items", "10",
			"--vault-id", "vlt_011CZkZDLs7fYzm1hXNPeRjv",
			"--include-archived=true",
			"--limit", "0",
			"--page", "page",
			"--beta", "message-batches-2024-09-24",
		)
	})
}

func TestBetaVaultsCredentialsDelete(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:vaults:credentials", "delete",
			"--vault-id", "vlt_011CZkZDLs7fYzm1hXNPeRjv",
			"--credential-id", "vcrd_011CZkZEMt8gZan2iYOQfSkw",
			"--beta", "message-batches-2024-09-24",
		)
	})
}

func TestBetaVaultsCredentialsArchive(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:vaults:credentials", "archive",
			"--vault-id", "vlt_011CZkZDLs7fYzm1hXNPeRjv",
			"--credential-id", "vcrd_011CZkZEMt8gZan2iYOQfSkw",
			"--beta", "message-batches-2024-09-24",
		)
	})
}
