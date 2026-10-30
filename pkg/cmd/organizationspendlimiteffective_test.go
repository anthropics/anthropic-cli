package cmd

import (
	"testing"

	"github.com/anthropics/anthropic-cli/v2/internal/mocktest"
)

func TestOrganizationSpendLimitsEffectiveList(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "string")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"organization:spend-limits:effective", "list",
			"--max-items", "10",
			"--limit", "1",
			"--page", "page",
			"--period", "daily",
			"--user-id", "string",
		)
	})
}
