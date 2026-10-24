package cmd

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"

	"github.com/anthropics/anthropic-cli/v2/internal/jsonview"
)

func TestStreamOutput(t *testing.T) {
	t.Setenv("PAGER", "cat")
	err := streamOutput("stream test", func(w *os.File) error {
		_, writeErr := w.WriteString("Hello world\n")
		return writeErr
	})
	if err != nil {
		t.Errorf("streamOutput failed: %v", err)
	}
}

func TestGetPagerArgs(t *testing.T) {
	t.Run("defaults to less", func(t *testing.T) {
		t.Setenv("PAGER", "")

		args, err := getPagerArgs()
		require.NoError(t, err)
		assert.Equal(t, []string{"less"}, args)
	})

	t.Run("accepts a pager command", func(t *testing.T) {
		t.Setenv("PAGER", "less")

		args, err := getPagerArgs()
		require.NoError(t, err)
		assert.Equal(t, []string{"less"}, args)
	})

	t.Run("splits command and arguments", func(t *testing.T) {
		t.Setenv("PAGER", "cat -n")

		args, err := getPagerArgs()
		require.NoError(t, err)
		assert.Equal(t, []string{"cat", "-n"}, args)
	})

	t.Run("rejects whitespace-only values", func(t *testing.T) {
		t.Setenv("PAGER", " \t")

		_, err := getPagerArgs()
		require.Error(t, err)
		assert.EqualError(t, err, "PAGER must contain a command")
	})
}

func TestStreamToPagerWithPipeSupportsArguments(t *testing.T) {
	t.Setenv("PAGER", "cat -n")

	err := streamToPagerWithPipe("stream test", func(w *os.File) error {
		_, err := w.WriteString("Hello world\n")
		return err
	})
	require.NoError(t, err)
}

func TestStreamToPagerWithPipeRejectsMissingPager(t *testing.T) {
	t.Setenv("PAGER", filepath.Join(t.TempDir(), "missing-pager")+" --flag")

	called := false
	err := streamToPagerWithPipe("stream test", func(w *os.File) error {
		called = true
		_, err := w.WriteString("Hello world\n")
		return err
	})
	require.Error(t, err)
	assert.False(t, called)
}

func TestWriteBinaryResponse(t *testing.T) {
	t.Run("write to explicit file", func(t *testing.T) {
		tmpDir := t.TempDir()
		outfile := tmpDir + "/output.txt"
		body := []byte("test content")
		resp := &http.Response{
			Body: io.NopCloser(bytes.NewReader(body)),
		}

		msg, err := writeBinaryResponse(resp, os.Stdout, outfile)

		require.NoError(t, err)
		assert.Contains(t, msg, outfile)

		content, err := os.ReadFile(outfile)
		require.NoError(t, err)
		assert.Equal(t, body, content)
	})

	t.Run("download message escapes the server's filename", func(t *testing.T) {
		t.Chdir(t.TempDir())
		stdoutIsTerminal(t)
		resp := &http.Response{
			Header: http.Header{"Content-Disposition": []string{"attachment; filename=\"a\x1b]52;c;aGVsbG8=\x07b.bin\""}},
			Body:   io.NopCloser(bytes.NewReader([]byte{0, 1, 2, 3})),
		}

		msg, err := writeBinaryResponse(resp, io.Discard, "")

		require.NoError(t, err)
		assert.Equal(t, `Wrote output to: a\x1b]52;c;aGVsbG8=\ab.bin`, msg)
	})

	t.Run("text body on a terminal is escaped", func(t *testing.T) {
		stdoutIsTerminal(t)
		var buf bytes.Buffer
		resp := &http.Response{Body: io.NopCloser(strings.NewReader("one\x1b[2J\ntwo"))}

		_, err := writeBinaryResponse(resp, &buf, "")

		require.NoError(t, err)
		assert.Equal(t, `one\x1b[2J`+"\ntwo", buf.String())
	})

	t.Run("write to stdout", func(t *testing.T) {
		t.Parallel()

		var buf bytes.Buffer
		body := []byte("stdout content")
		resp := &http.Response{
			Body: io.NopCloser(bytes.NewReader(body)),
		}
		msg, err := writeBinaryResponse(resp, &buf, "-")

		require.NoError(t, err)
		assert.Empty(t, msg)
		assert.Equal(t, body, buf.Bytes())
	})
}

func TestCreateDownloadFile(t *testing.T) {
	t.Run("creates file with filename from header", func(t *testing.T) {
		t.Chdir(t.TempDir())

		resp := &http.Response{
			Header: http.Header{
				"Content-Disposition": []string{`attachment; filename="test.txt"`},
			},
		}
		file, err := createDownloadFile(resp, []byte("test content"))
		require.NoError(t, err)
		defer file.Close()
		assert.Equal(t, "test.txt", filepath.Base(file.Name()))

		// Create a second file with the same name to ensure it doesn't clobber the first
		resp2 := &http.Response{
			Header: http.Header{
				"Content-Disposition": []string{`attachment; filename="test.txt"`},
			},
		}
		file2, err := createDownloadFile(resp2, []byte("second content"))
		require.NoError(t, err)
		defer file2.Close()
		assert.NotEqual(t, file.Name(), file2.Name(), "second file should have a different name")
		assert.Contains(t, filepath.Base(file2.Name()), "test")
	})

	t.Run("creates temp file when no header", func(t *testing.T) {
		t.Chdir(t.TempDir())

		resp := &http.Response{Header: http.Header{}}
		file, err := createDownloadFile(resp, []byte("test content"))
		require.NoError(t, err)
		defer file.Close()
		assert.Contains(t, filepath.Base(file.Name()), "file-")
	})

	t.Run("escapes control sequences in the filename", func(t *testing.T) {
		t.Chdir(t.TempDir())

		resp := &http.Response{
			Header: http.Header{"Content-Disposition": []string{"attachment; filename=\"a\x1b]52;c;aGVsbG8=\x07b.bin\""}},
		}
		file, err := createDownloadFile(resp, []byte("test content"))
		require.NoError(t, err)
		defer file.Close()

		assert.Equal(t, `a\x1b]52;c;aGVsbG8=\ab.bin`, filepath.Base(file.Name()))
	})

	t.Run("prevents directory traversal", func(t *testing.T) {
		t.Chdir(t.TempDir())

		resp := &http.Response{
			Header: http.Header{
				"Content-Disposition": []string{`attachment; filename="../../../etc/passwd"`},
			},
		}
		file, err := createDownloadFile(resp, []byte("test content"))
		require.NoError(t, err)
		defer file.Close()
		assert.Equal(t, "passwd", filepath.Base(file.Name()))
	})
}

func TestValidateBaseURL(t *testing.T) {
	t.Parallel()

	t.Run("ValidHTTPS", func(t *testing.T) {
		t.Parallel()

		require.NoError(t, ValidateBaseURL("https://api.example.com", "--base-url"))
	})

	t.Run("ValidHTTP", func(t *testing.T) {
		t.Parallel()

		require.NoError(t, ValidateBaseURL("http://localhost:8080", "--base-url"))
	})

	t.Run("Empty", func(t *testing.T) {
		t.Parallel()

		require.NoError(t, ValidateBaseURL("", "MY_BASE_URL"))
	})

	t.Run("MissingScheme", func(t *testing.T) {
		t.Parallel()

		err := ValidateBaseURL("localhost:8080", "MY_BASE_URL")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "MY_BASE_URL")
		assert.Contains(t, err.Error(), "missing a scheme")
	})

	t.Run("HostOnly", func(t *testing.T) {
		t.Parallel()

		err := ValidateBaseURL("api.example.com", "--base-url")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "--base-url")
	})
}

func TestFormatJSON(t *testing.T) {
	t.Parallel()

	t.Run("RawWithTransform", func(t *testing.T) {
		t.Parallel()

		res := gjson.Parse(`{"id":"abc123","name":"test"}`)
		formatted, err := formatJSON(res, ShowJSONOpts{Format: "raw", Stdout: os.Stdout, Transform: "id"})
		require.NoError(t, err)
		require.Equal(t, `"abc123"`+"\n", string(formatted))
	})

	t.Run("RawWithoutTransform", func(t *testing.T) {
		t.Parallel()

		res := gjson.Parse(`{"id":"abc123","name":"test"}`)
		formatted, err := formatJSON(res, ShowJSONOpts{Format: "raw", Stdout: os.Stdout})
		require.NoError(t, err)
		require.Equal(t, `{"id":"abc123","name":"test"}`+"\n", string(formatted))
	})

	t.Run("RawWithNestedTransform", func(t *testing.T) {
		t.Parallel()

		res := gjson.Parse(`{"data":{"items":[1,2,3]}}`)
		formatted, err := formatJSON(res, ShowJSONOpts{Format: "raw", Stdout: os.Stdout, Transform: "data.items"})
		require.NoError(t, err)
		require.Equal(t, "[1,2,3]\n", string(formatted))
	})

	t.Run("RawWithNonexistentTransform", func(t *testing.T) {
		t.Parallel()

		res := gjson.Parse(`{"id":"abc123"}`)
		formatted, err := formatJSON(res, ShowJSONOpts{Format: "raw", Stdout: os.Stdout, Transform: "missing"})
		require.NoError(t, err)
		// Transform path doesn't exist, so original result is returned
		require.Equal(t, `{"id":"abc123"}`+"\n", string(formatted))
	})

	t.Run("RawOutputString", func(t *testing.T) {
		t.Parallel()

		res := gjson.Parse(`{"id":"abc123","name":"test"}`)
		formatted, err := formatJSON(res, ShowJSONOpts{Format: "json", Stdout: os.Stdout, Transform: "id", RawOutput: true})
		require.NoError(t, err)
		require.Equal(t, "abc123\n", string(formatted))
	})

	t.Run("RawOutputNonString", func(t *testing.T) {
		t.Parallel()

		// --raw-output has no effect on non-string values
		res := gjson.Parse(`{"count":42}`)
		formatted, err := formatJSON(res, ShowJSONOpts{Format: "raw", Stdout: os.Stdout, Transform: "count", RawOutput: true})
		require.NoError(t, err)
		require.Equal(t, "42\n", string(formatted))
	})

	t.Run("RawOutputObject", func(t *testing.T) {
		t.Parallel()

		// --raw-output has no effect on objects
		res := gjson.Parse(`{"nested":{"a":1}}`)
		formatted, err := formatJSON(res, ShowJSONOpts{Format: "raw", Stdout: os.Stdout, Transform: "nested", RawOutput: true})
		require.NoError(t, err)
		require.Equal(t, `{"a":1}`+"\n", string(formatted))
	})
}

// stdoutIsTerminal stands a pty master in for os.Stdout, which also makes any writer count as feeding the pager.
func stdoutIsTerminal(t *testing.T) {
	t.Helper()

	pty := openTerminal(t)
	realStdout := os.Stdout
	os.Stdout = pty
	t.Cleanup(func() { os.Stdout = realStdout })
}

func openTerminal(t *testing.T) *os.File {
	t.Helper()

	pty, err := os.OpenFile("/dev/ptmx", os.O_RDWR, 0)
	if err != nil {
		t.Skipf("no pty available: %v", err)
	}
	t.Cleanup(func() { pty.Close() })
	return pty
}

// showToFile runs ShowJSON with a file as Stdout, so that what it writes can be read back.
func showToFile(t *testing.T, res gjson.Result, opts ShowJSONOpts) string {
	t.Helper()

	f, err := os.CreateTemp(t.TempDir(), "out")
	require.NoError(t, err)
	defer f.Close()
	opts.Stdout = f

	require.NoError(t, ShowJSON(res, opts))

	written, err := os.ReadFile(f.Name())
	require.NoError(t, err)
	return string(written)
}

// Forms of an OSC 52 clipboard write: JSON-encoded, raw, as the visible text it becomes, and that text as JSON.
const (
	encodedOSC52     = `\u001b]52;c;aGVsbG8=\u0007`
	rawOSC52         = "\x1b]52;c;aGVsbG8=\x07"
	escapedOSC52     = `\x1b]52;c;aGVsbG8=\a`
	jsonEscapedOSC52 = `\\x1b]52;c;aGVsbG8=\\a`
)

func assertNoLiveControls(t *testing.T, got string) {
	t.Helper()

	assert.NotContains(t, got, "\x1b]52")
	assert.NotContains(t, got, "\x1b[2J")
	assert.NotContains(t, got, "\x07")
	assert.NotContains(t, got, "\u009b")
	assert.NotContains(t, got, "\u202e")
}

func TestShowJSONWhenRenderingToTerminal(t *testing.T) {
	tests := []struct {
		name   string
		format string
		body   string
		want   string
	}{
		{"JSONEncodedEscape", "json", `{"a":"` + encodedOSC52 + `"}`, jsonEscapedOSC52},
		{"JSONRawBEL", "json", `{"a":"` + rawOSC52 + `"}`, jsonEscapedOSC52},
		{"JSONRawST", "json", "{\"a\":\"\x1b]52;c;aGVsbG8=\x1b\\\\\"}", `\\x1b]52;c;aGVsbG8=\\x1b\\`},
		{"JSONRawCSI", "json", "{\"a\":\"\x1b[2J\"}", `\\x1b[2J`},
		{"JSONRawC1", "json", "{\"a\":\"\u009b2J\"}", `\\u009b2J`},
		{"JSONLoneEscape", "json", "{\"a\":\"\x1bx\"}", `\\x1bx`},
		{"JSONLRawBEL", "jsonl", `{"a":"` + rawOSC52 + `"}`, jsonEscapedOSC52},
		{"RawRawBEL", "raw", `{"a":"` + rawOSC52 + `"}`, jsonEscapedOSC52},
		{"RawRawC1", "raw", "{\"a\":\"\u009b2J\"}", `\\u009b2J`},
		{"PrettyEncodedEscape", "pretty", `{"a":"` + encodedOSC52 + `"}`, escapedOSC52},
		{"PrettyRawBEL", "pretty", `{"a":"` + rawOSC52 + `"}`, escapedOSC52},
		{"YAMLRawBEL", "yaml", `{"a":"` + rawOSC52 + `"}`, escapedOSC52},
		{"YAMLEncodedC1", "yaml", `{"a":"\u009b2J"}`, `\u009b2J`},
		{"YAMLEncodedBidi", "yaml", `{"a":"\u202e"}`, `\u202e`},
		{"YAMLCRLFWhitespace", "yaml", "{\"a\":1,\r\n\"b\":2}", "b: 2"},
		{"RawCRLFWhitespace", "raw", "{\"a\":1,\r\n\"b\":2}", "{\"a\":1, \n\"b\":2}"},
		{"JSONLCRLFWhitespace", "jsonl", "{\"a\":1,\r\n\"b\":2}", `{"a":1,"b":2}`},
		{"JSONInvalidEscapeBeforeRaw", "json", `{"a":"\q` + rawOSC52 + `"}`, `\\q` + jsonEscapedOSC52},
		{"PrettyInvalidEscapeAfterEncoded", "pretty", `{"a":"` + encodedOSC52 + `\q"}`, escapedOSC52},
		{"RawUnterminatedString", "raw", `{"a":"` + rawOSC52, jsonEscapedOSC52},
		{"PrettyUnterminatedString", "pretty", `"` + encodedOSC52, escapedOSC52},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			stdoutIsTerminal(t)

			got := showToFile(t, gjson.Parse(test.body), ShowJSONOpts{Format: test.format})

			assertNoLiveControls(t, got)
			assert.Contains(t, got, test.want)
		})
	}

	t.Run("StdoutOptionIsTheTerminal", func(t *testing.T) {
		t.Parallel()

		res := gjson.Parse(`{"a":"` + rawOSC52 + `"}`)

		formatted, err := formatJSON(res, ShowJSONOpts{Format: "raw", Stdout: openTerminal(t)})

		require.NoError(t, err)
		assert.Equal(t, `{"a":"`+jsonEscapedOSC52+`"}`+"\n", string(formatted))
	})

	t.Run("JSONStaysValid", func(t *testing.T) {
		stdoutIsTerminal(t)
		body := "{\"a\":\"x\u009by\\n\u202e世界\"}"

		got := showToFile(t, gjson.Parse(body), ShowJSONOpts{Format: "json"})

		assert.True(t, json.Valid([]byte(got)))
		assert.Equal(t, `x\u009by`+"\n"+`\u202e世界`, gjson.Get(got, "a").Str)
	})

	t.Run("CleanJSONKeepsItsValues", func(t *testing.T) {
		stdoutIsTerminal(t)
		body := `{"a":"line\n\ttwo 世界 \"q\" \\"}`

		got := showToFile(t, gjson.Parse(body), ShowJSONOpts{Format: "raw"})

		assert.Equal(t, body+"\n", got)
	})

	t.Run("RawOutputString", func(t *testing.T) {
		stdoutIsTerminal(t)
		res := gjson.Parse(`{"id":"héllo` + encodedOSC52 + `\n\tline \"two\""}`)

		got := showToFile(t, res, ShowJSONOpts{Format: "json", Transform: "id", RawOutput: true})

		assert.Equal(t, `héllo`+escapedOSC52+"\n\tline \"two\"\n", got)
	})

	t.Run("RawOutputStringWithInvalidEscape", func(t *testing.T) {
		stdoutIsTerminal(t)
		res := gjson.Parse(`{"id":"` + encodedOSC52 + `\q"}`)

		got := showToFile(t, res, ShowJSONOpts{Format: "json", Transform: "id", RawOutput: true})

		assertNoLiveControls(t, got)
		assert.Contains(t, got, escapedOSC52)
	})

	t.Run("TransformThatDecodesAStringIntoJSON", func(t *testing.T) {
		stdoutIsTerminal(t)
		res := gjson.Parse(`{"a":"\"\\u001b]52;c;aGVsbG8=\\u0007\""}`)

		got := showToFile(t, res, ShowJSONOpts{Format: "pretty", Transform: "a.@fromstr"})

		assertNoLiveControls(t, got)
		assert.Contains(t, got, escapedOSC52)
	})
}

// Goes over OutputFormats, so a format added there is covered without changing this test.
func TestShowJSONOfEveryFormatIsSafeOnATerminal(t *testing.T) {
	body := `{"k` + encodedOSC52 + `":["` + encodedOSC52 + `","` + rawOSC52 + "\u009b\u202e\"]}"

	for _, format := range OutputFormats {
		t.Run(format, func(t *testing.T) {
			stdoutIsTerminal(t)

			got := showToFile(t, gjson.Parse(body), ShowJSONOpts{Format: format, Stderr: io.Discard, Title: "Result"})

			assertNoLiveControls(t, got)
			assert.Contains(t, got, `x1b]52;c;aGVsbG8=`)
		})
	}
}

func TestShowJSONIteratorWithForcedColor(t *testing.T) {
	t.Setenv("FORCE_COLOR", "1")
	iter := &sliceIterator[map[string]any]{items: []map[string]any{{"a": rawOSC52}}}

	got := captureShowJSONIterator(t, iter, "raw", "", -1)

	assertNoLiveControls(t, got)
	assert.Contains(t, got, jsonEscapedOSC52)
}

// main passes an API error's raw body to ShowJSON in the format of --format-error.
func TestShowJSONOfAPIErrorBody(t *testing.T) {
	body := `{"type":"error","error":{"message":"denied ` + encodedOSC52 + ` and ` + rawOSC52 + `"}}`

	tests := []struct {
		format string
		want   string
	}{
		{"auto", jsonEscapedOSC52},
		{"explore", jsonEscapedOSC52},
		{"json", jsonEscapedOSC52},
		{"jsonl", jsonEscapedOSC52},
		{"pretty", escapedOSC52},
		{"raw", jsonEscapedOSC52},
		{"yaml", escapedOSC52},
	}
	for _, test := range tests {
		t.Run(test.format, func(t *testing.T) {
			stdoutIsTerminal(t)

			got := showToFile(t, gjson.Parse(body), ShowJSONOpts{Format: test.format, Stderr: io.Discard, Title: "Error"})

			assertNoLiveControls(t, got)
			assert.Contains(t, got, test.want)
		})
	}
}

func TestShowJSONWithForcedColor(t *testing.T) {
	hostile := `{"a":"` + rawOSC52 + "\u009b2J\"}"

	tests := []struct {
		name string
		opts ShowJSONOpts
		want string
	}{
		{"JSON", ShowJSONOpts{Format: "json"}, `\u009b`},
		{"JSONL", ShowJSONOpts{Format: "jsonl"}, `\u009b`},
		{"Raw", ShowJSONOpts{Format: "raw"}, jsonEscapedOSC52 + `\\u009b2J`},
		{"RawOutputString", ShowJSONOpts{Format: "json", Transform: "a", RawOutput: true}, escapedOSC52 + `\u009b2J`},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Setenv("FORCE_COLOR", "1")

			got := showToFile(t, gjson.Parse(hostile), test.opts)

			assertNoLiveControls(t, got)
			assert.Contains(t, got, test.want)
		})
	}
}

func TestShowJSONWhenPiped(t *testing.T) {
	t.Parallel()

	if shouldUseColors(os.Stdout) || isTerminal(os.Stdout) {
		t.Skip("a terminal on the test's own stdout or FORCE_COLOR=1 counts as rendering to one")
	}

	hostile := `{"a":"` + rawOSC52 + "\u009b\"}"

	t.Run("RawIsByteExact", func(t *testing.T) {
		t.Parallel()

		got := showToFile(t, gjson.Parse(hostile), ShowJSONOpts{Format: "raw"})

		assert.Equal(t, hostile+"\n", got)
	})

	t.Run("JSONLIsByteExact", func(t *testing.T) {
		t.Parallel()

		got := showToFile(t, gjson.Parse(hostile), ShowJSONOpts{Format: "jsonl"})

		assert.Equal(t, hostile+"\n", got)
	})

	t.Run("RawOutputStringIsByteExact", func(t *testing.T) {
		t.Parallel()

		res := gjson.Parse(`{"id":"a\u001b[31m\tb"}`)

		got := showToFile(t, res, ShowJSONOpts{Format: "json", Transform: "id", RawOutput: true})

		assert.Equal(t, "a\x1b[31m\tb\n", got)
	})

	t.Run("BinaryResponseIsByteExact", func(t *testing.T) {
		t.Parallel()

		body := "one\x1b[2J\ntwo"
		var buf bytes.Buffer
		resp := &http.Response{Body: io.NopCloser(strings.NewReader(body))}

		_, err := writeBinaryResponse(resp, &buf, "")

		require.NoError(t, err)
		assert.Equal(t, body, buf.String())
	})

	t.Run("PrettyIsStillEscaped", func(t *testing.T) {
		t.Parallel()

		got := showToFile(t, gjson.Parse(hostile), ShowJSONOpts{Format: "pretty"})

		assertNoLiveControls(t, got)
		assert.Contains(t, got, escapedOSC52)
	})
}

func TestShowJSONIterator(t *testing.T) {
	t.Parallel()

	t.Run("RawMultipleItems", func(t *testing.T) {
		t.Parallel()

		iter := &sliceIterator[map[string]any]{items: []map[string]any{
			{"id": "abc", "name": "first"},
			{"id": "def", "name": "second"},
		}}
		captured := captureShowJSONIterator(t, iter, "raw", "", -1)
		assert.Equal(t, `{"id":"abc","name":"first"}`+"\n"+`{"id":"def","name":"second"}`+"\n", captured)
	})

	t.Run("RawWithTransform", func(t *testing.T) {
		t.Parallel()

		iter := &sliceIterator[map[string]any]{items: []map[string]any{
			{"id": "abc", "name": "first"},
			{"id": "def", "name": "second"},
		}}
		captured := captureShowJSONIterator(t, iter, "raw", "id", -1)
		assert.Equal(t, `"abc"`+"\n"+`"def"`+"\n", captured)
	})

	t.Run("LimitItems", func(t *testing.T) {
		t.Parallel()

		iter := &sliceIterator[map[string]any]{items: []map[string]any{
			{"id": "abc"},
			{"id": "def"},
			{"id": "ghi"},
		}}
		captured := captureShowJSONIterator(t, iter, "raw", "", 2)
		assert.Equal(t, `{"id":"abc"}`+"\n"+`{"id":"def"}`+"\n", captured)
	})
}

func TestExploreFallback(t *testing.T) {
	t.Parallel()

	t.Run("ShowJSONFallsBackToJsonOnNonTTY", func(t *testing.T) {
		t.Parallel()

		// os.Pipe() produces a *os.File that isn't a terminal, so explore should fall back.
		r, w, err := os.Pipe()
		require.NoError(t, err)
		defer r.Close()

		var stderr bytes.Buffer
		res := gjson.Parse(`{"id":"abc"}`)
		err = ShowJSON(res, ShowJSONOpts{
			Format: "explore",
			Stderr: &stderr,
			Stdout: w,
			Title:  "test",
		})
		w.Close()
		require.NoError(t, err)

		var buf bytes.Buffer
		_, _ = buf.ReadFrom(r)
		assert.Contains(t, buf.String(), `"id"`)
		assert.Contains(t, buf.String(), `"abc"`)
	})

	t.Run("ShowJSONIteratorFallsBackToJsonOnNonTTY", func(t *testing.T) {
		t.Parallel()

		iter := &sliceIterator[map[string]any]{items: []map[string]any{
			{"id": "abc"},
		}}
		captured := captureShowJSONIterator(t, iter, "explore", "", -1)
		assert.Contains(t, captured, `"id"`)
		assert.Contains(t, captured, `"abc"`)
	})

	t.Run("ShowJSONWarnsWhenExplicitFormatOnNonTTY", func(t *testing.T) {
		t.Parallel()

		r, w, err := os.Pipe()
		require.NoError(t, err)
		defer r.Close()

		var stderr bytes.Buffer
		res := gjson.Parse(`{"id":"abc"}`)
		err = ShowJSON(res, ShowJSONOpts{
			ExplicitFormat: true,
			Format:         "explore",
			Stderr:         &stderr,
			Stdout:         w,
			Title:          "test",
		})
		w.Close()
		require.NoError(t, err)

		assert.Equal(t, warningExploreNotSupported, stderr.String())
	})

	t.Run("ShowJSONSilentWhenDefaultFormatOnNonTTY", func(t *testing.T) {
		t.Parallel()

		r, w, err := os.Pipe()
		require.NoError(t, err)
		defer r.Close()

		var stderr bytes.Buffer
		res := gjson.Parse(`{"id":"abc"}`)
		err = ShowJSON(res, ShowJSONOpts{
			Format: "explore",
			Stderr: &stderr,
			Stdout: w,
			Title:  "test",
		})
		w.Close()
		require.NoError(t, err)

		assert.Empty(t, stderr.String(), "no warning expected when format was not explicit")
	})
}

// sliceIterator is a simple iterator over a slice for testing.
type sliceIterator[T any] struct {
	index int
	items []T
}

func (it *sliceIterator[T]) Next() bool {
	it.index++
	return it.index <= len(it.items)
}

func (it *sliceIterator[T]) Current() T {
	return it.items[it.index-1]
}

func (it *sliceIterator[T]) Err() error {
	return nil
}

var _ jsonview.Iterator[any] = (*sliceIterator[any])(nil)

// captureShowJSONIterator runs ShowJSONIterator and captures the output written to a file.
func captureShowJSONIterator[T any](t *testing.T, iter jsonview.Iterator[T], format, transform string, itemsToDisplay int64) string {
	t.Helper()

	r, w, err := os.Pipe()
	require.NoError(t, err)
	defer r.Close()

	err = ShowJSONIterator(iter, itemsToDisplay, ShowJSONOpts{
		Format:    format,
		Stderr:    io.Discard,
		Stdout:    w,
		Title:     "test",
		Transform: transform,
	})
	w.Close()
	require.NoError(t, err)

	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	return buf.String()
}
