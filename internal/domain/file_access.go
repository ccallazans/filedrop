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

func NewFileAccess(hash string, secret string, fileID uint) *FileAccess {
	return &FileAccess{
		Hash:   hash,
		Secret: secret,
		FileID: fileID,
	}
}
