package dto

import (
	"genesis-currency-api/internal/model"
)

type UserResponseDTO struct {
	ID        int64  `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type UserSaveRequestDTO struct {
	Email     string `json:"email" binding:"required,email"`
	FirstName string `json:"firstName" binding:"required,min=2,max=30"`
	LastName  string `json:"lastName" binding:"required,min=2,max=30"`
}

func ToDTO(entity model.User) UserResponseDTO {
	return UserResponseDTO{
		ID:        entity.ID,
		Email:     entity.Email,
		FirstName: entity.FirstName,
		LastName:  entity.LastName,
	}
}

func SaveRequestToModel(dto UserSaveRequestDTO) model.User {
	return model.User{
		Email:     dto.Email,
		FirstName: dto.FirstName,
		LastName:  dto.LastName,
	}
}
