package repository

import (
	"github.com/ccallazans/filedrop/internal/domain"
	"github.com/ccallazans/filedrop/internal/domain/repository"
	"gorm.io/gorm"
)

type fileRepository struct {
	db *gorm.DB
}

func NewFileRepository(db *gorm.DB) repository.IFile {
	return &fileRepository{
		db: db,
	}
}

func (r *fileRepository) Save(file *domain.File) error {

	result := r.db.Create(&file)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *fileRepository) FindById(id int) (*domain.File, error) {

	file := &domain.File{}

	result := r.db.First(file, "id = ?", id)
	if result.Error != nil {
		return &domain.File{}, result.Error
	}

	return file, nil
}
