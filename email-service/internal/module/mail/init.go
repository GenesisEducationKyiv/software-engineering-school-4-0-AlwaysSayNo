package mail

import (
	"context"
	"github.com/AlwaysSayNo/genesis-currency-api/email-service/internal/module/mail/dto"
	"github.com/AlwaysSayNo/genesis-currency-api/email-service/internal/module/mail/repository"
	"github.com/AlwaysSayNo/genesis-currency-api/email-service/internal/module/mail/service"
	"gorm.io/gorm"
)

type UserService interface {
	Save(ctx context.Context, userSaveDTO dto.UserSaveDTO) error
	ChangeUserSubscriptionStatus(ctx context.Context, email string, isSubscribed bool) error
	GetAllSubscribed() ([]dto.UserResponseDTO, error)
}

type CurrencyService interface {
	Save(ctx context.Context, currencyAddDTO dto.CurrencyAddDTO) error
	FindLatest() (*dto.CurrencyDTO, error)
}

type Module struct {
	UserService     UserService
	CurrencyService CurrencyService
}

func Init(db *gorm.DB) *Module {
	userRepository := repository.NewUserRepository(db)
	currencyRepository := repository.NewCurrencyRepository(db)

	userService := service.NewUserService(userRepository)
	currencyService := service.NewCurrencyService(currencyRepository)

	return &Module{
		UserService:     userService,
		CurrencyService: currencyService,
	}
}
