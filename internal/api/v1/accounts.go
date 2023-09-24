package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

func (a *api) Login(c echo.Context) error {
	type LoginRequest struct {
		Email    string `json:"email" validate:"required"`
		Password string `json:"password" validate:"required"`
	}

	var request LoginRequest
	err := c.Bind(&request)
	if err != nil {
		return err
	}

	err = validator.New().Struct(request)
	if err != nil {
		return err
	}

	token, err := a.accountService.Login(c.Request().Context(), request.Email, request.Password)
	if err != nil {
		return err
	}

	type LoginResponse struct {
		Token string `json:"token"`
	}

	return c.JSON(http.StatusOK, LoginResponse{token})
}

func (a *api) Register(c echo.Context) error {
	type RegisterRequest struct {
		FirstName string `json:"first_name" validate:"required"`
		LastName  string `json:"last_name" validate:"required"`
		Email     string `json:"email" validate:"required"`
		Password  string `json:"password" validate:"required"`
	}

	var request RegisterRequest
	err := c.Bind(&request)
	if err != nil {
		return err
	}

	err = validator.New().Struct(request)
	if err != nil {
		return err
	}

	user, err := a.accountService.Register(c.Request().Context(), request.FirstName, request.LastName, request.Email, request.Password)
	if err != nil {
		return err
	}

	type RegisterResponse struct {
		ID        string    `json:"id"`
		FirstName string    `json:"first_name"`
		LastName  string    `json:"last_name"`
		Email     string    `json:"email"`
		Role      uint      `json:"role"`
		CreatedAt time.Time `json:"created_at"`
	}

	c.Response().Header().Set("Location", fmt.Sprintf("/users/%s", user.ID))
	return c.JSON(http.StatusCreated, RegisterResponse{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Role:      user.RoleID,
		CreatedAt: user.CreatedAt,
	})
}

func (a *api) FindAccountByID(c echo.Context) error {
	id := c.Param("id")

	user, err := a.accountService.FindByID(c.Request().Context(), id)
	if err != nil {
		return err
	}

	type FindAccountByIDResponse struct {
		ID        string    `json:"id"`
		FirstName string    `json:"first_name"`
		LastName  string    `json:"last_name"`
		Email     string    `json:"email"`
		Role      uint      `json:"role"`
		CreatedAt time.Time `json:"created_at"`
	}

	return c.JSON(http.StatusOK, FindAccountByIDResponse{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Role:      user.RoleID,
		CreatedAt: user.CreatedAt,
	})
}
