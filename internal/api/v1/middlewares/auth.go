package middlewares

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/ccallazans/filedrop/internal/application/service"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func Auth(next echo.HandlerFunc, roles []int) echo.HandlerFunc {
	return func(c echo.Context) error {
		// extract the token from the Authorization header
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return fmt.Errorf("authorization header is empty") //&utils.BadRequestError{Message: "authorization header is empty"}
		}
		authHeaderParts := strings.Split(authHeader, " ")
		if len(authHeaderParts) != 2 || strings.ToLower(authHeaderParts[0]) != "bearer" {
			return fmt.Errorf("invalid authorization header format, expected 'Bearer <token>'") //&utils.BadRequestError{Message: "invalid authorization header format, expected 'Bearer <token>'"}
		}
		tokenString := authHeaderParts[1]

		// validate the token
		token, err := jwt.ParseWithClaims(tokenString, &service.JWTClaim{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SIGNING_KEY")), nil
		})
		if err != nil {
			return err //&utils.AuthenticationError{Message: "invalid token"}
		}

		claims, ok := token.Claims.(*service.JWTClaim)
		if !ok || !token.Valid {
			return fmt.Errorf("token expired")//&utils.AuthenticationError{Message: "token expired"}
		}

		// set the user context for the downstream handlers
		user := claims.User

		err = verifyRoles(user.Role, roles)
		if err != nil {
			return err//&utils.AuthorizationError{Message: "user not allowed"}
		}

		ctx := context.WithValue(c.Request().Context(), "user", &user)
		c.SetRequest(c.Request().WithContext(ctx))

		return next(c)
	}
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
		return errors.New("aaaaaaaaaaa")
	}

	return nil
}
