package cmd

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"maps"
	"mime/multipart"
	"net/http"
	"os"
	"reflect"
	"strings"
	"unicode/utf8"

	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/stainless-sdks/anthropic-cli/internal/apiform"
	"github.com/stainless-sdks/anthropic-cli/internal/apiquery"
	"github.com/stainless-sdks/anthropic-cli/internal/debugmiddleware"
	"github.com/stainless-sdks/anthropic-cli/internal/requestflag"

	"github.com/goccy/go-yaml"
	"github.com/urfave/cli/v3"
)

type BodyContentType int

const (
	EmptyBody BodyContentType = iota
	MultipartFormEncoded
	ApplicationJSON
	ApplicationOctetStream
)

type FileEmbedStyle int

const (
	EmbedText FileEmbedStyle = iota
	EmbedIOReader
)

func embedFiles(obj any, embedStyle FileEmbedStyle) (any, error) {
	v := reflect.ValueOf(obj)
	result, err := embedFilesValue(v, embedStyle)
	if err != nil {
		return nil, err
	}
	return result.Interface(), nil
}

// Replace "@file.txt" with the file's contents inside a value
func embedFilesValue(v reflect.Value, embedStyle FileEmbedStyle) (reflect.Value, error) {
	// Unwrap interface values to get the concrete type
	if v.Kind() == reflect.Interface {
		if v.IsNil() {
			return v, nil
		}
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.Map:
		if v.Len() == 0 {
			return v, nil
		}
		// Always create map[string]any to handle potential type changes when embedding files
		result := reflect.MakeMap(reflect.TypeOf(map[string]any{}))

		iter := v.MapRange()
		for iter.Next() {
			key := iter.Key()
			val := iter.Value()
			newVal, err := embedFilesValue(val, embedStyle)
			if err != nil {
				return reflect.Value{}, err
			}
			result.SetMapIndex(key, newVal)
		}
		return result, nil

	case reflect.Slice, reflect.Array:
		if v.Len() == 0 {
			return v, nil
		}
		// Use `[]any` to allow for types to change when embedding files
		result := reflect.MakeSlice(reflect.TypeOf([]any{}), v.Len(), v.Len())
		for i := 0; i < v.Len(); i++ {
			newVal, err := embedFilesValue(v.Index(i), embedStyle)
			if err != nil {
				return reflect.Value{}, err
			}
			result.Index(i).Set(newVal)
		}
		return result, nil

	case reflect.String:
		s := v.String()
		if literal, ok := strings.CutPrefix(s, "\\@"); ok {
			// Allow for escaped @ signs if you don't want them to be treated as files
			return reflect.ValueOf("@" + literal), nil
		}

		if embedStyle == EmbedText {
			if filename, ok := strings.CutPrefix(s, "@data://"); ok {
				// The "@data://" prefix is for files you explicitly want to upload
				// as base64-encoded (even if the file itself is plain text)
				content, err := os.ReadFile(filename)
				if err != nil {
					return v, err
				}
				return reflect.ValueOf(base64.StdEncoding.EncodeToString(content)), nil
			} else if filename, ok := strings.CutPrefix(s, "@file://"); ok {
				// The "@file://" prefix is for files that you explicitly want to
				// upload as a string literal with backslash escapes (not base64
				// encoded)
				content, err := os.ReadFile(filename)
				if err != nil {
					return v, err
				}
				return reflect.ValueOf(string(content)), nil
			} else if filename, ok := strings.CutPrefix(s, "@"); ok {
				content, err := os.ReadFile(filename)
				if err != nil {
					// If the string is "@username", it's probably supposed to be a
					// string literal and not a file reference. However, if the
					// string looks like "@file.txt" or "@/tmp/file", then it's
					// probably supposed to be a file.
					probablyFile := strings.Contains(filename, ".") || strings.Contains(filename, "/")
					if probablyFile {
						// Give a useful error message if the user tried to upload a
						// file, but the file couldn't be read (e.g. mistyped
						// filename or permission error)
						return v, err
					}
					// Fall back to the raw value if the user provided something
					// like "@username" that's not intended to be a file.
					return v, nil
				}
				// If the file looks like a plain text UTF8 file format, then use the contents directly.
				if isUTF8TextFile(content) {
					return reflect.ValueOf(string(content)), nil
				}
				// Otherwise it's a binary file, so encode it with base64
				return reflect.ValueOf(base64.StdEncoding.EncodeToString(content)), nil
			}
		} else {
			if filename, ok := strings.CutPrefix(s, "@"); ok {
				// Behavior is the same for @file, @data://file, and @file://file, except that
				// @username will be treated as a literal string if no "username" file exists
				expectsFile := true
				if withoutPrefix, ok := strings.CutPrefix(filename, "data://"); ok {
					filename = withoutPrefix
				} else if withoutPrefix, ok := strings.CutPrefix(filename, "file://"); ok {
					filename = withoutPrefix
				} else {
					expectsFile = strings.Contains(filename, ".") || strings.Contains(filename, "/")
				}

				file, err := os.Open(filename)
				if err != nil {
					if !expectsFile {
						// For strings that start with "@" and don't look like a filename, return the string
						return v, nil
					}
					return v, err
				}
				return reflect.ValueOf(file), nil
			}
		}
		return v, nil

	default:
		return v, nil
	}
}

// Guess whether a file's contents are binary (e.g. a .jpg or .mp3), as opposed
// to plain text (e.g. .txt or .md).
func isUTF8TextFile(content []byte) bool {
	// Go's DetectContentType follows https://mimesniff.spec.whatwg.org/ and
	// these are the sniffable content types that are plain text:
	textTypes := []string{
		"text/",
		"application/json",
		"application/xml",
		"application/javascript",
		"application/x-javascript",
		"application/ecmascript",
		"application/x-ecmascript",
	}

	contentType := http.DetectContentType(content)
	for _, prefix := range textTypes {
		if strings.HasPrefix(contentType, prefix) {
			return utf8.Valid(content)
		}
	}
	return false
}

func flagOptions(
	cmd *cli.Command,
	nestedFormat apiquery.NestedQueryFormat,
	arrayFormat apiquery.ArrayQueryFormat,
	bodyType BodyContentType,

	// This parameter is true if stdin is already in use to pass a binary parameter by using the special value
	// "-". In this case, we won't attempt to read it as a JSON/YAML blob for options setting.
	stdinInUse bool,
) ([]option.RequestOption, error) {
	var options []option.RequestOption
	if cmd.Bool("debug") {
		options = append(options, option.WithMiddleware(debugmiddleware.NewRequestLogger().Middleware()))
	}

	flagContents := requestflag.ExtractRequestContents(cmd)

	var bodyData any
	if isInputPiped() && !stdinInUse {
		var err error
		pipeData, err := io.ReadAll(os.Stdin)
		if err != nil {
			return nil, err
		}

		if err := yaml.Unmarshal(pipeData, &bodyData); err == nil {
			if bodyMap, ok := bodyData.(map[string]any); ok {
				if flagMap, ok := flagContents.Body.(map[string]any); ok {
					maps.Copy(bodyMap, flagMap)
				} else {
					bodyData = flagContents.Body
				}
			} else if flagMap, ok := flagContents.Body.(map[string]any); ok && len(flagMap) > 0 {
				return nil, fmt.Errorf("Cannot merge flags with a body that is not a map: %v", bodyData)
			}
		}
	} else {
		// No piped input, just use body flag values as a map
		bodyData = flagContents.Body
	}

	// Embed files passed as "@file.jpg" in the request body, headers, and query:
	embedStyle := EmbedText
	if bodyType == ApplicationOctetStream || bodyType == MultipartFormEncoded {
		embedStyle = EmbedIOReader
	}
	bodyData, err := embedFiles(bodyData, embedStyle)
	if err != nil {
		return nil, err
	}
	if headersWithFiles, err := embedFiles(flagContents.Headers, EmbedText); err != nil {
		return nil, err
	} else {
		flagContents.Headers = headersWithFiles.(map[string]any)
	}
	if queriesWithFiles, err := embedFiles(flagContents.Queries, EmbedText); err != nil {
		return nil, err
	} else {
		flagContents.Queries = queriesWithFiles.(map[string]any)
	}

	querySettings := apiquery.QuerySettings{
		NestedFormat: nestedFormat,
		ArrayFormat:  arrayFormat,
	}

	// Add query parameters:
	if values, err := apiquery.MarshalWithSettings(flagContents.Queries, querySettings); err != nil {
		return nil, err
	} else {
		for k, vs := range values {
			if len(vs) == 0 {
				options = append(options, option.WithQueryDel(k))
			} else {
				options = append(options, option.WithQuery(k, vs[0]))
				for _, v := range vs[1:] {
					options = append(options, option.WithQueryAdd(k, v))
				}
			}
		}
	}

	// Add header parameters
	if values, err := apiquery.MarshalWithSettings(flagContents.Headers, querySettings); err != nil {
		return nil, err
	} else {
		for k, vs := range values {
			if len(vs) == 0 {
				options = append(options, option.WithHeaderDel(k))
			} else {
				options = append(options, option.WithHeader(k, vs[0]))
				for _, v := range vs[1:] {
					options = append(options, option.WithHeaderAdd(k, v))
				}
			}
		}
	}

	switch bodyType {
	case EmptyBody:
		break
	case MultipartFormEncoded:
		buf := new(bytes.Buffer)
		writer := multipart.NewWriter(buf)

		// For multipart/form-encoded, we need a map structure
		bodyMap, ok := bodyData.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("Cannot send a non-map value to a form-encoded endpoint: %v\n", bodyData)
		}
		encodingFormat := apiform.FormatComma
		if err := apiform.MarshalWithSettings(bodyMap, writer, encodingFormat); err != nil {
			return nil, err
		}
		if err := writer.Close(); err != nil {
			return nil, err
		}
		options = append(options, option.WithRequestBody(writer.FormDataContentType(), buf))

	case ApplicationJSON:
		bodyBytes, err := json.Marshal(bodyData)
		if err != nil {
			return nil, err
		}
		options = append(options, option.WithRequestBody("application/json", bodyBytes))

	case ApplicationOctetStream:
		if bodyBytes, ok := bodyData.([]byte); ok {
			options = append(options, option.WithRequestBody("application/octet-stream", bodyBytes))
		} else if bodyStr, ok := bodyData.(string); ok {
			options = append(options, option.WithRequestBody("application/octet-stream", []byte(bodyStr)))
		} else {
			return nil, fmt.Errorf("Unsupported body for application/octet-stream: %v", bodyData)
		}

	default:
		panic("Invalid body content type!")
	}

	return options, nil
}
