package service

import (
	"context"
	"fmt"
	"github.com/AlwaysSayNo/genesis-currency-api/email-service/internal/module/mail/dto"
	"github.com/AlwaysSayNo/genesis-currency-api/email-service/internal/module/mail/model"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-AlwaysSayNo/pkg/apperrors"
)

type UserRepository interface {
	Create(ctx context.Context, user model.User) (*model.User, error)
	ExistsByEmail(email string) bool
	GetByEmail(email string) (*model.User, error)
	GetAllSubscribed() (*[]model.User, error)
}

type UserService struct {
	userRepository UserRepository
}

func NewUserService(userRepository UserRepository) *UserService {
	return &UserService{
		userRepository: userRepository,
	}
}

func (us *UserService) Save(ctx context.Context, userSaveDTO dto.UserSaveDTO) error {
	userModel := dto.UserSaveDTOToUser(&userSaveDTO)
	userModel.IsSubscribed = true

	// Check email uniqueness.
	if us.userRepository.ExistsByEmail(userModel.Email) {
		return apperrors.NewUserWithEmailExistsError()
	}

	if _, err := us.userRepository.Create(ctx, userModel); err != nil {
		return fmt.Errorf("saving user: %w", err)
	}

	return nil
}

func (us *UserService) ChangeUserSubscriptionStatus(ctx context.Context, email string, isSubscribed bool) error {
	// Check if user exists.
	if !us.userRepository.ExistsByEmail(email) {
		return fmt.Errorf("user doesn't exist %s", email)
	}

	user, err := us.userRepository.GetByEmail(email)
	if err != nil {
		return fmt.Errorf("fetching user by email: %w", err)
	}

	user.IsSubscribed = isSubscribed

	_, err = us.userRepository.Create(ctx, *user)
	if err != nil {
		return fmt.Errorf("updating user's isSubscribed status: %w", err)
	}

	return nil
}

func (us *UserService) GetAllSubscribed() ([]dto.UserResponseDTO, error) {
	users, err := us.userRepository.GetAllSubscribed()
	if err != nil {
		return nil, fmt.Errorf("getting all users from database: %w", err)
	}

	result := make([]dto.UserResponseDTO, 0, len(*users))
	for _, u := range *users {
		result = append(result, dto.UserToUserResponseDTO(&u))
	}

	return result, nil
}
