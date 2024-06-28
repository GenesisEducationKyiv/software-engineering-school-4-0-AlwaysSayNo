package job

import (
	"time"

	"github.com/robfig/cron/v3"
)

func StartAllJobs(cu CurrencyUpdater, es EmailSender) {
	scheduler := cron.New(cron.WithLocation(time.UTC))

	UpdateCurrencyJob(scheduler, cu)
	SendEmailsJob(scheduler, es)

	scheduler.Start()
}
