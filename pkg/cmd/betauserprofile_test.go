// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cmd

import (
	"testing"

	"github.com/anthropics/anthropic-cli/internal/mocktest"
)

func TestBetaUserProfilesCreate(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:user-profiles", "create",
			"--external-id", "user_12345",
			"--metadata", "{}",
			"--beta", "message-batches-2024-09-24",
		)
	})

	t.Run("piping data", func(t *testing.T) {
		// Test piping YAML data over stdin
		pipeData := []byte("" +
			"external_id: user_12345\n" +
			"metadata: {}\n")
		mocktest.TestRunMockTestWithPipeAndFlags(
			t, pipeData,
			"--api-key", "string",
			"beta:user-profiles", "create",
			"--beta", "message-batches-2024-09-24",
		)
	})
}

func TestBetaUserProfilesRetrieve(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:user-profiles", "retrieve",
			"--user-profile-id", "uprof_011CZkZCu8hGbp5mYRQgUmz9",
			"--beta", "message-batches-2024-09-24",
		)
	})
}

func TestBetaUserProfilesUpdate(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:user-profiles", "update",
			"--user-profile-id", "uprof_011CZkZCu8hGbp5mYRQgUmz9",
			"--external-id", "user_12345",
			"--metadata", "{foo: string}",
			"--beta", "message-batches-2024-09-24",
		)
	})

	t.Run("piping data", func(t *testing.T) {
		// Test piping YAML data over stdin
		pipeData := []byte("" +
			"external_id: user_12345\n" +
			"metadata:\n" +
			"  foo: string\n")
		mocktest.TestRunMockTestWithPipeAndFlags(
			t, pipeData,
			"--api-key", "string",
			"beta:user-profiles", "update",
			"--user-profile-id", "uprof_011CZkZCu8hGbp5mYRQgUmz9",
			"--beta", "message-batches-2024-09-24",
		)
	})
}

func TestBetaUserProfilesList(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:user-profiles", "list",
			"--max-items", "10",
			"--limit", "0",
			"--order", "asc",
			"--page", "page",
			"--beta", "message-batches-2024-09-24",
		)
	})
}

func TestBetaUserProfilesCreateEnrollmentURL(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:user-profiles", "create-enrollment-url",
			"--user-profile-id", "uprof_011CZkZCu8hGbp5mYRQgUmz9",
			"--beta", "message-batches-2024-09-24",
		)
	})
}
