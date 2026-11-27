package cmd

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os/exec"
	"sync/atomic"
	"testing"

	"github.com/anthropics/anthropic-cli/v2/internal/mocktest"
	"github.com/stretchr/testify/require"
)

func TestNullableStringFileInput(t *testing.T) {
	clearCredentialEnv(t)
	t.Setenv("ANTHROPIC_API_KEY", "test-key")
	const prompt = "You are a helpful assistant.\nKeep answers concise."
	path := mocktest.TestFile(t, prompt)

	for _, tc := range []struct {
		name string
		args []string
		want string
	}{
		{"file", []string{"--system", "@" + path}, fmt.Sprintf(`{"system":%q}`, prompt)},
		{"text file", []string{"--system", "@file://" + path}, fmt.Sprintf(`{"system":%q}`, prompt)},
		{"base64 file", []string{"--system", "@data://" + path}, fmt.Sprintf(`{"system":%q}`, base64.StdEncoding.EncodeToString([]byte(prompt)))},
		{"escaped reference", []string{"--system", `\@prompt.txt`}, `{"system":"@prompt.txt"}`},
		{"literal", []string{"--system", prompt}, fmt.Sprintf(`{"system":%q}`, prompt)},
		{"null", []string{"--system", "null"}, `{"system":null}`},
		{"omitted", nil, `{}`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var body []byte
			var requests atomic.Int32
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if requests.Add(1) != 1 || r.Method != http.MethodPost || r.URL.Path != "/v1/agents" {
					t.Errorf("unexpected request: %s %s (request %d)", r.Method, r.URL.Path, requests.Load())
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				var request map[string]json.RawMessage
				if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
					t.Error(err)
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				delete(request, "model")
				delete(request, "name")
				body, _ = json.Marshal(request)
				w.Header().Set("Content-Type", "application/json")
				fmt.Fprint(w, `{}`)
			}))
			t.Cleanup(srv.Close)

			args := []string{"run", "./cmd/ant", "--base-url", srv.URL,
				"beta:agents", "create", "--model", "claude-opus-5-5", "--name", "test"}
			cmd := exec.Command("go", append(args, tc.args...)...)
			cmd.Dir = "../.."
			out, err := cmd.CombinedOutput()
			srv.Close()
			require.NoError(t, err, "%s", out)
			require.EqualValues(t, 1, requests.Load())
			require.JSONEq(t, tc.want, string(body))
		})
	}
}

func TestEmbedFilesNestedNullableStrings(t *testing.T) {
	t.Parallel()
	const prompt = "Nested prompt contents."
	reference := "@" + mocktest.TestFile(t, prompt)
	var nullString *string

	embedded, err := embedFiles(map[string]any{
		"nested": []any{map[string]any{
			"system":      &reference,
			"description": nullString,
		}},
	}, EmbedText, nil)
	require.NoError(t, err)
	require.Equal(t, map[string]any{
		"nested": []any{map[string]any{
			"system":      prompt,
			"description": nullString,
		}},
	}, embedded)
}
