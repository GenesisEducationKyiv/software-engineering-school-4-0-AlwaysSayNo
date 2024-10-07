package currency

import (
	"github.com/AlwaysSayNo/genesis-currency-api/currency-rate/internal/module/user/decorator"

	"github.com/AlwaysSayNo/genesis-currency-api/currency-rate/internal/module/currency/api/handler"
	"github.com/AlwaysSayNo/genesis-currency-api/currency-rate/internal/module/currency/service"
)

type Module struct {
	Service service.Service
	Handler handler.Handler
}

func Init(currencyProvider service.Provider, producerClient decorator.ProducerClient) *Module {
	currencyService := service.NewService(currencyProvider, producerClient)
	currencyHandler := handler.NewHandler(currencyService)

	return &Module{
		Service: *currencyService,
		Handler: *currencyHandler,
	}
}
