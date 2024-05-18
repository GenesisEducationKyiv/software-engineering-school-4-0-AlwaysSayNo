package job

import (
	"fmt"
	"genesis-currency-api/internal/service"
	"github.com/go-co-op/gocron"
	"time"
)

func UpdateCurrency(currencyService *service.CurrencyService) {
	scheduler := gocron.NewScheduler(time.UTC)

	_, err := scheduler.Every(30).Minute().Do(currencyService.UpdateCurrencyRates)
	if err != nil {
		fmt.Println("Error scheduling task:", err)
		return
	}

	go scheduler.StartAsync()
}
