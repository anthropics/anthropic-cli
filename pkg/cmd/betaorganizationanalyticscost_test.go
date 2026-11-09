package cmd

import (
	"testing"

	"github.com/anthropics/anthropic-cli/v2/internal/mocktest"
)

func TestBetaOrganizationAnalyticsCostList(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "string")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--base-url", mocktest.ServerWithQuery(t, "beta=true"),
			"beta:organization:analytics:cost", "list",
			"--max-items", "10",
			"--starting-at", "'2026-09-01T00:00:00Z'",
			"--actor-id", "user_01WxYzAbCdEfGhJkMnPqRsTu",
			"--actor-type", "service_account_actor",
			"--agent-id", "agent_01VHd3EAGJ6PLkP9LM7DDu3f",
			"--api-key-id", "apikey_01Rj2N8SVvo6BePZj99NhmiT",
			"--claude-tag-category", "engaged",
			"--claude-tag-user-id", "U0123ABCDEF",
			"--context-window", "0-200k",
			"--deployment-id", "depl_01JFC6ihzZHtAuyHsaENdJhY",
			"--ending-at", "'2026-10-01T00:00:00Z'",
			"--environment-id", "env_017oTjuGedHtjhJauVC5xGJD",
			"--exclude-deleted-users=true",
			"--group-by", "actor",
			"--inference-geo", "global",
			"--limit", "1",
			"--model", "string",
			"--order", "asc",
			"--order-by", "amount",
			"--organization-id", "12345678-1234-5678-1234-567812345678",
			"--page", "page",
			"--product", "chat",
			"--rbac-group-id", "rbac_group_012rppKaSVsmTo6NqRDXQXNF",
			"--service-tier", "batch",
			"--session-id", "sesn_0187kPDyNCwibNNvfaPEApCJ",
			"--slack-channel-id", "C0123ABCDEF",
			"--speed", "fast",
			"--time-bucket-width", "1d",
			"--user-profile-id", "uprof_01M33NeYUQK2WV9PpYcsPHLi",
			"--workspace-id", "wrkspc_01JwQvzr7rXLA5AGx3HKfFUJ",
			"--beta", "message-batches-2024-09-24",
		)
	})
}
