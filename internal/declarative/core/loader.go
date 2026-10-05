package core

import (
	"cmp"
	"context"
	"errors"
	"fmt"
	"io/fs"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

// Loader turns a set of paths into a closed set of Sources: every reference
// reachable from the named paths is loaded too, so naming one file is enough
// to apply it together with everything it depends on.
type Loader struct {
	// registry supplies the kinds this loader understands.
	registry *Registry
	// root is the directory lockfile keys are relative to.
	root string
	// fetcher resolves URL-sourced resources. Nil disables URL support.
	fetcher URLFetcher

	sources map[string]*Source
	slots   map[string][]resolvedSlot
	// pins carries what the lockfile recorded per URL, so a floating ref
	// keeps resolving to the same commit until --upgrade.
	pins map[string]URLPin
}

// NewLoader returns a loader for the kinds in registry, keying resources
// relative to root. A nil fetcher disables URL references.
func NewLoader(registry *Registry, root string, fetcher URLFetcher) *Loader {
	return &Loader{
		registry: registry,
		root:     root,
		fetcher:  fetcher,
		sources:  map[string]*Source{},
		slots:    map[string][]resolvedSlot{},
		pins:     map[string]URLPin{},
	}
}

// Pin seeds what a previous apply recorded for each URL. Without this every
// apply would re-resolve a branch URL and churn the resource on every push to
// that branch.
func (l *Loader) Pin(pins map[string]URLPin) {
	maps.Copy(l.pins, pins)
}

// Sources returns everything loaded, keyed by lockfile key.
func (l *Loader) Sources() map[string]*Source { return l.sources }

// keyFor renders an absolute path as the canonical lockfile key: always
// relative, always prefixed, so it round-trips through the lockfile's key
// validation. filepath.Rel only fails across Windows volumes, where there is
// no relative path to give — "./" keeps the result loadable rather than
// writing state that the next run refuses to read. The key is also how a path
// is named in error messages, so what an error names matches the plan.
func (l *Loader) keyFor(abs string) string {
	rel, err := filepath.Rel(l.root, abs)
	if err != nil {
		rel = filepath.Base(abs)
	}
	rel = filepath.ToSlash(rel)
	if strings.HasPrefix(rel, "../") {
		return rel
	}
	return "./" + rel
}

// resourcePath maps a marker file to the directory it stands for, so a
// resource reached either way has one key.
func (l *Loader) resourcePath(path string) string {
	if dir, ok := l.registry.schema.ResourceDirOf(path); ok {
		return dir
	}
	return path
}

// pathForKey is keyFor's inverse.
func (l *Loader) pathForKey(key string) string {
	if isURL(key) {
		return key
	}
	return filepath.Join(l.root, filepath.FromSlash(strings.TrimPrefix(key, "./")))
}

func isURL(s string) bool {
	return strings.HasPrefix(s, "http://") || strings.HasPrefix(s, "https://")
}

func isGlob(s string) bool { return strings.ContainsAny(s, "*?[") }

// Add loads the given paths (files, directories or globs) and everything they
// reference. The context bounds URL fetches.
func (l *Loader) Add(ctx context.Context, paths []string) error {
	for _, p := range paths {
		matches, err := expandArg(p)
		if err != nil {
			return err
		}
		for _, m := range matches {
			if err := l.loadPath(ctx, m); err != nil {
				return err
			}
		}
	}
	return nil
}

// AddKeys loads resources named by lockfile keys, skipping any whose file has
// been removed — those become deletions, which the planner handles from the
// lockfile alone.
func (l *Loader) AddKeys(ctx context.Context, keys []string) error {
	for _, key := range keys {
		if isURL(key) {
			// A URL-sourced resource has no file of its own, so it is tracked
			// only for as long as something references it. Fetching it here
			// would make it permanently un-prunable — and would hit the network
			// on every apply for something nothing uses any more.
			continue
		}
		abs := l.pathForKey(key)
		if _, err := os.Stat(abs); err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				continue
			}
			return err
		}
		if err := l.loadPath(ctx, abs); err != nil {
			return err
		}
	}
	return nil
}

// expandArg turns a command-line argument into concrete paths. The shell
// usually does this, but quoted globs and globs from a config file reach us
// intact, and Windows shells don't glob at all.
func expandArg(arg string) ([]string, error) {
	abs, err := filepath.Abs(arg)
	if err != nil {
		return nil, err
	}
	if !isGlob(arg) {
		if _, err := os.Stat(abs); err != nil {
			return nil, fmt.Errorf("%s: %w", arg, err)
		}
		return []string{abs}, nil
	}
	matches, err := filepath.Glob(abs)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", arg, err)
	}
	if len(matches) == 0 {
		return nil, fmt.Errorf("%s: no files match", arg)
	}
	slices.Sort(matches)
	return matches, nil
}

// loadPath loads a path the user named: a file, a resource directory, or a
// directory to walk.
func (l *Loader) loadPath(ctx context.Context, abs string) (err error) {
	// A child visited before a failure may point back to an unresolved ancestor.
	// Restore the whole dependency closure while retaining earlier loaded paths.
	sources, slots := maps.Clone(l.sources), maps.Clone(l.slots)
	defer func() {
		if err != nil {
			clear(l.sources)
			maps.Copy(l.sources, sources)
			clear(l.slots)
			maps.Copy(l.slots, slots)
		}
	}()
	info, err := os.Stat(abs)
	if err != nil {
		return err
	}
	if info.IsDir() {
		if _, ok := l.registry.schema.IsResourceDir(abs); !ok {
			return l.walkDir(ctx, abs)
		}
	}
	_, err = l.loadOne(ctx, abs, modeNamed)
	return err
}

// walkDir picks up every resource under a directory. It is deliberately
// lenient: Classify is told the file was walked, so the Schema can pass over
// a file it does not recognise instead of failing on it.
func (l *Loader) walkDir(ctx context.Context, root string) error {
	return l.walkCandidates(root, func(path string) error {
		_, err := l.loadOne(ctx, path, modeWalked)
		return err
	})
}

// loadMode distinguishes a path the user named from one we found by walking a
// directory. Named paths are strict — if we can't tell what a file is, that's
// an error. Walked paths are lenient, because a repo will contain READMEs and
// CI config we must not mistake for resource definitions.
type loadMode int

const (
	modeNamed loadMode = iota
	modeWalked
)

// loadOne loads a single file or resource directory and recurses into its
// references. It is idempotent: the second call for a key returns the cached
// Source.
func (l *Loader) loadOne(ctx context.Context, abs string, mode loadMode) (*Source, error) {
	abs, err := filepath.Abs(abs)
	if err != nil {
		return nil, err
	}
	// Resolve to the canonical key before the cache check so that a resource
	// reached by its directory and by its defining file dedupe to one resource.
	key := l.keyFor(l.resourcePath(abs))
	if existing, ok := l.sources[key]; ok {
		return existing, nil
	}

	src, err := l.readSource(abs, mode)
	if err != nil || src == nil {
		return nil, err
	}
	// Included data becomes part of the body before anything else looks at
	// it, so references, hashing and the plan never know it lived elsewhere.
	if err := l.expandIncludes(src); err != nil {
		return nil, err
	}

	// Cached before resolving references so a cycle terminates here on the
	// second visit rather than recursing; TopoOrder reports it properly.
	l.sources[key] = src
	if err := l.resolveSlots(ctx, src); err != nil {
		return nil, err
	}
	return src, nil
}

// skipDirs are never descended into during a directory walk.
var skipDirs = map[string]bool{
	".git": true, "node_modules": true, "vendor": true,
	"__pycache__": true, ".venv": true, "venv": true, "dist": true, "build": true,
}

// walkCandidates visits everything under root that could be a resource: each
// resource directory (without descending into it) and each file with an
// extension the loader reads, skipping tooling directories and the lockfile.
// Whether a candidate really is a resource is for the visitor to decide.
func (l *Loader) walkCandidates(root string, visit func(path string) error) error {
	return filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			base := filepath.Base(path)
			if path != root && (skipDirs[base] || strings.HasPrefix(base, ".")) {
				return fs.SkipDir
			}
			if _, ok := l.registry.schema.IsResourceDir(path); ok {
				if err := visit(path); err != nil {
					return err
				}
				// Its own files are its payload, not more resources.
				return fs.SkipDir
			}
			return nil
		}
		switch strings.ToLower(filepath.Ext(path)) {
		case ".md", ".markdown", ".yml", ".yaml", ".json":
		default:
			return nil
		}
		if l.registry.isLockfileName(filepath.Base(path)) {
			return nil
		}
		return visit(path)
	})
}

// Found is a resource Discover noticed on disk.
type Found struct {
	Key  string
	Kind Kind
}

// Discover lists the resources under dir that this loader has not loaded,
// by the same conservative rules a directory walk applies — but without
// building bodies or following references, so one malformed file cannot
// hide its neighbours. It is for offering untracked files to the user, who
// decides; nothing here is added to the loader.
func (l *Loader) Discover(dir string) ([]Found, error) {
	var found []Found
	err := l.walkCandidates(dir, func(path string) error {
		key := l.keyFor(l.resourcePath(path))
		if _, loaded := l.sources[key]; loaded {
			return nil
		}
		if kind, ok := l.registry.schema.IsResourceDir(path); ok {
			found = append(found, Found{key, kind})
			return nil
		}
		c, err := l.open(path)
		if err != nil || c.unsupportedExt != nil {
			return nil
		}
		if kind, err := l.registry.schema.Classify(c, true); err == nil && kind != "" {
			found = append(found, Found{key, kind})
		}
		return nil
	})
	slices.SortFunc(found, func(a, b Found) int {
		return cmp.Or(cmp.Compare(l.registry.rank(a.Kind), l.registry.rank(b.Kind)), strings.Compare(a.Key, b.Key))
	})
	return found, err
}
