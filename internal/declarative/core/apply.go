package core

import (
	"context"
	"errors"
	"fmt"
)

// Applier executes a Plan. It persists state after every individual success,
// because the alternative — writing the lockfile once at the end — loses the
// IDs of everything created before a mid-run failure, and those resources are
// then invisible to the next apply and get created a second time.
type Applier struct {
	Client Client
	Lock   *Lockfile
	// Report is told about each resource as soon as it has been applied, so a
	// long apply shows progress. Optional.
	Report func(Applied)
}

// Applied is one resource's outcome, as it happens.
type Applied struct {
	Change *Change
	// ID is the resource's identifier after the operation, when it has one.
	ID string
	// Outcome is what was actually done: "created", "updated", "unchanged",
	// the kind's Destroy.Past, "already gone" — or "failed" or "interrupted",
	// with Err set.
	Outcome string
	Err     error
}

// Result summarizes what an apply actually did.
type Result struct {
	Created   int
	Updated   int
	Destroyed int
	Unchanged int
}

// Apply executes the plan in order and stops at the first failure. A plan with
// any blocked change is refused before anything is sent, with a nil Result.
// Otherwise the Result counts what was done, even alongside an error, and the
// lockfile has been saved.
func (a *Applier) Apply(ctx context.Context, plan *Plan) (*Result, error) {
	if blocked := plan.Blocked(); len(blocked) > 0 {
		return nil, blockedError(blocked)
	}

	result := &Result{}
	resolved := map[string]Target{}

	for _, change := range plan.Changes {
		if err := ctx.Err(); err != nil {
			return result, a.saveState(err)
		}
		target, err := a.applyOne(ctx, plan.builder, change, resolved, result)
		if err != nil {
			outcome := "failed"
			if errors.Is(err, context.Canceled) {
				// The write may or may not have landed; the next plan will say.
				outcome, err = "interrupted", fmt.Errorf("interrupted: %w", err)
			}
			id := ""
			if change.Entry != nil {
				id = change.Entry.ID
			}
			a.report(change, id, outcome, err)
			// Persist whatever succeeded before returning the failure.
			return result, a.saveState(fmt.Errorf("%s: %w", change.Key, err))
		}
		if change.Action != ActionDestroy {
			resolved[change.Key] = target
		}
	}

	return result, a.saveState(nil)
}

// saveState writes the lockfile and returns prior, which may be nil. A save
// failure is joined onto prior so the error that stopped the run still shows.
func (a *Applier) saveState(prior error) error {
	if err := a.Lock.Save(); err != nil {
		return errors.Join(prior, fmt.Errorf("writing %s: %w", a.Lock.Path, err))
	}
	return prior
}

func blockedError(blocked []*Change) error {
	errs := make([]error, 0, len(blocked))
	for _, c := range blocked {
		errs = append(errs, fmt.Errorf("%s: %w", c.Key, c.Blocked))
	}
	return errors.Join(errs...)
}

func (a *Applier) report(change *Change, id, outcome string, err error) {
	if a.Report != nil {
		a.Report(Applied{Change: change, ID: id, Outcome: outcome, Err: err})
	}
}

// applyOne carries out one planned change and returns the target it leaves
// behind for the resources that reference it.
func (a *Applier) applyOne(ctx context.Context, builder bodyBuilder, change *Change, resolved map[string]Target, result *Result) (Target, error) {
	switch change.Action {
	case ActionDestroy:
		return Target{}, a.destroy(ctx, change, result)
	case ActionNoop:
		result.Unchanged++
		if change.Hash != change.Entry.Hash {
			// Nothing to send, but the file was respelled: remember its new
			// fingerprint so the next plan need not rediscover that.
			change.Entry.Hash = change.Hash
		}
		a.report(change, change.Entry.ID, "unchanged", nil)
		return targetAfter(change), nil
	case ActionCreate, ActionUpdate:
		return a.write(ctx, builder, change, resolved, result)
	default:
		return Target{}, fmt.Errorf("unknown action %q", change.Action)
	}
}

// write creates or updates the resource change describes and commits the
// outcome.
func (a *Applier) write(ctx context.Context, builder bodyBuilder, change *Change, resolved map[string]Target, result *Result) (Target, error) {
	// Rebuild the body now that dependencies have real IDs and versions. The
	// plan's copy may hold "(known after apply)" placeholders.
	desired, err := builder.build(change.Source, resolved)
	if err != nil {
		return Target{}, err
	}
	if desired.Unresolved {
		// Every dependency has been applied by now, so a placeholder left in
		// the body means one of them came back without an ID or version.
		return Target{}, fmt.Errorf("%s: a referenced resource was applied but reported no ID or version to pin; re-run apply to read it back from the server", change.Key)
	}

	// Fingerprint before mutating: if this failed afterwards the resource
	// would exist with nothing in the lockfile pointing at it, and the next
	// apply would create a second one.
	hash, err := hashFor(change, desired)
	if err != nil {
		return Target{}, err
	}

	var obj map[string]any
	var id string
	switch change.Action {
	case ActionCreate:
		obj, id, err = a.create(ctx, change, desired)
	default:
		obj, id, err = a.update(ctx, change, desired)
	}
	if err != nil {
		return Target{}, err
	}
	return a.commit(change, obj, id, hash, result)
}

func (a *Applier) create(ctx context.Context, change *Change, desired *desiredBody) (map[string]any, string, error) {
	obj, err := a.Client.Create(ctx, change.Kind, Request{Body: desired.Body, Payload: change.payload})
	if err != nil {
		return nil, "", err
	}
	id, _ := obj["id"].(string)
	if id == "" {
		return nil, "", fmt.Errorf("the API returned no id for the created %s", change.Kind)
	}
	return obj, id, nil
}

func (a *Applier) update(ctx context.Context, change *Change, desired *desiredBody) (map[string]any, string, error) {
	id := change.Entry.ID
	body := deepCopyMap(desired.Body)
	applyClears(body, change)
	if change.spec.UpdateNeedsVersion {
		// The API uses this to reject an update racing another writer.
		// Send what we actually read, not what the lockfile remembers.
		v := remoteVersion(change.spec, change.Remote)
		if v == "" {
			return nil, "", fmt.Errorf(
				"cannot update %s: the server did not report a current version, and sending none would skip the concurrent-write check", id)
		}
		body["version"] = versionValue(change.spec, v)
	}
	obj, err := a.Client.Update(ctx, change.Kind, id, Request{Body: body, Payload: change.payload})
	if err != nil {
		return nil, "", err
	}
	return obj, id, nil
}

// commit is the tail every successful write shares: record the outcome, tell
// Report, save the lockfile before the next write, and return the Target that
// dependents will resolve against.
func (a *Applier) commit(change *Change, obj map[string]any, id, hash string, result *Result) (Target, error) {
	version := responseVersion(change.spec, obj)
	a.record(change, id, version, hash, obj)
	switch change.Action {
	case ActionCreate:
		result.Created++
		a.report(change, id, "created", nil)
	default:
		result.Updated++
		a.report(change, id, "updated", nil)
	}

	if err := a.saveState(nil); err != nil {
		return Target{}, err
	}
	return Target{
		Kind:      change.Kind,
		ID:        id,
		Version:   versionValue(change.spec, version),
		Versioned: change.spec.VersionField != "",
		Known:     true,
	}, nil
}

// hashFor is the desired-state fingerprint to record on success. It recomputes
// from the body actually sent, because the plan's copy may have held
// placeholder versions for dependencies that have since been applied.
//
// A payload-backed resource has no such dependency: its fingerprint is the
// payload, fixed when the plan collected it. Reusing the plan's hash is both
// cheaper and the only way to guarantee it describes what was uploaded.
func hashFor(change *Change, desired *desiredBody) (string, error) {
	if change.payload != nil {
		return change.Hash, nil
	}
	return hashBody(desired.Body)
}

// record writes the outcome into the in-memory lockfile.
func (a *Applier) record(change *Change, id, version, hash string, obj map[string]any) {
	entry := &LockEntry{
		Kind:    change.Kind,
		ID:      id,
		Version: version,
		Hash:    hash,
	}
	if change.Source != nil {
		entry.Revision = change.Source.Pin.Revision
		entry.Subpath = change.Source.Pin.Subpath
	}
	// A payload write answers with a version rather than the resource, so
	// fingerprinting it would cost another round trip — and a payload-backed
	// kind has a version to compare instead.
	if change.payload == nil && obj != nil {
		if h, err := hashBody(normalizeRemote(change.spec, obj)); err == nil {
			entry.RemoteHash = h
		}
	}
	a.Lock.Resources[change.Key] = entry
}

func (a *Applier) destroy(ctx context.Context, change *Change, result *Result) error {
	outcome := change.spec.Destroy.Past
	err := a.Client.Destroy(ctx, change.Kind, change.Entry.ID)
	switch {
	case errors.Is(err, ErrNotFound):
		outcome = "already gone"
	case err != nil:
		return err
	}
	delete(a.Lock.Resources, change.Key)
	result.Destroyed++
	a.report(change, change.Entry.ID, outcome, nil)
	return a.saveState(nil)
}

// applyClears turns "the file no longer declares this" into an explicit null,
// for the fields the API documents as clearable. Everything else follows the
// API's omit-to-preserve rule, so a field the server defaults for us is left
// alone rather than fought over on every run.
func applyClears(body map[string]any, change *Change) {
	for _, field := range change.Diff.clears() {
		body[field] = nil
	}
	clearRemovedMetadataKeys(change.spec, body, change.Remote)
}

// clearRemovedMetadataKeys sends null for metadata keys that used to exist and
// no longer appear locally. The bag is patched key by key, so dropping a key
// from the file would otherwise leave it on the server forever.
func clearRemovedMetadataKeys(spec KindSpec, body map[string]any, remote map[string]any) {
	field := spec.metadataField
	if field == "" {
		return
	}
	remoteBag, _ := remote[field].(map[string]any)
	if len(remoteBag) == 0 {
		return
	}
	desiredBag, _ := body[field].(map[string]any)
	if desiredBag == nil {
		desiredBag = map[string]any{}
	}
	removed := removedMetadataKeys(desiredBag, remoteBag)
	if len(removed) == 0 {
		return
	}
	for _, k := range removed {
		desiredBag[k] = nil
	}
	body[field] = desiredBag
}
