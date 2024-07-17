package service

import "context"

type Mailer interface {
	SendEmail(ctx context.Context, emails []string, subject, message string) error
}
