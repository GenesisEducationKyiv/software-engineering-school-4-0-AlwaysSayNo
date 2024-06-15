package errors

import (
	"log"

	"github.com/gin-gonic/gin"
)

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
