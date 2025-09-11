// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cmd

import (
	"bytes"
	"fmt"
	"golang.org/x/term"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/itchyny/json2yaml"
	"github.com/stainless-sdks/anthropic-cli/pkg/jsonflag"
	"github.com/stainless-sdks/anthropic-cli/pkg/jsonview"
	"github.com/tidwall/gjson"
	"github.com/tidwall/pretty"
	"github.com/urfave/cli/v3"
)

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
	body, query, header, err := jsonflag.ApplyMutations(body, query, header)
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

		if r.Body != nil || len(body) > 2 {
			r.Body = io.NopCloser(bytes.NewBuffer(body))
			r.ContentLength = int64(len(body))
			r.Header.Set("Content-Type", "application/json")
		}

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
	client := anthropic.NewClient(getDefaultRequestOptions(cmd)...)
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

	return isTerminal(w)
}

func ShowJSON(title, jsonText, format string) error {
	switch strings.ToLower(format) {
	case "auto":
		return ShowJSON(title, jsonText, "json")
	case "explore":
		return jsonview.ExploreJSON(title, jsonText)
	case "pretty":
		jsonview.DisplayJSON(title, jsonText)
		return nil
	case "json":
		prettyJSON := pretty.Pretty([]byte(jsonText))
		if shouldUseColors(os.Stdout) {
			fmt.Print(string(pretty.Color(prettyJSON, pretty.TerminalStyle)))
		} else {
			fmt.Print(string(prettyJSON))
		}
		return nil
	case "raw":
		fmt.Println(jsonText)
		return nil
	case "yaml":
		input := strings.NewReader(jsonText)
		var yaml strings.Builder
		if err := json2yaml.Convert(&yaml, input); err != nil {
			return err
		}
		fmt.Print(yaml.String())
		return nil
	default:
		return fmt.Errorf("Invalid format: %s, valid formats are: %s", format, strings.Join(OutputFormats[:], ", "))
	}
}
