package domain

type RoleEnum uint

const (
	ADMIN RoleEnum = 1
	USER  RoleEnum = 2
	GUEST RoleEnum = 3
)

type Role struct {
	ID   uint
}
