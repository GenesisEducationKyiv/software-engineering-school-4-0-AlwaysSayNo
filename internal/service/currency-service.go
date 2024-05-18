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
	currencyRate dto.CurrencyResponseDto
}

func NewCurrencyService() *CurrencyService {
	var currencyRate dto.CurrencyResponseDto
	return &CurrencyService{
		currencyRate,
	}
}

func (s *CurrencyService) GetCurrencyRate() dto.CurrencyResponseDto {
	return s.currencyRate
}

func getCurrencyRateFromApi() (*[]dto.CurrencyResponseDto, error) {
	apiResponses, err := callApi()
	if err != nil {
		return nil, err
	}

	var result []dto.CurrencyResponseDto
	updateDate := date.Format(time.Now())
	for _, r := range *apiResponses {
		tmp := dto.ApiCurrencyResponseToDTO(r)
		tmp.UpdateDate = updateDate
		result = append(result, tmp)
	}

	return &result, nil
}

func callApi() (*[]dto.ApiCurrencyResponseDto, error) {
	log.Println("Start calling external API")
	apiUrl := getApiUrl()

	resp, err := http.Get(apiUrl)
	if err != nil {
		return nil, errors.NewApiError("Something went wrong while calling external API", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.NewApiError(fmt.Sprintf("Unexpected status code: %d", resp.StatusCode), nil)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.NewInvalidStateError("Failed to read response body", err)
	}

	var apiResponses []dto.ApiCurrencyResponseDto
	if err := json.Unmarshal(body, &apiResponses); err != nil {
		return nil, errors.NewInvalidStateError("Failed to unmarshal response", err)
	}

	log.Println("Finish calling external API")

	return &apiResponses, nil
}

func getApiUrl() string {
	return "https://api.privatbank.ua/p24api/pubinfo?exchange&coursid=5"
}

func (s *CurrencyService) UpdateCurrencyRates() {
	log.Println("Start updating currency rates")

	currencyRates, err := getCurrencyRateFromApi()
	if err != nil {
		log.Panic("Failed to update currency rates")
		return
	}

	isUpdated := false
	for _, r := range *currencyRates {
		if r.FromCcy == "USD" {
			s.currencyRate = r
			isUpdated = true
			break
		}
	}

	if !isUpdated {
		log.Panicf("No currency %s was found", "UAH")
		return
	}

	log.Println("Finish updating currency rates")
}
