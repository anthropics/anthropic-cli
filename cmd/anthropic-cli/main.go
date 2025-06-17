// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/stainless-sdks/anthropic-cli/pkg/cmd"
)

func main() {
	app := cmd.Command
	if err := app.Run(context.Background(), os.Args); err != nil {
		var apierr *anthropic.Error
		if errors.As(err, &apierr) {
			fmt.Printf("%s\n", cmd.ColorizeJSON(apierr.RawJSON(), os.Stderr))
		} else {
			fmt.Printf("%s\n", err.Error())
		}
		os.Exit(1)
	}
}
