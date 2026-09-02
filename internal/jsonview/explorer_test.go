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

func TestGetSelectedContent_EmptyRowData(t *testing.T) {
	t.Parallel()

	// Empty containers produce TableViews with no rows. Pressing "p"
	// (print and exit) must not panic on them — mirror of the guard
	// exercised by TestNavigateForward_EmptyRowData.
	for _, raw := range []string{"[]", "{}"} {
		view, err := newTableView("", gjson.Parse(raw), false)
		require.NoError(t, err)

		viewer := &JSONViewer{
			stack: []JSONView{view},
			root:  "test",
			help:  help.New(),
		}

		require.NotPanics(t, func() {
			// With nothing to select, "p" falls back to printing the
			// container itself.
			require.Equal(t, raw, viewer.getSelectedContent())
		}, "getSelectedContent must not panic on %s", raw)
	}
}

// exhaustedIterator reports how many times Next is called after the data ran out.
type exhaustedIterator struct {
	nextCalls int
}

func (e *exhaustedIterator) Next() bool   { e.nextCalls++; return false }
func (e *exhaustedIterator) Current() any { return nil }
func (e *exhaustedIterator) Err() error   { return nil }

func TestLoadMoreData_DetachesExhaustedIterator(t *testing.T) {
	t.Parallel()

	view, err := newTableView("", gjson.Parse(`[{"id":1}]`), false)
	require.NoError(t, err)
	it := &exhaustedIterator{}
	view.iterator = it

	cmd := view.loadMoreData(false)
	require.NotNil(t, cmd)
	_ = cmd()

	// Once the iterator reports no more data, the view must stop polling it
	// on every subsequent cursor move at the last row.
	require.Nil(t, view.iterator, "exhausted iterator should be detached")
	require.Equal(t, 1, it.nextCalls)
	require.False(t, view.isLoading)
}

func TestStreamPreloadCount_FallbackWithoutTerminal(t *testing.T) {
	t.Parallel()

	// Test processes have no tty on stdout, so the terminal-size probe fails
	// and the fallback must be used.
	require.Equal(t, 20, streamPreloadCount())
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
