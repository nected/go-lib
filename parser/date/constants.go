package date

import (
	"bytes"
	"unicode"
)

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

var goFormatMap = map[string]byte{
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
			if value, ok := goFormatMap[key]; ok {
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
		if value, ok := goFormatMap[key]; ok {
			buffer.WriteByte(value)
		} else {
			buffer.WriteString(key)
		}
	}
	return buffer.String()
}

func convertToGoFormat(format string) string {
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
