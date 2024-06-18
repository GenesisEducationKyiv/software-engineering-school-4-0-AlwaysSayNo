package job

import (
	"genesis-currency-api/pkg/interface/service"
	"time"

	"github.com/robfig/cron/v3"
)

func StartAllJobs(cs service.CurrencyService, es service.EmailService) {
	scheduler := cron.New(cron.WithLocation(time.UTC))

	UpdateCurrencyJob(scheduler, cs)
	SendEmailsJob(scheduler, es)

	scheduler.Start()
}
