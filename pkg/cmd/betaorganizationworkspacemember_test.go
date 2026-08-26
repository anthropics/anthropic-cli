// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cmd

import (
	"testing"

	"github.com/anthropics/anthropic-cli/internal/mocktest"
)

func TestBetaOrganizationWorkspacesMembersRetrieve(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:organization:workspaces:members", "retrieve",
			"--workspace-id", "workspace_id",
			"--user-id", "user_id",
		)
	})
}

func TestBetaOrganizationWorkspacesMembersUpdate(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:organization:workspaces:members", "update",
			"--workspace-id", "workspace_id",
			"--user-id", "user_id",
			"--workspace-role", "workspace_admin",
		)
	})

	t.Run("piping data", func(t *testing.T) {
		// Test piping YAML data over stdin
		pipeData := []byte("workspace_role: workspace_admin")
		mocktest.TestRunMockTestWithPipeAndFlags(
			t, pipeData,
			"--api-key", "string",
			"beta:organization:workspaces:members", "update",
			"--workspace-id", "workspace_id",
			"--user-id", "user_id",
		)
	})
}

func TestBetaOrganizationWorkspacesMembersList(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:organization:workspaces:members", "list",
			"--max-items", "10",
			"--workspace-id", "workspace_id",
			"--after-id", "after_id",
			"--before-id", "before_id",
			"--limit", "1",
		)
	})
}

func TestBetaOrganizationWorkspacesMembersAdd(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:organization:workspaces:members", "add",
			"--workspace-id", "workspace_id",
			"--user-id", "user_01WCz1FkmYMm4gnmykNKUu3Q",
			"--workspace-role", "workspace_admin",
		)
	})

	t.Run("piping data", func(t *testing.T) {
		// Test piping YAML data over stdin
		pipeData := []byte("" +
			"user_id: user_01WCz1FkmYMm4gnmykNKUu3Q\n" +
			"workspace_role: workspace_admin\n")
		mocktest.TestRunMockTestWithPipeAndFlags(
			t, pipeData,
			"--api-key", "string",
			"beta:organization:workspaces:members", "add",
			"--workspace-id", "workspace_id",
		)
	})
}

func TestBetaOrganizationWorkspacesMembersRemove(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:organization:workspaces:members", "remove",
			"--workspace-id", "workspace_id",
			"--user-id", "user_id",
		)
	})
}
