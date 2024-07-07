package handler

import "github.com/gin-gonic/gin"

type Interface interface {
	SendEmails(ctx *gin.Context)
}

// RegisterRoutes registers routes for passed Handler.
func RegisterRoutes(r *gin.Engine, handler Interface) {
	routes := r.Group("/api/v1")
	routes.POST("/emails/send", handler.SendEmails)
}
