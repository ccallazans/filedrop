package utils

import (
	"context"
	"fmt"
	"log"

	"github.com/ccallazans/filedrop/internal/application/auth"
)

func GetContextValue(ctx context.Context, key string) (*interface{}, error) {

	value := ctx.Value(key)

	if value == nil {
		log.Printf("error getting value %s on context", key)
		return nil, &ErrorType{Type: InternalErr, Message: fmt.Sprintf("error getting value %s on context", key)}
	}

	return &value, nil
}

func GetContextUser(ctx context.Context) (*auth.JWTUser, error) {
	ctxValue := ctx.Value("user")
	if ctxValue == nil {
		return nil, &ErrorType{Type: InternalErr, Message: "could not retrieve user from context"}
	}

	jwtUser := ctxValue.(auth.JWTUser)
	
	return &jwtUser, nil
}
