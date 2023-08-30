package middlewares

import (
	"strings"

	"github.com/ccallazans/filedrop/internal/config"
	"github.com/ccallazans/filedrop/internal/utils"
	"github.com/labstack/echo/v4"
)

func AuthenticationMiddleware(next echo.HandlerFunc, roles []int) echo.HandlerFunc {
	return func(c echo.Context) error {
		loasda := config.NewLogger("middleware")
		loasda.Error("asaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
		// extract the token from the Authorization header
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return &utils.BadRequestError{Message: "authorization header is empty"}
		}
		authHeaderParts := strings.Split(authHeader, " ")
		if len(authHeaderParts) != 2 || strings.ToLower(authHeaderParts[0]) != "bearer" {
			return &utils.BadRequestError{Message: "invalid authorization header format, expected 'Bearer <token>'"}
		}
		// tokenString := authHeaderParts[1]

		// // validate the token
		// token, err := jwt.ParseWithClaims(tokenString, &usecase.JWTClaim{}, func(token *jwt.Token) (interface{}, error) {
		// 	return []byte(os.Getenv("JWT_SIGNING_KEY")), nil
		// })
		// if err != nil {
		// 	return api.ParseApiError(&utils.AuthenticationError{Message: "invalid token"})
		// }

		// claims, ok := token.Claims.(*usecase.JWTClaim)
		// if !ok || !token.Valid {
		// 	return api.ParseApiError(&utils.AuthenticationError{Message: "token expired"})
		// }

		// // set the user context for the downstream handlers
		// user := claims.User

		// err = verifyRoles(user.Role, roles)
		// if err != nil {
		// 	return api.ParseApiError(&utils.AuthorizationError{Message: "user not allowed"})
		// }

		// ctx := context.WithValue(c.Request().Context(), "user", &user)
		// c.SetRequest(c.Request().WithContext(ctx))

		return next(c)
	}
}

// func verifyRoles(userRole uint, roles []int) error {
// 	if len(roles) == 0 {
// 		return nil
// 	}

// 	roleMap := map[int]bool{}
// 	for _, role := range roles {
// 		roleMap[role] = true
// 	}

// 	if !roleMap[int(userRole)] {
// 		return errors.New("")
// 	}

// 	return nil
// }
