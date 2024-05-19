package controller

import (
	"genesis-currency-api/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type CurrencyController struct {
	currencyService *service.CurrencyService
}

func (c *CurrencyController) GetLatest(ctx *gin.Context) {
	result := c.currencyService.GetCurrencyRate()
	ctx.JSON(http.StatusOK, &result)
}

func CurrencyRegisterRoutes(r *gin.Engine, s *service.CurrencyService) {
	c := &CurrencyController{
		s,
	}

	routes := r.Group("/api/rate")
	routes.GET("/", c.GetLatest)
}
