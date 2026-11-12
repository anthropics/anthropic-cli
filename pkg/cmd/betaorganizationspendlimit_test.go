package cmd

import (
	"testing"

	"github.com/anthropics/anthropic-cli/v2/internal/mocktest"
)

func TestBetaOrganizationSpendLimitsRetrieve(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "string")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--base-url", mocktest.ServerWithQuery(t, "beta=true"),
			"beta:organization:spend-limits", "retrieve",
			"--spend-limit-id", "spend_limit_id",
		)
	})
}

func TestBetaOrganizationSpendLimitsList(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "string")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--base-url", mocktest.ServerWithQuery(t, "beta=true"),
			"beta:organization:spend-limits", "list",
			"--max-items", "10",
			"--limit", "1",
			"--page", "page",
			"--scope-type", "organization",
			"--workspace-id", "workspace_id",
			"--beta", "computer-use-2024-10-22",
		)
	})
}

func TestBetaOrganizationSpendLimitsDelete(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "string")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--base-url", mocktest.ServerWithQuery(t, "beta=true"),
			"beta:organization:spend-limits", "delete",
			"--spend-limit-id", "spend_limit_id",
		)
	})
}

func TestBetaOrganizationSpendLimitsSet(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "string")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--base-url", mocktest.ServerWithQuery(t, "beta=true"),
			"beta:organization:spend-limits", "set",
			"--amount", "50000",
			"--scope", "{type: user, user_id: user_01WxYzAbCdEfGhJkMnPqRsTu}",
			"--period", "monthly",
			"--beta", "computer-use-2024-10-22",
		)
	})

	t.Run("piping data", func(t *testing.T) {
		// Test piping YAML data over stdin
		pipeData := []byte("" +
			"amount: '50000'\n" +
			"scope:\n" +
			"  type: user\n" +
			"  user_id: user_01WxYzAbCdEfGhJkMnPqRsTu\n" +
			"period: monthly\n")
		mocktest.TestRunMockTestWithPipeAndFlags(
			t, pipeData,
			"--base-url", mocktest.ServerWithQuery(t, "beta=true"),
			"beta:organization:spend-limits", "set",
			"--beta", "computer-use-2024-10-22",
		)
	})
}
