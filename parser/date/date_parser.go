package date

import (
	"time"
	"unicode"
)

// newParseError creates a new ParseError.
// The provided value and valueElem are cloned to avoid escaping their values.
func newParseError(layout, value, layoutElem, valueElem, message string) *time.ParseError {
	valueCopy := cloneString(value)
	valueElemCopy := cloneString(valueElem)
	return &time.ParseError{layout, valueCopy, layoutElem, valueElemCopy, message}
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
//
// dd/mm/yyyy
// dd/mmm/yyyy
// yyyy-mm-dd
func parseCustomTime(input, format string) (time.Time, error) {
	if len(input) == 0 {
		return time.Time{}, newParseError(format, input, "", "", "empty input")
	}
	if len(input) != len(format) {
		return time.Time{}, newParseError(format, input, format, input, "input length does not match format")
	}
	timeStruct := NewTimeStruct()
	// var parsedString string
	var lastPos int
	for i := 0; i < len(format); i++ {
		if !unicode.IsLetter(rune(format[i])) && !unicode.IsNumber(rune(format[i])) {
			value := input[lastPos:i]
			key := format[lastPos:i]

			if err := timeStruct.Set(key, value); err != nil {
				return time.Time{}, newParseError(format, input, key, value, err.Error())
			}
			lastPos = i + 1
		}
	}
	if len(format) > lastPos {
		value := input[lastPos:]
		key := format[lastPos:]
		if err := timeStruct.Set(key, value); err != nil {
			return time.Time{}, newParseError(format, input, key, value, err.Error())
		}
	}
	return timeStruct.ToTime(), nil
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
