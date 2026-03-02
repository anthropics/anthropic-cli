package cmd

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStreamOutput(t *testing.T) {
	t.Setenv("PAGER", "cat")
	err := streamOutput("stream test", func(w *os.File) error {
		_, writeErr := w.WriteString("Hello world\n")
		return writeErr
	})
	if err != nil {
		t.Errorf("streamOutput failed: %v", err)
	}
}

func TestWriteBinaryResponse(t *testing.T) {
	t.Run("write to explicit file", func(t *testing.T) {
		tmpDir := t.TempDir()
		outfile := tmpDir + "/output.txt"
		body := []byte("test content")
		resp := &http.Response{
			Body: io.NopCloser(bytes.NewReader(body)),
		}

		msg, err := writeBinaryResponse(resp, outfile)

		require.NoError(t, err)
		assert.Contains(t, msg, outfile)

		content, err := os.ReadFile(outfile)
		require.NoError(t, err)
		assert.Equal(t, body, content)
	})

	t.Run("write to stdout", func(t *testing.T) {
		oldStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		body := []byte("stdout content")
		resp := &http.Response{
			Body: io.NopCloser(bytes.NewReader(body)),
		}
		msg, err := writeBinaryResponse(resp, "-")

		w.Close()
		os.Stdout = oldStdout

		require.NoError(t, err)
		assert.Empty(t, msg)

		var buf bytes.Buffer
		_, _ = buf.ReadFrom(r)
		assert.Equal(t, body, buf.Bytes())
	})
}

func TestCreateDownloadFile(t *testing.T) {
	t.Run("creates file with filename from header", func(t *testing.T) {
		tmpDir := t.TempDir()
		oldWd, _ := os.Getwd()
		os.Chdir(tmpDir)
		defer os.Chdir(oldWd)

		resp := &http.Response{
			Header: http.Header{
				"Content-Disposition": []string{`attachment; filename="test.txt"`},
			},
		}
		file, err := createDownloadFile(resp, []byte("test content"))
		require.NoError(t, err)
		defer file.Close()
		assert.Equal(t, "test.txt", filepath.Base(file.Name()))

		// Create a second file with the same name to ensure it doesn't clobber the first
		resp2 := &http.Response{
			Header: http.Header{
				"Content-Disposition": []string{`attachment; filename="test.txt"`},
			},
		}
		file2, err := createDownloadFile(resp2, []byte("second content"))
		require.NoError(t, err)
		defer file2.Close()
		assert.NotEqual(t, file.Name(), file2.Name(), "second file should have a different name")
		assert.Contains(t, filepath.Base(file2.Name()), "test")
	})

	t.Run("creates temp file when no header", func(t *testing.T) {
		tmpDir := t.TempDir()
		oldWd, _ := os.Getwd()
		os.Chdir(tmpDir)
		defer os.Chdir(oldWd)

		resp := &http.Response{Header: http.Header{}}
		file, err := createDownloadFile(resp, []byte("test content"))
		require.NoError(t, err)
		defer file.Close()
		assert.Contains(t, filepath.Base(file.Name()), "file-")
	})

	t.Run("prevents directory traversal", func(t *testing.T) {
		tmpDir := t.TempDir()
		oldWd, _ := os.Getwd()
		os.Chdir(tmpDir)
		defer os.Chdir(oldWd)

		resp := &http.Response{
			Header: http.Header{
				"Content-Disposition": []string{`attachment; filename="../../../etc/passwd"`},
			},
		}
		file, err := createDownloadFile(resp, []byte("test content"))
		require.NoError(t, err)
		defer file.Close()
		assert.Equal(t, "passwd", filepath.Base(file.Name()))
	})
}
