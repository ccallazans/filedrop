package utils

import "fmt"

const (
	ValidationErr = "ValidationErr"
	AuthenticationErr = "AuthenticationErr"
	BadRequestErr = "BadRequestErr"
	InternalErr   = "InternalErr"
)

type ErrorType struct {
	Type    string
	Message string
}

func (e *ErrorType) Error() string {
	return fmt.Sprintf("%s -> %s", e.Type, e.Message)
}
