package service

import (
	"encoding/json"
	"fmt"
	"genesis-currency-api/pkg/dto"
	"genesis-currency-api/pkg/errors"
	"genesis-currency-api/pkg/util/date"
	"io"
	"log"
	"net/http"
	"time"
)

type CurrencyService struct {
	currencyRates *[]dto.CurrencyResponseDto
}

func NewCurrencyService() *CurrencyService {
	var rates []dto.CurrencyResponseDto
	return &CurrencyService{
		&rates,
	}
}

func (s *CurrencyService) GetCurrencyRates() (*[]dto.CurrencyResponseDto, error) {
	return s.currencyRates, nil
}

func callCurrencyRates() (*[]dto.CurrencyResponseDto, error) {
	resp, err := http.Get(getApiUrl())
	if err != nil {
		return nil, errors.NewApiError("Something went wrong while calling external API", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.NewApiError(fmt.Sprintf("Unexpected status code: %d", resp.StatusCode), nil)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.NewApiError("Failed to read response body", err)
	}

	var apiResponses []dto.ApiCurrencyResponseDto
	if err := json.Unmarshal(body, &apiResponses); err != nil {
		return nil, errors.NewApiError("Failed to unmarshal response", err)
	}

	var result []dto.CurrencyResponseDto
	updateDate := date.Format(time.Now())
	for _, r := range apiResponses {
		tmp := dto.ApiCurrencyResponseToDTO(r)
		tmp.UpdateDate = updateDate
		result = append(result, tmp)
	}

	return &result, nil
}

func getApiUrl() string {
	return "https://api.privatbank.ua/p24api/pubinfo?exchange&coursid=5"
}

func (s *CurrencyService) UpdateCurrencyRates() {
	log.Println("Start updating currency rates")

	currencyRates, err := callCurrencyRates()
	if err != nil {
		log.Fatalf("Failed to update currency rates")
		return
	}
	s.currencyRates = currencyRates

	log.Println("Finish updating currency rates")
}
