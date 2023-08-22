package repository

import (
	"context"
	"fmt"

	"github.com/ccallazans/filedrop/internal/domain"
	"github.com/ccallazans/filedrop/internal/domain/repository"

	"gorm.io/gorm"
)

type fileAccessRepositoryImpl struct {
	db *gorm.DB
}

func NewFileAccessRepository(db *gorm.DB) repository.FileAccessRepository {
	return &fileAccessRepositoryImpl{
		db: db,
	}
}

func (r *fileAccessRepositoryImpl) DB() *gorm.DB {
	return r.db
}

func (r *fileAccessRepositoryImpl) FindByHash(ctx context.Context, hash string) (*domain.FileAccess, error) {
	fileAccess := &domain.FileAccess{}
	err := r.db.WithContext(ctx).Where("hash = ?", hash).First(fileAccess).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find fileAccess by hash: %w", err)
	}

	return fileAccess, nil
}

func (r *fileAccessRepositoryImpl) Save(ctx context.Context, item *domain.FileAccess) error {

	tr := r.db
	hasTransaction := ctx.Value("tx")
	if hasTransaction != nil {
		tr = hasTransaction.(*gorm.DB)
	}

	err := tr.WithContext(ctx).Save(item).Error
	if err != nil {
		return fmt.Errorf("failed to save fileAccess: %w", err)
	}

	return nil
}

func (r *fileAccessRepositoryImpl) DeleteByID(ctx context.Context, id string) error {

	tr := r.db
	hasTransaction := ctx.Value("tx")
	if hasTransaction != nil {
		tr = hasTransaction.(*gorm.DB)
	}

	fileAccess := &domain.FileAccess{}
	err := tr.WithContext(ctx).Delete(fileAccess, id).Error
	if err != nil {
		return fmt.Errorf("failed to delete fileAccess by id: %w", err)
	}

	return nil
}
