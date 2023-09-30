package middlewares

import (
	"context"
	"errors"
	"os"
	"strings"

	"github.com/ccallazans/filedrop/internal/application/service"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

func Auth(next echo.HandlerFunc, roles []int) echo.HandlerFunc {
	return func(c echo.Context) error {
		// extract the token from the Authorization header
		token, err := extractToken(c)
		if err != nil {
			log.Error(err)
			return err
		}

		// validate the token
		claims, nil := verifyToken(token)
		if err != nil {
			log.Error(err)
			return err
		}

		// set the user context for the downstream handlers
		user := claims.User

		err = verifyRoles(user.Role, roles)
		if err != nil {
			log.Error(err)
			return err //&utils.AuthorizationError{Message: "user not allowed"}
		}

		ctx := context.WithValue(c.Request().Context(), "user", &user)
		c.SetRequest(c.Request().WithContext(ctx))

		return next(c)
	}
}

func extractToken(c echo.Context) (string, error) {
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("authorization header is empty")
	}
	authHeaderParts := strings.Split(authHeader, " ")
	if len(authHeaderParts) != 2 || strings.ToLower(authHeaderParts[0]) != "bearer" {
		return "", errors.New("invalid authorization header format, expected 'Bearer <token>'")
	}
	tokenString := authHeaderParts[1]

	return tokenString, nil
}

func verifyToken(t string) (*service.JWTClaim, error) {
	parsedToken, err := jwt.ParseWithClaims(t, &service.JWTClaim{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SIGNING_KEY")), nil
	})
	if err != nil {
		return nil, errors.New("token invalid") // invalid token
	}

	claims, ok := parsedToken.Claims.(*service.JWTClaim)
	if !ok || !parsedToken.Valid {
		return nil, errors.New("token expired")
	}

	return claims, nil
}

func verifyRoles(userRole uint, roles []int) error {
	if len(roles) == 0 {
		return nil
	}

	roleMap := map[int]bool{}
	for _, role := range roles {
		roleMap[role] = true
	}

	if !roleMap[int(userRole)] {
		return errors.New("user do not have the role to access the resource")
	}

	return nil
}
