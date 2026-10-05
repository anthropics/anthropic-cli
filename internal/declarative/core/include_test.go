package core

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// loadGadget loads gadgets/g.md from a tree and returns its body as the rest
// of the engine will see it.
func loadGadget(t *testing.T, files map[string]string) (map[string]any, error) {
	t.Helper()
	root := writeTree(t, files)
	l := NewLoader(testRegistry(), root, nil)
	if err := l.Add(context.Background(), []string{filepath.Join(root, "gadgets/g.md")}); err != nil {
		return nil, err
	}
	return l.Sources()["./gadgets/g.md"].Body, nil
}

func TestIncludeSplicesFilesIntoTheListInPlace(t *testing.T) {
	body, err := loadGadget(t, map[string]string{
		"gadgets/g.md": `---
parts:
  - {name: first}
  - ./parts/pair.yml
  - ../shared/one.json
  - {name: last}
---
prose
`,
		// A sequence contributes each element; a single mapping contributes one.
		"gadgets/parts/pair.yml": "- name: a\n  size: 1\n- name: b\n",
		"shared/one.json":        `{"name": "c", "spec": {"deep": [1, 2]}}`,
	})
	require.NoError(t, err)
	assert.Equal(t, []any{
		map[string]any{"name": "first"},
		map[string]any{"name": "a", "size": uint64(1)},
		map[string]any{"name": "b"},
		map[string]any{"name": "c", "spec": map[string]any{"deep": []any{uint64(1), uint64(2)}}},
		map[string]any{"name": "last"},
	}, body["parts"])
}

func TestIncludeGlobIsAFilterInStableOrder(t *testing.T) {
	body, err := loadGadget(t, map[string]string{
		"gadgets/g.md":           "---\nparts:\n  - ./parts/*\n---\nprose\n",
		"gadgets/parts/b.yml":    "name: b\n",
		"gadgets/parts/a.json":   `[{"name": "a1"}, {"name": "a2"}]`,
		"gadgets/parts/README":   "not data\n",
		"gadgets/parts/notes.md": "not data either\n",
	})
	require.NoError(t, err)
	assert.Equal(t, []any{
		map[string]any{"name": "a1"},
		map[string]any{"name": "a2"},
		map[string]any{"name": "b"},
	}, body["parts"])
}

func TestIncludeLeavesAListWithNoPathsAlone(t *testing.T) {
	body, err := loadGadget(t, map[string]string{
		"gadgets/g.md": "---\nparts:\n  - {name: only}\n---\nprose\n",
	})
	require.NoError(t, err)
	assert.Equal(t, []any{map[string]any{"name": "only"}}, body["parts"])
}

func TestIncludedContentIsPartOfTheHash(t *testing.T) {
	files := map[string]string{
		"gadgets/g.md":       "---\nparts:\n  - ./data/p.yml\n---\nprose\n",
		"gadgets/data/p.yml": "name: v1\n",
	}
	before, err := loadGadget(t, files)
	require.NoError(t, err)
	files["gadgets/data/p.yml"] = "name: v2\n"
	after, err := loadGadget(t, files)
	require.NoError(t, err)

	h1, err := hashBody(before)
	require.NoError(t, err)
	h2, err := hashBody(after)
	require.NoError(t, err)
	assert.NotEqual(t, h1, h2, "editing an included file must plan an update of the resource that includes it")
}

func TestIncludeErrorsNameTheFieldEntryAndFile(t *testing.T) {
	for _, tc := range []struct {
		description string
		files       map[string]string
		wantErr     []string
	}{
		{
			description: "a path that matches nothing is an assertion, and fails naming the path as written",
			files:       map[string]string{"gadgets/g.md": "---\nparts:\n  - ./missing.yml\n---\np\n"},
			wantErr:     []string{"./gadgets/g.md: `parts`[0]", `"./missing.yml" matched nothing on disk`},
		},
		{
			description: "a glob that matches only non-data files contributes nothing, which is refused the same way",
			files:       map[string]string{"gadgets/g.md": "---\nparts:\n  - ./notes/*\n---\np\n", "gadgets/notes/a.md": "prose\n"},
			wantErr:     []string{`"./notes/*" matched nothing on disk`},
		},
		{
			description: "a scalar element cannot land in a list of objects, and the message says which element",
			files: map[string]string{
				"gadgets/g.md":  "---\nparts:\n  - {name: ok}\n  - ./p.yml\n---\np\n",
				"gadgets/p.yml": "- name: fine\n- just a string\n",
			},
			wantErr: []string{"`parts`[1]", "p.yml", "element 1", "must be a mapping"},
		},
		{
			description: "a directory named outright is refused with advice rather than silently read",
			files: map[string]string{
				"gadgets/g.md":      "---\nparts:\n  - ./dir\n---\np\n",
				"gadgets/dir/a.yml": "name: a\n",
			},
			wantErr: []string{"`parts`[0]", "is a directory"},
		},
		{
			description: "a data file placed where a walk would take it for a declaration is refused with advice",
			files:       map[string]string{"gadgets/g.md": "---\nparts:\n  - ./p.yml\n---\np\n", "gadgets/p.yml": "name: looks-like-a-gadget\n"},
			wantErr:     []string{"`parts`[0]", "gadgets/p.yml would itself be read as a resource declaration", "subdirectory"},
		},
		{
			description: "URLs are for resources, not data files",
			files:       map[string]string{"gadgets/g.md": "---\nparts:\n  - https://example.com/x.json\n---\np\n"},
			wantErr:     []string{"only files in the working tree can be included"},
		},
		{
			description: "the field itself must still be a list",
			files:       map[string]string{"gadgets/g.md": "---\nparts: ./p.yml\n---\np\n", "gadgets/p.yml": "name: a\n"},
			wantErr:     []string{"`parts` must be a list"},
		},
		{
			description: "unparseable data reports the file",
			files:       map[string]string{"gadgets/g.md": "---\nparts:\n  - ./p.json\n---\np\n", "gadgets/p.json": "{not json"},
			wantErr:     []string{"p.json", "parsing"},
		},
	} {
		t.Run(tc.description, func(t *testing.T) {
			_, err := loadGadget(t, tc.files)
			require.Error(t, err)
			for _, want := range tc.wantErr {
				assert.Contains(t, err.Error(), want)
			}
		})
	}
}

func TestIncludeRejectsAdditionalYAMLDocuments(t *testing.T) {
	for _, extra := range []string{"name: second\n", "- name: second\n", "null\n", ""} {
		t.Run(extra, func(t *testing.T) {
			_, err := loadGadget(t, map[string]string{
				"gadgets/g.md":            "---\nparts: [./data/parts.yaml]\n---\nbody\n",
				"gadgets/data/parts.yaml": "name: first\n---\n" + extra,
			})
			require.Error(t, err, "later documents must not be silently discarded")
			assert.Contains(t, err.Error(), "parts.yaml")
			assert.Contains(t, err.Error(), "multiple YAML documents")
			assert.Contains(t, err.Error(), "`parts`[0]")
		})
	}
}

func TestIncludedSingleDocumentMarkersAndScalarFences(t *testing.T) {
	body, err := loadGadget(t, map[string]string{
		"gadgets/g.md":            "---\nparts: [./data/parts.yaml]\n---\nbody\n",
		"gadgets/data/parts.yaml": "---\n- name: first\n  text: |\n    before\n    ---\n    after\n...\n# trailing comment\n",
	})
	require.NoError(t, err)
	assert.Equal(t, []any{map[string]any{"name": "first", "text": "before\n---\nafter\n"}}, body["parts"])
}
