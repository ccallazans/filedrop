package repository

import (
	"database/sql"

	"github.com/ccallazans/filedrop/internal/domain"
	"gorm.io/gorm"
)

type IUser interface {
	Begin() (*sql.DB, error)
	Rollback() error
	Commit() error
	Save(user *domain.User) error
	FindByEmail(email string) (*domain.User, error)
}
