package service

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"mime/multipart"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/ccallazans/filedrop/internal/domain"
	"github.com/ccallazans/filedrop/internal/domain/repository"
)

const (
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
	ctxUser, err := GetUserFromCtx(ctx)
	if err != nil {
		return "", err
	}

	location, err := uploadFileToS3(ctx, s.s3Client, multiPartFile)
	if err != nil {
		return "", err
	}

	hashedPassword, err := hash(password)
	if err != nil {
		return "", err
	}

	hashUrl := generateRandomHash(6)

	file := &domain.File{
		ID:       uuid.NewString(),
		Filename: multiPartFile.Filename,
		Password: hashedPassword,
		Location: location,
		Hash:     hashUrl,
		IsActive: true,
		UserID:   ctxUser.ID,
	}

	err = s.fileStore.Save(ctx, file)
	if err != nil {
		return "", err
	}

	return hashUrl, nil
}

func (s *FileService) DownloadFile(ctx context.Context, hash string, password string) (*s3.GetObjectOutput, string, error) {
	exists, err := s.fileStore.Exists(ctx, hash)
	if err != nil {
		return nil, "", err
	}

	if !exists {
		return nil, "", errors.New(ErrFileNotFound)
	}

	file, err := s.fileStore.FindByHash(ctx, hash)
	if err != nil {
		return nil, "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(file.Password), []byte(password))
	if err != nil {
		return nil, "", errors.New(ErrInvalidPassword)
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
