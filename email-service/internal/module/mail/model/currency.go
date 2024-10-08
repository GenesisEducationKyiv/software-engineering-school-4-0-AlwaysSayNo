package model

type Currency struct {
	ID     int64   `db:"id"`
	Number float64 `db:"number"`
	Date   string  `db:"date"`
}
