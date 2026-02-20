// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/stainless-sdks/anthropic-cli/internal/apiquery"
	"github.com/stainless-sdks/anthropic-cli/internal/requestflag"
	"github.com/tidwall/gjson"
	"github.com/urfave/cli/v3"
)

var betaSkillsVersionsCreate = cli.Command{
	Name:    "create",
	Usage:   "Create Skill Version",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:     "skill-id",
			Usage:    "Unique identifier for the skill.\n\nThe format and length of IDs may change over time.",
			Required: true,
		},
		&requestflag.Flag[any]{
			Name:     "file",
			Usage:    "Files to upload for the skill.\n\nAll files must be in the same top-level directory and must include a SKILL.md file at the root of that directory.",
			BodyPath: "files",
		},
		&requestflag.Flag[[]string]{
			Name:       "beta",
			Usage:      "Optional header to specify the beta version(s) you want to use.",
			HeaderPath: "anthropic-beta",
		},
	},
	Action:          handleBetaSkillsVersionsCreate,
	HideHelpCommand: true,
}

var betaSkillsVersionsRetrieve = cli.Command{
	Name:    "retrieve",
	Usage:   "Get Skill Version",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:     "skill-id",
			Usage:    "Unique identifier for the skill.\n\nThe format and length of IDs may change over time.",
			Required: true,
		},
		&requestflag.Flag[string]{
			Name:     "version",
			Usage:    "Version identifier for the skill.\n\nEach version is identified by a Unix epoch timestamp (e.g., \"1759178010641129\").",
			Required: true,
		},
		&requestflag.Flag[[]string]{
			Name:       "beta",
			Usage:      "Optional header to specify the beta version(s) you want to use.",
			HeaderPath: "anthropic-beta",
		},
	},
	Action:          handleBetaSkillsVersionsRetrieve,
	HideHelpCommand: true,
}

var betaSkillsVersionsList = cli.Command{
	Name:    "list",
	Usage:   "List Skill Versions",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:     "skill-id",
			Usage:    "Unique identifier for the skill.\n\nThe format and length of IDs may change over time.",
			Required: true,
		},
		&requestflag.Flag[any]{
			Name:      "limit",
			Usage:     "Number of items to return per page.\n\nDefaults to `20`. Ranges from `1` to `1000`.",
			QueryPath: "limit",
		},
		&requestflag.Flag[any]{
			Name:      "page",
			Usage:     "Optionally set to the `next_page` token from the previous response.",
			QueryPath: "page",
		},
		&requestflag.Flag[[]string]{
			Name:       "beta",
			Usage:      "Optional header to specify the beta version(s) you want to use.",
			HeaderPath: "anthropic-beta",
		},
	},
	Action:          handleBetaSkillsVersionsList,
	HideHelpCommand: true,
}

var betaSkillsVersionsDelete = cli.Command{
	Name:    "delete",
	Usage:   "Delete Skill Version",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:     "skill-id",
			Usage:    "Unique identifier for the skill.\n\nThe format and length of IDs may change over time.",
			Required: true,
		},
		&requestflag.Flag[string]{
			Name:     "version",
			Usage:    "Version identifier for the skill.\n\nEach version is identified by a Unix epoch timestamp (e.g., \"1759178010641129\").",
			Required: true,
		},
		&requestflag.Flag[[]string]{
			Name:       "beta",
			Usage:      "Optional header to specify the beta version(s) you want to use.",
			HeaderPath: "anthropic-beta",
		},
	},
	Action:          handleBetaSkillsVersionsDelete,
	HideHelpCommand: true,
}

func handleBetaSkillsVersionsCreate(ctx context.Context, cmd *cli.Command) error {
	client := anthropic.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("skill-id") && len(unusedArgs) > 0 {
		cmd.Set("skill-id", unusedArgs[0])
		unusedArgs = unusedArgs[1:]
	}
	if len(unusedArgs) > 0 {
		return fmt.Errorf("Unexpected extra arguments: %v", unusedArgs)
	}

	params := anthropic.BetaSkillVersionNewParams{}

	options, err := flagOptions(
		cmd,
		apiquery.NestedQueryFormatBrackets,
		apiquery.ArrayQueryFormatBrackets,
		MultipartFormEncoded,
		false,
	)
	if err != nil {
		return err
	}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Beta.Skills.Versions.New(
		ctx,
		cmd.Value("skill-id").(string),
		params,
		options...,
	)
	if err != nil {
		return err
	}

	obj := gjson.ParseBytes(res)
	format := cmd.Root().String("format")
	transform := cmd.Root().String("transform")
	return ShowJSON(os.Stdout, "beta:skills:versions create", obj, format, transform)
}

func handleBetaSkillsVersionsRetrieve(ctx context.Context, cmd *cli.Command) error {
	client := anthropic.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("version") && len(unusedArgs) > 0 {
		cmd.Set("version", unusedArgs[0])
		unusedArgs = unusedArgs[1:]
	}
	if len(unusedArgs) > 0 {
		return fmt.Errorf("Unexpected extra arguments: %v", unusedArgs)
	}

	params := anthropic.BetaSkillVersionGetParams{
		SkillID: cmd.Value("skill-id").(string),
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
	_, err = client.Beta.Skills.Versions.Get(
		ctx,
		cmd.Value("version").(string),
		params,
		options...,
	)
	if err != nil {
		return err
	}

	obj := gjson.ParseBytes(res)
	format := cmd.Root().String("format")
	transform := cmd.Root().String("transform")
	return ShowJSON(os.Stdout, "beta:skills:versions retrieve", obj, format, transform)
}

func handleBetaSkillsVersionsList(ctx context.Context, cmd *cli.Command) error {
	client := anthropic.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("skill-id") && len(unusedArgs) > 0 {
		cmd.Set("skill-id", unusedArgs[0])
		unusedArgs = unusedArgs[1:]
	}
	if len(unusedArgs) > 0 {
		return fmt.Errorf("Unexpected extra arguments: %v", unusedArgs)
	}

	params := anthropic.BetaSkillVersionListParams{}

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

	format := cmd.Root().String("format")
	transform := cmd.Root().String("transform")
	if format == "raw" {
		var res []byte
		options = append(options, option.WithResponseBodyInto(&res))
		_, err = client.Beta.Skills.Versions.List(
			ctx,
			cmd.Value("skill-id").(string),
			params,
			options...,
		)
		if err != nil {
			return err
		}
		obj := gjson.ParseBytes(res)
		return ShowJSON(os.Stdout, "beta:skills:versions list", obj, format, transform)
	} else {
		iter := client.Beta.Skills.Versions.ListAutoPaging(
			ctx,
			cmd.Value("skill-id").(string),
			params,
			options...,
		)
		return ShowJSONIterator(os.Stdout, "beta:skills:versions list", iter, format, transform)
	}
}

func handleBetaSkillsVersionsDelete(ctx context.Context, cmd *cli.Command) error {
	client := anthropic.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("version") && len(unusedArgs) > 0 {
		cmd.Set("version", unusedArgs[0])
		unusedArgs = unusedArgs[1:]
	}
	if len(unusedArgs) > 0 {
		return fmt.Errorf("Unexpected extra arguments: %v", unusedArgs)
	}

	params := anthropic.BetaSkillVersionDeleteParams{
		SkillID: cmd.Value("skill-id").(string),
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
	_, err = client.Beta.Skills.Versions.Delete(
		ctx,
		cmd.Value("version").(string),
		params,
		options...,
	)
	if err != nil {
		return err
	}

	obj := gjson.ParseBytes(res)
	format := cmd.Root().String("format")
	transform := cmd.Root().String("transform")
	return ShowJSON(os.Stdout, "beta:skills:versions delete", obj, format, transform)
}
