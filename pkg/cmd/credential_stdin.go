package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"
	"sync"

	"github.com/charmbracelet/x/term"
	"github.com/urfave/cli/v3"
)

var stdinCredentialFlags = []struct{ stdinFlag, target string }{
	{"api-key-stdin", "api-key"},
	{"auth-token-stdin", "auth-token"},
	{"aws-api-key-stdin", "aws-api-key"},
	{"foundry-api-key-stdin", "foundry-api-key"},
	{"foundry-auth-token-stdin", "foundry-auth-token"},
}

var directCredentialEnv = map[string]string{
	"api-key":            "ANTHROPIC_API_KEY",
	"auth-token":         "ANTHROPIC_AUTH_TOKEN",
	"aws-api-key":        "ANTHROPIC_AWS_API_KEY",
	"foundry-api-key":    "ANTHROPIC_FOUNDRY_API_KEY",
	"foundry-auth-token": "ANTHROPIC_FOUNDRY_AUTH_TOKEN",
}

// Vars so tests can substitute stdin.
var (
	credentialStdin      = func() io.Reader { return os.Stdin }
	credentialStdinIsTTY = func() bool { return term.IsTerminal(os.Stdin.Fd()) }
	stdinCredentialOnce  sync.Once
	stdinCredentialErr   error
	// stdinCredential is the credential read from stdin; stdin carries at most one.
	stdinCredential struct{ target, value string }
	// stdinConsumedByCredential names the *-stdin flag that drained stdin, so
	// flagOptions can refuse a second stdin consumer with a clear error.
	stdinConsumedByCredential string
)

// applyStdinCredential reads a credential from stdin when --api-key-stdin or
// --auth-token-stdin is set, for directCredential to return. It consumes all
// of stdin: a command using it must take its request body from flags.
func applyStdinCredential(cmd *cli.Command) error {
	stdinCredentialOnce.Do(func() {
		stdinCredential.target, stdinCredential.value, stdinCredentialErr = readStdinCredential(cmd, credentialStdin(), credentialStdinIsTTY())
	})
	return stdinCredentialErr
}

// directCredential returns the api-key or auth-token credential for this run,
// the two that skip profiles and federation: the value piped to
// --<target>-stdin, else the env var. applyStdinCredential must run first.
func directCredential(target string) (string, bool) {
	if stdinCredential.target == target {
		return stdinCredential.value, true
	}
	return os.LookupEnv(directCredentialEnv[target])
}

func readStdinCredential(cmd *cli.Command, in io.Reader, isTTY bool) (target, secret string, err error) {
	var chosen []struct{ stdinFlag, target string }
	for _, f := range stdinCredentialFlags {
		if cmd.Bool(f.stdinFlag) {
			chosen = append(chosen, f)
		}
	}
	switch len(chosen) {
	case 0:
		return "", "", nil
	case 1:
	default:
		return "", "", fmt.Errorf("--%s and --%s are mutually exclusive: stdin can carry only one credential", chosen[0].stdinFlag, chosen[1].stdinFlag)
	}
	f := chosen[0]
	if isTTY {
		return "", "", fmt.Errorf("--%s reads the credential from standard input; pipe it in (e.g. op read op://vault/item/credential | ant --%s ...)", f.stdinFlag, f.stdinFlag)
	}
	stdinConsumedByCredential = "--" + f.stdinFlag
	data, err := io.ReadAll(io.LimitReader(in, 64<<10))
	if err != nil {
		return "", "", fmt.Errorf("--%s: reading stdin: %w", f.stdinFlag, err)
	}
	secret = strings.TrimSpace(string(data))
	if secret == "" {
		return "", "", fmt.Errorf("--%s: stdin was empty", f.stdinFlag)
	}
	return f.target, secret, nil
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
	return directCredentialEnv[target]
}

// credentialEnvVar names the env var an api-key / auth-token value came from,
// or "" when it was piped to --<target>-stdin.
func credentialEnvVar(cmd *cli.Command, target string) string {
	if credentialFromStdin(cmd, target) {
		return ""
	}
	return directCredentialEnv[target]
}
