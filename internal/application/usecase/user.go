package usecase

import (
	"context"
	"strconv"

	"github.com/ccallazans/filedrop/internal/domain"
	"github.com/ccallazans/filedrop/internal/domain/repository"
	"github.com/ccallazans/filedrop/internal/utils"
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

	newUser, err := domain.NewUser(firstName, lastName, email, password)
	if err != nil {
		return nil, err
	}

	err = u.userStore.Save(ctx, newUser)
	if err != nil {
		return nil, err
	}

	return newUser, nil
}

func (u *UserUsecase) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
	parseID, err := strconv.ParseUint(id, 10, 0)
	if err != nil {
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
	return users, nil
}

func (u *UserUsecase) DeleteUserByID(ctx context.Context, id string) error {
	parseID, err := strconv.ParseUint(id, 10, 0)
	if err != nil {
		return &utils.BadRequestError{Message: "id should be an integer"}
	}

	_, err = u.GetUserByID(ctx, id)
	if err != nil {
		return err
	}

	err = u.userStore.DeleteByID(ctx, uint(parseID))
	if err != nil {
		return err
	}

	return nil
}
