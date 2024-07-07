package dto

import (
	"genesis-currency-api/internal/user/model"
)

type UserResponseDTO struct {
	ID    int64  `json:"id"`
	Email string `json:"email"`
}

type UserSaveRequestDTO struct {
	Email string `json:"email" binding:"required,email"`
}

func ToDTO(entity model.User) UserResponseDTO {
	return UserResponseDTO{
		ID:    entity.ID,
		Email: entity.Email,
	}
}

func SaveRequestToModel(dto UserSaveRequestDTO) model.User {
	return model.User{
		Email: dto.Email,
	}
}
