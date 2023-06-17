package service

import (
	"context"
	"mime/multipart"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/ccallazans/filedrop/internal/application/service"
)

type S3ClientService struct {
	s3CLient *s3.Client
	Bucket   string
}

func NewS3ClientService() service.IS3Client {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return nil
	}

	client := s3.NewFromConfig(cfg)

	return &S3ClientService{
		s3CLient: client,
		Bucket:   os.Getenv("AWS_S3_BUCKET"),
	}
}

func (s *S3ClientService) Save(key string, file *multipart.FileHeader) (string, error) {
	readFile, err := file.Open()
	if err != nil {
		return "", err
	}
	defer readFile.Close()
	
	uploader := manager.NewUploader(s.s3CLient)
	uploader.

	result, err := uploader.Upload(context.Background(), &s3.PutObjectInput{
		Bucket: &s.Bucket,
		Key:    &key,
		Body:   readFile,
	})
	if err != nil {
		return "", err
	}

	return result.Location, nil
}

func (s *S3ClientService) Get(key string) (*aws.WriteAtBuffer, error) {
	downloader := manager.NewDownloader(s.s3CLient)

	buf := aws.NewWriteAtBuffer([]byte{})
	_, err := downloader.Download(context.TODO(), buf, &s3.GetObjectInput{
		Bucket: &s.Bucket,
		Key:    &key,
	})
	if err != nil {
		return nil, err
	}

	buf.Bytes()

	return buf, nil
}

// func readFileInMemory(fileHeader *multipart.FileHeader) (*[]byte, error) {
// 	file, err := fileHeader.Open()
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer file.Close()

// 	// Read the file contents into memory
// 	fileBytes, err := io.ReadAll(file)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &fileBytes, nil
// }
