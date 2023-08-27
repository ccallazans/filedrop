package domain

import (
	"time"
)

type User struct {
	ID        uint
	FirstName string
	LastName  string
	Email     string
	Password  string
	RoleID    uint
	CreatedAt time.Time
	UpdatedAt time.Time

	Role Role
}

func NewUser(firstName string, lastName string, email string, password string) *User {
	return &User{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Password:  password,
		RoleID:    USER,

		Role: Role{ID: USER, Role: "USER"},
	}
}
