package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"
)

// federation holds the inputs for a JWT-bearer token exchange. The exchange
// itself is performed by the SDK's option.WithFederationTokenProvider; this
// type is only the CLI's flag/env carrier plus the "fully/partially/not
// configured" classification used by `auth status` and the precedence switch
// in getDefaultRequestOptions.
type federation struct {
	// Assertion is a literal signed JWT. Mutually exclusive with
	// AssertionFile.
	Assertion string
	// AssertionFile is a path to a file containing the signed JWT. When
	// set, the file is re-read on every AccessToken call so rotating
	// tokens (Kubernetes projected service account tokens, GitHub Actions
	// OIDC tokens) stay fresh.
	AssertionFile string
	// Rule is the tagged-ID of the OidcFederationRule governing this
	// exchange (prefix `fdrl_`). Required.
	Rule string
	// OrganizationID is the Anthropic org UUID whose credentials the
	// exchange should mint. Required.
	OrganizationID string
	// ServiceAccountID is an optional `svac_...` tagged-ID for
	// target_type=SERVICE_ACCOUNT federation rules. Leave empty for
	// user-target rules.
	ServiceAccountID string
}

// present reports whether enough of federation is populated to attempt an
// exchange. Assertion and AssertionFile are XOR; Rule and OrganizationID
// are both required; ServiceAccountID is always optional.
func (f federation) present() bool {
	hasAssertion := f.Assertion != "" || f.AssertionFile != ""
	return hasAssertion && f.Rule != "" && f.OrganizationID != ""
}

// AnySet reports whether any federation field is populated, even a partial
// (and therefore unusable) configuration. Used by callers to distinguish
// "the user tried to configure federation" from "nothing is set" — a partial
// config should surface a clear error rather than silently fall through to
// a lower-precedence credential source.
func (f federation) AnySet() bool {
	return f.Assertion != "" || f.AssertionFile != "" || f.Rule != "" ||
		f.OrganizationID != "" || f.ServiceAccountID != ""
}

// Missing returns the human-readable names of federation inputs that are
// required but unset. It returns nil when the configuration is either fully
// set (ready to exchange) or entirely unset (caller should fall through).
// Only a partial configuration yields a non-nil list.
func (f federation) Missing() []string {
	if !f.AnySet() || f.present() {
		return nil
	}
	var missing []string
	if f.Assertion == "" && f.AssertionFile == "" {
		missing = append(missing, "--identity-token[-file] / ANTHROPIC_IDENTITY_TOKEN[_FILE]")
	}
	if f.Rule == "" {
		missing = append(missing, "--federation-rule / ANTHROPIC_FEDERATION_RULE_ID")
	}
	if f.OrganizationID == "" {
		missing = append(missing, "--organization-id / ANTHROPIC_ORGANIZATION_ID")
	}
	return missing
}

// IdentityTokenFunc returns a callable that resolves the federation
// assertion (literal value or read from file each call). It mirrors the
// SDK's option.IdentityTokenFunc signature so callers can pass the result
// directly to option.WithFederationTokenProvider.
func (f federation) IdentityTokenFunc() (func(ctx context.Context) (string, error), error) {
	if f.Assertion != "" && f.AssertionFile != "" {
		return nil, fmt.Errorf("oauth: --identity-token and --identity-token-file are mutually exclusive")
	}
	if f.Assertion != "" {
		return func(context.Context) (string, error) { return f.Assertion, nil }, nil
	}
	if f.AssertionFile != "" {
		path := f.AssertionFile
		return func(context.Context) (string, error) {
			data, err := os.ReadFile(path)
			if err != nil {
				return "", fmt.Errorf("oauth: read identity token file %s: %w", path, err)
			}
			tok := strings.TrimSpace(string(data))
			if tok == "" {
				return "", fmt.Errorf("oauth: identity token file %s is empty", path)
			}
			return tok, nil
		}, nil
	}
	return nil, nil
}
