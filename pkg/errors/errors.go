package errors

import (
	"github.com/gin-gonic/gin"
	"log"
)

type ValidationError struct {
	customError
}

type DbError struct {
	customError
}

type UserWithEmailExistsError struct {
	customError
}

type ApiError struct {
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

func NewDbError(message string, cause error) *DbError {
	return &DbError{
		customError: customError{
			Message: message,
			Cause:   cause,
		},
	}
}

func NewApiError(message string, cause error) *ApiError {
	return &ApiError{
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
		log.Printf("Error %v happened while attaching error %v to current context", ctxErr, err)
	}
}
