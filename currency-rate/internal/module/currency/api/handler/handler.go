package handler

import (
	"context"
	"net/http"
	"time"

	sharcurrdto "github.com/AlwaysSayNo/genesis-currency-api/currency-rate/internal/shared/dto/currency"

	"github.com/AlwaysSayNo/genesis-currency-api/currency-rate/pkg/errors"

	"github.com/gin-gonic/gin"
)

const (
	DefaultRequestTime = 60
)

type Rater interface {
	GetCurrencyRate(ctx context.Context) (sharcurrdto.ResponseDTO, error)
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
	appCtx, appCancel := context.WithTimeout(ctx, DefaultRequestTime*time.Second)
	defer appCancel()

	if result, err := h.rater.GetCurrencyRate(appCtx); err != nil {
		errors.AttachToCtx(err, ctx)
	} else {
		ctx.String(http.StatusOK, "%f", result.Number)
	}
}
