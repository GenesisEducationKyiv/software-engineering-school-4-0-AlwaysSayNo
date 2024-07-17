package job

import (
	"log"
	"time"

	"github.com/robfig/cron/v3"
)

type WithCron struct {
	job  func()
	name string
	cron string
}

func StartAllJobs(jobs ...WithCron) *cron.Cron {
	scheduler := cron.New(cron.WithLocation(time.UTC))

	for _, j := range jobs {
		if id, err := scheduler.AddFunc(j.cron, j.job); err != nil {
			log.Printf("failed to start %s (id %d) scheduler: %v\n", j.name, id, err)
		}
	}

	scheduler.Start()

	return scheduler
}
