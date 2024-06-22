package dto

import "strconv"

type PrivateAPICurrencyResponseDTO struct {
	FromCcy string `json:"ccy"`
	BaseCcy string `json:"base_ccy"`
	Buy     string `json:"buy"`
	Sale    string `json:"sale"`
}

type GovUaAPICurrencyResponseDTO struct {
	FromCcy string  `json:"cc"`
	Rate    float64 `json:"rate"`
}

type JSDeliverAPICurrencyResponseDTO struct {
	Date string `json:"date"`
	Usd  struct {
		Uah float64 `json:"uah"`
	} `json:"usd"`
}

type CurrencyInfoDTO struct {
	FromCcy    string  `json:"fromCcy"`
	ToCcy      string  `json:"toCcy"`
	UpdateDate string  `json:"updateDate"`
	BuyRate    float64 `json:"buyRate"`
	SaleRate   float64 `json:"saleRate"`
}

type CurrencyResponseDTO struct {
	Number float64 `json:"number"`
}

type CachedCurrency struct {
	UpdateDate string `json:"updateDate"`
	CurrencyResponseDTO
}

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
