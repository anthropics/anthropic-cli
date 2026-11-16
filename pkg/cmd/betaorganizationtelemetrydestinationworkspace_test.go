package cmd

import (
	"testing"

	"github.com/anthropics/anthropic-cli/v2/internal/mocktest"
)

func TestBetaOrganizationTelemetryDestinationsWorkspacesList(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "string")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--base-url", mocktest.ServerWithQuery(t, "beta=true"),
			"beta:organization:telemetry-destinations:workspaces", "list",
			"--max-items", "10",
			"--destination-id", "destination_id",
			"--limit", "0",
			"--page", "page",
			"--beta", "computer-use-2024-10-22",
		)
	})
}

func TestBetaOrganizationTelemetryDestinationsWorkspacesAdd(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "string")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--base-url", mocktest.ServerWithQuery(t, "beta=true"),
			"beta:organization:telemetry-destinations:workspaces", "add",
			"--destination-id", "destination_id",
			"--workspace-id", "workspace_id",
			"--beta", "computer-use-2024-10-22",
		)
	})
}

func TestBetaOrganizationTelemetryDestinationsWorkspacesRemove(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "string")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--base-url", mocktest.ServerWithQuery(t, "beta=true"),
			"beta:organization:telemetry-destinations:workspaces", "remove",
			"--destination-id", "destination_id",
			"--workspace-id", "workspace_id",
			"--beta", "computer-use-2024-10-22",
		)
	})
}
