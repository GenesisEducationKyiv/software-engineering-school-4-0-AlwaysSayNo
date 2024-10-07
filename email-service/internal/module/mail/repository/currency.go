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

func (cr *CurrencyRepository) FindLatest(ctx context.Context) (*model.Currency, error) {
	var currency model.Currency

	result := cr.DB.WithContext(ctx).Last(&currency)

	return &currency, result.Error
}
