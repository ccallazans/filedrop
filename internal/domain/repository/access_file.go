package repository

import "github.com/ccallazans/filedrop/internal/domain/valueobject"

type IAccessFile interface {
	Save(accessFile *valueobject.AccessFile) error
	FindByHash(hash string) (*valueobject.AccessFile, error)
}
