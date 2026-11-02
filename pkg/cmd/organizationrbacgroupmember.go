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

var organizationRBACGroupsMembersList = cli.Command{
	Name:    "list",
	Usage:   "List members of an RBAC Group.",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "rbac-group-id",
			Usage:     "ID of the RBAC Group.",
			Required:  true,
			PathParam: "rbac_group_id",
		},
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
	Action:          handleOrganizationRBACGroupsMembersList,
	HideHelpCommand: true,
}

var organizationRBACGroupsMembersAdd = cli.Command{
	Name:    "add",
	Usage:   "Add a User to an RBAC Group. Membership of groups provisioned by an identity\nprovider (source type `\"scim\"`) cannot be modified via the API while an\norganization in the tenant uses SCIM provisioning.",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "rbac-group-id",
			Usage:     "ID of the RBAC Group.",
			Required:  true,
			PathParam: "rbac_group_id",
		},
		&requestflag.Flag[string]{
			Name:     "user-id",
			Usage:    "ID of the User.",
			Required: true,
			BodyPath: "user_id",
		},
	},
	Action:          handleOrganizationRBACGroupsMembersAdd,
	HideHelpCommand: true,
}

var organizationRBACGroupsMembersRemove = cli.Command{
	Name:    "remove",
	Usage:   "Remove a User from an RBAC Group. Membership of groups provisioned by an\nidentity provider (source type `\"scim\"`) cannot be modified via the API while an\norganization in the tenant uses SCIM provisioning.",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "rbac-group-id",
			Usage:     "ID of the RBAC Group.",
			Required:  true,
			PathParam: "rbac_group_id",
		},
		&requestflag.Flag[string]{
			Name:      "user-id",
			Usage:     "ID of the User.",
			Required:  true,
			PathParam: "user_id",
		},
	},
	Action:          handleOrganizationRBACGroupsMembersRemove,
	HideHelpCommand: true,
}

func handleOrganizationRBACGroupsMembersList(ctx context.Context, cmd *cli.Command) error {
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

	params := anthropic.OrganizationRBACGroupMemberListParams{}

	format := "explore"
	explicitFormat := cmd.Root().IsSet("format")
	if explicitFormat {
		format = cmd.Root().String("format")
	}
	transform := cmd.Root().String("transform")
	if format == "raw" {
		var res []byte
		options = append(options, option.WithResponseBodyInto(&res))
		_, err = client.Organization.RBACGroups.Members.List(
			ctx,
			cmd.Value("rbac-group-id").(string),
			params,
			options...,
		)
		if err != nil {
			return err
		}
		obj := gjson.ParseBytes(res)
		return ShowJSON(obj, ShowJSONOpts{
			ExplicitFormat: explicitFormat,
			Format:         format,
			RawOutput:      cmd.Root().Bool("raw-output"),
			Title:          "organization:rbac-groups:members list",
			Transform:      transform,
		})
	} else {
		iter := client.Organization.RBACGroups.Members.ListAutoPaging(
			ctx,
			cmd.Value("rbac-group-id").(string),
			params,
			options...,
		)
		maxItems := int64(-1)
		if cmd.IsSet("max-items") {
			maxItems = cmd.Value("max-items").(int64)
		}
		return ShowJSONIterator(iter, maxItems, ShowJSONOpts{
			ExplicitFormat: explicitFormat,
			Format:         format,
			RawOutput:      cmd.Root().Bool("raw-output"),
			Title:          "organization:rbac-groups:members list",
			Transform:      transform,
		})
	}
}

func handleOrganizationRBACGroupsMembersAdd(ctx context.Context, cmd *cli.Command) error {
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

	params := anthropic.OrganizationRBACGroupMemberAddParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Organization.RBACGroups.Members.Add(
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
		Title:          "organization:rbac-groups:members add",
		Transform:      transform,
	})
}

func handleOrganizationRBACGroupsMembersRemove(ctx context.Context, cmd *cli.Command) error {
	client := anthropic.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("user-id") && len(unusedArgs) > 0 {
		cmd.Set("user-id", unusedArgs[0])
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

	params := anthropic.OrganizationRBACGroupMemberRemoveParams{
		RBACGroupID: cmd.Value("rbac-group-id").(string),
	}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Organization.RBACGroups.Members.Remove(
		ctx,
		cmd.Value("user-id").(string),
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
		Title:          "organization:rbac-groups:members remove",
		Transform:      transform,
	})
}
