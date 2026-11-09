package cmd

import (
	"testing"

	"github.com/anthropics/anthropic-cli/v2/internal/mocktest"
	"github.com/anthropics/anthropic-cli/v2/internal/requestflag"
)

func TestBetaAgentsCreate(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "string")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--base-url", mocktest.ServerWithQuery(t, "beta=true"),
			"beta:agents", "create",
			"--model", "claude-opus-5-5",
			"--name", "My First Agent",
			"--description", "A general-purpose starter agent.",
			"--execution-identity", "{role_arn: arn:aws:iam::228601553692:role/T, type: aws_role}",
			"--mcp-server", "{name: example-mcp, type: url, url: https://example-server.modelcontextprotocol.io/sse}",
			"--metadata", "{foo: bar}",
			"--multiagent", "{type: multiagent_20261001, advisor: {type: disabled}, subagents: {type: enabled, inline_agents: {type: disabled}, predefined_agents: [agent_011CZkYqphY8vELVzwCUpqiQ]}, workflows: {type: enabled, inline_agents: {type: disabled}, predefined_agents: [{id: agent_011CZkYqphY8vELVzwCUpqiQ, type: agent, version: 1}]}}",
			"--skill", "{skill_id: xlsx, type: anthropic, version: '1'}",
			"--system", "You are a general-purpose agent that can research, write code, run commands, and use connected tools to complete the user's task end to end.",
			"--tool", "{type: agent_toolset_20260401, configs: [{name: bash, enabled: true, permission_policy: {type: always_allow}, type: bash}], default_config: {enabled: true, permission_policy: {type: always_allow}}}",
			"--beta", "message-batches-2024-09-24",
			"--workspace-id", "wrkspc_011CZkZaBF1tNoB5wlCeusgy",
		)
	})

	t.Run("inner flags", func(t *testing.T) {
		// Check that inner flags have been set up correctly
		requestflag.CheckInnerFlags(betaAgentsCreate)

		// Alternative argument passing style using inner flags
		mocktest.TestRunMockTestWithFlags(
			t,
			"--base-url", mocktest.ServerWithQuery(t, "beta=true"),
			"beta:agents", "create",
			"--model", "claude-opus-5-5",
			"--name", "My First Agent",
			"--description", "A general-purpose starter agent.",
			"--execution-identity", "{role_arn: arn:aws:iam::228601553692:role/T, type: aws_role}",
			"--mcp-server.name", "example-mcp",
			"--mcp-server.type", "url",
			"--mcp-server.url", "https://example-server.modelcontextprotocol.io/sse",
			"--metadata", "{foo: bar}",
			"--multiagent", "{type: multiagent_20261001, advisor: {type: disabled}, subagents: {type: enabled, inline_agents: {type: disabled}, predefined_agents: [agent_011CZkYqphY8vELVzwCUpqiQ]}, workflows: {type: enabled, inline_agents: {type: disabled}, predefined_agents: [{id: agent_011CZkYqphY8vELVzwCUpqiQ, type: agent, version: 1}]}}",
			"--skill", "{skill_id: xlsx, type: anthropic, version: '1'}",
			"--system", "You are a general-purpose agent that can research, write code, run commands, and use connected tools to complete the user's task end to end.",
			"--tool", "{type: agent_toolset_20260401, configs: [{name: bash, enabled: true, permission_policy: {type: always_allow}, type: bash}], default_config: {enabled: true, permission_policy: {type: always_allow}}}",
			"--beta", "message-batches-2024-09-24",
			"--workspace-id", "wrkspc_011CZkZaBF1tNoB5wlCeusgy",
		)
	})

	t.Run("piping data", func(t *testing.T) {
		// Test piping YAML data over stdin
		pipeData := []byte("" +
			"model: claude-opus-5-5\n" +
			"name: My First Agent\n" +
			"description: A general-purpose starter agent.\n" +
			"execution_identity:\n" +
			"  role_arn: arn:aws:iam::228601553692:role/T\n" +
			"  type: aws_role\n" +
			"mcp_servers:\n" +
			"  - name: example-mcp\n" +
			"    type: url\n" +
			"    url: https://example-server.modelcontextprotocol.io/sse\n" +
			"metadata:\n" +
			"  foo: bar\n" +
			"multiagent:\n" +
			"  type: multiagent_20261001\n" +
			"  advisor:\n" +
			"    type: disabled\n" +
			"  subagents:\n" +
			"    type: enabled\n" +
			"    inline_agents:\n" +
			"      type: disabled\n" +
			"    predefined_agents:\n" +
			"      - agent_011CZkYqphY8vELVzwCUpqiQ\n" +
			"  workflows:\n" +
			"    type: enabled\n" +
			"    inline_agents:\n" +
			"      type: disabled\n" +
			"    predefined_agents:\n" +
			"      - id: agent_011CZkYqphY8vELVzwCUpqiQ\n" +
			"        type: agent\n" +
			"        version: 1\n" +
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
			"        type: bash\n" +
			"    default_config:\n" +
			"      enabled: true\n" +
			"      permission_policy:\n" +
			"        type: always_allow\n")
		mocktest.TestRunMockTestWithPipeAndFlags(
			t, pipeData,
			"--base-url", mocktest.ServerWithQuery(t, "beta=true"),
			"beta:agents", "create",
			"--beta", "message-batches-2024-09-24",
			"--workspace-id", "wrkspc_011CZkZaBF1tNoB5wlCeusgy",
		)
	})
}

func TestBetaAgentsRetrieve(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "string")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--base-url", mocktest.ServerWithQuery(t, "beta=true"),
			"beta:agents", "retrieve",
			"--agent-id", "agent_011CZkYpogX7uDKUyvBTophP",
			"--version", "0",
			"--beta", "message-batches-2024-09-24",
			"--workspace-id", "wrkspc_011CZkZaBF1tNoB5wlCeusgy",
		)
	})
}

func TestBetaAgentsUpdate(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "string")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--base-url", mocktest.ServerWithQuery(t, "beta=true"),
			"beta:agents", "update",
			"--agent-id", "agent_011CZkYpogX7uDKUyvBTophP",
			"--description", "updated",
			"--execution-identity", "{role_arn: arn:aws:iam::228601553692:role/T, type: aws_role}",
			"--mcp-server", "[{name: example-mcp, type: url, url: https://example-server.modelcontextprotocol.io/sse}]",
			"--metadata", "{foo: string}",
			"--model", "{id: claude-opus-5-5, effort: low, inference_geo: inference_geo, speed: standard}",
			"--multiagent", "{agents: [agent_011CZkYqphY8vELVzwCUpqiQ, {type: self}], type: coordinator}",
			"--name", "name",
			"--skill", "[{skill_id: xlsx, type: anthropic, version: '1'}]",
			"--system", "You are a general-purpose agent that can research, write code, run commands, and use connected tools to complete the user's task end to end.",
			"--tool", "[{type: agent_toolset_20260401, configs: [{name: bash, enabled: true, permission_policy: {type: always_allow}, type: bash}], default_config: {enabled: true, permission_policy: {type: always_allow}}}]",
			"--version", "1",
			"--beta", "message-batches-2024-09-24",
			"--workspace-id", "wrkspc_011CZkZaBF1tNoB5wlCeusgy",
		)
	})

	t.Run("inner flags", func(t *testing.T) {
		// Check that inner flags have been set up correctly
		requestflag.CheckInnerFlags(betaAgentsUpdate)

		// Alternative argument passing style using inner flags
		mocktest.TestRunMockTestWithFlags(
			t,
			"--base-url", mocktest.ServerWithQuery(t, "beta=true"),
			"beta:agents", "update",
			"--agent-id", "agent_011CZkYpogX7uDKUyvBTophP",
			"--description", "updated",
			"--execution-identity", "{role_arn: arn:aws:iam::228601553692:role/T, type: aws_role}",
			"--mcp-server.name", "example-mcp",
			"--mcp-server.type", "url",
			"--mcp-server.url", "https://example-server.modelcontextprotocol.io/sse",
			"--metadata", "{foo: string}",
			"--model", "{id: claude-opus-5-5, effort: low, inference_geo: inference_geo, speed: standard}",
			"--multiagent", "{agents: [agent_011CZkYqphY8vELVzwCUpqiQ, {type: self}], type: coordinator}",
			"--name", "name",
			"--skill", "[{skill_id: xlsx, type: anthropic, version: '1'}]",
			"--system", "You are a general-purpose agent that can research, write code, run commands, and use connected tools to complete the user's task end to end.",
			"--tool", "[{type: agent_toolset_20260401, configs: [{name: bash, enabled: true, permission_policy: {type: always_allow}, type: bash}], default_config: {enabled: true, permission_policy: {type: always_allow}}}]",
			"--version", "1",
			"--beta", "message-batches-2024-09-24",
			"--workspace-id", "wrkspc_011CZkZaBF1tNoB5wlCeusgy",
		)
	})

	t.Run("piping data", func(t *testing.T) {
		// Test piping YAML data over stdin
		pipeData := []byte("" +
			"description: updated\n" +
			"execution_identity:\n" +
			"  role_arn: arn:aws:iam::228601553692:role/T\n" +
			"  type: aws_role\n" +
			"mcp_servers:\n" +
			"  - name: example-mcp\n" +
			"    type: url\n" +
			"    url: https://example-server.modelcontextprotocol.io/sse\n" +
			"metadata:\n" +
			"  foo: string\n" +
			"model:\n" +
			"  id: claude-opus-5-5\n" +
			"  effort: low\n" +
			"  inference_geo: inference_geo\n" +
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
			"        type: bash\n" +
			"    default_config:\n" +
			"      enabled: true\n" +
			"      permission_policy:\n" +
			"        type: always_allow\n" +
			"version: 1\n")
		mocktest.TestRunMockTestWithPipeAndFlags(
			t, pipeData,
			"--base-url", mocktest.ServerWithQuery(t, "beta=true"),
			"beta:agents", "update",
			"--agent-id", "agent_011CZkYpogX7uDKUyvBTophP",
			"--beta", "message-batches-2024-09-24",
			"--workspace-id", "wrkspc_011CZkZaBF1tNoB5wlCeusgy",
		)
	})
}

func TestBetaAgentsList(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "string")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--base-url", mocktest.ServerWithQuery(t, "beta=true"),
			"beta:agents", "list",
			"--max-items", "10",
			"--created-at-gte", "'2019-12-27T18:11:19.117Z'",
			"--created-at-lte", "'2019-12-27T18:11:19.117Z'",
			"--include-archived=true",
			"--limit", "0",
			"--page", "page",
			"--beta", "message-batches-2024-09-24",
			"--workspace-id", "wrkspc_011CZkZaBF1tNoB5wlCeusgy",
		)
	})
}

func TestBetaAgentsArchive(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "string")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--base-url", mocktest.ServerWithQuery(t, "beta=true"),
			"beta:agents", "archive",
			"--agent-id", "agent_011CZkYpogX7uDKUyvBTophP",
			"--beta", "message-batches-2024-09-24",
			"--workspace-id", "wrkspc_011CZkZaBF1tNoB5wlCeusgy",
		)
	})
}
