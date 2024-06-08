package job

import (
	"fmt"
	"genesis-currency-api/internal/service"
	"github.com/robfig/cron/v3"
	"log"
	"time"
)

// SendEmailsJob is a cron function to send emails to subscribed users.
// It is executed every day at 9:00
func SendEmailsJob(emailService *service.EmailService) {
	scheduler := cron.New(cron.WithLocation(time.UTC))

	_, err := scheduler.AddFunc("0 9 * * *", func() {
		log.Println("Start job: Send Emails")

		err := emailService.SendEmails()
		if err != nil {
			fmt.Printf("Error with: Send Emails")
			fmt.Println(err)
		} else {
			log.Println("Finish job: Send Emails")
		}
	})
	if err != nil {
		fmt.Println("Error scheduling task:", err)
		return
	}

	scheduler.Start()
}
