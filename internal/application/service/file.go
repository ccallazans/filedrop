package service

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
)

const (
	MAX_FILE_SIZE   = 1024 * 5
	ErrFileNotFound = "file not found"
)

type FileService struct {
	fileStore repository.FileStore
	userStore repository.UserStore
	s3Client  *s3.Client
}

func NewFileService(fileStore repository.FileStore, userStore repository.UserStore, s3Client *s3.Client) *FileService {
	return &FileService{
		fileStore: fileStore,
		userStore: userStore,
		s3Client:  s3Client,
	}
}

func (s *FileService) Upload(ctx context.Context, password string, multiPartFile *multipart.FileHeader) (string, error) {
	ctxUser, err := getContextUser(ctx)
	if err != nil {
		return "", err
	}

	tx := s.fileStore.DB().Begin()
	ctxTx := context.WithValue(ctx, "tx", tx)

	location, err := uploadFileToS3(ctxTx, s.s3Client, multiPartFile)
	if err != nil {
		return "", err
	}

	hash := generateRandomHash(6)
	file := domain.NewFile(multiPartFile.Filename, password, location, hash, ctxUser.ID)

	err = s.fileStore.Save(ctxTx, file)
	if err != nil {
		tx.Rollback()
		return "", err
	}

	tx.Commit()

	return hash, nil
}

func (s *FileService) DownloadFile(ctx context.Context, hash string, password string) (*s3.GetObjectOutput, string, error) {
	file, err := s.fileStore.FindByHash(ctx, hash)
	if err != nil {
		return nil, "", fmt.Errorf(ErrFileNotFound)
	}

	if password != file.Password {
		return nil, "", fmt.Errorf(ErrInvalidPassword)
	}

	bufferFile, err := s.s3Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(os.Getenv("AWS_BUCKET")),
		Key:    aws.String(strings.Split(file.Location, "/")[1]),
	})
	if err != nil {
		return nil, "", err
	}

	return bufferFile, file.Filename, nil
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
