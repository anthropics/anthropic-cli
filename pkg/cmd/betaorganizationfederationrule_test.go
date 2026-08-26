// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cmd

import (
	"testing"

	"github.com/anthropics/anthropic-cli/internal/mocktest"
	"github.com/anthropics/anthropic-cli/internal/requestflag"
)

func TestBetaOrganizationFederationRulesCreate(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:organization:federation:rules", "create",
			"--issuer-id", "issuer_id",
			"--match", "{audience: audience, claims: {foo: string}, condition: condition, subject_prefix: subject_prefix}",
			"--name", "x",
			"--oauth-scope", "x",
			"--target", "{service_account_id: svac_01SDCCSbTxrXDpWc1phhtcfK, type: service_account, service_account_name: service_account_name}",
			"--applies-to-all-workspaces=true",
			"--attributes", "{foo: string}",
			"--description", "description",
			"--token-lifetime-seconds", "60",
			"--workspace-id", "workspace_id",
			"--beta", "message-batches-2024-09-24",
		)
	})

	t.Run("inner flags", func(t *testing.T) {
		// Check that inner flags have been set up correctly
		requestflag.CheckInnerFlags(betaOrganizationFederationRulesCreate)

		// Alternative argument passing style using inner flags
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:organization:federation:rules", "create",
			"--issuer-id", "issuer_id",
			"--match.audience", "audience",
			"--match.claims", "{foo: string}",
			"--match.condition", "condition",
			"--match.subject-prefix", "subject_prefix",
			"--name", "x",
			"--oauth-scope", "x",
			"--target.service-account-id", "svac_01SDCCSbTxrXDpWc1phhtcfK",
			"--target.type", "service_account",
			"--target.service-account-name", "service_account_name",
			"--applies-to-all-workspaces=true",
			"--attributes", "{foo: string}",
			"--description", "description",
			"--token-lifetime-seconds", "60",
			"--workspace-id", "workspace_id",
			"--beta", "message-batches-2024-09-24",
		)
	})

	t.Run("piping data", func(t *testing.T) {
		// Test piping YAML data over stdin
		pipeData := []byte("" +
			"issuer_id: issuer_id\n" +
			"match:\n" +
			"  audience: audience\n" +
			"  claims:\n" +
			"    foo: string\n" +
			"  condition: condition\n" +
			"  subject_prefix: subject_prefix\n" +
			"name: x\n" +
			"oauth_scope: x\n" +
			"target:\n" +
			"  service_account_id: svac_01SDCCSbTxrXDpWc1phhtcfK\n" +
			"  type: service_account\n" +
			"  service_account_name: service_account_name\n" +
			"applies_to_all_workspaces: true\n" +
			"attributes:\n" +
			"  foo: string\n" +
			"description: description\n" +
			"token_lifetime_seconds: 60\n" +
			"workspace_id: workspace_id\n")
		mocktest.TestRunMockTestWithPipeAndFlags(
			t, pipeData,
			"--api-key", "string",
			"beta:organization:federation:rules", "create",
			"--beta", "message-batches-2024-09-24",
		)
	})
}

func TestBetaOrganizationFederationRulesRetrieve(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:organization:federation:rules", "retrieve",
			"--federation-rule-id", "federation_rule_id",
			"--beta", "message-batches-2024-09-24",
		)
	})
}

func TestBetaOrganizationFederationRulesUpdate(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:organization:federation:rules", "update",
			"--federation-rule-id", "federation_rule_id",
			"--applies-to-all-workspaces=true",
			"--attributes", "{foo: string}",
			"--description", "description",
			"--match", "{audience: audience, claims: {foo: string}, condition: condition, subject_prefix: subject_prefix}",
			"--name", "x",
			"--oauth-scope", "x",
			"--target", "{service_account_id: svac_01SDCCSbTxrXDpWc1phhtcfK, type: service_account, service_account_name: service_account_name}",
			"--token-lifetime-seconds", "60",
			"--workspace-id", "workspace_id",
			"--beta", "message-batches-2024-09-24",
		)
	})

	t.Run("inner flags", func(t *testing.T) {
		// Check that inner flags have been set up correctly
		requestflag.CheckInnerFlags(betaOrganizationFederationRulesUpdate)

		// Alternative argument passing style using inner flags
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:organization:federation:rules", "update",
			"--federation-rule-id", "federation_rule_id",
			"--applies-to-all-workspaces=true",
			"--attributes", "{foo: string}",
			"--description", "description",
			"--match.audience", "audience",
			"--match.claims", "{foo: string}",
			"--match.condition", "condition",
			"--match.subject-prefix", "subject_prefix",
			"--name", "x",
			"--oauth-scope", "x",
			"--target.service-account-id", "svac_01SDCCSbTxrXDpWc1phhtcfK",
			"--target.type", "service_account",
			"--target.service-account-name", "service_account_name",
			"--token-lifetime-seconds", "60",
			"--workspace-id", "workspace_id",
			"--beta", "message-batches-2024-09-24",
		)
	})

	t.Run("piping data", func(t *testing.T) {
		// Test piping YAML data over stdin
		pipeData := []byte("" +
			"applies_to_all_workspaces: true\n" +
			"attributes:\n" +
			"  foo: string\n" +
			"description: description\n" +
			"match:\n" +
			"  audience: audience\n" +
			"  claims:\n" +
			"    foo: string\n" +
			"  condition: condition\n" +
			"  subject_prefix: subject_prefix\n" +
			"name: x\n" +
			"oauth_scope: x\n" +
			"target:\n" +
			"  service_account_id: svac_01SDCCSbTxrXDpWc1phhtcfK\n" +
			"  type: service_account\n" +
			"  service_account_name: service_account_name\n" +
			"token_lifetime_seconds: 60\n" +
			"workspace_id: workspace_id\n")
		mocktest.TestRunMockTestWithPipeAndFlags(
			t, pipeData,
			"--api-key", "string",
			"beta:organization:federation:rules", "update",
			"--federation-rule-id", "federation_rule_id",
			"--beta", "message-batches-2024-09-24",
		)
	})
}

func TestBetaOrganizationFederationRulesList(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:organization:federation:rules", "list",
			"--max-items", "10",
			"--include-archived=true",
			"--issuer-id", "issuer_id",
			"--limit", "1",
			"--page", "page",
			"--beta", "message-batches-2024-09-24",
		)
	})
}

func TestBetaOrganizationFederationRulesArchive(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"beta:organization:federation:rules", "archive",
			"--federation-rule-id", "federation_rule_id",
			"--beta", "message-batches-2024-09-24",
		)
	})
}
