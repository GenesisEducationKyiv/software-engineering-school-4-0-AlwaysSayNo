package job

import "genesis-currency-api/internal/service"

func StartAllJobs(cs *service.CurrencyService, es *service.EmailService) {
	UpdateCurrencyJob(cs)
	SendEmailsJob(es)
}
