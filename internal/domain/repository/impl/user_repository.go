package repository

import (
	"context"
	"fmt"

	"github.com/ccallazans/filedrop/internal/domain"
	"github.com/ccallazans/filedrop/internal/domain/repository"
	"gorm.io/gorm"
)

type userRepositoryImpl struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) repository.UserRepository {
	return &userRepositoryImpl{
		db: db,
	}
}

func (r *userRepositoryImpl) DB() *gorm.DB {
	return r.db
}

func (r *userRepositoryImpl) FindByUUID(ctx context.Context, uuid string) (*domain.User, error) {
	user := &domain.User{}
	err := r.db.WithContext(ctx).Preload("Role").Where("uuid = ?", uuid).First(user).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find user by uuid: %w", err)
	}

	return user, nil
}

func (r *userRepositoryImpl) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	user := &domain.User{}
	err := r.db.WithContext(ctx).Preload("Role").Where("email = ?", email).First(user).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find user by email: %w", err)
	}

	return user, nil
}

func (r *userRepositoryImpl) Save(ctx context.Context, user *domain.User) error {

	tr := r.db
	hasTransaction := ctx.Value("tx")
	if hasTransaction != nil {
		tr = hasTransaction.(*gorm.DB)
	}

	err := tr.WithContext(ctx).Save(user).Error
	if err != nil {
		return fmt.Errorf("failed to save user: %w", err)
	}

	return nil
}

func (r *userRepositoryImpl) DeleteByUUID(ctx context.Context, uuid string) error {

	tr := r.db
	hasTransaction := ctx.Value("tx")
	if hasTransaction != nil {
		tr = hasTransaction.(*gorm.DB)
	}

	user := &domain.User{}
	err := tr.WithContext(ctx).Delete(user, uuid).Error
	if err != nil {
		return fmt.Errorf("failed to delete user by id: %w", err)
	}

	return nil
}
