package core

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLockfileKeyOutsideRootIsNotPruned(t *testing.T) {
	root := t.TempDir()
	lockPath := filepath.Join(root, testLockfileName)
	require.NoError(t, os.WriteFile(lockPath, []byte(`{
  "version": 1,
  "resources": {
    "../missing": {
      "kind": "gadget",
      "id": "victim-id",
      "hash": "poc"
    }
  }
}
`), 0o644))

	_, err := LoadLockfile(testRegistry(), lockPath)
	require.Error(t, err)
	require.ErrorContains(t, err, "outside the lockfile directory")

	lf := newLockfile(lockPath)
	lf.Resources["../missing"] = &LockEntry{Kind: "gadget", ID: "victim-id", Hash: "poc"}
	changes := (&Planner{Registry: testRegistry(), Lock: lf, Prune: true}).planDestroys(nil)
	require.Empty(t, changes)
}
