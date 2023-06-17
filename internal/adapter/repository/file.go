package repository

import (
	"database/sql"
	"log"

	"github.com/ccallazans/filedrop/internal/domain"
	"github.com/ccallazans/filedrop/internal/domain/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type fileRepository struct {
	DB *gorm.DB
}

func NewFileRepository(db *gorm.DB) repository.IFile {
	return &fileRepository{
		db: db,
	}
}

func (r *fileRepository) Begin() *gorm.DB {
	return r.db.Begin()
}

func (r *fileRepository) Rollback() error { return r.db.Rollback().Error }
func (r *fileRepository) Commit() error { return r.db.Commit().Error }

func (r *fileRepository) Save(file *domain.File) error {

	result := r.db.Create(&file)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *fileRepository) FindByUUID(uuid uuid.UUID) (*domain.File, error) {

	file := &domain.File{}

	result := r.db.First(file, "uuid = ?", uuid)
	if result.Error != nil {
		return &domain.File{}, result.Error
	}

	return file, nil
}
