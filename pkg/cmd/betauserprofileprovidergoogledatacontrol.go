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

var betaUserProfilesProvidersGoogleDataControlsRetrieve = cli.Command{
	Name:    "retrieve",
	Usage:   "Get the profile's encryption keys and storage buckets for Claude on Google\nCloud's Agent Platform, with each region's status.",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "user-profile-id",
			Usage:     "The ID of the user profile whose Google-managed data controls to retrieve (`uprof_...`).",
			Required:  true,
			PathParam: "user_profile_id",
		},
		&requestflag.Flag[[]string]{
			Name:       "beta",
			Usage:      "This endpoint is in beta: requests must send `google-managed-data-controls-2026-09-01` in this header.",
			HeaderPath: "anthropic-beta",
		},
		&requestflag.Flag[string]{
			Name:       "workspace-id",
			Usage:      "Optional header to select the Workspace for this request. The value is a Workspace ID (for example, `wrkspc_011CZkZaBF1tNoB5wlCeusgy`).\n\nOnly needed for credentials that can act on more than one Workspace. A credential that belongs to a specific Workspace may omit it; if sent, it must match that Workspace.",
			HeaderPath: "anthropic-workspace-id",
		},
	},
	Action:          handleBetaUserProfilesProvidersGoogleDataControlsRetrieve,
	HideHelpCommand: true,
}

var betaUserProfilesProvidersGoogleDataControlsUpdate = requestflag.WithInnerFlags(cli.Command{
	Name:    "update",
	Usage:   "Turn a control on or off, or change the key or bucket set for a region. Controls\nand regions that you leave out stay as they are.",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "user-profile-id",
			Usage:     "The ID of the user profile whose Google-managed data controls to update (`uprof_...`).",
			Required:  true,
			PathParam: "user_profile_id",
		},
		&requestflag.Flag[map[string]any]{
			Name:     "cmek",
			Usage:    "The customer-managed encryption key control. Leave it out to keep the control as it is. Send `cmek`, `storage`, or both.",
			BodyPath: "cmek",
		},
		&requestflag.Flag[map[string]any]{
			Name:     "storage",
			Usage:    "The customer-owned storage control. Leave it out to keep the control as it is. Send `cmek`, `storage`, or both.",
			BodyPath: "storage",
		},
		&requestflag.Flag[[]string]{
			Name:       "beta",
			Usage:      "This endpoint is in beta: requests must send `google-managed-data-controls-2026-09-01` in this header.",
			HeaderPath: "anthropic-beta",
		},
		&requestflag.Flag[string]{
			Name:       "workspace-id",
			Usage:      "Optional header to select the Workspace for this request. The value is a Workspace ID (for example, `wrkspc_011CZkZaBF1tNoB5wlCeusgy`).\n\nOnly needed for credentials that can act on more than one Workspace. A credential that belongs to a specific Workspace may omit it; if sent, it must match that Workspace.",
			HeaderPath: "anthropic-workspace-id",
		},
	},
	Action:          handleBetaUserProfilesProvidersGoogleDataControlsUpdate,
	HideHelpCommand: true,
}, map[string][]requestflag.HasOuterFlag{
	"cmek": {
		&requestflag.InnerFlag[string]{
			Name:       "cmek.type",
			Usage:      `Allowed values: "disabled", "enabled".`,
			InnerField: "type",
		},
		&requestflag.InnerFlag[map[string]any]{
			Name:       "cmek.external-key-ids",
			Usage:      "The External Key to use in each region, keyed by region code, such as `us` or `eu`. An External Key ID (`ekey_...`) sets that region's key, `null` removes it, and a region you leave out keeps its key. While the control is on, at least one region must have a key.",
			InnerField: "external_key_ids",
		},
	},
	"storage": {
		&requestflag.InnerFlag[string]{
			Name:       "storage.type",
			Usage:      `Allowed values: "disabled", "enabled".`,
			InnerField: "type",
		},
		&requestflag.InnerFlag[map[string]any]{
			Name:       "storage.buckets",
			Usage:      "The bucket to use in each region, keyed by region code, such as `us` or `eu`. An object with `project_number` and `bucket_name` sets that region's bucket, `null` removes it, and a region you leave out keeps its bucket. A bucket can serve only one region. While the control is on, at least one region must have a bucket.",
			InnerField: "buckets",
		},
	},
})

func handleBetaUserProfilesProvidersGoogleDataControlsRetrieve(ctx context.Context, cmd *cli.Command) error {
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

	params := anthropic.BetaUserProfileProviderGoogleDataControlGetParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Beta.UserProfiles.Providers.Google.DataControls.Get(
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
		Title:          "beta:user-profiles:providers:google:data-controls retrieve",
		Transform:      transform,
	})
}

func handleBetaUserProfilesProvidersGoogleDataControlsUpdate(ctx context.Context, cmd *cli.Command) error {
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

	params := anthropic.BetaUserProfileProviderGoogleDataControlUpdateParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Beta.UserProfiles.Providers.Google.DataControls.Update(
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
		Title:          "beta:user-profiles:providers:google:data-controls update",
		Transform:      transform,
	})
}
