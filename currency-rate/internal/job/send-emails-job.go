package job

import (
	"context"
	"log"
)

type EmailNotifier interface {
	Notify(ctx context.Context) error
}

// GetSendEmailsJob is a cron function to send emails to subscribed users.
// It is executed every day at 9:00
func GetSendEmailsJob(ctx context.Context, emailNotifier EmailNotifier) WithCron {
	job := func() {
		log.Println("Start job: SendEmails")

		if err := emailNotifier.Notify(ctx); err != nil {
			log.Printf("error with: Send Emails - %v\n", err)
		} else {
			log.Println("Finish job: SendEmails")
		}
	}

	return WithCron{
		job:  job,
		cron: "0 9 * * *",
		name: "SendEmailsJob",
	}
}
