package domain

const (
	ADMIN RoleEnum = 1
	USER  RoleEnum = 2
	GUEST RoleEnum = 3
)

type RoleEnum uint

type Role struct {
	ID   uint
	Role string
}
