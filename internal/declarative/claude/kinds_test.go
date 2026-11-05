package claude

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/anthropics/anthropic-cli/v2/internal/declarative/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAgentClearableFields(t *testing.T) {
	spec, ok := Registry().Spec(KindAgent)
	require.True(t, ok)
	assert.ElementsMatch(t,
		[]string{"description", "system", "skills", "tools", "mcp_servers", "multiagent"},
		spec.Clearable())
}

func TestAPlanWithholdsADeploymentsToken(t *testing.T) {
	h := newHarness(t, map[string]string{
		"agents/a.md":        "---\nmodel: m\n---\na\n",
		"environments/e.yml": "name: e\n",
		"deployments/d.yml": `agent: ../agents/a.md
environment_id: ../environments/e.yml
resources:
  - {type: github_repository, url: https://github.com/o/r, authorization_token: t}
`,
	})

	change, ok := h.plan().Change("./deployments/d.yml")

	require.True(t, ok)
	assert.True(t, change.IsSensitive("resources"))
}

func TestAChangedTimestampIsNotDrift(t *testing.T) {
	h := newHarness(t, map[string]string{
		"environments/e.yml":  "name: e\n",
		"memory_stores/m.yml": "name: m\n",
		"vaults/v.yml":        "display_name: v\n",
		"agents/a.md":         "---\nmodel: m\n---\na\n",
		"deployments/d.yml":   "agent: ../agents/a.md\nenvironment_id: ../environments/e.yml\n",
	})
	h.apply()
	for _, obj := range h.client.objects {
		obj["updated_at"] = "2026-01-01T00:00:00Z"
	}

	plan := h.plan()

	assert.False(t, plan.HasWork())
}

func TestNothingIsDriftWhileUpdatedAtHasNotMoved(t *testing.T) {
	h := newHarness(t, map[string]string{
		"environments/e.yml":  "name: e\n",
		"memory_stores/m.yml": "name: m\n",
		"vaults/v.yml":        "display_name: v\n",
		"agents/a.md":         "---\nmodel: m\n---\na\n",
		"deployments/d.yml":   "agent: ../agents/a.md\nenvironment_id: ../environments/e.yml\n",
	})
	h.client.now = "2026-01-01T00:00:00Z"
	h.apply()
	for _, obj := range h.client.objects {
		obj["added_by_the_api"] = "default"
		obj["metadata"] = map[string]any{"filled_in_by": "the API"}
	}

	plan := h.plan()

	assert.False(t, plan.HasWork())
}

func TestAChangedDescriptionIsDrift(t *testing.T) {
	h := newHarness(t, map[string]string{"environments/e.yml": "name: e\n"})
	h.apply()
	id := h.lock.Resources["./environments/e.yml"].ID
	h.client.objects[id]["description"] = "edited elsewhere"

	plan := h.plan()

	assert.Len(t, plan.Blocked(), 1)
}

func TestDeploymentVaultIDsResolvePathsAndKeepLiteralIDs(t *testing.T) {
	h := newHarness(t, map[string]string{
		"agents/a.md":        "---\nmodel: m\n---\na\n",
		"environments/e.yml": "name: e\n",
		"vaults/keys.yml":    "display_name: vendor keys\n",
		"deployments/d.yml": `agent: ../agents/a.md
environment_id: ../environments/e.yml
vault_ids:
  - ../vaults/keys.yml
  - vlt_011CZkYpogX7uDKUyvBTophP
`,
	})

	h.apply()

	deployment := h.client.objects[h.lock.Resources["./deployments/d.yml"].ID]
	keys := h.lock.Resources["./vaults/keys.yml"].ID
	assert.Equal(t, []any{keys, "vlt_011CZkYpogX7uDKUyvBTophP"}, deployment["vault_ids"])
}

func TestDeletingTheVaultIDsLineDetachesTheVaults(t *testing.T) {
	const required = "agent: ../agents/a.md\nenvironment_id: ../environments/e.yml\n"
	h := newHarness(t, map[string]string{
		"agents/a.md":        "---\nmodel: m\n---\na\n",
		"environments/e.yml": "name: e\n",
		"deployments/d.yml":  required + "vault_ids: [vlt_011CZkYpogX7uDKUyvBTophP]\n",
	})
	h.apply()
	require.NoError(t, os.WriteFile(filepath.Join(h.root, "deployments/d.yml"), []byte(required), 0o644))

	h.apply()

	deployment := h.client.objects[h.lock.Resources["./deployments/d.yml"].ID]
	assert.NotContains(t, deployment, "vault_ids")
}

func TestKindOrderIsTopological(t *testing.T) {
	// Destroy ordering and plan tie-breaking rely on a kind appearing after
	// everything it can reference.
	r := Registry()
	position := map[core.Kind]int{}
	for i, k := range r.Kinds() {
		position[k] = i
	}
	for _, k := range r.Kinds() {
		spec, _ := r.Spec(k)
		for _, slot := range spec.RefSlots() {
			for _, allowed := range slot.To {
				if allowed == k {
					continue // self-reference orders against itself
				}
				assert.Less(t, position[allowed], position[k],
					"%s references %s, so %s must be registered first", k, allowed, allowed)
			}
		}
	}
}
