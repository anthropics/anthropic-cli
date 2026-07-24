// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cmd

import (
	"testing"

	"github.com/anthropics/anthropic-cli/internal/mocktest"
)

func TestBetaTunnelsCreate(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:tunnels", "create",
			"--display-name", "x",
			"--beta", "message-batches-2024-09-24",
		)
	})

	t.Run("piping data", func(t *testing.T) {
		// Test piping YAML data over stdin
		pipeData := []byte("display_name: x")
		mocktest.TestRunMockTestWithPipeAndFlags(
			t, pipeData,
			"--api-key", "string",
			"beta:tunnels", "create",
			"--beta", "message-batches-2024-09-24",
		)
	})
}

func TestBetaTunnelsRetrieve(t *testing.T) {
	t.Skip("buildURL drops path-level query params (SDK-4349)")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:tunnels", "retrieve",
			"--tunnel-id", "tunnel_id",
			"--beta", "message-batches-2024-09-24",
		)
	})
}

func TestBetaTunnelsList(t *testing.T) {
	t.Skip("buildURL drops path-level query params (SDK-4349)")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:tunnels", "list",
			"--max-items", "10",
			"--include-archived=true",
			"--limit", "0",
			"--page", "page",
			"--beta", "message-batches-2024-09-24",
		)
	})
}

func TestBetaTunnelsArchive(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:tunnels", "archive",
			"--tunnel-id", "tunnel_id",
			"--beta", "message-batches-2024-09-24",
		)
	})
}

func TestBetaTunnelsRevealToken(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:tunnels", "reveal-token",
			"--tunnel-id", "tunnel_id",
			"--beta", "message-batches-2024-09-24",
		)
	})
}

func TestBetaTunnelsRotateToken(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:tunnels", "rotate-token",
			"--tunnel-id", "tunnel_id",
			"--reason", "reason",
			"--beta", "message-batches-2024-09-24",
		)
	})

	t.Run("piping data", func(t *testing.T) {
		// Test piping YAML data over stdin
		pipeData := []byte("reason: reason")
		mocktest.TestRunMockTestWithPipeAndFlags(
			t, pipeData,
			"--api-key", "string",
			"beta:tunnels", "rotate-token",
			"--tunnel-id", "tunnel_id",
			"--beta", "message-batches-2024-09-24",
		)
	})
}
