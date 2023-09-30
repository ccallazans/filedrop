package api

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

func (a *api) Login(c echo.Context) error {
	var request LoginRequest
	err := c.Bind(&request)
	if err != nil {
		return err
	}

	err = validator.New().Struct(request)
	if err != nil {
		return err
	}

	token, err := a.authService.Login(c.Request().Context(), request.Email, request.Password)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, LoginResponse{token})
}

func (a *api) Register(c echo.Context) error {
	var request RegisterRequest
	err := c.Bind(&request)
	if err != nil {
		return err
	}

	err = validator.New().Struct(request)
	if err != nil {
		return err
	}

	user, err := a.authService.Register(c.Request().Context(), request.FirstName, request.Email, request.Password)
	if err != nil {
		return err
	}

	c.Response().Header().Set("Location", fmt.Sprintf("/users/%s", user.ID))
	return c.JSON(http.StatusCreated, UserResponse{
		ID:        user.ID,
		FirstName: user.FirstName,
		Email:     user.Email,
		Role:      user.RoleID,
		CreatedAt: user.CreatedAt,
	})
}

func (a *api) FindAccountByID(c echo.Context) error {
	id := c.Param("id")

	user, err := a.authService.FindByID(c.Request().Context(), id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, UserResponse{
		ID:        user.ID,
		FirstName: user.FirstName,
		Email:     user.Email,
		Role:      user.RoleID,
		CreatedAt: user.CreatedAt,
	})
}
