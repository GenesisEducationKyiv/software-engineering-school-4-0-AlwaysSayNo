package controller

import (
	"genesis-currency-api/pkg/errors"
	"genesis-currency-api/pkg/request"
	"genesis-currency-api/pkg/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type EmailController interface {
	FindAll(ctx *gin.Context)
	Add(ctx *gin.Context)
}

type controller struct {
	emailService *service.EmailService
}

func (c *controller) FindAll(ctx *gin.Context) {
	result := c.emailService.FindAll()
	ctx.JSON(http.StatusOK, &result)
}

func (c *controller) Add(ctx *gin.Context) {
	var email request.Email
	err := ctx.ShouldBindJSON(&email)

	if err != nil {
		ctx.Error(errors.NewValidationError("", err))
		return
	}

	result := c.emailService.Save(email)
	ctx.JSON(http.StatusOK, &result)
}

func RegisterRoutes(r *gin.Engine, s *service.EmailService) {
	c := &controller{
		s,
	}

	routes := r.Group("/api/emails")
	routes.GET("/", c.FindAll)
	routes.POST("/", c.Add)
}
