package errors

import "fmt"

var (
	ErrEmptyData   = fmt.Errorf("data is empty")
	ErrInvalidData = fmt.Errorf("data is invalid")
)
