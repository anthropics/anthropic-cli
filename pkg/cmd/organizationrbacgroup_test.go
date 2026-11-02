package cmd

import (
	"testing"

	"github.com/anthropics/anthropic-cli/v2/internal/mocktest"
)

func TestOrganizationRBACGroupsCreate(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "string")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"organization:rbac-groups", "create",
			"--name", "Engineering",
		)
	})

	t.Run("piping data", func(t *testing.T) {
		// Test piping YAML data over stdin
		pipeData := []byte("name: Engineering")
		mocktest.TestRunMockTestWithPipeAndFlags(
			t, pipeData,
			"organization:rbac-groups", "create",
		)
	})
}

func TestOrganizationRBACGroupsRetrieve(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "string")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"organization:rbac-groups", "retrieve",
			"--rbac-group-id", "rbac_group_id",
		)
	})
}

func TestOrganizationRBACGroupsUpdate(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "string")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"organization:rbac-groups", "update",
			"--rbac-group-id", "rbac_group_id",
			"--name", "Engineering",
		)
	})

	t.Run("piping data", func(t *testing.T) {
		// Test piping YAML data over stdin
		pipeData := []byte("name: Engineering")
		mocktest.TestRunMockTestWithPipeAndFlags(
			t, pipeData,
			"organization:rbac-groups", "update",
			"--rbac-group-id", "rbac_group_id",
		)
	})
}

func TestOrganizationRBACGroupsList(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "string")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"organization:rbac-groups", "list",
			"--max-items", "10",
			"--limit", "1",
			"--page", "eyJjdXJzb3IiOiAicmJhY19ncm91cF8wMSJ9",
		)
	})
}

func TestOrganizationRBACGroupsDelete(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "string")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"organization:rbac-groups", "delete",
			"--rbac-group-id", "rbac_group_id",
		)
	})
}
