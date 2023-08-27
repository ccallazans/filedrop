package usecase

import (
	"context"
	"strconv"

	"github.com/ccallazans/filedrop/internal/domain"
	"github.com/ccallazans/filedrop/internal/domain/repository"
	"github.com/ccallazans/filedrop/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

type UserUsecase struct {
	userStore repository.UserStore
	fileStore repository.FileStore
}

func NewUserUsecase(userStore repository.UserStore, fileStore repository.FileStore) *UserUsecase {
	return &UserUsecase{
		userStore: userStore,
		fileStore: fileStore,
	}
}

func (u *UserUsecase) CreateUser(ctx context.Context, firstName string, lastName string, email string, password string) (*domain.User, error) {

	userExists, _ := u.userStore.FindByEmail(ctx, email)
	if userExists != nil {
		return nil, &utils.ConflictError{Message: "email already registered"}
	}

	hashedPassword, err := hashPassword(password)
	if err != nil {
		return nil, err
	}

	newUser := domain.NewUser(firstName, lastName, email, hashedPassword)

	err = u.userStore.Save(ctx, newUser)
	if err != nil {
		return nil, &utils.InternalError{}
	}

	return newUser, nil
}

func (u *UserUsecase) GetUserByID(ctx context.Context, id string) (*domain.User, error) {

	parseID, err := strconv.ParseUint(id, 10, 0)
	if err != nil {
		utils.Logger.Infof("error when parse id %s to uint", id)
		return nil, &utils.BadRequestError{Message: "id should be an integer"}
	}

	user, _ := u.userStore.FindByID(ctx, uint(parseID))
	if user == nil {
		return nil, &utils.NotFoundError{Message: "user does not exist"}
	}

	return user, nil
}

func (u *UserUsecase) GetAllUsers(ctx context.Context) ([]*domain.User, error) {
	users := u.userStore.FindAll(ctx)

	if len(users) == 0 {
		return nil, &utils.NoContentError{Message: "no users found"}
	}

	return users, nil
}

func (u *UserUsecase) DeleteUserByID(ctx context.Context, id string) error {
	parseID, err := strconv.ParseUint(id, 10, 0)
	if err != nil {
		utils.Logger.Infof("error when parse id %s to uint", id)
		return &utils.BadRequestError{Message: "id should be an integer"}
	}

	_, err = u.GetUserByID(ctx, id)
	if err != nil {
		return err
	}

	err = u.userStore.DeleteByID(ctx, uint(parseID))
	if err != nil {
		return &utils.InternalError{}
	}

	return nil
}

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		utils.Logger.Errorf("error when generate bycrypt hash: %w", err)
		return "", &utils.InternalError{}
	}

	return string(hashedPassword), nil
}
