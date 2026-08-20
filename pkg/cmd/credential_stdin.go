package cmd

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"

	"github.com/charmbracelet/x/term"
	"github.com/urfave/cli/v3"
)

// Credentials passed as argv are visible in shell history and process
// listings; --api-key / --auth-token are deprecated in favor of the env vars,
// the *-stdin flags below, or `ant auth login`. Values are the replacement
// each deprecation notice points at.
var argvCredentialFlags = map[string]string{
	"api-key":    "ANTHROPIC_API_KEY, --api-key-stdin, or 'ant auth login'",
	"auth-token": "ANTHROPIC_AUTH_TOKEN or --auth-token-stdin",
}

// No backticks in flag usage text: urfave/cli renders a backticked span as
// the flag's value placeholder.
const argvCredentialDeprecation = "Deprecated: a credential on the command line is visible in shell history and process listings. Use %s instead."

var stdinCredentialFlags = []struct{ stdinFlag, target string }{
	{"api-key-stdin", "api-key"},
	{"auth-token-stdin", "auth-token"},
}

// Vars so tests can substitute argv and stdin.
var (
	commandLineArgs        = func() []string { return os.Args[1:] }
	credentialStdin        = func() io.Reader { return os.Stdin }
	credentialStdinIsTTY   = func() bool { return term.IsTerminal(os.Stdin.Fd()) }
	argvCredentialWarnOnce sync.Once
	stdinCredentialOnce    sync.Once
	stdinCredentialErr     error
	// stdinConsumedByCredential names the *-stdin flag that drained stdin, so
	// flagOptions can refuse a second stdin consumer with a clear error.
	stdinConsumedByCredential string
)

// credentialFlagsOnCommandLine returns the credential flag names (without
// dashes) that appear literally in argv. urfave/cli can't tell a flag value
// that came from argv from one that came from its env source, so argv is
// inspected directly; scanning stops at "--".
func credentialFlagsOnCommandLine() []string {
	var found []string
	for _, arg := range commandLineArgs() {
		if arg == "--" {
			break
		}
		name := strings.TrimLeft(arg, "-")
		if name == arg {
			continue
		}
		name, _, _ = strings.Cut(name, "=")
		if _, ok := argvCredentialFlags[name]; ok {
			found = append(found, name)
		}
	}
	return found
}

func warnIfCredentialOnCommandLine(w io.Writer) {
	names := credentialFlagsOnCommandLine()
	if len(names) == 0 {
		return
	}
	argvCredentialWarnOnce.Do(func() {
		for _, name := range names {
			fmt.Fprintf(w, "Warning: --%s is deprecated and will be removed in a future release: "+
				"a credential on the command line is visible in shell history and process listings. Use %s instead.\n",
				name, argvCredentialFlags[name])
		}
	})
}

// applyStdinCredential reads a credential from stdin when --api-key-stdin or
// --auth-token-stdin is set and stores it on the corresponding credential
// flag, so the precedence chain in getDefaultRequestOptions (and `auth
// status`) treats it exactly like --api-key / --auth-token. It consumes all
// of stdin: a command using it must take its request body from flags.
func applyStdinCredential(cmd *cli.Command) error {
	stdinCredentialOnce.Do(func() {
		stdinCredentialErr = readStdinCredential(cmd, credentialStdin(), credentialStdinIsTTY(), credentialFlagsOnCommandLine())
	})
	return stdinCredentialErr
}

func readStdinCredential(cmd *cli.Command, in io.Reader, isTTY bool, onCommandLine []string) error {
	var chosen []struct{ stdinFlag, target string }
	for _, f := range stdinCredentialFlags {
		if cmd.Bool(f.stdinFlag) {
			chosen = append(chosen, f)
		}
	}
	switch len(chosen) {
	case 0:
		return nil
	case 1:
	default:
		return errors.New("--api-key-stdin and --auth-token-stdin are mutually exclusive: stdin can carry only one credential")
	}
	f := chosen[0]
	for _, name := range onCommandLine {
		if name == f.target {
			return fmt.Errorf("--%s cannot be combined with --%s", f.stdinFlag, f.target)
		}
	}
	if isTTY {
		return fmt.Errorf("--%s reads the credential from standard input; pipe it in (e.g. op read op://vault/item/credential | ant --%s ...)", f.stdinFlag, f.stdinFlag)
	}
	stdinConsumedByCredential = "--" + f.stdinFlag
	data, err := io.ReadAll(io.LimitReader(in, 64<<10))
	if err != nil {
		return fmt.Errorf("--%s: reading stdin: %w", f.stdinFlag, err)
	}
	secret := strings.TrimSpace(string(data))
	if secret == "" {
		return fmt.Errorf("--%s: stdin was empty", f.stdinFlag)
	}
	return cmd.Set(f.target, secret)
}

func credentialFromStdin(cmd *cli.Command, target string) bool {
	for _, f := range stdinCredentialFlags {
		if f.target == target && cmd.Bool(f.stdinFlag) {
			return true
		}
	}
	return false
}

// credentialSourceLabel names where an api-key / auth-token value came from
// for diagnostics, without printing the value.
func credentialSourceLabel(cmd *cli.Command, target string) string {
	if credentialFromStdin(cmd, target) {
		return "--" + target + "-stdin"
	}
	env := "ANTHROPIC_API_KEY"
	if target == "auth-token" {
		env = "ANTHROPIC_AUTH_TOKEN"
	}
	return fmt.Sprintf("--%s / %s", target, env)
}
