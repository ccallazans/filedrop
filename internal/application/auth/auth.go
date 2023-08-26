package auth

import (
	"github.com/golang-jwt/jwt/v5"
)

type JWTClaim struct {
	User JWTUser
	jwt.RegisteredClaims
}

type JWTUser struct {
	ID    uint
	UUID  string
	Email string
	Role  string
}

type UserCtxKey struct{}
