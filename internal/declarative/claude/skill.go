package claude

import (
	"cmp"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"mime"
	"os"
	"path"
	"path/filepath"
	"slices"
	"sort"
	"strings"

	"github.com/anthropics/anthropic-cli/internal/declarative/core"
)

// skillFileName is the file that marks a directory as a skill.
const skillFileName = "SKILL.md"

// loadSkillDir reads a skill directory: its request body and its bundle.
func loadSkillDir(dir string) (map[string]any, core.Payload, error) {
	skillMD := filepath.Join(dir, skillFileName)
	content, err := os.ReadFile(skillMD)
	if err != nil {
		return nil, nil, err
	}
	front, err := skillFrontmatter(content)
	if err != nil {
		return nil, nil, fmt.Errorf("%s %w", skillMD, err)
	}
	fields, err := core.ParseYAMLMap(front, skillMD+" frontmatter")
	if err != nil {
		return nil, nil, err
	}
	files, err := collectSkillFiles(dir)
	if err != nil {
		return nil, nil, err
	}

	// The server derives name and description from SKILL.md itself; the only
	// field sent beside the files is the display name. Older spellings of it
	// (display_title, title) are still accepted, and a skill with none is
	// shown under its `name`, or failing that its directory.
	name := firstString(fields, "name")
	if err := checkSkillName(name); err != nil {
		return nil, nil, fmt.Errorf("%s: %w", skillMD, err)
	}
	displayName := cmp.Or(firstString(fields, "display_name", "display_title", "title"), name, filepath.Base(dir))
	uploadDir := cmp.Or(name, filepath.Base(dir))

	// `name` rides in the body only so the plan can hold it against the name
	// the server has: it is the skill's identity across versions, and a
	// version that changes it is rejected, so a rename must be refused here
	// rather than fail midway through an upload.
	body := map[string]any{"display_name": displayName}
	if name != "" {
		body["name"] = name
	}
	return body, &skillBundle{UploadDir: uploadDir, DisplayName: displayName, Files: files}, nil
}

// maxSkillFrontmatter is the largest frontmatter block the API reads.
const maxSkillFrontmatter = 256_000

// skillFrontmatter returns the YAML between SKILL.md's fences, found the way
// the API finds it, so that a manifest the API would read as having no
// frontmatter is refused before anything uploads, and with the reason. The
// rules are narrower than the ones core reads other markdown by: only
// ASCII-blank lines before the opening fence (no byte-order mark), each fence
// exactly "---" and ending in a newline (no trailing spaces, and not the last
// bytes of the file), and something between them. Line endings inside the
// block come back as "\n" however the file spells them.
func skillFrontmatter(content []byte) ([]byte, error) {
	text := string(content)
	if strings.HasPrefix(text, "\ufeff") {
		return nil, errors.New("starts with a byte-order mark, which hides its frontmatter from the API; save it as UTF-8 without a BOM")
	}
	start := 0
	for start < len(text) {
		line, _, more := strings.Cut(text[start:], "\n")
		// Only ASCII blanks make a line skippable. A line holding anything else
		// (a non-breaking space, NEL) is not blank to every reader.
		if strings.Trim(line, " \t\r\v\f") != "" {
			break
		}
		start += len(line)
		if more {
			start++
		}
	}
	doc := text[start:]
	if !strings.HasPrefix(doc, "---\n") && !strings.HasPrefix(doc, "---\r\n") {
		return nil, errors.New(`has no frontmatter: it must open with a line of exactly "---"`)
	}
	end := strings.Index(doc[3:], "\n---\n")
	if end < 0 {
		end = strings.Index(doc[3:], "\n---\r\n")
	}
	switch {
	case end < 0:
		return nil, errors.New(`frontmatter is never closed: it needs a line of exactly "---", with no trailing spaces and a newline after it`)
	case end <= 1:
		return nil, errors.New("frontmatter is empty; it must at least declare `name`")
	case end+3 > maxSkillFrontmatter:
		return nil, fmt.Errorf("frontmatter is longer than the %d bytes the API reads", maxSkillFrontmatter)
	}
	// end indexes doc[3:], so the closing fence starts at end+3 of doc and the
	// YAML is what lies between the opening fence's newline and it.
	yaml := doc[4 : end+3]
	return []byte(strings.NewReplacer("\r\n", "\n", "\r", "\n").Replace(yaml)), nil
}

// checkSkillName refuses a manifest name that is not one path segment. The
// name is the directory the skill is uploaded under, and path.Join would
// otherwise turn "../../other" into a multipart filename outside that folder.
func checkSkillName(name string) error {
	if name == "" {
		return nil
	}
	if strings.Contains(name, "/") || strings.Contains(name, `\`) || strings.Contains(name, "..") {
		return fmt.Errorf("name %q cannot contain '/', '\\', or '..'; it is the directory the skill is uploaded under", name)
	}
	return nil
}

// firstString returns the first of keys whose value is a non-blank string, trimmed.
func firstString(m map[string]any, keys ...string) string {
	for _, k := range keys {
		s, _ := m[k].(string)
		if s = strings.TrimSpace(s); s != "" {
			return s
		}
	}
	return ""
}

// skillBundle is a skill's core.Payload: the files, the folder they upload
// under, and the one piece of metadata not derived from SKILL.md itself.
type skillBundle struct {
	// UploadDir is the top-level folder the bundle uploads under. The API
	// rejects an upload whose folder does not match the `name` in SKILL.md, and
	// a skill vendored from another repository frequently lives in a directory
	// named something else — so the declared name wins over the directory.
	UploadDir string
	// DisplayName is settable only when the skill is first created.
	DisplayName string
	// Files is the bundle's content, sorted by RelPath.
	Files []skillFile
}

// Fingerprint hashes paths as well as content, so renaming a file is a change
// even when the bytes are unmoved. It uses the collected file list rather than
// re-walking, so what is hashed is provably what gets uploaded.
//
// The result lands in the lockfile, so what is hashed and the labels written
// below are part of its format: changing either makes every recorded skill
// plan as changed.
func (b *skillBundle) Fingerprint() (string, error) {
	h := sha256.New()
	fmt.Fprintf(h, "dir:%s\x00title:%s\x00", b.UploadDir, b.DisplayName)
	for _, f := range b.Files {
		fmt.Fprintf(h, "path:%s\x00size:%d\x00", f.RelPath, f.Size)
		fh, err := os.Open(f.AbsPath)
		if err != nil {
			return "", err
		}
		_, err = io.Copy(h, fh)
		fh.Close()
		if err != nil {
			return "", err
		}
		h.Write([]byte{0})
	}
	return hex.EncodeToString(h.Sum(nil))[:32], nil
}

// Describe summarizes the bundle for plan output, e.g. "3 files".
func (b *skillBundle) Describe() string {
	if len(b.Files) == 1 {
		return "1 file"
	}
	return fmt.Sprintf("%d files", len(b.Files))
}

// skillFile is one member of a skill bundle.
type skillFile struct {
	// RelPath is the slash-separated path within the skill directory.
	RelPath string
	// AbsPath is where to read it from.
	AbsPath string
	// Size is the file's length in bytes; it is hashed alongside RelPath.
	Size int64
}

// UploadName is the multipart filename the API sees. The server reconstructs
// the tree from these, and requires every file to sit under one top-level
// folder, so the bundle's upload folder (skillBundle.UploadDir) is prefixed.
func (f skillFile) UploadName(uploadDir string) string {
	return path.Join(uploadDir, f.RelPath)
}

// ContentType is the MIME type the file uploads as, guessed from its extension.
func (f skillFile) ContentType() string {
	if ct := mime.TypeByExtension(filepath.Ext(f.RelPath)); ct != "" {
		return ct
	}
	return "application/octet-stream"
}

// bundleSkipDirs are never included in a skill bundle: they are tooling, not
// content.
var bundleSkipDirs = map[string]bool{
	".git": true, "node_modules": true, "vendor": true,
	"__pycache__": true, ".venv": true, "venv": true, "dist": true, "build": true,
}

// Limits on a skill bundle. These are client-side sanity checks: the point is
// to fail with "you pointed at your home directory" rather than to stream a
// gigabyte into a multipart request and let the server decide.
const (
	maxSkillFiles     = 1000
	maxSkillTotalSize = 32 << 20 // 32 MiB
	maxSkillFileSize  = 16 << 20 // 16 MiB
)

// collectSkillFiles lists a skill directory's files, sorted by path. It skips
// dotfiles, bundleSkipDirs and anything that is not a regular file, enforces
// the maxSkill* limits, and fails when there is no SKILL.md at the root.
//
// The content hash and the upload are both driven from this one list, so
// drift detection can never disagree with what was sent.
func collectSkillFiles(dir string) ([]skillFile, error) {
	var files []skillFile
	var total int64

	err := filepath.WalkDir(dir, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		name := d.Name()
		if p != dir && (strings.HasPrefix(name, ".") || (d.IsDir() && bundleSkipDirs[name])) {
			// Dotfiles (.git, .DS_Store, .gitignore) and tooling directories
			// are not skill content; uploading them is noise at best.
			if d.IsDir() {
				return fs.SkipDir
			}
			return nil
		}
		if d.IsDir() {
			return nil
		}
		info, err := d.Info()
		if err != nil {
			return err
		}
		if !info.Mode().IsRegular() {
			// Symlinks would let a skill bundle reach outside its directory.
			return nil
		}
		if info.Size() > maxSkillFileSize {
			return fmt.Errorf("%s is %s, over the %d MiB per-file limit",
				p, humanSize(info.Size()), maxSkillFileSize>>20)
		}
		rel, err := filepath.Rel(dir, p)
		if err != nil {
			return err
		}
		total += info.Size()
		if total > maxSkillTotalSize {
			return fmt.Errorf("skill directory %s exceeds the %d MiB total limit", dir, maxSkillTotalSize>>20)
		}
		files = append(files, skillFile{
			RelPath: filepath.ToSlash(rel),
			AbsPath: p,
			Size:    info.Size(),
		})
		if len(files) > maxSkillFiles {
			return fmt.Errorf("skill directory %s holds more than %d files", dir, maxSkillFiles)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	sort.Slice(files, func(i, j int) bool { return files[i].RelPath < files[j].RelPath })

	// The API wants the manifest spelled SKILL.md exactly; a file system that
	// ignores case would otherwise let skill.md through to be rejected there.
	i := slices.IndexFunc(files, func(f skillFile) bool { return strings.EqualFold(f.RelPath, skillFileName) })
	switch {
	case i < 0:
		return nil, fmt.Errorf("%s has no SKILL.md at its root", dir)
	case files[i].RelPath != skillFileName:
		return nil, fmt.Errorf("%s: the manifest must be spelled %s, in capitals; found %s", dir, skillFileName, files[i].RelPath)
	}
	return files, nil
}

// humanSize formats a byte count for an error message, e.g. "17.0 MiB".
func humanSize(n int64) string {
	switch {
	case n >= 1<<20:
		return fmt.Sprintf("%.1f MiB", float64(n)/(1<<20))
	case n >= 1<<10:
		return fmt.Sprintf("%.1f KiB", float64(n)/(1<<10))
	default:
		return fmt.Sprintf("%d B", n)
	}
}
