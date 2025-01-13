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

func createBetaModelsRetrieveSubcommand() Subcommand {
	var modelID *string = nil
	query := []byte("{}")
	header := []byte("{}")
	var flagSet = flag.NewFlagSet("beta.models.retrieve", flag.ExitOnError)

	flagSet.Func(
		"model-id",
		"",
		func(string string) error {
			modelID = &string
			return nil
		},
	)

	return Subcommand{
		flagSet: flagSet,
		handle: func(client *anthropic.Client) {
			res, err := client.Beta.Models.Get(
				context.TODO(),
				*modelID,
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

func createBetaModelsListSubcommand() Subcommand {
	query := []byte("{}")
	header := []byte("{}")
	var flagSet = flag.NewFlagSet("beta.models.list", flag.ExitOnError)

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

	return Subcommand{
		flagSet: flagSet,
		handle: func(client *anthropic.Client) {
			res, err := client.Beta.Models.List(
				context.TODO(),
				anthropic.BetaModelListParams{},
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
