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

var skillsVersionsCreate = cli.Command{
	Name:    "create",
	Usage:   "Create Skill Version",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "skill-id",
			Usage:     "Unique identifier for the skill.\n\nThe format and length of IDs may change over time.",
			Required:  true,
			PathParam: "skill_id",
		},
		&requestflag.Flag[[]string]{
			Name:      "file",
			Usage:     "Files to upload for the skill.\n\nAll files must be in the same top-level directory and must include a SKILL.md file at the root of that directory.",
			Required:  true,
			BodyPath:  "files",
			FileInput: true,
		},
	},
	Action:          handleSkillsVersionsCreate,
	HideHelpCommand: true,
}

var skillsVersionsRetrieve = cli.Command{
	Name:    "retrieve",
	Usage:   "Get Skill Version",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "skill-id",
			Usage:     "Unique identifier for the skill.\n\nThe format and length of IDs may change over time.",
			Required:  true,
			PathParam: "skill_id",
		},
		&requestflag.Flag[string]{
			Name:      "version",
			Usage:     "Identifies the skill version: a version ID, or — where the endpoint accepts it — the literal `latest` for the skill's most recent version.\n\nRequests carrying the `skills-2025-10-02` beta header address versions by their Unix epoch timestamp instead (e.g., \"1759178010641129\").",
			Required:  true,
			PathParam: "version",
		},
	},
	Action:          handleSkillsVersionsRetrieve,
	HideHelpCommand: true,
}

var skillsVersionsList = cli.Command{
	Name:    "list",
	Usage:   "List Skill Versions",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "skill-id",
			Usage:     "Unique identifier for the skill.\n\nThe format and length of IDs may change over time.",
			Required:  true,
			PathParam: "skill_id",
		},
		&requestflag.Flag[int64]{
			Name:      "limit",
			Usage:     "Number of results to return per page.\n\nRanges from `1` to `1000`. Defaults to `20`.",
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
	Action:          handleSkillsVersionsList,
	HideHelpCommand: true,
}

var skillsVersionsDelete = cli.Command{
	Name:    "delete",
	Usage:   "Delete Skill Version",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "skill-id",
			Usage:     "Unique identifier for the skill.\n\nThe format and length of IDs may change over time.",
			Required:  true,
			PathParam: "skill_id",
		},
		&requestflag.Flag[string]{
			Name:      "version",
			Usage:     "Identifies the skill version: a version ID, or — where the endpoint accepts it — the literal `latest` for the skill's most recent version.\n\nRequests carrying the `skills-2025-10-02` beta header address versions by their Unix epoch timestamp instead (e.g., \"1759178010641129\").",
			Required:  true,
			PathParam: "version",
		},
	},
	Action:          handleSkillsVersionsDelete,
	HideHelpCommand: true,
}

func handleSkillsVersionsCreate(ctx context.Context, cmd *cli.Command) error {
	client := anthropic.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("skill-id") && len(unusedArgs) > 0 {
		cmd.Set("skill-id", unusedArgs[0])
		unusedArgs = unusedArgs[1:]
	}
	if len(unusedArgs) > 0 {
		return fmt.Errorf("Unexpected extra arguments: %v", unusedArgs)
	}

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

	params := anthropic.SkillVersionNewParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Skills.Versions.New(
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
	explicitFormat := cmd.Root().IsSet("format")
	transform := cmd.Root().String("transform")
	return ShowJSON(obj, ShowJSONOpts{
		ExplicitFormat: explicitFormat,
		Format:         format,
		RawOutput:      cmd.Root().Bool("raw-output"),
		Title:          "skills:versions create",
		Transform:      transform,
	})
}

func handleSkillsVersionsRetrieve(ctx context.Context, cmd *cli.Command) error {
	client := anthropic.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("version") && len(unusedArgs) > 0 {
		cmd.Set("version", unusedArgs[0])
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

	params := anthropic.SkillVersionGetParams{
		SkillID: cmd.Value("skill-id").(string),
	}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Skills.Versions.Get(
		ctx,
		cmd.Value("version").(string),
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
		Title:          "skills:versions retrieve",
		Transform:      transform,
	})
}

func handleSkillsVersionsList(ctx context.Context, cmd *cli.Command) error {
	client := anthropic.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("skill-id") && len(unusedArgs) > 0 {
		cmd.Set("skill-id", unusedArgs[0])
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

	params := anthropic.SkillVersionListParams{}

	format := "explore"
	explicitFormat := cmd.Root().IsSet("format")
	if explicitFormat {
		format = cmd.Root().String("format")
	}
	transform := cmd.Root().String("transform")
	if format == "raw" {
		var res []byte
		options = append(options, option.WithResponseBodyInto(&res))
		_, err = client.Skills.Versions.List(
			ctx,
			cmd.Value("skill-id").(string),
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
			Title:          "skills:versions list",
			Transform:      transform,
		})
	} else {
		iter := client.Skills.Versions.ListAutoPaging(
			ctx,
			cmd.Value("skill-id").(string),
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
			Title:          "skills:versions list",
			Transform:      transform,
		})
	}
}

func handleSkillsVersionsDelete(ctx context.Context, cmd *cli.Command) error {
	client := anthropic.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("version") && len(unusedArgs) > 0 {
		cmd.Set("version", unusedArgs[0])
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

	params := anthropic.SkillVersionDeleteParams{
		SkillID: cmd.Value("skill-id").(string),
	}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Skills.Versions.Delete(
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
	explicitFormat := cmd.Root().IsSet("format")
	transform := cmd.Root().String("transform")
	return ShowJSON(obj, ShowJSONOpts{
		ExplicitFormat: explicitFormat,
		Format:         format,
		RawOutput:      cmd.Root().Bool("raw-output"),
		Title:          "skills:versions delete",
		Transform:      transform,
	})
}
