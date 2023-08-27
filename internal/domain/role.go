package domain

const (
	ADMIN = 1
	USER = 2
)

type Role struct {
	ID   uint
	Role string
}