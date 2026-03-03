// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cmd

import (
	"testing"

	"github.com/anthropics/anthropic-cli/internal/mocktest"
)

func TestBetaFilesList(t *testing.T) {
	mocktest.TestRunMockTestWithFlags(
		t,
		"beta:files", "list",
		"--api-key", "string",
		"--after-id", "after_id",
		"--before-id", "before_id",
		"--limit", "1",
		"--beta", "message-batches-2024-09-24",
	)
}

func TestBetaFilesDelete(t *testing.T) {
	mocktest.TestRunMockTestWithFlags(
		t,
		"beta:files", "delete",
		"--api-key", "string",
		"--file-id", "file_id",
		"--beta", "message-batches-2024-09-24",
	)
}

func TestBetaFilesDownload(t *testing.T) {
	mocktest.TestRunMockTestWithFlags(
		t,
		"beta:files", "download",
		"--api-key", "string",
		"--file-id", "file_id",
		"--beta", "message-batches-2024-09-24",
		"--output", "/dev/null",
	)
}

func TestBetaFilesRetrieveMetadata(t *testing.T) {
	mocktest.TestRunMockTestWithFlags(
		t,
		"beta:files", "retrieve-metadata",
		"--api-key", "string",
		"--file-id", "file_id",
		"--beta", "message-batches-2024-09-24",
	)
}

func TestBetaFilesUpload(t *testing.T) {
	mocktest.TestRunMockTestWithFlags(
		t,
		"beta:files", "upload",
		"--api-key", "string",
		"--file", "...",
		"--beta", "message-batches-2024-09-24",
	)
}
