package job

import (
	"context"
	"log"
)

type EmailSender interface {
	SendEmails(ctx context.Context) error
}

// GetSendEmailsJob is a cron function to send emails to subscribed users.
// It is executed every day at 9:00
func GetSendEmailsJob(ctx context.Context, emailSender EmailSender) WithCron {
	job := func() {
		log.Println("Start job: Send Emails")

		if err := emailSender.SendEmails(ctx); err != nil {
			log.Printf("error with: Send Emails - %v\n", err)
		} else {
			log.Println("Finish job: Send Emails")
		}
	}

	return WithCron{
		job:  job,
		cron: "0 9 * * *",
		name: "SendEmailsJob",
	}
}
