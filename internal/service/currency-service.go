package service

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"

	"genesis-currency-api/pkg/dto"
	"genesis-currency-api/pkg/errors"
	"genesis-currency-api/pkg/util/date"
	"github.com/spf13/viper"
)

type CurrencyService struct {
	currencyInfo dto.CurrencyInfoDTO
}

// NewCurrencyService is a factory function for CurrencyService
func NewCurrencyService() *CurrencyService {
	// A cache value for 3rd party API response.
	var currencyInfo dto.CurrencyInfoDTO
	c := &CurrencyService{
		currencyInfo,
	}

	err := c.UpdateCurrencyRates()
	if err != nil {
		log.Panic("error during creating CurrencyService: ", err)
	}

	return c
}

// GetCurrencyInfo returns extended information about currency rate.
// It is then used in email.
func (s *CurrencyService) GetCurrencyInfo() dto.CurrencyInfoDTO {
	return s.currencyInfo
}

// GetCurrencyRate returns short information about currency rate (sale rate).
// It is then used in API.
func (s *CurrencyService) GetCurrencyRate() dto.CurrencyResponseDto {
	return dto.InfoToResponseDto(&s.currencyInfo)
}

// getCurrencyRateFromAPI retrieves a full set of data from the 3rd party API call.
// Then it maps ApiCurrencyResponse to CurrencyInfoDTO and adds the time when call was made.
// Returns a list of CurrencyInfoDTO for all available from 3rd party API currencies.
func getCurrencyRateFromAPI() (*[]dto.CurrencyInfoDTO, error) {
	apiResponses, err := callAPI()
	if err != nil {
		return nil, err
	}

	infos := make([]dto.CurrencyInfoDTO, 0, len(*apiResponses))
	// Update (cache) time
	updateDate := date.Format(time.Now())
	for _, r := range *apiResponses {
		info := dto.APICurrencyResponseToInfoDTO(&r)
		info.UpdateDate = updateDate

		infos = append(infos, info)
	}

	return &infos, nil
}

// callAPI prepares and executes call to the 3rd party API.
// Returns all available from 3rd party API currencies with the original schema.
func callAPI() (*[]dto.APICurrencyResponseDTO, error) {
	log.Println("Start calling external API")

	apiURL, err := get3rdPartyURL()
	if err != nil {
		return nil, err
	}

	resp, err := http.Get(apiURL.String())
	if err != nil {
		return nil, errors.NewAPIError("Something went wrong while calling external API", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("error closing response body: %v", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.NewAPIError(fmt.Sprintf("Unexpected status code: %d", resp.StatusCode), nil)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.NewInvalidStateError("failed to read response body", err)
	}

	var apiResponses []dto.APICurrencyResponseDTO
	if err := json.Unmarshal(body, &apiResponses); err != nil {
		return nil, errors.NewInvalidStateError("failed to unmarshal response", err)
	}

	log.Println("Finish calling external API")

	return &apiResponses, nil
}

func get3rdPartyURL() (*url.URL, error) {
	parsedURL, err := url.ParseRequestURI(viper.Get("THIRD_PARTY_API").(string))
	if err != nil {
		return nil, errors.NewAPIError("Invalid URL", err)
	}

	return parsedURL, nil
}

// UpdateCurrencyRates is used to update #currencyInfo by calling 3rd party API.
// In this case #currencyInfo is a cache value of API response for currency USD.
func (s *CurrencyService) UpdateCurrencyRates() error {
	log.Println("Start updating currency rates")

	// Get list of 3rd party values.
	currencyRates, err := getCurrencyRateFromAPI()
	if err != nil {
		return err
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
		return errors.NewInvalidStateError("No currency UAH was found", nil)
	}

	log.Println("Finish updating currency rates")
	return nil
}
