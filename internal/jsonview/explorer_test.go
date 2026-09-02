package jsonview

import (
	"testing"

	"github.com/charmbracelet/bubbles/help"
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

func TestNewArrayOfObjectsTableView_DottedKeys(t *testing.T) {
	t.Parallel()

	// Keys containing "." or "\" are treated as gjson paths when looked up,
	// so column values must be looked up with the key escaped.
	array := gjson.Parse(`[{"a.b":"dotted","a\\b":"backslash","plain":"x"}]`).Array()
	view := newArrayOfObjectsTableView("", gjson.Parse(`[{"a.b":"dotted","a\\b":"backslash","plain":"x"}]`), array, false)

	row := view.table.Rows()[0]
	require.Len(t, row, 3)
	require.Equal(t, "dotted", row[0])
	require.Equal(t, "backslash", row[1])
	require.Equal(t, "x", row[2])
}

func TestNewObjectTableView_DottedKey(t *testing.T) {
	t.Parallel()

	data := gjson.Parse(`{"a.b":1,"plain":2}`)
	view := newObjectTableView("", data, false)

	require.Len(t, view.rowData, 2)
	require.Equal(t, "1", view.rowData[0].Raw)
	require.Equal(t, "2", view.rowData[1].Raw)
	require.Equal(t, "1", view.table.Rows()[0][1])
	require.Equal(t, "2", view.table.Rows()[1][1])
}

func TestFormatObject_DottedKey(t *testing.T) {
	t.Parallel()

	obj := gjson.Parse(`{"a.b":"v","plain":"w"}`)
	require.Equal(t, `{a.b:"v", plain:"w"}`, formatObject(obj))
}
