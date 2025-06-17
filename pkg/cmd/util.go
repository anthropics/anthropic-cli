package cmd

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/stainless-sdks/anthropic-cli/pkg/jsonflag"

	"github.com/tidwall/gjson"
	"github.com/tidwall/pretty"
	"github.com/urfave/cli/v3"
	"golang.org/x/term"
)

func getDefaultRequestOptions() []option.RequestOption {
	return []option.RequestOption{
		option.WithHeader("X-Stainless-Lang", "cli"),
		option.WithHeader("X-Stainless-Runtime", "cli"),
	}
}

type apiCommandContext struct {
	client anthropic.Client
	cmd    *cli.Command
}

func (c apiCommandContext) AsMiddleware() option.Middleware {
	body := getStdInput()
	if body == nil {
		body = []byte("{}")
	}
	var query = []byte("{}")
	var header = []byte("{}")

	// Apply JSON flag mutations
	body, query, header, err := jsonflag.Apply(body, query, header)
	if err != nil {
		log.Fatal(err)
	}

	debug := c.cmd.Bool("debug")

	return func(r *http.Request, mn option.MiddlewareNext) (*http.Response, error) {
		q := r.URL.Query()
		for key, values := range serializeQuery(query) {
			for _, value := range values {
				q.Set(key, value)
			}
		}
		r.URL.RawQuery = q.Encode()

		for key, values := range serializeHeader(header) {
			for _, value := range values {
				r.Header.Add(key, value)
			}
		}

		r.Body = io.NopCloser(bytes.NewBuffer(body))

		// Add debug logging if the --debug flag is set
		if debug {
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
		}

		return mn(r)
	}
}

func getAPICommandContext(cmd *cli.Command) *apiCommandContext {
	client := anthropic.NewClient(getDefaultRequestOptions()...)
	return &apiCommandContext{client, cmd}
}

func serializeQuery(params []byte) url.Values {
	serialized := url.Values{}

	var serialize func(value gjson.Result, path string)
	serialize = func(res gjson.Result, path string) {
		if res.IsObject() {
			for key, value := range res.Map() {
				newPath := path
				if len(newPath) == 0 {
					newPath += key
				} else {
					newPath = "[" + key + "]"
				}

				serialize(value, newPath)
			}
		} else if res.IsArray() {
			for _, value := range res.Array() {
				serialize(value, path)
			}
		} else {
			serialized.Add(path, res.String())
		}
	}
	serialize(gjson.GetBytes(params, "@this"), "")

	for key, values := range serialized {
		serialized.Set(key, strings.Join(values, ","))
	}

	return serialized
}

func serializeHeader(params []byte) http.Header {
	serialized := http.Header{}

	var serialize func(value gjson.Result, path string)
	serialize = func(res gjson.Result, path string) {
		if res.IsObject() {
			for key, value := range res.Map() {
				newPath := path
				if len(newPath) > 0 {
					newPath += "."
				}
				newPath += key

				serialize(value, newPath)
			}
		} else if res.IsArray() {
			for _, value := range res.Array() {
				serialize(value, path)
			}
		} else {
			serialized.Add(path, res.String())
		}
	}
	serialize(gjson.GetBytes(params, "@this"), "")

	return serialized
}

func getStdInput() []byte {
	if !isInputPiped() {
		return nil
	}
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	return data
}

func isInputPiped() bool {
	stat, _ := os.Stdin.Stat()
	return (stat.Mode() & os.ModeCharDevice) == 0
}

func isTerminal(w io.Writer) bool {
	switch v := w.(type) {
	case *os.File:
		return term.IsTerminal(int(v.Fd()))
	default:
		return false
	}
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

	if isTerminal(w) {
		return true
	}
	return false

}

func ColorizeJSON(input string, w io.Writer) string {
	if !shouldUseColors(w) {
		return input
	}
	return string(pretty.Color(pretty.Pretty([]byte(input)), nil))
}
