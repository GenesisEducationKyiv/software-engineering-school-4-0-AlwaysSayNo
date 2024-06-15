package job

import (
	"log"
	"time"

	"genesis-currency-api/internal/service"
	"github.com/robfig/cron/v3"
)

// SendEmailsJob is a cron function to send emails to subscribed users.
// It is executed every day at 9:00
func SendEmailsJob(emailService *service.EmailService) {
	scheduler := cron.New(cron.WithLocation(time.UTC))

	_, err := scheduler.AddFunc("0 9 * * *", func() {
		log.Println("Start job: Send Emails")

		if err := emailService.SendEmails(); err != nil {
			log.Printf("error with: Send Emails - %v\n", err)
		} else {
			log.Println("Finish job: Send Emails")
		}
	})
	if err != nil {
		log.Printf("failed to start UpdateCurrencyJob scheduler: %v\n", err)
		return
	}

	scheduler.Start()
}
