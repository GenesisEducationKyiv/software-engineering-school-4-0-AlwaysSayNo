package abstract

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-AlwaysSayNo/pkg/apperrors"
	"io"
	"log"
	"net/http"

	sharcurrdto "github.com/AlwaysSayNo/genesis-currency-api/currency-rate/internal/shared/dto/currency"
)

type CurrencyRater interface {
	GetCurrencyRate(ctx context.Context) (*sharcurrdto.ResponseDTO, error)
}

type Client struct {
	APIURL       string
	ProviderName string
	Next         CurrencyRater
}

// ProcessCurrencyResponseDTO based on the input parameters generates the output.
func (c *Client) ProcessCurrencyResponseDTO(
	ctx context.Context,
	responseDTO *sharcurrdto.ResponseDTO,
	err error,
) (*sharcurrdto.ResponseDTO, error) {
	if err == nil {
		log.Printf("Success response from %s: %v\n", c.ProviderName, *responseDTO)
		return responseDTO, nil
	} else if err != nil && c.Next == nil {
		return nil, fmt.Errorf("end of the currency rater chain: %w", err)
	}

	log.Printf("Error while calling %s: %v", c.ProviderName, err)
	log.Println("Try next currency getter")
	return c.Next.GetCurrencyRate(ctx)
}

// CallAPI prepares and executes call to the 3rd party API.
// Returns all available from 3rd party API currencies with the original schema.
func (c *Client) CallAPI(ctx context.Context, resp any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.APIURL, nil)
	if err != nil {
		return apperrors.NewAPIError(fmt.Sprintf("making GET request to %s", c.ProviderName), err)
	}

	httpResp, err := http.DefaultClient.Do(req)
	if err != nil {
		return apperrors.NewAPIError(fmt.Sprintf("doing GET request to %s", c.ProviderName), err)
	}
	defer func() {
		if err := httpResp.Body.Close(); err != nil {
			log.Printf("closing response body from %s: %v\n", c.ProviderName, err)
		}
	}()

	if httpResp.StatusCode != http.StatusOK {
		return apperrors.NewAPIError(fmt.Sprintf("unexpected response status code: %d\n", httpResp.StatusCode), nil)
	}

	body, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return apperrors.NewInvalidStateError("reading response body", err)
	}

	if err := json.Unmarshal(body, &resp); err != nil {
		return fmt.Errorf("unmarshalling response body: %w", err)
	}

	return nil
}

// SetNext sets a next rater into raters chain.
func (c *Client) SetNext(next CurrencyRater) {
	c.Next = next
}
