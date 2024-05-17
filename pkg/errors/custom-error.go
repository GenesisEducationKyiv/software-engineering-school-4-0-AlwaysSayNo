package errors

import "fmt"

type customError struct {
	Message string
	Cause   error
}

func (e *customError) Error() string {
	if e.Cause != nil && len(e.Message) != 0 {
		return fmt.Sprintf("%s: %v", e.Message, e.Cause)
	} else if e.Cause != nil {
		return e.Cause.Error()
	}
	return e.Message
}

func (e *customError) Unwrap() error {
	return e.Cause
}
