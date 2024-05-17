package service

import (
	"errors"
	"genesis-currency-api/pkg/dto"
	"genesis-currency-api/pkg/models"
	"gorm.io/gorm"
)

type EmailService struct {
	db *gorm.DB
}

func New(db *gorm.DB) *EmailService {
	return &EmailService{
		db: db,
	}
}

func (s *EmailService) Save(user dto.UserSaveRequestDTO) (dto.UserResponseDTO, error) {
	entity := dto.SaveRequestToModel(user)

	if result := s.db.Create(&entity); result.Error != nil {
		return dto.UserResponseDTO{}, errors.New(result.Error.Error())
	}

	return dto.ToDTO(entity), nil
}

func (s *EmailService) GetAll() ([]dto.UserResponseDTO, error) {
	var users []models.User

	if result := s.db.Find(&users); result.Error != nil {
		return nil, errors.New(result.Error.Error())
	}

	var result []dto.UserResponseDTO
	for _, u := range users {
		result = append(result, dto.ToDTO(u))
	}

	return result, nil
}
