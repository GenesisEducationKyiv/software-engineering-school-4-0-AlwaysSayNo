package handler

import "github.com/gin-gonic/gin"

type Interface interface {
	GetLatest(ctx *gin.Context)
}

// RegisterRoutes registers routes for passed Handler
func RegisterRoutes(r *gin.Engine, handler Interface) {
	routes := r.Group("/api/v1/rate")
	routes.GET("/", handler.GetLatest)
}
