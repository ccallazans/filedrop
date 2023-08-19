package usecase

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/ccallazans/filedrop/internal/application/auth"
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
		return err
	}

	newUser := &domain.User{
		UUID: uuid.NewString(),
		Email:    email,
		Password: hashedPassword,
		RoleID: domain.USER,
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

func generateJWT(user *domain.User) (string, error) {

	claims := &auth.JWTClaim{
		User: auth.JWTUser{
			ID: user.ID,
			UUID:  user.UUID,
			Email: user.Email,
			Role:  user.Role.Role,
		},
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "filedrop",
			Subject:   user.UUID,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(2 * time.Hour)),
			NotBefore: jwt.NewNumericDate(time.Now()),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
	if err != nil {
		log.Printf("error creating token: %s", err.Error())
		return "", &utils.ErrorType{Type: utils.InternalErr, Message: err.Error()}
	}

	return tokenString, nil
}

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("error hashing password: %s", err.Error())
		return "", &utils.ErrorType{Type: utils.InternalErr, Message: err.Error()}
	}

	return string(hashedPassword), nil
}
