package cmd

import (
	"testing"

	"github.com/anthropics/anthropic-cli/v2/internal/mocktest"
	"github.com/anthropics/anthropic-cli/v2/internal/requestflag"
)

func TestUserProfilesCreate(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "string")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"user-profiles", "create",
			"--access-type", "application",
			"--external-user-details", "{account_status: active, country: US, email_hash: x, entity_type: business, name_hash: x, onboarded_at: '2024-11-02T08:15:00Z', reference_id: user_12345}",
			"--metadata", "{}",
			"--name", "Example User",
			"--workspace-id", "wrkspc_011CZkZaBF1tNoB5wlCeusgy",
		)
	})

	t.Run("inner flags", func(t *testing.T) {
		// Check that inner flags have been set up correctly
		requestflag.CheckInnerFlags(userProfilesCreate)

		// Alternative argument passing style using inner flags
		mocktest.TestRunMockTestWithFlags(
			t,
			"user-profiles", "create",
			"--access-type", "application",
			"--external-user-details.account-status", "active",
			"--external-user-details.country", "US",
			"--external-user-details.email-hash", "x",
			"--external-user-details.entity-type", "business",
			"--external-user-details.name-hash", "x",
			"--external-user-details.onboarded-at", "2024-11-02T08:15:00Z",
			"--external-user-details.reference-id", "user_12345",
			"--metadata", "{}",
			"--name", "Example User",
			"--workspace-id", "wrkspc_011CZkZaBF1tNoB5wlCeusgy",
		)
	})

	t.Run("piping data", func(t *testing.T) {
		// Test piping YAML data over stdin
		pipeData := []byte("" +
			"access_type: application\n" +
			"external_user_details:\n" +
			"  account_status: active\n" +
			"  country: US\n" +
			"  email_hash: x\n" +
			"  entity_type: business\n" +
			"  name_hash: x\n" +
			"  onboarded_at: '2024-11-02T08:15:00Z'\n" +
			"  reference_id: user_12345\n" +
			"metadata: {}\n" +
			"name: Example User\n")
		mocktest.TestRunMockTestWithPipeAndFlags(
			t, pipeData,
			"user-profiles", "create",
			"--workspace-id", "wrkspc_011CZkZaBF1tNoB5wlCeusgy",
		)
	})
}

func TestUserProfilesRetrieve(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "string")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"user-profiles", "retrieve",
			"--user-profile-id", "uprof_011CZkZCu8hGbp5mYRQgUmz9",
			"--workspace-id", "wrkspc_011CZkZaBF1tNoB5wlCeusgy",
		)
	})
}

func TestUserProfilesUpdate(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "string")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"user-profiles", "update",
			"--user-profile-id", "uprof_011CZkZCu8hGbp5mYRQgUmz9",
			"--access-type", "application",
			"--external-user-details", "{account_status: suspended, country: SE, email_hash: x, entity_type: individual, name_hash: x, onboarded_at: '2019-12-27T18:11:19.117Z', reference_id: x}",
			"--metadata", "{foo: string}",
			"--name", "x",
			"--workspace-id", "wrkspc_011CZkZaBF1tNoB5wlCeusgy",
		)
	})

	t.Run("inner flags", func(t *testing.T) {
		// Check that inner flags have been set up correctly
		requestflag.CheckInnerFlags(userProfilesUpdate)

		// Alternative argument passing style using inner flags
		mocktest.TestRunMockTestWithFlags(
			t,
			"user-profiles", "update",
			"--user-profile-id", "uprof_011CZkZCu8hGbp5mYRQgUmz9",
			"--access-type", "application",
			"--external-user-details.account-status", "suspended",
			"--external-user-details.country", "SE",
			"--external-user-details.email-hash", "x",
			"--external-user-details.entity-type", "individual",
			"--external-user-details.name-hash", "x",
			"--external-user-details.onboarded-at", "2019-12-27T18:11:19.117Z",
			"--external-user-details.reference-id", "x",
			"--metadata", "{foo: string}",
			"--name", "x",
			"--workspace-id", "wrkspc_011CZkZaBF1tNoB5wlCeusgy",
		)
	})

	t.Run("piping data", func(t *testing.T) {
		// Test piping YAML data over stdin
		pipeData := []byte("" +
			"access_type: application\n" +
			"external_user_details:\n" +
			"  account_status: suspended\n" +
			"  country: SE\n" +
			"  email_hash: x\n" +
			"  entity_type: individual\n" +
			"  name_hash: x\n" +
			"  onboarded_at: '2019-12-27T18:11:19.117Z'\n" +
			"  reference_id: x\n" +
			"metadata:\n" +
			"  foo: string\n" +
			"name: x\n")
		mocktest.TestRunMockTestWithPipeAndFlags(
			t, pipeData,
			"user-profiles", "update",
			"--user-profile-id", "uprof_011CZkZCu8hGbp5mYRQgUmz9",
			"--workspace-id", "wrkspc_011CZkZaBF1tNoB5wlCeusgy",
		)
	})
}

func TestUserProfilesList(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "string")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"user-profiles", "list",
			"--max-items", "10",
			"--enforcement-status", "unrestricted",
			"--limit", "0",
			"--order", "asc",
			"--order-by", "created_at",
			"--page", "page",
			"--workspace-id", "wrkspc_011CZkZaBF1tNoB5wlCeusgy",
		)
	})
}

func TestUserProfilesDelete(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "string")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"user-profiles", "delete",
			"--user-profile-id", "uprof_011CZkZCu8hGbp5mYRQgUmz9",
			"--workspace-id", "wrkspc_011CZkZaBF1tNoB5wlCeusgy",
		)
	})
}
