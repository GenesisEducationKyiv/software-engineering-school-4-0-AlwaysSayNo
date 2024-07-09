package email

import (
	"context"

	"github.com/AlwaysSayNo/genesis-currency-api/currency-rate/internal/module/email/api/handler"
	"github.com/AlwaysSayNo/genesis-currency-api/currency-rate/internal/module/email/config"
	"github.com/AlwaysSayNo/genesis-currency-api/currency-rate/internal/module/email/service"
	sharcurrdto "github.com/AlwaysSayNo/genesis-currency-api/currency-rate/internal/shared/dto/currency"
	"github.com/AlwaysSayNo/genesis-currency-api/currency-rate/internal/shared/dto/user"
	"github.com/gin-gonic/gin"
)

type UserGetter interface {
	GetAll() ([]user.ResponseDTO, error)
}

type DatedCurrencyGetter interface {
	GetCachedCurrency(ctx context.Context) (sharcurrdto.CachedCurrency, error)
}

type Service interface {
	SendEmails(ctx context.Context) error
}

type Handler interface {
	SendEmails(ctx *gin.Context)
}

type Module struct {
	Service Service
	Handler Handler
}

func Init(
	userGetter UserGetter,
	datedCurrencyGetter DatedCurrencyGetter,
	config config.EmailServiceConfig,
) *Module {
	emailService := service.NewService(userGetter, datedCurrencyGetter, config)
	emailHandler := handler.NewHandler(emailService)

	return &Module{
		Service: emailService,
		Handler: emailHandler,
	}
}
