package notifier

import (
	"github.com/AlwaysSayNo/genesis-currency-api/currency-rate/internal/notifier/api/handler"
	"github.com/AlwaysSayNo/genesis-currency-api/currency-rate/internal/notifier/config"
	"github.com/gin-gonic/gin"
)

type Handler interface {
	SendEmails(ctx *gin.Context)
}

type Module struct {
	EmailNotifier *EmailNotifier
	Handler       Handler
}

func Init(mailClient MailClient,
	datedCurrencyGetter DatedCurrencyGetter,
	userGetter UserGetter,
	cnf config.EmailServiceConfig) *Module {
	emailNotifier := NewEmailNotifier(mailClient, datedCurrencyGetter, userGetter, cnf)
	emailHandler := handler.NewHandler(emailNotifier)

	return &Module{
		Handler:       emailHandler,
		EmailNotifier: emailNotifier,
	}
}
