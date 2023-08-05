package domain

import "github.com/google/uuid"

type User struct {
	ID       string
	UUID     string
	Email    string
	Password string
}

type UserRole string

const (
	USER  UserRole = "USER"
	ADMIN UserRole = "ADMIN"
)
