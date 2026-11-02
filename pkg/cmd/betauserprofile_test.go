package cmd

import (
	"testing"

	"github.com/anthropics/anthropic-cli/v2/internal/mocktest"
	"github.com/anthropics/anthropic-cli/v2/internal/requestflag"
)

func TestBetaUserProfilesCreate(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "string")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--base-url", mocktest.ServerWithQuery(t, "beta=true"),
			"beta:user-profiles", "create",
			"--access-type", "application",
			"--external-id", "x",
			"--external-user-details", "{account_status: active, country: US, email_hash: x, entity_type: business, name_hash: x, onboarded_at: '2024-11-02T08:15:00Z', reference_id: user_12345}",
			"--external-user-onboarded-at", "'2019-12-27T18:11:19.117Z'",
			"--metadata", "{}",
			"--name", "Example User",
			"--beta", "message-batches-2024-09-24",
			"--workspace-id", "wrkspc_011CZkZaBF1tNoB5wlCeusgy",
		)
	})

	t.Run("inner flags", func(t *testing.T) {
		// Check that inner flags have been set up correctly
		requestflag.CheckInnerFlags(betaUserProfilesCreate)

		// Alternative argument passing style using inner flags
		mocktest.TestRunMockTestWithFlags(
			t,
			"--base-url", mocktest.ServerWithQuery(t, "beta=true"),
			"beta:user-profiles", "create",
			"--access-type", "application",
			"--external-id", "x",
			"--external-user-details.account-status", "active",
			"--external-user-details.country", "US",
			"--external-user-details.email-hash", "x",
			"--external-user-details.entity-type", "business",
			"--external-user-details.name-hash", "x",
			"--external-user-details.onboarded-at", "2024-11-02T08:15:00Z",
			"--external-user-details.reference-id", "user_12345",
			"--external-user-onboarded-at", "'2019-12-27T18:11:19.117Z'",
			"--metadata", "{}",
			"--name", "Example User",
			"--beta", "message-batches-2024-09-24",
			"--workspace-id", "wrkspc_011CZkZaBF1tNoB5wlCeusgy",
		)
	})

	t.Run("piping data", func(t *testing.T) {
		// Test piping YAML data over stdin
		pipeData := []byte("" +
			"access_type: application\n" +
			"external_id: x\n" +
			"external_user_details:\n" +
			"  account_status: active\n" +
			"  country: US\n" +
			"  email_hash: x\n" +
			"  entity_type: business\n" +
			"  name_hash: x\n" +
			"  onboarded_at: '2024-11-02T08:15:00Z'\n" +
			"  reference_id: user_12345\n" +
			"external_user_onboarded_at: '2019-12-27T18:11:19.117Z'\n" +
			"metadata: {}\n" +
			"name: Example User\n")
		mocktest.TestRunMockTestWithPipeAndFlags(
			t, pipeData,
			"--base-url", mocktest.ServerWithQuery(t, "beta=true"),
			"beta:user-profiles", "create",
			"--beta", "message-batches-2024-09-24",
			"--workspace-id", "wrkspc_011CZkZaBF1tNoB5wlCeusgy",
		)
	})
}

func TestBetaUserProfilesRetrieve(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "string")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--base-url", mocktest.ServerWithQuery(t, "beta=true"),
			"beta:user-profiles", "retrieve",
			"--user-profile-id", "uprof_011CZkZCu8hGbp5mYRQgUmz9",
			"--beta", "message-batches-2024-09-24",
			"--workspace-id", "wrkspc_011CZkZaBF1tNoB5wlCeusgy",
		)
	})
}

func TestBetaUserProfilesUpdate(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "string")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--base-url", mocktest.ServerWithQuery(t, "beta=true"),
			"beta:user-profiles", "update",
			"--user-profile-id", "uprof_011CZkZCu8hGbp5mYRQgUmz9",
			"--access-type", "application",
			"--external-id", "x",
			"--external-user-details", "{account_status: suspended, country: SE, email_hash: x, entity_type: individual, name_hash: x, onboarded_at: '2019-12-27T18:11:19.117Z', reference_id: x}",
			"--external-user-onboarded-at", "'2019-12-27T18:11:19.117Z'",
			"--metadata", "{foo: string}",
			"--name", "x",
			"--beta", "message-batches-2024-09-24",
			"--workspace-id", "wrkspc_011CZkZaBF1tNoB5wlCeusgy",
		)
	})

	t.Run("inner flags", func(t *testing.T) {
		// Check that inner flags have been set up correctly
		requestflag.CheckInnerFlags(betaUserProfilesUpdate)

		// Alternative argument passing style using inner flags
		mocktest.TestRunMockTestWithFlags(
			t,
			"--base-url", mocktest.ServerWithQuery(t, "beta=true"),
			"beta:user-profiles", "update",
			"--user-profile-id", "uprof_011CZkZCu8hGbp5mYRQgUmz9",
			"--access-type", "application",
			"--external-id", "x",
			"--external-user-details.account-status", "suspended",
			"--external-user-details.country", "SE",
			"--external-user-details.email-hash", "x",
			"--external-user-details.entity-type", "individual",
			"--external-user-details.name-hash", "x",
			"--external-user-details.onboarded-at", "2019-12-27T18:11:19.117Z",
			"--external-user-details.reference-id", "x",
			"--external-user-onboarded-at", "'2019-12-27T18:11:19.117Z'",
			"--metadata", "{foo: string}",
			"--name", "x",
			"--beta", "message-batches-2024-09-24",
			"--workspace-id", "wrkspc_011CZkZaBF1tNoB5wlCeusgy",
		)
	})

	t.Run("piping data", func(t *testing.T) {
		// Test piping YAML data over stdin
		pipeData := []byte("" +
			"access_type: application\n" +
			"external_id: x\n" +
			"external_user_details:\n" +
			"  account_status: suspended\n" +
			"  country: SE\n" +
			"  email_hash: x\n" +
			"  entity_type: individual\n" +
			"  name_hash: x\n" +
			"  onboarded_at: '2019-12-27T18:11:19.117Z'\n" +
			"  reference_id: x\n" +
			"external_user_onboarded_at: '2019-12-27T18:11:19.117Z'\n" +
			"metadata:\n" +
			"  foo: string\n" +
			"name: x\n")
		mocktest.TestRunMockTestWithPipeAndFlags(
			t, pipeData,
			"--base-url", mocktest.ServerWithQuery(t, "beta=true"),
			"beta:user-profiles", "update",
			"--user-profile-id", "uprof_011CZkZCu8hGbp5mYRQgUmz9",
			"--beta", "message-batches-2024-09-24",
			"--workspace-id", "wrkspc_011CZkZaBF1tNoB5wlCeusgy",
		)
	})
}

func TestBetaUserProfilesList(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "string")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--base-url", mocktest.ServerWithQuery(t, "beta=true"),
			"beta:user-profiles", "list",
			"--max-items", "10",
			"--limit", "0",
			"--order", "asc",
			"--order-by", "created_at",
			"--page", "page",
			"--beta", "message-batches-2024-09-24",
			"--workspace-id", "wrkspc_011CZkZaBF1tNoB5wlCeusgy",
		)
	})
}

func TestBetaUserProfilesDelete(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "string")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--base-url", mocktest.ServerWithQuery(t, "beta=true"),
			"beta:user-profiles", "delete",
			"--user-profile-id", "uprof_011CZkZCu8hGbp5mYRQgUmz9",
			"--beta", "message-batches-2024-09-24",
			"--workspace-id", "wrkspc_011CZkZaBF1tNoB5wlCeusgy",
		)
	})
}

func TestBetaUserProfilesCreateEnrollmentURL(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "string")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--base-url", mocktest.ServerWithQuery(t, "beta=true"),
			"beta:user-profiles", "create-enrollment-url",
			"--user-profile-id", "uprof_011CZkZCu8hGbp5mYRQgUmz9",
			"--beta", "message-batches-2024-09-24",
			"--workspace-id", "wrkspc_011CZkZaBF1tNoB5wlCeusgy",
		)
	})
}
