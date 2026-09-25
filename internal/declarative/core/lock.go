package core

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

// lockfileSchemaVersion guards against a future format change silently
// corrupting state written by an older CLI.
const lockfileSchemaVersion = 1

// LockEntry is the recorded outcome of the last successful apply of one
// resource. It is the only thing standing between "this file changed" and "the
// resource changed", so every field here exists to answer a specific question
// the next plan has to ask.
type LockEntry struct {
	Kind Kind   `json:"kind"`
	ID   string `json:"id"`
	// Version is the resource version at last apply: an integer for some
	// kinds, an epoch string for others, empty for kinds that don't version.
	// Kept as a string so both round-trip losslessly.
	Version string `json:"version,omitempty"`
	// Hash fingerprints the desired state we last sent. A mismatch means the
	// local files changed.
	Hash string `json:"hash"`
	// RemoteHash fingerprints the normalized remote object as it looked right
	// after that apply. A mismatch means somebody changed it out of band —
	// this is the only drift signal available for kinds with no version field.
	RemoteHash string `json:"remote_hash,omitempty"`
	// Revision and Subpath are the commit a URL-sourced resource was pinned
	// to and where in it the resource lives.
	Revision string `json:"revision,omitempty"`
	Subpath  string `json:"subpath,omitempty"`

	// unknown holds fields written by a newer CLI. They are preserved on save:
	// the lockfile is committed and shared, so an older CLI running one apply
	// must not silently strip state a colleague's newer one depends on.
	unknown map[string]json.RawMessage
}

// lockEntryFields is LockEntry without its methods, so the custom marshalers
// can defer to the standard ones for the known fields.
type lockEntryFields LockEntry

// lockEntryKeys are the json names of LockEntry's known fields; keep them in
// step with its tags. Any other key in an entry is kept in unknown, and a key
// missing here would be read into unknown and would override the field on save.
var lockEntryKeys = []string{"kind", "id", "version", "hash", "remote_hash", "revision", "subpath"}

// UnmarshalJSON decodes the known fields and keeps any others in unknown.
func (e *LockEntry) UnmarshalJSON(data []byte) error {
	if err := json.Unmarshal(data, (*lockEntryFields)(e)); err != nil {
		return err
	}
	var all map[string]json.RawMessage
	if err := json.Unmarshal(data, &all); err != nil {
		return err
	}
	for _, k := range lockEntryKeys {
		delete(all, k)
	}
	if len(all) > 0 {
		e.unknown = all
	}
	return nil
}

// MarshalJSON writes the known fields and merges the unknown ones back in. It
// has a value receiver so that a LockEntry marshals the same way whether or
// not it is addressable.
func (e LockEntry) MarshalJSON() ([]byte, error) {
	known, err := json.Marshal(lockEntryFields(e))
	if err != nil || len(e.unknown) == 0 {
		return known, err
	}
	// Splice the unknown fields back in. Going through a map sorts every key,
	// which is still deterministic — and only ever happens to an entry a newer
	// CLI has touched.
	merged := map[string]json.RawMessage{}
	if err := json.Unmarshal(known, &merged); err != nil {
		return nil, err
	}
	maps.Copy(merged, e.unknown)
	return json.Marshal(merged)
}

// Origin says where a lockfile's resources live: which API host and, when
// the credentials that wrote it knew, which organization and workspace. It is
// recorded on the first write so a later run with different credentials can
// be refused before it mistakes "not yours" for "deleted".
type Origin struct {
	BaseURL        string `json:"base_url,omitempty"`
	OrganizationID string `json:"organization_id,omitempty"`
	WorkspaceID    string `json:"workspace_id,omitempty"`
}

// lockfileDoc is the on-disk shape.
type lockfileDoc struct {
	SchemaVersion int                   `json:"version"`
	Origin        *Origin               `json:"origin,omitempty"`
	Resources     map[string]*LockEntry `json:"resources"`
}

// Lockfile is the on-disk state file.
type Lockfile struct {
	Path      string
	Resources map[string]*LockEntry
	// Origin is where these resources live; nil until the first apply
	// records it, and in lockfiles written before it existed.
	Origin *Origin

	// existed records whether the file was on disk when we loaded it, so a
	// first apply can say it is creating one rather than silently
	// inventing state.
	existed bool
}

func newLockfile(path string) *Lockfile {
	return &Lockfile{Path: path, Resources: map[string]*LockEntry{}}
}

// Existed reports whether the lockfile was on disk when it was loaded.
func (lf *Lockfile) Existed() bool { return lf.existed }

// Root is the directory resource keys are relative to.
func (lf *Lockfile) Root() string { return filepath.Dir(lf.Path) }

// Get returns the entry recorded for key, if any.
func (lf *Lockfile) Get(key string) (*LockEntry, bool) {
	e, ok := lf.Resources[key]
	return e, ok
}

// Keys returns the recorded resource keys in sorted order.
func (lf *Lockfile) Keys() []string {
	return slices.Sorted(maps.Keys(lf.Resources))
}

// Pins returns what was recorded per URL key, for seeding the loader.
func (lf *Lockfile) Pins() map[string]URLPin {
	out := map[string]URLPin{}
	for k, e := range lf.Resources {
		if e.Revision != "" {
			out[k] = URLPin{Revision: e.Revision, Subpath: e.Subpath}
		}
	}
	return out
}

// FindLockfile walks up from start looking for a lockfile, stopping after the
// first directory that holds .git. When none is found it returns an empty
// lockfile rooted at start, which is what a first apply wants.
func FindLockfile(r *Registry, start string) (*Lockfile, error) {
	abs, err := filepath.Abs(start)
	if err != nil {
		return nil, err
	}
	for dir := abs; ; dir = filepath.Dir(dir) {
		candidate := filepath.Join(dir, r.LockfileName())
		if st, err := os.Stat(candidate); err == nil && !st.IsDir() {
			return LoadLockfile(r, candidate)
		}
		// Stop at the filesystem root, and don't wander out of the repository
		// looking for state that belongs to somebody else's project.
		if filepath.Dir(dir) == dir || isRepoRoot(dir) {
			break
		}
	}
	return newLockfile(filepath.Join(abs, r.LockfileName())), nil
}

func isRepoRoot(dir string) bool {
	_, err := os.Stat(filepath.Join(dir, ".git"))
	return err == nil
}

// LoadLockfile reads a lockfile from an explicit path. A missing file is not an
// error; a malformed one is.
func LoadLockfile(r *Registry, path string) (*Lockfile, error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	lf := newLockfile(abs)

	data, err := os.ReadFile(abs)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return lf, nil
		}
		return nil, err
	}
	lf.existed = true
	if len(bytes.TrimSpace(data)) == 0 {
		return lf, nil
	}

	var doc lockfileDoc
	if err := json.Unmarshal(data, &doc); err != nil {
		return nil, fmt.Errorf("%s: %w", path, err)
	}
	if doc.SchemaVersion > lockfileSchemaVersion {
		return nil, fmt.Errorf("%s: written by a newer ant (lockfile version %d, this build understands %d)", path, doc.SchemaVersion, lockfileSchemaVersion)
	}
	lf.Origin = doc.Origin
	for key, entry := range doc.Resources {
		if err := checkLockEntry(r, key, entry); err != nil {
			return nil, fmt.Errorf("%s: entry %q: %w", path, key, err)
		}
		if keyEscapesRoot(filepath.Dir(abs), key) {
			return nil, fmt.Errorf("%s: entry %q: key resolves outside the lockfile directory", path, key)
		}
		lf.Resources[key] = entry
	}
	return lf, nil
}

func checkLockEntry(r *Registry, key string, e *LockEntry) error {
	switch {
	case !isResourceKey(key):
		return fmt.Errorf("resource keys start with ./, ../ or https://")
	case e == nil:
		return fmt.Errorf("expected an object")
	case e.ID == "":
		return fmt.Errorf("missing `id`")
	case e.Kind == "":
		return fmt.Errorf("missing `kind`")
	case !r.Valid(e.Kind):
		return fmt.Errorf("unknown kind %q", e.Kind)
	}
	return nil
}

func isResourceKey(key string) bool {
	return strings.HasPrefix(key, "./") || strings.HasPrefix(key, "../") || isURL(key)
}

// keyEscapesRoot reports whether a relative lockfile key, joined to the
// lockfile's directory, leaves that directory. URL keys are not paths.
func keyEscapesRoot(root, key string) bool {
	if isURL(key) {
		return false
	}
	root = filepath.Clean(root)
	cleaned := filepath.Clean(filepath.Join(root, filepath.FromSlash(strings.TrimPrefix(key, "./"))))
	rel, err := filepath.Rel(root, cleaned)
	if err != nil {
		return true
	}
	return rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator))
}

// Save writes the lockfile atomically. encoding/json sorts map keys and emits
// struct fields in declaration order, so the output is byte-stable across runs
// — a lockfile that reorders itself produces noisy commits and pointless merge
// conflicts.
func (lf *Lockfile) Save() error {
	data, err := json.MarshalIndent(lockfileDoc{SchemaVersion: lockfileSchemaVersion, Origin: lf.Origin, Resources: lf.Resources}, "", "  ")
	if err != nil {
		return err
	}
	return writeFileAtomic(lf.Path, append(data, '\n'), 0o644)
}

// writeFileAtomic writes via a temp file in the same directory so a crash or a
// full disk can't truncate existing state.
func writeFileAtomic(path string, data []byte, perm os.FileMode) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	tmp, err := os.CreateTemp(dir, "."+filepath.Base(path)+".tmp-*")
	if err != nil {
		return err
	}
	tmpName := tmp.Name()
	defer os.Remove(tmpName)

	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Sync(); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	if err := os.Chmod(tmpName, perm); err != nil {
		return err
	}
	return os.Rename(tmpName, path)
}
