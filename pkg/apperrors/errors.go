package apperrors

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

// customError is used to create own types of apperrors and handle them appropriately.
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

type ValidationError struct {
	customError
}

type DBError struct {
	customError
}

type UserWithEmailExistsError struct {
	customError
}

type APIError struct {
	customError
}

type InvalidStateError struct {
	customError
}

type InvalidInputError struct {
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

func NewUserWithEmailExistsError() *UserWithEmailExistsError {
	return &UserWithEmailExistsError{
		customError: customError{},
	}
}

func NewDBError(message string, cause error) *DBError {
	return &DBError{
		customError: customError{
			Message: message,
			Cause:   cause,
		},
	}
}

func NewAPIError(message string, cause error) *APIError {
	return &APIError{
		customError: customError{
			Message: message,
			Cause:   cause,
		},
	}
}

func NewInvalidStateError(message string, cause error) *InvalidStateError {
	return &InvalidStateError{
		customError: customError{
			Message: message,
			Cause:   cause,
		},
	}
}

func NewInvalidInputError(message string, cause error) *InvalidInputError {
	return &InvalidInputError{
		customError: customError{
			Message: message,
			Cause:   cause,
		},
	}
}

func AttachToCtx(err error, ctx *gin.Context) {
	ctxErr := ctx.Error(err)
	if ctxErr != nil {
		log.Printf("error %v happened while attaching error %v to current context\f", ctxErr, err)
	}
}
