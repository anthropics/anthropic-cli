package cmd

import (
	"strings"
	"testing"

	"github.com/anthropics/anthropic-cli/v2/internal/mocktest"
)

func TestBetaOrganizationPluginMarketplacesRetrieve(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "string")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--base-url", mocktest.ServerWithQuery(t, "beta=true"),
			"beta:organization:plugin-marketplaces", "retrieve",
			"--marketplace-id", "marketplace_id",
			"--organization-id", "organization_id",
			"--beta", "computer-use-2024-10-22",
		)
	})
}

func TestBetaOrganizationPluginMarketplacesUpdate(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "string")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--base-url", mocktest.ServerWithQuery(t, "beta=true"),
			"beta:organization:plugin-marketplaces", "update",
			"--marketplace-id", "marketplace_id",
			"--default-installation-preference", "available",
			"--beta", "computer-use-2024-10-22",
		)
	})

	t.Run("piping data", func(t *testing.T) {
		// Test piping YAML data over stdin
		pipeData := []byte("default_installation_preference: available")
		mocktest.TestRunMockTestWithPipeAndFlags(
			t, pipeData,
			"--base-url", mocktest.ServerWithQuery(t, "beta=true"),
			"beta:organization:plugin-marketplaces", "update",
			"--marketplace-id", "marketplace_id",
			"--beta", "computer-use-2024-10-22",
		)
	})
}

func TestBetaOrganizationPluginMarketplacesList(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "string")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--base-url", mocktest.ServerWithQuery(t, "beta=true"),
			"beta:organization:plugin-marketplaces", "list",
			"--max-items", "10",
			"--limit", "1",
			"--organization-id", "organization_id",
			"--owner-type", "organization",
			"--page", "page",
			"--source", "directory",
			"--beta", "computer-use-2024-10-22",
		)
	})
}

func TestBetaOrganizationPluginMarketplacesValidateArchive(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "string")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--base-url", mocktest.ServerWithQuery(t, "beta=true"),
			"beta:organization:plugin-marketplaces", "validate-archive",
			"--archive", mocktest.TestFile(t, "Example data"),
			"--beta", "computer-use-2024-10-22",
		)
	})

	t.Run("piping data", func(t *testing.T) {
		testFile := mocktest.TestFile(t, "Example data")
		// Test piping YAML data over stdin
		pipeDataStr := "archive: Example data"
		pipeDataStr = strings.ReplaceAll(pipeDataStr, "Example data", testFile)
		pipeData := []byte(pipeDataStr)
		mocktest.TestRunMockTestWithPipeAndFlags(
			t, pipeData,
			"--base-url", mocktest.ServerWithQuery(t, "beta=true"),
			"beta:organization:plugin-marketplaces", "validate-archive",
			"--beta", "computer-use-2024-10-22",
		)
	})
}

func TestBetaOrganizationPluginMarketplacesValidateRepository(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "string")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--base-url", mocktest.ServerWithQuery(t, "beta=true"),
			"beta:organization:plugin-marketplaces", "validate-repository",
			"--repository-url", "https://github.com/example-org/example-marketplace",
			"--ref", "main",
			"--beta", "computer-use-2024-10-22",
		)
	})

	t.Run("piping data", func(t *testing.T) {
		// Test piping YAML data over stdin
		pipeData := []byte("" +
			"repository_url: https://github.com/example-org/example-marketplace\n" +
			"ref: main\n")
		mocktest.TestRunMockTestWithPipeAndFlags(
			t, pipeData,
			"--base-url", mocktest.ServerWithQuery(t, "beta=true"),
			"beta:organization:plugin-marketplaces", "validate-repository",
			"--beta", "computer-use-2024-10-22",
		)
	})
}
