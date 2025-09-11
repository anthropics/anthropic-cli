// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cmd

import (
	"compress/gzip"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	docs "github.com/urfave/cli-docs/v3"
	"github.com/urfave/cli/v3"
)

var (
	Command       *cli.Command
	OutputFormats = []string{"auto", "explore", "json", "pretty", "raw", "yaml"}
)

func init() {
	Command = &cli.Command{
		Name:    "anthropic-cli",
		Usage:   "CLI for the anthropic API",
		Version: Version,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "debug",
				Usage: "Enable debug logging",
			},
			&cli.StringFlag{
				Name:        "base-url",
				DefaultText: "url",
				Usage:       "Override the base URL for API requests",
			},
			&cli.StringFlag{
				Name:  "format",
				Usage: "The format for data output (one of: " + strings.Join(OutputFormats[:], ", ") + ")",
				Value: "auto",
				Validator: func(format string) error {
					if !slices.Contains(OutputFormats[:], strings.ToLower(format)) {
						return fmt.Errorf("format must be one of: %s", strings.Join(OutputFormats[:], ", "))
					}
					return nil
				},
			},
		},
		Commands: []*cli.Command{
			{
				Name:     "completions",
				Category: "API RESOURCE",
				Commands: []*cli.Command{
					&completionsCreate,
				},
			},
			{
				Name:     "messages",
				Category: "API RESOURCE",
				Commands: []*cli.Command{
					&messagesCreate,
					&messagesCountTokens,
				},
			},
			{
				Name:     "messages:batches",
				Category: "API RESOURCE",
				Commands: []*cli.Command{
					&messagesBatchesCreate,
					&messagesBatchesRetrieve,
					&messagesBatchesList,
					&messagesBatchesDelete,
					&messagesBatchesCancel,
				},
			},
			{
				Name:     "models",
				Category: "API RESOURCE",
				Commands: []*cli.Command{
					&modelsRetrieve,
					&modelsList,
				},
			},
			{
				Name:     "beta:models",
				Category: "API RESOURCE",
				Commands: []*cli.Command{
					&betaModelsRetrieve,
					&betaModelsList,
				},
			},
			{
				Name:     "beta:messages",
				Category: "API RESOURCE",
				Commands: []*cli.Command{
					&betaMessagesCreate,
					&betaMessagesCountTokens,
				},
			},
			{
				Name:     "beta:messages:batches",
				Category: "API RESOURCE",
				Commands: []*cli.Command{
					&betaMessagesBatchesCreate,
					&betaMessagesBatchesRetrieve,
					&betaMessagesBatchesList,
					&betaMessagesBatchesDelete,
					&betaMessagesBatchesCancel,
				},
			},
			{
				Name:     "beta:files",
				Category: "API RESOURCE",
				Commands: []*cli.Command{
					&betaFilesList,
					&betaFilesDelete,
					&betaFilesRetrieveMetadata,
					&betaFilesUpload,
				},
			},
			{
				Name:            "@manpages",
				Usage:           "Generate documentation for 'man'",
				UsageText:       "anthropic-cli @manpages [-o anthropic-cli.1] [--gzip]",
				Hidden:          true,
				Action:          generateManpages,
				HideHelpCommand: true,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "output",
						Aliases: []string{"o"},
						Usage:   "write manpages to the given folder",
						Value:   "man",
					},
					&cli.BoolFlag{
						Name:    "gzip",
						Aliases: []string{"z"},
						Usage:   "output gzipped manpage files to .gz",
						Value:   true,
					},
					&cli.BoolFlag{
						Name:    "text",
						Aliases: []string{"z"},
						Usage:   "output uncompressed text files",
						Value:   false,
					},
				},
			},
		},
		EnableShellCompletion:      true,
		ShellCompletionCommandName: "@completion",
		HideHelpCommand:            true,
	}
}

func generateManpages(ctx context.Context, c *cli.Command) error {
	manpage, err := docs.ToManWithSection(Command, 1)
	if err != nil {
		return err
	}
	dir := c.String("output")
	err = os.MkdirAll(filepath.Join(dir, "man1"), 0755)
	if err != nil {
		// handle error
	}
	if c.Bool("text") {
		file, err := os.Create(filepath.Join(dir, "man1", "anthropic-cli.1"))
		if err != nil {
			return err
		}
		defer file.Close()
		if _, err := file.WriteString(manpage); err != nil {
			return err
		}
	}
	if c.Bool("gzip") {
		file, err := os.Create(filepath.Join(dir, "man1", "anthropic-cli.1.gz"))
		if err != nil {
			return err
		}
		defer file.Close()
		gzWriter := gzip.NewWriter(file)
		defer gzWriter.Close()
		_, err = gzWriter.Write([]byte(manpage))
		if err != nil {
			return err
		}
	}
	fmt.Printf("Wrote manpages to %s\n", dir)
	return nil
}
