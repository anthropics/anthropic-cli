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

var messagesCreate = cli.Command{
	Name:  "create",
	Usage: "Send a structured list of input messages with text and/or image content, and the\nmodel will generate the next message in the conversation.",
	Flags: []cli.Flag{
		&jsonflag.JSONIntFlag{
			Name: "max-tokens",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "max_tokens",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "messages.content.id",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "messages.#.content.#.id",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "messages.content.name",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "messages.#.content.#.name",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "messages.content.type",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "messages.#.content.#.type",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "messages.content.cache_control.type",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "messages.#.content.#.cache_control.type",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "messages.content.content.encrypted_content",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "messages.#.content.#.content.#.encrypted_content",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "messages.content.content.title",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "messages.#.content.#.content.#.title",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "messages.content.content.type",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "messages.#.content.#.content.#.type",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "messages.content.content.url",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "messages.#.content.#.content.#.url",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "messages.content.content.page_age",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "messages.#.content.#.content.#.page_age",
			},
		},
		&jsonflag.JSONAnyFlag{
			Name: "messages.content.+content",
			Config: jsonflag.JSONConfig{
				Kind:     jsonflag.Body,
				Path:     "messages.#.content.#.content.-1",
				SetValue: map[string]interface{}{},
			},
			Value: map[string]interface{}{},
		},
		&jsonflag.JSONStringFlag{
			Name: "messages.content.content.error_code",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "messages.#.content.#.content.error_code",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "messages.content.tool_use_id",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "messages.#.content.#.tool_use_id",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "messages.content.text",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "messages.#.content.#.text",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "messages.content.citations.cited_text",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "messages.#.content.#.citations.#.cited_text",
			},
		},
		&jsonflag.JSONIntFlag{
			Name: "messages.content.citations.document_index",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "messages.#.content.#.citations.#.document_index",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "messages.content.citations.document_title",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "messages.#.content.#.citations.#.document_title",
			},
		},
		&jsonflag.JSONIntFlag{
			Name: "messages.content.citations.end_char_index",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "messages.#.content.#.citations.#.end_char_index",
			},
		},
		&jsonflag.JSONIntFlag{
			Name: "messages.content.citations.start_char_index",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "messages.#.content.#.citations.#.start_char_index",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "messages.content.citations.type",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "messages.#.content.#.citations.#.type",
			},
		},
		&jsonflag.JSONIntFlag{
			Name: "messages.content.citations.end_page_number",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "messages.#.content.#.citations.#.end_page_number",
			},
		},
		&jsonflag.JSONIntFlag{
			Name: "messages.content.citations.start_page_number",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "messages.#.content.#.citations.#.start_page_number",
			},
		},
		&jsonflag.JSONIntFlag{
			Name: "messages.content.citations.end_block_index",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "messages.#.content.#.citations.#.end_block_index",
			},
		},
		&jsonflag.JSONIntFlag{
			Name: "messages.content.citations.start_block_index",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "messages.#.content.#.citations.#.start_block_index",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "messages.content.citations.encrypted_index",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "messages.#.content.#.citations.#.encrypted_index",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "messages.content.citations.title",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "messages.#.content.#.citations.#.title",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "messages.content.citations.url",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "messages.#.content.#.citations.#.url",
			},
		},
		&jsonflag.JSONAnyFlag{
			Name: "messages.content.+citation",
			Config: jsonflag.JSONConfig{
				Kind:     jsonflag.Body,
				Path:     "messages.#.content.#.citations.-1",
				SetValue: map[string]interface{}{},
			},
			Value: map[string]interface{}{},
		},
		&jsonflag.JSONStringFlag{
			Name: "messages.content.source.data",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "messages.#.content.#.source.data",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "messages.content.source.media_type",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "messages.#.content.#.source.media_type",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "messages.content.source.type",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "messages.#.content.#.source.type",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "messages.content.source.url",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "messages.#.content.#.source.url",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "messages.content.content.text",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "messages.#.content.#.content.#.text",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "messages.content.content.cache_control.type",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "messages.#.content.#.content.#.cache_control.type",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "messages.content.content.citations.cited_text",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "messages.#.content.#.content.#.citations.#.cited_text",
			},
		},
		&jsonflag.JSONIntFlag{
			Name: "messages.content.content.citations.document_index",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "messages.#.content.#.content.#.citations.#.document_index",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "messages.content.content.citations.document_title",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "messages.#.content.#.content.#.citations.#.document_title",
			},
		},
		&jsonflag.JSONIntFlag{
			Name: "messages.content.content.citations.end_char_index",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "messages.#.content.#.content.#.citations.#.end_char_index",
			},
		},
		&jsonflag.JSONIntFlag{
			Name: "messages.content.content.citations.start_char_index",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "messages.#.content.#.content.#.citations.#.start_char_index",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "messages.content.content.citations.type",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "messages.#.content.#.content.#.citations.#.type",
			},
		},
		&jsonflag.JSONIntFlag{
			Name: "messages.content.content.citations.end_page_number",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "messages.#.content.#.content.#.citations.#.end_page_number",
			},
		},
		&jsonflag.JSONIntFlag{
			Name: "messages.content.content.citations.start_page_number",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "messages.#.content.#.content.#.citations.#.start_page_number",
			},
		},
		&jsonflag.JSONIntFlag{
			Name: "messages.content.content.citations.end_block_index",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "messages.#.content.#.content.#.citations.#.end_block_index",
			},
		},
		&jsonflag.JSONIntFlag{
			Name: "messages.content.content.citations.start_block_index",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "messages.#.content.#.content.#.citations.#.start_block_index",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "messages.content.content.citations.encrypted_index",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "messages.#.content.#.content.#.citations.#.encrypted_index",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "messages.content.content.citations.title",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "messages.#.content.#.content.#.citations.#.title",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "messages.content.content.citations.url",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "messages.#.content.#.content.#.citations.#.url",
			},
		},
		&jsonflag.JSONAnyFlag{
			Name: "messages.content.content.+citation",
			Config: jsonflag.JSONConfig{
				Kind:     jsonflag.Body,
				Path:     "messages.#.content.#.content.#.citations.-1",
				SetValue: map[string]interface{}{},
			},
			Value: map[string]interface{}{},
		},
		&jsonflag.JSONStringFlag{
			Name: "messages.content.content.source.data",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "messages.#.content.#.content.#.source.data",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "messages.content.content.source.media_type",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "messages.#.content.#.content.#.source.media_type",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "messages.content.content.source.type",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "messages.#.content.#.content.#.source.type",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "messages.content.content.source.url",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "messages.#.content.#.content.#.source.url",
			},
		},
		&jsonflag.JSONBoolFlag{
			Name: "messages.content.is_error",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "messages.#.content.#.is_error",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "messages.content.source.content",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "messages.#.content.#.source.content",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "messages.content.source.content.text",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "messages.#.content.#.source.content.#.text",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "messages.content.source.content.type",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "messages.#.content.#.source.content.#.type",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "messages.content.source.content.cache_control.type",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "messages.#.content.#.source.content.#.cache_control.type",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "messages.content.source.content.citations.cited_text",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "messages.#.content.#.source.content.#.citations.#.cited_text",
			},
		},
		&jsonflag.JSONIntFlag{
			Name: "messages.content.source.content.citations.document_index",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "messages.#.content.#.source.content.#.citations.#.document_index",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "messages.content.source.content.citations.document_title",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "messages.#.content.#.source.content.#.citations.#.document_title",
			},
		},
		&jsonflag.JSONIntFlag{
			Name: "messages.content.source.content.citations.end_char_index",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "messages.#.content.#.source.content.#.citations.#.end_char_index",
			},
		},
		&jsonflag.JSONIntFlag{
			Name: "messages.content.source.content.citations.start_char_index",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "messages.#.content.#.source.content.#.citations.#.start_char_index",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "messages.content.source.content.citations.type",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "messages.#.content.#.source.content.#.citations.#.type",
			},
		},
		&jsonflag.JSONIntFlag{
			Name: "messages.content.source.content.citations.end_page_number",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "messages.#.content.#.source.content.#.citations.#.end_page_number",
			},
		},
		&jsonflag.JSONIntFlag{
			Name: "messages.content.source.content.citations.start_page_number",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "messages.#.content.#.source.content.#.citations.#.start_page_number",
			},
		},
		&jsonflag.JSONIntFlag{
			Name: "messages.content.source.content.citations.end_block_index",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "messages.#.content.#.source.content.#.citations.#.end_block_index",
			},
		},
		&jsonflag.JSONIntFlag{
			Name: "messages.content.source.content.citations.start_block_index",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "messages.#.content.#.source.content.#.citations.#.start_block_index",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "messages.content.source.content.citations.encrypted_index",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "messages.#.content.#.source.content.#.citations.#.encrypted_index",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "messages.content.source.content.citations.title",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "messages.#.content.#.source.content.#.citations.#.title",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "messages.content.source.content.citations.url",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "messages.#.content.#.source.content.#.citations.#.url",
			},
		},
		&jsonflag.JSONAnyFlag{
			Name: "messages.content.source.content.+citation",
			Config: jsonflag.JSONConfig{
				Kind:     jsonflag.Body,
				Path:     "messages.#.content.#.source.content.#.citations.-1",
				SetValue: map[string]interface{}{},
			},
			Value: map[string]interface{}{},
		},
		&jsonflag.JSONStringFlag{
			Name: "messages.content.source.content.source.data",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "messages.#.content.#.source.content.#.source.data",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "messages.content.source.content.source.media_type",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "messages.#.content.#.source.content.#.source.media_type",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "messages.content.source.content.source.type",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "messages.#.content.#.source.content.#.source.type",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "messages.content.source.content.source.url",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "messages.#.content.#.source.content.#.source.url",
			},
		},
		&jsonflag.JSONAnyFlag{
			Name: "messages.content.source.+content",
			Config: jsonflag.JSONConfig{
				Kind:     jsonflag.Body,
				Path:     "messages.#.content.#.source.content.-1",
				SetValue: map[string]interface{}{},
			},
			Value: map[string]interface{}{},
		},
		&jsonflag.JSONBoolFlag{
			Name: "messages.content.citations.enabled",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "messages.#.content.#.citations.enabled",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "messages.content.context",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "messages.#.content.#.context",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "messages.content.title",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "messages.#.content.#.title",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "messages.content.signature",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "messages.#.content.#.signature",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "messages.content.thinking",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "messages.#.content.#.thinking",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "messages.content.data",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "messages.#.content.#.data",
			},
		},
		&jsonflag.JSONAnyFlag{
			Name: "messages.+content",
			Config: jsonflag.JSONConfig{
				Kind:     jsonflag.Body,
				Path:     "messages.#.content.-1",
				SetValue: map[string]interface{}{},
			},
			Value: map[string]interface{}{},
		},
		&jsonflag.JSONStringFlag{
			Name: "messages.role",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "messages.#.role",
			},
		},
		&jsonflag.JSONAnyFlag{
			Name: "+message",
			Config: jsonflag.JSONConfig{
				Kind:     jsonflag.Body,
				Path:     "messages.-1",
				SetValue: map[string]interface{}{},
			},
			Value: map[string]interface{}{},
		},
		&jsonflag.JSONStringFlag{
			Name: "model",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "model",
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
			Name: "service-tier",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "service_tier",
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
		&jsonflag.JSONStringFlag{
			Name: "system.text",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "system.#.text",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "system.type",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "system.#.type",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "system.cache_control.type",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "system.#.cache_control.type",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "system.citations.cited_text",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "system.#.citations.#.cited_text",
			},
		},
		&jsonflag.JSONIntFlag{
			Name: "system.citations.document_index",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "system.#.citations.#.document_index",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "system.citations.document_title",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "system.#.citations.#.document_title",
			},
		},
		&jsonflag.JSONIntFlag{
			Name: "system.citations.end_char_index",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "system.#.citations.#.end_char_index",
			},
		},
		&jsonflag.JSONIntFlag{
			Name: "system.citations.start_char_index",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "system.#.citations.#.start_char_index",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "system.citations.type",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "system.#.citations.#.type",
			},
		},
		&jsonflag.JSONIntFlag{
			Name: "system.citations.end_page_number",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "system.#.citations.#.end_page_number",
			},
		},
		&jsonflag.JSONIntFlag{
			Name: "system.citations.start_page_number",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "system.#.citations.#.start_page_number",
			},
		},
		&jsonflag.JSONIntFlag{
			Name: "system.citations.end_block_index",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "system.#.citations.#.end_block_index",
			},
		},
		&jsonflag.JSONIntFlag{
			Name: "system.citations.start_block_index",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "system.#.citations.#.start_block_index",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "system.citations.encrypted_index",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "system.#.citations.#.encrypted_index",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "system.citations.title",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "system.#.citations.#.title",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "system.citations.url",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "system.#.citations.#.url",
			},
		},
		&jsonflag.JSONAnyFlag{
			Name: "system.+citation",
			Config: jsonflag.JSONConfig{
				Kind:     jsonflag.Body,
				Path:     "system.#.citations.-1",
				SetValue: map[string]interface{}{},
			},
			Value: map[string]interface{}{},
		},
		&jsonflag.JSONAnyFlag{
			Name: "+system",
			Config: jsonflag.JSONConfig{
				Kind:     jsonflag.Body,
				Path:     "system.-1",
				SetValue: map[string]interface{}{},
			},
			Value: map[string]interface{}{},
		},
		&jsonflag.JSONFloatFlag{
			Name: "temperature",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "temperature",
			},
		},
		&jsonflag.JSONIntFlag{
			Name: "thinking.budget_tokens",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "thinking.budget_tokens",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "thinking.type",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "thinking.type",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "tool-choice.type",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "tool_choice.type",
			},
		},
		&jsonflag.JSONBoolFlag{
			Name: "tool-choice.disable_parallel_tool_use",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "tool_choice.disable_parallel_tool_use",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "tool-choice.name",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "tool_choice.name",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "tools.name",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "tools.#.name",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "tools.cache_control.type",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "tools.#.cache_control.type",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "tools.description",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "tools.#.description",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "tools.type",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "tools.#.type",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "tools.allowed_domains",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "tools.#.allowed_domains.#",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "tools.+allowed_domain",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "tools.#.allowed_domains.-1",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "tools.blocked_domains",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "tools.#.blocked_domains.#",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "tools.+blocked_domain",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "tools.#.blocked_domains.-1",
			},
		},
		&jsonflag.JSONIntFlag{
			Name: "tools.max_uses",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "tools.#.max_uses",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "tools.user_location.type",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "tools.#.user_location.type",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "tools.user_location.city",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "tools.#.user_location.city",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "tools.user_location.country",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "tools.#.user_location.country",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "tools.user_location.region",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "tools.#.user_location.region",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "tools.user_location.timezone",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "tools.#.user_location.timezone",
			},
		},
		&jsonflag.JSONAnyFlag{
			Name: "+tool",
			Config: jsonflag.JSONConfig{
				Kind:     jsonflag.Body,
				Path:     "tools.-1",
				SetValue: map[string]interface{}{},
			},
			Value: map[string]interface{}{},
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
	},
	Action:          handleMessagesCreate,
	HideHelpCommand: true,
}

var messagesCountTokens = cli.Command{
	Name:  "count_tokens",
	Usage: "Count the number of tokens in a Message.",
	Flags: []cli.Flag{
		&jsonflag.JSONStringFlag{
			Name: "messages.content.id",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "messages.#.content.#.id",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "messages.content.name",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "messages.#.content.#.name",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "messages.content.type",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "messages.#.content.#.type",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "messages.content.cache_control.type",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "messages.#.content.#.cache_control.type",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "messages.content.content.encrypted_content",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "messages.#.content.#.content.#.encrypted_content",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "messages.content.content.title",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "messages.#.content.#.content.#.title",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "messages.content.content.type",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "messages.#.content.#.content.#.type",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "messages.content.content.url",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "messages.#.content.#.content.#.url",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "messages.content.content.page_age",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "messages.#.content.#.content.#.page_age",
			},
		},
		&jsonflag.JSONAnyFlag{
			Name: "messages.content.+content",
			Config: jsonflag.JSONConfig{
				Kind:     jsonflag.Body,
				Path:     "messages.#.content.#.content.-1",
				SetValue: map[string]interface{}{},
			},
			Value: map[string]interface{}{},
		},
		&jsonflag.JSONStringFlag{
			Name: "messages.content.content.error_code",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "messages.#.content.#.content.error_code",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "messages.content.tool_use_id",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "messages.#.content.#.tool_use_id",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "messages.content.text",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "messages.#.content.#.text",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "messages.content.citations.cited_text",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "messages.#.content.#.citations.#.cited_text",
			},
		},
		&jsonflag.JSONIntFlag{
			Name: "messages.content.citations.document_index",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "messages.#.content.#.citations.#.document_index",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "messages.content.citations.document_title",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "messages.#.content.#.citations.#.document_title",
			},
		},
		&jsonflag.JSONIntFlag{
			Name: "messages.content.citations.end_char_index",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "messages.#.content.#.citations.#.end_char_index",
			},
		},
		&jsonflag.JSONIntFlag{
			Name: "messages.content.citations.start_char_index",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "messages.#.content.#.citations.#.start_char_index",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "messages.content.citations.type",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "messages.#.content.#.citations.#.type",
			},
		},
		&jsonflag.JSONIntFlag{
			Name: "messages.content.citations.end_page_number",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "messages.#.content.#.citations.#.end_page_number",
			},
		},
		&jsonflag.JSONIntFlag{
			Name: "messages.content.citations.start_page_number",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "messages.#.content.#.citations.#.start_page_number",
			},
		},
		&jsonflag.JSONIntFlag{
			Name: "messages.content.citations.end_block_index",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "messages.#.content.#.citations.#.end_block_index",
			},
		},
		&jsonflag.JSONIntFlag{
			Name: "messages.content.citations.start_block_index",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "messages.#.content.#.citations.#.start_block_index",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "messages.content.citations.encrypted_index",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "messages.#.content.#.citations.#.encrypted_index",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "messages.content.citations.title",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "messages.#.content.#.citations.#.title",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "messages.content.citations.url",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "messages.#.content.#.citations.#.url",
			},
		},
		&jsonflag.JSONAnyFlag{
			Name: "messages.content.+citation",
			Config: jsonflag.JSONConfig{
				Kind:     jsonflag.Body,
				Path:     "messages.#.content.#.citations.-1",
				SetValue: map[string]interface{}{},
			},
			Value: map[string]interface{}{},
		},
		&jsonflag.JSONStringFlag{
			Name: "messages.content.source.data",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "messages.#.content.#.source.data",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "messages.content.source.media_type",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "messages.#.content.#.source.media_type",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "messages.content.source.type",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "messages.#.content.#.source.type",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "messages.content.source.url",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "messages.#.content.#.source.url",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "messages.content.content.text",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "messages.#.content.#.content.#.text",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "messages.content.content.cache_control.type",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "messages.#.content.#.content.#.cache_control.type",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "messages.content.content.citations.cited_text",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "messages.#.content.#.content.#.citations.#.cited_text",
			},
		},
		&jsonflag.JSONIntFlag{
			Name: "messages.content.content.citations.document_index",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "messages.#.content.#.content.#.citations.#.document_index",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "messages.content.content.citations.document_title",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "messages.#.content.#.content.#.citations.#.document_title",
			},
		},
		&jsonflag.JSONIntFlag{
			Name: "messages.content.content.citations.end_char_index",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "messages.#.content.#.content.#.citations.#.end_char_index",
			},
		},
		&jsonflag.JSONIntFlag{
			Name: "messages.content.content.citations.start_char_index",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "messages.#.content.#.content.#.citations.#.start_char_index",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "messages.content.content.citations.type",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "messages.#.content.#.content.#.citations.#.type",
			},
		},
		&jsonflag.JSONIntFlag{
			Name: "messages.content.content.citations.end_page_number",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "messages.#.content.#.content.#.citations.#.end_page_number",
			},
		},
		&jsonflag.JSONIntFlag{
			Name: "messages.content.content.citations.start_page_number",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "messages.#.content.#.content.#.citations.#.start_page_number",
			},
		},
		&jsonflag.JSONIntFlag{
			Name: "messages.content.content.citations.end_block_index",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "messages.#.content.#.content.#.citations.#.end_block_index",
			},
		},
		&jsonflag.JSONIntFlag{
			Name: "messages.content.content.citations.start_block_index",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "messages.#.content.#.content.#.citations.#.start_block_index",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "messages.content.content.citations.encrypted_index",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "messages.#.content.#.content.#.citations.#.encrypted_index",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "messages.content.content.citations.title",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "messages.#.content.#.content.#.citations.#.title",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "messages.content.content.citations.url",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "messages.#.content.#.content.#.citations.#.url",
			},
		},
		&jsonflag.JSONAnyFlag{
			Name: "messages.content.content.+citation",
			Config: jsonflag.JSONConfig{
				Kind:     jsonflag.Body,
				Path:     "messages.#.content.#.content.#.citations.-1",
				SetValue: map[string]interface{}{},
			},
			Value: map[string]interface{}{},
		},
		&jsonflag.JSONStringFlag{
			Name: "messages.content.content.source.data",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "messages.#.content.#.content.#.source.data",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "messages.content.content.source.media_type",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "messages.#.content.#.content.#.source.media_type",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "messages.content.content.source.type",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "messages.#.content.#.content.#.source.type",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "messages.content.content.source.url",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "messages.#.content.#.content.#.source.url",
			},
		},
		&jsonflag.JSONBoolFlag{
			Name: "messages.content.is_error",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "messages.#.content.#.is_error",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "messages.content.source.content",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "messages.#.content.#.source.content",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "messages.content.source.content.text",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "messages.#.content.#.source.content.#.text",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "messages.content.source.content.type",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "messages.#.content.#.source.content.#.type",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "messages.content.source.content.cache_control.type",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "messages.#.content.#.source.content.#.cache_control.type",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "messages.content.source.content.citations.cited_text",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "messages.#.content.#.source.content.#.citations.#.cited_text",
			},
		},
		&jsonflag.JSONIntFlag{
			Name: "messages.content.source.content.citations.document_index",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "messages.#.content.#.source.content.#.citations.#.document_index",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "messages.content.source.content.citations.document_title",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "messages.#.content.#.source.content.#.citations.#.document_title",
			},
		},
		&jsonflag.JSONIntFlag{
			Name: "messages.content.source.content.citations.end_char_index",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "messages.#.content.#.source.content.#.citations.#.end_char_index",
			},
		},
		&jsonflag.JSONIntFlag{
			Name: "messages.content.source.content.citations.start_char_index",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "messages.#.content.#.source.content.#.citations.#.start_char_index",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "messages.content.source.content.citations.type",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "messages.#.content.#.source.content.#.citations.#.type",
			},
		},
		&jsonflag.JSONIntFlag{
			Name: "messages.content.source.content.citations.end_page_number",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "messages.#.content.#.source.content.#.citations.#.end_page_number",
			},
		},
		&jsonflag.JSONIntFlag{
			Name: "messages.content.source.content.citations.start_page_number",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "messages.#.content.#.source.content.#.citations.#.start_page_number",
			},
		},
		&jsonflag.JSONIntFlag{
			Name: "messages.content.source.content.citations.end_block_index",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "messages.#.content.#.source.content.#.citations.#.end_block_index",
			},
		},
		&jsonflag.JSONIntFlag{
			Name: "messages.content.source.content.citations.start_block_index",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "messages.#.content.#.source.content.#.citations.#.start_block_index",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "messages.content.source.content.citations.encrypted_index",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "messages.#.content.#.source.content.#.citations.#.encrypted_index",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "messages.content.source.content.citations.title",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "messages.#.content.#.source.content.#.citations.#.title",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "messages.content.source.content.citations.url",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "messages.#.content.#.source.content.#.citations.#.url",
			},
		},
		&jsonflag.JSONAnyFlag{
			Name: "messages.content.source.content.+citation",
			Config: jsonflag.JSONConfig{
				Kind:     jsonflag.Body,
				Path:     "messages.#.content.#.source.content.#.citations.-1",
				SetValue: map[string]interface{}{},
			},
			Value: map[string]interface{}{},
		},
		&jsonflag.JSONStringFlag{
			Name: "messages.content.source.content.source.data",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "messages.#.content.#.source.content.#.source.data",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "messages.content.source.content.source.media_type",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "messages.#.content.#.source.content.#.source.media_type",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "messages.content.source.content.source.type",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "messages.#.content.#.source.content.#.source.type",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "messages.content.source.content.source.url",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "messages.#.content.#.source.content.#.source.url",
			},
		},
		&jsonflag.JSONAnyFlag{
			Name: "messages.content.source.+content",
			Config: jsonflag.JSONConfig{
				Kind:     jsonflag.Body,
				Path:     "messages.#.content.#.source.content.-1",
				SetValue: map[string]interface{}{},
			},
			Value: map[string]interface{}{},
		},
		&jsonflag.JSONBoolFlag{
			Name: "messages.content.citations.enabled",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "messages.#.content.#.citations.enabled",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "messages.content.context",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "messages.#.content.#.context",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "messages.content.title",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "messages.#.content.#.title",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "messages.content.signature",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "messages.#.content.#.signature",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "messages.content.thinking",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "messages.#.content.#.thinking",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "messages.content.data",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "messages.#.content.#.data",
			},
		},
		&jsonflag.JSONAnyFlag{
			Name: "messages.+content",
			Config: jsonflag.JSONConfig{
				Kind:     jsonflag.Body,
				Path:     "messages.#.content.-1",
				SetValue: map[string]interface{}{},
			},
			Value: map[string]interface{}{},
		},
		&jsonflag.JSONStringFlag{
			Name: "messages.role",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "messages.#.role",
			},
		},
		&jsonflag.JSONAnyFlag{
			Name: "+message",
			Config: jsonflag.JSONConfig{
				Kind:     jsonflag.Body,
				Path:     "messages.-1",
				SetValue: map[string]interface{}{},
			},
			Value: map[string]interface{}{},
		},
		&jsonflag.JSONStringFlag{
			Name: "model",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "model",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "system",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "system",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "system.text",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "system.#.text",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "system.type",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "system.#.type",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "system.cache_control.type",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "system.#.cache_control.type",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "system.citations.cited_text",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "system.#.citations.#.cited_text",
			},
		},
		&jsonflag.JSONIntFlag{
			Name: "system.citations.document_index",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "system.#.citations.#.document_index",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "system.citations.document_title",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "system.#.citations.#.document_title",
			},
		},
		&jsonflag.JSONIntFlag{
			Name: "system.citations.end_char_index",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "system.#.citations.#.end_char_index",
			},
		},
		&jsonflag.JSONIntFlag{
			Name: "system.citations.start_char_index",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "system.#.citations.#.start_char_index",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "system.citations.type",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "system.#.citations.#.type",
			},
		},
		&jsonflag.JSONIntFlag{
			Name: "system.citations.end_page_number",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "system.#.citations.#.end_page_number",
			},
		},
		&jsonflag.JSONIntFlag{
			Name: "system.citations.start_page_number",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "system.#.citations.#.start_page_number",
			},
		},
		&jsonflag.JSONIntFlag{
			Name: "system.citations.end_block_index",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "system.#.citations.#.end_block_index",
			},
		},
		&jsonflag.JSONIntFlag{
			Name: "system.citations.start_block_index",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "system.#.citations.#.start_block_index",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "system.citations.encrypted_index",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "system.#.citations.#.encrypted_index",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "system.citations.title",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "system.#.citations.#.title",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "system.citations.url",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "system.#.citations.#.url",
			},
		},
		&jsonflag.JSONAnyFlag{
			Name: "system.+citation",
			Config: jsonflag.JSONConfig{
				Kind:     jsonflag.Body,
				Path:     "system.#.citations.-1",
				SetValue: map[string]interface{}{},
			},
			Value: map[string]interface{}{},
		},
		&jsonflag.JSONAnyFlag{
			Name: "+system",
			Config: jsonflag.JSONConfig{
				Kind:     jsonflag.Body,
				Path:     "system.-1",
				SetValue: map[string]interface{}{},
			},
			Value: map[string]interface{}{},
		},
		&jsonflag.JSONIntFlag{
			Name: "thinking.budget_tokens",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "thinking.budget_tokens",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "thinking.type",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "thinking.type",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "tool-choice.type",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "tool_choice.type",
			},
		},
		&jsonflag.JSONBoolFlag{
			Name: "tool-choice.disable_parallel_tool_use",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "tool_choice.disable_parallel_tool_use",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "tool-choice.name",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "tool_choice.name",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "tools.name",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "tools.#.name",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "tools.cache_control.type",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "tools.#.cache_control.type",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "tools.description",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "tools.#.description",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "tools.type",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "tools.#.type",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "tools.allowed_domains",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "tools.#.allowed_domains.#",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "tools.+allowed_domain",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "tools.#.allowed_domains.-1",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "tools.blocked_domains",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "tools.#.blocked_domains.#",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "tools.+blocked_domain",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "tools.#.blocked_domains.-1",
			},
		},
		&jsonflag.JSONIntFlag{
			Name: "tools.max_uses",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "tools.#.max_uses",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "tools.user_location.type",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "tools.#.user_location.type",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "tools.user_location.city",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "tools.#.user_location.city",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "tools.user_location.country",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "tools.#.user_location.country",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "tools.user_location.region",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "tools.#.user_location.region",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "tools.user_location.timezone",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "tools.#.user_location.timezone",
			},
		},
		&jsonflag.JSONAnyFlag{
			Name: "+tool",
			Config: jsonflag.JSONConfig{
				Kind:     jsonflag.Body,
				Path:     "tools.-1",
				SetValue: map[string]interface{}{},
			},
			Value: map[string]interface{}{},
		},
	},
	Action:          handleMessagesCountTokens,
	HideHelpCommand: true,
}

func handleMessagesCreate(ctx context.Context, cmd *cli.Command) error {
	cc := getAPICommandContext(cmd)
	params := anthropic.MessageNewParams{}
	stream := cc.client.Messages.NewStreaming(
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

func handleMessagesCountTokens(ctx context.Context, cmd *cli.Command) error {
	cc := getAPICommandContext(cmd)
	params := anthropic.MessageCountTokensParams{}
	res, err := cc.client.Messages.CountTokens(
		context.TODO(),
		params,
		option.WithMiddleware(cc.AsMiddleware()),
		option.WithRequestBody("application/json", cc.body),
	)
	if err != nil {
		return err
	}

	fmt.Printf("%s\n", colorizeJSON(res.RawJSON(), os.Stdout))
	return nil
}
