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

var messagesCreate = cli.Command{
	Name:  "create",
	Usage: "Send a structured list of input messages with text and/or image content, and the\nmodel will generate the next message in the conversation.",
	Flags: []cli.Flag{
		&cli.Int64Flag{
			Name:   "max-tokens",
			Action: getAPIFlagAction[int64]("body", "max_tokens"),
		},
		&cli.StringFlag{
			Name:   "messages.content.id",
			Action: getAPIFlagAction[string]("body", "messages.#.content.#.id"),
		},
		&cli.StringFlag{
			Name:   "messages.content.name",
			Action: getAPIFlagAction[string]("body", "messages.#.content.#.name"),
		},
		&cli.StringFlag{
			Name:   "messages.content.type",
			Action: getAPIFlagAction[string]("body", "messages.#.content.#.type"),
		},
		&cli.StringFlag{
			Name:   "messages.content.cache_control.type",
			Action: getAPIFlagAction[string]("body", "messages.#.content.#.cache_control.type"),
		},
		&cli.StringFlag{
			Name:   "messages.content.content.encrypted_content",
			Action: getAPIFlagAction[string]("body", "messages.#.content.#.content.#.encrypted_content"),
		},
		&cli.StringFlag{
			Name:   "messages.content.content.title",
			Action: getAPIFlagAction[string]("body", "messages.#.content.#.content.#.title"),
		},
		&cli.StringFlag{
			Name:   "messages.content.content.type",
			Action: getAPIFlagAction[string]("body", "messages.#.content.#.content.#.type"),
		},
		&cli.StringFlag{
			Name:   "messages.content.content.url",
			Action: getAPIFlagAction[string]("body", "messages.#.content.#.content.#.url"),
		},
		&cli.StringFlag{
			Name:   "messages.content.content.page_age",
			Action: getAPIFlagAction[string]("body", "messages.#.content.#.content.#.page_age"),
		},
		&cli.BoolFlag{
			Name:   "messages.content.+content",
			Action: getAPIFlagActionWithValue[bool]("body", "messages.#.content.#.content.-1", map[string]interface{}{}),
		},
		&cli.StringFlag{
			Name:   "messages.content.content.error_code",
			Action: getAPIFlagAction[string]("body", "messages.#.content.#.content.error_code"),
		},
		&cli.StringFlag{
			Name:   "messages.content.tool_use_id",
			Action: getAPIFlagAction[string]("body", "messages.#.content.#.tool_use_id"),
		},
		&cli.StringFlag{
			Name:   "messages.content.text",
			Action: getAPIFlagAction[string]("body", "messages.#.content.#.text"),
		},
		&cli.StringFlag{
			Name:   "messages.content.citations.cited_text",
			Action: getAPIFlagAction[string]("body", "messages.#.content.#.citations.#.cited_text"),
		},
		&cli.Int64Flag{
			Name:   "messages.content.citations.document_index",
			Action: getAPIFlagAction[int64]("body", "messages.#.content.#.citations.#.document_index"),
		},
		&cli.StringFlag{
			Name:   "messages.content.citations.document_title",
			Action: getAPIFlagAction[string]("body", "messages.#.content.#.citations.#.document_title"),
		},
		&cli.Int64Flag{
			Name:   "messages.content.citations.end_char_index",
			Action: getAPIFlagAction[int64]("body", "messages.#.content.#.citations.#.end_char_index"),
		},
		&cli.Int64Flag{
			Name:   "messages.content.citations.start_char_index",
			Action: getAPIFlagAction[int64]("body", "messages.#.content.#.citations.#.start_char_index"),
		},
		&cli.StringFlag{
			Name:   "messages.content.citations.type",
			Action: getAPIFlagAction[string]("body", "messages.#.content.#.citations.#.type"),
		},
		&cli.Int64Flag{
			Name:   "messages.content.citations.end_page_number",
			Action: getAPIFlagAction[int64]("body", "messages.#.content.#.citations.#.end_page_number"),
		},
		&cli.Int64Flag{
			Name:   "messages.content.citations.start_page_number",
			Action: getAPIFlagAction[int64]("body", "messages.#.content.#.citations.#.start_page_number"),
		},
		&cli.Int64Flag{
			Name:   "messages.content.citations.end_block_index",
			Action: getAPIFlagAction[int64]("body", "messages.#.content.#.citations.#.end_block_index"),
		},
		&cli.Int64Flag{
			Name:   "messages.content.citations.start_block_index",
			Action: getAPIFlagAction[int64]("body", "messages.#.content.#.citations.#.start_block_index"),
		},
		&cli.StringFlag{
			Name:   "messages.content.citations.encrypted_index",
			Action: getAPIFlagAction[string]("body", "messages.#.content.#.citations.#.encrypted_index"),
		},
		&cli.StringFlag{
			Name:   "messages.content.citations.title",
			Action: getAPIFlagAction[string]("body", "messages.#.content.#.citations.#.title"),
		},
		&cli.StringFlag{
			Name:   "messages.content.citations.url",
			Action: getAPIFlagAction[string]("body", "messages.#.content.#.citations.#.url"),
		},
		&cli.BoolFlag{
			Name:   "messages.content.+citation",
			Action: getAPIFlagActionWithValue[bool]("body", "messages.#.content.#.citations.-1", map[string]interface{}{}),
		},
		&cli.StringFlag{
			Name:   "messages.content.source.data",
			Action: getAPIFlagAction[string]("body", "messages.#.content.#.source.data"),
		},
		&cli.StringFlag{
			Name:   "messages.content.source.media_type",
			Action: getAPIFlagAction[string]("body", "messages.#.content.#.source.media_type"),
		},
		&cli.StringFlag{
			Name:   "messages.content.source.type",
			Action: getAPIFlagAction[string]("body", "messages.#.content.#.source.type"),
		},
		&cli.StringFlag{
			Name:   "messages.content.source.url",
			Action: getAPIFlagAction[string]("body", "messages.#.content.#.source.url"),
		},
		&cli.StringFlag{
			Name:   "messages.content.content.text",
			Action: getAPIFlagAction[string]("body", "messages.#.content.#.content.#.text"),
		},
		&cli.StringFlag{
			Name:   "messages.content.content.cache_control.type",
			Action: getAPIFlagAction[string]("body", "messages.#.content.#.content.#.cache_control.type"),
		},
		&cli.StringFlag{
			Name:   "messages.content.content.citations.cited_text",
			Action: getAPIFlagAction[string]("body", "messages.#.content.#.content.#.citations.#.cited_text"),
		},
		&cli.Int64Flag{
			Name:   "messages.content.content.citations.document_index",
			Action: getAPIFlagAction[int64]("body", "messages.#.content.#.content.#.citations.#.document_index"),
		},
		&cli.StringFlag{
			Name:   "messages.content.content.citations.document_title",
			Action: getAPIFlagAction[string]("body", "messages.#.content.#.content.#.citations.#.document_title"),
		},
		&cli.Int64Flag{
			Name:   "messages.content.content.citations.end_char_index",
			Action: getAPIFlagAction[int64]("body", "messages.#.content.#.content.#.citations.#.end_char_index"),
		},
		&cli.Int64Flag{
			Name:   "messages.content.content.citations.start_char_index",
			Action: getAPIFlagAction[int64]("body", "messages.#.content.#.content.#.citations.#.start_char_index"),
		},
		&cli.StringFlag{
			Name:   "messages.content.content.citations.type",
			Action: getAPIFlagAction[string]("body", "messages.#.content.#.content.#.citations.#.type"),
		},
		&cli.Int64Flag{
			Name:   "messages.content.content.citations.end_page_number",
			Action: getAPIFlagAction[int64]("body", "messages.#.content.#.content.#.citations.#.end_page_number"),
		},
		&cli.Int64Flag{
			Name:   "messages.content.content.citations.start_page_number",
			Action: getAPIFlagAction[int64]("body", "messages.#.content.#.content.#.citations.#.start_page_number"),
		},
		&cli.Int64Flag{
			Name:   "messages.content.content.citations.end_block_index",
			Action: getAPIFlagAction[int64]("body", "messages.#.content.#.content.#.citations.#.end_block_index"),
		},
		&cli.Int64Flag{
			Name:   "messages.content.content.citations.start_block_index",
			Action: getAPIFlagAction[int64]("body", "messages.#.content.#.content.#.citations.#.start_block_index"),
		},
		&cli.StringFlag{
			Name:   "messages.content.content.citations.encrypted_index",
			Action: getAPIFlagAction[string]("body", "messages.#.content.#.content.#.citations.#.encrypted_index"),
		},
		&cli.StringFlag{
			Name:   "messages.content.content.citations.title",
			Action: getAPIFlagAction[string]("body", "messages.#.content.#.content.#.citations.#.title"),
		},
		&cli.StringFlag{
			Name:   "messages.content.content.citations.url",
			Action: getAPIFlagAction[string]("body", "messages.#.content.#.content.#.citations.#.url"),
		},
		&cli.BoolFlag{
			Name:   "messages.content.content.+citation",
			Action: getAPIFlagActionWithValue[bool]("body", "messages.#.content.#.content.#.citations.-1", map[string]interface{}{}),
		},
		&cli.StringFlag{
			Name:   "messages.content.content.source.data",
			Action: getAPIFlagAction[string]("body", "messages.#.content.#.content.#.source.data"),
		},
		&cli.StringFlag{
			Name:   "messages.content.content.source.media_type",
			Action: getAPIFlagAction[string]("body", "messages.#.content.#.content.#.source.media_type"),
		},
		&cli.StringFlag{
			Name:   "messages.content.content.source.type",
			Action: getAPIFlagAction[string]("body", "messages.#.content.#.content.#.source.type"),
		},
		&cli.StringFlag{
			Name:   "messages.content.content.source.url",
			Action: getAPIFlagAction[string]("body", "messages.#.content.#.content.#.source.url"),
		},
		&cli.BoolFlag{
			Name:   "messages.content.is_error",
			Action: getAPIFlagAction[bool]("body", "messages.#.content.#.is_error"),
		},
		&cli.StringFlag{
			Name:   "messages.content.source.content",
			Action: getAPIFlagAction[string]("body", "messages.#.content.#.source.content"),
		},
		&cli.StringFlag{
			Name:   "messages.content.source.content.text",
			Action: getAPIFlagAction[string]("body", "messages.#.content.#.source.content.#.text"),
		},
		&cli.StringFlag{
			Name:   "messages.content.source.content.type",
			Action: getAPIFlagAction[string]("body", "messages.#.content.#.source.content.#.type"),
		},
		&cli.StringFlag{
			Name:   "messages.content.source.content.cache_control.type",
			Action: getAPIFlagAction[string]("body", "messages.#.content.#.source.content.#.cache_control.type"),
		},
		&cli.StringFlag{
			Name:   "messages.content.source.content.citations.cited_text",
			Action: getAPIFlagAction[string]("body", "messages.#.content.#.source.content.#.citations.#.cited_text"),
		},
		&cli.Int64Flag{
			Name:   "messages.content.source.content.citations.document_index",
			Action: getAPIFlagAction[int64]("body", "messages.#.content.#.source.content.#.citations.#.document_index"),
		},
		&cli.StringFlag{
			Name:   "messages.content.source.content.citations.document_title",
			Action: getAPIFlagAction[string]("body", "messages.#.content.#.source.content.#.citations.#.document_title"),
		},
		&cli.Int64Flag{
			Name:   "messages.content.source.content.citations.end_char_index",
			Action: getAPIFlagAction[int64]("body", "messages.#.content.#.source.content.#.citations.#.end_char_index"),
		},
		&cli.Int64Flag{
			Name:   "messages.content.source.content.citations.start_char_index",
			Action: getAPIFlagAction[int64]("body", "messages.#.content.#.source.content.#.citations.#.start_char_index"),
		},
		&cli.StringFlag{
			Name:   "messages.content.source.content.citations.type",
			Action: getAPIFlagAction[string]("body", "messages.#.content.#.source.content.#.citations.#.type"),
		},
		&cli.Int64Flag{
			Name:   "messages.content.source.content.citations.end_page_number",
			Action: getAPIFlagAction[int64]("body", "messages.#.content.#.source.content.#.citations.#.end_page_number"),
		},
		&cli.Int64Flag{
			Name:   "messages.content.source.content.citations.start_page_number",
			Action: getAPIFlagAction[int64]("body", "messages.#.content.#.source.content.#.citations.#.start_page_number"),
		},
		&cli.Int64Flag{
			Name:   "messages.content.source.content.citations.end_block_index",
			Action: getAPIFlagAction[int64]("body", "messages.#.content.#.source.content.#.citations.#.end_block_index"),
		},
		&cli.Int64Flag{
			Name:   "messages.content.source.content.citations.start_block_index",
			Action: getAPIFlagAction[int64]("body", "messages.#.content.#.source.content.#.citations.#.start_block_index"),
		},
		&cli.StringFlag{
			Name:   "messages.content.source.content.citations.encrypted_index",
			Action: getAPIFlagAction[string]("body", "messages.#.content.#.source.content.#.citations.#.encrypted_index"),
		},
		&cli.StringFlag{
			Name:   "messages.content.source.content.citations.title",
			Action: getAPIFlagAction[string]("body", "messages.#.content.#.source.content.#.citations.#.title"),
		},
		&cli.StringFlag{
			Name:   "messages.content.source.content.citations.url",
			Action: getAPIFlagAction[string]("body", "messages.#.content.#.source.content.#.citations.#.url"),
		},
		&cli.BoolFlag{
			Name:   "messages.content.source.content.+citation",
			Action: getAPIFlagActionWithValue[bool]("body", "messages.#.content.#.source.content.#.citations.-1", map[string]interface{}{}),
		},
		&cli.StringFlag{
			Name:   "messages.content.source.content.source.data",
			Action: getAPIFlagAction[string]("body", "messages.#.content.#.source.content.#.source.data"),
		},
		&cli.StringFlag{
			Name:   "messages.content.source.content.source.media_type",
			Action: getAPIFlagAction[string]("body", "messages.#.content.#.source.content.#.source.media_type"),
		},
		&cli.StringFlag{
			Name:   "messages.content.source.content.source.type",
			Action: getAPIFlagAction[string]("body", "messages.#.content.#.source.content.#.source.type"),
		},
		&cli.StringFlag{
			Name:   "messages.content.source.content.source.url",
			Action: getAPIFlagAction[string]("body", "messages.#.content.#.source.content.#.source.url"),
		},
		&cli.BoolFlag{
			Name:   "messages.content.source.+content",
			Action: getAPIFlagActionWithValue[bool]("body", "messages.#.content.#.source.content.-1", map[string]interface{}{}),
		},
		&cli.BoolFlag{
			Name:   "messages.content.citations.enabled",
			Action: getAPIFlagAction[bool]("body", "messages.#.content.#.citations.enabled"),
		},
		&cli.StringFlag{
			Name:   "messages.content.context",
			Action: getAPIFlagAction[string]("body", "messages.#.content.#.context"),
		},
		&cli.StringFlag{
			Name:   "messages.content.title",
			Action: getAPIFlagAction[string]("body", "messages.#.content.#.title"),
		},
		&cli.StringFlag{
			Name:   "messages.content.signature",
			Action: getAPIFlagAction[string]("body", "messages.#.content.#.signature"),
		},
		&cli.StringFlag{
			Name:   "messages.content.thinking",
			Action: getAPIFlagAction[string]("body", "messages.#.content.#.thinking"),
		},
		&cli.StringFlag{
			Name:   "messages.content.data",
			Action: getAPIFlagAction[string]("body", "messages.#.content.#.data"),
		},
		&cli.BoolFlag{
			Name:   "messages.+content",
			Action: getAPIFlagActionWithValue[bool]("body", "messages.#.content.-1", map[string]interface{}{}),
		},
		&cli.StringFlag{
			Name:   "messages.role",
			Action: getAPIFlagAction[string]("body", "messages.#.role"),
		},
		&cli.BoolFlag{
			Name:   "+message",
			Action: getAPIFlagActionWithValue[bool]("body", "messages.-1", map[string]interface{}{}),
		},
		&cli.StringFlag{
			Name:   "model",
			Action: getAPIFlagAction[string]("body", "model"),
		},
		&cli.StringFlag{
			Name:   "metadata.user_id",
			Action: getAPIFlagAction[string]("body", "metadata.user_id"),
		},
		&cli.StringFlag{
			Name:   "service-tier",
			Action: getAPIFlagAction[string]("body", "service_tier"),
		},
		&cli.StringFlag{
			Name:   "stop-sequences",
			Action: getAPIFlagAction[string]("body", "stop_sequences.#"),
		},
		&cli.StringFlag{
			Name:   "+stop_sequence",
			Action: getAPIFlagAction[string]("body", "stop_sequences.-1"),
		},
		&cli.StringFlag{
			Name:   "system.text",
			Action: getAPIFlagAction[string]("body", "system.#.text"),
		},
		&cli.StringFlag{
			Name:   "system.type",
			Action: getAPIFlagAction[string]("body", "system.#.type"),
		},
		&cli.StringFlag{
			Name:   "system.cache_control.type",
			Action: getAPIFlagAction[string]("body", "system.#.cache_control.type"),
		},
		&cli.StringFlag{
			Name:   "system.citations.cited_text",
			Action: getAPIFlagAction[string]("body", "system.#.citations.#.cited_text"),
		},
		&cli.Int64Flag{
			Name:   "system.citations.document_index",
			Action: getAPIFlagAction[int64]("body", "system.#.citations.#.document_index"),
		},
		&cli.StringFlag{
			Name:   "system.citations.document_title",
			Action: getAPIFlagAction[string]("body", "system.#.citations.#.document_title"),
		},
		&cli.Int64Flag{
			Name:   "system.citations.end_char_index",
			Action: getAPIFlagAction[int64]("body", "system.#.citations.#.end_char_index"),
		},
		&cli.Int64Flag{
			Name:   "system.citations.start_char_index",
			Action: getAPIFlagAction[int64]("body", "system.#.citations.#.start_char_index"),
		},
		&cli.StringFlag{
			Name:   "system.citations.type",
			Action: getAPIFlagAction[string]("body", "system.#.citations.#.type"),
		},
		&cli.Int64Flag{
			Name:   "system.citations.end_page_number",
			Action: getAPIFlagAction[int64]("body", "system.#.citations.#.end_page_number"),
		},
		&cli.Int64Flag{
			Name:   "system.citations.start_page_number",
			Action: getAPIFlagAction[int64]("body", "system.#.citations.#.start_page_number"),
		},
		&cli.Int64Flag{
			Name:   "system.citations.end_block_index",
			Action: getAPIFlagAction[int64]("body", "system.#.citations.#.end_block_index"),
		},
		&cli.Int64Flag{
			Name:   "system.citations.start_block_index",
			Action: getAPIFlagAction[int64]("body", "system.#.citations.#.start_block_index"),
		},
		&cli.StringFlag{
			Name:   "system.citations.encrypted_index",
			Action: getAPIFlagAction[string]("body", "system.#.citations.#.encrypted_index"),
		},
		&cli.StringFlag{
			Name:   "system.citations.title",
			Action: getAPIFlagAction[string]("body", "system.#.citations.#.title"),
		},
		&cli.StringFlag{
			Name:   "system.citations.url",
			Action: getAPIFlagAction[string]("body", "system.#.citations.#.url"),
		},
		&cli.BoolFlag{
			Name:   "system.+citation",
			Action: getAPIFlagActionWithValue[bool]("body", "system.#.citations.-1", map[string]interface{}{}),
		},
		&cli.BoolFlag{
			Name:   "+system",
			Action: getAPIFlagActionWithValue[bool]("body", "system.-1", map[string]interface{}{}),
		},
		&cli.FloatFlag{
			Name:   "temperature",
			Action: getAPIFlagAction[float64]("body", "temperature"),
		},
		&cli.Int64Flag{
			Name:   "thinking.budget_tokens",
			Action: getAPIFlagAction[int64]("body", "thinking.budget_tokens"),
		},
		&cli.StringFlag{
			Name:   "thinking.type",
			Action: getAPIFlagAction[string]("body", "thinking.type"),
		},
		&cli.StringFlag{
			Name:   "tool-choice.type",
			Action: getAPIFlagAction[string]("body", "tool_choice.type"),
		},
		&cli.BoolFlag{
			Name:   "tool-choice.disable_parallel_tool_use",
			Action: getAPIFlagAction[bool]("body", "tool_choice.disable_parallel_tool_use"),
		},
		&cli.StringFlag{
			Name:   "tool-choice.name",
			Action: getAPIFlagAction[string]("body", "tool_choice.name"),
		},
		&cli.StringFlag{
			Name:   "tools.name",
			Action: getAPIFlagAction[string]("body", "tools.#.name"),
		},
		&cli.StringFlag{
			Name:   "tools.cache_control.type",
			Action: getAPIFlagAction[string]("body", "tools.#.cache_control.type"),
		},
		&cli.StringFlag{
			Name:   "tools.description",
			Action: getAPIFlagAction[string]("body", "tools.#.description"),
		},
		&cli.StringFlag{
			Name:   "tools.type",
			Action: getAPIFlagAction[string]("body", "tools.#.type"),
		},
		&cli.StringFlag{
			Name:   "tools.allowed_domains",
			Action: getAPIFlagAction[string]("body", "tools.#.allowed_domains.#"),
		},
		&cli.StringFlag{
			Name:   "tools.+allowed_domain",
			Action: getAPIFlagAction[string]("body", "tools.#.allowed_domains.-1"),
		},
		&cli.StringFlag{
			Name:   "tools.blocked_domains",
			Action: getAPIFlagAction[string]("body", "tools.#.blocked_domains.#"),
		},
		&cli.StringFlag{
			Name:   "tools.+blocked_domain",
			Action: getAPIFlagAction[string]("body", "tools.#.blocked_domains.-1"),
		},
		&cli.Int64Flag{
			Name:   "tools.max_uses",
			Action: getAPIFlagAction[int64]("body", "tools.#.max_uses"),
		},
		&cli.StringFlag{
			Name:   "tools.user_location.type",
			Action: getAPIFlagAction[string]("body", "tools.#.user_location.type"),
		},
		&cli.StringFlag{
			Name:   "tools.user_location.city",
			Action: getAPIFlagAction[string]("body", "tools.#.user_location.city"),
		},
		&cli.StringFlag{
			Name:   "tools.user_location.country",
			Action: getAPIFlagAction[string]("body", "tools.#.user_location.country"),
		},
		&cli.StringFlag{
			Name:   "tools.user_location.region",
			Action: getAPIFlagAction[string]("body", "tools.#.user_location.region"),
		},
		&cli.StringFlag{
			Name:   "tools.user_location.timezone",
			Action: getAPIFlagAction[string]("body", "tools.#.user_location.timezone"),
		},
		&cli.BoolFlag{
			Name:   "+tool",
			Action: getAPIFlagActionWithValue[bool]("body", "tools.-1", map[string]interface{}{}),
		},
		&cli.Int64Flag{
			Name:   "top-k",
			Action: getAPIFlagAction[int64]("body", "top_k"),
		},
		&cli.FloatFlag{
			Name:   "top-p",
			Action: getAPIFlagAction[float64]("body", "top_p"),
		},
	},
	Before:          initAPICommand,
	Action:          handleMessagesCreate,
	HideHelpCommand: true,
}

var messagesCountTokens = cli.Command{
	Name:  "count_tokens",
	Usage: "Count the number of tokens in a Message.",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:   "messages.content.id",
			Action: getAPIFlagAction[string]("body", "messages.#.content.#.id"),
		},
		&cli.StringFlag{
			Name:   "messages.content.name",
			Action: getAPIFlagAction[string]("body", "messages.#.content.#.name"),
		},
		&cli.StringFlag{
			Name:   "messages.content.type",
			Action: getAPIFlagAction[string]("body", "messages.#.content.#.type"),
		},
		&cli.StringFlag{
			Name:   "messages.content.cache_control.type",
			Action: getAPIFlagAction[string]("body", "messages.#.content.#.cache_control.type"),
		},
		&cli.StringFlag{
			Name:   "messages.content.content.encrypted_content",
			Action: getAPIFlagAction[string]("body", "messages.#.content.#.content.#.encrypted_content"),
		},
		&cli.StringFlag{
			Name:   "messages.content.content.title",
			Action: getAPIFlagAction[string]("body", "messages.#.content.#.content.#.title"),
		},
		&cli.StringFlag{
			Name:   "messages.content.content.type",
			Action: getAPIFlagAction[string]("body", "messages.#.content.#.content.#.type"),
		},
		&cli.StringFlag{
			Name:   "messages.content.content.url",
			Action: getAPIFlagAction[string]("body", "messages.#.content.#.content.#.url"),
		},
		&cli.StringFlag{
			Name:   "messages.content.content.page_age",
			Action: getAPIFlagAction[string]("body", "messages.#.content.#.content.#.page_age"),
		},
		&cli.BoolFlag{
			Name:   "messages.content.+content",
			Action: getAPIFlagActionWithValue[bool]("body", "messages.#.content.#.content.-1", map[string]interface{}{}),
		},
		&cli.StringFlag{
			Name:   "messages.content.content.error_code",
			Action: getAPIFlagAction[string]("body", "messages.#.content.#.content.error_code"),
		},
		&cli.StringFlag{
			Name:   "messages.content.tool_use_id",
			Action: getAPIFlagAction[string]("body", "messages.#.content.#.tool_use_id"),
		},
		&cli.StringFlag{
			Name:   "messages.content.text",
			Action: getAPIFlagAction[string]("body", "messages.#.content.#.text"),
		},
		&cli.StringFlag{
			Name:   "messages.content.citations.cited_text",
			Action: getAPIFlagAction[string]("body", "messages.#.content.#.citations.#.cited_text"),
		},
		&cli.Int64Flag{
			Name:   "messages.content.citations.document_index",
			Action: getAPIFlagAction[int64]("body", "messages.#.content.#.citations.#.document_index"),
		},
		&cli.StringFlag{
			Name:   "messages.content.citations.document_title",
			Action: getAPIFlagAction[string]("body", "messages.#.content.#.citations.#.document_title"),
		},
		&cli.Int64Flag{
			Name:   "messages.content.citations.end_char_index",
			Action: getAPIFlagAction[int64]("body", "messages.#.content.#.citations.#.end_char_index"),
		},
		&cli.Int64Flag{
			Name:   "messages.content.citations.start_char_index",
			Action: getAPIFlagAction[int64]("body", "messages.#.content.#.citations.#.start_char_index"),
		},
		&cli.StringFlag{
			Name:   "messages.content.citations.type",
			Action: getAPIFlagAction[string]("body", "messages.#.content.#.citations.#.type"),
		},
		&cli.Int64Flag{
			Name:   "messages.content.citations.end_page_number",
			Action: getAPIFlagAction[int64]("body", "messages.#.content.#.citations.#.end_page_number"),
		},
		&cli.Int64Flag{
			Name:   "messages.content.citations.start_page_number",
			Action: getAPIFlagAction[int64]("body", "messages.#.content.#.citations.#.start_page_number"),
		},
		&cli.Int64Flag{
			Name:   "messages.content.citations.end_block_index",
			Action: getAPIFlagAction[int64]("body", "messages.#.content.#.citations.#.end_block_index"),
		},
		&cli.Int64Flag{
			Name:   "messages.content.citations.start_block_index",
			Action: getAPIFlagAction[int64]("body", "messages.#.content.#.citations.#.start_block_index"),
		},
		&cli.StringFlag{
			Name:   "messages.content.citations.encrypted_index",
			Action: getAPIFlagAction[string]("body", "messages.#.content.#.citations.#.encrypted_index"),
		},
		&cli.StringFlag{
			Name:   "messages.content.citations.title",
			Action: getAPIFlagAction[string]("body", "messages.#.content.#.citations.#.title"),
		},
		&cli.StringFlag{
			Name:   "messages.content.citations.url",
			Action: getAPIFlagAction[string]("body", "messages.#.content.#.citations.#.url"),
		},
		&cli.BoolFlag{
			Name:   "messages.content.+citation",
			Action: getAPIFlagActionWithValue[bool]("body", "messages.#.content.#.citations.-1", map[string]interface{}{}),
		},
		&cli.StringFlag{
			Name:   "messages.content.source.data",
			Action: getAPIFlagAction[string]("body", "messages.#.content.#.source.data"),
		},
		&cli.StringFlag{
			Name:   "messages.content.source.media_type",
			Action: getAPIFlagAction[string]("body", "messages.#.content.#.source.media_type"),
		},
		&cli.StringFlag{
			Name:   "messages.content.source.type",
			Action: getAPIFlagAction[string]("body", "messages.#.content.#.source.type"),
		},
		&cli.StringFlag{
			Name:   "messages.content.source.url",
			Action: getAPIFlagAction[string]("body", "messages.#.content.#.source.url"),
		},
		&cli.StringFlag{
			Name:   "messages.content.content.text",
			Action: getAPIFlagAction[string]("body", "messages.#.content.#.content.#.text"),
		},
		&cli.StringFlag{
			Name:   "messages.content.content.cache_control.type",
			Action: getAPIFlagAction[string]("body", "messages.#.content.#.content.#.cache_control.type"),
		},
		&cli.StringFlag{
			Name:   "messages.content.content.citations.cited_text",
			Action: getAPIFlagAction[string]("body", "messages.#.content.#.content.#.citations.#.cited_text"),
		},
		&cli.Int64Flag{
			Name:   "messages.content.content.citations.document_index",
			Action: getAPIFlagAction[int64]("body", "messages.#.content.#.content.#.citations.#.document_index"),
		},
		&cli.StringFlag{
			Name:   "messages.content.content.citations.document_title",
			Action: getAPIFlagAction[string]("body", "messages.#.content.#.content.#.citations.#.document_title"),
		},
		&cli.Int64Flag{
			Name:   "messages.content.content.citations.end_char_index",
			Action: getAPIFlagAction[int64]("body", "messages.#.content.#.content.#.citations.#.end_char_index"),
		},
		&cli.Int64Flag{
			Name:   "messages.content.content.citations.start_char_index",
			Action: getAPIFlagAction[int64]("body", "messages.#.content.#.content.#.citations.#.start_char_index"),
		},
		&cli.StringFlag{
			Name:   "messages.content.content.citations.type",
			Action: getAPIFlagAction[string]("body", "messages.#.content.#.content.#.citations.#.type"),
		},
		&cli.Int64Flag{
			Name:   "messages.content.content.citations.end_page_number",
			Action: getAPIFlagAction[int64]("body", "messages.#.content.#.content.#.citations.#.end_page_number"),
		},
		&cli.Int64Flag{
			Name:   "messages.content.content.citations.start_page_number",
			Action: getAPIFlagAction[int64]("body", "messages.#.content.#.content.#.citations.#.start_page_number"),
		},
		&cli.Int64Flag{
			Name:   "messages.content.content.citations.end_block_index",
			Action: getAPIFlagAction[int64]("body", "messages.#.content.#.content.#.citations.#.end_block_index"),
		},
		&cli.Int64Flag{
			Name:   "messages.content.content.citations.start_block_index",
			Action: getAPIFlagAction[int64]("body", "messages.#.content.#.content.#.citations.#.start_block_index"),
		},
		&cli.StringFlag{
			Name:   "messages.content.content.citations.encrypted_index",
			Action: getAPIFlagAction[string]("body", "messages.#.content.#.content.#.citations.#.encrypted_index"),
		},
		&cli.StringFlag{
			Name:   "messages.content.content.citations.title",
			Action: getAPIFlagAction[string]("body", "messages.#.content.#.content.#.citations.#.title"),
		},
		&cli.StringFlag{
			Name:   "messages.content.content.citations.url",
			Action: getAPIFlagAction[string]("body", "messages.#.content.#.content.#.citations.#.url"),
		},
		&cli.BoolFlag{
			Name:   "messages.content.content.+citation",
			Action: getAPIFlagActionWithValue[bool]("body", "messages.#.content.#.content.#.citations.-1", map[string]interface{}{}),
		},
		&cli.StringFlag{
			Name:   "messages.content.content.source.data",
			Action: getAPIFlagAction[string]("body", "messages.#.content.#.content.#.source.data"),
		},
		&cli.StringFlag{
			Name:   "messages.content.content.source.media_type",
			Action: getAPIFlagAction[string]("body", "messages.#.content.#.content.#.source.media_type"),
		},
		&cli.StringFlag{
			Name:   "messages.content.content.source.type",
			Action: getAPIFlagAction[string]("body", "messages.#.content.#.content.#.source.type"),
		},
		&cli.StringFlag{
			Name:   "messages.content.content.source.url",
			Action: getAPIFlagAction[string]("body", "messages.#.content.#.content.#.source.url"),
		},
		&cli.BoolFlag{
			Name:   "messages.content.is_error",
			Action: getAPIFlagAction[bool]("body", "messages.#.content.#.is_error"),
		},
		&cli.StringFlag{
			Name:   "messages.content.source.content",
			Action: getAPIFlagAction[string]("body", "messages.#.content.#.source.content"),
		},
		&cli.StringFlag{
			Name:   "messages.content.source.content.text",
			Action: getAPIFlagAction[string]("body", "messages.#.content.#.source.content.#.text"),
		},
		&cli.StringFlag{
			Name:   "messages.content.source.content.type",
			Action: getAPIFlagAction[string]("body", "messages.#.content.#.source.content.#.type"),
		},
		&cli.StringFlag{
			Name:   "messages.content.source.content.cache_control.type",
			Action: getAPIFlagAction[string]("body", "messages.#.content.#.source.content.#.cache_control.type"),
		},
		&cli.StringFlag{
			Name:   "messages.content.source.content.citations.cited_text",
			Action: getAPIFlagAction[string]("body", "messages.#.content.#.source.content.#.citations.#.cited_text"),
		},
		&cli.Int64Flag{
			Name:   "messages.content.source.content.citations.document_index",
			Action: getAPIFlagAction[int64]("body", "messages.#.content.#.source.content.#.citations.#.document_index"),
		},
		&cli.StringFlag{
			Name:   "messages.content.source.content.citations.document_title",
			Action: getAPIFlagAction[string]("body", "messages.#.content.#.source.content.#.citations.#.document_title"),
		},
		&cli.Int64Flag{
			Name:   "messages.content.source.content.citations.end_char_index",
			Action: getAPIFlagAction[int64]("body", "messages.#.content.#.source.content.#.citations.#.end_char_index"),
		},
		&cli.Int64Flag{
			Name:   "messages.content.source.content.citations.start_char_index",
			Action: getAPIFlagAction[int64]("body", "messages.#.content.#.source.content.#.citations.#.start_char_index"),
		},
		&cli.StringFlag{
			Name:   "messages.content.source.content.citations.type",
			Action: getAPIFlagAction[string]("body", "messages.#.content.#.source.content.#.citations.#.type"),
		},
		&cli.Int64Flag{
			Name:   "messages.content.source.content.citations.end_page_number",
			Action: getAPIFlagAction[int64]("body", "messages.#.content.#.source.content.#.citations.#.end_page_number"),
		},
		&cli.Int64Flag{
			Name:   "messages.content.source.content.citations.start_page_number",
			Action: getAPIFlagAction[int64]("body", "messages.#.content.#.source.content.#.citations.#.start_page_number"),
		},
		&cli.Int64Flag{
			Name:   "messages.content.source.content.citations.end_block_index",
			Action: getAPIFlagAction[int64]("body", "messages.#.content.#.source.content.#.citations.#.end_block_index"),
		},
		&cli.Int64Flag{
			Name:   "messages.content.source.content.citations.start_block_index",
			Action: getAPIFlagAction[int64]("body", "messages.#.content.#.source.content.#.citations.#.start_block_index"),
		},
		&cli.StringFlag{
			Name:   "messages.content.source.content.citations.encrypted_index",
			Action: getAPIFlagAction[string]("body", "messages.#.content.#.source.content.#.citations.#.encrypted_index"),
		},
		&cli.StringFlag{
			Name:   "messages.content.source.content.citations.title",
			Action: getAPIFlagAction[string]("body", "messages.#.content.#.source.content.#.citations.#.title"),
		},
		&cli.StringFlag{
			Name:   "messages.content.source.content.citations.url",
			Action: getAPIFlagAction[string]("body", "messages.#.content.#.source.content.#.citations.#.url"),
		},
		&cli.BoolFlag{
			Name:   "messages.content.source.content.+citation",
			Action: getAPIFlagActionWithValue[bool]("body", "messages.#.content.#.source.content.#.citations.-1", map[string]interface{}{}),
		},
		&cli.StringFlag{
			Name:   "messages.content.source.content.source.data",
			Action: getAPIFlagAction[string]("body", "messages.#.content.#.source.content.#.source.data"),
		},
		&cli.StringFlag{
			Name:   "messages.content.source.content.source.media_type",
			Action: getAPIFlagAction[string]("body", "messages.#.content.#.source.content.#.source.media_type"),
		},
		&cli.StringFlag{
			Name:   "messages.content.source.content.source.type",
			Action: getAPIFlagAction[string]("body", "messages.#.content.#.source.content.#.source.type"),
		},
		&cli.StringFlag{
			Name:   "messages.content.source.content.source.url",
			Action: getAPIFlagAction[string]("body", "messages.#.content.#.source.content.#.source.url"),
		},
		&cli.BoolFlag{
			Name:   "messages.content.source.+content",
			Action: getAPIFlagActionWithValue[bool]("body", "messages.#.content.#.source.content.-1", map[string]interface{}{}),
		},
		&cli.BoolFlag{
			Name:   "messages.content.citations.enabled",
			Action: getAPIFlagAction[bool]("body", "messages.#.content.#.citations.enabled"),
		},
		&cli.StringFlag{
			Name:   "messages.content.context",
			Action: getAPIFlagAction[string]("body", "messages.#.content.#.context"),
		},
		&cli.StringFlag{
			Name:   "messages.content.title",
			Action: getAPIFlagAction[string]("body", "messages.#.content.#.title"),
		},
		&cli.StringFlag{
			Name:   "messages.content.signature",
			Action: getAPIFlagAction[string]("body", "messages.#.content.#.signature"),
		},
		&cli.StringFlag{
			Name:   "messages.content.thinking",
			Action: getAPIFlagAction[string]("body", "messages.#.content.#.thinking"),
		},
		&cli.StringFlag{
			Name:   "messages.content.data",
			Action: getAPIFlagAction[string]("body", "messages.#.content.#.data"),
		},
		&cli.BoolFlag{
			Name:   "messages.+content",
			Action: getAPIFlagActionWithValue[bool]("body", "messages.#.content.-1", map[string]interface{}{}),
		},
		&cli.StringFlag{
			Name:   "messages.role",
			Action: getAPIFlagAction[string]("body", "messages.#.role"),
		},
		&cli.BoolFlag{
			Name:   "+message",
			Action: getAPIFlagActionWithValue[bool]("body", "messages.-1", map[string]interface{}{}),
		},
		&cli.StringFlag{
			Name:   "model",
			Action: getAPIFlagAction[string]("body", "model"),
		},
		&cli.StringFlag{
			Name:   "system",
			Action: getAPIFlagAction[string]("body", "system"),
		},
		&cli.StringFlag{
			Name:   "system.text",
			Action: getAPIFlagAction[string]("body", "system.#.text"),
		},
		&cli.StringFlag{
			Name:   "system.type",
			Action: getAPIFlagAction[string]("body", "system.#.type"),
		},
		&cli.StringFlag{
			Name:   "system.cache_control.type",
			Action: getAPIFlagAction[string]("body", "system.#.cache_control.type"),
		},
		&cli.StringFlag{
			Name:   "system.citations.cited_text",
			Action: getAPIFlagAction[string]("body", "system.#.citations.#.cited_text"),
		},
		&cli.Int64Flag{
			Name:   "system.citations.document_index",
			Action: getAPIFlagAction[int64]("body", "system.#.citations.#.document_index"),
		},
		&cli.StringFlag{
			Name:   "system.citations.document_title",
			Action: getAPIFlagAction[string]("body", "system.#.citations.#.document_title"),
		},
		&cli.Int64Flag{
			Name:   "system.citations.end_char_index",
			Action: getAPIFlagAction[int64]("body", "system.#.citations.#.end_char_index"),
		},
		&cli.Int64Flag{
			Name:   "system.citations.start_char_index",
			Action: getAPIFlagAction[int64]("body", "system.#.citations.#.start_char_index"),
		},
		&cli.StringFlag{
			Name:   "system.citations.type",
			Action: getAPIFlagAction[string]("body", "system.#.citations.#.type"),
		},
		&cli.Int64Flag{
			Name:   "system.citations.end_page_number",
			Action: getAPIFlagAction[int64]("body", "system.#.citations.#.end_page_number"),
		},
		&cli.Int64Flag{
			Name:   "system.citations.start_page_number",
			Action: getAPIFlagAction[int64]("body", "system.#.citations.#.start_page_number"),
		},
		&cli.Int64Flag{
			Name:   "system.citations.end_block_index",
			Action: getAPIFlagAction[int64]("body", "system.#.citations.#.end_block_index"),
		},
		&cli.Int64Flag{
			Name:   "system.citations.start_block_index",
			Action: getAPIFlagAction[int64]("body", "system.#.citations.#.start_block_index"),
		},
		&cli.StringFlag{
			Name:   "system.citations.encrypted_index",
			Action: getAPIFlagAction[string]("body", "system.#.citations.#.encrypted_index"),
		},
		&cli.StringFlag{
			Name:   "system.citations.title",
			Action: getAPIFlagAction[string]("body", "system.#.citations.#.title"),
		},
		&cli.StringFlag{
			Name:   "system.citations.url",
			Action: getAPIFlagAction[string]("body", "system.#.citations.#.url"),
		},
		&cli.BoolFlag{
			Name:   "system.+citation",
			Action: getAPIFlagActionWithValue[bool]("body", "system.#.citations.-1", map[string]interface{}{}),
		},
		&cli.BoolFlag{
			Name:   "+system",
			Action: getAPIFlagActionWithValue[bool]("body", "system.-1", map[string]interface{}{}),
		},
		&cli.Int64Flag{
			Name:   "thinking.budget_tokens",
			Action: getAPIFlagAction[int64]("body", "thinking.budget_tokens"),
		},
		&cli.StringFlag{
			Name:   "thinking.type",
			Action: getAPIFlagAction[string]("body", "thinking.type"),
		},
		&cli.StringFlag{
			Name:   "tool-choice.type",
			Action: getAPIFlagAction[string]("body", "tool_choice.type"),
		},
		&cli.BoolFlag{
			Name:   "tool-choice.disable_parallel_tool_use",
			Action: getAPIFlagAction[bool]("body", "tool_choice.disable_parallel_tool_use"),
		},
		&cli.StringFlag{
			Name:   "tool-choice.name",
			Action: getAPIFlagAction[string]("body", "tool_choice.name"),
		},
		&cli.StringFlag{
			Name:   "tools.name",
			Action: getAPIFlagAction[string]("body", "tools.#.name"),
		},
		&cli.StringFlag{
			Name:   "tools.cache_control.type",
			Action: getAPIFlagAction[string]("body", "tools.#.cache_control.type"),
		},
		&cli.StringFlag{
			Name:   "tools.description",
			Action: getAPIFlagAction[string]("body", "tools.#.description"),
		},
		&cli.StringFlag{
			Name:   "tools.type",
			Action: getAPIFlagAction[string]("body", "tools.#.type"),
		},
		&cli.StringFlag{
			Name:   "tools.allowed_domains",
			Action: getAPIFlagAction[string]("body", "tools.#.allowed_domains.#"),
		},
		&cli.StringFlag{
			Name:   "tools.+allowed_domain",
			Action: getAPIFlagAction[string]("body", "tools.#.allowed_domains.-1"),
		},
		&cli.StringFlag{
			Name:   "tools.blocked_domains",
			Action: getAPIFlagAction[string]("body", "tools.#.blocked_domains.#"),
		},
		&cli.StringFlag{
			Name:   "tools.+blocked_domain",
			Action: getAPIFlagAction[string]("body", "tools.#.blocked_domains.-1"),
		},
		&cli.Int64Flag{
			Name:   "tools.max_uses",
			Action: getAPIFlagAction[int64]("body", "tools.#.max_uses"),
		},
		&cli.StringFlag{
			Name:   "tools.user_location.type",
			Action: getAPIFlagAction[string]("body", "tools.#.user_location.type"),
		},
		&cli.StringFlag{
			Name:   "tools.user_location.city",
			Action: getAPIFlagAction[string]("body", "tools.#.user_location.city"),
		},
		&cli.StringFlag{
			Name:   "tools.user_location.country",
			Action: getAPIFlagAction[string]("body", "tools.#.user_location.country"),
		},
		&cli.StringFlag{
			Name:   "tools.user_location.region",
			Action: getAPIFlagAction[string]("body", "tools.#.user_location.region"),
		},
		&cli.StringFlag{
			Name:   "tools.user_location.timezone",
			Action: getAPIFlagAction[string]("body", "tools.#.user_location.timezone"),
		},
		&cli.BoolFlag{
			Name:   "+tool",
			Action: getAPIFlagActionWithValue[bool]("body", "tools.-1", map[string]interface{}{}),
		},
	},
	Before:          initAPICommand,
	Action:          handleMessagesCountTokens,
	HideHelpCommand: true,
}

func handleMessagesCreate(ctx context.Context, cmd *cli.Command) error {
	cc := getAPICommandContext(ctx, cmd)
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
	cc := getAPICommandContext(ctx, cmd)
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
