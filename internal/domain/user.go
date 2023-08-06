package domain

import (
	"time"
)

type User struct {
	ID        uint
	UUID      string
	Email     string
	Password  string
	RoleID    uint
	CreatedAt time.Time
	UpdatedAt time.Time
}
