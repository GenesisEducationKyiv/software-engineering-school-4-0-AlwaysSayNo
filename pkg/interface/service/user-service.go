package service

import "genesis-currency-api/pkg/dto"

type UserService interface {
	Save(user dto.UserSaveRequestDTO) (dto.UserResponseDTO, error)
	GetAll() ([]dto.UserResponseDTO, error)
}
