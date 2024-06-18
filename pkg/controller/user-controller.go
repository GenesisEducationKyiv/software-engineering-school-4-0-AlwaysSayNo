package controller

import (
	"genesis-currency-api/pkg/interface/service"
	"net/http"

	"genesis-currency-api/pkg/dto"
	"genesis-currency-api/pkg/errors"
	"github.com/gin-gonic/gin"
)

type UserControllerImpl struct {
	userService  service.UserService
	emailService service.EmailService
}

func (c *UserControllerImpl) Add(ctx *gin.Context) {
	email := ctx.PostForm("email")

	var saveDto dto.UserSaveRequestDTO
	saveDto.Email = email

	_, err := c.userService.Save(saveDto)
	if err != nil {
		errors.AttachToCtx(err, ctx)
		return
	}

	ctx.JSON(http.StatusOK, "E-mail додано")
}

// RegisterUserRoutes creates an instance of UserControllerImpl and registers routes for it.
func RegisterUserRoutes(r *gin.Engine, us service.UserService, es service.EmailService) {
	c := &UserControllerImpl{
		us,
		es,
	}

	routes := r.Group("/api")
	routes.POST("/subscribe", c.Add)
}
