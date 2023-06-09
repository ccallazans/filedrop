package repository

import "github.com/ccallazans/filedrop/internal/domain"

type IUser interface {
	Save(user *domain.User) error
	FindByEmail(email string) (domain.User, error)
}