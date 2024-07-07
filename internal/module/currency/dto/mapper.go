package dto

import (
	sharcurrdto "genesis-currency-api/internal/shared/dto/currency"
	"strconv"
)

func JSDeliverAPICurrencyResponseToDTO(dto *JSDeliverAPICurrencyResponseDTO) sharcurrdto.CurrencyResponseDTO {
	return sharcurrdto.CurrencyResponseDTO{
		Number: dto.Usd.Uah,
	}
}

func GovUaAPICurrencyResponseDTOToDTO(dto *GovUaAPICurrencyResponseDTO) sharcurrdto.CurrencyResponseDTO {
	return sharcurrdto.CurrencyResponseDTO{
		Number: dto.Rate,
	}
}

func PrivateAPICurrencyResponseToDTO(dto *PrivateAPICurrencyResponseDTO) sharcurrdto.CurrencyResponseDTO {
	sale, _ := strconv.ParseFloat(dto.Sale, 64)
	return sharcurrdto.CurrencyResponseDTO{
		Number: sale,
	}
}
