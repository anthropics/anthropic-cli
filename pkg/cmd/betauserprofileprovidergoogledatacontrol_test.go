package cmd

import (
	"testing"

	"github.com/anthropics/anthropic-cli/v2/internal/mocktest"
)

func TestBetaUserProfilesProvidersGoogleDataControlsRetrieve(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "string")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--base-url", mocktest.ServerWithQuery(t, "beta=true"),
			"beta:user-profiles:providers:google:data-controls", "retrieve",
			"--user-profile-id", "uprof_011CZkZCu8hGbp5mYRQgUmz9",
			"--beta", "computer-use-2024-10-22",
			"--workspace-id", "wrkspc_011CZkZaBF1tNoB5wlCeusgy",
		)
	})
}

func TestBetaUserProfilesProvidersGoogleDataControlsUpdate(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "string")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--base-url", mocktest.ServerWithQuery(t, "beta=true"),
			"beta:user-profiles:providers:google:data-controls", "update",
			"--user-profile-id", "uprof_011CZkZCu8hGbp5mYRQgUmz9",
			"--cmek", "{type: enabled, external_key_ids: {us: ekey_01SDCCSbTxrXDpWc1phhtcfK}}",
			"--storage", "{type: disabled}",
			"--beta", "computer-use-2024-10-22",
			"--workspace-id", "wrkspc_011CZkZaBF1tNoB5wlCeusgy",
		)
	})

	t.Run("piping data", func(t *testing.T) {
		// Test piping YAML data over stdin
		pipeData := []byte("" +
			"cmek:\n" +
			"  type: enabled\n" +
			"  external_key_ids:\n" +
			"    us: ekey_01SDCCSbTxrXDpWc1phhtcfK\n" +
			"storage:\n" +
			"  type: disabled\n")
		mocktest.TestRunMockTestWithPipeAndFlags(
			t, pipeData,
			"--base-url", mocktest.ServerWithQuery(t, "beta=true"),
			"beta:user-profiles:providers:google:data-controls", "update",
			"--user-profile-id", "uprof_011CZkZCu8hGbp5mYRQgUmz9",
			"--beta", "computer-use-2024-10-22",
			"--workspace-id", "wrkspc_011CZkZaBF1tNoB5wlCeusgy",
		)
	})
}
