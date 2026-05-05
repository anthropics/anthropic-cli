// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cmd

import (
	"testing"

	"github.com/anthropics/anthropic-cli/internal/mocktest"
	"github.com/anthropics/anthropic-cli/internal/requestflag"
)

func TestBetaAgentsCreate(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:agents", "create",
			"--model", "{id: claude-opus-4-6, speed: standard}",
			"--name", "My First Agent",
			"--description", "A general-purpose starter agent.",
			"--mcp-server", "{name: example-mcp, type: url, url: https://example-server.modelcontextprotocol.io/sse}",
			"--metadata", "{foo: bar}",
			"--multiagent", "{agents: [agent_011CZkYqphY8vELVzwCUpqiQ, {type: self}], type: coordinator}",
			"--skill", "{skill_id: xlsx, type: anthropic, version: '1'}",
			"--system", "You are a general-purpose agent that can research, write code, run commands, and use connected tools to complete the user's task end to end.",
			"--tool", "{type: agent_toolset_20260401, configs: [{name: bash, enabled: true, permission_policy: {type: always_allow}}], default_config: {enabled: true, permission_policy: {type: always_allow}}}",
			"--beta", "message-batches-2024-09-24",
		)
	})

	t.Run("inner flags", func(t *testing.T) {
		// Check that inner flags have been set up correctly
		requestflag.CheckInnerFlags(betaAgentsCreate)

		// Alternative argument passing style using inner flags
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:agents", "create",
			"--model", "{id: claude-opus-4-6, speed: standard}",
			"--name", "My First Agent",
			"--description", "A general-purpose starter agent.",
			"--mcp-server.name", "example-mcp",
			"--mcp-server.type", "url",
			"--mcp-server.url", "https://example-server.modelcontextprotocol.io/sse",
			"--metadata", "{foo: bar}",
			"--multiagent.agents", "[agent_011CZkYqphY8vELVzwCUpqiQ, {type: self}]",
			"--multiagent.type", "coordinator",
			"--skill", "{skill_id: xlsx, type: anthropic, version: '1'}",
			"--system", "You are a general-purpose agent that can research, write code, run commands, and use connected tools to complete the user's task end to end.",
			"--tool", "{type: agent_toolset_20260401, configs: [{name: bash, enabled: true, permission_policy: {type: always_allow}}], default_config: {enabled: true, permission_policy: {type: always_allow}}}",
			"--beta", "message-batches-2024-09-24",
		)
	})

	t.Run("piping data", func(t *testing.T) {
		// Test piping YAML data over stdin
		pipeData := []byte("" +
			"model:\n" +
			"  id: claude-opus-4-6\n" +
			"  speed: standard\n" +
			"name: My First Agent\n" +
			"description: A general-purpose starter agent.\n" +
			"mcp_servers:\n" +
			"  - name: example-mcp\n" +
			"    type: url\n" +
			"    url: https://example-server.modelcontextprotocol.io/sse\n" +
			"metadata:\n" +
			"  foo: bar\n" +
			"multiagent:\n" +
			"  agents:\n" +
			"    - agent_011CZkYqphY8vELVzwCUpqiQ\n" +
			"    - type: self\n" +
			"  type: coordinator\n" +
			"skills:\n" +
			"  - skill_id: xlsx\n" +
			"    type: anthropic\n" +
			"    version: '1'\n" +
			"system: >-\n" +
			"  You are a general-purpose agent that can research, write code, run commands,\n" +
			"  and use connected tools to complete the user's task end to end.\n" +
			"tools:\n" +
			"  - type: agent_toolset_20260401\n" +
			"    configs:\n" +
			"      - name: bash\n" +
			"        enabled: true\n" +
			"        permission_policy:\n" +
			"          type: always_allow\n" +
			"    default_config:\n" +
			"      enabled: true\n" +
			"      permission_policy:\n" +
			"        type: always_allow\n")
		mocktest.TestRunMockTestWithPipeAndFlags(
			t, pipeData,
			"--api-key", "string",
			"beta:agents", "create",
			"--beta", "message-batches-2024-09-24",
		)
	})
}

func TestBetaAgentsRetrieve(t *testing.T) {
	t.Skip("buildURL drops path-level query params (SDK-4349)")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:agents", "retrieve",
			"--agent-id", "agent_011CZkYpogX7uDKUyvBTophP",
			"--version", "0",
			"--beta", "message-batches-2024-09-24",
		)
	})
}

func TestBetaAgentsUpdate(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:agents", "update",
			"--agent-id", "agent_011CZkYpogX7uDKUyvBTophP",
			"--version", "1",
			"--description", "description",
			"--mcp-server", "[{name: example-mcp, type: url, url: https://example-server.modelcontextprotocol.io/sse}]",
			"--metadata", "{foo: string}",
			"--model", "{id: claude-opus-4-6, speed: standard}",
			"--multiagent", "{agents: [agent_011CZkYqphY8vELVzwCUpqiQ, {type: self}], type: coordinator}",
			"--name", "name",
			"--skill", "[{skill_id: xlsx, type: anthropic, version: '1'}]",
			"--system", "You are a general-purpose agent that can research, write code, run commands, and use connected tools to complete the user's task end to end.",
			"--tool", "[{type: agent_toolset_20260401, configs: [{name: bash, enabled: true, permission_policy: {type: always_allow}}], default_config: {enabled: true, permission_policy: {type: always_allow}}}]",
			"--beta", "message-batches-2024-09-24",
		)
	})

	t.Run("inner flags", func(t *testing.T) {
		// Check that inner flags have been set up correctly
		requestflag.CheckInnerFlags(betaAgentsUpdate)

		// Alternative argument passing style using inner flags
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:agents", "update",
			"--agent-id", "agent_011CZkYpogX7uDKUyvBTophP",
			"--version", "1",
			"--description", "description",
			"--mcp-server.name", "example-mcp",
			"--mcp-server.type", "url",
			"--mcp-server.url", "https://example-server.modelcontextprotocol.io/sse",
			"--metadata", "{foo: string}",
			"--model", "{id: claude-opus-4-6, speed: standard}",
			"--multiagent.agents", "[agent_011CZkYqphY8vELVzwCUpqiQ, {type: self}]",
			"--multiagent.type", "coordinator",
			"--name", "name",
			"--skill", "[{skill_id: xlsx, type: anthropic, version: '1'}]",
			"--system", "You are a general-purpose agent that can research, write code, run commands, and use connected tools to complete the user's task end to end.",
			"--tool", "[{type: agent_toolset_20260401, configs: [{name: bash, enabled: true, permission_policy: {type: always_allow}}], default_config: {enabled: true, permission_policy: {type: always_allow}}}]",
			"--beta", "message-batches-2024-09-24",
		)
	})

	t.Run("piping data", func(t *testing.T) {
		// Test piping YAML data over stdin
		pipeData := []byte("" +
			"version: 1\n" +
			"description: description\n" +
			"mcp_servers:\n" +
			"  - name: example-mcp\n" +
			"    type: url\n" +
			"    url: https://example-server.modelcontextprotocol.io/sse\n" +
			"metadata:\n" +
			"  foo: string\n" +
			"model:\n" +
			"  id: claude-opus-4-6\n" +
			"  speed: standard\n" +
			"multiagent:\n" +
			"  agents:\n" +
			"    - agent_011CZkYqphY8vELVzwCUpqiQ\n" +
			"    - type: self\n" +
			"  type: coordinator\n" +
			"name: name\n" +
			"skills:\n" +
			"  - skill_id: xlsx\n" +
			"    type: anthropic\n" +
			"    version: '1'\n" +
			"system: >-\n" +
			"  You are a general-purpose agent that can research, write code, run commands,\n" +
			"  and use connected tools to complete the user's task end to end.\n" +
			"tools:\n" +
			"  - type: agent_toolset_20260401\n" +
			"    configs:\n" +
			"      - name: bash\n" +
			"        enabled: true\n" +
			"        permission_policy:\n" +
			"          type: always_allow\n" +
			"    default_config:\n" +
			"      enabled: true\n" +
			"      permission_policy:\n" +
			"        type: always_allow\n")
		mocktest.TestRunMockTestWithPipeAndFlags(
			t, pipeData,
			"--api-key", "string",
			"beta:agents", "update",
			"--agent-id", "agent_011CZkYpogX7uDKUyvBTophP",
			"--beta", "message-batches-2024-09-24",
		)
	})
}

func TestBetaAgentsList(t *testing.T) {
	t.Skip("buildURL drops path-level query params (SDK-4349)")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:agents", "list",
			"--max-items", "10",
			"--created-at-gte", "'2019-12-27T18:11:19.117Z'",
			"--created-at-lte", "'2019-12-27T18:11:19.117Z'",
			"--include-archived=true",
			"--limit", "0",
			"--page", "page",
			"--beta", "message-batches-2024-09-24",
		)
	})
}

func TestBetaAgentsArchive(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:agents", "archive",
			"--agent-id", "agent_011CZkYpogX7uDKUyvBTophP",
			"--beta", "message-batches-2024-09-24",
		)
	})
}
