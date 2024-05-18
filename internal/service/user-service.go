package service

import (
	"genesis-currency-api/internal/model"
	"genesis-currency-api/pkg/dto"
	"genesis-currency-api/pkg/errors"
	"gorm.io/gorm"
)

type UserService struct {
	DB *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{
		DB: db,
	}
}

func (s *UserService) Save(user dto.UserSaveRequestDTO) (dto.UserResponseDTO, error) {
	entity := dto.SaveRequestToModel(user)

	if s.existsByEmail(entity.Email) {
		return dto.UserResponseDTO{}, errors.NewUserWithEmailExistsError()
	}

	if result := s.DB.Create(&entity); result.Error != nil {
		return dto.UserResponseDTO{}, errors.NewDbError("", result.Error)
	}

	return dto.ToDTO(entity), nil
}

func (s *UserService) GetAll() ([]dto.UserResponseDTO, error) {
	var users []model.User

	if result := s.DB.Find(&users); result.Error != nil {
		return nil, errors.NewDbError("", result.Error)
	}

	var result []dto.UserResponseDTO
	for _, u := range users {
		result = append(result, dto.ToDTO(u))
	}

	return result, nil
}

func (s *UserService) existsByEmail(email string) bool {
	var user model.User
	if result := s.DB.Where("email = ?", email).First(&user); result.Error != nil {
		return false
	}

	return true
}
