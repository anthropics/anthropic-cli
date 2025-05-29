// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cmd

import (
	"context"
	"fmt"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/urfave/cli/v3"
)

var completionsCreate = cli.Command{
	Name:  "create",
	Usage: "[Legacy] Create a Text Completion.",
	Flags: []cli.Flag{
		&cli.Int64Flag{
			Name:   "max-tokens-to-sample",
			Action: getAPIFlagAction[int64]("body", "max_tokens_to_sample"),
		},
		&cli.StringFlag{
			Name:   "model",
			Action: getAPIFlagAction[string]("body", "model"),
		},
		&cli.StringFlag{
			Name:   "prompt",
			Action: getAPIFlagAction[string]("body", "prompt"),
		},
		&cli.StringFlag{
			Name:   "metadata.user_id",
			Action: getAPIFlagAction[string]("body", "metadata.user_id"),
		},
		&cli.StringFlag{
			Name:   "stop-sequences",
			Action: getAPIFlagAction[string]("body", "stop_sequences.#"),
		},
		&cli.StringFlag{
			Name:   "+stop_sequence",
			Action: getAPIFlagAction[string]("body", "stop_sequences.-1"),
		},
		&cli.FloatFlag{
			Name:   "temperature",
			Action: getAPIFlagAction[float64]("body", "temperature"),
		},
		&cli.Int64Flag{
			Name:   "top-k",
			Action: getAPIFlagAction[int64]("body", "top_k"),
		},
		&cli.FloatFlag{
			Name:   "top-p",
			Action: getAPIFlagAction[float64]("body", "top_p"),
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
	Action:          handleCompletionsCreate,
	HideHelpCommand: true,
}

func handleCompletionsCreate(ctx context.Context, cmd *cli.Command) error {
	cc := getAPICommandContext(ctx, cmd)
	params := anthropic.CompletionNewParams{}
	stream := cc.client.Completions.NewStreaming(
		context.TODO(),
		params,
		option.WithMiddleware(cc.AsMiddleware()),
		option.WithRequestBody("application/json", cc.body),
	)
	for stream.Next() {
		fmt.Printf("%s\n", stream.Current().RawJSON())
	}
	if err := stream.Err(); err != nil {
		return err
	}
	return nil
}
