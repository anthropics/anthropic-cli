package cmd

import (
	"testing"

	"github.com/anthropics/anthropic-cli/v2/internal/mocktest"
	"github.com/anthropics/anthropic-cli/v2/internal/requestflag"
)

func TestOrganizationFederationRulesCreate(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "string")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"organization:federation:rules", "create",
			"--issuer-id", "fdis_01SDCCSbTxrXDpWc1phhtcfK",
			"--match", "{audience: audience, claims: {event_name: push, ref: refs/heads/main, repository_id: '123456789', repository_owner_id: '987654321'}, condition: condition, subject_prefix: subject_prefix}",
			"--name", "prod-deploy-pipeline",
			"--oauth-scope", "workspace:developer",
			"--target", "{service_account_id: svac_01SDCCSbTxrXDpWc1phhtcfK, type: service_account, service_account_name: service_account_name}",
			"--applies-to-all-workspaces=true",
			"--description", "description",
			"--token-lifetime-seconds", "60",
			"--workspace-id", "wrkspc_01JwQvzr7rXLA5AGx3HKfFUJ",
		)
	})

	t.Run("inner flags", func(t *testing.T) {
		// Check that inner flags have been set up correctly
		requestflag.CheckInnerFlags(organizationFederationRulesCreate)

		// Alternative argument passing style using inner flags
		mocktest.TestRunMockTestWithFlags(
			t,
			"organization:federation:rules", "create",
			"--issuer-id", "fdis_01SDCCSbTxrXDpWc1phhtcfK",
			"--match.audience", "audience",
			"--match.claims", "{event_name: push, ref: refs/heads/main, repository_id: '123456789', repository_owner_id: '987654321'}",
			"--match.condition", "condition",
			"--match.subject-prefix", "subject_prefix",
			"--name", "prod-deploy-pipeline",
			"--oauth-scope", "workspace:developer",
			"--target.service-account-id", "svac_01SDCCSbTxrXDpWc1phhtcfK",
			"--target.type", "service_account",
			"--target.service-account-name", "service_account_name",
			"--applies-to-all-workspaces=true",
			"--description", "description",
			"--token-lifetime-seconds", "60",
			"--workspace-id", "wrkspc_01JwQvzr7rXLA5AGx3HKfFUJ",
		)
	})

	t.Run("piping data", func(t *testing.T) {
		// Test piping YAML data over stdin
		pipeData := []byte("" +
			"issuer_id: fdis_01SDCCSbTxrXDpWc1phhtcfK\n" +
			"match:\n" +
			"  audience: audience\n" +
			"  claims:\n" +
			"    event_name: push\n" +
			"    ref: refs/heads/main\n" +
			"    repository_id: '123456789'\n" +
			"    repository_owner_id: '987654321'\n" +
			"  condition: condition\n" +
			"  subject_prefix: subject_prefix\n" +
			"name: prod-deploy-pipeline\n" +
			"oauth_scope: workspace:developer\n" +
			"target:\n" +
			"  service_account_id: svac_01SDCCSbTxrXDpWc1phhtcfK\n" +
			"  type: service_account\n" +
			"  service_account_name: service_account_name\n" +
			"applies_to_all_workspaces: true\n" +
			"description: description\n" +
			"token_lifetime_seconds: 60\n" +
			"workspace_id: wrkspc_01JwQvzr7rXLA5AGx3HKfFUJ\n")
		mocktest.TestRunMockTestWithPipeAndFlags(
			t, pipeData,
			"organization:federation:rules", "create",
		)
	})
}

func TestOrganizationFederationRulesRetrieve(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "string")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"organization:federation:rules", "retrieve",
			"--federation-rule-id", "federation_rule_id",
		)
	})
}

func TestOrganizationFederationRulesUpdate(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "string")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"organization:federation:rules", "update",
			"--federation-rule-id", "federation_rule_id",
			"--applies-to-all-workspaces=true",
			"--description", "description",
			"--match", "{audience: audience, claims: {foo: string}, condition: condition, subject_prefix: subject_prefix}",
			"--name", "x",
			"--oauth-scope", "x",
			"--target", "{service_account_id: svac_01SDCCSbTxrXDpWc1phhtcfK, type: service_account, service_account_name: service_account_name}",
			"--token-lifetime-seconds", "60",
			"--workspace-id", "workspace_id",
		)
	})

	t.Run("inner flags", func(t *testing.T) {
		// Check that inner flags have been set up correctly
		requestflag.CheckInnerFlags(organizationFederationRulesUpdate)

		// Alternative argument passing style using inner flags
		mocktest.TestRunMockTestWithFlags(
			t,
			"organization:federation:rules", "update",
			"--federation-rule-id", "federation_rule_id",
			"--applies-to-all-workspaces=true",
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
		)
	})

	t.Run("piping data", func(t *testing.T) {
		// Test piping YAML data over stdin
		pipeData := []byte("" +
			"applies_to_all_workspaces: true\n" +
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
			"organization:federation:rules", "update",
			"--federation-rule-id", "federation_rule_id",
		)
	})
}

func TestOrganizationFederationRulesList(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "string")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"organization:federation:rules", "list",
			"--max-items", "10",
			"--include-archived=true",
			"--issuer-id", "issuer_id",
			"--limit", "1",
			"--page", "page",
		)
	})
}

func TestOrganizationFederationRulesArchive(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "string")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"organization:federation:rules", "archive",
			"--federation-rule-id", "federation_rule_id",
		)
	})
}
