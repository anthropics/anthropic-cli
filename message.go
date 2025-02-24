// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
)

func createMessagesCreateSubcommand(initialBody []byte) Subcommand {
	query := []byte("{}")
	header := []byte("{}")
	body := initialBody
	var flagSet = flag.NewFlagSet("messages.create", flag.ExitOnError)

	flagSet.Func(
		"max-tokens",
		"",
		func(string string) error {
			int, err := parseInt(string)
			if err != nil {
				return err
			}
			var jsonErr error
			body, jsonErr = jsonSet(body, "max_tokens", int)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"messages.content.text",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "messages.#.content.#.text", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"messages.content.type",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "messages.#.content.#.type", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"messages.content.cache_control.type",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "messages.#.content.#.cache_control.type", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"messages.content.citations.cited_text",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "messages.#.content.#.citations.#.cited_text", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"messages.content.citations.document_index",
		"",
		func(string string) error {
			int, err := parseInt(string)
			if err != nil {
				return err
			}
			var jsonErr error
			body, jsonErr = jsonSet(body, "messages.#.content.#.citations.#.document_index", int)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"messages.content.citations.document_title",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "messages.#.content.#.citations.#.document_title", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"messages.content.citations.end_char_index",
		"",
		func(string string) error {
			int, err := parseInt(string)
			if err != nil {
				return err
			}
			var jsonErr error
			body, jsonErr = jsonSet(body, "messages.#.content.#.citations.#.end_char_index", int)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"messages.content.citations.start_char_index",
		"",
		func(string string) error {
			int, err := parseInt(string)
			if err != nil {
				return err
			}
			var jsonErr error
			body, jsonErr = jsonSet(body, "messages.#.content.#.citations.#.start_char_index", int)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"messages.content.citations.type",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "messages.#.content.#.citations.#.type", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"messages.content.citations.end_page_number",
		"",
		func(string string) error {
			int, err := parseInt(string)
			if err != nil {
				return err
			}
			var jsonErr error
			body, jsonErr = jsonSet(body, "messages.#.content.#.citations.#.end_page_number", int)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"messages.content.citations.start_page_number",
		"",
		func(string string) error {
			int, err := parseInt(string)
			if err != nil {
				return err
			}
			var jsonErr error
			body, jsonErr = jsonSet(body, "messages.#.content.#.citations.#.start_page_number", int)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"messages.content.citations.end_block_index",
		"",
		func(string string) error {
			int, err := parseInt(string)
			if err != nil {
				return err
			}
			var jsonErr error
			body, jsonErr = jsonSet(body, "messages.#.content.#.citations.#.end_block_index", int)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"messages.content.citations.start_block_index",
		"",
		func(string string) error {
			int, err := parseInt(string)
			if err != nil {
				return err
			}
			var jsonErr error
			body, jsonErr = jsonSet(body, "messages.#.content.#.citations.#.start_block_index", int)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.BoolFunc(
		"messages.content.+citation",
		"",
		func(_ string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "messages.#.content.#.citations.-1", map[string]interface{}{})
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"messages.content.source.data",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "messages.#.content.#.source.data", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"messages.content.source.media_type",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "messages.#.content.#.source.media_type", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"messages.content.source.type",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "messages.#.content.#.source.type", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"messages.content.id",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "messages.#.content.#.id", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"messages.content.name",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "messages.#.content.#.name", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"messages.content.tool_use_id",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "messages.#.content.#.tool_use_id", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"messages.content.content.text",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "messages.#.content.#.content.#.text", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"messages.content.content.type",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "messages.#.content.#.content.#.type", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"messages.content.content.cache_control.type",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "messages.#.content.#.content.#.cache_control.type", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"messages.content.content.citations.cited_text",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "messages.#.content.#.content.#.citations.#.cited_text", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"messages.content.content.citations.document_index",
		"",
		func(string string) error {
			int, err := parseInt(string)
			if err != nil {
				return err
			}
			var jsonErr error
			body, jsonErr = jsonSet(body, "messages.#.content.#.content.#.citations.#.document_index", int)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"messages.content.content.citations.document_title",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "messages.#.content.#.content.#.citations.#.document_title", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"messages.content.content.citations.end_char_index",
		"",
		func(string string) error {
			int, err := parseInt(string)
			if err != nil {
				return err
			}
			var jsonErr error
			body, jsonErr = jsonSet(body, "messages.#.content.#.content.#.citations.#.end_char_index", int)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"messages.content.content.citations.start_char_index",
		"",
		func(string string) error {
			int, err := parseInt(string)
			if err != nil {
				return err
			}
			var jsonErr error
			body, jsonErr = jsonSet(body, "messages.#.content.#.content.#.citations.#.start_char_index", int)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"messages.content.content.citations.type",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "messages.#.content.#.content.#.citations.#.type", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"messages.content.content.citations.end_page_number",
		"",
		func(string string) error {
			int, err := parseInt(string)
			if err != nil {
				return err
			}
			var jsonErr error
			body, jsonErr = jsonSet(body, "messages.#.content.#.content.#.citations.#.end_page_number", int)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"messages.content.content.citations.start_page_number",
		"",
		func(string string) error {
			int, err := parseInt(string)
			if err != nil {
				return err
			}
			var jsonErr error
			body, jsonErr = jsonSet(body, "messages.#.content.#.content.#.citations.#.start_page_number", int)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"messages.content.content.citations.end_block_index",
		"",
		func(string string) error {
			int, err := parseInt(string)
			if err != nil {
				return err
			}
			var jsonErr error
			body, jsonErr = jsonSet(body, "messages.#.content.#.content.#.citations.#.end_block_index", int)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"messages.content.content.citations.start_block_index",
		"",
		func(string string) error {
			int, err := parseInt(string)
			if err != nil {
				return err
			}
			var jsonErr error
			body, jsonErr = jsonSet(body, "messages.#.content.#.content.#.citations.#.start_block_index", int)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.BoolFunc(
		"messages.content.content.+citation",
		"",
		func(_ string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "messages.#.content.#.content.#.citations.-1", map[string]interface{}{})
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"messages.content.content.source.data",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "messages.#.content.#.content.#.source.data", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"messages.content.content.source.media_type",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "messages.#.content.#.content.#.source.media_type", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"messages.content.content.source.type",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "messages.#.content.#.content.#.source.type", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.BoolFunc(
		"messages.content.+content",
		"",
		func(_ string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "messages.#.content.#.content.-1", map[string]interface{}{})
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.BoolFunc(
		"messages.content.is_error",
		"",
		func(_ string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "messages.#.content.#.is_error", true)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"messages.content.source.content",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "messages.#.content.#.source.content", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"messages.content.source.content.text",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "messages.#.content.#.source.content.#.text", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"messages.content.source.content.type",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "messages.#.content.#.source.content.#.type", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"messages.content.source.content.cache_control.type",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "messages.#.content.#.source.content.#.cache_control.type", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"messages.content.source.content.citations.cited_text",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "messages.#.content.#.source.content.#.citations.#.cited_text", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"messages.content.source.content.citations.document_index",
		"",
		func(string string) error {
			int, err := parseInt(string)
			if err != nil {
				return err
			}
			var jsonErr error
			body, jsonErr = jsonSet(body, "messages.#.content.#.source.content.#.citations.#.document_index", int)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"messages.content.source.content.citations.document_title",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "messages.#.content.#.source.content.#.citations.#.document_title", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"messages.content.source.content.citations.end_char_index",
		"",
		func(string string) error {
			int, err := parseInt(string)
			if err != nil {
				return err
			}
			var jsonErr error
			body, jsonErr = jsonSet(body, "messages.#.content.#.source.content.#.citations.#.end_char_index", int)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"messages.content.source.content.citations.start_char_index",
		"",
		func(string string) error {
			int, err := parseInt(string)
			if err != nil {
				return err
			}
			var jsonErr error
			body, jsonErr = jsonSet(body, "messages.#.content.#.source.content.#.citations.#.start_char_index", int)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"messages.content.source.content.citations.type",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "messages.#.content.#.source.content.#.citations.#.type", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"messages.content.source.content.citations.end_page_number",
		"",
		func(string string) error {
			int, err := parseInt(string)
			if err != nil {
				return err
			}
			var jsonErr error
			body, jsonErr = jsonSet(body, "messages.#.content.#.source.content.#.citations.#.end_page_number", int)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"messages.content.source.content.citations.start_page_number",
		"",
		func(string string) error {
			int, err := parseInt(string)
			if err != nil {
				return err
			}
			var jsonErr error
			body, jsonErr = jsonSet(body, "messages.#.content.#.source.content.#.citations.#.start_page_number", int)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"messages.content.source.content.citations.end_block_index",
		"",
		func(string string) error {
			int, err := parseInt(string)
			if err != nil {
				return err
			}
			var jsonErr error
			body, jsonErr = jsonSet(body, "messages.#.content.#.source.content.#.citations.#.end_block_index", int)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"messages.content.source.content.citations.start_block_index",
		"",
		func(string string) error {
			int, err := parseInt(string)
			if err != nil {
				return err
			}
			var jsonErr error
			body, jsonErr = jsonSet(body, "messages.#.content.#.source.content.#.citations.#.start_block_index", int)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.BoolFunc(
		"messages.content.source.content.+citation",
		"",
		func(_ string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "messages.#.content.#.source.content.#.citations.-1", map[string]interface{}{})
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"messages.content.source.content.source.data",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "messages.#.content.#.source.content.#.source.data", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"messages.content.source.content.source.media_type",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "messages.#.content.#.source.content.#.source.media_type", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"messages.content.source.content.source.type",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "messages.#.content.#.source.content.#.source.type", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.BoolFunc(
		"messages.content.source.+content",
		"",
		func(_ string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "messages.#.content.#.source.content.-1", map[string]interface{}{})
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.BoolFunc(
		"messages.content.citations.enabled",
		"",
		func(_ string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "messages.#.content.#.citations.enabled", true)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"messages.content.context",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "messages.#.content.#.context", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"messages.content.title",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "messages.#.content.#.title", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"messages.content.signature",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "messages.#.content.#.signature", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"messages.content.thinking",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "messages.#.content.#.thinking", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"messages.content.data",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "messages.#.content.#.data", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.BoolFunc(
		"messages.+content",
		"",
		func(_ string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "messages.#.content.-1", map[string]interface{}{})
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"messages.role",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "messages.#.role", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.BoolFunc(
		"+message",
		"",
		func(_ string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "messages.-1", map[string]interface{}{})
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"model",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "model", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"metadata.user_id",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "metadata.user_id", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"stop-sequences",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "stop_sequences.#", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"+stop_sequence",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "stop_sequences.-1", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"system.text",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "system.#.text", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"system.type",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "system.#.type", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"system.cache_control.type",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "system.#.cache_control.type", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"system.citations.cited_text",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "system.#.citations.#.cited_text", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"system.citations.document_index",
		"",
		func(string string) error {
			int, err := parseInt(string)
			if err != nil {
				return err
			}
			var jsonErr error
			body, jsonErr = jsonSet(body, "system.#.citations.#.document_index", int)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"system.citations.document_title",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "system.#.citations.#.document_title", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"system.citations.end_char_index",
		"",
		func(string string) error {
			int, err := parseInt(string)
			if err != nil {
				return err
			}
			var jsonErr error
			body, jsonErr = jsonSet(body, "system.#.citations.#.end_char_index", int)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"system.citations.start_char_index",
		"",
		func(string string) error {
			int, err := parseInt(string)
			if err != nil {
				return err
			}
			var jsonErr error
			body, jsonErr = jsonSet(body, "system.#.citations.#.start_char_index", int)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"system.citations.type",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "system.#.citations.#.type", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"system.citations.end_page_number",
		"",
		func(string string) error {
			int, err := parseInt(string)
			if err != nil {
				return err
			}
			var jsonErr error
			body, jsonErr = jsonSet(body, "system.#.citations.#.end_page_number", int)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"system.citations.start_page_number",
		"",
		func(string string) error {
			int, err := parseInt(string)
			if err != nil {
				return err
			}
			var jsonErr error
			body, jsonErr = jsonSet(body, "system.#.citations.#.start_page_number", int)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"system.citations.end_block_index",
		"",
		func(string string) error {
			int, err := parseInt(string)
			if err != nil {
				return err
			}
			var jsonErr error
			body, jsonErr = jsonSet(body, "system.#.citations.#.end_block_index", int)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"system.citations.start_block_index",
		"",
		func(string string) error {
			int, err := parseInt(string)
			if err != nil {
				return err
			}
			var jsonErr error
			body, jsonErr = jsonSet(body, "system.#.citations.#.start_block_index", int)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.BoolFunc(
		"system.+citation",
		"",
		func(_ string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "system.#.citations.-1", map[string]interface{}{})
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.BoolFunc(
		"+system",
		"",
		func(_ string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "system.-1", map[string]interface{}{})
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"temperature",
		"",
		func(string string) error {
			float, err := parseFloat(string)
			if err != nil {
				return err
			}
			var jsonErr error
			body, jsonErr = jsonSet(body, "temperature", float)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"thinking.budget_tokens",
		"",
		func(string string) error {
			int, err := parseInt(string)
			if err != nil {
				return err
			}
			var jsonErr error
			body, jsonErr = jsonSet(body, "thinking.budget_tokens", int)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"thinking.type",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "thinking.type", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"tool-choice.type",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "tool_choice.type", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.BoolFunc(
		"tool-choice.disable_parallel_tool_use",
		"",
		func(_ string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "tool_choice.disable_parallel_tool_use", true)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"tool-choice.name",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "tool_choice.name", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"tools.name",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "tools.#.name", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"tools.type",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "tools.#.type", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"tools.cache_control.type",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "tools.#.cache_control.type", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"tools.description",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "tools.#.description", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.BoolFunc(
		"+tool",
		"",
		func(_ string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "tools.-1", map[string]interface{}{})
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"top-k",
		"",
		func(string string) error {
			int, err := parseInt(string)
			if err != nil {
				return err
			}
			var jsonErr error
			body, jsonErr = jsonSet(body, "top_k", int)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"top-p",
		"",
		func(string string) error {
			float, err := parseFloat(string)
			if err != nil {
				return err
			}
			var jsonErr error
			body, jsonErr = jsonSet(body, "top_p", float)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	return Subcommand{
		flagSet: flagSet,
		handle: func(client *anthropic.Client) {
			var err error
			body, err = jsonSet(body, "stream", true)
			if err != nil {
				fmt.Printf("%s\n", err)
				os.Exit(1)
			}
			stream := client.Messages.NewStreaming(
				context.TODO(),
				anthropic.MessageNewParams{},
				option.WithMiddleware(func(r *http.Request, mn option.MiddlewareNext) (*http.Response, error) {
					r.URL.RawQuery = serializeQuery(query).Encode()
					r.Header = serializeHeader(header)
					return mn(r)
				}),
				option.WithRequestBody("application/json", body),
			)
			for stream.Next() {
				fmt.Printf("%s\n", stream.Current().JSON.RawJSON())
			}
		},
	}
}

func createMessagesCountTokensSubcommand(initialBody []byte) Subcommand {
	query := []byte("{}")
	header := []byte("{}")
	body := initialBody
	var flagSet = flag.NewFlagSet("messages.count_tokens", flag.ExitOnError)

	flagSet.Func(
		"messages.content.text",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "messages.#.content.#.text", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"messages.content.type",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "messages.#.content.#.type", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"messages.content.cache_control.type",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "messages.#.content.#.cache_control.type", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"messages.content.citations.cited_text",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "messages.#.content.#.citations.#.cited_text", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"messages.content.citations.document_index",
		"",
		func(string string) error {
			int, err := parseInt(string)
			if err != nil {
				return err
			}
			var jsonErr error
			body, jsonErr = jsonSet(body, "messages.#.content.#.citations.#.document_index", int)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"messages.content.citations.document_title",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "messages.#.content.#.citations.#.document_title", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"messages.content.citations.end_char_index",
		"",
		func(string string) error {
			int, err := parseInt(string)
			if err != nil {
				return err
			}
			var jsonErr error
			body, jsonErr = jsonSet(body, "messages.#.content.#.citations.#.end_char_index", int)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"messages.content.citations.start_char_index",
		"",
		func(string string) error {
			int, err := parseInt(string)
			if err != nil {
				return err
			}
			var jsonErr error
			body, jsonErr = jsonSet(body, "messages.#.content.#.citations.#.start_char_index", int)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"messages.content.citations.type",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "messages.#.content.#.citations.#.type", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"messages.content.citations.end_page_number",
		"",
		func(string string) error {
			int, err := parseInt(string)
			if err != nil {
				return err
			}
			var jsonErr error
			body, jsonErr = jsonSet(body, "messages.#.content.#.citations.#.end_page_number", int)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"messages.content.citations.start_page_number",
		"",
		func(string string) error {
			int, err := parseInt(string)
			if err != nil {
				return err
			}
			var jsonErr error
			body, jsonErr = jsonSet(body, "messages.#.content.#.citations.#.start_page_number", int)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"messages.content.citations.end_block_index",
		"",
		func(string string) error {
			int, err := parseInt(string)
			if err != nil {
				return err
			}
			var jsonErr error
			body, jsonErr = jsonSet(body, "messages.#.content.#.citations.#.end_block_index", int)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"messages.content.citations.start_block_index",
		"",
		func(string string) error {
			int, err := parseInt(string)
			if err != nil {
				return err
			}
			var jsonErr error
			body, jsonErr = jsonSet(body, "messages.#.content.#.citations.#.start_block_index", int)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.BoolFunc(
		"messages.content.+citation",
		"",
		func(_ string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "messages.#.content.#.citations.-1", map[string]interface{}{})
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"messages.content.source.data",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "messages.#.content.#.source.data", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"messages.content.source.media_type",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "messages.#.content.#.source.media_type", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"messages.content.source.type",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "messages.#.content.#.source.type", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"messages.content.id",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "messages.#.content.#.id", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"messages.content.name",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "messages.#.content.#.name", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"messages.content.tool_use_id",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "messages.#.content.#.tool_use_id", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"messages.content.content.text",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "messages.#.content.#.content.#.text", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"messages.content.content.type",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "messages.#.content.#.content.#.type", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"messages.content.content.cache_control.type",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "messages.#.content.#.content.#.cache_control.type", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"messages.content.content.citations.cited_text",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "messages.#.content.#.content.#.citations.#.cited_text", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"messages.content.content.citations.document_index",
		"",
		func(string string) error {
			int, err := parseInt(string)
			if err != nil {
				return err
			}
			var jsonErr error
			body, jsonErr = jsonSet(body, "messages.#.content.#.content.#.citations.#.document_index", int)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"messages.content.content.citations.document_title",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "messages.#.content.#.content.#.citations.#.document_title", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"messages.content.content.citations.end_char_index",
		"",
		func(string string) error {
			int, err := parseInt(string)
			if err != nil {
				return err
			}
			var jsonErr error
			body, jsonErr = jsonSet(body, "messages.#.content.#.content.#.citations.#.end_char_index", int)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"messages.content.content.citations.start_char_index",
		"",
		func(string string) error {
			int, err := parseInt(string)
			if err != nil {
				return err
			}
			var jsonErr error
			body, jsonErr = jsonSet(body, "messages.#.content.#.content.#.citations.#.start_char_index", int)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"messages.content.content.citations.type",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "messages.#.content.#.content.#.citations.#.type", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"messages.content.content.citations.end_page_number",
		"",
		func(string string) error {
			int, err := parseInt(string)
			if err != nil {
				return err
			}
			var jsonErr error
			body, jsonErr = jsonSet(body, "messages.#.content.#.content.#.citations.#.end_page_number", int)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"messages.content.content.citations.start_page_number",
		"",
		func(string string) error {
			int, err := parseInt(string)
			if err != nil {
				return err
			}
			var jsonErr error
			body, jsonErr = jsonSet(body, "messages.#.content.#.content.#.citations.#.start_page_number", int)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"messages.content.content.citations.end_block_index",
		"",
		func(string string) error {
			int, err := parseInt(string)
			if err != nil {
				return err
			}
			var jsonErr error
			body, jsonErr = jsonSet(body, "messages.#.content.#.content.#.citations.#.end_block_index", int)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"messages.content.content.citations.start_block_index",
		"",
		func(string string) error {
			int, err := parseInt(string)
			if err != nil {
				return err
			}
			var jsonErr error
			body, jsonErr = jsonSet(body, "messages.#.content.#.content.#.citations.#.start_block_index", int)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.BoolFunc(
		"messages.content.content.+citation",
		"",
		func(_ string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "messages.#.content.#.content.#.citations.-1", map[string]interface{}{})
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"messages.content.content.source.data",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "messages.#.content.#.content.#.source.data", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"messages.content.content.source.media_type",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "messages.#.content.#.content.#.source.media_type", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"messages.content.content.source.type",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "messages.#.content.#.content.#.source.type", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.BoolFunc(
		"messages.content.+content",
		"",
		func(_ string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "messages.#.content.#.content.-1", map[string]interface{}{})
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.BoolFunc(
		"messages.content.is_error",
		"",
		func(_ string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "messages.#.content.#.is_error", true)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"messages.content.source.content",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "messages.#.content.#.source.content", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"messages.content.source.content.text",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "messages.#.content.#.source.content.#.text", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"messages.content.source.content.type",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "messages.#.content.#.source.content.#.type", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"messages.content.source.content.cache_control.type",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "messages.#.content.#.source.content.#.cache_control.type", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"messages.content.source.content.citations.cited_text",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "messages.#.content.#.source.content.#.citations.#.cited_text", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"messages.content.source.content.citations.document_index",
		"",
		func(string string) error {
			int, err := parseInt(string)
			if err != nil {
				return err
			}
			var jsonErr error
			body, jsonErr = jsonSet(body, "messages.#.content.#.source.content.#.citations.#.document_index", int)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"messages.content.source.content.citations.document_title",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "messages.#.content.#.source.content.#.citations.#.document_title", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"messages.content.source.content.citations.end_char_index",
		"",
		func(string string) error {
			int, err := parseInt(string)
			if err != nil {
				return err
			}
			var jsonErr error
			body, jsonErr = jsonSet(body, "messages.#.content.#.source.content.#.citations.#.end_char_index", int)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"messages.content.source.content.citations.start_char_index",
		"",
		func(string string) error {
			int, err := parseInt(string)
			if err != nil {
				return err
			}
			var jsonErr error
			body, jsonErr = jsonSet(body, "messages.#.content.#.source.content.#.citations.#.start_char_index", int)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"messages.content.source.content.citations.type",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "messages.#.content.#.source.content.#.citations.#.type", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"messages.content.source.content.citations.end_page_number",
		"",
		func(string string) error {
			int, err := parseInt(string)
			if err != nil {
				return err
			}
			var jsonErr error
			body, jsonErr = jsonSet(body, "messages.#.content.#.source.content.#.citations.#.end_page_number", int)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"messages.content.source.content.citations.start_page_number",
		"",
		func(string string) error {
			int, err := parseInt(string)
			if err != nil {
				return err
			}
			var jsonErr error
			body, jsonErr = jsonSet(body, "messages.#.content.#.source.content.#.citations.#.start_page_number", int)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"messages.content.source.content.citations.end_block_index",
		"",
		func(string string) error {
			int, err := parseInt(string)
			if err != nil {
				return err
			}
			var jsonErr error
			body, jsonErr = jsonSet(body, "messages.#.content.#.source.content.#.citations.#.end_block_index", int)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"messages.content.source.content.citations.start_block_index",
		"",
		func(string string) error {
			int, err := parseInt(string)
			if err != nil {
				return err
			}
			var jsonErr error
			body, jsonErr = jsonSet(body, "messages.#.content.#.source.content.#.citations.#.start_block_index", int)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.BoolFunc(
		"messages.content.source.content.+citation",
		"",
		func(_ string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "messages.#.content.#.source.content.#.citations.-1", map[string]interface{}{})
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"messages.content.source.content.source.data",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "messages.#.content.#.source.content.#.source.data", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"messages.content.source.content.source.media_type",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "messages.#.content.#.source.content.#.source.media_type", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"messages.content.source.content.source.type",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "messages.#.content.#.source.content.#.source.type", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.BoolFunc(
		"messages.content.source.+content",
		"",
		func(_ string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "messages.#.content.#.source.content.-1", map[string]interface{}{})
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.BoolFunc(
		"messages.content.citations.enabled",
		"",
		func(_ string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "messages.#.content.#.citations.enabled", true)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"messages.content.context",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "messages.#.content.#.context", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"messages.content.title",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "messages.#.content.#.title", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"messages.content.signature",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "messages.#.content.#.signature", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"messages.content.thinking",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "messages.#.content.#.thinking", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"messages.content.data",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "messages.#.content.#.data", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.BoolFunc(
		"messages.+content",
		"",
		func(_ string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "messages.#.content.-1", map[string]interface{}{})
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"messages.role",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "messages.#.role", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.BoolFunc(
		"+message",
		"",
		func(_ string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "messages.-1", map[string]interface{}{})
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"model",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "model", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"system",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "system", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"system.text",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "system.#.text", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"system.type",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "system.#.type", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"system.cache_control.type",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "system.#.cache_control.type", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"system.citations.cited_text",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "system.#.citations.#.cited_text", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"system.citations.document_index",
		"",
		func(string string) error {
			int, err := parseInt(string)
			if err != nil {
				return err
			}
			var jsonErr error
			body, jsonErr = jsonSet(body, "system.#.citations.#.document_index", int)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"system.citations.document_title",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "system.#.citations.#.document_title", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"system.citations.end_char_index",
		"",
		func(string string) error {
			int, err := parseInt(string)
			if err != nil {
				return err
			}
			var jsonErr error
			body, jsonErr = jsonSet(body, "system.#.citations.#.end_char_index", int)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"system.citations.start_char_index",
		"",
		func(string string) error {
			int, err := parseInt(string)
			if err != nil {
				return err
			}
			var jsonErr error
			body, jsonErr = jsonSet(body, "system.#.citations.#.start_char_index", int)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"system.citations.type",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "system.#.citations.#.type", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"system.citations.end_page_number",
		"",
		func(string string) error {
			int, err := parseInt(string)
			if err != nil {
				return err
			}
			var jsonErr error
			body, jsonErr = jsonSet(body, "system.#.citations.#.end_page_number", int)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"system.citations.start_page_number",
		"",
		func(string string) error {
			int, err := parseInt(string)
			if err != nil {
				return err
			}
			var jsonErr error
			body, jsonErr = jsonSet(body, "system.#.citations.#.start_page_number", int)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"system.citations.end_block_index",
		"",
		func(string string) error {
			int, err := parseInt(string)
			if err != nil {
				return err
			}
			var jsonErr error
			body, jsonErr = jsonSet(body, "system.#.citations.#.end_block_index", int)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"system.citations.start_block_index",
		"",
		func(string string) error {
			int, err := parseInt(string)
			if err != nil {
				return err
			}
			var jsonErr error
			body, jsonErr = jsonSet(body, "system.#.citations.#.start_block_index", int)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.BoolFunc(
		"system.+citation",
		"",
		func(_ string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "system.#.citations.-1", map[string]interface{}{})
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.BoolFunc(
		"+system",
		"",
		func(_ string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "system.-1", map[string]interface{}{})
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"thinking.budget_tokens",
		"",
		func(string string) error {
			int, err := parseInt(string)
			if err != nil {
				return err
			}
			var jsonErr error
			body, jsonErr = jsonSet(body, "thinking.budget_tokens", int)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"thinking.type",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "thinking.type", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"tool-choice.type",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "tool_choice.type", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.BoolFunc(
		"tool-choice.disable_parallel_tool_use",
		"",
		func(_ string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "tool_choice.disable_parallel_tool_use", true)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"tool-choice.name",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "tool_choice.name", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"tools.name",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "tools.#.name", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"tools.type",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "tools.#.type", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"tools.cache_control.type",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "tools.#.cache_control.type", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"tools.description",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "tools.#.description", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.BoolFunc(
		"+tool",
		"",
		func(_ string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "tools.-1", map[string]interface{}{})
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	return Subcommand{
		flagSet: flagSet,
		handle: func(client *anthropic.Client) {
			res, err := client.Messages.CountTokens(
				context.TODO(),
				anthropic.MessageCountTokensParams{},
				option.WithMiddleware(func(r *http.Request, mn option.MiddlewareNext) (*http.Response, error) {
					r.URL.RawQuery = serializeQuery(query).Encode()
					r.Header = serializeHeader(header)
					return mn(r)
				}),
				option.WithRequestBody("application/json", body),
			)
			if err != nil {
				fmt.Printf("%s\n", err)
				os.Exit(1)
			}

			fmt.Printf("%s\n", res.JSON.RawJSON())
		},
	}
}
