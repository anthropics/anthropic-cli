package cmd

import (
	"testing"

	"github.com/anthropics/anthropic-cli/v2/internal/mocktest"
)

func TestBetaEnvironmentsWorkRetrieve(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "string")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--base-url", mocktest.ServerWithQuery(t, "beta=true"),
			"beta:environments:work", "retrieve",
			"--environment-id", "env_011CZkZ9X2dpNyB7HsEFoRfW",
			"--work-id", "work_id",
			"--beta", "computer-use-2024-10-22",
			"--workspace-id", "wrkspc_011CZkZaBF1tNoB5wlCeusgy",
		)
	})
}

func TestBetaEnvironmentsWorkUpdate(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "string")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--base-url", mocktest.ServerWithQuery(t, "beta=true"),
			"beta:environments:work", "update",
			"--environment-id", "env_011CZkZ9X2dpNyB7HsEFoRfW",
			"--work-id", "work_id",
			"--metadata", "{foo: string}",
			"--beta", "computer-use-2024-10-22",
			"--workspace-id", "wrkspc_011CZkZaBF1tNoB5wlCeusgy",
		)
	})

	t.Run("piping data", func(t *testing.T) {
		// Test piping YAML data over stdin
		pipeData := []byte("" +
			"metadata:\n" +
			"  foo: string\n")
		mocktest.TestRunMockTestWithPipeAndFlags(
			t, pipeData,
			"--base-url", mocktest.ServerWithQuery(t, "beta=true"),
			"beta:environments:work", "update",
			"--environment-id", "env_011CZkZ9X2dpNyB7HsEFoRfW",
			"--work-id", "work_id",
			"--beta", "computer-use-2024-10-22",
			"--workspace-id", "wrkspc_011CZkZaBF1tNoB5wlCeusgy",
		)
	})
}

func TestBetaEnvironmentsWorkList(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "string")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--base-url", mocktest.ServerWithQuery(t, "beta=true"),
			"beta:environments:work", "list",
			"--max-items", "10",
			"--environment-id", "env_011CZkZ9X2dpNyB7HsEFoRfW",
			"--limit", "1",
			"--page", "page",
			"--beta", "computer-use-2024-10-22",
		)
	})
}

func TestBetaEnvironmentsWorkAck(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "string")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--base-url", mocktest.ServerWithQuery(t, "beta=true"),
			"beta:environments:work", "ack",
			"--environment-id", "env_011CZkZ9X2dpNyB7HsEFoRfW",
			"--work-id", "work_id",
			"--beta", "computer-use-2024-10-22",
		)
	})
}

func TestBetaEnvironmentsWorkHeartbeat(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "string")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--base-url", mocktest.ServerWithQuery(t, "beta=true"),
			"beta:environments:work", "heartbeat",
			"--environment-id", "env_011CZkZ9X2dpNyB7HsEFoRfW",
			"--work-id", "work_id",
			"--desired-ttl-seconds", "0",
			"--expected-last-heartbeat", "expected_last_heartbeat",
			"--beta", "computer-use-2024-10-22",
		)
	})
}

func TestBetaEnvironmentsWorkPoll(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "string")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--base-url", mocktest.ServerWithQuery(t, "beta=true"),
			"beta:environments:work", "poll",
			"--environment-id", "env_011CZkZ9X2dpNyB7HsEFoRfW",
			"--block-ms", "1",
			"--reclaim-older-than-ms", "1",
			"--beta", "computer-use-2024-10-22",
			"--anthropic-worker-id", "Anthropic-Worker-ID",
		)
	})
}

func TestBetaEnvironmentsWorkStats(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "string")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--base-url", mocktest.ServerWithQuery(t, "beta=true"),
			"beta:environments:work", "stats",
			"--environment-id", "env_011CZkZ9X2dpNyB7HsEFoRfW",
			"--beta", "computer-use-2024-10-22",
			"--workspace-id", "wrkspc_011CZkZaBF1tNoB5wlCeusgy",
		)
	})
}

func TestBetaEnvironmentsWorkStop(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "string")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--base-url", mocktest.ServerWithQuery(t, "beta=true"),
			"beta:environments:work", "stop",
			"--environment-id", "env_011CZkZ9X2dpNyB7HsEFoRfW",
			"--work-id", "work_id",
			"--force=true",
			"--beta", "computer-use-2024-10-22",
			"--workspace-id", "wrkspc_011CZkZaBF1tNoB5wlCeusgy",
		)
	})

	t.Run("piping data", func(t *testing.T) {
		// Test piping YAML data over stdin
		pipeData := []byte("force: true")
		mocktest.TestRunMockTestWithPipeAndFlags(
			t, pipeData,
			"--base-url", mocktest.ServerWithQuery(t, "beta=true"),
			"beta:environments:work", "stop",
			"--environment-id", "env_011CZkZ9X2dpNyB7HsEFoRfW",
			"--work-id", "work_id",
			"--beta", "computer-use-2024-10-22",
			"--workspace-id", "wrkspc_011CZkZaBF1tNoB5wlCeusgy",
		)
	})
}
