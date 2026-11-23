package cmd

import (
	"context"
	"fmt"

	"github.com/anthropics/anthropic-cli/v2/internal/apiquery"
	"github.com/anthropics/anthropic-cli/v2/internal/requestflag"
	"github.com/anthropics/anthropic-sdk-go/v2"
	"github.com/anthropics/anthropic-sdk-go/v2/option"
	"github.com/tidwall/gjson"
	"github.com/urfave/cli/v3"
)

var betaOrganizationExternalStorageCreate = requestflag.WithInnerFlags(cli.Command{
	Name:    "create",
	Usage:   "Register a bucket in the organization's own cloud account. It is saved as\n`pending_configuration` and verified when activated.",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[map[string]any]{
			Name:     "provider-config",
			Required: true,
			BodyPath: "provider_config",
		},
		&requestflag.Flag[*string]{
			Name:     "display-name",
			Usage:    "Human-friendly display name.",
			BodyPath: "display_name",
		},
		&requestflag.Flag[*string]{
			Name:     "user-profile-id",
			Usage:    "The user profile to bind the storage to. It can't be changed later. The storage can be attached only to Workspaces with this profile, and the bucket's setup must name the same ID.",
			BodyPath: "user_profile_id",
		},
		&requestflag.Flag[[]string]{
			Name:       "beta",
			Usage:      "This endpoint is in beta: requests must send `efs-2026-10-01` in this header.",
			HeaderPath: "anthropic-beta",
		},
	},
	Action:          handleBetaOrganizationExternalStorageCreate,
	HideHelpCommand: true,
}, map[string][]requestflag.HasOuterFlag{
	"provider-config": {
		&requestflag.InnerFlag[string]{
			Name:       "provider-config.bucket-name",
			Usage:      "The bucket's name, without `s3://`: 3 to 63 lower-case letters, digits and hyphens, starting and ending with a letter or digit.",
			InnerField: "bucket_name",
		},
		&requestflag.InnerFlag[string]{
			Name:       "provider-config.region",
			Usage:      "The bucket's AWS region. Activating the storage fails if Anthropic can't use buckets in that region.",
			InnerField: "region",
		},
		&requestflag.InnerFlag[string]{
			Name:       "provider-config.type",
			Usage:      `Allowed values: "aws", "gcp".`,
			InnerField: "type",
		},
		&requestflag.InnerFlag[string]{
			Name:       "provider-config.account-id",
			Usage:      "12-digit ID of the AWS account that owns the bucket.",
			InnerField: "account_id",
		},
		&requestflag.InnerFlag[string]{
			Name:       "provider-config.project-number",
			Usage:      "Project number, not project ID, of the Google Cloud project that holds the bucket.",
			InnerField: "project_number",
		},
	},
})

var betaOrganizationExternalStorageRetrieve = cli.Command{
	Name:    "retrieve",
	Usage:   "Get one External Storage in the organization.",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:        "external-storage-id",
			Usage:       "The External Storage ID, which starts with `estrg_`. Through Claude Platform on AWS, storage whose `provider_config.type` is `aws` has its bucket ARN as its ID instead: `arn:<partition>:s3:::<bucket_name>`, where `<partition>` is the partition of the bucket's region.",
			Required:    true,
			PathParam:   "external_storage_id",
			DataAliases: []string{"id"},
		},
		&requestflag.Flag[[]string]{
			Name:       "beta",
			Usage:      "This endpoint is in beta: requests must send `efs-2026-10-01` in this header.",
			HeaderPath: "anthropic-beta",
		},
	},
	Action:          handleBetaOrganizationExternalStorageRetrieve,
	HideHelpCommand: true,
}

var betaOrganizationExternalStorageUpdate = cli.Command{
	Name:    "update",
	Usage:   "Rename External Storage or clear its name, in any status.",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:        "external-storage-id",
			Usage:       "The External Storage ID, which starts with `estrg_`. Through Claude Platform on AWS, storage whose `provider_config.type` is `aws` has its bucket ARN as its ID instead: `arn:<partition>:s3:::<bucket_name>`, where `<partition>` is the partition of the bucket's region.",
			Required:    true,
			PathParam:   "external_storage_id",
			DataAliases: []string{"id"},
		},
		&requestflag.Flag[*string]{
			Name:     "display-name",
			Usage:    "Human-friendly display name; null clears it.",
			Required: true,
			BodyPath: "display_name",
		},
		&requestflag.Flag[[]string]{
			Name:       "beta",
			Usage:      "This endpoint is in beta: requests must send `efs-2026-10-01` in this header.",
			HeaderPath: "anthropic-beta",
		},
	},
	Action:          handleBetaOrganizationExternalStorageUpdate,
	HideHelpCommand: true,
}

var betaOrganizationExternalStorageList = cli.Command{
	Name:    "list",
	Usage:   "List the organization's External Storage, newest first.",
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
		&requestflag.Flag[[]string]{
			Name:      "user-profile-id",
			Usage:     "Returns only External Storage bound to one of these user profile IDs (`uprof_...`). An empty entry matches storage bound to no user profile. At most 10 entries, counting repeats.",
			QueryPath: "user_profile_ids",
		},
		&requestflag.Flag[[]string]{
			Name:       "beta",
			Usage:      "This endpoint is in beta: requests must send `efs-2026-10-01` in this header.",
			HeaderPath: "anthropic-beta",
		},
		&requestflag.Flag[int64]{
			Name:  "max-items",
			Usage: "The maximum number of items to return (use -1 for unlimited).",
		},
	},
	Action:          handleBetaOrganizationExternalStorageList,
	HideHelpCommand: true,
}

var betaOrganizationExternalStorageDelete = cli.Command{
	Name:    "delete",
	Usage:   "Delete External Storage. Storage that is pending configuration, available, or\nactive but never attached is deleted at once, unless an Enterprise Frontier\nSafeguards workspace policy lists it. Other storage gets a 400 that says what to\ndo next. On `draining` storage that 400 can mean this delete started a removal\nthat can't be undone. To check, get the storage: if `status.removal_started` is\ntrue, delete it again later to finish.",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:        "external-storage-id",
			Usage:       "The External Storage ID, which starts with `estrg_`. Through Claude Platform on AWS, storage whose `provider_config.type` is `aws` has its bucket ARN as its ID instead: `arn:<partition>:s3:::<bucket_name>`, where `<partition>` is the partition of the bucket's region.",
			Required:    true,
			PathParam:   "external_storage_id",
			DataAliases: []string{"id"},
		},
		&requestflag.Flag[[]string]{
			Name:       "beta",
			Usage:      "This endpoint is in beta: requests must send `efs-2026-10-01` in this header.",
			HeaderPath: "anthropic-beta",
		},
	},
	Action:          handleBetaOrganizationExternalStorageDelete,
	HideHelpCommand: true,
}

var betaOrganizationExternalStorageActivate = cli.Command{
	Name:    "activate",
	Usage:   "Verify pending External Storage and make it available. Storage that is already\navailable, active or draining comes back unchanged, without any check.",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:        "external-storage-id",
			Usage:       "The External Storage ID, which starts with `estrg_`. Through Claude Platform on AWS, storage whose `provider_config.type` is `aws` has its bucket ARN as its ID instead: `arn:<partition>:s3:::<bucket_name>`, where `<partition>` is the partition of the bucket's region.",
			Required:    true,
			PathParam:   "external_storage_id",
			DataAliases: []string{"id"},
		},
		&requestflag.Flag[[]string]{
			Name:       "beta",
			Usage:      "This endpoint is in beta: requests must send `efs-2026-10-01` in this header.",
			HeaderPath: "anthropic-beta",
		},
	},
	Action:          handleBetaOrganizationExternalStorageActivate,
	HideHelpCommand: true,
}

func handleBetaOrganizationExternalStorageCreate(ctx context.Context, cmd *cli.Command) error {
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

	params := anthropic.BetaOrganizationExternalStorageNewParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Beta.Organization.ExternalStorage.New(ctx, params, options...)
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
		Title:          "beta:organization:external-storage create",
		Transform:      transform,
	})
}

func handleBetaOrganizationExternalStorageRetrieve(ctx context.Context, cmd *cli.Command) error {
	client := anthropic.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("external-storage-id") && len(unusedArgs) > 0 {
		cmd.Set("external-storage-id", unusedArgs[0])
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

	params := anthropic.BetaOrganizationExternalStorageGetParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Beta.Organization.ExternalStorage.Get(
		ctx,
		cmd.Value("external-storage-id").(string),
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
		Title:          "beta:organization:external-storage retrieve",
		Transform:      transform,
	})
}

func handleBetaOrganizationExternalStorageUpdate(ctx context.Context, cmd *cli.Command) error {
	client := anthropic.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("external-storage-id") && len(unusedArgs) > 0 {
		cmd.Set("external-storage-id", unusedArgs[0])
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

	params := anthropic.BetaOrganizationExternalStorageUpdateParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Beta.Organization.ExternalStorage.Update(
		ctx,
		cmd.Value("external-storage-id").(string),
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
		Title:          "beta:organization:external-storage update",
		Transform:      transform,
	})
}

func handleBetaOrganizationExternalStorageList(ctx context.Context, cmd *cli.Command) error {
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

	params := anthropic.BetaOrganizationExternalStorageListParams{}

	format := "explore"
	explicitFormat := cmd.Root().IsSet("format")
	if explicitFormat {
		format = cmd.Root().String("format")
	}
	transform := cmd.Root().String("transform")
	if format == "raw" {
		var res []byte
		options = append(options, option.WithResponseBodyInto(&res))
		_, err = client.Beta.Organization.ExternalStorage.List(ctx, params, options...)
		if err != nil {
			return err
		}
		obj := gjson.ParseBytes(res)
		return ShowJSON(obj, ShowJSONOpts{
			ExplicitFormat: explicitFormat,
			Format:         format,
			RawOutput:      cmd.Root().Bool("raw-output"),
			Title:          "beta:organization:external-storage list",
			Transform:      transform,
		})
	} else {
		iter := client.Beta.Organization.ExternalStorage.ListAutoPaging(ctx, params, options...)
		maxItems := int64(-1)
		if cmd.IsSet("max-items") {
			maxItems = cmd.Value("max-items").(int64)
		}
		return ShowJSONIterator(iter, maxItems, ShowJSONOpts{
			ExplicitFormat: explicitFormat,
			Format:         format,
			RawOutput:      cmd.Root().Bool("raw-output"),
			Title:          "beta:organization:external-storage list",
			Transform:      transform,
		})
	}
}

func handleBetaOrganizationExternalStorageDelete(ctx context.Context, cmd *cli.Command) error {
	client := anthropic.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("external-storage-id") && len(unusedArgs) > 0 {
		cmd.Set("external-storage-id", unusedArgs[0])
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

	params := anthropic.BetaOrganizationExternalStorageDeleteParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Beta.Organization.ExternalStorage.Delete(
		ctx,
		cmd.Value("external-storage-id").(string),
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
		Title:          "beta:organization:external-storage delete",
		Transform:      transform,
	})
}

func handleBetaOrganizationExternalStorageActivate(ctx context.Context, cmd *cli.Command) error {
	client := anthropic.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("external-storage-id") && len(unusedArgs) > 0 {
		cmd.Set("external-storage-id", unusedArgs[0])
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

	params := anthropic.BetaOrganizationExternalStorageActivateParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Beta.Organization.ExternalStorage.Activate(
		ctx,
		cmd.Value("external-storage-id").(string),
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
		Title:          "beta:organization:external-storage activate",
		Transform:      transform,
	})
}
