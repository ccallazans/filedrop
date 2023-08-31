package usecase

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/ccallazans/filedrop/internal/domain"
	"github.com/ccallazans/filedrop/internal/domain/repository"
	"github.com/ccallazans/filedrop/internal/utils"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type JWTClaim struct {
	jwt.RegisteredClaims
	User JWTUser
}

type JWTUser struct {
	ID        uint
	FirstName string
	LastName  string
	Email     string
	Role      uint
}

type AuthUsecase struct {
	userStore repository.UserStore
}

func NewAuthUsecase(userStore repository.UserStore) *AuthUsecase {
	return &AuthUsecase{
		userStore: userStore,
	}
}

func (u *AuthUsecase) AuthUser(ctx context.Context, email string, password string) (string, error) {
	validUser, err := u.userStore.FindByEmail(ctx, email)
	if err != nil {
		return "", &utils.NotFoundError{Message: "user does not exist"}
	}

	err = bcrypt.CompareHashAndPassword([]byte(validUser.Password), []byte(password))
	if err != nil {
		return "", &utils.AuthenticationError{Message: "wrong password"}
	}

	token, err := generateJWT(validUser)
	if err != nil {
		return "", err
	}

	return token, nil
}

func generateJWT(user *domain.User) (string, error) {
	claims := JWTClaim{
		User: JWTUser{
			ID:        user.ID,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
			Role:      user.RoleID,
		},
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    os.Getenv("APP_URL"),
			Subject:   fmt.Sprint(user.ID),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(2 * time.Hour)),
			NotBefore: jwt.NewNumericDate(time.Now()),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SIGNING_KEY")))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
