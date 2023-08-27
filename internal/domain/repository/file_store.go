package repository

import (
	"context"
	"errors"

	"github.com/ccallazans/filedrop/internal/domain"
	"github.com/ccallazans/filedrop/internal/utils"
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

func NewPostgresFileStore(db *gorm.DB) *PostgresFileStore {
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
	err := tx.WithContext(ctx).Preload("User").Where("id = ?", id).First(file).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			utils.Logger.Errorf("error when find file by id %d: %w", id, err)
		}
		return nil, err
	}

	return file, nil
}

func (r *PostgresFileStore) Save(ctx context.Context, file *domain.File) error {

	tx := HasTransaction(ctx, r.db)
	err := tx.WithContext(ctx).Preload("User").Save(file).Error
	if err != nil {
		utils.Logger.Errorf("error when save file: %w", err)
		return err
	}

	return nil
}

func (r *PostgresFileStore) DeleteByID(ctx context.Context, id uint) error {
	file := &domain.File{}

	tx := HasTransaction(ctx, r.db)
	err := tx.WithContext(ctx).Preload("User").Delete(file, id).Error
	if err != nil {
		utils.Logger.Errorf("error when delete file by id %d: %w", id, err)
		return err
	}

	return nil
}
