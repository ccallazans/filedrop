package repository

import (
	"context"

	"github.com/ccallazans/filedrop/internal/domain"
	"gorm.io/gorm"
)

type UserRepository interface {
	DB() *gorm.DB
	FindByUUID(ctx context.Context, uuid string) (*domain.User, error)
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
	Save(ctx context.Context, user *domain.User) error
	DeleteByUUID(ctx context.Context, uuid string) error
}
