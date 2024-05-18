package controller

import (
	service2 "genesis-currency-api/internal/service"
	"genesis-currency-api/pkg/dto"
	"genesis-currency-api/pkg/errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserController struct {
	userService  *service2.UserService
	emailService *service2.EmailService
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

func (c *UserController) SendEmails(ctx *gin.Context) {
	err := c.emailService.SendEmails()
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, "")
}

func UserRegisterRoutes(r *gin.Engine, us *service2.UserService, es *service2.EmailService) {
	c := &UserController{
		us,
		es,
	}

	routes := r.Group("/api/emails")
	routes.GET("/", c.FindAll)
	routes.POST("/", c.Add)
	routes.POST("/send", c.SendEmails)
}
