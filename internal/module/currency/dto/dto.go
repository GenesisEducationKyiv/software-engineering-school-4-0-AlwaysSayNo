package dto

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

type CurrencyResponseDTO struct {
	Number float64 `json:"number"`
}

type CachedCurrency struct {
	UpdateDate string `json:"updateDate"`
	CurrencyResponseDTO
}
