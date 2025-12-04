// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cmd

import (
	"context"
	"fmt"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/stainless-sdks/anthropic-cli/internal/apiquery"
	"github.com/stainless-sdks/anthropic-cli/internal/requestflag"
	"github.com/tidwall/gjson"
	"github.com/urfave/cli/v3"
)

var betaFilesList = cli.Command{
	Name:  "list",
	Usage: "List Files",
	Flags: []cli.Flag{
		&requestflag.StringFlag{
			Name:  "after-id",
			Usage: "ID of the object to use as a cursor for pagination. When provided, returns the page of results immediately after this object.",
			Config: requestflag.RequestConfig{
				QueryPath: "after_id",
			},
		},
		&requestflag.StringFlag{
			Name:  "before-id",
			Usage: "ID of the object to use as a cursor for pagination. When provided, returns the page of results immediately before this object.",
			Config: requestflag.RequestConfig{
				QueryPath: "before_id",
			},
		},
		&requestflag.IntFlag{
			Name:  "limit",
			Usage: "Number of items to return per page.\n\nDefaults to `20`. Ranges from `1` to `1000`.",
			Value: requestflag.Value[int64](20),
			Config: requestflag.RequestConfig{
				QueryPath: "limit",
			},
		},
		&requestflag.YAMLSliceFlag{
			Name:  "beta",
			Usage: "Optional header to specify the beta version(s) you want to use.",
			Config: requestflag.RequestConfig{
				HeaderPath: "anthropic-beta",
			},
		},
	},
	Action:          handleBetaFilesList,
	HideHelpCommand: true,
}

var betaFilesDelete = cli.Command{
	Name:  "delete",
	Usage: "Delete File",
	Flags: []cli.Flag{
		&requestflag.StringFlag{
			Name:  "file-id",
			Usage: "ID of the File.",
		},
		&requestflag.YAMLSliceFlag{
			Name:  "beta",
			Usage: "Optional header to specify the beta version(s) you want to use.",
			Config: requestflag.RequestConfig{
				HeaderPath: "anthropic-beta",
			},
		},
	},
	Action:          handleBetaFilesDelete,
	HideHelpCommand: true,
}

var betaFilesDownload = cli.Command{
	Name:  "download",
	Usage: "Download File",
	Flags: []cli.Flag{
		&requestflag.StringFlag{
			Name:  "file-id",
			Usage: "ID of the File.",
		},
		&requestflag.YAMLSliceFlag{
			Name:  "beta",
			Usage: "Optional header to specify the beta version(s) you want to use.",
			Config: requestflag.RequestConfig{
				HeaderPath: "anthropic-beta",
			},
		},
	},
	Action:          handleBetaFilesDownload,
	HideHelpCommand: true,
}

var betaFilesRetrieveMetadata = cli.Command{
	Name:  "retrieve-metadata",
	Usage: "Get File Metadata",
	Flags: []cli.Flag{
		&requestflag.StringFlag{
			Name:  "file-id",
			Usage: "ID of the File.",
		},
		&requestflag.YAMLSliceFlag{
			Name:  "beta",
			Usage: "Optional header to specify the beta version(s) you want to use.",
			Config: requestflag.RequestConfig{
				HeaderPath: "anthropic-beta",
			},
		},
	},
	Action:          handleBetaFilesRetrieveMetadata,
	HideHelpCommand: true,
}

var betaFilesUpload = cli.Command{
	Name:  "upload",
	Usage: "Upload File",
	Flags: []cli.Flag{
		&requestflag.StringFlag{
			Name:  "file",
			Usage: "The file to upload",
			Config: requestflag.RequestConfig{
				BodyPath: "file",
			},
		},
		&requestflag.YAMLSliceFlag{
			Name:  "beta",
			Usage: "Optional header to specify the beta version(s) you want to use.",
			Config: requestflag.RequestConfig{
				HeaderPath: "anthropic-beta",
			},
		},
	},
	Action:          handleBetaFilesUpload,
	HideHelpCommand: true,
}

func handleBetaFilesList(ctx context.Context, cmd *cli.Command) error {
	client := anthropic.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if len(unusedArgs) > 0 {
		return fmt.Errorf("Unexpected extra arguments: %v", unusedArgs)
	}
	params := anthropic.BetaFileListParams{}

	options, err := flagOptions(
		cmd,
		apiquery.NestedQueryFormatBrackets,
		apiquery.ArrayQueryFormatComma,
		ApplicationJSON,
	)
	if err != nil {
		return err
	}
	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Beta.Files.List(
		ctx,
		params,
		options...,
	)
	if err != nil {
		return err
	}

	json := gjson.Parse(string(res))
	format := cmd.Root().String("format")
	transform := cmd.Root().String("transform")
	return ShowJSON("beta:files list", json, format, transform)
}

func handleBetaFilesDelete(ctx context.Context, cmd *cli.Command) error {
	client := anthropic.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("file-id") && len(unusedArgs) > 0 {
		cmd.Set("file-id", unusedArgs[0])
		unusedArgs = unusedArgs[1:]
	}
	if len(unusedArgs) > 0 {
		return fmt.Errorf("Unexpected extra arguments: %v", unusedArgs)
	}
	params := anthropic.BetaFileDeleteParams{}

	options, err := flagOptions(
		cmd,
		apiquery.NestedQueryFormatBrackets,
		apiquery.ArrayQueryFormatComma,
		ApplicationJSON,
	)
	if err != nil {
		return err
	}
	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Beta.Files.Delete(
		ctx,
		requestflag.CommandRequestValue[string](cmd, "file-id"),
		params,
		options...,
	)
	if err != nil {
		return err
	}

	json := gjson.Parse(string(res))
	format := cmd.Root().String("format")
	transform := cmd.Root().String("transform")
	return ShowJSON("beta:files delete", json, format, transform)
}

func handleBetaFilesDownload(ctx context.Context, cmd *cli.Command) error {
	client := anthropic.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("file-id") && len(unusedArgs) > 0 {
		cmd.Set("file-id", unusedArgs[0])
		unusedArgs = unusedArgs[1:]
	}
	if len(unusedArgs) > 0 {
		return fmt.Errorf("Unexpected extra arguments: %v", unusedArgs)
	}
	params := anthropic.BetaFileDownloadParams{}

	options, err := flagOptions(
		cmd,
		apiquery.NestedQueryFormatBrackets,
		apiquery.ArrayQueryFormatComma,
		ApplicationJSON,
	)
	if err != nil {
		return err
	}
	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Beta.Files.Download(
		ctx,
		requestflag.CommandRequestValue[string](cmd, "file-id"),
		params,
		options...,
	)
	if err != nil {
		return err
	}

	json := gjson.Parse(string(res))
	format := cmd.Root().String("format")
	transform := cmd.Root().String("transform")
	return ShowJSON("beta:files download", json, format, transform)
}

func handleBetaFilesRetrieveMetadata(ctx context.Context, cmd *cli.Command) error {
	client := anthropic.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("file-id") && len(unusedArgs) > 0 {
		cmd.Set("file-id", unusedArgs[0])
		unusedArgs = unusedArgs[1:]
	}
	if len(unusedArgs) > 0 {
		return fmt.Errorf("Unexpected extra arguments: %v", unusedArgs)
	}
	params := anthropic.BetaFileGetMetadataParams{}

	options, err := flagOptions(
		cmd,
		apiquery.NestedQueryFormatBrackets,
		apiquery.ArrayQueryFormatComma,
		ApplicationJSON,
	)
	if err != nil {
		return err
	}
	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Beta.Files.GetMetadata(
		ctx,
		requestflag.CommandRequestValue[string](cmd, "file-id"),
		params,
		options...,
	)
	if err != nil {
		return err
	}

	json := gjson.Parse(string(res))
	format := cmd.Root().String("format")
	transform := cmd.Root().String("transform")
	return ShowJSON("beta:files retrieve-metadata", json, format, transform)
}

func handleBetaFilesUpload(ctx context.Context, cmd *cli.Command) error {
	client := anthropic.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if len(unusedArgs) > 0 {
		return fmt.Errorf("Unexpected extra arguments: %v", unusedArgs)
	}
	params := anthropic.BetaFileUploadParams{}

	options, err := flagOptions(
		cmd,
		apiquery.NestedQueryFormatBrackets,
		apiquery.ArrayQueryFormatComma,
		MultipartFormEncoded,
	)
	if err != nil {
		return err
	}
	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Beta.Files.Upload(
		ctx,
		params,
		options...,
	)
	if err != nil {
		return err
	}

	json := gjson.Parse(string(res))
	format := cmd.Root().String("format")
	transform := cmd.Root().String("transform")
	return ShowJSON("beta:files upload", json, format, transform)
}
