package domain

import (
	"time"
)

type File struct {
	ID        uint
	Filename  string
	Size      string
	Location  string
	UserID    uint
	User      User
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewFile(filename string, size string, location string, userID uint) *File {
	return &File{
		Filename: filename,
		Size:     size,
		Location: location,
		UserID:   userID,
	}
}
