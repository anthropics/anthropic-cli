package core

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type cancellationClient struct {
	failAction Action
	failure    error
	calls      []Action
	onFailure  func()
}

func (c *cancellationClient) Get(context.Context, Kind, string) (map[string]any, error) {
	return nil, ErrNotFound
}
func (c *cancellationClient) operation(action Action, body map[string]any) (map[string]any, error) {
	c.calls = append(c.calls, action)
	if len(c.calls) == 2 && action == c.failAction {
		if c.onFailure != nil {
			c.onFailure()
		}
		return nil, c.failure
	}
	copy := deepCopyMap(body)
	copy["id"] = "wdgt_created"
	return copy, nil
}
func (c *cancellationClient) Create(_ context.Context, _ Kind, req Request) (map[string]any, error) {
	return c.operation(ActionCreate, req.Body)
}
func (c *cancellationClient) Update(_ context.Context, _ Kind, _ string, req Request) (map[string]any, error) {
	return c.operation(ActionUpdate, req.Body)
}
func (c *cancellationClient) Destroy(context.Context, Kind, string) error {
	_, err := c.operation(ActionDestroy, map[string]any{})
	return err
}

func cancellationPlan(t *testing.T, action Action) (*Registry, *Lockfile, *Plan) {
	t.Helper()
	registry := testRegistry()
	lock := newLockfile(filepath.Join(t.TempDir(), testLockfileName))
	spec := registry.specOrZero("widget")
	changes := []*Change{}
	for i, name := range []string{"first", "interrupted", "untouched"} {
		key := "./widgets/" + name + ".yml"
		change := &Change{Key: key, Kind: "widget", Action: ActionCreate, spec: spec, Source: &Source{Key: key, Kind: "widget", Body: map[string]any{"name": name}}, Hash: "hash"}
		if i == 1 {
			change.Action = action
			if action != ActionCreate {
				change.Entry = &LockEntry{Kind: "widget", ID: "wdgt_existing", Hash: "old"}
				lock.Resources[key] = change.Entry
			}
		}
		changes = append(changes, change)
	}
	return registry, lock, &Plan{Changes: changes, builder: bodyBuilder{registry: registry}}
}

func TestApplyPreservesCancellationCausesAndSavedProgress(t *testing.T) {
	for _, action := range []Action{ActionCreate, ActionUpdate, ActionDestroy} {
		for _, cause := range []error{context.Canceled, context.DeadlineExceeded, errors.New("ordinary failure")} {
			t.Run(string(action)+"/"+cause.Error(), func(t *testing.T) {
				registry, lock, plan := cancellationPlan(t, action)
				ctx, cancel := context.WithCancel(context.Background())
				defer cancel()
				failure := WithHint(fmt.Errorf("remote write: %w", cause), "check remote state before retrying")
				client := &cancellationClient{failAction: action, failure: failure}
				if errors.Is(cause, context.Canceled) {
					client.onFailure = cancel
				}
				var reported []Applied
				applier := Applier{Client: client, Lock: lock, Report: func(a Applied) { reported = append(reported, a) }}
				result, err := applier.Apply(ctx, plan)
				require.ErrorIs(t, err, cause)
				assert.ErrorIs(t, err, failure)
				assert.Equal(t, "check remote state before retrying", Hint(err))
				require.NotNil(t, result)
				assert.Equal(t, 1, result.Created)
				assert.Zero(t, result.Updated)
				assert.Zero(t, result.Destroyed)
				require.Len(t, reported, 2)
				assert.Equal(t, "created", reported[0].Outcome)
				expected := "failed"
				if errors.Is(cause, context.Canceled) {
					expected = "interrupted"
				}
				assert.Equal(t, expected, reported[1].Outcome)
				assert.ErrorIs(t, reported[1].Err, cause)
				assert.Len(t, client.calls, 2, "later resources must not be written")
				reloaded, loadErr := LoadLockfile(registry, lock.Path)
				require.NoError(t, loadErr)
				assert.Equal(t, "wdgt_created", reloaded.Resources[plan.Changes[0].Key].ID)
				assert.NotContains(t, reloaded.Resources, plan.Changes[2].Key)
				if action != ActionCreate {
					assert.Equal(t, "wdgt_existing", reloaded.Resources[plan.Changes[1].Key].ID)
				} else {
					assert.NotContains(t, reloaded.Resources, plan.Changes[1].Key)
				}
			})
		}
	}
}

func TestApplyJoinsSaveFailureWithoutLosingCancellation(t *testing.T) {
	_, lock, plan := cancellationPlan(t, ActionCreate)
	client := &cancellationClient{failAction: ActionCreate, failure: fmt.Errorf("transport: %w", context.Canceled)}
	file := filepath.Join(t.TempDir(), "not-a-directory")
	require.NoError(t, os.WriteFile(file, []byte("fixture"), 0600))
	client.onFailure = func() { lock.Path = filepath.Join(file, "lock.json") }
	applier := Applier{Client: client, Lock: lock}
	result, err := applier.Apply(context.Background(), plan)
	require.ErrorIs(t, err, context.Canceled)
	var pathError *os.PathError
	assert.ErrorAs(t, err, &pathError)
	assert.ErrorContains(t, err, "writing ")
	assert.Equal(t, 1, result.Created)
	assert.Len(t, client.calls, 2)
}

func TestApplyPrecancelledContextStillSavesWithoutRequests(t *testing.T) {
	registry, lock, plan := cancellationPlan(t, ActionCreate)
	client := &cancellationClient{}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	applier := Applier{Client: client, Lock: lock}
	result, err := applier.Apply(ctx, plan)
	assert.ErrorIs(t, err, context.Canceled)
	assert.Zero(t, result.Created)
	assert.Empty(t, client.calls)
	saved, loadErr := LoadLockfile(registry, lock.Path)
	require.NoError(t, loadErr)
	assert.Empty(t, saved.Resources)
}
