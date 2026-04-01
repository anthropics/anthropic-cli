package jsonview

import (
	"testing"

	"github.com/charmbracelet/bubbles/help"
	"github.com/tidwall/gjson"
)

func TestNavigateForward_EmptyRowData(t *testing.T) {
	// An empty JSON array produces a TableView with no rows.
	emptyArray := gjson.Parse("[]")
	view, err := newTableView("", emptyArray, false)
	if err != nil {
		t.Fatalf("newTableView: %v", err)
	}

	viewer := &JSONViewer{
		stack: []JSONView{view},
		root:  "test",
		help:  help.New(),
	}

	// Should return without panicking despite the empty data set.
	model, cmd := viewer.navigateForward()
	if model != viewer {
		t.Error("expected same viewer model returned")
	}
	if cmd != nil {
		t.Error("expected nil cmd")
	}

	// Stack should remain unchanged (no new view pushed).
	if len(viewer.stack) != 1 {
		t.Errorf("expected stack length 1, got %d", len(viewer.stack))
	}
}
