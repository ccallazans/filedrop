package service

import (
	"context"
	"errors"
)

var ErrGetUserFromContext = "Failed to get user from context"

func getContextUser(ctx context.Context) (*JWTUser, error) {
	ctxValue := ctx.Value("user")
	if ctxValue == nil {
		return nil, errors.New("error when get user from context")
	}

	jwtUser := ctxValue.(*JWTUser)

	return jwtUser, nil
}
