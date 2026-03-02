// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cmd

import (
	"testing"

	"github.com/anthropics/anthropic-cli/internal/mocktest"
)

func TestBetaModelsRetrieve(t *testing.T) {
	mocktest.TestRunMockTestWithFlags(
		t,
		"beta:models", "retrieve",
		"--api-key", "string",
		"--model-id", "model_id",
		"--beta", "message-batches-2024-09-24",
	)
}

func TestBetaModelsList(t *testing.T) {
	mocktest.TestRunMockTestWithFlags(
		t,
		"beta:models", "list",
		"--api-key", "string",
		"--after-id", "after_id",
		"--before-id", "before_id",
		"--limit", "1",
		"--beta", "message-batches-2024-09-24",
	)
}
