package cmd

import (
	"testing"

	"github.com/anthropics/anthropic-cli/v2/internal/mocktest"
)

func TestBetaOrganizationWebhooksCreate(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "string")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--base-url", mocktest.ServerWithQuery(t, "beta=true"),
			"beta:organization:webhooks", "create",
			"--enabled-event", "user_profile.enforcement_action_opened",
			"--enabled-event", "user_profile.enforcement_action_resolved",
			"--url", "https://example.com/anthropic-webhooks",
			"--description", "description",
			"--name", "Enforcement actions",
			"--beta", "computer-use-2024-10-22",
		)
	})

	t.Run("piping data", func(t *testing.T) {
		// Test piping YAML data over stdin
		pipeData := []byte("" +
			"enabled_events:\n" +
			"  - user_profile.enforcement_action_opened\n" +
			"  - user_profile.enforcement_action_resolved\n" +
			"url: https://example.com/anthropic-webhooks\n" +
			"description: description\n" +
			"name: Enforcement actions\n")
		mocktest.TestRunMockTestWithPipeAndFlags(
			t, pipeData,
			"--base-url", mocktest.ServerWithQuery(t, "beta=true"),
			"beta:organization:webhooks", "create",
			"--beta", "computer-use-2024-10-22",
		)
	})
}

func TestBetaOrganizationWebhooksRetrieve(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "string")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--base-url", mocktest.ServerWithQuery(t, "beta=true"),
			"beta:organization:webhooks", "retrieve",
			"--webhook-endpoint-id", "webhook_endpoint_id",
			"--beta", "computer-use-2024-10-22",
		)
	})
}

func TestBetaOrganizationWebhooksUpdate(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "string")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--base-url", mocktest.ServerWithQuery(t, "beta=true"),
			"beta:organization:webhooks", "update",
			"--webhook-endpoint-id", "webhook_endpoint_id",
			"--description", "description",
			"--enabled-event", "external_key.failed",
			"--name", "name",
			"--state", "disabled",
			"--url", "x",
			"--beta", "computer-use-2024-10-22",
		)
	})

	t.Run("piping data", func(t *testing.T) {
		// Test piping YAML data over stdin
		pipeData := []byte("" +
			"description: description\n" +
			"enabled_events:\n" +
			"  - external_key.failed\n" +
			"name: name\n" +
			"state: disabled\n" +
			"url: x\n")
		mocktest.TestRunMockTestWithPipeAndFlags(
			t, pipeData,
			"--base-url", mocktest.ServerWithQuery(t, "beta=true"),
			"beta:organization:webhooks", "update",
			"--webhook-endpoint-id", "webhook_endpoint_id",
			"--beta", "computer-use-2024-10-22",
		)
	})
}

func TestBetaOrganizationWebhooksList(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "string")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--base-url", mocktest.ServerWithQuery(t, "beta=true"),
			"beta:organization:webhooks", "list",
			"--max-items", "10",
			"--after-id", "after_id",
			"--before-id", "before_id",
			"--limit", "0",
			"--order", "asc",
			"--beta", "computer-use-2024-10-22",
		)
	})
}

func TestBetaOrganizationWebhooksDelete(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "string")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--base-url", mocktest.ServerWithQuery(t, "beta=true"),
			"beta:organization:webhooks", "delete",
			"--webhook-endpoint-id", "webhook_endpoint_id",
			"--beta", "computer-use-2024-10-22",
		)
	})
}

func TestBetaOrganizationWebhooksRegenerateSigningSecret(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "string")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--base-url", mocktest.ServerWithQuery(t, "beta=true"),
			"beta:organization:webhooks", "regenerate-signing-secret",
			"--webhook-endpoint-id", "webhook_endpoint_id",
			"--beta", "computer-use-2024-10-22",
		)
	})
}
