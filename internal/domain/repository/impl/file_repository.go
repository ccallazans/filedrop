package repository

import (
	"context"
	"fmt"

	"github.com/ccallazans/filedrop/internal/domain"
	"github.com/ccallazans/filedrop/internal/domain/repository"
	"gorm.io/gorm"
)

type fileRepositoryImpl struct {
	db *gorm.DB
}

func NewFileRepository(db *gorm.DB) repository.FileRepository {
	return &fileRepositoryImpl{
		db: db,
	}
}

func (r *fileRepositoryImpl) DB() *gorm.DB {
	return r.db
}

func (r *fileRepositoryImpl) FindByUUID(ctx context.Context, uuid string) (*domain.File, error) {
	file := &domain.File{}
	err := r.db.WithContext(ctx).Where("uuid = ?", uuid).First(file).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find file by uuid: %w", err)
	}

	return file, nil
}

func (r *fileRepositoryImpl) Save(ctx context.Context, item *domain.File) error {

	tr := r.db
	hasTransaction := ctx.Value("tx")
	if hasTransaction != nil {
		tr = hasTransaction.(*gorm.DB)
	}

	err := tr.WithContext(ctx).Save(item).Error
	if err != nil {
		return fmt.Errorf("failed to save item: %w", err)
	}

	return nil
}

func (r *fileRepositoryImpl) DeleteByUUID(ctx context.Context, uuid string) error {

	tr := r.db
	hasTransaction := ctx.Value("tx")
	if hasTransaction != nil {
		tr = hasTransaction.(*gorm.DB)
	}

	file := &domain.File{}
	err := tr.WithContext(ctx).Delete(file, uuid).Error
	if err != nil {
		return fmt.Errorf("failed to delete file by id: %w", err)
	}

	return nil
}
