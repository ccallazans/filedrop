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
	"gorm.io/gorm"
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

func (u *UploadUsecase) WithTrx(trxHandle *gorm.DB) UploadUsecase {
	u.userRepo = u.userRepo.WithTrx(trxHandle)
	fileRepo.WithTrx(trxHandle)
	accessFileRepo.WithTrx(trxHandle)
	return u
}

type UploadFileArgs struct {
	Lock       bool
	AccessCode string
	File       *multipart.FileHeader
}

func (u *UploadUsecase) UploadFile(args *UploadFileArgs) error {

	// if args.File.Size > MAX_FILE_SIZE {
	// 	return fmt.Errorf("max size allowed: %d, uploaded file: %d", MAX_FILE_SIZE, args.File.Size)
	// }

	txFileRepo, err := u.fileRepo.Begin()
	if err != nil {
		return err
	}
	txAccessFileRepo, err := u.accessFileRepo.Begin()
	if err != nil {
		u.fileRepo.Rollback()
		return err
	}

	u.fileRepo

	defer func() {
		if r := recover(); r != nil {
			u.fileRepo.Rollback()
			u.accessFileRepo.Rollback()
			fmt.Println("Recovered from panic during transaction:", r)
		} else if err != nil {
			u.fileRepo.Rollback()
			u.accessFileRepo.Rollback()
			fmt.Println("Rolling back transaction due to error:", err)
		} else {
			u.fileRepo.Commit()
			u.accessFileRepo.Commit()
		}
	}()

	location, err := u.s3Client.Save(args.File.Filename, args.File)
	if err != nil {
		return err
	}

	fileSize := strconv.FormatFloat(bytesToMegabytes(args.File.Size), 'f', -1, 64)
	saveFile := &domain.File{
		UUID:        uuid.New(),
		Filename:    args.File.Filename,
		Size:        fileSize,
		LocationURL: location,
		UserUUID:    uuid.MustParse("74318875-6aca-4bdc-a00b-4d2c5d38dd0f"),
	}

	hash, err := generateRandomHash(6)
	if err != nil {
		u.fileRepo.Rollback()
		return err
	}

	saveAccessFile := &valueobject.AccessFile{
		Hash:       hash,
		Lock:       args.Lock,
		AccessCode: args.AccessCode,
		FileUUID:   uuid.MustParse("74318875-6aca-4bdc-a00b-4d2c5d38dd0f"),
	}

	err = u.fileRepo.Save(saveFile)
	if err != nil {
		u.fileRepo.Rollback()
		return err
	}

	err = u.accessFileRepo.Save(saveAccessFile)
	if err != nil {
		fmt.Println("---------------------------------> AAAAAAAAAAAA")
		u.fileRepo.Rollback()
		return err
	}

	err = u.fileRepo.Commit()
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
