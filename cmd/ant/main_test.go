package main

import (
	"bytes"
	"context"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"testing"

	"github.com/anthropics/anthropic-cli/v2/pkg/cmd"
	"github.com/anthropics/anthropic-sdk-go/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli/v3"
)

// TestMain runs main in place of the tests in a process that runMain started.
func TestMain(m *testing.M) {
	if body, ok := os.LookupEnv("TEST_API_ERROR_BODY"); ok {
		status, _ := strconv.Atoi(os.Getenv("TEST_API_ERROR_STATUS"))
		request, _ := http.NewRequest(http.MethodPost, "https://example.com/things", nil)
		apierr := &anthropic.Error{
			StatusCode: status,
			Request:    request,
			Response:   &http.Response{StatusCode: status},
		}
		_ = apierr.UnmarshalJSON([]byte(body))
		cmd.Command.Action = func(context.Context, *cli.Command) error { return apierr }
		main()
		return
	}
	os.Exit(m.Run())
}

// runMain runs main in a new process, where the command fails with an API error that has the given HTTP status and body.
func runMain(t *testing.T, status int, body string, args ...string) (stdout string, stderr string, exitCode int) {
	t.Helper()
	executable, err := os.Executable()
	require.NoError(t, err)
	process := exec.Command(executable, args...)
	process.Env = append(os.Environ(), "FORCE_COLOR=0", "TEST_API_ERROR_STATUS="+strconv.Itoa(status), "TEST_API_ERROR_BODY="+body)
	var stdoutBuffer, stderrBuffer bytes.Buffer
	process.Stdout, process.Stderr = &stdoutBuffer, &stderrBuffer

	err = process.Run()

	var exitErr *exec.ExitError
	require.ErrorAs(t, err, &exitErr)
	return stdoutBuffer.String(), stderrBuffer.String(), exitErr.ExitCode()
}

func TestAPIErrorOutput(t *testing.T) {
	tests := []struct {
		name   string
		status int
		stderr string
	}{
		{"ErrorStatus", http.StatusTooManyRequests, "POST \"https://example.com/things\": 429 Too Many Requests\n"},
		{"ErrorEventInStream", http.StatusOK, "POST \"https://example.com/things\": overloaded_error: Overloaded\n"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			body := `{"error":{"type":"overloaded_error","message":"Overloaded"}}`

			stdout, stderr, exitCode := runMain(t, test.status, body, "--format-error", "json")

			assert.JSONEq(t, body, stdout)
			assert.Equal(t, test.stderr, stderr)
			assert.Equal(t, 1, exitCode)
		})
	}
}
