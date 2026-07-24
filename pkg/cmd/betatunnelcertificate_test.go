// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cmd

import (
	"testing"

	"github.com/anthropics/anthropic-cli/internal/mocktest"
)

func TestBetaTunnelsCertificatesCreate(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:tunnels:certificates", "create",
			"--tunnel-id", "tunnel_id",
			"--ca-certificate-pem", "ca_certificate_pem",
			"--beta", "message-batches-2024-09-24",
		)
	})

	t.Run("piping data", func(t *testing.T) {
		// Test piping YAML data over stdin
		pipeData := []byte("ca_certificate_pem: ca_certificate_pem")
		mocktest.TestRunMockTestWithPipeAndFlags(
			t, pipeData,
			"--api-key", "string",
			"beta:tunnels:certificates", "create",
			"--tunnel-id", "tunnel_id",
			"--beta", "message-batches-2024-09-24",
		)
	})
}

func TestBetaTunnelsCertificatesRetrieve(t *testing.T) {
	t.Skip("buildURL drops path-level query params (SDK-4349)")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:tunnels:certificates", "retrieve",
			"--tunnel-id", "tunnel_id",
			"--certificate-id", "certificate_id",
			"--beta", "message-batches-2024-09-24",
		)
	})
}

func TestBetaTunnelsCertificatesList(t *testing.T) {
	t.Skip("buildURL drops path-level query params (SDK-4349)")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:tunnels:certificates", "list",
			"--max-items", "10",
			"--tunnel-id", "tunnel_id",
			"--include-archived=true",
			"--limit", "0",
			"--page", "page",
			"--beta", "message-batches-2024-09-24",
		)
	})
}

func TestBetaTunnelsCertificatesArchive(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:tunnels:certificates", "archive",
			"--tunnel-id", "tunnel_id",
			"--certificate-id", "certificate_id",
			"--beta", "message-batches-2024-09-24",
		)
	})
}
