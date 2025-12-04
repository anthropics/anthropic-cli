// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cmd

import (
	"context"
	"fmt"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/stainless-sdks/anthropic-cli/internal/apiquery"
	"github.com/stainless-sdks/anthropic-cli/internal/requestflag"
	"github.com/tidwall/gjson"
	"github.com/urfave/cli/v3"
)

var betaSkillsCreate = cli.Command{
	Name:  "create",
	Usage: "Create Skill",
	Flags: []cli.Flag{
		&requestflag.StringFlag{
			Name:  "display-title",
			Usage: "Display title for the skill.\n\nThis is a human-readable label that is not included in the prompt sent to the model.",
			Config: requestflag.RequestConfig{
				BodyPath: "display_title",
			},
		},
		&requestflag.StringSliceFlag{
			Name:  "file",
			Usage: "Files to upload for the skill.\n\nAll files must be in the same top-level directory and must include a SKILL.md file at the root of that directory.",
			Config: requestflag.RequestConfig{
				BodyPath: "files",
			},
		},
		&requestflag.YAMLSliceFlag{
			Name:  "beta",
			Usage: "Optional header to specify the beta version(s) you want to use.",
			Config: requestflag.RequestConfig{
				HeaderPath: "anthropic-beta",
			},
		},
	},
	Action:          handleBetaSkillsCreate,
	HideHelpCommand: true,
}

var betaSkillsRetrieve = cli.Command{
	Name:  "retrieve",
	Usage: "Get Skill",
	Flags: []cli.Flag{
		&requestflag.StringFlag{
			Name:  "skill-id",
			Usage: "Unique identifier for the skill.\n\nThe format and length of IDs may change over time.",
		},
		&requestflag.YAMLSliceFlag{
			Name:  "beta",
			Usage: "Optional header to specify the beta version(s) you want to use.",
			Config: requestflag.RequestConfig{
				HeaderPath: "anthropic-beta",
			},
		},
	},
	Action:          handleBetaSkillsRetrieve,
	HideHelpCommand: true,
}

var betaSkillsList = cli.Command{
	Name:  "list",
	Usage: "List Skills",
	Flags: []cli.Flag{
		&requestflag.IntFlag{
			Name:  "limit",
			Usage: "Number of results to return per page.\n\nMaximum value is 100. Defaults to 20.",
			Value: requestflag.Value[int64](20),
			Config: requestflag.RequestConfig{
				QueryPath: "limit",
			},
		},
		&requestflag.StringFlag{
			Name:  "page",
			Usage: "Pagination token for fetching a specific page of results.\n\nPass the value from a previous response's `next_page` field to get the next page of results.",
			Config: requestflag.RequestConfig{
				QueryPath: "page",
			},
		},
		&requestflag.StringFlag{
			Name:  "source",
			Usage: "Filter skills by source.\n\nIf provided, only skills from the specified source will be returned:\n* `\"custom\"`: only return user-created skills\n* `\"anthropic\"`: only return Anthropic-created skills",
			Config: requestflag.RequestConfig{
				QueryPath: "source",
			},
		},
		&requestflag.YAMLSliceFlag{
			Name:  "beta",
			Usage: "Optional header to specify the beta version(s) you want to use.",
			Config: requestflag.RequestConfig{
				HeaderPath: "anthropic-beta",
			},
		},
	},
	Action:          handleBetaSkillsList,
	HideHelpCommand: true,
}

var betaSkillsDelete = cli.Command{
	Name:  "delete",
	Usage: "Delete Skill",
	Flags: []cli.Flag{
		&requestflag.StringFlag{
			Name:  "skill-id",
			Usage: "Unique identifier for the skill.\n\nThe format and length of IDs may change over time.",
		},
		&requestflag.YAMLSliceFlag{
			Name:  "beta",
			Usage: "Optional header to specify the beta version(s) you want to use.",
			Config: requestflag.RequestConfig{
				HeaderPath: "anthropic-beta",
			},
		},
	},
	Action:          handleBetaSkillsDelete,
	HideHelpCommand: true,
}

func handleBetaSkillsCreate(ctx context.Context, cmd *cli.Command) error {
	client := anthropic.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if len(unusedArgs) > 0 {
		return fmt.Errorf("Unexpected extra arguments: %v", unusedArgs)
	}
	params := anthropic.BetaSkillNewParams{}

	options, err := flagOptions(
		cmd,
		apiquery.NestedQueryFormatBrackets,
		apiquery.ArrayQueryFormatComma,
		MultipartFormEncoded,
	)
	if err != nil {
		return err
	}
	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Beta.Skills.New(
		ctx,
		params,
		options...,
	)
	if err != nil {
		return err
	}

	json := gjson.Parse(string(res))
	format := cmd.Root().String("format")
	transform := cmd.Root().String("transform")
	return ShowJSON("beta:skills create", json, format, transform)
}

func handleBetaSkillsRetrieve(ctx context.Context, cmd *cli.Command) error {
	client := anthropic.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("skill-id") && len(unusedArgs) > 0 {
		cmd.Set("skill-id", unusedArgs[0])
		unusedArgs = unusedArgs[1:]
	}
	if len(unusedArgs) > 0 {
		return fmt.Errorf("Unexpected extra arguments: %v", unusedArgs)
	}
	params := anthropic.BetaSkillGetParams{}

	options, err := flagOptions(
		cmd,
		apiquery.NestedQueryFormatBrackets,
		apiquery.ArrayQueryFormatComma,
		ApplicationJSON,
	)
	if err != nil {
		return err
	}
	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Beta.Skills.Get(
		ctx,
		requestflag.CommandRequestValue[string](cmd, "skill-id"),
		params,
		options...,
	)
	if err != nil {
		return err
	}

	json := gjson.Parse(string(res))
	format := cmd.Root().String("format")
	transform := cmd.Root().String("transform")
	return ShowJSON("beta:skills retrieve", json, format, transform)
}

func handleBetaSkillsList(ctx context.Context, cmd *cli.Command) error {
	client := anthropic.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if len(unusedArgs) > 0 {
		return fmt.Errorf("Unexpected extra arguments: %v", unusedArgs)
	}
	params := anthropic.BetaSkillListParams{}

	options, err := flagOptions(
		cmd,
		apiquery.NestedQueryFormatBrackets,
		apiquery.ArrayQueryFormatComma,
		ApplicationJSON,
	)
	if err != nil {
		return err
	}
	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Beta.Skills.List(
		ctx,
		params,
		options...,
	)
	if err != nil {
		return err
	}

	json := gjson.Parse(string(res))
	format := cmd.Root().String("format")
	transform := cmd.Root().String("transform")
	return ShowJSON("beta:skills list", json, format, transform)
}

func handleBetaSkillsDelete(ctx context.Context, cmd *cli.Command) error {
	client := anthropic.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("skill-id") && len(unusedArgs) > 0 {
		cmd.Set("skill-id", unusedArgs[0])
		unusedArgs = unusedArgs[1:]
	}
	if len(unusedArgs) > 0 {
		return fmt.Errorf("Unexpected extra arguments: %v", unusedArgs)
	}
	params := anthropic.BetaSkillDeleteParams{}

	options, err := flagOptions(
		cmd,
		apiquery.NestedQueryFormatBrackets,
		apiquery.ArrayQueryFormatComma,
		ApplicationJSON,
	)
	if err != nil {
		return err
	}
	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Beta.Skills.Delete(
		ctx,
		requestflag.CommandRequestValue[string](cmd, "skill-id"),
		params,
		options...,
	)
	if err != nil {
		return err
	}

	json := gjson.Parse(string(res))
	format := cmd.Root().String("format")
	transform := cmd.Root().String("transform")
	return ShowJSON("beta:skills delete", json, format, transform)
}
