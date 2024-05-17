package request

type Email struct {
	Email string `json:"email" binding:"required,email"`
}
