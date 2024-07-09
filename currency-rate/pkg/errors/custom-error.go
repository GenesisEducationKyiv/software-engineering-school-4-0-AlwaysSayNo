package errors

import "fmt"

// customError is used to create own types of errors and handle them appropriately.
// It implements interface for standard error.
type customError struct {
	Message string
	Cause   error
}

// Error returns formatted error message.
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
