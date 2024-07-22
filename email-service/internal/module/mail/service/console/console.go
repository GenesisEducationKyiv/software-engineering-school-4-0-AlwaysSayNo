package console

import (
	"context"
	"github.com/AlwaysSayNo/genesis-currency-api/email-service/internal/module/mail/config"
	"log"
)

type Mailer struct {
	cnf config.MailerConfig
}

func (m *Mailer) SendEmail(
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

func NewConsoleMailer(cnf config.MailerConfig) *Mailer {
	return &Mailer{cnf: cnf}
}
