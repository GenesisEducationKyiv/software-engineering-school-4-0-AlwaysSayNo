package handler

import (
	"net/http"

	sharcurrdto "genesis-currency-api/internal/shared/dto/currency"

	"genesis-currency-api/pkg/errors"

	"github.com/gin-gonic/gin"
)

type Rater interface {
	GetCurrencyRate() (sharcurrdto.ResponseDTO, error)
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
