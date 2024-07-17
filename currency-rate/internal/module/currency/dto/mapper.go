package dto

import (
	"strconv"

	sharcurrdto "github.com/AlwaysSayNo/genesis-currency-api/currency-rate/internal/shared/dto/currency"
)

func JSDeliverAPICurrencyResponseToDTO(dto *JSDeliverAPICurrencyResponseDTO) sharcurrdto.ResponseDTO {
	return sharcurrdto.ResponseDTO{
		Number: dto.Usd.Uah,
	}
}

func GovUaAPICurrencyResponseDTOToDTO(dto *GovUaAPICurrencyResponseDTO) sharcurrdto.ResponseDTO {
	return sharcurrdto.ResponseDTO{
		Number: dto.Rate,
	}
}

func PrivateAPICurrencyResponseToDTO(dto *PrivateAPICurrencyResponseDTO) sharcurrdto.ResponseDTO {
	sale, _ := strconv.ParseFloat(dto.Sale, 64)
	return sharcurrdto.ResponseDTO{
		Number: sale,
	}
}
