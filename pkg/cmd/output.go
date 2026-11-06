package cmd

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/anthropics/anthropic-cli/v2/internal/jsonview"
	"github.com/anthropics/anthropic-cli/v2/internal/termsafe"
	"github.com/tidwall/gjson"
)

// rendersToTerminal is true for the pipe that feeds the pager too, because the pager only runs on a terminal.
func rendersToTerminal(w io.Writer) bool {
	return isTerminal(w) || isTerminal(os.Stdout)
}

// displayedJSON returns the result to render: res after any transform, and with its keys and strings safe to print
// unless the output is piped. It runs after the transform because a transform can decode a string into new JSON.
func displayedJSON(res gjson.Result, opts ShowJSONOpts) gjson.Result {
	if opts.Transform != "" {
		transformed := res.Get(opts.Transform)
		if transformed.Exists() {
			res = transformed
		}
	}
	// Forced color says that piped output is headed for a display, such as `less -R`, and pretty is only for reading.
	if rendersToTerminal(opts.Stdout) || shouldUseColors(opts.Stdout) || strings.EqualFold(opts.Format, "pretty") {
		return gjson.Parse(termsafe.JSON(res.Raw))
	}
	return res
}

// safeIterator gives the explorer, which only runs on a terminal, items that are safe to print.
type safeIterator[T any] struct {
	jsonview.Iterator[T]
	current safeItem
	err     error
}

func (it *safeIterator[T]) Next() bool {
	if !it.Iterator.Next() {
		return false
	}
	res, err := itemJSON(it.Iterator.Current())
	it.current, it.err = safeItem(termsafe.JSON(res.Raw)), err
	return err == nil
}

func (it *safeIterator[T]) Current() safeItem {
	return it.current
}

func (it *safeIterator[T]) Err() error {
	return errors.Join(it.err, it.Iterator.Err())
}

type safeItem string

func (i safeItem) RawJSON() string {
	return string(i)
}

func writeBody(stdout io.Writer, body []byte) error {
	if rendersToTerminal(stdout) {
		stdout = termsafe.Writer{W: stdout}
	}
	_, err := stdout.Write(body)
	return err
}

// downloadFilename returns the name to save a download under, from a filename the server sent. The name gets printed
// as well as created on disk. It keeps only the last path component to prevent directory traversal, and escapes first
// because the backslashes that adds are path separators on Windows.
func downloadFilename(name string) string {
	return filepath.Base(termsafe.Escape(name))
}
