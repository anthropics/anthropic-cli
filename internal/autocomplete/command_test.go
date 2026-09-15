package autocomplete

import (
	"bytes"
	"context"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli/v3"
)

// runCompletion drives the same factory production attaches in cmd.go's init,
// then runs argv against it. The Writer is wired down the whole tree because
// urfave/cli v3's setupDefaults only defaults Writer to os.Stdout on commands
// whose Writer is nil — and it runs per-subcommand, not via inheritance — so
// the buffer has to be set on every node we care about capturing output from.
func runCompletion(t *testing.T, argv ...string) (string, error) {
	t.Helper()
	buf := &bytes.Buffer{}
	app := &cli.Command{
		Name:     "ant",
		Commands: []*cli.Command{BuildCompletionCommand("ant")},
	}
	attachWriter(app, buf)
	err := app.Run(context.Background(), append([]string{"ant"}, argv...))
	return buf.String(), err
}

func attachWriter(c *cli.Command, w io.Writer) {
	c.Writer = w
	for _, sub := range c.Commands {
		attachWriter(sub, w)
	}
}

func TestBuildCompletionCommand_RendersEveryShell(t *testing.T) {
	// One case per shell. wantSubstrs are the function-name + dispatcher
	// patterns that prove __APPNAME__ was substituted in every occurrence,
	// not just the obvious top-level one. The substring `__ant_<shell>` is
	// the form `____APPNAME___<shell>` collapses to after replacement —
	// regressions that only substitute the first match would leave the
	// inner occurrence intact and fail this assertion.
	cases := []struct {
		shell      CompletionStyle
		wantPrefix string
		wantSubstrs []string
	}{
		{CompletionStyleBash, "#!/bin/bash", []string{"__ant_bash_autocomplete", "complete -F __ant_bash_autocomplete ant"}},
		{CompletionStyleZsh, "#compdef ant", []string{"__ant_zsh_autocomplete", "compdef __ant_zsh_autocomplete ant"}},
		{CompletionStyleFish, "#!/usr/bin/env fish", []string{"__ant_fish_autocomplete", "complete -c ant"}},
		{CompletionStylePowershell, "", []string{"__ant_pwsh"}},
	}
	for _, tc := range cases {
		t.Run(string(tc.shell), func(t *testing.T) {
			out, err := runCompletion(t, "completion", string(tc.shell))
			require.NoError(t, err)
			require.NotEmpty(t, out)
			assert.NotContains(t, out, "__APPNAME__", "template token must be substituted")
			if tc.wantPrefix != "" {
				assert.True(t, strings.HasPrefix(out, tc.wantPrefix),
					"want prefix %q, got %q", tc.wantPrefix, out[:min(80, len(out))])
			}
			for _, s := range tc.wantSubstrs {
				assert.Contains(t, out, s)
			}
		})
	}
}

func TestBuildCompletionCommand_BareReturnsUsageError(t *testing.T) {
	out, err := runCompletion(t, "completion")
	require.Error(t, err)
	assert.Empty(t, out, "no script should be written when no shell is specified")
	assert.Contains(t, err.Error(), "ant completion <shell>")
	for _, s := range SupportedShells() {
		assert.Contains(t, err.Error(), string(s))
	}
}

func TestBuildCompletionCommand_UnknownShellReturnsError(t *testing.T) {
	out, err := runCompletion(t, "completion", "tcsh")
	require.Error(t, err)
	assert.Empty(t, out, "no script should be written for an unsupported shell")
	assert.Contains(t, err.Error(), `unsupported shell "tcsh"`)
	for _, s := range SupportedShells() {
		assert.Contains(t, err.Error(), string(s))
	}
}

func TestUsageError(t *testing.T) {
	bare := usageError("ant", "")
	assert.Contains(t, bare, "ant completion <shell>")
	assert.Contains(t, bare, "bash")
	assert.Contains(t, bare, "zsh")
	assert.Contains(t, bare, "fish")
	assert.Contains(t, bare, "pwsh")

	unknown := usageError("ant", "tcsh")
	assert.Contains(t, unknown, `"tcsh"`)
	assert.Contains(t, unknown, "supported:")
	assert.Contains(t, unknown, "bash")
}

func TestRenderCompletionScript_UnsupportedShell(t *testing.T) {
	_, err := RenderCompletionScript("tcsh", "ant")
	require.Error(t, err)
	assert.Contains(t, err.Error(), `unsupported shell "tcsh"`)
}

func TestCompletionDescription_MentionsEveryShell(t *testing.T) {
	desc := completionDescription("ant", SupportedShells())
	for _, s := range SupportedShells() {
		assert.Contains(t, desc, string(s), "description should mention %s", s)
	}
	// install hints must use the app name we passed in
	assert.Contains(t, desc, "ant completion bash")
}
