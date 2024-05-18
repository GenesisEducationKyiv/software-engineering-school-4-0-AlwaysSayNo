package controller

import (
	"genesis-currency-api/pkg/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type CurrencyController struct {
	currencyService *service.CurrencyService
}

func (c *CurrencyController) GetLatest(ctx *gin.Context) {
	result, err := c.currencyService.GetCurrencyRate()

	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, &result)
}

func CurrencyRegisterRoutes(r *gin.Engine, s *service.CurrencyService) {
	c := &CurrencyController{
		s,
	}

	routes := r.Group("/api/currency")
	routes.GET("/", c.GetLatest)
}
