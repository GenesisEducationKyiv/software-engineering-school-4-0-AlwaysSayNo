package cdn_jsdelivr

import (
	"encoding/json"
	"fmt"
	"genesis-currency-api/pkg/config"
	"genesis-currency-api/pkg/dto"
	apperrors "genesis-currency-api/pkg/errors"
	"genesis-currency-api/pkg/util/parser"
	"io"
	"log"
	"net/http"
)

type CurrencyRater interface {
	GetCurrencyRate() (*dto.CurrencyResponseDTO, error)
}

type Client struct {
	apiURL string
	next   CurrencyRater
}

// NewClient is a factory function for JsDelivr API Client
func NewClient(cnf config.CurrencyRaterConfig) (*Client, error) {
	if apiURL, err := parser.ParseURL(cnf.ThirdPartyAPICDNJSDeliver); err != nil {
		return nil, err
	} else {
		return &Client{
			apiURL: apiURL,
		}, err
	}
}

// GetCurrencyRate returns short information about currency rate.
func (s *Client) GetCurrencyRate() (*dto.CurrencyResponseDTO, error) {
	if responseDTO, err := s.getInternal(); err == nil {
		log.Printf("Success response from JsDelivr API: %v\n", *responseDTO)
		return responseDTO, nil
	} else if s.next == nil {
		return nil, fmt.Errorf("end of the currency rater chain: %w", err)
	} else {
		log.Printf("Error while calling JsDelivr API: %v", err)
	}

	log.Println("Try next currency getter")
	return s.next.GetCurrencyRate()
}

func (s *Client) getInternal() (*dto.CurrencyResponseDTO, error) {
	if data, err := s.callAPI(); err != nil {
		return nil, err
	} else {
		responseDTO := dto.JSDeliverAPICurrencyResponseToDTO(data)
		return &responseDTO, nil
	}
}

// callAPI prepares and executes call to the 3rd party API.
// Returns all available from 3rd party API currencies.
func (s *Client) callAPI() (*dto.JSDeliverAPICurrencyResponseDTO, error) {
	log.Println("Start calling JsDelivr API")

	resp, err := http.Get(s.apiURL)
	if err != nil {
		return nil, apperrors.NewAPIError("doing GET request to JsDelivr API", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("closing response body from JsDelivr API: %v\n", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, apperrors.NewAPIError(fmt.Sprintf("unexpected response status code: %d\n", resp.StatusCode), nil)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, apperrors.NewInvalidStateError("reading response body", err)
	}

	var apiResponse dto.JSDeliverAPICurrencyResponseDTO
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		return nil, apperrors.NewInvalidStateError("unmarshalling response body", err)
	}

	log.Println("Finish calling JsDelivr API")

	return &apiResponse, nil
}

func (s *Client) setNext(next CurrencyRater) {
	s.next = next
}
