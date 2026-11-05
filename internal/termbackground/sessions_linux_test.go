package termbackground_test

import (
	"context"
	"os"
	"testing"

	"github.com/anthropics/anthropic-sdk-go/v2"
	"github.com/stretchr/testify/require"

	"github.com/anthropics/anthropic-cli/v2/internal/sessions/live"
	"github.com/anthropics/anthropic-cli/v2/internal/sessions/tui"
)

const runsSessionsView = "runs-sessions-view"

func TestSessionsViewAsksTheTerminalForItsBackground(t *testing.T) {
	written := runOnLightTerminal(t, runsSessionsView)

	require.Contains(t, written, backgroundQuery)
}

// The generated TestMain does not know this program, so it runs here, before TestMain would run the tests again.
// Its context is already done, so the view closes as soon as it opens.
func init() {
	if os.Getenv(programEnvVar) != runsSessionsView {
		return
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_ = tui.Run(ctx, idleSession{}, tui.Options{})
	os.Exit(0)
}

type idleSession struct{}

func (idleSession) Updates() <-chan live.Update { return nil }
func (idleSession) Snapshot() []live.Event      { return nil }
func (idleSession) Session() *anthropic.BetaManagedAgentsSession {
	return &anthropic.BetaManagedAgentsSession{}
}
func (idleSession) SendMessage(context.Context, string) error           { return nil }
func (idleSession) Interrupt(context.Context) error                     { return nil }
func (idleSession) Confirm(context.Context, string, bool, string) error { return nil }
