package render

import (
	"bytes"
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/charmbracelet/x/ansi"

	"github.com/anthropics/anthropic-cli/internal/declarative/core"
	"github.com/stretchr/testify/assert"
)

// The renderer is the last thing between a diff and a terminal or CI log. Core
// withholds secret values from the tree; the renderer must not go around it by
// printing Desired or Remote for anything but a create's summary fields.
func TestSensitiveNodesRenderAsAPlaceholder(t *testing.T) {
	var buf bytes.Buffer
	r := &Renderer{Out: &buf}
	r.renderDiff(&core.Diff{Kind: core.DiffObject, Fields: map[string]*core.Diff{
		"resources": {Kind: core.DiffList, Items: []core.ItemDiff{
			{Before: -1, After: 0, Diff: &core.Diff{Kind: core.DiffSensitive}},
			{Before: 0, After: 1, Diff: &core.Diff{Kind: core.DiffObject, Fields: map[string]*core.Diff{
				"authorization_token": {Kind: core.DiffWriteOnly},
				"url":                 {Kind: core.DiffText, Before: "https://a", After: "https://b"},
			}}},
		}},
	}}, "")

	out := buf.String()
	assert.Contains(t, out, "(sensitive; changed)")
	assert.Contains(t, out, "write-only")
	assert.Contains(t, out, `"https://a" → "https://b"`)
	assert.Contains(t, out, "[0→1]")
}

// A one-word edit in a long prompt shows the words around it, not the prompt.
func TestLongTextShowsTheEditInContext(t *testing.T) {
	before := strings.Repeat("lorem ipsum dolor sit amet ", 20) + "review the code carefully" + strings.Repeat(" consectetur adipiscing elit", 20)
	after := strings.Replace(before, "carefully", "thoroughly and kindly", 1)

	var buf bytes.Buffer
	r := &Renderer{Out: &buf}
	r.renderDiff(&core.Diff{Kind: core.DiffObject, Fields: map[string]*core.Diff{
		"system": {Kind: core.DiffText, Before: before, After: after},
	}}, "")

	out := buf.String()
	assert.Contains(t, out, "[-carefully-]{+thoroughly and kindly+}")
	assert.Contains(t, out, "review the code ")
	assert.Contains(t, out, "words…")
	assert.Less(t, len(out), len(before)/2, "most of the unchanged prose is elided")

	buf.Reset()
	r.Verbose = true
	r.renderDiff(&core.Diff{Kind: core.DiffObject, Fields: map[string]*core.Diff{
		"system": {Kind: core.DiffText, Before: before, After: after},
	}}, "")
	assert.NotContains(t, buf.String(), "words…", "verbose keeps the whole text")
}

func TestElidePreservesUnicodeGraphemesAndCellWidth(t *testing.T) {
	for _, tc := range []struct {
		name, unit string
		cells      int
	}{
		{"ASCII", "a", 1}, {"accented", "é", 1}, {"combining", "e\u0301", 1},
		{"CJK", "界", 2}, {"joined emoji", "👩‍💻", 2},
	} {
		t.Run(tc.name, func(t *testing.T) {
			input := strings.Repeat(tc.unit, 12)
			want := strings.Repeat(tc.unit, 4/tc.cells) + "…" + strings.Repeat(tc.unit, 2/tc.cells)
			got := elide(input, 7)
			assert.True(t, utf8.ValidString(got), "invalid UTF-8: %q", got)
			assert.Equal(t, want, got)
			assert.LessOrEqual(t, ansi.StringWidth(got), 7)
			fitting := strings.Repeat(tc.unit, 3)
			assert.Equal(t, fitting, elide(fitting, 7))
		})
	}
}

func TestElidedPlanDetailsKeepValidUnicode(t *testing.T) {
	for _, text := range []string{strings.Repeat("é", 160), strings.Repeat("界", 90), strings.Repeat("👩‍💻", 90)} {
		for _, verbose := range []bool{false, true} {
			t.Run(string([]rune(text)[:1])+map[bool]string{false: " compact", true: " verbose"}[verbose], func(t *testing.T) {
				var buf bytes.Buffer
				renderer := &Renderer{Out: &buf, Verbose: verbose}
				renderer.renderDiff(&core.Diff{Kind: core.DiffObject, Fields: map[string]*core.Diff{
					"instructions": {Kind: core.DiffAdded, After: text},
					"metadata":     {Kind: core.DiffAdded, After: map[string]any{"description": text}},
				}}, "")
				output := buf.String()
				assert.True(t, utf8.ValidString(output), "plan contains invalid UTF-8")
				assert.NotContains(t, output, "�")
				if verbose {
					assert.Contains(t, output, text)
				} else {
					assert.Contains(t, output, "…")
				}
			})
		}
	}
}

func TestElideKeepsASCIIBoundaries(t *testing.T) {
	assert.Equal(t, "abcd…ij", elide("abcdefghij", 7))
	assert.Equal(t, "abcdefg", elide("abcdefg", 7))
	assert.Equal(t, "", elide("", 7))
	assert.Equal(t, "…", elide("abcdefg", 1))
}
