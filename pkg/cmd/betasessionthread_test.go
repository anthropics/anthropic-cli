// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cmd

import (
	"testing"

	"github.com/anthropics/anthropic-cli/internal/mocktest"
)

func TestBetaSessionsThreadsRetrieve(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:sessions:threads", "retrieve",
			"--session-id", "sesn_011CZkZAtmR3yMPDzynEDxu7",
			"--thread-id", "sthr_011CZkZVWa6oIjw0rgXZpnBt",
			"--beta", "message-batches-2024-09-24",
		)
	})
}

func TestBetaSessionsThreadsList(t *testing.T) {
	t.Skip("buildURL drops path-level query params (SDK-4349)")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:sessions:threads", "list",
			"--max-items", "10",
			"--session-id", "sesn_011CZkZAtmR3yMPDzynEDxu7",
			"--limit", "0",
			"--page", "page",
			"--beta", "message-batches-2024-09-24",
		)
	})
}

func TestBetaSessionsThreadsArchive(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:sessions:threads", "archive",
			"--session-id", "sesn_011CZkZAtmR3yMPDzynEDxu7",
			"--thread-id", "sthr_011CZkZVWa6oIjw0rgXZpnBt",
			"--beta", "message-batches-2024-09-24",
		)
	})
}
