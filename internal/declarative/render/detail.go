package render

import (
	"encoding/json"
	"fmt"
	"maps"
	"slices"
	"strings"

	"github.com/charmbracelet/x/ansi"

	"github.com/anthropics/anthropic-cli/internal/declarative/core"
)

// Rendering what an expanded row shows beneath it: the identifying fields of a
// create, and for an update the core.Diff tree. The tree follows the body's
// own shape, so the output reads like the file the user wrote: nested fields
// indent, list elements are labelled by index, and only what changed appears.

// HasDetails reports whether expanding would show anything, so a prompt can
// decide whether to offer it.
func (r *Renderer) HasDetails(plan *core.Plan) bool {
	for _, c := range plan.Changes {
		switch c.Action {
		case core.ActionCreate, core.ActionUpdate:
			return true
		case core.ActionDestroy:
			if len(c.Reasons) > 0 {
				return true
			}
		}
	}
	return false
}

// detail writes what sits beneath an expanded row: for a create the fields
// that identify it, for an update the diff tree, and the reasons behind any
// change that has them. It reports whether it wrote anything.
func (r *Renderer) detail(c *core.Change, indent string) bool {
	wrote := false
	switch c.Action {
	case core.ActionCreate:
		wrote = r.createDetail(c, indent)
	case core.ActionUpdate:
		r.renderDiff(c.Diff, indent)
		wrote = c.Diff != nil
	}
	mark, style := "· ", dimStyle
	if c.Drift {
		mark, style = "! ", updateStyle
	}
	for _, reason := range c.Reasons {
		r.line(indent + r.paint(style, mark+reason))
	}
	return wrote || len(c.Reasons) > 0
}

// createDetail shows the handful of fields that identify a new resource, or
// every field when verbose.
func (r *Renderer) createDetail(c *core.Change, indent string) bool {
	if what, ok := c.Upload(); ok {
		r.line(indent + r.paint(dimStyle, "upload ") + what)
		return true
	}
	keys := c.SummaryFields()
	if r.Verbose {
		keys = slices.Sorted(maps.Keys(c.Desired))
	}
	wrote := false
	for _, k := range keys {
		v, ok := c.Desired[k]
		if !ok || core.IsEmptyValue(v) {
			continue
		}
		val := r.value(v)
		if c.IsSensitive(k) {
			val = r.paint(dimStyle, "(sensitive)")
		}
		r.wrapLine(indent+r.paint(dimStyle, k+": "), indent+indentStep, val)
		wrote = true
	}
	return wrote
}

const (
	// contextWords is how much unchanged prose is kept either side of an
	// edit inside a long string.
	contextWords = 8
	// inlineTextLimit is the longest single-line before/after pair still
	// shown as `"a" → "b"`; past it the word diff takes over.
	inlineTextLimit = 40
)

// renderDiff writes the tree beneath a table row.
func (r *Renderer) renderDiff(d *core.Diff, indent string) {
	if d == nil {
		return
	}
	switch d.Kind {
	case core.DiffObject:
		r.renderFields(d, indent)
	default:
		// The planner always roots a Change's Diff at an object; anything
		// else was built by hand, so show it whole.
		r.line(indent + r.leaf(d))
	}
}

// renderFields writes an object's changed fields in name order, each behind
// its glyph.
func (r *Renderer) renderFields(d *core.Diff, indent string) {
	for _, name := range slices.Sorted(maps.Keys(d.Fields)) {
		label, child := collapse(name, d.Fields[name])
		r.renderNamed(r.glyphFor(child)+" "+label, child, indent)
	}
}

// collapse folds a chain of single-field objects into one dotted label:
// "multiagent.agents:" rather than a line for each level of nesting.
func collapse(label string, d *core.Diff) (string, *core.Diff) {
	for d.Kind == core.DiffObject && len(d.Fields) == 1 {
		for k, v := range d.Fields {
			label, d = label+"."+k, v
		}
	}
	return label, d
}

// renderNamed prints `label: value` on one line for a leaf, or `label:` and an
// indented block for an interior node.
func (r *Renderer) renderNamed(label string, d *core.Diff, indent string) {
	switch d.Kind {
	case core.DiffObject:
		r.line(indent + label + r.paint(dimStyle, ":"))
		r.renderFields(d, indent+indentStep)
	case core.DiffList:
		r.line(indent + label + r.paint(dimStyle, ":"))
		r.renderItems(d, indent+indentStep)
	default:
		// A long value wraps beneath itself, one level in from its key.
		r.wrapLine(indent+label+r.paint(dimStyle, ": "), indent+indentStep, r.leaf(d))
	}
}

// renderItems writes a list's changed elements, each labelled by index:
// "+ [i]" added, "- [i]" removed, "~ [i]" changed in place, "~ [a→b]" moved.
// An element that moved without changing is marked "(moved)".
func (r *Renderer) renderItems(d *core.Diff, indent string) {
	for _, it := range d.Items {
		var label string
		switch {
		case it.Before < 0:
			label = r.paint(createStyle, fmt.Sprintf("+ [%d]", it.After))
		case it.After < 0:
			label = r.paint(destroyStyle, fmt.Sprintf("- [%d]", it.Before))
		case it.Before != it.After:
			label = r.paint(updateStyle, "~") + fmt.Sprintf(" [%d→%d]", it.Before, it.After)
		default:
			label = r.paint(updateStyle, "~") + fmt.Sprintf(" [%d]", it.After)
		}
		if it.Diff == nil {
			r.line(indent + label + " " + r.paint(dimStyle, "(moved)"))
			continue
		}
		label, child := collapse(label, it.Diff)
		r.renderNamed(label, child, indent)
	}
}

// diffGlyph is the marker and colour for a diff node of kind k.
func diffGlyph(k core.DiffKind) (string, styleFunc) {
	switch k {
	case core.DiffAdded:
		return "+", createStyle
	case core.DiffRemoved:
		return "-", destroyStyle
	default:
		return "~", updateStyle
	}
}

// glyphFor is the marker in front of a field name.
func (r *Renderer) glyphFor(d *core.Diff) string {
	glyph, style := diffGlyph(d.Kind)
	return r.paint(style, glyph)
}

// leaf renders a terminal node's value(s) on one line.
func (r *Renderer) leaf(d *core.Diff) string {
	switch d.Kind {
	case core.DiffAdded:
		return r.paint(createStyle, r.value(d.After))
	case core.DiffRemoved:
		return r.paint(destroyStyle, r.value(d.Before)) + r.paint(dimStyle, " → (unset)")
	case core.DiffValue:
		return r.beforeAfter(d.Before, d.After)
	case core.DiffText:
		return r.text(fmt.Sprint(d.Before), fmt.Sprint(d.After))
	case core.DiffSensitive:
		return r.paint(dimStyle, "(sensitive; changed)")
	case core.DiffWriteOnly:
		return r.paint(dimStyle, "(write-only; sent again, cannot be compared)")
	default:
		return r.paint(dimStyle, "(changed)")
	}
}

// text renders a changed string. Short single-line values read best as a plain
// before → after; anything longer gets an inline word diff with the unchanged
// stretches elided, so the edit itself is what the eye lands on. A value the
// plan cannot know yet is not text to compare, so it is shown as the
// placeholder rather than word-diffed against the old value.
func (r *Renderer) text(before, after string) string {
	short := !strings.Contains(before, "\n") && !strings.Contains(after, "\n") &&
		len(before) <= inlineTextLimit && len(after) <= inlineTextLimit
	if short || after == core.KnownAfterApply {
		return r.beforeAfter(before, after)
	}
	runs, ok := diffWords(before, after)
	if !ok {
		return r.paint(dimStyle, fmt.Sprintf("(%d → %d characters of text)", len(before), len(after)))
	}
	if !r.Verbose {
		runs = elideContext(runs, contextWords, func(n int) string {
			return r.paint(dimStyle, fmt.Sprintf("…%d words…", n))
		})
	}
	var b strings.Builder
	for _, run := range runs {
		s := strings.ReplaceAll(run.text, "\n", r.paint(dimStyle, "⏎"))
		switch run.op {
		case wordDel:
			b.WriteString(r.highlight(deletedWordStyle, "[-", s, "-]"))
		case wordIns:
			b.WriteString(r.highlight(insertedWordStyle, "{+", s, "+}"))
		default:
			b.WriteString(s)
		}
	}
	return b.String()
}

// highlight shows an edit: with colour as a background, without it in the
// bracket notation `git diff --word-diff` uses, so a CI log stays legible.
func (r *Renderer) highlight(style styleFunc, opener, s, closer string) string {
	if r.Color {
		return style(s)
	}
	return opener + s + closer
}

// beforeAfter renders a change as `before → after`.
func (r *Renderer) beforeAfter(before, after any) string {
	return r.value(before) + r.paint(dimStyle, " → ") + r.value(after)
}

// value renders any JSON value compactly on one line, eliding the middle of
// anything long unless verbose.
func (r *Renderer) value(v any) string {
	var s string
	switch t := v.(type) {
	case nil:
		return r.paint(dimStyle, "(unset)")
	case string:
		if t == core.KnownAfterApply {
			return r.paint(dimStyle, t)
		}
		s = quote(t)
	default:
		b, err := json.Marshal(v)
		if err != nil {
			s = fmt.Sprint(v)
		} else {
			s = string(b)
		}
	}
	if !r.Verbose {
		s = elide(s, valueWidth)
	}
	return s
}

func quote(s string) string {
	b, _ := json.Marshal(s)
	return string(b)
}

// elide keeps the head and tail of a long value: an ID's prefix and its last
// characters are both worth seeing.
func elide(s string, width int) string {
	length := ansi.StringWidth(s)
	if length <= width {
		return s
	}
	head := (width - 1) * 2 / 3
	tail := width - 1 - head
	return ansi.Truncate(s, head, "") + "…" + ansi.TruncateLeft(s, length-tail, "")
}
