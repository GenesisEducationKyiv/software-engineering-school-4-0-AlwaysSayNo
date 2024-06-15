package job

import (
	"genesis-currency-api/internal/service"
	"github.com/robfig/cron/v3"
	"log"
)

// SendEmailsJob is a cron function to send emails to subscribed users.
// It is executed every day at 9:00
func SendEmailsJob(cron *cron.Cron, emailService *service.EmailService) {
	_, err := cron.AddFunc("0 9 * * *", func() {
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
}
