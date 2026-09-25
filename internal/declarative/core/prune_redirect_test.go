package core

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// recordingClient implements Client and records Get/Destroy calls so the PoC
// can show planDestroys never fetches, while Apply destroys whatever id the
// lockfile currently holds.
type recordingClient struct {
	gets     []destroyCall
	destroys []destroyCall
}

type destroyCall struct {
	Kind Kind
	ID   string
}

func (c *recordingClient) Get(_ context.Context, kind Kind, id string) (map[string]any, error) {
	c.gets = append(c.gets, destroyCall{Kind: kind, ID: id})
	return map[string]any{"name": "victim"}, nil
}

func (c *recordingClient) Create(context.Context, Kind, Request) (map[string]any, error) {
	return nil, nil
}

func (c *recordingClient) Update(context.Context, Kind, string, Request) (map[string]any, error) {
	return nil, nil
}

func (c *recordingClient) Destroy(_ context.Context, kind Kind, id string) error {
	c.destroys = append(c.destroys, destroyCall{Kind: kind, ID: id})
	return nil
}

// TestPocPruneRedirect proves HACK-ANTCLI-02: checkLockEntry accepts any
// nonempty id for a valid kind, and planDestroys (with Prune) copies that
// lockfile id into ActionDestroy with no fetch and no comparison to a
// previously recorded id. Apply then calls Client.Destroy(kind, victimID).
func TestPocPruneRedirect(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, testLockfileName)

	require.NoError(t, os.WriteFile(path, []byte(`{
  "version": 1,
  "resources": {
    "./gadgets/a.md": {"kind": "gadget", "id": "gdgt_ORIGINAL", "hash": "deadbeef"}
  }
}
`), 0o644))

	lf, err := LoadLockfile(testRegistry(), path)
	require.NoError(t, err)
	require.Contains(t, lf.Resources, "./gadgets/a.md")
	assert.Equal(t, "gdgt_ORIGINAL", lf.Resources["./gadgets/a.md"].ID)

	// Contributor edits the orphan lock entry's id from A to same-kind id B.
	spec, ok := testRegistry().Spec("gadget")
	require.True(t, ok)
	recorded, err := hashBody(normalizeRemote(spec, map[string]any{"name": "original"}))
	require.NoError(t, err)
	lf.Resources["./gadgets/a.md"].RemoteHash = recorded
	lf.Resources["./gadgets/a.md"].ID = "gdgt_VICTIM"
	require.NoError(t, lf.Save())

	reloaded, err := LoadLockfile(testRegistry(), path)
	require.NoError(t, err, "checkLockEntry must accept any nonempty id for a valid kind")
	assert.Equal(t, "gdgt_VICTIM", reloaded.Resources["./gadgets/a.md"].ID)

	// Do not create ./gadgets/a.md — the key is an orphan.
	loader := NewLoader(testRegistry(), root, nil)
	require.NoError(t, loader.AddKeys(context.Background(), reloaded.Keys()))
	assert.Empty(t, loader.Sources(), "missing file must be skipped by AddKeys")

	client := &recordingClient{}
	planner := Planner{
		Registry: testRegistry(),
		Lock:     reloaded,
		Prune:    true,
		Client:   client,
	}
	plan, err := planner.Plan(context.Background(), loader)
	require.NoError(t, err)

	require.Len(t, plan.Changes, 1)
	change := plan.Changes[0]
	assert.Equal(t, ActionDestroy, change.Action)
	require.NotNil(t, change.Entry)
	assert.Equal(t, Kind("gadget"), change.Kind)
	assert.Equal(t, "gdgt_VICTIM", change.Entry.ID)
	require.Error(t, change.Blocked)
	require.ErrorContains(t, change.Blocked, "not the resource the last apply recorded")
	require.Len(t, client.gets, 1)
	assert.Equal(t, destroyCall{Kind: "gadget", ID: "gdgt_VICTIM"}, client.gets[0])

	applier := Applier{Client: client, Lock: reloaded}
	_, err = applier.Apply(context.Background(), plan)
	require.Error(t, err)
	assert.Empty(t, client.destroys, "Apply must not Destroy an id that does not match the recorded remote hash")
}
