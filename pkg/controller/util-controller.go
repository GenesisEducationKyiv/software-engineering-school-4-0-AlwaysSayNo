package controller

import (
	"net/http"

	"genesis-currency-api/pkg/interface/service"

	"genesis-currency-api/pkg/errors"
	"github.com/gin-gonic/gin"
)

type UtilControllerImpl struct {
	userService  service.UserService
	emailService service.EmailService
}

func (c *UtilControllerImpl) FindAll(ctx *gin.Context) {
	result, err := c.userService.GetAll()
	if err != nil {
		errors.AttachToCtx(err, ctx)
		return
	}

	ctx.JSON(http.StatusOK, &result)
}

func (c *UtilControllerImpl) SendEmails(ctx *gin.Context) {
	err := c.emailService.SendEmails()
	if err != nil {
		errors.AttachToCtx(err, ctx)
		return
	}

	ctx.JSON(http.StatusOK, "Emails are successfully sent")
}

// RegisterUtilRoutes creates an instance of UtilControllerImpl and registers routes for it.
func RegisterUtilRoutes(r *gin.Engine, us service.UserService, es service.EmailService) {
	c := &UtilControllerImpl{
		us,
		es,
	}

	routes := r.Group("/api/util")
	routes.GET("/emails", c.FindAll)
	routes.POST("/emails/send", c.SendEmails)
}
