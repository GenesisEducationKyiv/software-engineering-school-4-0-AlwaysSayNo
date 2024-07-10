package config

import (
	"github.com/AlwaysSayNo/genesis-currency-api/common/pkg/envs"
)

type EmailServiceConfig struct {
	SMTPHost     string
	SMTPPort     string
	SMTPUser     string
	SMTPPassword string
	EmailSubject string
}

func LoadEmailServiceConfig() EmailServiceConfig {
	return EmailServiceConfig{
		SMTPHost:     envs.Get("SMTP_HOST"),
		SMTPPort:     envs.Get("SMTP_PORT"),
		SMTPUser:     envs.Get("SMTP_USER"),
		SMTPPassword: envs.Get("SMTP_PASSWORD"),
		EmailSubject: envs.Get("EMAIL_SUBJECT"),
	}
}
