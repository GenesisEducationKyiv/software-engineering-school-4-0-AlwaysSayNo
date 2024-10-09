package service

import (
	"context"
	"fmt"

	"github.com/AlwaysSayNo/genesis-currency-api/currency-rate/internal/module/user/dto"
	"github.com/AlwaysSayNo/genesis-currency-api/currency-rate/internal/module/user/model"
	userdto "github.com/AlwaysSayNo/genesis-currency-api/currency-rate/internal/shared/dto/user"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-AlwaysSayNo/pkg/apperrors"
)

type Repository interface {
	Save(ctx context.Context, user model.User) (*model.User, error)
	GetAll(ctx context.Context) (*[]model.User, error)
	Get(ctx context.Context, id int) (*model.User, error)
	ExistsByEmail(ctx context.Context, email string) bool
	ExistsById(ctx context.Context, id int) bool
	Update(ctx context.Context, user model.User) (*model.User, error)
}

type Service struct {
	userRepository Repository
}

// NewService is a factory function for Service
func NewService(userRepository Repository) *Service {
	return &Service{
		userRepository: userRepository,
	}
}

// Save saves user's information into the database. Only users with unique emails are saved.
// Returns dto.UserResponseDTO with additional information or error.
func (s *Service) Save(ctx context.Context, saveRequestDTO dto.SaveRequestDTO) (*userdto.ResponseDTO, error) {
	userModel := dto.SaveRequestToModel(saveRequestDTO)
	userModel.IsSubscribed = true

	// Check email uniqueness.
	if s.userRepository.ExistsByEmail(ctx, userModel.Email) {
		return nil, apperrors.NewUserWithEmailExistsError()
	}

	savedUser, err := s.userRepository.Save(ctx, userModel)
	if err != nil {
		return nil, fmt.Errorf("saving savedUser in database: %w", err)
	}

	userDTO := dto.ToDTO(*savedUser)

	return &userDTO, nil
}

// GetAll is used to get all available in database users' information.
// Returns all available dto.UserResponseDTO.
func (s *Service) GetAll(ctx context.Context) ([]userdto.ResponseDTO, error) {
	users, err := s.userRepository.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting all users from database: %w", err)
	}

	result := make([]userdto.ResponseDTO, 0, len(*users))
	for _, u := range *users {
		result = append(result, dto.ToDTO(u))
	}

	return result, nil
}

func (s *Service) ChangeSubscriptionStatus(ctx context.Context, id int, isSubscribed bool) (*userdto.ResponseDTO, error) {
	if !s.userRepository.ExistsById(ctx, id) {
		return nil, fmt.Errorf("user with such id doesn't exist")
	}

	user, err := s.userRepository.Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("getting user from db: %w", err)
	}

	user.IsSubscribed = isSubscribed

	savedUser, err := s.userRepository.Update(ctx, *user)
	userDTO := dto.ToDTO(*savedUser)

	return &userDTO, err
}
