// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cmd

import (
	"testing"

	"github.com/anthropics/anthropic-cli/internal/mocktest"
)

func TestBetaOrganizationWorkspacesRateLimitsList(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:organization:workspaces:rate-limits", "list",
			"--max-items", "10",
			"--workspace-id", "workspace_id",
			"--group-type", "batch",
			"--limit", "1",
			"--page", "page",
		)
	})
}
