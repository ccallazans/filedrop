package repository

import (
	"github.com/ccallazans/filedrop/internal/domain"
	"github.com/google/uuid"
)

type IFile interface {
	Save(file *domain.File) error
	FindByUUID(uuid uuid.UUID) (*domain.File, error)
}
