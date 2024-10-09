package repository

import (
	"context"
	"github.com/AlwaysSayNo/genesis-currency-api/email-service/internal/module/mail/model"
	"gorm.io/gorm"
)

type UserRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (ur *UserRepository) Create(ctx context.Context, user model.User) (*model.User, error) {
	result := ur.DB.WithContext(ctx).Create(&user)

	return &user, result.Error
}

func (ur *UserRepository) Update(ctx context.Context, user model.User) (*model.User, error) {
	result := ur.DB.WithContext(ctx).Save(&user)

	return &user, result.Error
}

func (ur *UserRepository) GetAllSubscribed(ctx context.Context) (*[]model.User, error) {
	var users []model.User

	result := ur.DB.WithContext(ctx).Where("is_subscribed = ?", true).Find(&users)

	return &users, result.Error
}

func (ur *UserRepository) GetByEmail(email string) (*model.User, error) {
	var user model.User
	result := ur.DB.Where("email = ?", email).First(&user)

	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}

func (ur *UserRepository) ExistsByEmail(email string) bool {
	var user model.User
	if result := ur.DB.Where("email = ?", email).First(&user); result.Error != nil {
		// result.Error - there is no user with such email
		return false
	}

	return true
}
