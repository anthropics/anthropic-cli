// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/stainless-sdks/anthropic-cli/pkg/jsonflag"
	"github.com/urfave/cli/v3"
)

var betaModelsRetrieve = cli.Command{
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
	Action:          handleBetaModelsRetrieve,
	HideHelpCommand: true,
}

var betaModelsList = cli.Command{
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
	Action:          handleBetaModelsList,
	HideHelpCommand: true,
}

func handleBetaModelsRetrieve(ctx context.Context, cmd *cli.Command) error {
	cc := getAPICommandContext(cmd)
	params := anthropic.BetaModelGetParams{}
	res, err := cc.client.Beta.Models.Get(
		context.TODO(),
		cmd.Value("model-id").(string),
		params,
		option.WithMiddleware(cc.AsMiddleware()),
	)
	if err != nil {
		return err
	}

	fmt.Printf("%s\n", ColorizeJSON(res.RawJSON(), os.Stdout))
	return nil
}

func handleBetaModelsList(ctx context.Context, cmd *cli.Command) error {
	cc := getAPICommandContext(cmd)
	params := anthropic.BetaModelListParams{}
	res, err := cc.client.Beta.Models.List(
		context.TODO(),
		params,
		option.WithMiddleware(cc.AsMiddleware()),
	)
	if err != nil {
		return err
	}

	fmt.Printf("%s\n", ColorizeJSON(res.RawJSON(), os.Stdout))
	return nil
}
