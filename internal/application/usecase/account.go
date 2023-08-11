package usecase

import (
	"context"
	"errors"
	"log"
	"os"
	"time"

	"github.com/ccallazans/filedrop/internal/domain"
	"github.com/ccallazans/filedrop/internal/domain/repository"
	"github.com/ccallazans/filedrop/internal/utils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AccountUsecase struct {
	userRepo repository.UserRepository
	fileRepo repository.FileRepository
}

func NewAccountUsecase(userRepo repository.UserRepository, fileRepo repository.FileRepository) AccountUsecase {
	return AccountUsecase{
		userRepo: userRepo,
		fileRepo: fileRepo,
	}
}

func (a *AccountUsecase) CreateUser(ctx context.Context, email string, password string) error {

	existingUser, err := a.userRepo.FindByEmail(ctx, email)
	if existingUser != nil {
		return &utils.ErrorType{Type: utils.ValidationErr, Message: "email already registered!"}
	}

	hashedPassword, err := hashPassword(password)
	if err != nil {
		log.Printf("error hashing password: %s", err.Error())
		return &utils.ErrorType{Type: utils.InternalErr, Message: err.Error()}
	}

	newUser := &domain.User{
		Email:    email,
		Password: hashedPassword,
	}

	err = a.userRepo.Save(ctx, newUser)
	if err != nil {
		log.Println("error saving user into database: %s", err.Error())
		return &utils.ErrorType{Type: utils.InternalErr, Message: err.Error()}
	}

	return nil
}

func (a *AccountUsecase) AuthUser(ctx context.Context, email string, password string) (string, error) {

	validUser, err := a.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return "", &utils.ErrorType{Type: utils.ValidationErr, Message: "email not registered!"}
	}

	err = bcrypt.CompareHashAndPassword([]byte(validUser.Password), []byte(password))
	if err != nil {
		return "", &utils.ErrorType{Type: utils.ValidationErr, Message: "wrong password!"}
	}

	token, err := generateJWT(validUser)
	if err != nil {
		return "", err
	}

	return token, nil
}

type JWTClaim struct {
	User JWTUser
	jwt.RegisteredClaims
}

type JWTUser struct {
	UUID  uuid.UUID
	Name  string
	Email string
	Role  domain.UserRole
}

func generateJWT(user *domain.User) (string, error) {

	claims := &JWTClaim{
		User: JWTUser{
			UUID:  user.UUID,
			Name:  user.Name,
			Email: user.Email,
			Role:  user.Role,
		},
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "localhost",
			Subject:   user.UUID.String(),
			Audience:  jwt.ClaimStrings{"localhost"},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			NotBefore: jwt.NewNumericDate(time.Now()),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SIGN_KEY")))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}
