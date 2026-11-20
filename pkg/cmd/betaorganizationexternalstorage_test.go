package cmd

import (
	"testing"

	"github.com/anthropics/anthropic-cli/v2/internal/mocktest"
)

func TestBetaOrganizationExternalStorageCreate(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "string")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--base-url", mocktest.ServerWithQuery(t, "beta=true"),
			"beta:organization:external-storage", "create",
			"--provider-config", "{account_id: '111122223333', bucket_name: acme-safety-retention, region: us-east-2, type: aws}",
			"--display-name", "x",
			"--user-profile-id", "uprof_01SDCCSbTxrXDpWc1phhtcfK",
			"--beta", "computer-use-2024-10-22",
		)
	})

	t.Run("piping data", func(t *testing.T) {
		// Test piping YAML data over stdin
		pipeData := []byte("" +
			"provider_config:\n" +
			"  account_id: '111122223333'\n" +
			"  bucket_name: acme-safety-retention\n" +
			"  region: us-east-2\n" +
			"  type: aws\n" +
			"display_name: x\n" +
			"user_profile_id: uprof_01SDCCSbTxrXDpWc1phhtcfK\n")
		mocktest.TestRunMockTestWithPipeAndFlags(
			t, pipeData,
			"--base-url", mocktest.ServerWithQuery(t, "beta=true"),
			"beta:organization:external-storage", "create",
			"--beta", "computer-use-2024-10-22",
		)
	})
}

func TestBetaOrganizationExternalStorageRetrieve(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "string")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--base-url", mocktest.ServerWithQuery(t, "beta=true"),
			"beta:organization:external-storage", "retrieve",
			"--external-storage-id", "external_storage_id",
			"--beta", "computer-use-2024-10-22",
		)
	})
}

func TestBetaOrganizationExternalStorageUpdate(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "string")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--base-url", mocktest.ServerWithQuery(t, "beta=true"),
			"beta:organization:external-storage", "update",
			"--external-storage-id", "external_storage_id",
			"--display-name", "x",
			"--beta", "computer-use-2024-10-22",
		)
	})

	t.Run("piping data", func(t *testing.T) {
		// Test piping YAML data over stdin
		pipeData := []byte("display_name: x")
		mocktest.TestRunMockTestWithPipeAndFlags(
			t, pipeData,
			"--base-url", mocktest.ServerWithQuery(t, "beta=true"),
			"beta:organization:external-storage", "update",
			"--external-storage-id", "external_storage_id",
			"--beta", "computer-use-2024-10-22",
		)
	})
}

func TestBetaOrganizationExternalStorageList(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "string")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--base-url", mocktest.ServerWithQuery(t, "beta=true"),
			"beta:organization:external-storage", "list",
			"--max-items", "10",
			"--limit", "1",
			"--page", "page",
			"--user-profile-id", "string",
			"--beta", "computer-use-2024-10-22",
		)
	})
}

func TestBetaOrganizationExternalStorageDelete(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "string")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--base-url", mocktest.ServerWithQuery(t, "beta=true"),
			"beta:organization:external-storage", "delete",
			"--external-storage-id", "external_storage_id",
			"--beta", "computer-use-2024-10-22",
		)
	})
}

func TestBetaOrganizationExternalStorageActivate(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "string")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--base-url", mocktest.ServerWithQuery(t, "beta=true"),
			"beta:organization:external-storage", "activate",
			"--external-storage-id", "external_storage_id",
			"--beta", "computer-use-2024-10-22",
		)
	})
}
