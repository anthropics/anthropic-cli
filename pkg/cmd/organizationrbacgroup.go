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

var organizationRBACGroupsCreate = cli.Command{
	Name:    "create",
	Usage:   "Create an RBAC Group in the Claude Enterprise tenant. Groups created via the API\nhave source type `\"direct\"`.",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:     "name",
			Usage:    "Name of the RBAC Group. Not uniqueness-enforced.",
			Required: true,
			BodyPath: "name",
		},
	},
	Action:          handleOrganizationRBACGroupsCreate,
	HideHelpCommand: true,
}

var organizationRBACGroupsRetrieve = cli.Command{
	Name:    "retrieve",
	Usage:   "Retrieve an RBAC Group by ID.",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:        "rbac-group-id",
			Usage:       "ID of the RBAC Group.",
			Required:    true,
			PathParam:   "rbac_group_id",
			DataAliases: []string{"id"},
		},
	},
	Action:          handleOrganizationRBACGroupsRetrieve,
	HideHelpCommand: true,
}

var organizationRBACGroupsUpdate = cli.Command{
	Name:    "update",
	Usage:   "Update an RBAC Group's name. Groups provisioned by an identity provider (source\ntype `\"scim\"`) cannot be modified via the API while an organization in the\ntenant uses SCIM provisioning.",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:        "rbac-group-id",
			Usage:       "ID of the RBAC Group.",
			Required:    true,
			PathParam:   "rbac_group_id",
			DataAliases: []string{"id"},
		},
		&requestflag.Flag[*string]{
			Name:     "name",
			Usage:    "Name of the RBAC Group. Not uniqueness-enforced.",
			BodyPath: "name",
		},
	},
	Action:          handleOrganizationRBACGroupsUpdate,
	HideHelpCommand: true,
}

var organizationRBACGroupsList = cli.Command{
	Name:    "list",
	Usage:   "List RBAC Groups in the Claude Enterprise tenant.",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[int64]{
			Name:      "limit",
			Usage:     "Number of items to return per page.\n\nDefaults to `20`. Ranges from `1` to `1000`.",
			Default:   20,
			QueryPath: "limit",
		},
		&requestflag.Flag[string]{
			Name:      "page",
			Usage:     "Optionally set to the `next_page` token from the previous response.",
			QueryPath: "page",
		},
		&requestflag.Flag[int64]{
			Name:  "max-items",
			Usage: "The maximum number of items to return (use -1 for unlimited).",
		},
	},
	Action:          handleOrganizationRBACGroupsList,
	HideHelpCommand: true,
}

var organizationRBACGroupsDelete = cli.Command{
	Name:    "delete",
	Usage:   "Delete an RBAC Group. Groups provisioned by an identity provider (source type\n`\"scim\"`) cannot be deleted via the API while an organization in the tenant uses\nSCIM provisioning.",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:        "rbac-group-id",
			Usage:       "ID of the RBAC Group.",
			Required:    true,
			PathParam:   "rbac_group_id",
			DataAliases: []string{"id"},
		},
	},
	Action:          handleOrganizationRBACGroupsDelete,
	HideHelpCommand: true,
}

func handleOrganizationRBACGroupsCreate(ctx context.Context, cmd *cli.Command) error {
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

	params := anthropic.OrganizationRBACGroupNewParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Organization.RBACGroups.New(ctx, params, options...)
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
		Title:          "organization:rbac-groups create",
		Transform:      transform,
	})
}

func handleOrganizationRBACGroupsRetrieve(ctx context.Context, cmd *cli.Command) error {
	client := anthropic.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("rbac-group-id") && len(unusedArgs) > 0 {
		cmd.Set("rbac-group-id", unusedArgs[0])
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
	_, err = client.Organization.RBACGroups.Get(ctx, cmd.Value("rbac-group-id").(string), options...)
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
		Title:          "organization:rbac-groups retrieve",
		Transform:      transform,
	})
}

func handleOrganizationRBACGroupsUpdate(ctx context.Context, cmd *cli.Command) error {
	client := anthropic.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("rbac-group-id") && len(unusedArgs) > 0 {
		cmd.Set("rbac-group-id", unusedArgs[0])
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

	params := anthropic.OrganizationRBACGroupUpdateParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Organization.RBACGroups.Update(
		ctx,
		cmd.Value("rbac-group-id").(string),
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
		Title:          "organization:rbac-groups update",
		Transform:      transform,
	})
}

func handleOrganizationRBACGroupsList(ctx context.Context, cmd *cli.Command) error {
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

	params := anthropic.OrganizationRBACGroupListParams{}

	format := "explore"
	explicitFormat := cmd.Root().IsSet("format")
	if explicitFormat {
		format = cmd.Root().String("format")
	}
	transform := cmd.Root().String("transform")
	if format == "raw" {
		var res []byte
		options = append(options, option.WithResponseBodyInto(&res))
		_, err = client.Organization.RBACGroups.List(ctx, params, options...)
		if err != nil {
			return err
		}
		obj := gjson.ParseBytes(res)
		return ShowJSON(obj, ShowJSONOpts{
			ExplicitFormat: explicitFormat,
			Format:         format,
			RawOutput:      cmd.Root().Bool("raw-output"),
			Title:          "organization:rbac-groups list",
			Transform:      transform,
		})
	} else {
		iter := client.Organization.RBACGroups.ListAutoPaging(ctx, params, options...)
		maxItems := int64(-1)
		if cmd.IsSet("max-items") {
			maxItems = cmd.Value("max-items").(int64)
		}
		return ShowJSONIterator(iter, maxItems, ShowJSONOpts{
			ExplicitFormat: explicitFormat,
			Format:         format,
			RawOutput:      cmd.Root().Bool("raw-output"),
			Title:          "organization:rbac-groups list",
			Transform:      transform,
		})
	}
}

func handleOrganizationRBACGroupsDelete(ctx context.Context, cmd *cli.Command) error {
	client := anthropic.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("rbac-group-id") && len(unusedArgs) > 0 {
		cmd.Set("rbac-group-id", unusedArgs[0])
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
	_, err = client.Organization.RBACGroups.Delete(ctx, cmd.Value("rbac-group-id").(string), options...)
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
		Title:          "organization:rbac-groups delete",
		Transform:      transform,
	})
}
