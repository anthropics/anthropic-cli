package cmd

import (
	"testing"

	"github.com/anthropics/anthropic-cli/v2/internal/mocktest"
	"github.com/anthropics/anthropic-cli/v2/internal/requestflag"
)

func TestOrganizationWorkspacesCreate(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "string")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"organization:workspaces", "create",
			"--name", "x",
			"--data-residency", "{allowed_inference_geos: unrestricted, default_inference_geo: global, workspace_geo: us}",
			"--display-color", "#6C5BB9",
			"--external-key-id", "ekey_01SDCCSbTxrXDpWc1phhtcfK",
			"--tags", "{env: prod, team: platform}",
			"--user-profile-id", "uprof_01SDCCSbTxrXDpWc1phhtcfK",
		)
	})

	t.Run("inner flags", func(t *testing.T) {
		// Check that inner flags have been set up correctly
		requestflag.CheckInnerFlags(organizationWorkspacesCreate)

		// Alternative argument passing style using inner flags
		mocktest.TestRunMockTestWithFlags(
			t,
			"organization:workspaces", "create",
			"--name", "x",
			"--data-residency.allowed-inference-geos", "unrestricted",
			"--data-residency.default-inference-geo", "global",
			"--data-residency.workspace-geo", "us",
			"--display-color", "#6C5BB9",
			"--external-key-id", "ekey_01SDCCSbTxrXDpWc1phhtcfK",
			"--tags", "{env: prod, team: platform}",
			"--user-profile-id", "uprof_01SDCCSbTxrXDpWc1phhtcfK",
		)
	})

	t.Run("piping data", func(t *testing.T) {
		// Test piping YAML data over stdin
		pipeData := []byte("" +
			"name: x\n" +
			"data_residency:\n" +
			"  allowed_inference_geos: unrestricted\n" +
			"  default_inference_geo: global\n" +
			"  workspace_geo: us\n" +
			"display_color: '#6C5BB9'\n" +
			"external_key_id: ekey_01SDCCSbTxrXDpWc1phhtcfK\n" +
			"tags:\n" +
			"  env: prod\n" +
			"  team: platform\n" +
			"user_profile_id: uprof_01SDCCSbTxrXDpWc1phhtcfK\n")
		mocktest.TestRunMockTestWithPipeAndFlags(
			t, pipeData,
			"organization:workspaces", "create",
		)
	})
}

func TestOrganizationWorkspacesRetrieve(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "string")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"organization:workspaces", "retrieve",
			"--workspace-id", "workspace_id",
		)
	})
}

func TestOrganizationWorkspacesUpdate(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "string")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"organization:workspaces", "update",
			"--workspace-id", "workspace_id",
			"--data-residency", "{allowed_inference_geos: unrestricted, default_inference_geo: global}",
			"--display-color", "#6C5BB9",
			"--external-key-id", "ekey_01SDCCSbTxrXDpWc1phhtcfK",
			"--name", "x",
			"--tags", "{env: prod, team: platform}",
			"--user-profile-id", "uprof_01SDCCSbTxrXDpWc1phhtcfK",
		)
	})

	t.Run("inner flags", func(t *testing.T) {
		// Check that inner flags have been set up correctly
		requestflag.CheckInnerFlags(organizationWorkspacesUpdate)

		// Alternative argument passing style using inner flags
		mocktest.TestRunMockTestWithFlags(
			t,
			"organization:workspaces", "update",
			"--workspace-id", "workspace_id",
			"--data-residency.allowed-inference-geos", "unrestricted",
			"--data-residency.default-inference-geo", "global",
			"--display-color", "#6C5BB9",
			"--external-key-id", "ekey_01SDCCSbTxrXDpWc1phhtcfK",
			"--name", "x",
			"--tags", "{env: prod, team: platform}",
			"--user-profile-id", "uprof_01SDCCSbTxrXDpWc1phhtcfK",
		)
	})

	t.Run("piping data", func(t *testing.T) {
		// Test piping YAML data over stdin
		pipeData := []byte("" +
			"data_residency:\n" +
			"  allowed_inference_geos: unrestricted\n" +
			"  default_inference_geo: global\n" +
			"display_color: '#6C5BB9'\n" +
			"external_key_id: ekey_01SDCCSbTxrXDpWc1phhtcfK\n" +
			"name: x\n" +
			"tags:\n" +
			"  env: prod\n" +
			"  team: platform\n" +
			"user_profile_id: uprof_01SDCCSbTxrXDpWc1phhtcfK\n")
		mocktest.TestRunMockTestWithPipeAndFlags(
			t, pipeData,
			"organization:workspaces", "update",
			"--workspace-id", "workspace_id",
		)
	})
}

func TestOrganizationWorkspacesList(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "string")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"organization:workspaces", "list",
			"--max-items", "10",
			"--after-id", "after_id",
			"--before-id", "before_id",
			"--id", "string",
			"--include-archived=true",
			"--include-default=true",
			"--limit", "1",
		)
	})
}

func TestOrganizationWorkspacesArchive(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "string")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"organization:workspaces", "archive",
			"--workspace-id", "workspace_id",
		)
	})
}
