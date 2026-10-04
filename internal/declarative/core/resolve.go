package core

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

// resolvedSlot is a reference slot with its entries resolved to keys.
type resolvedSlot struct {
	Ref     RefSlot
	Entries []slotEntry
}

// slotEntry is one resolved element of a reference slot, in declaration order.
type slotEntry struct {
	// TargetKey is the key of the resource this points at, or "" when the user
	// hand-wrote an API object we pass through untouched.
	TargetKey string
	// Pinned is false when the user asked for `version: latest`.
	Pinned bool
	// Value is the hand-written object, sent as is. Set only when TargetKey
	// is empty.
	Value any
	// Extra is merged into the encoded reference: keys the user wrote next
	// to `path` that configure the attachment rather than name the target.
	Extra map[string]any
}

// resolveSlots walks the reference slots for a source, expands globs, and
// recursively loads each referenced resource.
func (l *Loader) resolveSlots(ctx context.Context, src *Source) error {
	var resolved []resolvedSlot
	for _, slot := range l.registry.specOrZero(src.Kind).refs {
		one, err := l.resolveSlot(ctx, src, slot)
		if err != nil {
			return err
		}
		if one != nil {
			resolved = append(resolved, *one)
		}
	}
	l.slots[src.Key] = resolved
	return nil
}

// resolveSlot resolves every reference written in one slot of a source. It
// returns nil when the source does not set the slot at all.
func (l *Loader) resolveSlot(ctx context.Context, src *Source, slot RefSlot) (*resolvedSlot, error) {
	raw, present := getPath(src.Body, slot.Path)
	if !present {
		return nil, nil
	}

	var values []any
	if slot.List {
		list, ok := raw.([]any)
		if !ok {
			return nil, fmt.Errorf("%s: `%s` must be a list, got %T", src.Key, slot.Path, raw)
		}
		values = list
	} else {
		values = []any{raw}
	}

	out := resolvedSlot{Ref: slot}
	for i, v := range values {
		expr, err := parseRefValue(l.registry, v)
		if err != nil {
			return nil, fmt.Errorf("%s: `%s`[%d]: %w", src.Key, slot.Path, i, err)
		}
		if expr.isPassthrough() {
			out.Entries = append(out.Entries, slotEntry{Value: expr.Passthrough})
			continue
		}
		entries, err := l.resolveExpr(ctx, src, slot, expr)
		if err != nil {
			return nil, err
		}
		out.Entries = append(out.Entries, entries...)
	}
	return &out, nil
}

// resolveExpr turns one reference into slot entries. The reference is an
// inline ID, a URL or a path.
func (l *Loader) resolveExpr(ctx context.Context, src *Source, slot RefSlot, expr refExpr) ([]slotEntry, error) {
	where := fmt.Sprintf("%s: `%s`", src.Key, slot.Path)

	switch {
	case expr.InlineID != "":
		kind, _ := l.registry.KindForID(expr.InlineID)
		if !slices.Contains(slot.To, kind) {
			return nil, wrongKindError(where, expr.InlineID, kind, slot)
		}
		// An inline ID is an unmanaged resource: no file, no lockfile entry,
		// no dependency edge. We take it at face value.
		value, err := withExtra(inlineIDValue(l.registry, slot, expr.InlineID), expr.Extra)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", where, err)
		}
		return []slotEntry{{Value: value}}, nil

	case expr.URL != "":
		fetched, err := l.loadURL(ctx, expr.URL)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", where, err)
		}
		if !slices.Contains(slot.To, fetched.Kind) {
			return nil, wrongKindError(where, expr.URL, fetched.Kind, slot)
		}
		return []slotEntry{{TargetKey: fetched.Key, Pinned: expr.Pinned, Extra: expr.Extra}}, nil

	default:
		return l.resolvePathRef(ctx, src, slot, expr, where)
	}
}

// resolvePathRef loads every file a path or glob reference names.
func (l *Loader) resolvePathRef(ctx context.Context, src *Source, slot RefSlot, expr refExpr, where string) ([]slotEntry, error) {
	pattern := src.resolve(expr.Path)

	matches, err := l.matchRef(pattern, slot)
	if err != nil {
		return nil, fmt.Errorf("%s: %q: %w", where, expr.Path, err)
	}
	if len(matches) == 0 {
		return nil, fmt.Errorf("%s: %q matched nothing on disk", where, expr.Path)
	}

	var entries []slotEntry
	for _, m := range matches {
		target, err := l.loadOne(ctx, m, modeNamed)
		if err != nil {
			return nil, fmt.Errorf("%s → %s: %w", where, l.keyFor(m), err)
		}
		if target == nil {
			continue
		}
		if !slices.Contains(slot.To, target.Kind) {
			return nil, wrongKindError(where, target.Key, target.Kind, slot)
		}
		if target.Key == src.Key {
			return nil, fmt.Errorf("%s: a resource cannot reference itself by path; write the API's own self-reference form instead", where)
		}
		entries = append(entries, slotEntry{TargetKey: target.Key, Pinned: expr.Pinned, Extra: expr.Extra})
	}
	return entries, nil
}

// matchRef expands a reference pattern.
//
// A glob is a filter, not an assertion: `../skills/*` means "every skill in
// there", so a stray README or a non-skill subdirectory is skipped rather than
// turned into a confusing parse error. A path named outright is an assertion
// and must resolve, so it is passed through and left to fail loudly.
func (l *Loader) matchRef(pattern string, slot RefSlot) ([]string, error) {
	if !isGlob(pattern) {
		if _, err := os.Stat(pattern); err != nil {
			return nil, err
		}
		return []string{pattern}, nil
	}

	candidates, err := filepath.Glob(pattern)
	if err != nil {
		return nil, fmt.Errorf("bad pattern: %w", err)
	}
	var out []string
	for _, c := range candidates {
		if l.couldBe(c, slot.To) {
			out = append(out, c)
		}
	}
	slices.Sort(out)
	return out, nil
}

// couldBe reports whether a glob match plausibly denotes one of the given
// kinds.
func (l *Loader) couldBe(path string, allowed []Kind) bool {
	st, err := os.Stat(path)
	if err != nil {
		return false
	}
	probe := l.resourcePath(path)
	if kind, ok := l.registry.schema.IsResourceDir(probe); ok {
		return slices.Contains(allowed, kind)
	}
	if st.IsDir() {
		return false
	}
	c, err := l.open(path)
	if err != nil {
		return false
	}
	kind, err := l.registry.schema.Classify(c, true)
	if err != nil {
		// The file is recognisably a declaration with something wrong in it.
		// Keep it in the match set so loading reports the problem, instead of
		// letting the glob quietly drop it.
		return true
	}
	return slices.Contains(allowed, kind)
}

// loadURL materializes a URL-referenced resource as a Source keyed by the URL.
func (l *Loader) loadURL(ctx context.Context, rawURL string) (*Source, error) {
	if existing, ok := l.sources[rawURL]; ok {
		return existing, nil
	}
	if l.fetcher == nil {
		return nil, fmt.Errorf("%s: remote resources are not available here", rawURL)
	}
	fetched, err := l.fetcher.Fetch(ctx, rawURL, l.pins[rawURL])
	if err != nil {
		return nil, err
	}

	src, err := l.readDirSource(fetched.Dir)
	if err != nil {
		return nil, err
	}
	src.Key = rawURL
	src.URL = rawURL
	src.Pin = fetched.URLPin
	l.sources[rawURL] = src
	// A URL-sourced resource is a leaf: its own references are never
	// followed, so it has no slots and depends on nothing.
	l.slots[rawURL] = nil
	return src, nil
}

// inlineIDValue renders an inline-ID reference. There is no version to pin
// because we never fetched the resource.
func inlineIDValue(r *Registry, slot RefSlot, id string) any {
	kind, _ := r.KindForID(id)
	return slot.encode(Target{Kind: kind, ID: id, Known: true})
}

// wrongKindError refuses a reference whose target is not a kind the slot takes.
func wrongKindError(where, target string, kind Kind, slot RefSlot) error {
	return fmt.Errorf("%s: %s is a %s, but this field takes one of: %s", where, target, kind, kindNames(slot.To))
}

// kindNames joins kinds for an error message, as in "agent or skill".
func kindNames(kinds []Kind) string {
	names := make([]string, len(kinds))
	for i, k := range kinds {
		names[i] = string(k)
	}
	return strings.Join(names, " or ")
}
