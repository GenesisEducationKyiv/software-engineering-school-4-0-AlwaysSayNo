package job

import (
	"log"

	"github.com/robfig/cron/v3"
)

type CurrencyUpdater interface {
	UpdateCurrencyRates() error
}

// UpdateCurrencyJob is a cron function to update currency service cache.
// It is executed every hour.
func UpdateCurrencyJob(cron *cron.Cron, currencyUpdater CurrencyUpdater) {
	_, err := cron.AddFunc("0 * * * *", func() {
		log.Println("Start job: Update Currency Rates")

		if err := currencyUpdater.UpdateCurrencyRates(); err != nil {
			log.Printf("error with: Update Currency Rates - %v\n", err)
		} else {
			log.Println("Finish job: Update Currency Rates")
		}
	})
	if err != nil {
		log.Printf("failed to start UpdateCurrencyJob scheduler: %v\n", err)
		return
	}
}
