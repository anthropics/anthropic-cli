package cmd

import (
	"strings"
	"testing"

	"github.com/anthropics/anthropic-cli/v2/internal/mocktest"
)

func TestBetaOrganizationPluginsVersionsCreate(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "string")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--base-url", mocktest.ServerWithQuery(t, "beta=true"),
			"beta:organization:plugins:versions", "create",
			"--plugin-id", "plugin_id",
			"--file", mocktest.TestFile(t, "Example data"),
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
			"release_notes: release_notes\n"
		pipeDataStr = strings.ReplaceAll(pipeDataStr, "Example data", testFile)
		pipeData := []byte(pipeDataStr)
		mocktest.TestRunMockTestWithPipeAndFlags(
			t, pipeData,
			"--base-url", mocktest.ServerWithQuery(t, "beta=true"),
			"beta:organization:plugins:versions", "create",
			"--plugin-id", "plugin_id",
			"--beta", "computer-use-2024-10-22",
		)
	})
}

func TestBetaOrganizationPluginsVersionsRetrieve(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "string")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--base-url", mocktest.ServerWithQuery(t, "beta=true"),
			"beta:organization:plugins:versions", "retrieve",
			"--plugin-id", "plugin_id",
			"--version", "version",
			"--organization-id", "organization_id",
			"--beta", "computer-use-2024-10-22",
		)
	})
}

func TestBetaOrganizationPluginsVersionsList(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "string")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--base-url", mocktest.ServerWithQuery(t, "beta=true"),
			"beta:organization:plugins:versions", "list",
			"--max-items", "10",
			"--plugin-id", "plugin_id",
			"--limit", "1",
			"--organization-id", "organization_id",
			"--page", "page",
			"--beta", "computer-use-2024-10-22",
		)
	})
}

func TestBetaOrganizationPluginsVersionsDownload(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "string")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--base-url", mocktest.ServerWithQuery(t, "beta=true"),
			"beta:organization:plugins:versions", "download",
			"--plugin-id", "plugin_id",
			"--version", "version",
			"--organization-id", "organization_id",
			"--beta", "computer-use-2024-10-22",
			"--output", "/dev/null",
		)
	})
}
