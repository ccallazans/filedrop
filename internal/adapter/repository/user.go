package repository

import (
	"github.com/ccallazans/filedrop/internal/domain"
	"github.com/ccallazans/filedrop/internal/domain/repository"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) repository.IUser {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) Save(user *domain.User) error {

	result := r.db.Create(&user)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *userRepository) FindByEmail(email string) (*domain.User, error) {

	user := &domain.User{}

	result := r.db.First(user, "email = ?", email)
	if result.Error != nil {
		return &domain.User{}, result.Error
	}

	return user, nil
}
