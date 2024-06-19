package currency

import (
	"net/http"

	"genesis-currency-api/pkg/dto"

	"github.com/gin-gonic/gin"
)

type Rater interface {
	GetCurrencyRate() dto.CurrencyResponseDTO
}

type Handler struct {
	rater Rater
}

func NewHandler(rater Rater) *Handler {
	return &Handler{
		rater,
	}
}

func (c *Handler) GetLatest(ctx *gin.Context) {
	result := c.rater.GetCurrencyRate()
	ctx.String(http.StatusOK, "%f", result.Number)
}

// RegisterRoutes registers routes for passed Handler
func RegisterRoutes(r *gin.Engine, handler Handler) {
	routes := r.Group("/api/rate")
	routes.GET("/", handler.GetLatest)
}
