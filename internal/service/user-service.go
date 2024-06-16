package service

import (
	"genesis-currency-api/internal/model"
	"genesis-currency-api/pkg/dto"
	"genesis-currency-api/pkg/errors"
	"gorm.io/gorm"
)

type UserService interface {
	Save(user dto.UserSaveRequestDTO) (dto.UserResponseDTO, error)
	GetAll() ([]dto.UserResponseDTO, error)
}

type UserServiceImpl struct {
	DB *gorm.DB
}

// NewUserServiceImpl is a factory function for UserServiceImpl
func NewUserServiceImpl(db *gorm.DB) *UserServiceImpl {
	return &UserServiceImpl{
		DB: db,
	}
}

// Save saves user's information into the database. Only users with unique emails are saved.
// Returns UserResponseDTO with additional information or error.
func (s *UserServiceImpl) Save(user dto.UserSaveRequestDTO) (dto.UserResponseDTO, error) {
	entity := dto.SaveRequestToModel(user)

	// Check email uniqueness.
	if s.existsByEmail(entity.Email) {
		return dto.UserResponseDTO{}, errors.NewUserWithEmailExistsError()
	}

	if result := s.DB.Create(&entity); result.Error != nil {
		return dto.UserResponseDTO{}, errors.NewDBError("", result.Error)
	}

	return dto.ToDTO(entity), nil
}

// GetAll is used to get all available in database users' information.
// Returns all available UserResponseDTO.
func (s *UserServiceImpl) GetAll() ([]dto.UserResponseDTO, error) {
	var users []model.User

	if result := s.DB.Find(&users); result.Error != nil {
		return nil, errors.NewDBError("", result.Error)
	}

	result := make([]dto.UserResponseDTO, 0, len(users))
	for _, u := range users {
		result = append(result, dto.ToDTO(u))
	}

	return result, nil
}

// existsByEmail is used to check if user with such email already exists in database.
// Returns false if database responded with error, otherwise true.
func (s *UserServiceImpl) existsByEmail(email string) bool {
	var user model.User
	if result := s.DB.Where("email = ?", email).First(&user); result.Error != nil {
		// result.Error - there is no user with such email
		return false
	}

	return true
}
