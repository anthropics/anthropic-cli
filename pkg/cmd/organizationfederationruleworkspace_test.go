package cmd

import (
	"testing"

	"github.com/anthropics/anthropic-cli/v2/internal/mocktest"
)

func TestOrganizationFederationRulesWorkspacesList(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "string")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"organization:federation:rules:workspaces", "list",
			"--max-items", "10",
			"--federation-rule-id", "federation_rule_id",
			"--limit", "1",
			"--page", "page",
		)
	})
}

func TestOrganizationFederationRulesWorkspacesAdd(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "string")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"organization:federation:rules:workspaces", "add",
			"--federation-rule-id", "federation_rule_id",
			"--workspace-id", "wrkspc_01JwQvzr7rXLA5AGx3HKfFUJ",
		)
	})

	t.Run("piping data", func(t *testing.T) {
		// Test piping YAML data over stdin
		pipeData := []byte("workspace_id: wrkspc_01JwQvzr7rXLA5AGx3HKfFUJ")
		mocktest.TestRunMockTestWithPipeAndFlags(
			t, pipeData,
			"organization:federation:rules:workspaces", "add",
			"--federation-rule-id", "federation_rule_id",
		)
	})
}

func TestOrganizationFederationRulesWorkspacesRemove(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "string")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"organization:federation:rules:workspaces", "remove",
			"--federation-rule-id", "federation_rule_id",
			"--workspace-id", "workspace_id",
		)
	})
}
