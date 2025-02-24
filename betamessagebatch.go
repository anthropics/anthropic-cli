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

func createBetaMessagesBatchesCreateSubcommand(initialBody []byte) Subcommand {
	query := []byte("{}")
	header := []byte("{}")
	body := initialBody
	var flagSet = flag.NewFlagSet("beta.messages.batches.create", flag.ExitOnError)

	flagSet.Func(
		"requests.custom_id",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "requests.#.custom_id", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"requests.params.max_tokens",
		"",
		func(string string) error {
			int, err := parseInt(string)
			if err != nil {
				return err
			}
			var jsonErr error
			body, jsonErr = jsonSet(body, "requests.#.params.max_tokens", int)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"requests.params.messages.content.text",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "requests.#.params.messages.#.content.#.text", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"requests.params.messages.content.type",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "requests.#.params.messages.#.content.#.type", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"requests.params.messages.content.cache_control.type",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "requests.#.params.messages.#.content.#.cache_control.type", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"requests.params.messages.content.citations.cited_text",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "requests.#.params.messages.#.content.#.citations.#.cited_text", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"requests.params.messages.content.citations.document_index",
		"",
		func(string string) error {
			int, err := parseInt(string)
			if err != nil {
				return err
			}
			var jsonErr error
			body, jsonErr = jsonSet(body, "requests.#.params.messages.#.content.#.citations.#.document_index", int)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"requests.params.messages.content.citations.document_title",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "requests.#.params.messages.#.content.#.citations.#.document_title", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"requests.params.messages.content.citations.end_char_index",
		"",
		func(string string) error {
			int, err := parseInt(string)
			if err != nil {
				return err
			}
			var jsonErr error
			body, jsonErr = jsonSet(body, "requests.#.params.messages.#.content.#.citations.#.end_char_index", int)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"requests.params.messages.content.citations.start_char_index",
		"",
		func(string string) error {
			int, err := parseInt(string)
			if err != nil {
				return err
			}
			var jsonErr error
			body, jsonErr = jsonSet(body, "requests.#.params.messages.#.content.#.citations.#.start_char_index", int)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"requests.params.messages.content.citations.type",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "requests.#.params.messages.#.content.#.citations.#.type", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"requests.params.messages.content.citations.end_page_number",
		"",
		func(string string) error {
			int, err := parseInt(string)
			if err != nil {
				return err
			}
			var jsonErr error
			body, jsonErr = jsonSet(body, "requests.#.params.messages.#.content.#.citations.#.end_page_number", int)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"requests.params.messages.content.citations.start_page_number",
		"",
		func(string string) error {
			int, err := parseInt(string)
			if err != nil {
				return err
			}
			var jsonErr error
			body, jsonErr = jsonSet(body, "requests.#.params.messages.#.content.#.citations.#.start_page_number", int)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"requests.params.messages.content.citations.end_block_index",
		"",
		func(string string) error {
			int, err := parseInt(string)
			if err != nil {
				return err
			}
			var jsonErr error
			body, jsonErr = jsonSet(body, "requests.#.params.messages.#.content.#.citations.#.end_block_index", int)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"requests.params.messages.content.citations.start_block_index",
		"",
		func(string string) error {
			int, err := parseInt(string)
			if err != nil {
				return err
			}
			var jsonErr error
			body, jsonErr = jsonSet(body, "requests.#.params.messages.#.content.#.citations.#.start_block_index", int)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.BoolFunc(
		"requests.params.messages.content.+citation",
		"",
		func(_ string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "requests.#.params.messages.#.content.#.citations.-1", map[string]interface{}{})
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"requests.params.messages.content.source.data",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "requests.#.params.messages.#.content.#.source.data", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"requests.params.messages.content.source.media_type",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "requests.#.params.messages.#.content.#.source.media_type", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"requests.params.messages.content.source.type",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "requests.#.params.messages.#.content.#.source.type", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"requests.params.messages.content.id",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "requests.#.params.messages.#.content.#.id", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"requests.params.messages.content.name",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "requests.#.params.messages.#.content.#.name", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"requests.params.messages.content.tool_use_id",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "requests.#.params.messages.#.content.#.tool_use_id", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"requests.params.messages.content.content.text",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "requests.#.params.messages.#.content.#.content.#.text", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"requests.params.messages.content.content.type",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "requests.#.params.messages.#.content.#.content.#.type", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"requests.params.messages.content.content.cache_control.type",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "requests.#.params.messages.#.content.#.content.#.cache_control.type", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"requests.params.messages.content.content.citations.cited_text",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "requests.#.params.messages.#.content.#.content.#.citations.#.cited_text", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"requests.params.messages.content.content.citations.document_index",
		"",
		func(string string) error {
			int, err := parseInt(string)
			if err != nil {
				return err
			}
			var jsonErr error
			body, jsonErr = jsonSet(body, "requests.#.params.messages.#.content.#.content.#.citations.#.document_index", int)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"requests.params.messages.content.content.citations.document_title",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "requests.#.params.messages.#.content.#.content.#.citations.#.document_title", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"requests.params.messages.content.content.citations.end_char_index",
		"",
		func(string string) error {
			int, err := parseInt(string)
			if err != nil {
				return err
			}
			var jsonErr error
			body, jsonErr = jsonSet(body, "requests.#.params.messages.#.content.#.content.#.citations.#.end_char_index", int)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"requests.params.messages.content.content.citations.start_char_index",
		"",
		func(string string) error {
			int, err := parseInt(string)
			if err != nil {
				return err
			}
			var jsonErr error
			body, jsonErr = jsonSet(body, "requests.#.params.messages.#.content.#.content.#.citations.#.start_char_index", int)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"requests.params.messages.content.content.citations.type",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "requests.#.params.messages.#.content.#.content.#.citations.#.type", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"requests.params.messages.content.content.citations.end_page_number",
		"",
		func(string string) error {
			int, err := parseInt(string)
			if err != nil {
				return err
			}
			var jsonErr error
			body, jsonErr = jsonSet(body, "requests.#.params.messages.#.content.#.content.#.citations.#.end_page_number", int)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"requests.params.messages.content.content.citations.start_page_number",
		"",
		func(string string) error {
			int, err := parseInt(string)
			if err != nil {
				return err
			}
			var jsonErr error
			body, jsonErr = jsonSet(body, "requests.#.params.messages.#.content.#.content.#.citations.#.start_page_number", int)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"requests.params.messages.content.content.citations.end_block_index",
		"",
		func(string string) error {
			int, err := parseInt(string)
			if err != nil {
				return err
			}
			var jsonErr error
			body, jsonErr = jsonSet(body, "requests.#.params.messages.#.content.#.content.#.citations.#.end_block_index", int)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"requests.params.messages.content.content.citations.start_block_index",
		"",
		func(string string) error {
			int, err := parseInt(string)
			if err != nil {
				return err
			}
			var jsonErr error
			body, jsonErr = jsonSet(body, "requests.#.params.messages.#.content.#.content.#.citations.#.start_block_index", int)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.BoolFunc(
		"requests.params.messages.content.content.+citation",
		"",
		func(_ string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "requests.#.params.messages.#.content.#.content.#.citations.-1", map[string]interface{}{})
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"requests.params.messages.content.content.source.data",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "requests.#.params.messages.#.content.#.content.#.source.data", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"requests.params.messages.content.content.source.media_type",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "requests.#.params.messages.#.content.#.content.#.source.media_type", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"requests.params.messages.content.content.source.type",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "requests.#.params.messages.#.content.#.content.#.source.type", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.BoolFunc(
		"requests.params.messages.content.+content",
		"",
		func(_ string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "requests.#.params.messages.#.content.#.content.-1", map[string]interface{}{})
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.BoolFunc(
		"requests.params.messages.content.is_error",
		"",
		func(_ string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "requests.#.params.messages.#.content.#.is_error", true)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"requests.params.messages.content.source.content",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "requests.#.params.messages.#.content.#.source.content", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"requests.params.messages.content.source.content.text",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "requests.#.params.messages.#.content.#.source.content.#.text", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"requests.params.messages.content.source.content.type",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "requests.#.params.messages.#.content.#.source.content.#.type", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"requests.params.messages.content.source.content.cache_control.type",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "requests.#.params.messages.#.content.#.source.content.#.cache_control.type", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"requests.params.messages.content.source.content.citations.cited_text",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "requests.#.params.messages.#.content.#.source.content.#.citations.#.cited_text", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"requests.params.messages.content.source.content.citations.document_index",
		"",
		func(string string) error {
			int, err := parseInt(string)
			if err != nil {
				return err
			}
			var jsonErr error
			body, jsonErr = jsonSet(body, "requests.#.params.messages.#.content.#.source.content.#.citations.#.document_index", int)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"requests.params.messages.content.source.content.citations.document_title",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "requests.#.params.messages.#.content.#.source.content.#.citations.#.document_title", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"requests.params.messages.content.source.content.citations.end_char_index",
		"",
		func(string string) error {
			int, err := parseInt(string)
			if err != nil {
				return err
			}
			var jsonErr error
			body, jsonErr = jsonSet(body, "requests.#.params.messages.#.content.#.source.content.#.citations.#.end_char_index", int)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"requests.params.messages.content.source.content.citations.start_char_index",
		"",
		func(string string) error {
			int, err := parseInt(string)
			if err != nil {
				return err
			}
			var jsonErr error
			body, jsonErr = jsonSet(body, "requests.#.params.messages.#.content.#.source.content.#.citations.#.start_char_index", int)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"requests.params.messages.content.source.content.citations.type",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "requests.#.params.messages.#.content.#.source.content.#.citations.#.type", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"requests.params.messages.content.source.content.citations.end_page_number",
		"",
		func(string string) error {
			int, err := parseInt(string)
			if err != nil {
				return err
			}
			var jsonErr error
			body, jsonErr = jsonSet(body, "requests.#.params.messages.#.content.#.source.content.#.citations.#.end_page_number", int)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"requests.params.messages.content.source.content.citations.start_page_number",
		"",
		func(string string) error {
			int, err := parseInt(string)
			if err != nil {
				return err
			}
			var jsonErr error
			body, jsonErr = jsonSet(body, "requests.#.params.messages.#.content.#.source.content.#.citations.#.start_page_number", int)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"requests.params.messages.content.source.content.citations.end_block_index",
		"",
		func(string string) error {
			int, err := parseInt(string)
			if err != nil {
				return err
			}
			var jsonErr error
			body, jsonErr = jsonSet(body, "requests.#.params.messages.#.content.#.source.content.#.citations.#.end_block_index", int)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"requests.params.messages.content.source.content.citations.start_block_index",
		"",
		func(string string) error {
			int, err := parseInt(string)
			if err != nil {
				return err
			}
			var jsonErr error
			body, jsonErr = jsonSet(body, "requests.#.params.messages.#.content.#.source.content.#.citations.#.start_block_index", int)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.BoolFunc(
		"requests.params.messages.content.source.content.+citation",
		"",
		func(_ string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "requests.#.params.messages.#.content.#.source.content.#.citations.-1", map[string]interface{}{})
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"requests.params.messages.content.source.content.source.data",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "requests.#.params.messages.#.content.#.source.content.#.source.data", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"requests.params.messages.content.source.content.source.media_type",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "requests.#.params.messages.#.content.#.source.content.#.source.media_type", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"requests.params.messages.content.source.content.source.type",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "requests.#.params.messages.#.content.#.source.content.#.source.type", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.BoolFunc(
		"requests.params.messages.content.source.+content",
		"",
		func(_ string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "requests.#.params.messages.#.content.#.source.content.-1", map[string]interface{}{})
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.BoolFunc(
		"requests.params.messages.content.citations.enabled",
		"",
		func(_ string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "requests.#.params.messages.#.content.#.citations.enabled", true)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"requests.params.messages.content.context",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "requests.#.params.messages.#.content.#.context", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"requests.params.messages.content.title",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "requests.#.params.messages.#.content.#.title", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"requests.params.messages.content.signature",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "requests.#.params.messages.#.content.#.signature", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"requests.params.messages.content.thinking",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "requests.#.params.messages.#.content.#.thinking", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"requests.params.messages.content.data",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "requests.#.params.messages.#.content.#.data", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.BoolFunc(
		"requests.params.messages.+content",
		"",
		func(_ string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "requests.#.params.messages.#.content.-1", map[string]interface{}{})
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"requests.params.messages.role",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "requests.#.params.messages.#.role", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.BoolFunc(
		"requests.params.+message",
		"",
		func(_ string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "requests.#.params.messages.-1", map[string]interface{}{})
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"requests.params.model",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "requests.#.params.model", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"requests.params.metadata.user_id",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "requests.#.params.metadata.user_id", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"requests.params.stop_sequences",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "requests.#.params.stop_sequences.#", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"requests.params.+stop_sequence",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "requests.#.params.stop_sequences.-1", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.BoolFunc(
		"requests.params.stream",
		"",
		func(_ string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "requests.#.params.stream", true)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"requests.params.system.text",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "requests.#.params.system.#.text", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"requests.params.system.type",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "requests.#.params.system.#.type", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"requests.params.system.cache_control.type",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "requests.#.params.system.#.cache_control.type", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"requests.params.system.citations.cited_text",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "requests.#.params.system.#.citations.#.cited_text", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"requests.params.system.citations.document_index",
		"",
		func(string string) error {
			int, err := parseInt(string)
			if err != nil {
				return err
			}
			var jsonErr error
			body, jsonErr = jsonSet(body, "requests.#.params.system.#.citations.#.document_index", int)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"requests.params.system.citations.document_title",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "requests.#.params.system.#.citations.#.document_title", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"requests.params.system.citations.end_char_index",
		"",
		func(string string) error {
			int, err := parseInt(string)
			if err != nil {
				return err
			}
			var jsonErr error
			body, jsonErr = jsonSet(body, "requests.#.params.system.#.citations.#.end_char_index", int)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"requests.params.system.citations.start_char_index",
		"",
		func(string string) error {
			int, err := parseInt(string)
			if err != nil {
				return err
			}
			var jsonErr error
			body, jsonErr = jsonSet(body, "requests.#.params.system.#.citations.#.start_char_index", int)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"requests.params.system.citations.type",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "requests.#.params.system.#.citations.#.type", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"requests.params.system.citations.end_page_number",
		"",
		func(string string) error {
			int, err := parseInt(string)
			if err != nil {
				return err
			}
			var jsonErr error
			body, jsonErr = jsonSet(body, "requests.#.params.system.#.citations.#.end_page_number", int)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"requests.params.system.citations.start_page_number",
		"",
		func(string string) error {
			int, err := parseInt(string)
			if err != nil {
				return err
			}
			var jsonErr error
			body, jsonErr = jsonSet(body, "requests.#.params.system.#.citations.#.start_page_number", int)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"requests.params.system.citations.end_block_index",
		"",
		func(string string) error {
			int, err := parseInt(string)
			if err != nil {
				return err
			}
			var jsonErr error
			body, jsonErr = jsonSet(body, "requests.#.params.system.#.citations.#.end_block_index", int)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"requests.params.system.citations.start_block_index",
		"",
		func(string string) error {
			int, err := parseInt(string)
			if err != nil {
				return err
			}
			var jsonErr error
			body, jsonErr = jsonSet(body, "requests.#.params.system.#.citations.#.start_block_index", int)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.BoolFunc(
		"requests.params.system.+citation",
		"",
		func(_ string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "requests.#.params.system.#.citations.-1", map[string]interface{}{})
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.BoolFunc(
		"requests.params.+system",
		"",
		func(_ string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "requests.#.params.system.-1", map[string]interface{}{})
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"requests.params.temperature",
		"",
		func(string string) error {
			float, err := parseFloat(string)
			if err != nil {
				return err
			}
			var jsonErr error
			body, jsonErr = jsonSet(body, "requests.#.params.temperature", float)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"requests.params.thinking.budget_tokens",
		"",
		func(string string) error {
			int, err := parseInt(string)
			if err != nil {
				return err
			}
			var jsonErr error
			body, jsonErr = jsonSet(body, "requests.#.params.thinking.budget_tokens", int)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"requests.params.thinking.type",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "requests.#.params.thinking.type", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"requests.params.tool_choice.type",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "requests.#.params.tool_choice.type", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.BoolFunc(
		"requests.params.tool_choice.disable_parallel_tool_use",
		"",
		func(_ string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "requests.#.params.tool_choice.disable_parallel_tool_use", true)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"requests.params.tool_choice.name",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "requests.#.params.tool_choice.name", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"requests.params.tools.display_height_px",
		"",
		func(string string) error {
			int, err := parseInt(string)
			if err != nil {
				return err
			}
			var jsonErr error
			body, jsonErr = jsonSet(body, "requests.#.params.tools.#.display_height_px", int)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"requests.params.tools.display_width_px",
		"",
		func(string string) error {
			int, err := parseInt(string)
			if err != nil {
				return err
			}
			var jsonErr error
			body, jsonErr = jsonSet(body, "requests.#.params.tools.#.display_width_px", int)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"requests.params.tools.name",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "requests.#.params.tools.#.name", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"requests.params.tools.type",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "requests.#.params.tools.#.type", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"requests.params.tools.cache_control.type",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "requests.#.params.tools.#.cache_control.type", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"requests.params.tools.display_number",
		"",
		func(string string) error {
			int, err := parseInt(string)
			if err != nil {
				return err
			}
			var jsonErr error
			body, jsonErr = jsonSet(body, "requests.#.params.tools.#.display_number", int)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"requests.params.tools.input_schema.type",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "requests.#.params.tools.#.input_schema.type", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"requests.params.tools.description",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "requests.#.params.tools.#.description", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.BoolFunc(
		"requests.params.+tool",
		"",
		func(_ string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "requests.#.params.tools.-1", map[string]interface{}{})
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"requests.params.top_k",
		"",
		func(string string) error {
			int, err := parseInt(string)
			if err != nil {
				return err
			}
			var jsonErr error
			body, jsonErr = jsonSet(body, "requests.#.params.top_k", int)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"requests.params.top_p",
		"",
		func(string string) error {
			float, err := parseFloat(string)
			if err != nil {
				return err
			}
			var jsonErr error
			body, jsonErr = jsonSet(body, "requests.#.params.top_p", float)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.BoolFunc(
		"+request",
		"",
		func(_ string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "requests.-1", map[string]interface{}{})
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"betas",
		"",
		func(string string) error {
			var jsonErr error
			header, jsonErr = jsonSet(header, "anthropic-beta.#", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"+beta",
		"",
		func(string string) error {
			var jsonErr error
			header, jsonErr = jsonSet(header, "anthropic-beta.-1", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	return Subcommand{
		flagSet: flagSet,
		handle: func(client *anthropic.Client) {
			res, err := client.Beta.Messages.Batches.New(
				context.TODO(),
				anthropic.BetaMessageBatchNewParams{},
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

func createBetaMessagesBatchesRetrieveSubcommand() Subcommand {
	var messageBatchID *string = nil
	query := []byte("{}")
	header := []byte("{}")
	var flagSet = flag.NewFlagSet("beta.messages.batches.retrieve", flag.ExitOnError)

	flagSet.Func(
		"message-batch-id",
		"",
		func(string string) error {
			messageBatchID = &string
			return nil
		},
	)

	flagSet.Func(
		"betas",
		"",
		func(string string) error {
			var jsonErr error
			header, jsonErr = jsonSet(header, "anthropic-beta.#", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"+beta",
		"",
		func(string string) error {
			var jsonErr error
			header, jsonErr = jsonSet(header, "anthropic-beta.-1", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	return Subcommand{
		flagSet: flagSet,
		handle: func(client *anthropic.Client) {
			res, err := client.Beta.Messages.Batches.Get(
				context.TODO(),
				*messageBatchID,
				anthropic.BetaMessageBatchGetParams{},
				option.WithMiddleware(func(r *http.Request, mn option.MiddlewareNext) (*http.Response, error) {
					r.URL.RawQuery = serializeQuery(query).Encode()
					r.Header = serializeHeader(header)
					return mn(r)
				}),
			)
			if err != nil {
				fmt.Printf("%s\n", err)
				os.Exit(1)
			}

			fmt.Printf("%s\n", res.JSON.RawJSON())
		},
	}
}

func createBetaMessagesBatchesListSubcommand() Subcommand {
	query := []byte("{}")
	header := []byte("{}")
	var flagSet = flag.NewFlagSet("beta.messages.batches.list", flag.ExitOnError)

	flagSet.Func(
		"after-id",
		"",
		func(string string) error {
			var jsonErr error
			query, jsonErr = jsonSet(query, "after_id", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"before-id",
		"",
		func(string string) error {
			var jsonErr error
			query, jsonErr = jsonSet(query, "before_id", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"limit",
		"",
		func(string string) error {
			int, err := parseInt(string)
			if err != nil {
				return err
			}
			var jsonErr error
			query, jsonErr = jsonSet(query, "limit", int)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"betas",
		"",
		func(string string) error {
			var jsonErr error
			header, jsonErr = jsonSet(header, "anthropic-beta.#", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"+beta",
		"",
		func(string string) error {
			var jsonErr error
			header, jsonErr = jsonSet(header, "anthropic-beta.-1", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	return Subcommand{
		flagSet: flagSet,
		handle: func(client *anthropic.Client) {
			res, err := client.Beta.Messages.Batches.List(
				context.TODO(),
				anthropic.BetaMessageBatchListParams{},
				option.WithMiddleware(func(r *http.Request, mn option.MiddlewareNext) (*http.Response, error) {
					r.URL.RawQuery = serializeQuery(query).Encode()
					r.Header = serializeHeader(header)
					return mn(r)
				}),
			)
			if err != nil {
				fmt.Printf("%s\n", err)
				os.Exit(1)
			}

			fmt.Printf("%s\n", res.JSON.RawJSON())
		},
	}
}

func createBetaMessagesBatchesDeleteSubcommand(initialBody []byte) Subcommand {
	var messageBatchID *string = nil
	query := []byte("{}")
	header := []byte("{}")
	body := initialBody
	var flagSet = flag.NewFlagSet("beta.messages.batches.delete", flag.ExitOnError)

	flagSet.Func(
		"message-batch-id",
		"",
		func(string string) error {
			messageBatchID = &string
			return nil
		},
	)

	flagSet.Func(
		"betas",
		"",
		func(string string) error {
			var jsonErr error
			header, jsonErr = jsonSet(header, "anthropic-beta.#", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"+beta",
		"",
		func(string string) error {
			var jsonErr error
			header, jsonErr = jsonSet(header, "anthropic-beta.-1", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	return Subcommand{
		flagSet: flagSet,
		handle: func(client *anthropic.Client) {
			res, err := client.Beta.Messages.Batches.Delete(
				context.TODO(),
				*messageBatchID,
				anthropic.BetaMessageBatchDeleteParams{},
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

func createBetaMessagesBatchesCancelSubcommand(initialBody []byte) Subcommand {
	var messageBatchID *string = nil
	query := []byte("{}")
	header := []byte("{}")
	body := initialBody
	var flagSet = flag.NewFlagSet("beta.messages.batches.cancel", flag.ExitOnError)

	flagSet.Func(
		"message-batch-id",
		"",
		func(string string) error {
			messageBatchID = &string
			return nil
		},
	)

	flagSet.Func(
		"betas",
		"",
		func(string string) error {
			var jsonErr error
			header, jsonErr = jsonSet(header, "anthropic-beta.#", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"+beta",
		"",
		func(string string) error {
			var jsonErr error
			header, jsonErr = jsonSet(header, "anthropic-beta.-1", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	return Subcommand{
		flagSet: flagSet,
		handle: func(client *anthropic.Client) {
			res, err := client.Beta.Messages.Batches.Cancel(
				context.TODO(),
				*messageBatchID,
				anthropic.BetaMessageBatchCancelParams{},
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
