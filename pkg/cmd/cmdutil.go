package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"

	"github.com/anthropics/anthropic-cli/internal/jsonview"
	"github.com/anthropics/anthropic-sdk-go/option"

	"github.com/charmbracelet/x/term"
	"github.com/itchyny/json2yaml"
	"github.com/muesli/reflow/wrap"
	"github.com/tidwall/gjson"
	"github.com/tidwall/pretty"
	"github.com/urfave/cli/v3"
)

var OutputFormats = []string{"auto", "explore", "json", "jsonl", "pretty", "raw", "yaml"}

func getDefaultRequestOptions(cmd *cli.Command) []option.RequestOption {
	opts := []option.RequestOption{
		option.WithHeader("User-Agent", fmt.Sprintf("Anthropic/CLI %s", Version)),
		option.WithHeader("X-Stainless-Lang", "cli"),
		option.WithHeader("X-Stainless-Package-Version", Version),
		option.WithHeader("X-Stainless-Runtime", "cli"),
		option.WithHeader("X-Stainless-CLI-Command", cmd.FullName()),
	}

	// Override base URL if the --base-url flag is provided
	if baseURL := cmd.String("base-url"); baseURL != "" {
		opts = append(opts, option.WithBaseURL(baseURL))
	}

	return opts
}

var debugMiddlewareOption = option.WithMiddleware(
	func(r *http.Request, mn option.MiddlewareNext) (*http.Response, error) {
		logger := log.Default()

		if reqBytes, err := httputil.DumpRequest(r, true); err == nil {
			logger.Printf("Request Content:\n%s\n", reqBytes)
		}

		resp, err := mn(r)
		if err != nil {
			return resp, err
		}

		if respBytes, err := httputil.DumpResponse(resp, true); err == nil {
			logger.Printf("Response Content:\n%s\n", respBytes)
		}

		return resp, err
	},
)

func isInputPiped() bool {
	stat, _ := os.Stdin.Stat()
	return (stat.Mode() & os.ModeCharDevice) == 0
}

func isTerminal(w io.Writer) bool {
	switch v := w.(type) {
	case *os.File:
		return term.IsTerminal(v.Fd())
	default:
		return false
	}
}

func streamOutput(label string, generateOutput func(w *os.File) error) error {
	// For non-tty output (probably a pipe), write directly to stdout
	if !isTerminal(os.Stdout) {
		return streamToStdout(generateOutput)
	}

	// When streaming output on Unix-like systems, there's a special trick involving creating two socket pairs
	// that we prefer because it supports small buffer sizes which results in less pagination per buffer. The
	// constructs needed to run it don't exist on Windows builds, so we have this function broken up into
	// OS-specific files with conditional build comments. Under Windows (and in case our fancy constructs fail
	// on Unix), we fall back to using pipes (`streamToPagerWithPipe`), which are OS agnostic.
	//
	// Defined in either cmdutil_unix.go or cmdutil_windows.go.
	return streamOutputOSSpecific(label, generateOutput)
}

func streamToPagerWithPipe(label string, generateOutput func(w *os.File) error) error {
	r, w, err := os.Pipe()
	if err != nil {
		return err
	}
	defer r.Close()
	defer w.Close()

	pagerProgram := os.Getenv("PAGER")
	if pagerProgram == "" {
		pagerProgram = "less"
	}

	if _, err := exec.LookPath(pagerProgram); err != nil {
		return err
	}

	cmd := exec.Command(pagerProgram)
	cmd.Stdin = r
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = append(os.Environ(),
		"LESS=-X -r -P "+label,
		"MORE=-r -P "+label,
	)

	if err := cmd.Start(); err != nil {
		return err
	}

	if err := r.Close(); err != nil {
		return err
	}

	// If we would be streaming to a terminal and aren't forcing color one way
	// or the other, we should configure things to use color so the pager gets
	// colorized input.
	if isTerminal(os.Stdout) && os.Getenv("FORCE_COLOR") == "" {
		os.Setenv("FORCE_COLOR", "1")
	}

	if err := generateOutput(w); err != nil && !strings.Contains(err.Error(), "broken pipe") {
		return err
	}

	w.Close()
	return cmd.Wait()
}

func streamToStdout(generateOutput func(w *os.File) error) error {
	signal.Ignore(syscall.SIGPIPE)
	err := generateOutput(os.Stdout)
	if err != nil && strings.Contains(err.Error(), "broken pipe") {
		return nil
	}
	return err
}

func shouldUseColors(w io.Writer) bool {
	force, ok := os.LookupEnv("FORCE_COLOR")
	if ok {
		if force == "1" {
			return true
		}
		if force == "0" {
			return false
		}
	}
	return isTerminal(w)
}

func formatJSON(expectedOutput *os.File, title string, res gjson.Result, format string, transform string) ([]byte, error) {
	if format != "raw" && transform != "" {
		transformed := res.Get(transform)
		if transformed.Exists() {
			res = transformed
		}
	}
	switch strings.ToLower(format) {
	case "auto":
		return formatJSON(expectedOutput, title, res, "json", "")
	case "pretty":
		return []byte(jsonview.RenderJSON(title, res) + "\n"), nil
	case "json":
		prettyJSON := pretty.Pretty([]byte(res.Raw))
		if shouldUseColors(expectedOutput) {
			return pretty.Color(prettyJSON, pretty.TerminalStyle), nil
		} else {
			return prettyJSON, nil
		}
	case "jsonl":
		// @ugly is gjson syntax for "no whitespace", so it fits on one line
		oneLineJSON := res.Get("@ugly").Raw
		if shouldUseColors(expectedOutput) {
			bytes := append(pretty.Color([]byte(oneLineJSON), pretty.TerminalStyle), '\n')
			return bytes, nil
		} else {
			return []byte(oneLineJSON + "\n"), nil
		}
	case "raw":
		return []byte(res.Raw + "\n"), nil
	case "yaml":
		input := strings.NewReader(res.Raw)
		var yaml strings.Builder
		if err := json2yaml.Convert(&yaml, input); err != nil {
			return nil, err
		}
		_, err := expectedOutput.Write([]byte(yaml.String()))
		return nil, err
	default:
		return nil, fmt.Errorf("Invalid format: %s, valid formats are: %s", format, strings.Join(OutputFormats, ", "))
	}
}

// Display JSON to the user in various different formats
func ShowJSON(out *os.File, title string, res gjson.Result, format string, transform string) error {
	if format != "raw" && transform != "" {
		transformed := res.Get(transform)
		if transformed.Exists() {
			res = transformed
		}
	}

	switch strings.ToLower(format) {
	case "auto":
		return ShowJSON(out, title, res, "json", "")
	case "explore":
		return jsonview.ExploreJSON(title, res)
	default:
		bytes, err := formatJSON(out, title, res, format, transform)
		if err != nil {
			return err
		}

		_, err = out.Write(bytes)
		return err
	}
}

// Get the number of lines that would be output by writing the data to the terminal
func countTerminalLines(data []byte, terminalWidth int) int {
	return bytes.Count([]byte(wrap.String(string(data), terminalWidth)), []byte("\n"))
}

type HasRawJSON interface {
	RawJSON() string
}

// For an iterator over different value types, display its values to the user in
// different formats.
func ShowJSONIterator[T any](stdout *os.File, title string, iter jsonview.Iterator[T], format string, transform string) error {
	if format == "explore" {
		return jsonview.ExploreJSONStream(title, iter)
	}

	terminalWidth, terminalHeight, err := term.GetSize(os.Stdout.Fd())
	if err != nil {
		terminalWidth = 100
		terminalHeight = 40
	}

	// Decide whether or not to use a pager based on whether it's a short output or a long output
	usePager := false
	output := []byte{}
	numberOfNewlines := 0
	for iter.Next() {
		item := iter.Current()
		var obj gjson.Result
		if hasRaw, ok := any(item).(HasRawJSON); ok {
			obj = gjson.Parse(hasRaw.RawJSON())
		} else {
			jsonData, err := json.Marshal(item)
			if err != nil {
				return err
			}
			obj = gjson.ParseBytes(jsonData)
		}
		json, err := formatJSON(stdout, title, obj, format, transform)
		if err != nil {
			return err
		}

		output = append(output, json...)
		numberOfNewlines += countTerminalLines(json, terminalWidth)

		// If the output won't fit in the terminal window, stream it to a pager
		if numberOfNewlines >= terminalHeight-3 {
			usePager = true
			break
		}
	}

	if !usePager {
		_, err := stdout.Write(output)
		if err != nil {
			return err
		}

		return iter.Err()
	}

	return streamOutput(title, func(pager *os.File) error {
		// Write the output we used during the initial terminal size computation
		_, err := pager.Write(output)
		if err != nil {
			return err
		}

		for iter.Next() {
			item := iter.Current()
			var obj gjson.Result
			if hasRaw, ok := any(item).(HasRawJSON); ok {
				obj = gjson.Parse(hasRaw.RawJSON())
			} else {
				jsonData, err := json.Marshal(item)
				if err != nil {
					return err
				}
				obj = gjson.ParseBytes(jsonData)
			}
			if err := ShowJSON(pager, title, obj, format, transform); err != nil {
				return err
			}
		}
		return iter.Err()
	})
}
