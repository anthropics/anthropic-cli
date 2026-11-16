package cmd

import (
	"testing"

	"github.com/anthropics/anthropic-cli/v2/internal/mocktest"
)

func TestBetaAgentsVersionsList(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "string")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--base-url", mocktest.ServerWithQuery(t, "beta=true"),
			"beta:agents:versions", "list",
			"--max-items", "10",
			"--agent-id", "agent_011CZkYpogX7uDKUyvBTophP",
			"--limit", "0",
			"--page", "page",
			"--beta", "computer-use-2024-10-22",
			"--workspace-id", "wrkspc_011CZkZaBF1tNoB5wlCeusgy",
		)
	})
}
