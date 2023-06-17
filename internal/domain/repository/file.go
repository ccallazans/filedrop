package repository

import (
	"github.com/ccallazans/filedrop/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IFile interface {
	Begin() (*gorm.DB, error) 
	Rollback() error
	Commit() error
	Save(file *domain.File) error
	FindByUUID(uuid uuid.UUID) (*domain.File, error)
}
