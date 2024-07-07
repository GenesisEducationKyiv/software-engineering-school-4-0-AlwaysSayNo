package service

import (
	"fmt"
	"log"
	"time"

	"genesis-currency-api/internal/module/currency/util/date"
	sharcurrdto "genesis-currency-api/internal/shared/dto/currency"
)

type Provider interface {
	GetCurrencyRate() (*sharcurrdto.ResponseDTO, error)
}

type Service struct {
	currencyDTO      sharcurrdto.CachedCurrency
	currencyProvider Provider
}

// NewService is a factory function for Service
func NewService(currencyProvider Provider) *Service {
	// A cache value for 3rd party API response.
	return &Service{
		sharcurrdto.CachedCurrency{},
		currencyProvider,
	}
}

// GetCurrencyRate returns short information about currency rate.
func (s *Service) GetCurrencyRate() (sharcurrdto.ResponseDTO, error) {
	currencyDTO, err := s.getCurrencyDTO()
	if err != nil {
		return sharcurrdto.ResponseDTO{}, fmt.Errorf("getting currency rate: %w", err)
	}

	return currencyDTO.CurrencyResponseDTO, nil
}

func (s *Service) GetCachedCurrency() (sharcurrdto.CachedCurrency, error) {
	return s.getCurrencyDTO()
}

func (s *Service) getCurrencyDTO() (sharcurrdto.CachedCurrency, error) {
	if s.currencyDTO.UpdateDate == "" {
		if err := s.UpdateCurrencyRates(); err != nil {
			return s.currencyDTO, err
		}
	}

	return s.currencyDTO, nil
}

// UpdateCurrencyRates is used to update #currencyDTO by calling 3rd party API.
// In this case #currencyDTO is a cache value of API response for currency USD.
func (s *Service) UpdateCurrencyRates() error {
	log.Println("Start updating currency rates")

	currencyRate, err := s.currencyProvider.GetCurrencyRate()
	if err != nil {
		return err
	}

	// save updated data to cache
	s.currencyDTO = sharcurrdto.CachedCurrency{
		UpdateDate:          date.Format(time.Now()),
		CurrencyResponseDTO: *currencyRate,
	}

	log.Println("Finish updating currency rates")
	return nil
}
