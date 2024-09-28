package repository

import (
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

func (cr *CurrencyRepository) Add(currency model.Currency) (*model.Currency, error) {
	result := cr.DB.Create(&currency)

	return &currency, result.Error
}

func (cr *CurrencyRepository) FindLatest() (*model.Currency, error) {
	var currency model.Currency

	result := cr.DB.Last(&currency)

	return &currency, result.Error
}
