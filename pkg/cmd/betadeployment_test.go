// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cmd

import (
	"testing"

	"github.com/anthropics/anthropic-cli/internal/mocktest"
	"github.com/anthropics/anthropic-cli/internal/requestflag"
)

func TestBetaDeploymentsCreate(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:deployments", "create",
			"--agent", "string",
			"--environment-id", "x",
			"--initial-event", "{content: [{text: 'Where is my order #1234?', type: text}], type: user.message}",
			"--name", "x",
			"--description", "description",
			"--metadata", "{foo: string}",
			"--resource", "{file_id: file_011CNha8iCJcU1wXNR6q4V8w, type: file, mount_path: /uploads/receipt.pdf}",
			"--schedule", "{expression: 0 9 * * 1-5, timezone: America/Los_Angeles, type: cron}",
			"--vault-id", "string",
			"--beta", "message-batches-2024-09-24",
		)
	})

	t.Run("inner flags", func(t *testing.T) {
		// Check that inner flags have been set up correctly
		requestflag.CheckInnerFlags(betaDeploymentsCreate)

		// Alternative argument passing style using inner flags
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:deployments", "create",
			"--agent", "string",
			"--environment-id", "x",
			"--initial-event", "{content: [{text: 'Where is my order #1234?', type: text}], type: user.message}",
			"--name", "x",
			"--description", "description",
			"--metadata", "{foo: string}",
			"--resource", "{file_id: file_011CNha8iCJcU1wXNR6q4V8w, type: file, mount_path: /uploads/receipt.pdf}",
			"--schedule.expression", "0 9 * * 1-5",
			"--schedule.timezone", "America/Los_Angeles",
			"--schedule.type", "cron",
			"--vault-id", "string",
			"--beta", "message-batches-2024-09-24",
		)
	})

	t.Run("piping data", func(t *testing.T) {
		// Test piping YAML data over stdin
		pipeData := []byte("" +
			"agent: string\n" +
			"environment_id: x\n" +
			"initial_events:\n" +
			"  - content:\n" +
			"      - text: 'Where is my order #1234?'\n" +
			"        type: text\n" +
			"    type: user.message\n" +
			"name: x\n" +
			"description: description\n" +
			"metadata:\n" +
			"  foo: string\n" +
			"resources:\n" +
			"  - file_id: file_011CNha8iCJcU1wXNR6q4V8w\n" +
			"    type: file\n" +
			"    mount_path: /uploads/receipt.pdf\n" +
			"schedule:\n" +
			"  expression: 0 9 * * 1-5\n" +
			"  timezone: America/Los_Angeles\n" +
			"  type: cron\n" +
			"vault_ids:\n" +
			"  - string\n")
		mocktest.TestRunMockTestWithPipeAndFlags(
			t, pipeData,
			"--api-key", "string",
			"beta:deployments", "create",
			"--beta", "message-batches-2024-09-24",
		)
	})
}

func TestBetaDeploymentsRetrieve(t *testing.T) {
	t.Skip("buildURL drops path-level query params (SDK-4349)")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:deployments", "retrieve",
			"--deployment-id", "depl_011CZkZcDH3vPqd7xnEfwTai",
			"--beta", "message-batches-2024-09-24",
		)
	})
}

func TestBetaDeploymentsUpdate(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:deployments", "update",
			"--deployment-id", "depl_011CZkZcDH3vPqd7xnEfwTai",
			"--agent", "string",
			"--description", "description",
			"--environment-id", "environment_id",
			"--initial-event", "{content: [{text: 'Where is my order #1234?', type: text}], type: user.message}",
			"--metadata", "{foo: string}",
			"--name", "name",
			"--resource", "[{file_id: file_011CNha8iCJcU1wXNR6q4V8w, type: file, mount_path: /uploads/receipt.pdf}]",
			"--schedule", "{expression: 0 9 * * 1-5, timezone: America/Los_Angeles, type: cron}",
			"--vault-id", "[string]",
			"--beta", "message-batches-2024-09-24",
		)
	})

	t.Run("inner flags", func(t *testing.T) {
		// Check that inner flags have been set up correctly
		requestflag.CheckInnerFlags(betaDeploymentsUpdate)

		// Alternative argument passing style using inner flags
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:deployments", "update",
			"--deployment-id", "depl_011CZkZcDH3vPqd7xnEfwTai",
			"--agent", "string",
			"--description", "description",
			"--environment-id", "environment_id",
			"--initial-event", "{content: [{text: 'Where is my order #1234?', type: text}], type: user.message}",
			"--metadata", "{foo: string}",
			"--name", "name",
			"--resource", "[{file_id: file_011CNha8iCJcU1wXNR6q4V8w, type: file, mount_path: /uploads/receipt.pdf}]",
			"--schedule.expression", "0 9 * * 1-5",
			"--schedule.timezone", "America/Los_Angeles",
			"--schedule.type", "cron",
			"--vault-id", "[string]",
			"--beta", "message-batches-2024-09-24",
		)
	})

	t.Run("piping data", func(t *testing.T) {
		// Test piping YAML data over stdin
		pipeData := []byte("" +
			"agent: string\n" +
			"description: description\n" +
			"environment_id: environment_id\n" +
			"initial_events:\n" +
			"  - content:\n" +
			"      - text: 'Where is my order #1234?'\n" +
			"        type: text\n" +
			"    type: user.message\n" +
			"metadata:\n" +
			"  foo: string\n" +
			"name: name\n" +
			"resources:\n" +
			"  - file_id: file_011CNha8iCJcU1wXNR6q4V8w\n" +
			"    type: file\n" +
			"    mount_path: /uploads/receipt.pdf\n" +
			"schedule:\n" +
			"  expression: 0 9 * * 1-5\n" +
			"  timezone: America/Los_Angeles\n" +
			"  type: cron\n" +
			"vault_ids:\n" +
			"  - string\n")
		mocktest.TestRunMockTestWithPipeAndFlags(
			t, pipeData,
			"--api-key", "string",
			"beta:deployments", "update",
			"--deployment-id", "depl_011CZkZcDH3vPqd7xnEfwTai",
			"--beta", "message-batches-2024-09-24",
		)
	})
}

func TestBetaDeploymentsList(t *testing.T) {
	t.Skip("buildURL drops path-level query params (SDK-4349)")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:deployments", "list",
			"--max-items", "10",
			"--agent-id", "agent_id",
			"--created-at-gte", "'2019-12-27T18:11:19.117Z'",
			"--created-at-lte", "'2019-12-27T18:11:19.117Z'",
			"--include-archived=true",
			"--limit", "0",
			"--page", "page",
			"--status", "active",
			"--beta", "message-batches-2024-09-24",
		)
	})
}

func TestBetaDeploymentsArchive(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:deployments", "archive",
			"--deployment-id", "depl_011CZkZcDH3vPqd7xnEfwTai",
			"--beta", "message-batches-2024-09-24",
		)
	})
}

func TestBetaDeploymentsPause(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:deployments", "pause",
			"--deployment-id", "depl_011CZkZcDH3vPqd7xnEfwTai",
			"--beta", "message-batches-2024-09-24",
		)
	})
}

func TestBetaDeploymentsRun(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:deployments", "run",
			"--deployment-id", "depl_011CZkZcDH3vPqd7xnEfwTai",
			"--beta", "message-batches-2024-09-24",
		)
	})
}

func TestBetaDeploymentsUnpause(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:deployments", "unpause",
			"--deployment-id", "depl_011CZkZcDH3vPqd7xnEfwTai",
			"--beta", "message-batches-2024-09-24",
		)
	})
}
