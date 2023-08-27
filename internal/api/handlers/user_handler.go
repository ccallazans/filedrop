package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/ccallazans/filedrop/internal/application/usecase"
	"github.com/ccallazans/filedrop/internal/utils"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	userUsecase usecase.UserUsecase
}

func NewUserHandler(userUsecase usecase.UserUsecase) *UserHandler {
	return &UserHandler{
		userUsecase: userUsecase,
	}
}

func (h *UserHandler) CreateUser(c echo.Context) error {

	type CreateUserRequest struct {
		FirstName string `json:"first_name" validate:"required"`
		LastName  string `json:"last_name" validate:"required"`
		Email     string `json:"email" validate:"required"`
		Password  string `json:"password" validate:"required"`
	}

	var request CreateUserRequest
	err := c.Bind(&request)
	if err != nil {
		return ParseApiError(&utils.ValidationError{Message: "bad request"})
	}

	err = validator.New().Struct(request)
	if err != nil {
		return ParseApiError(&utils.ValidationError{Message: "bad request"})
	}

	user, err := h.userUsecase.CreateUser(c.Request().Context(), request.FirstName, request.LastName, request.Email, request.Password)
	if err != nil {
		return ParseApiError(err)
	}

	type CreateUserResponse struct {
		ID        uint      `json:"id"`
		FirstName string    `json:"first_name"`
		LastName  string    `json:"last_name"`
		Email     string    `json:"email"`
		Role      uint      `json:"role"`
		CreatedAt time.Time `json:"created_at"`
	}

	c.Response().Header().Set("Location", fmt.Sprintf("/users/%d", user.ID))
	return c.JSON(http.StatusCreated, CreateUserResponse{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Role:      user.RoleID,
		CreatedAt: user.CreatedAt,
	})
}

func (h *UserHandler) GetUserByID(c echo.Context) error {
	id := c.Param("id")

	user, err := h.userUsecase.GetUserByID(c.Request().Context(), id)
	if err != nil {
		return ParseApiError(err)
	}

	type GetUserByIDResponse struct {
		ID        uint      `json:"id"`
		FirstName string    `json:"first_name"`
		LastName  string    `json:"last_name"`
		Email     string    `json:"email"`
		Role      uint      `json:"role"`
		CreatedAt time.Time `json:"created_at"`
	}

	return c.JSON(http.StatusOK, GetUserByIDResponse{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Role:      user.RoleID,
		CreatedAt: user.CreatedAt,
	})
}

func (h *UserHandler) GetAllUsers(c echo.Context) error {

	users, err := h.userUsecase.GetAllUsers(c.Request().Context())
	if err != nil {
		c.NoContent(http.StatusNoContent)
	}

	type GetAllUsersResponse struct {
		ID        uint      `json:"id"`
		FirstName string    `json:"first_name"`
		LastName  string    `json:"last_name"`
		Email     string    `json:"email"`
		Role      uint      `json:"role"`
		CreatedAt time.Time `json:"created_at"`
	}

	usersResponse := []GetAllUsersResponse{}
	for _, user := range users {
		usersResponse = append(usersResponse, GetAllUsersResponse{
			ID:        user.ID,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
			Role:      user.RoleID,
			CreatedAt: user.CreatedAt,
		})
	}

	return c.JSON(http.StatusOK, usersResponse)
}

func (h *UserHandler) DeleteUserByID(c echo.Context) error {
	id := c.Param("id")

	err := h.userUsecase.DeleteUserByID(c.Request().Context(), id)
	if err != nil {
		return ParseApiError(err)
	}

	return c.NoContent(http.StatusOK)
}
