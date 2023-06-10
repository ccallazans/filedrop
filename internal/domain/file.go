package domain

import "github.com/google/uuid"

type File struct {
	UUID        uuid.UUID
	Filename    string
	Size        string
	LocationURL string
	UserUUID    uuid.UUID
}
