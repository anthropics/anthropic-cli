// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cmd

import (
	"testing"

	"github.com/anthropics/anthropic-cli/internal/mocktest"
)

func TestBetaSkillsCreate(t *testing.T) {
	t.Skip("CLI multipart serialization does not handle complex array elements (e.g. --file [null])")
	mocktest.TestRunMockTestWithFlags(
		t,
		"beta:skills", "create",
		"--api-key", "string",
		"--display-title", "display_title",
		"--file", "[null]",
		"--beta", "message-batches-2024-09-24",
	)
}

func TestBetaSkillsRetrieve(t *testing.T) {
	mocktest.TestRunMockTestWithFlags(
		t,
		"beta:skills", "retrieve",
		"--api-key", "string",
		"--skill-id", "skill_id",
		"--beta", "message-batches-2024-09-24",
	)
}

func TestBetaSkillsList(t *testing.T) {
	mocktest.TestRunMockTestWithFlags(
		t,
		"beta:skills", "list",
		"--api-key", "string",
		"--limit", "0",
		"--page", "page",
		"--source", "source",
		"--beta", "message-batches-2024-09-24",
	)
}

func TestBetaSkillsDelete(t *testing.T) {
	mocktest.TestRunMockTestWithFlags(
		t,
		"beta:skills", "delete",
		"--api-key", "string",
		"--skill-id", "skill_id",
		"--beta", "message-batches-2024-09-24",
	)
}
