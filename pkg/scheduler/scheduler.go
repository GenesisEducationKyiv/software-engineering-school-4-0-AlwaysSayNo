package scheduler

import (
	"github.com/robfig/cron/v3"
	"log"
	"time"
)

type WithCron struct {
	Job  func()
	Name string
	Cron string
}

func StartAllJobs(jobs ...WithCron) *cron.Cron {
	scheduler := cron.New(cron.WithLocation(time.UTC))

	for _, j := range jobs {
		if id, err := scheduler.AddFunc(j.Cron, j.Job); err != nil {
			log.Printf("failed to start %s (id %d) scheduler: %v\n", j.Name, id, err)
		}
	}

	scheduler.Start()

	return scheduler
}
