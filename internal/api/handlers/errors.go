package handlers

import (
	"net/http"
	"strings"

	"github.com/ccallazans/filedrop/internal/utils"
	"github.com/labstack/echo/v4"
)

type ErrorResponse struct {
	Status int    `json:"status"`
	Title  string `json:"title"`
	Detail string `json:"detail"`
}

func ParseApiError(err error) error {
	errorMessage := err.Error()
	types := strings.Split(errorMessage, " -> ")

	if len(types) != 2 {
		return echo.NewHTTPError(
			http.StatusInternalServerError,
			ErrorResponse{
				Status: http.StatusInternalServerError,
				Title:  "Internal server error",
				Detail: errorMessage,
			})
	}

	switch types[0] {
	case utils.AuthenticationErr:
		return echo.NewHTTPError(
			http.StatusUnauthorized,
			ErrorResponse{
				Status: http.StatusUnauthorized,
				Title:  "Unauthorized",
				Detail: types[1],
			})
	case utils.BadRequestErr:
		return echo.NewHTTPError(
			http.StatusBadRequest,
			ErrorResponse{
				Status: http.StatusBadRequest,
				Title:  "Bad Request",
				Detail: types[1],
			})
	case utils.ValidationErr:
		return echo.NewHTTPError(
			http.StatusBadRequest,
			ErrorResponse{
				Status: http.StatusBadRequest,
				Title:  "Bad Request",
				Detail: types[1],
			})
	case utils.InternalErr:
		return echo.NewHTTPError(
			http.StatusInternalServerError,
			ErrorResponse{
				Status: http.StatusInternalServerError,
				Title:  "Internal server error",
				Detail: types[1],
			})
	default:
		return echo.NewHTTPError(
			http.StatusInternalServerError,
			ErrorResponse{
				Status: http.StatusInternalServerError,
				Title:  "Internal server error",
				Detail: errorMessage,
			})
	}
}
