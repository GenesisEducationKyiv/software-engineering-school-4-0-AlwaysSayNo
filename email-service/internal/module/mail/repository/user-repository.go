package repository

import (
	"github.com/AlwaysSayNo/genesis-currency-api/email-service/internal/module/mail/model"
	"gorm.io/gorm"
)

type UserRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (ur *UserRepository) Create(user model.User) (*model.User, error) {
	result := ur.DB.Create(&user)

	return &user, result.Error
}

func (ur *UserRepository) GetAll() (*[]model.User, error) {
	var users []model.User

	result := ur.DB.Find(&users)

	return &users, result.Error
}
