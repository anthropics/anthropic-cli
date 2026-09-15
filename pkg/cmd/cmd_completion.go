package cmd

import "github.com/anthropics/anthropic-cli/internal/autocomplete"

// The completion command tree (`ant completion <shell>`) is defined in the
// autocomplete package so the goreleaser-driven release pipeline and the
// user-facing CLI share one source of truth for which shells are exposed,
// what the install hints say, and how the renderer is wired.
func init() {
	Command.Commands = append(Command.Commands, autocomplete.BuildCompletionCommand("ant"))
}
