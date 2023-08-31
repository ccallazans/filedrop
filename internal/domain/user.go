package domain

import (
	"time"

	"github.com/ccallazans/filedrop/internal/utils"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        uint
	FirstName string `validate:"required,lte=255"`
	LastName  string `validate:"required,lte=255"`
	Email     string `validate:"required,email"`
	Password  string `validate:"required,gte=6,lte=255"`
	RoleID    uint   `validate:"required"`
	CreatedAt time.Time
	UpdatedAt time.Time

	Role Role
}

func NewUser(firstName string, lastName string, email string, password string) (*User, error) {
	user := &User{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Password:  password,
		RoleID:    USER,
	}

	err := validator.New().Struct(user)
	if err != nil {
		return nil, &utils.ValidationError{Message: err.Error()}
	}

	user.Password, err = hashPassword(password)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}
