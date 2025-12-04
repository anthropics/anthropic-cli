package requestflag

import (
	"time"

	"github.com/goccy/go-yaml"
	"github.com/urfave/cli/v3"
)

type (
	YAMLFlag      = cli.FlagBase[requestValue[any], RequestConfig, requestValueCreator[any]]
	YAMLSlice     = cli.SliceBase[requestValue[any], RequestConfig, requestValueCreator[any]]
	YAMLSliceFlag = cli.FlagBase[[]requestValue[any], RequestConfig, YAMLSlice]

	StringFlag      = cli.FlagBase[requestValue[string], RequestConfig, requestValueCreator[string]]
	StringSlice     = cli.SliceBase[requestValue[string], RequestConfig, requestValueCreator[string]]
	StringSliceFlag = cli.FlagBase[[]requestValue[string], RequestConfig, StringSlice]

	IntFlag      = cli.FlagBase[requestValue[int64], RequestConfig, requestValueCreator[int64]]
	IntSlice     = cli.SliceBase[requestValue[int64], RequestConfig, requestValueCreator[int64]]
	IntSliceFlag = cli.FlagBase[[]requestValue[int64], RequestConfig, IntSlice]

	DateFlag      = cli.FlagBase[requestValue[string], RequestConfig, dateCreator]
	DateSlice     = cli.SliceBase[requestValue[string], RequestConfig, dateCreator]
	DateSliceFlag = cli.FlagBase[[]requestValue[string], RequestConfig, DateSlice]

	TimeFlag      = cli.FlagBase[requestValue[string], RequestConfig, timeCreator]
	TimeSlice     = cli.SliceBase[requestValue[string], RequestConfig, timeCreator]
	TimeSliceFlag = cli.FlagBase[[]requestValue[string], RequestConfig, TimeSlice]

	DateTimeFlag      = cli.FlagBase[requestValue[string], RequestConfig, dateTimeCreator]
	DateTimeSlice     = cli.SliceBase[requestValue[string], RequestConfig, dateTimeCreator]
	DateTimeSliceFlag = cli.FlagBase[[]requestValue[string], RequestConfig, DateTimeSlice]

	FloatFlag      = cli.FlagBase[requestValue[float64], RequestConfig, requestValueCreator[float64]]
	FloatSlice     = cli.SliceBase[requestValue[float64], RequestConfig, requestValueCreator[float64]]
	FloatSliceFlag = cli.FlagBase[[]requestValue[float64], RequestConfig, FloatSlice]

	BoolFlag      = cli.FlagBase[requestValue[bool], RequestConfig, requestValueCreator[bool]]
	BoolSlice     = cli.SliceBase[requestValue[bool], RequestConfig, requestValueCreator[bool]]
	BoolSliceFlag = cli.FlagBase[[]requestValue[bool], RequestConfig, BoolSlice]
)

type RequestConfig struct {
	BodyPath   string
	HeaderPath string
	QueryPath  string
	CookiePath string
}

type RequestValue interface {
	RequestConfig() RequestConfig
	RequestValue() any
}

type requestValue[T any | string | int64 | float64 | bool] struct {
	value  T
	config RequestConfig
}

func Value[T any | string | int64 | float64 | bool](val T) requestValue[T] {
	return requestValue[T]{val, RequestConfig{}}
}

func (s requestValue[T]) RequestConfig() RequestConfig {
	return s.config
}

func (s requestValue[T]) RequestValue() any {
	return s.value
}

func CommandRequestValue[T any | string | int64 | float64 | bool](cmd *cli.Command, name string) T {
	r := cmd.Value(name).(requestValue[T])
	return r.value
}

func CommandRequestValues[T any | string | int64 | float64 | bool](cmd *cli.Command, name string) []T {
	rs := cmd.Value(name).([]requestValue[T])
	values := make([]T, len(rs))
	for i, r := range rs {
		values[i] = r.value
	}
	return values
}

func CollectRequestValues(rs []RequestValue) []any {
	values := make([]any, len(rs))
	for i, r := range rs {
		values[i] = r.RequestValue()
	}
	return values
}

type requestValueCreator[T any | string | int64 | float64 | bool] struct {
	destination *requestValue[T]
}

func (s requestValueCreator[T]) Create(defaultValue requestValue[T], p *requestValue[T], c RequestConfig) cli.Value {
	p.value = defaultValue.value
	p.config = c
	return &requestValueCreator[T]{
		destination: p,
	}
}

func (s requestValueCreator[T]) ToString(val requestValue[T]) string {
	data, err := yaml.Marshal(val)
	if err != nil {
		return ""
	}
	return string(data)
}

func (s *requestValueCreator[T]) Set(str string) error {
	var isStringType bool
	var zeroVal T
	_, isStringType = any(zeroVal).(string)
	if isStringType {
		s.destination.value = any(str).(T)
	} else {
		var val T
		if err := yaml.Unmarshal([]byte(str), &val); err != nil {
			return err
		}
		s.destination.value = val
	}
	return nil
}

func (s *requestValueCreator[T]) Get() any {
	return *s.destination
}

func (s *requestValueCreator[T]) String() string {
	if s.destination != nil {
		return s.ToString(*s.destination)
	}
	return ""
}

func (s requestValueCreator[T]) IsBoolFlag() bool {
	var zero T
	_, ok := any(zero).(bool)
	return ok
}

func parseTimeWithLayouts(str string, layouts []string) (time.Time, error) {
	var t time.Time
	var err error
	for _, layout := range layouts {
		t, err = time.Parse(layout, str)
		if err == nil {
			break
		}
	}
	return t, err
}

// Value creator for date, time, and datetime types
type timeFormatCreator struct {
	requestValueCreator[string]
	inputFormats []string
	outputFormat string
}

func (s timeFormatCreator) Create(defaultValue requestValue[string], p *requestValue[string], c RequestConfig) cli.Value {
	*p = defaultValue
	p.config = c
	return &timeFormatCreator{
		requestValueCreator: requestValueCreator[string]{
			destination: p,
		},
		inputFormats: s.inputFormats,
		outputFormat: s.outputFormat,
	}
}

func (s *timeFormatCreator) Set(str string) error {
	t, err := parseTimeWithLayouts(str, s.inputFormats)
	s.destination.value = t.Format(s.outputFormat)
	return err
}

type dateCreator struct {
	timeFormatCreator
}

func (d dateCreator) Create(defaultValue requestValue[string], p *requestValue[string], c RequestConfig) cli.Value {
	return timeFormatCreator{
		requestValueCreator: requestValueCreator[string]{
			destination: p,
		},
		inputFormats: []string{
			"2006-01-02",
			"01/02/2006",
			"Jan 2, 2006",
			"January 2, 2006",
			"2-Jan-2006",
		},
		outputFormat: "2006-01-02",
	}.Create(defaultValue, p, c)
}

type timeCreator struct {
	timeFormatCreator
}

func (t timeCreator) Create(defaultValue requestValue[string], p *requestValue[string], c RequestConfig) cli.Value {
	return timeFormatCreator{
		requestValueCreator: requestValueCreator[string]{
			destination: p,
		},
		inputFormats: []string{
			"15:04:05",
			"3:04:05PM",
			"3:04 PM",
			"15:04",
			time.Kitchen,
		},
		outputFormat: "15:04:05",
	}.Create(defaultValue, p, c)
}

type dateTimeCreator struct {
	timeFormatCreator
}

func (dt dateTimeCreator) Create(defaultValue requestValue[string], p *requestValue[string], c RequestConfig) cli.Value {
	return timeFormatCreator{
		requestValueCreator: requestValueCreator[string]{
			destination: p,
		},
		inputFormats: []string{
			time.RFC3339,
			"2006-01-02T15:04:05Z07:00",
			"2006-01-02T15:04:05",
			"2006-01-02 15:04:05",
			time.RFC1123,
			time.RFC822,
			time.ANSIC,
		},
		outputFormat: time.RFC3339,
	}.Create(defaultValue, p, c)
}
