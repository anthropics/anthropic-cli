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

var betaOrganizationFederationIssuersCreate = requestflag.WithInnerFlags(cli.Command{
	Name:    "create",
	Usage:   "**Requires an OAuth access token with the `org:admin` scope**, from\n`ant auth login --scope org:admin` or a workload identity federation rule; Admin\nAPI keys are not accepted. See\n[Manage WIF with the Admin API](/docs/en/manage-claude/wif-admin-api).",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:     "issuer-url",
			Usage:    "The `iss` claim value to match against.",
			Required: true,
			BodyPath: "issuer_url",
		},
		&requestflag.Flag[string]{
			Name:     "name",
			Usage:    "Slug identifier (lowercase, digits, hyphens). Unique within the organization; a duplicate name returns 409.",
			Required: true,
			BodyPath: "name",
		},
		&requestflag.Flag[*bool]{
			Name:     "check-jti",
			Usage:    "Whether the jwt-bearer exchange enforces JTI single-use (replay protection) for tokens from this issuer. Defaults to true. Applies only to assertions carrying a `jti` claim; tokens without one are accepted without single-use enforcement.",
			BodyPath: "check_jti",
		},
		&requestflag.Flag[map[string]any]{
			Name:     "jwks",
			Usage:    "How signing keys are obtained. Defaults to OIDC discovery.",
			BodyPath: "jwks",
		},
		&requestflag.Flag[*int64]{
			Name:     "max-jwt-lifetime-seconds",
			Usage:    "Maximum allowed iat→exp spread for assertions from this issuer (1-176400 seconds, i.e. up to 49h). Defaults to 3600 (1h). Assertions must carry both `iat` and `exp`; a missing `iat` is rejected.",
			BodyPath: "max_jwt_lifetime_seconds",
		},
		&requestflag.Flag[[]string]{
			Name:       "beta",
			Usage:      "Optional header to specify the beta version(s) you want to use.",
			HeaderPath: "anthropic-beta",
		},
	},
	Action:          handleBetaOrganizationFederationIssuersCreate,
	HideHelpCommand: true,
}, map[string][]requestflag.HasOuterFlag{
	"jwks": {
		&requestflag.InnerFlag[string]{
			Name:       "jwks.type",
			Usage:      `Allowed values: "discovery", "explicit_url", "inline".`,
			InnerField: "type",
		},
		&requestflag.InnerFlag[*string]{
			Name:       "jwks.ca-cert-pem",
			Usage:      "Optional custom CA (PEM) for TLS verification of the JWKS fetch.",
			InnerField: "ca_cert_pem",
		},
		&requestflag.InnerFlag[*string]{
			Name:       "jwks.discovery-base",
			Usage:      "Set when the discovery URL differs from `issuer_url`.",
			InnerField: "discovery_base",
		},
		&requestflag.InnerFlag[any]{
			Name:       "jwks.keys",
			InnerField: "keys",
		},
		&requestflag.InnerFlag[string]{
			Name:       "jwks.url",
			Usage:      "JWKS endpoint.",
			InnerField: "url",
		},
	},
})

var betaOrganizationFederationIssuersRetrieve = cli.Command{
	Name:    "retrieve",
	Usage:   "**Requires an OAuth access token with the `org:admin` scope**, from\n`ant auth login --scope org:admin` or a workload identity federation rule; Admin\nAPI keys are not accepted. See\n[Manage WIF with the Admin API](/docs/en/manage-claude/wif-admin-api).",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "federation-issuer-id",
			Usage:     "ID of the federation issuer.",
			Required:  true,
			PathParam: "federation_issuer_id",
		},
		&requestflag.Flag[[]string]{
			Name:       "beta",
			Usage:      "Optional header to specify the beta version(s) you want to use.",
			HeaderPath: "anthropic-beta",
		},
	},
	Action:          handleBetaOrganizationFederationIssuersRetrieve,
	HideHelpCommand: true,
}

var betaOrganizationFederationIssuersUpdate = requestflag.WithInnerFlags(cli.Command{
	Name:    "update",
	Usage:   "**Requires an OAuth access token with the `org:admin` scope**, from\n`ant auth login --scope org:admin` or a workload identity federation rule; Admin\nAPI keys are not accepted. See\n[Manage WIF with the Admin API](/docs/en/manage-claude/wif-admin-api).",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "federation-issuer-id",
			Usage:     "ID of the federation issuer to update.",
			Required:  true,
			PathParam: "federation_issuer_id",
		},
		&requestflag.Flag[*bool]{
			Name:     "check-jti",
			Usage:    "Whether the jwt-bearer exchange enforces JTI single-use (replay protection) for tokens from this issuer. Applies only to assertions carrying a `jti` claim; tokens without one are accepted without single-use enforcement.",
			BodyPath: "check_jti",
		},
		&requestflag.Flag[*string]{
			Name:     "issuer-url",
			Usage:    "Replaces the `iss` claim value to match against. For discovery-mode issuers without a `discovery_base`, this is also the URL Anthropic fetches the OIDC discovery document and signing keys from, so changing it repoints the JWKS source. Changing the issuer URL to a well-known shared platform is rejected while any live rule under this issuer would not constrain tenant identity.",
			BodyPath: "issuer_url",
		},
		&requestflag.Flag[map[string]any]{
			Name:     "jwks",
			Usage:    "Replaces the entire JWKS configuration.",
			BodyPath: "jwks",
		},
		&requestflag.Flag[*bool]{
			Name:     "jwks-polling-disabled",
			Usage:    "Only `false` is accepted, to re-enable polling after the system pauses it. Polling is paused automatically; sending `true` is rejected.",
			BodyPath: "jwks_polling_disabled",
		},
		&requestflag.Flag[*int64]{
			Name:     "max-jwt-lifetime-seconds",
			Usage:    "Maximum allowed iat→exp spread for assertions from this issuer (1-176400 seconds, i.e. up to 49h). Assertions must carry both `iat` and `exp`; a missing `iat` is rejected.",
			BodyPath: "max_jwt_lifetime_seconds",
		},
		&requestflag.Flag[*string]{
			Name:     "name",
			Usage:    "Replaces the slug identifier (lowercase, digits, hyphens). Unique within the organization; a duplicate name returns 409.",
			BodyPath: "name",
		},
		&requestflag.Flag[[]string]{
			Name:       "beta",
			Usage:      "Optional header to specify the beta version(s) you want to use.",
			HeaderPath: "anthropic-beta",
		},
	},
	Action:          handleBetaOrganizationFederationIssuersUpdate,
	HideHelpCommand: true,
}, map[string][]requestflag.HasOuterFlag{
	"jwks": {
		&requestflag.InnerFlag[string]{
			Name:       "jwks.type",
			Usage:      `Allowed values: "discovery", "explicit_url", "inline".`,
			InnerField: "type",
		},
		&requestflag.InnerFlag[*string]{
			Name:       "jwks.ca-cert-pem",
			Usage:      "Optional custom CA (PEM) for TLS verification of the JWKS fetch.",
			InnerField: "ca_cert_pem",
		},
		&requestflag.InnerFlag[*string]{
			Name:       "jwks.discovery-base",
			Usage:      "Set when the discovery URL differs from `issuer_url`.",
			InnerField: "discovery_base",
		},
		&requestflag.InnerFlag[any]{
			Name:       "jwks.keys",
			InnerField: "keys",
		},
		&requestflag.InnerFlag[string]{
			Name:       "jwks.url",
			Usage:      "JWKS endpoint.",
			InnerField: "url",
		},
	},
})

var betaOrganizationFederationIssuersList = cli.Command{
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
	Action:          handleBetaOrganizationFederationIssuersList,
	HideHelpCommand: true,
}

var betaOrganizationFederationIssuersArchive = cli.Command{
	Name:    "archive",
	Usage:   "**Requires an OAuth access token with the `org:admin` scope**, from\n`ant auth login --scope org:admin` or a workload identity federation rule; Admin\nAPI keys are not accepted. See\n[Manage WIF with the Admin API](/docs/en/manage-claude/wif-admin-api).",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "federation-issuer-id",
			Usage:     "ID of the federation issuer to archive.",
			Required:  true,
			PathParam: "federation_issuer_id",
		},
		&requestflag.Flag[[]string]{
			Name:       "beta",
			Usage:      "Optional header to specify the beta version(s) you want to use.",
			HeaderPath: "anthropic-beta",
		},
	},
	Action:          handleBetaOrganizationFederationIssuersArchive,
	HideHelpCommand: true,
}

func handleBetaOrganizationFederationIssuersCreate(ctx context.Context, cmd *cli.Command) error {
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

	params := anthropic.BetaOrganizationFederationIssuerNewParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Beta.Organization.Federation.Issuers.New(ctx, params, options...)
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
		Title:          "beta:organization:federation:issuers create",
		Transform:      transform,
	})
}

func handleBetaOrganizationFederationIssuersRetrieve(ctx context.Context, cmd *cli.Command) error {
	client := anthropic.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("federation-issuer-id") && len(unusedArgs) > 0 {
		cmd.Set("federation-issuer-id", unusedArgs[0])
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

	params := anthropic.BetaOrganizationFederationIssuerGetParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Beta.Organization.Federation.Issuers.Get(
		ctx,
		cmd.Value("federation-issuer-id").(string),
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
		Title:          "beta:organization:federation:issuers retrieve",
		Transform:      transform,
	})
}

func handleBetaOrganizationFederationIssuersUpdate(ctx context.Context, cmd *cli.Command) error {
	client := anthropic.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("federation-issuer-id") && len(unusedArgs) > 0 {
		cmd.Set("federation-issuer-id", unusedArgs[0])
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

	params := anthropic.BetaOrganizationFederationIssuerUpdateParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Beta.Organization.Federation.Issuers.Update(
		ctx,
		cmd.Value("federation-issuer-id").(string),
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
		Title:          "beta:organization:federation:issuers update",
		Transform:      transform,
	})
}

func handleBetaOrganizationFederationIssuersList(ctx context.Context, cmd *cli.Command) error {
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

	params := anthropic.BetaOrganizationFederationIssuerListParams{}

	format := "explore"
	explicitFormat := cmd.Root().IsSet("format")
	if explicitFormat {
		format = cmd.Root().String("format")
	}
	transform := cmd.Root().String("transform")
	if format == "raw" {
		var res []byte
		options = append(options, option.WithResponseBodyInto(&res))
		_, err = client.Beta.Organization.Federation.Issuers.List(ctx, params, options...)
		if err != nil {
			return err
		}
		obj := gjson.ParseBytes(res)
		return ShowJSON(obj, ShowJSONOpts{
			ExplicitFormat: explicitFormat,
			Format:         format,
			RawOutput:      cmd.Root().Bool("raw-output"),
			Title:          "beta:organization:federation:issuers list",
			Transform:      transform,
		})
	} else {
		iter := client.Beta.Organization.Federation.Issuers.ListAutoPaging(ctx, params, options...)
		maxItems := int64(-1)
		if cmd.IsSet("max-items") {
			maxItems = cmd.Value("max-items").(int64)
		}
		return ShowJSONIterator(iter, maxItems, ShowJSONOpts{
			ExplicitFormat: explicitFormat,
			Format:         format,
			RawOutput:      cmd.Root().Bool("raw-output"),
			Title:          "beta:organization:federation:issuers list",
			Transform:      transform,
		})
	}
}

func handleBetaOrganizationFederationIssuersArchive(ctx context.Context, cmd *cli.Command) error {
	client := anthropic.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("federation-issuer-id") && len(unusedArgs) > 0 {
		cmd.Set("federation-issuer-id", unusedArgs[0])
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

	params := anthropic.BetaOrganizationFederationIssuerArchiveParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Beta.Organization.Federation.Issuers.Archive(
		ctx,
		cmd.Value("federation-issuer-id").(string),
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
		Title:          "beta:organization:federation:issuers archive",
		Transform:      transform,
	})
}
