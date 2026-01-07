package requestflag

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/goccy/go-yaml"
	"github.com/urfave/cli/v3"
)

// Flag [T] is a generic flag base which can be used to implement the most
// common interfaces used by urfave/cli. Additionally, it allows specifying
// where in an HTTP request the flag values should be placed (e.g. query, body, etc.).
type Flag[
	T []any | []map[string]any | []DateTimeValue | []DateValue | []TimeValue | []string |
		[]float64 | []int64 | []bool | any | map[string]any | DateTimeValue | DateValue | TimeValue |
		string | float64 | int64 | bool,
] struct {
	Name        string        // name of the flag
	Category    string        // category of the flag, if any
	DefaultText string        // default text of the flag for usage purposes
	HideDefault bool          // whether to hide the default value in output
	Usage       string        // usage string for help output
	Required    bool          // whether the flag is required or not
	Hidden      bool          // whether to hide the flag in help output
	Default     T             // default value for this flag if not set by from any source
	Aliases     []string      // aliases that are allowed for this flag
	Validator   func(T) error // custom function to validate this flag value

	QueryPath  string // location in the request query string to put this flag's value
	HeaderPath string // location in the request header to put this flag's value
	BodyPath   string // location in the request body to put this flag's value
	BodyRoot   bool   // if true, then use this value as the entire request body

	// unexported fields for internal use
	count      int       // number of times the flag has been set
	hasBeenSet bool      // whether the flag has been set from env or file
	applied    bool      // whether the flag has been applied to a flag set already
	value      cli.Value // value representing this flag's value
}

// Type assertions to verify we implement the relevant urfave/cli interfaces
var _ cli.CategorizableFlag = (*Flag[any])(nil)

// InRequest interface for flags that should be included in HTTP requests
type InRequest interface {
	GetQueryPath() string
	GetHeaderPath() string
	GetBodyPath() string
	IsBodyRoot() bool
}

func (f Flag[T]) GetQueryPath() string {
	return f.QueryPath
}

func (f Flag[T]) GetHeaderPath() string {
	return f.HeaderPath
}

func (f Flag[T]) GetBodyPath() string {
	return f.BodyPath
}

func (f Flag[T]) IsBodyRoot() bool {
	return f.BodyRoot
}

// The values that will be sent in different parts of a request.
type RequestContents struct {
	Queries map[string]any
	Headers map[string]any
	Body    any
}

// Extract query parameters, headers, and body values from command flags.
func ExtractRequestContents(cmd *cli.Command) RequestContents {
	bodyMap := make(map[string]any)
	res := RequestContents{
		Queries: make(map[string]any),
		Headers: make(map[string]any),
		Body:    bodyMap,
	}

	for _, flag := range cmd.Flags {
		if !flag.IsSet() {
			continue
		}

		value := flag.Get()
		if toSend, ok := flag.(InRequest); ok {
			if queryPath := toSend.GetQueryPath(); queryPath != "" {
				res.Queries[queryPath] = value
			}
			if headerPath := toSend.GetHeaderPath(); headerPath != "" {
				res.Headers[headerPath] = value
			}
			if toSend.IsBodyRoot() {
				res.Body = value
			} else if bodyPath := toSend.GetBodyPath(); bodyPath != "" {
				bodyMap[bodyPath] = value
			}
		}
	}
	return res
}

// Implementation of the cli.Flag interface
var _ cli.Flag = (*Flag[any])(nil) // Type assertion to ensure interface compliance

func (f *Flag[T]) PreParse() error {
	newVal := f.Default
	f.value = &cliValue[T]{newVal}

	// Validate the given default or values set from external sources as well
	if f.Validator != nil {
		if err := f.Validator(f.value.Get().(T)); err != nil {
			return err
		}
	}
	f.applied = true
	return nil
}

func (f *Flag[T]) PostParse() error {
	return nil
}

func (f *Flag[T]) Set(name string, val string) error {
	// Initialize flag if needed
	if !f.applied {
		if err := f.PreParse(); err != nil {
			return err
		}
		f.applied = true
	}

	f.count++

	// If this is the first time setting a slice type, reset it to empty
	// to avoid appending to the default value
	if f.count == 1 && f.value != nil {
		typ := reflect.TypeOf(f.Default)
		if typ != nil && typ.Kind() == reflect.Slice {
			// Create a new empty slice of the same type and set it
			emptySlice := reflect.MakeSlice(typ, 0, 0).Interface()
			f.value = &cliValue[T]{emptySlice.(T)}
		}
	}

	if err := f.value.Set(val); err != nil {
		return err
	}

	f.hasBeenSet = true

	if f.Validator != nil {
		if err := f.Validator(f.value.Get().(T)); err != nil {
			return err
		}
	}
	return nil
}

func (f *Flag[T]) Get() any {
	if f.value != nil {
		return f.value.Get()
	}
	return f.Default
}

func (f *Flag[T]) String() string {
	return cli.FlagStringer(f)
}

func (f *Flag[T]) IsSet() bool {
	return f.hasBeenSet
}

func (f *Flag[T]) Names() []string {
	return cli.FlagNames(f.Name, f.Aliases)
}

// Implementation for the cli.VisibleFlag interface
var _ cli.VisibleFlag = (*Flag[any])(nil) // Type assertion to ensure interface compliance

func (f *Flag[T]) IsVisible() bool {
	return !f.Hidden
}

func (f *Flag[T]) GetCategory() string {
	return f.Category
}

func (f *Flag[T]) SetCategory(c string) {
	f.Category = c
}

// Implementation for the cli.RequiredFlag interface
var _ cli.RequiredFlag = (*Flag[any])(nil) // Type assertion to ensure interface compliance

func (f *Flag[T]) IsRequired() bool {
	return f.Required
}

// Implementation for the cli.DocGenerationFlag interface
var _ cli.DocGenerationFlag = (*Flag[any])(nil) // Type assertion to ensure interface compliance

func (f *Flag[T]) TakesValue() bool {
	var t T
	return reflect.TypeOf(t) == nil || reflect.TypeOf(t).Kind() != reflect.Bool
}

func (f *Flag[T]) GetUsage() string {
	return f.Usage
}

func (f *Flag[T]) GetValue() string {
	if f.value == nil {
		return ""
	}
	return f.value.String()
}

func (f *Flag[T]) GetDefaultText() string {
	return f.DefaultText
}

func (f *Flag[T]) GetEnvVars() []string {
	return nil
}

func (f *Flag[T]) IsDefaultVisible() bool {
	return !f.HideDefault
}

func (f *Flag[T]) TypeName() string {
	ty := reflect.TypeOf(f.Default)
	if ty == nil {
		return ""
	}

	// Get base type name with special handling for built-in types
	getTypeName := func(t reflect.Type) string {
		switch t.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return "int"
		case reflect.Float32, reflect.Float64:
			return "float"
		case reflect.Bool:
			return "boolean"
		case reflect.String:
			switch t.Name() {
			case "DateTimeValue":
				return "datetime"
			case "DateValue":
				return "date"
			case "TimeValue":
				return "time"
			default:
				return "string"
			}
		default:
			if t.Name() == "" {
				return "any"
			}
			return strings.ToLower(t.Name())
		}
	}

	switch ty.Kind() {
	case reflect.Slice:
		elemType := ty.Elem()
		return getTypeName(elemType)
	case reflect.Map:
		keyType := ty.Key()
		valueType := ty.Elem()
		return fmt.Sprintf("%s=%s", getTypeName(keyType), getTypeName(valueType))
	default:
		return getTypeName(ty)
	}
}

// Implementation for the cli.DocGenerationMultiValueFlag interface
var _ cli.DocGenerationMultiValueFlag = (*Flag[any])(nil) // Type assertion to ensure interface compliance

func (f *Flag[T]) IsMultiValueFlag() bool {
	if reflect.TypeOf(f.Default) == nil {
		return false
	}
	kind := reflect.TypeOf(f.Default).Kind()
	return kind == reflect.Slice || kind == reflect.Map
}

func (f *Flag[T]) IsBoolFlag() bool {
	_, isBool := any(f.Default).(bool)
	return isBool
}

// Implementation for the cli.Countable interface
var _ cli.Countable = (*Flag[any])(nil) // Type assertion to ensure interface compliance

func (f *Flag[T]) Count() int {
	return f.count
}

// cliValue is a generic implementation of cli.Value for common types
type cliValue[
	T []any | []map[string]any | []DateTimeValue | []DateValue | []TimeValue | []string | []float64 |
		[]int64 | []bool | any | map[string]any | DateTimeValue | DateValue | TimeValue | string |
		float64 | int64 | bool,
] struct {
	value T
}

// Take an argument string for a single argument and convert it into a typed
// value for one of the supported CLI argument types
func parseCLIArg[
	T []any | []map[string]any | []DateTimeValue | []DateValue | []TimeValue | []string | []float64 |
		[]int64 | []bool | any | map[string]any | DateTimeValue | DateValue | TimeValue | string |
		float64 | int64 | bool,
](value string) (T, error) {
	var parsedValue any
	var err error

	var empty T
	switch any(empty).(type) {
	case string:
		parsedValue = value
	case int64:
		parsedValue, err = strconv.ParseInt(value, 0, 64)
	case float64:
		parsedValue, err = strconv.ParseFloat(value, 64)
	case bool:
		parsedValue, err = strconv.ParseBool(value)
	case DateTimeValue:
		var dt DateTimeValue
		err = (&dt).Parse(value)
		if err == nil {
			parsedValue = dt
		}

	case DateValue:
		var d DateValue
		err = (&d).Parse(value)
		if err == nil {
			parsedValue = d
		}

	case TimeValue:
		var t TimeValue
		err = (&t).Parse(value)
		if err == nil {
			parsedValue = t
		}

	default:
		var yamlValue T
		err = yaml.Unmarshal([]byte(value), &yamlValue)
		if err != nil {
			err = fmt.Errorf("failed to parse as YAML: %w", err)
		}
		parsedValue = yamlValue
	}

	// Nil needs to be handled specially because unmarshalling a YAML `null`
	// causes problems when doing type assertions.
	if parsedValue == nil {
		parsedValue = (*struct{})(nil)
	}

	if err == nil {
		if typedValue, ok := parsedValue.(T); ok {
			return typedValue, nil
		} else {
			expectedType := reflect.TypeFor[T]()
			err = fmt.Errorf("Couldn't convert %q (%v) to expected type %v", value, parsedValue, expectedType)
		}
	}
	return empty, err

}

// Parse the input string and set result as the cliValue's value
func (c *cliValue[T]) Set(value string) error {
	valueType := reflect.TypeOf(c.value)
	// When setting slice values, we append to the existing values
	// e.g. --foo 10 --foo 20 --foo 30 => [10, 20, 30]
	if valueType != nil && valueType.Kind() == reflect.Slice {
		elemType := valueType.Elem()

		var singleElem any
		var err error
		switch elemType.Kind() {
		case reflect.String:
			singleElem, err = parseCLIArg[string](value)
		case reflect.Int64:
			singleElem, err = parseCLIArg[int64](value)
		case reflect.Float64:
			singleElem, err = parseCLIArg[float64](value)
		case reflect.Bool:
			singleElem, err = parseCLIArg[bool](value)
		default:
			// Check for special types by name
			switch elemType.Name() {
			case "DateTimeValue":
				singleElem, err = parseCLIArg[DateTimeValue](value)
			case "DateValue":
				singleElem, err = parseCLIArg[DateValue](value)
			case "TimeValue":
				singleElem, err = parseCLIArg[TimeValue](value)
			default:
				// This handles []map[string]any
				if elemType.Kind() == reflect.Map && elemType.Key().Kind() == reflect.String {
					singleElem, err = parseCLIArg[map[string]any](value)
				} else {
					singleElem, err = parseCLIArg[any](value)
				}
			}
		}

		if err != nil {
			return err
		}

		// Append the new element to the slice
		sliceValue := reflect.ValueOf(c.value)
		if !sliceValue.IsValid() || sliceValue.IsNil() {
			// Create a new slice if the current one is nil
			sliceValue = reflect.MakeSlice(valueType, 0, 1)
		}

		// Append the new element
		newElem := reflect.ValueOf(singleElem)
		sliceValue = reflect.Append(sliceValue, newElem)

		// Set the updated slice back to c.value
		c.value = sliceValue.Interface().(T)
	} else {
		// For non-slice types, simply parse and set the value
		if parsedValue, err := parseCLIArg[T](value); err != nil {
			return err
		} else {
			c.value = parsedValue
		}
	}

	return nil
}

func (c *cliValue[T]) Get() any {
	return c.value
}

func (c *cliValue[T]) String() string {
	switch v := any(c.value).(type) {
	case string, int, int64, float64, bool, DateTimeValue, DateValue, TimeValue,
		[]string, []int, []int64, []float64, []bool, []DateTimeValue, []DateValue, []TimeValue:
		// For basic types, use standard string representation
		return fmt.Sprintf("%v", v)

	default:
		// For complex types, convert to YAML
		yamlBytes, err := yaml.MarshalWithOptions(c.value, yaml.Flow(true))
		if err != nil {
			// Fall back to standard format if YAML conversion fails
			return fmt.Sprintf("%v", c.value)
		}
		return string(yamlBytes)
	}
}

func (c *cliValue[T]) IsBoolFlag() bool {
	_, ok := any(c.value).(bool)
	return ok
}

// Time-related value types
type DateValue string
type DateTimeValue string
type TimeValue string

// String methods for time-related types
func (d DateValue) String() string {
	return string(d)
}

func (d DateTimeValue) String() string {
	return string(d)
}

func (t TimeValue) String() string {
	return string(t)
}

// parseTimeWithFormats attempts to parse a string using multiple formats
func parseTimeWithFormats(s string, formats []string) (time.Time, error) {
	var lastErr error
	for _, format := range formats {
		t, err := time.Parse(format, s)
		if err == nil {
			return t, nil
		}
		lastErr = err
	}
	return time.Time{}, lastErr
}

// Parse methods for time-related types
func (d *DateValue) Parse(s string) error {
	formats := []string{
		"2006-01-02",
		"01/02/2006",
		"Jan 2, 2006",
		"January 2, 2006",
		"2-Jan-2006",
	}

	t, err := parseTimeWithFormats(s, formats)
	if err != nil {
		return fmt.Errorf("unable to parse date: %v", err)
	}

	*d = DateValue(t.Format("2006-01-02"))
	return nil
}

func (d *DateTimeValue) Parse(s string) error {
	formats := []string{
		time.RFC3339,
		time.RFC3339Nano,
		"2006-01-02T15:04:05",
		"2006-01-02 15:04:05",
		time.RFC1123,
		time.RFC822,
		time.ANSIC,
	}

	t, err := parseTimeWithFormats(s, formats)
	if err != nil {
		return fmt.Errorf("unable to parse datetime: %v", err)
	}

	*d = DateTimeValue(t.Format(time.RFC3339))
	return nil
}

func (t *TimeValue) Parse(s string) error {
	formats := []string{
		"15:04:05",
		"15:04:05.999999999Z07:00",
		"3:04:05PM",
		"3:04 PM",
		"15:04",
		time.Kitchen,
	}

	parsedTime, err := parseTimeWithFormats(s, formats)
	if err != nil {
		return fmt.Errorf("unable to parse time: %v", err)
	}

	*t = TimeValue(parsedTime.Format("15:04:05"))
	return nil
}

// Allow setting inner fields on other flags (e.g. --foo.baz can set the "baz"
// field on the --foo flag)
type SettableInnerField interface {
	SetInnerField(string, any)
}

func (f *Flag[T]) SetInnerField(field string, val any) {
	if f.value == nil {
		f.value = &cliValue[T]{}
	}

	if settableInnerField, ok := f.value.(SettableInnerField); ok {
		settableInnerField.SetInnerField(field, val)
		f.hasBeenSet = true
	} else {
		panic(fmt.Sprintf("Cannot set inner field: %v", f.value))
	}
}

func (c *cliValue[T]) SetInnerField(field string, val any) {
	flagVal := c.value
	flagValReflect := reflect.ValueOf(flagVal)
	switch flagValReflect.Kind() {
	case reflect.Slice:
		if flagValReflect.Type().Elem().Kind() != reflect.Map {
			return
		}

		sliceLen := flagValReflect.Len()
		if sliceLen > 0 {
			// Check if the last element already has the InnerField
			lastElement := flagValReflect.Index(sliceLen - 1).Interface().(map[string]any)
			if _, hasInnerField := lastElement[field]; !hasInnerField {
				// Last element doesn't have the field, set it
				lastElement[field] = val
				return
			}
		}

		// Create a new map and append it to the slice
		newMap := map[string]any{field: val}
		switch sliceVal := any(c.value).(type) {
		case []map[string]any:
			c.value = any(append(sliceVal, newMap)).(T)
		case []any:
			c.value = any(append(sliceVal, newMap)).(T)
		}

	case reflect.Map:
		mapVal, ok := any(flagVal).(map[string]any)
		if !ok || mapVal == nil {
			mapVal = map[string]any{field: val}
			c.value = any(mapVal).(T)
		} else {
			mapVal[field] = val
		}
	}
}
