// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cmd

import (
	"testing"

	"github.com/anthropics/anthropic-cli/internal/mocktest"
)

func TestModelsRetrieve(t *testing.T) {
	mocktest.TestRunMockTestWithFlags(
		t,
		"models", "retrieve",
		"--api-key", "string",
		"--model-id", "model_id",
		"--beta", "message-batches-2024-09-24",
	)
}

func TestModelsList(t *testing.T) {
	mocktest.TestRunMockTestWithFlags(
		t,
		"models", "list",
		"--api-key", "string",
		"--after-id", "after_id",
		"--before-id", "before_id",
		"--limit", "1",
		"--beta", "message-batches-2024-09-24",
	)
}
