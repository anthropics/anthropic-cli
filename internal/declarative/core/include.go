package core

import (
	"bytes"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

// expandIncludes splices data files into the list fields that allow it. An
// entry that is a string names a file (or a glob over files) relative to the
// declaring file; the file's contents — one object, or a list of them — take
// the entry's place, in order. Anything that is not a string is left exactly
// where it was, so hand-written entries and included ones mix freely.
//
// This runs at load time, before references resolve and before anything is
// hashed: the rest of the engine sees an ordinary body, and a change to an
// included file is simply a change to the resource.
func (l *Loader) expandIncludes(src *Source) error {
	for _, field := range l.registry.specOrZero(src.Kind).includes {
		raw, present := getPath(src.Body, field)
		if !present || raw == nil {
			continue
		}
		list, ok := raw.([]any)
		if !ok {
			return fmt.Errorf("%s: `%s` must be a list, got %T", src.Key, field, raw)
		}
		where := fmt.Sprintf("%s: `%s`", src.Key, field)

		out := make([]any, 0, len(list))
		for i, entry := range list {
			path, isPath := entry.(string)
			if !isPath {
				out = append(out, entry)
				continue
			}
			items, err := l.includeFile(src, path)
			if err != nil {
				return fmt.Errorf("%s[%d]: %w", where, i, err)
			}
			out = append(out, items...)
		}
		setPath(src.Body, field, out)
	}
	return nil
}

// includeFile reads the objects one include entry contributes.
func (l *Loader) includeFile(src *Source, path string) ([]any, error) {
	path = strings.TrimSpace(path)
	if path == "" {
		return nil, fmt.Errorf("empty include path")
	}
	if isURL(path) {
		return nil, fmt.Errorf("%q: only files in the working tree can be included", path)
	}

	pattern := src.resolve(path)
	matches, err := matchInclude(pattern)
	if err != nil {
		return nil, fmt.Errorf("%q: %w", path, err)
	}
	if len(matches) == 0 {
		return nil, fmt.Errorf("%q matched nothing on disk", path)
	}

	var out []any
	for _, m := range matches {
		// A data file sitting where declarations live would also be picked
		// up as a resource of its own by a directory walk. Refuse the
		// placement outright rather than have `apply ./agents` quietly
		// create an agent out of a list of tool specs.
		if l.couldBe(m, l.registry.order) {
			return nil, fmt.Errorf("%s would itself be read as a resource declaration where it is; move it into a subdirectory (e.g. ./tools/) or out of the kind directory", l.keyFor(m))
		}
		items, err := readIncluded(m)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", l.keyFor(m), err)
		}
		out = append(out, items...)
	}
	return out, nil
}

// matchInclude expands an include pattern. As with references, a glob is a
// filter that keeps only data files; a path named outright is taken as
// written. Either way an empty result is for the caller to refuse.
func matchInclude(pattern string) ([]string, error) {
	if !isGlob(pattern) {
		info, err := os.Stat(pattern)
		if errors.Is(err, fs.ErrNotExist) {
			return nil, nil
		}
		if err != nil {
			return nil, err
		}
		if info.IsDir() {
			return nil, fmt.Errorf("is a directory; name the files inside it, or use a glob")
		}
		return []string{pattern}, nil
	}
	candidates, err := filepath.Glob(pattern)
	if err != nil {
		return nil, fmt.Errorf("bad pattern: %w", err)
	}
	var out []string
	for _, c := range candidates {
		if !isDataFile(c) {
			continue
		}
		if info, err := os.Stat(c); err == nil && !info.IsDir() {
			out = append(out, c)
		}
	}
	slices.Sort(out)
	return out, nil
}

// isDataFile reports whether path names a file an include can read.
func isDataFile(path string) bool {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".yml", ".yaml", ".json":
		return true
	}
	return false
}

// readIncluded decodes one data file into the objects it contributes: a
// document that is a single mapping contributes one, a sequence contributes
// each of its elements, and every element must itself be a mapping because
// that is what the list it lands in holds.
func readIncluded(path string) ([]any, error) {
	if !isDataFile(path) {
		return nil, fmt.Errorf("unsupported extension (expected .yml, .yaml or .json)")
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	if len(bytes.TrimSpace(data)) == 0 {
		return nil, nil
	}
	doc, err := decodeSingleYAMLDocument(data)
	if err != nil {
		return nil, fmt.Errorf("parsing: %w", err)
	}
	items, ok := doc.([]any)
	if !ok {
		items = []any{doc}
	}
	for i, item := range items {
		if _, ok := item.(map[string]any); !ok {
			return nil, fmt.Errorf("element %d is a %T; every included element must be a mapping", i, item)
		}
	}
	return items, nil
}
