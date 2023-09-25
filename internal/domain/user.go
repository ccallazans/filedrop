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
	Role      Role
	CreatedAt time.Time
	UpdatedAt time.Time
}
