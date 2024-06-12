package job

import (
	"fmt"
	"log"
	"time"

	"genesis-currency-api/internal/service"
	"github.com/robfig/cron/v3"
)

// UpdateCurrencyJob is a cron function to update currency service cache.
// It is executed every hour.
func UpdateCurrencyJob(currencyService *service.CurrencyService) {
	scheduler := cron.New(cron.WithLocation(time.UTC))
	_, err := scheduler.AddFunc("0 * * * *", func() {
		log.Println("Start job: Update Currency Rates")

		err := currencyService.UpdateCurrencyRates()
		if err != nil {
			fmt.Printf("Error with: Update Currency Rates")
			fmt.Println(err)
		} else {
			log.Println("Finish job: Update Currency Rates")
		}
	})
	if err != nil {
		return
	}
	scheduler.Start()
}
