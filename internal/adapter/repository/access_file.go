package repository

import (
	"github.com/ccallazans/filedrop/internal/domain/repository"
	"github.com/ccallazans/filedrop/internal/domain/valueobject"
	"gorm.io/gorm"
)

type accessFileRepository struct {
	db *gorm.DB
}

func NewAccessFileRepository(db *gorm.DB) repository.IAccessFile {
	return &accessFileRepository{
		db: db,
	}
}

func (r *accessFileRepository) Save(accessFile *valueobject.AccessFile) error {

	result := r.db.Create(&accessFile)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *accessFileRepository) FindByHash(hash string) (*valueobject.AccessFile, error) {

	accessFile := &valueobject.AccessFile{}

	result := r.db.First(accessFile, "hash = ?", hash)
	if result.Error != nil {
		return &valueobject.AccessFile{}, result.Error
	}

	return accessFile, nil

}
