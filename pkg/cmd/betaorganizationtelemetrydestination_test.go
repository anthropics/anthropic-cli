package cmd

import (
	"testing"

	"github.com/anthropics/anthropic-cli/v2/internal/mocktest"
	"github.com/anthropics/anthropic-cli/v2/internal/requestflag"
)

func TestBetaOrganizationTelemetryDestinationsCreate(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "string")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--base-url", mocktest.ServerWithQuery(t, "beta=true"),
			"beta:organization:telemetry-destinations", "create",
			"--endpoint-url", "https://gateway.example.com",
			"--workspaces", "{type: all}",
			"--display-name", "display_name",
			"--headers", "{foo: string}",
			"--traces", "{content_mode: metadata_only, enabled: true, endpoint_url: endpoint_url, max_content_bytes: 0}",
			"--beta", "computer-use-2024-10-22",
		)
	})

	t.Run("inner flags", func(t *testing.T) {
		// Check that inner flags have been set up correctly
		requestflag.CheckInnerFlags(betaOrganizationTelemetryDestinationsCreate)

		// Alternative argument passing style using inner flags
		mocktest.TestRunMockTestWithFlags(
			t,
			"--base-url", mocktest.ServerWithQuery(t, "beta=true"),
			"beta:organization:telemetry-destinations", "create",
			"--endpoint-url", "https://gateway.example.com",
			"--workspaces", "{type: all}",
			"--display-name", "display_name",
			"--headers", "{foo: string}",
			"--traces.content-mode", "metadata_only",
			"--traces.enabled=true",
			"--traces.endpoint-url", "endpoint_url",
			"--traces.max-content-bytes", "0",
			"--beta", "computer-use-2024-10-22",
		)
	})

	t.Run("piping data", func(t *testing.T) {
		// Test piping YAML data over stdin
		pipeData := []byte("" +
			"endpoint_url: https://gateway.example.com\n" +
			"workspaces:\n" +
			"  type: all\n" +
			"display_name: display_name\n" +
			"headers:\n" +
			"  foo: string\n" +
			"traces:\n" +
			"  content_mode: metadata_only\n" +
			"  enabled: true\n" +
			"  endpoint_url: endpoint_url\n" +
			"  max_content_bytes: 0\n")
		mocktest.TestRunMockTestWithPipeAndFlags(
			t, pipeData,
			"--base-url", mocktest.ServerWithQuery(t, "beta=true"),
			"beta:organization:telemetry-destinations", "create",
			"--beta", "computer-use-2024-10-22",
		)
	})
}

func TestBetaOrganizationTelemetryDestinationsRetrieve(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "string")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--base-url", mocktest.ServerWithQuery(t, "beta=true"),
			"beta:organization:telemetry-destinations", "retrieve",
			"--destination-id", "destination_id",
			"--beta", "computer-use-2024-10-22",
		)
	})
}

func TestBetaOrganizationTelemetryDestinationsUpdate(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "string")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--base-url", mocktest.ServerWithQuery(t, "beta=true"),
			"beta:organization:telemetry-destinations", "update",
			"--destination-id", "destination_id",
			"--display-name", "display_name",
			"--endpoint-url", "endpoint_url",
			"--headers", "{foo: string}",
			"--state", "{type: enabled}",
			"--traces", "{content_mode: metadata_only, enabled: true, endpoint_url: endpoint_url, max_content_bytes: 0}",
			"--workspaces", "{type: all}",
			"--beta", "computer-use-2024-10-22",
		)
	})

	t.Run("inner flags", func(t *testing.T) {
		// Check that inner flags have been set up correctly
		requestflag.CheckInnerFlags(betaOrganizationTelemetryDestinationsUpdate)

		// Alternative argument passing style using inner flags
		mocktest.TestRunMockTestWithFlags(
			t,
			"--base-url", mocktest.ServerWithQuery(t, "beta=true"),
			"beta:organization:telemetry-destinations", "update",
			"--destination-id", "destination_id",
			"--display-name", "display_name",
			"--endpoint-url", "endpoint_url",
			"--headers", "{foo: string}",
			"--state", "{type: enabled}",
			"--traces.content-mode", "metadata_only",
			"--traces.enabled=true",
			"--traces.endpoint-url", "endpoint_url",
			"--traces.max-content-bytes", "0",
			"--workspaces", "{type: all}",
			"--beta", "computer-use-2024-10-22",
		)
	})

	t.Run("piping data", func(t *testing.T) {
		// Test piping YAML data over stdin
		pipeData := []byte("" +
			"display_name: display_name\n" +
			"endpoint_url: endpoint_url\n" +
			"headers:\n" +
			"  foo: string\n" +
			"state:\n" +
			"  type: enabled\n" +
			"traces:\n" +
			"  content_mode: metadata_only\n" +
			"  enabled: true\n" +
			"  endpoint_url: endpoint_url\n" +
			"  max_content_bytes: 0\n" +
			"workspaces:\n" +
			"  type: all\n")
		mocktest.TestRunMockTestWithPipeAndFlags(
			t, pipeData,
			"--base-url", mocktest.ServerWithQuery(t, "beta=true"),
			"beta:organization:telemetry-destinations", "update",
			"--destination-id", "destination_id",
			"--beta", "computer-use-2024-10-22",
		)
	})
}

func TestBetaOrganizationTelemetryDestinationsList(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "string")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--base-url", mocktest.ServerWithQuery(t, "beta=true"),
			"beta:organization:telemetry-destinations", "list",
			"--max-items", "10",
			"--limit", "0",
			"--page", "page",
			"--workspace-id", "string",
			"--beta", "computer-use-2024-10-22",
		)
	})
}

func TestBetaOrganizationTelemetryDestinationsDelete(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "string")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--base-url", mocktest.ServerWithQuery(t, "beta=true"),
			"beta:organization:telemetry-destinations", "delete",
			"--destination-id", "destination_id",
			"--beta", "computer-use-2024-10-22",
		)
	})
}

func TestBetaOrganizationTelemetryDestinationsTest(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "string")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--base-url", mocktest.ServerWithQuery(t, "beta=true"),
			"beta:organization:telemetry-destinations", "test",
			"--destination-id", "destination_id",
			"--beta", "computer-use-2024-10-22",
		)
	})
}
