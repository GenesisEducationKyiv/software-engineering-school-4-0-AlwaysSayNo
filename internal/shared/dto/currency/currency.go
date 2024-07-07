package currency

type CurrencyResponseDTO struct {
	Number float64 `json:"number"`
}

type CachedCurrency struct {
	UpdateDate string `json:"updateDate"`
	CurrencyResponseDTO
}
