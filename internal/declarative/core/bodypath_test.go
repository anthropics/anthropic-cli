package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSplitPath(t *testing.T) {
	assert.Nil(t, splitPath(""))
	assert.Equal(t, []string{"a"}, splitPath("a"))
	assert.Equal(t, []string{"a", "b", "c"}, splitPath("a.b.c"))
}

func TestGetPath(t *testing.T) {
	body := map[string]any{
		"name":  "x",
		"empty": nil,
		"outer": map[string]any{"inner": map[string]any{"leaf": 7}},
	}

	got, ok := getPath(body, "outer.inner.leaf")
	assert.True(t, ok)
	assert.Equal(t, 7, got)

	// A key that is present with a null value is found, not missing.
	got, ok = getPath(body, "empty")
	assert.True(t, ok)
	assert.Nil(t, got)

	_, ok = getPath(body, "outer.missing")
	assert.False(t, ok)

	// Walking through a value that is not a map stops there.
	_, ok = getPath(body, "name.deeper")
	assert.False(t, ok)

	// The empty path is the body itself.
	got, ok = getPath(body, "")
	assert.True(t, ok)
	assert.Equal(t, body, got)
}

func TestSetPath(t *testing.T) {

	t.Run("creates the maps along the way", func(t *testing.T) {
		body := map[string]any{}
		setPath(body, "a.b.c", 1)
		assert.Equal(t, map[string]any{"a": map[string]any{"b": map[string]any{"c": 1}}}, body)
	})

	t.Run("overwrites a leaf and keeps its siblings", func(t *testing.T) {
		body := map[string]any{"a": map[string]any{"b": 1, "keep": true}}
		setPath(body, "a.b", 2)
		assert.Equal(t, map[string]any{"a": map[string]any{"b": 2, "keep": true}}, body)
	})

	t.Run("does nothing when a segment on the way is not a map", func(t *testing.T) {
		body := map[string]any{"a": "text"}
		setPath(body, "a.b", 1)
		assert.Equal(t, map[string]any{"a": "text"}, body)
	})

	t.Run("does nothing for the empty path", func(t *testing.T) {
		body := map[string]any{"a": 1}
		setPath(body, "", 2)
		assert.Equal(t, map[string]any{"a": 1}, body)
	})
}

func TestDeletePath(t *testing.T) {

	body := map[string]any{"a": map[string]any{"b": 1, "keep": 2}, "text": "x"}
	deletePath(body, "a.b")
	assert.Equal(t, map[string]any{"a": map[string]any{"keep": 2}, "text": "x"}, body)

	// Missing keys, non-map segments and the empty path leave the body alone.
	deletePath(body, "missing.b")
	deletePath(body, "text.b")
	deletePath(body, "")
	assert.Equal(t, map[string]any{"a": map[string]any{"keep": 2}, "text": "x"}, body)
}

func TestDeleteWildcardPath(t *testing.T) {

	t.Run("removes the field from every element", func(t *testing.T) {
		body := map[string]any{"resources": []any{
			map[string]any{"name": "a", "token": "s1"},
			map[string]any{"name": "b", "token": "s2"},
		}}
		deleteWildcardPath(body, "resources[].token")
		assert.Equal(t, map[string]any{"resources": []any{
			map[string]any{"name": "a"},
			map[string]any{"name": "b"},
		}}, body)
	})

	t.Run("handles an array inside an array", func(t *testing.T) {
		body := map[string]any{"a": []any{
			map[string]any{"b": []any{map[string]any{"c": 1, "keep": 2}}},
		}}
		deleteWildcardPath(body, "a[].b[].c")
		assert.Equal(t, map[string]any{"a": []any{
			map[string]any{"b": []any{map[string]any{"keep": 2}}},
		}}, body)
	})

	t.Run("without a wildcard it deletes the plain path", func(t *testing.T) {
		body := map[string]any{"a": map[string]any{"b": 1, "keep": 2}}
		deleteWildcardPath(body, "a.b")
		assert.Equal(t, map[string]any{"a": map[string]any{"keep": 2}}, body)
	})

	t.Run("skips a missing array, a non-array and non-map elements", func(t *testing.T) {
		body := map[string]any{
			"notList": map[string]any{"token": "s"},
			"mixed":   []any{"text", 3, map[string]any{"token": "s"}},
		}
		deleteWildcardPath(body, "missing[].token")
		deleteWildcardPath(body, "notList[].token")
		deleteWildcardPath(body, "mixed[].token")
		assert.Equal(t, map[string]any{
			"notList": map[string]any{"token": "s"},
			"mixed":   []any{"text", 3, map[string]any{}},
		}, body)
	})
}

func TestDeepCopyDoesNotShareMapsOrSlices(t *testing.T) {
	src := map[string]any{
		"list": []any{map[string]any{"k": "v"}, "x"},
		"map":  map[string]any{"n": 1},
	}
	cp := deepCopyMap(src)
	require.Equal(t, src, cp)

	cp["list"].([]any)[0].(map[string]any)["k"] = "changed"
	cp["list"].([]any)[1] = "changed"
	cp["map"].(map[string]any)["n"] = 2
	cp["added"] = true

	assert.Equal(t, map[string]any{
		"list": []any{map[string]any{"k": "v"}, "x"},
		"map":  map[string]any{"n": 1},
	}, src)
}

func TestDeepCopyMapIsNeverNil(t *testing.T) {
	cp := deepCopyMap(nil)
	require.NotNil(t, cp)
	assert.Empty(t, cp)
}
