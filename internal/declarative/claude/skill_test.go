package claude

import (
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCollectSkillFilesIsDeterministicAndSkipsJunk(t *testing.T) {
	root := writeTree(t, map[string]string{
		"s/SKILL.md":            "---\nname: s\n---\nbody\n",
		"s/refs/b.md":           "b\n",
		"s/refs/a.md":           "a\n",
		"s/.hidden":             "nope\n",
		"s/.git/config":         "nope\n",
		"s/node_modules/x/y.js": "nope\n",
	})
	files, err := collectSkillFiles(filepath.Join(root, "s"))
	require.NoError(t, err)

	var names []string
	for _, f := range files {
		names = append(names, f.RelPath)
	}
	assert.Equal(t, []string{"SKILL.md", "refs/a.md", "refs/b.md"}, names)
}

func TestCollectSkillFilesRequiresSkillMD(t *testing.T) {
	root := writeTree(t, map[string]string{"s/notes.md": "hi\n"})
	_, err := collectSkillFiles(filepath.Join(root, "s"))
	require.ErrorContains(t, err, "no SKILL.md")

	// A file system that ignores case finds skill.md for SKILL.md; the API
	// does not.
	root = writeTree(t, map[string]string{"s/skill.md": "---\nname: s\n---\n"})
	_, err = collectSkillFiles(filepath.Join(root, "s"))
	require.ErrorContains(t, err, "must be spelled SKILL.md")
}

// The fences are found the way the API finds them, so a manifest the API would
// read as having no frontmatter is refused here too, before an upload.
func TestSkillFrontmatterIsFoundTheWayTheAPIFindsIt(t *testing.T) {
	cases := map[string]struct {
		md   string
		yaml string // empty means refused
	}{
		"plain":               {"---\nname: a\n---\nbody", "name: a"},
		"crlf fences":         {"---\r\nname: a\r\n---\r\nbody", "\nname: a\n"},
		"leading blank lines": {"\n \t\n---\nname: a\n---\n", "name: a"},
		"body rule later":     {"---\nname: a\n---\nbody\n---\nmore", "name: a"},
		"empty body":          {"---\nname: a\n---\n", "name: a"},

		"bom":               {"\ufeff---\nname: a\n---\n", ""},
		"cr opener":         {"---\rname: a\n---\nbody", ""},
		"leading nbsp line": {"\u00a0\n---\nname: a\n---\n", ""},
		"leading nel line":  {"\u0085\n---\nname: a\n---\n", ""},
		"no final newline":  {"---\nname: a\n---", ""},
		"trailing spaces":   {"---\nname: a\n---  \nbody", ""},
		"four dashes":       {"----\nname: a\n---\n", ""},
		"cr only":           {"---\rname: a\r---\r", ""},
		"unterminated":      {"---\nname: a\nbody", ""},
		"not first":         {"intro\n---\nname: a\n---\n", ""},
		"adjacent fences":   {"---\n---\nbody", ""},
		"blank frontmatter": {"---\n\n---\nbody", ""},
		"no frontmatter":    {"# Just markdown\nname: nope\n", ""},
	}
	for label, c := range cases {
		got, err := skillFrontmatter([]byte(c.md))
		if c.yaml == "" {
			assert.Error(t, err, label)
			continue
		}
		if assert.NoError(t, err, label) {
			assert.Equal(t, c.yaml, string(got), label)
		}
	}
}

func TestSkillManifestProblemsAreNamedBeforeAnyUpload(t *testing.T) {
	for content, want := range map[string]string{
		"\ufeff---\nname: s\n---\nbody\n": "byte-order mark",
		"# no frontmatter\n":              "has no frontmatter",
		"---\nname: s\n---   \nbody\n":    "never closed",
		"---\n---\nbody\n":                "empty",
	} {
		root := writeTree(t, map[string]string{"s/SKILL.md": content})
		_, _, err := loadSkillDir(filepath.Join(root, "s"))
		require.ErrorContains(t, err, want)
		assert.ErrorContains(t, err, "SKILL.md", "the message names the file")
	}
}

func TestSkillNameCannotBeAPath(t *testing.T) {
	for _, name := range []string{"../../other", "/other", `..\other`, "foo/bar", "..", "foo..bar"} {
		root := writeTree(t, map[string]string{
			"s/SKILL.md": "---\nname: " + strconv.Quote(name) + "\n---\nbody\n",
		})
		_, _, err := loadSkillDir(filepath.Join(root, "s"))
		require.ErrorContains(t, err, "cannot contain")
		assert.ErrorContains(t, err, "SKILL.md")
	}

	root := writeTree(t, map[string]string{
		"s/SKILL.md": "---\nname: pr-writer\n---\nbody\n",
	})
	_, payload, err := loadSkillDir(filepath.Join(root, "s"))
	require.NoError(t, err)
	bundle := payload.(*skillBundle)
	assert.Equal(t, "pr-writer", bundle.UploadDir)
	assert.Equal(t, "pr-writer/SKILL.md", bundle.Files[0].UploadName(bundle.UploadDir))
}

func TestSkillBodyCarriesTheDeclaredName(t *testing.T) {
	root := writeTree(t, map[string]string{
		"named/SKILL.md":   "---\nname: pdf-tools\ndescription: d\n---\n",
		"unnamed/SKILL.md": "---\ndescription: d\n---\n",
	})
	body, _, err := loadSkillDir(filepath.Join(root, "named"))
	require.NoError(t, err)
	assert.Equal(t, map[string]any{"name": "pdf-tools", "display_name": "pdf-tools"}, body)

	// No name declared: nothing to hold the server to, and the directory
	// stands in for display and upload.
	body, payload, err := loadSkillDir(filepath.Join(root, "unnamed"))
	require.NoError(t, err)
	assert.Equal(t, map[string]any{"display_name": "unnamed"}, body)
	assert.Equal(t, "unnamed", payload.(*skillBundle).UploadDir)
}

func TestUploadNamePrefersTheDeclaredSkillName(t *testing.T) {
	// The API rejects an upload whose folder does not match SKILL.md's name,
	// which bites whenever a skill is vendored under a different directory.
	root := writeTree(t, map[string]string{
		"composition-patterns/SKILL.md": "---\nname: vercel-composition-patterns\n---\nbody\n",
		"pr-writer/SKILL.md":            "---\nname: pr-writer\n---\nbody\n",
	})
	for dir, want := range map[string]string{
		"composition-patterns": "vercel-composition-patterns",
		"pr-writer":            "pr-writer",
	} {
		_, payload, err := loadSkillDir(filepath.Join(root, dir))
		require.NoError(t, err)
		assert.Equal(t, want, payload.(*skillBundle).UploadDir, dir)
	}
}

func TestSkillHashTracksContentAndPaths(t *testing.T) {
	root := writeTree(t, map[string]string{
		"s/SKILL.md":  "---\nname: s\n---\nbody\n",
		"s/refs/a.md": "a\n",
	})
	dir := filepath.Join(root, "s")

	hashOf := func(title string) string {
		files, err := collectSkillFiles(dir)
		require.NoError(t, err)
		h, err := (&skillBundle{UploadDir: "s", DisplayName: title, Files: files}).Fingerprint()
		require.NoError(t, err)
		return h
	}

	base := hashOf("Title")
	assert.Equal(t, base, hashOf("Title"))
	assert.NotEqual(t, base, hashOf("Other"))

	// Renaming a file changes the bundle even though the bytes are unchanged.
	require.NoError(t, os.Rename(filepath.Join(dir, "refs/a.md"), filepath.Join(dir, "refs/b.md")))
	assert.NotEqual(t, base, hashOf("Title"))
}
