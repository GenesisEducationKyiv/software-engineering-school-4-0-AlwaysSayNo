package currency

import (
	"net/http"

	"genesis-currency-api/pkg/errors"

	"genesis-currency-api/pkg/dto"

	"github.com/gin-gonic/gin"
)

type Rater interface {
	GetCurrencyRate() (dto.CurrencyResponseDTO, error)
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
	if result, err := c.rater.GetCurrencyRate(); err != nil {
		errors.AttachToCtx(err, ctx)
	} else {
		ctx.String(http.StatusOK, "%f", result.Number)
	}
}

// RegisterRoutes registers routes for passed Handler
func RegisterRoutes(r *gin.Engine, handler Handler) {
	routes := r.Group("/api/v1/rate")
	routes.GET("/", handler.GetLatest)
}
