package job

import (
	"genesis-currency-api/internal/service"
	"github.com/robfig/cron/v3"
	"log"
	"time"
)

// UpdateCurrencyJob is a cron function to update currency service cache.
// It is executed every hour.
func UpdateCurrencyJob(currencyService *service.CurrencyService) {
	scheduler := cron.New(cron.WithLocation(time.UTC))
	scheduler.AddFunc("0 * * * *", func() {
		log.Println("Start job: Update Currency Rates")
		currencyService.UpdateCurrencyRates()
		log.Println("Finish job: Update Currency Rates")
	})
	scheduler.Start()
}
