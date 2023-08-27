package utils

type AuthenticationError struct {
	Message string
}

type AuthorizationError struct {
	Message string
}

type BadRequestError struct {
	Message string
}

type ValidationError struct {
	Message string
}

type ConflictError struct {
	Message string
}

type NoContentError struct {
	Message string
}

type NotFoundError struct {
	Message string
}

type InternalError struct {
	Message string
}

func (v *AuthenticationError) Error() string {
	return v.Message
}

func (v *AuthorizationError) Error() string {
	return v.Message
}

func (v *BadRequestError) Error() string {
	return v.Message
}

func (v *ValidationError) Error() string {
	return v.Message
}

func (v *ConflictError) Error() string {
	return v.Message
}

func (v *NoContentError) Error() string {
	return v.Message
}

func (v *NotFoundError) Error() string {
	return v.Message
}

func (v *InternalError) Error() string {
	return v.Message
}
