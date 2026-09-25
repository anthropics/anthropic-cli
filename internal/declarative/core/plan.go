package core

import (
	"cmp"
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"slices"
	"strings"
	"sync"
)

// Action is what apply will do to one resource.
type Action string

const (
	ActionCreate  Action = "create"
	ActionUpdate  Action = "update"
	ActionNoop    Action = "no change"
	ActionDestroy Action = "destroy"
)

// Change is one planned operation.
type Change struct {
	Key    string
	Kind   Kind
	Action Action

	// Source is the local definition. Nil for a destroy, whose file is gone.
	Source *Source
	// Entry is the lockfile record. Nil for a create.
	Entry *LockEntry

	// Desired is the request body apply will send, with references resolved.
	Desired map[string]any
	// Remote is the object as the server currently has it, when known.
	Remote map[string]any

	// Hash is the desired-state fingerprint to record on success.
	Hash string
	// Diff is the field-level detail shown under an update; nil otherwise.
	Diff *Diff
	// Reasons explains why this change exists, beyond the field diff — an
	// out-of-band edit, a dependency version bump, a recreated resource.
	Reasons []string
	// Sensitive marks dotted paths whose values must never be printed.
	Sensitive map[string]bool

	// Replaces is the ID of the tracked resource this create supersedes,
	// when one that was gone or archived has to be made again.
	Replaces string
	// Drift is set when the server's copy no longer matches what we recorded,
	// meaning somebody changed it outside this config.
	Drift bool
	// Blocked carries the reason apply must refuse to proceed.
	Blocked error
	// Unresolved marks a body containing "(known after apply)" placeholders.
	Unresolved bool

	// spec is the kind's behaviour, carried so a Change can answer questions
	// about itself without a registry in hand — which is what lets a formatter
	// live outside this package.
	spec KindSpec
	// payload is the upload content, set only for a kind whose content
	// travels as files rather than in the request body.
	payload Payload
}

// Destroy is what removing this resource actually does.
func (c *Change) Destroy() Destroy { return c.spec.Destroy }

// SummaryFields lists the fields that identify this resource in plan output.
func (c *Change) SummaryFields() []string { return c.spec.summary }

// Upload describes the content this change uploads instead of a JSON body,
// e.g. "3 files". Returns false when the change is an ordinary body write.
func (c *Change) Upload() (string, bool) {
	if c.payload == nil {
		return "", false
	}
	return c.payload.Describe(), true
}

// Plan is the full set of operations, in execution order.
type Plan struct {
	Changes  []*Change
	Warnings []string
	// LockfilePath is where state will be written.
	LockfilePath string
	// LockfileExisted reports whether state was found or is being started.
	LockfileExisted bool

	// builder rebuilds each request body at apply time, once the dependencies
	// that were "(known after apply)" during planning have real IDs.
	builder bodyBuilder
}

// Counts summarizes the plan for the confirmation prompt.
func (p *Plan) Counts() (create, update, destroy, noop int) {
	for _, c := range p.Changes {
		switch c.Action {
		case ActionCreate:
			create++
		case ActionUpdate:
			update++
		case ActionDestroy:
			destroy++
		case ActionNoop:
			noop++
		}
	}
	return
}

// HasWork reports whether anything would actually change.
func (p *Plan) HasWork() bool {
	c, u, d, _ := p.Counts()
	return c+u+d > 0
}

// Change returns the planned change for a resource key, or false when the
// plan has none.
func (p *Plan) Change(key string) (*Change, bool) {
	for _, c := range p.Changes {
		if c.Key == key {
			return c, true
		}
	}
	return nil, false
}

// Blocked returns the changes apply must refuse to perform.
func (p *Plan) Blocked() []*Change {
	var out []*Change
	for _, c := range p.Changes {
		if c.Blocked != nil {
			out = append(out, c)
		}
	}
	return out
}

// Planner diffs local definitions against the API and produces a Plan.
type Planner struct {
	// Registry supplies the kinds this planner understands.
	Registry *Registry
	Client   Client
	Lock     *Lockfile

	// Force proceeds even when a resource changed out of band.
	Force bool
	// Prune destroys resources that are in the lockfile but no longer on disk.
	Prune bool
	// Concurrency bounds parallel reads of remote state. Zero or less means 8.
	Concurrency int

	mu       sync.Mutex // guards warnings
	warnings []string
}

// Plan builds the execution plan for everything the loader has gathered.
func (p *Planner) Plan(ctx context.Context, loader *Loader) (*Plan, error) {
	order, err := loader.TopoOrder()
	if err != nil {
		return nil, err
	}

	remotes := p.fetchRemotes(ctx, order)
	if err := p.checkAllMissing(remotes); err != nil {
		return nil, err
	}

	plan := &Plan{
		LockfilePath:    p.Lock.Path,
		LockfileExisted: p.Lock.Existed(),
		builder:         bodyBuilder{registry: p.Registry, slots: loader.slots},
	}
	resolved := map[string]Target{}

	for _, key := range order {
		src := loader.Sources()[key]
		change, err := p.planOne(plan.builder, src, resolved, remotes[key])
		if err != nil {
			return nil, err
		}
		plan.Changes = append(plan.Changes, change)
		resolved[key] = targetAfter(change)
	}

	plan.Changes = append(plan.Changes, p.planDestroys(order)...)
	plan.Warnings = append([]string{}, p.warnings...)
	return plan, nil
}

// remoteResult is the outcome of reading one resource's current state.
type remoteResult struct {
	Object map[string]any
	Err    error
}

// fetchRemotes reads current state for every tracked resource in parallel.
// A plan over a few dozen resources should not take a few dozen round trips
// of latency.
func (p *Planner) fetchRemotes(ctx context.Context, keys []string) map[string]remoteResult {
	type job struct {
		key   string
		entry *LockEntry
	}
	var jobs []job
	for _, key := range keys {
		if entry, ok := p.Lock.Get(key); ok {
			jobs = append(jobs, job{key, entry})
		}
	}

	limit := p.Concurrency
	if limit <= 0 {
		limit = 8
	}
	results := make(map[string]remoteResult, len(jobs))
	var mu sync.Mutex
	var wg sync.WaitGroup
	sem := make(chan struct{}, limit)

	for _, j := range jobs {
		wg.Go(func() {
			sem <- struct{}{}
			defer func() { <-sem }()

			obj, err := p.Client.Get(ctx, j.entry.Kind, j.entry.ID)
			mu.Lock()
			results[j.key] = remoteResult{Object: obj, Err: err}
			mu.Unlock()
		})
	}
	wg.Wait()
	return results
}

// checkAllMissing refuses to plan when two or more resources are tracked and
// every one came back not found. One vanished resource is deletion; all of
// them at once is almost always credentials for a different organization or
// host, and offering to recreate everything under new IDs would be the wrong
// reading of it. --force says the reading is right after all.
func (p *Planner) checkAllMissing(remotes map[string]remoteResult) error {
	if p.Force || len(remotes) < 2 {
		return nil
	}
	for _, r := range remotes {
		if !errors.Is(r.Err, ErrNotFound) {
			return nil
		}
	}
	return fmt.Errorf("none of the %d resources in %s were found with these credentials — this usually means they belong to a different organization or profile than the one in use. If they really were all deleted, re-run with --force to recreate them",
		len(remotes), filepath.Base(p.Lock.Path))
}

// planOne decides what to do with one loaded resource. An untracked resource
// is a create. A tracked one that changed kind, or is gone or archived on the
// server, is a create too, blocked until the user says otherwise. Anything
// else is compared three ways: the local hash against the lockfile, the
// server's copy against the lockfile (drift), and the body against the
// server's copy, which catches a file that changed into what the server
// already holds. An update that touches an immutable field, or that would
// overwrite drift without --force, is planned but Blocked.
func (p *Planner) planOne(builder bodyBuilder, src *Source, resolved map[string]Target, remote remoteResult) (*Change, error) {
	desired, err := builder.build(src, resolved)
	if err != nil {
		return nil, err
	}

	change := &Change{
		spec:       p.Registry.specOrZero(src.Kind),
		Key:        src.Key,
		Kind:       src.Kind,
		Source:     src,
		Desired:    desired.Body,
		Sensitive:  desired.Sensitive,
		Unresolved: desired.Unresolved,
	}

	// A payload-backed kind fingerprints its content; everything else
	// fingerprints the body it will send.
	if src.Payload != nil {
		change.payload = src.Payload
		change.Hash, err = src.Payload.Fingerprint()
	} else {
		change.Hash, err = hashBody(desired.Body)
	}
	if err != nil {
		return nil, err
	}

	entry, tracked := p.Lock.Get(src.Key)
	if !tracked {
		change.Action = ActionCreate
		return change, nil
	}
	change.Entry = entry

	if entry.Kind != src.Kind {
		// Creating the new one would leave the old resource on the server with
		// nothing tracking it — invisible to every later apply, --prune
		// included.
		change.Action = ActionCreate
		change.Blocked = fmt.Errorf(
			"this file was a %s (%s) and is now a %s; remove it, apply with --prune, then add it back",
			entry.Kind, entry.ID, src.Kind)
		return change, nil
	}

	if remote.Err != nil {
		if errors.Is(remote.Err, ErrNotFound) {
			return p.planRecreate(change, "is no longer on the server"), nil
		}
		return nil, fmt.Errorf("reading %s (%s): %w", src.Key, entry.ID, remote.Err)
	}
	change.Remote = remote.Object

	if change.spec.IsArchived(remote.Object) {
		return p.planRecreate(change, "is archived and cannot be revived"), nil
	}

	detectDrift(change)

	if change.Hash == entry.Hash && !change.Drift && !change.Unresolved {
		change.Action = ActionNoop
		return change, nil
	}
	change.Action = ActionUpdate

	change.Diff = change.buildDiff()

	// The file changed, but into something the server already has — it was
	// reformatted, or respelled the way the server normalises it anyway.
	// There is nothing to send; the new fingerprint is recorded on apply so
	// the comparison is not repeated forever.
	if !change.Drift && !change.Unresolved && change.payload == nil && change.Diff == nil {
		change.Action = ActionNoop
		return change, nil
	}

	if err := change.checkImmutable(); err != nil {
		change.Blocked = err
	} else if change.Drift && !p.Force {
		change.Blocked = fmt.Errorf(
			"%s changed outside this config since the last apply; re-run with --force to overwrite it", entry.ID)
	}
	return change, nil
}

// planRecreate handles a tracked resource that is gone or archived. Making a
// new one is the only option the API leaves, but it mints a new ID — which
// silently breaks anything referencing the old one from outside this config,
// and is not something to do on the user's behalf without asking.
//
// So it is refused the same way an out-of-band edit is, and released by the
// same flag. Both are "the server no longer matches what we recorded"; that a
// recreate changes identity rather than just fields makes it the more
// consequential of the two, not the less.
func (p *Planner) planRecreate(change *Change, why string) *Change {
	entry := change.Entry
	change.Action = ActionCreate
	change.Entry = nil
	change.Replaces = entry.ID
	change.Reasons = append(change.Reasons,
		fmt.Sprintf("%s %s, so it would be created again with a new id", entry.ID, why))
	if !p.Force {
		change.Blocked = fmt.Errorf(
			"%s %s; re-run with --force to create a replacement, or remove %s from %s if it is gone for good",
			entry.ID, why, change.Key, filepath.Base(p.Lock.Path))
	}
	return change
}

// detectDrift compares the server's copy against what we recorded at the end of
// the last apply. A versioned kind has a version to compare; for everything
// else the recorded hash of the normalized object is the only signal available.
func detectDrift(change *Change) {
	entry, remote := change.Entry, change.Remote
	if v := remoteVersion(change.spec, remote); v != "" && entry.Version != "" && v != entry.Version {
		change.Drift = true
		change.Reasons = append(change.Reasons,
			fmt.Sprintf("server has version %s but the lockfile recorded %s", v, entry.Version))
		return
	}
	if entry.RemoteHash == "" {
		// Nothing to compare against — the last apply could not fingerprint
		// this object. No signal is better than a false one.
		return
	}
	current, err := hashBody(normalizeRemote(change.spec, remote))
	if err != nil {
		return
	}
	if current != entry.RemoteHash {
		change.Drift = true
		change.Reasons = append(change.Reasons, "the server's copy no longer matches what the last apply left behind")
	}
}

// targetAfter predicts what a resource will look like once its change is
// applied, for the resources that reference it. A resource about to change
// reports an unknown version. That cascades an update into everything that
// pins it, so no dependent is left pinned to the version being replaced.
func targetAfter(change *Change) Target {
	target := Target{Kind: change.Kind, Versioned: change.spec.VersionField != ""}
	if change.Action == ActionCreate {
		return target
	}
	target.Known = true
	target.ID = change.Entry.ID
	// Only an unchanged resource has a version yet; a changing one is assigned
	// its next version by the server.
	if change.Action == ActionNoop {
		target.Version = versionValue(change.spec, currentVersion(change))
	}
	return target
}

// planDestroys turns lockfile entries with no loaded source into destroys.
// The files are gone, so there is no reference graph to follow. Destroys run
// in reverse kind order instead, so a kind that can reference another is
// removed before the kind it points at. Within a kind the order is by key.
func (p *Planner) planDestroys(present []string) []*Change {
	var orphans []string
	for _, key := range p.Lock.Keys() {
		if slices.Contains(present, key) {
			continue
		}
		orphans = append(orphans, key)
	}
	if len(orphans) == 0 {
		return nil
	}
	if !p.Prune {
		for _, key := range orphans {
			p.warnf("%s is in the lockfile but was not loaded this run; pass --prune to remove it, or name it on the command line to keep it tracked", key)
		}
		return nil
	}

	// Reverse kind rank: the API refuses to remove a resource that something
	// still references.
	slices.SortFunc(orphans, func(a, b string) int {
		ra, rb := p.Registry.rank(p.Lock.Resources[a].Kind), p.Registry.rank(p.Lock.Resources[b].Kind)
		return cmp.Or(cmp.Compare(rb, ra), strings.Compare(a, b))
	})

	var changes []*Change
	for _, key := range orphans {
		if keyEscapesRoot(p.Lock.Root(), key) {
			p.warnf("%s resolves outside the lockfile directory; refusing to prune it", key)
			continue
		}
		entry := p.Lock.Resources[key]
		change := &Change{
			spec:   p.Registry.specOrZero(entry.Kind),
			Key:    key,
			Kind:   entry.Kind,
			Action: ActionDestroy,
			Entry:  entry,
		}
		changes = append(changes, change)
	}
	return changes
}

func (p *Planner) warnf(format string, args ...any) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.warnings = append(p.warnings, fmt.Sprintf(format, args...))
}
