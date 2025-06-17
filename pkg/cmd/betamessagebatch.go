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

var betaMessagesBatchesCreate = cli.Command{
	Name:  "create",
	Usage: "Send a batch of Message creation requests.",
	Flags: []cli.Flag{
		&jsonflag.JSONStringFlag{
			Name: "requests.custom_id",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.custom_id",
			},
		},
		&jsonflag.JSONIntFlag{
			Name: "requests.params.max_tokens",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.max_tokens",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "requests.params.messages.content.id",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.id",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "requests.params.messages.content.name",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.name",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "requests.params.messages.content.type",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.type",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "requests.params.messages.content.cache_control.type",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.cache_control.type",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "requests.params.messages.content.cache_control.ttl",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.cache_control.ttl",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "requests.params.messages.content.content.encrypted_content",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.content.#.encrypted_content",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "requests.params.messages.content.content.title",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.content.#.title",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "requests.params.messages.content.content.type",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.content.#.type",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "requests.params.messages.content.content.url",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.content.#.url",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "requests.params.messages.content.content.page_age",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.content.#.page_age",
			},
		},
		&jsonflag.JSONAnyFlag{
			Name: "requests.params.messages.content.+content",
			Config: jsonflag.JSONConfig{
				Kind:     jsonflag.Body,
				Path:     "requests.#.params.messages.#.content.#.content.-1",
				SetValue: map[string]interface{}{},
			},
			Value: map[string]interface{}{},
		},
		&jsonflag.JSONStringFlag{
			Name: "requests.params.messages.content.content.error_code",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.content.error_code",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "requests.params.messages.content.tool_use_id",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.tool_use_id",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "requests.params.messages.content.content.content.file_id",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.content.content.#.file_id",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "requests.params.messages.content.content.content.type",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.content.content.#.type",
			},
		},
		&jsonflag.JSONAnyFlag{
			Name: "requests.params.messages.content.content.+content",
			Config: jsonflag.JSONConfig{
				Kind:     jsonflag.Body,
				Path:     "requests.#.params.messages.#.content.#.content.content.-1",
				SetValue: map[string]interface{}{},
			},
			Value: map[string]interface{}{},
		},
		&jsonflag.JSONIntFlag{
			Name: "requests.params.messages.content.content.return_code",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.content.return_code",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "requests.params.messages.content.content.stderr",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.content.stderr",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "requests.params.messages.content.content.stdout",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.content.stdout",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "requests.params.messages.content.server_name",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.server_name",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "requests.params.messages.content.content",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.content",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "requests.params.messages.content.content.text",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.content.#.text",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "requests.params.messages.content.content.cache_control.type",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.content.#.cache_control.type",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "requests.params.messages.content.content.cache_control.ttl",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.content.#.cache_control.ttl",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "requests.params.messages.content.content.citations.cited_text",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.content.#.citations.#.cited_text",
			},
		},
		&jsonflag.JSONIntFlag{
			Name: "requests.params.messages.content.content.citations.document_index",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.content.#.citations.#.document_index",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "requests.params.messages.content.content.citations.document_title",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.content.#.citations.#.document_title",
			},
		},
		&jsonflag.JSONIntFlag{
			Name: "requests.params.messages.content.content.citations.end_char_index",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.content.#.citations.#.end_char_index",
			},
		},
		&jsonflag.JSONIntFlag{
			Name: "requests.params.messages.content.content.citations.start_char_index",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.content.#.citations.#.start_char_index",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "requests.params.messages.content.content.citations.type",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.content.#.citations.#.type",
			},
		},
		&jsonflag.JSONIntFlag{
			Name: "requests.params.messages.content.content.citations.end_page_number",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.content.#.citations.#.end_page_number",
			},
		},
		&jsonflag.JSONIntFlag{
			Name: "requests.params.messages.content.content.citations.start_page_number",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.content.#.citations.#.start_page_number",
			},
		},
		&jsonflag.JSONIntFlag{
			Name: "requests.params.messages.content.content.citations.end_block_index",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.content.#.citations.#.end_block_index",
			},
		},
		&jsonflag.JSONIntFlag{
			Name: "requests.params.messages.content.content.citations.start_block_index",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.content.#.citations.#.start_block_index",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "requests.params.messages.content.content.citations.encrypted_index",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.content.#.citations.#.encrypted_index",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "requests.params.messages.content.content.citations.title",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.content.#.citations.#.title",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "requests.params.messages.content.content.citations.url",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.content.#.citations.#.url",
			},
		},
		&jsonflag.JSONAnyFlag{
			Name: "requests.params.messages.content.content.+citation",
			Config: jsonflag.JSONConfig{
				Kind:     jsonflag.Body,
				Path:     "requests.#.params.messages.#.content.#.content.#.citations.-1",
				SetValue: map[string]interface{}{},
			},
			Value: map[string]interface{}{},
		},
		&jsonflag.JSONBoolFlag{
			Name: "requests.params.messages.content.is_error",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.is_error",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "requests.params.messages.content.text",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.text",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "requests.params.messages.content.citations.cited_text",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.citations.#.cited_text",
			},
		},
		&jsonflag.JSONIntFlag{
			Name: "requests.params.messages.content.citations.document_index",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.citations.#.document_index",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "requests.params.messages.content.citations.document_title",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.citations.#.document_title",
			},
		},
		&jsonflag.JSONIntFlag{
			Name: "requests.params.messages.content.citations.end_char_index",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.citations.#.end_char_index",
			},
		},
		&jsonflag.JSONIntFlag{
			Name: "requests.params.messages.content.citations.start_char_index",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.citations.#.start_char_index",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "requests.params.messages.content.citations.type",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.citations.#.type",
			},
		},
		&jsonflag.JSONIntFlag{
			Name: "requests.params.messages.content.citations.end_page_number",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.citations.#.end_page_number",
			},
		},
		&jsonflag.JSONIntFlag{
			Name: "requests.params.messages.content.citations.start_page_number",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.citations.#.start_page_number",
			},
		},
		&jsonflag.JSONIntFlag{
			Name: "requests.params.messages.content.citations.end_block_index",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.citations.#.end_block_index",
			},
		},
		&jsonflag.JSONIntFlag{
			Name: "requests.params.messages.content.citations.start_block_index",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.citations.#.start_block_index",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "requests.params.messages.content.citations.encrypted_index",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.citations.#.encrypted_index",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "requests.params.messages.content.citations.title",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.citations.#.title",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "requests.params.messages.content.citations.url",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.citations.#.url",
			},
		},
		&jsonflag.JSONAnyFlag{
			Name: "requests.params.messages.content.+citation",
			Config: jsonflag.JSONConfig{
				Kind:     jsonflag.Body,
				Path:     "requests.#.params.messages.#.content.#.citations.-1",
				SetValue: map[string]interface{}{},
			},
			Value: map[string]interface{}{},
		},
		&jsonflag.JSONStringFlag{
			Name: "requests.params.messages.content.source.data",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.source.data",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "requests.params.messages.content.source.media_type",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.source.media_type",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "requests.params.messages.content.source.type",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.source.type",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "requests.params.messages.content.source.url",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.source.url",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "requests.params.messages.content.source.file_id",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.source.file_id",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "requests.params.messages.content.content.source.data",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.content.#.source.data",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "requests.params.messages.content.content.source.media_type",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.content.#.source.media_type",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "requests.params.messages.content.content.source.type",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.content.#.source.type",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "requests.params.messages.content.content.source.url",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.content.#.source.url",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "requests.params.messages.content.content.source.file_id",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.content.#.source.file_id",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "requests.params.messages.content.source.content",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.source.content",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "requests.params.messages.content.source.content.text",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.source.content.#.text",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "requests.params.messages.content.source.content.type",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.source.content.#.type",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "requests.params.messages.content.source.content.cache_control.type",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.source.content.#.cache_control.type",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "requests.params.messages.content.source.content.cache_control.ttl",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.source.content.#.cache_control.ttl",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "requests.params.messages.content.source.content.citations.cited_text",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.source.content.#.citations.#.cited_text",
			},
		},
		&jsonflag.JSONIntFlag{
			Name: "requests.params.messages.content.source.content.citations.document_index",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.source.content.#.citations.#.document_index",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "requests.params.messages.content.source.content.citations.document_title",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.source.content.#.citations.#.document_title",
			},
		},
		&jsonflag.JSONIntFlag{
			Name: "requests.params.messages.content.source.content.citations.end_char_index",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.source.content.#.citations.#.end_char_index",
			},
		},
		&jsonflag.JSONIntFlag{
			Name: "requests.params.messages.content.source.content.citations.start_char_index",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.source.content.#.citations.#.start_char_index",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "requests.params.messages.content.source.content.citations.type",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.source.content.#.citations.#.type",
			},
		},
		&jsonflag.JSONIntFlag{
			Name: "requests.params.messages.content.source.content.citations.end_page_number",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.source.content.#.citations.#.end_page_number",
			},
		},
		&jsonflag.JSONIntFlag{
			Name: "requests.params.messages.content.source.content.citations.start_page_number",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.source.content.#.citations.#.start_page_number",
			},
		},
		&jsonflag.JSONIntFlag{
			Name: "requests.params.messages.content.source.content.citations.end_block_index",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.source.content.#.citations.#.end_block_index",
			},
		},
		&jsonflag.JSONIntFlag{
			Name: "requests.params.messages.content.source.content.citations.start_block_index",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.source.content.#.citations.#.start_block_index",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "requests.params.messages.content.source.content.citations.encrypted_index",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.source.content.#.citations.#.encrypted_index",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "requests.params.messages.content.source.content.citations.title",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.source.content.#.citations.#.title",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "requests.params.messages.content.source.content.citations.url",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.source.content.#.citations.#.url",
			},
		},
		&jsonflag.JSONAnyFlag{
			Name: "requests.params.messages.content.source.content.+citation",
			Config: jsonflag.JSONConfig{
				Kind:     jsonflag.Body,
				Path:     "requests.#.params.messages.#.content.#.source.content.#.citations.-1",
				SetValue: map[string]interface{}{},
			},
			Value: map[string]interface{}{},
		},
		&jsonflag.JSONStringFlag{
			Name: "requests.params.messages.content.source.content.source.data",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.source.content.#.source.data",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "requests.params.messages.content.source.content.source.media_type",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.source.content.#.source.media_type",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "requests.params.messages.content.source.content.source.type",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.source.content.#.source.type",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "requests.params.messages.content.source.content.source.url",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.source.content.#.source.url",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "requests.params.messages.content.source.content.source.file_id",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.source.content.#.source.file_id",
			},
		},
		&jsonflag.JSONAnyFlag{
			Name: "requests.params.messages.content.source.+content",
			Config: jsonflag.JSONConfig{
				Kind:     jsonflag.Body,
				Path:     "requests.#.params.messages.#.content.#.source.content.-1",
				SetValue: map[string]interface{}{},
			},
			Value: map[string]interface{}{},
		},
		&jsonflag.JSONBoolFlag{
			Name: "requests.params.messages.content.citations.enabled",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.citations.enabled",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "requests.params.messages.content.context",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.context",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "requests.params.messages.content.title",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.title",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "requests.params.messages.content.signature",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.signature",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "requests.params.messages.content.thinking",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.thinking",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "requests.params.messages.content.data",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.data",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "requests.params.messages.content.file_id",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.content.#.file_id",
			},
		},
		&jsonflag.JSONAnyFlag{
			Name: "requests.params.messages.+content",
			Config: jsonflag.JSONConfig{
				Kind:     jsonflag.Body,
				Path:     "requests.#.params.messages.#.content.-1",
				SetValue: map[string]interface{}{},
			},
			Value: map[string]interface{}{},
		},
		&jsonflag.JSONStringFlag{
			Name: "requests.params.messages.role",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.messages.#.role",
			},
		},
		&jsonflag.JSONAnyFlag{
			Name: "requests.params.+message",
			Config: jsonflag.JSONConfig{
				Kind:     jsonflag.Body,
				Path:     "requests.#.params.messages.-1",
				SetValue: map[string]interface{}{},
			},
			Value: map[string]interface{}{},
		},
		&jsonflag.JSONStringFlag{
			Name: "requests.params.model",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.model",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "requests.params.container",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.container",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "requests.params.mcp_servers.name",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.mcp_servers.#.name",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "requests.params.mcp_servers.type",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.mcp_servers.#.type",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "requests.params.mcp_servers.url",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.mcp_servers.#.url",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "requests.params.mcp_servers.authorization_token",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.mcp_servers.#.authorization_token",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "requests.params.mcp_servers.tool_configuration.allowed_tools",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.mcp_servers.#.tool_configuration.allowed_tools.#",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "requests.params.mcp_servers.tool_configuration.+allowed_tool",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.mcp_servers.#.tool_configuration.allowed_tools.-1",
			},
		},
		&jsonflag.JSONBoolFlag{
			Name: "requests.params.mcp_servers.tool_configuration.enabled",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.mcp_servers.#.tool_configuration.enabled",
			},
		},
		&jsonflag.JSONAnyFlag{
			Name: "requests.params.+mcp_server",
			Config: jsonflag.JSONConfig{
				Kind:     jsonflag.Body,
				Path:     "requests.#.params.mcp_servers.-1",
				SetValue: map[string]interface{}{},
			},
			Value: map[string]interface{}{},
		},
		&jsonflag.JSONStringFlag{
			Name: "requests.params.metadata.user_id",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.metadata.user_id",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "requests.params.service_tier",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.service_tier",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "requests.params.stop_sequences",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.stop_sequences.#",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "requests.params.+stop_sequence",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.stop_sequences.-1",
			},
		},
		&jsonflag.JSONBoolFlag{
			Name: "requests.params.stream",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.stream",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "requests.params.system.text",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.system.#.text",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "requests.params.system.type",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.system.#.type",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "requests.params.system.cache_control.type",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.system.#.cache_control.type",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "requests.params.system.cache_control.ttl",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.system.#.cache_control.ttl",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "requests.params.system.citations.cited_text",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.system.#.citations.#.cited_text",
			},
		},
		&jsonflag.JSONIntFlag{
			Name: "requests.params.system.citations.document_index",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.system.#.citations.#.document_index",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "requests.params.system.citations.document_title",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.system.#.citations.#.document_title",
			},
		},
		&jsonflag.JSONIntFlag{
			Name: "requests.params.system.citations.end_char_index",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.system.#.citations.#.end_char_index",
			},
		},
		&jsonflag.JSONIntFlag{
			Name: "requests.params.system.citations.start_char_index",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.system.#.citations.#.start_char_index",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "requests.params.system.citations.type",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.system.#.citations.#.type",
			},
		},
		&jsonflag.JSONIntFlag{
			Name: "requests.params.system.citations.end_page_number",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.system.#.citations.#.end_page_number",
			},
		},
		&jsonflag.JSONIntFlag{
			Name: "requests.params.system.citations.start_page_number",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.system.#.citations.#.start_page_number",
			},
		},
		&jsonflag.JSONIntFlag{
			Name: "requests.params.system.citations.end_block_index",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.system.#.citations.#.end_block_index",
			},
		},
		&jsonflag.JSONIntFlag{
			Name: "requests.params.system.citations.start_block_index",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.system.#.citations.#.start_block_index",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "requests.params.system.citations.encrypted_index",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.system.#.citations.#.encrypted_index",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "requests.params.system.citations.title",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.system.#.citations.#.title",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "requests.params.system.citations.url",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.system.#.citations.#.url",
			},
		},
		&jsonflag.JSONAnyFlag{
			Name: "requests.params.system.+citation",
			Config: jsonflag.JSONConfig{
				Kind:     jsonflag.Body,
				Path:     "requests.#.params.system.#.citations.-1",
				SetValue: map[string]interface{}{},
			},
			Value: map[string]interface{}{},
		},
		&jsonflag.JSONAnyFlag{
			Name: "requests.params.+system",
			Config: jsonflag.JSONConfig{
				Kind:     jsonflag.Body,
				Path:     "requests.#.params.system.-1",
				SetValue: map[string]interface{}{},
			},
			Value: map[string]interface{}{},
		},
		&jsonflag.JSONFloatFlag{
			Name: "requests.params.temperature",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.temperature",
			},
		},
		&jsonflag.JSONIntFlag{
			Name: "requests.params.thinking.budget_tokens",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.thinking.budget_tokens",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "requests.params.thinking.type",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.thinking.type",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "requests.params.tool_choice.type",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.tool_choice.type",
			},
		},
		&jsonflag.JSONBoolFlag{
			Name: "requests.params.tool_choice.disable_parallel_tool_use",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.tool_choice.disable_parallel_tool_use",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "requests.params.tool_choice.name",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.tool_choice.name",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "requests.params.tools.input_schema.type",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.tools.#.input_schema.type",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "requests.params.tools.input_schema.required",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.tools.#.input_schema.required.#",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "requests.params.tools.input_schema.+required",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.tools.#.input_schema.required.-1",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "requests.params.tools.name",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.tools.#.name",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "requests.params.tools.cache_control.type",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.tools.#.cache_control.type",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "requests.params.tools.cache_control.ttl",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.tools.#.cache_control.ttl",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "requests.params.tools.description",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.tools.#.description",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "requests.params.tools.type",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.tools.#.type",
			},
		},
		&jsonflag.JSONIntFlag{
			Name: "requests.params.tools.display_height_px",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.tools.#.display_height_px",
			},
		},
		&jsonflag.JSONIntFlag{
			Name: "requests.params.tools.display_width_px",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.tools.#.display_width_px",
			},
		},
		&jsonflag.JSONIntFlag{
			Name: "requests.params.tools.display_number",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.tools.#.display_number",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "requests.params.tools.allowed_domains",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.tools.#.allowed_domains.#",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "requests.params.tools.+allowed_domain",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.tools.#.allowed_domains.-1",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "requests.params.tools.blocked_domains",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.tools.#.blocked_domains.#",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "requests.params.tools.+blocked_domain",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.tools.#.blocked_domains.-1",
			},
		},
		&jsonflag.JSONIntFlag{
			Name: "requests.params.tools.max_uses",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.tools.#.max_uses",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "requests.params.tools.user_location.type",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.tools.#.user_location.type",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "requests.params.tools.user_location.city",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.tools.#.user_location.city",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "requests.params.tools.user_location.country",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.tools.#.user_location.country",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "requests.params.tools.user_location.region",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.tools.#.user_location.region",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "requests.params.tools.user_location.timezone",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.tools.#.user_location.timezone",
			},
		},
		&jsonflag.JSONAnyFlag{
			Name: "requests.params.+tool",
			Config: jsonflag.JSONConfig{
				Kind:     jsonflag.Body,
				Path:     "requests.#.params.tools.-1",
				SetValue: map[string]interface{}{},
			},
			Value: map[string]interface{}{},
		},
		&jsonflag.JSONIntFlag{
			Name: "requests.params.top_k",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.top_k",
			},
		},
		&jsonflag.JSONFloatFlag{
			Name: "requests.params.top_p",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "requests.#.params.top_p",
			},
		},
		&jsonflag.JSONAnyFlag{
			Name: "+request",
			Config: jsonflag.JSONConfig{
				Kind:     jsonflag.Body,
				Path:     "requests.-1",
				SetValue: map[string]interface{}{},
			},
			Value: map[string]interface{}{},
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
	Action:          handleBetaMessagesBatchesCreate,
	HideHelpCommand: true,
}

var betaMessagesBatchesRetrieve = cli.Command{
	Name:  "retrieve",
	Usage: "This endpoint is idempotent and can be used to poll for Message Batch\ncompletion. To access the results of a Message Batch, make a request to the\n`results_url` field in the response.",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name: "message-batch-id",
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
	Action:          handleBetaMessagesBatchesRetrieve,
	HideHelpCommand: true,
}

var betaMessagesBatchesList = cli.Command{
	Name:  "list",
	Usage: "List all Message Batches within a Workspace. Most recently created batches are\nreturned first.",
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
	Action:          handleBetaMessagesBatchesList,
	HideHelpCommand: true,
}

var betaMessagesBatchesDelete = cli.Command{
	Name:  "delete",
	Usage: "Delete a Message Batch.",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name: "message-batch-id",
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
	Action:          handleBetaMessagesBatchesDelete,
	HideHelpCommand: true,
}

var betaMessagesBatchesCancel = cli.Command{
	Name:  "cancel",
	Usage: "Batches may be canceled any time before processing ends. Once cancellation is\ninitiated, the batch enters a `canceling` state, at which time the system may\ncomplete any in-progress, non-interruptible requests before finalizing\ncancellation.",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name: "message-batch-id",
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
	Action:          handleBetaMessagesBatchesCancel,
	HideHelpCommand: true,
}

var betaMessagesBatchesResults = cli.Command{
	Name:  "results",
	Usage: "Streams the results of a Message Batch as a `.jsonl` file.",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name: "message-batch-id",
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
	Action:          handleBetaMessagesBatchesResults,
	HideHelpCommand: true,
}

func handleBetaMessagesBatchesCreate(ctx context.Context, cmd *cli.Command) error {
	cc := getAPICommandContext(cmd)
	params := anthropic.BetaMessageBatchNewParams{}
	res, err := cc.client.Beta.Messages.Batches.New(
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

func handleBetaMessagesBatchesRetrieve(ctx context.Context, cmd *cli.Command) error {
	cc := getAPICommandContext(cmd)
	params := anthropic.BetaMessageBatchGetParams{}
	res, err := cc.client.Beta.Messages.Batches.Get(
		context.TODO(),
		cmd.Value("message-batch-id").(string),
		params,
		option.WithMiddleware(cc.AsMiddleware()),
	)
	if err != nil {
		return err
	}

	fmt.Printf("%s\n", ColorizeJSON(res.RawJSON(), os.Stdout))
	return nil
}

func handleBetaMessagesBatchesList(ctx context.Context, cmd *cli.Command) error {
	cc := getAPICommandContext(cmd)
	params := anthropic.BetaMessageBatchListParams{}
	res, err := cc.client.Beta.Messages.Batches.List(
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

func handleBetaMessagesBatchesDelete(ctx context.Context, cmd *cli.Command) error {
	cc := getAPICommandContext(cmd)
	params := anthropic.BetaMessageBatchDeleteParams{}
	res, err := cc.client.Beta.Messages.Batches.Delete(
		context.TODO(),
		cmd.Value("message-batch-id").(string),
		params,
		option.WithMiddleware(cc.AsMiddleware()),
	)
	if err != nil {
		return err
	}

	fmt.Printf("%s\n", ColorizeJSON(res.RawJSON(), os.Stdout))
	return nil
}

func handleBetaMessagesBatchesCancel(ctx context.Context, cmd *cli.Command) error {
	cc := getAPICommandContext(cmd)
	params := anthropic.BetaMessageBatchCancelParams{}
	res, err := cc.client.Beta.Messages.Batches.Cancel(
		context.TODO(),
		cmd.Value("message-batch-id").(string),
		params,
		option.WithMiddleware(cc.AsMiddleware()),
	)
	if err != nil {
		return err
	}

	fmt.Printf("%s\n", ColorizeJSON(res.RawJSON(), os.Stdout))
	return nil
}

func handleBetaMessagesBatchesResults(ctx context.Context, cmd *cli.Command) error {
	cc := getAPICommandContext(cmd)
	params := anthropic.BetaMessageBatchResultsParams{}
	stream := cc.client.Beta.Messages.Batches.ResultsStreaming(
		context.TODO(),
		cmd.Value("message-batch-id").(string),
		params,
		option.WithMiddleware(cc.AsMiddleware()),
	)
	for stream.Next() {
		fmt.Printf("%s\n", stream.Current().RawJSON())
	}
	if err := stream.Err(); err != nil {
		return err
	}
	return nil
}
