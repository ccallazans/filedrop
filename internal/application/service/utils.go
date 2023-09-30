package service

import (
	"context"
	"errors"
)

func GetUserFromCtx(ctx context.Context) (*JWTUser, error) {
	user, ok := ctx.Value("user").(*JWTUser)
	if !ok {
		return nil, errors.New("user not found in context")
	}

	return user, nil
}
