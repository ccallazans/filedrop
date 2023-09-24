package domain

import (
	"time"

	"github.com/google/uuid"
)

type File struct {
	ID        string
	Filename  string
	Password  string
	Location  string
	Hash      string
	IsActive  bool
	UserID    string
	User      User
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewFile(filename string, password string, location string, hash string, userId string) *File {
	return &File{
		ID: uuid.NewString(),
		Filename: filename,
		Password: password,
		Location: location,
		Hash:     hash,
		IsActive: true,
		UserID:   userId,
	}
}
