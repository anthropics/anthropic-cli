package cmd

import (
	"testing"

	"github.com/anthropics/anthropic-cli/v2/internal/mocktest"
)

func TestBetaComplianceReviewItemsRetrieve(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "string")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--base-url", mocktest.ServerWithQuery(t, "beta=true"),
			"beta:compliance:review:items", "retrieve",
			"--item-id", "item_id",
			"--beta", "computer-use-2024-10-22",
			"--user-profile-id", "anthropic-user-profile-id",
		)
	})
}
