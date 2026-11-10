package cmd

import (
	"testing"

	"github.com/anthropics/anthropic-cli/v2/internal/mocktest"
	"github.com/anthropics/anthropic-cli/v2/internal/requestflag"
)

func TestBetaOrganizationClaudeTagSlackChannelsCreate(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "string")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--base-url", mocktest.ServerWithQuery(t, "beta=true"),
			"beta:organization:claude-tag:slack-channels", "create",
			"--slack", "{channel_id: xx, workspace_id: xx}",
			"--settings", "{enabled: true, environment_id: x, guest_policy: {type: allow}, instructions: instructions, model: x, reply_mode: {type: automatic}}",
			"--beta", "message-batches-2024-09-24",
		)
	})

	t.Run("inner flags", func(t *testing.T) {
		// Check that inner flags have been set up correctly
		requestflag.CheckInnerFlags(betaOrganizationClaudeTagSlackChannelsCreate)

		// Alternative argument passing style using inner flags
		mocktest.TestRunMockTestWithFlags(
			t,
			"--base-url", mocktest.ServerWithQuery(t, "beta=true"),
			"beta:organization:claude-tag:slack-channels", "create",
			"--slack.channel-id", "xx",
			"--slack.workspace-id", "xx",
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
			"slack:\n" +
			"  channel_id: xx\n" +
			"  workspace_id: xx\n" +
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
			"beta:organization:claude-tag:slack-channels", "create",
			"--beta", "message-batches-2024-09-24",
		)
	})
}

func TestBetaOrganizationClaudeTagSlackChannelsRetrieve(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "string")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--base-url", mocktest.ServerWithQuery(t, "beta=true"),
			"beta:organization:claude-tag:slack-channels", "retrieve",
			"--id", "slch_01SvpyRQK3nN6eH9uhRA7n1q",
			"--include-inherited=true",
			"--beta", "message-batches-2024-09-24",
		)
	})
}

func TestBetaOrganizationClaudeTagSlackChannelsUpdate(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "string")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--base-url", mocktest.ServerWithQuery(t, "beta=true"),
			"beta:organization:claude-tag:slack-channels", "update",
			"--id", "slch_01SvpyRQK3nN6eH9uhRA7n1q",
			"--settings", "{enabled: true, environment_id: x, guest_policy: {type: allow}, instructions: instructions, model: x, reply_mode: {type: automatic}}",
			"--beta", "message-batches-2024-09-24",
		)
	})

	t.Run("inner flags", func(t *testing.T) {
		// Check that inner flags have been set up correctly
		requestflag.CheckInnerFlags(betaOrganizationClaudeTagSlackChannelsUpdate)

		// Alternative argument passing style using inner flags
		mocktest.TestRunMockTestWithFlags(
			t,
			"--base-url", mocktest.ServerWithQuery(t, "beta=true"),
			"beta:organization:claude-tag:slack-channels", "update",
			"--id", "slch_01SvpyRQK3nN6eH9uhRA7n1q",
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
			"beta:organization:claude-tag:slack-channels", "update",
			"--id", "slch_01SvpyRQK3nN6eH9uhRA7n1q",
			"--beta", "message-batches-2024-09-24",
		)
	})
}

func TestBetaOrganizationClaudeTagSlackChannelsList(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "string")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--base-url", mocktest.ServerWithQuery(t, "beta=true"),
			"beta:organization:claude-tag:slack-channels", "list",
			"--max-items", "10",
			"--slack-workspace-id", "slack_workspace_id",
			"--include-inherited=true",
			"--limit", "0",
			"--page", "page",
			"--slack-channel-id", "slack_channel_id",
			"--beta", "message-batches-2024-09-24",
		)
	})
}

func TestBetaOrganizationClaudeTagSlackChannelsDelete(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "string")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--base-url", mocktest.ServerWithQuery(t, "beta=true"),
			"beta:organization:claude-tag:slack-channels", "delete",
			"--id", "slch_01SvpyRQK3nN6eH9uhRA7n1q",
			"--beta", "message-batches-2024-09-24",
		)
	})
}
