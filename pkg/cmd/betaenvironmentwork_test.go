// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cmd

import (
	"testing"

	"github.com/anthropics/anthropic-cli/internal/mocktest"
)

func TestBetaEnvironmentsWorkRetrieve(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:environments:work", "retrieve",
			"--environment-id", "env_011CZkZ9X2dpNyB7HsEFoRfW",
			"--work-id", "work_id",
			"--beta", "message-batches-2024-09-24",
		)
	})
}

func TestBetaEnvironmentsWorkUpdate(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:environments:work", "update",
			"--environment-id", "env_011CZkZ9X2dpNyB7HsEFoRfW",
			"--work-id", "work_id",
			"--metadata", "{foo: string}",
			"--beta", "message-batches-2024-09-24",
		)
	})

	t.Run("piping data", func(t *testing.T) {
		// Test piping YAML data over stdin
		pipeData := []byte("" +
			"metadata:\n" +
			"  foo: string\n")
		mocktest.TestRunMockTestWithPipeAndFlags(
			t, pipeData,
			"--api-key", "string",
			"beta:environments:work", "update",
			"--environment-id", "env_011CZkZ9X2dpNyB7HsEFoRfW",
			"--work-id", "work_id",
			"--beta", "message-batches-2024-09-24",
		)
	})
}

func TestBetaEnvironmentsWorkList(t *testing.T) {
	t.Skip("buildURL drops path-level query params (SDK-4349)")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:environments:work", "list",
			"--max-items", "10",
			"--environment-id", "env_011CZkZ9X2dpNyB7HsEFoRfW",
			"--limit", "1",
			"--page", "page",
			"--beta", "message-batches-2024-09-24",
		)
	})
}

func TestBetaEnvironmentsWorkAck(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:environments:work", "ack",
			"--environment-id", "env_011CZkZ9X2dpNyB7HsEFoRfW",
			"--work-id", "work_id",
			"--beta", "message-batches-2024-09-24",
		)
	})
}

func TestBetaEnvironmentsWorkHeartbeat(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:environments:work", "heartbeat",
			"--environment-id", "env_011CZkZ9X2dpNyB7HsEFoRfW",
			"--work-id", "work_id",
			"--desired-ttl-seconds", "0",
			"--expected-last-heartbeat", "expected_last_heartbeat",
			"--beta", "message-batches-2024-09-24",
		)
	})
}

func TestBetaEnvironmentsWorkPoll(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:environments:work", "poll",
			"--environment-id", "env_011CZkZ9X2dpNyB7HsEFoRfW",
			"--block-ms", "1",
			"--reclaim-older-than-ms", "1",
			"--beta", "message-batches-2024-09-24",
			"--anthropic-worker-id", "Anthropic-Worker-ID",
		)
	})
}

func TestBetaEnvironmentsWorkStats(t *testing.T) {
	t.Skip("buildURL drops path-level query params (SDK-4349)")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:environments:work", "stats",
			"--environment-id", "env_011CZkZ9X2dpNyB7HsEFoRfW",
			"--beta", "message-batches-2024-09-24",
		)
	})
}

func TestBetaEnvironmentsWorkStop(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:environments:work", "stop",
			"--environment-id", "env_011CZkZ9X2dpNyB7HsEFoRfW",
			"--work-id", "work_id",
			"--force=true",
			"--beta", "message-batches-2024-09-24",
		)
	})

	t.Run("piping data", func(t *testing.T) {
		// Test piping YAML data over stdin
		pipeData := []byte("force: true")
		mocktest.TestRunMockTestWithPipeAndFlags(
			t, pipeData,
			"--api-key", "string",
			"beta:environments:work", "stop",
			"--environment-id", "env_011CZkZ9X2dpNyB7HsEFoRfW",
			"--work-id", "work_id",
			"--beta", "message-batches-2024-09-24",
		)
	})
}
