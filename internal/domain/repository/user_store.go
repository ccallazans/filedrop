package repository

import (
	"context"

	"github.com/ccallazans/filedrop/internal/domain"
	"gorm.io/gorm"
)

type UserStore interface {
	DB() *gorm.DB
	FindAll(ctx context.Context) []*domain.User
	FindByID(ctx context.Context, id string) (*domain.User, error)
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
	Save(ctx context.Context, user *domain.User) error
	DeleteByID(ctx context.Context, id string) error
}

// IMPLEMENTATION

type PostgresUserStore struct {
	db *gorm.DB
}

func NewPostgresUserStore(db *gorm.DB) UserStore {
	return &PostgresUserStore{
		db: db,
	}
}

func (r *PostgresUserStore) DB() *gorm.DB {
	return r.db
}

func (r *PostgresUserStore) FindByID(ctx context.Context, id string) (*domain.User, error) {
	user := &domain.User{}

	tx := HasTransaction(ctx, r.db)
	err := tx.WithContext(ctx).Preload("Role").First(&user, "id = ?", id)
	if err.Error != nil {
		if err.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		
		return nil, err.Error
	}

	return user, nil
}

func (r *PostgresUserStore) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	user := &domain.User{}

	tx := HasTransaction(ctx, r.db)
	err := tx.WithContext(ctx).Preload("Role").First(&user, "email = ?", email)
	if err.Error != nil {
		if err.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}

		return nil, err.Error
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
	err := tx.WithContext(ctx).Preload("Role").Save(user).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *PostgresUserStore) DeleteByID(ctx context.Context, id string) error {
	user := &domain.User{}

	tx := HasTransaction(ctx, r.db)
	err := tx.WithContext(ctx).Delete(user, id).Error
	if err != nil {
		return err
	}

	return nil
}
