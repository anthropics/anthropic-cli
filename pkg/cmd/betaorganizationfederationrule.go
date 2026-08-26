// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cmd

import (
	"context"
	"fmt"

	"github.com/anthropics/anthropic-cli/internal/apiquery"
	"github.com/anthropics/anthropic-cli/internal/requestflag"
	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/tidwall/gjson"
	"github.com/urfave/cli/v3"
)

var betaOrganizationFederationRulesCreate = requestflag.WithInnerFlags(cli.Command{
	Name:    "create",
	Usage:   "**Requires an OAuth access token with the `org:admin` scope**, from\n`ant auth login --scope org:admin` or a workload identity federation rule; Admin\nAPI keys are not accepted. See\n[Manage WIF with the Admin API](/docs/en/manage-claude/wif-admin-api).",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:     "issuer-id",
			Usage:    "Tagged ID of the federation issuer.",
			Required: true,
			BodyPath: "issuer_id",
		},
		&requestflag.Flag[map[string]any]{
			Name:     "match",
			Usage:    "Does the incoming JWT qualify?\n\nAll populated fields must pass; omitted fields are skipped. At least one\nof `subject_prefix` (other than a wildcard-only value like `*`), `claims`,\nor `condition` is required; `audience` alone is not sufficient.",
			Required: true,
			BodyPath: "match",
		},
		&requestflag.Flag[string]{
			Name:     "name",
			Usage:    "Slug identifier (lowercase, digits, hyphens). Unique within the organization; a duplicate name returns 409.",
			Required: true,
			BodyPath: "name",
		},
		&requestflag.Flag[string]{
			Name:     "oauth-scope",
			Usage:    "Space-separated OAuth scopes. OAuth callers may only set `workspace:developer` or `workspace:inference`; other scopes (such as `org:admin`) require a Console session.",
			Required: true,
			BodyPath: "oauth_scope",
		},
		&requestflag.Flag[map[string]any]{
			Name:     "target",
			Usage:    "Bind to a fixed service account by ID.",
			Required: true,
			BodyPath: "target",
		},
		&requestflag.Flag[bool]{
			Name:     "applies-to-all-workspaces",
			Usage:    "When true, enable this rule for every workspace in the org (including workspaces created later).",
			BodyPath: "applies_to_all_workspaces",
		},
		&requestflag.Flag[map[string]any]{
			Name:     "attributes",
			Usage:    "CEL expressions `{name: expr}` extracting named values from claims. Not yet supported; any non-empty value is rejected with 400.",
			BodyPath: "attributes",
		},
		&requestflag.Flag[*string]{
			Name:     "description",
			Usage:    "Optional free-text description.",
			BodyPath: "description",
		},
		&requestflag.Flag[int64]{
			Name:     "token-lifetime-seconds",
			Usage:    "Lifetime in seconds for access tokens minted via this rule (60-86400). Defaults to 3600 (1h). Minted tokens are capped at `max(60, min(this value, 2 × remaining assertion validity))` seconds.",
			BodyPath: "token_lifetime_seconds",
		},
		&requestflag.Flag[*string]{
			Name:     "workspace-id",
			Usage:    "Tagged ID of the workspace to enable this rule for. Required unless `applies_to_all_workspaces` is true. Additional workspaces can be added via the `/federation_rules/{federation_rule_id}/workspaces` sub-resource.",
			BodyPath: "workspace_id",
		},
		&requestflag.Flag[[]string]{
			Name:       "beta",
			Usage:      "Optional header to specify the beta version(s) you want to use.",
			HeaderPath: "anthropic-beta",
		},
	},
	Action:          handleBetaOrganizationFederationRulesCreate,
	HideHelpCommand: true,
}, map[string][]requestflag.HasOuterFlag{
	"match": {
		&requestflag.InnerFlag[*string]{
			Name:       "match.audience",
			Usage:      "Exact match against the `aud` claim (any element if array). When omitted, the JWT's `aud` must still equal Anthropic's expected audience for the issuer; setting this field overrides that default.",
			InnerField: "audience",
		},
		&requestflag.InnerFlag[map[string]any]{
			Name:       "match.claims",
			Usage:      "Exact-match `{claim: value}` pairs against top-level claims. Only string-valued claims can be matched; use `condition` for non-string claims.",
			InnerField: "claims",
		},
		&requestflag.InnerFlag[*string]{
			Name:       "match.condition",
			Usage:      "CEL expression over claims for logic the structural fields can't express. Must evaluate to a boolean and may reference only the `claims` variable; a constant-true expression (such as `true`) is rejected with 400.",
			InnerField: "condition",
		},
		&requestflag.InnerFlag[*string]{
			Name:       "match.subject-prefix",
			Usage:      "Match the verified JWT `sub` claim. Exact match unless the value ends with `*`, in which case it is a prefix match. Example: `repo:my-org/my-repo:ref:refs/heads/main`.",
			InnerField: "subject_prefix",
		},
	},
	"target": {
		&requestflag.InnerFlag[string]{
			Name:       "target.service-account-id",
			Usage:      "Tagged ID of the service account to mint tokens for.",
			InnerField: "service_account_id",
		},
		&requestflag.InnerFlag[string]{
			Name:       "target.type",
			Usage:      `Allowed values: "service_account".`,
			InnerField: "type",
		},
		&requestflag.InnerFlag[*string]{
			Name:       "target.service-account-name",
			Usage:      "Service account's display name at read time. Ignored on writes.",
			InnerField: "service_account_name",
		},
	},
})

var betaOrganizationFederationRulesRetrieve = cli.Command{
	Name:    "retrieve",
	Usage:   "**Requires an OAuth access token with the `org:admin` scope**, from\n`ant auth login --scope org:admin` or a workload identity federation rule; Admin\nAPI keys are not accepted. See\n[Manage WIF with the Admin API](/docs/en/manage-claude/wif-admin-api).",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "federation-rule-id",
			Usage:     "ID of the federation rule.",
			Required:  true,
			PathParam: "federation_rule_id",
		},
		&requestflag.Flag[[]string]{
			Name:       "beta",
			Usage:      "Optional header to specify the beta version(s) you want to use.",
			HeaderPath: "anthropic-beta",
		},
	},
	Action:          handleBetaOrganizationFederationRulesRetrieve,
	HideHelpCommand: true,
}

var betaOrganizationFederationRulesUpdate = requestflag.WithInnerFlags(cli.Command{
	Name:    "update",
	Usage:   "**Requires an OAuth access token with the `org:admin` scope**, from\n`ant auth login --scope org:admin` or a workload identity federation rule; Admin\nAPI keys are not accepted. See\n[Manage WIF with the Admin API](/docs/en/manage-claude/wif-admin-api).",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "federation-rule-id",
			Usage:     "ID of the federation rule to update.",
			Required:  true,
			PathParam: "federation_rule_id",
		},
		&requestflag.Flag[*bool]{
			Name:     "applies-to-all-workspaces",
			Usage:    "When true, enables this rule for every workspace in the org (including workspaces created later). Setting `false` is rejected with 400 if no workspace would remain enabled; a rule with only a legacy `workspace_id` binding continues to mint.",
			BodyPath: "applies_to_all_workspaces",
		},
		&requestflag.Flag[map[string]any]{
			Name:     "attributes",
			Usage:    "Replaces the CEL expressions `{name: expr}` extracting named values from claims. Send null to clear them. Not yet supported; any non-empty value is rejected with 400.",
			BodyPath: "attributes",
		},
		&requestflag.Flag[*string]{
			Name:     "description",
			Usage:    "Replaces the description. Omit to leave unchanged; send `null` to clear (the field is stored as an empty string).",
			BodyPath: "description",
		},
		&requestflag.Flag[map[string]any]{
			Name:     "match",
			Usage:    "Does the incoming JWT qualify?\n\nAll populated fields must pass; omitted fields are skipped. At least one\nof `subject_prefix` (other than a wildcard-only value like `*`), `claims`,\nor `condition` is required; `audience` alone is not sufficient.",
			BodyPath: "match",
		},
		&requestflag.Flag[*string]{
			Name:     "name",
			Usage:    "Replaces the slug identifier (lowercase, digits, hyphens). Unique within the organization; a duplicate name returns 409.",
			BodyPath: "name",
		},
		&requestflag.Flag[*string]{
			Name:     "oauth-scope",
			Usage:    "Replaces the space-separated OAuth scopes granted on minted tokens. OAuth callers may only set `workspace:developer` or `workspace:inference`; other scopes (such as `org:admin`) require a Console session.",
			BodyPath: "oauth_scope",
		},
		&requestflag.Flag[map[string]any]{
			Name:     "target",
			Usage:    "Bind to a fixed service account by ID.",
			BodyPath: "target",
		},
		&requestflag.Flag[*int64]{
			Name:     "token-lifetime-seconds",
			Usage:    "Replaces the lifetime in seconds for access tokens minted via this rule (60-86400). Minted tokens are capped at `max(60, min(this value, 2 × remaining assertion validity))` seconds.",
			BodyPath: "token_lifetime_seconds",
		},
		&requestflag.Flag[*string]{
			Name:     "workspace-id",
			Usage:    "Replaces the existing single workspace enablement (the previous one is removed). Rejected with 400 if the rule is enabled for more than one workspace; use the `/federation_rules/{federation_rule_id}/workspaces` sub-resource instead.",
			BodyPath: "workspace_id",
		},
		&requestflag.Flag[[]string]{
			Name:       "beta",
			Usage:      "Optional header to specify the beta version(s) you want to use.",
			HeaderPath: "anthropic-beta",
		},
	},
	Action:          handleBetaOrganizationFederationRulesUpdate,
	HideHelpCommand: true,
}, map[string][]requestflag.HasOuterFlag{
	"match": {
		&requestflag.InnerFlag[*string]{
			Name:       "match.audience",
			Usage:      "Exact match against the `aud` claim (any element if array). When omitted, the JWT's `aud` must still equal Anthropic's expected audience for the issuer; setting this field overrides that default.",
			InnerField: "audience",
		},
		&requestflag.InnerFlag[map[string]any]{
			Name:       "match.claims",
			Usage:      "Exact-match `{claim: value}` pairs against top-level claims. Only string-valued claims can be matched; use `condition` for non-string claims.",
			InnerField: "claims",
		},
		&requestflag.InnerFlag[*string]{
			Name:       "match.condition",
			Usage:      "CEL expression over claims for logic the structural fields can't express. Must evaluate to a boolean and may reference only the `claims` variable; a constant-true expression (such as `true`) is rejected with 400.",
			InnerField: "condition",
		},
		&requestflag.InnerFlag[*string]{
			Name:       "match.subject-prefix",
			Usage:      "Match the verified JWT `sub` claim. Exact match unless the value ends with `*`, in which case it is a prefix match. Example: `repo:my-org/my-repo:ref:refs/heads/main`.",
			InnerField: "subject_prefix",
		},
	},
	"target": {
		&requestflag.InnerFlag[string]{
			Name:       "target.service-account-id",
			Usage:      "Tagged ID of the service account to mint tokens for.",
			InnerField: "service_account_id",
		},
		&requestflag.InnerFlag[string]{
			Name:       "target.type",
			Usage:      `Allowed values: "service_account".`,
			InnerField: "type",
		},
		&requestflag.InnerFlag[*string]{
			Name:       "target.service-account-name",
			Usage:      "Service account's display name at read time. Ignored on writes.",
			InnerField: "service_account_name",
		},
	},
})

var betaOrganizationFederationRulesList = cli.Command{
	Name:    "list",
	Usage:   "**Requires an OAuth access token with the `org:admin` scope**, from\n`ant auth login --scope org:admin` or a workload identity federation rule; Admin\nAPI keys are not accepted. See\n[Manage WIF with the Admin API](/docs/en/manage-claude/wif-admin-api).",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[bool]{
			Name:      "include-archived",
			Usage:     "Include archived resources. Defaults to false.",
			Default:   false,
			QueryPath: "include_archived",
		},
		&requestflag.Flag[string]{
			Name:      "issuer-id",
			Usage:     "Filter to rules referencing this federation issuer.",
			QueryPath: "issuer_id",
		},
		&requestflag.Flag[int64]{
			Name:      "limit",
			Usage:     "Number of results per page.",
			Default:   20,
			QueryPath: "limit",
		},
		&requestflag.Flag[string]{
			Name:      "page",
			Usage:     "Opaque cursor from a previous response's `next_page`.",
			QueryPath: "page",
		},
		&requestflag.Flag[[]string]{
			Name:       "beta",
			Usage:      "Optional header to specify the beta version(s) you want to use.",
			HeaderPath: "anthropic-beta",
		},
		&requestflag.Flag[int64]{
			Name:  "max-items",
			Usage: "The maximum number of items to return (use -1 for unlimited).",
		},
	},
	Action:          handleBetaOrganizationFederationRulesList,
	HideHelpCommand: true,
}

var betaOrganizationFederationRulesArchive = cli.Command{
	Name:    "archive",
	Usage:   "**Requires an OAuth access token with the `org:admin` scope**, from\n`ant auth login --scope org:admin` or a workload identity federation rule; Admin\nAPI keys are not accepted. See\n[Manage WIF with the Admin API](/docs/en/manage-claude/wif-admin-api).",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "federation-rule-id",
			Usage:     "ID of the federation rule to archive.",
			Required:  true,
			PathParam: "federation_rule_id",
		},
		&requestflag.Flag[[]string]{
			Name:       "beta",
			Usage:      "Optional header to specify the beta version(s) you want to use.",
			HeaderPath: "anthropic-beta",
		},
	},
	Action:          handleBetaOrganizationFederationRulesArchive,
	HideHelpCommand: true,
}

func handleBetaOrganizationFederationRulesCreate(ctx context.Context, cmd *cli.Command) error {
	client := anthropic.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()

	if len(unusedArgs) > 0 {
		return fmt.Errorf("Unexpected extra arguments: %v", unusedArgs)
	}

	options, err := flagOptions(
		cmd,
		apiquery.NestedQueryFormatBrackets,
		apiquery.ArrayQueryFormatBrackets,
		ApplicationJSON,
		false,
	)
	if err != nil {
		return err
	}

	params := anthropic.BetaOrganizationFederationRuleNewParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Beta.Organization.Federation.Rules.New(ctx, params, options...)
	if err != nil {
		return err
	}

	obj := gjson.ParseBytes(res)
	format := cmd.Root().String("format")
	explicitFormat := cmd.Root().IsSet("format")
	transform := cmd.Root().String("transform")
	return ShowJSON(obj, ShowJSONOpts{
		ExplicitFormat: explicitFormat,
		Format:         format,
		RawOutput:      cmd.Root().Bool("raw-output"),
		Title:          "beta:organization:federation:rules create",
		Transform:      transform,
	})
}

func handleBetaOrganizationFederationRulesRetrieve(ctx context.Context, cmd *cli.Command) error {
	client := anthropic.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("federation-rule-id") && len(unusedArgs) > 0 {
		cmd.Set("federation-rule-id", unusedArgs[0])
		unusedArgs = unusedArgs[1:]
	}
	if len(unusedArgs) > 0 {
		return fmt.Errorf("Unexpected extra arguments: %v", unusedArgs)
	}

	options, err := flagOptions(
		cmd,
		apiquery.NestedQueryFormatBrackets,
		apiquery.ArrayQueryFormatBrackets,
		EmptyBody,
		false,
	)
	if err != nil {
		return err
	}

	params := anthropic.BetaOrganizationFederationRuleGetParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Beta.Organization.Federation.Rules.Get(
		ctx,
		cmd.Value("federation-rule-id").(string),
		params,
		options...,
	)
	if err != nil {
		return err
	}

	obj := gjson.ParseBytes(res)
	format := "explore"
	explicitFormat := cmd.Root().IsSet("format")
	if explicitFormat {
		format = cmd.Root().String("format")
	}
	transform := cmd.Root().String("transform")
	return ShowJSON(obj, ShowJSONOpts{
		ExplicitFormat: explicitFormat,
		Format:         format,
		RawOutput:      cmd.Root().Bool("raw-output"),
		Title:          "beta:organization:federation:rules retrieve",
		Transform:      transform,
	})
}

func handleBetaOrganizationFederationRulesUpdate(ctx context.Context, cmd *cli.Command) error {
	client := anthropic.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("federation-rule-id") && len(unusedArgs) > 0 {
		cmd.Set("federation-rule-id", unusedArgs[0])
		unusedArgs = unusedArgs[1:]
	}
	if len(unusedArgs) > 0 {
		return fmt.Errorf("Unexpected extra arguments: %v", unusedArgs)
	}

	options, err := flagOptions(
		cmd,
		apiquery.NestedQueryFormatBrackets,
		apiquery.ArrayQueryFormatBrackets,
		ApplicationJSON,
		false,
	)
	if err != nil {
		return err
	}

	params := anthropic.BetaOrganizationFederationRuleUpdateParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Beta.Organization.Federation.Rules.Update(
		ctx,
		cmd.Value("federation-rule-id").(string),
		params,
		options...,
	)
	if err != nil {
		return err
	}

	obj := gjson.ParseBytes(res)
	format := cmd.Root().String("format")
	explicitFormat := cmd.Root().IsSet("format")
	transform := cmd.Root().String("transform")
	return ShowJSON(obj, ShowJSONOpts{
		ExplicitFormat: explicitFormat,
		Format:         format,
		RawOutput:      cmd.Root().Bool("raw-output"),
		Title:          "beta:organization:federation:rules update",
		Transform:      transform,
	})
}

func handleBetaOrganizationFederationRulesList(ctx context.Context, cmd *cli.Command) error {
	client := anthropic.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()

	if len(unusedArgs) > 0 {
		return fmt.Errorf("Unexpected extra arguments: %v", unusedArgs)
	}

	options, err := flagOptions(
		cmd,
		apiquery.NestedQueryFormatBrackets,
		apiquery.ArrayQueryFormatBrackets,
		EmptyBody,
		false,
	)
	if err != nil {
		return err
	}

	params := anthropic.BetaOrganizationFederationRuleListParams{}

	format := "explore"
	explicitFormat := cmd.Root().IsSet("format")
	if explicitFormat {
		format = cmd.Root().String("format")
	}
	transform := cmd.Root().String("transform")
	if format == "raw" {
		var res []byte
		options = append(options, option.WithResponseBodyInto(&res))
		_, err = client.Beta.Organization.Federation.Rules.List(ctx, params, options...)
		if err != nil {
			return err
		}
		obj := gjson.ParseBytes(res)
		return ShowJSON(obj, ShowJSONOpts{
			ExplicitFormat: explicitFormat,
			Format:         format,
			RawOutput:      cmd.Root().Bool("raw-output"),
			Title:          "beta:organization:federation:rules list",
			Transform:      transform,
		})
	} else {
		iter := client.Beta.Organization.Federation.Rules.ListAutoPaging(ctx, params, options...)
		maxItems := int64(-1)
		if cmd.IsSet("max-items") {
			maxItems = cmd.Value("max-items").(int64)
		}
		return ShowJSONIterator(iter, maxItems, ShowJSONOpts{
			ExplicitFormat: explicitFormat,
			Format:         format,
			RawOutput:      cmd.Root().Bool("raw-output"),
			Title:          "beta:organization:federation:rules list",
			Transform:      transform,
		})
	}
}

func handleBetaOrganizationFederationRulesArchive(ctx context.Context, cmd *cli.Command) error {
	client := anthropic.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("federation-rule-id") && len(unusedArgs) > 0 {
		cmd.Set("federation-rule-id", unusedArgs[0])
		unusedArgs = unusedArgs[1:]
	}
	if len(unusedArgs) > 0 {
		return fmt.Errorf("Unexpected extra arguments: %v", unusedArgs)
	}

	options, err := flagOptions(
		cmd,
		apiquery.NestedQueryFormatBrackets,
		apiquery.ArrayQueryFormatBrackets,
		EmptyBody,
		false,
	)
	if err != nil {
		return err
	}

	params := anthropic.BetaOrganizationFederationRuleArchiveParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Beta.Organization.Federation.Rules.Archive(
		ctx,
		cmd.Value("federation-rule-id").(string),
		params,
		options...,
	)
	if err != nil {
		return err
	}

	obj := gjson.ParseBytes(res)
	format := cmd.Root().String("format")
	explicitFormat := cmd.Root().IsSet("format")
	transform := cmd.Root().String("transform")
	return ShowJSON(obj, ShowJSONOpts{
		ExplicitFormat: explicitFormat,
		Format:         format,
		RawOutput:      cmd.Root().Bool("raw-output"),
		Title:          "beta:organization:federation:rules archive",
		Transform:      transform,
	})
}
