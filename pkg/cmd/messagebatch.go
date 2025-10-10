// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cmd

import (
	"context"
	"fmt"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/stainless-sdks/anthropic-cli/pkg/jsonflag"
	"github.com/tidwall/gjson"
	"github.com/urfave/cli/v3"
)

var messagesBatchesCreate = cli.Command{
	Name:  "create",
	Usage: "Send a batch of Message creation requests.",
	Flags: []cli.Flag{
		&jsonflag.JSONStringFlag{
			Name:  "requests.custom_id",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.custom_id",
			},
		},
		&jsonflag.JSONIntFlag{
			Name:  "requests.params.max_tokens",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.max_tokens",
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "requests.params.messages.content.text",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.text",
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "requests.params.messages.content.type",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.type",
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "requests.params.messages.content.cache_control.type",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.cache_control.type",
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "requests.params.messages.content.cache_control.ttl",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.cache_control.ttl",
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "requests.params.messages.content.citations.cited_text",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.citations.#.cited_text",
			},
		},
		&jsonflag.JSONIntFlag{
			Name:  "requests.params.messages.content.citations.document_index",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.citations.#.document_index",
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "requests.params.messages.content.citations.document_title",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.citations.#.document_title",
			},
		},
		&jsonflag.JSONIntFlag{
			Name:  "requests.params.messages.content.citations.end_char_index",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.citations.#.end_char_index",
			},
		},
		&jsonflag.JSONIntFlag{
			Name:  "requests.params.messages.content.citations.start_char_index",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.citations.#.start_char_index",
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "requests.params.messages.content.citations.type",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.citations.#.type",
			},
		},
		&jsonflag.JSONIntFlag{
			Name:  "requests.params.messages.content.citations.end_page_number",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.citations.#.end_page_number",
			},
		},
		&jsonflag.JSONIntFlag{
			Name:  "requests.params.messages.content.citations.start_page_number",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.citations.#.start_page_number",
			},
		},
		&jsonflag.JSONIntFlag{
			Name:  "requests.params.messages.content.citations.end_block_index",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.citations.#.end_block_index",
			},
		},
		&jsonflag.JSONIntFlag{
			Name:  "requests.params.messages.content.citations.start_block_index",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.citations.#.start_block_index",
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "requests.params.messages.content.citations.encrypted_index",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.citations.#.encrypted_index",
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "requests.params.messages.content.citations.title",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.citations.#.title",
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "requests.params.messages.content.citations.url",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.citations.#.url",
			},
		},
		&jsonflag.JSONIntFlag{
			Name:  "requests.params.messages.content.citations.search_result_index",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.citations.#.search_result_index",
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "requests.params.messages.content.citations.source",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.citations.#.source",
			},
		},
		&jsonflag.JSONAnyFlag{
			Name:  "requests.params.messages.content.+citation",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind:     jsonflag.Body,
				Path:     "requests.#.params.messages.#.content.#.citations.-1",
				SetValue: map[string]interface{}{},
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "requests.params.messages.content.source.data",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.source.data",
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "requests.params.messages.content.source.media_type",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.source.media_type",
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "requests.params.messages.content.source.type",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.source.type",
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "requests.params.messages.content.source.url",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.source.url",
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "requests.params.messages.content.source.content",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.source.content",
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "requests.params.messages.content.source.content.text",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.source.content.#.text",
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "requests.params.messages.content.source.content.type",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.source.content.#.type",
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "requests.params.messages.content.source.content.cache_control.type",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.source.content.#.cache_control.type",
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "requests.params.messages.content.source.content.cache_control.ttl",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.source.content.#.cache_control.ttl",
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "requests.params.messages.content.source.content.citations.cited_text",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.source.content.#.citations.#.cited_text",
			},
		},
		&jsonflag.JSONIntFlag{
			Name:  "requests.params.messages.content.source.content.citations.document_index",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.source.content.#.citations.#.document_index",
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "requests.params.messages.content.source.content.citations.document_title",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.source.content.#.citations.#.document_title",
			},
		},
		&jsonflag.JSONIntFlag{
			Name:  "requests.params.messages.content.source.content.citations.end_char_index",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.source.content.#.citations.#.end_char_index",
			},
		},
		&jsonflag.JSONIntFlag{
			Name:  "requests.params.messages.content.source.content.citations.start_char_index",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.source.content.#.citations.#.start_char_index",
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "requests.params.messages.content.source.content.citations.type",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.source.content.#.citations.#.type",
			},
		},
		&jsonflag.JSONIntFlag{
			Name:  "requests.params.messages.content.source.content.citations.end_page_number",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.source.content.#.citations.#.end_page_number",
			},
		},
		&jsonflag.JSONIntFlag{
			Name:  "requests.params.messages.content.source.content.citations.start_page_number",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.source.content.#.citations.#.start_page_number",
			},
		},
		&jsonflag.JSONIntFlag{
			Name:  "requests.params.messages.content.source.content.citations.end_block_index",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.source.content.#.citations.#.end_block_index",
			},
		},
		&jsonflag.JSONIntFlag{
			Name:  "requests.params.messages.content.source.content.citations.start_block_index",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.source.content.#.citations.#.start_block_index",
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "requests.params.messages.content.source.content.citations.encrypted_index",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.source.content.#.citations.#.encrypted_index",
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "requests.params.messages.content.source.content.citations.title",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.source.content.#.citations.#.title",
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "requests.params.messages.content.source.content.citations.url",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.source.content.#.citations.#.url",
			},
		},
		&jsonflag.JSONIntFlag{
			Name:  "requests.params.messages.content.source.content.citations.search_result_index",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.source.content.#.citations.#.search_result_index",
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "requests.params.messages.content.source.content.citations.source",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.source.content.#.citations.#.source",
			},
		},
		&jsonflag.JSONAnyFlag{
			Name:  "requests.params.messages.content.source.content.+citation",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind:     jsonflag.Body,
				Path:     "requests.#.params.messages.#.content.#.source.content.#.citations.-1",
				SetValue: map[string]interface{}{},
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "requests.params.messages.content.source.content.source.data",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.source.content.#.source.data",
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "requests.params.messages.content.source.content.source.media_type",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.source.content.#.source.media_type",
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "requests.params.messages.content.source.content.source.type",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.source.content.#.source.type",
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "requests.params.messages.content.source.content.source.url",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.source.content.#.source.url",
			},
		},
		&jsonflag.JSONAnyFlag{
			Name:  "requests.params.messages.content.source.+content",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind:     jsonflag.Body,
				Path:     "requests.#.params.messages.#.content.#.source.content.-1",
				SetValue: map[string]interface{}{},
			},
		},
		&jsonflag.JSONBoolFlag{
			Name:  "requests.params.messages.content.citations.enabled",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind:     jsonflag.Body,
				Path:     "requests.#.params.messages.#.content.#.citations.enabled",
				SetValue: true,
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "requests.params.messages.content.context",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.context",
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "requests.params.messages.content.title",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.title",
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "requests.params.messages.content.content.text",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.content.#.text",
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "requests.params.messages.content.content.type",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.content.#.type",
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "requests.params.messages.content.content.cache_control.type",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.content.#.cache_control.type",
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "requests.params.messages.content.content.cache_control.ttl",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.content.#.cache_control.ttl",
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "requests.params.messages.content.content.citations.cited_text",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.content.#.citations.#.cited_text",
			},
		},
		&jsonflag.JSONIntFlag{
			Name:  "requests.params.messages.content.content.citations.document_index",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.content.#.citations.#.document_index",
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "requests.params.messages.content.content.citations.document_title",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.content.#.citations.#.document_title",
			},
		},
		&jsonflag.JSONIntFlag{
			Name:  "requests.params.messages.content.content.citations.end_char_index",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.content.#.citations.#.end_char_index",
			},
		},
		&jsonflag.JSONIntFlag{
			Name:  "requests.params.messages.content.content.citations.start_char_index",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.content.#.citations.#.start_char_index",
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "requests.params.messages.content.content.citations.type",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.content.#.citations.#.type",
			},
		},
		&jsonflag.JSONIntFlag{
			Name:  "requests.params.messages.content.content.citations.end_page_number",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.content.#.citations.#.end_page_number",
			},
		},
		&jsonflag.JSONIntFlag{
			Name:  "requests.params.messages.content.content.citations.start_page_number",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.content.#.citations.#.start_page_number",
			},
		},
		&jsonflag.JSONIntFlag{
			Name:  "requests.params.messages.content.content.citations.end_block_index",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.content.#.citations.#.end_block_index",
			},
		},
		&jsonflag.JSONIntFlag{
			Name:  "requests.params.messages.content.content.citations.start_block_index",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.content.#.citations.#.start_block_index",
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "requests.params.messages.content.content.citations.encrypted_index",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.content.#.citations.#.encrypted_index",
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "requests.params.messages.content.content.citations.title",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.content.#.citations.#.title",
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "requests.params.messages.content.content.citations.url",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.content.#.citations.#.url",
			},
		},
		&jsonflag.JSONIntFlag{
			Name:  "requests.params.messages.content.content.citations.search_result_index",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.content.#.citations.#.search_result_index",
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "requests.params.messages.content.content.citations.source",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.content.#.citations.#.source",
			},
		},
		&jsonflag.JSONAnyFlag{
			Name:  "requests.params.messages.content.content.+citation",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind:     jsonflag.Body,
				Path:     "requests.#.params.messages.#.content.#.content.#.citations.-1",
				SetValue: map[string]interface{}{},
			},
		},
		&jsonflag.JSONAnyFlag{
			Name:  "requests.params.messages.content.+content",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind:     jsonflag.Body,
				Path:     "requests.#.params.messages.#.content.#.content.-1",
				SetValue: map[string]interface{}{},
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "requests.params.messages.content.source",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.source",
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "requests.params.messages.content.signature",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.signature",
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "requests.params.messages.content.thinking",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.thinking",
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "requests.params.messages.content.data",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.data",
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "requests.params.messages.content.id",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.id",
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "requests.params.messages.content.name",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.name",
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "requests.params.messages.content.tool_use_id",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.tool_use_id",
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "requests.params.messages.content.content.source.data",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.content.#.source.data",
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "requests.params.messages.content.content.source.media_type",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.content.#.source.media_type",
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "requests.params.messages.content.content.source.type",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.content.#.source.type",
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "requests.params.messages.content.content.source.url",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.content.#.source.url",
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "requests.params.messages.content.content.content.text",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.content.#.content.#.text",
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "requests.params.messages.content.content.content.type",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.content.#.content.#.type",
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "requests.params.messages.content.content.content.cache_control.type",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.content.#.content.#.cache_control.type",
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "requests.params.messages.content.content.content.cache_control.ttl",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.content.#.content.#.cache_control.ttl",
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "requests.params.messages.content.content.content.citations.cited_text",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.content.#.content.#.citations.#.cited_text",
			},
		},
		&jsonflag.JSONIntFlag{
			Name:  "requests.params.messages.content.content.content.citations.document_index",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.content.#.content.#.citations.#.document_index",
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "requests.params.messages.content.content.content.citations.document_title",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.content.#.content.#.citations.#.document_title",
			},
		},
		&jsonflag.JSONIntFlag{
			Name:  "requests.params.messages.content.content.content.citations.end_char_index",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.content.#.content.#.citations.#.end_char_index",
			},
		},
		&jsonflag.JSONIntFlag{
			Name:  "requests.params.messages.content.content.content.citations.start_char_index",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.content.#.content.#.citations.#.start_char_index",
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "requests.params.messages.content.content.content.citations.type",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.content.#.content.#.citations.#.type",
			},
		},
		&jsonflag.JSONIntFlag{
			Name:  "requests.params.messages.content.content.content.citations.end_page_number",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.content.#.content.#.citations.#.end_page_number",
			},
		},
		&jsonflag.JSONIntFlag{
			Name:  "requests.params.messages.content.content.content.citations.start_page_number",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.content.#.content.#.citations.#.start_page_number",
			},
		},
		&jsonflag.JSONIntFlag{
			Name:  "requests.params.messages.content.content.content.citations.end_block_index",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.content.#.content.#.citations.#.end_block_index",
			},
		},
		&jsonflag.JSONIntFlag{
			Name:  "requests.params.messages.content.content.content.citations.start_block_index",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.content.#.content.#.citations.#.start_block_index",
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "requests.params.messages.content.content.content.citations.encrypted_index",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.content.#.content.#.citations.#.encrypted_index",
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "requests.params.messages.content.content.content.citations.title",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.content.#.content.#.citations.#.title",
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "requests.params.messages.content.content.content.citations.url",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.content.#.content.#.citations.#.url",
			},
		},
		&jsonflag.JSONIntFlag{
			Name:  "requests.params.messages.content.content.content.citations.search_result_index",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.content.#.content.#.citations.#.search_result_index",
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "requests.params.messages.content.content.content.citations.source",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.content.#.content.#.citations.#.source",
			},
		},
		&jsonflag.JSONAnyFlag{
			Name:  "requests.params.messages.content.content.content.+citation",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind:     jsonflag.Body,
				Path:     "requests.#.params.messages.#.content.#.content.#.content.#.citations.-1",
				SetValue: map[string]interface{}{},
			},
		},
		&jsonflag.JSONAnyFlag{
			Name:  "requests.params.messages.content.content.+content",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind:     jsonflag.Body,
				Path:     "requests.#.params.messages.#.content.#.content.#.content.-1",
				SetValue: map[string]interface{}{},
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "requests.params.messages.content.content.source",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.content.#.source",
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "requests.params.messages.content.content.title",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.content.#.title",
			},
		},
		&jsonflag.JSONBoolFlag{
			Name:  "requests.params.messages.content.content.citations.enabled",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind:     jsonflag.Body,
				Path:     "requests.#.params.messages.#.content.#.content.#.citations.enabled",
				SetValue: true,
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "requests.params.messages.content.content.source.content",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.content.#.source.content",
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "requests.params.messages.content.content.source.content.text",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.content.#.source.content.#.text",
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "requests.params.messages.content.content.source.content.type",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.content.#.source.content.#.type",
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "requests.params.messages.content.content.source.content.cache_control.type",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.content.#.source.content.#.cache_control.type",
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "requests.params.messages.content.content.source.content.cache_control.ttl",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.content.#.source.content.#.cache_control.ttl",
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "requests.params.messages.content.content.source.content.citations.cited_text",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.content.#.source.content.#.citations.#.cited_text",
			},
		},
		&jsonflag.JSONIntFlag{
			Name:  "requests.params.messages.content.content.source.content.citations.document_index",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.content.#.source.content.#.citations.#.document_index",
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "requests.params.messages.content.content.source.content.citations.document_title",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.content.#.source.content.#.citations.#.document_title",
			},
		},
		&jsonflag.JSONIntFlag{
			Name:  "requests.params.messages.content.content.source.content.citations.end_char_index",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.content.#.source.content.#.citations.#.end_char_index",
			},
		},
		&jsonflag.JSONIntFlag{
			Name:  "requests.params.messages.content.content.source.content.citations.start_char_index",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.content.#.source.content.#.citations.#.start_char_index",
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "requests.params.messages.content.content.source.content.citations.type",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.content.#.source.content.#.citations.#.type",
			},
		},
		&jsonflag.JSONIntFlag{
			Name:  "requests.params.messages.content.content.source.content.citations.end_page_number",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.content.#.source.content.#.citations.#.end_page_number",
			},
		},
		&jsonflag.JSONIntFlag{
			Name:  "requests.params.messages.content.content.source.content.citations.start_page_number",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.content.#.source.content.#.citations.#.start_page_number",
			},
		},
		&jsonflag.JSONIntFlag{
			Name:  "requests.params.messages.content.content.source.content.citations.end_block_index",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.content.#.source.content.#.citations.#.end_block_index",
			},
		},
		&jsonflag.JSONIntFlag{
			Name:  "requests.params.messages.content.content.source.content.citations.start_block_index",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.content.#.source.content.#.citations.#.start_block_index",
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "requests.params.messages.content.content.source.content.citations.encrypted_index",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.content.#.source.content.#.citations.#.encrypted_index",
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "requests.params.messages.content.content.source.content.citations.title",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.content.#.source.content.#.citations.#.title",
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "requests.params.messages.content.content.source.content.citations.url",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.content.#.source.content.#.citations.#.url",
			},
		},
		&jsonflag.JSONIntFlag{
			Name:  "requests.params.messages.content.content.source.content.citations.search_result_index",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.content.#.source.content.#.citations.#.search_result_index",
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "requests.params.messages.content.content.source.content.citations.source",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.content.#.source.content.#.citations.#.source",
			},
		},
		&jsonflag.JSONAnyFlag{
			Name:  "requests.params.messages.content.content.source.content.+citation",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind:     jsonflag.Body,
				Path:     "requests.#.params.messages.#.content.#.content.#.source.content.#.citations.-1",
				SetValue: map[string]interface{}{},
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "requests.params.messages.content.content.source.content.source.data",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.content.#.source.content.#.source.data",
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "requests.params.messages.content.content.source.content.source.media_type",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.content.#.source.content.#.source.media_type",
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "requests.params.messages.content.content.source.content.source.type",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.content.#.source.content.#.source.type",
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "requests.params.messages.content.content.source.content.source.url",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.content.#.source.content.#.source.url",
			},
		},
		&jsonflag.JSONAnyFlag{
			Name:  "requests.params.messages.content.content.source.+content",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind:     jsonflag.Body,
				Path:     "requests.#.params.messages.#.content.#.content.#.source.content.-1",
				SetValue: map[string]interface{}{},
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "requests.params.messages.content.content.context",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.content.#.context",
			},
		},
		&jsonflag.JSONBoolFlag{
			Name:  "requests.params.messages.content.is_error",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind:     jsonflag.Body,
				Path:     "requests.#.params.messages.#.content.#.is_error",
				SetValue: true,
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "requests.params.messages.content.content.encrypted_content",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.content.#.encrypted_content",
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "requests.params.messages.content.content.url",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.content.#.url",
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "requests.params.messages.content.content.page_age",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.content.#.page_age",
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "requests.params.messages.content.content.error_code",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.content.error_code",
			},
		},
		&jsonflag.JSONAnyFlag{
			Name:  "requests.params.messages.+content",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind:     jsonflag.Body,
				Path:     "requests.#.params.messages.#.content.-1",
				SetValue: map[string]interface{}{},
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "requests.params.messages.role",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.role",
			},
		},
		&jsonflag.JSONAnyFlag{
			Name:  "requests.params.+message",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind:     jsonflag.Body,
				Path:     "requests.#.params.messages.-1",
				SetValue: map[string]interface{}{},
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "requests.params.model",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.model",
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "requests.params.metadata.user_id",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.metadata.user_id",
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "requests.params.service_tier",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.service_tier",
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "requests.params.stop_sequences",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.stop_sequences.#",
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "requests.params.+stop-sequence",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.stop_sequences.-1",
			},
		},
		&jsonflag.JSONBoolFlag{
			Name:  "requests.params.stream",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind:     jsonflag.Body,
				Path:     "requests.#.params.stream",
				SetValue: true,
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "requests.params.system.text",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.system.#.text",
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "requests.params.system.type",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.system.#.type",
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "requests.params.system.cache_control.type",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.system.#.cache_control.type",
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "requests.params.system.cache_control.ttl",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.system.#.cache_control.ttl",
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "requests.params.system.citations.cited_text",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.system.#.citations.#.cited_text",
			},
		},
		&jsonflag.JSONIntFlag{
			Name:  "requests.params.system.citations.document_index",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.system.#.citations.#.document_index",
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "requests.params.system.citations.document_title",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.system.#.citations.#.document_title",
			},
		},
		&jsonflag.JSONIntFlag{
			Name:  "requests.params.system.citations.end_char_index",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.system.#.citations.#.end_char_index",
			},
		},
		&jsonflag.JSONIntFlag{
			Name:  "requests.params.system.citations.start_char_index",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.system.#.citations.#.start_char_index",
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "requests.params.system.citations.type",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.system.#.citations.#.type",
			},
		},
		&jsonflag.JSONIntFlag{
			Name:  "requests.params.system.citations.end_page_number",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.system.#.citations.#.end_page_number",
			},
		},
		&jsonflag.JSONIntFlag{
			Name:  "requests.params.system.citations.start_page_number",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.system.#.citations.#.start_page_number",
			},
		},
		&jsonflag.JSONIntFlag{
			Name:  "requests.params.system.citations.end_block_index",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.system.#.citations.#.end_block_index",
			},
		},
		&jsonflag.JSONIntFlag{
			Name:  "requests.params.system.citations.start_block_index",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.system.#.citations.#.start_block_index",
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "requests.params.system.citations.encrypted_index",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.system.#.citations.#.encrypted_index",
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "requests.params.system.citations.title",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.system.#.citations.#.title",
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "requests.params.system.citations.url",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.system.#.citations.#.url",
			},
		},
		&jsonflag.JSONIntFlag{
			Name:  "requests.params.system.citations.search_result_index",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.system.#.citations.#.search_result_index",
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "requests.params.system.citations.source",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.system.#.citations.#.source",
			},
		},
		&jsonflag.JSONAnyFlag{
			Name:  "requests.params.system.+citation",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind:     jsonflag.Body,
				Path:     "requests.#.params.system.#.citations.-1",
				SetValue: map[string]interface{}{},
			},
		},
		&jsonflag.JSONAnyFlag{
			Name:  "requests.params.+system",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind:     jsonflag.Body,
				Path:     "requests.#.params.system.-1",
				SetValue: map[string]interface{}{},
			},
		},
		&jsonflag.JSONFloatFlag{
			Name:  "requests.params.temperature",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.temperature",
			},
		},
		&jsonflag.JSONIntFlag{
			Name:  "requests.params.thinking.budget_tokens",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.thinking.budget_tokens",
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "requests.params.thinking.type",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.thinking.type",
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "requests.params.tool_choice.type",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.tool_choice.type",
			},
		},
		&jsonflag.JSONBoolFlag{
			Name:  "requests.params.tool_choice.disable_parallel_tool_use",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind:     jsonflag.Body,
				Path:     "requests.#.params.tool_choice.disable_parallel_tool_use",
				SetValue: true,
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "requests.params.tool_choice.name",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.tool_choice.name",
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "requests.params.tools.name",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.tools.#.name",
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "requests.params.tools.cache_control.type",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.tools.#.cache_control.type",
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "requests.params.tools.cache_control.ttl",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.tools.#.cache_control.ttl",
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "requests.params.tools.description",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.tools.#.description",
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "requests.params.tools.type",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.tools.#.type",
			},
		},
		&jsonflag.JSONIntFlag{
			Name:  "requests.params.tools.max_characters",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.tools.#.max_characters",
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "requests.params.tools.allowed_domains",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.tools.#.allowed_domains.#",
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "requests.params.tools.+allowed-domain",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.tools.#.allowed_domains.-1",
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "requests.params.tools.blocked_domains",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.tools.#.blocked_domains.#",
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "requests.params.tools.+blocked-domain",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.tools.#.blocked_domains.-1",
			},
		},
		&jsonflag.JSONIntFlag{
			Name:  "requests.params.tools.max_uses",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.tools.#.max_uses",
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "requests.params.tools.user_location.type",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.tools.#.user_location.type",
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "requests.params.tools.user_location.city",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.tools.#.user_location.city",
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "requests.params.tools.user_location.country",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.tools.#.user_location.country",
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "requests.params.tools.user_location.region",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.tools.#.user_location.region",
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "requests.params.tools.user_location.timezone",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.tools.#.user_location.timezone",
			},
		},
		&jsonflag.JSONAnyFlag{
			Name:  "requests.params.+tool",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind:     jsonflag.Body,
				Path:     "requests.#.params.tools.-1",
				SetValue: map[string]interface{}{},
			},
		},
		&jsonflag.JSONIntFlag{
			Name:  "requests.params.top_k",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.top_k",
			},
		},
		&jsonflag.JSONFloatFlag{
			Name:  "requests.params.top_p",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.top_p",
			},
		},
		&jsonflag.JSONAnyFlag{
			Name:  "+request",
			Usage: "List of requests for prompt completion. Each is an individual request to create a Message.",
			Config: jsonflag.JSONConfig{
				Kind:     jsonflag.Body,
				Path:     "requests.-1",
				SetValue: map[string]interface{}{},
			},
		},
	},
	Action:          handleMessagesBatchesCreate,
	HideHelpCommand: true,
}

var messagesBatchesRetrieve = cli.Command{
	Name:  "retrieve",
	Usage: "This endpoint is idempotent and can be used to poll for Message Batch\ncompletion. To access the results of a Message Batch, make a request to the\n`results_url` field in the response.",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "message-batch-id",
			Usage: "ID of the Message Batch.",
		},
	},
	Action:          handleMessagesBatchesRetrieve,
	HideHelpCommand: true,
}

var messagesBatchesList = cli.Command{
	Name:  "list",
	Usage: "List all Message Batches within a Workspace. Most recently created batches are\nreturned first.",
	Flags: []cli.Flag{
		&jsonflag.JSONStringFlag{
			Name:  "after-id",
			Usage: "ID of the object to use as a cursor for pagination. When provided, returns the page of results immediately after this object.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Query,
				Path: "after_id",
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "before-id",
			Usage: "ID of the object to use as a cursor for pagination. When provided, returns the page of results immediately before this object.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Query,
				Path: "before_id",
			},
		},
		&jsonflag.JSONIntFlag{
			Name:  "limit",
			Usage: "Number of items to return per page.\n\nDefaults to `20`. Ranges from `1` to `1000`.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Query,
				Path: "limit",
			},
			Value: 20,
		},
	},
	Action:          handleMessagesBatchesList,
	HideHelpCommand: true,
}

var messagesBatchesDelete = cli.Command{
	Name:  "delete",
	Usage: "Delete a Message Batch.",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "message-batch-id",
			Usage: "ID of the Message Batch.",
		},
	},
	Action:          handleMessagesBatchesDelete,
	HideHelpCommand: true,
}

var messagesBatchesCancel = cli.Command{
	Name:  "cancel",
	Usage: "Batches may be canceled any time before processing ends. Once cancellation is\ninitiated, the batch enters a `canceling` state, at which time the system may\ncomplete any in-progress, non-interruptible requests before finalizing\ncancellation.",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "message-batch-id",
			Usage: "ID of the Message Batch.",
		},
	},
	Action:          handleMessagesBatchesCancel,
	HideHelpCommand: true,
}

var messagesBatchesResults = cli.Command{
	Name:  "results",
	Usage: "Streams the results of a Message Batch as a `.jsonl` file.",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "message-batch-id",
			Usage: "ID of the Message Batch.",
		},
	},
	Action:          handleMessagesBatchesResults,
	HideHelpCommand: true,
}

func handleMessagesBatchesCreate(ctx context.Context, cmd *cli.Command) error {
	cc := getAPICommandContext(cmd)
	unusedArgs := cmd.Args().Slice()
	if len(unusedArgs) > 0 {
		return fmt.Errorf("Unexpected extra arguments: %v", unusedArgs)
	}
	params := anthropic.MessageBatchNewParams{}
	var res []byte
	_, err := cc.client.Messages.Batches.New(
		ctx,
		params,
		option.WithMiddleware(cc.AsMiddleware()),
		option.WithResponseBodyInto(&res),
	)
	if err != nil {
		return err
	}

	json := gjson.Parse(string(res))
	format := cmd.Root().String("format")
	transform := cmd.Root().String("transform")
	return ShowJSON("messages:batches create", json, format, transform)
}

func handleMessagesBatchesRetrieve(ctx context.Context, cmd *cli.Command) error {
	cc := getAPICommandContext(cmd)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("message-batch-id") && len(unusedArgs) > 0 {
		cmd.Set("message-batch-id", unusedArgs[0])
		unusedArgs = unusedArgs[1:]
	}
	if len(unusedArgs) > 0 {
		return fmt.Errorf("Unexpected extra arguments: %v", unusedArgs)
	}
	var res []byte
	_, err := cc.client.Messages.Batches.Get(
		ctx,
		cmd.Value("message-batch-id").(string),
		option.WithMiddleware(cc.AsMiddleware()),
		option.WithResponseBodyInto(&res),
	)
	if err != nil {
		return err
	}

	json := gjson.Parse(string(res))
	format := cmd.Root().String("format")
	transform := cmd.Root().String("transform")
	return ShowJSON("messages:batches retrieve", json, format, transform)
}

func handleMessagesBatchesList(ctx context.Context, cmd *cli.Command) error {
	cc := getAPICommandContext(cmd)
	unusedArgs := cmd.Args().Slice()
	if len(unusedArgs) > 0 {
		return fmt.Errorf("Unexpected extra arguments: %v", unusedArgs)
	}
	params := anthropic.MessageBatchListParams{}
	var res []byte
	_, err := cc.client.Messages.Batches.List(
		ctx,
		params,
		option.WithMiddleware(cc.AsMiddleware()),
		option.WithResponseBodyInto(&res),
	)
	if err != nil {
		return err
	}

	json := gjson.Parse(string(res))
	format := cmd.Root().String("format")
	transform := cmd.Root().String("transform")
	return ShowJSON("messages:batches list", json, format, transform)
}

func handleMessagesBatchesDelete(ctx context.Context, cmd *cli.Command) error {
	cc := getAPICommandContext(cmd)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("message-batch-id") && len(unusedArgs) > 0 {
		cmd.Set("message-batch-id", unusedArgs[0])
		unusedArgs = unusedArgs[1:]
	}
	if len(unusedArgs) > 0 {
		return fmt.Errorf("Unexpected extra arguments: %v", unusedArgs)
	}
	var res []byte
	_, err := cc.client.Messages.Batches.Delete(
		ctx,
		cmd.Value("message-batch-id").(string),
		option.WithMiddleware(cc.AsMiddleware()),
		option.WithResponseBodyInto(&res),
	)
	if err != nil {
		return err
	}

	json := gjson.Parse(string(res))
	format := cmd.Root().String("format")
	transform := cmd.Root().String("transform")
	return ShowJSON("messages:batches delete", json, format, transform)
}

func handleMessagesBatchesCancel(ctx context.Context, cmd *cli.Command) error {
	cc := getAPICommandContext(cmd)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("message-batch-id") && len(unusedArgs) > 0 {
		cmd.Set("message-batch-id", unusedArgs[0])
		unusedArgs = unusedArgs[1:]
	}
	if len(unusedArgs) > 0 {
		return fmt.Errorf("Unexpected extra arguments: %v", unusedArgs)
	}
	var res []byte
	_, err := cc.client.Messages.Batches.Cancel(
		ctx,
		cmd.Value("message-batch-id").(string),
		option.WithMiddleware(cc.AsMiddleware()),
		option.WithResponseBodyInto(&res),
	)
	if err != nil {
		return err
	}

	json := gjson.Parse(string(res))
	format := cmd.Root().String("format")
	transform := cmd.Root().String("transform")
	return ShowJSON("messages:batches cancel", json, format, transform)
}

func handleMessagesBatchesResults(ctx context.Context, cmd *cli.Command) error {
	cc := getAPICommandContext(cmd)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("message-batch-id") && len(unusedArgs) > 0 {
		cmd.Set("message-batch-id", unusedArgs[0])
		unusedArgs = unusedArgs[1:]
	}
	if len(unusedArgs) > 0 {
		return fmt.Errorf("Unexpected extra arguments: %v", unusedArgs)
	}
	stream := cc.client.Messages.Batches.ResultsStreaming(
		ctx,
		cmd.Value("message-batch-id").(string),
		option.WithMiddleware(cc.AsMiddleware()),
	)
	for stream.Next() {
		fmt.Printf("%s\n", stream.Current().RawJSON())
	}
	return stream.Err()
}
