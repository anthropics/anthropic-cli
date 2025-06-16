// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cmd

import (
	"context"
	"fmt"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/stainless-sdks/anthropic-cli/pkg/jsonflag"
	"github.com/urfave/cli/v3"
)

var completionsCreate = cli.Command{
	Name:  "create",
	Usage: "[Legacy] Create a Text Completion.",
	Flags: []cli.Flag{
		&jsonflag.JSONIntFlag{
			Name: "max-tokens-to-sample",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "max_tokens_to_sample",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "model",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "model",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "prompt",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "prompt",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "metadata.user_id",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "metadata.user_id",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "stop-sequences",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "stop_sequences.#",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "+stop_sequence",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "stop_sequences.-1",
			},
		},
		&jsonflag.JSONFloatFlag{
			Name: "temperature",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "temperature",
			},
		},
		&jsonflag.JSONIntFlag{
			Name: "top-k",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "top_k",
			},
		},
		&jsonflag.JSONFloatFlag{
			Name: "top-p",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "top_p",
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
	Action:          handleCompletionsCreate,
	HideHelpCommand: true,
}

func handleCompletionsCreate(ctx context.Context, cmd *cli.Command) error {
	cc := getAPICommandContext(cmd)
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
