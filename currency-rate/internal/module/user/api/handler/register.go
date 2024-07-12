package handler

import "github.com/gin-gonic/gin"

type Interface interface {
	Add(ctx *gin.Context)
	FindAll(ctx *gin.Context)
	SendEmails(ctx *gin.Context)
}

// RegisterRoutes creates an instance of Handler and registers routes for it.
func RegisterRoutes(r *gin.Engine, handler Interface) {
	routes := r.Group("/api/v1/")
	routes.POST("/subscribe", handler.Add)
	routes.GET("/users", handler.FindAll)
	routes.POST("/users/emails/send", handler.SendEmails)
}
