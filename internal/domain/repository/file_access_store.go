package repository

import (
	"context"
	"errors"

	"github.com/ccallazans/filedrop/internal/domain"
	"github.com/ccallazans/filedrop/internal/utils"

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

func NewPostgresFileAccessStore(db *gorm.DB) *PostgresFileAccessStore {
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
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			utils.Logger.Errorf("error when find fileAccess by id %d: %w", id, err)
		}
		return nil, err
	}

	return fileAccess, nil
}

func (r *PostgresFileAccessStore) FindByHash(ctx context.Context, hash string) (*domain.FileAccess, error) {
	fileAccess := &domain.FileAccess{}

	tx := HasTransaction(ctx, r.db)
	err := tx.WithContext(ctx).Preload("File").Where("hash = ?", hash).First(fileAccess).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			utils.Logger.Errorf("error when find fileAccess by hash %s: %w", hash, err)
		}
		return nil, err
	}

	return fileAccess, nil
}

func (r *PostgresFileAccessStore) Save(ctx context.Context, fileAccess *domain.FileAccess) error {

	tx := HasTransaction(ctx, r.db)
	err := tx.WithContext(ctx).Preload("User").Save(fileAccess).Error
	if err != nil {
		utils.Logger.Errorf("error when save fileAccess: %w", err)
		return err
	}

	return nil
}

func (r *PostgresFileAccessStore) DeleteByID(ctx context.Context, id uint) error {
	fileAccess := &domain.FileAccess{}

	tx := HasTransaction(ctx, r.db)
	err := tx.WithContext(ctx).Delete(fileAccess, id).Error
	if err != nil {
		utils.Logger.Errorf("error when delete fileAccess by id %d: %w", id, err)
		return err
	}

	return nil
}
