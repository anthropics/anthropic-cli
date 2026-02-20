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

var completionsCreate = requestflag.WithInnerFlags(cli.Command{
	Name:    "create",
	Usage:   "[Legacy] Create a Text Completion.",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[int64]{
			Name:     "max-tokens-to-sample",
			Usage:    "The maximum number of tokens to generate before stopping.\n\nNote that our models may stop _before_ reaching this maximum. This parameter only specifies the absolute maximum number of tokens to generate.",
			Required: true,
			BodyPath: "max_tokens_to_sample",
		},
		&requestflag.Flag[string]{
			Name:     "model",
			Usage:    "The model that will complete your prompt.\\n\\nSee [models](https://docs.anthropic.com/en/docs/models-overview) for additional details and options.",
			Required: true,
			BodyPath: "model",
		},
		&requestflag.Flag[string]{
			Name:     "prompt",
			Usage:    "The prompt that you want Claude to complete.\n\nFor proper response generation you will need to format your prompt using alternating `\\n\\nHuman:` and `\\n\\nAssistant:` conversational turns. For example:\n\n```\n\"\\n\\nHuman: {userQuestion}\\n\\nAssistant:\"\n```\n\nSee [prompt validation](https://docs.claude.com/en/api/prompt-validation) and our guide to [prompt design](https://docs.claude.com/en/docs/intro-to-prompting) for more details.",
			Required: true,
			BodyPath: "prompt",
		},
		&requestflag.Flag[map[string]any]{
			Name:     "metadata",
			BodyPath: "metadata",
		},
		&requestflag.Flag[[]string]{
			Name:     "stop-sequence",
			Usage:    "Sequences that will cause the model to stop generating.\n\nOur models stop on `\"\\n\\nHuman:\"`, and may include additional built-in stop sequences in the future. By providing the stop_sequences parameter, you may include additional strings that will cause the model to stop generating.",
			BodyPath: "stop_sequences",
		},
		&requestflag.Flag[bool]{
			Name:     "stream",
			Usage:    "Whether to incrementally stream the response using server-sent events.\n\nSee [streaming](https://docs.claude.com/en/api/streaming) for details.",
			BodyPath: "stream",
		},
		&requestflag.Flag[float64]{
			Name:     "temperature",
			Usage:    "Amount of randomness injected into the response.\n\nDefaults to `1.0`. Ranges from `0.0` to `1.0`. Use `temperature` closer to `0.0` for analytical / multiple choice, and closer to `1.0` for creative and generative tasks.\n\nNote that even with `temperature` of `0.0`, the results will not be fully deterministic.",
			BodyPath: "temperature",
		},
		&requestflag.Flag[int64]{
			Name:     "top-k",
			Usage:    "Only sample from the top K options for each subsequent token.\n\nUsed to remove \"long tail\" low probability responses. [Learn more technical details here](https://towardsdatascience.com/how-to-sample-from-language-models-682bceb97277).\n\nRecommended for advanced use cases only. You usually only need to use `temperature`.",
			BodyPath: "top_k",
		},
		&requestflag.Flag[float64]{
			Name:     "top-p",
			Usage:    "Use nucleus sampling.\n\nIn nucleus sampling, we compute the cumulative distribution over all the options for each subsequent token in decreasing probability order and cut it off once it reaches a particular probability specified by `top_p`. You should either alter `temperature` or `top_p`, but not both.\n\nRecommended for advanced use cases only. You usually only need to use `temperature`.",
			BodyPath: "top_p",
		},
		&requestflag.Flag[[]string]{
			Name:       "beta",
			Usage:      "Optional header to specify the beta version(s) you want to use.",
			HeaderPath: "anthropic-beta",
		},
	},
	Action:          handleCompletionsCreate,
	HideHelpCommand: true,
}, map[string][]requestflag.HasOuterFlag{
	"metadata": {
		&requestflag.InnerFlag[any]{
			Name:       "metadata.user-id",
			Usage:      "An external identifier for the user who is associated with the request.\n\nThis should be a uuid, hash value, or other opaque identifier. Anthropic may use this id to help detect abuse. Do not include any identifying information such as name, email address, or phone number.",
			InnerField: "user_id",
		},
	},
})

func handleCompletionsCreate(ctx context.Context, cmd *cli.Command) error {
	client := anthropic.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()

	if len(unusedArgs) > 0 {
		return fmt.Errorf("Unexpected extra arguments: %v", unusedArgs)
	}

	params := anthropic.CompletionNewParams{}

	options, err := flagOptions(
		cmd,
		apiquery.NestedQueryFormatBrackets,
		apiquery.ArrayQueryFormatComma,
		ApplicationJSON,
		false,
	)
	if err != nil {
		return err
	}

	format := cmd.Root().String("format")
	transform := cmd.Root().String("transform")
	if cmd.Bool("stream") {
		stream := client.Completions.NewStreaming(ctx, params, options...)
		return ShowJSONIterator(os.Stdout, "completions create", stream, format, transform)
	} else {
		var res []byte
		options = append(options, option.WithResponseBodyInto(&res))
		_, err = client.Completions.New(ctx, params, options...)
		if err != nil {
			return err
		}

		obj := gjson.ParseBytes(res)
		return ShowJSON(os.Stdout, "completions create", obj, format, transform)
	}
}
