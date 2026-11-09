package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/anthropics/anthropic-cli/v2/internal/termsafe"
	"github.com/anthropics/anthropic-sdk-go/v2/config"
)

// revokeTimeout bounds one profile's revoke requests, so that logging out
// offline stays quick. It is a var so tests can shorten it.
var revokeTimeout = 5 * time.Second

// revokeAllTokens goes by the files in credentials/, not the profiles in
// configs/, so that tokens whose profile is gone still get a warning.
func revokeAllTokens(ctx context.Context, dir string) {
	entries, err := os.ReadDir(config.CredentialsDir(dir))
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		fmt.Fprintf(termsafe.Stderr, "warning: could not list the credentials whose tokens to revoke: %v\n", err)
	}
	for _, entry := range entries {
		if profile, ok := strings.CutSuffix(entry.Name(), ".json"); ok && !entry.IsDir() {
			revokeProfileTokens(ctx, dir, profile)
		}
	}
}

// revokeProfileTokens revokes on the server the tokens that logout is about to
// delete. A failure only warns, because logout has to work offline.
func revokeProfileTokens(ctx context.Context, dir, profile string) {
	if err := revokeStoredTokens(ctx, dir, profile); err != nil {
		fmt.Fprintf(termsafe.Stderr,
			"warning: could not revoke the tokens for profile %q on the server, so they stay valid until they expire: %v\n",
			profile, err)
	}
}

func revokeStoredTokens(ctx context.Context, dir, profile string) error {
	data, err := os.ReadFile(config.ProfileCredentialsPath(dir, profile))
	if errors.Is(err, fs.ErrNotExist) {
		return nil
	}
	if err != nil {
		return err
	}
	var creds config.Credentials
	// The parse error is dropped because it can quote the file.
	if json.Unmarshal(data, &creds) != nil {
		return errors.New("the credentials file could not be parsed")
	}
	cfg, err := loadProfileConfig(dir, profile)
	if err != nil {
		return err
	}
	// Only a login's tokens were issued to an OAuth client. A federation profile
	// caches a token that has none, and a profile with neither client_id nor
	// refresh token holds a hand-written bearer token.
	oauth := cfg.AuthenticationInfo.UserOAuth
	if oauth == nil || (oauth.ClientID == "" && creds.RefreshToken == "") {
		return nil
	}
	// The host and client that a refresh for this profile uses.
	baseURL := resolveBaseURL("", cfg)
	clientID := resolveClientID("", cfg)

	ctx, cancel := context.WithTimeout(ctx, revokeTimeout)
	defer cancel()
	var firstErr error
	for _, t := range []struct{ hint, token string }{
		{"refresh_token", creds.RefreshToken},
		{"access_token", creds.AccessToken},
	} {
		if t.token == "" {
			continue
		}
		if err := revokeToken(ctx, baseURL, clientID, t.token, t.hint); err != nil && firstErr == nil {
			firstErr = err
		}
	}
	return firstErr
}

// revokeToken revokes one token at /v1/oauth/token/revoke. It puts nothing
// from the request body or the response body in an error, only the status
// code, and passes a transport error on as net/http words it.
func revokeToken(ctx context.Context, baseURL, clientID, token, hint string) error {
	if err := config.RequireSecureTokenEndpoint(baseURL); err != nil {
		return err
	}
	body, err := json.Marshal(map[string]string{
		"token":           token,
		"token_type_hint": hint,
		"client_id":       clientID,
	})
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, "POST", baseURL+config.TokenEndpoint+"/revoke", bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := config.NewTokenHTTPClient()
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	// Drained so the next request can reuse the connection.
	_, _ = io.Copy(io.Discard, io.LimitReader(resp.Body, 4<<10))
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return fmt.Errorf("revoke endpoint returned HTTP %d", resp.StatusCode)
	}
	return nil
}
