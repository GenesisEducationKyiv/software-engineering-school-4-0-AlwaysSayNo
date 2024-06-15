package dto

import "strconv"

type APICurrencyResponseDTO struct {
	FromCcy string `json:"ccy"`
	BaseCcy string `json:"base_ccy"`
	Buy     string `json:"buy"`
	Sale    string `json:"sale"`
}

type CurrencyInfoDTO struct {
	FromCcy    string  `json:"fromCcy"`
	ToCcy      string  `json:"toCcy"`
	UpdateDate string  `json:"updateDate"`
	BuyRate    float64 `json:"buyRate"`
	SaleRate   float64 `json:"saleRate"`
}

type CurrencyResponseDto struct {
	Number float64 `json:"number"`
}

func APICurrencyResponseToInfoDTO(dto *APICurrencyResponseDTO) CurrencyInfoDTO {
	buy, _ := strconv.ParseFloat(dto.Buy, 64)
	sale, _ := strconv.ParseFloat(dto.Sale, 64)
	return CurrencyInfoDTO{
		FromCcy:  dto.FromCcy,
		ToCcy:    dto.BaseCcy,
		BuyRate:  buy,
		SaleRate: sale,
	}
}

func InfoToResponseDTO(dto *CurrencyInfoDTO) CurrencyResponseDto {
	return CurrencyResponseDto{
		Number: dto.SaleRate,
	}
}
