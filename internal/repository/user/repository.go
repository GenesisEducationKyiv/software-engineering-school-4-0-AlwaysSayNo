package user

import (
	"genesis-currency-api/internal/model"
	"gorm.io/gorm"
)

type Repository struct {
	DB *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{DB: db}
}

func (r *Repository) Create(user model.User) (*model.User, error) {
	result := r.DB.Create(&user)

	return &user, result.Error
}

func (r *Repository) GetAll() (*[]model.User, error) {
	var users []model.User

	result := r.DB.Find(&users)

	return &users, result.Error
}

// ExistsByEmail is used to check if user with such email already exists in database.
// Returns false if database responded with error, otherwise true.
func (r *Repository) ExistsByEmail(email string) bool {
	var user model.User
	if result := r.DB.Where("email = ?", email).First(&user); result.Error != nil {
		// result.Error - there is no user with such email
		return false
	}

	return true
}
