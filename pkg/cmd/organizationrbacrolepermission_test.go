package cmd

import (
	"testing"

	"github.com/anthropics/anthropic-cli/v2/internal/mocktest"
)

func TestOrganizationRBACRolesPermissionsList(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "string")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"organization:rbac-roles:permissions", "list",
			"--max-items", "10",
			"--rbac-role-id", "rbac_role_id",
			"--limit", "1",
			"--page", "eyJjdXJzb3IiOiAicmJhY19yb2xlXzAxIn0",
		)
	})
}
