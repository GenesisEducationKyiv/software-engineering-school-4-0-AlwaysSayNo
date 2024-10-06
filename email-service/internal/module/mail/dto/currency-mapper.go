package dto

import (
	"github.com/AlwaysSayNo/genesis-currency-api/email-service/internal/module/broker"
	"github.com/AlwaysSayNo/genesis-currency-api/email-service/internal/module/mail/model"
)

func CurrencyAddDTOToCurrency(dto *CurrencyAddDTO) model.Currency {
	return model.Currency{
		Number: dto.Number,
		Date:   dto.Date,
	}
}

func CurrencyToCurrencyDTO(dto *model.Currency) CurrencyDTO {
	return CurrencyDTO{
		Number: dto.Number,
		Date:   dto.Date,
	}
}

func CurrencyUpdateDataToCurrencyAddDTO(dto *broker.CurrencyUpdateData) CurrencyAddDTO {
	return CurrencyAddDTO{
		Number: dto.Number,
		Date:   dto.Date,
	}
}
