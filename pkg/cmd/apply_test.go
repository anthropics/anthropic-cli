package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/anthropics/anthropic-cli/v2/internal/declarative/acctest"
	"github.com/anthropics/anthropic-cli/v2/internal/declarative/claude"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli/v3"
)

func TestCanPromptRefusesWithoutAHuman(t *testing.T) {
	// A script that forgot --yes is asking for a plan, not for its resources
	// to be rewritten unattended.
	require.NoError(t, canPrompt(true, false))

	for _, tc := range []struct {
		name             string
		stdoutIsTerminal bool
		stdinIsPiped     bool
	}{
		{"output is redirected", false, false},
		{"input is piped", true, true},
		{"neither", false, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			err := canPrompt(tc.stdoutIsTerminal, tc.stdinIsPiped)
			require.Error(t, err)
			assert.Contains(t, err.Error(), "--yes")
			assert.Contains(t, err.Error(), "--dry-run")
		})
	}
}

// TestAccApplyCommand drives the real binary against the real API through a
// cold start, a converged re-plan and a prune. The reconciler has its own
// acceptance suite; this covers what only the command owns — flag wiring,
// credential resolution from the active profile, lockfile placement, the
// --yes gate, exit codes, and the rendered plan.
//
// It runs a subprocess because the package-level Command tree keeps parsed
// flag state between in-process runs.
func TestAccApplyCommand(t *testing.T) {
	acctest.SkipUnlessEnabled(t)
	bin := buildAnt(t)
	tree := t.TempDir()
	lockPath := filepath.Join(tree, claude.Registry().LockfileName())
	acctest.WriteTree(t, tree, acctest.Substitute(map[string]string{
		"skills/greeter/SKILL.md": "---\nname: greeter\ndescription: says hello\ndisplay_title: {{name}}\n---\nSay hello.\n",
		"agents/concierge.md":     "---\nname: {{name}}\nmodel: claude-sonnet-4-5\nskills:\n  - ../skills/greeter\n---\nBe helpful.\n",
		"environments/lobby.yml":  "name: {{name}}\n",
	}, acctest.UniqueName(t)))

	ant := func(args ...string) (string, error) {
		cmd := exec.Command(bin, append([]string{"apply", "--lock-file", lockPath}, args...)...)
		cmd.Dir = tree
		// apply creates its remote-skill cache on every run; keep that out of
		// the developer's own cache directory.
		cmd.Env = append(os.Environ(), "XDG_CACHE_HOME="+t.TempDir(), "NO_COLOR=1")
		out, err := cmd.CombinedOutput()
		return string(out), err
	}
	t.Cleanup(func() {
		// With the files gone every tracked key is an orphan, so --prune tears
		// down whatever the test created, whichever step it stopped at.
		for _, dir := range []string{"skills", "agents", "environments"} {
			_ = os.RemoveAll(filepath.Join(tree, dir))
		}
		if out, err := ant("--yes", "--prune"); err != nil {
			t.Logf("cleanup: %v\n%s", err, out)
		}
	})

	out, err := ant("--dry-run", tree)
	require.NoError(t, err, out)
	assert.Contains(t, out, "No lockfile yet")
	assert.Contains(t, out, "+ ./agents/concierge.md")
	assert.Contains(t, out, "Plan: 3 to create, 0 to update, 0 to remove, 0 unchanged.")
	assert.NoFileExists(t, lockPath, "--dry-run must not write state")

	out, err = ant(tree)
	require.Error(t, err, "without a terminal and without --yes, apply must refuse")
	assert.Contains(t, out, "--yes")
	assert.NoFileExists(t, lockPath)

	out, err = ant("--yes", tree)
	require.NoError(t, err, out)
	assert.Contains(t, out, "Apply complete: 3 created, 0 updated, 0 removed, 0 unchanged.")
	assert.FileExists(t, lockPath)

	// No paths: everything tracked is reconciled, and it has converged.
	out, err = ant()
	require.NoError(t, err, out)
	assert.Contains(t, out, "Found lockfile at")
	assert.Contains(t, out, "Everything is up to date.")

	require.NoError(t, os.RemoveAll(filepath.Join(tree, "agents")))
	out, err = ant("--yes")
	require.NoError(t, err, out)
	assert.Contains(t, out, "--prune", "an orphan is reported with the flag that removes it")
	assert.Contains(t, out, "Everything is up to date.")

	out, err = ant("--yes", "--prune")
	require.NoError(t, err, out)
	assert.Contains(t, out, "- ./agents/concierge.md")
	assert.Contains(t, out, "0 created, 0 updated, 1 removed")
	data, err := os.ReadFile(lockPath)
	require.NoError(t, err)
	assert.NotContains(t, string(data), "concierge")
	assert.Contains(t, string(data), "greeter")
}

// buildAnt compiles the CLI once for the whole test. go run would relink on
// every invocation and add its own exit-status line to the output under test.
// Tests run in their package directory, so the module root is two levels up.
func buildAnt(t *testing.T) string {
	t.Helper()
	bin := filepath.Join(t.TempDir(), "ant")
	build := exec.Command("go", "build", "-o", bin, "./cmd/ant")
	build.Dir = "../.."
	out, err := build.CombinedOutput()
	require.NoError(t, err, "go build: %s", out)
	return bin
}

// Errors from apply name files and lockfile entries from the repository, so
// they print escaped, keep their line breaks and still unwrap.
func TestApplyErrorEscapesControlCharacters(t *testing.T) {
	cause := errors.New("cause")
	err := error(applyError{fmt.Errorf("./agents/a\x1b[2J.md: %w\n  try again", cause)})
	assert.Equal(t, "./agents/a\\x1b[2J.md: cause\n      try again", err.Error())
	assert.ErrorIs(t, err, cause)

	forged := applyError{errors.New("./agents/a\nError: forged.md: bad")}
	assert.Equal(t, "./agents/a\n    Error: forged.md: bad", forged.Error())
}

// handleApply returns its errors escaped, whichever step they come from.
func TestHandleApplyEscapesItsError(t *testing.T) {
	path := filepath.Join(t.TempDir(), "a\x1b[2J.json")
	require.NoError(t, os.WriteFile(path, []byte("{"), 0o600))
	cmd := &cli.Command{Flags: []cli.Flag{&cli.StringFlag{Name: "lock-file"}}, Action: handleApply}
	err := cmd.Run(context.Background(), []string{"apply", "--lock-file", path})
	require.Error(t, err)
	assert.Contains(t, err.Error(), `a\x1b[2J.json`)
	assert.NotContains(t, err.Error(), "\x1b")
}

func setUserCacheDir(t *testing.T) string {
	t.Helper()
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("XDG_CACHE_HOME", filepath.Join(home, "cache"))
	t.Setenv("LocalAppData", home)
	base, err := os.UserCacheDir()
	require.NoError(t, err)
	return base
}

func unsetUserCacheDir(t *testing.T) {
	t.Helper()
	for _, name := range []string{"HOME", "XDG_CACHE_HOME", "LocalAppData"} {
		t.Setenv(name, "")
	}
	if _, err := os.UserCacheDir(); err == nil {
		t.Skip("this platform finds a user cache directory without HOME")
	}
}

func setTempDir(t *testing.T) string {
	t.Helper()
	tmp := t.TempDir()
	for _, name := range []string{"TMPDIR", "TMP", "TEMP"} {
		t.Setenv(name, tmp)
	}
	return tmp
}

// Without a user cache directory, another user on the machine must not be able
// to predict the cache's path, create it first, or read what is fetched into it.
func TestRemoteCacheDirWithoutAUserCacheDirIsPrivateToThisRun(t *testing.T) {
	unsetUserCacheDir(t)
	tmp := setTempDir(t)

	first, removeFirst, err := remoteCacheDir()
	require.NoError(t, err)
	second, removeSecond, err := remoteCacheDir()
	require.NoError(t, err)
	defer removeSecond()

	assert.Equal(t, tmp, filepath.Dir(first))
	assert.Equal(t, tmp, filepath.Dir(second))
	assert.NotEqual(t, first, second)
	assert.Equal(t, byte('.'), filepath.Base(first)[0], "a directory walk skips it")
	if runtime.GOOS != "windows" {
		for _, dir := range []string{first, second} {
			st, err := os.Stat(dir)
			require.NoError(t, err)
			assert.Equal(t, os.FileMode(0o700), st.Mode().Perm(), dir)
		}
	}

	removeFirst()
	assert.NoDirExists(t, first)
}

func TestRemoteCacheDirKeepsTheUserCacheBetweenRuns(t *testing.T) {
	base := setUserCacheDir(t)

	dir, remove, err := remoteCacheDir()
	require.NoError(t, err)
	remove()

	assert.Equal(t, filepath.Join(base, "anthropic", "apply-remote"), dir)
	assert.DirExists(t, dir)
}

func TestRemoteCacheDirInTheUserCacheIsPrivate(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("no permission bits")
	}
	base := setUserCacheDir(t)
	mode := func(dir string) os.FileMode {
		st, err := os.Stat(dir)
		require.NoError(t, err)
		return st.Mode().Perm()
	}

	dir, _, err := remoteCacheDir()
	require.NoError(t, err)
	require.Equal(t, filepath.Join(base, "anthropic", "apply-remote"), dir)
	assert.Equal(t, os.FileMode(0o700), mode(dir))

	require.NoError(t, os.Chmod(dir, 0o755))
	_, _, err = remoteCacheDir()
	require.NoError(t, err)
	assert.Equal(t, os.FileMode(0o700), mode(dir), "a cache an earlier version created is tightened")
}

func TestRemoteCacheDirIsPrivateToThisRunWhenTheUserCacheIsUnusable(t *testing.T) {
	base := setUserCacheDir(t)
	tmp := setTempDir(t)
	require.NoError(t, os.MkdirAll(base, 0o700))
	require.NoError(t, os.WriteFile(filepath.Join(base, "anthropic"), nil, 0o600))

	dir, remove, err := remoteCacheDir()
	require.NoError(t, err)
	assert.Equal(t, tmp, filepath.Dir(dir))

	remove()
	assert.NoDirExists(t, dir)
}
