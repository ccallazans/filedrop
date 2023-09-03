package repository

import (
	"context"

	"github.com/ccallazans/filedrop/internal/domain"
	"gorm.io/gorm"
)

type FileStore interface {
	DB() *gorm.DB
	FindByID(ctx context.Context, id uint) (*domain.File, error)
	Save(ctx context.Context, file *domain.File) error
	DeleteByID(ctx context.Context, id uint) error
}

// IMPLEMENTATION

type PostgresFileStore struct {
	db *gorm.DB
}

func NewPostgresFileStore(db *gorm.DB) FileStore {
	return &PostgresFileStore{
		db: db,
	}
}

func (r *PostgresFileStore) DB() *gorm.DB {
	return r.db
}

func (r *PostgresFileStore) FindByID(ctx context.Context, id uint) (*domain.File, error) {
	file := &domain.File{}

	tx := HasTransaction(ctx, r.db)
	err := tx.WithContext(ctx).Preload("User").Where("id = ?", id).Limit(1).Find(file).Error
	if err != nil {
		return nil, err
	}

	return file, nil
}

func (r *PostgresFileStore) Save(ctx context.Context, file *domain.File) error {
	tx := HasTransaction(ctx, r.db)
	err := tx.WithContext(ctx).Preload("User").Save(file).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *PostgresFileStore) DeleteByID(ctx context.Context, id uint) error {
	file := &domain.File{}

	tx := HasTransaction(ctx, r.db)
	err := tx.WithContext(ctx).Preload("User").Delete(file, id).Error
	if err != nil {
		return err
	}

	return nil
}
