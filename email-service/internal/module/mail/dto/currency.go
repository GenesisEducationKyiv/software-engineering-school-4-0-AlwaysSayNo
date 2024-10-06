package dto

type CurrencyDTO struct {
	Number float64 `json:"number"`
	Date   string  `json:"date"`
}

type CurrencyAddDTO struct {
	Number float64 `json:"number"`
	Date   string  `json:"date"`
}
