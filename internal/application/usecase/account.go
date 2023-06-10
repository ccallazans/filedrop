package usecase

import (
	"errors"
	"os"
	"time"

	"github.com/ccallazans/filedrop/internal/domain"
	"github.com/ccallazans/filedrop/internal/domain/repository"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AccountUsecase struct {
	userRepo repository.IUser
	fileRepo repository.IFile
}

func NewAccountUsecase(userRepo repository.IUser, fileRepo repository.IFile) AccountUsecase {
	return AccountUsecase{
		userRepo: userRepo,
		fileRepo: fileRepo,
	}
}

type CreateUserArgs struct {
	Name     string
	Email    string
	Password string
}

func (a *AccountUsecase) CreateUser(args CreateUserArgs) error {

	_, err := a.userRepo.FindByEmail(args.Email)
	if err == nil {
		return errors.New("email already registered")
	}

	args.Password, err = hashPassword(args.Password)
	if err != nil {
		return errors.New("error hashing password")
	}

	newUser := &domain.User{
		UUID:       uuid.New(),
		Name:     args.Name,
		Email:    args.Email,
		Password: args.Password,
		Role:     domain.USER,
	}

	err = a.userRepo.Save(newUser)
	if err != nil {
		return errors.New("error creating user")
	}

	return nil
}

type AuthUserArgs struct {
	Email    string
	Password string
}

func (a *AccountUsecase) AuthUser(args AuthUserArgs) (string, error) {

	validUser, err := a.userRepo.FindByEmail(args.Email)
	if err != nil {
		return "", errors.New("email not registered")
	}

	err = bcrypt.CompareHashAndPassword([]byte(validUser.Password), []byte(args.Password))
	if err != nil {
		return "", errors.New("wrong password")
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
