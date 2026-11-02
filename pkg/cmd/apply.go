package cmd

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/anthropics/anthropic-cli/v2/internal/debugmiddleware"
	"github.com/anthropics/anthropic-cli/v2/internal/declarative/claude"
	"github.com/anthropics/anthropic-cli/v2/internal/declarative/core"
	"github.com/anthropics/anthropic-cli/v2/internal/declarative/render"
	"github.com/anthropics/anthropic-sdk-go/v2"
	"github.com/anthropics/anthropic-sdk-go/v2/config"
	"github.com/anthropics/anthropic-sdk-go/v2/option"
	"github.com/charmbracelet/x/term"
	"github.com/urfave/cli/v3"
)

// applyCommand is `ant apply`. It is hand-written rather than generated from
// the API spec: it reconciles several endpoints instead of wrapping one. The
// generator does not write this file.
var applyCommand = cli.Command{
	Name:     "apply",
	Category: "CONFIGURATION",
	Usage:    "Create, update and remove resources so the API matches your files",
	UsageText: `ant apply [path...]

Reconciles agents, skills, environments, memory stores, vaults and deployments
declared in local files against the Claude Developer Platform, then records what
it did in claude-lock.json.

With no paths, apply reconciles everything already tracked in the lockfile.
Naming a path adds it to the lockfile and reconciles it, along with everything
it references. A directory is walked for resources; a glob is expanded.

Nothing is changed until you approve the plan.`,
	Suggest: true,
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:  "yes",
			Usage: "Apply without asking for confirmation",
		},
		&cli.BoolFlag{
			Name:  "dry-run",
			Usage: "Show the plan and exit without applying",
		},
		&cli.BoolFlag{
			Name:  "force",
			Usage: "Apply even where a resource was changed, archived or deleted outside this config",
		},
		&cli.StringFlag{
			Name:  "lock-file",
			Usage: "Use a specific lockfile instead of discovering claude-lock.json",
		},
		&cli.BoolFlag{
			Name:  "prune",
			Usage: "Remove resources that are in the lockfile but no longer declared on disk",
		},
		&cli.BoolFlag{
			Name:  "upgrade",
			Usage: "Re-resolve skills referenced by URL, picking up new commits on the branch they name",
		},
		&cli.BoolFlag{
			Name:    "verbose",
			Aliases: []string{"v"},
			Usage:   "Show unchanged resources and full field values",
		},
	},
	Action:          handleApply,
	HideHelpCommand: true,
}

// handleApply is the whole of `ant apply`. It loads the lockfile, confirms the
// credentials reach the place that lockfile records, loads the named and
// tracked sources, and plans. Once planned, the plan is always shown. A dry
// run, a blocked plan, a plan with no work or a declined prompt stops there;
// anything else is applied and the lockfile updated as it goes.
func handleApply(ctx context.Context, cmd *cli.Command) (err error) {
	// Errors can name files and lockfile entries from the repository, so
	// escape them before they are printed.
	defer func() {
		if err != nil {
			err = applyError{err}
		}
	}()
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()

	registry := claude.Registry()

	lock, err := loadApplyLockfile(registry, cmd)
	if err != nil {
		return err
	}
	root := lock.Root()

	renderer := &render.Renderer{
		Out:     os.Stdout,
		Color:   shouldUseColors(os.Stdout),
		Verbose: cmd.Bool("verbose"),
	}
	interactive := !cmd.Bool("yes") && !cmd.Bool("dry-run") && canPrompt(isTerminal(os.Stdout), isInputPiped()) == nil

	sdk := anthropic.NewClient(getDefaultRequestOptions(cmd)...)
	origin, err := identifyOrigin(ctx, cmd, sdk, lock)
	if err != nil {
		return err
	}
	renderer.Link = claude.ConsoleLinks(consoleURL(cmd), origin.WorkspaceID)
	if !lock.Existed() {
		renderer.FirstRun(lock.Path, describeOrigin(cmd, origin))
	}
	if err := checkLockOrigin(cmd, lock, origin); err != nil {
		return err
	}
	if lock.Origin == nil {
		// Recorded with the first successful write; a dry run leaves it be.
		lock.Origin = &origin
	}

	cacheDir, removeCacheDir, err := remoteCacheDir()
	if err != nil {
		return err
	}
	defer removeCacheDir()
	fetcher := core.NewGitHubFetcher(cacheDir)

	loader := core.NewLoader(registry, root, fetcher)
	if !cmd.Bool("upgrade") {
		// Without --upgrade a URL naming a branch keeps resolving to the commit
		// the last apply pinned, so a push to that branch doesn't silently
		// rewrite your agents.
		loader.Pin(lock.Pins())
	}

	paths := cmd.Args().Slice()
	if err := loader.Add(ctx, paths); err != nil {
		return err
	}
	// Everything already tracked still has to be reconciled, or applying one
	// file would look like every other resource had been removed.
	if err := loader.AddKeys(ctx, lock.Keys()); err != nil {
		return err
	}
	// With no paths named, offer whatever else under the root looks like a
	// resource — but only to a person at a terminal, who can say no. A
	// script gets exactly what it asked for and nothing it did not.
	if len(paths) == 0 && interactive {
		if err := offerUntracked(ctx, renderer, loader, lock); err != nil {
			return err
		}
	}
	if len(paths) == 0 && !lock.Existed() && len(loader.Sources()) == 0 {
		return fmt.Errorf("no lockfile found at %s and no paths given\n"+
			"  name the files to start tracking, e.g. `ant apply ./agents ./environments`\n"+
			"  (run at a terminal without --yes and apply offers what it finds under this directory)", lock.Path)
	}

	// No sources is still meaningful when the lockfile tracks something: that
	// is how a whole tree gets torn down, with every entry orphaned at once.
	if len(loader.Sources()) == 0 && len(lock.Resources) == 0 {
		return fmt.Errorf("no resources found; nothing to apply")
	}

	client := claude.NewClient(sdk, applyRequestOptions(cmd)...)

	planner := &core.Planner{
		Registry: registry,
		Client:   client,
		Lock:     lock,
		Force:    cmd.Bool("force"),
		Prune:    cmd.Bool("prune"),
	}
	plan, err := planner.Plan(ctx, loader)
	if err != nil {
		return err
	}

	// With nobody to ask, show the field-level detail outright: a dry run and
	// a CI log are read precisely for it. At a prompt it is an answer away
	// instead, so the table can be taken in first.
	prompting := interactive && plan.HasWork() && len(plan.Blocked()) == 0
	renderer.Expanded = !prompting || cmd.Bool("verbose")
	if isTerminal(os.Stdout) {
		renderer.Viewport = func() (int, int) {
			w, h, err := term.GetSize(os.Stdout.Fd())
			if err != nil {
				return 0, 0
			}
			return w, h
		}
	}
	renderer.RenderPlan(plan)

	if cmd.Bool("dry-run") {
		return nil
	}
	if len(plan.Blocked()) > 0 {
		return fmt.Errorf("refusing to apply")
	}
	if !plan.HasWork() {
		return nil
	}

	if !cmd.Bool("yes") {
		ok, err := confirmApply(renderer, plan)
		if err != nil {
			return err
		}
		if !ok {
			fmt.Fprintln(os.Stdout, "Aborted.")
			return nil
		}
	}

	applier := &core.Applier{
		Client: client,
		Lock:   lock,
		Report: renderer.Applied,
	}
	renderer.BeginApply(plan)
	result, err := applier.Apply(ctx, plan)
	if result != nil {
		renderer.RenderResult(result, lock.Path, err)
	}
	return err
}

// applyError prints an error with control characters escaped and continuation
// lines indented, so a newline in a file name can't forge a line of output. It
// unwraps to the original so errors.Is and errors.As still work.
type applyError struct{ err error }

func (e applyError) Error() string {
	return strings.ReplaceAll(render.EscapeLines(e.err.Error()), "\n", "\n    ")
}

func (e applyError) Unwrap() error { return e.err }

// loadApplyLockfile opens the lockfile --lock-file names, or else finds
// claude-lock.json by walking up from the working directory.
func loadApplyLockfile(registry *core.Registry, cmd *cli.Command) (*core.Lockfile, error) {
	if cmd.IsSet("lock-file") {
		path := cmd.String("lock-file")
		if strings.TrimSpace(path) == "" {
			// Falling back to discovery here would quietly point at a
			// different lockfile — or none — and create a second copy of every
			// resource. An unset variable in a CI invocation is exactly how
			// this happens.
			return nil, fmt.Errorf("--lock-file was given an empty value")
		}
		abs, err := filepath.Abs(path)
		if err != nil {
			return nil, err
		}
		return core.LoadLockfile(registry, abs)
	}
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	return core.FindLockfile(registry, wd)
}

// identifyOrigin says where these credentials act. The host is decided
// locally; the organization and workspace are the server's to say, which
// keeps the answer right for API keys and tokens no profile describes.
func identifyOrigin(ctx context.Context, cmd *cli.Command, sdk anthropic.Client, lock *core.Lockfile) (core.Origin, error) {
	origin := core.Origin{BaseURL: currentBaseURL(cmd)}
	if org, ws, err := claude.WhoAmI(ctx, sdk); err == nil {
		origin.OrganizationID, origin.WorkspaceID = org, ws
	} else if o := lock.Origin; o != nil && (o.OrganizationID != "" || o.WorkspaceID != "") {
		// Something is recorded to check against and we cannot check it:
		// stop rather than plan blind.
		return origin, fmt.Errorf("could not confirm these credentials reach the organization and workspace %s tracks: %w", filepath.Base(lock.Path), err)
	} else {
		fmt.Fprintf(os.Stderr, "warning: could not determine the organization and workspace for these credentials: %v\n", err)
	}
	return origin, nil
}

// currentBaseURL resolves the API host by the client's own precedence: the
// Google Cloud or AWS tier's own host when that tier is in use, else flag,
// then environment, then the profile in use, then the SDK default.
func currentBaseURL(cmd *cli.Command) string {
	if googleCloudInUse(cmd) {
		return googleCloudFromRoot(cmd.Root()).host()
	}
	if awsInUse(cmd) {
		return awsFromRoot(cmd.Root()).host()
	}
	if flag := cmd.String("base-url"); flag != "" {
		return flag
	}
	if env := os.Getenv("ANTHROPIC_BASE_URL"); env != "" {
		return env
	}
	if cfg := profileInUse(cmd); cfg != nil && cfg.BaseURL != "" {
		return cfg.BaseURL
	}
	return defaultBaseURL
}

// profileInUse returns the active profile's config only when that profile is
// what requests will authenticate with — by the client's own precedence, an
// API key, auth token, complete federation config, or Google Cloud or AWS
// settings from flags or the environment win over an implicit profile.
func profileInUse(cmd *cli.Command) *config.Config {
	cfg, explicit := loadProfileIfUsable(cmd)
	root := cmd.Root()
	fed := federationFromRoot(root)
	_, apiKeySet := directCredential("api-key")
	_, authTokenSet := directCredential("auth-token")
	switch {
	case cfg == nil, apiKeySet, authTokenSet:
		return nil
	case explicit:
		return cfg
	case fed.AnySet() && len(fed.Missing()) == 0:
		return nil
	case googleCloudFromRoot(root).AnySet(), awsFromRoot(root).AnySet():
		return nil
	}
	return cfg
}

// googleCloudInUse reports whether requests will authenticate through the
// Google Cloud tier, by the client's own precedence.
func googleCloudInUse(cmd *cli.Command) bool {
	root := cmd.Root()
	_, apiKeySet := directCredential("api-key")
	_, authTokenSet := directCredential("auth-token")
	if apiKeySet || authTokenSet || !googleCloudFromRoot(root).AnySet() {
		return false
	}
	if cfg, explicit := loadProfileIfUsable(cmd); cfg != nil && explicit {
		return false
	}
	return !federationFromRoot(root).AnySet()
}

// awsInUse reports whether requests will authenticate through the AWS tier,
// by the client's own precedence.
func awsInUse(cmd *cli.Command) bool {
	root := cmd.Root()
	_, apiKeySet := directCredential("api-key")
	_, authTokenSet := directCredential("auth-token")
	if apiKeySet || authTokenSet || !awsFromRoot(root).AnySet() || googleCloudFromRoot(root).AnySet() {
		return false
	}
	if cfg, explicit := loadProfileIfUsable(cmd); cfg != nil && explicit {
		return false
	}
	return !federationFromRoot(root).AnySet()
}

// consoleURL is the Console the profile in use logged in through, or "" when
// the run does not authenticate through a profile: links are then left off
// rather than guessed.
func consoleURL(cmd *cli.Command) string {
	if cfg := profileInUse(cmd); cfg != nil && cfg.AuthenticationInfo != nil && cfg.AuthenticationInfo.UserOAuth != nil {
		return cfg.AuthenticationInfo.UserOAuth.ConsoleURL
	}
	return ""
}

// describeOrigin puts the resolved origin into words for the first-run
// notice: how the run authenticates, and names for the organization and
// workspace when the profile's stored login knows them.
func describeOrigin(cmd *cli.Command, o core.Origin) render.OriginSummary {
	c := render.OriginSummary{Host: hostOf(o.BaseURL), Organization: o.OrganizationID, Workspace: o.WorkspaceID}
	root := cmd.Root()
	cfg := profileInUse(cmd)
	_, apiKeySet := directCredential("api-key")
	_, authTokenSet := directCredential("auth-token")
	switch {
	case cfg != nil:
		name, dir := activeProfile(cmd)
		c.Credentials = "profile " + name
		if creds, _, err := readCredentials(cfg, dir, name); err == nil {
			if creds.OrganizationUUID == o.OrganizationID {
				c.Organization = nameAndID(creds.OrganizationName, o.OrganizationID)
			}
			if creds.WorkspaceID == o.WorkspaceID {
				c.Workspace = nameAndID(creds.WorkspaceName, o.WorkspaceID)
			}
		}
	case apiKeySet:
		c.Credentials = "API key (" + credentialSourceLabel(root, "api-key") + ")"
	case authTokenSet:
		c.Credentials = "auth token (" + credentialSourceLabel(root, "auth-token") + ")"
	case googleCloudInUse(cmd):
		c.Credentials = "Google Cloud (Application Default Credentials)"
	case awsInUse(cmd) && awsFromRoot(root).APIKey != "":
		c.Credentials = "AWS (API key)"
	case awsInUse(cmd):
		c.Credentials = "AWS (AWS credential chain)"
	default:
		c.Credentials = "workload identity federation"
	}
	if c.Organization == "" {
		c.Organization = "(could not be determined)"
	}
	if c.Workspace == "" {
		c.Workspace = "(could not be determined)"
	}
	return c
}

func hostOf(base string) string {
	if u, err := url.Parse(base); err == nil && u.Host != "" {
		return u.Host
	}
	return base
}

func nameAndID(name, id string) string {
	switch {
	case name != "" && id != "":
		return name + " (" + id + ")"
	case name != "":
		return name
	default:
		return id
	}
}

// checkLockOrigin refuses a lockfile written against one host, organization or
// workspace when run with credentials for another, since every resource would
// read as deleted. The error says which profile would reach the right place.
func checkLockOrigin(cmd *cli.Command, lock *core.Lockfile, origin core.Origin) error {
	err := claude.CheckOrigin(lock, origin)
	if err == nil {
		return nil
	}
	var mismatch *claude.OriginMismatchError
	if errors.As(err, &mismatch) {
		mismatch.Lockfile = render.DisplayPath(mismatch.Lockfile)
	}
	if profile := profileFor(cmd, lock); profile != "" {
		return fmt.Errorf("%w\n  re-run with --profile %s, or use --lock-file to keep a separate lockfile for this target", err, profile)
	}
	return fmt.Errorf("%w\n  switch to credentials for that target, or use --lock-file to keep a separate lockfile for this one", err)
}

// profileFor names a configured profile that would reach the lockfile's
// origin, if there is exactly one obvious candidate.
func profileFor(cmd *cli.Command, lock *core.Lockfile) string {
	_, dir := activeProfile(cmd)
	entries, err := os.ReadDir(filepath.Join(dir, "configs"))
	if err != nil {
		return ""
	}
	var match string
	for _, e := range entries {
		name, ok := strings.CutSuffix(e.Name(), ".json")
		if !ok {
			continue
		}
		cfg, err := config.LoadProfile(dir, name)
		if err != nil || cfg == nil {
			continue
		}
		candidate := core.Origin{BaseURL: cfg.BaseURL, OrganizationID: cfg.OrganizationID, WorkspaceID: cfg.WorkspaceID}
		if candidate.BaseURL == "" {
			candidate.BaseURL = defaultBaseURL
		}
		if claude.CheckOrigin(lock, candidate) == nil {
			if match != "" {
				return "" // ambiguous; say nothing rather than guess
			}
			match = name
		}
	}
	return match
}

// remoteCacheDir is where URL-sourced skills are unpacked: an owner-only
// directory in the user cache, kept between runs, or failing that a new one in
// the temp directory, which the returned function removes.
func remoteCacheDir() (string, func(), error) {
	if base, err := os.UserCacheDir(); err == nil {
		dir := filepath.Join(base, "anthropic", "apply-remote")
		// MkdirAll keeps the mode of a directory that already exists.
		if os.MkdirAll(dir, 0o700) == nil && os.Chmod(dir, 0o700) == nil {
			return dir, func() {}, nil
		}
	}
	// A directory walk skips names that start with a dot.
	dir, err := os.MkdirTemp("", ".anthropic-apply-remote-")
	if err != nil {
		return "", nil, err
	}
	return dir, func() { os.RemoveAll(dir) }, nil
}

// offerUntracked lists resources under the lockfile root that nothing tracks
// yet and asks whether to include them in this apply. Declining changes
// nothing; accepting loads them exactly as if they had been named.
func offerUntracked(ctx context.Context, r *render.Renderer, loader *core.Loader, lock *core.Lockfile) error {
	found, err := loader.Discover(lock.Root())
	if err != nil || len(found) == 0 {
		return err
	}
	r.Untracked(lock.Root(), found, !lock.Existed())
	ok, err := r.Confirm(os.Stdin, fmt.Sprintf("Add %s to this apply?", itOrThese(len(found))))
	if err != nil || !ok {
		return err
	}
	keys := make([]string, len(found))
	for i, f := range found {
		keys[i] = f.Key
	}
	return loader.AddKeys(ctx, keys)
}

func itOrThese(n int) string {
	if n == 1 {
		return "it"
	}
	return fmt.Sprintf("these %d", n)
}

// applyRequestOptions attaches the --debug request logger. Other commands get
// it per request from flagOptions, which apply does not use, so the client
// the reconciler drives has to carry it itself.
func applyRequestOptions(cmd *cli.Command) []option.RequestOption {
	if cmd.Bool("debug") {
		return []option.RequestOption{option.WithMiddleware(debugmiddleware.NewRequestLogger().Middleware())}
	}
	return nil
}

// canPrompt reports whether there is a human to ask. Without one, apply must
// refuse rather than proceed: a script that forgot --yes is asking for a plan,
// not for its resources to be rewritten.
func canPrompt(stdoutIsTerminal, stdinIsPiped bool) error {
	if !stdoutIsTerminal || stdinIsPiped {
		return fmt.Errorf("cannot ask for confirmation without a terminal; re-run with --yes to apply, or --dry-run to see the plan only")
	}
	return nil
}

// confirmApply asks whether to apply the plan, and refuses when there is no
// terminal to ask. When the plan has details to show, "d" redraws it the other
// way and asks again. Any other answer but yes, including an empty line,
// declines.
func confirmApply(r *render.Renderer, plan *core.Plan) (bool, error) {
	if err := canPrompt(isTerminal(os.Stdout), isInputPiped()); err != nil {
		return false, err
	}
	reader := bufio.NewReader(os.Stdin)
	hasDetails := r.HasDetails(plan)
	for {
		r.Prompt(hasDetails)
		line, err := reader.ReadString('\n')
		if err != nil && !errors.Is(err, io.EOF) {
			return false, fmt.Errorf("reading confirmation: %w", err)
		}
		switch strings.ToLower(strings.TrimSpace(line)) {
		case "y", "yes":
			return true, nil
		case "d", "details":
			if !hasDetails {
				return false, nil
			}
			// Redraw the preview the other way: in place when it is all
			// still on screen, beneath otherwise. The one extra row is the
			// answer the terminal just echoed.
			r.Expanded = !r.Expanded
			r.Rewind(1)
			r.RenderPlan(plan)
			continue
		}
		return false, nil
	}
}
