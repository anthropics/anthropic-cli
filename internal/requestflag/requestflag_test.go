package requestflag

import (
	"fmt"
	"testing"
	"time"

	"github.com/goccy/go-yaml"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v3"
)

func TestDateValueParse(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{
			name:    "ISO format",
			input:   "2023-05-15",
			want:    "2023-05-15",
			wantErr: false,
		},
		{
			name:    "US format",
			input:   "05/15/2023",
			want:    "2023-05-15",
			wantErr: false,
		},
		{
			name:    "Short month format",
			input:   "May 15, 2023",
			want:    "2023-05-15",
			wantErr: false,
		},
		{
			name:    "Long month format",
			input:   "January 15, 2023",
			want:    "2023-01-15",
			wantErr: false,
		},
		{
			name:    "British format",
			input:   "15-Jan-2023",
			want:    "2023-01-15",
			wantErr: false,
		},
		{
			name:    "Invalid format",
			input:   "not a date",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var d DateValue
			err := d.Parse(tt.input)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, d.String())
			}
		})
	}
}

func TestDateTimeValueParse(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "RFC3339",
			input:   "2023-05-15T14:30:45Z",
			wantErr: false,
		},
		{
			name:    "ISO with timezone",
			input:   "2023-05-15T14:30:45+02:00",
			wantErr: false,
		},
		{
			name:    "ISO without timezone",
			input:   "2023-05-15T14:30:45",
			wantErr: false,
		},
		{
			name:    "Space separated",
			input:   "2023-05-15 14:30:45",
			wantErr: false,
		},
		{
			name:    "RFC1123",
			input:   "Mon, 15 May 2023 14:30:45 GMT",
			wantErr: false,
		},
		{
			name:    "RFC822",
			input:   "15 May 23 14:30 GMT",
			wantErr: false,
		},
		{
			name:    "ANSIC",
			input:   "Mon Jan 2 15:04:05 2006",
			wantErr: false,
		},
		{
			name:    "Invalid format",
			input:   "not a datetime",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var d DateTimeValue
			err := d.Parse(tt.input)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)

				// Parse the string back to ensure it's valid RFC3339
				_, parseErr := time.Parse(time.RFC3339, d.String())
				assert.NoError(t, parseErr)
			}
		})
	}
}

func TestTimeValueParse(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{
			name:    "24-hour format",
			input:   "14:30:45",
			want:    "14:30:45",
			wantErr: false,
		},
		{
			name:    "12-hour format with seconds",
			input:   "2:30:45PM",
			want:    "14:30:45",
			wantErr: false,
		},
		{
			name:    "12-hour format without seconds",
			input:   "2:30 PM",
			want:    "14:30:00",
			wantErr: false,
		},
		{
			name:    "24-hour without seconds",
			input:   "14:30",
			want:    "14:30:00",
			wantErr: false,
		},
		{
			name:    "Kitchen format",
			input:   "2:30PM",
			want:    "14:30:00",
			wantErr: false,
		},
		{
			name:    "Invalid format",
			input:   "not a time",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var tv TimeValue
			err := tv.Parse(tt.input)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, tv.String())
			}
		})
	}
}

func TestRequestParams(t *testing.T) {
	t.Run("map body type", func(t *testing.T) {
		// Create a mock command with flags
		cmd := &cli.Command{
			Name: "test",
		}

		// Create string flag with body path
		stringFlag := &Flag[string]{
			Name:       "string-flag",
			Default:    "default-string",
			BodyPath:   "string_field",
			value:      &cliValue[string]{value: "test-value"},
			hasBeenSet: true,
		}

		// Create int flag with header path
		intFlag := &Flag[int64]{
			Name:       "int-flag",
			Default:    42,
			HeaderPath: "X-Int-Value",
			value:      &cliValue[int64]{value: 99},
			hasBeenSet: true,
		}

		// Create bool flag with query path
		boolFlag := &Flag[bool]{
			Name:       "bool-flag",
			Default:    false,
			QueryPath:  "include_details",
			value:      &cliValue[bool]{value: true},
			hasBeenSet: true,
		}

		// Create date flag with multiple paths
		dateFlag := &Flag[DateValue]{
			Name:       "date-flag",
			Default:    DateValue("2023-01-01"),
			BodyPath:   "effective_date",
			HeaderPath: "X-Effective-Date",
			QueryPath:  "as_of_date",
			value:      &cliValue[DateValue]{value: DateValue("2023-05-15")},
			hasBeenSet: true,
		}

		// Create flag with no path
		noPathFlag := &Flag[string]{
			Name:       "no-path-flag",
			Default:    "no-path",
			value:      &cliValue[string]{value: "no-path-value"},
			hasBeenSet: true,
		}

		// Create unset flag
		unsetFlag := &Flag[string]{
			Name:       "unset-flag",
			Default:    "unset",
			BodyPath:   "should_not_appear",
			value:      &cliValue[string]{value: "unset-value"},
			hasBeenSet: false,
		}

		cmd.Flags = []cli.Flag{stringFlag, intFlag, boolFlag, dateFlag, noPathFlag, unsetFlag}

		// Test the RequestParams function
		contents := ExtractRequestContents(cmd)

		// Verify query parameters
		assert.Equal(t, true, contents.Queries["include_details"])
		assert.Equal(t, DateValue("2023-05-15"), contents.Queries["as_of_date"])
		assert.Len(t, contents.Queries, 2)

		// Verify headers
		assert.Equal(t, int64(99), contents.Headers["X-Int-Value"])
		assert.Equal(t, DateValue("2023-05-15"), contents.Headers["X-Effective-Date"])
		assert.Len(t, contents.Headers, 2)

		// Verify body
		bodyMap, ok := contents.Body.(map[string]any)
		assert.True(t, ok, "Expected body to be map[string]any, got %T", contents.Body)
		assert.Equal(t, "test-value", bodyMap["string_field"])
		assert.Equal(t, DateValue("2023-05-15"), bodyMap["effective_date"])
		assert.Len(t, bodyMap, 2)

		// Verify the unset flag didn't make it into the maps
		assert.NotContains(t, contents.Body, "should_not_appear")
	})

	t.Run("non-map body type", func(t *testing.T) {
		// Create a mock command with flags
		cmd := &cli.Command{
			Name: "test",
			Flags: []cli.Flag{
				&Flag[int64]{
					Name:     "int-body-flag",
					Default:  0,
					BodyRoot: true,
				},
			},
		}
		cmd.Set("int-body-flag", "42")

		contents := ExtractRequestContents(cmd)
		intBody, ok := contents.Body.(int64)
		assert.True(t, ok, "Expected body to be int64, got %T", contents.Body)
		assert.Equal(t, int64(42), intBody)
	})
}

func TestFlagSet(t *testing.T) {
	strFlag := &Flag[string]{
		Name:    "string-flag",
		Default: "default-string",
	}

	superstitiousIntFlag := &Flag[int64]{
		Name:    "int-flag",
		Default: 42,
		Validator: func(val int64) error {
			if val == 13 {
				return fmt.Errorf("Unlucky number!")
			}
			return nil
		},
	}

	boolFlag := &Flag[bool]{
		Name:    "bool-flag",
		Default: false,
	}

	// Test initialization and setting
	t.Run("PreParse initialization", func(t *testing.T) {
		assert.NoError(t, strFlag.PreParse())
		assert.True(t, strFlag.applied)
		assert.Equal(t, "default-string", strFlag.Get())
	})

	t.Run("Set string flag", func(t *testing.T) {
		assert.NoError(t, strFlag.Set("string-flag", "new-value"))
		assert.Equal(t, "new-value", strFlag.Get())
		assert.True(t, strFlag.IsSet())
	})

	t.Run("Set int flag with valid value", func(t *testing.T) {
		assert.NoError(t, superstitiousIntFlag.Set("int-flag", "100"))
		assert.Equal(t, int64(100), superstitiousIntFlag.Get())
		assert.True(t, superstitiousIntFlag.IsSet())
	})

	t.Run("Set int flag with invalid value", func(t *testing.T) {
		assert.Error(t, superstitiousIntFlag.Set("int-flag", "not-an-int"))
	})

	t.Run("Set int flag with validator failing", func(t *testing.T) {
		assert.Error(t, superstitiousIntFlag.Set("int-flag", "13"))
	})

	t.Run("Set bool flag", func(t *testing.T) {
		assert.NoError(t, boolFlag.Set("bool-flag", "true"))
		assert.Equal(t, true, boolFlag.Get())
		assert.True(t, boolFlag.IsSet())
	})

	t.Run("Set slice flag with multiple values", func(t *testing.T) {
		sliceFlag := &Flag[[]int64]{
			Name:    "slice-flag",
			Default: []int64{},
		}

		// Initialize the flag
		assert.NoError(t, sliceFlag.PreParse())

		// First set
		assert.NoError(t, sliceFlag.Set("slice-flag", "10"))

		// Subsequent setting should append, not replace
		assert.NoError(t, sliceFlag.Set("slice-flag", "20"))
		assert.NoError(t, sliceFlag.Set("slice-flag", "30"))

		// Verify that we have both values in the slice
		result := sliceFlag.Get()
		assert.Equal(t, []int64{10, 20, 30}, result)
		assert.True(t, sliceFlag.IsSet())
	})

	t.Run("Set slice flag with a nonempty default", func(t *testing.T) {
		sliceFlag := &Flag[[]int64]{
			Name:    "slice-flag",
			Default: []int64{99, 100},
		}

		assert.NoError(t, sliceFlag.PreParse())
		assert.NoError(t, sliceFlag.Set("slice-flag", "10"))
		assert.NoError(t, sliceFlag.Set("slice-flag", "20"))
		assert.NoError(t, sliceFlag.Set("slice-flag", "30"))

		// Verify that we have clobbered the default value instead of appending
		// to it.
		result := sliceFlag.Get()
		assert.Equal(t, []int64{10, 20, 30}, result)
		assert.True(t, sliceFlag.IsSet())
	})
}

func TestParseTimeWithFormats(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		formats  []string
		wantTime time.Time
		wantErr  bool
	}{
		{
			name:     "RFC3339 format",
			input:    "2023-05-15T14:30:45Z",
			formats:  []string{time.RFC3339},
			wantTime: time.Date(2023, 5, 15, 14, 30, 45, 0, time.UTC),
			wantErr:  false,
		},
		{
			name:     "Multiple formats - first matches",
			input:    "2023-05-15",
			formats:  []string{"2006-01-02", time.RFC3339},
			wantTime: time.Date(2023, 5, 15, 0, 0, 0, 0, time.UTC),
			wantErr:  false,
		},
		{
			name:     "Multiple formats - second matches",
			input:    "15/05/2023",
			formats:  []string{"2006-01-02", "02/01/2006"},
			wantTime: time.Date(2023, 5, 15, 0, 0, 0, 0, time.UTC),
			wantErr:  false,
		},
		{
			name:     "No matching format",
			input:    "not a date",
			formats:  []string{"2006-01-02", time.RFC3339},
			wantTime: time.Time{},
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseTimeWithFormats(tt.input, tt.formats)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.True(t, tt.wantTime.Equal(got), "Expected %v, got %v", tt.wantTime, got)
			}
		})
	}
}

func TestYamlHandling(t *testing.T) {
	// Test with any value
	t.Run("Parse YAML to any", func(t *testing.T) {
		cv := &cliValue[any]{}
		err := cv.Set("name: test\nvalue: 42\n")
		assert.NoError(t, err)

		// The value should be a map
		val, ok := cv.Get().(map[string]any)
		assert.True(t, ok, "Expected map[string]any, got %T", cv.Get())

		if ok {
			assert.Equal(t, "test", val["name"])
			assert.Equal(t, uint64(42), val["value"])
		}

		// The string representation should be valid YAML
		strVal := cv.String()
		var parsed map[string]any
		err = yaml.Unmarshal([]byte(strVal), &parsed)
		assert.NoError(t, err)
		assert.Equal(t, "test", parsed["name"])
		assert.Equal(t, uint64(42), parsed["value"])
	})

	// Test with array
	t.Run("Parse YAML array", func(t *testing.T) {
		cv := &cliValue[any]{}
		err := cv.Set("- item1\n- item2\n- item3\n")
		assert.NoError(t, err)

		// The value should be a slice
		val, ok := cv.Get().([]any)
		assert.True(t, ok, "Expected []any, got %T", cv.Get())

		if ok {
			assert.Len(t, val, 3)
			assert.Equal(t, "item1", val[0])
			assert.Equal(t, "item2", val[1])
			assert.Equal(t, "item3", val[2])
		}
	})

	// Test with invalid YAML
	t.Run("Parse invalid YAML", func(t *testing.T) {
		invalidYaml := `[not closed`
		cv := &cliValue[any]{}
		err := cv.Set(invalidYaml)
		assert.Error(t, err)
	})
}

func TestFlagTypeNames(t *testing.T) {
	tests := []struct {
		name     string
		flag     cli.DocGenerationFlag
		expected string
	}{
		{"string", &Flag[string]{}, "string"},
		{"int64", &Flag[int64]{}, "int"},
		{"float64", &Flag[float64]{}, "float"},
		{"bool", &Flag[bool]{}, "boolean"},
		{"string slice", &Flag[[]string]{}, "string"},
		{"date", &Flag[DateValue]{}, "date"},
		{"datetime", &Flag[DateTimeValue]{}, "datetime"},
		{"time", &Flag[TimeValue]{}, "time"},
		{"date slice", &Flag[[]DateValue]{}, "date"},
		{"datetime slice", &Flag[[]DateTimeValue]{}, "datetime"},
		{"time slice", &Flag[[]TimeValue]{}, "time"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			typeName := tt.flag.TypeName()
			assert.Equal(t, tt.expected, typeName, "Expected type name %q, got %q", tt.expected, typeName)
		})
	}
}
