// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cmd

import (
	"testing"

	"github.com/anthropics/anthropic-cli/internal/mocktest"
)

func TestBetaOrganizationAPIKeysRetrieve(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:organization:api-keys", "retrieve",
			"--api-key-id", "api_key_id",
		)
	})
}

func TestBetaOrganizationAPIKeysUpdate(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:organization:api-keys", "update",
			"--api-key-id", "api_key_id",
			"--name", "x",
			"--status", "active",
		)
	})

	t.Run("piping data", func(t *testing.T) {
		// Test piping YAML data over stdin
		pipeData := []byte("" +
			"name: x\n" +
			"status: active\n")
		mocktest.TestRunMockTestWithPipeAndFlags(
			t, pipeData,
			"--api-key", "string",
			"beta:organization:api-keys", "update",
			"--api-key-id", "api_key_id",
		)
	})
}

func TestBetaOrganizationAPIKeysList(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:organization:api-keys", "list",
			"--max-items", "10",
			"--after-id", "after_id",
			"--before-id", "before_id",
			"--created-by-user-id", "created_by_user_id",
			"--limit", "1",
			"--status", "active",
			"--workspace-id", "workspace_id",
		)
	})
}
