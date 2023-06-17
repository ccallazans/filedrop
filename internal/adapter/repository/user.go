package repository

import (
	"database/sql"
	"log"

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

func (r *userRepository) WithTrx(trxHandle *gorm.DB) *userRepository {
	if trxHandle == nil {
		log.Print("Transaction Database not found")
		return r
	}

	r.db = trxHandle
	return r
}

func (r *userRepository) Begin() (*sql.DB, error) {
	db, err := r.db.Begin().DB()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (r *userRepository) Rollback() error { return r.db.Rollback().Error }
func (r *userRepository) Commit() error   { return r.db.Commit().Error }

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
