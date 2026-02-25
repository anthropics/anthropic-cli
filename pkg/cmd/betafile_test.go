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
		"--file-id", "file_id",
		"--beta", "message-batches-2024-09-24",
	)
}

func TestBetaFilesRetrieveMetadata(t *testing.T) {
	mocktest.TestRunMockTestWithFlags(
		t,
		"beta:files", "retrieve-metadata",
		"--file-id", "file_id",
		"--beta", "message-batches-2024-09-24",
	)
}

func TestBetaFilesUpload(t *testing.T) {
	mocktest.TestRunMockTestWithFlags(
		t,
		"beta:files", "upload",
		"--file", "",
		"--beta", "message-batches-2024-09-24",
	)
}
