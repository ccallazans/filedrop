package usecase

import (
	"context"

	"github.com/ccallazans/filedrop/internal/utils"
)

func GetContextUser(ctx context.Context) (*JWTUser, error) {
	ctxValue := ctx.Value("user")
	if ctxValue == nil {
		utils.Logger.Error("error when retriever user from context")
		return nil, &utils.InternalError{}
	}

	jwtUser := ctxValue.(*JWTUser)

	return jwtUser, nil
}
