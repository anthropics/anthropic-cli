package cmd

import (
	"testing"

	"github.com/anthropics/anthropic-cli/v2/internal/mocktest"
	"github.com/anthropics/anthropic-cli/v2/internal/requestflag"
)

func TestBetaOrganizationClaudeTagSlackWorkspacesRetrieve(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "string")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--base-url", mocktest.ServerWithQuery(t, "beta=true"),
			"beta:organization:claude-tag:slack-workspaces", "retrieve",
			"--id", "slws_01Hq8TzW5mKc2YbX7vNd4RfJ",
			"--include-inherited=true",
			"--beta", "message-batches-2024-09-24",
		)
	})
}

func TestBetaOrganizationClaudeTagSlackWorkspacesUpdate(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "string")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--base-url", mocktest.ServerWithQuery(t, "beta=true"),
			"beta:organization:claude-tag:slack-workspaces", "update",
			"--id", "slws_01Hq8TzW5mKc2YbX7vNd4RfJ",
			"--settings", "{enabled: true, environment_id: x, guest_policy: {type: allow}, instructions: instructions, model: x, reply_mode: {type: automatic}}",
			"--beta", "message-batches-2024-09-24",
		)
	})

	t.Run("inner flags", func(t *testing.T) {
		// Check that inner flags have been set up correctly
		requestflag.CheckInnerFlags(betaOrganizationClaudeTagSlackWorkspacesUpdate)

		// Alternative argument passing style using inner flags
		mocktest.TestRunMockTestWithFlags(
			t,
			"--base-url", mocktest.ServerWithQuery(t, "beta=true"),
			"beta:organization:claude-tag:slack-workspaces", "update",
			"--id", "slws_01Hq8TzW5mKc2YbX7vNd4RfJ",
			"--settings.enabled=true",
			"--settings.environment-id", "x",
			"--settings.guest-policy", "{type: allow}",
			"--settings.instructions", "instructions",
			"--settings.model", "x",
			"--settings.reply-mode", "{type: automatic}",
			"--beta", "message-batches-2024-09-24",
		)
	})

	t.Run("piping data", func(t *testing.T) {
		// Test piping YAML data over stdin
		pipeData := []byte("" +
			"settings:\n" +
			"  enabled: true\n" +
			"  environment_id: x\n" +
			"  guest_policy:\n" +
			"    type: allow\n" +
			"  instructions: instructions\n" +
			"  model: x\n" +
			"  reply_mode:\n" +
			"    type: automatic\n")
		mocktest.TestRunMockTestWithPipeAndFlags(
			t, pipeData,
			"--base-url", mocktest.ServerWithQuery(t, "beta=true"),
			"beta:organization:claude-tag:slack-workspaces", "update",
			"--id", "slws_01Hq8TzW5mKc2YbX7vNd4RfJ",
			"--beta", "message-batches-2024-09-24",
		)
	})
}

func TestBetaOrganizationClaudeTagSlackWorkspacesList(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "string")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--base-url", mocktest.ServerWithQuery(t, "beta=true"),
			"beta:organization:claude-tag:slack-workspaces", "list",
			"--max-items", "10",
			"--include-inherited=true",
			"--limit", "0",
			"--page", "page",
			"--slack-workspace-id", "slack_workspace_id",
			"--beta", "message-batches-2024-09-24",
		)
	})
}
