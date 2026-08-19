// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cmd

import (
	"strings"
	"testing"

	"github.com/anthropics/anthropic-cli/internal/mocktest"
)

func TestFilesList(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"files", "list",
			"--max-items", "10",
			"--id", "[string]",
			"--limit", "1",
			"--page", "page",
		)
	})
}

func TestFilesDelete(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"files", "delete",
			"--file-id", "file_id",
		)
	})
}

func TestFilesDownload(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"files", "download",
			"--file-id", "file_id",
			"--output", "/dev/null",
		)
	})
}

func TestFilesRetrieveMetadata(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"files", "retrieve-metadata",
			"--file-id", "file_id",
		)
	})
}

func TestFilesUpload(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"files", "upload",
			"--file", mocktest.TestFile(t, "Example data"),
			"--expires-in-seconds", "3600",
		)
	})

	t.Run("piping data", func(t *testing.T) {
		testFile := mocktest.TestFile(t, "Example data")
		// Test piping YAML data over stdin
		pipeDataStr := "" +
			"file: Example data\n" +
			"expires_in_seconds: 3600\n"
		pipeDataStr = strings.ReplaceAll(pipeDataStr, "Example data", testFile)
		pipeData := []byte(pipeDataStr)
		mocktest.TestRunMockTestWithPipeAndFlags(
			t, pipeData,
			"--api-key", "string",
			"files", "upload",
		)
	})
}
