package model

type User struct {
	ID           int64  `db:"id"`
	Email        string `db:"email"`
	IsSubscribed bool   `db:"is_subscribed"`
}
