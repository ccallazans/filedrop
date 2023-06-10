package valueobject

import "github.com/google/uuid"

type AccessFile struct {
	Hash       string
	Lock       bool
	AccessCode string
	FileUUID     uuid.UUID
}