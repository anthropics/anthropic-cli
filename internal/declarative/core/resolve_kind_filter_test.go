package core

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReferenceGlobFiltersDeclaredKinds(t *testing.T) {
	root := writeTree(t, map[string]string{
		"mixed/a.yml":     "type: widget\nname: a\n",
		"mixed/b.yml":     "type: gadget\nname: b\n",
		"mixed/c.yml":     "type: widget\nname: c\n",
		"mixed/readme.md": "Ordinary prose\n",
	})
	for _, tc := range []struct {
		name    string
		allowed []Kind
		want    []string
	}{
		{"widget", []Kind{"widget"}, []string{"a.yml", "c.yml"}},
		{"gadget", []Kind{"gadget"}, []string{"b.yml"}},
		{"both", []Kind{"widget", "gadget"}, []string{"a.yml", "b.yml", "c.yml"}},
		{"neither", []Kind{"other"}, nil},
	} {
		t.Run(tc.name, func(t *testing.T) {
			loader := NewLoader(testRegistry(), root, nil)
			matches, err := loader.matchRef(filepath.Join(root, "mixed", "*"), RefSlot{Ref: Ref{To: tc.allowed}})
			require.NoError(t, err)
			var names []string
			for _, match := range matches {
				names = append(names, filepath.Base(match))
			}
			require.Equal(t, tc.want, names)
		})
	}
}

func TestReferenceGlobLoadsOnlyEligibleDependencies(t *testing.T) {
	root := writeTree(t, map[string]string{
		"gadgets/main.yml": "name: main\nwidgets: [../mixed/*.yml]\n",
		"mixed/a.yml":      "type: widget\nname: a\n",
		"mixed/b.yml":      "type: gadget\nname: b\nwidgets: [./missing.yml]\n",
		"mixed/c.yml":      "type: widget\nname: c\n",
	})
	loader := NewLoader(testRegistry(), root, nil)
	require.NoError(t, loader.Add(context.Background(), []string{filepath.Join(root, "gadgets/main.yml")}))
	require.Len(t, loader.Sources(), 3)
	require.Contains(t, loader.Sources(), "./mixed/a.yml")
	require.Contains(t, loader.Sources(), "./mixed/c.yml")
	require.NotContains(t, loader.Sources(), "./mixed/b.yml")
	order, err := loader.TopoOrder()
	require.NoError(t, err)
	require.Equal(t, "./gadgets/main.yml", order[len(order)-1])
}

func TestExplicitWrongKindReferenceStillFails(t *testing.T) {
	root := writeTree(t, map[string]string{
		"gadgets/main.yml": "name: main\nwidgets: [../mixed/b.yml]\n",
		"mixed/b.yml":      "type: gadget\nname: b\n",
	})
	loader := NewLoader(testRegistry(), root, nil)
	err := loader.Add(context.Background(), []string{filepath.Join(root, "gadgets/main.yml")})
	require.ErrorContains(t, err, "is a gadget, but this field takes one of: widget")
}

func TestReferenceGlobWithOnlyWrongKindsReportsNoMatch(t *testing.T) {
	root := writeTree(t, map[string]string{
		"gadgets/main.yml": "name: main\nwidgets: [../mixed/*.yml]\n",
		"mixed/b.yml":      "type: gadget\nname: b\n",
	})
	loader := NewLoader(testRegistry(), root, nil)
	err := loader.Add(context.Background(), []string{filepath.Join(root, "gadgets/main.yml")})
	require.ErrorContains(t, err, "matched nothing on disk")
}
