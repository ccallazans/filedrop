package domain

import "time"

type FileAccess struct {
	ID        uint
	Hash      string
	Secret    string
	FileID    uint
	CreatedAt time.Time
	UpdatedAt time.Time

	File File
}

func (FileAccess) TableName() string {
	return "file_access"
}
