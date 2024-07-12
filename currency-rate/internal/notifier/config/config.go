package config

import (
	"github.com/AlwaysSayNo/genesis-currency-api/common/pkg/envs"
)

type EmailServiceConfig struct {
	EmailSubject string
}

func LoadEmailServiceConfig() EmailServiceConfig {
	return EmailServiceConfig{
		EmailSubject: envs.Get("EMAIL_SUBJECT"),
	}
}
