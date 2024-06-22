package private

import (
	"encoding/json"
	"fmt"
	"genesis-currency-api/pkg/config"
	"genesis-currency-api/pkg/dto"
	"genesis-currency-api/pkg/errors"
	"genesis-currency-api/pkg/util/parser"
	"io"
	"log"
	"net/http"
)

type CurrencyRater interface {
	GetCurrencyRate() (*dto.CurrencyResponseDTO, error)
}

const (
	USD = "USD"
)

type Client struct {
	apiURL string
	next   CurrencyRater
}

// NewClient is a factory function for Private Bank API Client
func NewClient(cnf config.CurrencyServiceConfig) (*Client, error) {
	if apiURL, err := parser.ParseURL(cnf.ThirdPartyAPIPrivateBank); err != nil {
		return nil, err
	} else {
		return &Client{
			apiURL: apiURL,
		}, err
	}
}

// GetCurrencyRate returns information about currency rate.
func (s *Client) GetCurrencyRate() (*dto.CurrencyResponseDTO, error) {
	if responseDTO, err := s.getUSDCurrencyFromAPI(); err == nil {
		log.Printf("Success response from Private Bank API: %v\n", *responseDTO)
		return responseDTO, nil
	} else if s.next == nil {
		return nil, fmt.Errorf("end of the currency rater chain: %w", err)
	} else {
		log.Printf("Error while calling Private Bank API: %v", err)
	}

	log.Println("Try next currency getter")
	return s.next.GetCurrencyRate()
}

// getUSDCurrencyFromAPI retrieves a full set of data from the 3rd party API call.
// Then it looks for USD currency and in case of occurrence maps ApiCurrencyResponse to CurrencyResponseDTO.
// Returns a CurrencyResponseDTO from all available from 3rd party API currencies.
func (s *Client) getUSDCurrencyFromAPI() (*dto.CurrencyResponseDTO, error) {
	apiResponses, err := s.callAPI()
	if err != nil {
		return nil, err
	}

	for _, r := range *apiResponses {
		if r.FromCcy == USD {
			apiResponse := dto.PrivateAPICurrencyResponseToDTO(&r)
			return &apiResponse, nil
		}
	}

	return nil, errors.NewAPIError("No currency UAH was found", nil)
}

// callAPI prepares and executes call to the 3rd party API.
// Returns all available from 3rd party API currencies with the original schema.
func (s *Client) callAPI() (*[]dto.PrivateAPICurrencyResponseDTO, error) {
	log.Println("Start calling Private Bank API")

	resp, err := http.Get(s.apiURL)
	if err != nil {
		return nil, errors.NewAPIError("doing GET request to Private Bank API", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("closing response body from Private Bank API: %v\n", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.NewAPIError(fmt.Sprintf("unexpected response status code: %d\n", resp.StatusCode), nil)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.NewInvalidStateError("reading response body", err)
	}

	var apiResponses []dto.PrivateAPICurrencyResponseDTO
	if err := json.Unmarshal(body, &apiResponses); err != nil {
		return nil, errors.NewInvalidStateError("unmarshalling response body", err)
	}

	log.Println("Finish calling Private Bank API")

	return &apiResponses, nil
}

func (s *Client) SetNext(next CurrencyRater) {
	s.next = next
}
