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

var userProfilesCreate = requestflag.WithInnerFlags(cli.Command{
	Name:    "create",
	Usage:   "Register one of your end customers and get the ID to send on the requests you\nmake for them.",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:     "access-type",
			Usage:    "How the platform uses the API on behalf of the entity this profile represents. `application`: the platform sells a product that uses the API behind the scenes, and the profile represents an individual end-user of that product. `passthrough`: the platform resells raw inference, and the profile identifies the resold-to company.",
			BodyPath: "access_type",
		},
		&requestflag.Flag[map[string]any]{
			Name:     "external-user-details",
			BodyPath: "external_user_details",
		},
		&requestflag.Flag[map[string]any]{
			Name:     "metadata",
			Usage:    "Free-form key-value data to attach to this user profile. Maximum 16 keys, with keys up to 64 characters and values up to 512 characters. Values must be non-empty strings.",
			BodyPath: "metadata",
		},
		&requestflag.Flag[*string]{
			Name:     "name",
			Usage:    "Optional for all profiles. Real-world name of the entity this profile represents (company or individual); for a company the platform resells Claude access to (`access_type` `passthrough`), that company's name where known. Maximum 255 characters.",
			BodyPath: "name",
		},
		&requestflag.Flag[string]{
			Name:       "workspace-id",
			Usage:      "Optional header to select the Workspace for this request. The value is a Workspace ID (for example, `wrkspc_011CZkZaBF1tNoB5wlCeusgy`).\n\nOnly needed for credentials that can act on more than one Workspace. A credential that belongs to a specific Workspace may omit it; if sent, it must match that Workspace.",
			HeaderPath: "anthropic-workspace-id",
		},
	},
	Action:          handleUserProfilesCreate,
	HideHelpCommand: true,
}, map[string][]requestflag.HasOuterFlag{
	"external-user-details": {
		&requestflag.InnerFlag[*string]{
			Name:       "external-user-details.account-status",
			Usage:      "The status of the entity's account on the platform, as the platform states it: `active`; `suspended`, when the platform has restricted the account and may restore it; or `blocked`, when the platform has barred it. It records the platform's decision only. The statuses that Anthropic itself reports for the profile, such as its status in a program, are Anthropic's decisions and do not change with it.",
			InnerField: "account_status",
		},
		&requestflag.InnerFlag[*string]{
			Name:       "external-user-details.country",
			Usage:      "The country of the entity (not of the platform), as the platform determines it: an ISO 3166-1 alpha-2 code in upper case, for example `US`. Only the form, two uppercase ASCII letters, is checked.",
			InnerField: "country",
		},
		&requestflag.InnerFlag[*string]{
			Name:       "external-user-details.email-hash",
			Usage:      "A hash of the entity's email address, computed by the platform. Anthropic treats it as an opaque string and does not prescribe the hash function. 1 to 255 characters.",
			InnerField: "email_hash",
		},
		&requestflag.InnerFlag[*string]{
			Name:       "external-user-details.entity-type",
			Usage:      "What kind of entity the profile represents, as the platform states it: `individual`, `business`, `non_profit` or `government`.",
			InnerField: "entity_type",
		},
		&requestflag.InnerFlag[*string]{
			Name:       "external-user-details.name-hash",
			Usage:      "A hash of the entity's name, computed by the platform. Anthropic treats it as an opaque string and does not prescribe the hash function. 1 to 255 characters.",
			InnerField: "name_hash",
		},
		&requestflag.InnerFlag[any]{
			Name:       "external-user-details.onboarded-at",
			Usage:      "When the entity opened its account with the platform, in RFC 3339 format: for an `application` profile, when the end-user signed up; for a `passthrough` profile, when the company became the platform's customer. Must be a complete timestamp no more than 1 minute in the future.",
			InnerField: "onboarded_at",
		},
		&requestflag.InnerFlag[*string]{
			Name:       "external-user-details.reference-id",
			Usage:      "The platform's own reference for the entity, for example the key of the end-user's row in the platform's database. Not interpreted by Anthropic and not enforced unique. 1 to 255 characters.",
			InnerField: "reference_id",
		},
	},
})

var userProfilesRetrieve = cli.Command{
	Name:    "retrieve",
	Usage:   "Get a user profile by ID to check the details stored for it.",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:        "user-profile-id",
			Usage:       "The ID of the user profile to get (`uprof_...`).",
			Required:    true,
			PathParam:   "user_profile_id",
			DataAliases: []string{"id"},
		},
		&requestflag.Flag[string]{
			Name:       "workspace-id",
			Usage:      "Optional header to select the Workspace for this request. The value is a Workspace ID (for example, `wrkspc_011CZkZaBF1tNoB5wlCeusgy`).\n\nOnly needed for credentials that can act on more than one Workspace. A credential that belongs to a specific Workspace may omit it; if sent, it must match that Workspace.",
			HeaderPath: "anthropic-workspace-id",
		},
	},
	Action:          handleUserProfilesRetrieve,
	HideHelpCommand: true,
}

var userProfilesUpdate = requestflag.WithInnerFlags(cli.Command{
	Name:    "update",
	Usage:   "Change the details of a user profile, leaving the fields that you don't send as\nthey are.",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:        "user-profile-id",
			Usage:       "The ID of the user profile to update (`uprof_...`).",
			Required:    true,
			PathParam:   "user_profile_id",
			DataAliases: []string{"id"},
		},
		&requestflag.Flag[*string]{
			Name:     "access-type",
			Usage:    "How the platform uses the API on behalf of the entity this profile represents. `application`: the platform sells a product that uses the API behind the scenes, and the profile represents an individual end-user of that product. `passthrough`: the platform resells raw inference, and the profile identifies the resold-to company.",
			BodyPath: "access_type",
		},
		&requestflag.Flag[map[string]any]{
			Name:     "external-user-details",
			BodyPath: "external_user_details",
		},
		&requestflag.Flag[map[string]any]{
			Name:     "metadata",
			Usage:    "Key-value pairs to merge into the stored metadata. Keys provided overwrite existing values. To remove a key, set its value to an empty string. Keys not provided are left unchanged. Maximum 16 keys, with keys up to 64 characters and values up to 512 characters.",
			BodyPath: "metadata",
		},
		&requestflag.Flag[*string]{
			Name:     "name",
			Usage:    "If present, replaces the stored name. Omit to leave unchanged. Maximum 255 characters.",
			BodyPath: "name",
		},
		&requestflag.Flag[string]{
			Name:       "workspace-id",
			Usage:      "Optional header to select the Workspace for this request. The value is a Workspace ID (for example, `wrkspc_011CZkZaBF1tNoB5wlCeusgy`).\n\nOnly needed for credentials that can act on more than one Workspace. A credential that belongs to a specific Workspace may omit it; if sent, it must match that Workspace.",
			HeaderPath: "anthropic-workspace-id",
		},
	},
	Action:          handleUserProfilesUpdate,
	HideHelpCommand: true,
}, map[string][]requestflag.HasOuterFlag{
	"external-user-details": {
		&requestflag.InnerFlag[*string]{
			Name:       "external-user-details.account-status",
			Usage:      "The status of the entity's account on the platform, as the platform states it: `active`; `suspended`, when the platform has restricted the account and may restore it; or `blocked`, when the platform has barred it. It records the platform's decision only. The statuses that Anthropic itself reports for the profile, such as its status in a program, are Anthropic's decisions and do not change with it.",
			InnerField: "account_status",
		},
		&requestflag.InnerFlag[*string]{
			Name:       "external-user-details.country",
			Usage:      "The country of the entity (not of the platform), as the platform determines it: an ISO 3166-1 alpha-2 code in upper case, for example `US`. Only the form, two uppercase ASCII letters, is checked.",
			InnerField: "country",
		},
		&requestflag.InnerFlag[*string]{
			Name:       "external-user-details.email-hash",
			Usage:      "A hash of the entity's email address, computed by the platform. Anthropic treats it as an opaque string and does not prescribe the hash function. 1 to 255 characters.",
			InnerField: "email_hash",
		},
		&requestflag.InnerFlag[*string]{
			Name:       "external-user-details.entity-type",
			Usage:      "What kind of entity the profile represents, as the platform states it: `individual`, `business`, `non_profit` or `government`.",
			InnerField: "entity_type",
		},
		&requestflag.InnerFlag[*string]{
			Name:       "external-user-details.name-hash",
			Usage:      "A hash of the entity's name, computed by the platform. Anthropic treats it as an opaque string and does not prescribe the hash function. 1 to 255 characters.",
			InnerField: "name_hash",
		},
		&requestflag.InnerFlag[any]{
			Name:       "external-user-details.onboarded-at",
			Usage:      "When the entity opened its account with the platform, in RFC 3339 format: for an `application` profile, when the end-user signed up; for a `passthrough` profile, when the company became the platform's customer. Must be a complete timestamp no more than 1 minute in the future.",
			InnerField: "onboarded_at",
		},
		&requestflag.InnerFlag[*string]{
			Name:       "external-user-details.reference-id",
			Usage:      "The platform's own reference for the entity, for example the key of the end-user's row in the platform's database. Not interpreted by Anthropic and not enforced unique. 1 to 255 characters.",
			InnerField: "reference_id",
		},
	},
})

var userProfilesList = cli.Command{
	Name:    "list",
	Usage:   "List the user profiles in your organization, newest first by default.",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[[]string]{
			Name:      "enforcement-status",
			Usage:     "Return only the user profiles whose `enforcement_status` is one of these values: `unrestricted`, `pending`, `restricted` or `banned`. Repeat the parameter to match any of several. With it set, a page may hold fewer than `limit` profiles, or none, and still carry a `next_page`: the service examines a bounded number of profiles per request, and only a null `next_page` means the listing is complete. Rejected with `invalid_request_error` under the `user-profiles-2026-09-04` and earlier beta headers. Fails with `overloaded_error` when a profile's enforcement actions cannot be read for the moment.",
			QueryPath: "enforcement_status",
		},
		&requestflag.Flag[int64]{
			Name:      "limit",
			Usage:     "The maximum number of user profiles to return, from 1 to 100. Defaults to 20.",
			QueryPath: "limit",
		},
		&requestflag.Flag[string]{
			Name:      "order",
			Usage:     "The direction that a list of user profiles is sorted in.",
			QueryPath: "order",
		},
		&requestflag.Flag[string]{
			Name:      "order-by",
			Usage:     "Sort field for listing user profiles: `created_at` (default) or `name` (case-insensitive; profiles without a name sort last).",
			QueryPath: "order_by",
		},
		&requestflag.Flag[string]{
			Name:      "page",
			Usage:     "The cursor for the page to return, taken from `next_page` in a previous response.\n\nLeave it out to get the first page.",
			QueryPath: "page",
		},
		&requestflag.Flag[string]{
			Name:       "workspace-id",
			Usage:      "Optional header to select the Workspace for this request. The value is a Workspace ID (for example, `wrkspc_011CZkZaBF1tNoB5wlCeusgy`).\n\nOnly needed for credentials that can act on more than one Workspace. A credential that belongs to a specific Workspace may omit it; if sent, it must match that Workspace.",
			HeaderPath: "anthropic-workspace-id",
		},
		&requestflag.Flag[int64]{
			Name:  "max-items",
			Usage: "The maximum number of items to return (use -1 for unlimited).",
		},
	},
	Action:          handleUserProfilesList,
	HideHelpCommand: true,
}

var userProfilesDelete = cli.Command{
	Name:    "delete",
	Usage:   "Delete a user profile, for example when the end customer closes their account.",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:        "user-profile-id",
			Usage:       "The ID of the user profile to delete (`uprof_...`).",
			Required:    true,
			PathParam:   "user_profile_id",
			DataAliases: []string{"id"},
		},
		&requestflag.Flag[string]{
			Name:       "workspace-id",
			Usage:      "Optional header to select the Workspace for this request. The value is a Workspace ID (for example, `wrkspc_011CZkZaBF1tNoB5wlCeusgy`).\n\nOnly needed for credentials that can act on more than one Workspace. A credential that belongs to a specific Workspace may omit it; if sent, it must match that Workspace.",
			HeaderPath: "anthropic-workspace-id",
		},
	},
	Action:          handleUserProfilesDelete,
	HideHelpCommand: true,
}

func handleUserProfilesCreate(ctx context.Context, cmd *cli.Command) error {
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

	params := anthropic.UserProfileNewParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.UserProfiles.New(ctx, params, options...)
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
		Title:          "user-profiles create",
		Transform:      transform,
	})
}

func handleUserProfilesRetrieve(ctx context.Context, cmd *cli.Command) error {
	client := anthropic.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("user-profile-id") && len(unusedArgs) > 0 {
		cmd.Set("user-profile-id", unusedArgs[0])
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

	params := anthropic.UserProfileGetParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.UserProfiles.Get(
		ctx,
		cmd.Value("user-profile-id").(string),
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
		Title:          "user-profiles retrieve",
		Transform:      transform,
	})
}

func handleUserProfilesUpdate(ctx context.Context, cmd *cli.Command) error {
	client := anthropic.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("user-profile-id") && len(unusedArgs) > 0 {
		cmd.Set("user-profile-id", unusedArgs[0])
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

	params := anthropic.UserProfileUpdateParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.UserProfiles.Update(
		ctx,
		cmd.Value("user-profile-id").(string),
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
		Title:          "user-profiles update",
		Transform:      transform,
	})
}

func handleUserProfilesList(ctx context.Context, cmd *cli.Command) error {
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

	params := anthropic.UserProfileListParams{}

	format := "explore"
	explicitFormat := cmd.Root().IsSet("format")
	if explicitFormat {
		format = cmd.Root().String("format")
	}
	transform := cmd.Root().String("transform")
	if format == "raw" {
		var res []byte
		options = append(options, option.WithResponseBodyInto(&res))
		_, err = client.UserProfiles.List(ctx, params, options...)
		if err != nil {
			return err
		}
		obj := gjson.ParseBytes(res)
		return ShowJSON(obj, ShowJSONOpts{
			ExplicitFormat: explicitFormat,
			Format:         format,
			RawOutput:      cmd.Root().Bool("raw-output"),
			Title:          "user-profiles list",
			Transform:      transform,
		})
	} else {
		iter := client.UserProfiles.ListAutoPaging(ctx, params, options...)
		maxItems := int64(-1)
		if cmd.IsSet("max-items") {
			maxItems = cmd.Value("max-items").(int64)
		}
		return ShowJSONIterator(iter, maxItems, ShowJSONOpts{
			ExplicitFormat: explicitFormat,
			Format:         format,
			RawOutput:      cmd.Root().Bool("raw-output"),
			Title:          "user-profiles list",
			Transform:      transform,
		})
	}
}

func handleUserProfilesDelete(ctx context.Context, cmd *cli.Command) error {
	client := anthropic.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("user-profile-id") && len(unusedArgs) > 0 {
		cmd.Set("user-profile-id", unusedArgs[0])
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

	params := anthropic.UserProfileDeleteParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.UserProfiles.Delete(
		ctx,
		cmd.Value("user-profile-id").(string),
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
		Title:          "user-profiles delete",
		Transform:      transform,
	})
}
