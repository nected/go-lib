package date

import (
	"strconv"
	"time"
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
		t.monthName = value
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
	}
	return
}

// To Time
func (t *TimeStruct) ToTime() time.Time {
	return time.Date(t.year, time.Month(t.month), t.date, t.hour, t.minute, t.second, t.nanosec, time.UTC)
}
