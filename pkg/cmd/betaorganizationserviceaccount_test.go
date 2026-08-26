// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cmd

import (
	"testing"

	"github.com/anthropics/anthropic-cli/internal/mocktest"
)

func TestBetaOrganizationServiceAccountsCreate(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:organization:service-accounts", "create",
			"--name", "ci-deploy-bot",
			"--description", "description",
			"--organization-role", "admin",
			"--beta", "message-batches-2024-09-24",
		)
	})

	t.Run("piping data", func(t *testing.T) {
		// Test piping YAML data over stdin
		pipeData := []byte("" +
			"name: ci-deploy-bot\n" +
			"description: description\n" +
			"organization_role: admin\n")
		mocktest.TestRunMockTestWithPipeAndFlags(
			t, pipeData,
			"--api-key", "string",
			"beta:organization:service-accounts", "create",
			"--beta", "message-batches-2024-09-24",
		)
	})
}

func TestBetaOrganizationServiceAccountsRetrieve(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:organization:service-accounts", "retrieve",
			"--service-account-id", "service_account_id",
			"--beta", "message-batches-2024-09-24",
		)
	})
}

func TestBetaOrganizationServiceAccountsUpdate(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:organization:service-accounts", "update",
			"--service-account-id", "service_account_id",
			"--description", "description",
			"--organization-role", "admin",
			"--beta", "message-batches-2024-09-24",
		)
	})

	t.Run("piping data", func(t *testing.T) {
		// Test piping YAML data over stdin
		pipeData := []byte("" +
			"description: description\n" +
			"organization_role: admin\n")
		mocktest.TestRunMockTestWithPipeAndFlags(
			t, pipeData,
			"--api-key", "string",
			"beta:organization:service-accounts", "update",
			"--service-account-id", "service_account_id",
			"--beta", "message-batches-2024-09-24",
		)
	})
}

func TestBetaOrganizationServiceAccountsList(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:organization:service-accounts", "list",
			"--max-items", "10",
			"--include-archived=true",
			"--limit", "1",
			"--page", "page",
			"--beta", "message-batches-2024-09-24",
		)
	})
}

func TestBetaOrganizationServiceAccountsArchive(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:organization:service-accounts", "archive",
			"--service-account-id", "service_account_id",
			"--beta", "message-batches-2024-09-24",
		)
	})
}
