package service

import (
	"genesis-currency-api/pkg/util/date"
	"log"
	"time"

	"genesis-currency-api/pkg/dto"
)

type CurrencyServiceInterface interface {
	GetCurrencyRate() (*dto.CurrencyResponseDTO, error)
	UpdateCurrencyRates() error
}

type CurrencyProvider interface {
	GetCurrencyRate() (*dto.CurrencyResponseDTO, error)
}

type CurrencyService struct {
	currencyDTO      dto.CachedCurrency
	currencyProvider CurrencyProvider
}

// NewCurrencyService is a factory function for CurrencyService
func NewCurrencyService(currencyProvider CurrencyProvider) *CurrencyService {
	// A cache value for 3rd party API response.
	return &CurrencyService{
		dto.CachedCurrency{},
		currencyProvider,
	}
}

// GetCurrencyRate returns short information about currency rate.
func (s *CurrencyService) GetCurrencyRate() (dto.CurrencyResponseDTO, error) {
	data, err := s.getCurrencyDTO()
	return data.CurrencyResponseDTO, err
}

func (s *CurrencyService) getCurrencyDTO() (dto.CachedCurrency, error) {
	if s.currencyDTO.UpdateDate == "" {
		if err := s.UpdateCurrencyRates(); err != nil {
			return s.currencyDTO, err
		}
	}

	return s.currencyDTO, nil
}

// UpdateCurrencyRates is used to update #currencyDTO by calling 3rd party API.
// In this case #currencyDTO is a cache value of API response for currency USD.
func (s *CurrencyService) UpdateCurrencyRates() error {
	log.Println("Start updating currency rates")

	currencyRate, err := s.currencyProvider.GetCurrencyRate()
	if err != nil {
		return err
	}

	// save updated data to cache
	s.currencyDTO = dto.CachedCurrency{
		UpdateDate:          date.Format(time.Now()),
		CurrencyResponseDTO: *currencyRate,
	}

	log.Println("Finish updating currency rates")
	return nil
}
