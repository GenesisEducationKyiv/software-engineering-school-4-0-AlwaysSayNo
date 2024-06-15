package job

import (
	"genesis-currency-api/internal/service"
	"github.com/robfig/cron/v3"
	"time"
)

func StartAllJobs(cs *service.CurrencyService, es *service.EmailService) {
	scheduler := cron.New(cron.WithLocation(time.UTC))

	UpdateCurrencyJob(scheduler, cs)
	SendEmailsJob(scheduler, es)

	scheduler.Start()
}
