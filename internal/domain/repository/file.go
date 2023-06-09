package repository

import "github.com/ccallazans/filedrop/internal/domain"

type IFile interface {
	Save(file *domain.File) error
	FindById(id int) (domain.File, error)
}