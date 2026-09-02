// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cmd

import (
	"strings"
	"testing"

	"github.com/anthropics/anthropic-cli/internal/mocktest"
)

func TestSkillsVersionsCreate(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"skills:versions", "create",
			"--skill-id", "skill_id",
			"--file", mocktest.TestFile(t, "Example data"),
			"--workspace-id", "wrkspc_011CZkZaBF1tNoB5wlCeusgy",
		)
	})

	t.Run("piping data", func(t *testing.T) {
		testFile := mocktest.TestFile(t, "Example data")
		// Test piping YAML data over stdin
		pipeDataStr := "" +
			"files:\n" +
			"  - Example data\n"
		pipeDataStr = strings.ReplaceAll(pipeDataStr, "Example data", testFile)
		pipeData := []byte(pipeDataStr)
		mocktest.TestRunMockTestWithPipeAndFlags(
			t, pipeData,
			"--api-key", "string",
			"skills:versions", "create",
			"--skill-id", "skill_id",
			"--workspace-id", "wrkspc_011CZkZaBF1tNoB5wlCeusgy",
		)
	})
}

func TestSkillsVersionsRetrieve(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"skills:versions", "retrieve",
			"--skill-id", "skill_id",
			"--version", "version",
			"--workspace-id", "wrkspc_011CZkZaBF1tNoB5wlCeusgy",
		)
	})
}

func TestSkillsVersionsList(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"skills:versions", "list",
			"--max-items", "10",
			"--skill-id", "skill_id",
			"--limit", "1",
			"--page", "page",
			"--workspace-id", "wrkspc_011CZkZaBF1tNoB5wlCeusgy",
		)
	})
}

func TestSkillsVersionsDelete(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"skills:versions", "delete",
			"--skill-id", "skill_id",
			"--version", "version",
			"--workspace-id", "wrkspc_011CZkZaBF1tNoB5wlCeusgy",
		)
	})
}
