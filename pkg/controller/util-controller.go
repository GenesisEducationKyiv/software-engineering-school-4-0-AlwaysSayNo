package controller

import (
	"genesis-currency-api/internal/service"
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
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, &result)
}

func (c *UtilController) SendEmails(ctx *gin.Context) {
	err := c.emailService.SendEmails()
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, "Emails are successfully sent")
}

func UtilRegisterRoutes(r *gin.Engine, us *service.UserService, es *service.EmailService) {
	c := &UtilController{
		us,
		es,
	}

	routes := r.Group("/api/util")
	routes.GET("/emails", c.FindAll)
	routes.POST("/emails/send", c.SendEmails)
}
