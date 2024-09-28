package model

type Currency struct {
	ID     int64   `db:"id"`
	Number float64 `json:"number"`
	Date   string  `json:"date"`
}
