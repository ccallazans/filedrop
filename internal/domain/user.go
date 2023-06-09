package domain

import "github.com/google/uuid"

type User struct {
	UUID     uuid.UUID
	Name     string
	Email    string
	Password string
	Role     UserRole
}

type UserRole string

const (
	USER  UserRole = "USER"
	ADMIN UserRole = "ADMIN"
)
