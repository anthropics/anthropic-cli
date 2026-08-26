// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cmd

import (
	"testing"

	"github.com/anthropics/anthropic-cli/internal/mocktest"
)

func TestBetaOrganizationExternalKeysCreate(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:organization:external-keys", "create",
			"--provider-config", "{kms_arn: arn:aws:kms:us-east-1:111122223333:key/abcd1234-5678-90ab-cdef-000011112222, type: aws, region: us-east-1, role_arn: arn:aws:iam::111122223333:role/anthropic-cmek}",
			"--display-name", "x",
			"--geo", "us",
		)
	})

	t.Run("piping data", func(t *testing.T) {
		// Test piping YAML data over stdin
		pipeData := []byte("" +
			"provider_config:\n" +
			"  kms_arn: arn:aws:kms:us-east-1:111122223333:key/abcd1234-5678-90ab-cdef-000011112222\n" +
			"  type: aws\n" +
			"  region: us-east-1\n" +
			"  role_arn: arn:aws:iam::111122223333:role/anthropic-cmek\n" +
			"display_name: x\n" +
			"geo: us\n")
		mocktest.TestRunMockTestWithPipeAndFlags(
			t, pipeData,
			"--api-key", "string",
			"beta:organization:external-keys", "create",
		)
	})
}

func TestBetaOrganizationExternalKeysRetrieve(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:organization:external-keys", "retrieve",
			"--external-key-id", "external_key_id",
		)
	})
}

func TestBetaOrganizationExternalKeysUpdate(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:organization:external-keys", "update",
			"--external-key-id", "external_key_id",
			"--display-name", "x",
			"--geo", "us",
			"--provider-config", "{kms_arn: arn:aws:kms:us-east-1:111122223333:key/abcd1234-5678-90ab-cdef-000011112222, type: aws, region: us-east-1, role_arn: arn:aws:iam::111122223333:role/anthropic-cmek}",
		)
	})

	t.Run("piping data", func(t *testing.T) {
		// Test piping YAML data over stdin
		pipeData := []byte("" +
			"display_name: x\n" +
			"geo: us\n" +
			"provider_config:\n" +
			"  kms_arn: arn:aws:kms:us-east-1:111122223333:key/abcd1234-5678-90ab-cdef-000011112222\n" +
			"  type: aws\n" +
			"  region: us-east-1\n" +
			"  role_arn: arn:aws:iam::111122223333:role/anthropic-cmek\n")
		mocktest.TestRunMockTestWithPipeAndFlags(
			t, pipeData,
			"--api-key", "string",
			"beta:organization:external-keys", "update",
			"--external-key-id", "external_key_id",
		)
	})
}

func TestBetaOrganizationExternalKeysList(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:organization:external-keys", "list",
			"--max-items", "10",
			"--limit", "1",
			"--page", "page",
		)
	})
}

func TestBetaOrganizationExternalKeysDelete(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:organization:external-keys", "delete",
			"--external-key-id", "external_key_id",
		)
	})
}

func TestBetaOrganizationExternalKeysValidate(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:organization:external-keys", "validate",
			"--external-key-id", "external_key_id",
		)
	})
}
