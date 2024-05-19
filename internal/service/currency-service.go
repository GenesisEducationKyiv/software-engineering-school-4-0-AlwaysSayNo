package service

import (
	"encoding/json"
	"fmt"
	"genesis-currency-api/pkg/dto"
	"genesis-currency-api/pkg/errors"
	"genesis-currency-api/pkg/util/date"
	"github.com/spf13/viper"
	"io"
	"log"
	"net/http"
	"time"
)

type CurrencyService struct {
	currencyInfo dto.CurrencyInfoDto
}

// NewCurrencyService is a factory function for CurrencyService
func NewCurrencyService() *CurrencyService {
	// A cache value for 3rd party API response.
	var currencyInfo dto.CurrencyInfoDto
	c := &CurrencyService{
		currencyInfo,
	}
	c.UpdateCurrencyRates()

	return c
}

// GetCurrencyInfo returns extended information about currency rate.
// It is then used in email.
func (s *CurrencyService) GetCurrencyInfo() dto.CurrencyInfoDto {
	return s.currencyInfo
}

// GetCurrencyRate returns short information about currency rate (sale rate).
// It is then used in API.
func (s *CurrencyService) GetCurrencyRate() dto.CurrencyResponseDto {
	return dto.InfoToResponseDTO(&s.currencyInfo)
}

// getCurrencyRateFromApi retrieves a full set of data from the 3rd party API call.
// Then it maps ApiCurrencyResponse to CurrencyInfoDto and adds the time when call was made.
// Returns a list of CurrencyInfoDto for all available from 3rd party API currencies.
func getCurrencyRateFromApi() (*[]dto.CurrencyInfoDto, error) {
	apiResponses, err := callApi()
	if err != nil {
		return nil, err
	}

	var infos []dto.CurrencyInfoDto
	// Update (cache) time
	updateDate := date.Format(time.Now())
	for _, r := range *apiResponses {
		info := dto.ApiCurrencyResponseToInfoDTO(&r)
		info.UpdateDate = updateDate

		infos = append(infos, info)
	}

	return &infos, nil
}

// callApi prepares and executes call to the 3rd party API.
// Returns all available from 3rd party API currencies with the original schema.
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
	return viper.Get("THIRD_PARTY_API").(string)
}

// UpdateCurrencyRates is used to update #currencyInfo by calling 3rd party API.
// In this case #currencyInfo is a cache value of API response for currency USD.
func (s *CurrencyService) UpdateCurrencyRates() {
	log.Println("Start updating currency rates")

	// Get list of 3rd party values.
	currencyRates, err := getCurrencyRateFromApi()
	if err != nil {
		log.Panic("Failed to update currency rates: ", err)
		return
	}

	// Retrieve USD value only.
	isUpdated := false
	for _, r := range *currencyRates {
		if r.FromCcy == "USD" {
			s.currencyInfo = r
			isUpdated = true
			break
		}
	}

	// If there is no USD currency - raise a panic
	if !isUpdated {
		log.Panicf("No currency %s was found", "UAH")
		return
	}

	log.Println("Finish updating currency rates")
}
