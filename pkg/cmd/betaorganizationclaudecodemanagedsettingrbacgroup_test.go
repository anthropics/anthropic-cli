package cmd

import (
	"testing"

	"github.com/anthropics/anthropic-cli/v2/internal/mocktest"
)

func TestBetaOrganizationClaudeCodeManagedSettingsRBACGroupsRetrieve(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "string")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--base-url", mocktest.ServerWithQuery(t, "beta=true"),
			"beta:organization:claude-code:managed-settings:rbac-groups", "retrieve",
			"--rbac-group-id", "rbac_group_id",
			"--beta", "computer-use-2024-10-22",
		)
	})
}

func TestBetaOrganizationClaudeCodeManagedSettingsRBACGroupsList(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "string")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--base-url", mocktest.ServerWithQuery(t, "beta=true"),
			"beta:organization:claude-code:managed-settings:rbac-groups", "list",
			"--max-items", "10",
			"--limit", "1",
			"--page", "page",
			"--beta", "computer-use-2024-10-22",
		)
	})
}

func TestBetaOrganizationClaudeCodeManagedSettingsRBACGroupsDelete(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "string")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--base-url", mocktest.ServerWithQuery(t, "beta=true"),
			"beta:organization:claude-code:managed-settings:rbac-groups", "delete",
			"--rbac-group-id", "rbac_group_id",
			"--beta", "computer-use-2024-10-22",
		)
	})
}

func TestBetaOrganizationClaudeCodeManagedSettingsRBACGroupsSet(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "string")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--base-url", mocktest.ServerWithQuery(t, "beta=true"),
			"beta:organization:claude-code:managed-settings:rbac-groups", "set",
			"--rbac-group-id", "rbac_group_id",
			"--settings", "{foo: bar}",
			"--beta", "computer-use-2024-10-22",
		)
	})

	t.Run("piping data", func(t *testing.T) {
		// Test piping YAML data over stdin
		pipeData := []byte("" +
			"settings:\n" +
			"  foo: bar\n")
		mocktest.TestRunMockTestWithPipeAndFlags(
			t, pipeData,
			"--base-url", mocktest.ServerWithQuery(t, "beta=true"),
			"beta:organization:claude-code:managed-settings:rbac-groups", "set",
			"--rbac-group-id", "rbac_group_id",
			"--beta", "computer-use-2024-10-22",
		)
	})
}
