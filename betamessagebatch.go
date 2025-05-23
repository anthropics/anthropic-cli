// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package main

import (
	"context"
	"fmt"
	"os"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/urfave/cli/v3"
)

var betaMessagesBatchesCreate = cli.Command{
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
			Name:   "requests.params.messages.content.cache_control.ttl",
			Action: getAPIFlagAction[string]("body", "requests.#.params.messages.#.content.#.cache_control.ttl"),
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
			Name:   "requests.params.messages.content.content.content.file_id",
			Action: getAPIFlagAction[string]("body", "requests.#.params.messages.#.content.#.content.content.#.file_id"),
		},
		&cli.StringFlag{
			Name:   "requests.params.messages.content.content.content.type",
			Action: getAPIFlagAction[string]("body", "requests.#.params.messages.#.content.#.content.content.#.type"),
		},
		&cli.BoolFlag{
			Name:   "requests.params.messages.content.content.+content",
			Action: getAPIFlagActionWithValue[bool]("body", "requests.#.params.messages.#.content.#.content.content.-1", map[string]interface{}{}),
		},
		&cli.Int64Flag{
			Name:   "requests.params.messages.content.content.return_code",
			Action: getAPIFlagAction[int64]("body", "requests.#.params.messages.#.content.#.content.return_code"),
		},
		&cli.StringFlag{
			Name:   "requests.params.messages.content.content.stderr",
			Action: getAPIFlagAction[string]("body", "requests.#.params.messages.#.content.#.content.stderr"),
		},
		&cli.StringFlag{
			Name:   "requests.params.messages.content.content.stdout",
			Action: getAPIFlagAction[string]("body", "requests.#.params.messages.#.content.#.content.stdout"),
		},
		&cli.StringFlag{
			Name:   "requests.params.messages.content.server_name",
			Action: getAPIFlagAction[string]("body", "requests.#.params.messages.#.content.#.server_name"),
		},
		&cli.StringFlag{
			Name:   "requests.params.messages.content.content",
			Action: getAPIFlagAction[string]("body", "requests.#.params.messages.#.content.#.content"),
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
			Name:   "requests.params.messages.content.content.cache_control.ttl",
			Action: getAPIFlagAction[string]("body", "requests.#.params.messages.#.content.#.content.#.cache_control.ttl"),
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
		&cli.BoolFlag{
			Name:   "requests.params.messages.content.is_error",
			Action: getAPIFlagAction[bool]("body", "requests.#.params.messages.#.content.#.is_error"),
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
			Name:   "requests.params.messages.content.source.file_id",
			Action: getAPIFlagAction[string]("body", "requests.#.params.messages.#.content.#.source.file_id"),
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
		&cli.StringFlag{
			Name:   "requests.params.messages.content.content.source.file_id",
			Action: getAPIFlagAction[string]("body", "requests.#.params.messages.#.content.#.content.#.source.file_id"),
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
			Name:   "requests.params.messages.content.source.content.cache_control.ttl",
			Action: getAPIFlagAction[string]("body", "requests.#.params.messages.#.content.#.source.content.#.cache_control.ttl"),
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
		&cli.StringFlag{
			Name:   "requests.params.messages.content.source.content.source.file_id",
			Action: getAPIFlagAction[string]("body", "requests.#.params.messages.#.content.#.source.content.#.source.file_id"),
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
		&cli.StringFlag{
			Name:   "requests.params.messages.content.file_id",
			Action: getAPIFlagAction[string]("body", "requests.#.params.messages.#.content.#.file_id"),
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
			Name:   "requests.params.container",
			Action: getAPIFlagAction[string]("body", "requests.#.params.container"),
		},
		&cli.StringFlag{
			Name:   "requests.params.mcp_servers.name",
			Action: getAPIFlagAction[string]("body", "requests.#.params.mcp_servers.#.name"),
		},
		&cli.StringFlag{
			Name:   "requests.params.mcp_servers.type",
			Action: getAPIFlagAction[string]("body", "requests.#.params.mcp_servers.#.type"),
		},
		&cli.StringFlag{
			Name:   "requests.params.mcp_servers.url",
			Action: getAPIFlagAction[string]("body", "requests.#.params.mcp_servers.#.url"),
		},
		&cli.StringFlag{
			Name:   "requests.params.mcp_servers.authorization_token",
			Action: getAPIFlagAction[string]("body", "requests.#.params.mcp_servers.#.authorization_token"),
		},
		&cli.StringFlag{
			Name:   "requests.params.mcp_servers.tool_configuration.allowed_tools",
			Action: getAPIFlagAction[string]("body", "requests.#.params.mcp_servers.#.tool_configuration.allowed_tools.#"),
		},
		&cli.StringFlag{
			Name:   "requests.params.mcp_servers.tool_configuration.+allowed_tool",
			Action: getAPIFlagAction[string]("body", "requests.#.params.mcp_servers.#.tool_configuration.allowed_tools.-1"),
		},
		&cli.BoolFlag{
			Name:   "requests.params.mcp_servers.tool_configuration.enabled",
			Action: getAPIFlagAction[bool]("body", "requests.#.params.mcp_servers.#.tool_configuration.enabled"),
		},
		&cli.BoolFlag{
			Name:   "requests.params.+mcp_server",
			Action: getAPIFlagActionWithValue[bool]("body", "requests.#.params.mcp_servers.-1", map[string]interface{}{}),
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
			Name:   "requests.params.system.cache_control.ttl",
			Action: getAPIFlagAction[string]("body", "requests.#.params.system.#.cache_control.ttl"),
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
			Name:   "requests.params.tools.input_schema.type",
			Action: getAPIFlagAction[string]("body", "requests.#.params.tools.#.input_schema.type"),
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
			Name:   "requests.params.tools.cache_control.ttl",
			Action: getAPIFlagAction[string]("body", "requests.#.params.tools.#.cache_control.ttl"),
		},
		&cli.StringFlag{
			Name:   "requests.params.tools.description",
			Action: getAPIFlagAction[string]("body", "requests.#.params.tools.#.description"),
		},
		&cli.StringFlag{
			Name:   "requests.params.tools.type",
			Action: getAPIFlagAction[string]("body", "requests.#.params.tools.#.type"),
		},
		&cli.Int64Flag{
			Name:   "requests.params.tools.display_height_px",
			Action: getAPIFlagAction[int64]("body", "requests.#.params.tools.#.display_height_px"),
		},
		&cli.Int64Flag{
			Name:   "requests.params.tools.display_width_px",
			Action: getAPIFlagAction[int64]("body", "requests.#.params.tools.#.display_width_px"),
		},
		&cli.Int64Flag{
			Name:   "requests.params.tools.display_number",
			Action: getAPIFlagAction[int64]("body", "requests.#.params.tools.#.display_number"),
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
	Action:          handleBetaMessagesBatchesRetrieve,
	HideHelpCommand: true,
}

var betaMessagesBatchesList = cli.Command{
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
	Action:          handleBetaMessagesBatchesResults,
	HideHelpCommand: true,
}

func handleBetaMessagesBatchesCreate(ctx context.Context, cmd *cli.Command) error {
	cc := getAPICommandContext(ctx, cmd)
	params := anthropic.BetaMessageBatchNewParams{}
	res, err := cc.client.Beta.Messages.Batches.New(
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

func handleBetaMessagesBatchesRetrieve(ctx context.Context, cmd *cli.Command) error {
	cc := getAPICommandContext(ctx, cmd)
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

	fmt.Printf("%s\n", colorizeJSON(res.RawJSON(), os.Stdout))
	return nil
}

func handleBetaMessagesBatchesList(ctx context.Context, cmd *cli.Command) error {
	cc := getAPICommandContext(ctx, cmd)
	params := anthropic.BetaMessageBatchListParams{}
	res, err := cc.client.Beta.Messages.Batches.List(
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

func handleBetaMessagesBatchesDelete(ctx context.Context, cmd *cli.Command) error {
	cc := getAPICommandContext(ctx, cmd)
	params := anthropic.BetaMessageBatchDeleteParams{}
	res, err := cc.client.Beta.Messages.Batches.Delete(
		context.TODO(),
		cmd.Value("message-batch-id").(string),
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

func handleBetaMessagesBatchesCancel(ctx context.Context, cmd *cli.Command) error {
	cc := getAPICommandContext(ctx, cmd)
	params := anthropic.BetaMessageBatchCancelParams{}
	res, err := cc.client.Beta.Messages.Batches.Cancel(
		context.TODO(),
		cmd.Value("message-batch-id").(string),
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

func handleBetaMessagesBatchesResults(ctx context.Context, cmd *cli.Command) error {
	cc := getAPICommandContext(ctx, cmd)
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
