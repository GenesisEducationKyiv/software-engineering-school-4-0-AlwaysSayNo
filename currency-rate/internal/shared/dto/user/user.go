package user

type ResponseDTO struct {
	ID           int64  `json:"id"`
	Email        string `json:"email"`
	IsSubscribed bool   `json:"isSubscribed"`
}
