package usecase

import (
	"context"
	"fmt"
	"math/rand"
	"mime/multipart"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"

	"github.com/ccallazans/filedrop/internal/domain"
	"github.com/ccallazans/filedrop/internal/domain/repository"
	"github.com/ccallazans/filedrop/internal/utils"
)

const (
	MAX_FILE_SIZE = 1024 * 5
)

type FileUsecase struct {
	fileStore       repository.FileStore
	fileAccessStore repository.FileAccessStore
	userStore       repository.UserStore
	s3Client        *s3.Client
}

func NewFileUsecase(fileStore repository.FileStore, fileAccessStore repository.FileAccessStore, userStore repository.UserStore, s3Client *s3.Client) *FileUsecase {
	return &FileUsecase{
		fileStore:       fileStore,
		fileAccessStore: fileAccessStore,
		userStore:       userStore,
		s3Client:        s3Client,
	}
}

func (u *FileUsecase) UploadFile(ctx context.Context, secret string, multiPartFile *multipart.FileHeader) (string, error) {
	ctxUser, err := GetContextUser(ctx)
	if err != nil {
		return "", err
	}

	tx := u.fileStore.DB().Begin()
	ctxTx := context.WithValue(ctx, "tx", tx)

	location, err := uploadFileToS3(ctxTx, u.s3Client, multiPartFile)
	if err != nil {
		return "", err
	}

	file := domain.NewFile(multiPartFile.Filename, fmt.Sprintf("%d", multiPartFile.Size), location, ctxUser.ID)

	err = u.fileStore.Save(ctxTx, file)
	if err != nil {
		tx.Rollback()
		return "", err
	}

	fileAccess := domain.NewFileAccess(generateRandomHash(5), secret, file.ID)

	err = u.fileAccessStore.Save(ctxTx, fileAccess)
	if err != nil {
		tx.Rollback()
		return "", err
	}

	tx.Commit()

	return fileAccess.Hash, nil
}

func (u *FileUsecase) DownloadFile(ctx context.Context, hash string, secret string) (*s3.GetObjectOutput, error) {
	validAccessFile, err := u.fileAccessStore.FindByHash(ctx, hash)
	if err != nil {
		return nil, &utils.NotFoundError{Message: "hash access does not exist"}
	}

	if secret != validAccessFile.Secret {
		return nil, &utils.AuthenticationError{Message: "invalid secret"}
	}

	file, err := u.fileStore.FindByID(ctx, validAccessFile.FileID)
	if err != nil {
		return nil, err
	}

	bufferFile, err := u.s3Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(os.Getenv("AWS_BUCKET")),
		Key:    aws.String(strings.Split(file.Location, "/")[1]),
	})
	if err != nil {
		return nil, err
	}

	return bufferFile, nil
}

func uploadFileToS3(ctx context.Context, s3Client *s3.Client, fileHeader *multipart.FileHeader) (string, error) {
	openFile, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer openFile.Close()

	key := uuid.NewString()

	_, err = s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:             aws.String(os.Getenv("AWS_BUCKET")),
		Key:                aws.String(key),
		Body:               openFile,
		ContentDisposition: aws.String("attachment"),
	})
	if err != nil {
		return "", err
	}

	location := fmt.Sprintf("%s/%s", os.Getenv("AWS_BUCKET"), key)

	return location, nil
}

func generateRandomHash(length int) string {
	characters := []string{
		"0", "1", "2", "3", "4", "5", "6", "7", "8", "9",
		"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z",
	}

	hash := ""
	for i := 0; i < length; i++ {
		hash += characters[rand.Intn(len(characters))]
	}

	return hash
}
