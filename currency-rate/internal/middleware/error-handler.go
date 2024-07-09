package middleware

import (
	"net/http"

	"github.com/AlwaysSayNo/genesis-currency-api/currency-rate/pkg/errors"
	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Message string `json:"message"`
}

// ErrorHandler is a middleware for handling errors on a top level.
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors[0].Err

			switch e := err.(type) {
			case *errors.ValidationError:
				c.JSON(http.StatusBadRequest, ErrorResponse{Message: e.Error()})
			case *errors.UserWithEmailExistsError:
				c.JSON(http.StatusBadRequest, "Повертати, якщо e-mail вже є в базі даних")
			case *errors.InvalidInputError:
				c.JSON(http.StatusBadRequest, ErrorResponse{Message: e.Error()})
			default:
				c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "unknown error"})
			}

			c.Abort()
			return
		}
	}
}
