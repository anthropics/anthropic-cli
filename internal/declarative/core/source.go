package core

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/goccy/go-yaml"
)

// Source is one resource as declared on disk, or fetched from a URL.
// Body holds the declared fields verbatim; references inside it are still
// paths at this stage and get resolved to IDs later, once the graph is ordered.
type Source struct {
	// Key is the stable identity used in the lockfile and in plan output:
	// a root-relative slash path prefixed with "./", or the verbatim URL for
	// a URL-sourced resource.
	Key  string
	Kind Kind

	// Path is the absolute path of the defining file. Empty when the resource
	// came from a URL.
	Path string
	// Dir is set when the resource is declared by a directory rather than a
	// file.
	Dir string
	// Payload is the non-JSON content of a directory-backed resource.
	Payload Payload
	// URL is set when the resource was fetched from a URL rather than read
	// from the working tree, and Pin records what it resolved to.
	URL string
	Pin URLPin

	// Body is the declared request body, minus apply-only directives.
	Body map[string]any
}

// resolve makes a path written inside this resource absolute. A relative
// path is taken from the directory that declares the resource, or from the
// declaring file's directory.
func (s *Source) resolve(p string) string {
	if filepath.IsAbs(p) {
		return p
	}
	base := filepath.Dir(s.Path)
	if s.Dir != "" {
		base = s.Dir
	}
	return filepath.Join(base, p)
}

// Candidate is a file read and decoded once. The sniff that decides whether to
// take a file and the load that builds its Source look at the same bytes, so
// they cannot disagree and the file is not parsed twice.
type Candidate struct {
	// Path is the file's absolute path.
	Path string
	// Markdown reports whether the file has a prose body as well as fields.
	Markdown bool
	// Prose is the markdown after the frontmatter, empty for a YAML file.
	Prose []byte
	// Fields are the declared fields: the frontmatter of a markdown file, or
	// the whole document of a YAML one. The `type` field is removed once the
	// file is classified.
	Fields map[string]any
	// Name is the default for a resource's `name`: the filename without its
	// extension, unless the Schema refines it while classifying.
	Name string
	// HasFrontmatter tells an empty frontmatter block apart from no
	// frontmatter at all: only the latter means the file is prose.
	HasFrontmatter bool
	// DecodeErr holds a frontmatter or YAML failure. It is fatal for a path the
	// user named but merely disqualifying for one we walked into, so it travels
	// with the Candidate instead of aborting the read.
	DecodeErr error
	// unsupportedExt is also recorded rather than returned. A glob classifies
	// whatever it matched, so a file with an extension no builder reads is
	// still judged by its fields, and loading one that looks like a
	// declaration is what fails.
	unsupportedExt error
}

// readSource builds the Source declared at one path, a file or a resource
// directory, without caching or following references. Kind inference happens
// here so the rest of the engine only ever sees a typed Source. A nil Source
// and a nil error mean a walked path we declined to identify.
func (l *Loader) readSource(abs string, mode loadMode) (*Source, error) {
	if _, ok := l.registry.schema.IsResourceDir(abs); ok {
		return l.readDirSource(abs)
	}
	if dir, ok := l.registry.schema.ResourceDirOf(abs); ok {
		return l.readDirSource(dir)
	}
	info, err := os.Stat(abs)
	if err != nil {
		return nil, err
	}
	if info.IsDir() {
		// Walking a directory for the resources inside it happens in walkDir,
		// before we get here.
		return nil, fmt.Errorf("%s is a directory but declares no resource", l.keyFor(abs))
	}

	c, err := l.open(abs)
	if err != nil {
		return nil, err
	}
	if c.unsupportedExt != nil {
		return nil, c.unsupportedExt
	}
	kind, err := l.registry.schema.Classify(c, mode == modeWalked)
	if err != nil || kind == "" {
		return nil, err
	}
	delete(c.Fields, "type")

	build := l.registry.specOrZero(kind).Build
	if build == nil {
		return nil, fmt.Errorf("%s: %s resources cannot be declared in a file", l.keyFor(c.Path), kind)
	}
	body, err := build(c)
	if err != nil {
		return nil, err
	}
	return &Source{Key: l.keyFor(c.Path), Kind: kind, Path: c.Path, Body: body}, nil
}

// readDirSource builds the Source declared by a resource directory.
func (l *Loader) readDirSource(absDir string) (*Source, error) {
	if _, ok := l.registry.schema.IsResourceDir(absDir); !ok {
		return nil, fmt.Errorf("%s is a directory but declares no resource", l.keyFor(absDir))
	}
	kind, body, payload, err := l.registry.schema.LoadDir(absDir)
	if err != nil {
		return nil, err
	}
	return &Source{
		Key:     l.keyFor(absDir),
		Kind:    kind,
		Dir:     absDir,
		Payload: payload,
		Body:    body,
	}, nil
}

// open reads and decodes one file. It fails only when the bytes are
// unreachable; everything else is recorded on the Candidate for classify and
// the builders to judge.
func (l *Loader) open(abs string) (*Candidate, error) {
	content, err := os.ReadFile(abs)
	if err != nil {
		return nil, err
	}
	c := &Candidate{Path: abs, Name: defaultName(abs)}
	switch strings.ToLower(filepath.Ext(abs)) {
	case ".md", ".markdown":
		c.Markdown = true
		front, body, err := splitFrontmatter(content)
		if err != nil {
			c.DecodeErr = fmt.Errorf("%s: %w", l.keyFor(abs), err)
			return c, nil
		}
		c.HasFrontmatter, c.Prose = front != nil, body
		c.Fields, c.DecodeErr = ParseYAMLMap(front, l.keyFor(abs)+" frontmatter")
	case ".yml", ".yaml", ".json":
		c.Fields, c.DecodeErr = ParseYAMLMap(content, l.keyFor(abs))
	default:
		c.unsupportedExt = fmt.Errorf("%s: unsupported extension (expected .md, .yml, .yaml or .json)", l.keyFor(abs))
		c.Fields, _ = ParseYAMLMap(content, l.keyFor(abs))
	}
	return c, nil
}

// defaultName derives a resource name from a filename: "code-reviewer.md"
// becomes "code-reviewer".
func defaultName(absPath string) string {
	base := filepath.Base(absPath)
	return strings.TrimSuffix(base, filepath.Ext(base))
}

const frontmatterFence = "---"

// splitFrontmatter separates a YAML frontmatter block from the markdown body.
// A file with no leading fence is all body and no frontmatter.
func splitFrontmatter(content []byte) (front, body []byte, err error) {
	text := string(content)
	// Tolerate a UTF-8 BOM and leading blank lines before the opening fence.
	text = strings.TrimPrefix(text, "\ufeff")
	trimmed := strings.TrimLeft(text, "\r\n \t")
	if !strings.HasPrefix(trimmed, frontmatterFence) {
		return nil, []byte(text), nil
	}
	rest := trimmed[len(frontmatterFence):]
	// The opening fence must be alone on its line.
	nl := strings.IndexByte(rest, '\n')
	if nl < 0 || strings.TrimSpace(rest[:nl]) != "" {
		return nil, []byte(text), nil
	}
	rest = rest[nl+1:]

	// Find the closing fence at the start of a line. A line from strings.Lines
	// keeps its newline, so the body starts on the line after the fence.
	offset := 0
	for line := range strings.Lines(rest) {
		if strings.TrimRight(line, " \t\r\n") == frontmatterFence {
			return []byte(rest[:offset]), []byte(rest[offset+len(line):]), nil
		}
		offset += len(line)
	}
	return nil, nil, fmt.Errorf("frontmatter opened with %q but never closed", frontmatterFence)
}

// ParseYAMLMap decodes YAML into a map with string keys throughout. Every
// nested map must come out map[string]any, because request bodies are
// eventually run through json.Marshal in canonicalJSON, which rejects
// interface-keyed maps. goccy/go-yaml stringifies even non-string keys, so no
// normalizing pass is needed — TestYAMLKeysAreAlwaysStrings pins that so a
// dependency bump fails there rather than as an opaque marshal error while
// hashing.
func ParseYAMLMap(data []byte, what string) (map[string]any, error) {
	if len(bytes.TrimSpace(data)) == 0 {
		return map[string]any{}, nil
	}
	var raw any
	if err := yaml.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("parsing %s: %w", what, err)
	}
	m, ok := raw.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("parsing %s: expected a mapping, got %T", what, raw)
	}
	return m, nil
}
