package domain

import (
	"time"
)

type User struct {
	ID        string
	FirstName string
	Email     string
	Password  string
	RoleID    uint
	CreatedAt time.Time
	UpdatedAt time.Time
}
