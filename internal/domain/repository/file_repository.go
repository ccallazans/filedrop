package repository

import (
	"context"

	"github.com/ccallazans/filedrop/internal/domain"
	"gorm.io/gorm"
)

type FileRepository interface {
	DB() *gorm.DB
	FindByUUID(ctx context.Context, uuid string) (*domain.File, error)
	Save(ctx context.Context, item *domain.File) error
	DeleteByUUID(ctx context.Context, uuid string) error
}
