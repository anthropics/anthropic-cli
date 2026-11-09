package jsonview

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/anthropics/anthropic-cli/v2/internal/termsafe"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/tidwall/gjson"

	"github.com/stretchr/testify/require"
)

func TestNavigateForward_EmptyRowData(t *testing.T) {
	t.Parallel()

	// An empty JSON array produces a TableView with no rows.
	emptyArray := gjson.Parse("[]")
	view, err := newTableView("", emptyArray, false)
	require.NoError(t, err)

	viewer := &JSONViewer{
		stack: []JSONView{view},
		root:  "test",
		help:  help.New(),
	}

	// Should return without panicking despite the empty data set.
	model, cmd := viewer.navigateForward()
	require.Equal(t, model, viewer, "expected same viewer model returned")
	require.Nil(t, cmd)

	// Stack should remain unchanged (no new view pushed).
	require.Equal(t, 1, len(viewer.stack), "expected stack length 1, got %d", len(viewer.stack))
}

func TestPrintValue_EmptyRowData(t *testing.T) {
	t.Parallel()

	view, err := newTableView("", gjson.Parse("[]"), false)
	require.NoError(t, err)
	viewer := &JSONViewer{stack: []JSONView{view}, root: "test", help: help.New()}

	viewer.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'p'}})

	require.Empty(t, viewer.message)
}

// JSON-encoded control sequences that gjson decodes to live ones: an OSC 52 clipboard write and a CSI color change.
const (
	jsonOSC52 = `\u001b]52;c;aGVsbG8=\u0007`
	jsonCSI   = `\u001b[31m`
)

// safeJSON parses body as the output formats do before they render it.
func safeJSON(body string) gjson.Result {
	return gjson.Parse(termsafe.JSON(body))
}

func TestNewTableView_RendersSafeData(t *testing.T) {
	t.Parallel()

	t.Run("ArrayOfObjects", func(t *testing.T) {
		t.Parallel()

		data := safeJSON(`[{"na` + jsonOSC52 + `me": "val` + jsonCSI + `ue"}]`)

		view, err := newTableView("", data, false)

		require.NoError(t, err)
		require.Equal(t, `na\x1b]52;c;aGVsbG8=\ame`, view.table.Columns()[0].Title)
		require.Equal(t, `val\x1b[31mue`, view.table.Rows()[0][0])
	})

	t.Run("Object", func(t *testing.T) {
		t.Parallel()

		data := safeJSON(`{"na` + jsonOSC52 + `me": {"in` + jsonCSI + `ner": "val` + jsonCSI + `ue"}}`)

		view, err := newTableView("", data, false)

		require.NoError(t, err)
		row := view.table.Rows()[0]
		require.Equal(t, `na\x1b]52;c;aGVsbG8=\ame`, row[0])
		require.Equal(t, `{in\x1b[31mner:"val\x1b[31mue"}`, row[1])
	})

	t.Run("LongPreviewIsCutAfterEscaping", func(t *testing.T) {
		t.Parallel()

		data := safeJSON(`{"key": {"inner": "` + strings.Repeat("x", maxPreviewLength-4) + jsonOSC52 + `"}}`)

		view, err := newTableView("", data, false)

		require.NoError(t, err)
		require.Equal(t, `{inner:"`+strings.Repeat("x", maxPreviewLength-4)+`\x1b…"}`, view.table.Rows()[0][1])
	})

	t.Run("RawMode", func(t *testing.T) {
		t.Parallel()

		data := safeJSON("[\"val\u009b31mue\"]")

		view, err := newTableView("", data, true)

		require.NoError(t, err)
		require.Equal(t, `"val\\u009b31mue"`, view.table.Rows()[0][0])
	})

	t.Run("NonASCIIIsUntouched", func(t *testing.T) {
		t.Parallel()

		view, err := newTableView("", safeJSON("[\"héllo 世界 👨\u200d👩\u200d👧\"]"), false)

		require.NoError(t, err)
		require.Equal(t, "héllo 世界 👨\u200d👩\u200d👧", view.table.Rows()[0][0])
	})
}

func TestLoadMoreData_FindsCellsOfEscapedKeys(t *testing.T) {
	t.Parallel()

	data := safeJSON(`[{"id": 1, "na` + jsonOSC52 + `me": "first"}]`)
	view, err := newTableView("", data, false)
	require.NoError(t, err)
	view.iterator = &sliceIterator{items: []any{
		json.RawMessage(termsafe.JSON(`{"id": 2, "na` + jsonOSC52 + `me": "second` + jsonCSI + `"}`)),
	}}

	msg := view.loadMoreData(false)()

	require.Nil(t, msg)
	require.Equal(t, table.Row{"2", `second\x1b[31m`}, view.table.Rows()[1])
}

func TestLoadMoreData_KeepsRawJSONAsSent(t *testing.T) {
	t.Parallel()

	view, err := newTableView("", safeJSON(`["first"]`), true)
	require.NoError(t, err)
	view.iterator = &sliceIterator{items: []any{rawJSONItem{raw: `"a&b"`}}}

	msg := view.loadMoreData(true)()

	require.Nil(t, msg)
	require.Equal(t, table.Row{`"a&b"`}, view.table.Rows()[1])
}

func TestNewTableView_KeepsCellsOnOneLine(t *testing.T) {
	t.Parallel()

	view, err := newTableView("", safeJSON(`{"key\none": "one\ntwo\tthree"}`), false)

	require.NoError(t, err)
	require.Equal(t, table.Row{`key\none`, `one\ntwo\tthree`}, view.table.Rows()[0])
}

func TestNavigateForward_EscapesTitleAndTextView(t *testing.T) {
	t.Parallel()

	long := strings.Repeat("x", maxStringLength)
	data := safeJSON(`{"k` + jsonOSC52 + `ey": "` + long + jsonCSI + `\nline two"}`)
	view, err := newTableView("", data, false)
	require.NoError(t, err)
	viewer := &JSONViewer{stack: []JSONView{view}, root: "test", help: help.New()}

	viewer.navigateForward()

	require.Contains(t, viewer.buildTitle(viewer.current()), `k\\x1b]52;c;aGVsbG8=\\aey`)
	textView, ok := viewer.current().(*TextView)
	require.True(t, ok)
	textView.Resize(2*maxStringLength, 50)
	require.Contains(t, textView.View(), long+`\x1b[31m`)
	require.Contains(t, textView.View(), "\nline two")
}

func TestNavigateForward_MeasuresEscapedLength(t *testing.T) {
	t.Parallel()

	data := safeJSON(`["` + strings.Repeat(`\u0007`, maxStringLength/2) + `"]`)
	view, err := newTableView("", data, false)
	require.NoError(t, err)
	viewer := &JSONViewer{stack: []JSONView{view}, root: "test", help: help.New()}

	viewer.navigateForward()

	require.Len(t, viewer.stack, 2)
}

func TestPrintValue_EscapesControlSequencesAndKeepsNewlines(t *testing.T) {
	t.Parallel()

	view, err := newTableView("", safeJSON(`["value`+jsonOSC52+`\nline two"]`), false)
	require.NoError(t, err)
	viewer := &JSONViewer{stack: []JSONView{view}, root: "test", help: help.New()}

	viewer.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'p'}})

	require.Equal(t, `value\x1b]52;c;aGVsbG8=\a`+"\nline two", viewer.message)
}

func TestTableView_MovesCursorInShortWindow(t *testing.T) {
	t.Parallel()

	for name, size := range map[string]tea.WindowSizeMsg{
		"NoSize":  {},
		"SixRows": {Width: 80, Height: 6},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			view, err := newTableView("", safeJSON(`[1, 2, 3]`), false)
			require.NoError(t, err)
			viewer := &JSONViewer{stack: []JSONView{view}, root: "test", help: help.New()}
			viewer.Update(size)

			viewer.Update(tea.KeyMsg{Type: tea.KeyDown})
			viewer.Update(tea.KeyMsg{Type: tea.KeyDown})
			viewer.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'p'}})

			require.Equal(t, "3", viewer.message)
		})
	}
}

func TestTextView_ScrollsInShortWindow(t *testing.T) {
	t.Parallel()

	for name, size := range map[string]tea.WindowSizeMsg{
		"NoSize":   {},
		"FourRows": {Width: 80, Height: 4},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			view, err := newTableView("", safeJSON(`["one\ntwo\nthree"]`), false)
			require.NoError(t, err)
			viewer := &JSONViewer{stack: []JSONView{view}, root: "test", help: help.New()}
			viewer.Update(size)
			viewer.navigateForward()

			require.NotPanics(t, func() {
				viewer.Update(tea.KeyMsg{Type: tea.KeyDown})
				viewer.Update(tea.KeyMsg{Type: tea.KeyPgUp})
			})
		})
	}
}

type sliceIterator struct {
	index int
	items []any
}

func (it *sliceIterator) Next() bool {
	it.index++
	return it.index <= len(it.items)
}

func (it *sliceIterator) Current() any { return it.items[it.index-1] }
func (it *sliceIterator) Err() error   { return nil }

// rawJSONItem implements HasRawJSON, returning pre-built JSON.
type rawJSONItem struct {
	raw string
}

func (r rawJSONItem) RawJSON() string { return r.raw }

func TestMarshalItemsToJSONArray_WithHasRawJSON(t *testing.T) {
	t.Parallel()

	items := []any{
		rawJSONItem{raw: `{"id":1,"name":"alice"}`},
		rawJSONItem{raw: `{"id":2,"name":"bob"}`},
	}

	got, err := marshalItemsToJSONArray(items)
	require.NoError(t, err)
	require.JSONEq(t, `[{"id":1,"name":"alice"},{"id":2,"name":"bob"}]`, string(got))
}

func TestMarshalItemsToJSONArray_WithoutHasRawJSON(t *testing.T) {
	t.Parallel()

	items := []any{
		map[string]any{"id": 1, "name": "alice"},
		map[string]any{"id": 2, "name": "bob"},
	}

	got, err := marshalItemsToJSONArray(items)
	require.NoError(t, err)
	require.JSONEq(t, `[{"id":1,"name":"alice"},{"id":2,"name":"bob"}]`, string(got))
}
