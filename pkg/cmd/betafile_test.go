// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cmd

import (
	"strings"
	"testing"

	"github.com/anthropics/anthropic-cli/internal/mocktest"
)

func TestBetaFilesList(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:files", "list",
			"--max-items", "10",
			"--after-id", "after_id",
			"--before-id", "before_id",
			"--limit", "1",
			"--beta", "message-batches-2024-09-24",
		)
	})
}

func TestBetaFilesDelete(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:files", "delete",
			"--file-id", "file_id",
			"--beta", "message-batches-2024-09-24",
		)
	})
}

func TestBetaFilesDownload(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:files", "download",
			"--file-id", "file_id",
			"--beta", "message-batches-2024-09-24",
			"--output", "/dev/null",
		)
	})
}

func TestBetaFilesRetrieveMetadata(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:files", "retrieve-metadata",
			"--file-id", "file_id",
			"--beta", "message-batches-2024-09-24",
		)
	})
}

func TestBetaFilesUpload(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:files", "upload",
			"--file", mocktest.TestFile(t, "Example data"),
			"--beta", "message-batches-2024-09-24",
		)
	})

	t.Run("piping data", func(t *testing.T) {
		testFile := mocktest.TestFile(t, "Example data")
		// Test piping YAML data over stdin
		pipeDataStr := "file: Example data"
		pipeDataStr = strings.ReplaceAll(pipeDataStr, "Example data", testFile)
		pipeData := []byte(pipeDataStr)
		mocktest.TestRunMockTestWithPipeAndFlags(
			t, pipeData,
			"--api-key", "string",
			"beta:files", "upload",
			"--beta", "message-batches-2024-09-24",
		)
	})
}
