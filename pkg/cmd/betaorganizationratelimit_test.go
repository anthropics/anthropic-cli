// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cmd

import (
	"testing"

	"github.com/anthropics/anthropic-cli/internal/mocktest"
)

func TestBetaOrganizationRateLimitsList(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:organization:rate-limits", "list",
			"--max-items", "10",
			"--group-type", "batch",
			"--limit", "1",
			"--model", "model",
			"--page", "page",
		)
	})
}
