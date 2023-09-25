package service

import (
	"context"
	"errors"
)

var ErrGetUserFromContext = "Failed to get user from context"

func getContextUser(ctx context.Context) (*JWTClaim, error) {
	ctxValue := ctx.Value("user")
	if ctxValue == nil {
		return nil, errors.New("error when get user from context")
	}

	ctxUser := ctxValue.(*JWTClaim)

	return ctxUser, nil
}
