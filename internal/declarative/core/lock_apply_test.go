package core

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// A domain-neutral backend keeps the loader, planner, applier and on-disk
// lockfile real while recording exactly what the reconciler sends.
type lockUpdateClient struct {
	remote    map[string]any
	writes    []Request
	updateErr error
}

func (c *lockUpdateClient) Get(context.Context, Kind, string) (map[string]any, error) {
	return deepCopyMap(c.remote), nil
}
func (c *lockUpdateClient) Create(context.Context, Kind, Request) (map[string]any, error) {
	return nil, errors.New("unexpected create")
}
func (c *lockUpdateClient) Destroy(context.Context, Kind, string) error {
	return errors.New("unexpected destroy")
}
func (c *lockUpdateClient) Update(_ context.Context, _ Kind, id string, req Request) (map[string]any, error) {
	c.writes = append(c.writes, req)
	if c.updateErr != nil {
		return nil, c.updateErr
	}
	c.remote = deepCopyMap(req.Body)
	c.remote["id"] = id
	return deepCopyMap(c.remote), nil
}

func TestApplyUpdatesPreserveUnknownLockFields(t *testing.T) {
	for _, failWrite := range []bool{false, true} {
		t.Run(fmt.Sprintf("fail_write=%t", failWrite), func(t *testing.T) {
			registry := testRegistry()
			root := writeTree(t, map[string]string{"widgets/a.json": `{"name":"updated"}`})
			oldBody := map[string]any{"name": "original"}
			oldHash, err := hashBody(oldBody)
			require.NoError(t, err)
			lockPath := filepath.Join(root, testLockfileName)
			const fields = `"future_data":{"nested":[9007199254740993,null,{"enabled":false}]},"future_flag":false,"future_null":null,"future_text":"keep ☃"`
			raw := fmt.Sprintf(`{"version":1,"resources":{"./widgets/a.json":{"kind":"widget","id":"wdgt_01","hash":%q,%s}}}`, oldHash, fields)
			require.NoError(t, os.WriteFile(lockPath, []byte(raw), 0600))
			lock, err := LoadLockfile(registry, lockPath)
			require.NoError(t, err)
			oldEntry := lock.Resources["./widgets/a.json"]
			require.NoError(t, lock.Save())
			client := &lockUpdateClient{remote: map[string]any{"id": "wdgt_01", "name": "original"}}
			if failWrite {
				client.updateErr = errors.New("write rejected")
			}
			for _, name := range []string{"updated", "updated again"} {
				source := []byte(fmt.Sprintf(`{"name":%q}`, name))
				path := filepath.Join(root, "widgets", "a.json")
				require.NoError(t, os.WriteFile(path, source, 0600))
				loader := NewLoader(registry, root, nil)
				require.NoError(t, loader.Add(context.Background(), []string{path}))
				plan, err := (&Planner{Registry: registry, Client: client, Lock: lock}).Plan(context.Background(), loader)
				require.NoError(t, err)
				require.Empty(t, plan.Blocked())
				require.Len(t, plan.Changes, 1)
				require.Equal(t, ActionUpdate, plan.Changes[0].Action)
				result, err := (&Applier{Client: client, Lock: lock}).Apply(context.Background(), plan)
				if failWrite {
					require.ErrorContains(t, err, "write rejected")
					assert.Zero(t, result.Updated)
				} else {
					require.NoError(t, err)
					assert.Equal(t, 1, result.Updated)
				}
				reloaded, readErr := LoadLockfile(registry, lockPath)
				require.NoError(t, readErr)
				entry := reloaded.Resources["./widgets/a.json"]
				assert.Equal(t, len(oldEntry.unknown), len(entry.unknown))
				for key, expected := range oldEntry.unknown {
					actual, exists := entry.unknown[key]
					if assert.True(t, exists, "unknown field %q was discarded", key) {
						var want, got any
						wantDecoder := json.NewDecoder(bytes.NewReader(expected))
						gotDecoder := json.NewDecoder(bytes.NewReader(actual))
						wantDecoder.UseNumber()
						gotDecoder.UseNumber()
						require.NoError(t, wantDecoder.Decode(&want))
						require.NoError(t, gotDecoder.Decode(&got))
						assert.Equal(t, want, got, "unknown field %q changed", key)
					}
				}
				if !failWrite {
					expectedHash, hashErr := hashBody(map[string]any{"name": name})
					require.NoError(t, hashErr)
					assert.Equal(t, expectedHash, entry.Hash)
					assert.NotEmpty(t, entry.RemoteHash)
					assert.Equal(t, "wdgt_01", entry.ID)
					plan, planErr := (&Planner{Registry: registry, Client: client, Lock: reloaded}).Plan(context.Background(), loader)
					require.NoError(t, planErr)
					assert.False(t, plan.HasWork())
				}
				// Future bookkeeping is local-only; no unknown lock fields enter a request.
				assert.Equal(t, map[string]any{"name": name}, client.writes[len(client.writes)-1].Body)
				after, readErr := os.ReadFile(path)
				require.NoError(t, readErr)
				assert.Equal(t, source, after)
				lock = reloaded
			}
		})
	}
}

func TestRecordDoesNotTransferUnknownMetadataToReplacementResource(t *testing.T) {
	lock := newLockfile(filepath.Join(t.TempDir(), testLockfileName))
	old := &LockEntry{Kind: "widget", ID: "wdgt_old", unknown: map[string]json.RawMessage{"future": json.RawMessage(`{"a":1}`)}}
	lock.Resources["./widgets/a.json"] = old
	change := &Change{Key: "./widgets/a.json", Kind: "widget", Action: ActionCreate, Entry: nil, Replaces: old.ID}
	(&Applier{Lock: lock}).record(change, "wdgt_new", "", "hash", nil)
	assert.Empty(t, lock.Resources[change.Key].unknown)
	assert.Equal(t, json.RawMessage(`{"a":1}`), old.unknown["future"])
}

func TestRecordedFutureMetadataHasAnIndependentKeyMap(t *testing.T) {
	lock := newLockfile(filepath.Join(t.TempDir(), testLockfileName))
	old := &LockEntry{Kind: "widget", ID: "wdgt_01", unknown: map[string]json.RawMessage{"future": json.RawMessage(`1`)}}
	change := &Change{Key: "./widgets/a.json", Kind: "widget", Action: ActionUpdate, Entry: old}
	(&Applier{Lock: lock}).record(change, old.ID, "", "hash", nil)
	current := lock.Resources[change.Key]
	require.Equal(t, old.unknown, current.unknown)
	current.unknown["new_future_field"] = json.RawMessage(`true`)
	assert.NotContains(t, old.unknown, "new_future_field")
}
