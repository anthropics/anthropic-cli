// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cmd

import (
	"strings"
	"testing"

	"github.com/anthropics/anthropic-cli/internal/mocktest"
)

func TestSkillsCreate(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"skills", "create",
			"--file", mocktest.TestFile(t, "Example data"),
			"--display-name", "display_name",
			"--workspace-id", "wrkspc_011CZkZaBF1tNoB5wlCeusgy",
		)
	})

	t.Run("piping data", func(t *testing.T) {
		testFile := mocktest.TestFile(t, "Example data")
		// Test piping YAML data over stdin
		pipeDataStr := "" +
			"files:\n" +
			"  - Example data\n" +
			"display_name: display_name\n"
		pipeDataStr = strings.ReplaceAll(pipeDataStr, "Example data", testFile)
		pipeData := []byte(pipeDataStr)
		mocktest.TestRunMockTestWithPipeAndFlags(
			t, pipeData,
			"--api-key", "string",
			"skills", "create",
			"--workspace-id", "wrkspc_011CZkZaBF1tNoB5wlCeusgy",
		)
	})
}

func TestSkillsRetrieve(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"skills", "retrieve",
			"--skill-id", "skill_id",
			"--workspace-id", "wrkspc_011CZkZaBF1tNoB5wlCeusgy",
		)
	})
}

func TestSkillsList(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"skills", "list",
			"--max-items", "10",
			"--limit", "1",
			"--page", "page",
			"--source", "source",
			"--workspace-id", "wrkspc_011CZkZaBF1tNoB5wlCeusgy",
		)
	})
}

func TestSkillsDelete(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"skills", "delete",
			"--skill-id", "skill_id",
			"--workspace-id", "wrkspc_011CZkZaBF1tNoB5wlCeusgy",
		)
	})
}
