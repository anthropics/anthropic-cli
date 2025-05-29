// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/urfave/cli/v3"
)

var betaModelsRetrieve = cli.Command{
	Name:  "retrieve",
	Usage: "Get a specific model.",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name: "model-id",
		},
		&cli.StringFlag{
			Name:   "betas",
			Action: getAPIFlagAction[string]("header", "anthropic-beta.#"),
		},
		&cli.StringFlag{
			Name:   "+beta",
			Action: getAPIFlagAction[string]("header", "anthropic-beta.-1"),
		},
	},
	Before:          initAPICommand,
	Action:          handleBetaModelsRetrieve,
	HideHelpCommand: true,
}

var betaModelsList = cli.Command{
	Name:  "list",
	Usage: "List available models.",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:   "after-id",
			Action: getAPIFlagAction[string]("query", "after_id"),
		},
		&cli.StringFlag{
			Name:   "before-id",
			Action: getAPIFlagAction[string]("query", "before_id"),
		},
		&cli.Int64Flag{
			Name:   "limit",
			Action: getAPIFlagAction[int64]("query", "limit"),
		},
		&cli.StringFlag{
			Name:   "betas",
			Action: getAPIFlagAction[string]("header", "anthropic-beta.#"),
		},
		&cli.StringFlag{
			Name:   "+beta",
			Action: getAPIFlagAction[string]("header", "anthropic-beta.-1"),
		},
	},
	Before:          initAPICommand,
	Action:          handleBetaModelsList,
	HideHelpCommand: true,
}

func handleBetaModelsRetrieve(ctx context.Context, cmd *cli.Command) error {
	cc := getAPICommandContext(ctx, cmd)
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

	fmt.Printf("%s\n", colorizeJSON(res.RawJSON(), os.Stdout))
	return nil
}

func handleBetaModelsList(ctx context.Context, cmd *cli.Command) error {
	cc := getAPICommandContext(ctx, cmd)
	params := anthropic.BetaModelListParams{}
	res, err := cc.client.Beta.Models.List(
		context.TODO(),
		params,
		option.WithMiddleware(cc.AsMiddleware()),
	)
	if err != nil {
		return err
	}

	fmt.Printf("%s\n", colorizeJSON(res.RawJSON(), os.Stdout))
	return nil
}
