package abstract

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	sharcurrdto "genesis-currency-api/internal/shared/dto/currency"
	"genesis-currency-api/pkg/errors"
)

type CurrencyRater interface {
	GetCurrencyRate() (*sharcurrdto.ResponseDTO, error)
}

type Client struct {
	APIURL       string
	ProviderName string
	Next         CurrencyRater
}

// ProcessCurrencyResponseDTO based on the input parameters generates the output.
func (c *Client) ProcessCurrencyResponseDTO(responseDTO *sharcurrdto.ResponseDTO,
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
	return c.Next.GetCurrencyRate()
}

// CallAPI prepares and executes call to the 3rd party API.
// Returns all available from 3rd party API currencies with the original schema.
func (c *Client) CallAPI(resp any) error {
	httpResp, err := http.Get(c.APIURL)
	if err != nil {
		return errors.NewAPIError(fmt.Sprintf("doing GET request to %s", c.ProviderName), err)
	}
	defer func() {
		if err := httpResp.Body.Close(); err != nil {
			log.Printf("closing response body from %s: %v\n", c.ProviderName, err)
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

// SetNext sets a next rater into raters chain.
func (c *Client) SetNext(next CurrencyRater) {
	c.Next = next
}
