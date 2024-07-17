package cdnjsdelivr

import (
	"context"
	"fmt"

	"github.com/AlwaysSayNo/genesis-currency-api/currency-rate/internal/module/currency/api/external/rater/abstract"
	"github.com/AlwaysSayNo/genesis-currency-api/currency-rate/internal/module/currency/config"
	"github.com/AlwaysSayNo/genesis-currency-api/currency-rate/internal/module/currency/dto"
	"github.com/AlwaysSayNo/genesis-currency-api/currency-rate/internal/module/currency/util/parser"
	sharcurrdto "github.com/AlwaysSayNo/genesis-currency-api/currency-rate/internal/shared/dto/currency"
)

type CurrencyRater interface {
	GetCurrencyRate(ctx context.Context) (*sharcurrdto.ResponseDTO, error)
}

type Client struct {
	abstractClient abstract.Client
}

// NewClient is a factory function for JsDelivr API Client.
func NewClient(cnf config.CurrencyRaterConfig) (*Client, error) {
	apiURL, err := parser.ParseURL(cnf.ThirdPartyAPICDNJSDeliver)
	if err != nil {
		return nil, fmt.Errorf("parsing JsDelivr API URL")
	}

	c := &Client{
		abstractClient: abstract.Client{
			APIURL:       apiURL,
			ProviderName: "JsDelivr API",
		},
	}

	return c, nil
}

// GetCurrencyRate gets data from its API and processes it with abstract client.
func (c *Client) GetCurrencyRate(ctx context.Context) (*sharcurrdto.ResponseDTO, error) {
	responseDTO, err := c.getUSDCurrencyFromAPI(ctx)
	return c.abstractClient.ProcessCurrencyResponseDTO(ctx, responseDTO, err)
}

// getUSDCurrencyFromAPI calls JsDelivr API and maps it into dto.CurrencyResponseDTO.
func (c *Client) getUSDCurrencyFromAPI(ctx context.Context) (*sharcurrdto.ResponseDTO, error) {
	var apiResponse dto.JSDeliverAPICurrencyResponseDTO
	err := c.abstractClient.CallAPI(ctx, &apiResponse)
	if err != nil {
		return nil, fmt.Errorf("calling API: %w", err)
	}

	responseDTO := dto.JSDeliverAPICurrencyResponseToDTO(&apiResponse)
	return &responseDTO, nil
}

// SetNext sets a next rater into raters chain.
func (c *Client) SetNext(next CurrencyRater) {
	c.abstractClient.SetNext(next)
}
