package cmd

import (
	"context"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"

	"github.com/tidwall/gjson"
	"github.com/tidwall/pretty"
	"github.com/tidwall/sjson"
	"github.com/urfave/cli/v3"
	"golang.org/x/term"
)

func jsonSet(json []byte, path string, value interface{}) ([]byte, error) {
	keys := strings.Split(path, ".")
	path = ""
	for i := 0; i < len(keys); i++ {
		key := keys[i]
		if key == "#" {
			key = strconv.Itoa(len(gjson.GetBytes(json, path).Array()) - 1)
		}

		if len(path) > 0 {
			path += "."
		}
		path += key
	}
	return sjson.SetBytes(json, path, value)
}

type apiCommandKey string

type apiCommandContext struct {
	client anthropic.Client
	body   []byte
	query  []byte
	header []byte
}

func (c apiCommandContext) AsMiddleware() option.Middleware {
	return func(r *http.Request, mn option.MiddlewareNext) (*http.Response, error) {
		q := r.URL.Query()
		for key, values := range serializeQuery(c.query) {
			for _, value := range values {
				q.Add(key, value)
			}
		}
		r.URL.RawQuery = q.Encode()

		for key, values := range serializeHeader(c.header) {
			for _, value := range values {
				r.Header.Add(key, value)
			}
		}

		return mn(r)
	}
}

func initAPICommand(ctx context.Context, cmd *cli.Command) (context.Context, error) {
	client := anthropic.NewClient()
	body := getStdInput()
	if body == nil {
		body = []byte("{}")
	}
	var query = []byte("{}")
	var header = []byte("{}")

	return context.WithValue(ctx, apiCommandKey(cmd.Name), &apiCommandContext{client, body, query, header}), nil
}

func getAPICommandContext(ctx context.Context, cmd *cli.Command) *apiCommandContext {
	return ctx.Value(apiCommandKey(cmd.Name)).(*apiCommandContext)
}

func getAPIFlagAction[T any](kind string, path string) func(context.Context, *cli.Command, T) error {
	return func(ctx context.Context, cmd *cli.Command, value T) (err error) {
		commandContext := getAPICommandContext(ctx, cmd)
		var dest *[]byte
		switch kind {
		case "body":
			dest = &commandContext.body
		case "query":
			dest = &commandContext.query
		case "header":
			dest = &commandContext.header
		}
		*dest, err = jsonSet(*dest, path, value)
		return err
	}
}

func getAPIFlagActionWithValue[T any](kind string, path string, value interface{}) func(context.Context, *cli.Command, T) error {
	return func(ctx context.Context, cmd *cli.Command, unusedValue T) (err error) {
		commandContext := getAPICommandContext(ctx, cmd)
		var dest *[]byte
		switch kind {
		case "body":
			dest = &commandContext.body
		case "query":
			dest = &commandContext.query
		case "header":
			dest = &commandContext.header
		}
		*dest, err = jsonSet(*dest, path, value)
		return err
	}
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

func colorizeJSON(input string, w io.Writer) string {
	if !shouldUseColors(w) {
		return input
	}
	return string(pretty.Color(pretty.Pretty([]byte(input)), nil))
}
