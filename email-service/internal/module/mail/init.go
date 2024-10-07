package mail

import (
	emailhandl "github.com/AlwaysSayNo/genesis-currency-api/email-service/internal/module/mail/handler/email"
	"github.com/AlwaysSayNo/genesis-currency-api/email-service/internal/module/mail/repository"
	"github.com/AlwaysSayNo/genesis-currency-api/email-service/internal/module/mail/service"
	emailserv "github.com/AlwaysSayNo/genesis-currency-api/email-service/internal/module/mail/service/email"
	"github.com/AlwaysSayNo/genesis-currency-api/email-service/internal/module/mail/service/email/config"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handler interface {
	SendCurrencyPriceEmails(ctx *gin.Context)
}

type Module struct {
	UserRepository     repository.UserRepository
	CurrencyRepository repository.CurrencyRepository

	UserService     service.UserService
	CurrencyService service.CurrencyService
	EmailService    emailserv.Service

	EmailHandler emailhandl.Handler
}

func Init(db *gorm.DB, cnf config.EmailServiceConfig) *Module {
	userRepository := repository.NewUserRepository(db)
	currencyRepository := repository.NewCurrencyRepository(db)

	userService := service.NewUserService(userRepository)
	currencyService := service.NewCurrencyService(currencyRepository)
	emailService := emailserv.NewEmailService(userService, currencyService, cnf)

	emailHandler := emailhandl.NewHandler(emailService)

	return &Module{
		UserRepository:     *userRepository,
		CurrencyRepository: *currencyRepository,
		UserService:        *userService,
		CurrencyService:    *currencyService,
		EmailService:       *emailService,
		EmailHandler:       *emailHandler,
	}
}
