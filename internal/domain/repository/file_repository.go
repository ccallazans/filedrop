package repository

import (
	"context"

	"github.com/ccallazans/filedrop/internal/domain"
)

type FileRepository interface {
	FindByUUID(ctx context.Context, uuid string) (*domain.File, error)
	CreateFile(ctx context.Context, item *domain.File) error
	DeleteItem(ctx context.Context, uuid string) error
}
