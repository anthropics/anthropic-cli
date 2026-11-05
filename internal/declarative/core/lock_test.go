package core

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLockfileRoundTrips(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, testLockfileName)

	lf := newLockfile(path)
	lf.Origin = &Origin{BaseURL: "https://api.example.test", OrganizationID: "org-1"}
	lf.Resources["./gadgets/a.md"] = &LockEntry{
		Kind: "gadget", ID: "gdgt_01", Version: "4", Hash: "aaaa", RemoteHash: "bbbb",
	}
	lf.Resources["./widgets/w.yml"] = &LockEntry{
		Kind: "widget", ID: "wdgt_02", Hash: "dddd", RemoteHash: "eeee", Stamp: "2026-01-01T00:00:00Z",
	}
	lf.Resources["https://github.com/o/r/tree/main/s"] = &LockEntry{
		Kind: "widget", ID: "wdgt_01", Version: "1759178010641129", Hash: "cccc",
		Revision: "0123456789abcdef0123456789abcdef01234567", Subpath: "s",
	}
	require.NoError(t, lf.Save())

	reloaded, err := LoadLockfile(testRegistry(), path)
	require.NoError(t, err)
	assert.True(t, reloaded.Existed())
	assert.Equal(t, lf.Resources, reloaded.Resources)
	assert.Equal(t, lf.Origin, reloaded.Origin)

	// An epoch-string version must not come back as a number.
	assert.Equal(t, "1759178010641129", reloaded.Resources["https://github.com/o/r/tree/main/s"].Version)
}

func TestLockfileOutputIsStable(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, testLockfileName)

	build := func() *Lockfile {
		lf := newLockfile(path)
		for _, key := range []string{"./z.yml", "./a.yml", "./m.yml"} {
			lf.Resources[key] = &LockEntry{Kind: "widget", ID: "wdgt_" + key, Hash: "h"}
		}
		return lf
	}
	require.NoError(t, build().Save())
	first, err := os.ReadFile(path)
	require.NoError(t, err)

	require.NoError(t, build().Save())
	second, err := os.ReadFile(path)
	require.NoError(t, err)

	assert.Equal(t, string(first), string(second),
		"the lockfile must be byte-identical run to run, or every apply produces a diff")
}

// A field this build does not know about was written by a newer one a
// colleague ran. Dropping it on save would silently strip their state.
func TestLockfileKeepsFieldsItDoesNotUnderstand(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, testLockfileName)
	require.NoError(t, os.WriteFile(path, []byte(`{
  "version": 1,
  "resources": {
    "./widgets/w.yml": {"kind": "widget", "id": "wdgt_01", "hash": "h", "from_the_future": {"nested": [1, 2]}}
  }
}`), 0o644))

	lf, err := LoadLockfile(testRegistry(), path)
	require.NoError(t, err)
	require.NoError(t, lf.Save())

	data, err := os.ReadFile(path)
	require.NoError(t, err)
	assert.Contains(t, string(data), `"from_the_future": {`)
	assert.Contains(t, string(data), `"nested": [`)
}

func TestLockfileRejectsGarbage(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, testLockfileName)

	tests := []struct {
		name    string
		content string
		wantErr string
	}{
		{"not json", "lockfile:\n  version: 1\n", "invalid character"},
		{"unexpected key", `{"resources": {"whatever": {"kind": "widget", "id": "x"}}}`, "resource keys start with"},
		{"missing id", `{"resources": {"./a.md": {"kind": "gadget"}}}`, "missing `id`"},
		{"unknown kind", `{"resources": {"./a.md": {"kind": "sprocket", "id": "x"}}}`, "unknown kind"},
		{"future schema", `{"version": 99}`, "newer ant"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			require.NoError(t, os.WriteFile(path, []byte(tc.content), 0o644))
			_, err := LoadLockfile(testRegistry(), path)
			require.ErrorContains(t, err, tc.wantErr)
		})
	}
}

func TestMissingLockfileIsNotAnError(t *testing.T) {
	lf, err := LoadLockfile(testRegistry(), filepath.Join(t.TempDir(), testLockfileName))
	require.NoError(t, err)
	assert.False(t, lf.Existed())
	assert.Empty(t, lf.Resources)
}

func TestFindLockfileWalksUpButStopsAtTheRepoRoot(t *testing.T) {
	root := t.TempDir()
	require.NoError(t, os.MkdirAll(filepath.Join(root, "repo", ".git"), 0o755))
	require.NoError(t, os.MkdirAll(filepath.Join(root, "repo", "a", "b"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(root, testLockfileName), []byte(""), 0o644))

	// A lockfile above the repository belongs to somebody else's project.
	lf, err := FindLockfile(testRegistry(), filepath.Join(root, "repo", "a", "b"))
	require.NoError(t, err)
	assert.False(t, lf.Existed())

	// One inside the repository is found from a nested directory.
	inner := filepath.Join(root, "repo", testLockfileName)
	require.NoError(t, os.WriteFile(inner, []byte(""), 0o644))
	lf, err = FindLockfile(testRegistry(), filepath.Join(root, "repo", "a", "b"))
	require.NoError(t, err)
	assert.True(t, lf.Existed())
	assert.Equal(t, inner, lf.Path)
}

func TestSaveIsAtomic(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, testLockfileName)
	lf := newLockfile(path)
	lf.Resources["./a.md"] = &LockEntry{Kind: "gadget", ID: "gdgt_01", Hash: "h"}
	require.NoError(t, lf.Save())

	// No temp files left behind.
	entries, err := os.ReadDir(dir)
	require.NoError(t, err)
	require.Len(t, entries, 1)
	assert.Equal(t, testLockfileName, entries[0].Name())
}
