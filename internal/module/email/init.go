package email

import (
	"context"
	"genesis-currency-api/internal/module/email/api/handler"
	"genesis-currency-api/internal/module/email/config"
	"genesis-currency-api/internal/module/email/service"
	sharcurrdto "genesis-currency-api/internal/shared/dto/currency"
	"genesis-currency-api/internal/shared/dto/user"
	"github.com/gin-gonic/gin"
)

type UserGetter interface {
	GetAll() ([]user.ResponseDTO, error)
}

type DatedCurrencyGetter interface {
	GetCachedCurrency() (sharcurrdto.CachedCurrency, error)
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
