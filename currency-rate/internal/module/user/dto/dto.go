package dto

type SaveRequestDTO struct {
	Email string `json:"email" binding:"required,email"`
}
