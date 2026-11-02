package cmd

import (
	"testing"

	"github.com/anthropics/anthropic-cli/v2/internal/mocktest"
)

func TestOrganizationAnalyticsSummariesList(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "string")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"organization:analytics:summaries", "list",
			"--max-items", "10",
			"--starting-date", "'2019-12-27'",
			"--ending-date", "'2019-12-27'",
			"--filter", "string",
			"--limit", "1",
			"--page", "page",
		)
	})
}
