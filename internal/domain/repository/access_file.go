package repository

import (
	"database/sql"

	"github.com/ccallazans/filedrop/internal/domain/valueobject"
	"gorm.io/gorm"
)

type IAccessFile interface {
	Begin() (*sql.DB, error)
	Rollback() error
	Commit() error
	DB()
	Save(accessFile *valueobject.AccessFile) error
	FindByHash(hash string) (*valueobject.AccessFile, error)
}
