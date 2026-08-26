// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cmd

import (
	"testing"

	"github.com/anthropics/anthropic-cli/internal/mocktest"
)

func TestBetaOrganizationFederationIssuersCreate(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:organization:federation:issuers", "create",
			"--issuer-url", "x",
			"--name", "x",
			"--check-jti=true",
			"--jwks", "{type: discovery, ca_cert_pem: ca_cert_pem, discovery_base: discovery_base}",
			"--max-jwt-lifetime-seconds", "1",
			"--beta", "message-batches-2024-09-24",
		)
	})

	t.Run("piping data", func(t *testing.T) {
		// Test piping YAML data over stdin
		pipeData := []byte("" +
			"issuer_url: x\n" +
			"name: x\n" +
			"check_jti: true\n" +
			"jwks:\n" +
			"  type: discovery\n" +
			"  ca_cert_pem: ca_cert_pem\n" +
			"  discovery_base: discovery_base\n" +
			"max_jwt_lifetime_seconds: 1\n")
		mocktest.TestRunMockTestWithPipeAndFlags(
			t, pipeData,
			"--api-key", "string",
			"beta:organization:federation:issuers", "create",
			"--beta", "message-batches-2024-09-24",
		)
	})
}

func TestBetaOrganizationFederationIssuersRetrieve(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:organization:federation:issuers", "retrieve",
			"--federation-issuer-id", "federation_issuer_id",
			"--beta", "message-batches-2024-09-24",
		)
	})
}

func TestBetaOrganizationFederationIssuersUpdate(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:organization:federation:issuers", "update",
			"--federation-issuer-id", "federation_issuer_id",
			"--check-jti=true",
			"--issuer-url", "x",
			"--jwks", "{type: discovery, ca_cert_pem: ca_cert_pem, discovery_base: discovery_base}",
			"--jwks-polling-disabled=true",
			"--max-jwt-lifetime-seconds", "1",
			"--name", "x",
			"--beta", "message-batches-2024-09-24",
		)
	})

	t.Run("piping data", func(t *testing.T) {
		// Test piping YAML data over stdin
		pipeData := []byte("" +
			"check_jti: true\n" +
			"issuer_url: x\n" +
			"jwks:\n" +
			"  type: discovery\n" +
			"  ca_cert_pem: ca_cert_pem\n" +
			"  discovery_base: discovery_base\n" +
			"jwks_polling_disabled: true\n" +
			"max_jwt_lifetime_seconds: 1\n" +
			"name: x\n")
		mocktest.TestRunMockTestWithPipeAndFlags(
			t, pipeData,
			"--api-key", "string",
			"beta:organization:federation:issuers", "update",
			"--federation-issuer-id", "federation_issuer_id",
			"--beta", "message-batches-2024-09-24",
		)
	})
}

func TestBetaOrganizationFederationIssuersList(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:organization:federation:issuers", "list",
			"--max-items", "10",
			"--include-archived=true",
			"--limit", "1",
			"--page", "page",
			"--beta", "message-batches-2024-09-24",
		)
	})
}

func TestBetaOrganizationFederationIssuersArchive(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:organization:federation:issuers", "archive",
			"--federation-issuer-id", "federation_issuer_id",
			"--beta", "message-batches-2024-09-24",
		)
	})
}
