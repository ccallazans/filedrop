package api

import (
	"net/http"

	"github.com/ccallazans/filedrop/internal/utils"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

func (a *api) Signin(c echo.Context) error {
	type SigninRequest struct {
		Email    string `json:"email" validate:"required"`
		Password string `json:"password" validate:"required"`
	}

	var request SigninRequest
	err := c.Bind(&request)
	if err != nil {
		return &utils.BadRequestError{}
	}

	err = validator.New().Struct(request)
	if err != nil {
		return &utils.BadRequestError{}
	}

	token, err := a.authUsecase.AuthUser(c.Request().Context(), request.Email, request.Password)
	if err != nil {
		return err
	}

	type SigninResponse struct {
		Token string `json:"token"`
	}

	return c.JSON(http.StatusOK, SigninResponse{token})
}
