// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cmd

import (
	"context"
	"fmt"

	"github.com/anthropics/anthropic-cli/internal/apiquery"
	"github.com/anthropics/anthropic-cli/internal/requestflag"
	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/tidwall/gjson"
	"github.com/urfave/cli/v3"
)

var betaTunnelsCertificatesCreate = cli.Command{
	Name:    "create",
	Usage:   "The Tunnels API is in research preview. It requires the\n`anthropic-beta: mcp-tunnels-2026-06-22` header and may change without a\ndeprecation period. It supersedes the Admin API endpoints at\n`/v1/organizations/tunnels`, which remain available during a migration window.",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "tunnel-id",
			Required:  true,
			PathParam: "tunnel_id",
		},
		&requestflag.Flag[string]{
			Name:     "ca-certificate-pem",
			Usage:    "PEM-encoded X.509 CA certificate. Must contain exactly one certificate and no private-key material. Maximum 8KB.",
			Required: true,
			BodyPath: "ca_certificate_pem",
		},
		&requestflag.Flag[[]string]{
			Name:       "beta",
			Usage:      "Optional header to specify the beta version(s) you want to use.",
			HeaderPath: "anthropic-beta",
		},
	},
	Action:          handleBetaTunnelsCertificatesCreate,
	HideHelpCommand: true,
}

var betaTunnelsCertificatesRetrieve = cli.Command{
	Name:    "retrieve",
	Usage:   "The Tunnels API is in research preview. It requires the\n`anthropic-beta: mcp-tunnels-2026-06-22` header and may change without a\ndeprecation period. It supersedes the Admin API endpoints at\n`/v1/organizations/tunnels`, which remain available during a migration window.",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "tunnel-id",
			Required:  true,
			PathParam: "tunnel_id",
		},
		&requestflag.Flag[string]{
			Name:      "certificate-id",
			Required:  true,
			PathParam: "certificate_id",
		},
		&requestflag.Flag[[]string]{
			Name:       "beta",
			Usage:      "Optional header to specify the beta version(s) you want to use.",
			HeaderPath: "anthropic-beta",
		},
	},
	Action:          handleBetaTunnelsCertificatesRetrieve,
	HideHelpCommand: true,
}

var betaTunnelsCertificatesList = cli.Command{
	Name:    "list",
	Usage:   "The Tunnels API is in research preview. It requires the\n`anthropic-beta: mcp-tunnels-2026-06-22` header and may change without a\ndeprecation period. It supersedes the Admin API endpoints at\n`/v1/organizations/tunnels`, which remain available during a migration window.",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "tunnel-id",
			Required:  true,
			PathParam: "tunnel_id",
		},
		&requestflag.Flag[bool]{
			Name:      "include-archived",
			Usage:     "Whether to include archived certificates in the results. Defaults to false.",
			QueryPath: "include_archived",
		},
		&requestflag.Flag[int64]{
			Name:      "limit",
			Usage:     "Maximum number of certificates to return per page. Defaults to 20, maximum 1000.",
			QueryPath: "limit",
		},
		&requestflag.Flag[string]{
			Name:      "page",
			Usage:     "Opaque pagination cursor from a previous `list_tunnel_certificates` response.",
			QueryPath: "page",
		},
		&requestflag.Flag[[]string]{
			Name:       "beta",
			Usage:      "Optional header to specify the beta version(s) you want to use.",
			HeaderPath: "anthropic-beta",
		},
		&requestflag.Flag[int64]{
			Name:  "max-items",
			Usage: "The maximum number of items to return (use -1 for unlimited).",
		},
	},
	Action:          handleBetaTunnelsCertificatesList,
	HideHelpCommand: true,
}

var betaTunnelsCertificatesArchive = cli.Command{
	Name:    "archive",
	Usage:   "The Tunnels API is in research preview. It requires the\n`anthropic-beta: mcp-tunnels-2026-06-22` header and may change without a\ndeprecation period. It supersedes the Admin API endpoints at\n`/v1/organizations/tunnels`, which remain available during a migration window.",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "tunnel-id",
			Required:  true,
			PathParam: "tunnel_id",
		},
		&requestflag.Flag[string]{
			Name:      "certificate-id",
			Required:  true,
			PathParam: "certificate_id",
		},
		&requestflag.Flag[[]string]{
			Name:       "beta",
			Usage:      "Optional header to specify the beta version(s) you want to use.",
			HeaderPath: "anthropic-beta",
		},
	},
	Action:          handleBetaTunnelsCertificatesArchive,
	HideHelpCommand: true,
}

func handleBetaTunnelsCertificatesCreate(ctx context.Context, cmd *cli.Command) error {
	client := anthropic.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("tunnel-id") && len(unusedArgs) > 0 {
		cmd.Set("tunnel-id", unusedArgs[0])
		unusedArgs = unusedArgs[1:]
	}
	if len(unusedArgs) > 0 {
		return fmt.Errorf("Unexpected extra arguments: %v", unusedArgs)
	}

	options, err := flagOptions(
		cmd,
		apiquery.NestedQueryFormatBrackets,
		apiquery.ArrayQueryFormatBrackets,
		ApplicationJSON,
		false,
	)
	if err != nil {
		return err
	}

	params := anthropic.BetaTunnelCertificateNewParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Beta.Tunnels.Certificates.New(
		ctx,
		cmd.Value("tunnel-id").(string),
		params,
		options...,
	)
	if err != nil {
		return err
	}

	obj := gjson.ParseBytes(res)
	format := cmd.Root().String("format")
	explicitFormat := cmd.Root().IsSet("format")
	transform := cmd.Root().String("transform")
	return ShowJSON(obj, ShowJSONOpts{
		ExplicitFormat: explicitFormat,
		Format:         format,
		RawOutput:      cmd.Root().Bool("raw-output"),
		Title:          "beta:tunnels:certificates create",
		Transform:      transform,
	})
}

func handleBetaTunnelsCertificatesRetrieve(ctx context.Context, cmd *cli.Command) error {
	client := anthropic.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("certificate-id") && len(unusedArgs) > 0 {
		cmd.Set("certificate-id", unusedArgs[0])
		unusedArgs = unusedArgs[1:]
	}
	if len(unusedArgs) > 0 {
		return fmt.Errorf("Unexpected extra arguments: %v", unusedArgs)
	}

	options, err := flagOptions(
		cmd,
		apiquery.NestedQueryFormatBrackets,
		apiquery.ArrayQueryFormatBrackets,
		EmptyBody,
		false,
	)
	if err != nil {
		return err
	}

	params := anthropic.BetaTunnelCertificateGetParams{
		TunnelID: cmd.Value("tunnel-id").(string),
	}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Beta.Tunnels.Certificates.Get(
		ctx,
		cmd.Value("certificate-id").(string),
		params,
		options...,
	)
	if err != nil {
		return err
	}

	obj := gjson.ParseBytes(res)
	format := "explore"
	explicitFormat := cmd.Root().IsSet("format")
	if explicitFormat {
		format = cmd.Root().String("format")
	}
	transform := cmd.Root().String("transform")
	return ShowJSON(obj, ShowJSONOpts{
		ExplicitFormat: explicitFormat,
		Format:         format,
		RawOutput:      cmd.Root().Bool("raw-output"),
		Title:          "beta:tunnels:certificates retrieve",
		Transform:      transform,
	})
}

func handleBetaTunnelsCertificatesList(ctx context.Context, cmd *cli.Command) error {
	client := anthropic.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("tunnel-id") && len(unusedArgs) > 0 {
		cmd.Set("tunnel-id", unusedArgs[0])
		unusedArgs = unusedArgs[1:]
	}
	if len(unusedArgs) > 0 {
		return fmt.Errorf("Unexpected extra arguments: %v", unusedArgs)
	}

	options, err := flagOptions(
		cmd,
		apiquery.NestedQueryFormatBrackets,
		apiquery.ArrayQueryFormatBrackets,
		EmptyBody,
		false,
	)
	if err != nil {
		return err
	}

	params := anthropic.BetaTunnelCertificateListParams{}

	format := "explore"
	explicitFormat := cmd.Root().IsSet("format")
	if explicitFormat {
		format = cmd.Root().String("format")
	}
	transform := cmd.Root().String("transform")
	if format == "raw" {
		var res []byte
		options = append(options, option.WithResponseBodyInto(&res))
		_, err = client.Beta.Tunnels.Certificates.List(
			ctx,
			cmd.Value("tunnel-id").(string),
			params,
			options...,
		)
		if err != nil {
			return err
		}
		obj := gjson.ParseBytes(res)
		return ShowJSON(obj, ShowJSONOpts{
			ExplicitFormat: explicitFormat,
			Format:         format,
			RawOutput:      cmd.Root().Bool("raw-output"),
			Title:          "beta:tunnels:certificates list",
			Transform:      transform,
		})
	} else {
		iter := client.Beta.Tunnels.Certificates.ListAutoPaging(
			ctx,
			cmd.Value("tunnel-id").(string),
			params,
			options...,
		)
		maxItems := int64(-1)
		if cmd.IsSet("max-items") {
			maxItems = cmd.Value("max-items").(int64)
		}
		return ShowJSONIterator(iter, maxItems, ShowJSONOpts{
			ExplicitFormat: explicitFormat,
			Format:         format,
			RawOutput:      cmd.Root().Bool("raw-output"),
			Title:          "beta:tunnels:certificates list",
			Transform:      transform,
		})
	}
}

func handleBetaTunnelsCertificatesArchive(ctx context.Context, cmd *cli.Command) error {
	client := anthropic.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("certificate-id") && len(unusedArgs) > 0 {
		cmd.Set("certificate-id", unusedArgs[0])
		unusedArgs = unusedArgs[1:]
	}
	if len(unusedArgs) > 0 {
		return fmt.Errorf("Unexpected extra arguments: %v", unusedArgs)
	}

	options, err := flagOptions(
		cmd,
		apiquery.NestedQueryFormatBrackets,
		apiquery.ArrayQueryFormatBrackets,
		EmptyBody,
		false,
	)
	if err != nil {
		return err
	}

	params := anthropic.BetaTunnelCertificateArchiveParams{
		TunnelID: cmd.Value("tunnel-id").(string),
	}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Beta.Tunnels.Certificates.Archive(
		ctx,
		cmd.Value("certificate-id").(string),
		params,
		options...,
	)
	if err != nil {
		return err
	}

	obj := gjson.ParseBytes(res)
	format := cmd.Root().String("format")
	explicitFormat := cmd.Root().IsSet("format")
	transform := cmd.Root().String("transform")
	return ShowJSON(obj, ShowJSONOpts{
		ExplicitFormat: explicitFormat,
		Format:         format,
		RawOutput:      cmd.Root().Bool("raw-output"),
		Title:          "beta:tunnels:certificates archive",
		Transform:      transform,
	})
}
