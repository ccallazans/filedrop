package domain

import (
	"time"
)

type File struct {
	ID        string
	Filename  string
	Password  string
	Location  string
	Hash      string
	IsActive  bool
	UserID    string
	CreatedAt time.Time
	UpdatedAt time.Time
}
