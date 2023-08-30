package repository

import (
	"context"

	"github.com/ccallazans/filedrop/internal/domain"

	"gorm.io/gorm"
)

type FileAccessStore interface {
	DB() *gorm.DB
	FindByID(ctx context.Context, id uint) (*domain.FileAccess, error)
	FindByHash(ctx context.Context, hash string) (*domain.FileAccess, error)
	Save(ctx context.Context, fileAccess *domain.FileAccess) error
	DeleteByID(ctx context.Context, id uint) error
}

// IMPLEMENTATION

type PostgresFileAccessStore struct {
	db *gorm.DB
}

func NewPostgresFileAccessStore(db *gorm.DB) FileAccessStore {
	return &PostgresFileAccessStore{
		db: db,
	}
}

func (r *PostgresFileAccessStore) DB() *gorm.DB {
	return r.db
}

func (r *PostgresFileAccessStore) FindByID(ctx context.Context, id uint) (*domain.FileAccess, error) {
	fileAccess := &domain.FileAccess{}

	tx := HasTransaction(ctx, r.db)
	err := tx.WithContext(ctx).Preload("File").Where("id = ?", id).First(fileAccess).Error
	if err != nil {
		return nil, err
	}

	return fileAccess, nil
}

func (r *PostgresFileAccessStore) FindByHash(ctx context.Context, hash string) (*domain.FileAccess, error) {
	fileAccess := &domain.FileAccess{}

	tx := HasTransaction(ctx, r.db)
	err := tx.WithContext(ctx).Preload("File").Where("hash = ?", hash).First(fileAccess).Error
	if err != nil {
		return nil, err
	}

	return fileAccess, nil
}

func (r *PostgresFileAccessStore) Save(ctx context.Context, fileAccess *domain.FileAccess) error {
	tx := HasTransaction(ctx, r.db)
	err := tx.WithContext(ctx).Preload("User").Save(fileAccess).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *PostgresFileAccessStore) DeleteByID(ctx context.Context, id uint) error {
	fileAccess := &domain.FileAccess{}

	tx := HasTransaction(ctx, r.db)
	err := tx.WithContext(ctx).Delete(fileAccess, id).Error
	if err != nil {
		return err
	}

	return nil
}
