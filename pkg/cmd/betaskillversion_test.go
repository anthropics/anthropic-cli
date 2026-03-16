// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cmd

import (
	"testing"

	"github.com/anthropics/anthropic-cli/internal/mocktest"
)

func TestBetaSkillsVersionsCreate(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		t.Skip("CLI multipart serialization does not handle complex array elements (e.g. --file [null])")
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:skills:versions", "create",
			"--skill-id", "skill_id",
			"--file", "[Example data]",
			"--beta", "message-batches-2024-09-24",
		)
	})

	t.Run("piping data", func(t *testing.T) {
		// Test piping YAML data over stdin
		pipeData := []byte("" +
			"files:\n" +
			"  - Example data\n")
		mocktest.TestRunMockTestWithPipeAndFlags(
			t, pipeData,
			"--api-key", "string",
			"beta:skills:versions", "create",
			"--skill-id", "skill_id",
			"--beta", "message-batches-2024-09-24",
		)
	})
}

func TestBetaSkillsVersionsRetrieve(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:skills:versions", "retrieve",
			"--skill-id", "skill_id",
			"--version", "version",
			"--beta", "message-batches-2024-09-24",
		)
	})
}

func TestBetaSkillsVersionsList(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:skills:versions", "list",
			"--max-items", "10",
			"--skill-id", "skill_id",
			"--limit", "0",
			"--page", "page",
			"--beta", "message-batches-2024-09-24",
		)
	})
}

func TestBetaSkillsVersionsDelete(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:skills:versions", "delete",
			"--skill-id", "skill_id",
			"--version", "version",
			"--beta", "message-batches-2024-09-24",
		)
	})
}
