package controller

import (
	"net/http"

	"genesis-currency-api/pkg/errors"

	"genesis-currency-api/internal/service"
	"github.com/gin-gonic/gin"
)

type CurrencyController interface {
	GetLatest(ctx *gin.Context)
}

type CurrencyControllerImpl struct {
	currencyService service.CurrencyService
}

func (c *CurrencyControllerImpl) GetLatest(ctx *gin.Context) {
	result, err := c.currencyService.GetCurrencyRate()
	if err != nil {
		errors.AttachToCtx(err, ctx)
	}
	ctx.String(http.StatusOK, "%f", result.Number)
}

// RegisterCurrencyRoutes creates an instance of CurrencyControllerImpl and registers routes for it.
func RegisterCurrencyRoutes(r *gin.Engine, s service.CurrencyService) {
	c := &CurrencyControllerImpl{
		s,
	}

	routes := r.Group("/api/rate")
	routes.GET("/", c.GetLatest)
}
