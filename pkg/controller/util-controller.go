package controller

import (
	"genesis-currency-api/internal/service"
	"genesis-currency-api/pkg/errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UtilController struct {
	userService  *service.UserService
	emailService *service.EmailService
}

func (c *UtilController) FindAll(ctx *gin.Context) {
	result, err := c.userService.GetAll()

	if err != nil {
		errors.AttachToCtx(err, ctx)
		return
	}

	ctx.JSON(http.StatusOK, &result)
}

func (c *UtilController) SendEmails(ctx *gin.Context) {
	err := c.emailService.SendEmails()
	if err != nil {
		errors.AttachToCtx(err, ctx)
		return
	}

	ctx.JSON(http.StatusOK, "Emails are successfully sent")
}

// UtilRegisterRoutes creates an instance of UtilController and registers routes for it.
func UtilRegisterRoutes(r *gin.Engine, us *service.UserService, es *service.EmailService) {
	c := &UtilController{
		us,
		es,
	}

	routes := r.Group("/api/util")
	routes.GET("/emails", c.FindAll)
	routes.POST("/emails/send", c.SendEmails)
}
