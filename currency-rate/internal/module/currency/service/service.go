package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/AlwaysSayNo/genesis-currency-api/currency-rate/internal/module/currency/util/date"
	sharcurrdto "github.com/AlwaysSayNo/genesis-currency-api/currency-rate/internal/shared/dto/currency"
)

type Provider interface {
	GetCurrencyRate(ctx context.Context) (*sharcurrdto.ResponseDTO, error)
}

type Service struct {
	currencyDTO      sharcurrdto.CachedCurrency
	currencyProvider Provider
}

// NewService is a factory function for Service
func NewService(currencyProvider Provider) *Service {
	// A cache value for 3rd party API response.
	return &Service{
		currencyDTO:      sharcurrdto.CachedCurrency{},
		currencyProvider: currencyProvider,
	}
}

// GetCurrencyRate returns short information about currency rate.
func (s *Service) GetCurrencyRate(ctx context.Context) (sharcurrdto.ResponseDTO, error) {
	currencyDTO, err := s.getCurrencyDTO(ctx)
	if err != nil {
		return sharcurrdto.ResponseDTO{}, fmt.Errorf("getting currency rate: %w", err)
	}

	return currencyDTO.ResponseDTO, nil
}

func (s *Service) GetCachedCurrency(ctx context.Context) (sharcurrdto.CachedCurrency, error) {
	return s.getCurrencyDTO(ctx)
}

func (s *Service) getCurrencyDTO(ctx context.Context) (sharcurrdto.CachedCurrency, error) {
	if s.currencyDTO.UpdateDate == "" {
		if err := s.UpdateCurrencyRates(ctx); err != nil {
			return s.currencyDTO, err
		}
	}

	return s.currencyDTO, nil
}

// UpdateCurrencyRates is used to update #currencyDTO by calling 3rd party API.
// In this case #currencyDTO is a cache value of API response for currency USD.
func (s *Service) UpdateCurrencyRates(ctx context.Context) error {
	log.Println("Start updating currency rates")

	currencyRate, err := s.currencyProvider.GetCurrencyRate(ctx)
	if err != nil {
		return err
	}

	// save updated data to cache
	s.currencyDTO = sharcurrdto.CachedCurrency{
		UpdateDate:  date.Format(time.Now()),
		ResponseDTO: *currencyRate,
	}

	log.Println("Finish updating currency rates")

	return nil
}
