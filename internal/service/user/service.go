package user

import (
	"fmt"
	"genesis-currency-api/internal/repository/user"
	"genesis-currency-api/pkg/dto"
	apperrors "genesis-currency-api/pkg/errors"
)

//todo update docs

type ServiceInterface interface {
	Save(user dto.UserSaveRequestDTO) (dto.UserResponseDTO, error)
	GetAll() ([]dto.UserResponseDTO, error)
}

type Service struct {
	userRepository *user.Repository
}

// NewService is a factory function for Service
func NewService(userRepository *user.Repository) *Service {
	return &Service{
		userRepository: userRepository,
	}
}

// Save saves user's information into the database. Only users with unique emails are saved.
// Returns UserResponseDTO with additional information or error.
func (s *Service) Save(saveRequestDTO dto.UserSaveRequestDTO) (*dto.UserResponseDTO, error) {
	userModel := dto.SaveRequestToModel(saveRequestDTO)

	// Check email uniqueness.
	if s.userRepository.ExistsByEmail(userModel.Email) {
		return nil, apperrors.NewUserWithEmailExistsError()
	}

	savedUser, err := s.userRepository.Create(userModel)
	if err != nil {
		return nil, fmt.Errorf("saving savedUser in database: %w", err)
	}

	userDTO := dto.ToDTO(*savedUser)

	return &userDTO, nil
}

// GetAll is used to get all available in database users' information.
// Returns all available UserResponseDTO.
func (s *Service) GetAll() ([]dto.UserResponseDTO, error) {
	users, err := s.userRepository.GetAll()
	if err != nil {
		return nil, fmt.Errorf("getting all users from database: %w", err)
	}

	result := make([]dto.UserResponseDTO, 0, len(*users))
	for _, u := range *users {
		result = append(result, dto.ToDTO(u))
	}

	return result, nil
}
