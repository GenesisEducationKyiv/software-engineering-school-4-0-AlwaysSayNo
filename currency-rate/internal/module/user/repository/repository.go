package repository

import (
	"context"
	"github.com/AlwaysSayNo/genesis-currency-api/currency-rate/internal/module/user/model"
	"gorm.io/gorm"
)

type Repository struct {
	DB *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{DB: db}
}

func (r *Repository) Create(ctx context.Context, user model.User) (*model.User, error) {
	result := r.DB.WithContext(ctx).Create(&user)

	return &user, result.Error
}

func (r *Repository) GetAll(ctx context.Context) (*[]model.User, error) {
	var users []model.User

	result := r.DB.WithContext(ctx).Find(&users)

	return &users, result.Error
}

// ExistsByEmail is used to check if user with such email already exists in database.
// Returns false if database responded with error, otherwise true.
func (r *Repository) ExistsByEmail(ctx context.Context, email string) bool {
	var user model.User
	if result := r.DB.WithContext(ctx).Where("email = ?", email).First(&user); result.Error != nil {
		// result.Error - there is no user with such email
		return false
	}

	return true
}
