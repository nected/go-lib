package date

import (
	"strconv"
	"strings"
	"time"

	"github.com/nected/go-lib/parser/errors"
)

type TimeStruct struct {
	date      int
	month     int
	monthName string
	year      int
	hour      int
	minute    int
	second    int
	nanosec   int
}

func NewTimeStruct() *TimeStruct {
	return &TimeStruct{}
}

func (t *TimeStruct) Set(key string, value string) (err error) {
	switch key {
	case "dd":
		t.date, err = strconv.Atoi(value)
	case "MM":
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
	case "mm":
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
