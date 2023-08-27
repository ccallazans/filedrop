package repository

import (
	"context"
	"errors"

	"github.com/ccallazans/filedrop/internal/domain"
	"github.com/ccallazans/filedrop/internal/utils"
	"gorm.io/gorm"
)

type UserStore interface {
	DB() *gorm.DB
	FindAll(ctx context.Context) []*domain.User
	FindByID(ctx context.Context, id uint) (*domain.User, error)
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
	Save(ctx context.Context, user *domain.User) error
	DeleteByID(ctx context.Context, uuid uint) error
}

// IMPLEMENTATION

type PostgresUserStore struct {
	db *gorm.DB
}

func NewPostgresUserStore(db *gorm.DB) *PostgresUserStore {
	return &PostgresUserStore{
		db: db,
	}
}

func (r *PostgresUserStore) DB() *gorm.DB {
	return r.db
}

func (r *PostgresUserStore) FindByID(ctx context.Context, id uint) (*domain.User, error) {
	user := &domain.User{}

	tx := HasTransaction(ctx, r.db)
	err := tx.WithContext(ctx).Preload("Role").Where("id = ?", id).First(user).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			utils.Logger.Errorf("error when find user by id %d: %w", id, err)
		}
		return nil, err
	}

	return user, nil
}

func (r *PostgresUserStore) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	user := &domain.User{}

	tx := HasTransaction(ctx, r.db)
	err := tx.WithContext(ctx).Preload("Role").Where("email = ?", email).First(user).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			utils.Logger.Errorf("error when find user by email %s: %w", email, err)
		}
		return nil, err
	}

	return user, nil
}

func (r *PostgresUserStore) FindAll(ctx context.Context) []*domain.User {
	users := []*domain.User{}

	tx := HasTransaction(ctx, r.db)
	tx.WithContext(ctx).Preload("Role").Find(&users)

	return users
}

func (r *PostgresUserStore) Save(ctx context.Context, user *domain.User) error {

	tx := HasTransaction(ctx, r.db)
	err := tx.WithContext(ctx).Save(user).Error
	if err != nil {
		utils.Logger.Errorf("error when save user: %w", err)
		return err
	}

	return nil
}

func (r *PostgresUserStore) DeleteByID(ctx context.Context, id uint) error {
	user := &domain.User{}

	tx := HasTransaction(ctx, r.db)
	err := tx.WithContext(ctx).Delete(user, id).Error
	if err != nil {
		utils.Logger.Errorf("error when delete user by id %d: %w", id, err)
		return err
	}

	return nil
}
