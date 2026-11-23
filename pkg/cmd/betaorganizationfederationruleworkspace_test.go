package cmd

import (
	"testing"

	"github.com/anthropics/anthropic-cli/v2/internal/mocktest"
)

func TestBetaOrganizationFederationRulesWorkspacesList(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "string")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--base-url", mocktest.ServerWithQuery(t, "beta=true"),
			"beta:organization:federation:rules:workspaces", "list",
			"--max-items", "10",
			"--federation-rule-id", "federation_rule_id",
			"--limit", "1",
			"--page", "page",
			"--beta", "computer-use-2024-10-22",
		)
	})
}

func TestBetaOrganizationFederationRulesWorkspacesAdd(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "string")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--base-url", mocktest.ServerWithQuery(t, "beta=true"),
			"beta:organization:federation:rules:workspaces", "add",
			"--federation-rule-id", "federation_rule_id",
			"--workspace-id", "wrkspc_01JwQvzr7rXLA5AGx3HKfFUJ",
			"--beta", "computer-use-2024-10-22",
		)
	})

	t.Run("piping data", func(t *testing.T) {
		// Test piping YAML data over stdin
		pipeData := []byte("workspace_id: wrkspc_01JwQvzr7rXLA5AGx3HKfFUJ")
		mocktest.TestRunMockTestWithPipeAndFlags(
			t, pipeData,
			"--base-url", mocktest.ServerWithQuery(t, "beta=true"),
			"beta:organization:federation:rules:workspaces", "add",
			"--federation-rule-id", "federation_rule_id",
			"--beta", "computer-use-2024-10-22",
		)
	})
}

func TestBetaOrganizationFederationRulesWorkspacesRemove(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "string")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--base-url", mocktest.ServerWithQuery(t, "beta=true"),
			"beta:organization:federation:rules:workspaces", "remove",
			"--federation-rule-id", "federation_rule_id",
			"--workspace-id", "workspace_id",
			"--beta", "computer-use-2024-10-22",
		)
	})
}
