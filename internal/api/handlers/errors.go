package handlers

import (
	"net/http"

	"github.com/ccallazans/filedrop/internal/utils"
	"github.com/labstack/echo/v4"
)

type ErrorResponse struct {
	Status int    `json:"status"`
	Title  string `json:"title"`
	Detail string `json:"detail"`
}

func ParseApiError(err error) error {

	switch err.(type) {
	case *utils.AuthenticationError:
		return echo.NewHTTPError(
			http.StatusUnauthorized,
			ErrorResponse{
				Status: http.StatusUnauthorized,
				Title:  "Unauthorized",
				Detail: err.Error(),
			})
	case *utils.AuthorizationError:
		return echo.NewHTTPError(
			http.StatusForbidden,
			ErrorResponse{
				Status: http.StatusForbidden,
				Title:  "Forbidden",
				Detail: err.Error(),
			})

	//

	case *utils.BadRequestError:
		message := err.Error()
		if len(message) == 0 {
			message = "bad request"
		}

		return echo.NewHTTPError(
			http.StatusBadRequest,
			ErrorResponse{
				Status: http.StatusBadRequest,
				Title:  "Bad Request",
				Detail: message,
			})
	case *utils.ValidationError:
		return echo.NewHTTPError(
			http.StatusUnprocessableEntity,
			ErrorResponse{
				Status: http.StatusUnprocessableEntity,
				Title:  "Unprocessable Entity",
				Detail: err.Error(),
			})
	case *utils.ConflictError:
		return echo.NewHTTPError(
			http.StatusConflict,
			ErrorResponse{
				Status: http.StatusConflict,
				Title:  "Conflict",
				Detail: err.Error(),
			})

	//

	case *utils.NoContentError:
		return echo.NewHTTPError(
			http.StatusNoContent,
			ErrorResponse{
				Status: http.StatusNoContent,
				Title:  "No Content",
				Detail: err.Error(),
			})
	case *utils.NotFoundError:
		return echo.NewHTTPError(
			http.StatusNotFound,
			ErrorResponse{
				Status: http.StatusNotFound,
				Title:  "Not Found",
				Detail: err.Error(),
			})

	//

	case *utils.InternalError:
		message := err.Error()
		if len(message) == 0 {
			message = "internal server error"
		}

		return echo.NewHTTPError(
			http.StatusInternalServerError,
			ErrorResponse{
				Status: http.StatusInternalServerError,
				Title:  "Internal server error",
				Detail: message,
			})
	default:
		return echo.NewHTTPError(
			http.StatusInternalServerError,
			ErrorResponse{
				Status: http.StatusInternalServerError,
				Title:  "Internal server error",
				Detail: err.Error(),
			})
	}
}
