package cmd

import (
	"testing"

	"github.com/anthropics/anthropic-cli/v2/internal/mocktest"
)

func TestBetaOrganizationClaudeTagBundlesTargetsCreate(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "string")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--base-url", mocktest.ServerWithQuery(t, "beta=true"),
			"beta:organization:claude-tag:bundles:targets", "create",
			"--bundle-id", "capp_011CZkZjLP1cWxj5EtMmdgbq",
			"--target", "{slack_channel_id: slch_01SvpyRQK3nN6eH9uhRA7n1q, type: slack_channel}",
			"--beta", "computer-use-2024-10-22",
		)
	})

	t.Run("piping data", func(t *testing.T) {
		// Test piping YAML data over stdin
		pipeData := []byte("" +
			"slack_channel_id: slch_01SvpyRQK3nN6eH9uhRA7n1q\n" +
			"type: slack_channel\n")
		mocktest.TestRunMockTestWithPipeAndFlags(
			t, pipeData,
			"--base-url", mocktest.ServerWithQuery(t, "beta=true"),
			"beta:organization:claude-tag:bundles:targets", "create",
			"--bundle-id", "capp_011CZkZjLP1cWxj5EtMmdgbq",
			"--beta", "computer-use-2024-10-22",
		)
	})
}

func TestBetaOrganizationClaudeTagBundlesTargetsRetrieve(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "string")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--base-url", mocktest.ServerWithQuery(t, "beta=true"),
			"beta:organization:claude-tag:bundles:targets", "retrieve",
			"--bundle-id", "capp_011CZkZjLP1cWxj5EtMmdgbq",
			"--target-id", "btgt_01SvpyRQK3nN6eH9uhRA7n1qHq8TzW5mKc2YbX7vNd4RfJ",
			"--beta", "computer-use-2024-10-22",
		)
	})
}

func TestBetaOrganizationClaudeTagBundlesTargetsList(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "string")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--base-url", mocktest.ServerWithQuery(t, "beta=true"),
			"beta:organization:claude-tag:bundles:targets", "list",
			"--max-items", "10",
			"--bundle-id", "capp_011CZkZjLP1cWxj5EtMmdgbq",
			"--limit", "0",
			"--page", "page",
			"--beta", "computer-use-2024-10-22",
		)
	})
}

func TestBetaOrganizationClaudeTagBundlesTargetsDelete(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "string")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--base-url", mocktest.ServerWithQuery(t, "beta=true"),
			"beta:organization:claude-tag:bundles:targets", "delete",
			"--bundle-id", "capp_011CZkZjLP1cWxj5EtMmdgbq",
			"--target-id", "btgt_01SvpyRQK3nN6eH9uhRA7n1qHq8TzW5mKc2YbX7vNd4RfJ",
			"--beta", "computer-use-2024-10-22",
		)
	})
}
