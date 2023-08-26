package usecase

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"mime/multipart"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/ccallazans/filedrop/internal/domain"
	"github.com/ccallazans/filedrop/internal/domain/repository"
	"github.com/ccallazans/filedrop/internal/utils"
	"github.com/google/uuid"
)

const (
	MAX_FILE_SIZE = 1024 * 5
)

type UploadUsecase struct {
	fileRepo       repository.FileRepository
	fileAccessRepo repository.FileAccessRepository
	userRepo       repository.UserRepository
	s3Client       *s3.Client
}

func NewUploadUsecase(fileRepo repository.FileRepository, fileAccessRepo repository.FileAccessRepository, userRepo repository.UserRepository, s3Client *s3.Client) UploadUsecase {
	return UploadUsecase{
		fileRepo:       fileRepo,
		fileAccessRepo: fileAccessRepo,
		userRepo:       userRepo,
		s3Client:       s3Client,
	}
}

func (u *UploadUsecase) UploadFile(ctx context.Context, secret string, multiPartFile *multipart.FileHeader) (string, error) {

	ctxUser, err := utils.GetContextUser(ctx)
	if err != nil {
		log.Println("error uploading file")
		return "", err
	}
	// ctxUser := &domain.User{ID: 1}

	tx := u.fileRepo.DB().Begin()
	ctxTx := context.WithValue(ctx, "tx", tx)

	fileUUID := uuid.NewString()

	location, err := uploadFileToS3(ctx, u.s3Client, fileUUID, multiPartFile)
	if err != nil {
		log.Println("error uploadFileToS3: ", err)
		return "", err
	}

	file := &domain.File{
		UUID:     fileUUID,
		Filename: multiPartFile.Filename,
		Size:     fmt.Sprintf("%d", multiPartFile.Size),
		Location: location,
		UserID:   ctxUser.ID,
	}

	err = u.fileRepo.Save(ctxTx, file)
	if err != nil {
		tx.Rollback()
		return "", err
	}

	fileAccess := &domain.FileAccess{
		Hash:   generateRandomHash(5),
		Secret: secret,
		FileID: file.ID,
	}

	err = u.fileAccessRepo.Save(ctxTx, fileAccess)
	if err != nil {
		tx.Rollback()
		return "", err
	}

	tx.Commit()

	return fileAccess.Hash, nil
}

func uploadFileToS3(ctx context.Context, s3Client *s3.Client, fileUUID string, fileHeader *multipart.FileHeader) (string, error) {

	openFile, err := fileHeader.Open()
	if err != nil {
		log.Println(err)
		return "", &utils.ErrorType{Type: utils.InternalErr, Message: "could not open multipartfile!"}
	}
	defer openFile.Close()

	_, err = s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:             aws.String(os.Getenv("AWS_BUCKET")),
		Key:                aws.String(fileUUID),
		Body:               openFile,
		ContentDisposition: aws.String("attachment"),
	})
	if err != nil {
		log.Println(err)
		return "", &utils.ErrorType{Type: utils.InternalErr, Message: "could not upload file to s3 bucket"}
	}

	location := fmt.Sprintf("%s/%s", os.Getenv("AWS_BUCKET"), fileUUID)

	return location, nil
}

func (u *UploadUsecase) AccessFile(ctx context.Context, hash string, secret string) (*s3.GetObjectOutput, error) {

	validAccessFile, err := u.fileAccessRepo.FindByHash(ctx, hash)
	if err != nil {
		log.Println(err)
		return nil, &utils.ErrorType{Type: utils.ValidationErr, Message: "fileAccess do not exist!"}
	}

	if secret != validAccessFile.Secret {
		log.Println(err)
		return nil, &utils.ErrorType{Type: utils.ValidationErr, Message: "invalid secret!"}
	}

	file, err := u.fileRepo.FindByUUID(ctx, validAccessFile.File.UUID)
	if err != nil {
		log.Println(err)
		return nil, &utils.ErrorType{Type: utils.ValidationErr, Message: "file do not exist!"}
	}

	bufferFile, err := u.s3Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(os.Getenv("AWS_BUCKET")),
		Key:    aws.String(file.UUID),
	})

	if err != nil {
		log.Println(err)
		return nil, errors.New("error getting file from s3")
	}

	return bufferFile, nil
}

func generateRandomHash(length int) string {

	characters := []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9",
		"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z"}

	hash := ""
	for i := 0; i < length; i++ {
		hash += characters[rand.Intn(len(characters))]
	}

	return hash
}
