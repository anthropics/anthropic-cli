package cmd

import (
	"testing"

	"github.com/anthropics/anthropic-cli/v2/internal/mocktest"
)

func TestOrganizationRBACGroupsMembersList(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "string")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"organization:rbac-groups:members", "list",
			"--max-items", "10",
			"--rbac-group-id", "rbac_group_id",
			"--limit", "1",
			"--page", "eyJjdXJzb3IiOiAicmJhY19ncm91cF8wMSJ9",
		)
	})
}

func TestOrganizationRBACGroupsMembersAdd(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "string")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"organization:rbac-groups:members", "add",
			"--rbac-group-id", "rbac_group_id",
			"--user-id", "user_01WxYzAbCdEfGhJkMnPqRsTu",
		)
	})

	t.Run("piping data", func(t *testing.T) {
		// Test piping YAML data over stdin
		pipeData := []byte("user_id: user_01WxYzAbCdEfGhJkMnPqRsTu")
		mocktest.TestRunMockTestWithPipeAndFlags(
			t, pipeData,
			"organization:rbac-groups:members", "add",
			"--rbac-group-id", "rbac_group_id",
		)
	})
}

func TestOrganizationRBACGroupsMembersRemove(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "string")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"organization:rbac-groups:members", "remove",
			"--rbac-group-id", "rbac_group_id",
			"--user-id", "user_id",
		)
	})
}
