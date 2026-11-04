package cmd

import (
	"cmp"
	"context"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/anthropics/anthropic-sdk-go/v2/foundry"
	"github.com/anthropics/anthropic-sdk-go/v2/option"
	"github.com/urfave/cli/v3"
)

var foundryBaseURLNoticeOnce, foundryAPIKeyNoticeOnce sync.Once

// foundryPlatform holds the inputs that address Claude in Microsoft Foundry.
// The SDK's foundry package builds the URL and attaches the credential; this
// type is the CLI's flag/env carrier, like googleCloud and awsPlatform.
type foundryPlatform struct {
	// Resource is the name in https://{resource}.services.ai.azure.com/anthropic,
	// from --foundry-resource or ANTHROPIC_FOUNDRY_RESOURCE.
	Resource         string
	ResourceFromFlag bool
	// APIKey is from --foundry-api-key-stdin or ANTHROPIC_FOUNDRY_API_KEY.
	APIKey      string
	APIKeyPiped bool
	// AuthToken is a Microsoft Entra ID access token from
	// --foundry-auth-token-stdin or ANTHROPIC_FOUNDRY_AUTH_TOKEN.
	AuthToken      string
	AuthTokenPiped bool
	// BaseURL takes the place of Resource. Like the Google Cloud tier, this
	// tier never takes --base-url / ANTHROPIC_BASE_URL.
	BaseURL string
}

// foundryFromRoot builds a foundryPlatform from the root command's flag
// values, so a subcommand's same-named local flag can't shadow them.
// applyStdinCredential must run first for the API key and the token.
func foundryFromRoot(root *cli.Command) foundryPlatform {
	apiKey, _ := directCredential("foundry-api-key")
	authToken, _ := directCredential("foundry-auth-token")
	resourceFlag := root.String("foundry-resource")
	return foundryPlatform{
		Resource:         cmp.Or(resourceFlag, os.Getenv("ANTHROPIC_FOUNDRY_RESOURCE")),
		ResourceFromFlag: resourceFlag != "",
		APIKey:           apiKey,
		APIKeyPiped:      credentialFromStdin(root, "foundry-api-key"),
		AuthToken:        authToken,
		AuthTokenPiped:   credentialFromStdin(root, "foundry-auth-token"),
		BaseURL:          root.String("foundry-base-url"),
	}
}

// AnySet reports whether the user addressed Foundry at all. The base URL does
// not count, as in the Google Cloud tier, nor does ANTHROPIC_FOUNDRY_RESOURCE:
// Claude Code's Entra ID setup exports it with no key or token.
func (f foundryPlatform) AnySet() bool {
	return f.ResourceFromFlag || f.APIKey != "" || f.AuthToken != ""
}

// usesAuthToken reports which credential is sent when both are set: the piped
// one, else the token, which is the order Claude Code documents.
func (f foundryPlatform) usesAuthToken() bool {
	return f.AuthToken != "" && !f.APIKeyPiped
}

// noteIgnoredAPIKey tells the user, once, which credential is sent when the
// environment sets both.
func (f foundryPlatform) noteIgnoredAPIKey(w io.Writer) {
	if f.APIKey == "" || f.AuthToken == "" || f.APIKeyPiped || f.AuthTokenPiped {
		return
	}
	foundryAPIKeyNoticeOnce.Do(func() {
		fmt.Fprintln(w, "Note: ANTHROPIC_FOUNDRY_API_KEY and ANTHROPIC_FOUNDRY_AUTH_TOKEN are both set in your environment; using ANTHROPIC_FOUNDRY_AUTH_TOKEN. Run `ant auth status` for details.")
	})
}

// host is where this tier sends requests. Empty when neither a base URL nor a
// resource is set.
func (f foundryPlatform) host() string {
	if f.BaseURL != "" || f.Resource == "" {
		return f.BaseURL
	}
	return fmt.Sprintf("https://%s.services.ai.azure.com/anthropic", f.Resource)
}

// noteIgnoredBaseURL tells the user, once, that the generic base URL is not
// applied, for the same reason as the Google Cloud tier's notice.
func (f foundryPlatform) noteIgnoredBaseURL(w io.Writer, baseURL string) {
	if baseURL == "" {
		return
	}
	foundryBaseURLNoticeOnce.Do(func() {
		fmt.Fprintf(w,
			"Note: --base-url / ANTHROPIC_BASE_URL is not used with Foundry credentials; requests go to %s. To send them elsewhere: --foundry-base-url / ANTHROPIC_FOUNDRY_BASE_URL.\n",
			f.host())
	})
}

func (f foundryPlatform) requestOptions() ([]option.RequestOption, error) {
	if err := ValidateBaseURL(f.BaseURL, "--foundry-base-url / ANTHROPIC_FOUNDRY_BASE_URL"); err != nil {
		return nil, err
	}
	cfg := foundry.ClientConfig{
		Resource: f.Resource,
		BaseURL:  f.BaseURL,
	}
	if f.usesAuthToken() {
		cfg.AzureADTokenProvider = func(context.Context) (string, error) { return f.AuthToken, nil }
	} else {
		cfg.APIKey = f.APIKey
	}
	client, err := foundry.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("%w\n"+
			"  Set one credential and one endpoint.\n"+
			"  API key:        ANTHROPIC_FOUNDRY_API_KEY / --foundry-api-key-stdin\n"+
			"  Entra ID token: ANTHROPIC_FOUNDRY_AUTH_TOKEN / --foundry-auth-token-stdin\n"+
			"  Resource:       --foundry-resource / ANTHROPIC_FOUNDRY_RESOURCE\n"+
			"  Base URL:       --foundry-base-url / ANTHROPIC_FOUNDRY_BASE_URL", err)
	}
	return client.Options, nil
}
