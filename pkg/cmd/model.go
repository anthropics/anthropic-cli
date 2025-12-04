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

var modelsRetrieve = cli.Command{
	Name:  "retrieve",
	Usage: "Get a specific model.",
	Flags: []cli.Flag{
		&requestflag.StringFlag{
			Name:  "model-id",
			Usage: "Model identifier or alias.",
		},
		&requestflag.YAMLSliceFlag{
			Name:  "beta",
			Usage: "Optional header to specify the beta version(s) you want to use.",
			Config: requestflag.RequestConfig{
				HeaderPath: "anthropic-beta",
			},
		},
	},
	Action:          handleModelsRetrieve,
	HideHelpCommand: true,
}

var modelsList = cli.Command{
	Name:  "list",
	Usage: "List available models.",
	Flags: []cli.Flag{
		&requestflag.StringFlag{
			Name:  "after-id",
			Usage: "ID of the object to use as a cursor for pagination. When provided, returns the page of results immediately after this object.",
			Config: requestflag.RequestConfig{
				QueryPath: "after_id",
			},
		},
		&requestflag.StringFlag{
			Name:  "before-id",
			Usage: "ID of the object to use as a cursor for pagination. When provided, returns the page of results immediately before this object.",
			Config: requestflag.RequestConfig{
				QueryPath: "before_id",
			},
		},
		&requestflag.IntFlag{
			Name:  "limit",
			Usage: "Number of items to return per page.\n\nDefaults to `20`. Ranges from `1` to `1000`.",
			Value: requestflag.Value[int64](20),
			Config: requestflag.RequestConfig{
				QueryPath: "limit",
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
	Action:          handleModelsList,
	HideHelpCommand: true,
}

func handleModelsRetrieve(ctx context.Context, cmd *cli.Command) error {
	client := anthropic.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("model-id") && len(unusedArgs) > 0 {
		cmd.Set("model-id", unusedArgs[0])
		unusedArgs = unusedArgs[1:]
	}
	if len(unusedArgs) > 0 {
		return fmt.Errorf("Unexpected extra arguments: %v", unusedArgs)
	}
	params := anthropic.ModelGetParams{}

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
	_, err = client.Models.Get(
		ctx,
		requestflag.CommandRequestValue[string](cmd, "model-id"),
		params,
		options...,
	)
	if err != nil {
		return err
	}

	json := gjson.Parse(string(res))
	format := cmd.Root().String("format")
	transform := cmd.Root().String("transform")
	return ShowJSON("models retrieve", json, format, transform)
}

func handleModelsList(ctx context.Context, cmd *cli.Command) error {
	client := anthropic.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if len(unusedArgs) > 0 {
		return fmt.Errorf("Unexpected extra arguments: %v", unusedArgs)
	}
	params := anthropic.ModelListParams{}

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
	_, err = client.Models.List(
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
	return ShowJSON("models list", json, format, transform)
}
