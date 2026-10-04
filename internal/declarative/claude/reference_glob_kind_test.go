package claude

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/anthropics/anthropic-cli/internal/declarative/core"
	"github.com/stretchr/testify/require"
)

func TestAgentReferenceGlobSkipsOtherResourceKinds(t *testing.T) {
	root := writeTree(t, map[string]string{
		"agents/coordinator.md": "---\nmodel: m\nmultiagent:\n  type: coordinator\n  agents: [../mixed/*]\n---\ncoordinate\n",
		"mixed/a.yml":           "type: agent\nname: a\nmodel: m\n",
		"mixed/b.yml":           "type: environment\nname: b\n",
		"mixed/c.md":            "---\ntype: agent\nmodel: m\n---\nhelp\n",
		"mixed/skill/SKILL.md":  "---\nname: skill\n---\nbody\n",
	})
	loader := core.NewLoader(Registry(), root, nil)
	require.NoError(t, loader.Add(context.Background(), []string{filepath.Join(root, "agents/coordinator.md")}))
	require.Len(t, loader.Sources(), 3)
	require.NotContains(t, loader.Sources(), "./mixed/b.yml")
	require.NotContains(t, loader.Sources(), "./mixed/skill")
	roster := plannedField(t, loader, "./agents/coordinator.md", "multiagent").(map[string]any)["agents"].([]any)
	require.Len(t, roster, 2)
	for _, item := range roster {
		require.Equal(t, "agent", item.(map[string]any)["type"])
	}
}
