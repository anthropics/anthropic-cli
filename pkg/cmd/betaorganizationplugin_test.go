package cmd

import (
	"strings"
	"testing"

	"github.com/anthropics/anthropic-cli/v2/internal/mocktest"
)

func TestBetaOrganizationPluginsCreate(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "string")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--base-url", mocktest.ServerWithQuery(t, "beta=true"),
			"beta:organization:plugins", "create",
			"--file", mocktest.TestFile(t, "Example data"),
			"--marketplace-id", "marketplace_id",
			"--release-notes", "release_notes",
			"--beta", "computer-use-2024-10-22",
		)
	})

	t.Run("piping data", func(t *testing.T) {
		testFile := mocktest.TestFile(t, "Example data")
		// Test piping YAML data over stdin
		pipeDataStr := "" +
			"files:\n" +
			"  - Example data\n" +
			"marketplace_id: marketplace_id\n" +
			"release_notes: release_notes\n"
		pipeDataStr = strings.ReplaceAll(pipeDataStr, "Example data", testFile)
		pipeData := []byte(pipeDataStr)
		mocktest.TestRunMockTestWithPipeAndFlags(
			t, pipeData,
			"--base-url", mocktest.ServerWithQuery(t, "beta=true"),
			"beta:organization:plugins", "create",
			"--beta", "computer-use-2024-10-22",
		)
	})
}

func TestBetaOrganizationPluginsRetrieve(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "string")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--base-url", mocktest.ServerWithQuery(t, "beta=true"),
			"beta:organization:plugins", "retrieve",
			"--plugin-id", "plugin_id",
			"--organization-id", "organization_id",
			"--beta", "computer-use-2024-10-22",
		)
	})
}

func TestBetaOrganizationPluginsUpdate(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "string")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--base-url", mocktest.ServerWithQuery(t, "beta=true"),
			"beta:organization:plugins", "update",
			"--plugin-id", "plugin_id",
			"--served-version-id", "pluginver_01KaZmQpRsTuVwXyZ2b4c6d8",
			"--beta", "computer-use-2024-10-22",
		)
	})

	t.Run("piping data", func(t *testing.T) {
		// Test piping YAML data over stdin
		pipeData := []byte("served_version_id: pluginver_01KaZmQpRsTuVwXyZ2b4c6d8")
		mocktest.TestRunMockTestWithPipeAndFlags(
			t, pipeData,
			"--base-url", mocktest.ServerWithQuery(t, "beta=true"),
			"beta:organization:plugins", "update",
			"--plugin-id", "plugin_id",
			"--beta", "computer-use-2024-10-22",
		)
	})
}

func TestBetaOrganizationPluginsList(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "string")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--base-url", mocktest.ServerWithQuery(t, "beta=true"),
			"beta:organization:plugins", "list",
			"--max-items", "10",
			"--created-at-gt", "'2019-12-27T18:11:19.117Z'",
			"--created-at-gte", "'2019-12-27T18:11:19.117Z'",
			"--created-at-lt", "'2019-12-27T18:11:19.117Z'",
			"--created-at-lte", "'2019-12-27T18:11:19.117Z'",
			"--limit", "1",
			"--marketplace-id", "marketplace_id",
			"--organization-id", "organization_id",
			"--owner-type", "organization",
			"--owner-user-id", "owner_user_id",
			"--page", "page",
			"--beta", "computer-use-2024-10-22",
		)
	})
}

func TestBetaOrganizationPluginsDelete(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "string")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--base-url", mocktest.ServerWithQuery(t, "beta=true"),
			"beta:organization:plugins", "delete",
			"--plugin-id", "plugin_id",
			"--beta", "computer-use-2024-10-22",
		)
	})
}
