package handler

import (
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

var validate *validator.Validate

func (a *AccountHandler) CreateUser(c echo.Context) error {

	type CreateUserRequest struct {
		Name     string `json:"name" validate:"required"`
		Email    string `json:"email" validate:"required"`
		Password string `json:"password" validate:"required,email"`
	}

	var createUserRequest CreateUserRequest
	err := c.Bind(&createUserRequest)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	err = validate.Struct(createUserRequest)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	err = a.accountUsecase.CreateUser(usecase.CreateUserArgs{
		Name:     createUserRequest.Name,
		Email:    createUserRequest.Email,
		Password: createUserRequest.Password,
	})
	if err != nil {
		return c.JSON(http.StatusForbidden, err)
	}

	return c.JSON(http.StatusOK, createUserRequest)
}

func (a *AccountHandler) AuthUser(c echo.Context) error {

	type AuthUserRequest struct {
		Email    string `json:"email" validate:"required"`
		Password string `json:"password" validate:"required,email"`
	}

	var authUserRequest AuthUserRequest
	err := c.Bind(&authUserRequest)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	err = validate.Struct(authUserRequest)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	token, err := a.accountUsecase.AuthUser(usecase.AuthUserArgs{
		Email: authUserRequest.Email,
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