package service

import (
	"context"
	"github.com/AlwaysSayNo/genesis-currency-api/email-service/internal/module/mail/config"
)

type NetSMTPMailer struct {
	cnf config.MailerConfig
}

func (m *NetSMTPMailer) SendEmail(
	ctx context.Context, emails []string,
	subject, message string,
) error {
	return nil
}

func NewNetSMTPMailer(cnf config.MailerConfig) *NetSMTPMailer {
	return &NetSMTPMailer{cnf: cnf}
}
