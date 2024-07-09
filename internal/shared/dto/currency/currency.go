package currency

type ResponseDTO struct {
	Number float64 `json:"number"`
}

type CachedCurrency struct {
	UpdateDate string `json:"updateDate"`
	ResponseDTO
}
