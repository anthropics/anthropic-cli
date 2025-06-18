// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cmd

import (
	"github.com/urfave/cli/v3"
)

var Command = cli.Command{
	Name:    "anthropic-cli",
	Usage:   "CLI for the anthropic API",
	Version: Version,
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:  "debug",
			Usage: "Enable debug logging",
		},
	},
	Commands: []*cli.Command{
		{
			Name: "completions",
			Commands: []*cli.Command{
				&completionsCreate,
			},
		},

		{
			Name: "messages",
			Commands: []*cli.Command{
				&messagesCreate,
				&messagesCountTokens,
			},
		},

		{
			Name: "messages:batches",
			Commands: []*cli.Command{
				&messagesBatchesCreate,
				&messagesBatchesRetrieve,
				&messagesBatchesList,
				&messagesBatchesDelete,
				&messagesBatchesCancel,
			},
		},

		{
			Name: "models",
			Commands: []*cli.Command{
				&modelsRetrieve,
				&modelsList,
			},
		},

		{
			Name: "beta:models",
			Commands: []*cli.Command{
				&betaModelsRetrieve,
				&betaModelsList,
			},
		},

		{
			Name: "beta:messages",
			Commands: []*cli.Command{
				&betaMessagesCreate,
				&betaMessagesCountTokens,
			},
		},

		{
			Name: "beta:messages:batches",
			Commands: []*cli.Command{
				&betaMessagesBatchesCreate,
				&betaMessagesBatchesRetrieve,
				&betaMessagesBatchesList,
				&betaMessagesBatchesDelete,
				&betaMessagesBatchesCancel,
			},
		},

		{
			Name: "beta:files",
			Commands: []*cli.Command{
				&betaFilesList,
				&betaFilesDelete,
				&betaFilesRetrieveMetadata,
				&betaFilesUpload,
			},
		},
	},
	EnableShellCompletion: true,
	HideHelpCommand:       true,
}
