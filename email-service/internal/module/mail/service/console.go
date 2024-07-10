package service

import (
	"context"
	"github.com/AlwaysSayNo/genesis-currency-api/email-service/internal/module/mail/config"
	"log"
)

type ConsoleMailer struct {
	cnf config.MailerConfig
}

func (m *ConsoleMailer) SendEmail(
	_ context.Context, emails []string,
	subject, message string,
) error {
	log.Println(
		"Sending email:",
		"fromEmail", m.cnf.SMTPUser,
		"toEmails", emails,
		"subject", subject,
		"message", message,
	)
	return nil
}

func NewConsoleMailer(cnf config.MailerConfig) *ConsoleMailer {
	return &ConsoleMailer{cnf: cnf}
}
