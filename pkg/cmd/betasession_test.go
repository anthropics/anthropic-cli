// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cmd

import (
	"testing"

	"github.com/anthropics/anthropic-cli/internal/mocktest"
)

func TestBetaSessionsCreate(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:sessions", "create",
			"--agent", "agent_011CZkYpogX7uDKUyvBTophP",
			"--environment-id", "env_011CZkZ9X2dpNyB7HsEFoRfW",
			"--metadata", "{foo: string}",
			"--resource", "{file_id: file_011CNha8iCJcU1wXNR6q4V8w, type: file, mount_path: /uploads/receipt.pdf}",
			"--title", "Order #1234 inquiry",
			"--vault-id", "string",
			"--beta", "message-batches-2024-09-24",
		)
	})

	t.Run("piping data", func(t *testing.T) {
		// Test piping YAML data over stdin
		pipeData := []byte("" +
			"agent: agent_011CZkYpogX7uDKUyvBTophP\n" +
			"environment_id: env_011CZkZ9X2dpNyB7HsEFoRfW\n" +
			"metadata:\n" +
			"  foo: string\n" +
			"resources:\n" +
			"  - file_id: file_011CNha8iCJcU1wXNR6q4V8w\n" +
			"    type: file\n" +
			"    mount_path: /uploads/receipt.pdf\n" +
			"title: 'Order #1234 inquiry'\n" +
			"vault_ids:\n" +
			"  - string\n")
		mocktest.TestRunMockTestWithPipeAndFlags(
			t, pipeData,
			"--api-key", "string",
			"beta:sessions", "create",
			"--beta", "message-batches-2024-09-24",
		)
	})
}

func TestBetaSessionsRetrieve(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:sessions", "retrieve",
			"--session-id", "sesn_011CZkZAtmR3yMPDzynEDxu7",
			"--beta", "message-batches-2024-09-24",
		)
	})
}

func TestBetaSessionsUpdate(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:sessions", "update",
			"--session-id", "sesn_011CZkZAtmR3yMPDzynEDxu7",
			"--metadata", "{foo: string}",
			"--title", "Order #1234 inquiry",
			"--vault-id", "string",
			"--beta", "message-batches-2024-09-24",
		)
	})

	t.Run("piping data", func(t *testing.T) {
		// Test piping YAML data over stdin
		pipeData := []byte("" +
			"metadata:\n" +
			"  foo: string\n" +
			"title: 'Order #1234 inquiry'\n" +
			"vault_ids:\n" +
			"  - string\n")
		mocktest.TestRunMockTestWithPipeAndFlags(
			t, pipeData,
			"--api-key", "string",
			"beta:sessions", "update",
			"--session-id", "sesn_011CZkZAtmR3yMPDzynEDxu7",
			"--beta", "message-batches-2024-09-24",
		)
	})
}

func TestBetaSessionsList(t *testing.T) {
	t.Skip("buildURL drops path-level query params (SDK-4349)")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:sessions", "list",
			"--max-items", "10",
			"--agent-id", "agent_id",
			"--agent-version", "0",
			"--created-at-gt", "'2019-12-27T18:11:19.117Z'",
			"--created-at-gte", "'2019-12-27T18:11:19.117Z'",
			"--created-at-lt", "'2019-12-27T18:11:19.117Z'",
			"--created-at-lte", "'2019-12-27T18:11:19.117Z'",
			"--include-archived=true",
			"--limit", "0",
			"--order", "asc",
			"--page", "page",
			"--beta", "message-batches-2024-09-24",
		)
	})
}

func TestBetaSessionsDelete(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:sessions", "delete",
			"--session-id", "sesn_011CZkZAtmR3yMPDzynEDxu7",
			"--beta", "message-batches-2024-09-24",
		)
	})
}

func TestBetaSessionsArchive(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:sessions", "archive",
			"--session-id", "sesn_011CZkZAtmR3yMPDzynEDxu7",
			"--beta", "message-batches-2024-09-24",
		)
	})
}
