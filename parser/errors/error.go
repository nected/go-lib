package errors

import "fmt"

var (
	ErrInvalidInput      = fmt.Errorf("invalid input")
	ErrInvalidDateFormat = fmt.Errorf("invalid date format")
	ErrInvalidTimeZone   = fmt.Errorf("invalid time zone")
	ErrEmptyInput        = fmt.Errorf("empty input")
)
