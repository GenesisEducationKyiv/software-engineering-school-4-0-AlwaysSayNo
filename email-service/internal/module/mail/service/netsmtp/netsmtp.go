package netsmtp

import (
	"context"
	"fmt"
	"github.com/AlwaysSayNo/genesis-currency-api/common/pkg/apperrors"
	"github.com/AlwaysSayNo/genesis-currency-api/email-service/internal/module/mail/config"
	"log"
	"net/smtp"
)

type Mailer struct {
	cnf config.MailerConfig
}

func (m *Mailer) SendEmail(ctx context.Context, emails []string, subject, message string) error {
	mime := "MIME-version: 1.0;\r\nContent-Type: text/html; charset=\"UTF-8\";\r\n"

	mailMsg := []byte(fmt.Sprintf("%s\r\n%s\r\n%s", subject, mime, message))
	auth := smtp.PlainAuth("", m.cnf.SMTPUser, m.cnf.SMTPPassword, m.cnf.SMTPHost)

	done := make(chan error)
	go func() {
		done <- smtp.SendMail(m.cnf.SMTPHost+":"+m.cnf.SMTPPort, auth, m.cnf.SMTPUser, emails, mailMsg)
	}()

	select {
	case <-ctx.Done():
		return fmt.Errorf("email sending cancelled: %w", ctx.Err())
	case err := <-done:
		if err != nil {
			return apperrors.NewInvalidStateError("while sending email:", err)
		}
	}

	log.Println("Finish sending emails")

	return nil
}

func NewNetSMTPMailer(cnf config.MailerConfig) *Mailer {
	return &Mailer{cnf: cnf}
}
