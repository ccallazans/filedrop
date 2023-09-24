package domain

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        string
	FirstName string
	LastName  string
	Email     string
	Password  string
	RoleID    uint
	Role      Role
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewUser(firstName string, lastName string, email string, password string) (*User, error) {
	password, err := hashPassword(password)
	if err != nil {
		return nil, err
	}

	return &User{
		ID:        uuid.NewString(),
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Password:  password,
		RoleID:    uint(USER),
	}, nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}
