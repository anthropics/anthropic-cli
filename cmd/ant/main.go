// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"slices"

	"github.com/anthropics/anthropic-cli/pkg/cmd"
	"github.com/anthropics/anthropic-sdk-go"
	"github.com/tidwall/gjson"
	"github.com/urfave/cli/v3"
)

func main() {
	app := cmd.Command

	if slices.Contains(os.Args, "__complete") {
		prepareForAutocomplete(app)
	}

	if err := app.Run(context.Background(), os.Args); err != nil {
		exitCode := 1

		// Check if error has a custom exit code
		if exitErr, ok := err.(cli.ExitCoder); ok {
			exitCode = exitErr.ExitCode()
		}

		var apierr *anthropic.Error
		if errors.As(err, &apierr) {
			fmt.Fprintf(os.Stderr, "%s %q: %d %s\n", apierr.Request.Method, apierr.Request.URL, apierr.Response.StatusCode, http.StatusText(apierr.Response.StatusCode))
			format := app.String("format-error")
			json := gjson.Parse(apierr.RawJSON())
			show_err := cmd.ShowJSON(os.Stdout, "Error", json, format, app.String("transform-error"))
			if show_err != nil {
				// Just print the original error:
				fmt.Fprintf(os.Stderr, "%s\n", err.Error())
			}
		} else {
			fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		}
		os.Exit(exitCode)
	}
}

func prepareForAutocomplete(cmd *cli.Command) {
	// urfave/cli does not handle flag completions and will print an error if we inspect a command with invalid flags.
	// This skips that sort of validation
	cmd.SkipFlagParsing = true
	for _, child := range cmd.Commands {
		prepareForAutocomplete(child)
	}
}
