package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime"
	"net/http"
	"net/http/httputil"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"sync"
	"syscall"

	"github.com/anthropics/anthropic-cli/internal/jsonview"
	"github.com/anthropics/anthropic-sdk-go/config"
	"github.com/anthropics/anthropic-sdk-go/option"

	"github.com/charmbracelet/x/term"
	"github.com/itchyny/json2yaml"
	"github.com/muesli/reflow/wrap"
	"github.com/tidwall/gjson"
	"github.com/tidwall/pretty"
	"github.com/urfave/cli/v3"
)

var OutputFormats = []string{"auto", "explore", "json", "jsonl", "pretty", "raw", "yaml"}

// ValidateBaseURL checks that a base URL is correctly prefixed with a protocol scheme and produces a better
// error message than the person would see otherwise if it doesn't.
func ValidateBaseURL(value, source string) error {
	if value != "" && !strings.HasPrefix(value, "http://") && !strings.HasPrefix(value, "https://") {
		return fmt.Errorf("%s %q is missing a scheme (expected http:// or https://)", source, value)
	}
	return nil
}

// One-shot guards for stderr notices. Reassigned to a fresh sync.Once by
// tests via a reset helper.
var (
	multiAuthWarnOnce     sync.Once
	clientIDDefaultedOnce sync.Once
)

func getDefaultRequestOptions(cmd *cli.Command) []option.RequestOption {
	opts := []option.RequestOption{
		// The CLI walks the 5-tier credential chain itself (below), so opt
		// out of the SDK's environment-based autoloader (anthropic-go#758).
		// Without this, anthropic.NewClient would prepend DefaultClientOptions()
		// — which independently loads a profile via option.WithConfig and emits
		// its own shadow-warning — duplicating both the resolution work and the
		// diagnostic. With this marker, the SDK contributes only the production
		// base-URL default; warnIfMultipleAuthSources below is the sole
		// multi-auth diagnostic.
		option.WithoutEnvironmentDefaults(),
		option.WithHeader("User-Agent", fmt.Sprintf("Anthropic/CLI %s", Version)),
		option.WithHeader("X-Stainless-Lang", "cli"),
		option.WithHeader("X-Stainless-Package-Version", Version),
		option.WithHeader("X-Stainless-Runtime", "cli"),
		option.WithHeader("X-Stainless-CLI-Command", cmd.FullName()),
	}
	// Credential precedence mirrors the WIF User Guide's "Credential resolution" section:
	//   1. --api-key / ANTHROPIC_API_KEY         (flag or env; doc tiers 1+2)
	//   2. --auth-token / ANTHROPIC_AUTH_TOKEN   (flag or env; doc tiers 1+2)
	//   3. profile named by --profile / ANTHROPIC_PROFILE (explicit)
	//   4. ANTHROPIC_FEDERATION_RULE_ID + ANTHROPIC_ORGANIZATION_ID +
	//      ANTHROPIC_IDENTITY_TOKEN[_FILE]
	//   5. profile from active_config → "default" (implicit)
	// The explicit/implicit profile split means: a profile you named beats
	// federation env vars (you asked for it), but federation env vars beat a
	// profile that just happened to be lying around in active_config.
	apiKeySet := cmd.IsSet("api-key")
	authTokenSet := cmd.IsSet("auth-token")
	cfg, profileExplicit := loadProfileIfUsable(cmd)
	fed := federation{
		Assertion:        cmd.String("identity-token"),
		AssertionFile:    cmd.String("identity-token-file"),
		Rule:             cmd.String("federation-rule"),
		OrganizationID:   cmd.String("organization-id"),
		ServiceAccountID: cmd.String("service-account-id"),
	}
	fedAnySet := fed.AnySet()
	warnIfMultipleAuthSources(apiKeySet, authTokenSet, cfg != nil && profileExplicit, fedAnySet, cfg != nil && !profileExplicit)

	useProfile := func() {
		opts = append(opts, option.WithConfigQuiet(cfg))
		if cfg.AuthenticationInfo != nil && cfg.AuthenticationInfo.Type == config.AuthenticationTypeUserOAuth {
			// User-OAuth beta header. WithHeaderAdd appends — WithHeader
			// would overwrite (Header.Set), clobbering any --beta flag the
			// user passed. The SDK middleware appends its own value too;
			// the server selects the one matching the credential type.
			// TODO: drop once the SDK selects the header value per credential type.
			opts = append(opts, option.WithHeaderAdd("anthropic-beta", betaUserOAuth))
		}
	}

	switch {
	case apiKeySet:
		opts = append(opts, option.WithAPIKey(cmd.String("api-key")))
	case authTokenSet:
		opts = append(opts, option.WithAuthToken(cmd.String("auth-token")))
	case cfg != nil && profileExplicit:
		useProfile()
	case fedAnySet:
		opt, err := resolveOAuthOption(fed)
		if err != nil {
			// TODO: fatals on OAuth resolution error via os.Exit, bypassing
			// urfave/cli's error pipeline. Fixing properly requires returning
			// ([]option.RequestOption, error) from this helper and threading
			// through the ~30 codegen-emitted action handlers.
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		if opt != nil {
			opts = append(opts, opt...)
		}
	case cfg != nil:
		useProfile()
	}

	// Override base URL if the --base-url flag is provided. Goes last so it
	// wins over any base URL the loaded config may have set.
	if baseURL := cmd.String("base-url"); baseURL != "" {
		opts = append(opts, option.WithBaseURL(baseURL))
	}

	return opts
}

// warnIfMultipleAuthSources emits a one-shot stderr notice when more than one
// credential source is configured, naming the sources and the precedence
// winner. No secret values are printed. Order matches the User Guide's
// 5-tier precedence (explicit profile beats federation; implicit doesn't).
func warnIfMultipleAuthSources(apiKey, authToken, profileExplicit, federation, profileImplicit bool) {
	type src struct {
		on   bool
		name string
	}
	// The order of this slice MUST match the switch in getDefaultRequestOptions
	// (and credWinner in authStatus) — on[0] is reported as the winner. Both
	// orderings derive from the precedence comment block at the top of
	// getDefaultRequestOptions; if you reorder one, reorder all three.
	all := []src{
		{apiKey, "--api-key / ANTHROPIC_API_KEY"},
		{authToken, "--auth-token / ANTHROPIC_AUTH_TOKEN"},
		{profileExplicit, "profile from --profile / ANTHROPIC_PROFILE"},
		{federation, "federation env"},
		{profileImplicit, "active profile (active_config)"},
	}
	var on []string
	for _, s := range all {
		if s.on {
			on = append(on, s.name)
		}
	}
	if len(on) < 2 {
		return
	}
	multiAuthWarnOnce.Do(func() {
		fmt.Fprintf(os.Stderr,
			"Note: multiple auth sources configured (%s); using %s per precedence. Run `ant auth status` for details.\n",
			strings.Join(on, ", "), on[0])
	})
}

// profileIsExplicit reports whether the profile is named by --profile or
// ANTHROPIC_PROFILE (User Guide tier 3) rather than resolved from
// active_config / "default" (tier 5). --profile is a global flag (extras.go)
// with Sources: ANTHROPIC_PROFILE, so IsSet covers both; the LookupEnv branch
// is a defensive fallback for callers passing a Command not yet Run() (or nil).
func profileIsExplicit(cmd *cli.Command) bool {
	if cmd != nil && cmd.IsSet("profile") {
		return true
	}
	_, ok := os.LookupEnv("ANTHROPIC_PROFILE")
	return ok
}

// loadProfileIfUsable returns the active profile config only when the
// profile is actually usable — i.e. its credentials are present on disk
// for user_oauth profiles. After `ant auth logout` the credentials file
// is deleted but `configs/<profile>.json` is intentionally preserved so
// workspace_id/base_url survive a re-login; this check prevents that
// stale config from claiming the profile tier of credential precedence
// and blocking the fall-through to federation env vars / tokens files.
//
// For oidc_federation profiles the on-disk config is the authoritative
// description of how to mint tokens, so we trust it as-is; the SDK
// handles runtime resolution (identity_token source, rule ID, etc.).
//
// The second return value reports whether the profile was explicitly named
// (--profile / ANTHROPIC_PROFILE) vs implicitly resolved (active_config /
// "default") — these are distinct precedence tiers.
func loadProfileIfUsable(cmd *cli.Command) (*config.Config, bool) {
	explicit := profileIsExplicit(cmd)
	profile, dir := activeProfile(cmd)
	cfg, err := config.LoadProfile(dir, profile)
	if err != nil || cfg == nil || cfg.AuthenticationInfo == nil {
		return nil, explicit
	}
	if cfg.AuthenticationInfo.Type == config.AuthenticationTypeUserOAuth {
		credsPath := cfg.AuthenticationInfo.CredentialsPath
		if credsPath == "" {
			return nil, explicit
		}
		if _, err := os.Stat(credsPath); err != nil {
			return nil, explicit
		}
		// Belt-and-suspenders: profiles written before bootstrap always wrote
		// client_id (or hand-authored ones) may omit it. Fill in the prod
		// default so the SDK's refresh path doesn't fail. New bootstraps
		// always write it (cmd_auth.go), so this is a back-compat shim.
		// config.LoadProfile returns a fresh pointer per call, so mutating
		// the returned struct here doesn't leak across callers.
		if cfg.AuthenticationInfo.UserOAuth != nil && cfg.AuthenticationInfo.UserOAuth.ClientID == "" {
			cfg.AuthenticationInfo.UserOAuth.ClientID = oauthClientIDProd
			clientIDDefaultedOnce.Do(func() {
				fmt.Fprintln(os.Stderr,
					"Note: profile is missing client_id; defaulting to ant-cli prod client. Run `ant auth login` to persist.")
			})
		}
	}
	return cfg, explicit
}

// resolveOAuthOption returns request options for the federation credential
// tier. Routes through the SDK's option.WithFederationTokenProvider so the
// SDK middleware handles the jwt-bearer exchange, beta header, caching, and
// request-context cancellation.
func resolveOAuthOption(fed federation) ([]option.RequestOption, error) {
	if missing := fed.Missing(); len(missing) > 0 {
		return nil, fmt.Errorf("oauth: federation partially configured, missing: %s", strings.Join(missing, ", "))
	}
	idFunc, err := fed.IdentityTokenFunc()
	if err != nil {
		return nil, err
	}
	if idFunc == nil {
		return nil, nil
	}
	return []option.RequestOption{
		option.WithFederationTokenProvider(idFunc, option.FederationOptions{
			FederationRuleID: fed.Rule,
			OrganizationID:   fed.OrganizationID,
			ServiceAccountID: fed.ServiceAccountID,
		}),
	}, nil
}

var debugMiddlewareOption = option.WithMiddleware(
	func(r *http.Request, mn option.MiddlewareNext) (*http.Response, error) {
		logger := log.Default()

		if reqBytes, err := httputil.DumpRequest(r, true); err == nil {
			logger.Printf("Request Content:\n%s\n", reqBytes)
		}

		resp, err := mn(r)
		if err != nil {
			return resp, err
		}

		if respBytes, err := httputil.DumpResponse(resp, true); err == nil {
			logger.Printf("Response Content:\n%s\n", respBytes)
		}

		return resp, err
	},
)

// isInputPiped tries to check for input being piped into the CLI which tells us that we should try to read
// from stdin. This can be a bit tricky in some cases like when an stdin is connected to a pipe but nothing is
// being piped in (this may happen in some environments like Cursor's integration terminal or CI), which is
// why this function is a little more elaborate than it'd be otherwise.
func isInputPiped() bool {
	stat, err := os.Stdin.Stat()
	if err != nil {
		return false
	}

	mode := stat.Mode()

	// Regular file (redirect like < file.txt) — only if non-empty.
	//
	// Notably, on Unix the case like `< /dev/null` is handled below because `/dev/null` is not a regular
	// file. On Windows, NUL appears as a regular file with size 0, so it's also handled correctly.
	if mode.IsRegular() && stat.Size() > 0 {
		return true
	}

	// For pipes/sockets (e.g. `echo foo | stainlesscli`), use an OS-specific check to determine whether
	// data is actually available. Some environments like Cursor's integrated terminal connect stdin as a
	// pipe even when nothing is being piped.
	if mode&(os.ModeNamedPipe|os.ModeSocket) != 0 {
		// Defined in either cmdutil_unix.go or cmdutil_windows.go.
		return isPipedDataAvailableOSSpecific()
	}

	return false
}

func isTerminal(w io.Writer) bool {
	switch v := w.(type) {
	case *os.File:
		return term.IsTerminal(v.Fd())
	default:
		return false
	}
}

func streamOutput(label string, generateOutput func(w *os.File) error) error {
	// For non-tty output (probably a pipe), write directly to stdout
	if !isTerminal(os.Stdout) {
		return streamToStdout(generateOutput)
	}

	// When streaming output on Unix-like systems, there's a special trick involving creating two socket pairs
	// that we prefer because it supports small buffer sizes which results in less pagination per buffer. The
	// constructs needed to run it don't exist on Windows builds, so we have this function broken up into
	// OS-specific files with conditional build comments. Under Windows (and in case our fancy constructs fail
	// on Unix), we fall back to using pipes (`streamToPagerWithPipe`), which are OS agnostic.
	//
	// Defined in either cmdutil_unix.go or cmdutil_windows.go.
	return streamOutputOSSpecific(label, generateOutput)
}

func streamToPagerWithPipe(label string, generateOutput func(w *os.File) error) error {
	r, w, err := os.Pipe()
	if err != nil {
		return err
	}
	defer r.Close()
	defer w.Close()

	pagerProgram := os.Getenv("PAGER")
	if pagerProgram == "" {
		pagerProgram = "less"
	}

	if _, err := exec.LookPath(pagerProgram); err != nil {
		return err
	}

	cmd := exec.Command(pagerProgram)
	cmd.Stdin = r
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = append(os.Environ(),
		"LESS=-X -r -P "+label,
		"MORE=-r -P "+label,
	)

	if err := cmd.Start(); err != nil {
		return err
	}

	if err := r.Close(); err != nil {
		return err
	}

	// If we would be streaming to a terminal and aren't forcing color one way
	// or the other, we should configure things to use color so the pager gets
	// colorized input.
	if isTerminal(os.Stdout) && os.Getenv("FORCE_COLOR") == "" {
		os.Setenv("FORCE_COLOR", "1")
	}

	if err := generateOutput(w); err != nil && !strings.Contains(err.Error(), "broken pipe") {
		return err
	}

	w.Close()
	return cmd.Wait()
}

func streamToStdout(generateOutput func(w *os.File) error) error {
	signal.Ignore(syscall.SIGPIPE)
	err := generateOutput(os.Stdout)
	if err != nil && strings.Contains(err.Error(), "broken pipe") {
		return nil
	}
	return err
}

// writeBinaryResponse writes a binary response to stdout or a file.
//
// Takes in a stdout reference so we can test this function without overriding os.Stdout in tests.
func writeBinaryResponse(response *http.Response, stdout io.Writer, outfile string) (string, error) {
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	switch outfile {
	case "-", "/dev/stdout":
		_, err := stdout.Write(body)
		return "", err
	case "":
		// If output file is unspecified, then print to stdout for plain text or
		// if stdout is not a terminal:
		if !isTerminal(os.Stdout) || isUTF8TextFile(body) {
			_, err := stdout.Write(body)
			return "", err
		}

		// If response has a suggested filename in the content-disposition
		// header, then use that (with an optional suffix to ensure uniqueness):
		file, err := createDownloadFile(response, body)
		if err != nil {
			return "", err
		}
		defer file.Close()
		if _, err := file.Write(body); err != nil {
			return "", err
		}
		return fmt.Sprintf("Wrote output to: %s", file.Name()), nil
	default:
		if err := os.WriteFile(outfile, body, 0644); err != nil {
			return "", err
		}
		return fmt.Sprintf("Wrote output to: %s", outfile), nil
	}
}

// Return a writable file handle to a new file, which attempts to choose a good filename
// based on the Content-Disposition header or sniffing the MIME filetype of the response.
func createDownloadFile(response *http.Response, data []byte) (*os.File, error) {
	filename := "file"
	// If the header provided an output filename, use that
	disp := response.Header.Get("Content-Disposition")
	_, params, err := mime.ParseMediaType(disp)
	if err == nil {
		if dispFilename, ok := params["filename"]; ok {
			// Only use the last path component to prevent directory traversal
			filename = filepath.Base(dispFilename)
			// Try to create the file with exclusive flag to avoid race conditions
			file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644)
			if err == nil {
				return file, nil
			}
		}
	}

	// If file already exists, create a unique filename using CreateTemp
	ext := filepath.Ext(filename)
	if ext == "" {
		ext = guessExtension(data)
	}
	base := strings.TrimSuffix(filename, ext)
	return os.CreateTemp(".", base+"-*"+ext)
}

func guessExtension(data []byte) string {
	ct := http.DetectContentType(data)

	// Prefer common extensions over obscure ones
	switch ct {
	case "application/gzip":
		return ".gz"
	case "application/pdf":
		return ".pdf"
	case "application/zip":
		return ".zip"
	case "audio/mpeg":
		return ".mp3"
	case "image/bmp":
		return ".bmp"
	case "image/gif":
		return ".gif"
	case "image/jpeg":
		return ".jpg"
	case "image/png":
		return ".png"
	case "image/webp":
		return ".webp"
	case "video/mp4":
		return ".mp4"
	}

	exts, err := mime.ExtensionsByType(ct)
	if err == nil && len(exts) > 0 {
		return exts[0]
	} else if isUTF8TextFile(data) {
		return ".txt"
	} else {
		return ".bin"
	}
}

func shouldUseColors(w io.Writer) bool {
	force, ok := os.LookupEnv("FORCE_COLOR")
	if ok {
		if force == "1" {
			return true
		}
		if force == "0" {
			return false
		}
	}
	return isTerminal(w)
}

func formatJSON(res gjson.Result, opts ShowJSONOpts) ([]byte, error) {
	if opts.Transform != "" {
		transformed := res.Get(opts.Transform)
		if transformed.Exists() {
			res = transformed
		}
	}
	// Modeled after `jq -r` (`--raw-output`): if the result is a string, print it without JSON quotes so that
	// it's easier to pipe into other programs.
	if opts.RawOutput && res.Type == gjson.String {
		return []byte(res.Str + "\n"), nil
	}
	switch strings.ToLower(opts.Format) {
	case "auto":
		autoOpts := opts
		autoOpts.Format = "json"
		autoOpts.Transform = ""
		return formatJSON(res, autoOpts)
	case "pretty":
		return []byte(jsonview.RenderJSON(opts.Title, res) + "\n"), nil
	case "json":
		prettyJSON := pretty.Pretty([]byte(res.Raw))
		if shouldUseColors(opts.Stdout) {
			return pretty.Color(prettyJSON, pretty.TerminalStyle), nil
		} else {
			return prettyJSON, nil
		}
	case "jsonl":
		// @ugly is gjson syntax for "no whitespace", so it fits on one line
		oneLineJSON := res.Get("@ugly").Raw
		if shouldUseColors(opts.Stdout) {
			bytes := append(pretty.Color([]byte(oneLineJSON), pretty.TerminalStyle), '\n')
			return bytes, nil
		} else {
			return []byte(oneLineJSON + "\n"), nil
		}
	case "raw":
		return []byte(res.Raw + "\n"), nil
	case "yaml":
		input := strings.NewReader(res.Raw)
		var yaml strings.Builder
		if err := json2yaml.Convert(&yaml, input); err != nil {
			return nil, err
		}
		_, err := opts.Stdout.Write([]byte(yaml.String()))
		return nil, err
	default:
		return nil, fmt.Errorf("Invalid format: %s, valid formats are: %s", opts.Format, strings.Join(OutputFormats, ", "))
	}
}

const warningExploreNotSupported = "Warning: Output format 'explore' not supported for non-terminal output; falling back to 'json'\n"

// ShowJSONOpts configures how JSON output is displayed.
type ShowJSONOpts struct {
	ExplicitFormat bool      // true if the user explicitly passed --format
	Format         string    // output format (auto, explore, json, jsonl, pretty, raw, yaml)
	RawOutput      bool      // like jq -r: print strings without JSON quotes
	Stderr         io.Writer // stderr for warnings; injectable for testing; defaults to os.Stderr
	Stdout         *os.File  // stdout (or pager); injectable for testing; defaults to os.Stdout
	Title          string    // display title
	Transform      string    // GJSON path to extract before displaying
}

func (o *ShowJSONOpts) setDefaults() {
	if o.Stderr == nil {
		o.Stderr = os.Stderr
	}
	if o.Stdout == nil {
		o.Stdout = os.Stdout
	}
}

// ShowJSON displays a single JSON result to the user.
func ShowJSON(res gjson.Result, opts ShowJSONOpts) error {
	opts.setDefaults()

	switch strings.ToLower(opts.Format) {
	case "auto":
		autoOpts := opts
		autoOpts.Format = "json"
		return ShowJSON(res, autoOpts)
	case "explore":
		if !isTerminal(opts.Stdout) {
			if opts.ExplicitFormat {
				fmt.Fprint(opts.Stderr, warningExploreNotSupported)
			}
			jsonOpts := opts
			jsonOpts.Format = "json"
			return ShowJSON(res, jsonOpts)
		}
		if opts.Transform != "" {
			transformed := res.Get(opts.Transform)
			if transformed.Exists() {
				res = transformed
			}
		}
		return jsonview.ExploreJSON(opts.Title, res)
	default:
		bytes, err := formatJSON(res, opts)
		if err != nil {
			return err
		}

		_, err = opts.Stdout.Write(bytes)
		return err
	}
}

// Get the number of lines that would be output by writing the data to the terminal
func countTerminalLines(data []byte, terminalWidth int) int {
	return bytes.Count([]byte(wrap.String(string(data), terminalWidth)), []byte("\n"))
}

type hasRawJSON interface {
	RawJSON() string
}

// ShowJSONIterator displays an iterator of values to the user. Use itemsToDisplay = -1 for no limit.
func ShowJSONIterator[T any](iter jsonview.Iterator[T], itemsToDisplay int64, opts ShowJSONOpts) error {
	opts.setDefaults()

	if opts.Format == "explore" {
		if isTerminal(opts.Stdout) {
			return jsonview.ExploreJSONStream(opts.Title, iter)
		}
		if opts.ExplicitFormat {
			fmt.Fprint(opts.Stderr, warningExploreNotSupported)
		}
		opts.Format = "json"
	}

	terminalWidth, terminalHeight, err := term.GetSize(os.Stdout.Fd())
	if err != nil {
		terminalWidth = 100
		terminalHeight = 40
	}

	// Decide whether or not to use a pager based on whether it's a short output or a long output
	usePager := false
	output := []byte{}
	numberOfNewlines := 0
	// -1 is used to signal no limit of items to display
	for itemsToDisplay != 0 && iter.Next() {
		item := iter.Current()
		var obj gjson.Result
		if hasRaw, ok := any(item).(hasRawJSON); ok {
			obj = gjson.Parse(hasRaw.RawJSON())
		} else {
			jsonData, err := json.Marshal(item)
			if err != nil {
				return err
			}
			obj = gjson.ParseBytes(jsonData)
		}
		json, err := formatJSON(obj, opts)
		if err != nil {
			return err
		}

		output = append(output, json...)
		itemsToDisplay -= 1
		numberOfNewlines += countTerminalLines(json, terminalWidth)

		// If the output won't fit in the terminal window, stream it to a pager
		if numberOfNewlines >= terminalHeight-3 {
			usePager = true
			break
		}
	}

	if !usePager {
		_, err := opts.Stdout.Write(output)
		if err != nil {
			return err
		}

		return iter.Err()
	}

	return streamOutput(opts.Title, func(pager *os.File) error {
		_, err := pager.Write(output)
		if err != nil {
			return err
		}

		pagerOpts := opts
		pagerOpts.Stdout = pager

		for iter.Next() {
			if itemsToDisplay == 0 {
				break
			}
			item := iter.Current()
			var obj gjson.Result
			if hasRaw, ok := any(item).(hasRawJSON); ok {
				obj = gjson.Parse(hasRaw.RawJSON())
			} else {
				jsonData, err := json.Marshal(item)
				if err != nil {
					return err
				}
				obj = gjson.ParseBytes(jsonData)
			}
			if err := ShowJSON(obj, pagerOpts); err != nil {
				return err
			}
			itemsToDisplay -= 1
		}
		return iter.Err()
	})
}
