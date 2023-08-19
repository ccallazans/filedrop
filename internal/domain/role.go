package domain

const (
	ADMIN = 2
	USER = 1
)

type Role struct {
	ID   uint
	Role string
}