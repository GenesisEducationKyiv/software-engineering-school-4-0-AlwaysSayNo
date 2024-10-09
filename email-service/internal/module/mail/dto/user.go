package dto

type UserSaveDTO struct {
	Email string `db:"email"`
}

type UserResponseDTO struct {
	ID           int64  `db:"id"`
	Email        string `db:"email"`
	IsSubscribed bool   `db:"isSubscribed"`
}
