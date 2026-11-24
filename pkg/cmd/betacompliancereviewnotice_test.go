package cmd

import (
	"testing"

	"github.com/anthropics/anthropic-cli/v2/internal/mocktest"
)

func TestBetaComplianceReviewNoticesList(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "string")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--base-url", mocktest.ServerWithQuery(t, "beta=true"),
			"beta:compliance:review:notices", "list",
			"--max-items", "10",
			"--case-id", "case_id",
			"--created-at-gt", "'2019-12-27T18:11:19.117Z'",
			"--created-at-gte", "'2019-12-27T18:11:19.117Z'",
			"--created-at-lt", "'2019-12-27T18:11:19.117Z'",
			"--created-at-lte", "'2019-12-27T18:11:19.117Z'",
			"--limit", "0",
			"--order", "asc",
			"--page", "page",
			"--beta", "computer-use-2024-10-22",
		)
	})
}
