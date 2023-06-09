package usecase

import (
	"github.com/ccallazans/filedrop/internal/application/service"
	"github.com/ccallazans/filedrop/internal/domain/repository"
)

type UploadUsecase struct {
	userRepo repository.IUser
	fileRepo repository.IFile
	s3Client service.IS3Client
}

func NewUploadUsecase(userRepo repository.IUser, fileRepo repository.IFile, s3Client service.IS3Client) *UploadUsecase {
	return &UploadUsecase{
		userRepo: userRepo,
		fileRepo: fileRepo,
		s3Client: s3Client,
	}
}

type UploadFileArgs struct {
	// File
}

func (u *UploadUsecase) UploadFile() {
	// arquivo
	// publico
	// privado

	// gerar hash do arquivo
	// salvar usuario->hash e hash->url do arquivo
}

func (u *UploadUsecase) AccessFile() {
	// hash do arquivo
	// verificar permissao
	//

}

func SetAccessConditions() {

}
