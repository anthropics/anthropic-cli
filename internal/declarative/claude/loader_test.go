package claude

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/anthropics/anthropic-cli/internal/declarative/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// refFixture is a small but complete tree: two agents, two skills, an
// environment, a vault, a credential and a deployment that ties them together.
func refFixture(t *testing.T) string {
	t.Helper()
	return writeTree(t, map[string]string{
		"skills/pr-writer/SKILL.md": "---\nname: pr-writer\n---\nbody\n",
		"skills/changelog/SKILL.md": "---\nname: changelog\n---\nbody\n",
		"skills/README.md":          "not a skill\n",
		"agents/code-verifier.md":   "---\nmodel: m\n---\nverify\n",
		"agents/code-reviewer.md": `---
model: m
skills:
  - ../skills/*
  - type: anthropic
    skill_id: xlsx
multiagent:
  type: coordinator
  agents:
    - ./code-verifier.md
    - type: self
---
review
`,
		"environments/cloud.yml": "description: env\n",
		"deployments/nightly.yml": `agent: ../agents/code-reviewer.md
environment_id: ../environments/cloud.yml
`,
	})
}

func TestReferencesLoadTransitively(t *testing.T) {
	root := refFixture(t)
	l := core.NewLoader(Registry(), root, nil)
	// Naming only the deployment must pull in everything it depends on.
	require.NoError(t, l.Add(context.Background(), []string{filepath.Join(root, "deployments/nightly.yml")}))

	keys := make([]string, 0, len(l.Sources()))
	for k := range l.Sources() {
		keys = append(keys, k)
	}
	assert.ElementsMatch(t, []string{
		"./deployments/nightly.yml",
		"./agents/code-reviewer.md",
		"./agents/code-verifier.md",
		"./environments/cloud.yml",
		"./skills/pr-writer",
		"./skills/changelog",
	}, keys)
}

func TestGlobExpandsToSkillDirectoriesOnly(t *testing.T) {
	root := refFixture(t)
	l := core.NewLoader(Registry(), root, nil)
	require.NoError(t, l.Add(context.Background(), []string{filepath.Join(root, "agents/code-reviewer.md")}))

	// Assert on the body that would be sent, not on the loader's bookkeeping:
	// the glob must expand to the two skill directories, in a stable order,
	// with the hand-written Anthropic skill passed through untouched.
	skills := plannedField(t, l, "./agents/code-reviewer.md", "skills").([]any)
	require.Len(t, skills, 3, "two globbed skills plus the passthrough anthropic skill")
	assert.Equal(t, "custom", skills[0].(map[string]any)["type"])
	assert.Equal(t, "custom", skills[1].(map[string]any)["type"])
	assert.Equal(t, map[string]any{"type": "anthropic", "skill_id": "xlsx"}, skills[2])
}

func TestTopoOrderPlacesDependenciesFirst(t *testing.T) {
	root := refFixture(t)
	l := core.NewLoader(Registry(), root, nil)
	require.NoError(t, l.Add(context.Background(), []string{root}))

	order, err := l.TopoOrder()
	require.NoError(t, err)

	position := map[string]int{}
	for i, k := range order {
		position[k] = i
	}
	assert.Less(t, position["./skills/pr-writer"], position["./agents/code-reviewer.md"])
	assert.Less(t, position["./agents/code-verifier.md"], position["./agents/code-reviewer.md"])
	assert.Less(t, position["./agents/code-reviewer.md"], position["./deployments/nightly.yml"])
	assert.Less(t, position["./environments/cloud.yml"], position["./deployments/nightly.yml"])
}

func TestReferenceCycleIsReportedAsAPath(t *testing.T) {
	root := writeTree(t, map[string]string{
		"agents/a.md": "---\nmodel: m\nmultiagent:\n  type: coordinator\n  agents: [./b.md]\n---\na\n",
		"agents/b.md": "---\nmodel: m\nmultiagent:\n  type: coordinator\n  agents: [./a.md]\n---\nb\n",
	})
	l := core.NewLoader(Registry(), root, nil)
	require.NoError(t, l.Add(context.Background(), []string{filepath.Join(root, "agents/a.md")}))

	_, err := l.TopoOrder()
	require.ErrorContains(t, err, "reference cycle")
	require.ErrorContains(t, err, "→")
}

func TestSelfReferenceByPathIsRejected(t *testing.T) {
	root := writeTree(t, map[string]string{
		"agents/a.md": "---\nmodel: m\nmultiagent:\n  type: coordinator\n  agents: [./a.md]\n---\na\n",
	})
	l := core.NewLoader(Registry(), root, nil)
	err := l.Add(context.Background(), []string{filepath.Join(root, "agents/a.md")})
	require.ErrorContains(t, err, "cannot reference itself")
}

func TestWrongKindReferenceIsRejected(t *testing.T) {
	root := writeTree(t, map[string]string{
		"agents/a.md":            "---\nmodel: m\nskills: [../environments/cloud.yml]\n---\na\n",
		"environments/cloud.yml": "description: env\n",
	})
	l := core.NewLoader(Registry(), root, nil)
	err := l.Add(context.Background(), []string{filepath.Join(root, "agents/a.md")})
	require.ErrorContains(t, err, "is a environment, but this field takes one of: skill")
}

func TestMissingReferenceNamesThePath(t *testing.T) {
	root := writeTree(t, map[string]string{
		"agents/a.md": "---\nmodel: m\nskills: [../skills/nope]\n---\na\n",
	})
	l := core.NewLoader(Registry(), root, nil)
	err := l.Add(context.Background(), []string{filepath.Join(root, "agents/a.md")})
	require.ErrorContains(t, err, "../skills/nope")
}

func TestInlineIDReferenceCreatesNoDependency(t *testing.T) {
	const id = "agent_011CZkYpogX7uDKUyvBTophP"
	root := writeTree(t, map[string]string{
		"agents/a.md": "---\nmodel: m\nmultiagent:\n  type: coordinator\n  agents: [" + id + "]\n---\na\n",
	})
	l := core.NewLoader(Registry(), root, nil)
	require.NoError(t, l.Add(context.Background(), []string{filepath.Join(root, "agents/a.md")}))

	// An inline ID names an unmanaged resource: it is written through as-is,
	// with no version pinned, and pulls nothing else into the plan.
	assert.Len(t, l.Sources(), 1, "an inline ID must not load anything")
	roster := plannedField(t, l, "./agents/a.md", "multiagent").(map[string]any)["agents"].([]any)
	assert.Equal(t, map[string]any{"type": "agent", "id": id}, roster[0])
}

func TestKindForIDTellsIDsFromPaths(t *testing.T) {
	for id, want := range map[string]core.Kind{
		"agent_011CZkYpogX7uDKUyvBTophP":       KindAgent,
		"agent_label_011CZkYpogX7uDKUyvBTophP": KindAgent, // a label between prefix and suffix
		"env_011CZkYpogX7uDKUyvBTophP":         KindEnvironment,
		"agent_defs":                           "", // a directory
		"agent_011CZkYpogX7uDKUyvBToph":        "", // one character short
		"agent_Label_011CZkYpogX7uDKUyvBTophP": "", // labels are lowercase
		"agent_a_b_011CZkYpogX7uDKUyvBTophP":   "", // one label at most
		"nosuch_011CZkYpogX7uDKUyvBTophP":      "",
	} {
		got, ok := Registry().KindForID(id)
		assert.Equal(t, want, got, id)
		assert.Equal(t, want != "", ok, id)
	}
}

// A path whose first segment happens to start with an ID prefix is still a
// path. Misreading it as an inline ID would send the literal string to the API
// and never load the file.
func TestPathStartingWithAnIDPrefixIsStillAPath(t *testing.T) {
	root := writeTree(t, map[string]string{
		"agents/agent_defs/verifier.md": "---\ntype: agent\nmodel: m\n---\nverify\n",
		"agents/a.md":                   "---\nmodel: m\nmultiagent:\n  type: coordinator\n  agents: [agent_defs/verifier.md]\n---\na\n",
	})
	l := core.NewLoader(Registry(), root, nil)
	require.NoError(t, l.Add(context.Background(), []string{filepath.Join(root, "agents/a.md")}))
	assert.Contains(t, l.Sources(), "./agents/agent_defs/verifier.md")
}

func TestVersionLatestOptsOutOfPinning(t *testing.T) {
	root := writeTree(t, map[string]string{
		"skills/s/SKILL.md": "---\nname: s\n---\nbody\n",
		"agents/a.md":       "---\nmodel: m\nskills:\n  - path: ../skills/s\n    version: latest\n---\na\n",
	})
	l := core.NewLoader(Registry(), root, nil)
	require.NoError(t, l.Add(context.Background(), []string{filepath.Join(root, "agents/a.md")}))

	// Unpinned: the reference carries no version, so the server resolves it.
	skills := plannedField(t, l, "./agents/a.md", "skills").([]any)
	require.Len(t, skills, 1)
	assert.NotContains(t, skills[0].(map[string]any), "version")
}

func TestVersionOtherThanLatestIsRejected(t *testing.T) {
	root := writeTree(t, map[string]string{
		"skills/s/SKILL.md": "---\nname: s\n---\nbody\n",
		"agents/a.md":       "---\nmodel: m\nskills:\n  - path: ../skills/s\n    version: \"123\"\n---\na\n",
	})
	l := core.NewLoader(Registry(), root, nil)
	err := l.Add(context.Background(), []string{filepath.Join(root, "agents/a.md")})
	require.ErrorContains(t, err, `may only be "latest"`)
}

func TestGlobKeepsAMalformedDeclarationSoItIsReported(t *testing.T) {
	// A glob is a filter, so an unrecognizable file is skipped. But a file that
	// declares what it is and then fails to parse must stay in the match set —
	// dropping it silently turns a typo into a resource that quietly stops
	// being managed.
	root := writeTree(t, map[string]string{
		"agents/a.md":            "---\nmodel: m\nskills: [../skills/*]\n---\na\n",
		"skills/good/SKILL.md":   "---\nname: good\n---\nbody\n",
		"skills/notes.md":        "just prose, no frontmatter\n",
		"skills/broken/SKILL.md": "---\nname: [unclosed\n---\nbody\n",
	})
	l := core.NewLoader(Registry(), root, nil)
	err := l.Add(context.Background(), []string{filepath.Join(root, "agents/a.md")})

	require.Error(t, err, "the broken skill must surface, not be filtered out")
	assert.Contains(t, err.Error(), "broken")
}

func TestAgentToolsIncludeGeneratedSpecsAndWalkIgnoresTheDataFile(t *testing.T) {
	// The shape a framework produces: it writes the custom tool specs it will
	// execute to a data file, and the agent lists that file among the tools a
	// person wrote. Walking ./agents must not mistake the data file for an
	// agent, and the body sent must be the spliced list in declaration order.
	root := writeTree(t, map[string]string{
		"agents/ops.md": `---
model: m
tools:
  - type: agent_toolset_20260401
  - ./tools/ops.json
mcp_servers:
  - ./mcp/*.yml
---
ops
`,
		"agents/tools/ops.json": `[
  {"type": "custom", "name": "service_status", "description": "d", "input_schema": {"type": "object"}},
  {"type": "custom", "name": "restart_service", "description": "d", "input_schema": {"type": "object"}}
]`,
		"agents/mcp/docs.yml": "name: docs\ntype: url\nurl: https://mcp.example.com/sse\n",
	})
	l := core.NewLoader(Registry(), root, nil)
	require.NoError(t, l.Add(context.Background(), []string{filepath.Join(root, "agents")}))

	keys := make([]string, 0, len(l.Sources()))
	for k := range l.Sources() {
		keys = append(keys, k)
	}
	assert.Equal(t, []string{"./agents/ops.md"}, keys, "the included data files are not resources")

	tools := plannedField(t, l, "./agents/ops.md", "tools").([]any)
	require.Len(t, tools, 3)
	assert.Equal(t, map[string]any{"type": "agent_toolset_20260401"}, tools[0])
	assert.Equal(t, "service_status", tools[1].(map[string]any)["name"])
	assert.Equal(t, "restart_service", tools[2].(map[string]any)["name"])

	servers := plannedField(t, l, "./agents/ops.md", "mcp_servers").([]any)
	assert.Equal(t, []any{map[string]any{"name": "docs", "type": "url", "url": "https://mcp.example.com/sse"}}, servers)
}

func TestIncludedDataFileDirectlyInsideAKindDirectoryIsRefused(t *testing.T) {
	// agents/*.yml is where agents live, so a tools file dropped straight in
	// there would become an agent on the next `apply ./agents`. Say so at the
	// include site instead.
	root := writeTree(t, map[string]string{
		"agents/ops.md":    "---\nmodel: m\ntools:\n  - ./tools.yml\n---\nops\n",
		"agents/tools.yml": "- {type: custom, name: t, input_schema: {type: object}}\n",
		"agents/specs.yml": "type: custom\nname: t\n",
		"agents/other.md":  "---\nmodel: m\ntools:\n  - ./specs.yml\n---\nother\n",
	})
	l := core.NewLoader(Registry(), root, nil)
	// A YAML list is not a declaration, so that placement is merely untidy...
	require.NoError(t, l.Add(context.Background(), []string{filepath.Join(root, "agents/ops.md")}))
	// ...but a mapping there is one, and is refused.
	l = core.NewLoader(Registry(), root, nil)
	err := l.Add(context.Background(), []string{filepath.Join(root, "agents/other.md")})
	require.ErrorContains(t, err, "agents/specs.yml would itself be read as a resource declaration")
}

func TestDeploymentResourcesMixReferencesAndPassthroughs(t *testing.T) {
	root := writeTree(t, map[string]string{
		"agents/a.md":              "---\nmodel: m\n---\na\n",
		"environments/e.yml":       "name: e\n",
		"memory_stores/prefs.yml":  "description: what the customer likes\n",
		"memory_stores/policy.yml": "name: returns-policy\n",
		"deployments/d.yml": `agent: ../agents/a.md
environment_id: ../environments/e.yml
resources:
  - ../memory_stores/prefs.yml
  - path: ../memory_stores/policy.yml
    access: read_only
    instructions: Check before promising a refund.
  - {type: file, file_id: file_011CZ, mount_path: /mnt/data.csv}
`,
	})
	l := core.NewLoader(Registry(), root, nil)
	require.NoError(t, l.Add(context.Background(), []string{filepath.Join(root, "deployments/d.yml")}))

	assert.Contains(t, l.Sources(), "./memory_stores/prefs.yml", "a path in resources must load the store it names")
	resources := plannedField(t, l, "./deployments/d.yml", "resources").([]any)
	require.Len(t, resources, 3)
	assert.Equal(t, map[string]any{"type": "memory_store", "memory_store_id": core.KnownAfterApply}, resources[0])
	assert.Equal(t, map[string]any{
		"type": "memory_store", "memory_store_id": core.KnownAfterApply,
		"access": "read_only", "instructions": "Check before promising a refund.",
	}, resources[1], "keys written alongside `path` ride along on the encoded reference")
	assert.Equal(t, map[string]any{"type": "file", "file_id": "file_011CZ", "mount_path": "/mnt/data.csv"}, resources[2])
}

func TestExtraKeysCannotOverrideTheReference(t *testing.T) {
	root := writeTree(t, map[string]string{
		"agents/a.md":             "---\nmodel: m\n---\na\n",
		"memory_stores/prefs.yml": "name: prefs\n",
		"deployments/d.yml":       "agent: ../agents/a.md\nresources:\n  - {path: ../memory_stores/prefs.yml, memory_store_id: memstore_01other}\n",
	})
	l := core.NewLoader(Registry(), root, nil)
	require.NoError(t, l.Add(context.Background(), []string{filepath.Join(root, "deployments/d.yml")}))
	_, err := planBody(t, l, "./deployments/d.yml")
	require.ErrorContains(t, err, "memory_store_id")
}

func TestExtraKeysOnABareIDReferenceAreRejected(t *testing.T) {
	root := writeTree(t, map[string]string{
		"agents/a.md":        "---\nmodel: m\n---\na\n",
		"environments/e.yml": "name: e\n",
		"deployments/d.yml":  "agent: ../agents/a.md\nenvironment_id: {path: ../environments/e.yml, region: eu}\n",
	})
	l := core.NewLoader(Registry(), root, nil)
	require.NoError(t, l.Add(context.Background(), []string{filepath.Join(root, "deployments/d.yml")}))
	_, err := planBody(t, l, "./deployments/d.yml")
	require.ErrorContains(t, err, "bare ID")
}

func TestMemoryStoreDirectoryImpliesKind(t *testing.T) {
	root := writeTree(t, map[string]string{
		"memory_stores/prefs.yml": "description: d\n",
	})
	l := core.NewLoader(Registry(), root, nil)
	require.NoError(t, l.Add(context.Background(), []string{root}))
	src := l.Sources()["./memory_stores/prefs.yml"]
	require.NotNil(t, src)
	assert.Equal(t, KindMemoryStore, src.Kind)
	assert.Equal(t, "prefs", src.Body["name"], "a store with no name takes the filename")
}

func TestEveryFileKindTakesAMarkdownBody(t *testing.T) {
	root := writeTree(t, map[string]string{
		"agents/a.md":            "---\nmodel: m\n---\nBe brief.\n",
		"environments/e.md":      "---\nconfig: {type: cloud}\n---\nSandbox for the nightly review.\n",
		"memory_stores/prefs.md": "---\n---\nWhat the customer likes.\n",
		"deployments/d.md":       "---\nagent: ../agents/a.md\nenvironment_id: ../environments/e.md\n---\nReview any open pull requests.\n",
	})
	l := core.NewLoader(Registry(), root, nil)
	require.NoError(t, l.Add(context.Background(), []string{filepath.Join(root, "deployments/d.md")}))

	assert.Equal(t, "Be brief.", plannedField(t, l, "./agents/a.md", "system"))
	assert.Equal(t, "Sandbox for the nightly review.", plannedField(t, l, "./environments/e.md", "description"))
	assert.Equal(t, []any{map[string]any{
		"type":    "user.message",
		"content": []any{map[string]any{"type": "text", "text": "Review any open pull requests."}},
	}}, plannedField(t, l, "./deployments/d.md", "initial_events"), "a deployment's body is the message every run opens with")

	require.NoError(t, l.Add(context.Background(), []string{filepath.Join(root, "memory_stores/prefs.md")}))
	assert.Equal(t, "What the customer likes.", plannedField(t, l, "./memory_stores/prefs.md", "description"))
}

func TestProseFieldAndBodyTogetherAreRefusedByName(t *testing.T) {
	root := writeTree(t, map[string]string{
		"deployments/d.md": "---\ninitial_events: []\n---\nReview any open pull requests.\n",
	})
	l := core.NewLoader(Registry(), root, nil)
	err := l.Add(context.Background(), []string{filepath.Join(root, "deployments/d.md")})
	require.ErrorContains(t, err, "also sets `initial_events`")
}

func TestAgentFrontmatterScalarFenceKeepsPlannedFields(t *testing.T) {
	root := writeTree(t, map[string]string{
		"agents/example.md": "---\ndescription: |\n  first rule\n  ---\n  last rule\nmodel: chosen-model\n---\nTask body\n",
	})
	loader := core.NewLoader(Registry(), root, nil)
	require.NoError(t, loader.Add(context.Background(), []string{filepath.Join(root, "agents/example.md")}))
	assert.Equal(t, "chosen-model", plannedField(t, loader, "./agents/example.md", "model"))
	assert.Equal(t, "first rule\n---\nlast rule\n", plannedField(t, loader, "./agents/example.md", "description"))
}
