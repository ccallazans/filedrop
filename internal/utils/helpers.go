package utils

import (
	"context"

	"github.com/ccallazans/filedrop/internal/application/auth"
)

const (
	ERROR_GET_USER_CONTEXT = "Could not retrieve user from context"
)

func GetContextUser(ctx context.Context) (*auth.JWTUser, error) {
	ctxValue := ctx.Value("user")
	if ctxValue == nil {
		return nil, &ErrorType{Type: InternalErr, Message: ERROR_GET_USER_CONTEXT}
	}

	jwtUser := ctxValue.(*auth.JWTUser)

	return jwtUser, nil
}
