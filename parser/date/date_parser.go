package date

import (
	"bytes"
	"time"
	"unicode"

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
	// c = carbon.NewCarbon()
}

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
func parseCustomTime(input, format string) (time.Time, error) {
	if len(input) == 0 {
		return time.Time{}, newParseError(format, input, "", "", errors.ErrEmptyInput.Error())
	}
	c := carbon.NewCarbon()
	if len(format) == 0 || format == "" {
		c = carbon.Parse(input)
		if c.Error == nil {
			return c.StdTime(), nil
		}
	}

	layout := format2layout(format)

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
	// format = format2layout(format)

	if result := t.Format(format); result != "" && result != format {
		return result, nil
	}
	timeStruct := NewTimeStruct()
	timeStruct.FromTime(t)
	return timeStruct.Format(format)
}

// common formatting symbols
var formats = map[byte]string{
	'd': "02",      // Day:    Day of the month, 2 digits with leading zeros. Eg: 01 to 31.
	'D': "Mon",     // Day:    A textual representation of a day, three letters. Eg: Mon through Sun.
	'j': "2",       // Day:    Day of the month without leading zeros. Eg: 1 to 31.
	'l': "Monday",  // Day:    A full textual representation of the day of the week. Eg: Sunday through Saturday.
	'F': "January", // Month:  A full textual representation of a month, such as January or March. Eg: January through December.
	'm': "01",      // Month:  Numeric representation of a month, with leading zeros. Eg: 01 through 12.
	'M': "Jan",     // Month:  A short textual representation of a month, three letters. Eg: Jan through Dec.
	'n': "1",       // Month:  Numeric representation of a month, without leading zeros. Eg: 1 through 12.
	'Y': "2006",    // Year:   A full numeric representation of a year, 4 digits. Eg: 1999 or 2003.
	'y': "06",      // Year:   A two digit representation of a year. Eg: 99 or 03.
	'a': "pm",      // Time:   Lowercase morning or afternoon sign. Eg: am or pm.
	'A': "PM",      // Time:   Uppercase morning or afternoon sign. Eg: AM or PM.
	'g': "3",       // Time:   12-hour format of an hour without leading zeros. Eg: 1 through 12.
	'h': "03",      // Time:   12-hour format of an hour with leading zeros. Eg: 01 through 12.
	'H': "15",      // Time:   24-hour format of an hour with leading zeros. Eg: 00 through 23.
	'i': "04",      // Time:   Minutes with leading zeros. Eg: 00 to 59.
	's': "05",      // Time:   Seconds with leading zeros. Eg: 00 through 59.
	'O': "-0700",   // Zone:   Difference to Greenwich time (GMT) in hours. Eg: +0200.
	'P': "-07:00",  // Zone:   Difference to Greenwich time (GMT) with colon between hours and minutes. Eg: +02:00.
	'T': "MST",     // Zone:   Timezone abbreviation. Eg: UTC, EST, MDT ...

	'U': "timestamp",      // Timestamp with second. Eg: 1699677240.
	'V': "timestampMilli", // TimestampMilli with second. Eg: 1596604455666.
	'X': "timestampMicro", // TimestampMicro with second. Eg: 1596604455666666.
	'Z': "timestampNano",  // TimestampNano with second. Eg: 1596604455666666666.
}

var formatMap = map[string]byte{
	"d":    'j',
	"dd":   'd',
	"ddd":  'D',
	"DDD":  'D',
	"dddd": 'l',
	"DDDD": 'l',
	"m":    'n',
	"mm":   'm',
	"mmm":  'M',
	"mmmm": 'F',
	"yy":   'y',
	"yyyy": 'Y',
	"h":    'g',
	"hh":   'h',
	"H":    'H',
	"HH":   'H',
	"M":    'i',
	"MM":   'i',
	"o":    'O',
	"p":    'P',
	"s":    's',
	"ss":   's',
	"t":    'a',
	"tt":   'a',
	"T":    'A',
	"TT":   'A',
	"Z":    'T',
}

func formatMap2format(format string) string {
	buffer := bytes.NewBuffer(nil)
	lastPos := 0
	for i := 0; i < len(format); i++ {
		char := rune(format[i])
		if !unicode.IsLetter(char) && !unicode.IsNumber(char) {
			key := format[lastPos:i]
			if value, ok := formatMap[key]; ok {
				buffer.WriteByte(value)
			} else {
				buffer.WriteString(key)
			}
			separator := format[i]
			// if len(separator) > 0 {
			buffer.WriteByte(separator)
			// }
			lastPos = i + 1
		}
	}
	if len(format) > lastPos {
		key := format[lastPos:]
		if value, ok := formatMap[key]; ok {
			buffer.WriteByte(value)
		} else {
			buffer.WriteString(key)
		}
	}
	return buffer.String()
}

func format2layout(format string) string {
	format = formatMap2format(format)
	buffer := bytes.NewBuffer(nil)
	for i := 0; i < len(format); i++ {
		if layout, ok := formats[format[i]]; ok {
			buffer.WriteString(layout)
		} else {
			switch format[i] {
			case '\\': // raw output, no parse
				buffer.WriteByte(format[i+1])
				i++
				continue
			default:
				buffer.WriteByte(format[i])
			}
		}
	}
	return buffer.String()
}
