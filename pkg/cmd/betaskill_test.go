// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cmd

import (
	"strings"
	"testing"

	"github.com/anthropics/anthropic-cli/internal/mocktest"
)

func TestBetaSkillsCreate(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		t.Skip("CLI multipart serialization does not handle complex array elements (e.g. --file [null])")
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:skills", "create",
			"--display-title", "display_title",
			"--file", mocktest.TestFile(t, "[Example data]"),
			"--beta", "message-batches-2024-09-24",
		)
	})

	t.Run("piping data", func(t *testing.T) {
		testFile := mocktest.TestFile(t, "Example data")
		// Test piping YAML data over stdin
		pipeDataStr := "" +
			"display_title: display_title\n" +
			"files:\n" +
			"  - Example data\n"
		pipeDataStr = strings.ReplaceAll(pipeDataStr, "Example data", testFile)
		pipeData := []byte(pipeDataStr)
		mocktest.TestRunMockTestWithPipeAndFlags(
			t, pipeData,
			"--api-key", "string",
			"beta:skills", "create",
			"--beta", "message-batches-2024-09-24",
		)
	})
}

func TestBetaSkillsRetrieve(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:skills", "retrieve",
			"--skill-id", "skill_id",
			"--beta", "message-batches-2024-09-24",
		)
	})
}

func TestBetaSkillsList(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:skills", "list",
			"--max-items", "10",
			"--limit", "0",
			"--page", "page",
			"--source", "source",
			"--beta", "message-batches-2024-09-24",
		)
	})
}

func TestBetaSkillsDelete(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:skills", "delete",
			"--skill-id", "skill_id",
			"--beta", "message-batches-2024-09-24",
		)
	})
}
