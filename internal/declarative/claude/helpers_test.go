package claude

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/anthropics/anthropic-cli/v2/internal/declarative/core"
	"github.com/stretchr/testify/require"
)

// writeTree materializes a map of relative path to contents under a temp dir.
func writeTree(t *testing.T, files map[string]string) string {
	t.Helper()
	root := t.TempDir()
	for rel, content := range files {
		full := filepath.Join(root, filepath.FromSlash(rel))
		require.NoError(t, os.MkdirAll(filepath.Dir(full), 0o755))
		require.NoError(t, os.WriteFile(full, []byte(content), 0o644))
	}
	return root
}

// deepCopyMap clones a decoded JSON object so the fake client hands out copies
// rather than aliases of its own state.
func deepCopyMap(m map[string]any) map[string]any {
	out := make(map[string]any, len(m))
	for k, v := range m {
		switch t := v.(type) {
		case map[string]any:
			out[k] = deepCopyMap(t)
		case []any:
			list := make([]any, len(t))
			for i, item := range t {
				if nested, ok := item.(map[string]any); ok {
					list[i] = deepCopyMap(nested)
					continue
				}
				list[i] = item
			}
			out[k] = list
		default:
			out[k] = v
		}
	}
	return out
}

// plannedField builds the plan for a loaded tree and returns one field of the
// request body a resource would be created with. Tests assert on this rather
// than on the loader's internal bookkeeping: it is what actually reaches the
// API, and it survives the engine being rearranged underneath.
func plannedField(t *testing.T, l *core.Loader, key, field string) any {
	t.Helper()
	body, err := planBody(t, l, key)
	require.NoError(t, err)
	require.Contains(t, body, field, "%s has no %s", key, field)
	return body[field]
}

// planBody plans a loaded tree and returns the request body one resource would
// be created with, or the planning error.
func planBody(t *testing.T, l *core.Loader, key string) (map[string]any, error) {
	t.Helper()
	lock, err := core.LoadLockfile(Registry(), filepath.Join(t.TempDir(), Registry().LockfileName()))
	require.NoError(t, err)

	planner := &core.Planner{Registry: Registry(), Client: newFakeClient(), Lock: lock}
	plan, err := planner.Plan(context.Background(), l)
	if err != nil {
		return nil, err
	}
	for _, c := range plan.Changes {
		if c.Key == key {
			return c.Desired, nil
		}
	}
	t.Fatalf("no change for %s", key)
	return nil, nil
}

// fakeClient is an in-memory stand-in for the API. It reproduces the behaviors
// the reconciler actually depends on: agents bump an integer version on every
// update, skills mint a new opaque version per upload, and, once now is set, a
// write stamps updated_at.
type fakeClient struct {
	objects map[string]map[string]any
	seq     int
	// calls records every mutating operation, in order, so tests can assert on
	// sequencing rather than just end state.
	calls []string
	// failOn makes the named operation ("create environment") return an error.
	failOn string
	// now, if set, is the updated_at a write leaves behind.
	now string
}

func newFakeClient() *fakeClient {
	return &fakeClient{objects: map[string]map[string]any{}}
}

func (f *fakeClient) nextID(kind core.Kind) string {
	f.seq++
	spec, _ := Registry().Spec(kind)
	return fmt.Sprintf("%s_%02d", spec.IDPrefix, f.seq)
}

func (f *fakeClient) record(op string) error {
	f.calls = append(f.calls, op)
	if f.failOn == op {
		return fmt.Errorf("injected failure for %q", op)
	}
	return nil
}

func (f *fakeClient) Get(_ context.Context, kind core.Kind, id string) (map[string]any, error) {
	obj, ok := f.objects[id]
	if !ok {
		return nil, fmt.Errorf("%w: %s", core.ErrNotFound, id)
	}
	return deepCopyMap(obj), nil
}

func (f *fakeClient) Create(_ context.Context, kind core.Kind, req core.Request) (map[string]any, error) {
	if err := f.record("create " + string(kind)); err != nil {
		return nil, err
	}
	if kind == KindSkill {
		obj := map[string]any{
			"id":                f.nextID(KindSkill),
			"type":              "skill",
			"display_name":      req.Body["display_name"],
			"latest_version_id": strconv.Itoa(1000 + f.seq),
		}
		// The real client reads the name back from the latest version.
		if name, ok := req.Body["name"]; ok {
			obj["name"] = name
		}
		f.objects[obj["id"].(string)] = obj
		return deepCopyMap(obj), nil
	}
	obj := deepCopyMap(req.Body)
	obj["id"] = f.nextID(kind)
	obj["type"] = string(kind)
	if spec, _ := Registry().Spec(kind); spec.VersionField == "version" {
		obj["version"] = int64(1)
	}
	if f.now != "" {
		obj["updated_at"] = f.now
	}
	f.objects[obj["id"].(string)] = obj
	return deepCopyMap(obj), nil
}

func (f *fakeClient) Update(_ context.Context, kind core.Kind, id string, req core.Request) (map[string]any, error) {
	op := "update " + string(kind)
	if kind == KindSkill {
		op = "create skill version"
	}
	if err := f.record(op); err != nil {
		return nil, err
	}
	obj, ok := f.objects[id]
	if !ok {
		return nil, fmt.Errorf("%w: %s", core.ErrNotFound, id)
	}
	if kind == KindSkill {
		f.seq++
		version := strconv.Itoa(1000 + f.seq)
		obj["latest_version_id"] = version
		return map[string]any{"type": "skill_version", "skill_id": id, "version": version}, nil
	}
	for k, v := range req.Body {
		if v == nil {
			delete(obj, k)
			continue
		}
		// Metadata is patched key by key: a null value deletes that key rather
		// than replacing the whole bag.
		if k == "metadata" {
			patch, _ := v.(map[string]any)
			bag, _ := obj["metadata"].(map[string]any)
			if bag == nil {
				bag = map[string]any{}
			}
			for mk, mv := range patch {
				if mv == nil {
					delete(bag, mk)
					continue
				}
				bag[mk] = mv
			}
			obj["metadata"] = bag
			continue
		}
		obj[k] = v
	}
	if spec, _ := Registry().Spec(kind); spec.VersionField == "version" {
		current, _ := obj["version"].(int64)
		obj["version"] = current + 1
	}
	if f.now != "" {
		obj["updated_at"] = f.now
	}
	return deepCopyMap(obj), nil
}

func (f *fakeClient) Destroy(_ context.Context, kind core.Kind, id string) error {
	if err := f.record("destroy " + string(kind) + " " + id); err != nil {
		return err
	}
	delete(f.objects, id)
	return nil
}

// harness bundles everything a plan-and-apply round needs.
type harness struct {
	t      *testing.T
	root   string
	client *fakeClient
	lock   *core.Lockfile
}

func newHarness(t *testing.T, files map[string]string) *harness {
	t.Helper()
	root := writeTree(t, files)
	lock, err := core.LoadLockfile(Registry(), filepath.Join(root, Registry().LockfileName()))
	require.NoError(t, err)
	return &harness{t: t, root: root, client: newFakeClient(), lock: lock}
}

// reloadLock re-reads state from disk, proving the lockfile actually
// round-trips rather than only living in memory.
func (h *harness) reloadLock() {
	h.t.Helper()
	lock, err := core.LoadLockfile(Registry(), h.lock.Path)
	require.NoError(h.t, err)
	h.lock = lock
}

func (h *harness) plan(opts ...func(*core.Planner)) *core.Plan {
	h.t.Helper()
	loader := core.NewLoader(Registry(), h.root, nil)
	require.NoError(h.t, loader.Add(context.Background(), []string{h.root}))
	require.NoError(h.t, loader.AddKeys(context.Background(), h.lock.Keys()))

	planner := &core.Planner{Registry: Registry(), Client: h.client, Lock: h.lock}
	for _, opt := range opts {
		opt(planner)
	}
	plan, err := planner.Plan(context.Background(), loader)
	require.NoError(h.t, err)
	return plan
}

func (h *harness) apply(opts ...func(*core.Planner)) *core.Result {
	h.t.Helper()
	plan := h.plan(opts...)
	applier := &core.Applier{Client: h.client, Lock: h.lock}
	res, err := applier.Apply(context.Background(), plan)
	require.NoError(h.t, err)
	h.reloadLock()
	return res
}

func withPrune(p *core.Planner) { p.Prune = true }
func withForce(p *core.Planner) { p.Force = true }

func actionFor(plan *core.Plan, key string) core.Action {
	for _, c := range plan.Changes {
		if c.Key == key {
			return c.Action
		}
	}
	return ""
}

const twoAgentsAndASkill = `---
model: claude-sonnet-5
skills:
  - ../skills/pr-writer
multiagent:
  type: coordinator
  agents:
    - ./code-verifier.md
---
review
`

func basicTree() map[string]string {
	return map[string]string{
		"skills/pr-writer/SKILL.md": "---\nname: pr-writer\ndescription: d\n---\nhow to write\n",
		"agents/code-verifier.md":   "---\nmodel: claude-sonnet-5\n---\nverify\n",
		"agents/code-reviewer.md":   twoAgentsAndASkill,
		"environments/cloud.yml":    "description: sandbox\n",
	}
}
