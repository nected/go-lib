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

func (t *TimeStruct) Set(key string, value string) (err error) {
	switch key {
	case "dd":
		t.date, err = strconv.Atoi(value)
	case "mm":
		t.month, err = strconv.Atoi(value)
	case "mmm":
		month, err := parseMonthName(value)
		if err != nil {
			return err
		}
		t.month = int(month)
	case "yy", "yyyy":
		t.year, err = strconv.Atoi(value)
	case "hh", "HH":
		t.hour, err = strconv.Atoi(value)
	case "MM":
		t.minute, err = strconv.Atoi(value)
	case "ss":
		t.second, err = strconv.Atoi(value)
	case "ns":
		t.nanosec, err = strconv.Atoi(value)
	default:
		err = errors.ErrInvalidDateFormat
	}
	return
}

// To Time
func (t *TimeStruct) ToTime() time.Time {
	return time.Date(t.year, time.Month(t.month), t.date, t.hour, t.minute, t.second, t.nanosec, time.UTC)
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
	case "ddd":
		return time.Weekday(time.Date(t.year, time.Month(t.month), t.date, t.hour, t.minute, t.second, t.nanosec, t.timeZone).Weekday()).String()[:3]
	case "DDD", "dddd":
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
	case "H":
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
		return strconv.Itoa(t.nanosec)
	default:
		return ""
	}
}

// function to parse month name to time.Month
// all comparison should be case insensitive
// returns parser error if month name is invalid
// example: Jan -> 1, January -> 1
func parseMonthName(monthName string) (time.Month, error) {
	monthName = strings.ToLower(monthName)
	switch monthName {
	case "jan", "january":
		return time.January, nil
	case "feb", "february":
		return time.February, nil
	case "mar", "march":
		return time.March, nil
	case "apr", "april":
		return time.April, nil
	case "may":
		return time.May, nil
	case "jun", "june":
		return time.June, nil
	case "jul", "july":
		return time.July, nil
	case "aug", "august":
		return time.August, nil
	case "sep", "september":
		return time.September, nil
	case "oct", "october":
		return time.October, nil
	case "nov", "november":
		return time.November, nil
	case "dec", "december":
		return time.December, nil
	default:
		return 0, errors.ErrInvalidDateFormat
	}
}
