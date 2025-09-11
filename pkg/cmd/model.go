// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cmd

import (
	"context"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/stainless-sdks/anthropic-cli/pkg/jsonflag"
	"github.com/urfave/cli/v3"
)

var modelsRetrieve = cli.Command{
	Name:  "retrieve",
	Usage: "Get a specific model.",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name: "model-id",
		},
		&jsonflag.JSONStringFlag{
			Name: "betas",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Header,
				Path: "anthropic-beta.#",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "+beta",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Header,
				Path: "anthropic-beta.-1",
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
		&jsonflag.JSONStringFlag{
			Name: "after-id",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Query,
				Path: "after_id",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "before-id",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Query,
				Path: "before_id",
			},
		},
		&jsonflag.JSONIntFlag{
			Name: "limit",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Query,
				Path: "limit",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "betas",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Header,
				Path: "anthropic-beta.#",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "+beta",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Header,
				Path: "anthropic-beta.-1",
			},
		},
	},
	Action:          handleModelsList,
	HideHelpCommand: true,
}

func handleModelsRetrieve(ctx context.Context, cmd *cli.Command) error {
	cc := getAPICommandContext(cmd)
	params := anthropic.ModelGetParams{}
	res, err := cc.client.Models.Get(
		context.TODO(),
		cmd.Value("model-id").(string),
		params,
		option.WithMiddleware(cc.AsMiddleware()),
	)
	if err != nil {
		return err
	}

	format := cmd.Root().String("format")
	return ShowJSON("models retrieve", res.RawJSON(), format)
}

func handleModelsList(ctx context.Context, cmd *cli.Command) error {
	cc := getAPICommandContext(cmd)
	params := anthropic.ModelListParams{}
	res, err := cc.client.Models.List(
		context.TODO(),
		params,
		option.WithMiddleware(cc.AsMiddleware()),
	)
	if err != nil {
		return err
	}

	format := cmd.Root().String("format")
	return ShowJSON("models list", res.RawJSON(), format)
}
