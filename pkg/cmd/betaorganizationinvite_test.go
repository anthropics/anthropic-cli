// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cmd

import (
	"testing"

	"github.com/anthropics/anthropic-cli/internal/mocktest"
)

func TestBetaOrganizationInvitesCreate(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:organization:invites", "create",
			"--email", "user@emaildomain.com",
			"--role", "user",
			"--rbac-group-id", "string",
		)
	})

	t.Run("piping data", func(t *testing.T) {
		// Test piping YAML data over stdin
		pipeData := []byte("" +
			"email: user@emaildomain.com\n" +
			"role: user\n" +
			"rbac_group_ids:\n" +
			"  - string\n")
		mocktest.TestRunMockTestWithPipeAndFlags(
			t, pipeData,
			"--api-key", "string",
			"beta:organization:invites", "create",
		)
	})
}

func TestBetaOrganizationInvitesRetrieve(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:organization:invites", "retrieve",
			"--invite-id", "invite_id",
		)
	})
}

func TestBetaOrganizationInvitesList(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:organization:invites", "list",
			"--max-items", "10",
			"--after-id", "after_id",
			"--before-id", "before_id",
			"--email", "dev@stainless.com",
			"--limit", "1",
			"--role", "string",
			"--status", "accepted",
		)
	})
}

func TestBetaOrganizationInvitesDelete(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:organization:invites", "delete",
			"--invite-id", "invite_id",
		)
	})
}
