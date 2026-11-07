package cmd

import (
	"cmp"
	"context"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/anthropics/anthropic-sdk-go/v2/bedrock"
	"github.com/anthropics/anthropic-sdk-go/v2/option"
	"github.com/urfave/cli/v3"
)

var bedrockBaseURLNoticeOnce sync.Once

// bedrockPlatform holds the inputs that address Claude in Amazon Bedrock, at
// the Messages API endpoint Bedrock serves at bedrock-mantle.{region}.api.aws.
// The SDK's bedrock package builds the URL and signs the request; this type is
// the CLI's flag/env carrier, like awsPlatform and foundryPlatform.
type bedrockPlatform struct {
	// Region picks the endpoint and the signing region, from --bedrock-region
	// or ANTHROPIC_BEDROCK_REGION. When empty the SDK falls back to
	// AWS_REGION, then AWS_DEFAULT_REGION.
	Region string
	// APIKey is a Bedrock API key from --bedrock-api-key-stdin or
	// AWS_BEARER_TOKEN_BEDROCK. When empty, requests are signed with AWS
	// credentials from the default credential chain.
	APIKey      string
	APIKeyPiped bool
	// BaseURL replaces the regional endpoint. Like the Google Cloud tier, this
	// tier never takes --base-url / ANTHROPIC_BASE_URL.
	BaseURL string
}

// bedrockFromRoot builds a bedrockPlatform from the root command's flag
// values, so a subcommand's same-named local flag can't shadow them.
// applyStdinCredential must run first for the API key.
func bedrockFromRoot(root *cli.Command) bedrockPlatform {
	apiKey, _ := directCredential("bedrock-api-key")
	return bedrockPlatform{
		Region:      root.String("bedrock-region"),
		APIKey:      apiKey,
		APIKeyPiped: credentialFromStdin(root, "bedrock-api-key"),
		BaseURL:     root.String("bedrock-base-url"),
	}
}

// AnySet reports whether the user addressed Bedrock at all: a region or a
// piped key. AWS_REGION and AWS_PROFILE do not count, as in the AWS tier, nor
// does AWS_BEARER_TOKEN_BEDROCK: other Bedrock clients read it too, so it is
// exported on machines whose `ant` is logged in elsewhere. Nor does the base
// URL, as in the Google Cloud tier.
func (b bedrockPlatform) AnySet() bool {
	return b.Region != "" || b.APIKeyPiped
}

// host is where this tier sends requests. Empty when neither a base URL nor a
// region is set anywhere, since the SDK then has no endpoint to derive.
func (b bedrockPlatform) host() string {
	region := cmp.Or(b.Region, os.Getenv("AWS_REGION"), os.Getenv("AWS_DEFAULT_REGION"))
	if b.BaseURL != "" || region == "" {
		return b.BaseURL
	}
	return fmt.Sprintf("https://bedrock-mantle.%s.api.aws/anthropic", region)
}

// noteIgnoredBaseURL tells the user, once, that the generic base URL is not
// applied, for the same reason as the Google Cloud tier's notice.
func (b bedrockPlatform) noteIgnoredBaseURL(w io.Writer, baseURL string) {
	if baseURL == "" {
		return
	}
	bedrockBaseURLNoticeOnce.Do(func() {
		fmt.Fprintf(w,
			"Note: --base-url / ANTHROPIC_BASE_URL is not used with Amazon Bedrock credentials; requests go to %s. To send them elsewhere: --bedrock-base-url / ANTHROPIC_BEDROCK_MANTLE_BASE_URL.\n",
			cmp.Or(b.host(), "the regional endpoint"))
	})
}

func (b bedrockPlatform) requestOptions(ctx context.Context) ([]option.RequestOption, error) {
	if err := ValidateBaseURL(b.BaseURL, "--bedrock-base-url / ANTHROPIC_BEDROCK_MANTLE_BASE_URL"); err != nil {
		return nil, err
	}
	client, err := bedrock.NewMantleClient(ctx, bedrock.MantleClientConfig{
		AWSRegion: b.Region,
		APIKey:    b.APIKey,
		BaseURL:   b.BaseURL,
	})
	if err != nil {
		return nil, fmt.Errorf("%w\n"+
			"  Set a region, and a Bedrock API key or AWS credentials.\n"+
			"  Region:          --bedrock-region / ANTHROPIC_BEDROCK_REGION / AWS_REGION / AWS_DEFAULT_REGION\n"+
			"  API key:         AWS_BEARER_TOKEN_BEDROCK / --bedrock-api-key-stdin\n"+
			"  AWS credentials: the AWS default credential chain (AWS_PROFILE, environment keys, an IAM role)\n"+
			"  Base URL:        --bedrock-base-url / ANTHROPIC_BEDROCK_MANTLE_BASE_URL", err)
	}
	return client.Options, nil
}
