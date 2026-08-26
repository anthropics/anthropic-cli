// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cmd

import (
	"testing"

	"github.com/anthropics/anthropic-cli/internal/mocktest"
	"github.com/anthropics/anthropic-cli/internal/requestflag"
)

func TestBetaOrganizationWorkspacesCreate(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:organization:workspaces", "create",
			"--name", "x",
			"--data-residency", "{allowed_inference_geos: unrestricted, default_inference_geo: global, workspace_geo: us}",
			"--display-color", "#6C5BB9",
			"--external-key-id", "ekey_01SDCCSbTxrXDpWc1phhtcfK",
			"--tags", "{env: prod, team: platform}",
			"--beta", "message-batches-2024-09-24",
		)
	})

	t.Run("inner flags", func(t *testing.T) {
		// Check that inner flags have been set up correctly
		requestflag.CheckInnerFlags(betaOrganizationWorkspacesCreate)

		// Alternative argument passing style using inner flags
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:organization:workspaces", "create",
			"--name", "x",
			"--data-residency.allowed-inference-geos", "unrestricted",
			"--data-residency.default-inference-geo", "global",
			"--data-residency.workspace-geo", "us",
			"--display-color", "#6C5BB9",
			"--external-key-id", "ekey_01SDCCSbTxrXDpWc1phhtcfK",
			"--tags", "{env: prod, team: platform}",
			"--beta", "message-batches-2024-09-24",
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
			"  team: platform\n")
		mocktest.TestRunMockTestWithPipeAndFlags(
			t, pipeData,
			"--api-key", "string",
			"beta:organization:workspaces", "create",
			"--beta", "message-batches-2024-09-24",
		)
	})
}

func TestBetaOrganizationWorkspacesRetrieve(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:organization:workspaces", "retrieve",
			"--workspace-id", "workspace_id",
		)
	})
}

func TestBetaOrganizationWorkspacesUpdate(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:organization:workspaces", "update",
			"--workspace-id", "workspace_id",
			"--data-residency", "{allowed_inference_geos: unrestricted, default_inference_geo: global}",
			"--display-color", "#6C5BB9",
			"--external-key-id", "ekey_01SDCCSbTxrXDpWc1phhtcfK",
			"--name", "x",
			"--tags", "{env: prod, team: platform}",
		)
	})

	t.Run("inner flags", func(t *testing.T) {
		// Check that inner flags have been set up correctly
		requestflag.CheckInnerFlags(betaOrganizationWorkspacesUpdate)

		// Alternative argument passing style using inner flags
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:organization:workspaces", "update",
			"--workspace-id", "workspace_id",
			"--data-residency.allowed-inference-geos", "unrestricted",
			"--data-residency.default-inference-geo", "global",
			"--display-color", "#6C5BB9",
			"--external-key-id", "ekey_01SDCCSbTxrXDpWc1phhtcfK",
			"--name", "x",
			"--tags", "{env: prod, team: platform}",
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
			"  team: platform\n")
		mocktest.TestRunMockTestWithPipeAndFlags(
			t, pipeData,
			"--api-key", "string",
			"beta:organization:workspaces", "update",
			"--workspace-id", "workspace_id",
		)
	})
}

func TestBetaOrganizationWorkspacesList(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:organization:workspaces", "list",
			"--max-items", "10",
			"--after-id", "after_id",
			"--before-id", "before_id",
			"--include-archived=true",
			"--limit", "1",
		)
	})
}

func TestBetaOrganizationWorkspacesArchive(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:organization:workspaces", "archive",
			"--workspace-id", "workspace_id",
		)
	})
}
