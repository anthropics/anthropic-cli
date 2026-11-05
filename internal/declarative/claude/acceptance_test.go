package claude

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/anthropics/anthropic-cli/v2/internal/declarative/acctest"
	"github.com/anthropics/anthropic-cli/v2/internal/declarative/core"
	"github.com/anthropics/anthropic-sdk-go/v2"
	"github.com/anthropics/anthropic-sdk-go/v2/option"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Acceptance tests drive whole plan/apply lifecycles through the real
// SDK-backed client against the real API — whatever the active `ant auth
// login` profile (or ANTHROPIC_API_KEY / ANTHROPIC_BASE_URL) points at. They
// create and destroy real resources, so they are skipped unless ANT_ACC is set:
//
//	ANT_ACC=1 go test ./internal/declarative/claude ./pkg/cmd -run TestAcc -v
//
// Every step is followed by an automatic convergence check (re-plan must be
// empty, re-apply must write nothing), and every test ends by pruning what it
// created — see acctest.Run.

// runAcceptance runs steps against the live API through the SDK-backed client,
// skipping unless ANT_ACC is set.
func runAcceptance(t *testing.T, steps []acctest.Step) {
	t.Helper()
	acctest.SkipUnlessEnabled(t)
	// The SDK's default credential chain: ANTHROPIC_API_KEY, then
	// ANTHROPIC_AUTH_TOKEN, then federation env, then the profile `ant auth
	// login` wrote. ANTHROPIC_BASE_URL overrides the host. No retries: a 409
	// on a stale version will not get fresher, and a flake should be seen.
	client := NewClient(anthropic.NewClient(option.WithMaxRetries(0)))
	acctest.Run(t, Registry(), client, steps)
}

// accTree declares one of every kind with every reference slot in use. Names
// carry {{name}} because the server requires some (a skill's display name) to
// be unique across the organization.
func accTree() map[string]string {
	return map[string]string{
		"skills/pr-writer/SKILL.md":         prWriterFile("{{name}}", "how to write a PR"),
		"skills/pr-writer/examples/good.md": "a good PR\n",
		"agents/code-verifier.md":           verifierFile("verify"),
		"agents/code-reviewer.md": `---
name: {{name}}-reviewer
model: claude-sonnet-4-5
skills:
  - ../skills/pr-writer
multiagent:
  type: coordinator
  agents:
    - ./code-verifier.md
---
review
`,
		"environments/cloud.yml": "name: {{name}}\ndescription: sandbox\n",
		"deployments/nightly.yml": `name: {{name}}
agent: ../agents/code-reviewer.md
environment_id: ../environments/cloud.yml
initial_events:
  - type: user.message
    content: [{type: text, text: Review any open pull requests.}]
`,
	}
}

// prWriterFile is accTree's skill file with the given display_title and body.
func prWriterFile(title, body string) string {
	return "---\nname: pr-writer\ndescription: writes PRs\ndisplay_title: " + title + "\n---\n" + body + "\n"
}

// verifierFile is accTree's sub-agent file with the given system prompt.
func verifierFile(prompt string) string {
	return "---\nname: {{name}}-verifier\nmodel: claude-sonnet-4-5\n---\n" + prompt + "\n"
}

// accTreeKeys are the lockfile keys accTree produces.
var accTreeKeys = []string{
	"./agents/code-reviewer.md", "./agents/code-verifier.md",
	"./deployments/nightly.yml", "./environments/cloud.yml", "./skills/pr-writer",
}

// TestAccLifecycle is the end-to-end path most changes should be checked
// against: cold start, no-op, cascade, drift, teardown.
func TestAccLifecycle(t *testing.T) {
	runAcceptance(t, []acctest.Step{
		{
			Name:  "cold start creates everything in dependency order",
			Write: accTree(),
			ExpectWrites: []string{
				"create skill", "create environment",
				"create agent", "create agent", "create deployment",
			},
			ExpectKeys: accTreeKeys,
			Check: func(t *testing.T, fx *acctest.Fixture) {
				reviewer := fx.Remote(t, "./agents/code-reviewer.md")
				roster := reviewer["multiagent"].(map[string]any)["agents"].([]any)
				require.Len(t, roster, 1)
				assert.Equal(t, fx.ID(t, "./agents/code-verifier.md"), roster[0].(map[string]any)["id"])

				skills := reviewer["skills"].([]any)
				require.Len(t, skills, 1)
				entry := fx.Entry(t, "./skills/pr-writer")
				assert.Equal(t, entry.ID, skills[0].(map[string]any)["skill_id"])
				assert.Equal(t, entry.Version, fmt.Sprint(skills[0].(map[string]any)["version"]),
					"a skill reference is pinned to the version the lockfile recorded")

				deployment := fx.Remote(t, "./deployments/nightly.yml")
				assert.Equal(t, fx.ID(t, "./environments/cloud.yml"), deployment["environment_id"])
			},
		},
		{
			Name:  "editing a sub-agent re-pins its coordinator and the deployment",
			Write: map[string]string{"agents/code-verifier.md": verifierFile("verify harder")},
			ExpectActions: map[string]core.Action{
				"./agents/code-verifier.md": core.ActionUpdate,
				"./agents/code-reviewer.md": core.ActionUpdate,
				"./deployments/nightly.yml": core.ActionUpdate,
				"./skills/pr-writer":        core.ActionNoop,
				"./environments/cloud.yml":  core.ActionNoop,
			},
			ExpectWrites: []string{"update agent", "update agent", "update deployment"},
			Check: func(t *testing.T, fx *acctest.Fixture) {
				verifier := fx.Entry(t, "./agents/code-verifier.md")
				reviewer := fx.Remote(t, "./agents/code-reviewer.md")
				roster := reviewer["multiagent"].(map[string]any)["agents"].([]any)
				assert.Equal(t, verifier.Version, fmt.Sprint(roster[0].(map[string]any)["version"]),
					"the roster must point at the version just written")
			},
		},
		{
			Name: "an out-of-band edit blocks the plan",
			BeforePlan: func(t *testing.T, fx *acctest.Fixture) {
				id := fx.ID(t, "./environments/cloud.yml")
				_, err := fx.Unrecorded.Update(context.Background(), KindEnvironment, id,
					core.Request{Body: map[string]any{"description": "edited in the console"}})
				require.NoError(t, err)
			},
			ExpectBlocked: map[string]string{"./environments/cloud.yml": "--force"},
		},
		{
			Name:          "--force reasserts the file and converges",
			Force:         true,
			ExpectActions: map[string]core.Action{"./environments/cloud.yml": core.ActionUpdate},
			ExpectWrites:  []string{"update environment"},
			Check: func(t *testing.T, fx *acctest.Fixture) {
				assert.Equal(t, "sandbox", fx.Remote(t, "./environments/cloud.yml")["description"])
			},
		},
		{
			Name:       "removing files without --prune only warns",
			Remove:     []string{"deployments", "agents"},
			ExpectNoop: true,
			ExpectKeys: accTreeKeys,
			Check: func(t *testing.T, fx *acctest.Fixture) {
				require.NotEmpty(t, fx.Plan.Warnings)
				assert.Contains(t, fx.Plan.Warnings[0], "--prune")
			},
		},
		{
			Name:  "prune destroys dependents before what they reference",
			Prune: true,
			ExpectWrites: []string{
				"destroy deployment", "destroy agent", "destroy agent",
			},
			ExpectKeys: []string{"./environments/cloud.yml", "./skills/pr-writer"},
		},
		{
			Name:         "prune with nothing left empties the lockfile",
			Remove:       []string{"environments", "skills"},
			Prune:        true,
			ExpectWrites: []string{"destroy environment", "destroy skill"},
			ExpectKeys:   []string{},
		},
	})
}

// TestAccOutOfBandAgentEditBlocksTheApply covers the versioned-kind drift path:
// the server's integer version is the signal, and the stale one must not be
// sent.
func TestAccOutOfBandAgentEditBlocksTheApply(t *testing.T) {
	runAcceptance(t, []acctest.Step{
		{Write: accTree()},
		{
			Name: "console edit bumps the version and blocks",
			BeforePlan: func(t *testing.T, fx *acctest.Fixture) {
				id := fx.ID(t, "./agents/code-verifier.md")
				remote := fx.Remote(t, "./agents/code-verifier.md")
				_, err := fx.Unrecorded.Update(context.Background(), KindAgent, id, core.Request{Body: map[string]any{
					"version":     remote["version"],
					"description": "edited in the console",
				}})
				require.NoError(t, err)
			},
			ExpectBlocked: map[string]string{"./agents/code-verifier.md": "changed outside this config"},
			Check: func(t *testing.T, fx *acctest.Fixture) {
				c := fx.Change(t, "./agents/code-verifier.md")
				assert.True(t, c.Drift)
				assert.Contains(t, c.Reasons[0], "its version moved from")
			},
		},
		{
			Name:  "--force sends the version it read, and dependents follow",
			Force: true,
			ExpectActions: map[string]core.Action{
				"./agents/code-verifier.md": core.ActionUpdate,
				"./agents/code-reviewer.md": core.ActionUpdate,
			},
			Check: func(t *testing.T, fx *acctest.Fixture) {
				remote := fx.Remote(t, "./agents/code-verifier.md")
				assert.Empty(t, remote["description"], "--force reverts the console edit by clearing description")
			},
		},
	})
}

// TestAccArchivedResourceNeedsForceToBeRecreated covers a resource archived out
// of band. Archiving is a one-way door, so the only repair is a new resource
// with a new ID, gated behind the same flag as an overwrite.
func TestAccArchivedResourceNeedsForceToBeRecreated(t *testing.T) {
	var oldID string
	runAcceptance(t, []acctest.Step{
		{Write: accTree()},
		{
			Name: "archived environment is refused",
			BeforePlan: func(t *testing.T, fx *acctest.Fixture) {
				oldID = fx.ID(t, "./environments/cloud.yml")
				require.NoError(t, fx.Unrecorded.Destroy(context.Background(), KindEnvironment, oldID))
			},
			ExpectActions: map[string]core.Action{"./environments/cloud.yml": core.ActionCreate},
			ExpectBlocked: map[string]string{"./environments/cloud.yml": "archived"},
		},
		{
			Name:  "--force recreates it and repoints the deployment",
			Force: true,
			ExpectActions: map[string]core.Action{
				"./environments/cloud.yml":  core.ActionCreate,
				"./deployments/nightly.yml": core.ActionUpdate,
			},
			ExpectWrites: []string{"create environment", "update deployment"},
			Check: func(t *testing.T, fx *acctest.Fixture) {
				newID := fx.ID(t, "./environments/cloud.yml")
				assert.NotEqual(t, oldID, newID)
				assert.Equal(t, newID, fx.Remote(t, "./deployments/nightly.yml")["environment_id"])
			},
		},
	})
}

// TestAccDeletedResourceNeedsForceToBeRecreated checks that a tracked resource
// that has vanished gets the same treatment as an archived one. It uses a skill
// because that is the kind the API deletes outright rather than archives.
func TestAccDeletedResourceNeedsForceToBeRecreated(t *testing.T) {
	var oldID string
	runAcceptance(t, []acctest.Step{
		{Write: map[string]string{
			"skills/s/SKILL.md": "---\nname: s\ndescription: d\ndisplay_title: {{name}}\n---\nbody\n",
		}},
		{
			Name: "deleted skill is refused, naming the old ID",
			BeforePlan: func(t *testing.T, fx *acctest.Fixture) {
				oldID = fx.ID(t, "./skills/s")
				require.NoError(t, fx.Unrecorded.Destroy(context.Background(), KindSkill, oldID))
			},
			ExpectActions: map[string]core.Action{"./skills/s": core.ActionCreate},
			ExpectBlocked: map[string]string{"./skills/s": "no longer on the server"},
			Check: func(t *testing.T, fx *acctest.Fixture) {
				assert.ErrorContains(t, fx.Change(t, "./skills/s").Blocked, oldID)
			},
		},
		{
			Name:         "--force recreates it under a new ID",
			Force:        true,
			ExpectWrites: []string{"create skill"},
			Check: func(t *testing.T, fx *acctest.Fixture) {
				assert.NotEqual(t, oldID, fx.ID(t, "./skills/s"))
			},
		},
	})
}

// TestAccSkillEditsMintVersionsThatReferrersRepinTo checks that a skill edit is
// a new version, that references to it re-pin, and that destroying it clears
// every version first.
func TestAccSkillEditsMintVersionsThatReferrersRepinTo(t *testing.T) {
	var versions []string
	runAcceptance(t, []acctest.Step{
		{
			Write: accTree(),
			Check: func(t *testing.T, fx *acctest.Fixture) {
				versions = []string{fx.Entry(t, "./skills/pr-writer").Version}
			},
		},
		{
			// One run: the new skill version is minted, recorded, and sent to
			// the agent that pins it — the agent update must never carry a
			// placeholder version.
			Name:  "adding a file uploads a version and re-pins the agent in the same run",
			Write: map[string]string{"skills/pr-writer/reference/style.md": "be brief\n"},
			ExpectActions: map[string]core.Action{
				"./skills/pr-writer":        core.ActionUpdate,
				"./agents/code-reviewer.md": core.ActionUpdate,
			},
			ExpectWrites: []string{"update skill", "update agent", "update deployment"},
			Check: func(t *testing.T, fx *acctest.Fixture) {
				entry := fx.Entry(t, "./skills/pr-writer")
				require.NotEmpty(t, entry.Version, "the lock records the version the update minted")
				assert.NotContains(t, versions, entry.Version)
				versions = append(versions, entry.Version)
				assert.Equal(t, entry.Version, fmt.Sprint(fx.Remote(t, "./skills/pr-writer")["latest_version_id"]))

				skills := fx.Remote(t, "./agents/code-reviewer.md")["skills"].([]any)
				assert.Equal(t, entry.Version, fmt.Sprint(skills[0].(map[string]any)["version"]),
					"the agent is pinned to the version just minted")
			},
		},
		{
			Name:         "a third edit mints a third version for teardown to delete",
			Write:        map[string]string{"skills/pr-writer/reference/tone.md": "be kind\n"},
			ExpectWrites: []string{"update skill", "update agent", "update deployment"},
		},
		{
			// SKILL.md's display_title is sent as the skill's display_name,
			// which is the field the plan refuses to change.
			Name:          "changing a skill's display name is refused, not faked",
			Write:         map[string]string{"skills/pr-writer/SKILL.md": prWriterFile("{{name}} renamed", "how to write a PR")},
			ExpectBlocked: map[string]string{"./skills/pr-writer": "display_name cannot be changed"},
		},
		{
			// The API holds every version to the first one's `name`, so an
			// upload under a new name can only fail; the plan says so first.
			Name: "renaming the skill in SKILL.md is refused before anything uploads",
			Write: map[string]string{"skills/pr-writer/SKILL.md": strings.Replace(
				prWriterFile("{{name}}", "how to write a PR"), "name: pr-writer", "name: pr-author", 1)},
			ExpectBlocked: map[string]string{"./skills/pr-writer": "name cannot be changed"},
		},
		{
			Name:       "reverting the name and pruning the rest leaves only the skill",
			Write:      map[string]string{"skills/pr-writer/SKILL.md": prWriterFile("{{name}}", "how to write a PR")},
			Remove:     []string{"agents", "deployments", "environments"},
			Prune:      true,
			ExpectKeys: []string{"./skills/pr-writer"},
		},
		{
			Name:         "pruning the skill clears every version and destroys it",
			Remove:       []string{"skills"},
			Prune:        true,
			ExpectWrites: []string{"destroy skill"},
			ExpectKeys:   []string{},
		},
	})
}

// TestAccDroppedFieldsAreClearedOnTheServer checks that deleting a line from a
// file actually clears the field server-side for fields marked Clearable, and
// that dropping a metadata key deletes it rather than leaving it behind. The
// remote has to be read back explicitly: if the server ignores a null, apply
// still records the response as applied and the next plan is clean, so the
// convergence check passes.
func TestAccDroppedFieldsAreClearedOnTheServer(t *testing.T) {
	runAcceptance(t, []acctest.Step{
		{
			Write: map[string]string{
				"environments/e.yml": "name: {{name}}\ndescription: to be cleared\nmetadata:\n  owner: platform\n  team: infra\n",
				"agents/a.md":        "---\nname: {{name}}\nmodel: claude-sonnet-4-5\ndescription: to be cleared\n---\nto be cleared too\n",
			},
		},
		{
			Name: "dropping a description, a system prompt and a metadata key clears each on the server",
			Write: map[string]string{
				"environments/e.yml": "name: {{name}}\nmetadata:\n  owner: platform\n",
				"agents/a.md":        "---\nname: {{name}}\nmodel: claude-sonnet-4-5\n---\n",
			},
			ExpectWrites: []string{"update environment", "update agent"},
			Check: func(t *testing.T, fx *acctest.Fixture) {
				agent := fx.Remote(t, "./agents/a.md")
				assert.Empty(t, agent["description"], "agent description: dropping the key clears")
				assert.Empty(t, agent["system"], "agent system: dropping the key clears")

				environment := fx.Remote(t, "./environments/e.yml")
				assert.Empty(t, environment["description"], "environment description: dropping the key clears")
				meta := environment["metadata"].(map[string]any)
				assert.Equal(t, "platform", meta["owner"])
				assert.NotContains(t, meta, "team")
			},
		},
	})
}

// TestAccUnpinnedReferenceDoesNotCascade checks that `version: latest` opts a
// reference out of pinning, so editing the target must not rewrite the
// referrer — and the referrer must stay converged afterwards, whatever the
// server echoes back for a versionless reference.
func TestAccUnpinnedReferenceDoesNotCascade(t *testing.T) {
	tree := accTree()
	tree["agents/code-reviewer.md"] = `---
name: {{name}}-reviewer
model: claude-sonnet-4-5
skills:
  - path: ../skills/pr-writer
    version: latest
multiagent:
  type: coordinator
  agents:
    - path: ./code-verifier.md
      version: latest
---
review
`
	tree["deployments/nightly.yml"] = strings.Replace(tree["deployments/nightly.yml"],
		"agent: ../agents/code-reviewer.md", "agent: {path: ../agents/code-reviewer.md, version: latest}", 1)
	runAcceptance(t, []acctest.Step{
		{
			Name:       "versionless references are accepted",
			Write:      tree,
			ExpectKeys: accTreeKeys,
			Check: func(t *testing.T, fx *acctest.Fixture) {
				reviewer := fx.Remote(t, "./agents/code-reviewer.md")
				skill := reviewer["skills"].([]any)[0].(map[string]any)
				assert.Equal(t, fx.ID(t, "./skills/pr-writer"), skill["skill_id"])
				subAgent := reviewer["multiagent"].(map[string]any)["agents"].([]any)[0].(map[string]any)
				assert.Equal(t, fx.ID(t, "./agents/code-verifier.md"), subAgent["id"])
			},
		},
		{
			Name: "editing the targets leaves every referrer alone",
			Write: map[string]string{
				"skills/pr-writer/SKILL.md": prWriterFile("{{name}}", "how to write a better PR"),
				"agents/code-verifier.md":   verifierFile("verify harder"),
			},
			ExpectActions: map[string]core.Action{
				"./skills/pr-writer":        core.ActionUpdate,
				"./agents/code-verifier.md": core.ActionUpdate,
				"./agents/code-reviewer.md": core.ActionNoop,
				"./deployments/nightly.yml": core.ActionNoop,
			},
			ExpectWrites: []string{"update skill", "update agent"},
		},
		{
			Name:       "everything is still converged on the next run",
			ExpectNoop: true,
			Check: func(t *testing.T, fx *acctest.Fixture) {
				// What "latest" means is the server's call, and it differs by
				// slot: a versionless skill reference floats, but a versionless
				// agent reference is resolved to a concrete version when the
				// referrer is written — so an unpinned sub-agent stays on the
				// version it had when the coordinator was last applied.
				reviewer := fx.Remote(t, "./agents/code-reviewer.md")
				skill := reviewer["skills"].([]any)[0].(map[string]any)
				assert.Equal(t, "latest", skill["version"], "skill reference floats")
				subAgent := reviewer["multiagent"].(map[string]any)["agents"].([]any)[0].(map[string]any)
				assert.NotEqual(t, fx.Entry(t, "./agents/code-verifier.md").Version, fmt.Sprint(subAgent["version"]),
					"agent reference was fixed server-side and did not follow the edit")
			},
		},
	})
}

// TestAccServerNormalizationIsNotAChange checks that the server's normalizing
// of what it is sent — a model given as a string comes back as an object,
// omitted lists come back empty or defaulted — never reads as a change. Lists
// the file stops declaring must actually be cleared.
func TestAccServerNormalizationIsNotAChange(t *testing.T) {
	const withToolsAndServer = `---
name: {{name}}
model: claude-sonnet-4-5
description: has tools
tools:
  - type: agent_toolset_20260401
  - type: mcp_toolset
    mcp_server_name: docs
mcp_servers:
  - type: url
    name: docs
    url: https://example.com/mcp
---
prompt
`
	runAcceptance(t, []acctest.Step{
		{
			Name:  "a string model, tools and an MCP server are accepted",
			Write: map[string]string{"agents/a.md": withToolsAndServer},
			Check: func(t *testing.T, fx *acctest.Fixture) {
				remote := fx.Remote(t, "./agents/a.md")
				assert.NotEmpty(t, remote["tools"])
				assert.NotEmpty(t, remote["mcp_servers"])
			},
		},
		{
			Name:       "the same model written as an object is not a change",
			Write:      map[string]string{"agents/a.md": strings.Replace(withToolsAndServer, "model: claude-sonnet-4-5", "model: {id: claude-sonnet-4-5}", 1)},
			ExpectNoop: true,
		},
		{
			Name:         "dropping tools, mcp_servers and description clears them",
			Write:        map[string]string{"agents/a.md": "---\nname: {{name}}\nmodel: claude-sonnet-4-5\n---\nprompt\n"},
			ExpectWrites: []string{"update agent"},
			Check: func(t *testing.T, fx *acctest.Fixture) {
				remote := fx.Remote(t, "./agents/a.md")
				assert.Empty(t, remote["tools"], "tools: dropping the key clears")
				assert.Empty(t, remote["mcp_servers"], "mcp_servers: dropping the key clears")
				assert.Empty(t, remote["description"], "description: dropping the key clears")
				removed := map[string]bool{}
				for name, d := range fx.Change(t, "./agents/a.md").Diff.Fields {
					if d.Kind == core.DiffRemoved {
						removed[name] = true
					}
				}
				assert.Equal(t, map[string]bool{"tools": true, "mcp_servers": true, "description": true}, removed,
					"the plan names exactly the fields the file dropped")
			},
		},
	})
}

// TestAccIDReferencesReachUnmanagedResources covers references that name a
// resource this config does not manage by giving its ID instead of a path. The
// loader tells the two apart by the ID's shape, so this only works if real IDs
// have the shape it expects.
func TestAccIDReferencesReachUnmanagedResources(t *testing.T) {
	runAcceptance(t, []acctest.Step{
		{
			Name: "the targets are created first",
			Write: map[string]string{
				"skills/s/SKILL.md":      "---\nname: s\ndescription: d\ndisplay_title: {{name}}\n---\nbody\n",
				"agents/sub.md":          "---\nname: {{name}}-sub\nmodel: claude-sonnet-4-5\n---\nsub\n",
				"environments/cloud.yml": "name: {{name}}\n",
			},
		},
		{
			Name: "referrers name them by ID and only the referrers are created",
			BeforePlan: func(t *testing.T, fx *acctest.Fixture) {
				// Write is fixed when the step is built, and these files name
				// IDs the first step's apply minted.
				coord := fmt.Sprintf(`---
name: {{name}}-coord
model: claude-sonnet-4-5
skills:
  - %s
multiagent:
  type: coordinator
  agents:
    - %s
---
coordinate
`, fx.ID(t, "./skills/s"), fx.ID(t, "./agents/sub.md"))
				deployment := fmt.Sprintf(`name: {{name}}
agent: ../agents/coord.md
environment_id: %s
initial_events:
  - type: user.message
    content: [{type: text, text: go}]
`, fx.ID(t, "./environments/cloud.yml"))
				acctest.WriteTree(t, fx.Root, acctest.Substitute(map[string]string{
					"agents/coord.md":   coord,
					"deployments/d.yml": deployment,
				}, fx.UniqueName))
			},
			ExpectWrites: []string{"create agent", "create deployment"},
			Check: func(t *testing.T, fx *acctest.Fixture) {
				coord := fx.Remote(t, "./agents/coord.md")
				assert.Equal(t, fx.ID(t, "./skills/s"), coord["skills"].([]any)[0].(map[string]any)["skill_id"])
				assert.Equal(t, fx.ID(t, "./agents/sub.md"),
					coord["multiagent"].(map[string]any)["agents"].([]any)[0].(map[string]any)["id"])
				assert.Equal(t, fx.ID(t, "./environments/cloud.yml"), fx.Remote(t, "./deployments/d.yml")["environment_id"])
			},
		},
		{
			Name:  "editing an ID-referenced target does not touch the referrer",
			Write: map[string]string{"agents/sub.md": "---\nname: {{name}}-sub\nmodel: claude-sonnet-4-5\n---\nsub, edited\n"},
			ExpectActions: map[string]core.Action{
				"./agents/sub.md":   core.ActionUpdate,
				"./agents/coord.md": core.ActionNoop,
			},
			ExpectWrites: []string{"update agent"},
		},
	})
}
