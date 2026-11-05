package termbackground_test

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"testing"

	"github.com/charmbracelet/lipgloss"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
	"golang.org/x/sys/unix"

	"github.com/anthropics/anthropic-cli/v2/internal/jsonview"
	"github.com/anthropics/anthropic-cli/v2/internal/termbackground"
)

const (
	thisPackagePath = "github.com/anthropics/anthropic-cli/v2/internal/termbackground"
	bubbleTeaPath   = "github.com/charmbracelet/bubbletea"

	programEnvVar      = "TERMBACKGROUND_TEST_PROGRAM"
	onlyStartsUp       = "only-starts-up"
	callsDetect        = "calls-detect"
	exploresJSON       = "explores-json"
	exploresJSONStream = "explores-json-stream"

	backgroundQuery      = "\x1b]11;?\x1b\\\x1b[6n"
	lightBackgroundReply = "\x1b]11;rgb:ffff/ffff/ffff\x1b\\\x1b[1;1R"
	quitHelp             = "quit"
	quitKey              = "q"
)

// TestMain also stands in for the programs that the tests run on a terminal of their own.
func TestMain(m *testing.M) {
	switch os.Getenv(programEnvVar) {
	case onlyStartsUp:
	case callsDetect:
		termbackground.Detect()
		fmt.Printf("dark background: %t", lipgloss.HasDarkBackground())
	case exploresJSON:
		_ = jsonview.ExploreJSON("title", gjson.Parse(`{"key": "value"}`))
	case exploresJSONStream:
		_ = jsonview.ExploreJSONStream("title", &noItems{})
	default:
		os.Exit(m.Run())
	}
}

type noItems struct{}

func (*noItems) Next() bool   { return false }
func (*noItems) Err() error   { return nil }
func (*noItems) Current() any { return nil }

func TestStartingUpAsksTheTerminalNothing(t *testing.T) {
	skipWhereBubbleTeaAsksAtStartUp(t)

	written := runOnLightTerminal(t, onlyStartsUp)

	require.NotContains(t, written, "\x1b")
}

func TestExploreViewAsksTheTerminalForItsBackground(t *testing.T) {
	skipWhereBubbleTeaAsksAtStartUp(t)

	for _, program := range []string{exploresJSON, exploresJSONStream} {
		t.Run(program, func(t *testing.T) {
			written := runOnLightTerminal(t, program)

			require.Contains(t, written, backgroundQuery)
		})
	}
}

func TestDetectLearnsTheBackgroundFromTheTerminal(t *testing.T) {
	written := runOnLightTerminal(t, callsDetect)

	require.Contains(t, written, "dark background: false")
}

func skipWhereBubbleTeaAsksAtStartUp(t *testing.T) {
	t.Helper()

	if thisPackagePath > bubbleTeaPath {
		t.Skip("Go initializes Bubble Tea before this package, whose import path sorts after Bubble Tea's")
	}
}

// runOnLightTerminal runs one of TestMain's programs in the foreground of a new terminal, which answers a background
// query as a white terminal does and quits the explore view once it is drawn. It returns what the program wrote.
func runOnLightTerminal(t *testing.T, program string) string {
	t.Helper()

	terminal, err := os.OpenFile("/dev/ptmx", os.O_RDWR, 0)
	if err != nil {
		t.Skipf("no pty available: %v", err)
	}
	defer terminal.Close()
	require.NoError(t, unix.IoctlSetPointerInt(int(terminal.Fd()), unix.TIOCSPTLCK, 0))
	require.NoError(t, unix.IoctlSetWinsize(int(terminal.Fd()), unix.TIOCSWINSZ, &unix.Winsize{Row: 24, Col: 80}))
	number, err := unix.IoctlGetInt(int(terminal.Fd()), unix.TIOCGPTN)
	require.NoError(t, err)
	programSide, err := os.OpenFile(fmt.Sprintf("/dev/pts/%d", number), os.O_RDWR|syscall.O_NOCTTY, 0)
	require.NoError(t, err)

	cmd := exec.Command(os.Args[0])
	cmd.Env = append(os.Environ(), programEnvVar+"="+program, "TERM=xterm", "CI=")
	cmd.Stdin, cmd.Stdout, cmd.Stderr = programSide, programSide, programSide
	cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true, Setctty: true}
	require.NoError(t, cmd.Start())
	require.NoError(t, programSide.Close())

	var written strings.Builder
	chunk := make([]byte, 1024)
	pressedQuit := false
	for {
		n, err := terminal.Read(chunk)
		written.Write(chunk[:n])
		if err != nil {
			break
		}
		if strings.HasSuffix(written.String(), backgroundQuery) {
			_, err := terminal.WriteString(lightBackgroundReply)
			require.NoError(t, err)
		}
		if !pressedQuit && strings.Contains(written.String(), quitHelp) {
			_, err := terminal.WriteString(quitKey)
			require.NoError(t, err)
			pressedQuit = true
		}
	}
	require.NoError(t, cmd.Wait())
	return written.String()
}
