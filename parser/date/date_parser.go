package date

import (
	"time"

	carbon "github.com/golang-module/carbon/v2"
	"github.com/nected/go-lib/parser/errors"
)

// var c carbon.Carbon

func init() {
	carbon.SetDefault(carbon.Default{
		Layout:       carbon.RFC3339Layout,
		Timezone:     carbon.UTC,
		WeekStartsAt: carbon.Monday,
		Locale:       "en", // value range: translate file name in the lang directory, excluding file suffix
	})
}

// newParseError creates a new ParseError.
// The provided value and valueElem are cloned to avoid escaping their values.
func newParseError(layout, value, layoutElem, valueElem, message string) *time.ParseError {
	valueCopy := cloneString(value)
	valueElemCopy := cloneString(valueElem)
	return &time.ParseError{Layout: layout, Value: valueCopy, LayoutElem: layoutElem, ValueElem: valueElemCopy, Message: message}
}

// cloneString returns a string copy of s.
// Do not use strings.Clone to avoid dependency on strings package.
func cloneString(s string) string {
	return string([]byte(s))
}

func Parse(input, format string) (time.Time, error) {
	val, err := time.Parse(format, input)
	if _, ok := err.(*time.ParseError); ok {
		return parseCustomTime(input, format)
	}
	return val, nil
}

func Format(t time.Time, format string) (string, error) {
	return formatCustomTime(t, format)
}

// parse date in following formats
func parseCustomTime(input, format string) (time.Time, error) {
	var c carbon.Carbon
	if len(input) == 0 {
		return time.Time{}, newParseError(format, input, "", "", errors.ErrEmptyInput.Error())
	}
	if len(format) == 0 || format == "" {
		c = carbon.Parse(input)
		if c.Error == nil {
			return c.StdTime(), nil
		}
	}

	input, layout := handleTimeString(input, format)

	c = carbon.ParseByLayout(input, layout)
	if c.Error != nil {
		return time.Time{}, newParseError(format, input, "", "", errors.ErrInvalidInput.Error())
	}

	return c.StdTime(), nil
}

func formatCustomTime(t time.Time, format string) (string, error) {
	if len(format) == 0 || format == "" {
		format = time.RFC3339
	}

	if result := t.Format(format); result != "" && result != format {
		return result, nil
	}

	timeStruct := NewTimeStruct()
	timeStruct.FromTime(t)
	return timeStruct.Format(format)
}
