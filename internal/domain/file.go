package domain

import (
	"time"
)

type File struct {
	ID        uint
	UUID      string
	Filename  string
	Size      string
	Location  string
	UserID    uint
	CreatedAt time.Time
	UpdatedAt time.Time
}
