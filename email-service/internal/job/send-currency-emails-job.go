package job

import (
	"context"
	"github.com/AlwaysSayNo/genesis-currency-api/email-service/internal/module/mail/service"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-AlwaysSayNo/pkg/scheduler"
	"log"
)

type CurrencyEmailSender interface {
	SendCurrencyPriceEmails(ctx context.Context, mailTransport service.Mailer) error
}

// GetSendCurrencyEmailsJob is a cron function to send emails to subscribed users.
// It is executed every day at 9:00
func GetSendCurrencyEmailsJob(ctx context.Context, currencyEmailSender CurrencyEmailSender, mailTransport service.Mailer) scheduler.WithCron {
	job := func() {
		log.Println("Start job: SendCurrencyEmailsJob")

		if err := currencyEmailSender.SendCurrencyPriceEmails(ctx, mailTransport); err != nil {
			log.Printf("error with: Send Emails - %v\n", err)
		} else {
			log.Println("Finish job: SendCurrencyEmailsJob")
		}
	}

	return scheduler.WithCron{
		Job:  job,
		Cron: "0 9 * * *",
		Name: "SendCurrencyEmailsJob",
	}
}
