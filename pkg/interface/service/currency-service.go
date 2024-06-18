package service

import "genesis-currency-api/pkg/dto"

type CurrencyService interface {
	GetCurrencyInfo() dto.CurrencyInfoDTO
	GetCurrencyRate() dto.CurrencyResponseDto
	UpdateCurrencyRates() error
}
