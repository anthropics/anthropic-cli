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

	t.Run("rejects reserved Windows device names", func(t *testing.T) {
		tmpDir := t.TempDir()
		oldWd, _ := os.Getwd()
		os.Chdir(tmpDir)
		defer os.Chdir(oldWd)

		reservedNames := []string{"CON", "PRN", "AUX", "NUL", "COM1", "LPT1"}
		for _, name := range reservedNames {
			resp := &http.Response{
				Header: http.Header{
					"Content-Disposition": []string{`attachment; filename="` + name + `"`},
				},
			}
			file, err := createDownloadFile(resp, []byte("test content"))
			require.NoError(t, err, "should not error for reserved name %s", name)
			defer file.Close()
			// Should fall back to default filename, not the reserved name
			assert.NotEqual(t, name, filepath.Base(file.Name()), "reserved name %s should be rejected", name)
			assert.Contains(t, filepath.Base(file.Name()), "file-", "should use temp file naming for reserved name %s", name)
		}
	})

	t.Run("rejects reserved names with extensions", func(t *testing.T) {
		tmpDir := t.TempDir()
		oldWd, _ := os.Getwd()
		os.Chdir(tmpDir)
		defer os.Chdir(oldWd)

		// CON.txt is still reserved on Windows
		resp := &http.Response{
			Header: http.Header{
				"Content-Disposition": []string{`attachment; filename="CON.txt"`},
			},
		}
		file, err := createDownloadFile(resp, []byte("test content"))
		require.NoError(t, err)
		defer file.Close()
		assert.NotEqual(t, "CON.txt", filepath.Base(file.Name()))
		assert.Contains(t, filepath.Base(file.Name()), "file-")
	})

	t.Run("rejects path traversal in filename", func(t *testing.T) {
		tmpDir := t.TempDir()
		oldWd, _ := os.Getwd()
		os.Chdir(tmpDir)
		defer os.Chdir(oldWd)

		// Test path traversal patterns - filepath.Base should extract safe portion
		// and validation should reject the original patterns before Base extracts them
		resp := &http.Response{
			Header: http.Header{
				"Content-Disposition": []string{`attachment; filename="../../etc/passwd"`},
			},
		}
		file, err := createDownloadFile(resp, []byte("test content"))
		require.NoError(t, err)
		defer file.Close()
		// filepath.Base extracts "passwd" which is valid
		assert.Equal(t, "passwd", filepath.Base(file.Name()))

		// Test that embedded .. patterns are rejected by validation (filepath.Base doesn't strip these)
		resp2 := &http.Response{
			Header: http.Header{
				"Content-Disposition": []string{`attachment; filename="file..txt"`},
			},
		}
		file2, err := createDownloadFile(resp2, []byte("test content"))
		require.NoError(t, err)
		defer file2.Close()
		// "file..txt" contains ".." and should fall back to temp file
		assert.NotEqual(t, "file..txt", filepath.Base(file2.Name()), "filename with embedded .. should be rejected")
		assert.Contains(t, filepath.Base(file2.Name()), "file-", "should use temp file naming for rejected pattern")
	})
}

func TestIsValidFilename(t *testing.T) {
	t.Run("accepts valid filenames", func(t *testing.T) {
		validNames := []string{
			"file.txt",
			"document.pdf",
			"data.json",
			"archive.tar.gz",
			"README",
			".gitignore",
		}
		for _, name := range validNames {
			err := isValidFilename(name)
			assert.NoError(t, err, "filename %s should be valid", name)
		}
	})

	t.Run("rejects empty filename", func(t *testing.T) {
		err := isValidFilename("")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot be empty")
	})

	t.Run("rejects reserved Windows device names", func(t *testing.T) {
		reservedNames := []string{
			"CON", "PRN", "AUX", "NUL",
			"COM1", "COM2", "COM9",
			"LPT1", "LPT2", "LPT9",
		}
		for _, name := range reservedNames {
			err := isValidFilename(name)
			assert.Error(t, err, "reserved name %s should be rejected", name)
			assert.Contains(t, err.Error(), "reserved device name")
		}
	})

	t.Run("rejects reserved names case-insensitively", func(t *testing.T) {
		err := isValidFilename("con")
		assert.Error(t, err)
		err = isValidFilename("Con")
		assert.Error(t, err)
		err = isValidFilename("CON")
		assert.Error(t, err)
	})

	t.Run("rejects reserved names with extensions", func(t *testing.T) {
		err := isValidFilename("CON.txt")
		assert.Error(t, err)
		err = isValidFilename("COM1.pdf")
		assert.Error(t, err)
	})

	t.Run("rejects path traversal patterns", func(t *testing.T) {
		traversalNames := []string{
			"../file.txt",
			"file../name.txt",
			"file..txt",
			"dir/file.txt",
			"dir\\file.txt",
		}
		for _, name := range traversalNames {
			err := isValidFilename(name)
			assert.Error(t, err, "filename %s with path traversal should be rejected", name)
			assert.Contains(t, err.Error(), "path traversal")
		}
	})
}
