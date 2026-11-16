package web

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/anthropics/anthropic-cli/v2/internal/sessions/live"
)

// startServe runs Serve until the test ends, with the temp directory pointed
// at tempDir, and returns what it logs.
func startServe(t *testing.T, tempDir string, open func(path string) error) <-chan string {
	t.Helper()
	if runtime.GOOS == "windows" {
		t.Skip("TMPDIR does not move the temp directory, and file modes say nothing about access")
	}
	t.Setenv("TMPDIR", tempDir)
	ctx, cancel := context.WithCancel(context.Background())
	logged, done := make(chan string, 4), make(chan error, 1)
	go func() {
		done <- Serve(ctx, &fakeSession{updates: make(chan live.ThreadUpdate)}, Options{
			Open: open,
			Logf: func(format string, args ...any) { logged <- fmt.Sprintf(format, args...) },
		})
	}()
	t.Cleanup(func() {
		cancel()
		assert.NoError(t, <-done)
	})
	return logged
}

func viewerURL(t *testing.T, logged <-chan string) string {
	t.Helper()
	url, found := strings.CutPrefix(<-logged, "Session viewer: ")
	require.True(t, found)
	return url
}

func get(t *testing.T, url string) string {
	t.Helper()
	resp, err := http.Get(url)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	return readBody(t, resp)
}

func TestServeOpensAnOwnerOnlyFileThatLeadsToTheViewer(t *testing.T) {
	opened := make(chan string, 1)

	logged := startServe(t, t.TempDir(), func(path string) error {
		opened <- path
		return nil
	})
	url, path := viewerURL(t, logged), <-opened

	_, bootCode, _ := strings.Cut(url, "boot=")
	require.NotEmpty(t, bootCode)
	assert.NotContains(t, path, bootCode, "other users can read the opener's arguments")
	info, err := os.Stat(path)
	require.NoError(t, err)
	assert.Equal(t, os.FileMode(0o600), info.Mode())
	page, err := os.ReadFile(path)
	require.NoError(t, err)
	assert.Contains(t, string(page), `content="0;url=`+url+`"`)
}

func TestServeRemovesTheFileOnceItsBootCodeIsSpent(t *testing.T) {
	opened := make(chan string, 1)
	logged := startServe(t, t.TempDir(), func(path string) error {
		opened <- path
		return nil
	})
	url, path := viewerURL(t, logged), <-opened

	assert.NotContains(t, get(t, url+"wrong"), `"token":"`)
	assert.FileExists(t, path, "the browser may not have read it yet")

	assert.Contains(t, get(t, url), `"token":"`)
	assert.NoFileExists(t, path)
}

func TestServeWithNoBrowserToOpenWritesNoFile(t *testing.T) {
	tempDir := t.TempDir()
	url := viewerURL(t, startServe(t, tempDir, nil))

	assert.Contains(t, get(t, url), `"token":"`)
	files, err := os.ReadDir(tempDir)
	require.NoError(t, err)
	assert.Empty(t, files)
}

func TestServeRemovesTheFileWhenNothingOpensIt(t *testing.T) {
	tempDir := t.TempDir()
	logged := startServe(t, tempDir, func(string) error { return errors.New("no opener") })
	url := viewerURL(t, logged)

	assert.Contains(t, <-logged, "could not open a browser (no opener)")
	files, err := os.ReadDir(tempDir)
	require.NoError(t, err)
	assert.Empty(t, files)
	assert.Contains(t, get(t, url), `"token":"`)
}

func TestServeOpensNothingWhenItCannotWriteTheFile(t *testing.T) {
	logged := startServe(t, filepath.Join(t.TempDir(), "missing"), func(path string) error {
		t.Errorf("opened %s", path)
		return errors.New("opened")
	})
	url := viewerURL(t, logged)

	assert.Contains(t, <-logged, "could not open a browser")
	assert.Contains(t, get(t, url), `"token":"`)
}
