package usecase

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"mime/multipart"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/ccallazans/filedrop/internal/application/service"
	"github.com/ccallazans/filedrop/internal/domain"
	"github.com/ccallazans/filedrop/internal/domain/repository"
	"github.com/ccallazans/filedrop/internal/domain/valueobject"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

const (
	MAX_FILE_SIZE = 1024 * 5
)

type UploadUsecase struct {
	userRepo       repository.IUser
	fileRepo       repository.IFile
	accessFileRepo repository.IAccessFile
	s3Client       service.IS3Client
}

func NewUploadUsecase(userRepo repository.IUser, fileRepo repository.IFile, accessFileRepo repository.IAccessFile, s3Client service.IS3Client) UploadUsecase {
	return UploadUsecase{
		userRepo:       userRepo,
		fileRepo:       fileRepo,
		accessFileRepo: accessFileRepo,
		s3Client:       s3Client,
	}
}

type UploadFileArgs struct {
	Lock       bool
	AccessCode string
	File       *multipart.FileHeader
}

func (u *UploadUsecase) UploadFile(args *UploadFileArgs) error {

	if args.File.Size > MAX_FILE_SIZE {
		return fmt.Errorf("max size allowed: %d, uploaded file: %d", MAX_FILE_SIZE, args.File.Size)
	}

	readFile, err := args.File.Open()
	if err != nil {
		return err
	}
	defer readFile.Close()

	location, err := u.s3Client.Save(args.File.Filename, &readFile)
	if err != nil {
		return err
	}

	fileSize := strconv.FormatFloat(bytesToMegabytes(args.File.Size), 'f', -1, 64)

	saveFile := &domain.File{
		UUID:        uuid.New(),
		Filename:    args.File.Filename,
		Size:        fileSize,
		LocationURL: location,
	}

	hash, err := generateRandomHash(6)
	if err != nil {
		return err
	}

	saveAccessFile := &valueobject.AccessFile{
		Hash:       hash,
		Lock:       args.Lock,
		AccessCode: args.AccessCode,
		FileUUID:   saveFile.UUID,
	}

	err = u.fileRepo.Save(saveFile)
	if err != nil {
		return err
	}

	err = u.accessFileRepo.Save(saveAccessFile)
	if err != nil {
		return err
	}

	return nil
}

type AccessFileArgs struct {
	AccessCode string
}

func (u *UploadUsecase) AccessFile(hash string, args AccessFileArgs) (*aws.WriteAtBuffer, error) {

	validAccessFile, err := u.accessFileRepo.FindByHash(hash)
	if err != nil {
		return nil, errors.New("invalid file")
	}

	if validAccessFile.Lock {
		err = bcrypt.CompareHashAndPassword([]byte(validAccessFile.AccessCode), []byte(args.AccessCode))
		if err != nil {
			return nil, errors.New("invalid access code")
		}
	}

	file, err := u.fileRepo.FindByUUID(validAccessFile.FileUUID)
	if err != nil {
		return nil, errors.New("could not find file")
	}

	bufferFile, err := u.s3Client.Get(file.Filename)
	if err != nil {
		return nil, errors.New("could not download file")
	}

	return bufferFile, nil
}

func generateRandomHash(length int) (string, error) {
	// Calculate the number of bytes needed to generate the hash
	byteLength := (length * 3) / 4

	// Create a byte slice to hold the random bytes
	bytes := make([]byte, byteLength)

	// Read random bytes from the crypto/rand package
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	// Encode the random bytes to base64
	hash := base64.URLEncoding.EncodeToString(bytes)

	// Remove any padding characters from the base64 string
	hash = hash[:length]

	return hash, nil
}

func bytesToMegabytes(bytes int64) float64 {
	return float64(bytes) / (1024 * 1024)
}
