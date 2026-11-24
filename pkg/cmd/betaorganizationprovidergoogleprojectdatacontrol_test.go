package cmd

import (
	"testing"

	"github.com/anthropics/anthropic-cli/v2/internal/mocktest"
)

func TestBetaOrganizationProvidersGoogleProjectsDataControlsRetrieve(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "string")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--base-url", mocktest.ServerWithQuery(t, "beta=true"),
			"beta:organization:providers:google:projects:data-controls", "retrieve",
			"--project-number", "123456789012",
			"--beta", "computer-use-2024-10-22",
		)
	})
}

func TestBetaOrganizationProvidersGoogleProjectsDataControlsUpdate(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "string")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--base-url", mocktest.ServerWithQuery(t, "beta=true"),
			"beta:organization:providers:google:projects:data-controls", "update",
			"--project-number", "123456789012",
			"--cmek", "{type: disabled}",
			"--storage", "{type: disabled}",
			"--beta", "computer-use-2024-10-22",
		)
	})

	t.Run("piping data", func(t *testing.T) {
		// Test piping YAML data over stdin
		pipeData := []byte("" +
			"cmek:\n" +
			"  type: disabled\n" +
			"storage:\n" +
			"  type: disabled\n")
		mocktest.TestRunMockTestWithPipeAndFlags(
			t, pipeData,
			"--base-url", mocktest.ServerWithQuery(t, "beta=true"),
			"beta:organization:providers:google:projects:data-controls", "update",
			"--project-number", "123456789012",
			"--beta", "computer-use-2024-10-22",
		)
	})
}
