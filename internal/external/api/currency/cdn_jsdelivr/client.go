package cdn_jsdelivr

import (
	"fmt"
	"genesis-currency-api/internal/external/api/currency/abstract"
	"genesis-currency-api/pkg/config"
	"genesis-currency-api/pkg/dto"
	"genesis-currency-api/pkg/util/parser"
)

//todo update docs

type CurrencyRater interface {
	GetCurrencyRate() (*dto.CurrencyResponseDTO, error)
}

type Client struct {
	abstractClient abstract.Client
}

// NewClient is a factory function for JsDelivr API Client
func NewClient(cnf config.CurrencyRaterConfig) (*Client, error) {
	apiURL, err := parser.ParseURL(cnf.ThirdPartyAPICDNJSDeliver)
	if err != nil {
		return nil, fmt.Errorf("parsing JsDelivr API URL")
	}

	c := &Client{
		abstractClient: abstract.Client{
			ApiURL:       apiURL,
			ProviderName: "JsDelivr API",
		},
	}

	return c, nil
}

// GetCurrencyRate returns short information about currency rate.
func (c *Client) GetCurrencyRate() (*dto.CurrencyResponseDTO, error) {
	responseDTO, err := c.getUSDCurrencyFromAPI()
	return c.abstractClient.ProcessCurrencyResponseDTO(responseDTO, err)
}

func (c *Client) getUSDCurrencyFromAPI() (*dto.CurrencyResponseDTO, error) {
	var apiResponse dto.JSDeliverAPICurrencyResponseDTO
	err := c.abstractClient.CallAPI(&apiResponse)
	if err != nil {
		return nil, fmt.Errorf("calling API: %w", err)
	}

	responseDTO := dto.JSDeliverAPICurrencyResponseToDTO(&apiResponse)
	return &responseDTO, nil
}

func (c *Client) setNext(next CurrencyRater) {
	c.abstractClient.Next = next
}
