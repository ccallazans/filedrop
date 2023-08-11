package auth

import (
	"github.com/golang-jwt/jwt/v5"
)

type JWTClaim struct {
	User JWTUser
	jwt.RegisteredClaims
}

type JWTUser struct {
	UUID  string
	Email string
	Role  string
}
