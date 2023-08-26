package handlers

import (
	"net/http"

	"github.com/ccallazans/filedrop/internal/application/usecase"
	"github.com/ccallazans/filedrop/internal/utils"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	authUsecase usecase.AccountUsecase
}

func NewAuthHandler(authUsecase usecase.AccountUsecase) *AuthHandler {
	return &AuthHandler{
		authUsecase: authUsecase,
	}
}

func (h *AuthHandler) Signin(c echo.Context) error {

	type SigninRequest struct {
		Email    string `json:"email" validate:"required"`
		Password string `json:"password" validate:"required"`
	}

	var request SigninRequest
	err := c.Bind(&request)
	if err != nil {
		return ParseApiError(&utils.ErrorType{Type: utils.BadRequestErr, Message: err.Error()})
	}

	err = validator.New().Struct(request)
	if err != nil {
		return ParseApiError(&utils.ErrorType{Type: utils.BadRequestErr, Message: err.Error()})
	}

	token, err := h.authUsecase.AuthUser(c.Request().Context(), request.Email, request.Password)
	if err != nil {
		return ParseApiError(err)
	}

	type SigninResponse struct {
		Token string `json:"token"`
	}

	return c.JSON(http.StatusOK, SigninResponse{token})
}

func (h *AuthHandler) Register(c echo.Context) error {

	type RegisterRequest struct {
		Email    string `json:"email" validate:"required"`
		Password string `json:"password" validate:"required"`
	}

	var request RegisterRequest
	err := c.Bind(&request)
	if err != nil {
		return ParseApiError(err)
	}

	err = validator.New().Struct(request)
	if err != nil {
		return ParseApiError(err)
	}

	err = h.authUsecase.CreateUser(c.Request().Context(), request.Email, request.Password)
	if err != nil {
		return ParseApiError(err)
	}

	type RegisterReponse struct {
		Email string `json:"email"`
	}

	return c.JSON(http.StatusCreated, RegisterReponse{Email: request.Email})
}
