package handlers

import (
	"net/http"

	"github.com/ccallazans/filedrop/internal/application/usecase"
	"github.com/ccallazans/filedrop/internal/utils"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	authUsecase usecase.AuthUsecase
}

func NewAuthHandler(authUsecase usecase.AuthUsecase) *AuthHandler {
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
		return ParseApiError(&utils.BadRequestError{})
	}

	err = validator.New().Struct(request)
	if err != nil {
		return ParseApiError(&utils.BadRequestError{})
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
