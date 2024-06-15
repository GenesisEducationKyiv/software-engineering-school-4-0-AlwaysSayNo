package job

import (
	"genesis-currency-api/internal/service"
	"github.com/robfig/cron/v3"
	"log"
)

// UpdateCurrencyJob is a cron function to update currency service cache.
// It is executed every hour.
func UpdateCurrencyJob(cron *cron.Cron, currencyService *service.CurrencyService) {
	_, err := cron.AddFunc("0 * * * *", func() {
		log.Println("Start job: Update Currency Rates")

		if err := currencyService.UpdateCurrencyRates(); err != nil {
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
