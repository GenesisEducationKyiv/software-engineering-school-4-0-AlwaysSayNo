package controller

import (
	"net/http"

	"genesis-currency-api/internal/service"
	"github.com/gin-gonic/gin"
)

type CurrencyController struct {
	currencyService *service.CurrencyService
}

func (c *CurrencyController) GetLatest(ctx *gin.Context) {
	result := c.currencyService.GetCurrencyRate()
	ctx.String(http.StatusOK, "%f", result.Number)
}

// RegisterCurrencyRoutes creates an instance of CurrencyController and registers routes for it.
func RegisterCurrencyRoutes(r *gin.Engine, s *service.CurrencyService) {
	c := &CurrencyController{
		s,
	}

	routes := r.Group("/api/rate")
	routes.GET("/", c.GetLatest)
}
