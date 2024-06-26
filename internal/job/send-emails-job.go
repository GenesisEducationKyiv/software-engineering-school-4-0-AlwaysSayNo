package job

import (
	"log"

	"github.com/robfig/cron/v3"
)

type EmailSender interface {
	SendEmails() error
}

// SendEmailsJob is a cron function to send emails to subscribed users.
// It is executed every day at 9:00
func SendEmailsJob(cron *cron.Cron, emailSender EmailSender) {
	_, err := cron.AddFunc("0 9 * * *", func() {
		log.Println("Start job: Send Emails")

		if err := emailSender.SendEmails(); err != nil {
			log.Printf("error with: Send Emails - %v\n", err)
		} else {
			log.Println("Finish job: Send Emails")
		}
	})
	if err != nil {
		log.Printf("failed to start UpdateCurrencyJob scheduler: %v\n", err)
		return
	}
}
