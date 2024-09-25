package config

import (
	"github.com/AlwaysSayNo/genesis-currency-api/currency-rate/pkg/envs"
)

type EmailServiceConfig struct {
	EmailSubject string
}

func LoadEmailServiceConfig() EmailServiceConfig {
	return EmailServiceConfig{
		EmailSubject: envs.Get("CURRENCY_SERVICE_EMAIL_SUBJECT"),
	}
}
