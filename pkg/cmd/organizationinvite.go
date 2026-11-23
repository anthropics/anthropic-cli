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

var organizationInvitesCreate = cli.Command{
	Name:    "create",
	Usage:   "Invite a user to join the organization by email.",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:     "email",
			Usage:    "Email of the User.",
			Required: true,
			BodyPath: "email",
		},
		&requestflag.Flag[string]{
			Name:     "role",
			Usage:    "Role for the invited User.\n\nThe accepted values depend on the organization type. Console and API organizations accept `user`, `developer`, `billing`, and `claude_code_user`; `admin` cannot be assigned through the API. Claude Enterprise organizations accept `user` and `managed`; inviting with the `managed` role requires the `write:rbac_groups` scope in addition to `write:members`.\n\nAllowed values: \"billing\", \"claude_code_user\", \"developer\", \"managed\", \"user\".",
			Required: true,
			BodyPath: "role",
		},
		&requestflag.Flag[[]string]{
			Name:     "rbac-group-id",
			Usage:    "RBAC group IDs to assign to the User when the Invite is accepted. A non-empty array is accepted only for a Claude Enterprise organization with RBAC groups, and requires the key to carry the `write:rbac_groups` scope.",
			BodyPath: "rbac_group_ids",
		},
	},
	Action:          handleOrganizationInvitesCreate,
	HideHelpCommand: true,
}

var organizationInvitesRetrieve = cli.Command{
	Name:    "retrieve",
	Usage:   "Retrieve an invite by ID.",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "invite-id",
			Usage:     "ID of the Invite.",
			Required:  true,
			PathParam: "invite_id",
		},
	},
	Action:          handleOrganizationInvitesRetrieve,
	HideHelpCommand: true,
}

var organizationInvitesList = cli.Command{
	Name:    "list",
	Usage:   "List the organization's invites.",
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
		&requestflag.Flag[string]{
			Name:      "email",
			Usage:     "Filter by the email address the Invite was sent to. Matches the same way as the Users list's `email` filter (normalized, case-insensitive).",
			QueryPath: "email",
		},
		&requestflag.Flag[int64]{
			Name:      "limit",
			Usage:     "Number of items to return per page.\n\nDefaults to `20`. Ranges from `1` to `1000`.",
			Default:   20,
			QueryPath: "limit",
		},
		&requestflag.Flag[[]string]{
			Name:      "role",
			Usage:     "Filter to items whose `role` equals one of the supplied values. Repeatable; values are OR'ed together.\n\nAccepted values depend on the organization type: Console and API organizations accept `user`, `developer`, `billing`, `admin`, and `claude_code_user`; Claude Enterprise organizations accept `user`, `owner`, `primary_owner`, `membership_admin`, and `managed`.",
			QueryPath: "roles",
		},
		&requestflag.Flag[[]string]{
			Name:      "status",
			Usage:     "Filter by Invite status. Repeatable; values are OR'ed together. Omit to return `pending`, `accepted`, and `expired` Invites alike.",
			QueryPath: "statuses",
		},
		&requestflag.Flag[int64]{
			Name:  "max-items",
			Usage: "The maximum number of items to return (use -1 for unlimited).",
		},
	},
	Action:          handleOrganizationInvitesList,
	HideHelpCommand: true,
}

var organizationInvitesDelete = cli.Command{
	Name:    "delete",
	Usage:   "Delete a pending invite.",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "invite-id",
			Usage:     "ID of the Invite.",
			Required:  true,
			PathParam: "invite_id",
		},
	},
	Action:          handleOrganizationInvitesDelete,
	HideHelpCommand: true,
}

func handleOrganizationInvitesCreate(ctx context.Context, cmd *cli.Command) error {
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

	params := anthropic.OrganizationInviteNewParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Organization.Invites.New(ctx, params, options...)
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
		Title:          "organization:invites create",
		Transform:      transform,
	})
}

func handleOrganizationInvitesRetrieve(ctx context.Context, cmd *cli.Command) error {
	client := anthropic.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("invite-id") && len(unusedArgs) > 0 {
		cmd.Set("invite-id", unusedArgs[0])
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
	_, err = client.Organization.Invites.Get(ctx, cmd.Value("invite-id").(string), options...)
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
		Title:          "organization:invites retrieve",
		Transform:      transform,
	})
}

func handleOrganizationInvitesList(ctx context.Context, cmd *cli.Command) error {
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

	params := anthropic.OrganizationInviteListParams{}

	format := "explore"
	explicitFormat := cmd.Root().IsSet("format")
	if explicitFormat {
		format = cmd.Root().String("format")
	}
	transform := cmd.Root().String("transform")
	if format == "raw" {
		var res []byte
		options = append(options, option.WithResponseBodyInto(&res))
		_, err = client.Organization.Invites.List(ctx, params, options...)
		if err != nil {
			return err
		}
		obj := gjson.ParseBytes(res)
		return ShowJSON(obj, ShowJSONOpts{
			ExplicitFormat: explicitFormat,
			Format:         format,
			RawOutput:      cmd.Root().Bool("raw-output"),
			Title:          "organization:invites list",
			Transform:      transform,
		})
	} else {
		iter := client.Organization.Invites.ListAutoPaging(ctx, params, options...)
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
			Title:          "organization:invites list",
			Transform:      transform,
		})
	}
}

func handleOrganizationInvitesDelete(ctx context.Context, cmd *cli.Command) error {
	client := anthropic.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("invite-id") && len(unusedArgs) > 0 {
		cmd.Set("invite-id", unusedArgs[0])
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
	_, err = client.Organization.Invites.Delete(ctx, cmd.Value("invite-id").(string), options...)
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
		Title:          "organization:invites delete",
		Transform:      transform,
	})
}
