package repository

import (
	"context"

	"github.com/ccallazans/filedrop/internal/domain"
	"gorm.io/gorm"
)

type FileAccessRepository interface {
	DB() *gorm.DB
	FindByHash(ctx context.Context, hash string) (*domain.FileAccess, error)
	Save(ctx context.Context, item *domain.FileAccess) error
	DeleteByID(ctx context.Context, id string) error
}
