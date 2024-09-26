package config

import (
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-AlwaysSayNo/pkg/envs"
)

type EmailServiceConfig struct {
	EmailSubject string
}

func LoadEmailServiceConfig() EmailServiceConfig {
	return EmailServiceConfig{
		EmailSubject: envs.Get("CURRENCY_SERVICE_EMAIL_SUBJECT"),
	}
}
