package cmd

import (
	"testing"

	"github.com/anthropics/anthropic-cli/v2/internal/mocktest"
	"github.com/anthropics/anthropic-cli/v2/internal/requestflag"
)

func TestBetaComplianceReviewCasesRetrieve(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "string")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--base-url", mocktest.ServerWithQuery(t, "beta=true"),
			"beta:compliance:review:cases", "retrieve",
			"--case-id", "case_id",
			"--beta", "computer-use-2024-10-22",
		)
	})
}

func TestBetaComplianceReviewCasesList(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "string")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--base-url", mocktest.ServerWithQuery(t, "beta=true"),
			"beta:compliance:review:cases", "list",
			"--max-items", "10",
			"--created-at-gt", "'2019-12-27T18:11:19.117Z'",
			"--created-at-gte", "'2019-12-27T18:11:19.117Z'",
			"--created-at-lt", "'2019-12-27T18:11:19.117Z'",
			"--created-at-lte", "'2019-12-27T18:11:19.117Z'",
			"--limit", "0",
			"--order", "asc",
			"--order-by", "created_at",
			"--page", "page",
			"--status", "open",
			"--updated-at-gt", "'2019-12-27T18:11:19.117Z'",
			"--updated-at-gte", "'2019-12-27T18:11:19.117Z'",
			"--updated-at-lt", "'2019-12-27T18:11:19.117Z'",
			"--updated-at-lte", "'2019-12-27T18:11:19.117Z'",
			"--user-profile-id", "user_profile_id",
			"--beta", "computer-use-2024-10-22",
		)
	})
}

func TestBetaComplianceReviewCasesRespond(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "string")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--base-url", mocktest.ServerWithQuery(t, "beta=true"),
			"beta:compliance:review:cases", "respond",
			"--case-id", "case_id",
			"--determination", "no_violation",
			"--policy-category", "cyber",
			"--verdict", "{determination: violation, item_id: item_id, policy_category: cyber}",
			"--beta", "computer-use-2024-10-22",
		)
	})

	t.Run("inner flags", func(t *testing.T) {
		// Check that inner flags have been set up correctly
		requestflag.CheckInnerFlags(betaComplianceReviewCasesRespond)

		// Alternative argument passing style using inner flags
		mocktest.TestRunMockTestWithFlags(
			t,
			"--base-url", mocktest.ServerWithQuery(t, "beta=true"),
			"beta:compliance:review:cases", "respond",
			"--case-id", "case_id",
			"--determination", "no_violation",
			"--policy-category", "cyber",
			"--verdict.determination", "violation",
			"--verdict.item-id", "item_id",
			"--verdict.policy-category", "cyber",
			"--beta", "computer-use-2024-10-22",
		)
	})

	t.Run("piping data", func(t *testing.T) {
		// Test piping YAML data over stdin
		pipeData := []byte("" +
			"determination: no_violation\n" +
			"policy_categories:\n" +
			"  - cyber\n" +
			"verdicts:\n" +
			"  - determination: violation\n" +
			"    item_id: item_id\n" +
			"    policy_category: cyber\n")
		mocktest.TestRunMockTestWithPipeAndFlags(
			t, pipeData,
			"--base-url", mocktest.ServerWithQuery(t, "beta=true"),
			"beta:compliance:review:cases", "respond",
			"--case-id", "case_id",
			"--beta", "computer-use-2024-10-22",
		)
	})
}
