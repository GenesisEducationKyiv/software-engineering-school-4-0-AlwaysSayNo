package controller

import (
	"genesis-currency-api/internal/service"
	"genesis-currency-api/pkg/dto"
	"genesis-currency-api/pkg/errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserController struct {
	userService  *service.UserService
	emailService *service.EmailService
}

func (c *UserController) Add(ctx *gin.Context) {
	var saveDto dto.UserSaveRequestDTO
	err := ctx.ShouldBindJSON(&saveDto)

	if err != nil {
		ctx.Error(errors.NewValidationError("", err))
		return
	}

	result, err := c.userService.Save(saveDto)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, &result)
}

func UserRegisterRoutes(r *gin.Engine, us *service.UserService, es *service.EmailService) {
	c := &UserController{
		us,
		es,
	}

	routes := r.Group("/api")
	routes.POST("/subscribe", c.Add)
}
