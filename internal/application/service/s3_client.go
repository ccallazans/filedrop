package service

import (
	"mime/multipart"

	"github.com/aws/aws-sdk-go/aws"
)

type IS3Client interface {
	Save(key string, file *multipart.File) (string, error)
	Get(key string) (*aws.WriteAtBuffer, error)
}
