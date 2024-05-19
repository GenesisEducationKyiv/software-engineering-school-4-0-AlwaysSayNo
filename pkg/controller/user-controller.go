package controller

import (
	"genesis-currency-api/internal/service"
	"genesis-currency-api/pkg/dto"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserController struct {
	userService  *service.UserService
	emailService *service.EmailService
}

func (c *UserController) Add(ctx *gin.Context) {
	email := ctx.PostForm("email")

	var saveDto dto.UserSaveRequestDTO
	saveDto.Email = email

	_, err := c.userService.Save(saveDto)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, "E-mail додано")
}

// UserRegisterRoutes creates an instance of UserController and registers routes for it.
func UserRegisterRoutes(r *gin.Engine, us *service.UserService, es *service.EmailService) {
	c := &UserController{
		us,
		es,
	}

	routes := r.Group("/api")
	routes.POST("/subscribe", c.Add)
}
