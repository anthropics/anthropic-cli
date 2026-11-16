package cmd

import (
	"testing"

	"github.com/anthropics/anthropic-cli/v2/internal/mocktest"
)

func TestBetaSessionsThreadsRetrieve(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "string")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--base-url", mocktest.ServerWithQuery(t, "beta=true"),
			"beta:sessions:threads", "retrieve",
			"--session-id", "sesn_011CZkZAtmR3yMPDzynEDxu7",
			"--thread-id", "sthr_011CZkZVWa6oJjw1rgXZpnBt",
			"--beta", "computer-use-2024-10-22",
			"--workspace-id", "wrkspc_011CZkZaBF1tNoB5wlCeusgy",
		)
	})
}

func TestBetaSessionsThreadsList(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "string")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--base-url", mocktest.ServerWithQuery(t, "beta=true"),
			"beta:sessions:threads", "list",
			"--max-items", "10",
			"--session-id", "sesn_011CZkZAtmR3yMPDzynEDxu7",
			"--limit", "0",
			"--page", "page",
			"--status", "running",
			"--beta", "computer-use-2024-10-22",
			"--workspace-id", "wrkspc_011CZkZaBF1tNoB5wlCeusgy",
		)
	})
}

func TestBetaSessionsThreadsArchive(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "string")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--base-url", mocktest.ServerWithQuery(t, "beta=true"),
			"beta:sessions:threads", "archive",
			"--session-id", "sesn_011CZkZAtmR3yMPDzynEDxu7",
			"--thread-id", "sthr_011CZkZVWa6oJjw1rgXZpnBt",
			"--beta", "computer-use-2024-10-22",
			"--workspace-id", "wrkspc_011CZkZaBF1tNoB5wlCeusgy",
		)
	})
}
