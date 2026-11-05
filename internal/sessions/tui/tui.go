// Package tui is the terminal front end of `ant beta:sessions connect`: a
// scrolling transcript above a composer, with a status bar on top and key
// hints below. Events arrive from internal/sessions/live and are folded into
// transcript nodes by internal/sessions/timeline.
//
// Files: model.go owns state and the bubbletea loop, input.go the keys and
// the composer/approval strip, transcript.go the scrolling body, chrome.go the
// fixed rows around it, and markdown.go, format.go and styles.go the helpers
// they share.
package tui

import (
	"context"
	"errors"

	"github.com/anthropics/anthropic-sdk-go/v2"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/anthropics/anthropic-cli/v2/internal/sessions/live"
	"github.com/anthropics/anthropic-cli/v2/internal/termbackground"
)

// Run takes over the terminal until the user detaches or ctx ends.
func Run(ctx context.Context, session Session, opts Options) error {
	termbackground.Detect()
	program := tea.NewProgram(newModel(ctx, session, opts), tea.WithAltScreen(), tea.WithMouseCellMotion(), tea.WithContext(ctx))
	_, err := program.Run()
	if errors.Is(err, tea.ErrProgramKilled) {
		return nil // ctx cancelled or interrupted: a detach, not a failure
	}
	return err
}

// Options configures what the TUI shows.
type Options struct {
	// Verbose starts with tool bodies and token usage expanded (Ctrl-O toggles it).
	Verbose bool
	// ConsoleURL is the session's page in the web Console, shown in the footer.
	ConsoleURL string
}

// Session is what the TUI needs from an attached session; *live.Conn is one.
// The three senders block on the network and are only called from tea.Cmds.
type Session interface {
	Updates() <-chan live.Update
	Snapshot() []live.Event
	Session() *anthropic.BetaManagedAgentsSession
	SendMessage(ctx context.Context, text string) error
	Interrupt(ctx context.Context) error
	Confirm(ctx context.Context, toolUseID string, allow bool, denyMessage string) error
}

var _ Session = (*live.Conn)(nil)
