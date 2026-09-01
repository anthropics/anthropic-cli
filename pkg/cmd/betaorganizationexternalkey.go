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

var betaOrganizationExternalKeysCreate = requestflag.WithInnerFlags(cli.Command{
	Name:    "create",
	Usage:   "Create an external key config owned by the caller's organization.",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[map[string]any]{
			Name:     "provider-config",
			Usage:    "KMS provider identity and auth coordinates.",
			Required: true,
			BodyPath: "provider_config",
		},
		&requestflag.Flag[*string]{
			Name:     "display-name",
			Usage:    "Human-friendly display name.",
			BodyPath: "display_name",
		},
		&requestflag.Flag[string]{
			Name:     "geo",
			Usage:    "Data residency geo. Only `us` is supported.",
			BodyPath: "geo",
		},
	},
	Action:          handleBetaOrganizationExternalKeysCreate,
	HideHelpCommand: true,
}, map[string][]requestflag.HasOuterFlag{
	"provider-config": {
		&requestflag.InnerFlag[string]{
			Name:       "provider-config.type",
			Usage:      `Allowed values: "aws", "gcp", "azure".`,
			InnerField: "type",
		},
		&requestflag.InnerFlag[*string]{
			Name:       "provider-config.client-id",
			Usage:      "Azure AD application (client) ID. Omit to use Anthropic's multitenant app. Provide only if using a single-tenant app registration in the customer's directory.",
			InnerField: "client_id",
		},
		&requestflag.InnerFlag[string]{
			Name:       "provider-config.key-name",
			Usage:      "Full resource name of the Cloud KMS key.",
			InnerField: "key_name",
		},
		&requestflag.InnerFlag[string]{
			Name:       "provider-config.kms-arn",
			Usage:      "Full ARN of the AWS KMS key. On Claude Platform on AWS the key must be a single-Region key in your organization's own AWS account; cross-account keys, multi-Region keys, and alias ARNs are rejected.",
			InnerField: "kms_arn",
		},
		&requestflag.InnerFlag[*string]{
			Name:       "provider-config.region",
			Usage:      "AWS region. Derived from `kms_arn` if omitted.",
			InnerField: "region",
		},
		&requestflag.InnerFlag[*string]{
			Name:       "provider-config.role-arn",
			Usage:      "IAM role ARN. Deprecated — Anthropic reaches the KMS key through its own intermediate role (or, on Claude Platform on AWS, with credentials AWS issues for the Workspace); this field is ignored.",
			InnerField: "role_arn",
		},
		&requestflag.InnerFlag[string]{
			Name:       "provider-config.tenant-id",
			Usage:      "Azure AD tenant ID.",
			InnerField: "tenant_id",
		},
		&requestflag.InnerFlag[string]{
			Name:       "provider-config.vault-uri",
			Usage:      "Key Vault data-plane URI — `https://{vault-name}.vault.azure.net` or `https://{hsm-name}.managedhsm.azure.net`.",
			InnerField: "vault_uri",
		},
	},
})

var betaOrganizationExternalKeysRetrieve = cli.Command{
	Name:    "retrieve",
	Usage:   "Retrieve a single external key config in the caller's organization by ID.",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "external-key-id",
			Usage:     "ID of the External Key.",
			Required:  true,
			PathParam: "external_key_id",
		},
	},
	Action:          handleBetaOrganizationExternalKeysRetrieve,
	HideHelpCommand: true,
}

var betaOrganizationExternalKeysUpdate = requestflag.WithInnerFlags(cli.Command{
	Name:    "update",
	Usage:   "Partially update an external key config. Omitted fields are left unchanged.",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "external-key-id",
			Usage:     "ID of the External Key.",
			Required:  true,
			PathParam: "external_key_id",
		},
		&requestflag.Flag[*string]{
			Name:     "display-name",
			Usage:    "Human-friendly display name.",
			BodyPath: "display_name",
		},
		&requestflag.Flag[*string]{
			Name:     "geo",
			Usage:    "Data residency geo. Only `us` is supported.",
			BodyPath: "geo",
		},
		&requestflag.Flag[map[string]any]{
			Name:     "provider-config",
			Usage:    "KMS provider identity and auth coordinates.",
			BodyPath: "provider_config",
		},
	},
	Action:          handleBetaOrganizationExternalKeysUpdate,
	HideHelpCommand: true,
}, map[string][]requestflag.HasOuterFlag{
	"provider-config": {
		&requestflag.InnerFlag[string]{
			Name:       "provider-config.type",
			Usage:      `Allowed values: "aws", "gcp", "azure".`,
			InnerField: "type",
		},
		&requestflag.InnerFlag[*string]{
			Name:       "provider-config.client-id",
			Usage:      "Azure AD application (client) ID. Omit to use Anthropic's multitenant app. Provide only if using a single-tenant app registration in the customer's directory.",
			InnerField: "client_id",
		},
		&requestflag.InnerFlag[string]{
			Name:       "provider-config.key-name",
			Usage:      "Full resource name of the Cloud KMS key.",
			InnerField: "key_name",
		},
		&requestflag.InnerFlag[string]{
			Name:       "provider-config.kms-arn",
			Usage:      "Full ARN of the AWS KMS key. On Claude Platform on AWS the key must be a single-Region key in your organization's own AWS account; cross-account keys, multi-Region keys, and alias ARNs are rejected.",
			InnerField: "kms_arn",
		},
		&requestflag.InnerFlag[*string]{
			Name:       "provider-config.region",
			Usage:      "AWS region. Derived from `kms_arn` if omitted.",
			InnerField: "region",
		},
		&requestflag.InnerFlag[*string]{
			Name:       "provider-config.role-arn",
			Usage:      "IAM role ARN. Deprecated — Anthropic reaches the KMS key through its own intermediate role (or, on Claude Platform on AWS, with credentials AWS issues for the Workspace); this field is ignored.",
			InnerField: "role_arn",
		},
		&requestflag.InnerFlag[string]{
			Name:       "provider-config.tenant-id",
			Usage:      "Azure AD tenant ID.",
			InnerField: "tenant_id",
		},
		&requestflag.InnerFlag[string]{
			Name:       "provider-config.vault-uri",
			Usage:      "Key Vault data-plane URI — `https://{vault-name}.vault.azure.net` or `https://{hsm-name}.managedhsm.azure.net`.",
			InnerField: "vault_uri",
		},
	},
})

var betaOrganizationExternalKeysList = cli.Command{
	Name:    "list",
	Usage:   "List external key configs in the caller's organization.",
	Suggest: true,
	Flags: []cli.Flag{
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
		&requestflag.Flag[int64]{
			Name:  "max-items",
			Usage: "The maximum number of items to return (use -1 for unlimited).",
		},
	},
	Action:          handleBetaOrganizationExternalKeysList,
	HideHelpCommand: true,
}

var betaOrganizationExternalKeysDelete = cli.Command{
	Name:    "delete",
	Usage:   "Delete an external key config.",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "external-key-id",
			Usage:     "ID of the External Key.",
			Required:  true,
			PathParam: "external_key_id",
		},
	},
	Action:          handleBetaOrganizationExternalKeysDelete,
	HideHelpCommand: true,
}

var betaOrganizationExternalKeysValidate = cli.Command{
	Name:    "validate",
	Usage:   "Validate an external key config against the customer's KMS.",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "external-key-id",
			Usage:     "ID of the External Key.",
			Required:  true,
			PathParam: "external_key_id",
		},
	},
	Action:          handleBetaOrganizationExternalKeysValidate,
	HideHelpCommand: true,
}

func handleBetaOrganizationExternalKeysCreate(ctx context.Context, cmd *cli.Command) error {
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

	params := anthropic.BetaOrganizationExternalKeyNewParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Beta.Organization.ExternalKeys.New(ctx, params, options...)
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
		Title:          "beta:organization:external-keys create",
		Transform:      transform,
	})
}

func handleBetaOrganizationExternalKeysRetrieve(ctx context.Context, cmd *cli.Command) error {
	client := anthropic.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("external-key-id") && len(unusedArgs) > 0 {
		cmd.Set("external-key-id", unusedArgs[0])
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

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Beta.Organization.ExternalKeys.Get(ctx, cmd.Value("external-key-id").(string), options...)
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
		Title:          "beta:organization:external-keys retrieve",
		Transform:      transform,
	})
}

func handleBetaOrganizationExternalKeysUpdate(ctx context.Context, cmd *cli.Command) error {
	client := anthropic.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("external-key-id") && len(unusedArgs) > 0 {
		cmd.Set("external-key-id", unusedArgs[0])
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

	params := anthropic.BetaOrganizationExternalKeyUpdateParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Beta.Organization.ExternalKeys.Update(
		ctx,
		cmd.Value("external-key-id").(string),
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
		Title:          "beta:organization:external-keys update",
		Transform:      transform,
	})
}

func handleBetaOrganizationExternalKeysList(ctx context.Context, cmd *cli.Command) error {
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

	params := anthropic.BetaOrganizationExternalKeyListParams{}

	format := "explore"
	explicitFormat := cmd.Root().IsSet("format")
	if explicitFormat {
		format = cmd.Root().String("format")
	}
	transform := cmd.Root().String("transform")
	if format == "raw" {
		var res []byte
		options = append(options, option.WithResponseBodyInto(&res))
		_, err = client.Beta.Organization.ExternalKeys.List(ctx, params, options...)
		if err != nil {
			return err
		}
		obj := gjson.ParseBytes(res)
		return ShowJSON(obj, ShowJSONOpts{
			ExplicitFormat: explicitFormat,
			Format:         format,
			RawOutput:      cmd.Root().Bool("raw-output"),
			Title:          "beta:organization:external-keys list",
			Transform:      transform,
		})
	} else {
		iter := client.Beta.Organization.ExternalKeys.ListAutoPaging(ctx, params, options...)
		maxItems := int64(-1)
		if cmd.IsSet("max-items") {
			maxItems = cmd.Value("max-items").(int64)
		}
		return ShowJSONIterator(iter, maxItems, ShowJSONOpts{
			ExplicitFormat: explicitFormat,
			Format:         format,
			RawOutput:      cmd.Root().Bool("raw-output"),
			Title:          "beta:organization:external-keys list",
			Transform:      transform,
		})
	}
}

func handleBetaOrganizationExternalKeysDelete(ctx context.Context, cmd *cli.Command) error {
	client := anthropic.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("external-key-id") && len(unusedArgs) > 0 {
		cmd.Set("external-key-id", unusedArgs[0])
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

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Beta.Organization.ExternalKeys.Delete(ctx, cmd.Value("external-key-id").(string), options...)
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
		Title:          "beta:organization:external-keys delete",
		Transform:      transform,
	})
}

func handleBetaOrganizationExternalKeysValidate(ctx context.Context, cmd *cli.Command) error {
	client := anthropic.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("external-key-id") && len(unusedArgs) > 0 {
		cmd.Set("external-key-id", unusedArgs[0])
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

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Beta.Organization.ExternalKeys.Validate(ctx, cmd.Value("external-key-id").(string), options...)
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
		Title:          "beta:organization:external-keys validate",
		Transform:      transform,
	})
}
