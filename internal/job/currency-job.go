package job

import (
	"context"
	"log"
)

type CurrencyUpdater interface {
	UpdateCurrencyRates() error
}

// GetUpdateCurrencyJob is a cron function to update currency service cache.
// It is executed every hour.
func GetUpdateCurrencyJob(ctx context.Context, currencyUpdater CurrencyUpdater) WithCron {
	job := func() {
		log.Println("Start job: Update Currency Rates")

		if err := currencyUpdater.UpdateCurrencyRates(); err != nil {
			log.Printf("error with: Update Currency Rates - %v\n", err)
		} else {
			log.Println("Finish job: Update Currency Rates")
		}
	}

	return WithCron{
		job:  job,
		cron: "0 * * * *",
		name: "UpdateCurrencyJob",
	}
}
