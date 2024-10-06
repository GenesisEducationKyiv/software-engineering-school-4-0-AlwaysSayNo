package service

import (
	"context"
	"fmt"
	"github.com/AlwaysSayNo/genesis-currency-api/email-service/internal/module/mail/dto"
	"github.com/AlwaysSayNo/genesis-currency-api/email-service/internal/module/mail/model"
)

type CurrencyRepository interface {
	Add(ctx context.Context, currency model.Currency) (*model.Currency, error)
	FindLatest() (*model.Currency, error)
}

type CurrencyService struct {
	currencyRepository CurrencyRepository
}

func NewCurrencyService(currencyRepository CurrencyRepository) *CurrencyService {
	return &CurrencyService{
		currencyRepository: currencyRepository,
	}
}

func (cs *CurrencyService) Save(ctx context.Context, currencyAddDTO dto.CurrencyAddDTO) error {
	currencyModel := dto.CurrencyAddDTOToCurrency(&currencyAddDTO)

	if _, err := cs.currencyRepository.Add(ctx, currencyModel); err != nil {
		return fmt.Errorf("saving currency: %w", err)
	}

	return nil
}

func (cs *CurrencyService) FindLatest() (*dto.CurrencyDTO, error) {
	currencyModel, err := cs.currencyRepository.FindLatest()
	if err != nil {
		return nil, fmt.Errorf("finding latest currency: %w", err)
	}

	currencyDTO := dto.CurrencyToCurrencyDTO(currencyModel)

	return &currencyDTO, nil
}
