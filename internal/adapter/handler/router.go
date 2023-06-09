package handler

import (
	"github.com/ccallazans/filedrop/internal/application/usecase"
	"github.com/ccallazans/filedrop/internal/domain/repository"
	"gorm.io/gorm"
)

func NewRouter(db *gorm.DB) {

	// Repositories
	userRepository := repository.IUser
	fileRepository := repository.IFile

	// Usecases
	accountUsecase := usecase.NewAccountUsecase(userRepository, fileRepository)
	uploadUsecase := usecase.NewUploadUsecase(userRepository, fileRepository)
}
