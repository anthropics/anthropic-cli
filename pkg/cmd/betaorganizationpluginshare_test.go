package cmd

import (
	"testing"

	"github.com/anthropics/anthropic-cli/v2/internal/mocktest"
)

func TestBetaOrganizationPluginsSharesList(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "string")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--base-url", mocktest.ServerWithQuery(t, "beta=true"),
			"beta:organization:plugins:shares", "list",
			"--max-items", "10",
			"--plugin-id", "plugin_id",
			"--limit", "1",
			"--organization-id", "organization_id",
			"--page", "page",
			"--target-type", "organization",
			"--beta", "computer-use-2024-10-22",
		)
	})
}
