package job

import (
	"context"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-AlwaysSayNo/pkg/scheduler"
	"log"
)

type CurrencyUpdater interface {
	UpdateCurrencyRates(ctx context.Context) error
}

// GetUpdateCurrencyJob is a cron function to update currency service cache.
// It is executed every hour.
func GetUpdateCurrencyJob(ctx context.Context, currencyUpdater CurrencyUpdater) scheduler.WithCron {
	job := func() {
		log.Println("Start job: Update Currency Rates")

		if err := currencyUpdater.UpdateCurrencyRates(ctx); err != nil {
			log.Printf("error with: Update Currency Rates - %v\n", err)
		} else {
			log.Println("Finish job: Update Currency Rates")
		}
	}

	return scheduler.WithCron{
		Job:  job,
		Cron: "0 * * * *",
		Name: "UpdateCurrencyJob",
	}
}
