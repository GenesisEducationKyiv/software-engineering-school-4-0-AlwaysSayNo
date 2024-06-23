package private

import (
	"fmt"
	"genesis-currency-api/internal/external/api/currency/abstract"
	"genesis-currency-api/pkg/config"
	"genesis-currency-api/pkg/dto"
	"genesis-currency-api/pkg/errors"
	"genesis-currency-api/pkg/util/parser"
)

//todo update docs

type CurrencyRater interface {
	GetCurrencyRate() (*dto.CurrencyResponseDTO, error)
}

const (
	USD = "USD"
)

type Client struct {
	abstractClient abstract.Client
}

// NewClient is a factory function for Private Bank API Client
func NewClient(cnf config.CurrencyRaterConfig) (*Client, error) {
	apiURL, err := parser.ParseURL(cnf.ThirdPartyAPIPrivateBank)
	if err != nil {
		return nil, err
	}

	c := &Client{
		abstractClient: abstract.Client{
			ApiURL:       apiURL,
			ProviderName: "Private Bank API",
		},
	}

	return c, nil
}

// GetCurrencyRate returns information about currency rate.
func (c *Client) GetCurrencyRate() (*dto.CurrencyResponseDTO, error) {
	responseDTO, err := c.getUSDCurrencyFromAPI()
	return c.abstractClient.ProcessCurrencyResponseDTO(responseDTO, err)
}

// getUSDCurrencyFromAPI retrieves a full set of data from the 3rd party API call.
// Then it looks for USD currency and in case of occurrence maps ApiCurrencyResponse to CurrencyResponseDTO.
// Returns a CurrencyResponseDTO from all available from 3rd party API currencies.
func (c *Client) getUSDCurrencyFromAPI() (*dto.CurrencyResponseDTO, error) {
	var apiResponses []dto.PrivateAPICurrencyResponseDTO
	err := c.abstractClient.CallAPI(&apiResponses)
	if err != nil {
		return nil, fmt.Errorf("calling API: %w", err)
	}

	for _, r := range apiResponses {
		if r.FromCcy == USD {
			apiResponse := dto.PrivateAPICurrencyResponseToDTO(&r)
			return &apiResponse, nil
		}
	}

	return nil, errors.NewAPIError("No currency USD was found", nil)
}

func (c *Client) SetNext(next CurrencyRater) {
	c.abstractClient.Next = next
}
