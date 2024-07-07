package private

import (
	"fmt"
	"genesis-currency-api/internal/module/currency/api/external/rater/abstract"
	"genesis-currency-api/internal/module/currency/config"
	"genesis-currency-api/internal/module/currency/dto"
	"genesis-currency-api/internal/module/currency/util/parser"
	sharcurrdto "genesis-currency-api/internal/shared/dto/currency"
	"genesis-currency-api/pkg/errors"
)

type CurrencyRater interface {
	GetCurrencyRate() (*sharcurrdto.CurrencyResponseDTO, error)
}

const (
	USD = "USD"
)

type Client struct {
	abstractClient abstract.Client
}

// NewClient is a factory function for Private Bank API Client.
func NewClient(cnf config.CurrencyRaterConfig) (*Client, error) {
	apiURL, err := parser.ParseURL(cnf.ThirdPartyAPIPrivateBank)
	if err != nil {
		return nil, err
	}

	c := &Client{
		abstractClient: abstract.Client{
			APIURL:       apiURL,
			ProviderName: "Private Bank API",
		},
	}

	return c, nil
}

// GetCurrencyRate gets data from its API and processes it with abstract client.
func (c *Client) GetCurrencyRate() (*sharcurrdto.CurrencyResponseDTO, error) {
	responseDTO, err := c.getUSDCurrencyFromAPI()
	return c.abstractClient.ProcessCurrencyResponseDTO(responseDTO, err)
}

// getUSDCurrencyFromAPI retrieves a full set of data from the 3rd party API call.
// Then it looks for USD currency and in case of occurrence maps response to dto.CurrencyResponseDTO.
// Returns a dto.CurrencyResponseDTO from all available from 3rd party API currencies.
func (c *Client) getUSDCurrencyFromAPI() (*sharcurrdto.CurrencyResponseDTO, error) {
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

// SetNext sets a next rater into raters chain.
func (c *Client) SetNext(next CurrencyRater) {
	c.abstractClient.SetNext(next)
}
