package usecase

import (
	"context"
	"errors"
)

func GetContextUser(ctx context.Context) (*JWTUser, error) {
	ctxValue := ctx.Value("user")
	if ctxValue == nil {
		// utils.Logger.Error("error when retriever user from context")
		return nil, errors.New("error when get user from context")
	}

	jwtUser := ctxValue.(*JWTUser)

	return jwtUser, nil
}
