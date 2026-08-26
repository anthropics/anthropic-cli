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

var betaOrganizationWorkspacesCreate = requestflag.WithInnerFlags(cli.Command{
	Name:    "create",
	Usage:   "Create Workspace",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:     "name",
			Usage:    "Name of the Workspace.",
			Required: true,
			BodyPath: "name",
		},
		&requestflag.Flag[map[string]any]{
			Name:     "data-residency",
			BodyPath: "data_residency",
		},
		&requestflag.Flag[*string]{
			Name:     "display-color",
			Usage:    "Hex color code representing the Workspace in the Anthropic Console.",
			BodyPath: "display_color",
		},
		&requestflag.Flag[*string]{
			Name:     "external-key-id",
			Usage:    "ID of the customer-managed encryption key (CMEK) configuration to use for this\nWorkspace. Setting this field requires CMEK to be enabled for your\norganization. When set, data stored for this Workspace is encrypted with the\nreferenced key. Create key configurations with the External Keys API. This\nfield is write-once: once a key is attached to a Workspace it cannot be\ndetached or replaced. To rotate key material, rotate the underlying key on\nyour cloud KMS; the `external_key_id` stays the same.",
			BodyPath: "external_key_id",
		},
		&requestflag.Flag[map[string]any]{
			Name:     "tags",
			Usage:    "User-defined tags as string key-value pairs. Keys may not begin with `anthropic`.",
			BodyPath: "tags",
		},
		&requestflag.Flag[[]string]{
			Name:       "beta",
			Usage:      "Optional header to specify the beta version(s) you want to use.",
			HeaderPath: "anthropic-beta",
		},
	},
	Action:          handleBetaOrganizationWorkspacesCreate,
	HideHelpCommand: true,
}, map[string][]requestflag.HasOuterFlag{
	"data-residency": {
		&requestflag.InnerFlag[any]{
			Name:       "data-residency.allowed-inference-geos",
			Usage:      "Permitted inference geo values. Defaults to 'unrestricted' if omitted, which allows all geos. Use the string 'unrestricted' to allow all geos, or a list of specific geos.",
			InnerField: "allowed_inference_geos",
		},
		&requestflag.InnerFlag[*string]{
			Name:       "data-residency.default-inference-geo",
			Usage:      "Default inference geo applied when requests omit the parameter. Defaults to 'global' if omitted. Must be a member of `allowed_inference_geos` unless `allowed_inference_geos` is `\"unrestricted\"`.",
			InnerField: "default_inference_geo",
		},
		&requestflag.InnerFlag[*string]{
			Name:       "data-residency.workspace-geo",
			Usage:      "Geographic region for workspace data storage. Immutable after creation. Defaults to 'us' if omitted.",
			InnerField: "workspace_geo",
		},
	},
})

var betaOrganizationWorkspacesRetrieve = cli.Command{
	Name:    "retrieve",
	Usage:   "Get Workspace",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "workspace-id",
			Usage:     "ID of the Workspace.",
			Required:  true,
			PathParam: "workspace_id",
		},
	},
	Action:          handleBetaOrganizationWorkspacesRetrieve,
	HideHelpCommand: true,
}

var betaOrganizationWorkspacesUpdate = requestflag.WithInnerFlags(cli.Command{
	Name:    "update",
	Usage:   "Update Workspace",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "workspace-id",
			Required:  true,
			PathParam: "workspace_id",
		},
		&requestflag.Flag[map[string]any]{
			Name:     "data-residency",
			BodyPath: "data_residency",
		},
		&requestflag.Flag[string]{
			Name:     "display-color",
			Usage:    "Hex color code representing the Workspace in the Anthropic Console.",
			BodyPath: "display_color",
		},
		&requestflag.Flag[string]{
			Name:     "external-key-id",
			Usage:    "ID of the customer-managed encryption key (CMEK) configuration to use for this\nWorkspace. Setting this field requires CMEK to be enabled for your\norganization. When set, data stored for this Workspace is encrypted with the\nreferenced key. Create key configurations with the External Keys API. This\nfield is write-once: once a key is attached to a Workspace it cannot be\ndetached or replaced. To rotate key material, rotate the underlying key on\nyour cloud KMS; the `external_key_id` stays the same.",
			BodyPath: "external_key_id",
		},
		&requestflag.Flag[string]{
			Name:     "name",
			Usage:    "Name of the Workspace.",
			BodyPath: "name",
		},
		&requestflag.Flag[map[string]any]{
			Name:     "tags",
			Usage:    "User-defined tags as string key-value pairs. Keys may not begin with `anthropic`.",
			BodyPath: "tags",
		},
	},
	Action:          handleBetaOrganizationWorkspacesUpdate,
	HideHelpCommand: true,
}, map[string][]requestflag.HasOuterFlag{
	"data-residency": {
		&requestflag.InnerFlag[any]{
			Name:       "data-residency.allowed-inference-geos",
			Usage:      "Permitted inference geo values. Use 'unrestricted' to allow all geos, or a list of specific geos.",
			InnerField: "allowed_inference_geos",
		},
		&requestflag.InnerFlag[*string]{
			Name:       "data-residency.default-inference-geo",
			Usage:      "Default inference geo applied when requests omit the parameter. Must be a member of `allowed_inference_geos` unless `allowed_inference_geos` is `\"unrestricted\"`.",
			InnerField: "default_inference_geo",
		},
	},
})

var betaOrganizationWorkspacesList = cli.Command{
	Name:    "list",
	Usage:   "List Workspaces",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "after-id",
			Usage:     "ID of the object to use as a cursor for pagination. When provided, returns the page of results immediately after this object.",
			QueryPath: "after_id",
		},
		&requestflag.Flag[string]{
			Name:      "before-id",
			Usage:     "ID of the object to use as a cursor for pagination. When provided, returns the page of results immediately before this object.",
			QueryPath: "before_id",
		},
		&requestflag.Flag[bool]{
			Name:      "include-archived",
			Usage:     "Whether to include Workspaces that have been archived in the response",
			Default:   false,
			QueryPath: "include_archived",
		},
		&requestflag.Flag[int64]{
			Name:      "limit",
			Usage:     "Number of items to return per page.\n\nDefaults to `20`. Ranges from `1` to `1000`.",
			Default:   20,
			QueryPath: "limit",
		},
		&requestflag.Flag[int64]{
			Name:  "max-items",
			Usage: "The maximum number of items to return (use -1 for unlimited).",
		},
	},
	Action:          handleBetaOrganizationWorkspacesList,
	HideHelpCommand: true,
}

var betaOrganizationWorkspacesArchive = cli.Command{
	Name:    "archive",
	Usage:   "Archive Workspace",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "workspace-id",
			Required:  true,
			PathParam: "workspace_id",
		},
	},
	Action:          handleBetaOrganizationWorkspacesArchive,
	HideHelpCommand: true,
}

func handleBetaOrganizationWorkspacesCreate(ctx context.Context, cmd *cli.Command) error {
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

	params := anthropic.BetaOrganizationWorkspaceNewParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Beta.Organization.Workspaces.New(ctx, params, options...)
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
		Title:          "beta:organization:workspaces create",
		Transform:      transform,
	})
}

func handleBetaOrganizationWorkspacesRetrieve(ctx context.Context, cmd *cli.Command) error {
	client := anthropic.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("workspace-id") && len(unusedArgs) > 0 {
		cmd.Set("workspace-id", unusedArgs[0])
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
	_, err = client.Beta.Organization.Workspaces.Get(ctx, cmd.Value("workspace-id").(string), options...)
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
		Title:          "beta:organization:workspaces retrieve",
		Transform:      transform,
	})
}

func handleBetaOrganizationWorkspacesUpdate(ctx context.Context, cmd *cli.Command) error {
	client := anthropic.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("workspace-id") && len(unusedArgs) > 0 {
		cmd.Set("workspace-id", unusedArgs[0])
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

	params := anthropic.BetaOrganizationWorkspaceUpdateParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Beta.Organization.Workspaces.Update(
		ctx,
		cmd.Value("workspace-id").(string),
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
		Title:          "beta:organization:workspaces update",
		Transform:      transform,
	})
}

func handleBetaOrganizationWorkspacesList(ctx context.Context, cmd *cli.Command) error {
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

	params := anthropic.BetaOrganizationWorkspaceListParams{}

	format := "explore"
	explicitFormat := cmd.Root().IsSet("format")
	if explicitFormat {
		format = cmd.Root().String("format")
	}
	transform := cmd.Root().String("transform")
	if format == "raw" {
		var res []byte
		options = append(options, option.WithResponseBodyInto(&res))
		_, err = client.Beta.Organization.Workspaces.List(ctx, params, options...)
		if err != nil {
			return err
		}
		obj := gjson.ParseBytes(res)
		return ShowJSON(obj, ShowJSONOpts{
			ExplicitFormat: explicitFormat,
			Format:         format,
			RawOutput:      cmd.Root().Bool("raw-output"),
			Title:          "beta:organization:workspaces list",
			Transform:      transform,
		})
	} else {
		iter := client.Beta.Organization.Workspaces.ListAutoPaging(ctx, params, options...)
		maxItems := int64(-1)
		if cmd.IsSet("max-items") {
			maxItems = cmd.Value("max-items").(int64)
		} else if cmd.IsSet("limit") {
			// notably, `limit` is still sent, so results are truncated server side, but this will stop further auto-iteration
			maxItems = cmd.Value("limit").(int64)
		}
		return ShowJSONIterator(iter, maxItems, ShowJSONOpts{
			ExplicitFormat: explicitFormat,
			Format:         format,
			RawOutput:      cmd.Root().Bool("raw-output"),
			Title:          "beta:organization:workspaces list",
			Transform:      transform,
		})
	}
}

func handleBetaOrganizationWorkspacesArchive(ctx context.Context, cmd *cli.Command) error {
	client := anthropic.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("workspace-id") && len(unusedArgs) > 0 {
		cmd.Set("workspace-id", unusedArgs[0])
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
	_, err = client.Beta.Organization.Workspaces.Archive(ctx, cmd.Value("workspace-id").(string), options...)
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
		Title:          "beta:organization:workspaces archive",
		Transform:      transform,
	})
}
