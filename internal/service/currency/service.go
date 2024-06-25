package currency

import (
	"genesis-currency-api/pkg/util/date"
	"log"
	"time"

	"genesis-currency-api/pkg/dto"
)

type ServiceInterface interface {
	GetCurrencyRate() (*dto.CurrencyResponseDTO, error)
	UpdateCurrencyRates() error
}

type Provider interface {
	GetCurrencyRate() (*dto.CurrencyResponseDTO, error)
}

type Service struct {
	currencyDTO      dto.CachedCurrency
	currencyProvider Provider
}

// NewService is a factory function for Service
func NewService(currencyProvider Provider) *Service {
	// A cache value for 3rd party API response.
	return &Service{
		dto.CachedCurrency{},
		currencyProvider,
	}
}

// GetCurrencyRate returns short information about currency rate.
func (s *Service) GetCurrencyRate() (dto.CurrencyResponseDTO, error) {
	currencyDTO, err := s.getCurrencyDTO()
	return currencyDTO.CurrencyResponseDTO, err
}

func (s *Service) GetCachedCurrency() (dto.CachedCurrency, error) {
	return s.getCurrencyDTO()
}

func (s *Service) getCurrencyDTO() (dto.CachedCurrency, error) {
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
	s.currencyDTO = dto.CachedCurrency{
		UpdateDate:          date.Format(time.Now()),
		CurrencyResponseDTO: *currencyRate,
	}

	log.Println("Finish updating currency rates")
	return nil
}
