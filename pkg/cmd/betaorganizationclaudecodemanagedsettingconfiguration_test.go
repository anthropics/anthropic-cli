package cmd

import (
	"testing"

	"github.com/anthropics/anthropic-cli/v2/internal/mocktest"
)

func TestBetaOrganizationClaudeCodeManagedSettingsConfigurationRetrieve(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "string")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--base-url", mocktest.ServerWithQuery(t, "beta=true"),
			"beta:organization:claude-code:managed-settings:configuration", "retrieve",
			"--beta", "computer-use-2024-10-22",
		)
	})
}

func TestBetaOrganizationClaudeCodeManagedSettingsConfigurationUpdate(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "string")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--base-url", mocktest.ServerWithQuery(t, "beta=true"),
			"beta:organization:claude-code:managed-settings:configuration", "update",
			"--mode", "{type: organization}",
			"--rbac-group-order", "string",
			"--rbac-group-order-version", "rbac_group_order_version",
			"--beta", "computer-use-2024-10-22",
		)
	})

	t.Run("piping data", func(t *testing.T) {
		// Test piping YAML data over stdin
		pipeData := []byte("" +
			"mode:\n" +
			"  type: organization\n" +
			"rbac_group_order:\n" +
			"  - string\n" +
			"rbac_group_order_version: rbac_group_order_version\n")
		mocktest.TestRunMockTestWithPipeAndFlags(
			t, pipeData,
			"--base-url", mocktest.ServerWithQuery(t, "beta=true"),
			"beta:organization:claude-code:managed-settings:configuration", "update",
			"--beta", "computer-use-2024-10-22",
		)
	})
}
