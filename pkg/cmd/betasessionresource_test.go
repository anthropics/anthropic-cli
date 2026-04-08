// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cmd

import (
	"testing"

	"github.com/anthropics/anthropic-cli/internal/mocktest"
)

func TestBetaSessionsResourcesRetrieve(t *testing.T) {
	t.Skip("prism can't find endpoint with beta only tag")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:sessions:resources", "retrieve",
			"--session-id", "sesn_011CZkZAtmR3yMPDzynEDxu7",
			"--resource-id", "sesrsc_011CZkZBJq5dWxk9fVLNcPht",
			"--beta", "message-batches-2024-09-24",
		)
	})
}

func TestBetaSessionsResourcesUpdate(t *testing.T) {
	t.Skip("prism can't find endpoint with beta only tag")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:sessions:resources", "update",
			"--session-id", "sesn_011CZkZAtmR3yMPDzynEDxu7",
			"--resource-id", "sesrsc_011CZkZBJq5dWxk9fVLNcPht",
			"--authorization-token", "ghp_exampletoken",
			"--beta", "message-batches-2024-09-24",
		)
	})

	t.Run("piping data", func(t *testing.T) {
		// Test piping YAML data over stdin
		pipeData := []byte("authorization_token: ghp_exampletoken")
		mocktest.TestRunMockTestWithPipeAndFlags(
			t, pipeData,
			"--api-key", "string",
			"beta:sessions:resources", "update",
			"--session-id", "sesn_011CZkZAtmR3yMPDzynEDxu7",
			"--resource-id", "sesrsc_011CZkZBJq5dWxk9fVLNcPht",
			"--beta", "message-batches-2024-09-24",
		)
	})
}

func TestBetaSessionsResourcesList(t *testing.T) {
	t.Skip("prism can't find endpoint with beta only tag")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:sessions:resources", "list",
			"--max-items", "10",
			"--session-id", "sesn_011CZkZAtmR3yMPDzynEDxu7",
			"--limit", "0",
			"--page", "page",
			"--beta", "message-batches-2024-09-24",
		)
	})
}

func TestBetaSessionsResourcesDelete(t *testing.T) {
	t.Skip("prism can't find endpoint with beta only tag")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:sessions:resources", "delete",
			"--session-id", "sesn_011CZkZAtmR3yMPDzynEDxu7",
			"--resource-id", "sesrsc_011CZkZBJq5dWxk9fVLNcPht",
			"--beta", "message-batches-2024-09-24",
		)
	})
}

func TestBetaSessionsResourcesAdd(t *testing.T) {
	t.Skip("prism can't find endpoint with beta only tag")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:sessions:resources", "add",
			"--session-id", "sesn_011CZkZAtmR3yMPDzynEDxu7",
			"--file-id", "file_011CNha8iCJcU1wXNR6q4V8w",
			"--type", "file",
			"--mount-path", "/uploads/receipt.pdf",
			"--beta", "message-batches-2024-09-24",
		)
	})

	t.Run("piping data", func(t *testing.T) {
		// Test piping YAML data over stdin
		pipeData := []byte("" +
			"file_id: file_011CNha8iCJcU1wXNR6q4V8w\n" +
			"type: file\n" +
			"mount_path: /uploads/receipt.pdf\n")
		mocktest.TestRunMockTestWithPipeAndFlags(
			t, pipeData,
			"--api-key", "string",
			"beta:sessions:resources", "add",
			"--session-id", "sesn_011CZkZAtmR3yMPDzynEDxu7",
			"--beta", "message-batches-2024-09-24",
		)
	})
}
