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
		Type:         envs.Get("MAILER_TYPE"),
		SMTPHost:     envs.Get("SMTP_HOST"),
		SMTPPort:     envs.Get("SMTP_PORT"),
		SMTPUser:     envs.Get("SMTP_USER"),
		SMTPPassword: envs.Get("SMTP_PASSWORD"),
	}
}
