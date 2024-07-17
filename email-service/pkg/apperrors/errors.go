package apperrors

type InvalidStateError struct {
	customError
}

func NewInvalidStateError(message string, cause error) *InvalidStateError {
	return &InvalidStateError{
		customError: customError{
			Message: message,
			Cause:   cause,
		},
	}
}
