// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cmd

import (
	"testing"

	"github.com/anthropics/anthropic-cli/internal/mocktest"
)

func TestBetaSessionsEventsList(t *testing.T) {
	t.Skip("buildURL drops path-level query params (SDK-4349)")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:sessions:events", "list",
			"--max-items", "10",
			"--session-id", "sesn_011CZkZAtmR3yMPDzynEDxu7",
			"--created-at-gt", "'2019-12-27T18:11:19.117Z'",
			"--created-at-gte", "'2019-12-27T18:11:19.117Z'",
			"--created-at-lt", "'2019-12-27T18:11:19.117Z'",
			"--created-at-lte", "'2019-12-27T18:11:19.117Z'",
			"--limit", "0",
			"--order", "asc",
			"--page", "page",
			"--type", "string",
			"--beta", "message-batches-2024-09-24",
		)
	})
}

func TestBetaSessionsEventsSend(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:sessions:events", "send",
			"--session-id", "sesn_011CZkZAtmR3yMPDzynEDxu7",
			"--event", "{content: [{text: 'Where is my order #1234?', type: text}], type: user.message}",
			"--beta", "message-batches-2024-09-24",
		)
	})

	t.Run("piping data", func(t *testing.T) {
		// Test piping YAML data over stdin
		pipeData := []byte("" +
			"events:\n" +
			"  - content:\n" +
			"      - text: 'Where is my order #1234?'\n" +
			"        type: text\n" +
			"    type: user.message\n")
		mocktest.TestRunMockTestWithPipeAndFlags(
			t, pipeData,
			"--api-key", "string",
			"beta:sessions:events", "send",
			"--session-id", "sesn_011CZkZAtmR3yMPDzynEDxu7",
			"--beta", "message-batches-2024-09-24",
		)
	})
}

func TestBetaSessionsEventsStream(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:sessions:events", "stream",
			"--max-items", "10",
			"--session-id", "sesn_011CZkZAtmR3yMPDzynEDxu7",
			"--beta", "message-batches-2024-09-24",
		)
	})
}
