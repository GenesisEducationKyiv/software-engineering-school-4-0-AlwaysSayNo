package controller

import (
	"genesis-currency-api/pkg/dto"
	"genesis-currency-api/pkg/errors"
	"genesis-currency-api/pkg/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserController struct {
	userService *service.UserService
}

func (c *UserController) FindAll(ctx *gin.Context) {
	result, err := c.userService.GetAll()

	if err != nil {
		ctx.Error(errors.NewDbError("", err))
		return
	}

	ctx.JSON(http.StatusOK, &result)
}

func (c *UserController) Add(ctx *gin.Context) {
	var dto dto.UserSaveRequestDTO
	err := ctx.ShouldBindJSON(&dto)

	if err != nil {
		ctx.Error(errors.NewValidationError("", err))
		return
	}

	result, err := c.userService.Save(dto)

	if err != nil {
		ctx.Error(errors.NewUserWithEmailExistsErrorError())
		return
	}

	ctx.JSON(http.StatusOK, &result)
}

func UserRegisterRoutes(r *gin.Engine, s *service.UserService) {
	c := &UserController{
		s,
	}

	routes := r.Group("/api/emails")
	routes.GET("/", c.FindAll)
	routes.POST("/", c.Add)
}
