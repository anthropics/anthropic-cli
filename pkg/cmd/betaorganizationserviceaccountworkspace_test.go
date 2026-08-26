// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cmd

import (
	"testing"

	"github.com/anthropics/anthropic-cli/internal/mocktest"
)

func TestBetaOrganizationServiceAccountsWorkspacesList(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:organization:service-accounts:workspaces", "list",
			"--max-items", "10",
			"--service-account-id", "service_account_id",
			"--limit", "1",
			"--page", "page",
			"--beta", "message-batches-2024-09-24",
		)
	})
}

func TestBetaOrganizationServiceAccountsWorkspacesAdd(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:organization:service-accounts:workspaces", "add",
			"--service-account-id", "service_account_id",
			"--workspace-id", "workspace_id",
			"--workspace-role", "workspace_admin",
			"--beta", "message-batches-2024-09-24",
		)
	})

	t.Run("piping data", func(t *testing.T) {
		// Test piping YAML data over stdin
		pipeData := []byte("" +
			"workspace_id: workspace_id\n" +
			"workspace_role: workspace_admin\n")
		mocktest.TestRunMockTestWithPipeAndFlags(
			t, pipeData,
			"--api-key", "string",
			"beta:organization:service-accounts:workspaces", "add",
			"--service-account-id", "service_account_id",
			"--beta", "message-batches-2024-09-24",
		)
	})
}

func TestBetaOrganizationServiceAccountsWorkspacesRemove(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:organization:service-accounts:workspaces", "remove",
			"--service-account-id", "service_account_id",
			"--workspace-id", "workspace_id",
			"--beta", "message-batches-2024-09-24",
		)
	})
}
