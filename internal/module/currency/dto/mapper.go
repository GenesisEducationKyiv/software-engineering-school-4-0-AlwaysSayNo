package dto

import "strconv"

func JSDeliverAPICurrencyResponseToDTO(dto *JSDeliverAPICurrencyResponseDTO) CurrencyResponseDTO {
	return CurrencyResponseDTO{
		Number: dto.Usd.Uah,
	}
}

func GovUaAPICurrencyResponseDTOToDTO(dto *GovUaAPICurrencyResponseDTO) CurrencyResponseDTO {
	return CurrencyResponseDTO{
		Number: dto.Rate,
	}
}

func PrivateAPICurrencyResponseToDTO(dto *PrivateAPICurrencyResponseDTO) CurrencyResponseDTO {
	sale, _ := strconv.ParseFloat(dto.Sale, 64)
	return CurrencyResponseDTO{
		Number: sale,
	}
}
