package abstract

import (
	"encoding/json"
	"fmt"
	"genesis-currency-api/pkg/dto"
	"genesis-currency-api/pkg/errors"
	"io"
	"log"
	"net/http"
)

//todo update docs

type CurrencyRater interface {
	GetCurrencyRate() (*dto.CurrencyResponseDTO, error)
}

type Client struct {
	ApiURL       string
	ProviderName string
	Next         CurrencyRater
}

// ProcessCurrencyResponseDTO returns information about currency rate.
func (s *Client) ProcessCurrencyResponseDTO(responseDTO *dto.CurrencyResponseDTO, err error) (*dto.CurrencyResponseDTO, error) {
	if err == nil {
		log.Printf("Success response from %s: %v\n", s.ProviderName, *responseDTO)
		return responseDTO, nil
	} else if s.Next == nil {
		return nil, fmt.Errorf("end of the currency rater chain: %w", err)
	} else {
		log.Printf("Error while calling %s: %v", s.ProviderName, err)
	}

	log.Println("Try next currency getter")
	return s.Next.GetCurrencyRate()
}

// CallAPI prepares and executes call to the 3rd party API.
// Returns all available from 3rd party API currencies with the original schema.
func (s *Client) CallAPI(resp any) error {
	httpResp, err := http.Get(s.ApiURL)
	if err != nil {
		return errors.NewAPIError(fmt.Sprintf("doing GET request to %s", s.ProviderName), err)
	}
	defer func() {
		if err := httpResp.Body.Close(); err != nil {
			log.Printf("closing response body from %s: %v\n", s.ProviderName, err)
		}
	}()

	if httpResp.StatusCode != http.StatusOK {
		return errors.NewAPIError(fmt.Sprintf("unexpected response status code: %d\n", httpResp.StatusCode), nil)
	}

	body, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return errors.NewInvalidStateError("reading response body", err)
	}

	if err := json.Unmarshal(body, &resp); err != nil {
		return fmt.Errorf("unmarshalling response body: %w", err)
	}

	return nil
}

func (s *Client) SetNext(next CurrencyRater) {
	s.Next = next
}
