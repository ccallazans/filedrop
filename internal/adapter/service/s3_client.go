package service

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3ClientService struct {
	s3CLient *s3.Client
	Bucket   string
}

func NewS3ClientService() *S3ClientService {
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

func (s *S3ClientService) Save(key *string, file io.Reader) {
	ctx := context.Background()

	uploader := manager.NewUploader(s.s3CLient)

	result, err := uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket: &s.Bucket,
		Key:    key,
		Body:   file,
	})
	if err != nil {
		panic(err)
	}

	log.Println("File Uploaded Successfully, URL : ", result.Location)
}

func (s *S3ClientService) Get(key string) {
	downloader := manager.NewDownloader(s.s3CLient)

	file, err := os.Create(key)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	numBytes, err := downloader.Download(context.TODO(), file, &s3.GetObjectInput{
		Bucket: &s.Bucket,
		Key:    &key,
	})

	fmt.Println("File donwloaded!: ", numBytes)
}
