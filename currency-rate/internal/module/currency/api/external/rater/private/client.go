package private

import (
	"context"
	"fmt"
	"github.com/AlwaysSayNo/genesis-currency-api/currency-rate/pkg/apperrors"

	"github.com/AlwaysSayNo/genesis-currency-api/currency-rate/internal/module/currency/api/external/rater/abstract"
	"github.com/AlwaysSayNo/genesis-currency-api/currency-rate/internal/module/currency/config"
	"github.com/AlwaysSayNo/genesis-currency-api/currency-rate/internal/module/currency/dto"
	"github.com/AlwaysSayNo/genesis-currency-api/currency-rate/internal/module/currency/util/parser"
	sharcurrdto "github.com/AlwaysSayNo/genesis-currency-api/currency-rate/internal/shared/dto/currency"
)

type CurrencyRater interface {
	GetCurrencyRate(ctx context.Context) (*sharcurrdto.ResponseDTO, error)
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
func (c *Client) GetCurrencyRate(ctx context.Context) (*sharcurrdto.ResponseDTO, error) {
	responseDTO, err := c.getUSDCurrencyFromAPI(ctx)
	return c.abstractClient.ProcessCurrencyResponseDTO(ctx, responseDTO, err)
}

// getUSDCurrencyFromAPI retrieves a full set of data from the 3rd party API call.
// Then it looks for USD currency and in case of occurrence maps response to dto.CurrencyResponseDTO.
// Returns a dto.CurrencyResponseDTO from all available from 3rd party API currencies.
func (c *Client) getUSDCurrencyFromAPI(ctx context.Context) (*sharcurrdto.ResponseDTO, error) {
	var apiResponses []dto.PrivateAPICurrencyResponseDTO
	err := c.abstractClient.CallAPI(ctx, &apiResponses)
	if err != nil {
		return nil, fmt.Errorf("calling API: %w", err)
	}

	for _, r := range apiResponses {
		if r.FromCcy == USD {
			apiResponse := dto.PrivateAPICurrencyResponseToDTO(&r)
			return &apiResponse, nil
		}
	}

	return nil, apperrors.NewAPIError("No currency USD was found", nil)
}

// SetNext sets a next rater into raters chain.
func (c *Client) SetNext(next CurrencyRater) {
	c.abstractClient.SetNext(next)
}
