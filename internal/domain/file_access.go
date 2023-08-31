package domain

import "time"

type FileAccess struct {
	ID        uint
	Hash      string
	Secret    string
	FileID    uint
	File      File
	CreatedAt time.Time
	UpdatedAt time.Time
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
