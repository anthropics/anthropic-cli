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

var betaComplianceReviewItemsRetrieve = cli.Command{
	Name:    "retrieve",
	Usage:   "Read one message that Anthropic flagged for your organization to review, with\nits request and response.",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:        "item-id",
			Usage:       "The ID of the item. For a message, its ID as the client received it or as a review case lists it, such as `msg_...`.",
			Required:    true,
			PathParam:   "item_id",
			DataAliases: []string{"id"},
		},
		&requestflag.Flag[[]string]{
			Name:       "beta",
			Usage:      "Optional header to specify the beta version(s) you want to use.",
			HeaderPath: "anthropic-beta",
		},
		&requestflag.Flag[string]{
			Name:        "user-profile-id",
			Usage:       "The ID of the user profile that the message was recorded with (`uprof_...`). Send it for a message recorded with a user profile, and leave it out for a message recorded with none.",
			HeaderPath:  "anthropic-user-profile-id",
			DataAliases: []string{"user_profile_id"},
		},
	},
	Action:          handleBetaComplianceReviewItemsRetrieve,
	HideHelpCommand: true,
}

func handleBetaComplianceReviewItemsRetrieve(ctx context.Context, cmd *cli.Command) error {
	client := anthropic.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("item-id") && len(unusedArgs) > 0 {
		cmd.Set("item-id", unusedArgs[0])
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

	params := anthropic.BetaComplianceReviewItemGetParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Beta.Compliance.Review.Items.Get(
		ctx,
		cmd.Value("item-id").(string),
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
		Title:          "beta:compliance:review:items retrieve",
		Transform:      transform,
	})
}
