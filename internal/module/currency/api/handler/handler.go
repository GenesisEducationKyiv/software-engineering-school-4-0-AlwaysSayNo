package handler

import (
	sharcurrdto "genesis-currency-api/internal/shared/dto/currency"
	"net/http"

	"genesis-currency-api/pkg/errors"

	"github.com/gin-gonic/gin"
)

type Rater interface {
	GetCurrencyRate() (sharcurrdto.CurrencyResponseDTO, error)
}

type Handler struct {
	rater Rater
}

func NewHandler(rater Rater) *Handler {
	return &Handler{
		rater,
	}
}

func (h *Handler) GetLatest(ctx *gin.Context) {
	if result, err := h.rater.GetCurrencyRate(); err != nil {
		errors.AttachToCtx(err, ctx)
	} else {
		ctx.String(http.StatusOK, "%f", result.Number)
	}
}
