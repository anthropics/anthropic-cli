// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package main

import (
	"context"
	"log"
	"os"

	"github.com/stainless-sdks/anthropic-cli/pkg/cmd"
)

func main() {
	app := cmd.Command
	if err := app.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
