// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cmd

import (
	"testing"

	"github.com/anthropics/anthropic-cli/internal/mocktest"
)

func TestBetaSkillsVersionsCreate(t *testing.T) {
	t.Skip("prism binary unsupported")
	mocktest.TestRunMockTestWithFlags(
		t,
		"beta:skills:versions", "create",
		"--api-key", "string",
		"--skill-id", "skill_id",
		"--file", "[null]",
		"--beta", "message-batches-2024-09-24",
	)
}

func TestBetaSkillsVersionsRetrieve(t *testing.T) {
	mocktest.TestRunMockTestWithFlags(
		t,
		"beta:skills:versions", "retrieve",
		"--api-key", "string",
		"--skill-id", "skill_id",
		"--version", "version",
		"--beta", "message-batches-2024-09-24",
	)
}

func TestBetaSkillsVersionsList(t *testing.T) {
	mocktest.TestRunMockTestWithFlags(
		t,
		"beta:skills:versions", "list",
		"--api-key", "string",
		"--skill-id", "skill_id",
		"--limit", "0",
		"--page", "page",
		"--beta", "message-batches-2024-09-24",
	)
}

func TestBetaSkillsVersionsDelete(t *testing.T) {
	mocktest.TestRunMockTestWithFlags(
		t,
		"beta:skills:versions", "delete",
		"--api-key", "string",
		"--skill-id", "skill_id",
		"--version", "version",
		"--beta", "message-batches-2024-09-24",
	)
}
