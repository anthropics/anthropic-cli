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
			Name:  "max-tokens-to-sample",
			Usage: "The maximum number of tokens to generate before stopping.\n\nNote that our models may stop _before_ reaching this maximum. This parameter only specifies the absolute maximum number of tokens to generate.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "max_tokens_to_sample",
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "model",
			Usage: "The model that will complete your prompt.\\n\\nSee [models](https://docs.anthropic.com/en/docs/models-overview) for additional details and options.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "model",
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "prompt",
			Usage: "The prompt that you want Claude to complete.\n\nFor proper response generation you will need to format your prompt using alternating `\\n\\nHuman:` and `\\n\\nAssistant:` conversational turns. For example:\n\n```\n\"\\n\\nHuman: {userQuestion}\\n\\nAssistant:\"\n```\n\nSee [prompt validation](https://docs.anthropic.com/en/api/prompt-validation) and our guide to [prompt design](https://docs.anthropic.com/en/docs/intro-to-prompting) for more details.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "prompt",
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "metadata.user_id",
			Usage: "An external identifier for the user who is associated with the request.\n\nThis should be a uuid, hash value, or other opaque identifier. Anthropic may use this id to help detect abuse. Do not include any identifying information such as name, email address, or phone number.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "metadata.user_id",
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "stop-sequences",
			Usage: "Sequences that will cause the model to stop generating.\n\nOur models stop on `\"\\n\\nHuman:\"`, and may include additional built-in stop sequences in the future. By providing the stop_sequences parameter, you may include additional strings that will cause the model to stop generating.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "stop_sequences.#",
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "+stop-sequence",
			Usage: "Sequences that will cause the model to stop generating.\n\nOur models stop on `\"\\n\\nHuman:\"`, and may include additional built-in stop sequences in the future. By providing the stop_sequences parameter, you may include additional strings that will cause the model to stop generating.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "stop_sequences.-1",
			},
		},
		&jsonflag.JSONFloatFlag{
			Name:  "temperature",
			Usage: "Amount of randomness injected into the response.\n\nDefaults to `1.0`. Ranges from `0.0` to `1.0`. Use `temperature` closer to `0.0` for analytical / multiple choice, and closer to `1.0` for creative and generative tasks.\n\nNote that even with `temperature` of `0.0`, the results will not be fully deterministic.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "temperature",
			},
		},
		&jsonflag.JSONIntFlag{
			Name:  "top-k",
			Usage: "Only sample from the top K options for each subsequent token.\n\nUsed to remove \"long tail\" low probability responses. [Learn more technical details here](https://towardsdatascience.com/how-to-sample-from-language-models-682bceb97277).\n\nRecommended for advanced use cases only. You usually only need to use `temperature`.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "top_k",
			},
		},
		&jsonflag.JSONFloatFlag{
			Name:  "top-p",
			Usage: "Use nucleus sampling.\n\nIn nucleus sampling, we compute the cumulative distribution over all the options for each subsequent token in decreasing probability order and cut it off once it reaches a particular probability specified by `top_p`. You should either alter `temperature` or `top_p`, but not both.\n\nRecommended for advanced use cases only. You usually only need to use `temperature`.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "top_p",
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "betas",
			Usage: "Optional header to specify the beta version(s) you want to use.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Header,
				Path: "anthropic-beta.#",
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "+beta",
			Usage: "Optional header to specify the beta version(s) you want to use.",
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
	unusedArgs := cmd.Args().Slice()
	if len(unusedArgs) > 0 {
		return fmt.Errorf("Unexpected extra arguments: %v", unusedArgs)
	}
	params := anthropic.CompletionNewParams{}
	stream := cc.client.Completions.NewStreaming(
		ctx,
		params,
		option.WithMiddleware(cc.AsMiddleware()),
	)
	for stream.Next() {
		fmt.Printf("%s\n", stream.Current().RawJSON())
	}
	return stream.Err()
}
