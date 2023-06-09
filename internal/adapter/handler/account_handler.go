package handler

import (
	"log"
	"net/http"

	"github.com/ccallazans/filedrop/internal/application/usecase"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type AccountHandler struct {
	accountUsecase usecase.AccountUsecase
}

func NewAccountHandler(accountUsecase usecase.AccountUsecase) *AccountHandler {
	return &AccountHandler{
		accountUsecase: accountUsecase,
	}
}

func (a *AccountHandler) CreateUser(c echo.Context) error {

	type CreateUserRequest struct {
		Name     string `json:"name" validate:"required"`
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
	}

	var createUserRequest CreateUserRequest
	err := c.Bind(&createUserRequest)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	err = validator.New().Struct(&createUserRequest)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	err = a.accountUsecase.CreateUser(usecase.CreateUserArgs{
		Name:     createUserRequest.Name,
		Email:    createUserRequest.Email,
		Password: createUserRequest.Password,
	})
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusForbidden, err.Error())
	}

	return c.JSON(http.StatusOK, createUserRequest)
}

func (a *AccountHandler) AuthUser(c echo.Context) error {

	type AuthUserRequest struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
	}

	var authUserRequest AuthUserRequest
	err := c.Bind(&authUserRequest)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	err = validator.New().Struct(authUserRequest)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	token, err := a.accountUsecase.AuthUser(usecase.AuthUserArgs{
		Email:    authUserRequest.Email,
		Password: authUserRequest.Password,
	})
	if err != nil {
		return c.JSON(http.StatusForbidden, err)
	}

	type AuthUserResponse struct {
		Token string `json:"token"`
	}

	return c.JSON(http.StatusOK, AuthUserResponse{Token: token})
}
