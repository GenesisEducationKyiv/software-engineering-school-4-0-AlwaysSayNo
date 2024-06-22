package parser

import (
	"genesis-currency-api/pkg/errors"
	"net/url"
)

func ParseURL(rawURL string) (string, error) {
	parsedURL, err := url.ParseRequestURI(rawURL)
	if err != nil {
		return "", errors.NewAPIError("Invalid URL", err)
	}

	return parsedURL.String(), nil
}
