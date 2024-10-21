package date

import (
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/nected/go-lib/parser/errors"
)

type TimeStruct struct {
	date     int
	month    int
	year     int
	hour     int
	minute   int
	second   int
	nanosec  int
	timeZone *time.Location
}

func NewTimeStruct() *TimeStruct {
	return &TimeStruct{}
}

// From Time
func (t *TimeStruct) FromTime(time time.Time) {
	t.date = time.Day()
	t.month = int(time.Month())
	t.year = time.Year()
	t.hour = time.Hour()
	t.minute = time.Minute()
	t.second = time.Second()
	t.nanosec = time.Nanosecond()
	t.timeZone = time.Location()
}

// Format
func (t *TimeStruct) Format(format string) (string, error) {
	var result strings.Builder
	var lastPos int
	for i := 0; i < len(format); i++ {
		if !unicode.IsLetter(rune(format[i])) && !unicode.IsNumber(rune(format[i])) {
			key := format[lastPos:i]
			if key == "" {
				result.WriteByte(format[i])
				lastPos = i + 1
				continue
			}
			val := t.getValue(key)
			if val == "" {
				return "", newParseError(format, val, key, val, errors.ErrInvalidDateFormat.Error())
			}
			result.WriteString(t.getValue(key))
			result.WriteByte(format[i])
			lastPos = i + 1
		}
	}
	if len(format) > lastPos {
		key := format[lastPos:]
		val := t.getValue(key)
		if val == "" {
			return "", newParseError(format, val, key, val, errors.ErrInvalidDateFormat.Error())
		}
		result.WriteString(t.getValue(key))
	}
	return result.String(), nil
}

func (t *TimeStruct) getValue(key string) string {
	switch key {
	case "d":
		return strconv.Itoa(t.date)
	case "dd":
		dateStr := strconv.Itoa(t.date)
		if len(dateStr) == 1 {
			dateStr = "0" + dateStr
		}
		return dateStr
	case "ddd", "DDD":
		return time.Weekday(time.Date(t.year, time.Month(t.month), t.date, t.hour, t.minute, t.second, t.nanosec, t.timeZone).Weekday()).String()[:3]
	case "dddd", "DDDD":
		return time.Weekday(time.Date(t.year, time.Month(t.month), t.date, t.hour, t.minute, t.second, t.nanosec, t.timeZone).Weekday()).String()
	case "m":
		return strconv.Itoa(t.month)
	case "mm":
		monthStr := strconv.Itoa(t.month)
		if len(monthStr) == 1 {
			monthStr = "0" + monthStr
		}
		return monthStr
	case "mmm":
		return time.Month(t.month).String()[:3]
	case "mmmm":
		return time.Month(t.month).String()
	case "y":
		return strconv.Itoa(t.year % 100)
	case "yy":
		yearStr := strconv.Itoa(t.year % 100)
		if len(yearStr) == 1 {
			yearStr = "0" + yearStr
		}
		return yearStr
	case "yyyy":
		yearStr := strconv.Itoa(t.year)
		prefix := ""
		for i := 0; i < 4-len(yearStr); i++ {
			prefix += "0"
		}
		return prefix + yearStr
	case "h":
		return strconv.Itoa(t.hour % 12)
	case "hh":
		hourStr := strconv.Itoa(t.hour % 12)
		if len(hourStr) == 1 {
			hourStr = "0" + hourStr
		}
		return hourStr
	case "H", "g":
		return strconv.Itoa(t.hour)
	case "HH":
		hourStr := strconv.Itoa(t.hour)
		if len(hourStr) == 1 {
			hourStr = "0" + hourStr
		}
		return hourStr
	case "M":
		return strconv.Itoa(t.minute)
	case "MM":
		minuteStr := strconv.Itoa(t.minute)
		if len(minuteStr) == 1 {
			minuteStr = "0" + minuteStr
		}
		return minuteStr
	case "s":
		return strconv.Itoa(t.second)
	case "ss":
		secondString := strconv.Itoa(t.second)
		if len(secondString) == 1 {
			secondString = "0" + secondString
		}
		return secondString
	case "ns":
		nsLength := 9
		ns := strconv.Itoa(t.nanosec)
		if len(ns) < nsLength {
			len := nsLength - len(ns)
			for i := 0; i < len; i++ {
				ns = "0" + ns
			}
		}
		return ns
	case "t":
		if t.hour < 12 {
			return "a"
		}
		return "p"
	case "T":
		if t.hour < 12 {
			return "A"
		}
		return "P"
	case "tt":
		if t.hour < 12 {
			return "am"
		}
		return "pm"
	case "TT":
		if t.hour < 12 {
			return "AM"
		}
		return "PM"
	case "Z":
		return t.timeZone.String()
	case "l": // millisconds three digits
		return strconv.Itoa(t.nanosec / 1000000)
	case "L": // milliseconds two digits
		return strconv.Itoa(t.nanosec / 10000000)
	case "o": // timezone offset as +HHMM or -HHMM
		format := time.Date(t.year, time.Month(t.month), t.date, t.hour, t.minute, t.second, t.nanosec, t.timeZone).Format("-0700")
		return format
	case "p": // timezone offset with colon as +HH:MM or -HH:MM
		format := time.Date(t.year, time.Month(t.month), t.date, t.hour, t.minute, t.second, t.nanosec, t.timeZone).Format("-07:00")
		return format
	default:
		return ""
	}
}
