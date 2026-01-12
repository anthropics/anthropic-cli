package jsonflag

import (
	"fmt"
	"strconv"
	"time"

	"github.com/urfave/cli/v3"
)

type JSONConfig struct {
	Kind MutationKind
	Path string
	// For boolean flags that set a specific value when present
	SetValue any
}

type JSONValueCreator[T any] struct{}

func (c JSONValueCreator[T]) Create(val T, dest *T, config JSONConfig) cli.Value {
	*dest = val
	return &jsonValue[T]{
		destination: dest,
		config:      config,
	}
}

func (c JSONValueCreator[T]) ToString(val T) string {
	switch v := any(val).(type) {
	case string:
		if v == "" {
			return v
		}
		return fmt.Sprintf("%q", v)
	case bool:
		return strconv.FormatBool(v)
	case int:
		return strconv.Itoa(v)
	case float64:
		return strconv.FormatFloat(v, 'g', -1, 64)
	case time.Time:
		return v.Format(time.RFC3339)
	default:
		return fmt.Sprintf("%v", v)
	}
}

type jsonValue[T any] struct {
	destination *T
	config      JSONConfig
}

func (v *jsonValue[T]) Set(val string) error {
	var parsed T
	var err error

	// If SetValue is configured, use that value instead of parsing the input
	if v.config.SetValue != nil {
		// For boolean flags with SetValue, register the configured value
		if _, isBool := any(parsed).(bool); isBool {
			globalRegistry.Mutate(v.config.Kind, v.config.Path, v.config.SetValue)
			*v.destination = any(true).(T) // Set the flag itself to true
			return nil
		}
		// For any flags with SetValue, register the configured value
		globalRegistry.Mutate(v.config.Kind, v.config.Path, v.config.SetValue)
		*v.destination = any(v.config.SetValue).(T)
		return nil
	}

	switch any(parsed).(type) {
	case string:
		parsed = any(val).(T)
	case bool:
		boolVal, parseErr := strconv.ParseBool(val)
		if parseErr != nil {
			return fmt.Errorf("invalid boolean value %q: %w", val, parseErr)
		}
		parsed = any(boolVal).(T)
	case int:
		intVal, parseErr := strconv.Atoi(val)
		if parseErr != nil {
			return fmt.Errorf("invalid integer value %q: %w", val, parseErr)
		}
		parsed = any(intVal).(T)
	case float64:
		floatVal, parseErr := strconv.ParseFloat(val, 64)
		if parseErr != nil {
			return fmt.Errorf("invalid float value %q: %w", val, parseErr)
		}
		parsed = any(floatVal).(T)
	case time.Time:
		// Try common datetime formats
		formats := []string{
			time.RFC3339,
			"2006-01-02T15:04:05Z07:00",
			"2006-01-02T15:04:05",
			"2006-01-02 15:04:05",
			"2006-01-02",
			"15:04:05",
			"15:04",
		}
		var timeVal time.Time
		var parseErr error
		for _, format := range formats {
			timeVal, parseErr = time.Parse(format, val)
			if parseErr == nil {
				break
			}
		}
		if parseErr != nil {
			return fmt.Errorf("invalid datetime value %q: %w", val, parseErr)
		}
		parsed = any(timeVal).(T)
	case any:
		// For `any`, store the string value directly
		parsed = any(val).(T)
	default:
		return fmt.Errorf("unsupported type for JSON flag")
	}

	*v.destination = parsed
	globalRegistry.Mutate(v.config.Kind, v.config.Path, parsed)
	return err
}

func (v *jsonValue[T]) Get() any {
	if v.destination != nil {
		return *v.destination
	}
	var zero T
	return zero
}

func (v *jsonValue[T]) String() string {
	if v.destination != nil {
		switch val := any(*v.destination).(type) {
		case string:
			return val
		case bool:
			return strconv.FormatBool(val)
		case int:
			return strconv.Itoa(val)
		case float64:
			return strconv.FormatFloat(val, 'g', -1, 64)
		case time.Time:
			return val.Format(time.RFC3339)
		default:
			return fmt.Sprintf("%v", val)
		}
	}
	var zero T
	switch any(zero).(type) {
	case string:
		return ""
	case bool:
		return "false"
	case int:
		return "0"
	case float64:
		return "0"
	case time.Time:
		return ""
	default:
		return fmt.Sprintf("%v", zero)
	}
}

func (v *jsonValue[T]) IsBoolFlag() bool {
	return v.config.SetValue != nil
}

// JSONDateValueCreator is a specialized creator for date-only values
type JSONDateValueCreator struct{}

func (c JSONDateValueCreator) Create(val time.Time, dest *time.Time, config JSONConfig) cli.Value {
	*dest = val
	return &jsonDateValue{
		destination: dest,
		config:      config,
	}
}

func (c JSONDateValueCreator) ToString(val time.Time) string {
	return val.Format("2006-01-02")
}

type jsonDateValue struct {
	destination *time.Time
	config      JSONConfig
}

func (v *jsonDateValue) Set(val string) error {
	// Try date-only formats first, then fall back to datetime formats
	formats := []string{
		"2006-01-02",
		"01/02/2006",
		"Jan 2, 2006",
		"January 2, 2006",
		"2-Jan-2006",
		time.RFC3339,
		"2006-01-02T15:04:05Z07:00",
		"2006-01-02T15:04:05",
		"2006-01-02 15:04:05",
	}

	var timeVal time.Time
	var parseErr error
	for _, format := range formats {
		timeVal, parseErr = time.Parse(format, val)
		if parseErr == nil {
			break
		}
	}
	if parseErr != nil {
		return fmt.Errorf("invalid date value %q: %w", val, parseErr)
	}

	*v.destination = timeVal
	globalRegistry.Mutate(v.config.Kind, v.config.Path, timeVal.Format("2006-01-02"))
	return nil
}

func (v *jsonDateValue) Get() any {
	if v.destination != nil {
		return *v.destination
	}
	return time.Time{}
}

func (v *jsonDateValue) String() string {
	if v.destination != nil {
		return v.destination.Format("2006-01-02")
	}
	return ""
}

func (v *jsonDateValue) IsBoolFlag() bool {
	return false
}

type JSONStringFlag = cli.FlagBase[string, JSONConfig, JSONValueCreator[string]]
type JSONBoolFlag = cli.FlagBase[bool, JSONConfig, JSONValueCreator[bool]]
type JSONIntFlag = cli.FlagBase[int, JSONConfig, JSONValueCreator[int]]
type JSONFloatFlag = cli.FlagBase[float64, JSONConfig, JSONValueCreator[float64]]
type JSONDatetimeFlag = cli.FlagBase[time.Time, JSONConfig, JSONValueCreator[time.Time]]
type JSONDateFlag = cli.FlagBase[time.Time, JSONConfig, JSONDateValueCreator]
type JSONAnyFlag = cli.FlagBase[any, JSONConfig, JSONValueCreator[any]]
