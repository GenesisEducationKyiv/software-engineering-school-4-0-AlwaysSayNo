package config

import "github.com/AlwaysSayNo/genesis-currency-api/email-service/pkg/envs"

type MailerConfig struct {
	Type         string
	SMTPHost     string
	SMTPPort     string
	SMTPUser     string
	SMTPPassword string
	EmailSubject string
}

func LoadMailerConfig() MailerConfig {
	return MailerConfig{
		Type:         envs.Get("EMAIL_SERVICE_MAILER_TYPE"),
		SMTPHost:     envs.Get("EMAIL_SERVICE_SMTP_HOST"),
		SMTPPort:     envs.Get("EMAIL_SERVICE_SMTP_PORT"),
		SMTPUser:     envs.Get("EMAIL_SERVICE_SMTP_USER"),
		SMTPPassword: envs.Get("EMAIL_SERVICE_SMTP_PASSWORD"),
	}
}
