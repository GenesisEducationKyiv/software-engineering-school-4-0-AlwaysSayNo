package errors

type ValidationError struct {
	customError
}

type DbError struct {
	customError
}

func NewValidationError(message string, cause error) *ValidationError {
	return &ValidationError{
		customError: customError{
			Message: message,
			Cause:   cause,
		},
	}
}

func NewDbError(message string, cause error) *DbError {
	return &DbError{
		customError: customError{
			Message: message,
			Cause:   cause,
		},
	}
}
