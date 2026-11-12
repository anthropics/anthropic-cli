package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/anthropics/anthropic-cli/v2/internal/apiquery"
	"github.com/anthropics/anthropic-cli/v2/internal/requestflag"
	"github.com/anthropics/anthropic-cli/v2/internal/sessions/live"
	"github.com/anthropics/anthropic-cli/v2/internal/sessions/tui"
	"github.com/anthropics/anthropic-cli/v2/internal/sessions/web"
	"github.com/anthropics/anthropic-sdk-go/v2"
	"github.com/urfave/cli/v3"
)

var sessionsConnect = cli.Command{
	Name:      "connect",
	Usage:     "Attach an interactive terminal to a session: follow its transcript live, send messages, approve tools",
	ArgsUsage: "<session-id>",
	Suggest:   true,
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "session-id", Usage: "Session to attach to. May instead be given as the first positional argument."},
		&cli.BoolFlag{Name: "verbose", Aliases: []string{"v"}, Usage: "Start with detail expanded (tool bodies, token usage, status events). Toggle with Ctrl-O."},
		&cli.BoolFlag{Name: "web", Usage: "Serve the session viewer on loopback and open it in a browser instead of the terminal view."},
		&cli.BoolFlag{Name: "no-browser", Usage: "With --web, print the URL without opening a browser."},
		&requestflag.Flag[[]string]{
			Name:       "beta",
			Usage:      "Optional header to specify the beta version(s) you want to use.",
			HeaderPath: "anthropic-beta",
		},
	},
	Action:          handleSessionsConnect,
	HideHelpCommand: true,
}

// Registered by lookup rather than by editing the generated tree so a regen
// can't drop it; TestSessionsConnectRegistered catches a renamed group.
func init() {
	for _, group := range Command.Commands {
		if group.Name == "beta:sessions" {
			group.Commands = append(group.Commands, &sessionsConnect)
			return
		}
	}
}

// handleSessionsConnect validates the arguments for the chosen view before
// opening any connection, then hands the session to that view.
func handleSessionsConnect(ctx context.Context, cmd *cli.Command) error {
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("session-id") && len(unusedArgs) > 0 {
		cmd.Set("session-id", unusedArgs[0])
		unusedArgs = unusedArgs[1:]
	}
	if len(unusedArgs) > 0 {
		return fmt.Errorf("Unexpected extra arguments: %v", unusedArgs)
	}
	sessionID := cmd.String("session-id")
	if sessionID == "" {
		return fmt.Errorf("missing session id\nUsage: ant beta:sessions connect <session-id>")
	}
	useWeb := cmd.Bool("web")
	if !useWeb && (!isTerminal(os.Stdout) || !isTerminal(os.Stdin)) {
		return fmt.Errorf("connect needs an interactive terminal (or --web); for scripts use `ant beta:sessions:events stream` and `ant beta:sessions:events send`")
	}
	// Checked before flagOptions so its request logger never attaches.
	if !useWeb && cmd.Bool("debug") {
		return fmt.Errorf("--debug would write over the full-screen view; use --web, or `ant beta:sessions:events stream --debug` to inspect traffic")
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
	client := anthropic.NewClient(append(getDefaultRequestOptions(cmd), options...)...)
	sessionURL := sessionConsoleURL(cmd) + "/workspaces/default/sessions/" + sessionID

	if useWeb {
		return connectWeb(ctx, cmd, client, sessionID, sessionURL)
	}
	return connectTerminal(ctx, cmd, client, sessionID, sessionURL)
}

// connectWeb follows every thread of the session, since the browser viewer
// shows child threads too.
func connectWeb(ctx context.Context, cmd *cli.Command, client anthropic.Client, sessionID, sessionURL string) error {
	threads, err := live.OpenThreads(ctx, client, sessionID)
	if err != nil {
		return err
	}
	defer threads.Close()
	opts := web.Options{
		ConsoleURL: sessionURL,
		Logf:       func(format string, args ...any) { fmt.Fprintf(os.Stderr, format+"\n", args...) },
	}
	if !cmd.Bool("no-browser") {
		opts.Open = openBrowser
	}
	return web.Serve(ctx, threads, opts)
}

// connectTerminal follows the session's primary thread in the full-screen view.
func connectTerminal(ctx context.Context, cmd *cli.Command, client anthropic.Client, sessionID, sessionURL string) error {
	conn, err := live.Open(ctx, client, sessionID)
	if err != nil {
		return err
	}
	defer conn.Close()
	return tui.Run(ctx, conn, tui.Options{Verbose: cmd.Bool("verbose"), ConsoleURL: sessionURL})
}

// sessionConsoleURL is the Console the active profile signed in through, so
// links open on the same deployment the credentials belong to. Unlike
// consoleURL it falls back to the default Console, since a session link is
// always shown.
func sessionConsoleURL(cmd *cli.Command) string {
	cfg, _ := loadProfileIfUsable(cmd)
	return resolveConsoleURL("", cfg)
}
