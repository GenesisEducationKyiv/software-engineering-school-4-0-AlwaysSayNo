package controller

import (
	"genesis-currency-api/pkg/dto"
	"genesis-currency-api/pkg/errors"
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
	result, err := c.emailService.GetAll()

	if err != nil {
		ctx.Error(errors.NewDbError("", err))
		return
	}

	ctx.JSON(http.StatusOK, &result)
}

func (c *controller) Add(ctx *gin.Context) {
	var dto dto.UserSaveRequestDTO
	err := ctx.ShouldBindJSON(&dto)

	if err != nil {
		ctx.Error(errors.NewValidationError("", err))
		return
	}

	result, err := c.emailService.Save(dto)

	if err != nil {
		ctx.Error(errors.NewDbError("", err))
		return
	}

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
