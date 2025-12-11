// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/stainless-sdks/anthropic-cli/pkg/cmd"
	"github.com/tidwall/gjson"
)

func main() {
	app := cmd.Command
	if err := app.Run(context.Background(), os.Args); err != nil {
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
		os.Exit(1)
	}
}
