package currency

import (
	"context"

	"github.com/AlwaysSayNo/genesis-currency-api/currency-rate/internal/module/currency/api/handler"
	"github.com/AlwaysSayNo/genesis-currency-api/currency-rate/internal/module/currency/service"
	sharcurrdto "github.com/AlwaysSayNo/genesis-currency-api/currency-rate/internal/shared/dto/currency"
	"github.com/gin-gonic/gin"
)

type Provider interface {
	GetCurrencyRate(ctx context.Context) (*sharcurrdto.ResponseDTO, error)
}

type Service interface {
	GetCurrencyRate(ctx context.Context) (sharcurrdto.ResponseDTO, error)
	GetCachedCurrency(ctx context.Context) (sharcurrdto.CachedCurrency, error)
	UpdateCurrencyRates(ctx context.Context) error
}

type Handler interface {
	GetLatest(ctx *gin.Context)
}

type Module struct {
	Service Service
	Handler Handler
}

func Init(currencyProvider Provider) *Module {
	currencyService := service.NewService(currencyProvider)
	currencyHandler := handler.NewHandler(currencyService)

	return &Module{
		Service: currencyService,
		Handler: currencyHandler,
	}
}
