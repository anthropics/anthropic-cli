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

var betaOrganizationWorkspacesMembersRetrieve = cli.Command{
	Name:    "retrieve",
	Usage:   "Get Workspace Member",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "workspace-id",
			Usage:     "ID of the Workspace.",
			Required:  true,
			PathParam: "workspace_id",
		},
		&requestflag.Flag[string]{
			Name:      "user-id",
			Usage:     "ID of the User.",
			Required:  true,
			PathParam: "user_id",
		},
	},
	Action:          handleBetaOrganizationWorkspacesMembersRetrieve,
	HideHelpCommand: true,
}

var betaOrganizationWorkspacesMembersUpdate = cli.Command{
	Name:    "update",
	Usage:   "Update Workspace Member",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "workspace-id",
			Usage:     "ID of the Workspace.",
			Required:  true,
			PathParam: "workspace_id",
		},
		&requestflag.Flag[string]{
			Name:      "user-id",
			Usage:     "ID of the User.",
			Required:  true,
			PathParam: "user_id",
		},
		&requestflag.Flag[string]{
			Name:     "workspace-role",
			Usage:    `Allowed values: "workspace_admin", "workspace_billing", "workspace_developer", "workspace_restricted_developer", "workspace_user".`,
			Required: true,
			BodyPath: "workspace_role",
		},
	},
	Action:          handleBetaOrganizationWorkspacesMembersUpdate,
	HideHelpCommand: true,
}

var betaOrganizationWorkspacesMembersList = cli.Command{
	Name:    "list",
	Usage:   "List Workspace Members",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "workspace-id",
			Usage:     "ID of the Workspace.",
			Required:  true,
			PathParam: "workspace_id",
		},
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
	Action:          handleBetaOrganizationWorkspacesMembersList,
	HideHelpCommand: true,
}

var betaOrganizationWorkspacesMembersAdd = cli.Command{
	Name:    "add",
	Usage:   "Create Workspace Member",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "workspace-id",
			Usage:     "ID of the Workspace.",
			Required:  true,
			PathParam: "workspace_id",
		},
		&requestflag.Flag[string]{
			Name:     "user-id",
			Usage:    "ID of the User.",
			Required: true,
			BodyPath: "user_id",
		},
		&requestflag.Flag[string]{
			Name:     "workspace-role",
			Usage:    `Allowed values: "workspace_admin", "workspace_developer", "workspace_restricted_developer", "workspace_user".`,
			Required: true,
			BodyPath: "workspace_role",
		},
	},
	Action:          handleBetaOrganizationWorkspacesMembersAdd,
	HideHelpCommand: true,
}

var betaOrganizationWorkspacesMembersRemove = cli.Command{
	Name:    "remove",
	Usage:   "Delete Workspace Member",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "workspace-id",
			Usage:     "ID of the Workspace.",
			Required:  true,
			PathParam: "workspace_id",
		},
		&requestflag.Flag[string]{
			Name:      "user-id",
			Usage:     "ID of the User.",
			Required:  true,
			PathParam: "user_id",
		},
	},
	Action:          handleBetaOrganizationWorkspacesMembersRemove,
	HideHelpCommand: true,
}

func handleBetaOrganizationWorkspacesMembersRetrieve(ctx context.Context, cmd *cli.Command) error {
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

	params := anthropic.BetaOrganizationWorkspaceMemberGetParams{
		WorkspaceID: cmd.Value("workspace-id").(string),
	}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Beta.Organization.Workspaces.Members.Get(
		ctx,
		cmd.Value("user-id").(string),
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
		Title:          "beta:organization:workspaces:members retrieve",
		Transform:      transform,
	})
}

func handleBetaOrganizationWorkspacesMembersUpdate(ctx context.Context, cmd *cli.Command) error {
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
		ApplicationJSON,
		false,
	)
	if err != nil {
		return err
	}

	params := anthropic.BetaOrganizationWorkspaceMemberUpdateParams{
		WorkspaceID: cmd.Value("workspace-id").(string),
	}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Beta.Organization.Workspaces.Members.Update(
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
		Title:          "beta:organization:workspaces:members update",
		Transform:      transform,
	})
}

func handleBetaOrganizationWorkspacesMembersList(ctx context.Context, cmd *cli.Command) error {
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

	params := anthropic.BetaOrganizationWorkspaceMemberListParams{}

	format := "explore"
	explicitFormat := cmd.Root().IsSet("format")
	if explicitFormat {
		format = cmd.Root().String("format")
	}
	transform := cmd.Root().String("transform")
	if format == "raw" {
		var res []byte
		options = append(options, option.WithResponseBodyInto(&res))
		_, err = client.Beta.Organization.Workspaces.Members.List(
			ctx,
			cmd.Value("workspace-id").(string),
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
			Title:          "beta:organization:workspaces:members list",
			Transform:      transform,
		})
	} else {
		iter := client.Beta.Organization.Workspaces.Members.ListAutoPaging(
			ctx,
			cmd.Value("workspace-id").(string),
			params,
			options...,
		)
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
			Title:          "beta:organization:workspaces:members list",
			Transform:      transform,
		})
	}
}

func handleBetaOrganizationWorkspacesMembersAdd(ctx context.Context, cmd *cli.Command) error {
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

	params := anthropic.BetaOrganizationWorkspaceMemberAddParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Beta.Organization.Workspaces.Members.Add(
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
		Title:          "beta:organization:workspaces:members add",
		Transform:      transform,
	})
}

func handleBetaOrganizationWorkspacesMembersRemove(ctx context.Context, cmd *cli.Command) error {
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

	params := anthropic.BetaOrganizationWorkspaceMemberRemoveParams{
		WorkspaceID: cmd.Value("workspace-id").(string),
	}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Beta.Organization.Workspaces.Members.Remove(
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
		Title:          "beta:organization:workspaces:members remove",
		Transform:      transform,
	})
}
