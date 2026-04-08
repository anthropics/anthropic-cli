// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cmd

import (
	"testing"

	"github.com/anthropics/anthropic-cli/internal/mocktest"
)

func TestBetaAgentsVersionsList(t *testing.T) {
	t.Skip("buildURL drops path-level query params (SDK-4349)")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:agents:versions", "list",
			"--max-items", "10",
			"--agent-id", "agent_011CZkYpogX7uDKUyvBTophP",
			"--limit", "0",
			"--page", "page",
			"--beta", "message-batches-2024-09-24",
		)
	})
}
