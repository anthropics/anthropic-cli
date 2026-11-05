package termbackground_test

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestImportPathSortsBeforeBubbleTea(t *testing.T) {
	require.Less(t, thisPackagePath, bubbleTeaPath,
		"Go would initialize Bubble Tea before this package's init in termbackground.go, so every command would ask the "+
			"terminal for its background color at start-up, and the tests that catch that would skip.")
}
