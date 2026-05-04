package cmd

import (
	"bufio"
	"bytes"
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"io"
	"io/fs"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/anthropics/anthropic-sdk-go/config"
	"github.com/urfave/cli/v3"
)

const (
	oauthClientIDProd = "41077d10-94b8-4194-be48-d251e9eb21b4"
	defaultConsoleURL = "https://platform.claude.com"
	defaultBaseURL    = "https://api.anthropic.com"
	oauthScope        = "user:profile user:inference user:developer"

	// betaUserOAuth is the anthropic-beta header value for user_oauth
	// credentials (interactive PKCE login). Required on the
	// authorization_code grant, the refresh_token grant, and on
	// authenticated requests that use a user_oauth access token.
	betaUserOAuth = "oauth-2025-04-20"
)

// tokenResponse mirrors the /v1/oauth/token response for the interactive
// authorization_code grant driven by `ant auth login`.
type tokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
	// Organization, Account, and Workspace come from oauth-server's
	// token_metadata_extension hook (api/oauth_server/oauth/grants/extensions.py)
	// on the AuthorizationCode and RefreshToken grants. Workspace is
	// populated when the authorization was bound to one (either via the
	// authorize URL workspace_id param or via the Console picker), and
	// carries the tagged `wrkspc_...` form so it round-trips with the
	// CLI flag, profile config, and anthropic-workspace-id header.
	Organization tokenOrganization `json:"organization"`
	Account      tokenAccount      `json:"account"`
	Workspace    tokenWorkspace    `json:"workspace"`
}

type tokenOrganization struct {
	UUID string `json:"uuid"`
	Name string `json:"name"`
}

type tokenAccount struct {
	UUID         string `json:"uuid"`
	EmailAddress string `json:"email_address"`
}

type tokenWorkspace struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func init() {
	Command.Commands = append(Command.Commands, &cli.Command{
		Name:    "auth",
		Usage:   "Manage authentication",
		Suggest: true,
		Commands: []*cli.Command{
			{
				Name:   "login",
				Usage:  "Log in via browser (OAuth)",
				Action: authLogin,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "console-url",
						Usage:   "Base URL for the OAuth authorize page (default: " + defaultConsoleURL + ")",
						Sources: cli.EnvVars("ANTHROPIC_CONSOLE_URL"),
					},
					&cli.StringFlag{
						Name:    "base-url",
						Usage:   "API base URL for token exchange (default: " + defaultBaseURL + ")",
						Sources: cli.EnvVars("ANTHROPIC_BASE_URL"),
						// Hidden so --help lists --base-url once (the root
						// command also defines it as a global flag). Kept as a
						// local flag for the ANTHROPIC_BASE_URL env source,
						// which the global flag doesn't provide. No Value: —
						// resolveBaseURL falls through to the profile's
						// base_url then defaultBaseURL.
						Hidden: true,
					},
					&cli.BoolFlag{
						Name:  "no-browser",
						Usage: "Print the authorize URL instead of opening a browser",
					},
					&cli.IntFlag{
						Name:  "callback-port",
						Usage: "Local port for the OAuth callback (0 = ephemeral). Use a fixed port when port-forwarding from a remote host.",
						Value: 0,
					},
					&cli.DurationFlag{
						Name:  "timeout",
						Usage: "How long to wait for the browser callback",
						Value: 5 * time.Minute,
					},
					&cli.StringFlag{
						Name:    "client-id",
						Usage:   "OAuth client_id (overrides the profile's stored value; defaults to the ant-cli prod client)",
						Sources: cli.EnvVars("ANTHROPIC_OAUTH_CLIENT_ID"),
					},
					&cli.StringFlag{
						Name:  "scope",
						Usage: "OAuth scope to request (space-separated; overrides the profile's stored scope or the default)",
					},
					&cli.StringFlag{
						Name:  "workspace-id",
						Usage: "Workspace to bind the access token to (optional). If omitted, Console shows a workspace picker after org selection. Find IDs under Settings → Workspaces in the Console (resolved from --console-url / profile / default).",
					},
					&cli.BoolFlag{
						Name:  "debug",
						Usage: "Print the token exchange status and Request-Id to stderr",
					},
				},
			},
			{
				Name:   "logout",
				Usage:  "Remove stored credentials",
				Action: authLogout,
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:  "all",
						Usage: "Remove all profiles and credentials",
					},
				},
			},
			{
				Name:   "status",
				Usage:  "Show current authentication status",
				Action: authStatus,
			},
			{
				Name:  "print-credentials",
				Usage: "Print credentials for the active profile, refreshing the access token if it is expired or near expiry",
				Description: `By default, prints the active profile's credentials as a JSON object
(the same shape as credentials/<profile>.json on disk).

With --access-token, prints only the bare access token — suitable for
command substitution into an Authorization header:

    curl https://api.anthropic.com/v1/messages \
      -H "Authorization: Bearer $(ant auth print-credentials --access-token)" ...

With --env, prints KEY=value lines (ANTHROPIC_AUTH_TOKEN, plus
ANTHROPIC_BASE_URL when the profile sets one) for use as a .env file.`,
				Action: authPrintCredentials,
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:  "access-token",
						Usage: "Print only the bare access token instead of the full JSON object",
					},
					&cli.BoolFlag{
						Name:  "env",
						Usage: "Print credentials as KEY=value lines (.env file format)",
					},
				},
			},
		},
	})
}

func authLogin(ctx context.Context, c *cli.Command) error {
	profile, dir := activeProfile(c)
	debug := c.Bool("debug")

	// `auth login` writes configs/<profile>.json only when the profile doesn't
	// exist yet. Re-login on an existing profile produces credentials only —
	// the config is user-authored intent and login mustn't rewrite it.
	// TOCTOU: another process could create the profile between this stat and
	// the SaveProfile call below (after the network round-trip). Accepted —
	// `auth login` is interactive and not expected to run concurrently against
	// the same profile; the race window is bounded by user approval latency.
	_, statErr := os.Stat(config.ProfilePath(dir, profile))
	bootstrapping := errors.Is(statErr, fs.ErrNotExist)
	if bootstrapping {
		fmt.Fprintf(os.Stderr, "Creating profile %q.\n", profile)
	}

	// Load any existing profile up front so we can (a) reuse its stored scope
	// and client_id for the authorize request when their flags aren't given,
	// and (b) hint the stored organization_id on the authorize URL.
	prev, prevErr := loadProfileConfig(dir, profile)
	if prevErr != nil && !errors.Is(prevErr, fs.ErrNotExist) {
		fmt.Fprintf(os.Stderr, "warning: ignoring unreadable profile config: %v\n", prevErr)
	}
	requestedScope := resolveRequestedScope(c.String("scope"), prev)
	clientID := resolveClientID(c.String("client-id"), prev)
	consoleURL := resolveConsoleURL(c.String("console-url"), prev)
	baseURL := resolveBaseURL(c.String("base-url"), prev)
	workspaceID := resolveWorkspaceID(c.String("workspace-id"), prev)
	// workspaceID may be empty here — Console will show a workspace picker
	// after the org selection step, and the resolved workspace comes back
	// in the token response (tok.Workspace.ID).

	verifier, err := randomURLSafeString(64)
	if err != nil {
		return fmt.Errorf("generate PKCE verifier: %w", err)
	}
	challenge := pkceChallengeS256(verifier)

	state, err := randomURLSafeString(32)
	if err != nil {
		return fmt.Errorf("generate state: %w", err)
	}

	wantOrg := ""
	if prev != nil {
		wantOrg = prev.OrganizationID
	}
	buildAuthorizeURL := func(redirectURI string) string {
		params := url.Values{
			"client_id":             {clientID},
			"redirect_uri":          {redirectURI},
			"response_type":         {"code"},
			"scope":                 {requestedScope},
			"state":                 {state},
			"code_challenge":        {challenge},
			"code_challenge_method": {"S256"},
		}
		// Omit when empty so Console renders the picker step.
		if workspaceID != "" {
			params.Set("workspace_id", workspaceID)
		}
		if wantOrg != "" {
			// Console's OAuthConsentPage reads ?orgUUID= and auto-switches to
			// that org (skipping the switcher) when the signed-in account is a
			// member, or shows a no-access error when not. The mismatch guard
			// below against the token response's organization.uuid remains the
			// CLI-side backstop.
			params.Set("orgUUID", wantOrg)
		}
		return fmt.Sprintf("%s/oauth/authorize?%s", consoleURL, params.Encode())
	}

	var code, redirectURI string
	switch loginMode(c) {
	case loginModeHeadless:
		// The authorize server renders the code on a public page the user
		// copies into this terminal. No local port is bound — the flow for
		// environments with no usable localhost at all (sandboxed
		// containers, CI shells).
		redirectURI = manualRedirectURIForAuthorize(consoleURL)
		if redirectURI == "" {
			return fmt.Errorf("--no-browser requires a valid --console-url (got %q)", consoleURL)
		}
		authorizeURL := buildAuthorizeURL(redirectURI)
		fmt.Fprintf(os.Stderr, "Open this URL to authorize:\n\n  %s\n\n", authorizeURL)
		code, err = readManualCode(ctx, os.Stderr, os.Stdin, state)
		if err != nil {
			return err
		}

	case loginModeSSHForward:
		// Bind the listener but don't open a browser — the user is on a
		// remote host reached via `ssh -L <N>:localhost:<N>` and opens the
		// URL on whatever machine their SSH session originates from; the
		// callback comes back through the tunnel.
		listener, rURI, lErr := startLoopbackListener(c.Int("callback-port"))
		if lErr != nil {
			return lErr
		}
		defer listener.Close()
		redirectURI = rURI
		authorizeURL := buildAuthorizeURL(redirectURI)
		fmt.Fprintf(os.Stderr, "Authorization URL (open on the machine where your SSH session originates):\n\n  %s\n\n", authorizeURL)
		fmt.Fprintf(os.Stderr, "Waiting for callback on %s ...\n", redirectURI)
		code, err = waitForCallback(ctx, listener, state, c.Duration("timeout"))
		if err != nil {
			return err
		}

	default: // loginModeDefault
		listener, loopbackRedirectURI, lErr := startLoopbackListener(c.Int("callback-port"))
		if lErr != nil {
			return lErr
		}
		defer listener.Close()
		loopbackURL := buildAuthorizeURL(loopbackRedirectURI)

		if err := openBrowser(loopbackURL); err != nil {
			fmt.Fprintf(os.Stderr, "(couldn't launch browser automatically: %v)\n", err)
		}

		// Print a paste-able URL so users whose browser can't reach this host
		// (SSH, container) can complete in one step. The printed URL must use
		// the manual redirect (Console renders the code on a page) — the
		// loopback URL would 302 to 127.0.0.1 with nothing to paste.
		manualRedirectURI := manualRedirectURIForAuthorize(consoleURL)
		if manualRedirectURI == "" {
			// No manual landing page for this consoleURL; fall back to
			// printing the loopback URL and waiting for the callback only.
			fmt.Fprintf(os.Stderr,
				"Opening your browser to authenticate. If it doesn't open, visit:\n\n  %s\n\n"+
					"Waiting for authentication...\n",
				loopbackURL)
			redirectURI = loopbackRedirectURI
			code, err = waitForCallback(ctx, listener, state, c.Duration("timeout"))
			if err != nil {
				return err
			}
			break
		}

		manualURL := buildAuthorizeURL(manualRedirectURI)
		fmt.Fprintf(os.Stderr,
			"Opening your browser to authenticate...\n"+
				"If your browser doesn't open or can't reach this machine, visit the URL below and paste the code here:\n\n"+
				"  %s\n\n"+
				"Waiting for authentication...\nCode: ",
			manualURL)

		// Race the loopback callback against a code pasted on stdin. First
		// result wins; the loser is cancelled via raceCtx. A stdin error (EOF,
		// closed pipe — common when non-interactive) is not a race result: it
		// just means the stdin path is unavailable, so the callback path
		// proceeds alone. A callback error is surfaced as before.
		raceCtx, raceCancel := context.WithCancel(ctx)
		defer raceCancel()
		type raceResult struct {
			code, redirectURI string
			err               error
		}
		resultCh := make(chan raceResult, 2)
		timeout := c.Duration("timeout")
		go func() {
			cbCode, cbErr := waitForCallback(raceCtx, listener, state, timeout)
			resultCh <- raceResult{code: cbCode, redirectURI: loopbackRedirectURI, err: cbErr}
		}()
		go func() {
			// readManualCode returns on raceCtx cancel, but its inner
			// bufio.ReadString goroutine remains blocked on stdin until
			// EOF or process exit — there's no portable way to interrupt
			// a blocking read on os.Stdin. Acceptable here: login is the
			// last thing this process does before exit.
			pasted, pErr := readManualCode(raceCtx, io.Discard, stdin, state)
			if pErr != nil {
				return
			}
			resultCh <- raceResult{code: pasted, redirectURI: manualRedirectURI}
		}()
		r := <-resultCh
		raceCancel()
		if r.err != nil {
			return r.err
		}
		if r.redirectURI == loopbackRedirectURI {
			// Callback won; wipe the dangling "Code: " prompt.
			fmt.Fprint(os.Stderr, "\r      \r")
		}
		code, redirectURI = r.code, r.redirectURI
	}

	tok, err := exchangeCode(ctx, baseURL, clientID, code, verifier, redirectURI, state, debug)
	if err != nil {
		return fmt.Errorf("exchange code for token: %w", err)
	}

	gotOrg := tok.Organization.UUID
	if wantOrg != "" && gotOrg != "" && wantOrg != gotOrg {
		// Console's OAuthConsentPage enforces ?orgUUID (auto-switches members,
		// blocks non-members), so a post-exchange mismatch should only happen
		// if the authorize URL was hand-edited or via a Console bug. Warn and
		// proceed rather than discarding a valid token.
		fmt.Fprintf(os.Stderr,
			"⚠  Authorized into organization %s (%s), but profile %q has organization_id=%s.\n"+
				"   Console enforces the requested org, so this likely means the authorize URL was edited.\n"+
				"   Run `ant profile set organization_id %s --profile %s` to retarget the profile, or re-login.\n",
			gotOrg, tok.Organization.Name, profile, wantOrg, gotOrg, profile)
	}

	// Effective workspace: prefer what the token is actually bound to (the
	// authoritative source — backend mints the token against this), fall
	// back to the flag value when the token response omits it (older
	// backend that doesn't yet emit the workspace block).
	effectiveWorkspaceID := tok.Workspace.ID
	if effectiveWorkspaceID == "" {
		effectiveWorkspaceID = workspaceID
	}
	if effectiveWorkspaceID == "" {
		return fmt.Errorf("no workspace bound to the issued token. " +
			"Pass --workspace-id, set workspace_id on the profile, " +
			"or pick a workspace in the Console consent page.")
	}

	if bootstrapping {
		cfg := &config.Config{
			AuthenticationInfo: &config.AuthenticationInfo{
				Type: config.AuthenticationTypeUserOAuth,
				// client_id is always written: the profile file is shared
				// with the SDKs, whose refresh path requires it; and the
				// refresh_token is bound to this client, so it's a property
				// of the credential — not a default that might change.
				UserOAuth: &config.UserOAuth{ClientID: clientID},
			},
			OrganizationID: gotOrg,
			WorkspaceID:    effectiveWorkspaceID,
		}
		// scope/base_url/console_url are only-if-explicitly-set: those are
		// runtime defaults that the CLI/SDK fills in.
		if c.IsSet("scope") {
			scope := tok.Scope
			if scope == "" {
				scope = requestedScope
			}
			cfg.AuthenticationInfo.UserOAuth.Scope = scope
		}
		if c.IsSet("base-url") && baseURL != defaultBaseURL {
			cfg.BaseURL = baseURL
		}
		if c.IsSet("console-url") {
			cfg.AuthenticationInfo.UserOAuth.ConsoleURL = consoleURL
		}
		if err := config.SaveProfile(dir, profile, cfg); err != nil {
			return fmt.Errorf("write config: %w", err)
		}
	} else if prev != nil {
		// Re-login never rewrites the profile config — creds record the
		// ephemeral fact (what this token is bound to); config records
		// durable intent set via `ant profile set`. When they diverge,
		// surface it so the user can decide; `auth status` shows the
		// drift until they act.
		switch {
		case prev.WorkspaceID == effectiveWorkspaceID:
			// No drift.
		case prev.WorkspaceID == "":
			fmt.Fprintf(os.Stderr,
				"→ Token bound to workspace %q (%s).\n"+
					"  To make this the default for profile %q: ant profile set workspace_id %s --profile %s\n",
				tok.Workspace.Name, effectiveWorkspaceID, profile, effectiveWorkspaceID, profile)
		default:
			fmt.Fprintf(os.Stderr,
				"⚠ Token bound to workspace %q (%s), but profile %q targets %s.\n"+
					"  To retarget this profile:      ant profile set workspace_id %s --profile %s\n"+
					"  Or create a separate profile:  ant auth login --profile <new> --workspace-id %s\n",
				tok.Workspace.Name, effectiveWorkspaceID, profile, prev.WorkspaceID,
				effectiveWorkspaceID, profile, effectiveWorkspaceID)
		}
	}

	expiresAt := time.Now().Add(time.Duration(tok.ExpiresIn) * time.Second)
	credsPath := config.ProfileCredentialsPath(dir, profile)
	if err := config.WriteCredentials(credsPath, config.Credentials{
		AccessToken:      tok.AccessToken,
		RefreshToken:     tok.RefreshToken,
		ExpiresAt:        &expiresAt,
		Scope:            tok.Scope,
		OrganizationUUID: tok.Organization.UUID,
		OrganizationName: tok.Organization.Name,
		AccountEmail:     tok.Account.EmailAddress,
		WorkspaceID:      tok.Workspace.ID,
		WorkspaceName:    tok.Workspace.Name,
	}); err != nil {
		return fmt.Errorf("write credentials: %w", err)
	}

	// Decide whether this login should also become the active profile:
	//   - --profile/ANTHROPIC_PROFILE explicitly given → always activate
	//     (the user named a target; make subsequent commands use it).
	//   - profile came from active_config / "default" → only write
	//     active_config when none exists yet (don't retarget by accident).
	prevActive := ""
	if data, err := os.ReadFile(config.ActiveConfigPath(dir)); err == nil {
		prevActive = strings.TrimSpace(string(data))
	}
	wantActivate := c.IsSet("profile") || prevActive == ""
	activated := false
	if wantActivate && prevActive != profile {
		if err := config.SetActiveProfile(dir, profile); err != nil {
			return fmt.Errorf("set active profile: %w", err)
		}
		activated = true
	}

	if who := loggedInSummary(tok.Organization.Name, tok.Account.EmailAddress); who != "" {
		fmt.Fprintf(os.Stderr, "✓ Logged in %s (profile %q).\n", who, profile)
	} else {
		fmt.Fprintf(os.Stderr, "✓ Logged in as profile %q.\n", profile)
	}
	if tok.Organization.UUID != "" {
		fmt.Fprintf(os.Stderr, "  organization: %s (%s)\n", tok.Organization.Name, tok.Organization.UUID)
	}
	if activated {
		if prevActive == "" {
			fmt.Fprintf(os.Stderr, "  → set as active profile\n")
		} else {
			fmt.Fprintf(os.Stderr, "  → set as active profile (was %q)\n", prevActive)
		}
	}
	fmt.Fprintf(os.Stderr, "  config:       %s\n  credentials:  %s\n",
		config.ProfilePath(dir, profile), credsPath)
	return nil
}

func authLogout(ctx context.Context, c *cli.Command) error {
	dir := config.DefaultDir()
	if c.Bool("all") {
		// Remove only what we own, never RemoveAll(dir) — DefaultDir() respects
		// ANTHROPIC_CONFIG_DIR, so a misconfigured env var pointing at $HOME
		// would otherwise be catastrophic.
		for _, p := range []string{
			config.ProfilesDir(dir),
			config.CredentialsDir(dir),
			config.ActiveConfigPath(dir),
		} {
			if err := os.RemoveAll(p); err != nil && !os.IsNotExist(err) {
				return err
			}
		}
		fmt.Fprintln(os.Stderr, "✓ Removed all profiles and credentials.")
		return nil
	}
	profile, _ := activeProfile(c)
	// Remove only the credential file: configs/<profile>.json is intentionally
	// left in place so workspace_id/base_url survive a re-login (per the
	// credentials-file spec, the config half is non-secret).
	path := config.ProfileCredentialsPath(dir, profile)
	if err := os.Remove(path); err != nil {
		if os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "Not logged in (profile %q).\n", profile)
			return nil
		}
		return err
	}
	fmt.Fprintf(os.Stderr, "✓ Logged out of profile %q.\n", profile)
	return nil
}

// storedCredentials is the subset of credentials/<profile>.json that the CLI
// reads back for status. Scope/Organization*/AccountEmail
// reflect what the current token was granted at mint time (written on every
// login), as opposed to the profile config's organization_id/scope which
// record user intent.
type storedCredentials struct {
	AccessToken      string `json:"access_token"`
	ExpiresAt        int64  `json:"expires_at"`
	Scope            string `json:"scope"`
	OrganizationUUID string `json:"organization_uuid"`
	OrganizationName string `json:"organization_name"`
	AccountEmail     string `json:"account_email"`
	WorkspaceID      string `json:"workspace_id"`
	WorkspaceName    string `json:"workspace_name"`
}

// loggedInSummary renders the human-readable "to <org> as <email>" fragment
// shared by `auth login`'s success line and `auth status`'s headline.
// Degrades gracefully when one or both pieces are missing (older creds files
// may lack them).
func loggedInSummary(orgName, accountEmail string) string {
	switch {
	case orgName != "" && accountEmail != "":
		return fmt.Sprintf("to %s as %s", orgName, accountEmail)
	case accountEmail != "":
		return "as " + accountEmail
	case orgName != "":
		return "to " + orgName
	default:
		return ""
	}
}

// resolveCredentialsPath returns the path readCredentials/readFullCredentials
// load from: the profile config's explicit credentials_path if set, else the
// SDK default credentials/<profile>.json.
func resolveCredentialsPath(cfg *config.Config, dir, profile string) string {
	if cfg.AuthenticationInfo != nil && cfg.AuthenticationInfo.CredentialsPath != "" {
		return cfg.AuthenticationInfo.CredentialsPath
	}
	return config.ProfileCredentialsPath(dir, profile)
}

// readCredentials loads the on-disk token for cfg/profile, returning the path
// it read from alongside the parsed contents.
func readCredentials(cfg *config.Config, dir, profile string) (storedCredentials, string, error) {
	credsPath := resolveCredentialsPath(cfg, dir, profile)
	data, err := os.ReadFile(credsPath)
	if err != nil {
		return storedCredentials{}, credsPath, err
	}
	var creds storedCredentials
	if err := json.Unmarshal(data, &creds); err != nil {
		return storedCredentials{}, credsPath, fmt.Errorf("parse credentials %s: %w", credsPath, err)
	}
	return creds, credsPath, nil
}

func authStatus(ctx context.Context, c *cli.Command) error {
	dir := config.DefaultDir()
	profile, profileSource := activeProfileWithSource(c, dir)

	out := os.Stdout
	fmt.Fprintf(out, "Active profile:  %s (%s)\n", profile, profileSource)
	fmt.Fprintf(out, "Config dir:      %s\n", dir)
	fmt.Fprintf(out, "Profile config:  %s\n", config.ProfilePath(dir, profile))
	fmt.Fprintf(out, "Credentials:     %s\n", config.ProfileCredentialsPath(dir, profile))

	cfg, cfgErr := loadProfileConfig(dir, profile)
	var creds storedCredentials
	var credsErr error
	if cfgErr == nil && cfg != nil {
		creds, _, credsErr = readCredentials(cfg, dir, profile)
	}
	profileTokenPresent := cfgErr == nil && cfg != nil && credsErr == nil && creds.AccessToken != ""

	root := c.Root()
	apiKeySet := root.IsSet("api-key")
	authTokenSet := root.IsSet("auth-token")
	fed := federationFromRoot(root)
	fedReady := fed.AnySet() && len(fed.Missing()) == 0
	fedMissing := fed.Missing()

	// Credential tier mirrors getDefaultRequestOptions exactly. Note that
	// partial federation config (fed.AnySet() but not fedReady) does NOT
	// claim the slot — it fails loudly at request time instead. The warning
	// below surfaces that state explicitly.
	profileExplicit := profileIsExplicit(c)
	credWinner := 0
	switch {
	case apiKeySet:
		credWinner = 1
	case authTokenSet:
		credWinner = 2
	case profileTokenPresent && profileExplicit:
		credWinner = 3
	case fedReady:
		credWinner = 4
	case profileTokenPresent:
		credWinner = 5
	}

	fmt.Fprintln(out)
	fmt.Fprintln(out, "Credentials")
	if credWinner == 3 || credWinner == 5 {
		if who := loggedInSummary(creds.OrganizationName, creds.AccountEmail); who != "" {
			fmt.Fprintf(out, "  Logged in %s\n", who)
		}
	}
	switch {
	case credWinner == 0 && cfgErr != nil && errors.Is(cfgErr, os.ErrNotExist):
		fmt.Fprintf(out, "  (profile %q not configured — run `ant auth login` to set it up)\n", profile)
	case credWinner == 0 && cfgErr != nil:
		fmt.Fprintf(out, "  (profile config unreadable: %v)\n", cfgErr)
	case credWinner == 0 && cfgErr == nil && cfg != nil && credsErr != nil && os.IsNotExist(credsErr):
		fmt.Fprintf(out, "  (profile %q configured but not logged in — run `ant auth login`)\n", profile)
	case credWinner == 0 && cfgErr == nil && cfg != nil && credsErr != nil:
		fmt.Fprintf(out, "  (credentials unreadable: %v)\n", credsErr)
	case credWinner == 0:
		fmt.Fprintln(out, "  (no credential configured — set ANTHROPIC_API_KEY or run `ant auth login`)")
	}
	if apiKeySet {
		writeRow(out, credWinner == 1, "--api-key / ANTHROPIC_API_KEY", formatSecret(root.String("api-key"), true))
	}
	if authTokenSet {
		writeRow(out, credWinner == 2, "--auth-token / ANTHROPIC_AUTH_TOKEN", formatSecret(root.String("auth-token"), true))
	}
	if profileTokenPresent {
		authType := "unknown"
		if cfg.AuthenticationInfo != nil {
			authType = string(cfg.AuthenticationInfo.Type)
		}
		profileLabel := fmt.Sprintf("Profile (%s)", authType)
		if profileExplicit {
			profileLabel += " [via --profile/ANTHROPIC_PROFILE]"
		} else {
			profileLabel += " [via active_config]"
		}
		writeRow(out, credWinner == 3 || credWinner == 5, profileLabel, formatSecret(creds.AccessToken, true))
		exp := time.Unix(creds.ExpiresAt, 0)
		writeDetail(out, "expires", fmt.Sprintf("%s (%s)", exp.Format(time.RFC3339), formatRemaining(time.Until(exp))))
		// Prefer the credentials file's scope (what the live token was granted)
		// over the config file's scope (what was requested at bootstrap).
		scope := creds.Scope
		if scope == "" && cfg.AuthenticationInfo != nil && cfg.AuthenticationInfo.UserOAuth != nil {
			scope = cfg.AuthenticationInfo.UserOAuth.Scope
		}
		if scope != "" {
			writeDetail(out, "scope", scope)
		}
		if creds.OrganizationUUID != "" {
			writeDetail(out, "organization", fmt.Sprintf("%s (%s)", creds.OrganizationName, creds.OrganizationUUID))
		}
		if cfg.AuthenticationInfo != nil && cfg.AuthenticationInfo.UserOAuth != nil {
			if uo := cfg.AuthenticationInfo.UserOAuth; uo.ClientID != "" {
				writeDetail(out, "client_id", uo.ClientID)
			}
		}
		if cfg.AuthenticationInfo != nil && cfg.AuthenticationInfo.OIDCFederation != nil {
			fed := cfg.AuthenticationInfo.OIDCFederation
			if fed.FederationRuleID != "" {
				writeDetail(out, "federation_rule_id", fed.FederationRuleID)
			}
			if fed.ServiceAccountID != "" {
				writeDetail(out, "service_account_id", fed.ServiceAccountID)
			}
			if fed.Scope != "" {
				writeDetail(out, "scope", fed.Scope)
			}
		}
	}
	if fedReady {
		writeRow(out, credWinner == 4, "Federation (jwt-bearer)", "see 'Federation inputs' below")
	} else if fed.AnySet() {
		writeRow(out, false, "Federation (jwt-bearer)", "partial — missing required inputs")
	}

	// Surface the surprising-override case: the user ran `ant auth login` but
	// has a credential env var set that silently beats the profile on every
	// request.
	if profileTokenPresent && (apiKeySet || authTokenSet) {
		overrideEnv := "ANTHROPIC_API_KEY"
		if !apiKeySet {
			overrideEnv = "ANTHROPIC_AUTH_TOKEN"
		}
		fmt.Fprintln(out)
		fmt.Fprintf(out, "⚠  %s is set in your environment and overrides the logged-in profile.\n", overrideEnv)
		fmt.Fprintf(out, "   Unset it to use the profile:  unset %s\n", overrideEnv)
	}

	// Surface partial federation config: the user set some federation inputs
	// but forgot one of the required ones. The next real API call will error
	// out — warn now so the user can fix it before they hit that.
	if len(fedMissing) > 0 {
		fmt.Fprintln(out)
		fmt.Fprintln(out, "⚠  Federation is partially configured and will not be used until fully set.")
		fmt.Fprintln(out, "   Missing:")
		for _, m := range fedMissing {
			fmt.Fprintf(out, "     - %s\n", m)
		}
	}

	fmt.Fprintln(out)
	fmt.Fprintln(out, "Base URL")
	flagBaseURL := root.String("base-url")
	envBaseURL := os.Getenv("ANTHROPIC_BASE_URL")
	profileBaseURL := ""
	if cfg != nil {
		profileBaseURL = cfg.BaseURL
	}
	baseWinner := 4 // SDK default when nothing else is set
	switch {
	case flagBaseURL != "":
		baseWinner = 1
	case envBaseURL != "":
		baseWinner = 2
	case profileBaseURL != "":
		baseWinner = 3
	}
	if flagBaseURL != "" {
		writeRow(out, baseWinner == 1, "--base-url flag", flagBaseURL)
	}
	if envBaseURL != "" {
		writeRow(out, baseWinner == 2, "ANTHROPIC_BASE_URL env", envBaseURL)
	}
	if profileBaseURL != "" {
		writeRow(out, baseWinner == 3, "Profile base_url", profileBaseURL)
	}
	if baseWinner == 4 {
		writeRow(out, true, "SDK default", defaultBaseURL)
	}

	fmt.Fprintln(out)
	fmt.Fprintln(out, "Organization")
	flagOrg := root.String("organization-id")
	orgFromFlag := root.IsSet("organization-id")
	profileOrg := ""
	if cfg != nil {
		profileOrg = cfg.OrganizationID
	}
	orgWinner := 3
	switch {
	case orgFromFlag:
		orgWinner = 1
	case profileOrg != "":
		orgWinner = 2
	}
	if orgFromFlag {
		writeRow(out, orgWinner == 1, "--organization-id / ANTHROPIC_ORGANIZATION_ID", flagOrg)
	}
	if profileOrg != "" {
		writeRow(out, orgWinner == 2, "Profile organization_id", profileOrg)
	}
	if orgWinner == 3 {
		if creds.OrganizationUUID != "" {
			writeRow(out, true, "Token organization (from credentials)", fmt.Sprintf("%s (%s)", creds.OrganizationName, creds.OrganizationUUID))
		} else {
			writeRow(out, true, "Derived from credential (server-side)", "(determined at request time)")
		}
	}

	fmt.Fprintln(out)
	fmt.Fprintln(out, "Workspace")
	profileWs := ""
	if cfg != nil {
		profileWs = cfg.WorkspaceID
	}
	// The credentials file's workspace_id is what the token is actually
	// bound to (set at mint time); the profile's workspace_id is user
	// intent. Show both when they diverge so silent drift is visible.
	credsWs := creds.WorkspaceID
	credsWsName := creds.WorkspaceName
	switch {
	case credsWs != "" && profileWs != "" && credsWs == profileWs:
		label := credsWs
		if credsWsName != "" {
			label = fmt.Sprintf("%s (%s)", credsWs, credsWsName)
		}
		writeRow(out, true, "Workspace", label)
	case credsWs != "" && profileWs != "" && credsWs != profileWs:
		writeRow(out, false, "Profile workspace_id", profileWs)
		credsLabel := credsWs
		if credsWsName != "" {
			credsLabel = fmt.Sprintf("%s (%s)", credsWs, credsWsName)
		}
		writeRow(out, true, "Active token workspace", credsLabel+" — drift, re-login to reconcile")
	case credsWs != "":
		credsLabel := credsWs
		if credsWsName != "" {
			credsLabel = fmt.Sprintf("%s (%s)", credsWs, credsWsName)
		}
		writeRow(out, true, "Active token workspace", credsLabel)
	case profileWs != "":
		writeRow(out, true, "Profile workspace_id", profileWs)
	default:
		writeRow(out, true, "Server-side default", "(determined at request time)")
	}

	// Federation inputs section — only when at least one is set.
	fedFlags := []struct {
		flag, env string
		secret    bool
	}{
		{"identity-token", "ANTHROPIC_IDENTITY_TOKEN", true},
		{"identity-token-file", "ANTHROPIC_IDENTITY_TOKEN_FILE", false},
		{"federation-rule", "ANTHROPIC_FEDERATION_RULE_ID", false},
		{"organization-id", "ANTHROPIC_ORGANIZATION_ID", false},
		{"service-account-id", "ANTHROPIC_SERVICE_ACCOUNT_ID", false},
	}
	anyFed := false
	for _, f := range fedFlags {
		if root.IsSet(f.flag) {
			anyFed = true
			break
		}
	}
	if anyFed {
		fmt.Fprintln(out)
		fmt.Fprintln(out, "Federation inputs")
		for _, f := range fedFlags {
			if !root.IsSet(f.flag) {
				continue
			}
			val := root.String(f.flag)
			display := val
			if f.secret {
				display = formatSecret(val, true)
			}
			fmt.Fprintf(out, "  --%-22s / %-32s %s\n", f.flag, f.env, display)
		}
	}

	return nil
}

// activeProfileWithSource mirrors activeProfile but returns a human-readable
// description of which tier resolved the active profile, for use in status
// output. The global --profile flag's Sources include ANTHROPIC_PROFILE, so
// IsSet covers both; the LookupEnv branch is a defensive fallback for callers
// that pass a Command not yet Run() (or nil).
func activeProfileWithSource(c *cli.Command, dir string) (profile, source string) {
	const explicitSrc = "from --profile / ANTHROPIC_PROFILE"
	if c != nil && c.IsSet("profile") {
		return c.String("profile"), explicitSrc
	}
	if p, ok := os.LookupEnv("ANTHROPIC_PROFILE"); ok {
		return p, explicitSrc
	}
	activePath := config.ActiveConfigPath(dir)
	if data, err := os.ReadFile(activePath); err == nil {
		if p := strings.TrimSpace(string(data)); p != "" {
			return p, "from active_config file"
		}
	}
	return "default", "fallback — no active_config set"
}

// federationFromRoot builds a federation from the root command's flag values
// so authStatus can reuse the type's "fully/partially/not configured"
// classification for status display.
func federationFromRoot(root *cli.Command) federation {
	return federation{
		Assertion:        root.String("identity-token"),
		AssertionFile:    root.String("identity-token-file"),
		Rule:             root.String("federation-rule"),
		OrganizationID:   root.String("organization-id"),
		ServiceAccountID: root.String("service-account-id"),
	}
}

// writeRow emits a source row for a status section. The active-row marker
// mirrors the indentation of the "(active) *" prefix so non-active rows
// line up under the active one.
func writeRow(w io.Writer, active bool, label, value string) {
	marker := "           * "
	if active {
		marker = "  (active) * "
	}
	fmt.Fprintf(w, "%s%-46s %s\n", marker, label, value)
}

// writeDetail emits an indented "key: value" line underneath a row, for
// multi-line source details (token expiry, scope, etc.). Aligned so the
// key column sits under the row's label column.
func writeDetail(w io.Writer, key, value string) {
	fmt.Fprintf(w, "               %-15s %s\n", key+":", value)
}

// formatSecret masks sensitive values so token prefixes are visible for
// disambiguation ("which key is this?") without leaking the full secret.
func formatSecret(s string, set bool) string {
	if !set || s == "" {
		return "(unset)"
	}
	if len(s) <= 16 {
		return "(set, " + fmt.Sprint(len(s)) + " chars)"
	}
	return s[:16] + "..."
}

// formatRemaining renders a time.Duration as a short human string for token
// expiry display.
func formatRemaining(d time.Duration) string {
	if d <= 0 {
		return "expired " + (-d).Round(time.Second).String() + " ago"
	}
	d = d.Round(time.Second)
	if d >= time.Hour {
		return "in " + d.Round(time.Minute).String()
	}
	return "in " + d.String()
}

// printCredentialsRefreshThreshold is the remaining-lifetime cutoff below
// which print-credentials refreshes before printing. Matches the SDK's
// advisory-refresh threshold so a printed token has enough headroom to
// survive the caller's subsequent request.
const printCredentialsRefreshThreshold = 120 * time.Second

func authPrintCredentials(ctx context.Context, c *cli.Command) error {
	if c.Bool("access-token") && c.Bool("env") {
		return fmt.Errorf("--access-token and --env are mutually exclusive")
	}

	profile, dir := activeProfile(c)
	cfg, err := loadProfileConfig(dir, profile)
	if err != nil {
		return fmt.Errorf("not logged in (profile %q): %w", profile, err)
	}
	creds, credsPath, err := readFullCredentials(cfg, dir, profile)
	if err != nil {
		return fmt.Errorf("read credentials: %w", err)
	}

	if creds.ExpiresAt != nil && time.Until(*creds.ExpiresAt) <= printCredentialsRefreshThreshold {
		expiry := "expires at " + creds.ExpiresAt.Format(time.RFC3339)
		if time.Until(*creds.ExpiresAt) <= 0 {
			expiry = "expired at " + creds.ExpiresAt.Format(time.RFC3339)
		}
		switch {
		case creds.RefreshToken == "":
			fmt.Fprintf(os.Stderr,
				"warning: token for profile %q %s and no refresh_token is stored; run `ant auth login`\n",
				profile, expiry)
		default:
			baseURL := resolveBaseURL("", cfg)
			clientID := resolveClientID("", cfg)
			tok, err := refreshAccessToken(ctx, baseURL, clientID, creds.RefreshToken)
			if err != nil {
				// Warn-and-proceed rather than error: the on-disk token may
				// still have a few seconds of validity (we refresh at the 120s
				// advisory threshold, not at hard expiry), and a transient 5xx
				// shouldn't turn $(ant auth print-credentials --access-token)
				// into an empty interpolation when a working token exists.
				fmt.Fprintf(os.Stderr,
					"warning: token for profile %q %s and refresh failed: %v\n",
					profile, expiry, err)
				break
			}
			expiresAt := time.Now().Add(time.Duration(tok.ExpiresIn) * time.Second)
			creds.AccessToken = tok.AccessToken
			creds.ExpiresAt = &expiresAt
			if tok.RefreshToken != "" {
				creds.RefreshToken = tok.RefreshToken
			}
			if tok.Scope != "" {
				creds.Scope = tok.Scope
			}
			if tok.Organization.UUID != "" {
				creds.OrganizationUUID = tok.Organization.UUID
				creds.OrganizationName = tok.Organization.Name
			}
			if tok.Account.EmailAddress != "" {
				creds.AccountEmail = tok.Account.EmailAddress
			}
			if err := config.WriteCredentials(credsPath, creds); err != nil {
				return fmt.Errorf("write refreshed credentials: %w", err)
			}
		}
	}

	switch {
	case c.Bool("access-token"):
		fmt.Fprintln(os.Stdout, creds.AccessToken)
	case c.Bool("env"):
		fmt.Fprintf(os.Stdout, "ANTHROPIC_AUTH_TOKEN=%s\n", creds.AccessToken)
		if cfg.BaseURL != "" {
			fmt.Fprintf(os.Stdout, "ANTHROPIC_BASE_URL=%s\n", cfg.BaseURL)
		}
	default:
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(creds)
	}
	return nil
}

// readFullCredentials loads credentials/<profile>.json into the SDK's
// [config.Credentials] (whose MarshalJSON is the canonical wire shape that
// print-credentials emits).
func readFullCredentials(cfg *config.Config, dir, profile string) (config.Credentials, string, error) {
	credsPath := resolveCredentialsPath(cfg, dir, profile)
	data, err := os.ReadFile(credsPath)
	if err != nil {
		return config.Credentials{}, credsPath, err
	}
	var creds config.Credentials
	if err := json.Unmarshal(data, &creds); err != nil {
		return config.Credentials{}, credsPath, fmt.Errorf("parse credentials %s: %w", credsPath, err)
	}
	return creds, credsPath, nil
}

// refreshAccessToken redeems a refresh_token at /v1/oauth/token. Unlike
// [exchangeCode], this sends a JSON body with the oauth-2025-04-20 beta —
// the Python oauth-server requires that beta on the refresh_token grant for
// user_oauth credentials. The federation beta is intentionally absent so the
// request is not gateway-routed to api-go's jwt-bearer-only handler.
func refreshAccessToken(ctx context.Context, baseURL, clientID, refreshToken string) (*tokenResponse, error) {
	body, err := json.Marshal(map[string]string{
		"grant_type":    "refresh_token",
		"refresh_token": refreshToken,
		"client_id":     clientID,
	})
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, "POST", baseURL+"/v1/oauth/token", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("anthropic-beta", betaUserOAuth)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read refresh response: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("token endpoint returned %s (request_id=%s): %s",
			resp.Status, resp.Header.Get("Request-Id"), strings.TrimSpace(string(respBody)))
	}
	var tok tokenResponse
	if err := json.Unmarshal(respBody, &tok); err != nil {
		return nil, fmt.Errorf("parse refresh response: %w", err)
	}
	if tok.AccessToken == "" {
		return nil, fmt.Errorf("token endpoint returned 200 with empty access_token (request_id=%s)",
			resp.Header.Get("Request-Id"))
	}
	return &tok, nil
}

func waitForCallback(ctx context.Context, listener net.Listener, wantState string, timeout time.Duration) (string, error) {
	type result struct {
		code string
		err  error
	}
	ch := make(chan result, 1)
	// Only the first result matters; subsequent handler entries (browser
	// prefetch, favicon, etc.) must not block on a full channel.
	send := func(r result) {
		select {
		case ch <- r:
		default:
		}
	}

	srv := &http.Server{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// The registered redirect_uri is .../callback and the server
			// enforces that, so anything else hitting this listener didn't
			// come from the OAuth flow. state still gates the result; this
			// just avoids treating stray local traffic as a failed auth.
			if r.URL.Path != "/callback" {
				http.NotFound(w, r)
				return
			}
			q := r.URL.Query()
			if e := q.Get("error"); e != "" {
				desc := q.Get("error_description")
				send(result{err: fmt.Errorf("authorization denied: %s: %s", e, desc)})
				writeCallbackPage(w, http.StatusBadRequest, "✗ Authorization failed", e+": "+desc)
				return
			}
			if q.Get("state") != wantState {
				send(result{err: fmt.Errorf("state mismatch (possible CSRF)")})
				writeCallbackPage(w, http.StatusBadRequest, "✗ State mismatch", "The callback state did not match the expected value. This may be a CSRF attempt — do not retry in this browser session.")
				return
			}
			code := q.Get("code")
			if code == "" {
				send(result{err: fmt.Errorf("no code in callback")})
				writeCallbackPage(w, http.StatusBadRequest, "✗ Missing code", "The authorization callback did not include a code parameter.")
				return
			}
			send(result{code: code})
			writeCallbackPage(w, http.StatusOK, "✓ Authorized", "You can close this window and return to the terminal.")
		}),
	}
	go srv.Serve(listener)
	defer func() {
		shutCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		_ = srv.Shutdown(shutCtx)
	}()

	select {
	case r := <-ch:
		return r.code, r.err
	case <-ctx.Done():
		return "", ctx.Err()
	case <-time.After(timeout):
		return "", fmt.Errorf("timed out waiting for authorization")
	}
}

// exchangeCode redeems an authorization_code at /v1/oauth/token. It uses
// form-encoding and no anthropic-beta header so the request lands on the
// Python oauth_server (the established Claude Code path). The
// oidc-federation-2026-04-01 beta header would route to api-go's userauth
// handler, which currently only implements the jwt-bearer grant and rejects
// authorization_code (api-go/core/betagating/interceptor.go).
//
// The Python oauth_server's AuthCodeTokenRequest model requires `state` on
// the token exchange (not just the authorize redirect) as a bound CSRF
// check across both legs — omitting it returns a 400 with
// `oauth_request_parse_error`.
func exchangeCode(ctx context.Context, baseURL, clientID, code, verifier, redirectURI, state string, debug bool) (*tokenResponse, error) {
	form := url.Values{
		"grant_type":    {"authorization_code"},
		"code":          {code},
		"code_verifier": {verifier},
		"client_id":     {clientID},
		"redirect_uri":  {redirectURI},
		"state":         {state},
	}
	req, err := http.NewRequestWithContext(ctx, "POST", baseURL+"/v1/oauth/token",
		strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read token response: %w", err)
	}
	if debug {
		fmt.Fprintf(os.Stderr, "[debug] token exchange POST %s/v1/oauth/token -> %s (request_id=%s)\n",
			baseURL, resp.Status, resp.Header.Get("Request-Id"))
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("token endpoint returned %s (request_id=%s): %s",
			resp.Status, resp.Header.Get("Request-Id"), strings.TrimSpace(string(respBody)))
	}
	var tok tokenResponse
	if err := json.Unmarshal(respBody, &tok); err != nil {
		return nil, fmt.Errorf("parse token response: %w", err)
	}
	if tok.AccessToken == "" {
		return nil, fmt.Errorf("token endpoint returned 200 with empty access_token (request_id=%s)", resp.Header.Get("Request-Id"))
	}
	return &tok, nil
}

// activeProfile resolves the profile to operate on for a command:
// --profile flag → ANTHROPIC_PROFILE env → active_config file → "default".
// Delegates to activeProfileWithSource so the two callers agree on
// resolution order — only the source-description string is dropped here.
// Nil-safe: a nil c falls through to env/active_config resolution (see the
// nil guard in activeProfileWithSource).
func activeProfile(c *cli.Command) (profile, dir string) {
	dir = config.DefaultDir()
	profile, _ = activeProfileWithSource(c, dir)
	return profile, dir
}

// resolveRequestedScope decides which OAuth scope to request on the
// authorize URL: --scope flag → profile's stored UserOAuth.Scope →
// oauthScope default. The stored scope records what the previous login was
// granted, so re-login implicitly re-requests the same set unless the user
// narrows it explicitly.
func resolveRequestedScope(flag string, prev *config.Config) string {
	if flag != "" {
		return flag
	}
	if prev != nil && prev.AuthenticationInfo != nil &&
		prev.AuthenticationInfo.UserOAuth != nil &&
		prev.AuthenticationInfo.UserOAuth.Scope != "" {
		return prev.AuthenticationInfo.UserOAuth.Scope
	}
	return oauthScope
}

// resolveWorkspaceID picks the workspace_id sent on /oauth/authorize:
// --workspace-id flag → profile's stored WorkspaceID → ANTHROPIC_WORKSPACE_ID
// → "" (Console renders the picker step).
//
// The env var is fill-missing only — it does not override the profile's
// stored value, so a stray env can't silently steer the token mint away
// from the profile's intent. Unlike the sibling resolvers there is no
// hardcoded default; empty is a valid return.
func resolveWorkspaceID(flag string, prev *config.Config) string {
	if v := strings.TrimSpace(flag); v != "" {
		return v
	}
	if prev != nil && prev.WorkspaceID != "" {
		return prev.WorkspaceID
	}
	return strings.TrimSpace(os.Getenv("ANTHROPIC_WORKSPACE_ID"))
}

// resolveClientID decides which OAuth client_id to use on the authorize URL
// and token exchange: --client-id flag → profile's stored UserOAuth.ClientID
// → oauthClientIDProd default.
func resolveClientID(flag string, prev *config.Config) string {
	if flag != "" {
		return flag
	}
	if prev != nil && prev.AuthenticationInfo != nil &&
		prev.AuthenticationInfo.UserOAuth != nil &&
		prev.AuthenticationInfo.UserOAuth.ClientID != "" {
		return prev.AuthenticationInfo.UserOAuth.ClientID
	}
	return oauthClientIDProd
}

// resolveConsoleURL decides which base URL to use for the /oauth/authorize
// page: --console-url flag → profile's stored UserOAuth.ConsoleURL →
// defaultConsoleURL. Trailing slashes are stripped so the caller can append
// paths directly.
func resolveConsoleURL(flag string, prev *config.Config) string {
	if flag != "" {
		return strings.TrimRight(flag, "/")
	}
	if prev != nil && prev.AuthenticationInfo != nil &&
		prev.AuthenticationInfo.UserOAuth != nil &&
		prev.AuthenticationInfo.UserOAuth.ConsoleURL != "" {
		return strings.TrimRight(prev.AuthenticationInfo.UserOAuth.ConsoleURL, "/")
	}
	return defaultConsoleURL
}

// resolveBaseURL decides which API base URL the /v1/oauth/token exchange
// posts to: --base-url flag → profile's stored top-level BaseURL →
// defaultBaseURL. The token exchange must hit the same deployment as the
// Console that issued the code, or the client_id won't be recognised.
func resolveBaseURL(flag string, prev *config.Config) string {
	if flag != "" {
		return strings.TrimRight(flag, "/")
	}
	if prev != nil && prev.BaseURL != "" {
		return strings.TrimRight(prev.BaseURL, "/")
	}
	return defaultBaseURL
}

// loadProfileConfig loads configs/<profile>.json via the SDK's
// config.LoadProfile, bypassing ANTHROPIC_PROFILE / active_config resolution
// so a CLI-supplied --profile flag wins. dir is taken explicitly so callers
// that already resolved it (via activeProfile) pass the same value they're
// writing under, keeping read/write paths trivially in sync.
func loadProfileConfig(dir, profile string) (*config.Config, error) {
	return config.LoadProfile(dir, profile)
}

func randomURLSafeString(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func pkceChallengeS256(verifier string) string {
	h := sha256.Sum256([]byte(verifier))
	return base64.RawURLEncoding.EncodeToString(h[:])
}

// writeCallbackPage renders a styled HTML page for the loopback callback
// response, HTML-escaping heading and body so interpolated OAuth error
// strings can't inject markup or script into the user's browser.
func writeCallbackPage(w http.ResponseWriter, status int, heading, body string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	h := html.EscapeString(heading)
	b := html.EscapeString(body)
	fmt.Fprintf(w, `<html><body style="font-family:system-ui;text-align:center;padding:4em"><h2>%s</h2><p>%s</p></body></html>`, h, b)
}

// loginModeEnum tags the three operational modes of `ant auth login`. Mode
// is derived from the (--no-browser, --callback-port) flag combination so
// the dispatch logic is in one place and testable without end-to-end setup.
type loginModeEnum int

const (
	// loginModeDefault: bind an ephemeral (or user-selected) loopback port
	// and auto-open the browser. The normal desktop case.
	loginModeDefault loginModeEnum = iota
	// loginModeSSHForward: --no-browser AND an explicit --callback-port.
	// Bind the port locally but don't open a browser — the user tunnels
	// that port back to their laptop via `ssh -L` and opens the URL there.
	loginModeSSHForward
	// loginModeHeadless: --no-browser with no port. No listener at all;
	// the authorize server renders the code on a public page that the
	// user pastes back into stdin.
	loginModeHeadless
)

// loginMode picks the login flow from the parsed command flags. Extracted
// so the selection rules are unit-testable without standing up a full
// login ceremony.
func loginMode(c *cli.Command) loginModeEnum {
	if !c.Bool("no-browser") {
		return loginModeDefault
	}
	if c.IsSet("callback-port") {
		return loginModeSSHForward
	}
	return loginModeHeadless
}

// startLoopbackListener binds a TCP listener on 127.0.0.1:<port> (port 0
// chooses an ephemeral port) and returns the listener along with the
// redirect URI that should be used in the authorize request. The caller
// is responsible for Close-ing the listener.
func startLoopbackListener(port int) (net.Listener, string, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	if err != nil {
		return nil, "", fmt.Errorf("start callback listener: %w", err)
	}
	actual := listener.Addr().(*net.TCPAddr).Port
	return listener, fmt.Sprintf("http://localhost:%d/callback", actual), nil
}

// manualRedirectURIForAuthorize derives the console-hosted "show the code"
// callback page used for --no-browser login. There is intentionally no
// client-side allowlist of console hosts — the OAuth server validates
// redirect_uri against the client's registered list, so any
// well-formed-but-wrong console URL fails server-side with a clear error.
// Derivation avoids baking non-prod hostnames into the binary.
//
// The ?app=anthropic-cli query param is required: the console landing page
// at /oauth/code/callback uses it to render the right product name and copy
// for the displayed code, and the registered redirect_uri on the
// `anthropic-cli` OAuth client includes the query string verbatim, so an
// authorize call without it is rejected with "redirect_uri not supported by
// client".
func manualRedirectURIForAuthorize(consoleURL string) string {
	consoleURL = strings.TrimRight(consoleURL, "/")
	u, err := url.Parse(consoleURL)
	if err != nil || u.Host == "" {
		return ""
	}
	return consoleURL + "/oauth/code/callback?app=anthropic-cli"
}

// readManualCode prompts on `out` and reads a single line from `in`,
// returning the trimmed authorization code. The read runs in a goroutine so
// the parent select can return on context cancellation; cancellation does
// not actually abort the underlying read, but for this short-lived login
// flow that is acceptable.
//
// Console's /oauth/code/callback page renders the code as "<code>#<state>"
// so the pasted value carries the CSRF state for client-side verification
// (parity with the loopback callback's state check). A bare code with no
// "#" is also accepted — state is still sent on the token exchange and
// validated server-side.
func readManualCode(ctx context.Context, out io.Writer, in io.Reader, wantState string) (string, error) {
	if in == nil {
		in = os.Stdin
	}
	fmt.Fprintln(out, "After authorizing, the page will display a code. Paste it here:")
	fmt.Fprint(out, "Code: ")

	type result struct {
		line string
		err  error
	}
	ch := make(chan result, 1)
	go func() {
		line, err := bufio.NewReader(in).ReadString('\n')
		if err != nil && (err != io.EOF || line == "") {
			ch <- result{err: err}
			return
		}
		ch <- result{line: strings.TrimSpace(line)}
	}()

	select {
	case <-ctx.Done():
		return "", ctx.Err()
	case r := <-ch:
		if r.err != nil {
			return "", fmt.Errorf("read code from stdin: %w", r.err)
		}
		code, pastedState, hasState := strings.Cut(r.line, "#")
		if hasState && pastedState != wantState {
			return "", errors.New("pasted state does not match — code may be from a different login attempt")
		}
		if code == "" {
			return "", errors.New("empty code pasted")
		}
		return code, nil
	}
}

// stdin is a var so tests can pipe a pasted code into the default-mode
// login race.
var stdin io.Reader = os.Stdin

// openBrowser is a var so tests can stub the browser launch when exercising
// the default (non --no-browser) login path.
var openBrowser = func(url string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	default:
		cmd = exec.Command("xdg-open", url)
	}
	if err := cmd.Start(); err != nil {
		return err
	}
	// Reap the child so it doesn't linger as a zombie until process exit.
	go func() { _ = cmd.Wait() }()
	return nil
}
