package repository

import (
	"context"
	"github.com/AlwaysSayNo/genesis-currency-api/email-service/internal/module/mail/model"
	"gorm.io/gorm"
)

type CurrencyRepository struct {
	DB *gorm.DB
}

func NewCurrencyRepository(db *gorm.DB) *CurrencyRepository {
	return &CurrencyRepository{
		DB: db,
	}
}

func (cr *CurrencyRepository) Add(ctx context.Context, currency model.Currency) (*model.Currency, error) {
	result := cr.DB.WithContext(ctx).Create(&currency)

	return &currency, result.Error
}

func (cr *CurrencyRepository) FindLatest() (*model.Currency, error) {
	var currency model.Currency

	result := cr.DB.Last(&currency)

	return &currency, result.Error
}

func (r *CurrencyRepository) ExistsByEmail(email string) bool {
	var user model.User
	if result := r.DB.Where("email = ?", email).First(&user); result.Error != nil {
		// result.Error - there is no user with such email
		return false
	}

	return true
}
