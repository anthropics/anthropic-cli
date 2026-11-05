package core

import "slices"

// Kind is a resource type. Values are defined by the domain, not here.
type Kind string

// Destroy names how a kind is removed, in the two tenses plan and apply
// output need.
type Destroy struct {
	Verb string // "archive"
	Past string // "archived"
}

// KindSpec is everything core needs to know about one kind. A domain builds
// one per kind and hands the set over as a Registry; keeping them together
// means a reader sees every fact about a kind in one place, and adding a
// kind is one literal rather than an edit in each phase of the engine.
type KindSpec struct {
	// Kind is the kind this spec describes.
	Kind Kind

	// IDPrefix is the leading token of the API's tagged IDs ("agent" in
	// "agent_01…"), used to tell an inline ID reference from a file path.
	IDPrefix string

	// VersionField is the response field holding the resource's version, if it
	// has one, which references pin and updates echo back. Empty means the
	// kind is unversioned. Drift is detected from Stamp and a fingerprint of the
	// whole object, for versioned and unversioned kinds alike.
	VersionField string
	// VersionIsInt renders the version as a number rather than a string when
	// it is sent back to the API.
	VersionIsInt bool
	// UpdateNeedsVersion means update requires the caller to echo the current
	// version back, under `version` in the request body, which is how the API
	// rejects a racing writer.
	UpdateNeedsVersion bool

	// Stamp returns what the server says about when the resource last
	// changed. While it equals the value recorded at the last apply, nobody
	// has changed the resource. Unless StampMovesOnlyOnEdit is set, a moved
	// stamp proves nothing, since the server can move it on its own, so the
	// fingerprint then decides. Nil means the kind has none.
	Stamp func(remote map[string]any) string
	// StampMovesOnlyOnEdit means the server moves the Stamp only when someone
	// edits the resource, as with a version, so a moved Stamp is itself a change.
	StampMovesOnlyOnEdit bool
	// StampName is the name of the Stamp's field, such as "version", for the
	// message that reports a moved Stamp. Set it with StampMovesOnlyOnEdit.
	StampName string

	// Destroy is what removal actually does, for plan output: most APIs
	// offer archive, delete, or only one of the two.
	Destroy Destroy

	// Fields describes the resource's fields, one entry per field, keyed by a
	// dotted path into the body. Everything the reconciler needs to know about
	// a field lives in its entry rather than being spread across parallel
	// lists, so one line says everything about a field such as `description`.
	Fields Fields

	// Build turns a parsed file into a request body. Required for every kind a
	// Schema can classify from a file. Leave it nil for a kind whose content
	// travels as a Payload: core then refuses to load one from a file, and
	// reads a write's version from the version object the upload answers with.
	Build func(c *Candidate) (map[string]any, error)

	// Derived from Fields by NewRegistry.
	computed      []string
	writeOnly     []string
	immutable     []string
	clearable     []string
	summary       []string
	refs          []RefSlot
	includes      []string
	shorthand     map[string]string
	matchBy       map[string][]string
	archivedField string
	metadataField string
	payloadBacked bool
}

// Field describes one field of a resource, after the manner of a Terraform
// attribute. The zero value is an ordinary field: sent as written, compared as
// written, preserved when omitted.
type Field struct {
	// Computed marks a value the server owns and the config does not set —
	// timestamps, IDs, status. It is stripped before comparing, so a missing
	// Computed marking shows up as drift that never clears.
	Computed bool
	// WriteOnly marks a value the server never returns. It cannot be compared,
	// and it must never be printed. A "[]" segment in the path means "every
	// element of this array".
	WriteOnly bool
	// Immutable marks a value fixed when the resource is created. Changing it
	// needs a replace, which core refuses to do on the user's behalf.
	Immutable bool
	// Clearable marks a value that deleting from the file should actually
	// remove, by sending an explicit null. Without it the field follows the
	// API's usual omit-to-preserve rule, which is the safe default: a server
	// that fills in its own value would otherwise be fought with forever.
	Clearable bool
	// Archived marks the server-owned field whose presence means the resource
	// was archived out of band. Archiving is a one-way door, so the planner
	// offers a replacement rather than an update.
	Archived bool
	// Shorthand names the sub-field a bare scalar stands for when the API
	// accepts either form and always answers with the object:
	// `model: x` and `model: {id: x, …}` then compare equal.
	Shorthand string
	// Summary marks a field worth showing under a create in the plan. The
	// whole body would bury the plan in prose; these identify the resource.
	Summary bool
	// MatchBy names the members that identify an element of this list field,
	// most telling first, so a plan can pair elements up before diffing them.
	// Unset means id, then name; position is always the last resort.
	MatchBy []string

	// Metadata marks a string-to-string bag the server patches key by key, so
	// core sends null for keys the file no longer declares.
	Metadata bool
	// Ref marks a field holding a reference to another resource, written in
	// the config as a file path and sent to the API as an ID.
	Ref *Ref
	// Include marks a list field whose entries may name a data file — YAML or
	// JSON, by path or glob relative to the declaring file — whose contents
	// are spliced into the list where the entry stood. It is for lists of
	// hand-written API objects that some other tool generates (a framework
	// emitting the custom tool specs it executes, say), so the definition
	// stays declarative without being copied in by hand. The file is data,
	// not a resource: it gets no lockfile entry, but what it contributes is
	// part of the body, so editing it plans an update like any other edit.
	Include bool
}

// Fields describes a kind's fields, keyed by dotted path.
type Fields map[string]Field

// Ref describes a reference field: what it may point at, and how the resolved
// target is written into the request body.
type Ref struct {
	// To restricts what the reference may name.
	To []Kind
	// List marks a field holding several references rather than one.
	List bool
	// As renders a resolved target into the body. Nil writes the bare ID.
	As func(Target) any
}

// encode renders a resolved target into the body: through As when the
// field has one, otherwise as the bare ID.
func (r Ref) encode(t Target) any {
	if r.As == nil {
		return t.idOrPlaceholder()
	}
	return r.As(t)
}

// paths returns the fields matching a predicate, in a stable order.
func (f Fields) paths(pred func(Field) bool) []string {
	var out []string
	for path, field := range f {
		if pred(field) {
			out = append(out, path)
		}
	}
	slices.Sort(out)
	return out
}

// refSlots turns the reference fields into the slots the loader resolves.
func (f Fields) refSlots() []RefSlot {
	var out []RefSlot
	for _, path := range f.paths(func(x Field) bool { return x.Ref != nil }) {
		out = append(out, RefSlot{Path: path, Ref: *f[path].Ref})
	}
	return out
}

// Clearable lists the fields that deleting from a file actually removes.
func (s KindSpec) Clearable() []string { return s.clearable }

// RefSlots lists the reference fields, in a stable order.
func (s KindSpec) RefSlots() []RefSlot { return s.refs }

// stamp returns obj's Stamp, or "" if the kind has none.
func (s KindSpec) stamp(obj map[string]any) string {
	if s.Stamp == nil || obj == nil {
		return ""
	}
	return s.Stamp(obj)
}

// IsArchived reports whether the server's copy carries the kind's tombstone.
// Archiving is a one-way door, so the only way forward is a new resource.
func (s KindSpec) IsArchived(remote map[string]any) bool {
	if s.archivedField == "" {
		return false
	}
	return !IsEmptyValue(remote[s.archivedField])
}
