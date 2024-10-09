package service

import (
	"context"
	"fmt"
	producerclient "github.com/AlwaysSayNo/genesis-currency-api/currency-rate/internal/mail"
	"log"
	"time"

	"github.com/AlwaysSayNo/genesis-currency-api/currency-rate/internal/module/currency/util/date"
	sharcurrdto "github.com/AlwaysSayNo/genesis-currency-api/currency-rate/internal/shared/dto/currency"
)

type Provider interface {
	GetCurrencyRate(ctx context.Context) (*sharcurrdto.ResponseDTO, error)
}

type ProducerClient interface {
	SendEvent(ctx context.Context, eventType string, data any) error
}

type CurrencyUpdateData struct {
	Number float64 `json:"number"`
	Date   string  `json:"data"`
}

type Service struct {
	currencyDTO      sharcurrdto.CachedCurrency
	currencyProvider Provider
	producerClient   ProducerClient
}

// NewService is a factory function for Service
func NewService(currencyProvider Provider, producerClient ProducerClient) *Service {
	// A cache value for 3rd party API response.
	return &Service{
		currencyDTO:      sharcurrdto.CachedCurrency{},
		currencyProvider: currencyProvider,
		producerClient:   producerClient,
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

	//todo maybe save in db through repository (because it might cause a lot of messages in queue)
	currencyRate, err := s.currencyProvider.GetCurrencyRate(ctx)
	if err != nil {
		return err
	}

	// save updated data to cache
	s.currencyDTO = sharcurrdto.CachedCurrency{
		UpdateDate:  date.Format(time.Now()),
		ResponseDTO: *currencyRate,
	}

	s.sendUpdateCurrencyRatesMessage(ctx, s.currencyDTO)

	log.Println("Finish updating currency rates")

	return nil
}

func (s *Service) sendUpdateCurrencyRatesMessage(ctx context.Context, currencyDTO sharcurrdto.CachedCurrency) {
	data := CurrencyUpdateData{
		Number: currencyDTO.Number,
		Date:   currencyDTO.UpdateDate,
	}

	if err := s.producerClient.SendEvent(ctx, producerclient.CurrencyUpdatedEvent, data); err != nil {
		return
	}
}
