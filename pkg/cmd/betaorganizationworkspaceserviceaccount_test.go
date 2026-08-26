// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cmd

import (
	"testing"

	"github.com/anthropics/anthropic-cli/internal/mocktest"
)

func TestBetaOrganizationWorkspacesServiceAccountsRetrieve(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:organization:workspaces:service-accounts", "retrieve",
			"--workspace-id", "workspace_id",
			"--service-account-id", "service_account_id",
			"--beta", "message-batches-2024-09-24",
		)
	})
}

func TestBetaOrganizationWorkspacesServiceAccountsUpdate(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:organization:workspaces:service-accounts", "update",
			"--workspace-id", "workspace_id",
			"--service-account-id", "service_account_id",
			"--workspace-role", "workspace_admin",
			"--beta", "message-batches-2024-09-24",
		)
	})

	t.Run("piping data", func(t *testing.T) {
		// Test piping YAML data over stdin
		pipeData := []byte("workspace_role: workspace_admin")
		mocktest.TestRunMockTestWithPipeAndFlags(
			t, pipeData,
			"--api-key", "string",
			"beta:organization:workspaces:service-accounts", "update",
			"--workspace-id", "workspace_id",
			"--service-account-id", "service_account_id",
			"--beta", "message-batches-2024-09-24",
		)
	})
}

func TestBetaOrganizationWorkspacesServiceAccountsList(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:organization:workspaces:service-accounts", "list",
			"--max-items", "10",
			"--workspace-id", "workspace_id",
			"--limit", "1",
			"--page", "page",
			"--beta", "message-batches-2024-09-24",
		)
	})
}

func TestBetaOrganizationWorkspacesServiceAccountsAdd(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:organization:workspaces:service-accounts", "add",
			"--workspace-id", "workspace_id",
			"--service-account-id", "service_account_id",
			"--workspace-role", "workspace_admin",
			"--beta", "message-batches-2024-09-24",
		)
	})

	t.Run("piping data", func(t *testing.T) {
		// Test piping YAML data over stdin
		pipeData := []byte("" +
			"service_account_id: service_account_id\n" +
			"workspace_role: workspace_admin\n")
		mocktest.TestRunMockTestWithPipeAndFlags(
			t, pipeData,
			"--api-key", "string",
			"beta:organization:workspaces:service-accounts", "add",
			"--workspace-id", "workspace_id",
			"--beta", "message-batches-2024-09-24",
		)
	})
}

func TestBetaOrganizationWorkspacesServiceAccountsRemove(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:organization:workspaces:service-accounts", "remove",
			"--workspace-id", "workspace_id",
			"--service-account-id", "service_account_id",
			"--beta", "message-batches-2024-09-24",
		)
	})
}
