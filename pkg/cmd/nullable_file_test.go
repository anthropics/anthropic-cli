package cmd

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os/exec"
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
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
			require.JSONEq(t, tc.want, string(body))
		})
	}
}
