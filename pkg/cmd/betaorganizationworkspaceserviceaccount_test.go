package cmd

import (
	"testing"

	"github.com/anthropics/anthropic-cli/v2/internal/mocktest"
)

func TestBetaOrganizationWorkspacesServiceAccountsRetrieve(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "string")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--base-url", mocktest.ServerWithQuery(t, "beta=true"),
			"beta:organization:workspaces:service-accounts", "retrieve",
			"--workspace-id", "workspace_id",
			"--service-account-id", "service_account_id",
			"--beta", "computer-use-2024-10-22",
		)
	})
}

func TestBetaOrganizationWorkspacesServiceAccountsUpdate(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "string")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--base-url", mocktest.ServerWithQuery(t, "beta=true"),
			"beta:organization:workspaces:service-accounts", "update",
			"--workspace-id", "workspace_id",
			"--service-account-id", "service_account_id",
			"--workspace-role", "workspace_admin",
			"--beta", "computer-use-2024-10-22",
		)
	})

	t.Run("piping data", func(t *testing.T) {
		// Test piping YAML data over stdin
		pipeData := []byte("workspace_role: workspace_admin")
		mocktest.TestRunMockTestWithPipeAndFlags(
			t, pipeData,
			"--base-url", mocktest.ServerWithQuery(t, "beta=true"),
			"beta:organization:workspaces:service-accounts", "update",
			"--workspace-id", "workspace_id",
			"--service-account-id", "service_account_id",
			"--beta", "computer-use-2024-10-22",
		)
	})
}

func TestBetaOrganizationWorkspacesServiceAccountsList(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "string")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--base-url", mocktest.ServerWithQuery(t, "beta=true"),
			"beta:organization:workspaces:service-accounts", "list",
			"--max-items", "10",
			"--workspace-id", "workspace_id",
			"--limit", "1",
			"--page", "page",
			"--beta", "computer-use-2024-10-22",
		)
	})
}

func TestBetaOrganizationWorkspacesServiceAccountsAdd(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "string")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--base-url", mocktest.ServerWithQuery(t, "beta=true"),
			"beta:organization:workspaces:service-accounts", "add",
			"--workspace-id", "workspace_id",
			"--service-account-id", "svac_01SDCCSbTxrXDpWc1phhtcfK",
			"--workspace-role", "workspace_admin",
			"--beta", "computer-use-2024-10-22",
		)
	})

	t.Run("piping data", func(t *testing.T) {
		// Test piping YAML data over stdin
		pipeData := []byte("" +
			"service_account_id: svac_01SDCCSbTxrXDpWc1phhtcfK\n" +
			"workspace_role: workspace_admin\n")
		mocktest.TestRunMockTestWithPipeAndFlags(
			t, pipeData,
			"--base-url", mocktest.ServerWithQuery(t, "beta=true"),
			"beta:organization:workspaces:service-accounts", "add",
			"--workspace-id", "workspace_id",
			"--beta", "computer-use-2024-10-22",
		)
	})
}

func TestBetaOrganizationWorkspacesServiceAccountsRemove(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "string")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--base-url", mocktest.ServerWithQuery(t, "beta=true"),
			"beta:organization:workspaces:service-accounts", "remove",
			"--workspace-id", "workspace_id",
			"--service-account-id", "service_account_id",
			"--beta", "computer-use-2024-10-22",
		)
	})
}
