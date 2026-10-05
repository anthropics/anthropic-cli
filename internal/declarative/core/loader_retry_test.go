package core

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoaderRetriesFailedReferences(t *testing.T) {
	for _, mode := range []string{"named", "keys", "directory"} {
		t.Run(mode, func(t *testing.T) {
			root := writeTree(t, map[string]string{
				"widgets/stable.yml": "name: stable\n",
				"widgets/first.yml":  "name: first\n",
				"gadgets/one.yml":    "name: one\nwidgets: [../widgets/first.yml, ../widgets/missing.yml]\n",
			})
			l := NewLoader(testRegistry(), root, nil)
			require.NoError(t, l.Add(context.Background(), []string{filepath.Join(root, "widgets/stable.yml")}))
			original := l.Sources()["./widgets/stable.yml"]
			exposed := l.Sources()
			load := func() error {
				switch mode {
				case "keys":
					return l.AddKeys(context.Background(), []string{"./gadgets/one.yml"})
				case "directory":
					return l.Add(context.Background(), []string{filepath.Join(root, "gadgets")})
				default:
					return l.Add(context.Background(), []string{filepath.Join(root, "gadgets/one.yml")})
				}
			}
			for range 2 {
				require.Error(t, load(), "a repeated call must not treat an unresolved source as loaded")
				require.Len(t, l.Sources(), 1)
				require.Len(t, exposed, 1, "rollback preserves the Sources map already held by callers")
				require.Same(t, original, l.Sources()["./widgets/stable.yml"])
				assert.NotContains(t, l.slots, "./gadgets/one.yml")
				assert.NotContains(t, l.slots, "./widgets/first.yml")
			}
			writeTreeAt(t, root, map[string]string{"widgets/missing.yml": "name: missing\n"})
			require.NoError(t, load())
			require.Len(t, l.Sources(), 4)
			order, err := l.TopoOrder()
			require.NoError(t, err)
			require.Equal(t, "./gadgets/one.yml", order[len(order)-1])
			desired, err := (bodyBuilder{registry: l.registry, slots: l.slots}).build(l.Sources()["./gadgets/one.yml"], map[string]Target{
				"./widgets/first.yml":   {ID: "wdgt_first", Known: true},
				"./widgets/missing.yml": {ID: "wdgt_missing", Known: true},
			})
			require.NoError(t, err)
			assert.Equal(t, []any{"wdgt_first", "wdgt_missing"}, desired.Body["widgets"])
			assert.False(t, desired.Unresolved)
		})
	}
}

func TestLoaderRetryReadsCorrectedDeclaration(t *testing.T) {
	root := writeTree(t, map[string]string{"gadgets/one.yml": "widgets: [../widgets/missing.yml]\n"})
	l := NewLoader(testRegistry(), root, nil)
	path := filepath.Join(root, "gadgets/one.yml")
	require.Error(t, l.Add(context.Background(), []string{path}))
	writeTreeAt(t, root, map[string]string{"gadgets/one.yml": "name: corrected\nwidgets: []\n"})
	require.NoError(t, l.Add(context.Background(), []string{path}))
	assert.Equal(t, "corrected", l.Sources()["./gadgets/one.yml"].Body["name"])
	assert.Empty(t, l.dependencies("./gadgets/one.yml"))
}

func TestLoaderFailureRetainsEarlierExplicitArguments(t *testing.T) {
	root := writeTree(t, map[string]string{"widgets/one.yml": "name: one\n", "gadgets/broken.yml": "widgets: [missing.yml]\n"})
	l := NewLoader(testRegistry(), root, nil)
	err := l.Add(context.Background(), []string{filepath.Join(root, "widgets/one.yml"), filepath.Join(root, "gadgets/broken.yml")})
	require.Error(t, err)
	assert.Len(t, l.Sources(), 1)
	assert.Contains(t, l.Sources(), "./widgets/one.yml")
}

func TestLoaderRetryRestoresCyclesAfterAnUnresolvedDependency(t *testing.T) {
	registry := NewRegistry(testSchema{}, KindSpec{Kind: "gadget", IDPrefix: "gdgt", Build: buildNamed,
		Fields: Fields{"parts": {Ref: &Ref{To: []Kind{"gadget"}, List: true, As: EncodeID}}},
	})
	root := writeTree(t, map[string]string{
		"gadgets/a.yml": "parts: [b.yml, missing.yml]\n",
		"gadgets/b.yml": "parts: [a.yml]\n",
	})
	l := NewLoader(registry, root, nil)
	path := filepath.Join(root, "gadgets/a.yml")
	require.Error(t, l.Add(context.Background(), []string{path}))
	assert.Empty(t, l.Sources(), "a successfully visited child may still reference a failed ancestor")
	assert.Empty(t, l.slots)
	writeTreeAt(t, root, map[string]string{"gadgets/missing.yml": "name: missing\n"})
	require.NoError(t, l.Add(context.Background(), []string{path}))
	_, err := l.TopoOrder()
	require.ErrorContains(t, err, "reference cycle")
	assert.Contains(t, err.Error(), "./gadgets/a.yml")
	assert.Contains(t, err.Error(), "./gadgets/b.yml")
}
