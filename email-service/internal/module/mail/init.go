package mail

import (
	"context"
	"github.com/AlwaysSayNo/genesis-currency-api/email-service/internal/module/mail/dto"
	emailhandl "github.com/AlwaysSayNo/genesis-currency-api/email-service/internal/module/mail/handler/email"
	"github.com/AlwaysSayNo/genesis-currency-api/email-service/internal/module/mail/repository"
	"github.com/AlwaysSayNo/genesis-currency-api/email-service/internal/module/mail/service"
	emailserv "github.com/AlwaysSayNo/genesis-currency-api/email-service/internal/module/mail/service/email"
	"github.com/AlwaysSayNo/genesis-currency-api/email-service/internal/module/mail/service/email/config"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserService interface {
	Save(ctx context.Context, userSaveDTO dto.UserSaveDTO) error
	ChangeUserSubscriptionStatus(ctx context.Context, email string, isSubscribed bool) error
	GetAllSubscribed(ctx context.Context) ([]dto.UserResponseDTO, error)
}

type CurrencyService interface {
	Save(ctx context.Context, currencyAddDTO dto.CurrencyAddDTO) error
	FindLatest(ctx context.Context) (*dto.CurrencyDTO, error)
}

type EmailService interface {
	SendCurrencyPriceEmails(ctx context.Context, mailTransport service.Mailer) error
}

type EmailHandler interface {
	SendCurrencyPriceEmails(ctx *gin.Context)
}

type Module struct {
	UserService     UserService
	CurrencyService CurrencyService
	EmailService    EmailService
	EmailHandler    EmailHandler
}

func Init(db *gorm.DB, cnf config.EmailServiceConfig) *Module {
	userRepository := repository.NewUserRepository(db)
	currencyRepository := repository.NewCurrencyRepository(db)

	userService := service.NewUserService(userRepository)
	currencyService := service.NewCurrencyService(currencyRepository)
	emailService := emailserv.NewEmailService(userService, currencyService, cnf)

	emailHandler := emailhandl.NewHandler(emailService)

	return &Module{
		UserService:     userService,
		CurrencyService: currencyService,
		EmailService:    emailService,
		EmailHandler:    emailHandler,
	}
}
