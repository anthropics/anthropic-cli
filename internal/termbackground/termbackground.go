// Package termbackground decides when the terminal is asked for its background color.
package termbackground

import (
	"os"

	"github.com/charmbracelet/lipgloss"
)

// Bubble Tea's package init asks in every command: it writes a query to the terminal, reads the input waiting there
// and waits up to five seconds for a reply. Lip Gloss skips the query once an answer is set. Go initializes ready
// packages in import path order, so this runs first only in a module whose path sorts before Bubble Tea's.
func init() {
	lipgloss.SetHasDarkBackground(true)
}

// Detect asks the terminal for its background color. Call it before running a Bubble Tea program, because a
// running program reads the terminal's reply itself.
func Detect() {
	lipgloss.SetHasDarkBackground(lipgloss.NewRenderer(os.Stdout).HasDarkBackground())
}
