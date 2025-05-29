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

var messagesBatchesCreate = cli.Command{
	Name:  "create",
	Usage: "Send a batch of Message creation requests.",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:   "requests.custom_id",
			Action: getAPIFlagAction[string]("body", "requests.#.custom_id"),
		},
		&cli.Int64Flag{
			Name:   "requests.params.max_tokens",
			Action: getAPIFlagAction[int64]("body", "requests.#.params.max_tokens"),
		},
		&cli.StringFlag{
			Name:   "requests.params.messages.content.id",
			Action: getAPIFlagAction[string]("body", "requests.#.params.messages.#.content.#.id"),
		},
		&cli.StringFlag{
			Name:   "requests.params.messages.content.name",
			Action: getAPIFlagAction[string]("body", "requests.#.params.messages.#.content.#.name"),
		},
		&cli.StringFlag{
			Name:   "requests.params.messages.content.type",
			Action: getAPIFlagAction[string]("body", "requests.#.params.messages.#.content.#.type"),
		},
		&cli.StringFlag{
			Name:   "requests.params.messages.content.cache_control.type",
			Action: getAPIFlagAction[string]("body", "requests.#.params.messages.#.content.#.cache_control.type"),
		},
		&cli.StringFlag{
			Name:   "requests.params.messages.content.content.encrypted_content",
			Action: getAPIFlagAction[string]("body", "requests.#.params.messages.#.content.#.content.#.encrypted_content"),
		},
		&cli.StringFlag{
			Name:   "requests.params.messages.content.content.title",
			Action: getAPIFlagAction[string]("body", "requests.#.params.messages.#.content.#.content.#.title"),
		},
		&cli.StringFlag{
			Name:   "requests.params.messages.content.content.type",
			Action: getAPIFlagAction[string]("body", "requests.#.params.messages.#.content.#.content.#.type"),
		},
		&cli.StringFlag{
			Name:   "requests.params.messages.content.content.url",
			Action: getAPIFlagAction[string]("body", "requests.#.params.messages.#.content.#.content.#.url"),
		},
		&cli.StringFlag{
			Name:   "requests.params.messages.content.content.page_age",
			Action: getAPIFlagAction[string]("body", "requests.#.params.messages.#.content.#.content.#.page_age"),
		},
		&cli.BoolFlag{
			Name:   "requests.params.messages.content.+content",
			Action: getAPIFlagActionWithValue[bool]("body", "requests.#.params.messages.#.content.#.content.-1", map[string]interface{}{}),
		},
		&cli.StringFlag{
			Name:   "requests.params.messages.content.content.error_code",
			Action: getAPIFlagAction[string]("body", "requests.#.params.messages.#.content.#.content.error_code"),
		},
		&cli.StringFlag{
			Name:   "requests.params.messages.content.tool_use_id",
			Action: getAPIFlagAction[string]("body", "requests.#.params.messages.#.content.#.tool_use_id"),
		},
		&cli.StringFlag{
			Name:   "requests.params.messages.content.text",
			Action: getAPIFlagAction[string]("body", "requests.#.params.messages.#.content.#.text"),
		},
		&cli.StringFlag{
			Name:   "requests.params.messages.content.citations.cited_text",
			Action: getAPIFlagAction[string]("body", "requests.#.params.messages.#.content.#.citations.#.cited_text"),
		},
		&cli.Int64Flag{
			Name:   "requests.params.messages.content.citations.document_index",
			Action: getAPIFlagAction[int64]("body", "requests.#.params.messages.#.content.#.citations.#.document_index"),
		},
		&cli.StringFlag{
			Name:   "requests.params.messages.content.citations.document_title",
			Action: getAPIFlagAction[string]("body", "requests.#.params.messages.#.content.#.citations.#.document_title"),
		},
		&cli.Int64Flag{
			Name:   "requests.params.messages.content.citations.end_char_index",
			Action: getAPIFlagAction[int64]("body", "requests.#.params.messages.#.content.#.citations.#.end_char_index"),
		},
		&cli.Int64Flag{
			Name:   "requests.params.messages.content.citations.start_char_index",
			Action: getAPIFlagAction[int64]("body", "requests.#.params.messages.#.content.#.citations.#.start_char_index"),
		},
		&cli.StringFlag{
			Name:   "requests.params.messages.content.citations.type",
			Action: getAPIFlagAction[string]("body", "requests.#.params.messages.#.content.#.citations.#.type"),
		},
		&cli.Int64Flag{
			Name:   "requests.params.messages.content.citations.end_page_number",
			Action: getAPIFlagAction[int64]("body", "requests.#.params.messages.#.content.#.citations.#.end_page_number"),
		},
		&cli.Int64Flag{
			Name:   "requests.params.messages.content.citations.start_page_number",
			Action: getAPIFlagAction[int64]("body", "requests.#.params.messages.#.content.#.citations.#.start_page_number"),
		},
		&cli.Int64Flag{
			Name:   "requests.params.messages.content.citations.end_block_index",
			Action: getAPIFlagAction[int64]("body", "requests.#.params.messages.#.content.#.citations.#.end_block_index"),
		},
		&cli.Int64Flag{
			Name:   "requests.params.messages.content.citations.start_block_index",
			Action: getAPIFlagAction[int64]("body", "requests.#.params.messages.#.content.#.citations.#.start_block_index"),
		},
		&cli.StringFlag{
			Name:   "requests.params.messages.content.citations.encrypted_index",
			Action: getAPIFlagAction[string]("body", "requests.#.params.messages.#.content.#.citations.#.encrypted_index"),
		},
		&cli.StringFlag{
			Name:   "requests.params.messages.content.citations.title",
			Action: getAPIFlagAction[string]("body", "requests.#.params.messages.#.content.#.citations.#.title"),
		},
		&cli.StringFlag{
			Name:   "requests.params.messages.content.citations.url",
			Action: getAPIFlagAction[string]("body", "requests.#.params.messages.#.content.#.citations.#.url"),
		},
		&cli.BoolFlag{
			Name:   "requests.params.messages.content.+citation",
			Action: getAPIFlagActionWithValue[bool]("body", "requests.#.params.messages.#.content.#.citations.-1", map[string]interface{}{}),
		},
		&cli.StringFlag{
			Name:   "requests.params.messages.content.source.data",
			Action: getAPIFlagAction[string]("body", "requests.#.params.messages.#.content.#.source.data"),
		},
		&cli.StringFlag{
			Name:   "requests.params.messages.content.source.media_type",
			Action: getAPIFlagAction[string]("body", "requests.#.params.messages.#.content.#.source.media_type"),
		},
		&cli.StringFlag{
			Name:   "requests.params.messages.content.source.type",
			Action: getAPIFlagAction[string]("body", "requests.#.params.messages.#.content.#.source.type"),
		},
		&cli.StringFlag{
			Name:   "requests.params.messages.content.source.url",
			Action: getAPIFlagAction[string]("body", "requests.#.params.messages.#.content.#.source.url"),
		},
		&cli.StringFlag{
			Name:   "requests.params.messages.content.content.text",
			Action: getAPIFlagAction[string]("body", "requests.#.params.messages.#.content.#.content.#.text"),
		},
		&cli.StringFlag{
			Name:   "requests.params.messages.content.content.cache_control.type",
			Action: getAPIFlagAction[string]("body", "requests.#.params.messages.#.content.#.content.#.cache_control.type"),
		},
		&cli.StringFlag{
			Name:   "requests.params.messages.content.content.citations.cited_text",
			Action: getAPIFlagAction[string]("body", "requests.#.params.messages.#.content.#.content.#.citations.#.cited_text"),
		},
		&cli.Int64Flag{
			Name:   "requests.params.messages.content.content.citations.document_index",
			Action: getAPIFlagAction[int64]("body", "requests.#.params.messages.#.content.#.content.#.citations.#.document_index"),
		},
		&cli.StringFlag{
			Name:   "requests.params.messages.content.content.citations.document_title",
			Action: getAPIFlagAction[string]("body", "requests.#.params.messages.#.content.#.content.#.citations.#.document_title"),
		},
		&cli.Int64Flag{
			Name:   "requests.params.messages.content.content.citations.end_char_index",
			Action: getAPIFlagAction[int64]("body", "requests.#.params.messages.#.content.#.content.#.citations.#.end_char_index"),
		},
		&cli.Int64Flag{
			Name:   "requests.params.messages.content.content.citations.start_char_index",
			Action: getAPIFlagAction[int64]("body", "requests.#.params.messages.#.content.#.content.#.citations.#.start_char_index"),
		},
		&cli.StringFlag{
			Name:   "requests.params.messages.content.content.citations.type",
			Action: getAPIFlagAction[string]("body", "requests.#.params.messages.#.content.#.content.#.citations.#.type"),
		},
		&cli.Int64Flag{
			Name:   "requests.params.messages.content.content.citations.end_page_number",
			Action: getAPIFlagAction[int64]("body", "requests.#.params.messages.#.content.#.content.#.citations.#.end_page_number"),
		},
		&cli.Int64Flag{
			Name:   "requests.params.messages.content.content.citations.start_page_number",
			Action: getAPIFlagAction[int64]("body", "requests.#.params.messages.#.content.#.content.#.citations.#.start_page_number"),
		},
		&cli.Int64Flag{
			Name:   "requests.params.messages.content.content.citations.end_block_index",
			Action: getAPIFlagAction[int64]("body", "requests.#.params.messages.#.content.#.content.#.citations.#.end_block_index"),
		},
		&cli.Int64Flag{
			Name:   "requests.params.messages.content.content.citations.start_block_index",
			Action: getAPIFlagAction[int64]("body", "requests.#.params.messages.#.content.#.content.#.citations.#.start_block_index"),
		},
		&cli.StringFlag{
			Name:   "requests.params.messages.content.content.citations.encrypted_index",
			Action: getAPIFlagAction[string]("body", "requests.#.params.messages.#.content.#.content.#.citations.#.encrypted_index"),
		},
		&cli.StringFlag{
			Name:   "requests.params.messages.content.content.citations.title",
			Action: getAPIFlagAction[string]("body", "requests.#.params.messages.#.content.#.content.#.citations.#.title"),
		},
		&cli.StringFlag{
			Name:   "requests.params.messages.content.content.citations.url",
			Action: getAPIFlagAction[string]("body", "requests.#.params.messages.#.content.#.content.#.citations.#.url"),
		},
		&cli.BoolFlag{
			Name:   "requests.params.messages.content.content.+citation",
			Action: getAPIFlagActionWithValue[bool]("body", "requests.#.params.messages.#.content.#.content.#.citations.-1", map[string]interface{}{}),
		},
		&cli.StringFlag{
			Name:   "requests.params.messages.content.content.source.data",
			Action: getAPIFlagAction[string]("body", "requests.#.params.messages.#.content.#.content.#.source.data"),
		},
		&cli.StringFlag{
			Name:   "requests.params.messages.content.content.source.media_type",
			Action: getAPIFlagAction[string]("body", "requests.#.params.messages.#.content.#.content.#.source.media_type"),
		},
		&cli.StringFlag{
			Name:   "requests.params.messages.content.content.source.type",
			Action: getAPIFlagAction[string]("body", "requests.#.params.messages.#.content.#.content.#.source.type"),
		},
		&cli.StringFlag{
			Name:   "requests.params.messages.content.content.source.url",
			Action: getAPIFlagAction[string]("body", "requests.#.params.messages.#.content.#.content.#.source.url"),
		},
		&cli.BoolFlag{
			Name:   "requests.params.messages.content.is_error",
			Action: getAPIFlagAction[bool]("body", "requests.#.params.messages.#.content.#.is_error"),
		},
		&cli.StringFlag{
			Name:   "requests.params.messages.content.source.content",
			Action: getAPIFlagAction[string]("body", "requests.#.params.messages.#.content.#.source.content"),
		},
		&cli.StringFlag{
			Name:   "requests.params.messages.content.source.content.text",
			Action: getAPIFlagAction[string]("body", "requests.#.params.messages.#.content.#.source.content.#.text"),
		},
		&cli.StringFlag{
			Name:   "requests.params.messages.content.source.content.type",
			Action: getAPIFlagAction[string]("body", "requests.#.params.messages.#.content.#.source.content.#.type"),
		},
		&cli.StringFlag{
			Name:   "requests.params.messages.content.source.content.cache_control.type",
			Action: getAPIFlagAction[string]("body", "requests.#.params.messages.#.content.#.source.content.#.cache_control.type"),
		},
		&cli.StringFlag{
			Name:   "requests.params.messages.content.source.content.citations.cited_text",
			Action: getAPIFlagAction[string]("body", "requests.#.params.messages.#.content.#.source.content.#.citations.#.cited_text"),
		},
		&cli.Int64Flag{
			Name:   "requests.params.messages.content.source.content.citations.document_index",
			Action: getAPIFlagAction[int64]("body", "requests.#.params.messages.#.content.#.source.content.#.citations.#.document_index"),
		},
		&cli.StringFlag{
			Name:   "requests.params.messages.content.source.content.citations.document_title",
			Action: getAPIFlagAction[string]("body", "requests.#.params.messages.#.content.#.source.content.#.citations.#.document_title"),
		},
		&cli.Int64Flag{
			Name:   "requests.params.messages.content.source.content.citations.end_char_index",
			Action: getAPIFlagAction[int64]("body", "requests.#.params.messages.#.content.#.source.content.#.citations.#.end_char_index"),
		},
		&cli.Int64Flag{
			Name:   "requests.params.messages.content.source.content.citations.start_char_index",
			Action: getAPIFlagAction[int64]("body", "requests.#.params.messages.#.content.#.source.content.#.citations.#.start_char_index"),
		},
		&cli.StringFlag{
			Name:   "requests.params.messages.content.source.content.citations.type",
			Action: getAPIFlagAction[string]("body", "requests.#.params.messages.#.content.#.source.content.#.citations.#.type"),
		},
		&cli.Int64Flag{
			Name:   "requests.params.messages.content.source.content.citations.end_page_number",
			Action: getAPIFlagAction[int64]("body", "requests.#.params.messages.#.content.#.source.content.#.citations.#.end_page_number"),
		},
		&cli.Int64Flag{
			Name:   "requests.params.messages.content.source.content.citations.start_page_number",
			Action: getAPIFlagAction[int64]("body", "requests.#.params.messages.#.content.#.source.content.#.citations.#.start_page_number"),
		},
		&cli.Int64Flag{
			Name:   "requests.params.messages.content.source.content.citations.end_block_index",
			Action: getAPIFlagAction[int64]("body", "requests.#.params.messages.#.content.#.source.content.#.citations.#.end_block_index"),
		},
		&cli.Int64Flag{
			Name:   "requests.params.messages.content.source.content.citations.start_block_index",
			Action: getAPIFlagAction[int64]("body", "requests.#.params.messages.#.content.#.source.content.#.citations.#.start_block_index"),
		},
		&cli.StringFlag{
			Name:   "requests.params.messages.content.source.content.citations.encrypted_index",
			Action: getAPIFlagAction[string]("body", "requests.#.params.messages.#.content.#.source.content.#.citations.#.encrypted_index"),
		},
		&cli.StringFlag{
			Name:   "requests.params.messages.content.source.content.citations.title",
			Action: getAPIFlagAction[string]("body", "requests.#.params.messages.#.content.#.source.content.#.citations.#.title"),
		},
		&cli.StringFlag{
			Name:   "requests.params.messages.content.source.content.citations.url",
			Action: getAPIFlagAction[string]("body", "requests.#.params.messages.#.content.#.source.content.#.citations.#.url"),
		},
		&cli.BoolFlag{
			Name:   "requests.params.messages.content.source.content.+citation",
			Action: getAPIFlagActionWithValue[bool]("body", "requests.#.params.messages.#.content.#.source.content.#.citations.-1", map[string]interface{}{}),
		},
		&cli.StringFlag{
			Name:   "requests.params.messages.content.source.content.source.data",
			Action: getAPIFlagAction[string]("body", "requests.#.params.messages.#.content.#.source.content.#.source.data"),
		},
		&cli.StringFlag{
			Name:   "requests.params.messages.content.source.content.source.media_type",
			Action: getAPIFlagAction[string]("body", "requests.#.params.messages.#.content.#.source.content.#.source.media_type"),
		},
		&cli.StringFlag{
			Name:   "requests.params.messages.content.source.content.source.type",
			Action: getAPIFlagAction[string]("body", "requests.#.params.messages.#.content.#.source.content.#.source.type"),
		},
		&cli.StringFlag{
			Name:   "requests.params.messages.content.source.content.source.url",
			Action: getAPIFlagAction[string]("body", "requests.#.params.messages.#.content.#.source.content.#.source.url"),
		},
		&cli.BoolFlag{
			Name:   "requests.params.messages.content.source.+content",
			Action: getAPIFlagActionWithValue[bool]("body", "requests.#.params.messages.#.content.#.source.content.-1", map[string]interface{}{}),
		},
		&cli.BoolFlag{
			Name:   "requests.params.messages.content.citations.enabled",
			Action: getAPIFlagAction[bool]("body", "requests.#.params.messages.#.content.#.citations.enabled"),
		},
		&cli.StringFlag{
			Name:   "requests.params.messages.content.context",
			Action: getAPIFlagAction[string]("body", "requests.#.params.messages.#.content.#.context"),
		},
		&cli.StringFlag{
			Name:   "requests.params.messages.content.title",
			Action: getAPIFlagAction[string]("body", "requests.#.params.messages.#.content.#.title"),
		},
		&cli.StringFlag{
			Name:   "requests.params.messages.content.signature",
			Action: getAPIFlagAction[string]("body", "requests.#.params.messages.#.content.#.signature"),
		},
		&cli.StringFlag{
			Name:   "requests.params.messages.content.thinking",
			Action: getAPIFlagAction[string]("body", "requests.#.params.messages.#.content.#.thinking"),
		},
		&cli.StringFlag{
			Name:   "requests.params.messages.content.data",
			Action: getAPIFlagAction[string]("body", "requests.#.params.messages.#.content.#.data"),
		},
		&cli.BoolFlag{
			Name:   "requests.params.messages.+content",
			Action: getAPIFlagActionWithValue[bool]("body", "requests.#.params.messages.#.content.-1", map[string]interface{}{}),
		},
		&cli.StringFlag{
			Name:   "requests.params.messages.role",
			Action: getAPIFlagAction[string]("body", "requests.#.params.messages.#.role"),
		},
		&cli.BoolFlag{
			Name:   "requests.params.+message",
			Action: getAPIFlagActionWithValue[bool]("body", "requests.#.params.messages.-1", map[string]interface{}{}),
		},
		&cli.StringFlag{
			Name:   "requests.params.model",
			Action: getAPIFlagAction[string]("body", "requests.#.params.model"),
		},
		&cli.StringFlag{
			Name:   "requests.params.metadata.user_id",
			Action: getAPIFlagAction[string]("body", "requests.#.params.metadata.user_id"),
		},
		&cli.StringFlag{
			Name:   "requests.params.service_tier",
			Action: getAPIFlagAction[string]("body", "requests.#.params.service_tier"),
		},
		&cli.StringFlag{
			Name:   "requests.params.stop_sequences",
			Action: getAPIFlagAction[string]("body", "requests.#.params.stop_sequences.#"),
		},
		&cli.StringFlag{
			Name:   "requests.params.+stop_sequence",
			Action: getAPIFlagAction[string]("body", "requests.#.params.stop_sequences.-1"),
		},
		&cli.BoolFlag{
			Name:   "requests.params.stream",
			Action: getAPIFlagAction[bool]("body", "requests.#.params.stream"),
		},
		&cli.StringFlag{
			Name:   "requests.params.system.text",
			Action: getAPIFlagAction[string]("body", "requests.#.params.system.#.text"),
		},
		&cli.StringFlag{
			Name:   "requests.params.system.type",
			Action: getAPIFlagAction[string]("body", "requests.#.params.system.#.type"),
		},
		&cli.StringFlag{
			Name:   "requests.params.system.cache_control.type",
			Action: getAPIFlagAction[string]("body", "requests.#.params.system.#.cache_control.type"),
		},
		&cli.StringFlag{
			Name:   "requests.params.system.citations.cited_text",
			Action: getAPIFlagAction[string]("body", "requests.#.params.system.#.citations.#.cited_text"),
		},
		&cli.Int64Flag{
			Name:   "requests.params.system.citations.document_index",
			Action: getAPIFlagAction[int64]("body", "requests.#.params.system.#.citations.#.document_index"),
		},
		&cli.StringFlag{
			Name:   "requests.params.system.citations.document_title",
			Action: getAPIFlagAction[string]("body", "requests.#.params.system.#.citations.#.document_title"),
		},
		&cli.Int64Flag{
			Name:   "requests.params.system.citations.end_char_index",
			Action: getAPIFlagAction[int64]("body", "requests.#.params.system.#.citations.#.end_char_index"),
		},
		&cli.Int64Flag{
			Name:   "requests.params.system.citations.start_char_index",
			Action: getAPIFlagAction[int64]("body", "requests.#.params.system.#.citations.#.start_char_index"),
		},
		&cli.StringFlag{
			Name:   "requests.params.system.citations.type",
			Action: getAPIFlagAction[string]("body", "requests.#.params.system.#.citations.#.type"),
		},
		&cli.Int64Flag{
			Name:   "requests.params.system.citations.end_page_number",
			Action: getAPIFlagAction[int64]("body", "requests.#.params.system.#.citations.#.end_page_number"),
		},
		&cli.Int64Flag{
			Name:   "requests.params.system.citations.start_page_number",
			Action: getAPIFlagAction[int64]("body", "requests.#.params.system.#.citations.#.start_page_number"),
		},
		&cli.Int64Flag{
			Name:   "requests.params.system.citations.end_block_index",
			Action: getAPIFlagAction[int64]("body", "requests.#.params.system.#.citations.#.end_block_index"),
		},
		&cli.Int64Flag{
			Name:   "requests.params.system.citations.start_block_index",
			Action: getAPIFlagAction[int64]("body", "requests.#.params.system.#.citations.#.start_block_index"),
		},
		&cli.StringFlag{
			Name:   "requests.params.system.citations.encrypted_index",
			Action: getAPIFlagAction[string]("body", "requests.#.params.system.#.citations.#.encrypted_index"),
		},
		&cli.StringFlag{
			Name:   "requests.params.system.citations.title",
			Action: getAPIFlagAction[string]("body", "requests.#.params.system.#.citations.#.title"),
		},
		&cli.StringFlag{
			Name:   "requests.params.system.citations.url",
			Action: getAPIFlagAction[string]("body", "requests.#.params.system.#.citations.#.url"),
		},
		&cli.BoolFlag{
			Name:   "requests.params.system.+citation",
			Action: getAPIFlagActionWithValue[bool]("body", "requests.#.params.system.#.citations.-1", map[string]interface{}{}),
		},
		&cli.BoolFlag{
			Name:   "requests.params.+system",
			Action: getAPIFlagActionWithValue[bool]("body", "requests.#.params.system.-1", map[string]interface{}{}),
		},
		&cli.FloatFlag{
			Name:   "requests.params.temperature",
			Action: getAPIFlagAction[float64]("body", "requests.#.params.temperature"),
		},
		&cli.Int64Flag{
			Name:   "requests.params.thinking.budget_tokens",
			Action: getAPIFlagAction[int64]("body", "requests.#.params.thinking.budget_tokens"),
		},
		&cli.StringFlag{
			Name:   "requests.params.thinking.type",
			Action: getAPIFlagAction[string]("body", "requests.#.params.thinking.type"),
		},
		&cli.StringFlag{
			Name:   "requests.params.tool_choice.type",
			Action: getAPIFlagAction[string]("body", "requests.#.params.tool_choice.type"),
		},
		&cli.BoolFlag{
			Name:   "requests.params.tool_choice.disable_parallel_tool_use",
			Action: getAPIFlagAction[bool]("body", "requests.#.params.tool_choice.disable_parallel_tool_use"),
		},
		&cli.StringFlag{
			Name:   "requests.params.tool_choice.name",
			Action: getAPIFlagAction[string]("body", "requests.#.params.tool_choice.name"),
		},
		&cli.StringFlag{
			Name:   "requests.params.tools.name",
			Action: getAPIFlagAction[string]("body", "requests.#.params.tools.#.name"),
		},
		&cli.StringFlag{
			Name:   "requests.params.tools.cache_control.type",
			Action: getAPIFlagAction[string]("body", "requests.#.params.tools.#.cache_control.type"),
		},
		&cli.StringFlag{
			Name:   "requests.params.tools.description",
			Action: getAPIFlagAction[string]("body", "requests.#.params.tools.#.description"),
		},
		&cli.StringFlag{
			Name:   "requests.params.tools.type",
			Action: getAPIFlagAction[string]("body", "requests.#.params.tools.#.type"),
		},
		&cli.StringFlag{
			Name:   "requests.params.tools.allowed_domains",
			Action: getAPIFlagAction[string]("body", "requests.#.params.tools.#.allowed_domains.#"),
		},
		&cli.StringFlag{
			Name:   "requests.params.tools.+allowed_domain",
			Action: getAPIFlagAction[string]("body", "requests.#.params.tools.#.allowed_domains.-1"),
		},
		&cli.StringFlag{
			Name:   "requests.params.tools.blocked_domains",
			Action: getAPIFlagAction[string]("body", "requests.#.params.tools.#.blocked_domains.#"),
		},
		&cli.StringFlag{
			Name:   "requests.params.tools.+blocked_domain",
			Action: getAPIFlagAction[string]("body", "requests.#.params.tools.#.blocked_domains.-1"),
		},
		&cli.Int64Flag{
			Name:   "requests.params.tools.max_uses",
			Action: getAPIFlagAction[int64]("body", "requests.#.params.tools.#.max_uses"),
		},
		&cli.StringFlag{
			Name:   "requests.params.tools.user_location.type",
			Action: getAPIFlagAction[string]("body", "requests.#.params.tools.#.user_location.type"),
		},
		&cli.StringFlag{
			Name:   "requests.params.tools.user_location.city",
			Action: getAPIFlagAction[string]("body", "requests.#.params.tools.#.user_location.city"),
		},
		&cli.StringFlag{
			Name:   "requests.params.tools.user_location.country",
			Action: getAPIFlagAction[string]("body", "requests.#.params.tools.#.user_location.country"),
		},
		&cli.StringFlag{
			Name:   "requests.params.tools.user_location.region",
			Action: getAPIFlagAction[string]("body", "requests.#.params.tools.#.user_location.region"),
		},
		&cli.StringFlag{
			Name:   "requests.params.tools.user_location.timezone",
			Action: getAPIFlagAction[string]("body", "requests.#.params.tools.#.user_location.timezone"),
		},
		&cli.BoolFlag{
			Name:   "requests.params.+tool",
			Action: getAPIFlagActionWithValue[bool]("body", "requests.#.params.tools.-1", map[string]interface{}{}),
		},
		&cli.Int64Flag{
			Name:   "requests.params.top_k",
			Action: getAPIFlagAction[int64]("body", "requests.#.params.top_k"),
		},
		&cli.FloatFlag{
			Name:   "requests.params.top_p",
			Action: getAPIFlagAction[float64]("body", "requests.#.params.top_p"),
		},
		&cli.BoolFlag{
			Name:   "+request",
			Action: getAPIFlagActionWithValue[bool]("body", "requests.-1", map[string]interface{}{}),
		},
	},
	Before:          initAPICommand,
	Action:          handleMessagesBatchesCreate,
	HideHelpCommand: true,
}

var messagesBatchesRetrieve = cli.Command{
	Name:  "retrieve",
	Usage: "This endpoint is idempotent and can be used to poll for Message Batch\ncompletion. To access the results of a Message Batch, make a request to the\n`results_url` field in the response.",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name: "message-batch-id",
		},
	},
	Before:          initAPICommand,
	Action:          handleMessagesBatchesRetrieve,
	HideHelpCommand: true,
}

var messagesBatchesList = cli.Command{
	Name:  "list",
	Usage: "List all Message Batches within a Workspace. Most recently created batches are\nreturned first.",
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
	},
	Before:          initAPICommand,
	Action:          handleMessagesBatchesList,
	HideHelpCommand: true,
}

var messagesBatchesDelete = cli.Command{
	Name:  "delete",
	Usage: "Delete a Message Batch.",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name: "message-batch-id",
		},
	},
	Before:          initAPICommand,
	Action:          handleMessagesBatchesDelete,
	HideHelpCommand: true,
}

var messagesBatchesCancel = cli.Command{
	Name:  "cancel",
	Usage: "Batches may be canceled any time before processing ends. Once cancellation is\ninitiated, the batch enters a `canceling` state, at which time the system may\ncomplete any in-progress, non-interruptible requests before finalizing\ncancellation.",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name: "message-batch-id",
		},
	},
	Before:          initAPICommand,
	Action:          handleMessagesBatchesCancel,
	HideHelpCommand: true,
}

var messagesBatchesResults = cli.Command{
	Name:  "results",
	Usage: "Streams the results of a Message Batch as a `.jsonl` file.",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name: "message-batch-id",
		},
	},
	Before:          initAPICommand,
	Action:          handleMessagesBatchesResults,
	HideHelpCommand: true,
}

func handleMessagesBatchesCreate(ctx context.Context, cmd *cli.Command) error {
	cc := getAPICommandContext(ctx, cmd)
	params := anthropic.MessageBatchNewParams{}
	res, err := cc.client.Messages.Batches.New(
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

func handleMessagesBatchesRetrieve(ctx context.Context, cmd *cli.Command) error {
	cc := getAPICommandContext(ctx, cmd)
	res, err := cc.client.Messages.Batches.Get(
		context.TODO(),
		cmd.Value("message-batch-id").(string),
		option.WithMiddleware(cc.AsMiddleware()),
	)
	if err != nil {
		return err
	}

	fmt.Printf("%s\n", colorizeJSON(res.RawJSON(), os.Stdout))
	return nil
}

func handleMessagesBatchesList(ctx context.Context, cmd *cli.Command) error {
	cc := getAPICommandContext(ctx, cmd)
	params := anthropic.MessageBatchListParams{}
	res, err := cc.client.Messages.Batches.List(
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

func handleMessagesBatchesDelete(ctx context.Context, cmd *cli.Command) error {
	cc := getAPICommandContext(ctx, cmd)
	res, err := cc.client.Messages.Batches.Delete(
		context.TODO(),
		cmd.Value("message-batch-id").(string),
		option.WithMiddleware(cc.AsMiddleware()),
		option.WithRequestBody("application/json", cc.body),
	)
	if err != nil {
		return err
	}

	fmt.Printf("%s\n", colorizeJSON(res.RawJSON(), os.Stdout))
	return nil
}

func handleMessagesBatchesCancel(ctx context.Context, cmd *cli.Command) error {
	cc := getAPICommandContext(ctx, cmd)
	res, err := cc.client.Messages.Batches.Cancel(
		context.TODO(),
		cmd.Value("message-batch-id").(string),
		option.WithMiddleware(cc.AsMiddleware()),
		option.WithRequestBody("application/json", cc.body),
	)
	if err != nil {
		return err
	}

	fmt.Printf("%s\n", colorizeJSON(res.RawJSON(), os.Stdout))
	return nil
}

func handleMessagesBatchesResults(ctx context.Context, cmd *cli.Command) error {
	cc := getAPICommandContext(ctx, cmd)
	stream := cc.client.Messages.Batches.ResultsStreaming(
		context.TODO(),
		cmd.Value("message-batch-id").(string),
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
