package claude

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSkillFrontmatterUsesEarliestClosingFence(t *testing.T) {
	for _, opener := range []string{"\n", "\r\n"} {
		for _, closer := range []string{"\n", "\r\n"} {
			t.Run(strings.ReplaceAll(opener+closer, "\n", "LF"), func(t *testing.T) {
				markdown := "---" + opener + "name: selected" + closer + "---" + closer +
					"Ordinary body text\n---\nname: wrong\n---\r\n"
				before := []byte(markdown)
				front, err := skillFrontmatter(before)
				require.NoError(t, err)
				assert.Equal(t, "name: selected", strings.TrimSpace(string(front)))
				assert.Equal(t, markdown, string(before))
			})
		}
	}
}

func TestSkillFrontmatterLimitExcludesBodyAfterCRLFFence(t *testing.T) {
	markdown := "---\nname: selected\n---\r\n" + strings.Repeat("x", maxSkillFrontmatter+10) + "\n---\n"
	front, err := skillFrontmatter([]byte(markdown))
	require.NoError(t, err)
	assert.Equal(t, "name: selected", strings.TrimSpace(string(front)))
}

func TestSkillDirectoryKeepsBodyWithMixedLineEndings(t *testing.T) {
	markdown := "---\r\nname: selected\r\ndisplay_name: Display\r\n---\r\n" +
		strings.Repeat("Use this skill.\n", maxSkillFrontmatter/10) + "\n---\nThis is a body separator, not frontmatter.\n"
	root := writeTree(t, map[string]string{"skill/SKILL.md": markdown, "skill/reference.txt": "reference"})
	body, payload, err := loadSkillDir(filepath.Join(root, "skill"))
	require.NoError(t, err)
	assert.Equal(t, map[string]any{"name": "selected", "display_name": "Display"}, body)
	bundle, ok := payload.(*skillBundle)
	require.True(t, ok)
	assert.Equal(t, "selected", bundle.UploadDir)
	require.Len(t, bundle.Files, 2)
	data, err := os.ReadFile(bundle.Files[0].AbsPath)
	require.NoError(t, err)
	assert.Equal(t, markdown, string(data))
	fingerprint, err := bundle.Fingerprint()
	require.NoError(t, err)
	assert.NotEmpty(t, fingerprint)
}
