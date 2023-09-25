package service

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/ccallazans/filedrop/internal/domain"
	"github.com/ccallazans/filedrop/internal/domain/repository"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserNotFound           = "user not found"
	ErrEmailAlreadyRegistered = "email already registered"
	ErrInvalidPassword        = "wrong password"
	ErrCreatingUser           = "an error occured when creating user"
)

type JWTClaim struct {
	jwt.RegisteredClaims
	User JWTUser
}

type JWTUser struct {
	ID        string
	FirstName string
	Email     string
	Role      uint
}

type AuthService struct {
	userStore repository.UserStore
}

func NewAuthService(userStore repository.UserStore) *AuthService {
	return &AuthService{
		userStore: userStore,
	}
}

func (s *AuthService) Login(ctx context.Context, email string, password string) (string, error) {
	validUser, err := s.userStore.FindByEmail(ctx, email)
	if validUser == nil {
		return "", fmt.Errorf(ErrUserNotFound)
	}

	err = bcrypt.CompareHashAndPassword([]byte(validUser.Password), []byte(password))
	if err != nil {
		return "", fmt.Errorf(ErrInvalidPassword)
	}

	token, err := generateJWT(validUser)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *AuthService) Register(ctx context.Context, firstName string, email string, password string) (*domain.User, error) {
	userExists, _ := s.userStore.FindByEmail(ctx, email)
	if userExists != nil {
		return nil, fmt.Errorf(ErrEmailAlreadyRegistered)
	}

	hashedPassword, err := hash(password)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		ID:        uuid.NewString(),
		FirstName: firstName,
		Email:     email,
		Password:  hashedPassword,
		RoleID:    uint(domain.USER),
	}

	err = s.userStore.Save(ctx, user)
	if err != nil {
		return nil, fmt.Errorf(ErrCreatingUser)
	}

	return user, nil
}

func (s *AuthService) FindByID(ctx context.Context, id string) (*domain.User, error) {
	user, _ := s.userStore.FindByID(ctx, id)
	if user == nil {
		return nil, fmt.Errorf(ErrUserNotFound)
	}

	return user, nil
}

func hash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}

func generateJWT(user *domain.User) (string, error) {
	claims := JWTClaim{
		User: JWTUser{
			ID:        user.ID,
			FirstName: user.FirstName,
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
