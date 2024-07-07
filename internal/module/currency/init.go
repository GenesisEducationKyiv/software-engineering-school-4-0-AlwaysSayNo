package currency

import (
	"genesis-currency-api/internal/module/currency/api/handler"
	"genesis-currency-api/internal/module/currency/service"
	sharcurrdto "genesis-currency-api/internal/shared/dto/currency"
	"github.com/gin-gonic/gin"
)

type Provider interface {
	GetCurrencyRate() (*sharcurrdto.ResponseDTO, error)
}

type Service interface {
	GetCurrencyRate() (sharcurrdto.ResponseDTO, error)
	GetCachedCurrency() (sharcurrdto.CachedCurrency, error)
	UpdateCurrencyRates() error
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
