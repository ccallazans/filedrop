package service

import "io"

type IS3Client interface {
	Save(key *string, file io.Reader)
	Get(key string)
}
