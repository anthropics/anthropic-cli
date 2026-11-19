package cmd

import (
	"testing"

	"github.com/anthropics/anthropic-cli/v2/internal/mocktest"
)

func TestBetaOrganizationClaudeCodeManagedSettingsRBACGroupDefaultRetrieve(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "string")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--base-url", mocktest.ServerWithQuery(t, "beta=true"),
			"beta:organization:claude-code:managed-settings:rbac-group-default", "retrieve",
			"--beta", "computer-use-2024-10-22",
		)
	})
}

func TestBetaOrganizationClaudeCodeManagedSettingsRBACGroupDefaultUpdate(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "string")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--base-url", mocktest.ServerWithQuery(t, "beta=true"),
			"beta:organization:claude-code:managed-settings:rbac-group-default", "update",
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
			"beta:organization:claude-code:managed-settings:rbac-group-default", "update",
			"--beta", "computer-use-2024-10-22",
		)
	})
}
