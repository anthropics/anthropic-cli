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

var betaOrganizationProvidersGoogleProjectsDataControlsRetrieve = cli.Command{
	Name:    "retrieve",
	Usage:   "Get a Google Cloud project's encryption keys and storage buckets for Claude on\nGoogle Cloud's Agent Platform, with each region's status.",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "project-number",
			Usage:     "The Google Cloud project's number, such as `123456789012`, not its project ID.",
			Required:  true,
			PathParam: "project_number",
		},
		&requestflag.Flag[[]string]{
			Name:       "beta",
			Usage:      "This endpoint is in beta: requests must send `google-managed-data-controls-2026-09-01` in this header.",
			HeaderPath: "anthropic-beta",
		},
	},
	Action:          handleBetaOrganizationProvidersGoogleProjectsDataControlsRetrieve,
	HideHelpCommand: true,
}

var betaOrganizationProvidersGoogleProjectsDataControlsUpdate = requestflag.WithInnerFlags(cli.Command{
	Name:    "update",
	Usage:   "Turn a control on or off, or change the key or bucket set for a region. Controls\nand regions that you leave out stay as they are.",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "project-number",
			Usage:     "The Google Cloud project's number, such as `123456789012`, not its project ID.",
			Required:  true,
			PathParam: "project_number",
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
	},
	Action:          handleBetaOrganizationProvidersGoogleProjectsDataControlsUpdate,
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

func handleBetaOrganizationProvidersGoogleProjectsDataControlsRetrieve(ctx context.Context, cmd *cli.Command) error {
	client := anthropic.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("project-number") && len(unusedArgs) > 0 {
		cmd.Set("project-number", unusedArgs[0])
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

	params := anthropic.BetaOrganizationProviderGoogleProjectDataControlGetParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Beta.Organization.Providers.Google.Projects.DataControls.Get(
		ctx,
		cmd.Value("project-number").(string),
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
		Title:          "beta:organization:providers:google:projects:data-controls retrieve",
		Transform:      transform,
	})
}

func handleBetaOrganizationProvidersGoogleProjectsDataControlsUpdate(ctx context.Context, cmd *cli.Command) error {
	client := anthropic.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("project-number") && len(unusedArgs) > 0 {
		cmd.Set("project-number", unusedArgs[0])
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

	params := anthropic.BetaOrganizationProviderGoogleProjectDataControlUpdateParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Beta.Organization.Providers.Google.Projects.DataControls.Update(
		ctx,
		cmd.Value("project-number").(string),
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
		Title:          "beta:organization:providers:google:projects:data-controls update",
		Transform:      transform,
	})
}
