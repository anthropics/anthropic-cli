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

var betaFilesList = cli.Command{
	Name:  "list",
	Usage: "List Files",
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
	Action:          handleBetaFilesList,
	HideHelpCommand: true,
}

var betaFilesDelete = cli.Command{
	Name:  "delete",
	Usage: "Delete File",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name: "file-id",
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
	Action:          handleBetaFilesDelete,
	HideHelpCommand: true,
}

var betaFilesDownload = cli.Command{
	Name:  "download",
	Usage: "Download File",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name: "file-id",
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
	Action:          handleBetaFilesDownload,
	HideHelpCommand: true,
}

var betaFilesRetrieveMetadata = cli.Command{
	Name:  "retrieve_metadata",
	Usage: "Get File Metadata",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name: "file-id",
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
	Action:          handleBetaFilesRetrieveMetadata,
	HideHelpCommand: true,
}

var betaFilesUpload = cli.Command{
	Name:  "upload",
	Usage: "Upload File",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:   "file",
			Action: getAPIFlagAction[string]("body", "file"),
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
	Action:          handleBetaFilesUpload,
	HideHelpCommand: true,
}

func handleBetaFilesList(ctx context.Context, cmd *cli.Command) error {
	cc := getAPICommandContext(ctx, cmd)
	params := anthropic.BetaFileListParams{}
	res, err := cc.client.Beta.Files.List(
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

func handleBetaFilesDelete(ctx context.Context, cmd *cli.Command) error {
	cc := getAPICommandContext(ctx, cmd)
	params := anthropic.BetaFileDeleteParams{}
	res, err := cc.client.Beta.Files.Delete(
		context.TODO(),
		cmd.Value("file-id").(string),
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

func handleBetaFilesDownload(ctx context.Context, cmd *cli.Command) error {
	cc := getAPICommandContext(ctx, cmd)
	params := anthropic.BetaFileDownloadParams{}
	res, err := cc.client.Beta.Files.Download(
		context.TODO(),
		cmd.Value("file-id").(string),
		params,
		option.WithMiddleware(cc.AsMiddleware()),
	)
	if err != nil {
		return err
	}

	fmt.Printf("%s\n", colorizeJSON(res.RawJSON(), os.Stdout))
	return nil
}

func handleBetaFilesRetrieveMetadata(ctx context.Context, cmd *cli.Command) error {
	cc := getAPICommandContext(ctx, cmd)
	params := anthropic.BetaFileGetMetadataParams{}
	res, err := cc.client.Beta.Files.GetMetadata(
		context.TODO(),
		cmd.Value("file-id").(string),
		params,
		option.WithMiddleware(cc.AsMiddleware()),
	)
	if err != nil {
		return err
	}

	fmt.Printf("%s\n", colorizeJSON(res.RawJSON(), os.Stdout))
	return nil
}

func handleBetaFilesUpload(ctx context.Context, cmd *cli.Command) error {
	cc := getAPICommandContext(ctx, cmd)
	params := anthropic.BetaFileUploadParams{}
	res, err := cc.client.Beta.Files.Upload(
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
