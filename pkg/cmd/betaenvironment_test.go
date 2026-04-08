// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cmd

import (
	"testing"

	"github.com/anthropics/anthropic-cli/internal/mocktest"
	"github.com/anthropics/anthropic-cli/internal/requestflag"
)

func TestBetaEnvironmentsCreate(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:environments", "create",
			"--name", "python-data-analysis",
			"--config", "{type: cloud, networking: {type: limited, allow_mcp_servers: true, allow_package_managers: true, allowed_hosts: [api.example.com]}, packages: {apt: [string], cargo: [string], gem: [string], go: [string], npm: [string], pip: [pandas, numpy], type: packages}}",
			"--description", "Python environment with data-analysis packages.",
			"--metadata", "{foo: string}",
			"--beta", "message-batches-2024-09-24",
		)
	})

	t.Run("inner flags", func(t *testing.T) {
		// Check that inner flags have been set up correctly
		requestflag.CheckInnerFlags(betaEnvironmentsCreate)

		// Alternative argument passing style using inner flags
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:environments", "create",
			"--name", "python-data-analysis",
			"--config.type", "cloud",
			"--config.networking", "{type: limited, allow_mcp_servers: true, allow_package_managers: true, allowed_hosts: [api.example.com]}",
			"--config.packages", "{apt: [string], cargo: [string], gem: [string], go: [string], npm: [string], pip: [pandas, numpy], type: packages}",
			"--description", "Python environment with data-analysis packages.",
			"--metadata", "{foo: string}",
			"--beta", "message-batches-2024-09-24",
		)
	})

	t.Run("piping data", func(t *testing.T) {
		// Test piping YAML data over stdin
		pipeData := []byte("" +
			"name: python-data-analysis\n" +
			"config:\n" +
			"  type: cloud\n" +
			"  networking:\n" +
			"    type: limited\n" +
			"    allow_mcp_servers: true\n" +
			"    allow_package_managers: true\n" +
			"    allowed_hosts:\n" +
			"      - api.example.com\n" +
			"  packages:\n" +
			"    apt:\n" +
			"      - string\n" +
			"    cargo:\n" +
			"      - string\n" +
			"    gem:\n" +
			"      - string\n" +
			"    go:\n" +
			"      - string\n" +
			"    npm:\n" +
			"      - string\n" +
			"    pip:\n" +
			"      - pandas\n" +
			"      - numpy\n" +
			"    type: packages\n" +
			"description: Python environment with data-analysis packages.\n" +
			"metadata:\n" +
			"  foo: string\n")
		mocktest.TestRunMockTestWithPipeAndFlags(
			t, pipeData,
			"--api-key", "string",
			"beta:environments", "create",
			"--beta", "message-batches-2024-09-24",
		)
	})
}

func TestBetaEnvironmentsRetrieve(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:environments", "retrieve",
			"--environment-id", "env_011CZkZ9X2dpNyB7HsEFoRfW",
			"--beta", "message-batches-2024-09-24",
		)
	})
}

func TestBetaEnvironmentsUpdate(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:environments", "update",
			"--environment-id", "env_011CZkZ9X2dpNyB7HsEFoRfW",
			"--config", "{type: cloud, networking: {type: limited, allow_mcp_servers: true, allow_package_managers: true, allowed_hosts: [api.example.com]}, packages: {apt: [string], cargo: [string], gem: [string], go: [string], npm: [string], pip: [pandas, numpy], type: packages}}",
			"--description", "Python environment with data-analysis packages.",
			"--metadata", "{foo: string}",
			"--name", "x",
			"--beta", "message-batches-2024-09-24",
		)
	})

	t.Run("inner flags", func(t *testing.T) {
		// Check that inner flags have been set up correctly
		requestflag.CheckInnerFlags(betaEnvironmentsUpdate)

		// Alternative argument passing style using inner flags
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:environments", "update",
			"--environment-id", "env_011CZkZ9X2dpNyB7HsEFoRfW",
			"--config.type", "cloud",
			"--config.networking", "{type: limited, allow_mcp_servers: true, allow_package_managers: true, allowed_hosts: [api.example.com]}",
			"--config.packages", "{apt: [string], cargo: [string], gem: [string], go: [string], npm: [string], pip: [pandas, numpy], type: packages}",
			"--description", "Python environment with data-analysis packages.",
			"--metadata", "{foo: string}",
			"--name", "x",
			"--beta", "message-batches-2024-09-24",
		)
	})

	t.Run("piping data", func(t *testing.T) {
		// Test piping YAML data over stdin
		pipeData := []byte("" +
			"config:\n" +
			"  type: cloud\n" +
			"  networking:\n" +
			"    type: limited\n" +
			"    allow_mcp_servers: true\n" +
			"    allow_package_managers: true\n" +
			"    allowed_hosts:\n" +
			"      - api.example.com\n" +
			"  packages:\n" +
			"    apt:\n" +
			"      - string\n" +
			"    cargo:\n" +
			"      - string\n" +
			"    gem:\n" +
			"      - string\n" +
			"    go:\n" +
			"      - string\n" +
			"    npm:\n" +
			"      - string\n" +
			"    pip:\n" +
			"      - pandas\n" +
			"      - numpy\n" +
			"    type: packages\n" +
			"description: Python environment with data-analysis packages.\n" +
			"metadata:\n" +
			"  foo: string\n" +
			"name: x\n")
		mocktest.TestRunMockTestWithPipeAndFlags(
			t, pipeData,
			"--api-key", "string",
			"beta:environments", "update",
			"--environment-id", "env_011CZkZ9X2dpNyB7HsEFoRfW",
			"--beta", "message-batches-2024-09-24",
		)
	})
}

func TestBetaEnvironmentsList(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:environments", "list",
			"--max-items", "10",
			"--include-archived=true",
			"--limit", "1",
			"--page", "page",
			"--beta", "message-batches-2024-09-24",
		)
	})
}

func TestBetaEnvironmentsDelete(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:environments", "delete",
			"--environment-id", "env_011CZkZ9X2dpNyB7HsEFoRfW",
			"--beta", "message-batches-2024-09-24",
		)
	})
}

func TestBetaEnvironmentsArchive(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:environments", "archive",
			"--environment-id", "env_011CZkZ9X2dpNyB7HsEFoRfW",
			"--beta", "message-batches-2024-09-24",
		)
	})
}
