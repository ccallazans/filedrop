package middlewares

import (
	"context"
	"errors"
	"os"
	"strings"

	"github.com/ccallazans/filedrop/internal/api/handlers"
	"github.com/ccallazans/filedrop/internal/application/auth"
	"github.com/ccallazans/filedrop/internal/utils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func AuthenticationMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		// extract the token from the Authorization header
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return handlers.ParseApiError(&utils.ErrorType{Type: utils.BadRequestErr, Message: "authorization header is empty"})
		}
		authHeaderParts := strings.Split(authHeader, " ")
		if len(authHeaderParts) != 2 || strings.ToLower(authHeaderParts[0]) != "bearer" {
			return handlers.ParseApiError(&utils.ErrorType{Type: utils.BadRequestErr, Message: "invalid authorization header format, expected 'Bearer <token>'"})
		}
		tokenString := authHeaderParts[1]

		// validate the token
		token, err := jwt.ParseWithClaims(tokenString, &auth.JWTClaim{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SIGNING_KEY")), nil
		})
		if err != nil {
			return handlers.ParseApiError(&utils.ErrorType{Type: utils.AuthenticationErr, Message: "invalid token"})
		}

		claims, ok := token.Claims.(*auth.JWTClaim)
		if !ok || !token.Valid {
			return handlers.ParseApiError(&utils.ErrorType{Type: utils.BadRequestErr, Message: "token expired"})
		}

		// set the user context for the downstream handlers
		user := claims.User

		ctx := context.WithValue(c.Request().Context(), "user", &user)
		c.SetRequest(c.Request().WithContext(ctx))

		return next(c)
	}
}

func AuthorizationMiddleware(next echo.HandlerFunc, roles []string) echo.HandlerFunc {
	return func(c echo.Context) error {

		user, err := utils.GetContextUser(c.Request().Context())
		if err != nil {
			return err
		}

		for _, role := range roles {
			if user.Role == role {
				return next(c)
			}
		}

		return errors.New("user not authorized")
	}
}
