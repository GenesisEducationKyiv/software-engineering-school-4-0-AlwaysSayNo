package parser

import (
	"fmt"
	"net/url"
)

func ParseURL(rawURL string) (string, error) {
	parsedURL, err := url.ParseRequestURI(rawURL)
	if err != nil {
		return "", fmt.Errorf("parsing raw API URL: %w", err)
	}

	return parsedURL.String(), nil
}
