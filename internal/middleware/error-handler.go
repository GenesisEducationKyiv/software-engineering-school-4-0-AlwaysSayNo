package middleware

import (
	"genesis-currency-api/pkg/errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ErrorResponse struct {
	Message string `json:"message"`
}

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors[0].Err

			switch e := err.(type) {
			case *errors.ValidationError:
				c.JSON(http.StatusBadRequest, ErrorResponse{Message: e.Error()})
			case *errors.UserWithEmailExistsError:
				c.JSON(http.StatusBadRequest, ErrorResponse{Message: "User with such email is already subscribed"})
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
